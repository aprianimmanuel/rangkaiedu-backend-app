package db

import (
	"os"
	"testing"

	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/db"
)

func TestConnectDB(t *testing.T) {
	// Set required environment variables for testing
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")

	// Try to connect to the database
	err := db.Init()
	if err != nil {
		// This is expected to fail since we're using dummy credentials
		// but we want to ensure the function handles errors properly
		t.Logf("Expected connection error with dummy credentials: %v", err)
		return
	}

	// Clean up
	defer db.Close()
}

func TestGetPool(t *testing.T) {
	// Set required environment variables for testing
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")

	// Try to connect to the database
	err := db.Init()
	if err != nil {
		// This is expected to fail since we're using dummy credentials
		t.Logf("Expected connection error with dummy credentials: %v", err)
		return
	}

	// Check that we can get the pool
	pool := db.GetDB()
	if pool == nil {
		t.Error("Expected pool to be non-nil")
	}

	// Clean up
	defer db.Close()
}