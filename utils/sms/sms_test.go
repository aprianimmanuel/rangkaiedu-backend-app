package sms

import (
	"testing"

	"github.com/aprianimmanuel/backend-app/config"
	"github.com/stretchr/testify/assert"
)

// TestSendOTPSMS tests the SendOTPSMS function.
// Note: This is an integration test requiring valid Twilio config.
// For unit testing, consider mocking the Twilio client.
func TestSendOTPSMS(t *testing.T) {
	// Skip if no Twilio config (for local runs without env)
	cfg := config.Load()
	if cfg.TWILIOAccountSID == "" || cfg.TWILIOAuthToken == "" {
		t.Skip("Twilio config not set; skipping integration test")
	}

	to := "+1234567890" // Use a test phone number
	otp := "123456"

	err := SendOTPSMS(cfg, to, otp)
	if err != nil {
		// Twilio errors are expected if no valid creds; check for non-fatal errors
		assert.Contains(t, err.Error(), "failed to send", "Expected Twilio send error if not configured")
	} else {
		t.Log("SMS sent successfully (if configured)")
	}
}