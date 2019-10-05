package memory

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/infra"
)

type messageRepository struct {
	conn *redis.Pool
}

func NewMessageRepository() *messageRepository {
	return &messageRepository{
		conn: infra.RedisPool,
	}
}
func (r *messageRepository) Publish(key string, message *entity.Message) error {
	return errors.New("not implemented")
}

func (r *messageRepository) Subscribe(key string) (string, error) {
	return "", errors.New("not implemented")
}

func (r *messageRepository) SetCounter(value int64) error {
	return errors.New("not implemented")
}
