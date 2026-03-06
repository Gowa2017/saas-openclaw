package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gowa/saas-openclaw/backend/internal/infrastructure/platform"
	"github.com/gowa/saas-openclaw/backend/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key-for-jwt"

func init() {
	gin.SetMode(gin.TestMode)
}

func TestPlatformAuth(t *testing.T) {
	tests := []struct {
		name           string
		setupToken     func() string
		setupServer    func() *httptest.Server
		expectStatus   int
		expectUser     bool
		expectErrCode  string
	}{
		{
			name: "successful authentication",
			setupToken: func() string {
				token, _ := jwt.CreateTestToken(testSecret, "user-123", "test@example.com", "tenant-456", time.Now().Add(time.Hour))
				return token
			},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(platform.UserInfo{
						ID:       "user-123",
						Name:     "Test User",
						Email:    "test@example.com",
						TenantID: "tenant-456",
						Role:     "admin",
					})
				}))
			},
			expectStatus: http.StatusOK,
			expectUser:   true,
		},
		{
			name: "missing token",
			setupToken: func() string {
				return ""
			},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
			expectStatus:  http.StatusUnauthorized,
			expectErrCode: "MISSING_TOKEN",
		},
		{
			name: "invalid token",
			setupToken: func() string {
				return "invalid-token"
			},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
			expectStatus:  http.StatusUnauthorized,
			expectErrCode: "INVALID_TOKEN",
		},
		{
			name: "expired token",
			setupToken: func() string {
				token, _ := jwt.CreateTestToken(testSecret, "user-123", "test@example.com", "tenant-456", time.Now().Add(-time.Hour))
				return token
			},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
			expectStatus:  http.StatusUnauthorized,
			expectErrCode: "INVALID_TOKEN",
		},
		{
			name: "platform returns unauthorized",
			setupToken: func() string {
				token, _ := jwt.CreateTestToken(testSecret, "user-123", "test@example.com", "tenant-456", time.Now().Add(time.Hour))
				return token
			},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
				}))
			},
			expectStatus:  http.StatusUnauthorized,
			expectErrCode: "USER_INFO_FAILED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()
			defer server.Close()

			validator := jwt.NewValidator(testSecret)
			platformClient := platform.NewClient(server.URL)

			router := gin.New()
			router.Use(PlatformAuth(validator, platformClient))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"success": true})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			token := tt.setupToken()
			if token != "" {
				req.Header.Set("X-Platform-Token", token)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectStatus, w.Code)

			if tt.expectErrCode != "" {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				errData := resp["error"].(map[string]interface{})
				assert.Equal(t, tt.expectErrCode, errData["code"])
			}

			if tt.expectUser {
				// Verify user is set in context by making a second request that checks the context
				router2 := gin.New()
				router2.Use(PlatformAuth(validator, platformClient))
				router2.GET("/check-user", func(c *gin.Context) {
					user := GetCurrentUser(c)
					assert.NotNil(t, user)
					assert.Equal(t, "user-123", user.ID)
					c.JSON(http.StatusOK, gin.H{"user": user})
				})

				req2 := httptest.NewRequest(http.MethodGet, "/check-user", nil)
				req2.Header.Set("X-Platform-Token", token)
				w2 := httptest.NewRecorder()
				router2.ServeHTTP(w2, req2)
			}
		})
	}
}

