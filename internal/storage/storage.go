package storage

import (
	"context"
	"time"
)

// Storage defines the interface for rate limiting storage backends
type Storage interface {
	// IncrementAndCheck increments the counter for the given key and checks if it's within limits
	// Returns true if the request is allowed, false if the limit has been exceeded
	IncrementAndCheck(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
}
