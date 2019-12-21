package memory

import (
	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/external"
	"golang.org/x/xerrors"
	"time"
)

type gameStateRepository struct {
	conn *redis.Pool
	keyTTL time.Duration
}

func NewGameStateRepository() *gameStateRepository {
	return &gameStateRepository{
		conn: external.RedisPool,
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

	_, err = conn.Do("EXPIRE", key, r.keyTTL)
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

		_, err = conn.Do("EXPIRE", key, r.keyTTL)
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

	_, err = conn.Do("EXPIRE", key, r.keyTTL)
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

	_, err = conn.Do("EXPIRE", key, r.keyTTL)
	if err != nil {
		return "", xerrors.Errorf("error in getCurrentWord expire: %w", err)
	}

	return value, nil
}
