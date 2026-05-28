package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAPIKeyMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		required   bool
		key        string
		header     string
		headerName string
		wantStatus int
	}{
		{"disabled passes without key", false, "secret", "", "X-API-Key", http.StatusOK},
		{"missing key rejected", true, "secret", "", "X-API-Key", http.StatusUnauthorized},
		{"wrong key rejected", true, "secret", "wrong", "X-API-Key", http.StatusForbidden},
		{"valid key accepted", true, "secret", "secret", "X-API-Key", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(APIKeyMiddleware(tt.required, tt.key, tt.headerName))
			r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.header != "" {
				req.Header.Set(tt.headerName, tt.header)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
