package middleware

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/aprianimmanuel/rangkaiedu-backend/utils/db"
)

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Set environment variables for test database
	os.Setenv("DB_TEST_NAME", "rangkaiedu_test")
	os.Setenv("DB_NAME", "rangkaiedu_test") // Use the test database name
	os.Setenv("API_PORT", "8080")
	os.Setenv("JWT_SECRET", "test-jwt-secret-change-in-production")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "rangkaiedudev1")
	os.Setenv("DB_PASSWORD", "12d1q23wxm19wkc1fsdcq23")
	os.Setenv("DB_SSLMODE", "disable")

	// Initialize the database connection pool
	if err := db.Init(); err != nil {
		panic(err)
	}

	// Run tests
	code := m.Run()

	// Close the database connection pool
	db.Close()

	// Clean up
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_TEST_NAME")
	os.Unsetenv("API_PORT")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_SSLMODE")

	os.Exit(code)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Set up a protected route for testing
	protected := r.Group("/api/protected")
	protected.Use(AuthRequired())
	{
		protected.GET("/profile", func(c *gin.Context) {
			userInterface, exists := c.Get("user")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
				return
			}

			user, ok := userInterface.(UserContext)
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id":    user.ID,
				"email": user.Email,
				"phone": user.Phone,
				"role":  user.Role,
			})
		})

		// Admin-only route
		admin := protected.Group("/admin")
		admin.Use(RoleRequired(RoleAdmin))
		{
			admin.GET("/users", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Admin access granted"})
			})
		}

		// Teacher or admin route
		content := protected.Group("/content")
		content.Use(RolesRequired(RoleAdmin, RoleTeacher))
		{
			content.POST("/lessons", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Content creator access granted"})
			})
		}
	}

	// Manually set up auth routes to avoid import cycle
	// Note: These routes are not actually used in the tests since we're testing
	// the middleware functionality, not the auth controllers themselves.
	// The tests use the auth system through direct database operations and
	// token generation.
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented in tests"})
		})
		auth.POST("/send-otp", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented in tests"})
		})
		auth.POST("/verify-otp", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented in tests"})
		})
		auth.POST("/login", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented in tests"})
		})
	}

	return r
}

func createTestPool(t *testing.T) *sql.DB {
	// Use the global database connection pool
	pool := db.GetDB()
	if pool == nil {
		t.Fatalf("Database pool is not initialized")
	}
	return pool
}

func cleanupUser(t *testing.T, pool *sql.DB, email string) {
	_, err := pool.Exec("DELETE FROM users WHERE email = $1", email)
	if err != nil {
		t.Logf("Cleanup warning for %s: %v", email, err)
	}
}

func cleanupOTP(t *testing.T, pool *sql.DB, identifier string) {
	_, err := pool.Exec("DELETE FROM otps WHERE identifier = $1", identifier)
	if err != nil {
		t.Logf("Cleanup warning for OTP %s: %v", identifier, err)
	}
}

func TestAuthIntegration_ValidTokenFromAuthSystem(t *testing.T) {
	r := setupTestRouter()
	pool := createTestPool(t)
	defer pool.Close()

	email := "integration@example.com"

	// Register a user
	regBody, _ := json.Marshal(map[string]interface{}{
		"name":     "Integration Test User",
		"email":    email,
		"phone":    "1234567890",
		"password": "StrongPass1!",
		"role":     "student",
	})

	wReg := httptest.NewRecorder()
	reqReg, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(regBody))
	r.ServeHTTP(wReg, reqReg)

	if wReg.Code != http.StatusCreated {
		t.Fatalf("Registration failed: %d", wReg.Code)
	}

	// Send OTP
	sendBody, _ := json.Marshal(map[string]interface{}{
		"identifier": email,
		"type":       "email",
	})

	wSend := httptest.NewRecorder()
	reqSend, _ := http.NewRequest("POST", "/api/auth/send-otp", bytes.NewBuffer(sendBody))
	r.ServeHTTP(wSend, reqSend)

	if wSend.Code != http.StatusOK {
		t.Fatalf("Send OTP failed: %d", wSend.Code)
	}

	// Get OTP from database
	var otp string
	err := pool.QueryRow("SELECT otp FROM otps WHERE identifier = $1", email).Scan(&otp)
	if err != nil {
		t.Fatalf("Failed to get OTP: %v", err)
	}

	// Verify OTP to get token
	verifyBody, _ := json.Marshal(map[string]interface{}{
		"identifier": email,
		"otp":        otp,
	})

	wVerify := httptest.NewRecorder()
	reqVerify, _ := http.NewRequest("POST", "/api/auth/verify-otp", bytes.NewBuffer(verifyBody))
	r.ServeHTTP(wVerify, reqVerify)

	if wVerify.Code != http.StatusOK {
		t.Fatalf("Verify OTP failed: %d", wVerify.Code)
	}

	var verifyResp map[string]interface{}
	if err := json.Unmarshal(wVerify.Body.Bytes(), &verifyResp); err != nil {
		t.Fatalf("Failed to parse verify response: %v", err)
	}

	token, ok := verifyResp["token"].(string)
	if !ok || token == "" {
		t.Fatal("Token not found in verify response")
	}

	// Use token to access protected route
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/protected/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var profileResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &profileResp); err != nil {
		t.Fatalf("Failed to parse profile response: %v", err)
	}

	if profileResp["email"] != email {
		t.Errorf("Expected email %s, got %v", email, profileResp["email"])
	}

	// Cleanup
	cleanupUser(t, pool, email)
	cleanupOTP(t, pool, email)
}

