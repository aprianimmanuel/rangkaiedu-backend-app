package testutils

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/password"
)

// SetupTestDB sets up a test database connection
func SetupTestDB(t *testing.T) *sql.DB {
	// In a real test environment, you would use a test database
	// For now, we'll use the existing database connection
	db, err := sql.Open("pgx", "postgres://rangkaiedudev1:12d1q23wxm19wkc1fsdcq23@db:5432/rangkaiedu_test?sslmode=disable")
	require.NoError(t, err)
	return db
}

// CreateTestUser creates a test user for testing
func CreateTestUser(t *testing.T, db *sql.DB) *models.User {
	userID := uuid.New()
	user := &models.User{
		ID:           userID,
		Name:         "Test User",
		Email:        sql.NullString{String: fmt.Sprintf("test%d@example.com", time.Now().UnixNano()), Valid: true},
		Phone:        fmt.Sprintf("123456789%d", time.Now().UnixNano()%100),
		Role:         "student",
		IsMFAEnabled: false,
	}

	err := models.CreateUser(db, user)
	require.NoError(t, err)
	return user
}

// CreateSocialUser creates a social auth user for testing
func CreateSocialUser(t *testing.T, db *sql.DB, email, googleID, facebookID string) {
	_, err := db.Exec("INSERT INTO users (name, email, phone, role, password_hash, google_id, facebook_id) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		"Social User", email, fmt.Sprintf("123456789%d", time.Now().UnixNano()%100), "student", "", googleID, facebookID)
	require.NoError(t, err)
}

// CreateNormalUser creates a normal user with password for testing
func CreateNormalUser(t *testing.T, db *sql.DB, email, pass string) {
	hash, err := password.HashPassword(pass)
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO users (name, email, phone, role, password_hash, google_id, facebook_id) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		"Normal User", email, fmt.Sprintf("123456789%d", (time.Now().UnixNano()+1)%100), "student", hash, "", "")
	require.NoError(t, err)
}

// CleanupUser removes a user by email
func CleanupUser(t *testing.T, db *sql.DB, email string) {
	_, err := db.Exec("DELETE FROM users WHERE email = $1", email)
	if err != nil {
		t.Logf("Cleanup warning for %s: %v", email, err)
	}
}

// IsSocialAuthUser checks if a user is a social auth user
func IsSocialAuthUser(t *testing.T, db *sql.DB, email string) bool {
	var hasGoogleID, hasFacebookID sql.NullString
	err := db.QueryRow("SELECT google_id, facebook_id FROM users WHERE email = $1", email).Scan(&hasGoogleID, &hasFacebookID)
	if err != nil {
		t.Fatalf("Failed to query user: %v", err)
	}

	return (hasGoogleID.Valid && hasGoogleID.String != "") || (hasFacebookID.Valid && hasFacebookID.String != "")
}