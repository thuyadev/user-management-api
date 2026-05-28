package repositories

import (
	"testing"

	"user-management-api/database"
	"user-management-api/models"

	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := database.ConnectSQLite()
	if err != nil {
		t.Fatalf("connect sqlite: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func setupCategoryRepo(t *testing.T) CategoryRepository {
	t.Helper()
	return NewCategoryRepository(newTestDB(t))
}

func TestCategoryRepositoryCRUD(t *testing.T) {
	repo := setupCategoryRepo(t)

	category := &models.Category{Name: "Books", Description: "All books"}
	if err := repo.Create(category); err != nil {
		t.Fatalf("create: %v", err)
	}

	found, err := repo.FindByID(category.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if found.Name != "Books" {
		t.Errorf("expected Books, got %s", found.Name)
	}

	found.Name = "Updated Books"
	if err := repo.Update(found); err != nil {
		t.Fatalf("update: %v", err)
	}

	list, total, err := repo.List(1, 10, "updated")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Errorf("expected 1 result, got total=%d len=%d", total, len(list))
	}

	if err := repo.Delete(category.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err = repo.FindByID(category.ID)
	if err == nil {
		t.Error("expected error after delete")
	}
}
