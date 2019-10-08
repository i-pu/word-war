package memory

import (
	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/infra"
)

type counterRepository struct {
	conn *redis.Pool
	// 部屋名固定
	roomName string
}

func NewCounterRepository() *counterRepository {
	return &counterRepository{
		conn: infra.RedisPool,
	}
}

const columnKey = "counter"

// roomID は CA ではどこに書くべき????
func (r *counterRepository) IncrCounter() (int64, error) {
	// <部屋名> / <カラム名> に格納
	// TODO: ラッパー書いたほうがいいかも
	key := "room1" + "/" + columnKey
	conn := r.conn.Get()

	value, err := redis.Int64(conn.Do("GET", key))

	if err != nil {
		return -1, err
	}

	conn.Do("SET", key, value+1)

	return value + 1, nil
}

func (r *counterRepository) SetCounter(value int64) error {
	key := "room1" + "/" + columnKey
	conn := r.conn.Get()

	_, err := conn.Do("SET", key, value)

	return err
}

func (r *counterRepository) GetCounter() (int64, error) {
	key := "room1" + "/" + columnKey
	conn := r.conn.Get()

	value, err := redis.Int64(conn.Do("GET", key))

	return value, err
}
