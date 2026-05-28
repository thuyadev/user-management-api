package middleware

import (
	"crypto/subtle"

	"user-management-api/utils"

	"github.com/gin-gonic/gin"
)

func BasicAuthMiddleware(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, ok := c.Request.BasicAuth()
		if !ok ||
			subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			c.Header("WWW-Authenticate", `Basic realm="Swagger UI"`)
			utils.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
