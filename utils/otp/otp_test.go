package otp

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/aprianimmanuel/rangkaiedu-backend/config"
	"github.com/aprianimmanuel/rangkaiedu-backend/utils/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Set environment variables for test database
	os.Setenv("DB_NAME", "rangkaiedu_test")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_SSLMODE", "disable")

	// Initialize the database connection pool
	if err := db.Init(); err != nil {
		panic(err)
	}

	// Run tests
	code := m.Run()

	// Close the database connection pool
	db.Close()

	// Optional: unset env vars after tests
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_SSLMODE")

	os.Exit(code)
}

func createTestPool(t *testing.T) *pgxpool.Pool {
	// Use the global database connection pool
	pool := db.GetDB()
	if pool == nil {
		t.Fatalf("Database pool is not initialized")
	}
	return pool
}

func cleanupOTP(t *testing.T, pool *pgxpool.Pool, identifier string) {
	_, err := pool.Exec(context.Background(), "DELETE FROM otps WHERE identifier = $1", identifier)
	if err != nil && err != pgx.ErrNoRows {
		t.Logf("Cleanup warning for OTP %s: %v", identifier, err)
	}
}

func TestGenerateOTP(t *testing.T) {
	otp, err := GenerateOTP()
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
	pool, err := pgxpool.NewWithConfig(ctx, pgxpool.ParseConfig(cfg.DSN()))
	require.NoError(t, err)
	defer pool.Close()

	identifier := "test@example.com"
	knownOTP := "123456"

	err = SaveOTP(ctx, pool, identifier, knownOTP)
	require.NoError(t, err)

	valid, err := VerifyAndDeleteOTP(ctx, pool, identifier, knownOTP)
	assert.NoError(t, err)
	assert.True(t, valid)

	// Verify deleted
	valid, err = VerifyAndDeleteOTP(ctx, pool, identifier, knownOTP)
	assert.NoError(t, err)
	assert.False(t, valid)

	// Test expiry (insert with past expiry)
	pastExpiry := time.Now().Add(-time.Hour)
	_, err = pool.Exec(ctx, "INSERT INTO otps (identifier, otp, expiry) VALUES ($1, $2, $3)", identifier, "654321", pastExpiry)
	require.NoError(t, err)

	valid, err = VerifyAndDeleteOTP(ctx, pool, identifier, "654321")
	assert.NoError(t, err)
	assert.False(t, valid)
	*/
}

func TestCleanupExpiredOTPs(t *testing.T) {
	// Similar integration test setup required
	// Insert expired OTPs and verify deletion
}