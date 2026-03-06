package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gowa/saas-openclaw/backend/pkg/jwt"
)

// Admin context keys
const (
	AdminContextKey = "admin"
	AdminIDKey      = "adminID"
	AdminUsernameKey = "adminUsername"
	AdminRoleKey    = "adminRole"
)

// AdminContext represents the authenticated admin information
type AdminContext struct {
	ID       string
	Username string
	Role     string
}

// AdminAuth creates a middleware that validates admin JWT tokens
func AdminAuth(validator *jwt.AdminValidator, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from request
		token := extractAdminToken(c)
		if token == "" {
			logger.Debug("Missing admin authentication token")
			respondWithError(c, http.StatusUnauthorized, "MISSING_TOKEN", "Authentication token is required")
			c.Abort()
			return
		}

		// Validate token
		claims, err := validator.ValidateAdminToken(token)
		if err != nil {
			logger.Debug("Invalid admin token", zap.Error(err))
			respondWithError(c, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
			c.Abort()
			return
		}

		// Inject admin context
		adminCtx := &AdminContext{
			ID:       claims.AdminID,
			Username: claims.Username,
			Role:     claims.Role,
		}

		c.Set(AdminContextKey, adminCtx)
		c.Set(AdminIDKey, claims.AdminID)
		c.Set(AdminUsernameKey, claims.Username)
		c.Set(AdminRoleKey, claims.Role)

		c.Next()
	}
}

// extractAdminToken extracts the JWT token from the request
func extractAdminToken(c *gin.Context) string {
	// Try Authorization: Bearer header first
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// Try X-Admin-Token header as fallback
	if token := c.GetHeader("X-Admin-Token"); token != "" {
		return token
	}

	return ""
}

// GetAdminContext retrieves the current admin from the context
func GetAdminContext(c *gin.Context) *AdminContext {
	if admin, exists := c.Get(AdminContextKey); exists {
		if adminCtx, ok := admin.(*AdminContext); ok {
			return adminCtx
		}
	}
	return nil
}

// GetAdminID retrieves the current admin ID from the context
func GetAdminID(c *gin.Context) string {
	if adminID, exists := c.Get(AdminIDKey); exists {
		if id, ok := adminID.(string); ok {
			return id
		}
	}
	return ""
}

// GetAdminRole retrieves the current admin role from the context
func GetAdminRole(c *gin.Context) string {
	if role, exists := c.Get(AdminRoleKey); exists {
		if r, ok := role.(string); ok {
			return r
		}
	}
	return ""
}

// RequireAdminRole creates a middleware that checks if the admin has the required role
func RequireAdminRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		admin := GetAdminContext(c)
		if admin == nil {
			respondWithError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Admin not authenticated")
			c.Abort()
			return
		}

		for _, role := range roles {
			if admin.Role == role {
				c.Next()
				return
			}
		}

		respondWithError(c, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
		c.Abort()
	}
}

// RequireSuperAdmin is a convenience middleware for super_admin only routes
func RequireSuperAdmin() gin.HandlerFunc {
	return RequireAdminRole("super_admin")
}
