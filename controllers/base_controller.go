package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/aprianimmanuel/backend-app/config"
	"github.com/aprianimmanuel/backend-app/middleware"
)

// GetDBConnection returns a database connection pool
func GetDBConnection() (*pgxpool.Pool, error) {
	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		log.Printf("Failed to parse database config: %v", err)
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}
	
	// Configure pool settings
	poolConfig.MaxConns = 20
	poolConfig.MinConns = 4
	poolConfig.MaxConnLifetime = 5 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
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