package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gowa/saas-openclaw/backend/internal/infrastructure/platform"
	"github.com/gowa/saas-openclaw/backend/pkg/jwt"
)

// Context keys for user information
const (
	UserContextKey  = "user"
	UserIDKey       = "userID"
	UserEmailKey    = "userEmail"
	UserTenantIDKey = "userTenantID"
	UserNameKey     = "userName"
	UserRoleKey     = "userRole"
)

// UserContext represents the authenticated user information
type UserContext struct {
	ID       string
	Name     string
	Email    string
	TenantID string
	Role     string
}

// PlatformAuth creates a middleware that validates JWT tokens and fetches user info from platform
func PlatformAuth(validator *jwt.Validator, platformClient *platform.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from request
		token := extractToken(c)
		if token == "" {
			respondWithError(c, http.StatusUnauthorized, "MISSING_TOKEN", "Authentication token is required")
			c.Abort()
			return
		}

		// Validate JWT token
		_, err := validator.ValidateToken(token)
		if err != nil {
			respondWithError(c, http.StatusUnauthorized, "INVALID_TOKEN", err.Error())
			c.Abort()
			return
		}

		// Fetch user info from platform
		userInfo, err := platformClient.GetUserInfo(c.Request.Context(), token)
		if err != nil {
			respondWithError(c, http.StatusUnauthorized, "USER_INFO_FAILED", "Failed to get user information")
			c.Abort()
			return
		}

		// Inject user context
		userCtx := &UserContext{
			ID:       userInfo.ID,
			Name:     userInfo.Name,
			Email:    userInfo.Email,
			TenantID: userInfo.TenantID,
			Role:     userInfo.Role,
		}

		c.Set(UserContextKey, userCtx)
		c.Set(UserIDKey, userInfo.ID)
		c.Set(UserEmailKey, userInfo.Email)
		c.Set(UserTenantIDKey, userInfo.TenantID)
		c.Set(UserNameKey, userInfo.Name)
		c.Set(UserRoleKey, userInfo.Role)

		c.Next()
	}
}

// JWTAuth creates a lightweight middleware that only validates JWT tokens
func JWTAuth(validator *jwt.Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from request
		token := extractToken(c)
		if token == "" {
			respondWithError(c, http.StatusUnauthorized, "MISSING_TOKEN", "Authentication token is required")
			c.Abort()
			return
		}

		// Validate JWT token
		claims, err := validator.ValidateToken(token)
		if err != nil {
			respondWithError(c, http.StatusUnauthorized, "INVALID_TOKEN", err.Error())
			c.Abort()
			return
		}

		// Inject user context from JWT claims
		userCtx := &UserContext{
			ID:       claims.UserID,
			Email:    claims.Email,
			TenantID: claims.TenantID,
		}

		c.Set(UserContextKey, userCtx)
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)
		c.Set(UserTenantIDKey, claims.TenantID)

		c.Next()
	}
}

// extractToken extracts the JWT token from the request
func extractToken(c *gin.Context) string {
	// Try X-Platform-Token header first (primary method)
	if token := c.GetHeader("X-Platform-Token"); token != "" {
		return token
	}

	// Try Authorization: Bearer header (fallback)
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	return ""
}

// respondWithError sends an error response in the standard API format
func respondWithError(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, gin.H{
		"data": nil,
		"error": gin.H{
			"code":    code,
			"message": message,
		},
		"meta": nil,
	})
}

// GetCurrentUser retrieves the current user from the context
func GetCurrentUser(c *gin.Context) *UserContext {
	if user, exists := c.Get(UserContextKey); exists {
		if userCtx, ok := user.(*UserContext); ok {
			return userCtx
		}
	}
	return nil
}

// GetCurrentUserID retrieves the current user ID from the context
func GetCurrentUserID(c *gin.Context) string {
	if userID, exists := c.Get(UserIDKey); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// GetCurrentUserTenantID retrieves the current user's tenant ID from the context
func GetCurrentUserTenantID(c *gin.Context) string {
	if tenantID, exists := c.Get(UserTenantIDKey); exists {
		if id, ok := tenantID.(string); ok {
			return id
		}
	}
	return ""
}

// RequireRole creates a middleware that checks if the user has the required role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetCurrentUser(c)
		if user == nil {
			respondWithError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			c.Abort()
			return
		}

		for _, role := range roles {
			if user.Role == role {
				c.Next()
				return
			}
		}

		respondWithError(c, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
		c.Abort()
	}
}

// RequireTenant creates a middleware that ensures the user belongs to a tenant
func RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetCurrentUser(c)
		if user == nil {
			respondWithError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			c.Abort()
			return
		}

		if user.TenantID == "" {
			respondWithError(c, http.StatusForbidden, "NO_TENANT", "User does not belong to any tenant")
			c.Abort()
			return
		}

		c.Next()
	}
}
