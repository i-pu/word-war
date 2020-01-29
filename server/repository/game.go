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
	"github.com/go-redis/redis"
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/external"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type GameRepository interface {
	// util
	Lock(key string) error
	Unlock(key string) error

	// word
	InitWord(roomID string, word string) error
	UpdateCurrentMessage(message *entity.Message) error
	GetCurrentMessage(roomID string) (*entity.Message, error)
	ContainWord(word string) bool

	// user
	AddPlayer(player *entity.Player) error
	GetUserIDs(roomID string) ([]string, error)
	DeletePlayer(player *entity.Player) error

	// counter
	IncrCounter(roomID string) (int64, error)
	SetCounter(roomID string, value int64) error
	GetCounter(roomID string) (int64, error)

	// message
	Publish(message *entity.Message) error
	Subscribe(ctx context.Context, roomID string) (<-chan *entity.Message, <-chan error)

	// room candidates
	GetRoomCandidates() ([]string, error)
	AddRoomCandidate(roomID string) error
	DeleteRoomCandidate(roomID string) error

	// room
	DeleteRoom(roomID string) error
	CleanGame(player *entity.Player) error

	// score
	GetScore(player *entity.Player) (*entity.Result, error)
	SetScore(player *entity.Player, score int64) error
	IncrScoreBy(player *entity.Player, by int64) error

	// rating
	GetLatestRating(userID string) (int64, error)
	SetRating(userID string, rating int64) error
	AddRatingHistory(userID string, rating int64) error
}

type gameRepository struct {
	firestore  *firebase.App
	conn       *redis.Client
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
		conn:       external.RedisClient,
		keyTTL:     time.Minute * 10,
		dictionary: &dictionary,
	}
}

func (r *gameRepository) InitWord(roomID string, word string) error {
	key := roomID + ":currentWord"
	err := r.conn.Set(key, word, r.keyTTL).Err()
	if err != nil {
		return xerrors.Errorf("error in InitWord setnx: %w", err)
	}
	return err
}

func (r *gameRepository) LockCurrentWord(roomID string) error {
	key := roomID + ":currentWord:lock"

	for {
		ok, err := r.conn.SetNX(key, "locking", r.keyTTL).Result()
		if err != nil {
			return xerrors.Errorf("error in lockCurrentWord: %w", err)
		}
		if ok {
			break
		}
	}
	return nil
}

// updateして r.conn.Do("delete roomID+"currentWord:lock")
func (r *gameRepository) UnlockCurrentWord(roomID string) error {
	key := roomID + ":currentWord:lock"
	err := r.conn.Del(key).Err()
	if err != nil {
		return xerrors.Errorf("error in UnlockCurrentWord: %w", err)
	}
	return nil
}

// <roomID>:currentWord
func (r *gameRepository) UpdateCurrentMessage(message *entity.Message) error {
	key := message.RoomID + ":currentWord"
	err := r.conn.Set(key, message.Message, r.keyTTL).Err()
	if err != nil {
		return xerrors.Errorf("error in Set(%s, %s, %d): %w", key, message.Message, r.keyTTL, err)
	}
	return nil
}

func (r *gameRepository) GetCurrentMessage(roomID string) (*entity.Message, error) {
	key := roomID + ":currentWord"
	word, err := r.conn.Get(key).Result()
	if err != nil {
		return nil, xerrors.Errorf("error in Get(%s): %w", key, err)
	}
	return &entity.Message{RoomID: roomID, UserID: "fixme", Message: word}, nil
}

func (r *gameRepository) AddPlayer(player *entity.Player) error {
	key := player.RoomID + ":users"
	err := r.conn.HSet(key, player.UserID, 0).Err()
	log.Debugf("%+v", player)

	if err != nil {
		return xerrors.Errorf("error in AddPlayer %+v: %w", player, err)
	}
	return nil
}

