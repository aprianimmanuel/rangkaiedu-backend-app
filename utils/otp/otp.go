// Package otp provides utilities for generating, storing, and verifying one-time passwords (OTPs)
// for the Rangkai Edu authentication system. OTPs are 6-digit codes with a 10-minute expiry,
// stored in the database for secure verification.

package otp

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"time"
)

// OTPExpiryDuration defines the validity period for OTPs (10 minutes).
const OTPExpiryDuration = 10 * time.Minute

// GenerateOTP creates a secure 6-digit random OTP code using crypto/rand.
func GenerateOTP() (string, error) {
	// Generate a random number between 100000 and 999999
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", fmt.Errorf("failed to generate random OTP: %w", err)
	}
	otp := 100000 + n.Int64()
	return fmt.Sprintf("%06d", otp), nil
}

// SaveOTP stores the OTP in the database with an expiry time.
func SaveOTP(ctx context.Context, db *sql.DB, identifier, otp string) error {
	log.Printf("Saving OTP for identifier: %s", identifier)
	expiry := time.Now().Add(OTPExpiryDuration)
	_, err := db.ExecContext(ctx,
		"INSERT INTO otps (identifier, otp, expiry) VALUES (?, ?, ?) ON CONFLICT (identifier, otp) DO NOTHING",
		identifier, otp, expiry,
	)
	if err != nil {
		log.Printf("Failed to save OTP for identifier %s: %v", identifier, err)
		return fmt.Errorf("failed to save OTP: %w", err)
	}
	log.Printf("Successfully saved OTP for identifier: %s", identifier)
	return nil
}

// VerifyAndDeleteOTP verifies the OTP for the identifier, checks expiry, and deletes it if valid.
func VerifyAndDeleteOTP(ctx context.Context, db *sql.DB, identifier, otp string) (bool, error) {
	log.Printf("Verifying OTP for identifier: %s", identifier)
	var count int
	err := db.QueryRowContext(ctx,
		"DELETE FROM otps WHERE identifier = ? AND otp = ? AND expiry > NOW() RETURNING 1",
		identifier, otp,
	).Scan(&count)
	if err == sql.ErrNoRows {
		log.Printf("OTP not found or expired for identifier: %s", identifier)
		return false, nil // Invalid or expired OTP
	}
	if err != nil {
		log.Printf("Failed to verify OTP for identifier %s: %v", identifier, err)
		return false, fmt.Errorf("failed to verify OTP: %w", err)
	}
	log.Printf("Successfully verified OTP for identifier: %s", identifier)
	return count > 0, nil
}

// CleanupExpiredOTPs can be called periodically to remove old OTPs (optional, for maintenance).
func CleanupExpiredOTPs(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, "DELETE FROM otps WHERE expiry <= NOW()")
	if err != nil {
		return fmt.Errorf("failed to cleanup expired OTPs: %w", err)
	}
	return nil
}