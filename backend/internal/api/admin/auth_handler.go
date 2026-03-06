// Package admin provides API handlers for admin endpoints
package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gowa/saas-openclaw/backend/internal/service"
)

// AuthHandler handles admin authentication API endpoints
type AuthHandler struct {
	authService *service.AdminAuthService
	logger      *zap.Logger
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *service.AdminAuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login handles POST /v1/admin/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Debug("Invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"data": nil,
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Username and password are required",
			},
			"meta": nil,
		})
		return
	}

	// Call service
	resp, err := h.authService.Login(c.Request.Context(), &service.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{
				"data": nil,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid credentials",
				},
				"meta": nil,
			})
			return
		}

		h.logger.Error("Login failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": nil,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "An internal error occurred",
			},
			"meta": nil,
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"token":     resp.Token,
			"expiresAt": resp.ExpiresAt,
			"adminId":   resp.AdminID,
			"username":  resp.Username,
			"role":      resp.Role,
		},
		"error": nil,
		"meta":  nil,
	})
}

// RegisterRoutes registers admin authentication routes
func RegisterRoutes(router *gin.RouterGroup, handler *AuthHandler) {
	auth := router.Group("/admin/auth")
	{
		auth.POST("/login", handler.Login)
	}
}
