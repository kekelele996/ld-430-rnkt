package router

import (
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

func registerShareLink(api *gin.RouterGroup, h Handlers, jwtManager *util.JWTManager) {
	// Public endpoint: visitors resolve a share token without login.
	api.GET("/shared/:token", h.ShareLink.Resolve)

	// Owner-only management endpoints.
	collections := api.Group("/collections")
	collections.Use(middleware.Auth(jwtManager))
	collections.POST("/:id/share-links", h.ShareLink.Create)
	collections.GET("/:id/share-links", h.ShareLink.List)
	collections.DELETE("/:id/share-links/:linkId", h.ShareLink.Revoke)
}
