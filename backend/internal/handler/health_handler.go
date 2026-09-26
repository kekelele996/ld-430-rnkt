package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

// HealthHandler serves liveness and readiness endpoints.
type HealthHandler struct {
	startedAt time.Time
	mongo     *mongo.Client
	redis     *redis.Client
}

// NewHealthHandler creates a HealthHandler with the dependencies required for readiness.
func NewHealthHandler(mongoClient *mongo.Client, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{startedAt: time.Now(), mongo: mongoClient, redis: redisClient}
}

// Check handles GET /healthz and only reports process liveness.
func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data": gin.H{
			"status":     "healthy",
			"uptime_sec": int(time.Since(h.startedAt).Seconds()),
		},
	})
}

// Ready handles GET /readyz and verifies MongoDB and Redis connectivity.
func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if h.mongo != nil {
		if err := h.mongo.Ping(ctx, nil); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"code":    50300,
				"message": "mongodb unavailable",
				"data":    nil,
			})
			return
		}
	}
	if h.redis != nil {
		if err := h.redis.Ping(ctx).Err(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"code":    50300,
				"message": "redis unavailable",
				"data":    nil,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ready",
		"data": gin.H{
			"status": "ready",
		},
	})
}
