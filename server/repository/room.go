package repository

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"math/rand"
	"os"
	"time"
	"unicode/utf8"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/external"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type RoomRepository interface {
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
	GetPlayer(roomID string, userID string) *entity.Player

	// message
	Publish(message *entity.Message) error
	Subscribe(ctx context.Context, roomID string) (<-chan *entity.Message, <-chan error)

	// room candidates
	GetRoomCandidateIDs() ([]string, error)
	AddRoomCandidateID(roomID string) error
	DeleteRoomCandidateID(roomID string) error

	// room
	GetRoom(roomID string) (*entity.Room, error)
	CreateRoom() (*entity.Room, error)
	CleanRoom(player *entity.Room) error

	// timer
	PublishTimer(room *entity.Room) error
	SubscribeTimer(room *entity.Room) (context.Context, error)

	// score
	GetScore(player *entity.Player) (*entity.Result, error)
	SetScore(player *entity.Player, score int64) error
	IncrScoreBy(player *entity.Player, by int64) error

	// rating
	GetLatestRating(userID string) (int64, error)
	SetRating(userID string, rating int64) error
	AddRatingHistory(userID string, rating int64) error
}

type roomRepository struct {
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

type firestoreUserHistory struct {
	Date   time.Time `firestore:"date"`
	Rating int64     `firestore:"rating"`
}

type firestoreUser struct {
	History []firestoreUserHistory `firestore:"history"`
	Name    string                 `firestore:"name"`
	Rating  int64                  `firestore:"rating"`
}

func NewRoomRepository() *roomRepository {
	dictionary := map[string]struct{}{}

	if dicPath, ok := os.LookupEnv("DIC_PATH"); !ok {
		log.WithError(xerrors.Errorf("NewRoomRepository LookupEnv(DIC_PATH) is false.")).Warn()
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
				log.WithError(xerrors.Errorf("NewRoomRepository utf8.ValidString is false : %w", err)).Fatalf("line: %s", line)
			}
		}
		if err := scanner.Err(); err != nil {
			log.WithError(xerrors.Errorf("NewRoomRepository scanner error: %w", err)).Fatal("")
		}
	}

	return &roomRepository{
		firestore:  external.FirebaseApp,
		conn:       external.RedisClient,
		keyTTL:     time.Minute * 10,
		dictionary: &dictionary,
	}
}

func (r *roomRepository) InitWord(roomID string, word string) error {
	key := roomID + ":currentWord"
	err := r.conn.Set(key, word, r.keyTTL).Err()
	if err != nil {
		return xerrors.Errorf("error in InitWord setnx: %w", err)
	}
	return err
}