func (r *gameRepository) GetUserIDs(roomID string) ([]string, error) {
	key := roomID + ":users"
	users, err := r.conn.HKeys(key).Result()
	if err != nil {
		return nil, xerrors.Errorf("error in GetPlayerIDs(%s): %w", roomID, err)
	}
	return users, nil
}

func (r *gameRepository) DeletePlayer(player *entity.Player) error {
	key := player.RoomID + ":users"
	err := r.conn.HDel(key, player.UserID).Err()
	if err != nil {
		return xerrors.Errorf("error in DeletePlayer %+v: %w", player, err)
	}
	return nil
}

func (r *gameRepository) CleanGame(player *entity.Player) error {
	p := r.conn.TxPipeline()

	key := player.RoomID + ":users"
	if err := p.HDel(key, player.UserID).Err(); err != nil {
		return err
	}

	users, err := p.HKeys(player.RoomID).Result()
	if err != nil {
		return err
	}

	if len(users) == 0 {
		// delete wildcard
		for _, k := range p.Keys(player.RoomID + ":*").Val() {
			if err := p.Del(k).Err(); err != nil {
				return err
			}
		}

		if err := p.HDel(roomLockKey, player.RoomID).Err(); err != nil {
			return err
		}
	}

	_, err = p.Exec()

	if err != nil {
		return xerrors.Errorf("error in CleanGame: %w", err)
	}
	return nil
}

func (r *gameRepository) IncrCounter(roomID string) (int64, error) {
	key := roomID + ":counter"
	value, err := r.conn.Incr(key).Result()
	if err != nil {
		return -1, xerrors.Errorf("error in incr counter: %w", err)
	}

	err = r.conn.Expire(key, r.keyTTL).Err()
	if err != nil {
		return -1, xerrors.Errorf("error in IncrCounter expire: %w", err)
	}
	return value, nil
}

func (r *gameRepository) SetCounter(roomID string, value int64) error {
	key := roomID + ":counter"
	err := r.conn.Set(key, value, r.keyTTL).Err()
	if err != nil {
		return xerrors.Errorf("error in set counter: %w", err)
	}
	return nil
}

func (r *gameRepository) GetCounter(roomID string) (int64, error) {
	key := roomID + ":counter"
	value, err := r.conn.Get(key).Int64()
	if err != nil {
		return -1, xerrors.Errorf("error in main method: %w", err)
	}
	err = r.conn.Expire(key, r.keyTTL).Err()
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

	rep, err := r.conn.Publish(message.RoomID+":message", mesBytes).Result()

	if err != nil {
		return xerrors.Errorf("error in redis publish: %w", err)
	}

	err = r.conn.Expire(message.RoomID+":message", r.keyTTL).Err()
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
	pubSub := r.conn.Subscribe(roomID + ":message")
	_, err := pubSub.Receive()
	if err != nil {
		panic(err)
	}

	pubSubChan := pubSub.Channel()
	messageChan := make(chan *entity.Message)
	errChan := make(chan error)
	go func() {
		defer func() {
			if err := pubSub.Close(); err != nil {
				panic(xerrors.Errorf("error in subscribe: %w", err))
			}
			close(messageChan)
			close(errChan)
		}()

		for {
			select {
			case <-ctx.Done():
				log.WithFields(log.Fields{
					"roomId": roomID,
				}).Infof("game finish")
				return
			case v, ok := <-pubSubChan:
				if !ok {
					return
				}
				if v != nil {
					log.WithFields(log.Fields{
						"v.String()": v.String(),
						"v.Payload":  v.Payload,
						"v.Channel":  v.Channel,
						"v.Pattern":  v.Pattern,
					}).Debug("value of pubSub message:")
				} else {
					log.Debug("value of pubSub message is nil!!!")
				}
				var message entity.Message
				if err := json.Unmarshal([]byte(v.Payload), &message); err != nil {
					errChan <- xerrors.Errorf("error in json.Unmarshal(%s): %w", v.String(), err)
				}
				log.WithFields(log.Fields{
					"roomId":  message.RoomID,
					"userId":  message.UserID,
					"message": message.Message,
				}).Info("send message:")
				messageChan <- &message
			}
		}
	}()

	return messageChan, errChan
}

