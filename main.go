package main

import (
	"log"

	"github.com/aprianimmanuel/rangkaiedu-backend/internal/app"
)

func main() {
	// Create new application instance
	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Initialize application components
	if err := application.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Setup routes
	application.SetupRoutes()

	// Run the application
	if err := application.Run(); err != nil {
		log.Fatalf("Application failed to run: %v", err)
	}
}