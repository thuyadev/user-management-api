package routes

import (
	"user-management-api/middleware"
	"user-management-api/utils"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(api *gin.RouterGroup, h *Handlers, cfg *utils.Config) {
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/login", h.Auth.Login)
		authRoutes.POST("/register", h.Auth.Register)
		authRoutes.GET("/roles", h.Auth.Roles)
		authRoutes.GET("/me", middleware.AuthMiddleware(cfg.JWTSecret), h.Auth.Me)
	}
}
