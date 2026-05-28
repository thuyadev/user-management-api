package routes

import (
	"user-management-api/middleware"
	"user-management-api/policies"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(admin *gin.RouterGroup, h *Handlers) {
	perm := middleware.PermissionMiddleware(policies.PermUsersManage)

	admin.GET("/users", perm, h.User.List)
	admin.POST("/users", perm, h.User.Create)
	admin.GET("/users/:id", perm, h.User.Get)
	admin.PUT("/users/:id", perm, h.User.Update)
	admin.DELETE("/users/:id", perm, h.User.Delete)
}
