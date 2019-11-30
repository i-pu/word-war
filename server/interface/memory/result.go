package memory

import (
	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/external"
)

type resultRepository struct {
	conn *redis.Pool
}

func NewResultRepository() *resultRepository {
	return &resultRepository{
		conn: external.RedisPool,
	}
}

// redis result repo のkeyの命名規則
// <roomID>:<userID>:score

func (r *resultRepository) Get(roomID string, userID string) (*entity.Result, error) {
	conn := r.conn.Get()
	key := roomID + ":" + userID + ":" + "score"
	score, err := redis.Int64(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return &entity.Result{UserID: userID, Score: score}, nil
}
func (r *resultRepository) Set(result *entity.Result) error {
	conn := r.conn.Get()
	key := result.RoomID + ":" + result.UserID + ":" + "score"
	_, err := conn.Do("SET", key, result.Score)
	if err != nil {
		return err
	}
	return nil
}
func (r *resultRepository) IncrBy(roomID string, userID string, by int64) error {
	conn := r.conn.Get()
	key := roomID + ":" + userID + ":" + "score"
	_, err := conn.Do("INCRBY", key, by)
	if err != nil {
		return err
	}
	return nil
}
