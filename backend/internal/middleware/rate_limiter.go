package middleware

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/assethub/assethub/internal/client"
	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

// RateLimiter enforces per-IP request limits backed by Redis.
type RateLimiter struct {
	redis  *client.RedisClient
	logger *slog.Logger
}

// NewRateLimiter creates a RateLimiter.
func NewRateLimiter(redis *client.RedisClient, logger *slog.Logger) *RateLimiter {
	return &RateLimiter{redis: redis, logger: logger}
}

// Limit creates a middleware limiting a named bucket to max requests per window.
func (r *RateLimiter) Limit(bucket string, max int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("ratelimit:%s:%s", bucket, c.ClientIP())
		count, err := r.redis.IncrWithTTL(c.Request.Context(), key, window)
		if err != nil {
			util.Fail(c, 429, constants.CodeRateLimited, constants.MsgRateLimited)
			return
		}
		if count > max {
			r.logger.Warn(constants.LogRateLimited, "bucket", bucket, "ip", c.ClientIP(), "count", count)
			util.Fail(c, 429, constants.CodeRateLimited, constants.MsgRateLimited)
			return
		}
		c.Next()
	}
}
