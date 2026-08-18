package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/oopsla5xx/oops-api-v1/internal/config"
)

// NewClient creates a Redis client and verifies connectivity.
func NewClient(ctx context.Context, cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}