func (r *gameRepository) ContainWord(word string) bool {
	_, ok := (*r.dictionary)[word]
	return ok
}

func (r *gameRepository) Lock(key string) error {
	// while lockできなかったらloop
	for {
		ok, err := r.conn.SetNX(key+":lock", "locking", r.keyTTL).Result()
		if err != nil {
			return xerrors.Errorf("error in lockRoomCandidate: %w", err)
		}
		if ok {
			break
		}
	}
	return nil
}

func (r *gameRepository) Unlock(key string) error {
	err := r.conn.Del(key + ":lock").Err()
	if err != nil {
		return xerrors.Errorf("error in UnlockCurrentWord: %w", err)
	}
	return nil
}

func (r *gameRepository) GetRoomCandidates() ([]string, error) {
	value, err := r.conn.HKeys(roomLockKey).Result()
	if err != nil {
		return []string{}, xerrors.Errorf("error in GetRoomCandidates(): %w", err)
	}

	return value, nil
}

func (r *gameRepository) AddRoomCandidate(roomID string) error {
	_, err := r.conn.HSet(roomLockKey, roomID, 0).Result()
	if err != nil {
		return xerrors.Errorf("error in AddRoomCandidate HSet(%s, %s, %d): %w", roomLockKey, roomID, 0, err)
	}

	log.WithFields(log.Fields{
		roomID: roomID,
	}).Debug()

	return nil
}

func (r *gameRepository) DeleteRoomCandidate(roomID string) error {
	_, err := r.conn.HDel(roomLockKey, roomID).Result()
	if err != nil {
		return xerrors.Errorf("error in DeleteRoomCandidate HDel(%s, %s): %w", roomLockKey, roomID, err)
	}

	log.WithFields(log.Fields{
		roomID: roomID,
	}).Debug()

	return nil
}

// delete <roomID>:**
func (r *gameRepository) DeleteRoom(roomID string) error {
	// https://blog.morugu.com/entry/2018/01/06/233402
	_, err := r.conn.Eval("return redis.call('DEL', unpack(redis.call('KEYS', ARGV[1])))", []string{}, roomID+":*").Result()
	if err != nil {
		return xerrors.Errorf("error in DeleteRoom Eval: %w", err)
	}

	log.WithFields(log.Fields{
		roomID: roomID,
	}).Debug()

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

func (r *gameRepository) GetScore(player *entity.Player) (*entity.Result, error) {
	key := player.RoomID + ":" + player.UserID + ":" + "score"
	score, err := r.conn.Get(key).Int64()
	if err != nil {
		return nil, err
	}

	_, err = r.conn.Expire(key, r.keyTTL).Result()
	if err != nil {
		return nil, xerrors.Errorf("error in GetScore expire: %w", err)
	}

	return &entity.Result{UserID: player.UserID, Score: score}, nil
}

func (r *gameRepository) SetScore(player *entity.Player, score int64) error {
	key := player.RoomID + ":" + player.UserID + ":" + "score"
	_, err := r.conn.Set(key, score, r.keyTTL).Result()
	if err != nil {
		return xerrors.Errorf("error in Set(%s, %d, %d): %w", key, score, r.keyTTL, err)
	}

	return nil
}

func (r *gameRepository) IncrScoreBy(player *entity.Player, by int64) error {
	key := player.RoomID + ":" + player.UserID + ":" + "score"
	_, err := r.conn.IncrBy(key, by).Result()
	if err != nil {
		return xerrors.Errorf("error in IncrBy(%s, %d): %w", key, by, err)
	}
	_, err = r.conn.Expire(key, r.keyTTL).Result()
	if err != nil {
		return xerrors.Errorf("error in Expire(%s, %d): %w", key, r.keyTTL, err)
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
