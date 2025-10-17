// Package sms provides utilities for sending SMS messages, specifically OTP codes,
// for the Rangkai Edu authentication system using multiple SMS providers.

package sms

import (
	"fmt"

	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
)

// SendOTPSMS sends an SMS containing the OTP code to the specified phone number
// using the configured SMS provider.
func SendOTPSMS(cfg *config.Config, to, otp string) error {
	// Get the primary SMS provider
	provider, err := cfg.GetPrimarySMSProvider()
	if err != nil {
		return fmt.Errorf("no SMS provider configured: %w", err)
	}

	// Send SMS based on provider type
	switch provider.Type {
	case "twilio":
		return sendTwilioOTPSMS(cfg, provider, to, otp)
	case "sns":
		return sendSNSOTPSMS(cfg, provider, to, otp)
	default:
		return fmt.Errorf("unsupported SMS provider type: %s", provider.Type)
	}
}

// sendTwilioOTPSMS sends an SMS using Twilio
func sendTwilioOTPSMS(cfg *config.Config, provider config.SMSProviderConfig, to, otp string) error {
	// Get Twilio-specific configuration
	// Since provider is already the correct type, we can use it directly

	// Decrypt credentials if needed
	// if cfg.IsEncrypted(provider.AccountSID) {
	// 	if decrypted, err := cfg.DecryptCredential(provider.AccountSID); err == nil {
	// 		provider.AccountSID = decrypted
	// 	}
	// }
	// if cfg.IsEncrypted(provider.AuthToken) {
	// 	if decrypted, err := cfg.DecryptCredential(provider.AuthToken); err == nil {
	// 		provider.AuthToken = decrypted
	// 	}
	// }

	// Set default values if not provided
	// if provider.Settings.FromNumber == "" {
	// 	provider.Settings.FromNumber = "+1234567890"
	// }

	// Create Twilio client
	// client := twilio.NewRestClientWithParams(twilio.ClientParams{
	// 	Username: provider.AccountSID,
	// 	Password: provider.AuthToken,
	// })

	// Create message parameters
	// params := &openapi.CreateMessageParams{}
	// params.SetTo(to)
	// params.SetFrom(provider.Settings.FromNumber)
	// params.SetBody(fmt.Sprintf("Your Rangkai Edu OTP code is: %s. This code expires in 10 minutes. Do not share it with anyone. If you didn't request this, please ignore this message.", otp))

	// Send the SMS
	// resp, err := client.Api.CreateMessage(params)
	// if err != nil {
	// 	return fmt.Errorf("failed to send OTP SMS: %w", err)
	// }

	// Check for Twilio errors
	// if resp.ErrorCode != nil {
	// 	return fmt.Errorf("Twilio error: %s", *resp.ErrorMessage)
	// }

	return nil
}

// sendSNSOTPSMS sends an SMS using Amazon SNS
func sendSNSOTPSMS(cfg *config.Config, provider config.SMSProviderConfig, to, otp string) error {
	// Get SNS-specific configuration
	// Since provider is already the correct type, we can use it directly

	// Decrypt credentials if needed
	// if cfg.IsEncrypted(provider.AccessKeyID) {
	// 	if decrypted, err := cfg.DecryptCredential(provider.AccessKeyID); err == nil {
	// 		provider.AccessKeyID = decrypted
	// 	}
	// }
	// if cfg.IsEncrypted(provider.SecretAccessKey) {
	// 	if decrypted, err := cfg.DecryptCredential(provider.SecretAccessKey); err == nil {
	// 		provider.SecretAccessKey = decrypted
	// 	}
	// }

	// Set default values if not provided
	// if provider.Settings.FromNumber == "" {
	// 	provider.Settings.FromNumber = "+1234567890"
	// }
	// if provider.Region == "" {
	// 	provider.Region = "us-east-1"
	// }

	// TODO: Implement SNS API call
	// This would require the AWS Go SDK
	// For now, return a placeholder implementation
	return fmt.Errorf("SNS implementation not yet implemented")
}

// SendTestSMS sends a test SMS to verify the SMS configuration
func SendTestSMS(cfg *config.Config, to string) error {
	// Get the primary SMS provider
	provider, err := cfg.GetPrimarySMSProvider()
	if err != nil {
		return fmt.Errorf("no SMS provider configured: %w", err)
	}

	// Set default values if not provided
	if provider.Settings.FromNumber == "" {
		provider.Settings.FromNumber = "+1234567890"
	}

	// Send test SMS using the same provider logic
	switch provider.Type {
	case "twilio":
		return sendTwilioOTPSMS(cfg, provider, to, "TEST123")
	case "sns":
		return sendSNSOTPSMS(cfg, provider, to, "TEST123")
	default:
		return fmt.Errorf("unsupported SMS provider type: %s", provider.Type)
	}
}

// ValidatePhoneFormat validates that a phone number is in the correct format
func ValidatePhoneFormat(phone string) bool {
	// Basic validation - phone number should start with + and contain only digits and spaces
	if len(phone) < 10 {
		return false
	}
	
	if phone[0] != '+' {
		return false
	}
	
	// Check remaining characters are digits or spaces
	for _, char := range phone[1:] {
		if char != ' ' && (char < '0' || char > '9') {
			return false
		}
	}
	
	return true
}

// FormatPhoneNumber formats a phone number to E.164 format
func FormatPhoneNumber(phone string) string {
	// Remove all non-digit characters except the leading +
	var formatted string
	for _, char := range phone {
		if char == '+' {
			formatted += "+"
		} else if char >= '0' && char <= '9' {
			formatted += string(char)
		}
	}
	
	return formatted
}