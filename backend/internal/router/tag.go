package router

import (
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

func registerTag(api *gin.RouterGroup, h Handlers, jwtManager *util.JWTManager) {
	tags := api.Group("/tags")
	tags.Use(middleware.Auth(jwtManager))
	tags.GET("", h.Tag.TagCloud)

	manage := tags.Group("")
	manage.Use(middleware.Roles("Admin", "Moderator"))
	manage.POST("", h.Tag.Create)
}
