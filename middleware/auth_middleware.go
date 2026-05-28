package middleware

import (
	"strings"

	"user-management-api/models"
	"user-management-api/utils"

	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "user_id"
const ContextUserRoleKey = "user_role"
const ContextUserEmailKey = "user_email"

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "Unauthenticated")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "Unauthenticated")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1], jwtSecret)
		if err != nil {
			utils.Unauthorized(c, "Unauthenticated")
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUserRoleKey, claims.Role)
		c.Set(ContextUserEmailKey, claims.Email)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextUserRoleKey)
		if !exists || role.(string) != models.RoleAdmin {
			utils.Forbidden(c, "Admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetUserID(c *gin.Context) uint {
	userID, _ := c.Get(ContextUserIDKey)
	return userID.(uint)
}
