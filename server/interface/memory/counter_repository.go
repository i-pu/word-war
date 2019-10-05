package memory

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/infra"
)

type counterRepository struct {
	conn *redis.Pool
}

func NewCounterRepository() *counterRepository {
	return &counterRepository{
		conn: infra.RedisPool,
	}
}

func (r *counterRepository) IncrCounter() (int64, error) {
	// conn := r.conn.Get()
	return 3, errors.New("not implemented")
}

func (r *counterRepository) SetCounter(value int64) error {
	return errors.New("not implemented")
}

func (r *counterRepository) GetCounter() (int64, error) {
	return 3, errors.New("not implemented")
}
