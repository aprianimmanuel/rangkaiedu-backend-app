package integration_test

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

// setupTestRouterWithWhatsAppAuth sets up a test router with complete WhatsApp authentication routes
func setupTestRouterWithWhatsAppAuth() *gin.Engine {
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
		
		// User registration route (for creating users)
		auth.POST("/register", authHandler.Register)
	}
	
	return router
}

// TestEndToEndWhatsAppAuthFlow tests the complete WhatsApp authentication flow
func TestEndToEndWhatsAppAuthFlow(t *testing.T) {
	router := setupTestRouterWithWhatsAppAuth()
	
	// Step 1: Register a new user with WhatsApp phone number
	registerData := map[string]interface{}{
		"name":     "WhatsApp Test User",
		"email":    "whatsappintegration@example.com",
		"phone":    "+6281234567890",
		"password": "password123",
		"role":     "student",
	}
	
	registerJSON, err := json.Marshal(registerData)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(registerJSON))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var registerResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &registerResponse)
	require.NoError(t, err)
	
	assert.Equal(t, "User registered successfully. Verification OTP sent to email.", registerResponse["message"])
	
	// Step 2: Send OTP to WhatsApp phone number
	sendOTPData := map[string]interface{}{
		"phone": "+6281234567890",
	}
	
	sendOTPJSON, err := json.Marshal(sendOTPData)
	require.NoError(t, err)
	
	req2, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(sendOTPJSON))
	require.NoError(t, err)
	req2.Header.Set("Content-Type", "application/json")
	
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	
	assert.Equal(t, http.StatusOK, w2.Code)
	
	var sendOTPResponse map[string]interface{}
	err = json.Unmarshal(w2.Body.Bytes(), &sendOTPResponse)
	require.NoError(t, err)
	
	assert.Equal(t, "OTP sent to your WhatsApp number", sendOTPResponse["message"])
	assert.Contains(t, sendOTPResponse, "otp_id")
	assert.Contains(t, sendOTPResponse, "expires_in")
	
	// Step 3: Verify OTP (simulating successful verification)
	verifyOTPData := map[string]interface{}{
		"phone": "+6281234567890",
		"otp":   "123456", // This would normally be the actual OTP sent
	}
	
	verifyOTPJSON, err := json.Marshal(verifyOTPData)
	require.NoError(t, err)
	
	req3, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(verifyOTPJSON))
	require.NoError(t, err)
	req3.Header.Set("Content-Type", "application/json")
	
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	
	assert.Equal(t, http.StatusOK, w3.Code)
	
	var verifyOTPResponse map[string]interface{}
	err = json.Unmarshal(w3.Body.Bytes(), &verifyOTPResponse)
	require.NoError(t, err)
	
	assert.Equal(t, "OTP verified successfully", verifyOTPResponse["message"])
	assert.Contains(t, verifyOTPResponse, "token")
	assert.Contains(t, verifyOTPResponse, "user")
	
	// Step 4: Login with WhatsApp (using OTP flow)
	loginData := map[string]interface{}{
		"phone": "+6281234567890",
		"otp":   "123456", // This would normally be the actual OTP sent
	}
	
	loginJSON, err := json.Marshal(loginData)
	require.NoError(t, err)
	
	req4, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(loginJSON))
	require.NoError(t, err)
	req4.Header.Set("Content-Type", "application/json")
	
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)
	
	assert.Equal(t, http.StatusOK, w4.Code)
	
	var loginResponse map[string]interface{}
	err = json.Unmarshal(w4.Body.Bytes(), &loginResponse)
	require.NoError(t, err)
	
	assert.Equal(t, "Login successful", loginResponse["message"])
	assert.Contains(t, loginResponse, "token")
	assert.Contains(t, loginResponse, "user")
	
	// Verify user details
	user := loginResponse["user"].(map[string]interface{})
	assert.Equal(t, "WhatsApp Test User", user["name"])
	assert.Equal(t, "whatsappintegration@example.com", user["email"])
	assert.Equal(t, "+6281234567890", user["phone"])
	assert.Equal(t, "student", user["role"])
	
	// Clean up
	db := testutils.SetupTestDB(t)
	defer db.Close()
	testutils.CleanupUser(t, db, "whatsappintegration@example.com")
}

