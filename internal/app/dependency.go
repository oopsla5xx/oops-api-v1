package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/oopsla5xx/oops-api-v1/internal/config"
	"github.com/oopsla5xx/oops-api-v1/internal/infrastructure/database"
	"github.com/oopsla5xx/oops-api-v1/internal/infrastructure/redis"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/constants"
)

type dependencies struct {
	db    *pgxpool.Pool
	redis *goredis.Client
	log   *zap.Logger
	cfg   *config.Config
}

func newDependencies(cfg *config.Config, log *zap.Logger) (*dependencies, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultDBTimeout)
	defer cancel()

	db, err := database.NewPool(ctx, &cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}
	log.Info("database connected")

	redisClient, err := redis.NewClient(ctx, &cfg.Redis)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("redis: %w", err)
	}
	log.Info("redis connected")

	return &dependencies{
		db:    db,
		redis: redisClient,
		log:   log,
		cfg:   cfg,
	}, nil
}
