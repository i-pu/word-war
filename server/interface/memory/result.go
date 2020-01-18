package memory

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/external"
	"golang.org/x/xerrors"
	"time"
)

type firestoreUserHistory struct {
	Date time.Time  `firestore:"date"`
	Rating int64 	`firestore:"rating"`
}

type firestoreUser struct {
	History []firestoreUserHistory `firestore:"history"`
	Name    string                 `firestore:"name"`
	Rating  int64                  `firestore:"rating"`
}

// TODO: gameの一時的なスコアをresultリポジトリが管理しており、gameStateが管理すべきな気がするのでrefactorの対象かも
type resultRepository struct {
	firestore *firebase.App
	conn      *redis.Pool
	keyTTL    time.Duration
}

func NewResultRepository() *resultRepository {
	return &resultRepository{
		firestore: external.FirebaseApp,
		conn:      external.RedisPool,
		keyTTL:    time.Minute * 10,
	}
}

// redis
// GetScore
// SetScore
// IncrScoreBy

// firestore
// GetLatestRating
// SetRating

// redis result repo のkeyの命名規則
// <roomID>:<userID>:score

func (r *resultRepository) GetScore(roomID string, userID string) (*entity.Result, error) {
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

func (r *resultRepository) SetScore(roomID string, userID string, score int64) error {
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

func (r *resultRepository) IncrScoreBy(roomID string, userID string, by int64) error {
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

func (r *resultRepository) GetLatestRating(userID string) (int64, error) {
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

func (r *resultRepository) SetRating(userID string, rating int64) error {
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

func (r *resultRepository) AddRatingHistory(userID string, rating int64) error {
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
