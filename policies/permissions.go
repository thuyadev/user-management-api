package policies

import "user-management-api/models"

// Permission names follow a resource.action pattern (similar to Spatie / Laravel Policies).
const (
	PermUsersManage      = "users.manage"
	PermCategoriesView   = "categories.view"
	PermCategoriesManage = "categories.manage"
	PermProductsView     = "products.view"
	PermProductsManage   = "products.manage"
	PermLogsView         = "logs.view"
)

var rolePermissions = map[string][]string{
	models.RoleAdmin: {
		PermUsersManage,
		PermCategoriesView,
		PermCategoriesManage,
		PermProductsView,
		PermProductsManage,
		PermLogsView,
	},
	models.RoleUser: {
		PermCategoriesView,
		PermProductsView,
	},
}

// PermissionsForRole returns a copy of permissions assigned to the role.
func PermissionsForRole(role string) []string {
	perms, ok := rolePermissions[role]
	if !ok {
		return nil
	}
	out := make([]string, len(perms))
	copy(out, perms)
	return out
}

// HasPermission checks whether a role is allowed to perform an action.
func HasPermission(role, permission string) bool {
	for _, p := range rolePermissions[role] {
		if p == permission {
			return true
		}
	}
	return false
}

// AllRoles returns supported role names and their permissions (for docs/seeding).
func AllRoles() map[string][]string {
	out := make(map[string][]string, len(rolePermissions))
	for role, perms := range rolePermissions {
		copied := make([]string, len(perms))
		copy(copied, perms)
		out[role] = copied
	}
	return out
}
