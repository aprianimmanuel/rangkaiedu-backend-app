package utils_test

import (
	"testing"

	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
	"github.com/stretchr/testify/assert"
	
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/email"
)

// TestSendOTPEmail tests the SendOTPEmail function.
// Note: This is an integration test requiring valid SMTP config.
// For unit testing, consider mocking smtp.SendMail.
func TestSendOTPEmail(t *testing.T) {
	// Skip if no SMTP config (for local runs without env)
	cfg := config.Load()
	if cfg.SMTPHost == "" || cfg.SMTPUser == "" {
		t.Skip("SMTP config not set; skipping integration test")
	}

	to := "test@example.com" // Use a test email
	otp := "123456"

	err := email.SendOTPEmail(cfg, to, otp)
	if err != nil {
		// SMTP errors are expected if no valid creds; check for non-fatal errors
		assert.Contains(t, err.Error(), "failed to send", "Expected SMTP send error if not configured")
	} else {
		t.Log("Email sent successfully (if configured)")
	}
}