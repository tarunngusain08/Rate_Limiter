package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr, password string, db int) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisStore{client: client}
}

func (s *RedisStore) IncrementAndCheck(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	pipe := s.client.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return true, err // Fail open
	}

	return incr.Val() <= int64(limit), nil
}

func (s *RedisStore) Close() error {
	return s.client.Close()
}
