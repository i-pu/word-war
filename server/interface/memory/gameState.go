package memory

import (
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/external"
	"golang.org/x/xerrors"
)

type gameStateRepository struct {
	conn   *redis.Pool
	keyTTL time.Duration
}

func NewGameStateRepository() *gameStateRepository {
	return &gameStateRepository{
		conn:   external.RedisPool,
		keyTTL: time.Minute * 10,
	}
}
func (r *gameStateRepository) InitWord(roomID string, word string) error {
	key := roomID + ":currentWord"
	conn := r.conn.Get()
	_, err := conn.Do("SETNX", key, word)
	if err != nil {
		return xerrors.Errorf("error in InitWord setnx: %w", err)
	}

	_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
	if err != nil {
		return xerrors.Errorf("error in InitWord expire: %w", err)
	}

	return err
}

func (r *gameStateRepository) LockCurrentWord(roomID string) error {
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
func (r *gameStateRepository) UnlockCurrentWord(roomID string) error {
	key := roomID + ":currentWord:lock"
	conn := r.conn.Get()
	_, err := conn.Do("DEL", key)
	if err != nil {
		return xerrors.Errorf("error in UnlockCurrentWord: %w", err)
	}
	return nil
}

// <roomID>:currentWord
func (r *gameStateRepository) UpdateCurrentWord(roomID string, word string) error {
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

func (r *gameStateRepository) GetCurrentWord(roomID string) (string, error) {
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

func (r *gameStateRepository) AddUser(roomID string, userID string) error {
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

func (r *gameStateRepository) GetUsers(roomID string) ([]string, error) {
	key := roomID + ":users"
	conn := r.conn.Get()

	users, err := redis.Strings(conn.Do("HKEYS", key))
	if err != nil {
		return nil, xerrors.Errorf("error in GetUsers HKEYS %s, 0, -1: %w", key, err)
	}
	return users, nil
}
