package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"user-management-api/utils"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret"

	token, _, err := utils.GenerateToken(1, "admin@test.com", "admin", secret, 24)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "missing header",
			authHeader: "",
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "Unauthenticated",
		},
		{
			name:       "invalid format",
			authHeader: "Token abc",
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "Unauthenticated",
		},
		{
			name:       "invalid token",
			authHeader: "Bearer invalid.token.here",
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "Unauthenticated",
		},
		{
			name:       "valid token",
			authHeader: "Bearer " + token,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(AuthMiddleware(secret))
			r.GET("/protected", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status: got %d want %d body=%s", w.Code, tt.wantStatus, w.Body.String())
			}

			if tt.wantMsg != "" {
				var resp struct {
					Message string `json:"message"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp.Message != tt.wantMsg {
					t.Errorf("message: got %q want %q", resp.Message, tt.wantMsg)
				}
			}
		})
	}
}

func TestAdminMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		role       string
		setRole    bool
		wantStatus int
	}{
		{"admin allowed", "admin", true, http.StatusOK},
		{"user denied", "user", true, http.StatusForbidden},
		{"missing role denied", "", false, http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(func(c *gin.Context) {
				if tt.setRole {
					c.Set(ContextUserRoleKey, tt.role)
				}
				c.Next()
			}, AdminMiddleware())
			r.GET("/admin", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/admin", nil))

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
