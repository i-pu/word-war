package memory

import (
	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/external"
	"golang.org/x/xerrors"
	"time"
)

type counterRepository struct {
	conn *redis.Pool
	keyTTL time.Duration
}

func NewCounterRepository() *counterRepository {
	return &counterRepository{
		conn: external.RedisPool,
		keyTTL: time.Minute * 10, // 10m
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
		return -1, xerrors.Errorf("error in incr counter: %w", err)
	}

	_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
	if err != nil {
		return -1, xerrors.Errorf("error in IncrCounter expire: %w", err)
	}
	return value, nil
}

func (r *counterRepository) SetCounter(roomID string, value int64) error {
	key := roomID + ":counter"
	conn := r.conn.Get()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return xerrors.Errorf("error in set counter: %w", err)
	}

	_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
	if err != nil {
		return xerrors.Errorf("error in set counter expire: %w", err)
	}

	return nil
}

func (r *counterRepository) GetCounter(roomID string) (int64, error) {
	key := roomID + ":counter"
	conn := r.conn.Get()

	value, err := redis.Int64(conn.Do("GET", key))
	if err != nil {
		return -1, xerrors.Errorf("error in main method: %w", err)
	}
	_, err = conn.Do("EXPIRE", key, int64(r.keyTTL.Seconds()))
	if err != nil {
		return -1, xerrors.Errorf("error in GetCounter expire: %w", err)
	}

	return value, nil
}
