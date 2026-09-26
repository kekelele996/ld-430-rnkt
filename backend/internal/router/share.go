package router

import (
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

func registerShare(api *gin.RouterGroup, h Handlers, jwtManager *util.JWTManager) {
	// Public: visitors open a shared collection with the token, no login.
	// Revoked or expired tokens are rejected inside the service with 403.
	api.GET("/shared/:token", h.Share.Resolve)

	links := api.Group("/collections")
	links.Use(middleware.Auth(jwtManager))
	links.POST("/:id/share-links", h.Share.Create)
	links.GET("/:id/share-links", h.Share.List)
	links.DELETE("/:id/share-links/:linkId", h.Share.Revoke)
}