func (r *roomRepository) LockCurrentWord(roomID string) error {
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
func (r *roomRepository) UnlockCurrentWord(roomID string) error {
	key := roomID + ":currentWord:lock"
	err := r.conn.Del(key).Err()
	if err != nil {
		return xerrors.Errorf("error in UnlockCurrentWord: %w", err)
	}
	return nil
}

// <roomID>:currentWord
func (r *roomRepository) UpdateCurrentMessage(message *entity.Message) error {
	key := message.RoomID + ":currentWord"
	err := r.conn.Set(key, message.Message, r.keyTTL).Err()
	if err != nil {
		return xerrors.Errorf("error in Set(%s, %s, %d): %w", key, message.Message, r.keyTTL, err)
	}
	return nil
}

func (r *roomRepository) GetCurrentMessage(roomID string) (*entity.Message, error) {
	key := roomID + ":currentWord"
	word, err := r.conn.Get(key).Result()
	if err != nil {
		return nil, xerrors.Errorf("error in Get(%s): %w", key, err)
	}
	return &entity.Message{RoomID: roomID, UserID: "fixme", Message: word}, nil
}

func (r *roomRepository) AddPlayer(player *entity.Player) error {
	key := player.RoomID + ":users"
	err := r.conn.HSet(key, player.UserID, 0).Err()

	if err != nil {
		return xerrors.Errorf("error in AddPlayer %+v: %w", player, err)
	}

	log.Debugf("AddPlayer: %+v", player)
	return nil
}

func (r *roomRepository) GetUserIDs(roomID string) ([]string, error) {
	key := roomID + ":users"
	users, err := r.conn.HKeys(key).Result()
	if err != nil {
		return nil, xerrors.Errorf("error in GetPlayerIDs(%s): %w", roomID, err)
	}
	return users, nil
}

// TODO: ゆくゆくはPlayerの情報をredisに保存するかもしれない
func (r *roomRepository) GetPlayer(roomID string, userID string) *entity.Player {
	return &entity.Player{
		UserID: userID,
		RoomID: roomID,
	}
}

// delete '<room>:*'
func (r *roomRepository) CleanRoom(room *entity.Room) error {
	log.Debugf("called CleanRoom. room: %+v", room)
	// delete wildcard
	for _, k := range r.conn.Keys(room.RoomID + ":*").Val() {
		if err := r.conn.Del(k).Err(); err != nil {
			return err
		}
	}

	if err := r.conn.HDel(roomLockKey, room.RoomID).Err(); err != nil {
		return err
	}

	log.Debugf("done CleanRoom. room: %+v", room)
	return nil
}

func (r *roomRepository) Publish(message *entity.Message) error {
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

func (r *roomRepository) Subscribe(ctx context.Context, roomID string) (<-chan *entity.Message, <-chan error) {
	// TODO: Subscribeのより良いやり方あるかも
	// ctx: 親のcontextで親のcontextが終了するとgo func()内でctx.Done()により終了する
	// roomID: どこの部屋のイベントをsubscribeするか
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

func (r *roomRepository) ContainWord(word string) bool {
	_, ok := (*r.dictionary)[word]
	return ok
}

func (r *roomRepository) Lock(key string) error {
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

func (r *roomRepository) Unlock(key string) error {
	err := r.conn.Del(key + ":lock").Err()
	if err != nil {
		return xerrors.Errorf("error in UnlockCurrentWord: %w", err)
	}
	return nil
}

func (r *roomRepository) GetRoomCandidateIDs() ([]string, error) {
	value, err := r.conn.HKeys(roomLockKey).Result()
	if err != nil {
		return []string{}, xerrors.Errorf("error in GetRoomCandidateIDs(): %w", err)
	}

	return value, nil
}

func (r *roomRepository) AddRoomCandidateID(roomID string) error {
	_, err := r.conn.HSet(roomLockKey, roomID, 0).Result()
	if err != nil {
		return xerrors.Errorf("error in AddRoomCandidateID HSet(%s, %s, %d): %w", roomLockKey, roomID, 0, err)
	}

	log.WithFields(log.Fields{
		roomID: roomID,
	}).Debug()

	return nil
}

func (r *roomRepository) DeleteRoomCandidateID(roomID string) error {
	_, err := r.conn.HDel(roomLockKey, roomID).Result()
	if err != nil {
		return xerrors.Errorf("error in DeleteRoomCandidateID HDel(%s, %s): %w", roomLockKey, roomID, err)
	}

	log.WithFields(log.Fields{
		roomID: roomID,
	}).Debug()

	return nil
}

// roomID -> room
func (r *roomRepository) GetRoom(roomID string) (*entity.Room, error) {
	mes, err := r.GetCurrentMessage(roomID)

	if err != nil {
		return nil, xerrors.Errorf("GetRoom(%s) error: %w", roomID, err)
	}

	return &entity.Room{RoomID: roomID, CurrentMessage: mes}, nil
}

func (r *roomRepository) CreateRoom() (*entity.Room, error) {
	roomID := fmt.Sprintf("%d", rand.Intn(90000)+10000)

	initialMessage := &entity.Message{
		Message: "しりとり",
		UserID:  "unknown",
		RoomID:  roomID,
	}

	room := &entity.Room{RoomID: roomID, CurrentMessage: initialMessage}

	err := r.InitWord(room.RoomID, room.CurrentMessage.Message)
	if err != nil {
		return nil, xerrors.Errorf("error in CreateRoom %+v: %w", room, err)
	}

	log.Debugf("CreateRoom ID: %v", roomID)

	return room, nil
}

func (r *roomRepository) GetScore(player *entity.Player) (*entity.Result, error) {
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

func (r *roomRepository) SetScore(player *entity.Player, score int64) error {
	key := player.RoomID + ":" + player.UserID + ":" + "score"
	_, err := r.conn.Set(key, score, r.keyTTL).Result()
	if err != nil {
		return xerrors.Errorf("error in Set(%s, %d, %d): %w", key, score, r.keyTTL, err)
	}

	return nil
}

func (r *roomRepository) IncrScoreBy(player *entity.Player, by int64) error {
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

func (r *roomRepository) GetLatestRating(userID string) (int64, error) {
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

func (r *roomRepository) SetRating(userID string, rating int64) error {
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

func (r *roomRepository) AddRatingHistory(userID string, rating int64) error {
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

func (r *roomRepository) PublishTimer(room *entity.Room) error {
	key := room.RoomID + ":timer"

	err := r.conn.Publish(key, 0).Err()

	if err != nil {
		return xerrors.Errorf("error in redis publish: %w", err)
	}

	err = r.conn.Expire(key, r.keyTTL).Err()
	if err != nil {
		return xerrors.Errorf("error in Publish expire: %w", err)
	}

	log.WithFields(log.Fields{
		"roomId": room.RoomID,
	}).Info("PublishTimer")

	return nil
}

func (r *roomRepository) SubscribeTimer(room *entity.Room) (context.Context, error) {
	pubSub := r.conn.Subscribe(room.RoomID + ":timer")
	_, err := pubSub.Receive()
	if err != nil {
		panic(err)
	}

	pubSubChan := pubSub.Channel()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer func() {
			if err := pubSub.Close(); err != nil {
				panic(xerrors.Errorf("error in subscribe: %w", err))
			}
		}()

		for {
			select {
			case v, ok := <-pubSubChan:
				if !ok {
					log.Warn("pubSubChan is closed")
					return
				}
				if v != nil {
					log.Debug(v.Payload)
				} else {
					log.Debug("v is nil")
				}
				cancel()
				return
			}
		}
	}()
	return ctx, nil
}
