package repository

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"time"
	"unicode/utf8"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/external"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type GameRepository interface {
	InitWord(roomID string, word string) error
	LockCurrentWord(roomID string) error
	UnlockCurrentWord(roomID string) error
	UpdateCurrentWord(roomID string, word string) error
	GetCurrentWord(roomID string) (string, error)
	AddUser(roomID string, userID string) error
	GetUsers(roomID string) ([]string, error)
	DeleteUser(roomID string, userID string) error
	LockRoomUsers(roomID string) error
	UnlockRoomUsers(roomID string) error

	// counter
	IncrCounter(roomID string) (int64, error)
	SetCounter(roomID string, value int64) error
	GetCounter(roomID string) (int64, error)

	// message
	Publish(message *entity.Message) error
	Subscribe(ctx context.Context, roomID string) (<-chan *entity.Message, <-chan error)
	ContainWord(word string) bool

	// room
	Lock() error
	Unlock() error
	GetRoomCandidates() ([]string, error)
	AddRoomCandidate(roomID string) error
	DeleteRoomCandidate(roomID string) error
	DeleteRoom(roomID string) error

	// result
	GetScore(roomID string, userID string) (*entity.Result, error)
	SetScore(roomID string, userID string, score int64) error
	IncrScoreBy(roomID string, userID string, by int64) error
	GetLatestRating(userID string) (int64, error)
	SetRating(userID string, rating int64) error
	AddRatingHistory(userID string, rating int64) error
}

type gameRepository struct {
	firestore  *firebase.App
	conn       *redis.Pool
	keyTTL     time.Duration
	dictionary *map[string]struct{}
}

var (
	CLIENT      = "client"
	roomLockKey = "roomCandidates"
)

type messageInRedis struct {
	from string `json:"from" validate:"required"`
	*entity.Message
}

func NewGameRepository() *gameRepository {
	dictionary := map[string]struct{}{}

	if dicPath, ok := os.LookupEnv("DIC_PATH"); !ok {
		log.WithError(xerrors.Errorf("NewGameRepository LookupEnv(DIC_PATH) is false.")).Warn()
		dictionary["りんご"] = struct{}{}
		dictionary["ごま"] = struct{}{}
		dictionary["まり"] = struct{}{}
	} else {
		fp, err := os.Open(dicPath)
		if err != nil {
			panic(err)
		}
		defer fp.Close()

		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			line := scanner.Text()
			dictionary[line] = struct{}{}
			if !utf8.ValidString(line) {
				log.WithError(xerrors.Errorf("NewGameRepository utf8.ValidString is false : %w", err)).Fatalf("line: %s", line)
			}
		}
		if err := scanner.Err(); err != nil {
			log.WithError(xerrors.Errorf("NewGameRepository scanner error: %w", err)).Fatal("")
		}
	}

	return &gameRepository{
		firestore:  external.FirebaseApp,
		conn:       external.RedisPool,
		keyTTL:     time.Minute * 10,
		dictionary: &dictionary,
	}
}

func (r *gameRepository) InitWord(roomID string, word string) error {
	key := roomID + ":currentWord"
	conn := r.conn.Get()
	_, err := conn.Do("SET", key, word)
	if err != nil {
		return xerrors.Errorf("error in InitWord setnx: %w", err)
	}

	_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
	if err != nil {
		return xerrors.Errorf("error in InitWord expire: %w", err)
	}

	return err
}

func (r *gameRepository) LockCurrentWord(roomID string) error {
	key := roomID + ":currentWord:lock"
	conn := r.conn.Get()
	// while lockできなかったらloop
	for {
		res, err := conn.Do("SETNX", key, "locking")
		if err != nil {
			return xerrors.Errorf("error in lockCurrentWord: %w", err)
		}

		_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
		if err != nil {
			return xerrors.Errorf("error in lockCurrentWord expire: %w", err)
		}

		if res != 0 {
			break
		}
	}
	return nil
}

// updateして r.conn.Do("delete roomID+"currentWord:lock")
func (r *gameRepository) UnlockCurrentWord(roomID string) error {
	key := roomID + ":currentWord:lock"
	conn := r.conn.Get()
	_, err := conn.Do("DEL", key)
	if err != nil {
		return xerrors.Errorf("error in UnlockCurrentWord: %w", err)
	}
	return nil
}

