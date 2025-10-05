// Package password provides utility functions for password hashing and verification
// using the bcrypt algorithm. This package is designed for secure password handling
// in the Rangkai Edu authentication system.
//
// Security Considerations:
// - Uses bcrypt algorithm which is specifically designed for password hashing
// - Implements a cost factor of 12, providing a good balance between security and performance
// - Automatically handles salt generation for each password hash
// - Resistant to rainbow table and brute force attacks
//
// Usage:
//  hashedPassword, err := password.HashPassword("mySecretPassword")
//  if err != nil {
//      // Handle error
//  }
//
//  isValid := password.CheckPasswordHash("mySecretPassword", hashedPassword)
package password

import (
	"golang.org/x/crypto/bcrypt"
)

// bcryptCost defines the cost factor for bcrypt hashing.
// A cost of 12 is recommended as it provides a good balance between security and performance.
// Higher values increase security but also increase computation time exponentially.
const bcryptCost = 12

// HashPassword generates a bcrypt hash for the given password.
// It automatically generates a salt for each password and includes it in the hash.
//
// Parameters:
//   - password: The plain text password to hash
//
// Returns:
//   - string: The bcrypt hash of the password
//   - error: Any error that occurred during hashing
//
// Example:
//   hashedPassword, err := HashPassword("mySecretPassword")
//   if err != nil {
//       // Handle error appropriately
//       log.Fatal("Error hashing password:", err)
//   }
func HashPassword(password string) (string, error) {
	// Generate a bcrypt hash with the specified cost factor
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	
	// Convert the byte slice to a string and return
	return string(bytes), nil
}

// CheckPasswordHash compares a plain text password with a bcrypt hash.
// It returns true if the password matches the hash, false otherwise.
//
// Parameters:
//   - password: The plain text password to verify
//   - hash: The bcrypt hash to compare against
//
// Returns:
//   - bool: True if the password matches the hash, false otherwise
//
// Example:
//   isValid := CheckPasswordHash("mySecretPassword", hashedPassword)
//   if isValid {
//       // Password is correct
//   } else {
//       // Password is incorrect
//   }
func CheckPasswordHash(password, hash string) bool {
	// Compare the password with the hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	
	// If there's no error, the password is valid
	return err == nil
}