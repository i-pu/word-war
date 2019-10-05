package memory

import (
	"github.com/gomodule/redigo/redis"
	"github.com/i-pu/word-war/server/domain/entity"
	"log"
	"os"
)

type messageRepository struct {
	conn *redis.Pool
}

func NewMessageRepository() *messageRepository {
	redisPool := &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", os.Getenv("REDIS_URL")+":6379")
			if err != nil {
				log.Printf("ERROR: fail init redis pool: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}

	return &messageRepository{
		conn: redisPool,
	}
}

func (r *messageRepository) Publish(message *entity.Message) error {
	return nil
}

func (r *messageRepository) Subscribe(key string) (string, error) {
	return "", nil
}

func (r *messageRepository) Set(key string, value string) error {
	return nil
}
func (r *messageRepository) Get(key string) (string, error) {
	return "", nil
}