// <roomID>:currentWord
func (r *gameRepository) UpdateCurrentWord(roomID string, word string) error {
	key := roomID + ":currentWord"
	conn := r.conn.Get()

	_, err := conn.Do("SET", key, word)
	if err != nil {
		return xerrors.Errorf("error in UpdateCurrentWord: %w", err)
	}

	_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
	if err != nil {
		return xerrors.Errorf("error in UpdateCurrentWord expire: %w", err)
	}

	return nil
}

func (r *gameRepository) GetCurrentWord(roomID string) (string, error) {
	// lock()
	// defer unlock()
	key := roomID + ":currentWord"
	conn := r.conn.Get()

	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", xerrors.Errorf("error in getCurrentWord: %w", err)
	}

	_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
	if err != nil {
		return "", xerrors.Errorf("error in getCurrentWord expire: %w", err)
	}

	return value, nil
}

func (r *gameRepository) AddUser(roomID string, userID string) error {
	key := roomID + ":users"
	conn := r.conn.Get()

	_, err := conn.Do("HSET", key, userID, 0)
	log.WithFields(log.Fields{
		roomID: roomID,
		userID: userID,
	}).Debug()
	if err != nil {
		return xerrors.Errorf("error in AddUser HSET %s %s %d: %w", key, userID, 0, err)
	}
	return nil
}

func (r *gameRepository) GetUsers(roomID string) ([]string, error) {
	key := roomID + ":users"
	conn := r.conn.Get()

	users, err := redis.Strings(conn.Do("HKEYS", key))
	if err != nil {
		return nil, xerrors.Errorf("error in GetUsers HKEYS %s, 0, -1: %w", key, err)
	}
	return users, nil
}

func (r *gameRepository) DeleteUser(roomID string, userID string) error {
	key := roomID + ":users"
	conn := r.conn.Get()

	if _, err := conn.Do("HDEL", key, userID); err != nil {
		return xerrors.Errorf("error in GetUsers HDEL %s, %s: %w", key, userID)
	}
	return nil
}

func (r *gameRepository) LockRoomUsers(roomID string) error {
	// prevent for removing <roomID>:**
	key := "lock:" + roomID
	conn := r.conn.Get()
	_, err := conn.Do("SET", key, 0)
	if err != nil {
		return xerrors.Errorf("error in InitWord setnx: %w", err)
	}

	_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
	if err != nil {
		return xerrors.Errorf("error in LockRoomUsers() expire: %w", err)
	}

	return err
}

func (r *gameRepository) UnlockRoomUsers(roomID string) error {
	// prevent for removing <roomID>:**
	key := "lock:" + roomID
	conn := r.conn.Get()
	_, err := conn.Do("DEL", key)
	if err != nil {
		return xerrors.Errorf("error in UnlockCurrentWord: %w", err)
	}
	return nil
}

// redis counter repo の命名規則
// incr counter
// 将来は <roomID>:counter になるかも

func (r *gameRepository) IncrCounter(roomID string) (int64, error) {
	key := roomID + ":counter"
	conn := r.conn.Get()

	value, err := redis.Int64(conn.Do("Incr", key))
	if err != nil {
		return -1, xerrors.Errorf("error in incr counter: %w", err)
	}

	_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
	if err != nil {
		return -1, xerrors.Errorf("error in IncrCounter expire: %w", err)
	}
	return value, nil
}

func (r *gameRepository) SetCounter(roomID string, value int64) error {
	key := roomID + ":counter"
	conn := r.conn.Get()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return xerrors.Errorf("error in set counter: %w", err)
	}

	_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
	if err != nil {
		return xerrors.Errorf("error in set counter expire: %w", err)
	}

	return nil
}

func (r *gameRepository) GetCounter(roomID string) (int64, error) {
	key := roomID + ":counter"
	conn := r.conn.Get()

	value, err := redis.Int64(conn.Do("GET", key))
	if err != nil {
		return -1, xerrors.Errorf("error in main method: %w", err)
	}
	_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
	if err != nil {
		return -1, xerrors.Errorf("error in GetCounter expire: %w", err)
	}

	return value, nil
}

// redis message repo の命名規則
// publish message '{"userID": "7141-1414-1414...", "message": "hello"}'
// subscribe messae
// roomID:message になるかも