func TestJWTAuth(t *testing.T) {
	tests := []struct {
		name          string
		setupToken    func() string
		expectStatus  int
		expectUser    bool
		expectErrCode string
	}{
		{
			name: "successful authentication",
			setupToken: func() string {
				token, _ := jwt.CreateTestToken(testSecret, "user-123", "test@example.com", "tenant-456", time.Now().Add(time.Hour))
				return token
			},
			expectStatus: http.StatusOK,
			expectUser:   true,
		},
		{
			name: "missing token",
			setupToken: func() string {
				return ""
			},
			expectStatus:  http.StatusUnauthorized,
			expectErrCode: "MISSING_TOKEN",
		},
		{
			name: "invalid token",
			setupToken: func() string {
				return "invalid-token"
			},
			expectStatus:  http.StatusUnauthorized,
			expectErrCode: "INVALID_TOKEN",
		},
		{
			name: "expired token",
			setupToken: func() string {
				token, _ := jwt.CreateTestToken(testSecret, "user-123", "test@example.com", "tenant-456", time.Now().Add(-time.Hour))
				return token
			},
			expectStatus:  http.StatusUnauthorized,
			expectErrCode: "INVALID_TOKEN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := jwt.NewValidator(testSecret)

			router := gin.New()
			router.Use(JWTAuth(validator))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"success": true})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			token := tt.setupToken()
			if token != "" {
				req.Header.Set("X-Platform-Token", token)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectStatus, w.Code)

			if tt.expectErrCode != "" {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				errData := resp["error"].(map[string]interface{})
				assert.Equal(t, tt.expectErrCode, errData["code"])
			}
		})
	}
}

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name        string
		setupHeader func(req *http.Request)
		expectToken string
	}{
		{
			name: "X-Platform-Token header",
			setupHeader: func(req *http.Request) {
				req.Header.Set("X-Platform-Token", "platform-token-123")
			},
			expectToken: "platform-token-123",
		},
		{
			name: "Authorization Bearer header",
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer bearer-token-456")
			},
			expectToken: "bearer-token-456",
		},
		{
			name: "X-Platform-Token takes precedence",
			setupHeader: func(req *http.Request) {
				req.Header.Set("X-Platform-Token", "platform-token")
				req.Header.Set("Authorization", "Bearer bearer-token")
			},
			expectToken: "platform-token",
		},
		{
			name: "no token",
			setupHeader: func(req *http.Request) {
				// No headers
			},
			expectToken: "",
		},
		{
			name: "malformed Authorization header",
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "InvalidFormat")
			},
			expectToken: "",
		},
		{
			name: "case insensitive bearer",
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "bearer lower-case-token")
			},
			expectToken: "lower-case-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			tt.setupHeader(req)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			token := extractToken(c)
			assert.Equal(t, tt.expectToken, token)
		})
	}
}

func TestGetCurrentUser(t *testing.T) {
	t.Run("user exists in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		user := &UserContext{
			ID:       "user-123",
			Name:     "Test User",
			Email:    "test@example.com",
			TenantID: "tenant-456",
			Role:     "admin",
		}
		c.Set(UserContextKey, user)

		result := GetCurrentUser(c)
		assert.Equal(t, user, result)
	})

	t.Run("user not in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		result := GetCurrentUser(c)
		assert.Nil(t, result)
	})
}

func TestGetCurrentUserID(t *testing.T) {
	t.Run("user ID exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(UserIDKey, "user-123")

		result := GetCurrentUserID(c)
		assert.Equal(t, "user-123", result)
	})

	t.Run("user ID not in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		result := GetCurrentUserID(c)
		assert.Empty(t, result)
	})
}

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name         string
		userRole     string
		requireRoles []string
		expectStatus int
	}{
		{
			name:         "user has required role",
			userRole:     "admin",
			requireRoles: []string{"admin", "superadmin"},
			expectStatus: http.StatusOK,
		},
		{
			name:         "user does not have required role",
			userRole:     "user",
			requireRoles: []string{"admin", "superadmin"},
			expectStatus: http.StatusForbidden,
		},
		{
			name:         "no user in context",
			userRole:     "",
			requireRoles: []string{"admin"},
			expectStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			// Set up user in context if role is specified
			if tt.userRole != "" {
				router.Use(func(c *gin.Context) {
					c.Set(UserContextKey, &UserContext{Role: tt.userRole})
					c.Set(UserRoleKey, tt.userRole)
					c.Next()
				})
			}

			router.Use(RequireRole(tt.requireRoles...))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"success": true})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectStatus, w.Code)
		})
	}
}

func TestRequireTenant(t *testing.T) {
	tests := []struct {
		name         string
		tenantID     string
		expectStatus int
	}{
		{
			name:         "user has tenant",
			tenantID:     "tenant-123",
			expectStatus: http.StatusOK,
		},
		{
			name:         "user has no tenant",
			tenantID:     "",
			expectStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set(UserContextKey, &UserContext{TenantID: tt.tenantID})
				c.Set(UserTenantIDKey, tt.tenantID)
				c.Next()
			})
			router.Use(RequireTenant())
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"success": true})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectStatus, w.Code)
		})
	}
}

func TestRespondWithError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	respondWithError(c, http.StatusUnauthorized, "TEST_ERROR", "Test error message")

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Nil(t, resp["data"])
	assert.NotNil(t, resp["error"])
	assert.Nil(t, resp["meta"])

	errData := resp["error"].(map[string]interface{})
	assert.Equal(t, "TEST_ERROR", errData["code"])
	assert.Equal(t, "Test error message", errData["message"])
}
