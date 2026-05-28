package routes

import (
	"user-management-api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterHealthRoutes(router *gin.Engine) {
	router.GET("/health", controllers.HealthCheck)
}