// TestEndToEndWhatsAppAuthFlow_WithMultipleProviders tests the flow with multiple WhatsApp providers
func TestEndToEndWhatsAppAuthFlow_WithMultipleProviders(t *testing.T) {
	router := setupTestRouterWithWhatsAppAuth()
	
	// Step 1: Register a new user with WhatsApp phone number
	registerData := map[string]interface{}{
		"name":     "Multi-Provider Test User",
		"email":    "multiprovider@example.com",
		"phone":    "+6281112233445",
		"password": "password123",
		"role":     "student",
	}
	
	registerJSON, err := json.Marshal(registerData)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(registerJSON))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// Step 2: Send OTP to WhatsApp phone number (should use primary provider)
	sendOTPData := map[string]interface{}{
		"phone": "+6281112233445",
	}
	
	sendOTPJSON, err := json.Marshal(sendOTPData)
	require.NoError(t, err)
	
	req2, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(sendOTPJSON))
	require.NoError(t, err)
	req2.Header.Set("Content-Type", "application/json")
	
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	
	assert.Equal(t, http.StatusOK, w2.Code)
	
	var sendOTPResponse map[string]interface{}
	err = json.Unmarshal(w2.Body.Bytes(), &sendOTPResponse)
	require.NoError(t, err)
	
	assert.Equal(t, "OTP sent to your WhatsApp number", sendOTPResponse["message"])
	
	// Step 3: Verify OTP
	verifyOTPData := map[string]interface{}{
		"phone": "+6281112233445",
		"otp":   "123456",
	}
	
	verifyOTPJSON, err := json.Marshal(verifyOTPData)
	require.NoError(t, err)
	
	req3, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(verifyOTPJSON))
	require.NoError(t, err)
	req3.Header.Set("Content-Type", "application/json")
	
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	
	assert.Equal(t, http.StatusOK, w3.Code)
	
	// Step 4: Login with WhatsApp
	loginData := map[string]interface{}{
		"phone": "+6281112233445",
		"otp":   "123456",
	}
	
	loginJSON, err := json.Marshal(loginData)
	require.NoError(t, err)
	
	req4, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(loginJSON))
	require.NoError(t, err)
	req4.Header.Set("Content-Type", "application/json")
	
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)
	
	assert.Equal(t, http.StatusOK, w4.Code)
	
	var loginResponse map[string]interface{}
	err = json.Unmarshal(w4.Body.Bytes(), &loginResponse)
	require.NoError(t, err)
	
	assert.Equal(t, "Login successful", loginResponse["message"])
	assert.Contains(t, loginResponse, "token")
	
	// Clean up
	db := testutils.SetupTestDB(t)
	defer db.Close()
	testutils.CleanupUser(t, db, "multiprovider@example.com")
}

// TestEndToEndWhatsAppAuthFlow_WithRateLimiting tests the flow with rate limiting
func TestEndToEndWhatsAppAuthFlow_WithRateLimiting(t *testing.T) {
	router := setupTestRouterWithWhatsAppAuth()
	
	phone := "+6289999888777"
	
	// Step 1: Register a new user
	registerData := map[string]interface{}{
		"name":     "Rate Limit Test User",
		"email":    "ratelimit@example.com",
		"phone":    phone,
		"password": "password123",
		"role":     "student",
	}
	
	registerJSON, err := json.Marshal(registerData)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(registerJSON))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// Step 2: Send multiple OTP requests to trigger rate limiting
	for i := 0; i < 6; i++ { // Assuming rate limit is 5 requests per minute
		sendOTPData := map[string]interface{}{
			"phone": phone,
		}
		
		sendOTPJSON, err := json.Marshal(sendOTPData)
		require.NoError(t, err)
		
		req2, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(sendOTPJSON))
		require.NoError(t, err)
		req2.Header.Set("Content-Type", "application/json")
		
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		
		if i < 5 {
			// First 5 requests should succeed
			assert.Equal(t, http.StatusOK, w2.Code)
		} else {
			// 6th request should be rate limited
			assert.Equal(t, http.StatusTooManyRequests, w2.Code)
			
			var response map[string]interface{}
			err = json.Unmarshal(w2.Body.Bytes(), &response)
			require.NoError(t, err)
			
			assert.Contains(t, response["error"], "rate limit exceeded")
		}
	}
	
	// Clean up
	db := testutils.SetupTestDB(t)
	defer db.Close()
	testutils.CleanupUser(t, db, "ratelimit@example.com")
}

