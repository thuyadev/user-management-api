package routes

import (
	"net/http"
	"strings"

	"user-management-api/middleware"
	"user-management-api/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterSwaggerRoutes(router *gin.Engine, cfg *utils.Config) {
	if !cfg.SwaggerEnabled {
		return
	}

	path := strings.TrimSuffix(cfg.SwaggerPath, "/")
	if path == "" {
		path = "/swagger"
	}

	basicAuth := middleware.BasicAuthMiddleware(cfg.SwaggerUser, cfg.SwaggerPassword)
	handler := ginSwagger.WrapHandler(swaggerFiles.Handler)

	router.GET(path, basicAuth, func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, path+"/index.html")
	})

	router.GET(path+"/*any", basicAuth, handler)
}
