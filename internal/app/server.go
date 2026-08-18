package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/oopsla5xx/oops-api-v1/internal/config"
)

type Server struct {
	httpServer *http.Server
	db         *pgxpool.Pool
	redis      *goredis.Client
	log        *zap.Logger
	cfg        *config.Config
}

func NewServer(cfg *config.Config, log *zap.Logger) (*Server, error) {
	deps, err := newDependencies(cfg, log)
	if err != nil {
		return nil, fmt.Errorf("initialize dependencies: %w", err)
	}

	router := newRouter(deps)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:           router,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
	}

	return &Server{
		httpServer: httpServer,
		db:         deps.db,
		redis:      deps.redis,
		log:        log,
		cfg:        cfg,
	}, nil
}

// Run starts the HTTP server and blocks until ctx is cancelled (SIGTERM/SIGINT).
func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		s.log.Info("server starting", zap.String("addr", s.httpServer.Addr))
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		s.log.Info("shutdown signal received")
	}

	return s.shutdown()
}

func (s *Server) shutdown() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.Server.ShutdownTimeout)
	defer cancel()

	s.log.Info("shutting down server", zap.Duration("timeout", s.cfg.Server.ShutdownTimeout))

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	s.db.Close()
	s.log.Info("database connection pool closed")

	if err := s.redis.Close(); err != nil {
		s.log.Warn("redis close error", zap.Error(err))
	}
	s.log.Info("redis connection closed")

	s.log.Info("server stopped gracefully")
	return nil
}
