package external

import (
	"os"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

var RedisPool *redis.Pool

// redis のセッティング
func InitRedis() {
	RedisPool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", os.Getenv("REDIS_URL")+":6379")
			if err != nil {
				log.Fatalf("failed init redis %v", err)
			}
			return conn, err
		},
	}
}
