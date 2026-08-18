package constants

import "time"

const (
	DefaultRequestTimeout  = 30 * time.Second
	DefaultShutdownTimeout = 30 * time.Second
	DefaultDBTimeout       = 10 * time.Second
	DefaultRedisTimeout    = 5 * time.Second
)
