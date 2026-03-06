package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// AdminClaims represents the JWT claims for admin authentication
type AdminClaims struct {
	AdminID  string `json:"admin_id"`
	Username string `json:"username"`
	Role     string `json:"role"` // "admin" or "super_admin"
	jwt.RegisteredClaims
}

// TokenIssuer is the issuer identifier for admin tokens
const TokenIssuer = "saas-openclaw-admin"

// AdminValidator handles admin JWT token validation
type AdminValidator struct {
	secret string
}

// NewAdminValidator creates a new admin JWT validator
func NewAdminValidator(secret string) *AdminValidator {
	return &AdminValidator{secret: secret}
}

// ValidateAdminToken validates an admin JWT token string and returns the claims
func (v *AdminValidator) ValidateAdminToken(tokenString string) (*AdminClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(v.secret), nil
	})

	if err != nil {
		return nil, mapJWTError(err)
	}

	if claims, ok := token.Claims.(*AdminClaims); ok && token.Valid {
		// Additional validation for required fields
		if claims.AdminID == "" {
			return nil, errors.New("admin_id claim is required")
		}
		if claims.Username == "" {
			return nil, errors.New("username claim is required")
		}
		if claims.Role == "" {
			return nil, errors.New("role claim is required")
		}
		return claims, nil
	}

	return nil, errors.New("invalid admin token")
}

// GenerateAdminToken generates a JWT token for admin authentication
func GenerateAdminToken(secret string, claims *AdminClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
