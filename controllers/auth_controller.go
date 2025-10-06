package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/aprianimmanuel/rangkaiedu-backend/config"
	"github.com/aprianimmanuel/rangkaiedu-backend/utils/email"
	"github.com/aprianimmanuel/rangkaiedu-backend/utils/otp"
	"github.com/aprianimmanuel/rangkaiedu-backend/utils/password"
	"github.com/aprianimmanuel/rangkaiedu-backend/utils/sms"

	"github.com/golang-jwt/jwt/v4"
	"log"
)

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

func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Optional password strength check if provided
	if req.Password != "" {
		if len(req.Password) < 8 || !hasMinimumComplexity(req.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long and contain a mix of letters, numbers, and symbols"})
			return
		}
	}

	// Database connection
	cfg := config.Load()
	log.Printf("Attempting to connect to database with DSN: %s", cfg.DSN())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the global database connection pool
	pool, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to get database connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	log.Printf("Successfully connected to database")

	// Check for duplicate email
	var existingEmail string
	log.Printf("Checking for duplicate email: %s", req.Email)
	err = pool.QueryRow(ctx, "SELECT email FROM users WHERE email = $1", req.Email).Scan(&existingEmail)
	if err == nil {
		log.Printf("Email already exists: %s", req.Email)
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}
	if err != pgx.ErrNoRows {
		log.Printf("Database error checking email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking email"})
		return
	}
	log.Printf("Email is unique: %s", req.Email)

	// Check for duplicate phone
	var existingPhone string
	err = pool.QueryRow(ctx, "SELECT phone FROM users WHERE phone = $1", req.Phone).Scan(&existingPhone)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Phone number already exists"})
		return
	}
	if err != pgx.ErrNoRows {
		log.Printf("Database error checking phone: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking phone"})
		return
	}

	var passwordHash sql.NullString
	if req.Password != "" {
		ph, err := password.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
			return
		}
		passwordHash = sql.NullString{String: ph, Valid: true}
	}

	// Insert the new user (password_hash can be NULL if optional)
	_, err = pool.Exec(ctx,
		"INSERT INTO users (name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5)",
		req.Name, req.Email, req.Phone, req.Role, passwordHash.String,
	)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	// Generate and send initial OTP for email after registration
	generatedOTP, err := otp.GenerateOTP()
	if err != nil {
		log.Printf("Failed to generate OTP for new user %s: %v", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification OTP"})
		return
	}

	if err := otp.SaveOTP(ctx, pool, req.Email, generatedOTP); err != nil {
		log.Printf("Failed to save OTP for new user %s: %v", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification OTP"})
		return
	}

	if err := email.SendOTPEmail(cfg, req.Email, generatedOTP); err != nil {
		log.Printf("Failed to send OTP email to %s: %v", req.Email, err)
		// Don't fail registration on send error, but log
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Verification OTP sent to email.",
		"user": gin.H{
			"id":    "", // ID is UUID, not returned here for security
			"name":  req.Name,
			"email": req.Email,
			"role":  req.Role,
		},
	})
}

// hasMinimumComplexity performs basic password strength validation
func hasMinimumComplexity(pwd string) bool {
	hasUpper := false
	hasLower := false
	hasDigit := false
	// hasSymbol := false  // Not currently used but kept for future expansion

	for _, char := range pwd {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		// case !('A' <= char && char <= 'Z' || 'a' <= char && char <= 'z' || '0' <= char && char <= '9'):
		// 	hasSymbol = true
		}
	}

	return hasUpper && hasLower && hasDigit // Basic: upper, lower, digit. Symbol optional for now
}

func SendOTPHandler(c *gin.Context) {
	var req SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the global database connection pool
	pool, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to get database connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	// Check if user exists for the identifier
	var exists bool
	var query string
	if req.Type == "email" {
		query = "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	} else {
		query = "SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)"
	}
	err = pool.QueryRow(ctx, query, req.Identifier).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking user existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found for the provided identifier"})
		return
	}

	// Generate and save OTP
	generatedOTP, err := otp.GenerateOTP()
	if err != nil {
		log.Printf("Failed to generate OTP: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}

	if err := otp.SaveOTP(ctx, pool, req.Identifier, generatedOTP); err != nil {
		log.Printf("Failed to save OTP: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save OTP"})
		return
	}

	// Send OTP
	var sendErr error
	if req.Type == "email" {
		sendErr = email.SendOTPEmail(cfg, req.Identifier, generatedOTP)
	} else {
		sendErr = sms.SendOTPSMS(cfg, req.Identifier, generatedOTP)
	}

	if sendErr != nil {
		log.Printf("Failed to send OTP via %s to %s: %v", req.Type, req.Identifier, sendErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func VerifyOTPHandler(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the global database connection pool
	pool, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to get database connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	// Verify OTP
	valid, err := otp.VerifyAndDeleteOTP(ctx, pool, req.Identifier, req.OTP)
	if err != nil {
		log.Printf("OTP verification error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Verification error"})
		return
	}
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	// Query user by identifier (email or phone)
	var userID, email, phone, role string
	var query string
	if isEmail(req.Identifier) {
		query = "SELECT id, email, phone, role FROM users WHERE email = $1"
		err = pool.QueryRow(ctx, query, req.Identifier).Scan(&userID, &email, &phone, &role)
	} else {
		query = "SELECT id, email, phone, role FROM users WHERE phone = $1"
		err = pool.QueryRow(ctx, query, req.Identifier).Scan(&userID, &email, &phone, &role)
	}
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	if err != nil {
		log.Printf("Database error querying user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Generate JWT
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"phone": phone,
		"role":  role,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iss":   "rangkai-edu-backend", // Issuer
		"aud":   "rangkai-edu-frontend", // Audience
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Printf("Failed to generate JWT token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"email": email,
			"phone": phone,
			"role":  role,
		},
	})
}

// isEmail is a simple email validation helper
func isEmail(identifier string) bool {
	return len(identifier) > 0 && identifier[0] != '+' // Basic: phone starts with +, email doesn't
}

// LoginHandler is deprecated in favor of OTP-based authentication. Kept for backward compatibility.
// LoginHandler is deprecated in favor of OTP-based authentication. Kept for backward compatibility.
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the global database connection pool
	pool, err := GetDBConnection()
	if err != nil {
		log.Printf("Failed to get database connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	// Query user by email to get phone number for consistent claims
	var userID, email, phone, role, passwordHash string
	err = pool.QueryRow(ctx, "SELECT id, email, phone, role, password_hash FROM users WHERE email = $1", req.Email).Scan(&userID, &email, &phone, &role, &passwordHash)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	if err != nil {
		log.Printf("Database error querying user %s: %v", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !password.CheckPasswordHash(req.Password, passwordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT with consistent claims
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"phone": phone,
		"role":  role,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iss":   "rangkai-edu-backend", // Issuer
		"aud":   "rangkai-edu-frontend", // Audience
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Printf("Failed to generate JWT token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"email": email,
			"phone": phone,
			"role":  role,
		},
	})
}
