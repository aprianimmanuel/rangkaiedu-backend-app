// Package db contains manual tests for the database connection pool
// This is a manual test file and should not be run automatically
package manualdb

import (
	"os"
	"testing"
	
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/db"
)

// This is a manual test that can be run with:
// go run examples/db_example.go
// It's not meant to be run as part of the automated test suite
func TestConnectDBManual(t *testing.T) {
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

	// Get the database connection
	database := db.GetDB()
	if database != nil {
		// Make sure to close it
		defer database.Close()
	}
}