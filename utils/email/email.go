// Package email provides utilities for sending emails, specifically OTP codes,
// for the Rangkai Edu authentication system using Go's net/smtp package.

package email

import (
	"fmt"
	"net/smtp"

	"github.com/aprianimmanuel/rangkaiedu-backend/config"
)

// SendOTPEmail sends an email containing the OTP code to the specified recipient.
func SendOTPEmail(cfg *config.Config, to, otp string) error {
	// If SMTP host is not configured, skip sending email
	if cfg.SMTPHost == "" {
		return nil
	}

	from := cfg.SMTPUser
	password := cfg.SMTPPass
	host := cfg.SMTPHost
	port := cfg.SMTPPort

	// SMTP authentication
	auth := smtp.PlainAuth("", from, password, host)

	// Message template
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Your Rangkai Edu OTP Code\r\n"+
		"\r\n"+
		"Your OTP code is: %s\r\n"+
		"This code expires in 10 minutes. Do not share it with anyone.\r\n"+
		"If you didn't request this, please ignore this email.\r\n", to, otp))

	// Send the email
	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	return nil
}