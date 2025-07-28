package limiter

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"Rate-Limiter/internal/config"
	"Rate-Limiter/internal/storage"
)

type RequestRateLimiter struct {
	primaryStore  storage.Storage
	fallbackStore storage.Storage
	config        *config.Config
	metrics       atomic.Value // stores *Metrics
}

func NewRequestRateLimiter(primary, fallback storage.Storage, config *config.Config) *RequestRateLimiter {
	r := &RequestRateLimiter{
		primaryStore:  primary,
		fallbackStore: fallback,
		config:        config,
	}
	r.metrics.Store(&Metrics{})
	return r
}

func (r *RequestRateLimiter) Check(ctx context.Context, req *Request) (*Decision, error) {
	start := time.Now()
	metrics := r.GetMetrics()
	metrics.RequestsTotal++

	store := r.primaryStore
	if _, err := store.IncrementAndCheck(ctx, "health", 1, time.Second); err != nil {
		store = r.fallbackStore
	}

	// Get user-specific limits or fall back to defaults
	rps := r.config.RateLimits.DefaultRPS
	if userLimit, exists := r.config.RateLimits.UserLimits[req.UserID]; exists {
		rps = userLimit.RequestsPerSecond
	}

	key := fmt.Sprintf("rate:%s:%s", req.UserID, req.Endpoint)
	allowed, err := store.IncrementAndCheck(ctx, key, rps, time.Second)
	if err != nil {
		return &Decision{Allowed: true}, err // Fail open
	}

	decision := &Decision{Allowed: true, StatusCode: 200}
	if !allowed {
		metrics.RequestsDenied++
		decision = &Decision{
			Allowed:    false,
			StatusCode: 429,
			Reason:     "rate_limit_exceeded",
			RetryAfter: time.Second,
		}
	}

	metrics.AverageLatency = (metrics.AverageLatency + time.Since(start)) / 2
	r.metrics.Store(metrics)
	return decision, nil
}

func (r *RequestRateLimiter) GetMetrics() *Metrics {
	return r.metrics.Load().(*Metrics)
}
