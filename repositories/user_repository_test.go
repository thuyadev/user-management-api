package repositories

import (
	"testing"

	"user-management-api/database"
	"user-management-api/models"
)

func newTestUserRepo(t *testing.T) UserRepository {
	t.Helper()
	db, err := database.ConnectSQLite()
	if err != nil {
		t.Fatalf("failed to connect sqlite: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migration failed: %v", err)
	}
	return NewUserRepository(db)
}

func TestUserRepositoryListAndSearch(t *testing.T) {
	repo := newTestUserRepo(t)

	users := []models.User{
		{Name: "Alice", Email: "alice@test.com", Password: "hash", Role: models.RoleUser},
		{Name: "Bob", Email: "bob@test.com", Password: "hash", Role: models.RoleUser},
	}
	for i := range users {
		if err := repo.Create(&users[i]); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	list, total, err := repo.List(1, 10, "alice")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(list) != 1 || list[0].Name != "Alice" {
		t.Errorf("unexpected list result: %+v", list)
	}
}

func TestUserRepositoryFindByEmail(t *testing.T) {
	repo := newTestUserRepo(t)

	user := &models.User{Name: "Test", Email: "test@test.com", Password: "hash", Role: models.RoleUser}
	if err := repo.Create(user); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	found, err := repo.FindByEmail("test@test.com")
	if err != nil {
		t.Fatalf("find failed: %v", err)
	}
	if found.Email != "test@test.com" {
		t.Errorf("expected test@test.com, got %s", found.Email)
	}
}
