// Package service provides business logic services for the application
package service

import (
	"context"
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/gowa/saas-openclaw/backend/internal/domain/user"
	"github.com/gowa/saas-openclaw/backend/pkg/jwt"
)

// Common errors for admin authentication
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenGeneration    = errors.New("failed to generate token")
)

// AdminAuthRepository defines the interface for admin authentication data access
type AdminAuthRepository interface {
	GetByUsername(username string) (*user.AdminUser, error)
	VerifyPassword(admin *user.AdminUser, password string) bool
}

// LoginRequest represents the admin login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the admin login response
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
	AdminID   string `json:"adminId"`
	Username  string `json:"username"`
	Role      string `json:"role"`
}

// AdminAuthService handles admin authentication business logic
type AdminAuthService struct {
	adminRepo AdminAuthRepository
	jwtSecret string
	tokenExp  time.Duration
	logger    *zap.Logger
}

// NewAdminAuthService creates a new AdminAuthService
func NewAdminAuthService(
	adminRepo AdminAuthRepository,
	jwtSecret string,
	tokenExp time.Duration,
	logger *zap.Logger,
) *AdminAuthService {
	return &AdminAuthService{
		adminRepo: adminRepo,
		jwtSecret: jwtSecret,
		tokenExp:  tokenExp,
		logger:    logger,
	}
}

// Login authenticates an admin user and returns a JWT token
// The ctx parameter is reserved for future use (e.g., request timeout, tracing)
func (s *AdminAuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 1. Find user by username
	admin, err := s.adminRepo.GetByUsername(req.Username)
	if err != nil {
		// Log failure without exposing whether username exists
		s.logLoginFailure(req.Username, "user not found")
		return nil, ErrInvalidCredentials
	}

	// 2. Verify password
	if !s.adminRepo.VerifyPassword(admin, req.Password) {
		s.logLoginFailure(req.Username, "invalid password")
		return nil, ErrInvalidCredentials
	}

	// 3. Generate token
	expiresAt := time.Now().Add(s.tokenExp)
	claims := &jwt.AdminClaims{
		AdminID:  admin.ID,
		Username: admin.Username,
		Role:     string(admin.Role),
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(expiresAt),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			NotBefore: jwtlib.NewNumericDate(time.Now()),
			Issuer:    jwt.TokenIssuer,
		},
	}

	token, err := jwt.GenerateAdminToken(s.jwtSecret, claims)
	if err != nil {
		s.logger.Error("Failed to generate admin token",
			zap.String("adminId", admin.ID),
			zap.Error(err),
		)
		return nil, ErrTokenGeneration
	}

	// 4. Log successful login
	s.logger.Info("Admin login successful",
		zap.String("adminId", admin.ID),
		zap.String("username", admin.Username),
		zap.String("role", string(admin.Role)),
	)

	// 5. Return response
	return &LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt.Unix(),
		AdminID:   admin.ID,
		Username:  admin.Username,
		Role:      string(admin.Role),
	}, nil
}

// logLoginFailure logs a failed login attempt without exposing sensitive details
func (s *AdminAuthService) logLoginFailure(username, reason string) {
	s.logger.Warn("Admin login failed",
		zap.String("username", username),
		zap.String("reason", reason),
	)
}

// GenerateToken generates a JWT token for an admin user
func (s *AdminAuthService) GenerateToken(admin *user.AdminUser) (string, error) {
	expiresAt := time.Now().Add(s.tokenExp)
	claims := &jwt.AdminClaims{
		AdminID:  admin.ID,
		Username: admin.Username,
		Role:     string(admin.Role),
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(expiresAt),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			NotBefore: jwtlib.NewNumericDate(time.Now()),
			Issuer:    jwt.TokenIssuer,
		},
	}

	return jwt.GenerateAdminToken(s.jwtSecret, claims)
}
