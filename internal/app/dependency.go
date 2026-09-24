package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/newrelic/go-agent/v3/newrelic"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/oopsla5xx/oops-api-v1/internal/config"
	"github.com/oopsla5xx/oops-api-v1/internal/infrastructure/database"
	nrapp "github.com/oopsla5xx/oops-api-v1/internal/infrastructure/newrelic"
	"github.com/oopsla5xx/oops-api-v1/internal/infrastructure/redis"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/constants"
)

type dependencies struct {
	db       *pgxpool.Pool
	redis    *goredis.Client
	newRelic *newrelic.Application
	log      *zap.Logger
	cfg      *config.Config
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

	nrApp, err := nrapp.NewApplication(&cfg.NewRelic)
	if err != nil {
		db.Close()
		_ = redisClient.Close()
		return nil, fmt.Errorf("new relic: %w", err)
	}
	if nrApp != nil {
		log.Info("new relic agent enabled", zap.String("app_name", cfg.NewRelic.AppName))
	}

	return &dependencies{
		db:       db,
		redis:    redisClient,
		newRelic: nrApp,
		log:      log,
		cfg:      cfg,
	}, nil
}
