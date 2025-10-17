package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// User represents the user model with MFA and social provider support
type User struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	Email          sql.NullString `json:"email"`
	Phone          string         `json:"phone"`
	Role           string         `json:"role"`
	PasswordHash   sql.NullString `json:"password_hash"`
	GoogleID       sql.NullString `json:"google_id"`
	FacebookID      sql.NullString `json:"facebook_id"`
	MFASecret      sql.NullString `json:"mfa_secret"`
	IsMFAEnabled   bool           `json:"is_mfa_enabled"`
	MFABackupCodes sql.NullString  `json:"mfa_backup_codes"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	LastLogin      sql.NullTime   `json:"last_login"`
}

// UserUpsertRequest represents the request for upsert operations
type UserUpsertRequest struct {
	Provider       string   `json:"provider" binding:"required,oneof=google facebook whatsapp email"`
	Identifier     string   `json:"identifier" binding:"required"` // email, phone, or provider token
	Role           string   `json:"role" binding:"required,oneof=admin teacher student parent"`
	Name           string   `json:"name" binding:"required,max=255"`
	Password       string   `json:"password" binding:"omitempty,min=8"` // optional for email provider
	Otp            string   `json:"otp" binding:"omitempty,len=6"`      // optional for OTP verification
	GoogleID       string   `json:"google_id"`                         // optional for Google
	FacebookID     string   `json:"facebook_id"`                       // optional for Facebook
}

// UserResponse represents the response for user operations
type UserResponse struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	Phone          string   `json:"phone"`
	Role           string   `json:"role"`
	IsMFAEnabled   bool     `json:"is_mfa_enabled"`
}

// LoginSuccessResponse represents the successful login response
type LoginSuccessResponse struct {
	Status      string      `json:"status"`
	Token       string      `json:"token"`
	User        UserResponse `json:"user"`
	RequiresMFA bool        `json:"requires_mfa"`
}

// RegistrationSuccessResponse represents the successful registration response
type RegistrationSuccessResponse struct {
	Status   string        `json:"status"`
	User     UserResponse  `json:"user"`
	NextStep string        `json:"next_step"`
}

// ErrorResponse represents the error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// OTPRequest represents the OTP request
type OTPRequest struct {
	Provider   string `json:"provider" binding:"required,oneof=whatsapp email"`
	Identifier string `json:"identifier" binding:"required"` // phone number or email
}

// OTPVerificationRequest represents the OTP verification request
type OTPVerificationRequest struct {
	Provider   string `json:"provider" binding:"required,oneof=whatsapp email"`
	Identifier string `json:"identifier" binding:"required"` // phone number or email
	OTP        string `json:"otp" binding:"required,len=6"`
}

// MFASetupRequest represents the MFA setup request
type MFASetupRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	MFAType  string `json:"mfa_type" binding:"required,oneof=totp"`
}

// MFAVerifyRequest represents the MFA verification request
type MFAVerifyRequest struct {
	UserID     string `json:"user_id" binding:"required"`
	Token      string `json:"token" binding:"required,len=6"`
	BackupCode string `json:"backup_code"` // optional
}

// MFASetupResponse represents the MFA setup response
type MFASetupResponse struct {
	Success     bool     `json:"success"`
	QRCode      string   `json:"qr_code"`
	Secret      string   `json:"secret"`
	BackupCodes []string `json:"backup_codes"`
}

// MFAVerifyResponse represents the MFA verification response
type MFAVerifyResponse struct {
	Success bool        `json:"success"`
	Token   string      `json:"token"`
	User    UserResponse `json:"user"`
}

// FindUserByID finds a user by ID
func FindUserByID(db *sql.DB, id uuid.UUID) (*User, error) {
	query := `
		SELECT id, name, email, phone, role, password_hash, google_id, facebook_id,
		       mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login
		FROM users WHERE id = $1
	`
	
	var user User
	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.PasswordHash,
		&user.GoogleID, &user.FacebookID, &user.MFASecret, &user.IsMFAEnabled,
		&user.MFABackupCodes, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// FindUserByEmail finds a user by email
func FindUserByEmail(db *sql.DB, email string) (*User, error) {
	query := `
		SELECT id, name, email, phone, role, password_hash, google_id, facebook_id,
		       mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login
		FROM users WHERE email = $1
	`
	
	var user User
	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.PasswordHash,
		&user.GoogleID, &user.FacebookID, &user.MFASecret, &user.IsMFAEnabled,
		&user.MFABackupCodes, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// FindUserByPhone finds a user by phone number
func FindUserByPhone(db *sql.DB, phone string) (*User, error) {
	query := `
		SELECT id, name, email, phone, role, password_hash, google_id, facebook_id,
		       mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login
		FROM users WHERE phone = $1
	`
	
	var user User
	err := db.QueryRow(query, phone).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.PasswordHash,
		&user.GoogleID, &user.FacebookID, &user.MFASecret, &user.IsMFAEnabled,
		&user.MFABackupCodes, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// FindUserByGoogleID finds a user by Google ID
func FindUserByGoogleID(db *sql.DB, googleID string) (*User, error) {
	query := `
		SELECT id, name, email, phone, role, password_hash, google_id, facebook_id,
		       mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login
		FROM users WHERE google_id = $1
	`
	
	var user User
	err := db.QueryRow(query, googleID).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.PasswordHash,
		&user.GoogleID, &user.FacebookID, &user.MFASecret, &user.IsMFAEnabled,
		&user.MFABackupCodes, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// FindUserByFacebookID finds a user by Facebook ID
func FindUserByFacebookID(db *sql.DB, facebookID string) (*User, error) {
	query := `
		SELECT id, name, email, phone, role, password_hash, google_id, facebook_id,
		       mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login
		FROM users WHERE facebook_id = $1
	`
	
	var user User
	err := db.QueryRow(query, facebookID).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.PasswordHash,
		&user.GoogleID, &user.FacebookID, &user.MFASecret, &user.IsMFAEnabled,
		&user.MFABackupCodes, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// getNullStringValue returns the string value if the sql.NullString is valid, otherwise returns nil
func getNullStringValue(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}

// CreateUser creates a new user
func CreateUser(db *sql.DB, user *User) error {
	query := `
		INSERT INTO users (name, email, phone, role, password_hash, google_id, facebook_id,
		                   mfa_secret, is_mfa_enabled, mfa_backup_codes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	
	var id uuid.UUID
	var createdAt, updatedAt time.Time

	err := db.QueryRow(query, user.Name,
		getNullStringValue(user.Email), user.Phone, user.Role,
		getNullStringValue(user.PasswordHash), getNullStringValue(user.GoogleID), getNullStringValue(user.FacebookID),
		getNullStringValue(user.MFASecret), user.IsMFAEnabled, getNullStringValue(user.MFABackupCodes)).Scan(&id, &createdAt, &updatedAt)
	
	if err != nil {
		return err
	}
	
	user.ID = id
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt
	
	return nil
}

// UpdateUser updates an existing user
func UpdateUser(db *sql.DB, user *User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, phone = $3, role = $4, password_hash = $5,
		    google_id = $6, facebook_id = $7, mfa_secret = $8, is_mfa_enabled = $9,
		    mfa_backup_codes = $10, updated_at = CURRENT_TIMESTAMP
		WHERE id = $11
	`
	
	_, err := db.Exec(query, user.Name,
		getNullStringValue(user.Email), user.Phone, user.Role,
		getNullStringValue(user.PasswordHash), getNullStringValue(user.GoogleID), getNullStringValue(user.FacebookID),
		getNullStringValue(user.MFASecret), user.IsMFAEnabled, getNullStringValue(user.MFABackupCodes), user.ID)
	
	return err
}

// UpdateLastLogin updates the last login timestamp for a user
func UpdateLastLogin(db *sql.DB, userID uuid.UUID) error {
	query := `UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := db.Exec(query, userID)
	return err
}

// ToResponse converts a User to a UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:           u.ID.String(),
		Name:         u.Name,
		Email:        u.Email.String,
		Phone:        u.Phone,
		Role:         u.Role,
		IsMFAEnabled: u.IsMFAEnabled,
	}
}