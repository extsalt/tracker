package pubsub

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// NewClient creates a new Redis client connection
func NewClient() (*redis.Client, error) {
	option := redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	redisClient := redis.NewClient(&option)

	ctx := context.Background()
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return redisClient, nil
}
