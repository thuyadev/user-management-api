package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"user-management-api/controllers"
	"user-management-api/database"
	"user-management-api/models"
	"user-management-api/repositories"
	"user-management-api/routes"
	"user-management-api/services"
	"user-management-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	TestJWTSecret = "test-jwt-secret"
	TestAPIKey    = "test-api-key"
	TestPassword  = "password123"
)

type TestApp struct {
	Router *gin.Engine
	DB     *gorm.DB
	Admin  *models.User
	User   *models.User
}

func SetupTestApp(t *testing.T) *TestApp {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := database.ConnectSQLite()
	if err != nil {
		t.Fatalf("connect sqlite: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	hash, err := utils.HashPassword(TestPassword)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	admin := &models.User{
		Name:     "Admin User",
		Email:    "admin@test.com",
		Password: hash,
		Role:     models.RoleAdmin,
	}
	regularUser := &models.User{
		Name:     "John Doe",
		Email:    "user@test.com",
		Password: hash,
		Role:     models.RoleUser,
	}

	userRepo := repositories.NewUserRepository(db)
	if err := userRepo.Create(admin); err != nil {
		t.Fatalf("seed admin: %v", err)
	}
	if err := userRepo.Create(regularUser); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	cfg := &utils.Config{
		JWTSecret:      TestJWTSecret,
		JWTExpiryHours: 24,
		APIKeyRequired: true,
		APIKey:         TestAPIKey,
		APIKeyHeader:   "X-API-Key",
		AIEnabled:      false,
		SwaggerEnabled: false,
	}

	logSvc := NewMockLogService()
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)

	authService := services.NewAuthService(userRepo, logSvc, cfg.JWTSecret, cfg.JWTExpiryHours)
	userService := services.NewUserService(userRepo, logSvc)
	categoryService := services.NewCategoryService(categoryRepo, logSvc)
	productService := services.NewProductService(productRepo, categoryRepo, logSvc)
	aiService := services.NewAIService(cfg)

	handlers := &routes.Handlers{
		Auth:     controllers.NewAuthController(authService),
		User:     controllers.NewUserController(userService),
		Category: controllers.NewCategoryController(categoryService, aiService),
		Product:  controllers.NewProductController(productService, aiService),
		Log:      controllers.NewLogController(logSvc),
	}

	router := gin.New()
	routes.Setup(router, handlers, cfg)

	return &TestApp{
		Router: router,
		DB:     db,
		Admin:  admin,
		User:   regularUser,
	}
}

func (a *TestApp) Request(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reader *bytes.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		reader = bytes.NewReader(payload)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", TestAPIKey)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	a.Router.ServeHTTP(w, req)
	return w
}

func (a *TestApp) LoginToken(t *testing.T, email, password string) string {
	t.Helper()

	w := a.Request(http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    email,
		"password": password,
	}, "")

	if w.Code != http.StatusOK {
		t.Fatalf("login %s: status %d body %s", email, w.Code, w.Body.String())
	}

	var resp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if resp.Data.Token == "" {
		t.Fatal("expected token in login response")
	}
	return resp.Data.Token
}

func DecodeJSON(t *testing.T, w *httptest.ResponseRecorder, dest interface{}) {
	t.Helper()
	if err := json.Unmarshal(w.Body.Bytes(), dest); err != nil {
		t.Fatalf("decode json: %v body=%s", err, w.Body.String())
	}
}

func APIResponseMessage(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	var resp struct {
		Message string `json:"message"`
	}
	DecodeJSON(t, w, &resp)
	return resp.Message
}
