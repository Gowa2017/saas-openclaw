package platform

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetUserInfo(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *httptest.Server
		token      string
		expectErr  bool
		errType    error
	}{
		{
			name: "successful user info fetch",
			setup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "/api/user/info", r.URL.Path)
					assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(UserInfo{
						ID:       "user-123",
						Name:     "Test User",
						Email:    "test@example.com",
						TenantID: "tenant-456",
						Role:     "admin",
					})
				}))
			},
			token:     "test-token",
			expectErr: false,
		},
		{
			name: "unauthorized",
			setup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
				}))
			},
			token:     "invalid-token",
			expectErr: true,
			errType:   ErrUnauthorized,
		},
		{
			name: "user not found",
			setup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}))
			},
			token:     "test-token",
			expectErr: true,
			errType:   ErrUserNotFound,
		},
		{
			name: "platform unavailable",
			setup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusServiceUnavailable)
				}))
			},
			token:     "test-token",
			expectErr: true,
			errType:   ErrPlatformUnavailable,
		},
		{
			name: "invalid response - missing user id",
			setup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]string{
						"name":     "Test User",
						"email":    "test@example.com",
						"tenantId": "tenant-456",
					})
				}))
			},
			token:     "test-token",
			expectErr: true,
			errType:   ErrInvalidResponse,
		},
		{
			name: "invalid JSON response",
			setup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte("invalid-json"))
				}))
			},
			token:     "test-token",
			expectErr: true,
			errType:   ErrInvalidResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setup()
			defer server.Close()

			client := NewClient(server.URL)
			userInfo, err := client.GetUserInfo(context.Background(), tt.token)

			if tt.expectErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.errType)
				assert.Nil(t, userInfo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, userInfo)
				assert.Equal(t, "user-123", userInfo.ID)
				assert.Equal(t, "Test User", userInfo.Name)
				assert.Equal(t, "test@example.com", userInfo.Email)
				assert.Equal(t, "tenant-456", userInfo.TenantID)
				assert.Equal(t, "admin", userInfo.Role)
			}
		})
	}
}

func TestClient_GetUserInfo_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, WithTimeout(10*time.Millisecond))
	userInfo, err := client.GetUserInfo(context.Background(), "test-token")

	assert.Error(t, err)
	assert.Nil(t, userInfo)
}

func TestClient_GetUserInfo_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	userInfo, err := client.GetUserInfo(ctx, "test-token")

	assert.Error(t, err)
	assert.Nil(t, userInfo)
}

func TestClient_HealthCheck(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *httptest.Server
		expectErr bool
	}{
		{
			name: "healthy",
			setup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "/health", r.URL.Path)
					w.WriteHeader(http.StatusOK)
				}))
			},
			expectErr: false,
		},
		{
			name: "unhealthy",
			setup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusServiceUnavailable)
				}))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setup()
			defer server.Close()

			client := NewClient(server.URL)
			err := client.HealthCheck(context.Background())

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	t.Run("default configuration", func(t *testing.T) {
		client := NewClient("https://platform.example.com")
		assert.NotNil(t, client)
		assert.Equal(t, "https://platform.example.com", client.baseURL)
		assert.Equal(t, 5*time.Second, client.timeout)
	})

	t.Run("with custom timeout", func(t *testing.T) {
		client := NewClient("https://platform.example.com", WithTimeout(10*time.Second))
		assert.Equal(t, 10*time.Second, client.timeout)
	})

	t.Run("with custom HTTP client", func(t *testing.T) {
		customHTTPClient := &http.Client{Timeout: 15 * time.Second}
		client := NewClient("https://platform.example.com", WithHTTPClient(customHTTPClient))
		assert.Equal(t, customHTTPClient, client.httpClient)
	})
}

func TestClient_ServerDown(t *testing.T) {
	client := NewClient("http://localhost:12345") // Non-existent server
	userInfo, err := client.GetUserInfo(context.Background(), "test-token")

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrPlatformUnavailable)
	assert.Nil(t, userInfo)
}

func TestUserInfo_JSONSerialization(t *testing.T) {
	user := UserInfo{
		ID:       "user-123",
		Name:     "Test User",
		Email:    "test@example.com",
		TenantID: "tenant-456",
		Role:     "admin",
	}

	data, err := json.Marshal(user)
	require.NoError(t, err)

	var decoded UserInfo
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, user, decoded)
}
