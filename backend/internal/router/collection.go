package router

import (
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

func registerCollection(api *gin.RouterGroup, h Handlers, jwtManager *util.JWTManager) {
	collections := api.Group("/collections")
	collections.Use(middleware.Auth(jwtManager))
	collections.POST("", h.Collection.Create)
	collections.GET("", h.Collection.List)
	collections.GET("/:id", h.Collection.Get)
	collections.PUT("/:id", h.Collection.Update)
	collections.POST("/:id/assets", h.Collection.AddAsset)
	collections.DELETE("/:id/assets/:assetId", h.Collection.RemoveAsset)
	collections.POST("/:id/members", h.Collection.AddCollaborator)
}
