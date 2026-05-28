package middleware

import (
	"crypto/subtle"
	"net/http"

	"user-management-api/utils"

	"github.com/gin-gonic/gin"
)

// APIKeyMiddleware rejects requests without a valid X-API-Key header.
// This blocks casual browser access and ensures only trusted clients (Postman, mobile, backend) call the API.
func APIKeyMiddleware(required bool, expectedKey, headerName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !required || expectedKey == "" {
			c.Next()
			return
		}

		if headerName == "" {
			headerName = "X-API-Key"
		}

		provided := c.GetHeader(headerName)
		if provided == "" {
			utils.Error(c, http.StatusUnauthorized, "API key required", gin.H{
				"header": headerName,
				"hint":   "Send your API key in the request header",
			})
			c.Abort()
			return
		}

		if !secureCompare(provided, expectedKey) {
			utils.Error(c, http.StatusForbidden, "Invalid API key", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func secureCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
