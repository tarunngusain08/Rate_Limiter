package storage

import (
	"context"
	"sync"
	"time"
)

type MemoryStore struct {
	counters sync.Map
	mutex    sync.RWMutex
}

type counter struct {
	count     int64
	expiresAt time.Time
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) IncrementAndCheck(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	if c, exists := s.counters.Load(key); exists {
		count := c.(*counter)
		if now.After(count.expiresAt) {
			count.count = 1
			count.expiresAt = now.Add(window)
			return true, nil
		}
		count.count++
		return count.count <= int64(limit), nil
	}

	s.counters.Store(key, &counter{
		count:     1,
		expiresAt: now.Add(window),
	})
	return true, nil
}
