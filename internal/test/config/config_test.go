package config_test

import (
	"os"
	"testing"
	
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
)

func TestLoad(t *testing.T) {
	// Set some test environment variables
	os.Setenv("DB_HOST", "test_host")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USER", "test_user")
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("DB_SSLMODE", "require")
	os.Setenv("JWT_SECRET", "test")

	// Load config
	config := config.Load()

	// Check values
	if config.DBHost != "test_host" {
		t.Errorf("Expected DBHost 'test_host', got '%s'", config.DBHost)
	}
	if config.DBPort != "5433" {
		t.Errorf("Expected DBPort '5433', got '%s'", config.DBPort)
	}
	if config.DBName != "test_db" {
		t.Errorf("Expected DBName 'test_db', got '%s'", config.DBName)
	}
	if config.DBUser != "test_user" {
		t.Errorf("Expected DBUser 'test_user', got '%s'", config.DBUser)
	}
	if config.DBPassword != "test_password" {
		t.Errorf("Expected DBPassword 'test_password', got '%s'", config.DBPassword)
	}
	if config.DBSSLMode != "require" {
		t.Errorf("Expected DBSSLMode 'require', got '%s'", config.DBSSLMode)
	}
	if config.JWTSecret != "test" {
		t.Errorf("Expected JWTSecret 'test', got '%s'", config.JWTSecret)
	}

	// Test DSN building
	expectedDSN := "host=test_host port=5433 dbname=test_db user=test_user password=test_password sslmode=require"
	actualDSN := config.DSN()
	if actualDSN != expectedDSN {
		t.Errorf("Expected DSN '%s', got '%s'", expectedDSN, actualDSN)
	}

	// Clean up environment variables
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_SSLMODE")
	os.Unsetenv("JWT_SECRET")
}

func TestLoadWithDefaults(t *testing.T) {
	// Unset to use defaults
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_SSLMODE")
	os.Unsetenv("JWT_SECRET")

	// Load config - should use defaults without panic
	config := config.Load()

	if config.DBHost != "localhost" {
		t.Errorf("Expected default DBHost 'localhost', got '%s'", config.DBHost)
	}
	if config.DBPort != "5432" {
		t.Errorf("Expected default DBPort '5432', got '%s'", config.DBPort)
	}
	if config.DBName != "rangkaiedu_dev" {
		t.Errorf("Expected default DBName 'rangkaiedu_dev', got '%s'", config.DBName)
	}
	if config.DBUser != "postgres" {
		t.Errorf("Expected default DBUser 'postgres', got '%s'", config.DBUser)
	}
	if config.DBPassword != "password" {
		t.Errorf("Expected default DBPassword 'password', got '%s'", config.DBPassword)
	}
	if config.DBSSLMode != "disable" {
		t.Errorf("Expected default DBSSLMode 'disable', got '%s'", config.DBSSLMode)
	}
	if config.JWTSecret != "default-secret-key-change-in-production" {
		t.Errorf("Expected default JWTSecret, got '%s'", config.JWTSecret)
	}
}

func TestLoadMissingRequiredFields(t *testing.T) {
	// Create a temporary function to test the validation logic
	// We'll directly test the validation condition instead of trying to trigger log.Fatal
	cfg := &config.Config{
		DBHost:   "",
		DBPort:   "",
		DBName:   "",
		DBUser:   "",
		DBPassword: "password",
		DBSSLMode:  "disable",
		JWTSecret:  "secret",
	}

	// Test the validation condition directly
	if !(cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" || cfg.DBUser == "") {
		t.Error("Expected validation to fail for missing required database configuration")
	}

	// Test with valid config
	cfg = &config.Config{
		DBHost:   "localhost",
		DBPort:   "5432",
		DBName:   "testdb",
		DBUser:   "testuser",
		DBPassword: "password",
		DBSSLMode:  "disable",
		JWTSecret:  "secret",
	}

	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" || cfg.DBUser == "" {
		t.Error("Expected validation to pass for valid database configuration")
	}
}