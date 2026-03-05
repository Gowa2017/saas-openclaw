// Package user defines user-related domain entities
package user

// Role represents the role of a tenant user
type Role string

const (
	// RoleTenantUser is the default role for tenant users
	RoleTenantUser Role = "user"
	// RoleTenantAdmin is the admin role for tenant users
	RoleTenantAdmin Role = "admin"
)

// IsValid checks if the role is valid
func (r Role) IsValid() bool {
	switch r {
	case RoleTenantUser, RoleTenantAdmin:
		return true
	default:
		return false
	}
}

// AdminRole represents the role of an admin user
type AdminRole string

const (
	// AdminRoleSuperAdmin has full system access
	AdminRoleSuperAdmin AdminRole = "super_admin"
	// AdminRoleAdmin has standard admin access
	AdminRoleAdmin AdminRole = "admin"
)

// IsValid checks if the admin role is valid
func (r AdminRole) IsValid() bool {
	switch r {
	case AdminRoleSuperAdmin, AdminRoleAdmin:
		return true
	default:
		return false
	}
}
