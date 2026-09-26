package router

import (
	"time"

	"github.com/assethub/assethub/internal/config"
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

func registerDownload(api *gin.RouterGroup, h Handlers, cfg *config.Config, jwtManager *util.JWTManager, rateLimiter *middleware.RateLimiter) {
	downloads := api.Group("/downloads")
	downloads.Use(middleware.Auth(jwtManager))
	downloads.GET("", h.Download.List)

	assets := api.Group("/assets")
	assets.Use(middleware.Auth(jwtManager))
	assets.POST("/:id/download", rateLimiter.Limit("download", int64(cfg.Rate.DownloadPerMinute), time.Minute), h.Download.Download)
}