func (r *gameRepository) Publish(message *entity.Message) error {
	mesInRed := messageInRedis{
		from:    CLIENT,
		Message: message,
	}
	mesBytes, err := json.Marshal(&mesInRed)
	if err != nil {
		return xerrors.Errorf("error in json.Marshal: %w", err)
	}

	conn := r.conn.Get()
	defer func() {
		if err := conn.Close(); err != nil {
			panic(xerrors.Errorf("error client.Close: %w", err))
		}
	}()

	rep, err := conn.Do("PUBLISH", message.RoomID+":message", mesBytes)
	if err != nil {
		return xerrors.Errorf("error in redis publish: %w", err)
	}

	_, err = conn.Do("EXPIRE", message.RoomID+":message", int64(r.keyTTL.Seconds()))
	if err != nil {
		return xerrors.Errorf("error in Publish expire: %w", err)
	}

	log.WithFields(log.Fields{
		"rep":    rep,
		"roomId": message.RoomID,
	}).Info("publish reply")

	return nil
}

// Subscribeのより良いやり方あるかも
// ctx: 親のcontextで親のcontextが終了するとgo func()内でctx.Done()により終了する
// roomID: どこの部屋のイベントをsubscribeするか
func (r *gameRepository) Subscribe(ctx context.Context, roomID string) (<-chan *entity.Message, <-chan error) {
	ch := make(chan *entity.Message)
	errCh := make(chan error)
	go func() {
		defer close(ch)
		defer close(errCh)

		conn := r.conn.Get()
		defer func() {
			if err := conn.Close(); err != nil {
				panic(xerrors.Errorf("error client.Close: %w", err))
			}
		}()

		psc := redis.PubSubConn{Conn: conn}
		err := psc.Subscribe(roomID + ":message")
		if err != nil {
			errCh <- xerrors.Errorf("error in subscribe: %w", err)
		}

		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				var message entity.Message

				if err := json.Unmarshal(v.Data, &message); err != nil {
					errCh <- xerrors.Errorf("error in json.Unmarshal: %w", err)
				}
				select {
				case <-ctx.Done():
					log.Info("parent ctx done!")
					return
				default:
					log.WithFields(log.Fields{
						"roomId":  message.RoomID,
						"userId":  message.UserID,
						"message": message.Message,
					}).Info("send message:")
					ch <- &message
				}
			case redis.Subscription:
				log.WithFields(log.Fields{
					"channel": v.Channel,
					"kind":    v.Kind,
					"count":   v.Count,
				}).Info("redis subscription:")
				select {
				case <-ctx.Done():
					log.Info("parent ctx done!")
					return
				default:
					continue
				}
			case error:
				select {
				case <-ctx.Done():
					log.Info("parent ctx done!")
					return
				default:
					errCh <- xerrors.Errorf("error in subscribe: %w", v.Error())
				}
			}
		}
	}()
	return ch, errCh
}

func (r *gameRepository) ContainWord(word string) bool {
	_, ok := (*r.dictionary)[word]
	return ok
}

func (r *gameRepository) Lock() error {
	key := roomLockKey + ":lock"
	conn := r.conn.Get()
	// while lockできなかったらloop
	for {
		res, err := conn.Do("SETNX", key, "locking")
		if err != nil {
			return xerrors.Errorf("error in lockRoomCandidate: %w", err)
		}

		_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
		if err != nil {
			return xerrors.Errorf("error in lockRoomCandidate expire: %w", err)
		}

		if res != 0 {
			break
		}
	}
	return nil
}

func (r *gameRepository) Unlock() error {
	key := roomLockKey + ":lock"
	conn := r.conn.Get()
	_, err := conn.Do("DEL", key)
	if err != nil {
		return xerrors.Errorf("error in UnlockCurrentWord: %w", err)
	}
	return nil
}

func (r *gameRepository) GetRoomCandidates() ([]string, error) {
	// lock()
	// defer unlock()
	conn := r.conn.Get()

	value, err := redis.Strings(conn.Do("HKEYS", roomLockKey))
	if err != nil {
		return []string{}, xerrors.Errorf("error in GetRoomCandidates(): %w", err)
	}

	_, err = conn.Do("EXPIRE", roomLockKey, int64(r.keyTTL.Seconds()))
	if err != nil {
		return []string{}, xerrors.Errorf("error in GetRoomCandidates() expire: %w", err)
	}

	return value, nil
}

func (r *gameRepository) AddRoomCandidate(roomID string) error {
	conn := r.conn.Get()

	_, err := conn.Do("HSET", roomLockKey, roomID, 0)
	log.WithFields(log.Fields{
		roomID: roomID,
	}).Debug()
	if err != nil {
		return xerrors.Errorf("error in AddRoomCandidate HSET %s %s %d: %w", roomLockKey, roomID, 0, err)
	}
	return nil
}

