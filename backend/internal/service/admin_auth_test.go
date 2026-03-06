package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/gowa/saas-openclaw/backend/internal/domain/user"
	"github.com/gowa/saas-openclaw/backend/pkg/jwt"
)

// MockAdminAuthRepository is a mock implementation of AdminAuthRepository
type MockAdminAuthRepository struct {
	mock.Mock
}

func (m *MockAdminAuthRepository) GetByUsername(username string) (*user.AdminUser, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.AdminUser), args.Error(1)
}

func (m *MockAdminAuthRepository) VerifyPassword(admin *user.AdminUser, password string) bool {
	args := m.Called(admin, password)
	return args.Bool(0)
}

func TestAdminAuthService_Login(t *testing.T) {
	secret := "test-admin-secret-key"
	tokenExp := 24 * time.Hour
	logger := zap.NewNop()

	tests := []struct {
		name        string
		setupMock   func(mockRepo *MockAdminAuthRepository)
		req         *LoginRequest
		wantErr     error
		validateRes func(t *testing.T, res *LoginResponse)
	}{
		{
			name: "successful login",
			setupMock: func(mockRepo *MockAdminAuthRepository) {
				admin := &user.AdminUser{
					ID:           "admin-123",
					Username:     "testadmin",
					PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // hash for "password123"
					Name:         "Test Admin",
					Email:        "admin@test.com",
					Role:         user.AdminRoleAdmin,
				}
				mockRepo.On("GetByUsername", "testadmin").Return(admin, nil)
				mockRepo.On("VerifyPassword", admin, "password123").Return(true)
			},
			req: &LoginRequest{
				Username: "testadmin",
				Password: "password123",
			},
			wantErr: nil,
			validateRes: func(t *testing.T, res *LoginResponse) {
				assert.NotEmpty(t, res.Token)
				assert.Equal(t, "admin-123", res.AdminID)
				assert.Equal(t, "testadmin", res.Username)
				assert.Equal(t, "admin", res.Role)
				assert.Greater(t, res.ExpiresAt, time.Now().Unix())

				// Verify token is valid
				validator := jwt.NewAdminValidator(secret)
				claims, err := validator.ValidateAdminToken(res.Token)
				require.NoError(t, err)
				assert.Equal(t, "admin-123", claims.AdminID)
				assert.Equal(t, jwt.TokenIssuer, claims.Issuer)
			},
		},
		{
			name: "successful login as super_admin",
			setupMock: func(mockRepo *MockAdminAuthRepository) {
				admin := &user.AdminUser{
					ID:           "super-123",
					Username:     "superadmin",
					PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
					Name:         "Super Admin",
					Email:        "super@test.com",
					Role:         user.AdminRoleSuperAdmin,
				}
				mockRepo.On("GetByUsername", "superadmin").Return(admin, nil)
				mockRepo.On("VerifyPassword", admin, "password123").Return(true)
			},
			req: &LoginRequest{
				Username: "superadmin",
				Password: "password123",
			},
			wantErr: nil,
			validateRes: func(t *testing.T, res *LoginResponse) {
				assert.Equal(t, "super_admin", res.Role)
			},
		},
		{
			name: "user not found",
			setupMock: func(mockRepo *MockAdminAuthRepository) {
				mockRepo.On("GetByUsername", "nonexistent").Return(nil, errors.New("user not found"))
			},
			req: &LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "invalid password",
			setupMock: func(mockRepo *MockAdminAuthRepository) {
				admin := &user.AdminUser{
					ID:           "admin-123",
					Username:     "testadmin",
					PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
					Name:         "Test Admin",
					Email:        "admin@test.com",
					Role:         user.AdminRoleAdmin,
				}
				mockRepo.On("GetByUsername", "testadmin").Return(admin, nil)
				mockRepo.On("VerifyPassword", admin, "wrongpassword").Return(false)
			},
			req: &LoginRequest{
				Username: "testadmin",
				Password: "wrongpassword",
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "empty username",
			setupMock: func(mockRepo *MockAdminAuthRepository) {
				mockRepo.On("GetByUsername", "").Return(nil, errors.New("username is required"))
			},
			req: &LoginRequest{
				Username: "",
				Password: "password123",
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "empty password",
			setupMock: func(mockRepo *MockAdminAuthRepository) {
				admin := &user.AdminUser{
					ID:           "admin-123",
					Username:     "testadmin",
					PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
					Name:         "Test Admin",
					Email:        "admin@test.com",
					Role:         user.AdminRoleAdmin,
				}
				mockRepo.On("GetByUsername", "testadmin").Return(admin, nil)
				mockRepo.On("VerifyPassword", admin, "").Return(false)
			},
			req: &LoginRequest{
				Username: "testadmin",
				Password: "",
			},
			wantErr: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAdminAuthRepository)
			tt.setupMock(mockRepo)

			service := NewAdminAuthService(mockRepo, secret, tokenExp, logger)
			res, err := service.Login(context.Background(), tt.req)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
				if tt.validateRes != nil {
					tt.validateRes(t, res)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAdminAuthService_GenerateToken(t *testing.T) {
	secret := "test-admin-secret-key"
	tokenExp := 24 * time.Hour
	logger := zap.NewNop()

	mockRepo := new(MockAdminAuthRepository)
	service := NewAdminAuthService(mockRepo, secret, tokenExp, logger)

	admin := &user.AdminUser{
		ID:       "admin-123",
		Username: "testadmin",
		Role:     user.AdminRoleAdmin,
	}

	token, err := service.GenerateToken(admin)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token
	validator := jwt.NewAdminValidator(secret)
	claims, err := validator.ValidateAdminToken(token)
	require.NoError(t, err)
	assert.Equal(t, "admin-123", claims.AdminID)
	assert.Equal(t, "testadmin", claims.Username)
	assert.Equal(t, "admin", claims.Role)
}
