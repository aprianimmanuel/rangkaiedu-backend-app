package testutils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/handlers"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/test/mocks"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/whatsapp"
)

// WhatsAppTestConfig holds configuration for WhatsApp tests
type WhatsAppTestConfig struct {
	TestDB         *sql.DB
	Router         *gin.Engine
	AuthHandler    *handlers.AuthHandler
	MockWhatsApp   *mocks.MockWhatsAppAPI
	TestUser       *TestUser
}

// TestUser represents a test user for WhatsApp authentication tests
type TestUser struct {
	ID       string
	Name     string
	Email    string
	Phone    string
	Role     string
	Password string
}

// NewWhatsAppTestConfig creates a new WhatsApp test configuration
func NewWhatsAppTestConfig(t *testing.T) *WhatsAppTestConfig {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Initialize test database configuration
	// cfg := config.LoadTest()
	
	// Initialize database connection pool for testing
	// if err := InitTestDB(); err != nil {
	// 	// If database initialization fails, continue with tests but log the error
	// 	// The mock handlers will work without database connection
	// }
	
	// Create a new router
	router := gin.New()
	
	// Create AuthHandler instance
	authHandler := handlers.NewAuthHandler()
	
	// Setup authentication routes
	auth := router.Group("/api/auth")
	{
		// WhatsApp OTP sending route
		auth.POST("/whatsapp/send-otp", authHandler.SendOTP)
		
		// WhatsApp OTP verification route
		auth.POST("/whatsapp/verify", authHandler.VerifyOTP)
		
		// WhatsApp login route
		auth.POST("/whatsapp/login", authHandler.LoginWithWhatsApp)
		
		// User registration route (for creating users)
		auth.POST("/register", authHandler.Register)
	}
	
	// Create mock WhatsApp API
	mockWhatsApp := mocks.NewMockWhatsAppAPI()
	
	// Create test user
	testUser := &TestUser{
		ID:       "test-whatsapp-user-id",
		Name:     "WhatsApp Test User",
		Email:    "whatsapptest@example.com",
		Phone:    "+6281234567890",
		Role:     "student",
		Password: "password123",
	}
	
	return &WhatsAppTestConfig{
		TestDB:       SetupTestDB(t),
		Router:       router,
		AuthHandler:  authHandler,
		MockWhatsApp: mockWhatsApp,
		TestUser:     testUser,
	}
}

// SetupTestUser creates a test user in the database
func (wtc *WhatsAppTestConfig) SetupTestUser(t *testing.T) {
	_, err := wtc.TestDB.Exec("INSERT INTO users (id, name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5, $6)",
		wtc.TestUser.ID, wtc.TestUser.Name, wtc.TestUser.Email, wtc.TestUser.Phone, wtc.TestUser.Role, "hashed_password")
	require.NoError(t, err)
}

// CleanupTestUser removes a test user from the database
func (wtc *WhatsAppTestConfig) CleanupTestUser(t *testing.T) {
	CleanupUser(t, wtc.TestDB, wtc.TestUser.Email)
}

// SaveOTP saves a test OTP to the database
func (wtc *WhatsAppTestConfig) SaveOTP(t *testing.T, identifier, otp string) {
	_, err := wtc.TestDB.Exec("INSERT INTO otps (identifier, otp, expiry) VALUES ($1, $2, $3)",
		identifier, otp, time.Now().Add(10*time.Minute))
	require.NoError(t, err)
}

// CleanupOTP removes a test OTP from the database
func (wtc *WhatsAppTestConfig) CleanupOTP(t *testing.T, identifier string) {
	_, err := wtc.TestDB.Exec("DELETE FROM otps WHERE identifier = $1", identifier)
	require.NoError(t, err)
}

// MakeWhatsAppOTPRequest sends a WhatsApp OTP request
func (wtc *WhatsAppTestConfig) MakeWhatsAppOTPRequest(t *testing.T, phone string) *httptest.ResponseRecorder {
	testData := map[string]interface{}{
		"phone": phone,
	}
	
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/send-otp", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	wtc.Router.ServeHTTP(w, req)
	
	return w
}

// MakeWhatsAppVerifyRequest sends a WhatsApp OTP verification request
func (wtc *WhatsAppTestConfig) MakeWhatsAppVerifyRequest(t *testing.T, phone, otp string) *httptest.ResponseRecorder {
	testData := map[string]interface{}{
		"phone": phone,
		"otp":   otp,
	}
	
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/verify", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	wtc.Router.ServeHTTP(w, req)
	
	return w
}

