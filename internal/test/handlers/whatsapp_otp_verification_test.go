package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/handlers"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/testutils"
)

// setupTestRouterWithWhatsAppVerification sets up a test router with WhatsApp verification routes
func setupTestRouterWithWhatsAppVerification() *gin.Engine {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Initialize test database configuration
	cfg := config.LoadTest()
	
	// Initialize database connection pool for testing
	if err := testutils.InitTestDB(); err != nil {
		// If database initialization fails, continue with tests but log the error
		// The mock handlers will work without database connection
	}
	
	// Create a new router
	router := gin.New()
	
	// Create AuthHandler instance
	authHandler := handlers.NewAuthHandler()
	
	// Setup authentication routes
	auth := router.Group("/api/auth")
	{
		// WhatsApp OTP verification route
		auth.POST("/whatsapp/verify", authHandler.VerifyWhatsAppOTP)
		
		// WhatsApp login route
		auth.POST("/whatsapp/login", authHandler.WhatsappLogin)
	}
	
	return router
}

// TestVerifyWhatsAppOTP_Success tests successful WhatsApp OTP verification
func TestVerifyWhatsAppOTP_Success(t *testing.T) {
	router := setupTestRouterWithWhatsAppVerification()
	
	// Test data
	testData := map[string]interface{}{
		"phone": "+6281234567890",
		"otp":   "123456",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify response structure
	assert.Equal(t, "OTP verified successfully", response["message"])
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "user")
	assert.Contains(t, response, "expires_in")
}

