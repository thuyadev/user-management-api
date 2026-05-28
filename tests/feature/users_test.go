package feature_test

import (
	"fmt"
	"net/http"
	"testing"

	"user-management-api/tests/testutil"
)

func TestUsersAdminCRUD(t *testing.T) {
	app := testutil.SetupTestApp(t)
	adminToken := app.LoginToken(t, app.Admin.Email, testutil.TestPassword)

	createResp := app.Request(http.MethodPost, "/api/v1/admin/users", map[string]string{
		"name":     "Created User",
		"email":    "created@test.com",
		"password": testutil.TestPassword,
		"role":     "user",
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

	listResp := app.Request(http.MethodGet, "/api/v1/admin/users?page=1&per_page=10", nil, adminToken)
	if listResp.Code != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", listResp.Code)
	}

	getResp := app.Request(http.MethodGet, fmt.Sprintf("/api/v1/admin/users/%d", created.Data.ID), nil, adminToken)
	if getResp.Code != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", getResp.Code)
	}

	updateResp := app.Request(http.MethodPut, fmt.Sprintf("/api/v1/admin/users/%d", created.Data.ID), map[string]string{
		"name": "Updated User",
	}, adminToken)
	if updateResp.Code != http.StatusOK {
		t.Fatalf("update: expected 200, got %d", updateResp.Code)
	}

	deleteResp := app.Request(http.MethodDelete, fmt.Sprintf("/api/v1/admin/users/%d", created.Data.ID), nil, adminToken)
	if deleteResp.Code != http.StatusOK {
		t.Fatalf("delete: expected 200, got %d", deleteResp.Code)
	}
}

func TestUsersForbiddenForRegularUser(t *testing.T) {
	app := testutil.SetupTestApp(t)
	userToken := app.LoginToken(t, app.User.Email, testutil.TestPassword)

	w := app.Request(http.MethodGet, "/api/v1/admin/users?page=1&per_page=10", nil, userToken)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d body=%s", w.Code, w.Body.String())
	}
}