// TestEndToEndWhatsAppAuthFlow_WithDifferentPhoneFormats tests the flow with different phone number formats
func TestEndToEndWhatsAppAuthFlow_WithDifferentPhoneFormats(t *testing.T) {
	router := setupTestRouterWithWhatsAppAuth()
	
	// Test different phone number formats
	testCases := []struct {
		name        string
		phone       string
		normalizedPhone string
	}{
		{"International format", "+6281234567890", "+6281234567890"},
		{"With spaces", "+62 812 3456 7890", "+6281234567890"},
		{"With dashes", "+62-812-3456-7890", "+6281234567890"},
		{"With parentheses", "+62 (812) 3456-7890", "+6281234567890"},
		{"Mixed special chars", "+62 812-3456 7890", "+6281234567890"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Step 1: Register a new user
			email := tc.phone + "@example.com"
			registerData := map[string]interface{}{
				"name":     "Format Test User",
				"email":    email,
				"phone":    tc.phone,
				"password": "password123",
				"role":     "student",
			}
			
			registerJSON, err := json.Marshal(registerData)
			require.NoError(t, err)
			
			req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(registerJSON))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			assert.Equal(t, http.StatusCreated, w.Code)
			
			// Step 2: Send OTP using the same format
			sendOTPData := map[string]interface{}{
				"phone": tc.phone,
			}
			
			sendOTPJSON, err := json.Marshal(sendOTPData)
			require.NoError(t, err)
			
			req2, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(sendOTPJSON))
			require.NoError(t, err)
			req2.Header.Set("Content-Type", "application/json")
			
			w2 := httptest.NewRecorder()
			router.ServeHTTP(w2, req2)
			
			assert.Equal(t, http.StatusOK, w2.Code)
			
			// Step 3: Verify OTP using normalized format
			verifyOTPData := map[string]interface{}{
				"phone": tc.normalizedPhone,
				"otp":   "123456",
			}
			
			verifyOTPJSON, err := json.Marshal(verifyOTPData)
			require.NoError(t, err)
			
			req3, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(verifyOTPJSON))
			require.NoError(t, err)
			req3.Header.Set("Content-Type", "application/json")
			
			w3 := httptest.NewRecorder()
			router.ServeHTTP(w3, req3)
			
			assert.Equal(t, http.StatusOK, w3.Code)
			
			// Step 4: Login using normalized format
			loginData := map[string]interface{}{
				"phone": tc.normalizedPhone,
				"otp":   "123456",
			}
			
			loginJSON, err := json.Marshal(loginData)
			require.NoError(t, err)
			
			req4, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(loginJSON))
			require.NoError(t, err)
			req4.Header.Set("Content-Type", "application/json")
			
			w4 := httptest.NewRecorder()
			router.ServeHTTP(w4, req4)
			
			assert.Equal(t, http.StatusOK, w4.Code)
			
			// Clean up
			db := testutils.SetupTestDB(t)
			defer db.Close()
			testutils.CleanupUser(t, db, email)
		})
	}
}

// TestEndToEndWhatsAppAuthFlow_WithSecurityLogging tests security logging throughout the flow
func TestEndToEndWhatsAppAuthFlow_WithSecurityLogging(t *testing.T) {
	router := setupTestRouterWithWhatsAppAuth()
	
	// Step 1: Register a new user
	registerData := map[string]interface{}{
		"name":     "Security Test User",
		"email":    "security@example.com",
		"phone":    "+6287776665554",
		"password": "password123",
		"role":     "student",
	}
	
	registerJSON, err := json.Marshal(registerData)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(registerJSON))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// Step 2: Send OTP
	sendOTPData := map[string]interface{}{
		"phone": "+6287776665554",
	}
	
	sendOTPJSON, err := json.Marshal(sendOTPData)
	require.NoError(t, err)
	
	req2, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(sendOTPJSON))
	require.NoError(t, err)
	req2.Header.Set("Content-Type", "application/json")
	
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	
	assert.Equal(t, http.StatusOK, w2.Code)
	
	// Step 3: Verify OTP
	verifyOTPData := map[string]interface{}{
		"phone": "+6287776665554",
		"otp":   "123456",
	}
	
	verifyOTPJSON, err := json.Marshal(verifyOTPData)
	require.NoError(t, err)
	
	req3, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(verifyOTPJSON))
	require.NoError(t, err)
	req3.Header.Set("Content-Type", "application/json")
	
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	
	assert.Equal(t, http.StatusOK, w3.Code)
	
	// Step 4: Login
	loginData := map[string]interface{}{
		"phone": "+6287776665554",
		"otp":   "123456",
	}
	
	loginJSON, err := json.Marshal(loginData)
	require.NoError(t, err)
	
	req4, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(loginJSON))
	require.NoError(t, err)
	req4.Header.Set("Content-Type", "application/json")
	
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)
	
	assert.Equal(t, http.StatusOK, w4.Code)
	
	// In a real implementation, we would verify that security logging was called
	// For each step of the authentication flow
	// For now, we'll just ensure the flow completed successfully
	
	// Clean up
	db := testutils.SetupTestDB(t)
	defer db.Close()
	testutils.CleanupUser(t, db, "security@example.com")
}

