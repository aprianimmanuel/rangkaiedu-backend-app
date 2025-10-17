package utils_test

import (
	"testing"
	
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/password"
)

// TestHashPassword tests the HashPassword function
func TestHashPassword(t *testing.T) {
	plainPassword := "testPassword123"
	
	hashedPassword, err := password.HashPassword(plainPassword)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if hashedPassword == "" {
		t.Fatal("Expected hashed password, got empty string")
	}
	
	// The hashed password should be different from the original password
	if hashedPassword == plainPassword {
		t.Fatal("Hashed password should be different from original password")
	}
}

// TestCheckPasswordHash tests the CheckPasswordHash function
func TestCheckPasswordHash(t *testing.T) {
	plainPassword := "testPassword123"
	
	// Hash the password
	hashedPassword, err := password.HashPassword(plainPassword)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}
	
	// Check with correct password
	if !password.CheckPasswordHash(plainPassword, hashedPassword) {
		t.Fatal("Expected password to match hash")
	}
	
	// Check with incorrect password
	if password.CheckPasswordHash("wrongPassword", hashedPassword) {
		t.Fatal("Expected wrong password to not match hash")
	}
}

// TestHashPasswordWithEmptyString tests edge case with empty password
func TestHashPasswordWithEmptyString(t *testing.T) {
	plainPassword := ""
	
	hashedPassword, err := password.HashPassword(plainPassword)
	if err != nil {
		t.Fatalf("Expected no error for empty password, got %v", err)
	}
	
	if hashedPassword == "" {
		t.Fatal("Expected hashed password for empty string, got empty string")
	}
	
	// Check if empty password matches its hash
	if !password.CheckPasswordHash(plainPassword, hashedPassword) {
		t.Fatal("Expected empty password to match its hash")
	}
}