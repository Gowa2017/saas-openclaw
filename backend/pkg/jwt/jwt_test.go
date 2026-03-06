package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key-for-jwt-validation"

func TestValidator_ValidateToken(t *testing.T) {
	validator := NewValidator(testSecret)

	tests := []struct {
		name        string
		setupToken  func() string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid token",
			setupToken: func() string {
				token, _ := CreateTestToken(testSecret, "user-123", "test@example.com", "tenant-456", time.Now().Add(time.Hour))
				return token
			},
			expectError: false,
		},
		{
			name: "expired token",
			setupToken: func() string {
				token, _ := CreateTestToken(testSecret, "user-123", "test@example.com", "tenant-456", time.Now().Add(-time.Hour))
				return token
			},
			expectError: true,
			errorMsg:    "token has expired",
		},
		{
			name: "empty token",
			setupToken: func() string {
				return ""
			},
			expectError: true,
			errorMsg:    "token is empty",
		},
		{
			name: "malformed token",
			setupToken: func() string {
				return "not-a-valid-token"
			},
			expectError: true,
			errorMsg:    "token is malformed",
		},
		{
			name: "invalid signature",
			setupToken: func() string {
				token, _ := CreateTestToken("wrong-secret", "user-123", "test@example.com", "tenant-456", time.Now().Add(time.Hour))
				return token
			},
			expectError: true,
			errorMsg:    "token signature is invalid",
		},
		{
			name: "token not valid yet",
			setupToken: func() string {
				claims := &PlatformClaims{
					UserID:   "user-123",
					Email:    "test@example.com",
					TenantID: "tenant-456",
					RegisteredClaims: jwt.RegisteredClaims{
						NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
					},
				}
				token, _ := GenerateToken(testSecret, claims)
				return token
			},
			expectError: true,
			errorMsg:    "token is not valid yet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()
			claims, err := validator.ValidateToken(token)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, "user-123", claims.UserID)
				assert.Equal(t, "test@example.com", claims.Email)
				assert.Equal(t, "tenant-456", claims.TenantID)
			}
		})
	}
}

func TestValidator_ValidateToken_MissingUserID(t *testing.T) {
	validator := NewValidator(testSecret)

	// Create token without user_id
	claims := &PlatformClaims{
		Email:    "test@example.com",
		TenantID: "tenant-456",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := GenerateToken(testSecret, claims)
	require.NoError(t, err)

	result, err := validator.ValidateToken(token)
	assert.Error(t, err)
	assert.Equal(t, "user_id claim is required", err.Error())
	assert.Nil(t, result)
}

func TestIsTokenExpired(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "expired error from jwt library",
			err:      jwt.ErrTokenExpired,
			expected: true,
		},
		{
			name:     "custom expired error message",
			err:      assert.AnError,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: IsTokenExpired checks the error message for "token has expired"
			result := IsTokenExpired(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsTokenExpired_CustomError(t *testing.T) {
	// Test with actual expired token
	validator := NewValidator(testSecret)
	token, _ := CreateTestToken(testSecret, "user-123", "test@example.com", "tenant-456", time.Now().Add(-time.Hour))
	_, err := validator.ValidateToken(token)
	assert.True(t, IsTokenExpired(err))
}

func TestNewValidator(t *testing.T) {
	validator := NewValidator("my-secret")
	assert.NotNil(t, validator)
}

func TestGenerateToken(t *testing.T) {
	claims := &PlatformClaims{
		UserID:   "user-123",
		Email:    "test@example.com",
		TenantID: "tenant-456",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token, err := GenerateToken(testSecret, claims)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify the token can be parsed
	validator := NewValidator(testSecret)
	parsedClaims, err := validator.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, claims.UserID, parsedClaims.UserID)
}

func TestCreateTestToken(t *testing.T) {
	token, err := CreateTestToken(testSecret, "user-123", "test@example.com", "tenant-456", time.Now().Add(time.Hour))
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify the token
	validator := NewValidator(testSecret)
	claims, err := validator.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "tenant-456", claims.TenantID)
}
