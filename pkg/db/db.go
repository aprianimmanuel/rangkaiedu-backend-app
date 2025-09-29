package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Config holds database configuration.
type Config struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

// Load loads the configuration from environment variables.
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		Name:     getEnv("DB_NAME", "rangkaiedu_dev"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"), // Change in production
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Validate required fields
	if cfg.Host == "" || cfg.Port == "" || cfg.Name == "" || cfg.User == "" {
		log.Fatal("Missing required database configuration. Please set DB_HOST, DB_PORT, DB_NAME, DB_USER, and DB_PASSWORD in .env")
	}

	return cfg
}

// getEnv returns the value of the environment variable.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// DSN returns the PostgreSQL DSN string.
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.Name, c.User, c.Password, c.SSLMode)
}

// Pool is the global connection pool.
var Pool *pgxpool.Pool

// Init initializes the global connection pool.
func Init() error {
	cfg := Load()
	config, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return fmt.Errorf("invalid DSN: %w", err)
	}

	// Configure pool settings
	config.MaxConns = 20
	config.MinConns = 4
	config.MaxConnLifetime = 5 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the pool with a ping
	if err := Pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection pool initialized successfully")
	return nil
}

// Close closes the connection pool.
func Close() error {
	if Pool != nil {
		Pool.Close()
	}
	return nil
}

// GetDB returns the global pool.
func GetDB() *pgxpool.Pool {
	return Pool
}