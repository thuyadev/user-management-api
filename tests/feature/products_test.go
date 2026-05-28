package feature_test

import (
	"fmt"
	"net/http"
	"testing"

	"user-management-api/tests/testutil"
)

func TestProductsViewAndManage(t *testing.T) {
	app := testutil.SetupTestApp(t)
	adminToken := app.LoginToken(t, app.Admin.Email, testutil.TestPassword)
	userToken := app.LoginToken(t, app.User.Email, testutil.TestPassword)

	categoryResp := app.Request(http.MethodPost, "/api/v1/admin/categories", map[string]string{
		"name":        "Electronics",
		"description": "Electronic devices",
	}, adminToken)
	if categoryResp.Code != http.StatusCreated {
		t.Fatalf("create category: expected 201, got %d", categoryResp.Code)
	}

	var category struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	testutil.DecodeJSON(t, categoryResp, &category)

	createResp := app.Request(http.MethodPost, "/api/v1/admin/products", map[string]interface{}{
		"name":        "Wireless Mouse",
		"description": "Ergonomic wireless mouse",
		"price":       29.99,
		"stock":       50,
		"category_id": category.Data.ID,
	}, adminToken)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create product: expected 201, got %d body=%s", createResp.Code, createResp.Body.String())
	}

	var product struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	testutil.DecodeJSON(t, createResp, &product)

	listResp := app.Request(http.MethodGet, "/api/v1/admin/products?page=1&per_page=10", nil, userToken)
	if listResp.Code != http.StatusOK {
		t.Fatalf("user list: expected 200, got %d", listResp.Code)
	}

	getResp := app.Request(http.MethodGet, fmt.Sprintf("/api/v1/admin/products/%d", product.Data.ID), nil, userToken)
	if getResp.Code != http.StatusOK {
		t.Fatalf("user get: expected 200, got %d", getResp.Code)
	}

	forbiddenResp := app.Request(http.MethodPost, "/api/v1/admin/products", map[string]interface{}{
		"name":        "Blocked Product",
		"price":       9.99,
		"stock":       1,
		"category_id": category.Data.ID,
	}, userToken)
	if forbiddenResp.Code != http.StatusForbidden {
		t.Fatalf("user create: expected 403, got %d", forbiddenResp.Code)
	}

	updateResp := app.Request(http.MethodPut, fmt.Sprintf("/api/v1/admin/products/%d", product.Data.ID), map[string]interface{}{
		"name":  "Premium Wireless Mouse",
		"price": 39.99,
	}, adminToken)
	if updateResp.Code != http.StatusOK {
		t.Fatalf("update: expected 200, got %d", updateResp.Code)
	}

	deleteResp := app.Request(http.MethodDelete, fmt.Sprintf("/api/v1/admin/products/%d", product.Data.ID), nil, adminToken)
	if deleteResp.Code != http.StatusOK {
		t.Fatalf("delete: expected 200, got %d", deleteResp.Code)
	}
}

func TestProductAIDescription(t *testing.T) {
	app := testutil.SetupTestApp(t)
	adminToken := app.LoginToken(t, app.Admin.Email, testutil.TestPassword)

	w := app.Request(http.MethodPost, "/api/v1/admin/products/ai/description", map[string]string{
		"name":     "Wireless Mouse",
		"category": "Electronics",
	}, adminToken)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}

	var resp struct {
		Data struct {
			Description string `json:"description"`
		} `json:"data"`
	}
	testutil.DecodeJSON(t, w, &resp)
	if resp.Data.Description == "" {
		t.Error("expected generated description")
	}
}
