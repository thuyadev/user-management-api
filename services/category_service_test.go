package services

import (
	"testing"

	"user-management-api/database"
	"user-management-api/models"
	"user-management-api/repositories"
)

func setupCategoryService(t *testing.T) (CategoryService, repositories.UserRepository) {
	t.Helper()
	db, err := database.ConnectSQLite()
	if err != nil {
		t.Fatalf("failed to connect sqlite: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	categoryRepo := repositories.NewCategoryRepository(db)
	userRepo := repositories.NewUserRepository(db)
	logSvc := newMockLogService()
	return NewCategoryService(categoryRepo, logSvc), userRepo
}

func TestCategoryServiceCRUD(t *testing.T) {
	categorySvc, userRepo := setupCategoryService(t)

	admin := &models.User{Name: "Admin", Email: "admin@test.com", Password: "hash", Role: models.RoleAdmin}
	if err := userRepo.Create(admin); err != nil {
		t.Fatalf("create admin failed: %v", err)
	}

	created, err := categorySvc.Create(models.CreateCategoryRequest{
		Name:        "Electronics",
		Description: "Electronic items",
	}, admin.ID)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	got, err := categorySvc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.Name != "Electronics" {
		t.Errorf("expected Electronics, got %s", got.Name)
	}

	updated, err := categorySvc.Update(created.ID, models.UpdateCategoryRequest{Name: "Gadgets"}, admin.ID)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Name != "Gadgets" {
		t.Errorf("expected Gadgets, got %s", updated.Name)
	}

	categories, total, err := categorySvc.List(1, 10, "")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 category, got %d", total)
	}
	if len(categories) != 1 {
		t.Errorf("expected 1 category in page, got %d", len(categories))
	}

	if err := categorySvc.Delete(created.ID, admin.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err = categorySvc.GetByID(created.ID)
	if err != ErrCategoryNotFound {
		t.Errorf("expected ErrCategoryNotFound, got %v", err)
	}
}
