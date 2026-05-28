package feature_test

import (
	"fmt"
	"net/http"
	"testing"

	"user-management-api/tests/testutil"
)

func TestCategoriesViewAndManage(t *testing.T) {
	app := testutil.SetupTestApp(t)
	adminToken := app.LoginToken(t, app.Admin.Email, testutil.TestPassword)
	userToken := app.LoginToken(t, app.User.Email, testutil.TestPassword)

	createResp := app.Request(http.MethodPost, "/api/v1/admin/categories", map[string]string{
		"name":        "Electronics",
		"description": "Electronic devices",
	}, adminToken)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d body=%s", createResp.Code, createResp.Body.String())
	}

	var created struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	testutil.DecodeJSON(t, createResp, &created)

	listResp := app.Request(http.MethodGet, "/api/v1/admin/categories?page=1&per_page=10", nil, userToken)
	if listResp.Code != http.StatusOK {
		t.Fatalf("user list: expected 200, got %d", listResp.Code)
	}

	getResp := app.Request(http.MethodGet, fmt.Sprintf("/api/v1/admin/categories/%d", created.Data.ID), nil, userToken)
	if getResp.Code != http.StatusOK {
		t.Fatalf("user get: expected 200, got %d", getResp.Code)
	}

	forbiddenResp := app.Request(http.MethodPost, "/api/v1/admin/categories", map[string]string{
		"name":        "Blocked",
		"description": "Should fail",
	}, userToken)
	if forbiddenResp.Code != http.StatusForbidden {
		t.Fatalf("user create: expected 403, got %d", forbiddenResp.Code)
	}

	updateResp := app.Request(http.MethodPut, fmt.Sprintf("/api/v1/admin/categories/%d", created.Data.ID), map[string]string{
		"name": "Gadgets",
	}, adminToken)
	if updateResp.Code != http.StatusOK {
		t.Fatalf("update: expected 200, got %d", updateResp.Code)
	}

	deleteResp := app.Request(http.MethodDelete, fmt.Sprintf("/api/v1/admin/categories/%d", created.Data.ID), nil, adminToken)
	if deleteResp.Code != http.StatusOK {
		t.Fatalf("delete: expected 200, got %d", deleteResp.Code)
	}
}

func TestCategoryAISuggest(t *testing.T) {
	app := testutil.SetupTestApp(t)
	adminToken := app.LoginToken(t, app.Admin.Email, testutil.TestPassword)

	w := app.Request(http.MethodPost, "/api/v1/admin/categories/ai/suggest", map[string]string{
		"keywords": "fitness gym equipment",
	}, adminToken)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}

	var resp struct {
		Data struct {
			Name string `json:"name"`
		} `json:"data"`
	}
	testutil.DecodeJSON(t, w, &resp)
	if resp.Data.Name == "" {
		t.Error("expected suggested category name")
	}
}