// MakeWhatsAppLoginRequest sends a WhatsApp login request
func (wtc *WhatsAppTestConfig) MakeWhatsAppLoginRequest(t *testing.T, phone, otp string) *httptest.ResponseRecorder {
	testData := map[string]interface{}{
		"phone": phone,
		"otp":   otp,
	}
	
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/auth/whatsapp/login", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	wtc.Router.ServeHTTP(w, req)
	
	return w
}

// MakeRegisterRequest sends a user registration request
func (wtc *WhatsAppTestConfig) MakeRegisterRequest(t *testing.T, userData map[string]interface{}) *httptest.ResponseRecorder {
	jsonData, err := json.Marshal(userData)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	wtc.Router.ServeHTTP(w, req)
	
	return w
}

// AssertSuccessfulOTPResponse asserts that a response indicates successful OTP sending
func (wtc *WhatsAppTestConfig) AssertSuccessfulOTPResponse(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "OTP sent to your WhatsApp number", response["message"])
	assert.Contains(t, response, "otp_id")
	assert.Contains(t, response, "expires_in")
}

// AssertSuccessfulVerifyResponse asserts that a response indicates successful OTP verification
func (wtc *WhatsAppTestConfig) AssertSuccessfulVerifyResponse(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "OTP verified successfully", response["message"])
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "user")
}

// AssertSuccessfulLoginResponse asserts that a response indicates successful login
func (wtc *WhatsAppTestConfig) AssertSuccessfulLoginResponse(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "Login successful", response["message"])
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "user")
}

