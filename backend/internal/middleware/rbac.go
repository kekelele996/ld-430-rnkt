package middleware

import (
	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

// Roles returns a middleware requiring one of the given roles.
func Roles(roles ...constants.Role) gin.HandlerFunc {
	allowed := make(map[constants.Role]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		role := c.GetString(ContextRoleKey)
		if !allowed[constants.Role(role)] {
			util.Fail(c, 403, constants.CodeForbidden, constants.MsgForbidden)
			return
		}
		c.Next()
	}
}
