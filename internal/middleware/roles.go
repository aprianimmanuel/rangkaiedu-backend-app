package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/monitoring"
)

// Role represents a user role in the system
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleTeacher  Role = "teacher"
	RoleStudent  Role = "student"
)

// RoleRequired is a middleware that validates JWT tokens and checks if user has required role
func RoleRequired(requiredRole Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First, ensure AuthRequired middleware has been run
		userInterface, exists := c.Get("user")
		if !exists {
			// Log security event for unauthorized access attempt
			monitoring.LogAuthViolation(c, "", "", c.Request.URL.Path, "UNAUTHORIZED_ACCESS", map[string]interface{}{
				"ip": c.ClientIP(),
				"required_role": requiredRole,
				"error": "No user context",
			})
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Type assert to UserContext
		user, ok := userInterface.(UserContext)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
			c.Abort()
			return
		}

		// Check if user has the required role
		if Role(user.Role) != requiredRole {
			// Log security event for insufficient permissions
			monitoring.LogAuthViolation(c, user.ID, user.Email, c.Request.URL.Path, "INSUFFICIENT_PERMISSIONS", map[string]interface{}{
				"ip": c.ClientIP(),
				"user_role": user.Role,
				"required_role": requiredRole,
			})
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// Continue with the next handler
		c.Next()
	}
}

// RolesRequired is a middleware that validates JWT tokens and checks if user has any of the required roles
func RolesRequired(allowedRoles ...Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First, ensure AuthRequired middleware has been run
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Type assert to UserContext
		user, ok := userInterface.(UserContext)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
			c.Abort()
			return
		}

		// Check if user has any of the allowed roles
		userRole := Role(user.Role)
		allowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// Continue with the next handler
		c.Next()
	}
}