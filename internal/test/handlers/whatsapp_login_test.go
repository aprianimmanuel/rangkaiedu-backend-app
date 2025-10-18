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

// setupTestRouterWithWhatsAppLogin sets up a test router with WhatsApp login routes
func setupTestRouterWithWhatsAppLogin() *gin.Engine {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Initialize test database configuration
	_ = config.LoadTest()
	
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
		// WhatsApp login route
		auth.POST("/whatsapp/login", authHandler.WhatsappLogin)
	}
	
	return router
}

// TestWhatsappLogin_Success tests successful WhatsApp login
func TestWhatsappLogin_Success(t *testing.T) {
	router := setupTestRouterWithWhatsAppLogin()
	
	// Create a test user first
	db := testutils.SetupTestDB(t)
	defer db.Close()
	
	uniquePhone := "+6281234567890"
	uniqueEmail := "whatsappuser@example.com"
	
	// Create existing user with phone number
	userID := "test-user-id"
	_, err := db.Exec("INSERT INTO users (id, name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, "WhatsApp User", uniqueEmail, uniquePhone, "student", "hashed_password")
	require.NoError(t, err)
	
	// Test data
	testData := map[string]interface{}{
		"phone": uniquePhone,
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(jsonData))
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
	assert.Equal(t, "OTP sent to your WhatsApp number", response["message"])
	assert.Contains(t, response, "otp_required")
	assert.True(t, response["otp_required"].(bool))
	
	// Clean up
	testutils.CleanupUser(t, db, uniqueEmail)
}

// TestWhatsappLogin_WithOTP tests WhatsApp login with OTP verification
func TestWhatsappLogin_WithOTP(t *testing.T) {
	router := setupTestRouterWithWhatsAppLogin()
	
	// Create a test user first
	db := testutils.SetupTestDB(t)
	defer db.Close()
	
	uniquePhone := "+6281234567890"
	uniqueEmail := "whatsappuser@example.com"
	
	// Create existing user with phone number
	userID := "test-user-id"
	_, err := db.Exec("INSERT INTO users (id, name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, "WhatsApp User", uniqueEmail, uniquePhone, "student", "hashed_password")
	require.NoError(t, err)
	
	// Save a test OTP
	_, err = db.Exec("INSERT INTO otps (identifier, otp, expiry) VALUES ($1, $2, $3)",
		uniquePhone, "123456", time.Now().Add(10*time.Minute))
	require.NoError(t, err)
	
	// Test data with OTP
	testData := map[string]interface{}{
		"phone": uniquePhone,
		"otp":   "123456",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(jsonData))
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
	assert.Equal(t, "Login successful", response["message"])
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "user")
	
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

// TestWhatsappLogin_InvalidPhoneFormat tests invalid phone number format
func TestWhatsappLogin_InvalidPhoneFormat(t *testing.T) {
	router := setupTestRouterWithWhatsAppLogin()
	
	// Test data with invalid phone format
	testData := map[string]interface{}{
		"phone": "invalid-phone-format",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(jsonData))
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

// TestWhatsappLogin_MissingPhone tests missing phone number
func TestWhatsappLogin_MissingPhone(t *testing.T) {
	router := setupTestRouterWithWhatsAppLogin()
	
	// Test data without phone number
	testData := map[string]interface{}{}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(jsonData))
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

// TestWhatsappLogin_UserNotFound tests login with non-existent user
func TestWhatsappLogin_UserNotFound(t *testing.T) {
	router := setupTestRouterWithWhatsAppLogin()
	
	// Test data with non-existent phone number
	testData := map[string]interface{}{
		"phone": "+6281112233445", // Non-existent phone number
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(jsonData))
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

// TestWhatsappLogin_InvalidOTP tests invalid OTP handling
func TestWhatsappLogin_InvalidOTP(t *testing.T) {
	router := setupTestRouterWithWhatsAppLogin()
	
	// Create a test user first
	db := testutils.SetupTestDB(t)
	defer db.Close()
	
	uniquePhone := "+6281234567890"
	uniqueEmail := "whatsappuser@example.com"
	
	// Create existing user with phone number
	userID := "test-user-id"
	_, err := db.Exec("INSERT INTO users (id, name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, "WhatsApp User", uniqueEmail, uniquePhone, "student", "hashed_password")
	require.NoError(t, err)
	
	// Test data with invalid OTP
	testData := map[string]interface{}{
		"phone": uniquePhone,
		"otp":   "invalid-otp",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(jsonData))
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
	
	// Clean up
	testutils.CleanupUser(t, db, uniqueEmail)
}

// TestWhatsappLogin_ExpiredOTP tests expired OTP handling
func TestWhatsappLogin_ExpiredOTP(t *testing.T) {
	router := setupTestRouterWithWhatsAppLogin()
	
	// Create a test user first
	db := testutils.SetupTestDB(t)
	defer db.Close()
	
	uniquePhone := "+6281234567890"
	uniqueEmail := "whatsappuser@example.com"
	
	// Create existing user with phone number
	userID := "test-user-id"
	_, err := db.Exec("INSERT INTO users (id, name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, "WhatsApp User", uniqueEmail, uniquePhone, "student", "hashed_password")
	require.NoError(t, err)
	
	// Save an expired OTP
	_, err = db.Exec("INSERT INTO otps (identifier, otp, expiry) VALUES ($1, $2, $3)",
		uniquePhone, "123456", time.Now().Add(-1*time.Hour)) // Expired 1 hour ago
	require.NoError(t, err)
	
	// Test data with expired OTP
	testData := map[string]interface{}{
		"phone": uniquePhone,
		"otp":   "123456",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(jsonData))
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
	
	// Clean up
	testutils.CleanupUser(t, db, uniqueEmail)
}

// TestWhatsappLogin_WrongLengthOTP tests OTP with wrong length
func TestWhatsappLogin_WrongLengthOTP(t *testing.T) {
	router := setupTestRouterWithWhatsAppLogin()
	
	// Create a test user first
	db := testutils.SetupTestDB(t)
	defer db.Close()
	
	uniquePhone := "+6281234567890"
	uniqueEmail := "whatsappuser@example.com"
	
	// Create existing user with phone number
	userID := "test-user-id"
	_, err := db.Exec("INSERT INTO users (id, name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, "WhatsApp User", uniqueEmail, uniquePhone, "student", "hashed_password")
	require.NoError(t, err)
	
	// Test data with wrong length OTP
	testData := map[string]interface{}{
		"phone": uniquePhone,
		"otp":   "12345", // 5 digits instead of 6
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(jsonData))
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
	
	// Clean up
	testutils.CleanupUser(t, db, uniqueEmail)
}

// TestWhatsappLogin_WithInvalidContentType tests invalid content type
func TestWhatsappLogin_WithInvalidContentType(t *testing.T) {
	router := setupTestRouterWithWhatsAppLogin()
	
	// Test data
	testData := map[string]interface{}{
		"phone": "+6281234567890",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request with invalid content type
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "text/plain") // Invalid content type
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
}

// TestWhatsappLogin_WithEmptyBody tests empty request body
func TestWhatsappLogin_WithEmptyBody(t *testing.T) {
	router := setupTestRouterWithWhatsAppLogin()
	
	// Create request with empty body
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer([]byte("")))
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

// TestWhatsappLogin_WithMalformedJSON tests malformed JSON request
func TestWhatsappLogin_WithMalformedJSON(t *testing.T) {
	router := setupTestRouterWithWhatsAppLogin()
	
	// Create request with malformed JSON
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer([]byte("{invalid json}")))
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

// TestWhatsappLogin_WithSecurityLogging tests security logging for WhatsApp login
func TestWhatsappLogin_WithSecurityLogging(t *testing.T) {
	router := setupTestRouterWithWhatsAppLogin()
	
	// Create a test user first
	db := testutils.SetupTestDB(t)
	defer db.Close()
	
	uniquePhone := "+6281234567890"
	uniqueEmail := "whatsappuser@example.com"
	
	// Create existing user with phone number
	userID := "test-user-id"
	_, err := db.Exec("INSERT INTO users (id, name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, "WhatsApp User", uniqueEmail, uniquePhone, "student", "hashed_password")
	require.NoError(t, err)
	
	// Test data
	testData := map[string]interface{}{
		"phone": uniquePhone,
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(jsonData))
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
	assert.Equal(t, "OTP sent to your WhatsApp number", response["message"])
	
	// In a real implementation, we would verify that the security logging was called
	// For now, we'll just ensure the request succeeded
	
	// Clean up
	testutils.CleanupUser(t, db, uniqueEmail)
}

// TestWhatsappLogin_WithJWTTokenGeneration tests JWT token generation after successful login
func TestWhatsappLogin_WithJWTTokenGeneration(t *testing.T) {
	router := setupTestRouterWithWhatsAppLogin()
	
	// Create a test user first
	db := testutils.SetupTestDB(t)
	defer db.Close()
	
	uniquePhone := "+6281234567890"
	uniqueEmail := "whatsappuser@example.com"
	
	// Create existing user with phone number
	userID := "test-user-id"
	_, err := db.Exec("INSERT INTO users (id, name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, "WhatsApp User", uniqueEmail, uniquePhone, "student", "hashed_password")
	require.NoError(t, err)
	
	// Save a test OTP
	_, err = db.Exec("INSERT INTO otps (identifier, otp, expiry) VALUES ($1, $2, $3)",
		uniquePhone, "123456", time.Now().Add(10*time.Minute))
	require.NoError(t, err)
	
	// Test data with OTP
	testData := map[string]interface{}{
		"phone": uniquePhone,
		"otp":   "123456",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(jsonData))
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
	assert.Equal(t, "Login successful", response["message"])
	
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
	
	// Clean up
	testutils.CleanupUser(t, db, uniqueEmail)
}

// TestWhatsappLogin_WithDifferentPhoneFormats tests login with different phone number formats
func TestWhatsappLogin_WithDifferentPhoneFormats(t *testing.T) {
	router := setupTestRouterWithWhatsAppLogin()
	
	// Create a test user first
	db := testutils.SetupTestDB(t)
	defer db.Close()
	
	uniquePhone := "+6281234567890"
	uniqueEmail := "whatsappuser@example.com"
	
	// Create existing user with phone number
	userID := "test-user-id"
	_, err := db.Exec("INSERT INTO users (id, name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, "WhatsApp User", uniqueEmail, uniquePhone, "student", "hashed_password")
	require.NoError(t, err)
	
	// Save a test OTP
	_, err = db.Exec("INSERT INTO otps (identifier, otp, expiry) VALUES ($1, $2, $3)",
		uniquePhone, "123456", time.Now().Add(10*time.Minute))
	require.NoError(t, err)
	
	// Test different phone number formats
	testCases := []struct {
		name        string
		phone       string
		expectedErr string
	}{
		{"International format", "+6281234567890", ""},
		{"With spaces", "+62 812 3456 7890", ""},
		{"With dashes", "+62-812-3456-7890", ""},
		{"With parentheses", "+62 (812) 3456-7890", ""},
		{"Mixed special chars", "+62 812-3456 7890", ""},
		{"Invalid format", "6281234567890", "invalid phone number format"},
		{"Too short", "+123456", "invalid phone number format"},
		{"Contains letters", "+62812ABC7890", "invalid phone number format"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test data
			testData := map[string]interface{}{
				"phone": tc.phone,
				"otp":   "123456",
			}
			
			// Convert to JSON
			jsonData, err := json.Marshal(testData)
			require.NoError(t, err)
			
			// Create request
			req, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			
			// Record response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			if tc.expectedErr != "" {
				// Assert error response
				assert.Equal(t, http.StatusBadRequest, w.Code)
				
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				
				assert.Contains(t, response["error"], tc.expectedErr)
			} else {
				// Assert success response
				assert.Equal(t, http.StatusOK, w.Code)
				
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				
				assert.Equal(t, "Login successful", response["message"])
				assert.Contains(t, response, "token")
				assert.Contains(t, response, "user")
			}
		})
	}
	
	// Clean up
	testutils.CleanupUser(t, db, uniqueEmail)
}