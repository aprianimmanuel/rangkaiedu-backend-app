package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/repositories"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/email"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/otp"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/password"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/sms"
	"github.com/golang-jwt/jwt/v4"
)

// Custom error types
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type DuplicateError struct {
	Message string
}

func (e *DuplicateError) Error() string {
	return e.Message
}

type DatabaseError struct {
	Message string
	Err     error
}

func (e *DatabaseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

type ProcessingError struct {
	Message string
	Err     error
}

func (e *ProcessingError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

// Request types
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,max=255"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,max=20"`
	Password string `json:"password" binding:"omitempty,min=8"` // Optional password
	Role     string `json:"role" binding:"required,oneof=admin teacher student"`
}

type SendOTPRequest struct {
	Identifier string `json:"identifier" binding:"required"` // email or phone
	Type       string `json:"type" binding:"required,oneof=email phone"`
}

type VerifyOTPRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	OTP        string `json:"otp" binding:"required,len=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Response types
type AuthResponse struct {
	Status  string            `json:"status"`
	Token   string            `json:"token,omitempty"`
	User    models.UserResponse `json:"user,omitempty"`
	Message string            `json:"message,omitempty"`
}

// UserService defines the interface for user business logic
type UserService interface {
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByPhone(phone string) (*models.User, error)
	GetUserByGoogleID(googleID string) (*models.User, error)
	GetUserByFacebookID(facebookID string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	UpdateLastLogin(userID uuid.UUID) error
	
	// Authentication methods
	RegisterUser(req RegisterRequest) (*AuthResponse, error)
	SendOTP(req SendOTPRequest, db *sql.DB, cfg *config.Config) error
	VerifyOTP(req VerifyOTPRequest, db *sql.DB, cfg *config.Config) (*AuthResponse, error)
	Login(req LoginRequest, db *sql.DB, cfg *config.Config) (*AuthResponse, error)
}

// userService implements UserService interface
type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetUserByID gets a user by ID
func (s *userService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// GetUserByEmail gets a user by email
func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}

// GetUserByPhone gets a user by phone number
func (s *userService) GetUserByPhone(phone string) (*models.User, error) {
	return s.userRepo.FindByPhone(phone)
}

// GetUserByGoogleID gets a user by Google ID
func (s *userService) GetUserByGoogleID(googleID string) (*models.User, error) {
	return s.userRepo.FindByGoogleID(googleID)
}

// GetUserByFacebookID gets a user by Facebook ID
func (s *userService) GetUserByFacebookID(facebookID string) (*models.User, error) {
	return s.userRepo.FindByFacebookID(facebookID)
}

// CreateUser creates a new user
func (s *userService) CreateUser(user *models.User) error {
	return s.userRepo.Create(user)
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

// UpdateLastLogin updates the last login timestamp for a user
func (s *userService) UpdateLastLogin(userID uuid.UUID) error {
	return s.userRepo.UpdateLastLogin(userID)
}

// RegisterUser handles user registration
func (s *userService) RegisterUser(req RegisterRequest) (*AuthResponse, error) {
	// Optional password strength check if provided
	if req.Password != "" {
		if len(req.Password) < 8 || !hasMinimumComplexity(req.Password) {
			return nil, &ValidationError{Message: "Password must be at least 8 characters long and contain a mix of letters, numbers, and symbols"}
		}
	}

	// Check for duplicate email
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, &DuplicateError{Message: "Email already exists"}
	}
	
	// Check for duplicate phone
	existingUser, err = s.userRepo.FindByPhone(req.Phone)
	if err == nil && existingUser != nil {
		return nil, &DuplicateError{Message: "Phone number already exists"}
	}

	// Create user object
	user := &models.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        sql.NullString{String: req.Email, Valid: true},
		Phone:        req.Phone,
		Role:         req.Role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Set password hash if provided
	if req.Password != "" {
		hash, err := password.HashPassword(req.Password)
		if err != nil {
			return nil, &ProcessingError{Message: "Failed to process password", Err: err}
		}
		user.PasswordHash = sql.NullString{String: hash, Valid: true}
	}

	// Create user in database
	if err := s.userRepo.Create(user); err != nil {
		return nil, &DatabaseError{Message: "Failed to create user", Err: err}
	}

	// Convert to response
	userResponse := user.ToResponse()

	return &AuthResponse{
		Status:  "success",
		Message: "User registered successfully",
		User:    userResponse,
	}, nil
}

// SendOTP handles sending OTP to users
func (s *userService) SendOTP(req SendOTPRequest, db *sql.DB, cfg *config.Config) error {
	// Check if user exists for the identifier
	var exists bool
	var query string
	if req.Type == "email" {
		query = "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	} else {
		query = "SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)"
	}
	
	err := db.QueryRowContext(context.Background(), query, req.Identifier).Scan(&exists)
	if err != nil {
		return &DatabaseError{Message: "Database error checking user existence", Err: err}
	}
	
	if !exists {
		return &NotFoundError{Message: "User not found for the provided identifier"}
	}

	// Generate and save OTP
	generatedOTP, err := otp.GenerateOTP()
	if err != nil {
		return &ProcessingError{Message: "Failed to generate OTP", Err: err}
	}

	if err := otp.SaveOTP(context.Background(), db, req.Identifier, generatedOTP); err != nil {
		return &DatabaseError{Message: "Failed to save OTP", Err: err}
	}

	// Send OTP
	var sendErr error
	if req.Type == "email" {
		sendErr = email.SendOTPEmail(cfg, req.Identifier, generatedOTP)
	} else {
		sendErr = sms.SendOTPSMS(cfg, req.Identifier, generatedOTP)
	}

	if sendErr != nil {
		// Don't fail the operation on send error, but log it
		log.Printf("Warning: Failed to send OTP via %s to %s: %v", req.Type, req.Identifier, sendErr)
	}

	return nil
}

// VerifyOTP handles OTP verification
func (s *userService) VerifyOTP(req VerifyOTPRequest, db *sql.DB, cfg *config.Config) (*AuthResponse, error) {
	// Verify OTP
	valid, err := otp.VerifyAndDeleteOTP(context.Background(), db, req.Identifier, req.OTP)
	if err != nil {
		return nil, &DatabaseError{Message: "OTP verification error", Err: err}
	}
	
	if !valid {
		return nil, &UnauthorizedError{Message: "Invalid or expired OTP"}
	}

	// Query user by identifier (email or phone)
	var user *models.User
	var queryErr error
	
	if isEmail(req.Identifier) {
		user, queryErr = s.userRepo.FindByEmail(req.Identifier)
	} else {
		user, queryErr = s.userRepo.FindByPhone(req.Identifier)
	}
	
	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return nil, &UnauthorizedError{Message: "User not found"}
		}
		return nil, &DatabaseError{Message: "Database error querying user", Err: queryErr}
	}

	// Generate JWT
	tokenString, err := generateJWT(user, cfg)
	if err != nil {
		return nil, &ProcessingError{Message: "Failed to generate token", Err: err}
	}

	// Convert to response
	userResponse := user.ToResponse()

	return &AuthResponse{
		Status: "success",
		Token:  tokenString,
		User:   userResponse,
	}, nil
}

// Login handles traditional email/password login
func (s *userService) Login(req LoginRequest, db *sql.DB, cfg *config.Config) (*AuthResponse, error) {
	// Query user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &UnauthorizedError{Message: "User not found"}
		}
		return nil, &DatabaseError{Message: "Database error querying user", Err: err}
	}

	// Check if user is a social authentication user
	if isSocialAuthUser(user.GoogleID.String, user.FacebookID.String) {
		return nil, &ForbiddenError{Message: "This account was created using social authentication. Please use Google or Facebook to login."}
	}

	// Check if user has a password set (for non-social users)
	if !user.PasswordHash.Valid || user.PasswordHash.String == "" {
		return nil, &UnauthorizedError{Message: "No password set for this account. Please use a different authentication method."}
	}

	// Verify password
	if !password.CheckPasswordHash(req.Password, user.PasswordHash.String) {
		return nil, &UnauthorizedError{Message: "Invalid credentials"}
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		log.Printf("Warning: Failed to update last login for user %s: %v", user.ID, err)
	}

	// Generate JWT
	tokenString, err := generateJWT(user, cfg)
	if err != nil {
		return nil, &ProcessingError{Message: "Failed to generate token", Err: err}
	}

	// Convert to response
	userResponse := user.ToResponse()

	return &AuthResponse{
		Status: "success",
		Token:  tokenString,
		User:   userResponse,
	}, nil
}

// Helper functions
func hasMinimumComplexity(pwd string) bool {
	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range pwd {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

func isEmail(identifier string) bool {
	return len(identifier) > 0 && identifier[0] != '+' // Basic: phone starts with +, email doesn't
}

func isSocialAuthUser(googleID, facebookID string) bool {
	return googleID != "" || facebookID != ""
}

func generateJWT(user *models.User, cfg *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email.String,
		"phone": user.Phone,
		"role":  user.Role,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iss":   "rangkai-edu-backend", // Issuer
		"aud":   "rangkai-edu-frontend", // Audience
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}