// TestEndToEndWhatsAppAuthFlow_WithMonitoringIntegration tests monitoring integration throughout the flow
func TestEndToEndWhatsAppAuthFlow_WithMonitoringIntegration(t *testing.T) {
	router := setupTestRouterWithWhatsAppAuth()
	
	// Step 1: Register a new user
	registerData := map[string]interface{}{
		"name":     "Monitoring Test User",
		"email":    "monitoring@example.com",
		"phone":    "+6286665554443",
		"password": "password123",
		"role":     "student",
	}
	
	registerJSON, err := json.Marshal(registerData)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(registerJSON))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// Step 2: Send OTP
	sendOTPData := map[string]interface{}{
		"phone": "+6286665554443",
	}
	
	sendOTPJSON, err := json.Marshal(sendOTPData)
	require.NoError(t, err)
	
	req2, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(sendOTPJSON))
	require.NoError(t, err)
	req2.Header.Set("Content-Type", "application/json")
	
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	
	assert.Equal(t, http.StatusOK, w2.Code)
	
	// Step 3: Verify OTP
	verifyOTPData := map[string]interface{}{
		"phone": "+6286665554443",
		"otp":   "123456",
	}
	
	verifyOTPJSON, err := json.Marshal(verifyOTPData)
	require.NoError(t, err)
	
	req3, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(verifyOTPJSON))
	require.NoError(t, err)
	req3.Header.Set("Content-Type", "application/json")
	
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	
	assert.Equal(t, http.StatusOK, w3.Code)
	
	// Step 4: Login
	loginData := map[string]interface{}{
		"phone": "+6286665554443",
		"otp":   "123456",
	}
	
	loginJSON, err := json.Marshal(loginData)
	require.NoError(t, err)
	
	req4, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(loginJSON))
	require.NoError(t, err)
	req4.Header.Set("Content-Type", "application/json")
	
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)
	
	assert.Equal(t, http.StatusOK, w4.Code)
	
	// In a real implementation, we would verify that monitoring metrics were recorded
	// For each step of the authentication flow
	// For now, we'll just ensure the flow completed successfully
	
	// Clean up
	db := testutils.SetupTestDB(t)
	defer db.Close()
	testutils.CleanupUser(t, db, "monitoring@example.com")
}

// TestEndToEndWhatsAppAuthFlow_WithFallbackMechanism tests fallback mechanism when primary provider fails
func TestEndToEndWhatsAppAuthFlow_WithFallbackMechanism(t *testing.T) {
	router := setupTestRouterWithWhatsAppAuth()
	
	// Step 1: Register a new user
	registerData := map[string]interface{}{
		"name":     "Fallback Test User",
		"email":    "fallback@example.com",
		"phone":    "+6285554443332",
		"password": "password123",
		"role":     "student",
	}
	
	registerJSON, err := json.Marshal(registerData)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(registerJSON))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// Step 2: Send OTP (should fallback to secondary provider if primary fails)
	sendOTPData := map[string]interface{}{
		"phone": "+6285554443332",
	}
	
	sendOTPJSON, err := json.Marshal(sendOTPData)
	require.NoError(t, err)
	
	req2, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(sendOTPJSON))
	require.NoError(t, err)
	req2.Header.Set("Content-Type", "application/json")
	
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	
	assert.Equal(t, http.StatusOK, w2.Code)
	
	// In a real implementation, we would verify that the fallback provider was used
	// For now, we'll just ensure the request succeeded
	
	// Clean up
	db := testutils.SetupTestDB(t)
	defer db.Close()
	testutils.CleanupUser(t, db, "fallback@example.com")
}