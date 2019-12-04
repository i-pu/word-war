package memory

import (
	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/external"
)

type gameRepository struct {
	conn *redis.Pool
}

func NewGameRepository() *gameRepository {
	return &gameRepository{
		conn: external.RedisPool,
	}
}

// <roomID>:currentWord
func (r *gameRepository) UpdateCurrentWord(roomID string, word string) error {
	key := roomID + ":currentWord"
	conn := r.conn.Get()

	_, err := conn.Do("SET", key, word)

	return err
}

func (r *counterRepository) GetCurrentWord(roomID string) (string, error) {
	key := roomID + ":currentWord"
	conn := r.conn.Get()

	value, err := redis.String(conn.Do("GET", key))

	return value, err
}
