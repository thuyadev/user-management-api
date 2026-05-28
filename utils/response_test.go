package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestResponseHelpers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		handler    func(*gin.Context)
		wantStatus int
		wantBody   map[string]interface{}
	}{
		{
			name: "success",
			handler: func(c *gin.Context) {
				Success(c, http.StatusOK, "OK", gin.H{"id": 1})
			},
			wantStatus: http.StatusOK,
			wantBody: map[string]interface{}{
				"success": true,
				"message": "OK",
			},
		},
		{
			name: "unauthorized",
			handler: func(c *gin.Context) {
				Unauthorized(c, "Unauthenticated")
			},
			wantStatus: http.StatusUnauthorized,
			wantBody: map[string]interface{}{
				"success": false,
				"message": "Unauthenticated",
			},
		},
		{
			name: "forbidden",
			handler: func(c *gin.Context) {
				Forbidden(c, "Forbidden action")
			},
			wantStatus: http.StatusForbidden,
			wantBody: map[string]interface{}{
				"success": false,
				"message": "Forbidden action",
			},
		},
		{
			name: "not found",
			handler: func(c *gin.Context) {
				NotFound(c, "Resource not found")
			},
			wantStatus: http.StatusNotFound,
			wantBody: map[string]interface{}{
				"success": false,
				"message": "Resource not found",
			},
		},
		{
			name: "validation error",
			handler: func(c *gin.Context) {
				ValidationError(c, "email is required")
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantBody: map[string]interface{}{
				"success": false,
				"message": "Validation failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			tt.handler(c)

			if w.Code != tt.wantStatus {
				t.Errorf("status: got %d want %d", w.Code, tt.wantStatus)
			}

			var body map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode body: %v", err)
			}

			for key, want := range tt.wantBody {
				if body[key] != want {
					t.Errorf("%s: got %v want %v", key, body[key], want)
				}
			}
		})
	}
}

func TestPaginatedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	Paginated(c, []string{"a", "b"}, 11, 2, 5)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Success bool `json:"success"`
		Data    []string `json:"data"`
		Meta    struct {
			Total      int64 `json:"total"`
			Page       int   `json:"page"`
			PerPage    int   `json:"per_page"`
			TotalPages int   `json:"total_pages"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}

	if !resp.Success {
		t.Error("expected success true")
	}
	if len(resp.Data) != 2 {
		t.Errorf("data: got %d items want 2", len(resp.Data))
	}
	if resp.Meta.Total != 11 {
		t.Errorf("total: got %d want 11", resp.Meta.Total)
	}
	if resp.Meta.Page != 2 {
		t.Errorf("page: got %d want 2", resp.Meta.Page)
	}
	if resp.Meta.PerPage != 5 {
		t.Errorf("per_page: got %d want 5", resp.Meta.PerPage)
	}
	if resp.Meta.TotalPages != 3 {
		t.Errorf("total_pages: got %d want 3", resp.Meta.TotalPages)
	}
}
