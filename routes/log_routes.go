package routes

import (
	"user-management-api/middleware"
	"user-management-api/policies"

	"github.com/gin-gonic/gin"
)

func RegisterLogRoutes(admin *gin.RouterGroup, h *Handlers) {
	admin.GET("/logs/stats/events", middleware.PermissionMiddleware(policies.PermLogsView), h.Log.EventStats)
	admin.GET("/logs/stats/daily", middleware.PermissionMiddleware(policies.PermLogsView), h.Log.DailyStats)
	admin.GET("/logs", middleware.PermissionMiddleware(policies.PermLogsView), h.Log.List)
}
