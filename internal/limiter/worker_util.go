package limiter

import (
	"Rate-Limiter/internal/config"
	"context"
	"sync/atomic"
	"time"
)

type WorkerUtilLimiter struct {
	activeWorkers int64
	maxWorkers    int64
	config        *config.Config
	metrics       atomic.Value // stores *Metrics
}

func NewWorkerUtilLimiter(maxWorkers int64, config *config.Config) *WorkerUtilLimiter {
	l := &WorkerUtilLimiter{
		maxWorkers: maxWorkers,
		config:     config,
	}
	l.metrics.Store(&Metrics{})
	return l
}

func (w *WorkerUtilLimiter) Check(ctx context.Context, req *Request) (*Decision, error) {
	start := time.Now()
	metrics := w.GetMetrics()
	metrics.RequestsTotal++

	utilization := float64(atomic.LoadInt64(&w.activeWorkers)) / float64(w.maxWorkers)

	var decision *Decision
	switch {
	case utilization >= w.config.WorkerThresholds.Emergency && req.Priority != Critical:
		metrics.RequestsDenied++
		decision = &Decision{
			Allowed:    false,
			StatusCode: 503,
			Reason:     "emergency_worker_utilization",
			RetryAfter: time.Second * 30,
		}
	case utilization >= w.config.WorkerThresholds.High && req.Priority == TestMode:
		metrics.RequestsDenied++
		decision = &Decision{
			Allowed:    false,
			StatusCode: 429,
			Reason:     "high_worker_utilization",
			RetryAfter: time.Second * 10,
		}
	default:
		atomic.AddInt64(&w.activeWorkers, 1)
		decision = &Decision{Allowed: true, StatusCode: 200}
	}

	metrics.AverageLatency = (metrics.AverageLatency + time.Since(start)) / 2
	w.metrics.Store(metrics)
	return decision, nil
}

func (w *WorkerUtilLimiter) GetMetrics() *Metrics {
	return w.metrics.Load().(*Metrics)
}

func (w *WorkerUtilLimiter) ReleaseWorker() {
	atomic.AddInt64(&w.activeWorkers, -1)
}
