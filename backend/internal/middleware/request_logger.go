package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestLogger logs every request with request_id, method, path, status and latency.
func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.NewString()
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		logger.Info("http request start",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)
		c.Next()
		logger.Info("http request done",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}
