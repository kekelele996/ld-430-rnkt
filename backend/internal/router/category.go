package router

import (
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

func registerCategory(api *gin.RouterGroup, h Handlers, jwtManager *util.JWTManager) {
	categories := api.Group("/categories")
	categories.Use(middleware.Auth(jwtManager))
	categories.GET("", h.Category.Tree)

	manage := categories.Group("")
	manage.Use(middleware.Roles("Admin", "Moderator"))
	manage.POST("", h.Category.Create)
	manage.PUT("/:id", h.Category.Update)
	manage.DELETE("/:id", h.Category.Delete)
}
