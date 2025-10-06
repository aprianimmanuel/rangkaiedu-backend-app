package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/aprianimmanuel/rangkaiedu-backend/config"
	"github.com/aprianimmanuel/rangkaiedu-backend/utils/db"
	"github.com/aprianimmanuel/rangkaiedu-backend/routes"
)

func main() {
	// Initialize the database connection pool
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	// Welcome route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Rangkai Edu Backend API",
		})
	})

	// Setup auth routes
	routes.SetupAuthRoutes(r)

	// Setup class routes
	routes.SetupClassRoutes(r)

	// Setup subject routes
	routes.SetupSubjectRoutes(r)

	// Setup material routes
	routes.SetupMaterialRoutes(r)

	// Setup health check routes
	routes.SetupHealthRoutes(r)

	// Run the server
	r.Run(":8080")
}