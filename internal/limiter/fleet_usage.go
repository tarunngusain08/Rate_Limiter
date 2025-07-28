package limiter

import (
	"Rate-Limiter/internal/config"
	"Rate-Limiter/internal/storage"
	"context"
)

// NewFleetUsageLimiter is a stub implementation.
// Replace with actual logic as needed.
func NewFleetUsageLimiter(primary, fallback storage.Storage, config *config.Config) RateLimiter {
	// TODO: Implement actual fleet usage limiter
	return &dummyLimiter{}
}

type dummyLimiter struct{}

func (d *dummyLimiter) Check(ctx context.Context, req *Request) (*Decision, error) {
	return &Decision{Allowed: true, StatusCode: 200}, nil
}
func (d *dummyLimiter) GetMetrics() *Metrics {
	return &Metrics{}
}
