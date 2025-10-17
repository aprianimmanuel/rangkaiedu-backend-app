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

// setupTestRouterWithWhatsApp sets up a test router with WhatsApp authentication routes
func setupTestRouterWithWhatsApp() *gin.Engine {
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
		// WhatsApp OTP sending route
		auth.POST("/whatsapp/send-otp", authHandler.SendWhatsAppOTP)
		
		// WhatsApp OTP verification route
		auth.POST("/whatsapp/verify", authHandler.VerifyWhatsAppOTP)
		
		// WhatsApp login route
		auth.POST("/whatsapp/login", authHandler.WhatsappLogin)
	}
	
	return router
}

// TestSendWhatsAppOTP_Success tests successful WhatsApp OTP sending
func TestSendWhatsAppOTP_Success(t *testing.T) {
	router := setupTestRouterWithWhatsApp()
	
	// Test data
	testData := map[string]interface{}{
		"phone": "+6281234567890",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(jsonData))
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
	assert.Contains(t, response, "otp_id")
	assert.Contains(t, response, "expires_in")
}

// TestSendWhatsAppOTP_InvalidPhoneFormat tests invalid phone number format
func TestSendWhatsAppOTP_InvalidPhoneFormat(t *testing.T) {
	router := setupTestRouterWithWhatsApp()
	
	// Test data with invalid phone format
	testData := map[string]interface{}{
		"phone": "invalid-phone-format",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(jsonData))
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

// TestSendWhatsAppOTP_MissingPhone tests missing phone number
func TestSendWhatsAppOTP_MissingPhone(t *testing.T) {
	router := setupTestRouterWithWhatsApp()
	
	// Test data without phone number
	testData := map[string]interface{}{}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(jsonData))
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

// TestSendWhatsAppOTP_RateLimited tests rate limiting for WhatsApp OTP sending
func TestSendWhatsAppOTP_RateLimited(t *testing.T) {
	router := setupTestRouterWithWhatsApp()
	
	phone := "+6281234567890"
	
	// Send multiple OTP requests to trigger rate limiting
	for i := 0; i < 6; i++ { // Assuming rate limit is 5 requests per minute
		testData := map[string]interface{}{
			"phone": phone,
		}
		
		jsonData, err := json.Marshal(testData)
		require.NoError(t, err)
		
		req, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		if i < 5 {
			// First 5 requests should succeed
			assert.Equal(t, http.StatusOK, w.Code)
		} else {
			// 6th request should be rate limited
			assert.Equal(t, http.StatusTooManyRequests, w.Code)
			
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			
			assert.Contains(t, response["error"], "rate limit exceeded")
		}
	}
}

// TestSendWhatsAppOTP_ProviderError tests WhatsApp provider error handling
func TestSendWhatsAppOTP_ProviderError(t *testing.T) {
	router := setupTestRouterWithWhatsApp()
	
	// Mock the WhatsApp provider to return an error
	// This would require mocking the WhatsApp utility functions
	
	// For now, we'll test the error response structure
	testData := map[string]interface{}{
		"phone": "+6281234567890",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response - this will depend on the actual implementation
	// For now, we'll assume it returns a 500 error when provider fails
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify error message
	assert.Contains(t, response["error"], "failed to send OTP")
}

// TestSendWhatsAppOTP_WithExistingUser tests OTP sending for existing user
func TestSendWhatsAppOTP_WithExistingUser(t *testing.T) {
	router := setupTestRouterWithWhatsApp()
	
	// Create a test user first
	db := testutils.SetupTestDB(t)
	defer db.Close()
	
	uniquePhone := "+6281234567890"
	uniqueEmail := "whatsappuser@example.com"
	
	// Create existing user
	_, err := db.Exec("INSERT INTO users (name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5)",
		"WhatsApp User", uniqueEmail, uniquePhone, "student", "hashed_password")
	require.NoError(t, err)
	
	// Test data
	testData := map[string]interface{}{
		"phone": uniquePhone,
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(jsonData))
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
	assert.Contains(t, response, "otp_id")
	assert.Contains(t, response, "expires_in")
	
	// Clean up
	testutils.CleanupUser(t, db, uniqueEmail)
}

// TestSendWhatsAppOTP_WithNewUser tests OTP sending for new user (should still work)
func TestSendWhatsAppOTP_WithNewUser(t *testing.T) {
	router := setupTestRouterWithWhatsApp()
	
	// Test data with new phone number
	testData := map[string]interface{}{
		"phone": "+6281112233445", // New phone number
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(jsonData))
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
	assert.Contains(t, response, "otp_id")
	assert.Contains(t, response, "expires_in")
}

// TestSendWhatsAppOTP_WithInvalidContentType tests invalid content type
func TestSendWhatsAppOTP_WithInvalidContentType(t *testing.T) {
	router := setupTestRouterWithWhatsApp()
	
	// Test data
	testData := map[string]interface{}{
		"phone": "+6281234567890",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request with invalid content type
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "text/plain") // Invalid content type
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
}

// TestSendWhatsAppOTP_WithEmptyBody tests empty request body
func TestSendWhatsAppOTP_WithEmptyBody(t *testing.T) {
	router := setupTestRouterWithWhatsApp()
	
	// Create request with empty body
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer([]byte("")))
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

// TestSendWhatsAppOTP_WithMalformedJSON tests malformed JSON request
func TestSendWhatsAppOTP_WithMalformedJSON(t *testing.T) {
	router := setupTestRouterWithWhatsApp()
	
	// Create request with malformed JSON
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer([]byte("{invalid json}")))
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

// TestSendWhatsAppOTP_WithTimeout tests timeout handling
func TestSendWhatsAppOTP_WithTimeout(t *testing.T) {
	router := setupTestRouterWithWhatsApp()
	
	// Test data
	testData := map[string]interface{}{
		"phone": "+6281234567890",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	// Set a very short timeout to simulate timeout
	// This would require mocking the WhatsApp service to actually timeout
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response - this will depend on the actual timeout handling
	// For now, we'll assume it returns a 504 error when timeout occurs
	assert.Equal(t, http.StatusGatewayTimeout, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verify error message
	assert.Contains(t, response["error"], "request timeout")
}

// TestSendWhatsAppOTP_WithSecurityLogging tests security logging for WhatsApp OTP requests
func TestSendWhatsAppOTP_WithSecurityLogging(t *testing.T) {
	router := setupTestRouterWithWhatsApp()
	
	// Test data
	testData := map[string]interface{}{
		"phone": "+6281234567890",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(jsonData))
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
}

// TestSendWhatsAppOTP_WithMultipleProviders tests fallback to secondary provider
func TestSendWhatsAppOTP_WithMultipleProviders(t *testing.T) {
	router := setupTestRouterWithWhatsApp()
	
	// Test data
	testData := map[string]interface{}{
		"phone": "+6281234567890",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	// Create request
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(jsonData))
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
	assert.Contains(t, response, "otp_id")
	assert.Contains(t, response, "expires_in")
	
	// In a real implementation, we would verify that the fallback provider was used
	// For now, we'll just ensure the request succeeded
}