func TestAuthIntegration_AdminRouteAccess(t *testing.T) {
	r := setupTestRouter()
	pool := createTestPool(t)
	defer pool.Close()

	email := "admin@example.com"

	// Register an admin user
	regBody, _ := json.Marshal(map[string]interface{}{
		"name":     "Admin User",
		"email":    email,
		"phone":    "1234567891",
		"password": "StrongPass1!",
		"role":     "admin",
	})

	wReg := httptest.NewRecorder()
	reqReg, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(regBody))
	r.ServeHTTP(wReg, reqReg)

	if wReg.Code != http.StatusCreated {
		t.Fatalf("Registration failed: %d", wReg.Code)
	}

	// Send OTP
	sendBody, _ := json.Marshal(map[string]interface{}{
		"identifier": email,
		"type":       "email",
	})

	wSend := httptest.NewRecorder()
	reqSend, _ := http.NewRequest("POST", "/api/auth/send-otp", bytes.NewBuffer(sendBody))
	r.ServeHTTP(wSend, reqSend)

	if wSend.Code != http.StatusOK {
		t.Fatalf("Send OTP failed: %d", wSend.Code)
	}

	// Get OTP from database
	var otp string
	err := pool.QueryRow("SELECT otp FROM otps WHERE identifier = $1", email).Scan(&otp)
	if err != nil {
		t.Fatalf("Failed to get OTP: %v", err)
	}

	// Verify OTP to get token
	verifyBody, _ := json.Marshal(map[string]interface{}{
		"identifier": email,
		"otp":        otp,
	})

	wVerify := httptest.NewRecorder()
	reqVerify, _ := http.NewRequest("POST", "/api/auth/verify-otp", bytes.NewBuffer(verifyBody))
	r.ServeHTTP(wVerify, reqVerify)

	if wVerify.Code != http.StatusOK {
		t.Fatalf("Verify OTP failed: %d", wVerify.Code)
	}

	var verifyResp map[string]interface{}
	if err := json.Unmarshal(wVerify.Body.Bytes(), &verifyResp); err != nil {
		t.Fatalf("Failed to parse verify response: %v", err)
	}

	token, ok := verifyResp["token"].(string)
	if !ok || token == "" {
		t.Fatal("Token not found in verify response")
	}

	// Use token to access admin route
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/protected/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Cleanup
	cleanupUser(t, pool, email)
	cleanupOTP(t, pool, email)
}

