package router

import (
	"time"

	"github.com/assethub/assethub/internal/config"
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

func registerAsset(api *gin.RouterGroup, h Handlers, cfg *config.Config, jwtManager *util.JWTManager, rateLimiter *middleware.RateLimiter) {
	assets := api.Group("/assets")
	assets.Use(middleware.Auth(jwtManager))
	assets.GET("", h.Asset.List)
	assets.GET("/hot", h.Asset.Hot)
	assets.GET("/:id", h.Asset.Get)
	assets.GET("/:id/reviews", h.Review.ListByAsset)

	assets.POST("/upload", rateLimiter.Limit("upload", int64(cfg.Rate.UploadPerHour), time.Hour), h.Asset.Upload)
	assets.PUT("/:id", h.Asset.Update)
	assets.POST("/:id/publish", h.Asset.Publish)
	assets.POST("/:id/archive", h.Asset.Archive)
}
