package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/aprianimmanuel/rangkaiedu-backend/controllers"
)

// SetupHealthRoutes sets up all health check related routes
func SetupHealthRoutes(r *gin.Engine) {
	// Basic health check endpoint - fast response for load balancers
	r.GET("/health", controllers.HealthCheckBasic)
	
	// Standard health check endpoint - comprehensive health status
	r.GET("/health/check", controllers.HealthCheck)
	
	// Detailed health check endpoint - includes system information
	r.GET("/health/detailed", controllers.HealthCheckDetailed)
	
	// Database-specific health check endpoint
	r.GET("/health/database", controllers.HealthCheckDatabase)
	
	// Legacy health check endpoint for backward compatibility
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "pong",
		})
	})
}