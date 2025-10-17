package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"github.com/google/uuid"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByPhone(phone string) (*models.User, error)
	FindByGoogleID(googleID string) (*models.User, error)
	FindByFacebookID(facebookID string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	UpdateLastLogin(userID uuid.UUID) error
}

// userRepository implements UserRepository interface
type userRepository struct {
	db *sql.DB
}
// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
// scanMFABackupCodes handles scanning PostgreSQL TEXT[] array into sql.NullString
func scanMFABackupCodes(src interface{}) (sql.NullString, error) {
	if src == nil {
		return sql.NullString{Valid: false}, nil
	}
	
	switch v := src.(type) {
	case []byte:
		// Try JSON first, then PostgreSQL array format
		str := string(v)
		result, err := tryParseArrayFormats(str)
		if err != nil {
			return sql.NullString{}, err
		}
		jsonData, _ := json.Marshal(result)
		return sql.NullString{
			String: string(jsonData),
			Valid:  true,
		}, nil
	case string:
		// Try JSON first, then PostgreSQL array format
		result, err := tryParseArrayFormats(v)
		if err != nil {
			return sql.NullString{}, err
		}
		jsonData, _ := json.Marshal(result)
		return sql.NullString{
			String: string(jsonData),
			Valid:  true,
		}, nil
	case []string:
		jsonData, _ := json.Marshal(v)
		return sql.NullString{
			String: string(jsonData),
			Valid:  true,
		}, nil
	default:
		return sql.NullString{}, fmt.Errorf("unsupported type for mfa_backup_codes: %T", src)
	}
}

// tryParseArrayFormats tries to parse the string as JSON first, then as PostgreSQL array
func tryParseArrayFormats(str string) ([]string, error) {
	// Try JSON format first: ["code1","code2"]
	var result []string
	if err := json.Unmarshal([]byte(str), &result); err == nil {
		return result, nil
	}
	
	// Try PostgreSQL array format: {"code1","code2"}
	return parsePostgreSQLArray(str)
}

// parsePostgreSQLArray parses PostgreSQL array string format {"element1","element2"} into []string
func parsePostgreSQLArray(arrStr string) ([]string, error) {
	if len(arrStr) == 0 || arrStr == "{}" {
		return nil, nil
	}
	
	// Remove the curly braces
	if arrStr[0] == '{' && arrStr[len(arrStr)-1] == '}' {
		arrStr = arrStr[1 : len(arrStr)-1]
	}
	
	// Split by comma, but handle escaped quotes if needed
	var elements []string
	var current strings.Builder
	inQuotes := false
	
	for i := 0; i < len(arrStr); i++ {
		char := arrStr[i]
		
		if char == '"' {
			inQuotes = !inQuotes
			current.WriteByte(char)
		} else if char == ',' && !inQuotes {
			elements = append(elements, current.String())
			current.Reset()
		} else {
			current.WriteByte(char)
		}
	}
	
	// Add the last element
	if current.Len() > 0 {
		elements = append(elements, current.String())
	}
	
	// Remove quotes from each element
	var result []string
	for _, elem := range elements {
		if len(elem) >= 2 && elem[0] == '"' && elem[len(elem)-1] == '"' {
			result = append(result, elem[1:len(elem)-1])
		} else {
			result = append(result, elem)
		}
	}
	
	return result, nil
}
// valueMFABackupCodes handles converting sql.NullString to a database driver value
func valueMFABackupCodes(codes sql.NullString) interface{} {
	if !codes.Valid || codes.String == "" {
		return nil
	}
	
	// Parse the JSON array and convert to PostgreSQL array format
	var result []string
	if err := json.Unmarshal([]byte(codes.String), &result); err == nil {
		return formatAsPostgreSQLArray(result)
	}
	
	// If parsing fails, return as-is
	return codes.String
}

// formatAsPostgreSQLArray formats a slice as PostgreSQL array string
func formatAsPostgreSQLArray(elements []string) string {
	if len(elements) == 0 {
		return "{}"
	}
	
	var builder strings.Builder
	builder.WriteByte('{')
	
	for i, elem := range elements {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteByte('"')
		builder.WriteString(elem)
		builder.WriteByte('"')
	}
	
	builder.WriteByte('}')
	return builder.String()
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, name, email, phone, role, password_hash, google_id, facebook_id,
		       mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login
		FROM users WHERE id = $1
	`
	
	var user models.User
	var mfaBackupCodes sql.NullString
	
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.PasswordHash,
		&user.GoogleID, &user.FacebookID, &user.MFASecret, &user.IsMFAEnabled,
		&mfaBackupCodes, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	
	if err != nil {
		return nil, err
	}
	
	user.MFABackupCodes = mfaBackupCodes
	
	return &user, nil
}

// FindByEmail finds a user by email
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, name, email, phone, role, password_hash, google_id, facebook_id,
		       mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login
		FROM users WHERE email = $1
	`
	
	var user models.User
	var mfaBackupCodes sql.NullString
	
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.PasswordHash,
		&user.GoogleID, &user.FacebookID, &user.MFASecret, &user.IsMFAEnabled,
		&mfaBackupCodes, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	
	if err != nil {
		return nil, err
	}
	
	user.MFABackupCodes = mfaBackupCodes
	
	return &user, nil
}

