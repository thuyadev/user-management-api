package routes

import (
	"user-management-api/middleware"
	"user-management-api/policies"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(admin *gin.RouterGroup, h *Handlers) {
	view := middleware.PermissionMiddleware(policies.PermProductsView)
	manage := middleware.PermissionMiddleware(policies.PermProductsManage)

	admin.GET("/products", view, h.Product.List)
	admin.GET("/products/:id", view, h.Product.Get)
	admin.POST("/products", manage, h.Product.Create)
	admin.PUT("/products/:id", manage, h.Product.Update)
	admin.DELETE("/products/:id", manage, h.Product.Delete)
	admin.POST("/products/ai/description", manage, h.Product.GenerateDescription)
}
