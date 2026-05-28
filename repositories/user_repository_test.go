package repositories

import (
	"testing"
	"time"

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

	listByEmail, totalByEmail, err := repo.List(1, 10, "alice@test.com")
	if err != nil {
		t.Fatalf("list by email failed: %v", err)
	}
	if totalByEmail != 0 || len(listByEmail) != 0 {
		t.Errorf("expected name-only search, email should not match: total=%d", totalByEmail)
	}
}

func TestUserRepositoryListOrdersByCreatedAtDesc(t *testing.T) {
	repo := newTestUserRepo(t)

	first := &models.User{Name: "First", Email: "first@test.com", Password: "hash", Role: models.RoleUser}
	second := &models.User{Name: "Second", Email: "second@test.com", Password: "hash", Role: models.RoleUser}
	if err := repo.Create(first); err != nil {
		t.Fatalf("create first: %v", err)
	}
	time.Sleep(5 * time.Millisecond)
	if err := repo.Create(second); err != nil {
		t.Fatalf("create second: %v", err)
	}

	list, total, err := repo.List(1, 10, "")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if total != 2 || len(list) != 2 {
		t.Fatalf("expected 2 users, got total=%d len=%d", total, len(list))
	}
	if list[0].Name != "Second" {
		t.Errorf("expected newest first, got %s", list[0].Name)
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
