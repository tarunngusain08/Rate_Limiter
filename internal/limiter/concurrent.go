package limiter

import (
	"context"
	"sync/atomic"
	"time"
)

type ConcurrentLimiter struct {
	maxConcurrent int64
	current       int64
	metrics       atomic.Value // stores *Metrics
}

func NewConcurrentLimiter(maxConcurrent int64) *ConcurrentLimiter {
	l := &ConcurrentLimiter{
		maxConcurrent: maxConcurrent,
	}
	l.metrics.Store(&Metrics{})
	return l
}

func (c *ConcurrentLimiter) Check(ctx context.Context, req *Request) (*Decision, error) {
	start := time.Now()
	metrics := c.GetMetrics()
	metrics.RequestsTotal++

	if atomic.LoadInt64(&c.current) >= c.maxConcurrent {
		metrics.RequestsDenied++
		metrics.AverageLatency = (metrics.AverageLatency + time.Since(start)) / 2
		c.metrics.Store(metrics)
		return &Decision{
			Allowed:    false,
			StatusCode: 429,
			Reason:     "concurrent_limit_exceeded",
			RetryAfter: time.Second,
		}, nil
	}

	atomic.AddInt64(&c.current, 1)
	defer atomic.AddInt64(&c.current, -1)

	metrics.AverageLatency = (metrics.AverageLatency + time.Since(start)) / 2
	c.metrics.Store(metrics)
	return &Decision{Allowed: true, StatusCode: 200}, nil
}

func (c *ConcurrentLimiter) GetMetrics() *Metrics {
	return c.metrics.Load().(*Metrics)
}
