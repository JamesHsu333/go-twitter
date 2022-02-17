package redis

import (
	"time"

	"github.com/JamesHsu333/go-twitter/config"
	"github.com/go-redis/redis/v8"
)

// Return new redis client
func NewRedisClient(cfg *config.Config) *redis.Client {
	redisHost := cfg.Redis.RedisAddr

	if redisHost == "" {
		redisHost = ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		MinIdleConns: cfg.Redis.MinIdleConns,
		PoolSize:     cfg.Redis.PoolSize,
		PoolTimeout:  time.Duration(cfg.Redis.PoolTimeout) * time.Second,
		Password:     cfg.Redis.RedisPassword,
		DB:           cfg.Redis.DB,
	})

	return client
}
