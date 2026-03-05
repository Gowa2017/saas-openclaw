package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTenantUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name string
		role Role
		want bool
	}{
		{
			name: "admin user",
			role: RoleTenantAdmin,
			want: true,
		},
		{
			name: "regular user",
			role: RoleTenantUser,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &TenantUser{
				ID:       "test-id",
				TenantID: "tenant-id",
				Name:     "Test User",
				Email:    "test@example.com",
				Role:     tt.role,
			}
			assert.Equal(t, tt.want, user.IsAdmin())
		})
	}
}

func TestTenantUser_SetRole(t *testing.T) {
	tests := []struct {
		name    string
		role    Role
		want    bool
		wantErr bool
	}{
		{
			name: "set admin role",
			role: RoleTenantAdmin,
			want: true,
		},
		{
			name: "set user role",
			role: RoleTenantUser,
			want: true,
		},
		{
			name:    "set invalid role",
			role:    Role("invalid"),
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &TenantUser{}
			result := user.SetRole(tt.role)
			assert.Equal(t, tt.want, result)
			if !tt.wantErr {
				assert.Equal(t, tt.role, user.Role)
			}
		})
	}
}

func TestTenantUser_StructFields(t *testing.T) {
	now := time.Now()
	user := &TenantUser{
		ID:        "test-id",
		TenantID:  "tenant-id",
		Name:      "Test User",
		Email:     "test@example.com",
		Role:      RoleTenantUser,
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, "test-id", user.ID)
	assert.Equal(t, "tenant-id", user.TenantID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, RoleTenantUser, user.Role)
	assert.Equal(t, now, user.CreatedAt)
	assert.Equal(t, now, user.UpdatedAt)
}

func TestRole_IsValid(t *testing.T) {
	tests := []struct {
		name string
		role Role
		want bool
	}{
		{
			name: "user role",
			role: RoleTenantUser,
			want: true,
		},
		{
			name: "admin role",
			role: RoleTenantAdmin,
			want: true,
		},
		{
			name: "invalid role",
			role: Role("invalid"),
			want: false,
		},
		{
			name: "empty role",
			role: Role(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.role.IsValid())
		})
	}
}

func TestTenantUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    *TenantUser
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid user",
			user: &TenantUser{
				TenantID: "tenant-123",
				Name:     "Test User",
				Email:    "test@example.com",
				Role:     RoleTenantUser,
			},
			wantErr: false,
		},
		{
			name: "valid user with empty role",
			user: &TenantUser{
				TenantID: "tenant-123",
				Name:     "Test User",
				Email:    "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "missing tenant ID",
			user: &TenantUser{
				Name:  "Test User",
				Email: "test@example.com",
			},
			wantErr: true,
			errMsg:  "tenant ID is required",
		},
		{
			name: "missing name",
			user: &TenantUser{
				TenantID: "tenant-123",
				Email:    "test@example.com",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "missing email",
			user: &TenantUser{
				TenantID: "tenant-123",
				Name:     "Test User",
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "invalid role",
			user: &TenantUser{
				TenantID: "tenant-123",
				Name:     "Test User",
				Email:    "test@example.com",
				Role:     Role("invalid"),
			},
			wantErr: true,
			errMsg:  "invalid user role",
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
