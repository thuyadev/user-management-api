package policies

import (
	"testing"

	"user-management-api/models"
)

func TestPermissionsForRoleAdmin(t *testing.T) {
	perms := PermissionsForRole(models.RoleAdmin)
	if len(perms) != 6 {
		t.Fatalf("expected 6 admin permissions, got %d", len(perms))
	}
	if !HasPermission(models.RoleAdmin, PermCategoriesManage) {
		t.Error("admin should manage categories")
	}
}

func TestPermissionsForRoleUser(t *testing.T) {
	perms := PermissionsForRole(models.RoleUser)
	if len(perms) != 2 {
		t.Fatalf("expected 2 user permissions, got %d", len(perms))
	}
	if HasPermission(models.RoleUser, PermCategoriesManage) {
		t.Error("user should not manage categories")
	}
	if !HasPermission(models.RoleUser, PermProductsView) {
		t.Error("user should view products")
	}
}

func TestHasPermissionUnknownRole(t *testing.T) {
	if HasPermission("guest", PermProductsView) {
		t.Error("unknown role should have no permissions")
	}
}
