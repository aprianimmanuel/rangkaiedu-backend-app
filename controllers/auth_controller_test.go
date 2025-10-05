package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/aprianimmanuel/backend-app/config"
	"github.com/aprianimmanuel/backend-app/controllers"
	"github.com/aprianimmanuel/backend-app/utils/password"
)

func TestMain(m *testing.M) {
	// Set environment variables for test database
	os.Setenv("DB_NAME", "rangkaiedu_test")
	os.Setenv("JWT_SECRET", "test-jwt-secret-change-in-production")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_SSLMODE", "disable")

	// Run tests
	code := m.Run()

	// Optional: unset env vars after tests
	os.Unsetenv("DB_NAME")
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

	// Manually set up auth routes to avoid import cycle
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controllers.RegisterHandler)
		auth.POST("/send-otp", controllers.SendOTPHandler)
		auth.POST("/verify-otp", controllers.VerifyOTPHandler)
		auth.POST("/login", controllers.LoginHandler)
	}

	return r
}

func createTestPool(t *testing.T) *pgxpool.Pool {
	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, pgxpool.ParseConfig(cfg.DSN()))
	if err != nil {
		t.Fatalf("Failed to create test pool: %v", err)
	}
	return pool
}

func cleanupUser(t *testing.T, pool *pgxpool.Pool, email string) {
	_, err := pool.Exec(context.Background(), "DELETE FROM users WHERE email = $1", email)
	if err != nil && err != pgx.ErrNoRows {
		t.Logf("Cleanup warning for %s: %v", email, err)
	}
}

func cleanupOTP(t *testing.T, pool *pgxpool.Pool, identifier string) {
	_, err := pool.Exec(context.Background(), "DELETE FROM otps WHERE identifier = $1", identifier)
	if err != nil && err != pgx.ErrNoRows {
		t.Logf("Cleanup warning for OTP %s: %v", identifier, err)
	}
}

