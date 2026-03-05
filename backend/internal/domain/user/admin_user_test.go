package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminUser_IsSuperAdmin(t *testing.T) {
	tests := []struct {
		name string
		role AdminRole
		want bool
	}{
		{
			name: "super admin",
			role: AdminRoleSuperAdmin,
			want: true,
		},
		{
			name: "regular admin",
			role: AdminRoleAdmin,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &AdminUser{
				ID:       "test-id",
				Username: "admin",
				Name:     "Admin User",
				Email:    "admin@example.com",
				Role:     tt.role,
			}
			assert.Equal(t, tt.want, user.IsSuperAdmin())
		})
	}
}

func TestAdminUser_SetRole(t *testing.T) {
	tests := []struct {
		name    string
		role    AdminRole
		want    bool
		wantErr bool
	}{
		{
			name: "set super admin role",
			role: AdminRoleSuperAdmin,
			want: true,
		},
		{
			name: "set admin role",
			role: AdminRoleAdmin,
			want: true,
		},
		{
			name:    "set invalid role",
			role:    AdminRole("invalid"),
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &AdminUser{}
			result := user.SetRole(tt.role)
			assert.Equal(t, tt.want, result)
			if !tt.wantErr {
				assert.Equal(t, tt.role, user.Role)
			}
		})
	}
}

func TestAdminUser_StructFields(t *testing.T) {
	now := time.Now()
	user := &AdminUser{
		ID:           "test-id",
		Username:     "admin",
		PasswordHash: "hashed-password",
		Name:         "Admin User",
		Email:        "admin@example.com",
		Role:         AdminRoleAdmin,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	assert.Equal(t, "test-id", user.ID)
	assert.Equal(t, "admin", user.Username)
	assert.Equal(t, "hashed-password", user.PasswordHash)
	assert.Equal(t, "Admin User", user.Name)
	assert.Equal(t, "admin@example.com", user.Email)
	assert.Equal(t, AdminRoleAdmin, user.Role)
	assert.Equal(t, now, user.CreatedAt)
	assert.Equal(t, now, user.UpdatedAt)
}

func TestAdminRole_IsValid(t *testing.T) {
	tests := []struct {
		name string
		role AdminRole
		want bool
	}{
		{
			name: "super admin role",
			role: AdminRoleSuperAdmin,
			want: true,
		},
		{
			name: "admin role",
			role: AdminRoleAdmin,
			want: true,
		},
		{
			name: "invalid role",
			role: AdminRole("invalid"),
			want: false,
		},
		{
			name: "empty role",
			role: AdminRole(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.role.IsValid())
		})
	}
}

func TestAdminUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    *AdminUser
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid admin",
			user: &AdminUser{
				Username: "admin",
				Name:     "Admin User",
				Email:    "admin@example.com",
				Role:     AdminRoleAdmin,
			},
			wantErr: false,
		},
		{
			name: "valid admin with empty role",
			user: &AdminUser{
				Username: "admin",
				Name:     "Admin User",
				Email:    "admin@example.com",
			},
			wantErr: false,
		},
		{
			name: "missing username",
			user: &AdminUser{
				Name:  "Admin User",
				Email: "admin@example.com",
			},
			wantErr: true,
			errMsg:  "username is required",
		},
		{
			name: "missing name",
			user: &AdminUser{
				Username: "admin",
				Email:    "admin@example.com",
			},
			wantErr: true,
			errMsg:  "admin name is required",
		},
		{
			name: "missing email",
			user: &AdminUser{
				Username: "admin",
				Name:     "Admin User",
			},
			wantErr: true,
			errMsg:  "admin email is required",
		},
		{
			name: "invalid role",
			user: &AdminUser{
				Username: "admin",
				Name:     "Admin User",
				Email:    "admin@example.com",
				Role:     AdminRole("invalid"),
			},
			wantErr: true,
			errMsg:  "invalid admin role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
