
package repositories

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/repositories"
)

func TestUserRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	// Test successful user retrieval
	userID := uuid.New()
	expectedUser := &models.User{
		ID:              userID,
		Name:            "Test User",
		Email:           sql.NullString{String: "test@example.com", Valid: true},
		Phone:           "1234567890",
		Role:            "student",
		PasswordHash:    sql.NullString{String: "hashed_password", Valid: true},
		GoogleID:        sql.NullString{String: "google123", Valid: true},
		FacebookID:      sql.NullString{String: "facebook123", Valid: true},
		MFASecret:       sql.NullString{String: "mfa_secret", Valid: true},
		IsMFAEnabled:    true,
		MFABackupCodes:  sql.NullString{String: "[\"code1\",\"code2\"]", Valid: true},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		LastLogin:       sql.NullTime{Time: time.Now(), Valid: true},
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "phone", "role", "password_hash", "google_id", "facebook_id",
		"mfa_secret", "is_mfa_enabled", "mfa_backup_codes", "created_at", "updated_at", "last_login",
	}).AddRow(
		expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Phone, expectedUser.Role,
		expectedUser.PasswordHash, expectedUser.GoogleID, expectedUser.FacebookID,
		expectedUser.MFASecret, expectedUser.IsMFAEnabled, "{\"code1\",\"code2\"}",
		expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.LastLogin,
	)

	mock.ExpectQuery("SELECT id, name, email, phone, role, password_hash, google_id, facebook_id, mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := repo.FindByID(userID)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Phone, user.Phone)
	assert.Equal(t, expectedUser.Role, user.Role)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	userID := uuid.New()

	mock.ExpectQuery("SELECT id, name, email, phone, role, password_hash, google_id, facebook_id, mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.FindByID(userID)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrNoRows, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	// Test successful user retrieval by email
	expectedUser := &models.User{
		ID:              uuid.New(),
		Name:            "Test User",
		Email:           sql.NullString{String: "test@example.com", Valid: true},
		Phone:           "1234567890",
		Role:            "student",
		PasswordHash:    sql.NullString{String: "hashed_password", Valid: true},
		GoogleID:        sql.NullString{String: "google123", Valid: true},
		FacebookID:      sql.NullString{String: "facebook123", Valid: true},
		MFASecret:       sql.NullString{String: "mfa_secret", Valid: true},
		IsMFAEnabled:    true,
		MFABackupCodes:  sql.NullString{String: "[\"code1\",\"code2\"]", Valid: true},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		LastLogin:       sql.NullTime{Time: time.Now(), Valid: true},
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "phone", "role", "password_hash", "google_id", "facebook_id",
		"mfa_secret", "is_mfa_enabled", "mfa_backup_codes", "created_at", "updated_at", "last_login",
	}).AddRow(
		expectedUser.ID, expectedUser.Name, expectedUser.Email.String, expectedUser.Phone, expectedUser.Role,
		expectedUser.PasswordHash.String, expectedUser.GoogleID.String, expectedUser.FacebookID.String,
		expectedUser.MFASecret.String, expectedUser.IsMFAEnabled, "{\"code1\",\"code2\"}",
		expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.LastLogin.Time,
	)

	mock.ExpectQuery("SELECT id, name, email, phone, role, password_hash, google_id, facebook_id, mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login FROM users WHERE email = \\$1").
		WithArgs(expectedUser.Email.String).
		WillReturnRows(rows)

	user, err := repo.FindByEmail(expectedUser.Email.String)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	email := "nonexistent@example.com"

	mock.ExpectQuery("SELECT id, name, email, phone, role, password_hash, google_id, facebook_id, mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login FROM users WHERE email = \\$1").
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.FindByEmail(email)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrNoRows, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByPhone(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	// Test successful user retrieval by phone
	expectedUser := &models.User{
		ID:              uuid.New(),
		Name:            "Test User",
		Email:           sql.NullString{String: "test@example.com", Valid: true},
		Phone:           "1234567890",
		Role:            "student",
		PasswordHash:    sql.NullString{String: "hashed_password", Valid: true},
		GoogleID:        sql.NullString{String: "google123", Valid: true},
		FacebookID:      sql.NullString{String: "facebook123", Valid: true},
		MFASecret:       sql.NullString{String: "mfa_secret", Valid: true},
		IsMFAEnabled:    true,
		MFABackupCodes:  sql.NullString{String: "[\"code1\",\"code2\"]", Valid: true},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		LastLogin:       sql.NullTime{Time: time.Now(), Valid: true},
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "phone", "role", "password_hash", "google_id", "facebook_id",
		"mfa_secret", "is_mfa_enabled", "mfa_backup_codes", "created_at", "updated_at", "last_login",
	}).AddRow(
		expectedUser.ID, expectedUser.Name, expectedUser.Email.String, expectedUser.Phone, expectedUser.Role,
		expectedUser.PasswordHash.String, expectedUser.GoogleID.String, expectedUser.FacebookID.String,
		expectedUser.MFASecret.String, expectedUser.IsMFAEnabled, "{\"code1\",\"code2\"}",
		expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.LastLogin.Time,
	)

	mock.ExpectQuery("SELECT id, name, email, phone, role, password_hash, google_id, facebook_id, mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login FROM users WHERE phone = \\$1").
		WithArgs(expectedUser.Phone).
		WillReturnRows(rows)

	user, err := repo.FindByPhone(expectedUser.Phone)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Phone, user.Phone)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByGoogleID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	// Test successful user retrieval by Google ID
	googleID := "google123"
	expectedUser := &models.User{
		ID:              uuid.New(),
		Name:            "Test User",
		Email:           sql.NullString{String: "test@example.com", Valid: true},
		Phone:           "1234567890",
		Role:            "student",
		PasswordHash:    sql.NullString{String: "hashed_password", Valid: true},
		GoogleID:        sql.NullString{String: googleID, Valid: true},
		FacebookID:      sql.NullString{String: "", Valid: false},
		MFASecret:       sql.NullString{String: "mfa_secret", Valid: true},
		IsMFAEnabled:    true,
		MFABackupCodes:  sql.NullString{String: "[\"code1\",\"code2\"]", Valid: true},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		LastLogin:       sql.NullTime{Time: time.Now(), Valid: true},
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "phone", "role", "password_hash", "google_id", "facebook_id",
		"mfa_secret", "is_mfa_enabled", "mfa_backup_codes", "created_at", "updated_at", "last_login",
	}).AddRow(
		expectedUser.ID, expectedUser.Name, expectedUser.Email.String, expectedUser.Phone, expectedUser.Role,
		expectedUser.PasswordHash.String, expectedUser.GoogleID.String, expectedUser.FacebookID.String,
		expectedUser.MFASecret.String, expectedUser.IsMFAEnabled, "{\"code1\",\"code2\"}",
		expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.LastLogin.Time,
	)

	mock.ExpectQuery("SELECT id, name, email, phone, role, password_hash, google_id, facebook_id, mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login FROM users WHERE google_id = \\$1").
		WithArgs(googleID).
		WillReturnRows(rows)

	user, err := repo.FindByGoogleID(googleID)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, googleID, user.GoogleID.String)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByFacebookID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	// Test successful user retrieval by Facebook ID
	facebookID := "facebook123"
	expectedUser := &models.User{
		ID:              uuid.New(),
		Name:            "Test User",
		Email:           sql.NullString{String: "test@example.com", Valid: true},
		Phone:           "1234567890",
		Role:            "student",
		PasswordHash:    sql.NullString{String: "hashed_password", Valid: true},
		GoogleID:        sql.NullString{String: "", Valid: false},
		FacebookID:      sql.NullString{String: facebookID, Valid: true},
		MFASecret:       sql.NullString{String: "mfa_secret", Valid: true},
		IsMFAEnabled:    true,
		MFABackupCodes:  sql.NullString{String: "[\"code1\",\"code2\"]", Valid: true},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		LastLogin:       sql.NullTime{Time: time.Now(), Valid: true},
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "phone", "role", "password_hash", "google_id", "facebook_id",
		"mfa_secret", "is_mfa_enabled", "mfa_backup_codes", "created_at", "updated_at", "last_login",
	}).AddRow(
		expectedUser.ID, expectedUser.Name, expectedUser.Email.String, expectedUser.Phone, expectedUser.Role,
		expectedUser.PasswordHash.String, expectedUser.GoogleID.String, expectedUser.FacebookID.String,
		expectedUser.MFASecret.String, expectedUser.IsMFAEnabled, "{\"code1\",\"code2\"}",
		expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.LastLogin.Time,
	)

	mock.ExpectQuery("SELECT id, name, email, phone, role, password_hash, google_id, facebook_id, mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login FROM users WHERE facebook_id = \\$1").
		WithArgs(facebookID).
		WillReturnRows(rows)

	user, err := repo.FindByFacebookID(facebookID)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, facebookID, user.FacebookID.String)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	// Test successful user creation
	user := &models.User{
		Name:           "Test User",
		Email:          sql.NullString{String: "test@example.com", Valid: true},
		Phone:          "1234567890",
		Role:           "student",
		PasswordHash:   sql.NullString{String: "hashed_password", Valid: true},
		GoogleID:       sql.NullString{String: "google123", Valid: true},
		FacebookID:     sql.NullString{String: "facebook123", Valid: true},
		MFASecret:      sql.NullString{String: "mfa_secret", Valid: true},
		IsMFAEnabled:   true,
		MFABackupCodes:  sql.NullString{String: "[\"code1\",\"code2\"]", Valid: true},
	}

	userID := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(userID, createdAt, updatedAt)

	mock.ExpectQuery("INSERT INTO users \\(name, email, phone, role, password_hash, google_id, facebook_id, mfa_secret, is_mfa_enabled, mfa_backup_codes\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, \\$7, \\$8, \\$9, \\$10\\) RETURNING id, created_at, updated_at").
		WithArgs(user.Name, user.Email.String, user.Phone, user.Role, user.PasswordHash.String, user.GoogleID.String, user.FacebookID.String, user.MFASecret.String, user.IsMFAEnabled, "{\"code1\",\"code2\"}").
		WillReturnRows(rows)

	err = repo.Create(user)
	require.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, createdAt, user.CreatedAt)
	assert.Equal(t, updatedAt, user.UpdatedAt)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	// Test successful user update
	userID := uuid.New()
	user := &models.User{
		ID:             userID,
		Name:           "Updated Test User",
		Email:          sql.NullString{String: "updated@example.com", Valid: true},
		Phone:          "0987654321",
		Role:           "teacher",
		PasswordHash:   sql.NullString{String: "updated_hashed_password", Valid: true},
		GoogleID:       sql.NullString{String: "updated_google123", Valid: true},
		FacebookID:     sql.NullString{String: "updated_facebook123", Valid: true},
		MFASecret:      sql.NullString{String: "updated_mfa_secret", Valid: true},
		IsMFAEnabled:   false,
		MFABackupCodes: sql.NullString{String: "[\"updated_code1\",\"updated_code2\"]", Valid: true},
	}

	mock.ExpectExec("UPDATE users SET name = \\$1, email = \\$2, phone = \\$3, role = \\$4, password_hash = \\$5, google_id = \\$6, facebook_id = \\$7, mfa_secret = \\$8, is_mfa_enabled = \\$9, mfa_backup_codes = \\$10, updated_at = CURRENT_TIMESTAMP WHERE id = \\$11").
		WithArgs(user.Name, user.Email.String, user.Phone, user.Role, user.PasswordHash.String, user.GoogleID.String, user.FacebookID.String, user.MFASecret.String, user.IsMFAEnabled, "{\"updated_code1\",\"updated_code2\"}", user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(user)
	require.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_UpdateLastLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	// Test successful last login update
	userID := uuid.New()

	mock.ExpectExec("UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = \\$1").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateLastLogin(userID)
	require.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}