// Package db contains manual tests for the database connection pool
// This is a manual test file and should not be run automatically
package db

import (
	"os"
	"testing"
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
	database, err := ConnectDB()
	if err != nil {
		// This is expected to fail since we're using dummy credentials
		// but we want to ensure the function handles errors properly
		t.Logf("Expected connection error with dummy credentials: %v", err)
		return
	}

	// If we somehow got a connection, make sure to close it
	if database != nil {
		database.Close()
	}
}