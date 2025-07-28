package limiter

import (
	"context"
)

// RateLimiter defines the interface for all rate limiting implementations
type RateLimiter interface {
	Check(ctx context.Context, req *Request) (*Decision, error)
	GetMetrics() *Metrics
}
