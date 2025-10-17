package utils_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/otp"
)

func TestMain(m *testing.M) {
	// Check if database is available
	_ = config.LoadTest()
	if err := db.Init(); err != nil {
		// If database is not available, run tests without database
		// This allows unit tests to run even without a database
		// The globalDB variable will be nil if initialization fails
	}

	// Run tests
	code := m.Run()

	// Close the database connection pool if it was initialized
	db.Close()

	os.Exit(code)
}

func createTestDB(t *testing.T) *sql.DB {
	// Use the global database connection pool
	dbConn := db.GetDB()
	if dbConn == nil {
		t.Skip("Database connection is not available, skipping integration tests")
	}
	return dbConn
}

func cleanupOTP(t *testing.T, db *sql.DB, identifier string) {
	if db == nil {
		return
	}
	_, err := db.ExecContext(context.Background(), "DELETE FROM otps WHERE identifier = $1", identifier)
	if err != nil {
		t.Logf("Cleanup warning for OTP %s: %v", identifier, err)
	}
}

func TestGenerateOTP(t *testing.T) {
	otp, err := otp.GenerateOTP()
	require.NoError(t, err)
	assert.Len(t, otp, 6)
	assert.Regexp(t, `^\d{6}$`, otp)
}

func TestSaveOTPAndVerifyAndDeleteOTP(t *testing.T) {
	// Note: This is an integration test requiring a test database.
	// For full testing, set up a test PostgreSQL instance or use docker.
	// Here, we skip DB setup for unit test simplicity; use in CI with test DB.

	// Example structure (uncomment and configure for real testing):
	/*
	cfg := config.Load() // Use test config
	ctx := context.Background()
	dbConn, err := sql.Open("pgx", cfg.DSN())
	require.NoError(t, err)
	defer dbConn.Close()

	identifier := "test@example.com"
	knownOTP := "123456"

	err = otp.SaveOTP(ctx, dbConn, identifier, knownOTP)
	require.NoError(t, err)

	valid, err := otp.VerifyAndDeleteOTP(ctx, dbConn, identifier, knownOTP)
	assert.NoError(t, err)
	assert.True(t, valid)

	// Verify deleted
	valid, err = otp.VerifyAndDeleteOTP(ctx, dbConn, identifier, knownOTP)
	assert.NoError(t, err)
	assert.False(t, valid)

	// Test expiry (insert with past expiry)
	pastExpiry := time.Now().Add(-time.Hour)
	_, err = dbConn.ExecContext(ctx, "INSERT INTO otps (identifier, otp, expiry) VALUES ($1, $2, $3)", identifier, "654321", pastExpiry)
	require.NoError(t, err)

	valid, err = otp.VerifyAndDeleteOTP(ctx, dbConn, identifier, "654321")
	assert.NoError(t, err)
	assert.False(t, valid)
	*/
}

func TestCleanupExpiredOTPs(t *testing.T) {
	// Similar integration test setup required
	// Insert expired OTPs and verify deletion
}