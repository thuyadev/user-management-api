package feature_test

import (
	"net/http"
	"testing"

	"user-management-api/tests/testutil"
)

func TestHealthCheck(t *testing.T) {
	app := testutil.SetupTestApp(t)

	w := app.Request(http.MethodGet, "/health", nil, "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Success bool `json:"success"`
		Message string `json:"message"`
		Data    struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	testutil.DecodeJSON(t, w, &resp)

	if !resp.Success {
		t.Error("expected success true")
	}
	if resp.Data.Status != "healthy" {
		t.Errorf("expected healthy status, got %s", resp.Data.Status)
	}
}