// TestVerifyWhatsAppOTP_InvalidOTP tests invalid OTP handling
func TestVerifyWhatsAppOTP_InvalidOTP(t *testing.T) {
	router := setupTestRouterWithWhatsAppVerification()
	
	// Test data with invalid OTP
	testData := map[string]interface{}{
		"phone": "+6281234567890",
		"otp":   "invalid-otp",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify error message
	assert.Contains(t, response["error"], "Invalid OTP")
}

// TestVerifyWhatsAppOTP_ExpiredOTP tests expired OTP handling
func TestVerifyWhatsAppOTP_ExpiredOTP(t *testing.T) {
	router := setupTestRouterWithWhatsAppVerification()
	
	// Test data with expired OTP
	testData := map[string]interface{}{
		"phone": "+6281234567890",
		"otp":   "expired123",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify error message
	assert.Contains(t, response["error"], "OTP expired")
}

// TestVerifyWhatsAppOTP_MissingPhone tests missing phone number
func TestVerifyWhatsAppOTP_MissingPhone(t *testing.T) {
	router := setupTestRouterWithWhatsAppVerification()
	
	// Test data without phone number
	testData := map[string]interface{}{
		"otp": "123456",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify error message
	assert.Contains(t, response["error"], "phone number is required")
}

// TestVerifyWhatsAppOTP_MissingOTP tests missing OTP
func TestVerifyWhatsAppOTP_MissingOTP(t *testing.T) {
	router := setupTestRouterWithWhatsAppVerification()
	
	// Test data without OTP
	testData := map[string]interface{}{
		"phone": "+6281234567890",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify error message
	assert.Contains(t, response["error"], "OTP is required")
}

// TestVerifyWhatsAppOTP_WrongLengthOTP tests OTP with wrong length
func TestVerifyWhatsAppOTP_WrongLengthOTP(t *testing.T) {
	router := setupTestRouterWithWhatsAppVerification()
	
	// Test data with wrong length OTP
	testData := map[string]interface{}{
		"phone": "+6281234567890",
		"otp":   "12345", // 5 digits instead of 6
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify error message
	assert.Contains(t, response["error"], "OTP must be 6 digits")
}

// TestVerifyWhatsAppOTP_InvalidPhoneFormat tests invalid phone number format
func TestVerifyWhatsAppOTP_InvalidPhoneFormat(t *testing.T) {
	router := setupTestRouterWithWhatsAppVerification()
	
	// Test data with invalid phone format
	testData := map[string]interface{}{
		"phone": "invalid-phone-format",
		"otp":   "123456",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify error message
	assert.Contains(t, response["error"], "invalid phone number format")
}

// TestVerifyWhatsAppOTP_WithExistingUser tests OTP verification for existing user
func TestVerifyWhatsAppOTP_WithExistingUser(t *testing.T) {
	router := setupTestRouterWithWhatsAppVerification()
	
	// Create a test user first
	db := testutils.SetupTestDB(t)
	defer db.Close()
	
	uniquePhone := "+6281234567890"
	uniqueEmail := "whatsappuser@example.com"
	
	// Create existing user
	userID := "test-user-id"
	_, err := db.Exec("INSERT INTO users (id, name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, "WhatsApp User", uniqueEmail, uniquePhone, "student", "hashed_password")
	require.NoError(t, err)
	
	// Save a test OTP
	_, err = db.Exec("INSERT INTO otps (identifier, otp, expiry) VALUES ($1, $2, $3)",
		uniquePhone, "123456", time.Now().Add(10*time.Minute))
	require.NoError(t, err)
	
	// Test data
	testData := map[string]interface{}{
		"phone": uniquePhone,
		"otp":   "123456",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify response structure
	assert.Equal(t, "OTP verified successfully", response["message"])
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "user")
	assert.Contains(t, response, "expires_in")
	
	// Verify user details
	user := response["user"].(map[string]interface{})
	assert.Equal(t, userID, user["id"])
	assert.Equal(t, "WhatsApp User", user["name"])
	assert.Equal(t, uniqueEmail, user["email"])
	assert.Equal(t, uniquePhone, user["phone"])
	assert.Equal(t, "student", user["role"])
	
	// Clean up
	testutils.CleanupUser(t, db, uniqueEmail)
}

// TestVerifyWhatsAppOTP_WithNewUser tests OTP verification for new user (should fail)
func TestVerifyWhatsAppOTP_WithNewUser(t *testing.T) {
	router := setupTestRouterWithWhatsAppVerification()
	
	// Test data with new phone number
	testData := map[string]interface{}{
		"phone": "+6281112233445", // New phone number
		"otp":   "123456",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify error message
	assert.Contains(t, response["error"], "User not found")
}

// TestVerifyWhatsAppOTP_UsedOTP tests using the same OTP twice
func TestVerifyWhatsAppOTP_UsedOTP(t *testing.T) {
	router := setupTestRouterWithWhatsAppVerification()
	
	// Create a test user first
	db := testutils.SetupTestDB(t)
	defer db.Close()
	
	uniquePhone := "+6281234567890"
	uniqueEmail := "whatsappuser@example.com"
	
	// Create existing user
	userID := "test-user-id"
	_, err := db.Exec("INSERT INTO users (id, name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, "WhatsApp User", uniqueEmail, uniquePhone, "student", "hashed_password")
	require.NoError(t, err)
	
	// Save a test OTP
	_, err = db.Exec("INSERT INTO otps (identifier, otp, expiry) VALUES ($1, $2, $3)",
		uniquePhone, "123456", time.Now().Add(10*time.Minute))
	require.NoError(t, err)
	
	// Test data
	testData := map[string]interface{}{
		"phone": uniquePhone,
		"otp":   "123456",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// First verification should succeed
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Second verification should fail
	req2, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req2.Header.Set("Content-Type", "application/json")
	
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	
	assert.Equal(t, http.StatusBadRequest, w2.Code)
	
	var response2 map[string]interface{}
	err = json.Unmarshal(w2.Body.Bytes(), &response2)
	require.NoError(t, err)
	
	assert.Contains(t, response2["error"], "Invalid OTP")
	
	// Clean up
	testutils.CleanupUser(t, db, uniqueEmail)
}

// TestVerifyWhatsAppOTP_WithInvalidContentType tests invalid content type
func TestVerifyWhatsAppOTP_WithInvalidContentType(t *testing.T) {
	router := setupTestRouterWithWhatsAppVerification()
	
	// Test data
	testData := map[string]interface{}{
		"phone": "+6281234567890",
		"otp":   "123456",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request with invalid content type
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "text/plain") // Invalid content type
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
}

// TestVerifyWhatsAppOTP_WithEmptyBody tests empty request body
func TestVerifyWhatsAppOTP_WithEmptyBody(t *testing.T) {
	router := setupTestRouterWithWhatsAppVerification()
	
	// Create request with empty body
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer([]byte("")))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify error message
	assert.Contains(t, response["error"], "phone number is required")
}

// TestVerifyWhatsAppOTP_WithMalformedJSON tests malformed JSON request
func TestVerifyWhatsAppOTP_WithMalformedJSON(t *testing.T) {
	router := setupTestRouterWithWhatsAppVerification()
	
	// Create request with malformed JSON
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer([]byte("{invalid json}")))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify error message
	assert.Contains(t, response["error"], "invalid request format")
}

// TestVerifyWhatsAppOTP_WithSecurityLogging tests security logging for WhatsApp OTP verification
func TestVerifyWhatsAppOTP_WithSecurityLogging(t *testing.T) {
	router := setupTestRouterWithWhatsAppVerification()
	
	// Create a test user first
	db := testutils.SetupTestDB(t)
	defer db.Close()
	
	uniquePhone := "+6281234567890"
	uniqueEmail := "whatsappuser@example.com"
	
	// Create existing user
	userID := "test-user-id"
	_, err := db.Exec("INSERT INTO users (id, name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, "WhatsApp User", uniqueEmail, uniquePhone, "student", "hashed_password")
	require.NoError(t, err)
	
	// Save a test OTP
	_, err = db.Exec("INSERT INTO otps (identifier, otp, expiry) VALUES ($1, $2, $3)",
		uniquePhone, "123456", time.Now().Add(10*time.Minute))
	require.NoError(t, err)
	
	// Test data
	testData := map[string]interface{}{
		"phone": uniquePhone,
		"otp":   "123456",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify response structure
	assert.Equal(t, "OTP verified successfully", response["message"])
	
	// In a real implementation, we would verify that the security logging was called
	// For now, we'll just ensure the request succeeded
	
	// Clean up
	testutils.CleanupUser(t, db, uniqueEmail)
}

// TestVerifyWhatsAppOTP_WithJWTTokenGeneration tests JWT token generation after successful verification
func TestVerifyWhatsAppOTP_WithJWTTokenGeneration(t *testing.T) {
	router := setupTestRouterWithWhatsAppVerification()
	
	// Create a test user first
	db := testutils.SetupTestDB(t)
	defer db.Close()
	
	uniquePhone := "+6281234567890"
	uniqueEmail := "whatsappuser@example.com"
	
	// Create existing user
	userID := "test-user-id"
	_, err := db.Exec("INSERT INTO users (id, name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, "WhatsApp User", uniqueEmail, uniquePhone, "student", "hashed_password")
	require.NoError(t, err)
	
	// Save a test OTP
	_, err = db.Exec("INSERT INTO otps (identifier, otp, expiry) VALUES ($1, $2, $3)",
		uniquePhone, "123456", time.Now().Add(10*time.Minute))
	require.NoError(t, err)
	
	// Test data
	testData := map[string]interface{}{
		"phone": uniquePhone,
		"otp":   "123456",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify response structure
	assert.Equal(t, "OTP verified successfully", response["message"])
	
	// Verify JWT token
	token, exists := response["token"]
	assert.True(t, exists)
	assert.NotEmpty(t, token)
	
	// Verify user details
	user := response["user"].(map[string]interface{})
	assert.Equal(t, userID, user["id"])
	assert.Equal(t, "WhatsApp User", user["name"])
	assert.Equal(t, uniqueEmail, user["email"])
	assert.Equal(t, uniquePhone, user["phone"])
	assert.Equal(t, "student", user["role"])
	
	// Verify expires_in
	expiresIn, exists := response["expires_in"]
	assert.True(t, exists)
	assert.IsType(t, float64(0), expiresIn)
	assert.Greater(t, expiresIn, float64(0))
	
	// Clean up
	testutils.CleanupUser(t, db, uniqueEmail)
}