package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/gowa/saas-openclaw/backend/internal/domain/user"
	"github.com/gowa/saas-openclaw/backend/internal/service"
)

// mockAdminAuthRepository implements service.AdminAuthRepository for testing
type mockAdminAuthRepository struct {
	getByUsernameFunc func(username string) (*user.AdminUser, error)
	verifyPasswordFunc func(admin *user.AdminUser, password string) bool
}

func (m *mockAdminAuthRepository) GetByUsername(username string) (*user.AdminUser, error) {
	return m.getByUsernameFunc(username)
}

func (m *mockAdminAuthRepository) VerifyPassword(admin *user.AdminUser, password string) bool {
	return m.verifyPasswordFunc(admin, password)
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	tests := []struct {
		name          string
		reqBody       interface{}
		mockRepoSetup func() *mockAdminAuthRepository
		wantStatus    int
		wantErr       bool
		errCode       string
	}{
		{
			name: "successful login",
			reqBody: LoginRequest{
				Username: "testadmin",
				Password: "password123",
			},
			mockRepoSetup: func() *mockAdminAuthRepository {
				return &mockAdminAuthRepository{
					getByUsernameFunc: func(username string) (*user.AdminUser, error) {
						return &user.AdminUser{
							ID:           "admin-123",
							Username:     "testadmin",
							PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // hash for "password123"
							Name:         "Test Admin",
							Email:        "admin@test.com",
							Role:         user.AdminRoleAdmin,
						}, nil
					},
					verifyPasswordFunc: func(admin *user.AdminUser, password string) bool {
						return password == "password123"
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "missing username",
			reqBody: map[string]interface{}{
				"password": "password123",
			},
			mockRepoSetup: func() *mockAdminAuthRepository {
				return &mockAdminAuthRepository{}
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
			errCode:    "INVALID_REQUEST",
		},
		{
			name: "missing password",
			reqBody: map[string]interface{}{
				"username": "testadmin",
			},
			mockRepoSetup: func() *mockAdminAuthRepository {
				return &mockAdminAuthRepository{}
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
			errCode:    "INVALID_REQUEST",
		},
		{
			name: "invalid credentials - user not found",
			reqBody: LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			mockRepoSetup: func() *mockAdminAuthRepository {
				return &mockAdminAuthRepository{
					getByUsernameFunc: func(username string) (*user.AdminUser, error) {
						return nil, assert.AnError
					},
				}
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    true,
			errCode:    "UNAUTHORIZED",
		},
		{
			name: "invalid credentials - wrong password",
			reqBody: LoginRequest{
				Username: "testadmin",
				Password: "wrongpassword",
			},
			mockRepoSetup: func() *mockAdminAuthRepository {
				return &mockAdminAuthRepository{
					getByUsernameFunc: func(username string) (*user.AdminUser, error) {
						return &user.AdminUser{
							ID:           "admin-123",
							Username:     "testadmin",
							PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
							Name:         "Test Admin",
							Email:        "admin@test.com",
							Role:         user.AdminRoleAdmin,
						}, nil
					},
					verifyPasswordFunc: func(admin *user.AdminUser, password string) bool {
						return false // Wrong password
					},
				}
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    true,
			errCode:    "UNAUTHORIZED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock repository
			mockRepo := tt.mockRepoSetup()

			// Create service with mock
			authService := service.NewAdminAuthService(mockRepo, "test-secret-key", 24*time.Hour, logger)

			// Create handler
			handler := NewAuthHandler(authService, logger)

			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.reqBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/auth/login", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Execute
			handler.Login(c)

			// Assert
			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tt.wantErr {
				assert.Nil(t, response["data"])
				errData, ok := response["error"].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, tt.errCode, errData["code"])
			} else {
				assert.Nil(t, response["error"])
				data, ok := response["data"].(map[string]interface{})
				require.True(t, ok)
				assert.NotEmpty(t, data["token"])
				assert.NotEmpty(t, data["adminId"])
				assert.NotEmpty(t, data["username"])
				assert.NotEmpty(t, data["role"])
			}
		})
	}
}

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	v1 := router.Group("/v1")

	handler := &AuthHandler{}
	RegisterRoutes(v1, handler)

	// Test that route is registered
	routes := router.Routes()
	found := false
	for _, route := range routes {
		if route.Path == "/v1/admin/auth/login" && route.Method == http.MethodPost {
			found = true
			break
		}
	}
	assert.True(t, found, "Route POST /v1/admin/auth/login should be registered")
}
