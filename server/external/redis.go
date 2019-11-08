package external

import (
	"log"
	"os"

	"github.com/gomodule/redigo/redis"
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
				log.Printf("ERROR: fail init redis pool: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
}
