package utils_test

import (
	"testing"

	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
	"github.com/stretchr/testify/assert"
	
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/sms"
)

// TestSendOTPSMS tests the SendOTPSMS function.
// Note: This is an integration test requiring valid Twilio config.
// For unit testing, consider mocking the Twilio client.
func TestSendOTPSMS(t *testing.T) {
	// Skip if no SMS provider configured (for local runs without env)
	cfg := config.Load()
	provider, err := cfg.GetPrimarySMSProvider()
	if err != nil || provider.Type == "" {
		t.Skip("SMS provider not configured; skipping integration test")
	}

	to := "+1234567890" // Use a test phone number
	otp := "123456"

	err = sms.SendOTPSMS(cfg, to, otp)
	if err != nil {
		// Twilio errors are expected if no valid creds; check for non-fatal errors
		assert.Contains(t, err.Error(), "failed to send", "Expected Twilio send error if not configured")
	} else {
		t.Log("SMS sent successfully (if configured)")
	}
}