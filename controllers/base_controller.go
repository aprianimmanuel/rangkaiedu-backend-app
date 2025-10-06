package controllers

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/aprianimmanuel/backend-app/middleware"
	"github.com/aprianimmanuel/backend-app/pkg/db"
)

// GetDBConnection returns a database connection pool
func GetDBConnection() (*pgxpool.Pool, error) {
	// Use the global database connection pool
	pool := db.GetDB()
	if pool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}
	return pool, nil
}

// GetCurrentUser returns the current authenticated user from the context
func GetCurrentUser(c *gin.Context) (*middleware.UserContext, error) {
	userInterface, exists := c.Get("user")
	if !exists {
		return nil, fmt.Errorf("user not authenticated")
	}

	user, ok := userInterface.(middleware.UserContext)
	if !ok {
		return nil, fmt.Errorf("invalid user context")
	}

	return &user, nil
}

// SendErrorResponse sends a standardized error response
func SendErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message, "code": code})
}

// SendSuccessResponse sends a standardized success response
func SendSuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}