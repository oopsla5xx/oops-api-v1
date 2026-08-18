// @title           Oops API
// @version         1.0
// @description     Oops AI-native Software Development Workspace API
// @host            localhost:8080
// @BasePath        /api/v1
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/oopsla5xx/oops-api-v1/internal/app"
	"github.com/oopsla5xx/oops-api-v1/internal/config"
	"github.com/oopsla5xx/oops-api-v1/internal/infrastructure/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	zapLog := logger.New(cfg.App.Env)
	defer func() { _ = zapLog.Sync() }()

	srv, err := app.NewServer(cfg, zapLog)
	if err != nil {
		zapLog.Fatal("failed to initialize server", zap.Error(err))
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := srv.Run(ctx); err != nil {
		zapLog.Fatal("server error", zap.Error(err))
	}
}
