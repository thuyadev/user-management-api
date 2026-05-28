package main

import (
	"log"
	"os"

	"user-management-api/controllers"
	"user-management-api/database"
	"user-management-api/database/seeders"
	_ "user-management-api/docs"
	"user-management-api/repositories"
	"user-management-api/routes"
	"user-management-api/services"
	"user-management-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title           User Management API
// @version         1.0
// @description     REST API with JWT auth, role permissions, categories, products, and MongoDB activity logs.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  admin@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-UMA-a394985d00e67ddf
// @description API key required for all /api/v1 routes

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token. Format: Bearer {token}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := utils.LoadConfig()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	if os.Getenv("SKIP_SEED") != "true" {
		if err := seeders.RunAll(db, cfg); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
	}

	mongoClient, err := database.ConnectMongo(cfg)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}
	defer mongoClient.Disconnect(nil)

	logCollection := database.GetMongoCollection(mongoClient, cfg.MongoDB, "user_logs")

	userRepo := repositories.NewUserRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	logRepo := repositories.NewLogRepository(logCollection)

	logService := services.NewLogService(logRepo)
	authService := services.NewAuthService(userRepo, logService, cfg.JWTSecret, cfg.JWTExpiryHours)
	userService := services.NewUserService(userRepo, logService)
	categoryService := services.NewCategoryService(categoryRepo, logService)
	productService := services.NewProductService(productRepo, categoryRepo, logService)
	aiService := services.NewAIService(cfg)

	handlers := &routes.Handlers{
		Auth:     controllers.NewAuthController(authService),
		User:     controllers.NewUserController(userService),
		Category: controllers.NewCategoryController(categoryService, aiService),
		Product:  controllers.NewProductController(productService, aiService),
		Log:      controllers.NewLogController(logService),
	}

	router := gin.Default()
	if cfg.APIKeyRequired && cfg.APIKey == "" {
		log.Fatal("API_KEY_REQUIRED is true but API_KEY is not set")
	}
	if cfg.SwaggerEnabled && (cfg.SwaggerUser == "" || cfg.SwaggerPassword == "") {
		log.Fatal("SWAGGER_ENABLED is true but SWAGGER_USER or SWAGGER_PASSWORD is not set")
	}

	routes.Setup(router, handlers, cfg)

	if cfg.SwaggerEnabled {
		log.Printf("Swagger UI available at http://localhost:%s%s/index.html (basic auth required)", cfg.AppPort, cfg.SwaggerPath)
	}

	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
