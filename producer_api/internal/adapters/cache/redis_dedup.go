package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisDedup implements ports.DedupStore using Redis SETNX with a TTL.
type RedisDedup struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisDedup(client *redis.Client, ttl time.Duration) *RedisDedup {
	return &RedisDedup{client: client, ttl: ttl}
}

// IsDuplicate atomically records key. It returns true when the key already
// existed (i.e. SETNX did not set it), meaning the message is a duplicate.
func (r *RedisDedup) IsDuplicate(ctx context.Context, key string) (bool, error) {
	set, err := r.client.SetNX(ctx, key, 1, r.ttl).Result()
	if err != nil {
		return false, fmt.Errorf("redis setnx: %w", err)
	}
	// set == true  -> first time we see it -> not a duplicate
	// set == false -> key already present  -> duplicate
	return !set, nil
}