func (r *gameRepository) DeleteRoomCandidate(roomID string) error {
	conn := r.conn.Get()

	_, err := conn.Do("HDEL", roomLockKey, roomID)

	log.WithFields(log.Fields{
		roomID: roomID,
	}).Debug()

	if err != nil {
		return xerrors.Errorf("error in DeleteRoomCandidate HDEL %s %s %d: %w", roomLockKey, roomID, 0, err)
	}
	return nil
}

// delete <roomID>:**
func (r *gameRepository) DeleteRoom(roomID string) error {
	conn := r.conn.Get()
	// https://blog.morugu.com/entry/2018/01/06/233402
	_, err := conn.Do("EVAL", "return redis.call('DEL', unpack(redis.call('KEYS', ARGV[1])))", 0, roomID+":*")

	log.WithFields(log.Fields{
		roomID: roomID,
	}).Debug()

	if err != nil {
		return xerrors.Errorf("error in DeleteRoom HDEL %s %s %d: %w", roomLockKey, roomID, 0, err)
	}
	return nil
}

type firestoreUserHistory struct {
	Date   time.Time `firestore:"date"`
	Rating int64     `firestore:"rating"`
}

type firestoreUser struct {
	History []firestoreUserHistory `firestore:"history"`
	Name    string                 `firestore:"name"`
	Rating  int64                  `firestore:"rating"`
}

// redis result repo のkeyの命名規則
// <roomID>:<userID>:score

func (r *gameRepository) GetScore(roomID string, userID string) (*entity.Result, error) {
	conn := r.conn.Get()
	key := roomID + ":" + userID + ":" + "score"
	score, err := redis.Int64(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
	if err != nil {
		return nil, xerrors.Errorf("error in GetScore expire: %w", err)
	}

	return &entity.Result{UserID: userID, Score: score}, nil
}

func (r *gameRepository) SetScore(roomID string, userID string, score int64) error {
	conn := r.conn.Get()
	key := roomID + ":" + userID + ":" + "score"
	_, err := conn.Do("SET", key, score)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
	if err != nil {
		return xerrors.Errorf("error in SetScore expire: %w", err)
	}
	return nil
}

func (r *gameRepository) IncrScoreBy(roomID string, userID string, by int64) error {
	conn := r.conn.Get()
	key := roomID + ":" + userID + ":" + "score"
	_, err := conn.Do("INCRBY", key, by)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
	if err != nil {
		return xerrors.Errorf("error in Incr expire: %w", err)
	}
	return nil
}

func (r *gameRepository) GetLatestRating(userID string) (int64, error) {
	// TODO get only users.<id>.rating
	ctx := context.Background()
	client := external.GetFirestore()
	defer func() {
		if err := client.Close(); err != nil {
			panic(xerrors.Errorf("error client.Close: %w", err))
		}
	}()
	snapshot, err := client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		return 0, xerrors.Errorf("GetLatestRating: %w", err)
	}
	data := snapshot.Data()
	rating, ok := data["rating"].(int64)
	if !ok {
		return 0, xerrors.Errorf("Failed to cast interface{} to int64: %w", err)
	}
	return rating, nil
}

func (r *gameRepository) SetRating(userID string, rating int64) error {
	ctx := context.Background()
	client := external.GetFirestore()
	defer func() {
		if err := client.Close(); err != nil {
			panic(xerrors.Errorf("error client.Close: %w", err))
		}
	}()
	_, err := client.Collection("users").Doc(userID).Set(ctx, map[string]interface{}{
		"rating": rating,
	}, firestore.MergeAll)
	if err != nil {
		return xerrors.Errorf("SetRating: %w", err)
	}
	return nil
}

func (r *gameRepository) AddRatingHistory(userID string, rating int64) error {
	ctx := context.Background()
	client := external.GetFirestore()
	defer func() {
		if err := client.Close(); err != nil {
			panic(xerrors.Errorf("error client.Close: %w", err))
		}
	}()

	user := client.Collection("users").Doc(userID)
	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		snapshot, err := user.Get(ctx)
		if err != nil {
			return xerrors.Errorf("AddRatingHistory: %w", err)
		}
		var data firestoreUser
		if err = snapshot.DataTo(&data); err != nil {
			return xerrors.Errorf("error in AddRatingHistory DataTo: %w", err)
		}

		h := firestoreUserHistory{Date: time.Now(), Rating: rating}
		data.History = append(data.History, h)

		_, err = user.Set(ctx, map[string]interface{}{
			"history": data.History,
		}, firestore.MergeAll)

		if err != nil {
			return xerrors.Errorf("AddRatingHistory: %w", err)
		}

		return nil
	})

	if err != nil {
		return xerrors.Errorf("RunTransaction: %w", err)
	}

	return nil
}