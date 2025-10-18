package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/middleware"
)

func TestAuthRequired_ValidToken(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create a valid token
	_ = config.Load()
	claims := jwt.MapClaims{
		"sub":   "test-user-id",
		"email": "test@example.com",
		"phone": "1234567890",
		"role":  "student",
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iss":   "rangkai-edu-backend", // Issuer
		"aud":   "rangkai-edu-frontend", // Audience
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Create test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	// Create a test handler to verify the user context is set
	var capturedUser middleware.UserContext
	testHandler := func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			t.Error("User not found in context")
			return
		}

		user, ok := userInterface.(middleware.UserContext)
		if !ok {
			t.Error("User context has wrong type")
			return
		}

		capturedUser = user
		c.Status(http.StatusOK)
	}

	// Apply middleware and handler
	middleware.AuthRequired()(c)
	if c.IsAborted() {
		t.Error("Request was aborted unexpectedly")
		return
	}

	// Continue with test handler
	testHandler(c)

	// Verify user context
	if capturedUser.ID != "test-user-id" {
		t.Errorf("Expected user ID 'test-user-id', got '%s'", capturedUser.ID)
	}
	if capturedUser.Email != "test@example.com" {
		t.Errorf("Expected user email 'test@example.com', got '%s'", capturedUser.Email)
	}
	if capturedUser.Phone != "1234567890" {
		t.Errorf("Expected user phone '1234567890', got '%s'", capturedUser.Phone)
	}
	if capturedUser.Role != "student" {
		t.Errorf("Expected user role 'student', got '%s'", capturedUser.Role)
	}
}

func TestAuthRequired_MissingToken(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create test context without authorization header
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	// Apply middleware
	middleware.AuthRequired()(c)

	// Verify response
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if !c.IsAborted() {
		t.Error("Request should be aborted")
	}

	var response map[string]interface{}
	if err := c.ShouldBindJSON(&response); err == nil {
		if errorMsg, ok := response["error"]; !ok || errorMsg != "Authorization header is required" {
			t.Errorf("Expected error message 'Authorization header is required', got '%v'", errorMsg)
		}
	}
}

func TestAuthRequired_InvalidHeaderFormat(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create test context with invalid authorization header
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "InvalidFormat")

	// Apply middleware
	middleware.AuthRequired()(c)

	// Verify response
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if !c.IsAborted() {
		t.Error("Request should be aborted")
	}

	var response map[string]interface{}
	if err := c.ShouldBindJSON(&response); err == nil {
		if errorMsg, ok := response["error"]; !ok || errorMsg != "Invalid authorization header format" {
			t.Errorf("Expected error message 'Invalid authorization header format', got '%v'", errorMsg)
		}
	}
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create test context with invalid token
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid.token.here")

	// Apply middleware
	middleware.AuthRequired()(c)

	// Verify response
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if !c.IsAborted() {
		t.Error("Request should be aborted")
	}

	var response map[string]interface{}
	if err := c.ShouldBindJSON(&response); err == nil {
		if errorMsg, ok := response["error"]; !ok || errorMsg != "Invalid token" {
			t.Errorf("Expected error message 'Invalid token', got '%v'", errorMsg)
		}
	}
}

func TestAuthRequired_ExpiredToken(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create an expired token
	_ = config.Load()
	claims := jwt.MapClaims{
		"sub":   "test-user-id",
		"email": "test@example.com",
		"phone": "1234567890",
		"role":  "student",
		"iat":   time.Now().Add(-2 * time.Hour).Unix(),
		"exp":   time.Now().Add(-1 * time.Hour).Unix(),
		"iss":   "rangkai-edu-backend", // Issuer
		"aud":   "rangkai-edu-frontend", // Audience
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Create test context with expired token
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	// Apply middleware
	middleware.AuthRequired()(c)

	// Verify response
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if !c.IsAborted() {
		t.Error("Request should be aborted")
	}

	var response map[string]interface{}
	if err := c.ShouldBindJSON(&response); err == nil {
		if errorMsg, ok := response["error"]; !ok || errorMsg != "Invalid token" {
			t.Errorf("Expected error message 'Invalid token', got '%v'", errorMsg)
		}
	}
}