func TestRegisterHandler_Success(t *testing.T) {
	r := setupTestRouter()
	pool := createTestPool(t)
	defer pool.Close()

	email := "test@example.com"
	reqBody, _ := json.Marshal(map[string]string{
		"name":     "Test User",
		"email":    email,
		"phone":    "1234567890",
		"password": "StrongPass1!",
		"role":     "student",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(reqBody))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if msg, ok := resp["message"]; !ok || msg != "User registered successfully" {
		t.Errorf("Expected message 'User registered successfully', got %v", msg)
	}

	// Verify user in DB
	var name, storedEmail, phone, role, hash string
	err := pool.QueryRow(context.Background(), "SELECT name, email, phone, role, password_hash FROM users WHERE email = $1", email).Scan(&name, &storedEmail, &phone, &role, &hash)
	if err != nil {
		t.Fatalf("Failed to query user: %v", err)
	}
	if name != "Test User" || storedEmail != email || phone != "1234567890" || role != "student" {
		t.Errorf("User data mismatch: got name=%s, email=%s, phone=%s, role=%s", name, storedEmail, phone, role)
	}
	if !password.CheckPasswordHash("StrongPass1!", hash) {
		t.Error("Password hash verification failed")
	}

	cleanupUser(t, pool, email)
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	r := setupTestRouter()

	reqBody := []byte(`{"invalid": "json"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(reqBody))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if msg, ok := resp["error"]; !ok || msg == "" {
		t.Errorf("Expected error message, got %v", msg)
	}
}

func TestRegisterHandler_WeakPassword(t *testing.T) {
	r := setupTestRouter()
	pool := createTestPool(t)
	defer pool.Close()

	email := "weak@example.com"
	reqBody, _ := json.Marshal(map[string]string{
		"name":     "Weak User",
		"email":    email,
		"phone":    "1234567890",
		"password": "weakpass", // No upper, digit, etc.
		"role":     "student",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(reqBody))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if msg, ok := resp["error"]; !ok || msg == "" {
		t.Errorf("Expected password error, got %v", msg)
	}

	// Verify no user created
	var count int
	pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if count != 0 {
		t.Error("User should not be created with weak password")
		cleanupUser(t, pool, email)
	}
}

func TestRegisterHandler_DuplicateEmail(t *testing.T) {
	r := setupTestRouter()
	pool := createTestPool(t)
	defer pool.Close()

	email := "duplicate@example.com"

	// First registration
	reqBody1, _ := json.Marshal(map[string]string{
		"name":     "First User",
		"email":    email,
		"phone":    "1234567890",
		"password": "StrongPass1!",
		"role":     "student",
	})
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(reqBody1))
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusCreated {
		t.Fatalf("First registration failed: %d", w1.Code)
	}

	// Second registration (duplicate)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(reqBody1))
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("Expected status %d for duplicate, got %d", http.StatusConflict, w2.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	if msg, ok := resp["error"]; !ok || msg != "Email already exists" {
		t.Errorf("Expected 'Email already exists', got %v", msg)
	}

	cleanupUser(t, pool, email)
}

func TestLoginHandler_Success(t *testing.T) {
	r := setupTestRouter()
	pool := createTestPool(t)
	defer pool.Close()

	email := "login@example.com"
	password := "StrongPass1!"

	// First register
	regBody, _ := json.Marshal(map[string]string{
		"name":     "Login User",
		"email":    email,
		"phone":    "1234567890",
		"password": password,
		"role":     "student",
	})
	wReg := httptest.NewRecorder()
	reqReg, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(regBody))
	r.ServeHTTP(wReg, reqReg)
	if wReg.Code != http.StatusCreated {
		t.Fatalf("Registration failed: %d", wReg.Code)
	}

	// Now login
	loginBody, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(loginBody))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if token, ok := resp["token"]; !ok || token == "" {
		t.Error("Expected token in response")
	}

	// Verify token is JWT (basic check)
	if len(resp["token"].(string)) < 10 {
		t.Error("Token seems invalid")
	}

	cleanupUser(t, pool, email)
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	r := setupTestRouter()
	pool := createTestPool(t)
	defer pool.Close()

	email := "invalid@example.com"
	password := "WrongPass!"

	// Register first
	regBody, _ := json.Marshal(map[string]string{
		"name":     "Invalid User",
		"email":    email,
		"phone":    "1234567890",
		"password": "StrongPass1!",
		"role":     "student",
	})
	wReg := httptest.NewRecorder()
	reqReg, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(regBody))
	r.ServeHTTP(wReg, reqReg)

	// Login with wrong password
	loginBody, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(loginBody))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if msg, ok := resp["error"]; !ok || msg != "Invalid credentials" {
		t.Errorf("Expected 'Invalid credentials', got %v", msg)
	}

	cleanupUser(t, pool, email)
}

func TestLoginHandler_NonExistentUser(t *testing.T) {
	r := setupTestRouter()

	loginBody, _ := json.Marshal(map[string]string{
		"email":    "nonexistent@example.com",
		"password": "SomePass1!",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(loginBody))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if msg, ok := resp["error"]; !ok || msg != "User not found" {
		t.Errorf("Expected 'User not found', got %v", msg)
	}
}

func TestLoginHandler_NoPasswordSet(t *testing.T) {
	r := setupTestRouter()
	pool := createTestPool(t)
	defer pool.Close()

	email := "nopass@example.com"

	regBody, _ := json.Marshal(map[string]interface{}{
		"name":  "No Pass User",
		"email": email,
		"phone": "1234567890",
		"role":  "student",
		// no password
	})
	wReg := httptest.NewRecorder()
	reqReg, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(regBody))
	r.ServeHTTP(wReg, reqReg)
	if wReg.Code != http.StatusCreated {
		t.Fatalf("Registration failed: %d", wReg.Code)
	}

	loginBody, _ := json.Marshal(map[string]interface{}{
		"email":    email,
		"password": "SomePass1!",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(loginBody))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if msg, ok := resp["error"]; !ok || msg != "Invalid credentials" {
		t.Errorf("Expected 'Invalid credentials', got %v", msg)
	}

	var hash sql.NullString
	err := pool.QueryRow(context.Background(), "SELECT password_hash FROM users WHERE email = $1", email).Scan(&hash)
	if err != nil {
		t.Fatalf("Failed to query hash: %v", err)
	}
	if hash.Valid {
		t.Error("Password hash should not be set")
	}

	cleanupUser(t, pool, email)
}

func TestSendOTPHandler(t *testing.T) {
	r := setupTestRouter()
	pool := createTestPool(t)
	defer pool.Close()

	t.Run("Success Email", func(t *testing.T) {
		email := "sendotp@example.com"

		regBody, _ := json.Marshal(map[string]interface{}{
			"name":     "SendOTP User",
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

		sendBody, _ := json.Marshal(map[string]interface{}{
			"identifier": email,
			"type":       "email",
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/send-otp", bytes.NewBuffer(sendBody))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if msg, ok := resp["message"]; !ok || msg != "OTP sent successfully" {
			t.Errorf("Expected 'OTP sent successfully', got %v", msg)
		}

		var otpStr string
		var expiry time.Time
		err := pool.QueryRow(context.Background(), "SELECT otp, expiry FROM otps WHERE identifier = $1", email).Scan(&otpStr, &expiry)
		if err != nil {
			t.Fatalf("Failed to query OTP: %v", err)
		}
		if len(otpStr) != 6 || !expiry.After(time.Now()) {
			t.Error("OTP not saved correctly")
		}

		cleanupUser(t, pool, email)
		cleanupOTP(t, pool, email)
	})

	t.Run("Success Phone", func(t *testing.T) {
		email := "phone@example.com"
		phoneNum := "1234567891"

		regBody, _ := json.Marshal(map[string]interface{}{
			"name":     "Phone User",
			"email":    email,
			"phone":    phoneNum,
			"password": "StrongPass1!",
			"role":     "student",
		})
		wReg := httptest.NewRecorder()
		reqReg, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(regBody))
		r.ServeHTTP(wReg, reqReg)
		if wReg.Code != http.StatusCreated {
			t.Fatalf("Registration failed: %d", wReg.Code)
		}

		sendBody, _ := json.Marshal(map[string]interface{}{
			"identifier": phoneNum,
			"type":       "phone",
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/send-otp", bytes.NewBuffer(sendBody))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var otpStr string
		err := pool.QueryRow(context.Background(), "SELECT otp FROM otps WHERE identifier = $1", phoneNum).Scan(&otpStr)
		if err != nil {
			t.Fatalf("Failed to query OTP: %v", err)
		}
		if len(otpStr) != 6 {
			t.Error("OTP not saved correctly")
		}

		cleanupUser(t, pool, email)
		cleanupOTP(t, pool, phoneNum)
	})

	t.Run("User Not Found", func(t *testing.T) {
		sendBody, _ := json.Marshal(map[string]interface{}{
			"identifier": "nonexistent@example.com",
			"type":       "email",
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/send-otp", bytes.NewBuffer(sendBody))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if msg, ok := resp["error"]; !ok || msg != "User not found for the provided identifier" {
			t.Errorf("Expected 'User not found for the provided identifier', got %v", msg)
		}
	})

	t.Run("Invalid Type", func(t *testing.T) {
		email := "invalidtype@example.com"

		regBody, _ := json.Marshal(map[string]interface{}{
			"name":     "Invalid Type User",
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

		sendBody, _ := json.Marshal(map[string]interface{}{
			"identifier": email,
			"type":       "invalid",
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/send-otp", bytes.NewBuffer(sendBody))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		cleanupUser(t, pool, email)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		reqBody := []byte(`{"invalid": "json"}`)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/send-otp", bytes.NewBuffer(reqBody))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestVerifyOTPHandler(t *testing.T) {
	r := setupTestRouter()
	pool := createTestPool(t)
	defer pool.Close()

	t.Run("Success", func(t *testing.T) {
		email := "verify@example.com"

		regBody, _ := json.Marshal(map[string]interface{}{
			"name":     "Verify User",
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

		var otp string
		err := pool.QueryRow(context.Background(), "SELECT otp FROM otps WHERE identifier = $1", email).Scan(&otp)
		if err != nil {
			t.Fatalf("Failed to get OTP: %v", err)
		}

		verifyBody, _ := json.Marshal(map[string]interface{}{
			"identifier": email,
			"otp":        otp,
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/verify-otp", bytes.NewBuffer(verifyBody))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		token, ok := resp["token"].(string)
		if !ok || token == "" {
			t.Error("Expected token in response")
		}
		user, ok := resp["user"].(map[string]interface{})
		if !ok || user["email"] != email {
			t.Error("Expected user data in response")
		}

		var count int
		err = pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM otps WHERE identifier = $1", email).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to check deletion: %v", err)
		}
		if count != 0 {
			t.Error("OTP should be deleted after verification")
		}

		cleanupUser(t, pool, email)
	})

	t.Run("Invalid OTP", func(t *testing.T) {
		email := "invalidotp@example.com"

		regBody, _ := json.Marshal(map[string]interface{}{
			"name":     "Invalid OTP User",
			"email":    email,
			"phone":    "1234567893",
			"password": "StrongPass1!",
			"role":     "student",
		})
		wReg := httptest.NewRecorder()
		reqReg, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(regBody))
		r.ServeHTTP(wReg, reqReg)
		if wReg.Code != http.StatusCreated {
			t.Fatalf("Registration failed: %d", wReg.Code)
		}

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

		var otp string
		err := pool.QueryRow(context.Background(), "SELECT otp FROM otps WHERE identifier = $1", email).Scan(&otp)
		if err != nil {
			t.Fatalf("Failed to get OTP: %v", err)
		}

		verifyBody, _ := json.Marshal(map[string]interface{}{
			"identifier": email,
			"otp":        "wrongotp",
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/verify-otp", bytes.NewBuffer(verifyBody))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if msg, ok := resp["error"]; !ok || msg != "Invalid or expired OTP" {
			t.Errorf("Expected 'Invalid or expired OTP', got %v", msg)
		}

		var count int
		err = pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM otps WHERE identifier = $1", email).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to check: %v", err)
		}
		if count != 1 {
			t.Error("OTP should remain after invalid verification")
		}

		cleanupUser(t, pool, email)
		cleanupOTP(t, pool, email)
	})

	t.Run("Expired OTP", func(t *testing.T) {
		email := "expired@example.com"

		regBody, _ := json.Marshal(map[string]interface{}{
			"name":     "Expired User",
			"email":    email,
			"phone":    "1234567894",
			"password": "StrongPass1!",
			"role":     "student",
		})
		wReg := httptest.NewRecorder()
		reqReg, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(regBody))
		r.ServeHTTP(wReg, reqReg)
		if wReg.Code != http.StatusCreated {
			t.Fatalf("Registration failed: %d", wReg.Code)
		}

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

		_, err := pool.Exec(context.Background(), "UPDATE otps SET expiry = NOW() - INTERVAL '1 hour' WHERE identifier = $1", email)
		if err != nil {
			t.Fatalf("Failed to set expired: %v", err)
		}

		var otp string
		err = pool.QueryRow(context.Background(), "SELECT otp FROM otps WHERE identifier = $1", email).Scan(&otp)
		if err != nil {
			t.Fatalf("Failed to get OTP: %v", err)
		}

		verifyBody, _ := json.Marshal(map[string]interface{}{
			"identifier": email,
			"otp":        otp,
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/verify-otp", bytes.NewBuffer(verifyBody))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if msg, ok := resp["error"]; !ok || msg != "Invalid or expired OTP" {
			t.Errorf("Expected 'Invalid or expired OTP', got %v", msg)
		}

		var count int
		err = pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM otps WHERE identifier = $1", email).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to check: %v", err)
		}
		if count != 0 {
			t.Error("Expired OTP should be deleted")
		}

		cleanupUser(t, pool, email)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		reqBody := []byte(`{"invalid": "json"}`)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/verify-otp", bytes.NewBuffer(reqBody))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Non-existent User", func(t *testing.T) {
		verifyBody, _ := json.Marshal(map[string]interface{}{
			"identifier": "nonexistent@example.com",
			"otp":        "123456",
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/verify-otp", bytes.NewBuffer(verifyBody))
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if msg, ok := resp["error"]; !ok || msg != "User not found" {
			t.Errorf("Expected 'User not found', got %v", msg)
		}
	})
}

func TestEndToEndAuthFlow(t *testing.T) {
	r := setupTestRouter()
	pool := createTestPool(t)
	defer pool.Close()

	email := "e2e@example.com"

	// Register
	regBody, _ := json.Marshal(map[string]interface{}{
		"name":     "E2E User",
		"email":    email,
		"phone":    "1234567895",
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

	// Get OTP
	var otp string
	err := pool.QueryRow(context.Background(), "SELECT otp FROM otps WHERE identifier = $1", email).Scan(&otp)
	if err != nil {
		t.Fatalf("Failed to get OTP: %v", err)
	}

	// Verify OTP
	verifyBody, _ := json.Marshal(map[string]interface{}{
		"identifier": email,
		"otp":        otp,
	})
	wVerify := httptest.NewRecorder()
	reqVerify, _ := http.NewRequest("POST", "/api/auth/verify-otp", bytes.NewBuffer(verifyBody))
	r.ServeHTTP(wVerify, reqVerify)

	if wVerify.Code != http.StatusOK {
		t.Errorf("Expected status %d for verify, got %d", http.StatusOK, wVerify.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(wVerify.Body.Bytes(), &resp)
	token, ok := resp["token"].(string)
	if !ok || token == "" {
		t.Error("Expected token in verify response")
	}

	// Cleanup
	cleanupUser(t, pool, email)
	cleanupOTP(t, pool, email)
}