package memory

import (
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
	conn := r.conn.Get()
	score, err := redis.Int64(conn.Do("GET", userID+":"+"score"))
	if err != nil {
		return nil, err
	}
	return &entity.Result{UserID: userID, Score: score}, nil
}
func (r *resultRepository) Set(userID string, result *entity.Result) error {
	conn := r.conn.Get()
	_, err := conn.Do("SET", userID+":"+"score", result.Score)
	if err != nil {
		return err
	}
	return nil
}
func (r *resultRepository) IncrBy(userID string, by int64) error {
	conn := r.conn.Get()
	_, err := conn.Do("INCRBY", userID+":"+"score", by)
	if err != nil {
		return err
	}
	return nil
}
