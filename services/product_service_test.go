package services

import (
	"testing"

	"user-management-api/database"
	"user-management-api/models"
	"user-management-api/repositories"
)

func setupProductService(t *testing.T) (ProductService, repositories.UserRepository, repositories.CategoryRepository) {
	t.Helper()
	db, err := database.ConnectSQLite()
	if err != nil {
		t.Fatalf("failed to connect sqlite: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	userRepo := repositories.NewUserRepository(db)
	logSvc := newMockLogService()
	return NewProductService(productRepo, categoryRepo, logSvc), userRepo, categoryRepo
}

func TestProductServiceCRUD(t *testing.T) {
	productSvc, userRepo, categoryRepo := setupProductService(t)

	admin := &models.User{Name: "Admin", Email: "admin@test.com", Password: "hash", Role: models.RoleAdmin}
	if err := userRepo.Create(admin); err != nil {
		t.Fatalf("create admin failed: %v", err)
	}

	category := &models.Category{Name: "Books", Description: "Book category"}
	if err := categoryRepo.Create(category); err != nil {
		t.Fatalf("create category failed: %v", err)
	}

	created, err := productSvc.Create(models.CreateProductRequest{
		Name:        "Go Book",
		Description: "Learn Go",
		Price:       29.99,
		Stock:       10,
		CategoryID:  category.ID,
	}, admin.ID)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	got, err := productSvc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.Name != "Go Book" {
		t.Errorf("expected Go Book, got %s", got.Name)
	}

	newPrice := 39.99
	updated, err := productSvc.Update(created.ID, models.UpdateProductRequest{
		Name:  "Advanced Go",
		Price: &newPrice,
	}, admin.ID)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Name != "Advanced Go" {
		t.Errorf("expected Advanced Go, got %s", updated.Name)
	}

	products, total, err := productSvc.List(1, 10, "", category.ID)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 product, got %d", total)
	}
	if len(products) != 1 {
		t.Errorf("expected 1 product in page, got %d", len(products))
	}

	if err := productSvc.Delete(created.ID, admin.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err = productSvc.GetByID(created.ID)
	if err != ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestProductServiceCreateInvalidCategory(t *testing.T) {
	productSvc, userRepo, _ := setupProductService(t)

	admin := &models.User{Name: "Admin", Email: "admin@test.com", Password: "hash", Role: models.RoleAdmin}
	if err := userRepo.Create(admin); err != nil {
		t.Fatalf("create admin failed: %v", err)
	}

	_, err := productSvc.Create(models.CreateProductRequest{
		Name:       "Item",
		Price:      10,
		Stock:      1,
		CategoryID: 999,
	}, admin.ID)
	if err != ErrCategoryNotFound {
		t.Errorf("expected ErrCategoryNotFound, got %v", err)
	}
}
