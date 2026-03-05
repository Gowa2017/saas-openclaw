package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gowa/saas-openclaw/backend/internal/domain/user"
)

// TestTenantUserRepository_Create tests user creation
// Note: This is a unit test that requires a database connection.
// For integration tests with testcontainers, see tenant_user_repository_integration_test.go
func TestTenantUserRepository_Create_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	// Integration test placeholder - requires database
}

// TestTenantUserRepository_GetByID tests getting user by ID
func TestTenantUserRepository_GetByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	// Integration test placeholder - requires database
}

// TestTenantUserRepository_GetByEmail tests getting user by email
func TestTenantUserRepository_GetByEmail_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	// Integration test placeholder - requires database
}

// TestTenantUserRepository_Update tests updating user
func TestTenantUserRepository_Update_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	// Integration test placeholder - requires database
}

// TestTenantUserRepository_Delete tests deleting user
func TestTenantUserRepository_Delete_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	// Integration test placeholder - requires database
}

// Unit tests for domain validation within repository context

func TestTenantUser_Validation(t *testing.T) {
	tests := []struct {
		name  string
		user  *user.TenantUser
		valid bool
	}{
		{
			name: "valid user",
			user: &user.TenantUser{
				TenantID: "tenant-123",
				Name:     "Test User",
				Email:    "test@example.com",
				Role:     user.RoleTenantUser,
			},
			valid: true,
		},
		{
			name: "admin user",
			user: &user.TenantUser{
				TenantID: "tenant-123",
				Name:     "Admin User",
				Email:    "admin@example.com",
				Role:     user.RoleTenantAdmin,
			},
			valid: true,
		},
		{
			name: "missing tenant ID",
			user: &user.TenantUser{
				Name:  "Test User",
				Email: "test@example.com",
			},
			valid: false,
		},
		{
			name: "missing name",
			user: &user.TenantUser{
				TenantID: "tenant-123",
				Email:    "test@example.com",
			},
			valid: false,
		},
		{
			name: "missing email",
			user: &user.TenantUser{
				TenantID: "tenant-123",
				Name:     "Test User",
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

func TestTenantUser_RoleDefaults(t *testing.T) {
	// Test that the default role is properly set
	u := &user.TenantUser{
		TenantID: "tenant-123",
		Name:     "Test User",
		Email:    "test@example.com",
	}

	// Role should be empty string (zero value) by default
	require.Equal(t, user.Role(""), u.Role)

	// After setting a valid role
	ok := u.SetRole(user.RoleTenantUser)
	require.True(t, ok)
	assert.Equal(t, user.RoleTenantUser, u.Role)
}

func TestRepositoryErrors(t *testing.T) {
	// Test error types
	t.Run("ErrNotFound", func(t *testing.T) {
		assert.True(t, IsNotFoundError(ErrNotFound))
		assert.False(t, IsNotFoundError(ErrValidation))
	})

	t.Run("WrapNotFoundError", func(t *testing.T) {
		err := WrapNotFoundError(nil, "tenant")
		assert.Error(t, err)
		assert.True(t, IsNotFoundError(err))
	})
}
