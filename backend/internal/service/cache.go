package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/assethub/assethub/internal/client"
)

// CacheClient is a thin Redis-backed JSON cache used for hot assets and search results.
type CacheClient struct {
	redis *client.RedisClient
}

// NewCacheClient creates a CacheClient.
func NewCacheClient(redis *client.RedisClient) *CacheClient {
	return &CacheClient{redis: redis}
}

// GetJSON fetches a JSON value into dst. Returns false when missing.
func (c *CacheClient) GetJSON(ctx context.Context, key string, dst interface{}) (bool, error) {
	raw, err := c.redis.Get(ctx, key)
	if err != nil {
		return false, nil
	}
	if err := json.Unmarshal([]byte(raw), dst); err != nil {
		return false, nil
	}
	return true, nil
}

// SetJSON stores a JSON value with TTL.
func (c *CacheClient) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, key, raw, ttl)
}
