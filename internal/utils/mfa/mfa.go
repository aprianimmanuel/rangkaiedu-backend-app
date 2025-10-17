// Package mfa provides utilities for Multi-Factor Authentication (MFA) implementation
// using Time-based One-Time Passwords (TOTP) for the Rangkai Edu authentication system.
// This includes TOTP secret generation, QR code creation, backup code management,
// and secure validation following RFC 6238 standards.

package mfa

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/monitoring"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/db"
)

const (
	// BackupCodeCount defines the number of backup codes to generate
	BackupCodeCount = 6
	// BackupCodeLength defines the length of each backup code
	BackupCodeLength = 6
	// TOTPIssuer defines the issuer name for TOTP configuration
	TOTPIssuer = "Rangkai Edu"
)

// GenerateTOTPSecret generates a secure TOTP secret key
func GenerateTOTPSecret() (string, error) {
	// Generate a random secret key
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP secret: %w", err)
	}
	
	// Encode as base32 string (TOTP requires base32)
	return base32.StdEncoding.EncodeToString(secret), nil
}

// GenerateBackupCodes generates secure backup codes for MFA recovery
func GenerateBackupCodes() ([]string, error) {
	codes := make([]string, BackupCodeCount)
	
	for i := 0; i < BackupCodeCount; i++ {
	code, err := GenerateRandomCode(BackupCodeLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate backup code %d: %w", i, err)
	}
	codes[i] = code
}
	
	return codes, nil
}

// GenerateRandomCode generates a random alphanumeric code of specified length
func GenerateRandomCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}
	
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := make([]byte, length)
	
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		code[i] = charset[n.Int64()]
	}
	
	return string(code), nil
}

// ValidateTOTP validates a TOTP token against the stored secret
func ValidateTOTP(secret, token string) bool {
	// Remove any whitespace from the token
	token = strings.TrimSpace(token)
	
	// Validate TOTP token
	valid := totp.Validate(token, secret)
	
	return valid
}

// ValidateBackupCode validates a backup code against the stored codes
func ValidateBackupCode(backupCodes []string, code string) bool {
	// Remove any whitespace from the code
	code = strings.TrimSpace(code)
	
	// Check if the code exists in the backup codes
	for _, storedCode := range backupCodes {
		if storedCode == code {
			return true
		}
	}
	
	return false
}

// CreateQRCode generates a QR code image for TOTP setup
func CreateQRCode(secret, userEmail string) (string, error) {
	// Create TOTP key
	key, err := otp.NewKeyFromURL(fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s",
		TOTPIssuer,
		userEmail,
		secret,
		TOTPIssuer,
	))
	if err != nil {
		return "", fmt.Errorf("failed to create TOTP key: %w", err)
	}
	
	// Generate QR code
	qrCode, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}
	
	// Encode as base64 for JSON response
	return base64.StdEncoding.EncodeToString(qrCode), nil
}

