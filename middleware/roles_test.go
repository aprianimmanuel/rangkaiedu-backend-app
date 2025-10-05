package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRoleRequired_ValidRole(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create test context with user context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	// Set user context with required role
	user := UserContext{
		ID:    "test-user-id",
		Email: "test@example.com",
		Phone: "1234567890",
		Role:  "admin",
	}
	c.Set("user", user)

	// Create a test handler
	testHandler := func(c *gin.Context) {
		c.Status(http.StatusOK)
	}

	// Apply middleware and handler
	RoleRequired(RoleAdmin)(c)
	if c.IsAborted() {
		t.Error("Request was aborted unexpectedly")
		return
	}

	// Continue with test handler
	testHandler(c)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestRoleRequired_InvalidRole(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create test context with user context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	// Set user context with wrong role
	user := UserContext{
		ID:    "test-user-id",
		Email: "test@example.com",
		Phone: "1234567890",
		Role:  "student",
	}
	c.Set("user", user)

	// Apply middleware
	RoleRequired(RoleAdmin)(c)

	// Verify response
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}

	if !c.IsAborted() {
		t.Error("Request should be aborted")
	}

	var response map[string]interface{}
	if err := c.ShouldBindJSON(&response); err == nil {
		if errorMsg, ok := response["error"]; !ok || errorMsg != "Insufficient permissions" {
			t.Errorf("Expected error message 'Insufficient permissions', got '%v'", errorMsg)
		}
	}
}

func TestRoleRequired_NoUserContext(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create test context without user context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	// Apply middleware
	RoleRequired(RoleAdmin)(c)

	// Verify response
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if !c.IsAborted() {
		t.Error("Request should be aborted")
	}

	var response map[string]interface{}
	if err := c.ShouldBindJSON(&response); err == nil {
		if errorMsg, ok := response["error"]; !ok || errorMsg != "User not authenticated" {
			t.Errorf("Expected error message 'User not authenticated', got '%v'", errorMsg)
		}
	}
}

func TestRoleRequired_InvalidUserContext(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create test context with invalid user context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	// Set invalid user context
	c.Set("user", "invalid-context")

	// Apply middleware
	RoleRequired(RoleAdmin)(c)

	// Verify response
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if !c.IsAborted() {
		t.Error("Request should be aborted")
	}

	var response map[string]interface{}
	if err := c.ShouldBindJSON(&response); err == nil {
		if errorMsg, ok := response["error"]; !ok || errorMsg != "Invalid user context" {
			t.Errorf("Expected error message 'Invalid user context', got '%v'", errorMsg)
		}
	}
}

func TestRolesRequired_ValidRole(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create test context with user context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	// Set user context with one of the allowed roles
	user := UserContext{
		ID:    "test-user-id",
		Email: "test@example.com",
		Phone: "1234567890",
		Role:  "teacher",
	}
	c.Set("user", user)

	// Create a test handler
	testHandler := func(c *gin.Context) {
		c.Status(http.StatusOK)
	}

	// Apply middleware and handler (allowing admin or teacher)
	RolesRequired(RoleAdmin, RoleTeacher)(c)
	if c.IsAborted() {
		t.Error("Request was aborted unexpectedly")
		return
	}

	// Continue with test handler
	testHandler(c)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestRolesRequired_InvalidRole(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create test context with user context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	// Set user context with role not in allowed list
	user := UserContext{
		ID:    "test-user-id",
		Email: "test@example.com",
		Phone: "1234567890",
		Role:  "student",
	}
	c.Set("user", user)

	// Apply middleware (allowing only admin or teacher)
	RolesRequired(RoleAdmin, RoleTeacher)(c)

	// Verify response
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}

	if !c.IsAborted() {
		t.Error("Request should be aborted")
	}

	var response map[string]interface{}
	if err := c.ShouldBindJSON(&response); err == nil {
		if errorMsg, ok := response["error"]; !ok || errorMsg != "Insufficient permissions" {
			t.Errorf("Expected error message 'Insufficient permissions', got '%v'", errorMsg)
		}
	}
}

func TestRolesRequired_NoUserContext(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create test context without user context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	// Apply middleware
	RolesRequired(RoleAdmin, RoleTeacher)(c)

	// Verify response
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if !c.IsAborted() {
		t.Error("Request should be aborted")
	}

	var response map[string]interface{}
	if err := c.ShouldBindJSON(&response); err == nil {
		if errorMsg, ok := response["error"]; !ok || errorMsg != "User not authenticated" {
			t.Errorf("Expected error message 'User not authenticated', got '%v'", errorMsg)
		}
	}
}

func TestRolesRequired_InvalidUserContext(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create test context with invalid user context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	// Set invalid user context
	c.Set("user", 12345)

	// Apply middleware
	RolesRequired(RoleAdmin, RoleTeacher)(c)

	// Verify response
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if !c.IsAborted() {
		t.Error("Request should be aborted")
	}

	var response map[string]interface{}
	if err := c.ShouldBindJSON(&response); err == nil {
		if errorMsg, ok := response["error"]; !ok || errorMsg != "Invalid user context" {
			t.Errorf("Expected error message 'Invalid user context', got '%v'", errorMsg)
		}
	}
}