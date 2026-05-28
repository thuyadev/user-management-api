package routes

import (
	"net/http"

	"user-management-api/middleware"
	"user-management-api/utils"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, h *Handlers, cfg *utils.Config) {
	cors := middleware.CORSMiddleware(cfg.CORSAllowedOrigins, cfg.APIKeyHeader)

	router.Use(middleware.SecurityHeaders(), cors)

	// Gin only runs route middleware for registered methods. Browsers send OPTIONS
	// preflight before POST/PUT, so register explicit OPTIONS handlers.
	router.OPTIONS("/api/v1/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	RegisterHealthRoutes(router)
	RegisterSwaggerRoutes(router, cfg)

	api := router.Group("/api/v1")
	api.Use(middleware.APIKeyMiddleware(cfg.APIKeyRequired, cfg.APIKey, cfg.APIKeyHeader))

	RegisterAuthRoutes(api, h, cfg)

	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	RegisterUserRoutes(admin, h)
	RegisterCategoryRoutes(admin, h)
	RegisterProductRoutes(admin, h)
	RegisterLogRoutes(admin, h)
}
