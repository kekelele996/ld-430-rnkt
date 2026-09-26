package router

import (
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

func registerReview(api *gin.RouterGroup, h Handlers, jwtManager *util.JWTManager) {
	reviews := api.Group("/assets")
	reviews.Use(middleware.Auth(jwtManager), middleware.Roles("Admin", "Moderator"))
	reviews.POST("/:id/reviews", h.Review.Review)

	admin := api.Group("")
	admin.Use(middleware.Auth(jwtManager), middleware.Roles("Admin"))
	admin.GET("/audit-logs", h.Audit.List)
}
