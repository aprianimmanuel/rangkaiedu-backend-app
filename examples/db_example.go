// Package main demonstrates how to use the database package
// This is an example file and should not be used in production
package main

import (
	"context"
	"log"
	"os"

	"github.com/aprianimmanuel/backend-app/pkg/db"
)

func main() {
	// Set required environment variables for testing
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")

	// Connect to the database
	database, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Get the connection pool
	pool := database.GetPool()

	// Test the connection with a simple query
	var version string
	err = pool.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		log.Println("Query failed:", err)
		return
	}

	log.Println("Database connection successful!")
	log.Println("PostgreSQL version:", version)
}