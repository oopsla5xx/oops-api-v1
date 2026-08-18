package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	App       AppConfig
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	CORS      CORSConfig
	RateLimit RateLimitConfig
}

type AppConfig struct {
	Name  string `env:"APP_NAME,required"`
	Env   string `env:"APP_ENV,required"`
	Debug bool   `env:"APP_DEBUG,required"`
}

type ServerConfig struct {
	Host              string        `env:"SERVER_HOST,required"`
	Port              string        `env:"SERVER_PORT,required"`
	ShutdownTimeout   time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT,required"`
	ReadHeaderTimeout time.Duration `env:"SERVER_READ_HEADER_TIMEOUT,required"`
	ReadTimeout       time.Duration `env:"SERVER_READ_TIMEOUT,required"`
	WriteTimeout      time.Duration `env:"SERVER_WRITE_TIMEOUT,required"`
}

type DatabaseConfig struct {
	DSN             string        `env:"DATABASE_DSN,required"`
	MaxConns        int32         `env:"DATABASE_MAX_CONNS,required"`
	MinConns        int32         `env:"DATABASE_MIN_CONNS,required"`
	MaxConnLifetime time.Duration `env:"DATABASE_MAX_CONN_LIFETIME,required"`
	MaxConnIdleTime time.Duration `env:"DATABASE_MAX_CONN_IDLE_TIME,required"`
}

type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR,required"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB,required"`
}

type CORSConfig struct {
	AllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS,required"`
}

type RateLimitConfig struct {
	RequestsPerSecond float64 `env:"RATE_LIMIT_RPS,required"`
	Burst             int     `env:"RATE_LIMIT_BURST,required"`
}

// Load reads all configuration from environment variables.
// Every variable tagged required must be explicitly set — there are no silent defaults.
// A missing variable or a malformed value causes Load to return an error immediately.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
