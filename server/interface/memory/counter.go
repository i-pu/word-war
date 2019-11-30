package memory

import (
	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/external"
)

type counterRepository struct {
	conn *redis.Pool
}

func NewCounterRepository() *counterRepository {
	return &counterRepository{
		conn: external.RedisPool,
	}
}

// redis counter repo の命名規則
// incr counter
// 将来は <roomID>:counter になるかも

func (r *counterRepository) IncrCounter(roomID string) (int64, error) {
	key := roomID + ":counter"
	conn := r.conn.Get()

	value, err := redis.Int64(conn.Do("Incr", key))

	if err != nil {
		return -1, err
	}

	return value, nil
}

func (r *counterRepository) SetCounter(roomID string, value int64) error {
	key := roomID + ":counter"
	conn := r.conn.Get()

	_, err := conn.Do("SET", key, value)

	return err
}

func (r *counterRepository) GetCounter(roomID string) (int64, error) {
	key := roomID + ":counter"
	conn := r.conn.Get()

	value, err := redis.Int64(conn.Do("GET", key))

	return value, err
}
