package handlers_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/db"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/testutils"
)

// TestCreateUser tests user creation
func TestCreateUser(t *testing.T) {
	// Initialize test database configuration
	_ = config.LoadTest()
	
	// Initialize database connection pool for testing
	if err := db.Init(); err != nil {
		// If database initialization fails, skip the test
		t.Skipf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	
	db := testutils.SetupTestDB(t)
	defer db.Close()
	// Use unique email for each test run to avoid duplicate key conflicts
	uniqueEmail := fmt.Sprintf("test%d@example.com", time.Now().UnixNano())

	user := &models.User{
		Name:         "Test User",
		Email:        sql.NullString{String: uniqueEmail, Valid: true},
		Phone:        fmt.Sprintf("123456789%d", time.Now().UnixNano()%100),
		Role:         "student",
		IsMFAEnabled: false,
	}

	err := models.CreateUser(db, user)
	require.NoError(t, err)

	// Verify user was created
	var retrievedUser models.User
	err = db.QueryRow("SELECT id, name, email, phone, role, is_mfa_enabled FROM users WHERE id = $1", user.ID).Scan(
		&retrievedUser.ID, &retrievedUser.Name, &retrievedUser.Email, &retrievedUser.Phone, &retrievedUser.Role, &retrievedUser.IsMFAEnabled)
	require.NoError(t, err)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, user.Name, retrievedUser.Name)
	assert.Equal(t, user.Email.String, retrievedUser.Email.String)
	assert.Equal(t, user.Phone, retrievedUser.Phone)
	assert.Equal(t, user.Role, retrievedUser.Role)
	assert.Equal(t, user.IsMFAEnabled, retrievedUser.IsMFAEnabled)
}

// TestGetUser tests user retrieval
func TestGetUser(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// Use unique email for each test run to avoid duplicate key conflicts
	uniqueEmail := fmt.Sprintf("getuser%d@example.com", time.Now().UnixNano())

	// Create a test user first
	user := &models.User{
		Name:         "Test User",
		Email:        sql.NullString{String: uniqueEmail, Valid: true},
		Phone:        fmt.Sprintf("123456789%d", (time.Now().UnixNano()+1)%100),
		Role:         "student",
		IsMFAEnabled: false,
	}

	err := models.CreateUser(db, user)
	require.NoError(t, err)

	// Test FindUserByID function
	retrievedUser, err := models.FindUserByID(db, user.ID)
	require.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, user.Name, retrievedUser.Name)
	assert.Equal(t, user.Email.String, retrievedUser.Email.String)
	assert.Equal(t, user.Phone, retrievedUser.Phone)
	assert.Equal(t, user.Role, retrievedUser.Role)
	assert.Equal(t, user.IsMFAEnabled, retrievedUser.IsMFAEnabled)
}

// TestGetUserByEmail tests user retrieval by email
func TestGetUserByEmail(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// Use unique email for each test run to avoid duplicate key conflicts
	uniqueEmail := fmt.Sprintf("getemail%d@example.com", time.Now().UnixNano())

	// Create a test user first
	user := &models.User{
		Name:         "Test User",
		Email:        sql.NullString{String: uniqueEmail, Valid: true},
		Phone:        fmt.Sprintf("123456789%d", (time.Now().UnixNano()+2)%100),
		Role:         "student",
		IsMFAEnabled: false,
	}

	err := models.CreateUser(db, user)
	require.NoError(t, err)

	// Test FindUserByEmail function
	retrievedUser, err := models.FindUserByEmail(db, uniqueEmail)
	require.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, user.Name, retrievedUser.Name)
	assert.Equal(t, user.Email.String, retrievedUser.Email.String)
	assert.Equal(t, user.Phone, retrievedUser.Phone)
	assert.Equal(t, user.Role, retrievedUser.Role)
	assert.Equal(t, user.IsMFAEnabled, retrievedUser.IsMFAEnabled)
}

// TestUpdateUser tests user updates
func TestUpdateUser(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// Use unique email for each test run to avoid duplicate key conflicts
	uniqueEmail := fmt.Sprintf("updateuser%d@example.com", time.Now().UnixNano())

	// Create a test user first
	user := &models.User{
		Name:         "Test User",
		Email:        sql.NullString{String: uniqueEmail, Valid: true},
		Phone:        fmt.Sprintf("123456789%d", (time.Now().UnixNano()+3)%100),
		Role:         "student",
		IsMFAEnabled: false,
	}

	err := models.CreateUser(db, user)
	require.NoError(t, err)

	// Update user
	user.Name = "Updated Test User"
	user.IsMFAEnabled = true

	err = models.UpdateUser(db, user)
	require.NoError(t, err)

	// Verify user was updated
	var updatedUser models.User
	err = db.QueryRow("SELECT id, name, email, phone, role, is_mfa_enabled FROM users WHERE id = $1", user.ID).Scan(
		&updatedUser.ID, &updatedUser.Name, &updatedUser.Email, &updatedUser.Phone, &updatedUser.Role, &updatedUser.IsMFAEnabled)
	require.NoError(t, err)
	assert.Equal(t, "Updated Test User", updatedUser.Name)
	assert.Equal(t, true, updatedUser.IsMFAEnabled)
}

// TestDeleteUser tests user deletion
func TestDeleteUser(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	// Use unique email for each test run to avoid duplicate key conflicts
	uniqueEmail := fmt.Sprintf("deleteuser%d@example.com", time.Now().UnixNano())

	// Create a test user first
	user := &models.User{
		Name:         "Test User",
		Email:        sql.NullString{String: uniqueEmail, Valid: true},
		Phone:        fmt.Sprintf("123456789%d", (time.Now().UnixNano()+4)%100),
		Role:         "student",
		IsMFAEnabled: false,
	}

	err := models.CreateUser(db, user)
	require.NoError(t, err)

	// Delete user (using direct SQL since there's no DeleteUser function in models)
	_, err = db.Exec("DELETE FROM users WHERE id = $1", user.ID)
	require.NoError(t, err)

	// Verify user was deleted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", user.ID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

// MFA tests are commented out due to import issues
// TODO: Fix import issues and re-enable MFA tests