package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gowa/saas-openclaw/backend/internal/domain/user"
)

// TestAdminUserRepository_Create tests admin creation
// Note: This is a unit test that requires a database connection.
// For integration tests with testcontainers, see admin_user_repository_integration_test.go
func TestAdminUserRepository_Create_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	// Integration test placeholder - requires database
}

// TestAdminUserRepository_GetByID tests getting admin by ID
func TestAdminUserRepository_GetByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	// Integration test placeholder - requires database
}

// TestAdminUserRepository_GetByUsername tests getting admin by username
func TestAdminUserRepository_GetByUsername_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	// Integration test placeholder - requires database
}

// TestAdminUserRepository_GetByEmail tests getting admin by email
func TestAdminUserRepository_GetByEmail_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	// Integration test placeholder - requires database
}

// TestAdminUserRepository_Update tests updating admin
func TestAdminUserRepository_Update_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	// Integration test placeholder - requires database
}

// TestAdminUserRepository_Delete tests deleting admin
func TestAdminUserRepository_Delete_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	// Integration test placeholder - requires database
}

// Unit tests for domain validation within repository context

func TestAdminUser_Validation(t *testing.T) {
	tests := []struct {
		name  string
		user  *user.AdminUser
		valid bool
	}{
		{
			name: "valid admin",
			user: &user.AdminUser{
				Username: "admin",
				Name:     "Admin User",
				Email:    "admin@example.com",
				Role:     user.AdminRoleAdmin,
			},
			valid: true,
		},
		{
			name: "super admin",
			user: &user.AdminUser{
				Username: "superadmin",
				Name:     "Super Admin",
				Email:    "superadmin@example.com",
				Role:     user.AdminRoleSuperAdmin,
			},
			valid: true,
		},
		{
			name: "missing username",
			user: &user.AdminUser{
				Name:  "Admin User",
				Email: "admin@example.com",
			},
			valid: false,
		},
		{
			name: "missing name",
			user: &user.AdminUser{
				Username: "admin",
				Email:    "admin@example.com",
			},
			valid: false,
		},
		{
			name: "missing email",
			user: &user.AdminUser{
				Username: "admin",
				Name:     "Admin User",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestAdminUser_RoleDefaults(t *testing.T) {
	// Test that the default role is properly set
	u := &user.AdminUser{
		Username: "admin",
		Name:     "Admin User",
		Email:    "admin@example.com",
	}

	// Role should be empty string (zero value) by default
	require.Equal(t, user.AdminRole(""), u.Role)

	// After setting a valid role
	ok := u.SetRole(user.AdminRoleAdmin)
	require.True(t, ok)
	assert.Equal(t, user.AdminRoleAdmin, u.Role)
}

func TestAdminUserRepository_VerifyPassword(t *testing.T) {
	// Note: This test verifies the password verification logic without database
	// The actual bcrypt hashing is tested in integration tests

	t.Run("nil user returns false", func(t *testing.T) {
		repo := &AdminUserRepository{}
		assert.False(t, repo.VerifyPassword(nil, "password"))
	})

	t.Run("empty password hash returns false", func(t *testing.T) {
		repo := &AdminUserRepository{}
		u := &user.AdminUser{PasswordHash: ""}
		assert.False(t, repo.VerifyPassword(u, "password"))
	})

	t.Run("empty password returns false", func(t *testing.T) {
		repo := &AdminUserRepository{}
		u := &user.AdminUser{PasswordHash: "somehash"}
		assert.False(t, repo.VerifyPassword(u, ""))
	})
}
