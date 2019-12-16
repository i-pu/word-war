package memory

import (
	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/external"
	"log"
)

type gameStateRepository struct {
	conn *redis.Pool
}

func NewGameStateRepository() *gameStateRepository {
	return &gameStateRepository{
		conn: external.RedisPool,
	}
}
func (r *gameStateRepository) InitWord(roomID string, word string) error {
	key := roomID + ":currentWord"
	conn := r.conn.Get()
	_, err := conn.Do("SETNX", key, word)
	return err
}

func (r *gameStateRepository) LockCurrentWord(roomID string) error {
	key := roomID + ":currentWord:lock"
	conn := r.conn.Get()
	// while lockできなかったらloop
	for {
		res, err := conn.Do("SETNX", key, "locking")
		if err != nil {
			log.Println("lockCurrentWord error:", err)
			return err
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
		log.Println("UnlockCurrentWord error:", err)
		return err
	}
	return nil
}

// <roomID>:currentWord
func (r *gameStateRepository) UpdateCurrentWord(roomID string, word string) error {
	key := roomID + ":currentWord"
	conn := r.conn.Get()

	_, err := conn.Do("SET", key, word)

	if err != nil {
		log.Println("UpdateCurrentWord error:", err)
		return err
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
		log.Println("GetCurrentWord error")
		return "", err
	}

	return value, nil
}
