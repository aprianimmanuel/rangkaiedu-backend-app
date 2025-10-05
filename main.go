package main

import (
	"github.com/gin-gonic/gin"
	"github.com/aprianimmanuel/backend-app/routes"
)

func main() {
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

	// Run the server
	r.Run(":8080")
}