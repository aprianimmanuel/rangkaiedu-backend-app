package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration.
type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string // "disable" for local dev, "require" for production
}

// Load loads the configuration from environment variables.
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := &Config{
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBName:    getEnv("DB_NAME", "rangkaiedu_dev"),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"), // Change this in production
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Validate required fields
	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" || cfg.DBUser == "" {
		log.Fatal("Missing required database configuration. Please set DB_HOST, DB_PORT, DB_NAME, DB_USER, and DB_PASSWORD in .env")
	}

	return cfg
}

// getEnv returns the value of the environment variable named by the key.
// If the variable is not found, it returns the provided default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// DSN returns the PostgreSQL DSN string.
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBName, c.DBUser, c.DBPassword, c.DBSSLMode)
}