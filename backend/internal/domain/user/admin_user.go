package user

import (
	"errors"
	"time"
)

// Errors for admin user validation
var (
	ErrUsernameRequired  = errors.New("username is required")
	ErrPasswordRequired  = errors.New("password is required")
	ErrAdminNameRequired = errors.New("admin name is required")
	ErrAdminEmailRequired = errors.New("admin email is required")
	ErrInvalidAdminRole  = errors.New("invalid admin role")
)

// AdminUser represents a platform administrator
type AdminUser struct {
	ID           string     `json:"id" db:"ID"`
	Username     string     `json:"username" db:"Username"`
	PasswordHash string     `json:"-" db:"PasswordHash"` // Never expose password hash in JSON
	Name         string     `json:"name" db:"Name"`
	Email        string     `json:"email" db:"Email"`
	Role         AdminRole  `json:"role" db:"Role"`
	CreatedAt    time.Time  `json:"createdAt" db:"CreatedAt"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"UpdatedAt"`
}

// IsSuperAdmin returns true if the admin has super admin role
func (u *AdminUser) IsSuperAdmin() bool {
	return u.Role == AdminRoleSuperAdmin
}

// SetRole sets the admin role with validation
func (u *AdminUser) SetRole(role AdminRole) bool {
	if role.IsValid() {
		u.Role = role
		return true
	}
	return false
}

// Validate validates the admin user fields (without password)
func (u *AdminUser) Validate() error {
	if u.Username == "" {
		return ErrUsernameRequired
	}
	if u.Name == "" {
		return ErrAdminNameRequired
	}
	if u.Email == "" {
		return ErrAdminEmailRequired
	}
	if u.Role != "" && !u.Role.IsValid() {
		return ErrInvalidAdminRole
	}
	return nil
}
