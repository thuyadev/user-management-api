package feature_test

import (
	"net/http"
	"testing"

	"user-management-api/models"
	"user-management-api/tests/testutil"
)

func TestAuthLoginSuccess(t *testing.T) {
	app := testutil.SetupTestApp(t)

	w := app.Request(http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    app.Admin.Email,
		"password": testutil.TestPassword,
	}, "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			Token       string   `json:"token"`
			Role        string   `json:"role"`
			Permissions []string `json:"permissions"`
			User        struct {
				Email string `json:"email"`
			} `json:"user"`
		} `json:"data"`
	}
	testutil.DecodeJSON(t, w, &resp)

	if !resp.Success || resp.Data.Token == "" {
		t.Error("expected success with token")
	}
	if resp.Data.Role != models.RoleAdmin {
		t.Errorf("expected admin role, got %s", resp.Data.Role)
	}
	if resp.Data.User.Email != app.Admin.Email {
		t.Errorf("expected %s, got %s", app.Admin.Email, resp.Data.User.Email)
	}
}

func TestAuthLoginInvalidCredentials(t *testing.T) {
	app := testutil.SetupTestApp(t)

	w := app.Request(http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    app.Admin.Email,
		"password": "wrong-password",
	}, "")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
	if msg := testutil.APIResponseMessage(t, w); msg != "Invalid email or password" {
		t.Errorf("expected invalid credentials message, got %q", msg)
	}
}

func TestAuthRegisterSuccess(t *testing.T) {
	app := testutil.SetupTestApp(t)

	w := app.Request(http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "New User",
		"email":    "newuser@test.com",
		"password": testutil.TestPassword,
	}, "")

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}

	var resp struct {
		Data struct {
			Role string `json:"role"`
			User struct {
				Email string `json:"email"`
			} `json:"user"`
		} `json:"data"`
	}
	testutil.DecodeJSON(t, w, &resp)

	if resp.Data.Role != models.RoleUser {
		t.Errorf("expected user role, got %s", resp.Data.Role)
	}
	if resp.Data.User.Email != "newuser@test.com" {
		t.Errorf("unexpected email %s", resp.Data.User.Email)
	}
}

func TestAuthMeRequiresValidToken(t *testing.T) {
	app := testutil.SetupTestApp(t)

	t.Run("missing token", func(t *testing.T) {
		w := app.Request(http.MethodGet, "/api/v1/auth/me", nil, "")
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
		if msg := testutil.APIResponseMessage(t, w); msg != "Unauthenticated" {
			t.Errorf("expected Unauthenticated, got %q", msg)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		w := app.Request(http.MethodGet, "/api/v1/auth/me", nil, "invalid.token.here")
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
		if msg := testutil.APIResponseMessage(t, w); msg != "Unauthenticated" {
			t.Errorf("expected Unauthenticated, got %q", msg)
		}
	})

	t.Run("valid token", func(t *testing.T) {
		token := app.LoginToken(t, app.Admin.Email, testutil.TestPassword)
		w := app.Request(http.MethodGet, "/api/v1/auth/me", nil, token)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
		}

		var resp struct {
			Data struct {
				Email       string   `json:"email"`
				Role        string   `json:"role"`
				Permissions []string `json:"permissions"`
			} `json:"data"`
		}
		testutil.DecodeJSON(t, w, &resp)

		if resp.Data.Email != app.Admin.Email {
			t.Errorf("expected %s, got %s", app.Admin.Email, resp.Data.Email)
		}
		if resp.Data.Role != models.RoleAdmin {
			t.Errorf("expected admin, got %s", resp.Data.Role)
		}
	})
}

func TestAuthRolesEndpoint(t *testing.T) {
	app := testutil.SetupTestApp(t)

	w := app.Request(http.MethodGet, "/api/v1/auth/roles", nil, "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Data map[string][]string `json:"data"`
	}
	testutil.DecodeJSON(t, w, &resp)

	if len(resp.Data[models.RoleAdmin]) == 0 {
		t.Error("expected admin permissions")
	}
	if len(resp.Data[models.RoleUser]) == 0 {
		t.Error("expected user permissions")
	}
}
