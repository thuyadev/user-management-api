package repositories

import (
	"testing"

	"user-management-api/models"
)

func setupProductRepo(t *testing.T) (ProductRepository, CategoryRepository) {
	t.Helper()
	db := newTestDB(t)
	return NewProductRepository(db), NewCategoryRepository(db)
}

func TestProductRepositoryListWithCategoryFilter(t *testing.T) {
	productRepo, categoryRepo := setupProductRepo(t)

	category := &models.Category{Name: "Electronics", Description: "Devices"}
	if err := categoryRepo.Create(category); err != nil {
		t.Fatalf("create category: %v", err)
	}

	products := []models.Product{
		{Name: "Phone", Description: "Smartphone", Price: 699.99, Stock: 10, CategoryID: category.ID},
		{Name: "Laptop", Description: "Notebook", Price: 999.99, Stock: 5, CategoryID: category.ID},
	}
	for i := range products {
		if err := productRepo.Create(&products[i]); err != nil {
			t.Fatalf("create product: %v", err)
		}
	}

	list, total, err := productRepo.List(1, 10, "phone", category.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(list) != 1 || list[0].Name != "Phone" {
		t.Errorf("unexpected list: %+v", list)
	}

	found, err := productRepo.FindByID(products[0].ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if found.Category == nil || found.Category.Name != "Electronics" {
		t.Errorf("expected preloaded category, got %+v", found.Category)
	}
}
