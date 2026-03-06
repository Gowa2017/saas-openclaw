package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// PlatformClaims represents the JWT claims for platform authentication
type PlatformClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	TenantID string `json:"tenant_id"`
	jwt.RegisteredClaims
}

// Validator handles JWT token validation
type Validator struct {
	secret string
}

// NewValidator creates a new JWT validator
func NewValidator(secret string) *Validator {
	return &Validator{secret: secret}
}

// ValidateToken validates a JWT token string and returns the claims
func (v *Validator) ValidateToken(tokenString string) (*PlatformClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &PlatformClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(v.secret), nil
	})

	if err != nil {
		return nil, mapJWTError(err)
	}

	if claims, ok := token.Claims.(*PlatformClaims); ok && token.Valid {
		// Additional validation for required fields
		if claims.UserID == "" {
			return nil, errors.New("user_id claim is required")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenerateToken generates a JWT token for testing purposes
func GenerateToken(secret string, claims *PlatformClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// mapJWTError maps jwt library errors to more descriptive errors
func mapJWTError(err error) error {
	if errors.Is(err, jwt.ErrTokenExpired) {
		return errors.New("token has expired")
	}
	if errors.Is(err, jwt.ErrTokenNotValidYet) {
		return errors.New("token is not valid yet")
	}
	if errors.Is(err, jwt.ErrTokenMalformed) {
		return errors.New("token is malformed")
	}
	if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return errors.New("token signature is invalid")
	}
	return errors.New("failed to parse token")
}

// IsTokenExpired checks if the error is due to token expiration
func IsTokenExpired(err error) bool {
	return errors.Is(err, jwt.ErrTokenExpired) || err.Error() == "token has expired"
}

// CreateTestToken creates a test token with the given claims for testing
func CreateTestToken(secret string, userID, email, tenantID string, expiresAt time.Time) (string, error) {
	claims := &PlatformClaims{
		UserID:   userID,
		Email:    email,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	return GenerateToken(secret, claims)
}
