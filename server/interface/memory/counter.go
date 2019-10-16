package memory

import (
	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/infra"
)

type counterRepository struct {
	conn *redis.Pool
	// 部屋名固定
	// roomName string
}

func NewCounterRepository() *counterRepository {
	return &counterRepository{
		conn: infra.RedisPool,
	}
}

// redis counter repo の命名規則
// incr counter
// 将来は <roomID>:counter になるかも

func (r *counterRepository) IncrCounter() (int64, error) {
	key := "counter"
	conn := r.conn.Get()

	value, err := redis.Int64(conn.Do("Incr", key))

	if err != nil {
		return -1, err
	}

	return value, nil
}

func (r *counterRepository) SetCounter(value int64) error {
	key := "counter"
	conn := r.conn.Get()

	_, err := conn.Do("SET", key, value)

	return err
}

func (r *counterRepository) GetCounter() (int64, error) {
	key := "counter"
	conn := r.conn.Get()

	value, err := redis.Int64(conn.Do("GET", key))

	return value, err
}
