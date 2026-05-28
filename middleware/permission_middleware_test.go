package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"user-management-api/policies"

	"github.com/gin-gonic/gin"
)

func TestPermissionMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		role       string
		permission string
		wantStatus int
	}{
		{
			name:       "admin can manage users",
			role:       "admin",
			permission: policies.PermUsersManage,
			wantStatus: http.StatusOK,
		},
		{
			name:       "user cannot manage users",
			role:       "user",
			permission: policies.PermUsersManage,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "user can view categories",
			role:       "user",
			permission: policies.PermCategoriesView,
			wantStatus: http.StatusOK,
		},
		{
			name:       "user cannot manage categories",
			role:       "user",
			permission: policies.PermCategoriesManage,
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(func(c *gin.Context) {
				c.Set(ContextUserRoleKey, tt.role)
				c.Next()
			}, PermissionMiddleware(tt.permission))
			r.GET("/resource", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/resource", nil))

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d body=%s", w.Code, tt.wantStatus, w.Body.String())
			}

			if tt.wantStatus == http.StatusForbidden {
				var resp struct {
					Message string `json:"message"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp.Message != "You do not have permission to perform this action" {
					t.Errorf("unexpected message %q", resp.Message)
				}
			}
		})
	}
}