// FindByPhone finds a user by phone number
func (r *userRepository) FindByPhone(phone string) (*models.User, error) {
	query := `
		SELECT id, name, email, phone, role, password_hash, google_id, facebook_id,
		       mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login
		FROM users WHERE phone = $1
	`
	
	var user models.User
	var mfaBackupCodes sql.NullString
	
	err := r.db.QueryRow(query, phone).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.PasswordHash,
		&user.GoogleID, &user.FacebookID, &user.MFASecret, &user.IsMFAEnabled,
		&mfaBackupCodes, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	
	if err != nil {
		return nil, err
	}
	
	user.MFABackupCodes = mfaBackupCodes
	
	return &user, nil
}

// FindByGoogleID finds a user by Google ID
func (r *userRepository) FindByGoogleID(googleID string) (*models.User, error) {
	query := `
		SELECT id, name, email, phone, role, password_hash, google_id, facebook_id,
		       mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login
		FROM users WHERE google_id = $1
	`
	
	var user models.User
	var mfaBackupCodes sql.NullString
	
	err := r.db.QueryRow(query, googleID).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.PasswordHash,
		&user.GoogleID, &user.FacebookID, &user.MFASecret, &user.IsMFAEnabled,
		&mfaBackupCodes, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	
	if err != nil {
		return nil, err
	}
	
	user.MFABackupCodes = mfaBackupCodes
	
	return &user, nil
}

// FindByFacebookID finds a user by Facebook ID
func (r *userRepository) FindByFacebookID(facebookID string) (*models.User, error) {
	query := `
		SELECT id, name, email, phone, role, password_hash, google_id, facebook_id,
		       mfa_secret, is_mfa_enabled, mfa_backup_codes, created_at, updated_at, last_login
		FROM users WHERE facebook_id = $1
	`
	
	var user models.User
	var mfaBackupCodes sql.NullString
	
	err := r.db.QueryRow(query, facebookID).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.PasswordHash,
		&user.GoogleID, &user.FacebookID, &user.MFASecret, &user.IsMFAEnabled,
		&mfaBackupCodes, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	
	if err != nil {
		return nil, err
	}
	
	user.MFABackupCodes = mfaBackupCodes
	
	return &user, nil
}
// Create creates a new user
func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (name, email, phone, role, password_hash, google_id, facebook_id,
		                   mfa_secret, is_mfa_enabled, mfa_backup_codes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	
	var id uuid.UUID
	var createdAt, updatedAt time.Time
	
	mfaBackupCodesValue := valueMFABackupCodes(user.MFABackupCodes)
	
	err := r.db.QueryRow(query, user.Name, user.Email, user.Phone, user.Role,
		user.PasswordHash, user.GoogleID, user.FacebookID,
		user.MFASecret, user.IsMFAEnabled, mfaBackupCodesValue).Scan(&id, &createdAt, &updatedAt)
	
	if err != nil {
		return err
	}
	
	user.ID = id
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt
	
	return nil
}

// Update updates an existing user
func (r *userRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, phone = $3, role = $4, password_hash = $5,
		    google_id = $6, facebook_id = $7, mfa_secret = $8, is_mfa_enabled = $9,
		    mfa_backup_codes = $10, updated_at = CURRENT_TIMESTAMP
		WHERE id = $11
	`
	
	mfaBackupCodesValue := valueMFABackupCodes(user.MFABackupCodes)
	
	_, err := r.db.Exec(query, user.Name, user.Email, user.Phone, user.Role,
		user.PasswordHash, user.GoogleID, user.FacebookID,
		user.MFASecret, user.IsMFAEnabled, mfaBackupCodesValue, user.ID)
	
	return err
}

// UpdateLastLogin updates the last login timestamp for a user
func (r *userRepository) UpdateLastLogin(userID uuid.UUID) error {
	query := `UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}