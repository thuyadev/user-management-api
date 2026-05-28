package middleware

import (
	"user-management-api/auth"
	"user-management-api/utils"

	"github.com/gin-gonic/gin"
)

func PermissionMiddleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString(ContextUserRoleKey)
		if !auth.HasPermission(role, permission) {
			utils.Forbidden(c, "You do not have permission to perform this action")
			c.Abort()
			return
		}
		c.Next()
	}
}