// AssertErrorResponse asserts that a response contains an error
func (wtc *WhatsAppTestConfig) AssertErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedError string) {
	assert.Equal(t, expectedStatus, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Contains(t, response["error"], expectedError)
}

// GenerateTestOTP generates a test OTP
func (wtc *WhatsAppTestConfig) GenerateTestOTP() string {
	return "123456"
}

// GenerateTestPhoneNumbers generates various test phone number formats
func (wtc *WhatsAppTestConfig) GenerateTestPhoneNumbers() []string {
	return []string{
		"+6281234567890",      // International format
		"+62 812 3456 7890",   // With spaces
		"+62-812-3456-7890",   // With dashes
		"+62 (812) 3456-7890", // With parentheses
		"+62 812-3456 7890",   // Mixed special chars
	}
}

// GenerateInvalidPhoneNumbers generates various invalid phone number formats
func (wtc *WhatsAppTestConfig) GenerateInvalidPhoneNumbers() []string {
	return []string{
		"6281234567890",       // Missing plus
		"+123456",             // Too short
		"+62812ABC7890",       // Contains letters
		"+62812@#$7890",       // Contains special chars
		"",                    // Empty string
		"null",                // Null value
		"<script>alert('xss')</script>", // XSS attempt
	}
}

// ValidatePhoneFormat validates a phone number format using the WhatsApp utility
func (wtc *WhatsAppTestConfig) ValidatePhoneFormat(phone string) bool {
	return whatsapp.ValidatePhoneFormat(phone)
}

// FormatPhoneNumber formats a phone number using the WhatsApp utility
func (wtc *WhatsAppTestConfig) FormatPhoneNumber(phone string) string {
	return whatsapp.FormatPhoneNumber(phone)
}

// CreateTestUserWithPhone creates a test user with a specific phone number
func (wtc *WhatsAppTestConfig) CreateTestUserWithPhone(t *testing.T, phone string) *TestUser {
	email := fmt.Sprintf("%s@example.com", strings.ReplaceAll(phone, "+", ""))
	
	user := &TestUser{
		ID:       fmt.Sprintf("test-user-%s", phone),
		Name:     fmt.Sprintf("Test User %s", phone),
		Email:    email,
		Phone:    phone,
		Role:     "student",
		Password: "password123",
	}
	
	_, err := wtc.TestDB.Exec("INSERT INTO users (id, name, email, phone, role, password_hash) VALUES ($1, $2, $3, $4, $5, $6)",
		user.ID, user.Name, user.Email, user.Phone, user.Role, "hashed_password")
	require.NoError(t, err)
	
	return user
}

// CreateMultipleTestUsers creates multiple test users with different phone numbers
func (wtc *WhatsAppTestConfig) CreateMultipleTestUsers(t *testing.T, phoneNumbers []string) []*TestUser {
	var users []*TestUser
	
	for _, phone := range phoneNumbers {
		user := wtc.CreateTestUserWithPhone(t, phone)
		users = append(users, user)
	}
	
	return users
}

// CleanupMultipleTestUsers cleans up multiple test users
func (wtc *WhatsAppTestConfig) CleanupMultipleTestUsers(t *testing.T, users []*TestUser) {
	for _, user := range users {
		CleanupUser(t, wtc.TestDB, user.Email)
	}
}

// SimulateRateLimiting simulates rate limiting behavior
func (wtc *WhatsAppTestConfig) SimulateRateLimiting(t *testing.T, phone string, shouldRateLimit bool) {
	wtc.MockWhatsApp.SetShouldRateLimit(shouldRateLimit)
	
	// Send multiple requests to trigger rate limiting
	for i := 0; i < 6; i++ {
		w := wtc.MakeWhatsAppOTPRequest(t, phone)
		
		if shouldRateLimit && i >= 5 {
			wtc.AssertErrorResponse(t, w, http.StatusTooManyRequests, "rate limit exceeded")
		} else {
			wtc.AssertSuccessfulOTPResponse(t, w)
		}
	}
}

// SimulateAPIFailure simulates API failure behavior
func (wtc *WhatsAppTestConfig) SimulateAPIFailure(t *testing.T, phone string, shouldFail bool) {
	wtc.MockWhatsApp.SetShouldFail(shouldFail)
	
	w := wtc.MakeWhatsAppOTPRequest(t, phone)
	
	if shouldFail {
		wtc.AssertErrorResponse(t, w, http.StatusInternalServerError, "failed to send WhatsApp OTP")
	} else {
		wtc.AssertSuccessfulOTPResponse(t, w)
	}
}

// SimulateAPILatency simulates API call latency
func (wtc *WhatsAppTestConfig) SimulateAPILatency(t *testing.T, phone string, latency time.Duration) {
	wtc.MockWhatsApp.SetAPICallLatency(latency)
	
	startTime := time.Now()
	w := wtc.MakeWhatsAppOTPRequest(t, phone)
	elapsedTime := time.Since(startTime)
	
	wtc.AssertSuccessfulOTPResponse(t, w)
	
	// Check that the actual latency is close to the simulated latency
	assert.True(t, elapsedTime >= latency, "Actual latency should be at least the simulated latency")
}

// TestEndToEndFlow tests the complete WhatsApp authentication flow
func (wtc *WhatsAppTestConfig) TestEndToEndFlow(t *testing.T) {
	// Step 1: Register user
	userData := map[string]interface{}{
		"name":     wtc.TestUser.Name,
		"email":    wtc.TestUser.Email,
		"phone":    wtc.TestUser.Phone,
		"password": wtc.TestUser.Password,
		"role":     wtc.TestUser.Role,
	}
	
	w1 := wtc.MakeRegisterRequest(t, userData)
	assert.Equal(t, http.StatusCreated, w1.Code)
	
	// Step 2: Send OTP
	w2 := wtc.MakeWhatsAppOTPRequest(t, wtc.TestUser.Phone)
	wtc.AssertSuccessfulOTPResponse(t, w2)
	
	// Step 3: Verify OTP
	otp := wtc.GenerateTestOTP()
	wtc.SaveOTP(t, wtc.TestUser.Phone, otp)
	
	w3 := wtc.MakeWhatsAppVerifyRequest(t, wtc.TestUser.Phone, otp)
	wtc.AssertSuccessfulVerifyResponse(t, w3)
	
	// Step 4: Login
	w4 := wtc.MakeWhatsAppLoginRequest(t, wtc.TestUser.Phone, otp)
	wtc.AssertSuccessfulLoginResponse(t, w4)
	
	// Cleanup
	wtc.CleanupTestUser(t)
	wtc.CleanupOTP(t, wtc.TestUser.Phone)
}

// TestPhoneValidationFlow tests phone number validation throughout the flow
func (wtc *WhatsAppTestConfig) TestPhoneValidationFlow(t *testing.T) {
	validPhones := wtc.GenerateTestPhoneNumbers()
	invalidPhones := wtc.GenerateInvalidPhoneNumbers()
	
	// Test valid phone numbers
	for _, phone := range validPhones {
		t.Run(fmt.Sprintf("Valid phone: %s", phone), func(t *testing.T) {
			// user := wtc.CreateTestUserWithPhone(t, phone)
					
			w := wtc.MakeWhatsAppOTPRequest(t, phone)
			wtc.AssertSuccessfulOTPResponse(t, w)
					
			// wtc.CleanupTestUser(t)
		})
	}
	
	// Test invalid phone numbers
	for _, phone := range invalidPhones {
		t.Run(fmt.Sprintf("Invalid phone: %s", phone), func(t *testing.T) {
			w := wtc.MakeWhatsAppOTPRequest(t, phone)
			wtc.AssertErrorResponse(t, w, http.StatusBadRequest, "invalid phone number format")
		})
	}
}

// TestRateLimitingFlow tests rate limiting throughout the flow
func (wtc *WhatsAppTestConfig) TestRateLimitingFlow(t *testing.T) {
	phone := "+6289999888777"
	// user := wtc.CreateTestUserWithPhone(t, phone)
	
	// Test rate limiting
	wtc.SimulateRateLimiting(t, phone, true)
	
	// Cleanup
	// wtc.CleanupTestUser(t)
}

// TestAPIFailureFlow tests API failure handling throughout the flow
func (wtc *WhatsAppTestConfig) TestAPIFailureFlow(t *testing.T) {
	phone := "+6288887776666"
	// user := wtc.CreateTestUserWithPhone(t, phone)
	
	// Test API failure
	wtc.SimulateAPIFailure(t, phone, true)
	
	// Cleanup
	// wtc.CleanupTestUser(t)
}

// TestAPILatencyFlow tests API latency handling throughout the flow
func (wtc *WhatsAppTestConfig) TestAPILatencyFlow(t *testing.T) {
	phone := "+6287776665555"
	// user := wtc.CreateTestUserWithPhone(t, phone)
	
	// Test API latency
	wtc.SimulateAPILatency(t, phone, 500*time.Millisecond)
	
	// Cleanup
	// wtc.CleanupTestUser(t)
}

// GetMockWhatsAppStats returns statistics from the mock WhatsApp API
func (wtc *WhatsAppTestConfig) GetMockWhatsAppStats() map[string]interface{} {
	return wtc.MockWhatsApp.GetStatistics()
}

// GetMockWhatsAppMessages returns all messages sent through the mock WhatsApp API
func (wtc *WhatsAppTestConfig) GetMockWhatsAppMessages() []mocks.SentMessage {
	return wtc.MockWhatsApp.GetSentMessages()
}

// GetMockWhatsAppMessagesByProvider returns messages sent by a specific provider
func (wtc *WhatsAppTestConfig) GetMockWhatsAppMessagesByProvider(provider string) []mocks.SentMessage {
	return wtc.MockWhatsApp.GetSentMessagesByProvider(provider)
}

// ClearMockWhatsAppMessages clears all messages from the mock WhatsApp API
func (wtc *WhatsAppTestConfig) ClearMockWhatsAppMessages() {
	wtc.MockWhatsApp.ClearSentMessages()
}

// SetMockWhatsAppHealthStatus sets the health status of a provider in the mock API
func (wtc *WhatsAppTestConfig) SetMockWhatsAppHealthStatus(provider string, healthy bool) {
	wtc.MockWhatsApp.SetHealthStatus(provider, healthy)
}

// GetMockWhatsAppHealthStatus returns the health status of all providers
func (wtc *WhatsAppTestConfig) GetMockWhatsAppHealthStatus() map[string]bool {
	return wtc.MockWhatsApp.GetHealthStatus()
}

// Close closes the test database connection
func (wtc *WhatsAppTestConfig) Close() {
	if wtc.TestDB != nil {
		wtc.TestDB.Close()
	}
}

// WhatsAppTestSuite represents a complete test suite for WhatsApp authentication
type WhatsAppTestSuite struct {
	Config *WhatsAppTestConfig
}

// NewWhatsAppTestSuite creates a new WhatsApp test suite
func NewWhatsAppTestSuite(t *testing.T) *WhatsAppTestSuite {
	config := NewWhatsAppTestConfig(t)
	return &WhatsAppTestSuite{
		Config: config,
	}
}

// RunAllTests runs all WhatsApp authentication tests
func (suite *WhatsAppTestSuite) RunAllTests(t *testing.T) {
	t.Run("EndToEndFlow", suite.Config.TestEndToEndFlow)
	t.Run("PhoneValidationFlow", suite.Config.TestPhoneValidationFlow)
	t.Run("RateLimitingFlow", suite.Config.TestRateLimitingFlow)
	t.Run("APIFailureFlow", suite.Config.TestAPIFailureFlow)
	t.Run("APILatencyFlow", suite.Config.TestAPILatencyFlow)
}

// Close closes the test suite
func (suite *WhatsAppTestSuite) Close() {
	suite.Config.Close()
}