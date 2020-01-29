package external

import (
	"os"

	"github.com/go-redis/redis"
	"golang.org/x/xerrors"
)

var RedisClient *redis.Client

// redis のセッティング
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:        os.Getenv("REDIS_URL") + ":6379",
		MaxConnAge:  12000,
		IdleTimeout: 80,
	})
}

func HealthCheck() error {
	_, err := RedisClient.Ping().Result()
	if err != nil {
		return xerrors.Errorf("Redis init error: %w", err)
	}
	return nil
}
