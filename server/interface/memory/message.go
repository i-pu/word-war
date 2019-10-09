package memory

import (
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/infra"
)

type messageRepository struct {
	conn *redis.Pool
	// roomName  string
	columnKey string
}

func NewMessageRepository() *messageRepository {
	return &messageRepository{
		conn: infra.RedisPool,
		// roomName:  "room1",
	}
}

const messageKey = "message"

func (r *messageRepository) Publish(roomID string, message *entity.Message) error {
	fmt.Println("[Message/Publish] %+v %+v", roomID, message)

	// <https://godoc.org/github.com/gomodule/redigo/redis#pkg-examples>
	// Pub Sub どうやるんだ
	return errors.New("not implemented")
}

func (r *messageRepository) Subscribe(roomID string) (string, error) {
	return "", errors.New("not implemented")
}

func (r *messageRepository) SetCounter(value int64) error {
	return errors.New("not implemented")
}
