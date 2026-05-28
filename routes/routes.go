package routes

import (
	"user-management-api/auth"
	"user-management-api/controllers"
	"user-management-api/middleware"
	"user-management-api/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	Auth     *controllers.AuthController
	User     *controllers.UserController
	Category *controllers.CategoryController
	Product  *controllers.ProductController
	Log      *controllers.LogController
}

func Setup(router *gin.Engine, h *Handlers, cfg *utils.Config) {
	router.Use(middleware.SecurityHeaders())
	router.GET("/health", controllers.HealthCheck)

	if cfg.SwaggerEnabled {
		swagger := router.Group(cfg.SwaggerPath)
		swagger.Use(middleware.BasicAuthMiddleware(cfg.SwaggerUser, cfg.SwaggerPassword))
		swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	api := router.Group("/api/v1")
	api.Use(
		middleware.CORSMiddleware(cfg.CORSAllowedOrigins),
		middleware.APIKeyMiddleware(cfg.APIKeyRequired, cfg.APIKey, cfg.APIKeyHeader),
	)

	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/login", h.Auth.Login)
		authRoutes.POST("/register", h.Auth.Register)
		authRoutes.GET("/roles", h.Auth.Roles)
		authRoutes.GET("/me", middleware.AuthMiddleware(cfg.JWTSecret), h.Auth.Me)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Users — admin only
		admin.GET("/users", middleware.PermissionMiddleware(auth.PermUsersManage), h.User.List)
		admin.POST("/users", middleware.PermissionMiddleware(auth.PermUsersManage), h.User.Create)
		admin.GET("/users/:id", middleware.PermissionMiddleware(auth.PermUsersManage), h.User.Get)
		admin.PUT("/users/:id", middleware.PermissionMiddleware(auth.PermUsersManage), h.User.Update)
		admin.DELETE("/users/:id", middleware.PermissionMiddleware(auth.PermUsersManage), h.User.Delete)

		// Categories — view: admin + user, manage: admin only
		admin.GET("/categories", middleware.PermissionMiddleware(auth.PermCategoriesView), h.Category.List)
		admin.GET("/categories/:id", middleware.PermissionMiddleware(auth.PermCategoriesView), h.Category.Get)
		admin.POST("/categories", middleware.PermissionMiddleware(auth.PermCategoriesManage), h.Category.Create)
		admin.PUT("/categories/:id", middleware.PermissionMiddleware(auth.PermCategoriesManage), h.Category.Update)
		admin.DELETE("/categories/:id", middleware.PermissionMiddleware(auth.PermCategoriesManage), h.Category.Delete)
		admin.POST("/categories/ai/suggest", middleware.PermissionMiddleware(auth.PermCategoriesManage), h.Category.SuggestName)

		// Products — view: admin + user, manage: admin only
		admin.GET("/products", middleware.PermissionMiddleware(auth.PermProductsView), h.Product.List)
		admin.GET("/products/:id", middleware.PermissionMiddleware(auth.PermProductsView), h.Product.Get)
		admin.POST("/products", middleware.PermissionMiddleware(auth.PermProductsManage), h.Product.Create)
		admin.PUT("/products/:id", middleware.PermissionMiddleware(auth.PermProductsManage), h.Product.Update)
		admin.DELETE("/products/:id", middleware.PermissionMiddleware(auth.PermProductsManage), h.Product.Delete)
		admin.POST("/products/ai/description", middleware.PermissionMiddleware(auth.PermProductsManage), h.Product.GenerateDescription)

		// Logs — admin only
		admin.GET("/logs", middleware.PermissionMiddleware(auth.PermLogsView), h.Log.List)
	}
}
