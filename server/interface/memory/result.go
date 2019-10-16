package memory

import (
	"errors"

	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/infra"
)

type resultRepository struct {
	conn *redis.Pool
	// roomID  string
}

func NewResultRepository() *resultRepository {
	return &resultRepository{
		conn: infra.RedisPool,
		// roomName:  "room1",
		// columnKey: "result",
	}
}

// redis result repo の命名規則
// <userID>:score
// 将来は <roomID>:<userID>:scoreみたいな感じになるかも

func (r *resultRepository) Get(userID string) (*entity.Result, error) {
	return nil, errors.New("unimplemented")
}
func (r *resultRepository) Set(userID string, result *entity.Result) error {
	return errors.New("unimplemented")
}
func (r *resultRepository) IncrBy(userID string, by int64) error {
	conn := r.conn.Get()
	_, err := conn.Do("INCRBY", userID+":"+"score", by)
	if err != nil {
		return err
	}
	return nil
}
