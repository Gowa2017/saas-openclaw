package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/gowa/saas-openclaw/backend/pkg/jwt"
)

func TestAdminAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	secret := "test-admin-secret-key"
	validator := jwt.NewAdminValidator(secret)

	tests := []struct {
		name       string
		setupToken func() string
		wantStatus int
		wantAdmin  bool
	}{
		{
			name: "valid token",
			setupToken: func() string {
				claims := &jwt.AdminClaims{
					AdminID:  "admin-123",
					Username: "testadmin",
					Role:     "admin",
					RegisteredClaims: jwtlib.RegisteredClaims{
						ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
						IssuedAt:  jwtlib.NewNumericDate(time.Now()),
						Issuer:    jwt.TokenIssuer,
					},
				}
				token, _ := jwt.GenerateAdminToken(secret, claims)
				return token
			},
			wantStatus: http.StatusOK,
			wantAdmin:  true,
		},
		{
			name: "valid super_admin token",
			setupToken: func() string {
				claims := &jwt.AdminClaims{
					AdminID:  "super-123",
					Username: "superadmin",
					Role:     "super_admin",
					RegisteredClaims: jwtlib.RegisteredClaims{
						ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
						IssuedAt:  jwtlib.NewNumericDate(time.Now()),
						Issuer:    jwt.TokenIssuer,
					},
				}
				token, _ := jwt.GenerateAdminToken(secret, claims)
				return token
			},
			wantStatus: http.StatusOK,
			wantAdmin:  true,
		},
		{
			name: "missing token",
			setupToken: func() string {
				return ""
			},
			wantStatus: http.StatusUnauthorized,
			wantAdmin:  false,
		},
		{
			name: "expired token",
			setupToken: func() string {
				claims := &jwt.AdminClaims{
					AdminID:  "admin-123",
					Username: "testadmin",
					Role:     "admin",
					RegisteredClaims: jwtlib.RegisteredClaims{
						ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(-1 * time.Hour)),
						IssuedAt:  jwtlib.NewNumericDate(time.Now().Add(-2 * time.Hour)),
						Issuer:    jwt.TokenIssuer,
					},
				}
				token, _ := jwt.GenerateAdminToken(secret, claims)
				return token
			},
			wantStatus: http.StatusUnauthorized,
			wantAdmin:  false,
		},
		{
			name: "invalid token",
			setupToken: func() string {
				return "invalid-token-string"
			},
			wantStatus: http.StatusUnauthorized,
			wantAdmin:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router
			router := gin.New()
			router.Use(AdminAuth(validator, logger))
			router.GET("/test", func(c *gin.Context) {
				admin := GetAdminContext(c)
				if admin != nil {
					c.JSON(http.StatusOK, gin.H{
						"adminId":  admin.ID,
						"username": admin.Username,
						"role":     admin.Role,
					})
				} else {
					c.JSON(http.StatusOK, gin.H{"message": "ok"})
				}
			})

			// Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			token := tt.setupToken()
			if token != "" {
				req.Header.Set("Authorization", "Bearer "+token)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantAdmin {
				var response map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)
				_ = response // Just verify the request succeeded
			}
		})
	}
}

func TestAdminAuth_TokenExtraction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	secret := "test-secret"
	validator := jwt.NewAdminValidator(secret)

	claims := &jwt.AdminClaims{
		AdminID:  "admin-123",
		Username: "testadmin",
		Role:     "admin",
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			Issuer:    jwt.TokenIssuer,
		},
	}
	token, _ := jwt.GenerateAdminToken(secret, claims)

	tests := []struct {
		name   string
		header string
		value  string
	}{
		{
			name:   "Authorization Bearer header",
			header: "Authorization",
			value:  "Bearer " + token,
		},
		{
			name:   "X-Admin-Token header",
			header: "X-Admin-Token",
			value:  token,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(AdminAuth(validator, logger))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "ok"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set(tt.header, tt.value)

			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestRequireAdminRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		adminRole  string
		require    []string
		wantStatus int
	}{
		{
			name:       "admin role required - has admin",
			adminRole:  "admin",
			require:    []string{"admin"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "super_admin role required - has super_admin",
			adminRole:  "super_admin",
			require:    []string{"super_admin"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "any admin role - has admin",
			adminRole:  "admin",
			require:    []string{"admin", "super_admin"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "super_admin required - has admin",
			adminRole:  "admin",
			require:    []string{"super_admin"},
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "no admin context",
			adminRole:  "", // No context
			require:    []string{"admin"},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			// Setup context if role provided
			if tt.adminRole != "" {
				router.Use(func(c *gin.Context) {
					adminCtx := &AdminContext{
						ID:       "admin-123",
						Username: "testadmin",
						Role:     tt.adminRole,
					}
					c.Set(AdminContextKey, adminCtx)
					c.Set(AdminRoleKey, tt.adminRole)
					c.Next()
				})
			}

			router.Use(RequireAdminRole(tt.require...))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "ok"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestRequireSuperAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		adminRole  string
		wantStatus int
	}{
		{
			name:       "super_admin allowed",
			adminRole:  "super_admin",
			wantStatus: http.StatusOK,
		},
		{
			name:       "admin forbidden",
			adminRole:  "admin",
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				adminCtx := &AdminContext{
					ID:       "admin-123",
					Username: "testadmin",
					Role:     tt.adminRole,
				}
				c.Set(AdminContextKey, adminCtx)
				c.Set(AdminRoleKey, tt.adminRole)
				c.Next()
			})
			router.Use(RequireSuperAdmin())
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "ok"})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestGetAdminContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns admin context when set", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			adminCtx := &AdminContext{
				ID:       "admin-123",
				Username: "testadmin",
				Role:     "admin",
			}
			c.Set(AdminContextKey, adminCtx)
			c.Next()
		})
		router.GET("/test", func(c *gin.Context) {
			admin := GetAdminContext(c)
			require.NotNil(t, admin)
			assert.Equal(t, "admin-123", admin.ID)
			assert.Equal(t, "testadmin", admin.Username)
			assert.Equal(t, "admin", admin.Role)
			c.JSON(http.StatusOK, gin.H{})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)
	})

	t.Run("returns nil when not set", func(t *testing.T) {
		router := gin.New()
		router.GET("/test", func(c *gin.Context) {
			admin := GetAdminContext(c)
			assert.Nil(t, admin)
			c.JSON(http.StatusOK, gin.H{})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)
	})
}
