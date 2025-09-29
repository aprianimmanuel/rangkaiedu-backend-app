// Package main demonstrates how to use the configuration package
// This is an example file and should not be used in production
package main

import (
	"fmt"
	"log"

	"github.com/aprianimmanuel/backend-app/config"
)

func main() {
	// Load configuration
	dbConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	// Print configuration
	fmt.Printf("Database Config:\n")
	fmt.Printf("  Host: %s\n", dbConfig.Host)
	fmt.Printf("  Port: %s\n", dbConfig.Port)
	fmt.Printf("  Database: %s\n", dbConfig.Database)
	fmt.Printf("  Username: %s\n", dbConfig.Username)
	fmt.Printf("  SSL Mode: %s\n", dbConfig.SSLMode)

	// Build and print connection string
	dsn := dbConfig.BuildDSN()
	fmt.Printf("\nConnection String:\n%s\n", dsn)
}