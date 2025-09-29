package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Set some test environment variables
	os.Setenv("DB_HOST", "test_host")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USER", "test_user")
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("DB_SSLMODE", "require")

	// Load config
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check values
	if config.Host != "test_host" {
		t.Errorf("Expected host 'test_host', got '%s'", config.Host)
	}
	if config.Port != "5433" {
		t.Errorf("Expected port '5433', got '%s'", config.Port)
	}
	if config.Database != "test_db" {
		t.Errorf("Expected database 'test_db', got '%s'", config.Database)
	}
	if config.Username != "test_user" {
		t.Errorf("Expected username 'test_user', got '%s'", config.Username)
	}
	if config.Password != "test_password" {
		t.Errorf("Expected password 'test_password', got '%s'", config.Password)
	}
	if config.SSLMode != "require" {
		t.Errorf("Expected sslmode 'require', got '%s'", config.SSLMode)
	}

	// Test DSN building
	expectedDSN := "postgres://test_user:test_password@test_host:5433/test_db?sslmode=require"
	actualDSN := config.BuildDSN()
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
}

func TestLoadConfigWithDefaults(t *testing.T) {
	// Set only required environment variables
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USER", "test_user")

	// Load config
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check default values
	if config.Host != "localhost" {
		t.Errorf("Expected default host 'localhost', got '%s'", config.Host)
	}
	if config.Port != "5432" {
		t.Errorf("Expected default port '5432', got '%s'", config.Port)
	}
	if config.Database != "test_db" {
		t.Errorf("Expected database 'test_db', got '%s'", config.Database)
	}
	if config.Username != "test_user" {
		t.Errorf("Expected username 'test_user', got '%s'", config.Username)
	}
	if config.SSLMode != "disable" {
		t.Errorf("Expected default sslmode 'disable', got '%s'", config.SSLMode)
	}

	// Clean up environment variables
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_USER")
}

func TestLoadConfigMissingRequiredFields(t *testing.T) {
	// Don't set required environment variables
	_, err := LoadConfig()
	if err == nil {
		t.Error("Expected error for missing database name, but got none")
	}

	// Set database name but not username
	os.Setenv("DB_NAME", "test_db")
	_, err = LoadConfig()
	if err == nil {
		t.Error("Expected error for missing username, but got none")
	}

	// Clean up environment variables
	os.Unsetenv("DB_NAME")
}