package limiter

import (
	"context"
	"sync"
	"time"

	"Rate-Limiter/internal/config"
	"Rate-Limiter/internal/constants"
	"Rate-Limiter/internal/storage"
)

type Manager struct {
	limiters    []RateLimiter
	redisStore  storage.Storage
	memoryStore storage.Storage
	metrics     *Metrics
	mutex       sync.RWMutex
}

func NewManager(config *config.Config) *Manager {
	redisStore := storage.NewRedisStore(config.Redis.Addr, config.Redis.Password, config.Redis.DB)
	memoryStore := storage.NewMemoryStore()

	return &Manager{
		limiters: []RateLimiter{
			NewRequestRateLimiter(redisStore, memoryStore, config),
			NewWorkerUtilLimiter(constants.DefaultWorkerCount, config),
			NewFleetUsageLimiter(redisStore, memoryStore, config),
		},
		redisStore:  redisStore,
		memoryStore: memoryStore,
		metrics:     &Metrics{},
	}
}

func (m *Manager) CheckAllLimiters(ctx context.Context, req *Request) (*Decision, error) {
	start := time.Now()

	for _, limiter := range m.limiters {
		decision, err := limiter.Check(ctx, req)
		if err != nil {
			continue // Fail open
		}
		if !decision.Allowed {
			m.updateMetrics(start, false)
			return decision, nil
		}
	}

	m.updateMetrics(start, true)
	return &Decision{Allowed: true, StatusCode: 200}, nil
}

func (m *Manager) updateMetrics(start time.Time, allowed bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.metrics.RequestsTotal++
	if !allowed {
		m.metrics.RequestsDenied++
	}
	m.metrics.AverageLatency = (m.metrics.AverageLatency + time.Since(start)) / 2
}

// Add metrics aggregation
func (m *Manager) GetMetrics() *Metrics {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	aggregated := &Metrics{}
	for _, limiter := range m.limiters {
		metrics := limiter.GetMetrics()
		aggregated.RequestsTotal += metrics.RequestsTotal
		aggregated.RequestsDenied += metrics.RequestsDenied
		// Average the latencies
		if metrics.AverageLatency > 0 {
			aggregated.AverageLatency += metrics.AverageLatency
		}
	}

	if len(m.limiters) > 0 {
		aggregated.AverageLatency /= time.Duration(len(m.limiters))
	}

	return aggregated
}