func TestAuthIntegration_ForbiddenAccess(t *testing.T) {
	r := setupTestRouter()
	pool := createTestPool(t)
	defer pool.Close()

	email := "student@example.com"

	// Register a student user (not admin)
	regBody, _ := json.Marshal(map[string]interface{}{
		"name":     "Student User",
		"email":    email,
		"phone":    "1234567892",
		"password": "StrongPass1!",
		"role":     "student",
	})

	wReg := httptest.NewRecorder()
	reqReg, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(regBody))
	r.ServeHTTP(wReg, reqReg)

	if wReg.Code != http.StatusCreated {
		t.Fatalf("Registration failed: %d", wReg.Code)
	}

	// Send OTP
	sendBody, _ := json.Marshal(map[string]interface{}{
		"identifier": email,
		"type":       "email",
	})

	wSend := httptest.NewRecorder()
	reqSend, _ := http.NewRequest("POST", "/api/auth/send-otp", bytes.NewBuffer(sendBody))
	r.ServeHTTP(wSend, reqSend)

	if wSend.Code != http.StatusOK {
		t.Fatalf("Send OTP failed: %d", wSend.Code)
	}

	// Get OTP from database
	var otp string
	err := pool.QueryRow("SELECT otp FROM otps WHERE identifier = $1", email).Scan(&otp)
	if err != nil {
		t.Fatalf("Failed to get OTP: %v", err)
	}

	// Verify OTP to get token
	verifyBody, _ := json.Marshal(map[string]interface{}{
		"identifier": email,
		"otp":        otp,
	})

	wVerify := httptest.NewRecorder()
	reqVerify, _ := http.NewRequest("POST", "/api/auth/verify-otp", bytes.NewBuffer(verifyBody))
	r.ServeHTTP(wVerify, reqVerify)

	if wVerify.Code != http.StatusOK {
		t.Fatalf("Verify OTP failed: %d", wVerify.Code)
	}

	var verifyResp map[string]interface{}
	if err := json.Unmarshal(wVerify.Body.Bytes(), &verifyResp); err != nil {
		t.Fatalf("Failed to parse verify response: %v", err)
	}

	token, ok := verifyResp["token"].(string)
	if !ok || token == "" {
		t.Fatal("Token not found in verify response")
	}

	// Use token to access admin route (should be forbidden)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/protected/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}

	// Cleanup
	cleanupUser(t, pool, email)
	cleanupOTP(t, pool, email)
}

func TestAuthIntegration_TeacherAccess(t *testing.T) {
	r := setupTestRouter()
	pool := createTestPool(t)
	defer pool.Close()

	email := "teacher@example.com"

	// Register a teacher user
	regBody, _ := json.Marshal(map[string]interface{}{
		"name":     "Teacher User",
		"email":    email,
		"phone":    "1234567893",
		"password": "StrongPass1!",
		"role":     "teacher",
	})

	wReg := httptest.NewRecorder()
	reqReg, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(regBody))
	r.ServeHTTP(wReg, reqReg)

	if wReg.Code != http.StatusCreated {
		t.Fatalf("Registration failed: %d", wReg.Code)
	}

	// Send OTP
	sendBody, _ := json.Marshal(map[string]interface{}{
		"identifier": email,
		"type":       "email",
	})

	wSend := httptest.NewRecorder()
	reqSend, _ := http.NewRequest("POST", "/api/auth/send-otp", bytes.NewBuffer(sendBody))
	r.ServeHTTP(wSend, reqSend)

	if wSend.Code != http.StatusOK {
		t.Fatalf("Send OTP failed: %d", wSend.Code)
	}

	// Get OTP from database
	var otp string
	err := pool.QueryRow("SELECT otp FROM otps WHERE identifier = $1", email).Scan(&otp)
	if err != nil {
		t.Fatalf("Failed to get OTP: %v", err)
	}

	// Verify OTP to get token
	verifyBody, _ := json.Marshal(map[string]interface{}{
		"identifier": email,
		"otp":        otp,
	})

	wVerify := httptest.NewRecorder()
	reqVerify, _ := http.NewRequest("POST", "/api/auth/verify-otp", bytes.NewBuffer(verifyBody))
	r.ServeHTTP(wVerify, reqVerify)

	if wVerify.Code != http.StatusOK {
		t.Fatalf("Verify OTP failed: %d", wVerify.Code)
	}

	var verifyResp map[string]interface{}
	if err := json.Unmarshal(wVerify.Body.Bytes(), &verifyResp); err != nil {
		t.Fatalf("Failed to parse verify response: %v", err)
	}

	token, ok := verifyResp["token"].(string)
	if !ok || token == "" {
		t.Fatal("Token not found in verify response")
	}

	// Use token to access content route (teachers should have access)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/protected/content/lessons", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Cleanup
	cleanupUser(t, pool, email)
	cleanupOTP(t, pool, email)
}

func TestAuthIntegration_InvalidToken(t *testing.T) {
	r := setupTestRouter()

	// Try to access protected route with invalid token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/protected/profile", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthIntegration_MalformedToken(t *testing.T) {
	r := setupTestRouter()

	// Try to access protected route with malformed token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/protected/profile", nil)
	req.Header.Set("Authorization", "Bearer")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthIntegration_MissingToken(t *testing.T) {
	r := setupTestRouter()

	// Try to access protected route without token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/protected/profile", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}