package services

import (
	"testing"

	"user-management-api/database"
	"user-management-api/models"
	"user-management-api/repositories"
	"user-management-api/utils"
)

func setupUserService(t *testing.T) (UserService, repositories.UserRepository) {
	t.Helper()
	db, err := database.ConnectSQLite()
	if err != nil {
		t.Fatalf("failed to connect sqlite: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	logSvc := newMockLogService()
	return NewUserService(userRepo, logSvc), userRepo
}

func TestUserServiceCRUD(t *testing.T) {
	userSvc, userRepo := setupUserService(t)

	hash, _ := utils.HashPassword("password123")
	admin := &models.User{Name: "Admin", Email: "admin@test.com", Password: hash, Role: models.RoleAdmin}
	if err := userRepo.Create(admin); err != nil {
		t.Fatalf("create admin failed: %v", err)
	}

	created, err := userSvc.Create(models.CreateUserRequest{
		Name:     "Test User",
		Email:    "user@test.com",
		Password: "password123",
		Role:     models.RoleUser,
	}, admin.ID)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	got, err := userSvc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.Name != "Test User" {
		t.Errorf("expected Test User, got %s", got.Name)
	}

	updated, err := userSvc.Update(created.ID, models.UpdateUserRequest{Name: "Updated User"}, admin.ID)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Name != "Updated User" {
		t.Errorf("expected Updated User, got %s", updated.Name)
	}

	users, total, err := userSvc.List(1, 10, "")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if total < 2 {
		t.Errorf("expected at least 2 users, got %d", total)
	}
	if len(users) < 2 {
		t.Errorf("expected at least 2 users in page, got %d", len(users))
	}

	if err := userSvc.Delete(created.ID, admin.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err = userSvc.GetByID(created.ID)
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}
