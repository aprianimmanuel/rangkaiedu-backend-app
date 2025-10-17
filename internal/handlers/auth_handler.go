package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/monitoring"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/db"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/email"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/otp"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/password"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/sms"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/whatsapp"
	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v4"
	"log"
)

// AuthHandler handles authentication-related operations
type AuthHandler struct{}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,max=255"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,max=20"`
	Password string `json:"password" binding:"omitempty,min=8"` // Optional password
	Role     string `json:"role" binding:"required,oneof=admin teacher student"`
}

type SendOTPRequest struct {
	Identifier string `json:"identifier" binding:"required"` // email, phone, or WhatsApp phone
	Type       string `json:"type" binding:"required,oneof=email phone whatsapp"`
}

type VerifyOTPRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	OTP        string `json:"otp" binding:"required,len=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"omitempty,min=8"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Log security event for invalid request data using security logger
		monitoring.LogAuthFailure(c.Request.Context(), req.Email, c.ClientIP(), c.Request.UserAgent(), map[string]interface{}{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
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
	_ = config.Load()
	log.Printf("Attempting to connect to database with DSN: %s", config.Load().DSN())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the global database connection pool
	pool := db.GetDB()
	if pool == nil {
		log.Printf("Failed to get database connection: database pool is not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	log.Printf("Successfully connected to database")

	// Check for duplicate email
	var existingEmail string
	log.Printf("Checking for duplicate email: %s", req.Email)
	err := pool.QueryRowContext(ctx, "SELECT email FROM users WHERE email = $1", req.Email).Scan(&existingEmail)
	if err == nil {
		log.Printf("Email already exists: %s", req.Email)
		// Log security event for duplicate email attempt
		monitoring.LogAuthFailure(c.Request.Context(), req.Email, c.ClientIP(), c.Request.UserAgent(), map[string]interface{}{
			"error": "Email already exists",
		})
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}
	if err != sql.ErrNoRows {
		log.Printf("Database error checking email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking email"})
		return
	}
	log.Printf("Email is unique: %s", req.Email)

	// Check for duplicate phone
	var existingPhone string
	err = pool.QueryRowContext(ctx, "SELECT phone FROM users WHERE phone = $1", req.Phone).Scan(&existingPhone)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Phone number already exists"})
		return
	}
	if err != sql.ErrNoRows {
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
	_, err = pool.ExecContext(ctx,
		"INSERT INTO users (name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5)",
		req.Name, req.Email, req.Phone, req.Role, passwordHash.String,
	)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		// Log security event for failed user creation
		monitoring.LogAuthFailure(c.Request.Context(), req.Email, c.ClientIP(), c.Request.UserAgent(), map[string]interface{}{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	// Log successful user registration
	monitoring.LogAuthSuccess(c.Request.Context(), "", req.Email, c.ClientIP(), c.Request.UserAgent(), map[string]interface{}{
		"name": req.Name,
		"role": req.Role,
	})

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

	if err := email.SendOTPEmail(config.Load(), req.Email, generatedOTP); err != nil {
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

// SendOTP handles OTP sending
func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the global database connection pool
	pool := db.GetDB()
	if pool == nil {
		log.Printf("Failed to get database connection: database pool is not initialized")
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
	err := pool.QueryRowContext(ctx, query, req.Identifier).Scan(&exists)
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
	} else if req.Type == "whatsapp" {
		sendErr = whatsapp.SendOTPWithRateLimit(cfg, req.Identifier, generatedOTP)
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

// VerifyOTP handles OTP verification
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the global database connection pool
	pool := db.GetDB()
	if pool == nil {
		log.Printf("Failed to get database connection: database pool is not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	// Verify OTP
	valid, err := otp.VerifyAndDeleteOTP(ctx, pool, req.Identifier, req.OTP)
	if err != nil {
		log.Printf("Failed to verify OTP: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify OTP"})
		return
	}
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the global database connection pool
	pool := db.GetDB()
	if pool == nil {
		log.Printf("Failed to get database connection: database pool is not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	// Check if user exists
	var user struct {
		ID           string
		Name         string
		Email        string
		Phone        string
		Role         string
		PasswordHash sql.NullString
		GoogleID     sql.NullString
		FacebookID   sql.NullString
	}
	err := pool.QueryRowContext(ctx, "SELECT id, name, email, phone, role, password_hash, google_id, facebook_id FROM users WHERE email = $1", req.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.PasswordHash, &user.GoogleID, &user.FacebookID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Log security event for login attempt with non-existent email
			monitoring.LogAuthFailure(ctx, req.Email, c.ClientIP(), c.Request.UserAgent(), map[string]interface{}{
				"error": "User not found",
			})
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		log.Printf("Database error checking user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if user is a social auth user (has GoogleID or FacebookID)
	isSocialAuthUser := (user.GoogleID.Valid && user.GoogleID.String != "") || (user.FacebookID.Valid && user.FacebookID.String != "")

	// If user is a social auth user, they should not be able to login with password
	if isSocialAuthUser && req.Password != "" {
		// Log security event for social auth user trying to login with password
		monitoring.LogAuthFailure(ctx, req.Email, c.ClientIP(), c.Request.UserAgent(), map[string]interface{}{
			"error": "Social auth user attempting password login",
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Social auth users cannot login with password"})
		return
	}

	// If no password is provided, initiate OTP challenge
	if req.Password == "" {
		// Generate and save OTP
		generatedOTP, err := otp.GenerateOTP()
		if err != nil {
			log.Printf("Failed to generate OTP: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
			return
		}

		if err := otp.SaveOTP(ctx, pool, req.Email, generatedOTP); err != nil {
			log.Printf("Failed to save OTP: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save OTP"})
			return
		}

		// Send OTP via email
		if err := email.SendOTPEmail(cfg, req.Email, generatedOTP); err != nil {
			log.Printf("Failed to send OTP email to %s: %v", req.Email, err)
			// Don't fail the request on send error, but log
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "OTP sent to your email",
			"otp_required": true,
		})
		return
	}

	// Verify password if provided and user is not a social auth user
	if req.Password != "" && user.PasswordHash.Valid && user.PasswordHash.String != "" {
		if !password.CheckPassword(req.Password, user.PasswordHash.String) {
			// Log security event for failed password attempt
			monitoring.LogAuthFailure(ctx, req.Email, c.ClientIP(), c.Request.UserAgent(), map[string]interface{}{
				"error": "Invalid password",
			})
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hour expiration
	})

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Log successful login
	monitoring.LogAuthSuccess(ctx, user.ID, user.Email, c.ClientIP(), c.Request.UserAgent(), map[string]interface{}{
		"name": user.Name,
		"role": user.Role,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"phone": user.Phone,
			"role":  user.Role,
		},
	})
}

// LoginWithWhatsApp handles user login with WhatsApp OTP
func (h *AuthHandler) LoginWithWhatsApp(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the global database connection pool
	pool := db.GetDB()
	if pool == nil {
		log.Printf("Failed to get database connection: database pool is not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	// Check if user exists
	var user struct {
		ID           string
		Name         string
		Email        string
		Phone        string
		Role         string
		PasswordHash sql.NullString
		GoogleID     sql.NullString
		FacebookID   sql.NullString
	}
	err := pool.QueryRowContext(ctx, "SELECT id, name, email, phone, role, password_hash, google_id, facebook_id FROM users WHERE email = $1", req.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.PasswordHash, &user.GoogleID, &user.FacebookID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Log security event for login attempt with non-existent email
			monitoring.LogAuthFailure(ctx, req.Email, c.ClientIP(), c.Request.UserAgent(), map[string]interface{}{
				"error": "User not found",
			})
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		log.Printf("Database error checking user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if user is a social auth user (has GoogleID or FacebookID)
	isSocialAuthUser := (user.GoogleID.Valid && user.GoogleID.String != "") || (user.FacebookID.Valid && user.FacebookID.String != "")

	// If user is a social auth user, they should not be able to login with password
	if isSocialAuthUser && req.Password != "" {
		// Log security event for social auth user trying to login with password
		monitoring.LogAuthFailure(ctx, req.Email, c.ClientIP(), c.Request.UserAgent(), map[string]interface{}{
			"error": "Social auth user attempting password login",
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Social auth users cannot login with password"})
		return
	}

	// If no password is provided, initiate WhatsApp OTP challenge
	if req.Password == "" {
		// Generate and save OTP
		generatedOTP, err := otp.GenerateOTP()
		if err != nil {
			log.Printf("Failed to generate OTP: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
			return
		}

		if err := otp.SaveOTP(ctx, pool, user.Phone, generatedOTP); err != nil {
			log.Printf("Failed to save OTP: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save OTP"})
			return
		}

		// Send OTP via WhatsApp
		if err := whatsapp.SendOTPWithRateLimit(cfg, user.Phone, generatedOTP); err != nil {
			log.Printf("Failed to send OTP WhatsApp to %s: %v", user.Phone, err)
			// Don't fail the request on send error, but log
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "OTP sent to your WhatsApp number",
			"otp_required": true,
		})
		return
	}

	// Verify password if provided and user is not a social auth user
	if req.Password != "" && user.PasswordHash.Valid && user.PasswordHash.String != "" {
		if !password.CheckPassword(req.Password, user.PasswordHash.String) {
			// Log security event for failed password attempt
			monitoring.LogAuthFailure(ctx, req.Email, c.ClientIP(), c.Request.UserAgent(), map[string]interface{}{
				"error": "Invalid password",
			})
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hour expiration
	})

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Log successful login
	monitoring.LogAuthSuccess(ctx, user.ID, user.Email, c.ClientIP(), c.Request.UserAgent(), map[string]interface{}{
		"name": user.Name,
		"role": user.Role,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"phone": user.Phone,
			"role":  user.Role,
		},
	})
}

// hasMinimumComplexity checks if password meets minimum complexity requirements
func hasMinimumComplexity(password string) bool {
	// At least 8 characters
	if len(password) < 8 {
		return false
	}

	// Contains at least one letter
	hasLetter := false
	hasNumber := false
	hasSymbol := false

	for _, char := range password {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			hasLetter = true
		} else if char >= '0' && char <= '9' {
			hasNumber = true
		} else {
			hasSymbol = true
		}
	}

	return hasLetter && (hasNumber || hasSymbol)
}
