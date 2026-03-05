package user

import (
	"errors"
	"time"
)

// Errors for tenant user validation
var (
	ErrTenantIDRequired = errors.New("tenant ID is required")
	ErrNameRequired     = errors.New("name is required")
	ErrEmailRequired    = errors.New("email is required")
	ErrInvalidRole      = errors.New("invalid user role")
)

// TenantUser represents a user belonging to a tenant in the multi-tenant system
type TenantUser struct {
	ID        string    `json:"id" db:"ID"`
	TenantID  string    `json:"tenantId" db:"TenantID"`
	Name      string    `json:"name" db:"Name"`
	Email     string    `json:"email" db:"Email"`
	Role      Role      `json:"role" db:"Role"`
	CreatedAt time.Time `json:"createdAt" db:"CreatedAt"`
	UpdatedAt time.Time `json:"updatedAt" db:"UpdatedAt"`
}

// IsAdmin returns true if the user has admin role
func (u *TenantUser) IsAdmin() bool {
	return u.Role == RoleTenantAdmin
}

// SetRole sets the user role with validation
func (u *TenantUser) SetRole(role Role) bool {
	if role.IsValid() {
		u.Role = role
		return true
	}
	return false
}

// Validate validates the tenant user fields
func (u *TenantUser) Validate() error {
	if u.TenantID == "" {
		return ErrTenantIDRequired
	}
	if u.Name == "" {
		return ErrNameRequired
	}
	if u.Email == "" {
		return ErrEmailRequired
	}
	if u.Role != "" && !u.Role.IsValid() {
		return ErrInvalidRole
	}
	return nil
}
