package middleware

import (
	"strings"

	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

// ContextUserIDKey is the gin context key for the authenticated user id.
const ContextUserIDKey = "user_id"

// ContextRoleKey is the gin context key for the authenticated user role.
const ContextRoleKey = "user_role"

// ContextEmailKey is the gin context key for the authenticated user email.
const ContextEmailKey = "user_email"

// Auth returns a gin middleware that validates the JWT bearer token.
func Auth(jwtManager *util.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			util.Fail(c, 401, constants.CodeUnauthorized, constants.MsgUnauthorized)
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtManager.Verify(token)
		if err != nil {
			util.Fail(c, 401, constants.CodeUnauthorized, constants.MsgUnauthorized)
			return
		}
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextRoleKey, claims.Role)
		c.Set(ContextEmailKey, claims.Email)
		c.Next()
	}
}
