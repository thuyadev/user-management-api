package services

import (
	"testing"

	"user-management-api/policies"
	"user-management-api/database"
	"user-management-api/models"
	"user-management-api/repositories"
	"user-management-api/utils"
)

func newTestUserRepo(t *testing.T) repositories.UserRepository {
	t.Helper()
	db, err := database.ConnectSQLite()
	if err != nil {
		t.Fatalf("failed to connect sqlite: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migration failed: %v", err)
	}
	return repositories.NewUserRepository(db)
}

func TestAuthServiceLogin(t *testing.T) {
	userRepo := newTestUserRepo(t)
	logSvc := newMockLogService()
	authSvc := NewAuthService(userRepo, logSvc, "secret", 24)

	hash, err := utils.HashPassword("password123")
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}

	user := &models.User{
		Name:     "Admin",
		Email:    "admin@test.com",
		Password: hash,
		Role:     models.RoleAdmin,
	}
	if err := userRepo.Create(user); err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	resp, err := authSvc.Login(models.LoginRequest{
		Email:    "admin@test.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if resp.Token == "" {
		t.Error("expected token")
	}
	if resp.User.Email != "admin@test.com" {
		t.Errorf("expected admin@test.com, got %s", resp.User.Email)
	}
	if resp.Role != models.RoleAdmin {
		t.Errorf("expected role admin, got %s", resp.Role)
	}
	if resp.ExpiresAt.IsZero() {
		t.Error("expected expires_at")
	}
	if resp.ExpiresIn <= 0 {
		t.Error("expected positive expires_in")
	}
	if len(resp.Permissions) != 6 {
		t.Errorf("expected 6 permissions, got %d", len(resp.Permissions))
	}
	if !policies.HasPermission(resp.Role, policies.PermUsersManage) {
		t.Error("admin login should include users.manage")
	}

	_, err = authSvc.Login(models.LoginRequest{
		Email:    "admin@test.com",
		Password: "wrong",
	})
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthServiceRegister(t *testing.T) {
	userRepo := newTestUserRepo(t)
	logSvc := newMockLogService()
	authSvc := NewAuthService(userRepo, logSvc, "secret", 24)

	resp, err := authSvc.Register(models.RegisterRequest{
		Name:     "New User",
		Email:    "new@test.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if resp.Token == "" {
		t.Error("expected token")
	}
	if resp.User.Email != "new@test.com" {
		t.Errorf("expected new@test.com, got %s", resp.User.Email)
	}
	if resp.Role != models.RoleUser {
		t.Errorf("expected role user, got %s", resp.Role)
	}
	if len(resp.Permissions) != 2 {
		t.Errorf("expected 2 permissions, got %d", len(resp.Permissions))
	}

	_, err = authSvc.Register(models.RegisterRequest{
		Name:     "Duplicate",
		Email:    "new@test.com",
		Password: "password123",
	})
	if err != ErrEmailTaken {
		t.Errorf("expected ErrEmailTaken, got %v", err)
	}
}
