package jwt

import (
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminClaims_ValidateAdminToken(t *testing.T) {
	secret := "test-admin-secret-key"
	validator := NewAdminValidator(secret)

	tests := []struct {
		name      string
		setupFunc func() string // Returns token string
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid admin token",
			setupFunc: func() string {
				claims := &AdminClaims{
					AdminID:  "admin-123",
					Username: "testadmin",
					Role:     "admin",
					RegisteredClaims: jwtlib.RegisteredClaims{
						ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
						IssuedAt:  jwtlib.NewNumericDate(time.Now()),
						Issuer:    TokenIssuer,
					},
				}
				token, _ := GenerateAdminToken(secret, claims)
				return token
			},
			wantErr: false,
		},
		{
			name: "valid super_admin token",
			setupFunc: func() string {
				claims := &AdminClaims{
					AdminID:  "super-123",
					Username: "superadmin",
					Role:     "super_admin",
					RegisteredClaims: jwtlib.RegisteredClaims{
						ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
						IssuedAt:  jwtlib.NewNumericDate(time.Now()),
						Issuer:    TokenIssuer,
					},
				}
				token, _ := GenerateAdminToken(secret, claims)
				return token
			},
			wantErr: false,
		},
		{
			name: "expired token",
			setupFunc: func() string {
				claims := &AdminClaims{
					AdminID:  "admin-123",
					Username: "testadmin",
					Role:     "admin",
					RegisteredClaims: jwtlib.RegisteredClaims{
						ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(-1 * time.Hour)),
						IssuedAt:  jwtlib.NewNumericDate(time.Now().Add(-2 * time.Hour)),
						Issuer:    TokenIssuer,
					},
				}
				token, _ := GenerateAdminToken(secret, claims)
				return token
			},
			wantErr: true,
			errMsg:  "token has expired",
		},
		{
			name: "empty token",
			setupFunc: func() string {
				return ""
			},
			wantErr: true,
			errMsg:  "token is empty",
		},
		{
			name: "invalid signature",
			setupFunc: func() string {
				claims := &AdminClaims{
					AdminID:  "admin-123",
					Username: "testadmin",
					Role:     "admin",
					RegisteredClaims: jwtlib.RegisteredClaims{
						ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
						IssuedAt:  jwtlib.NewNumericDate(time.Now()),
						Issuer:    TokenIssuer,
					},
				}
				token, _ := GenerateAdminToken("wrong-secret", claims)
				return token
			},
			wantErr: true,
			errMsg:  "token signature is invalid",
		},
		{
			name: "missing admin_id",
			setupFunc: func() string {
				claims := &AdminClaims{
					AdminID:  "",
					Username: "testadmin",
					Role:     "admin",
					RegisteredClaims: jwtlib.RegisteredClaims{
						ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
						IssuedAt:  jwtlib.NewNumericDate(time.Now()),
						Issuer:    TokenIssuer,
					},
				}
				token, _ := GenerateAdminToken(secret, claims)
				return token
			},
			wantErr: true,
			errMsg:  "admin_id claim is required",
		},
		{
			name: "missing username",
			setupFunc: func() string {
				claims := &AdminClaims{
					AdminID:  "admin-123",
					Username: "",
					Role:     "admin",
					RegisteredClaims: jwtlib.RegisteredClaims{
						ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
						IssuedAt:  jwtlib.NewNumericDate(time.Now()),
						Issuer:    TokenIssuer,
					},
				}
				token, _ := GenerateAdminToken(secret, claims)
				return token
			},
			wantErr: true,
			errMsg:  "username claim is required",
		},
		{
			name: "missing role",
			setupFunc: func() string {
				claims := &AdminClaims{
					AdminID:  "admin-123",
					Username: "testadmin",
					Role:     "",
					RegisteredClaims: jwtlib.RegisteredClaims{
						ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
						IssuedAt:  jwtlib.NewNumericDate(time.Now()),
						Issuer:    TokenIssuer,
					},
				}
				token, _ := GenerateAdminToken(secret, claims)
				return token
			},
			wantErr: true,
			errMsg:  "role claim is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupFunc()
			claims, err := validator.ValidateAdminToken(token)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, claims)
			} else {
				require.NoError(t, err)
				require.NotNil(t, claims)
				assert.NotEmpty(t, claims.AdminID)
				assert.NotEmpty(t, claims.Username)
				assert.NotEmpty(t, claims.Role)
			}
		})
	}
}

func TestGenerateAdminToken(t *testing.T) {
	secret := "test-secret"

	t.Run("generates valid token", func(t *testing.T) {
		claims := &AdminClaims{
			AdminID:  "admin-123",
			Username: "testadmin",
			Role:     "admin",
			RegisteredClaims: jwtlib.RegisteredClaims{
				ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwtlib.NewNumericDate(time.Now()),
				Issuer:    TokenIssuer,
			},
		}

		token, err := GenerateAdminToken(secret, claims)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify token can be parsed
		validator := NewAdminValidator(secret)
		parsedClaims, err := validator.ValidateAdminToken(token)
		require.NoError(t, err)
		assert.Equal(t, claims.AdminID, parsedClaims.AdminID)
		assert.Equal(t, claims.Username, parsedClaims.Username)
		assert.Equal(t, claims.Role, parsedClaims.Role)
		assert.Equal(t, TokenIssuer, parsedClaims.Issuer)
	})
}

func TestAdminValidator_WithDifferentSecrets(t *testing.T) {
	secret1 := "secret-one"
	secret2 := "secret-two"

	validator1 := NewAdminValidator(secret1)

	// Create token with secret1
	claims := &AdminClaims{
		AdminID:  "admin-123",
		Username: "testadmin",
		Role:     "admin",
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			Issuer:    TokenIssuer,
		},
	}
	token, _ := GenerateAdminToken(secret1, claims)

	// Should validate with same secret
	_, err := validator1.ValidateAdminToken(token)
	require.NoError(t, err)

	// Should fail with different secret
	validator2 := NewAdminValidator(secret2)
	_, err = validator2.ValidateAdminToken(token)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "signature is invalid")
}
