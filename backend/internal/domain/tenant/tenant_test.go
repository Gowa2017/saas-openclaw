package tenant

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTenant_IsActive(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   bool
	}{
		{
			name:   "active tenant",
			status: StatusActive,
			want:   true,
		},
		{
			name:   "inactive tenant",
			status: StatusInactive,
			want:   false,
		},
		{
			name:   "suspended tenant",
			status: StatusSuspended,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tenant := &Tenant{
				ID:     "test-id",
				Name:   "Test Tenant",
				Status: tt.status,
			}
			assert.Equal(t, tt.want, tenant.IsActive())
		})
	}
}

func TestTenant_StructFields(t *testing.T) {
	now := time.Now()
	tenant := &Tenant{
		ID:        "test-id",
		Name:      "Test Tenant",
		Status:    StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, "test-id", tenant.ID)
	assert.Equal(t, "Test Tenant", tenant.Name)
	assert.Equal(t, StatusActive, tenant.Status)
	assert.Equal(t, now, tenant.CreatedAt)
	assert.Equal(t, now, tenant.UpdatedAt)
}

func TestStatus_Constants(t *testing.T) {
	assert.Equal(t, Status("active"), StatusActive)
	assert.Equal(t, Status("inactive"), StatusInactive)
	assert.Equal(t, Status("suspended"), StatusSuspended)
}

func TestStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   bool
	}{
		{
			name:   "active status",
			status: StatusActive,
			want:   true,
		},
		{
			name:   "inactive status",
			status: StatusInactive,
			want:   true,
		},
		{
			name:   "suspended status",
			status: StatusSuspended,
			want:   true,
		},
		{
			name:   "invalid status",
			status: Status("invalid"),
			want:   false,
		},
		{
			name:   "empty status",
			status: Status(""),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.IsValid())
		})
	}
}

func TestTenant_Validate(t *testing.T) {
	tests := []struct {
		name    string
		tenant  *Tenant
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid tenant",
			tenant: &Tenant{
				Name:   "Test Tenant",
				Status: StatusActive,
			},
			wantErr: false,
		},
		{
			name: "valid tenant with empty status",
			tenant: &Tenant{
				Name: "Test Tenant",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			tenant: &Tenant{
				Status: StatusActive,
			},
			wantErr: true,
			errMsg:  "tenant name is required",
		},
		{
			name: "invalid status",
			tenant: &Tenant{
				Name:   "Test Tenant",
				Status: Status("invalid"),
			},
			wantErr: true,
			errMsg:  "invalid tenant status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tenant.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
