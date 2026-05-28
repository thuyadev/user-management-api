package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware restricts which browser origins may call the API.
// Non-browser clients (Postman, curl, server-to-server) are not affected by CORS.
func CORSMiddleware(allowedOrigins []string, apiKeyHeader string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		o = strings.TrimSpace(o)
		if o != "" {
			allowed[o] = true
		}
	}

	defaultAllowHeaders := buildAllowHeaders(apiKeyHeader)

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		if len(allowed) == 0 || !allowed[origin] {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Vary", "Origin")

		// Echo requested headers on preflight so custom headers (e.g. X-Api-Key) are allowed.
		if requested := c.GetHeader("Access-Control-Request-Headers"); requested != "" {
			c.Header("Access-Control-Allow-Headers", requested)
		} else {
			c.Header("Access-Control-Allow-Headers", defaultAllowHeaders)
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func buildAllowHeaders(apiKeyHeader string) string {
	headers := []string{"Content-Type", "Authorization", "X-Api-Key", "X-API-Key"}
	if apiKeyHeader != "" && apiKeyHeader != "X-Api-Key" && apiKeyHeader != "X-API-Key" {
		headers = append(headers, apiKeyHeader)
	}
	return strings.Join(headers, ", ")
}
