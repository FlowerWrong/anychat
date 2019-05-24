package db

import (
	"sync"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

func initRedisClient() error {
	redisOptions, err := redis.ParseURL(viper.Get("redis_url").(string))
	redisClient = redis.NewClient(redisOptions)

	_, err = redisClient.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

// Redis return redis client
func Redis() *redis.Client {
	if redisClient == nil {
		redisOnce.Do(func() {
			err := initRedisClient()
			if err != nil {
				panic(err)
			}
		})
	}
	return redisClient
}