// LogMFAEvent logs MFA events for security auditing
func LogMFAEvent(userID string, eventType string, success bool, ipAddress, userAgent string) error {
	ctx := context.Background()
	
	query := `
		INSERT INTO mfa_events (user_id, event_type, ip_address, user_agent, success)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	_, err := db.GetDB().ExecContext(ctx, query, userID, eventType, ipAddress, userAgent, success)
	if err != nil {
		log.Printf("Failed to log MFA event: %v", err)
		return fmt.Errorf("failed to log MFA event: %w", err)
	}
	
	// Log security event for MFA
	monitoring.LogAuthMFA(ctx, userID, "", ipAddress, eventType, success, map[string]interface{}{
		"user_agent": userAgent,
	})
	
	log.Printf("MFA event logged: user_id=%s, event_type=%s, success=%t", userID, eventType, success)
	return nil
}

// IsMFAEnabled checks if MFA is enabled for a user
func IsMFAEnabled(userID string) (bool, error) {
	ctx := context.Background()
	
	var enabled bool
	query := "SELECT is_mfa_enabled FROM users WHERE id = $1"
	
	err := db.GetDB().QueryRowContext(ctx, query, userID).Scan(&enabled)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("user not found")
		}
		return false, fmt.Errorf("failed to check MFA status: %w", err)
	}
	
	return enabled, nil
}

// EnableMFA enables MFA for a user
func EnableMFA(userID string) error {
	ctx := context.Background()
	
	query := "UPDATE users SET is_mfa_enabled = TRUE WHERE id = $1"
	
	_, err := db.GetDB().ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to enable MFA: %w", err)
	}
	
	// Log security event for MFA enable
	monitoring.LogAuthMFA(ctx, userID, "", "", "enable", true, nil)
	
	log.Printf("MFA enabled for user: %s", userID)
	return nil
}

// DisableMFA disables MFA for a user
func DisableMFA(userID string) error {
	ctx := context.Background()
	
	query := "UPDATE users SET is_mfa_enabled = FALSE, mfa_secret = NULL, mfa_backup_codes = NULL WHERE id = $1"
	
	_, err := db.GetDB().ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to disable MFA: %w", err)
	}
	
	// Log security event for MFA disable
	monitoring.LogAuthMFA(ctx, userID, "", "", "disable", true, nil)
	
	log.Printf("MFA disabled for user: %s", userID)
	return nil
}

// GetUserMFAInfo retrieves MFA information for a user
func GetUserMFAInfo(userID string) (bool, []string, error) {
	ctx := context.Background()
	
	var enabled bool
	var backupCodes []string
	
	query := "SELECT is_mfa_enabled, mfa_backup_codes FROM users WHERE id = $1"
	
	err := db.GetDB().QueryRowContext(ctx, query, userID).Scan(&enabled, &backupCodes)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil, fmt.Errorf("user not found")
		}
		return false, nil, fmt.Errorf("failed to get MFA info: %w", err)
	}
	
	return enabled, backupCodes, nil
}

// UpdateUserMFASecret updates the MFA secret and backup codes for a user
func UpdateUserMFASecret(userID string, secret string, backupCodes []string) error {
	ctx := context.Background()
	
	query := `
		UPDATE users 
		SET mfa_secret = $1, mfa_backup_codes = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $3
	`
	
	_, err := db.GetDB().ExecContext(ctx, query, secret, backupCodes, userID)
	if err != nil {
		return fmt.Errorf("failed to update MFA secret: %w", err)
	}
	
	log.Printf("MFA secret updated for user: %s", userID)
	return nil
}

// ConsumeBackupCode removes a used backup code from the user's backup codes
func ConsumeBackupCode(userID string, code string) error {
	ctx := context.Background()
	
	// Get current backup codes
	var backupCodes []string
	query := "SELECT mfa_backup_codes FROM users WHERE id = $1"
	
	err := db.GetDB().QueryRowContext(ctx, query, userID).Scan(&backupCodes)
	if err != nil {
		return fmt.Errorf("failed to get backup codes: %w", err)
	}
	
	// Remove the used code
	newBackupCodes := make([]string, 0, len(backupCodes)-1)
	for _, bc := range backupCodes {
		if bc != code {
			newBackupCodes = append(newBackupCodes, bc)
		}
	}
	
	// Update the database
	query = "UPDATE users SET mfa_backup_codes = $1 WHERE id = $2"
	_, err = db.GetDB().ExecContext(ctx, query, newBackupCodes, userID)
	if err != nil {
		return fmt.Errorf("failed to update backup codes: %w", err)
	}
	
	// Log security event for backup code consumption
	monitoring.LogAuthMFA(ctx, userID, "", "", "backup_code_consumed", true, map[string]interface{}{
		"remaining_codes": len(newBackupCodes),
	})
	
	log.Printf("Backup code consumed for user: %s", userID)
	return nil
}

// GetClientIP retrieves the client IP address from the HTTP request
func GetClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header first (for reverse proxies)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}
	
	// Check for X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}
	
	// Fall back to RemoteAddr
	return strings.Split(r.RemoteAddr, ":")[0]
}

// GetUserAgent retrieves the user agent from the HTTP request
func GetUserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}