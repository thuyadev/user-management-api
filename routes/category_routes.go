package routes

import (
	"user-management-api/middleware"
	"user-management-api/policies"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(admin *gin.RouterGroup, h *Handlers) {
	view := middleware.PermissionMiddleware(policies.PermCategoriesView)
	manage := middleware.PermissionMiddleware(policies.PermCategoriesManage)

	admin.GET("/categories", view, h.Category.List)
	admin.GET("/categories/:id", view, h.Category.Get)
	admin.POST("/categories", manage, h.Category.Create)
	admin.PUT("/categories/:id", manage, h.Category.Update)
	admin.DELETE("/categories/:id", manage, h.Category.Delete)
	admin.POST("/categories/ai/suggest", manage, h.Category.SuggestName)
}
