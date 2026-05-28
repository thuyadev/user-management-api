package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"user-management-api/utils"

	"github.com/gin-gonic/gin"
)

func TestSetupCORSPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	cfg := &utils.Config{
		CORSAllowedOrigins: []string{"http://localhost:3000"},
		APIKeyHeader:       "X-UMA-a394985d00e67ddf",
		APIKeyRequired:     false,
		SwaggerEnabled:     false,
	}
	Setup(router, &Handlers{}, cfg)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/login", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "content-type, x-api-key, x-uma-a394985d00e67ddf")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("expected allow-origin, got %q", got)
	}
	allowHeaders := strings.ToLower(w.Header().Get("Access-Control-Allow-Headers"))
	if !strings.Contains(allowHeaders, "x-api-key") {
		t.Fatalf("expected x-api-key in allow-headers, got %q", allowHeaders)
	}
}
