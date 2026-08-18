// Package tests provides helpers for integration tests that require
// a real database or Redis connection. Import only from _test.go files.
package tests

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/oopsla5xx/oops-api-v1/internal/config"
	"github.com/oopsla5xx/oops-api-v1/internal/infrastructure/database"
)

// NewTestDB returns a pgxpool connected to the test database.
// The test is skipped automatically if DATABASE_DSN is not set,
// so unit tests run fine without a database available.
// The pool is closed when the test ends.
func NewTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		t.Skip("DATABASE_DSN not set — skipping integration test")
	}

	cfg := &config.DatabaseConfig{
		DSN:             dsn,
		MaxConns:        5,
		MinConns:        1,
		MaxConnLifetime: 30 * time.Minute,
		MaxConnIdleTime: 10 * time.Minute,
	}

	pool, err := database.NewPool(context.Background(), cfg)
	if err != nil {
		t.Fatalf("connect test database: %v", err)
	}

	t.Cleanup(pool.Close)
	return pool
}

// Truncate removes all rows from the given tables and resets their sequences.
// Call this in t.Cleanup (or at the start of each test) to isolate test state.
//
// Example:
//
//	pool := tests.NewTestDB(t)
//	t.Cleanup(func() { tests.Truncate(t, pool, "users", "workspaces") })
func Truncate(t *testing.T, pool *pgxpool.Pool, tables ...string) {
	t.Helper()

	if len(tables) == 0 {
		return
	}

	query := fmt.Sprintf(
		"TRUNCATE TABLE %s RESTART IDENTITY CASCADE",
		strings.Join(tables, ", "),
	)
	if _, err := pool.Exec(context.Background(), query); err != nil {
		t.Fatalf("truncate tables %v: %v", tables, err)
	}
}
