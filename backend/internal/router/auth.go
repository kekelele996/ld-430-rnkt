package router

import (
	"time"

	"github.com/assethub/assethub/internal/config"
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

func registerAuth(api *gin.RouterGroup, h Handlers, cfg *config.Config, jwtManager *util.JWTManager, rateLimiter *middleware.RateLimiter) {
	auth := api.Group("/auth")
	// Sensitive authentication endpoints are rate limited per IP.
	auth.POST("/register", rateLimiter.Limit("auth-register", 10, time.Minute), h.Auth.Register)
	auth.POST("/login", rateLimiter.Limit("auth-login", 10, time.Minute), h.Auth.Login)

	protected := api.Group("")
	protected.Use(middleware.Auth(jwtManager))
	protected.GET("/auth/me", h.Auth.Me)
}
