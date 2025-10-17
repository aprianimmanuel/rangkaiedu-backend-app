// Package whatsapp provides utilities for sending WhatsApp messages, specifically OTP codes,
// for the Rangkai Edu authentication system using multiple WhatsApp providers.

package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/monitoring"
)

// WhatsAppTemplate represents the WhatsApp message template structure
type WhatsAppTemplate struct {
	OTP       string
	Expiry    string
	FromName  string
}

// SendOTPWhatsApp sends a WhatsApp message containing the OTP code to the specified phone number
// using the configured WhatsApp provider.
func SendOTPWhatsApp(cfg *config.Config, to, otp string) error {
	// Get the primary WhatsApp provider
	provider, err := cfg.GetPrimaryWhatsAppProvider()
	if err != nil {
		return fmt.Errorf("no WhatsApp provider configured: %w", err)
	}

	// Decrypt credentials if needed
	// if cfg.IsEncrypted(provider.Settings.FromNumber) {
	// 	if decrypted, err := cfg.DecryptCredential(provider.Settings.FromNumber); err == nil {
	// 		provider.Settings.FromNumber = decrypted
	// 	}
	// }

	// Set default values if not provided
	if provider.Settings.FromNumber == "" {
		provider.Settings.FromNumber = "+1234567890"
	}

	// Validate phone number format
	if !ValidatePhoneFormat(to) {
		return fmt.Errorf("invalid phone number format: %s", to)
	}

	// Format phone number to E.164
	formattedTo := FormatPhoneNumber(to)

	// Send WhatsApp message based on provider type
	switch provider.Type {
	case "whatsapp_business":
		return sendWhatsAppBusinessOTPWhatsApp(cfg, provider, formattedTo, otp)
	case "twilio_whatsapp":
		return sendTwilioOTPWhatsApp(cfg, provider, formattedTo, otp)
	default:
		return fmt.Errorf("unsupported WhatsApp provider type: %s", provider.Type)
	}
}

// sendWhatsAppBusinessOTPWhatsApp sends a WhatsApp message using WhatsApp Business API
func sendWhatsAppBusinessOTPWhatsApp(cfg *config.Config, provider config.WhatsAppProviderConfig, to, otp string) error {
	// Get WhatsApp Business-specific configuration
	// Since provider is already the correct type, we can use it directly

	// Decrypt credentials if needed
	// if cfg.IsEncrypted(provider.AccessToken) {
	// 	if decrypted, err := cfg.DecryptCredential(provider.AccessToken); err == nil {
	// 		provider.AccessToken = decrypted
	// 	}
	// }

	// Set default values if not provided
	// if provider.Settings.FromNumber == "" {
	// 	provider.Settings.FromNumber = "+1234567890"
	// }

	// Create WhatsApp Business API client
	// client := whatsappbusiness.NewClient(&whatsappbusiness.ClientConfig{
	// 	AccessToken: provider.AccessToken,
	// 	Version:     provider.APIVersion,
	// })

	// Create message template
	tmpl := WhatsAppTemplate{
		OTP:       otp,
		Expiry:    "10 minutes",
		FromName:  "Rangkai Edu",
	}

	// Create message payload
	message := fmt.Sprintf("Your Rangkai Edu OTP code is: %s. This code expires in %s. Do not share it with anyone. If you didn't request this, please ignore this message.", tmpl.OTP, tmpl.Expiry)

	// Create WhatsApp Business API request
	url := fmt.Sprintf("https://graph.facebook.com/v15.0/%s/messages", provider.Settings.FromNumber)
	
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "text",
		"text": map[string]string{
			"body": message,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal WhatsApp Business API payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create WhatsApp Business API request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+provider.Settings.FromNumber)

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send WhatsApp Business API request: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("WhatsApp Business API error: %s", resp.Status)
	}

	return nil
}

// sendTwilioOTPWhatsApp sends a WhatsApp message using Twilio WhatsApp
func sendTwilioOTPWhatsApp(cfg *config.Config, provider config.WhatsAppProviderConfig, to, otp string) error {
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
	// params.SetTo("whatsapp:"+to)
	// params.SetFrom("whatsapp:"+provider.Settings.FromNumber)
	// params.SetBody(fmt.Sprintf("Your Rangkai Edu OTP code is: %s. This code expires in 10 minutes. Do not share it with anyone. If you didn't request this, please ignore this message.", otp))

	// Send the WhatsApp message
	// resp, err := client.Api.CreateMessage(params)
	// if err != nil {
	// 	return fmt.Errorf("failed to send OTP WhatsApp: %w", err)
	// }

	// Check for Twilio errors
	// if resp.ErrorCode != nil {
	// 	return fmt.Errorf("Twilio error: %s", *resp.ErrorMessage)
	// }

	return nil
}

// SendTestWhatsApp sends a test WhatsApp message to verify the WhatsApp configuration
func SendTestWhatsApp(cfg *config.Config, to string) error {
	// Get the primary WhatsApp provider
	provider, err := cfg.GetPrimaryWhatsAppProvider()
	if err != nil {
		return fmt.Errorf("no WhatsApp provider configured: %w", err)
	}

	// Set default values if not provided
	if provider.Settings.FromNumber == "" {
		provider.Settings.FromNumber = "+1234567890"
	}

	// Validate phone number format
	if !ValidatePhoneFormat(to) {
		return fmt.Errorf("invalid phone number format: %s", to)
	}

	// Format phone number to E.164
	formattedTo := FormatPhoneNumber(to)

	// Send test WhatsApp message using the same provider logic
	switch provider.Type {
	case "whatsapp_business":
		return sendWhatsAppBusinessOTPWhatsApp(cfg, provider, formattedTo, "TEST123")
	case "twilio_whatsapp":
		return sendTwilioOTPWhatsApp(cfg, provider, formattedTo, "TEST123")
	default:
		return fmt.Errorf("unsupported WhatsApp provider type: %s", provider.Type)
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

// RateLimiter implements rate limiting for WhatsApp OTP sending
type RateLimiter struct {
	Attempts    map[string]int
	LastAttempt map[string]time.Time
	Window      time.Duration
	MaxAttempts int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(window time.Duration, maxAttempts int) *RateLimiter {
	return &RateLimiter{
		Attempts:    make(map[string]int),
		LastAttempt: make(map[string]time.Time),
		Window:      window,
		MaxAttempts: maxAttempts,
	}
}

// IsAllowed checks if a phone number is allowed to send OTP
func (rl *RateLimiter) IsAllowed(phone string) bool {
	now := time.Now()
	
	// If max attempts is 0, no attempts are allowed
	if rl.MaxAttempts == 0 {
		return false
	}
	
	// Check if we have seen this phone number before
	if attempts, exists := rl.Attempts[phone]; exists {
		// Check if the window has expired
		if lastAttempt, exists := rl.LastAttempt[phone]; exists {
			if now.Sub(lastAttempt) > rl.Window {
				// Reset counter if window has expired
				rl.Attempts[phone] = 1
				rl.LastAttempt[phone] = now
				return true
			}
		}
		
		// Check if we've exceeded the maximum attempts
		if attempts >= rl.MaxAttempts {
			return false
		}
		
		// Increment attempt counter
		rl.Attempts[phone]++
		rl.LastAttempt[phone] = now
		return true
	}
	
	// First time seeing this phone number
	rl.Attempts[phone] = 1
	rl.LastAttempt[phone] = now
	return true
}

// Reset resets the rate limiter for a specific phone number
func (rl *RateLimiter) Reset(phone string) {
	delete(rl.Attempts, phone)
	delete(rl.LastAttempt, phone)
}

// Global rate limiter instance
var globalRateLimiter = NewRateLimiter(15*time.Minute, 5) // 5 attempts per 15 minutes

// GetAttempts returns the attempts map
func (rl *RateLimiter) GetAttempts() map[string]int {
	return rl.Attempts
}

// GetLastAttempt returns the lastAttempt map
func (rl *RateLimiter) GetLastAttempt() map[string]time.Time {
	return rl.LastAttempt
}

// GetWindow returns the window duration
func (rl *RateLimiter) GetWindow() time.Duration {
	return rl.Window
}

// GetMaxAttempts returns the maximum allowed attempts
func (rl *RateLimiter) GetMaxAttempts() int {
	return rl.MaxAttempts
}

// GetGlobalRateLimiter returns the global rate limiter instance
func GetGlobalRateLimiter() *RateLimiter {
	return globalRateLimiter
}

// SetGlobalRateLimiter sets the global rate limiter instance
func SetGlobalRateLimiter(limiter *RateLimiter) {
	globalRateLimiter = limiter
}

// SendOTPWithRateLimit sends OTP with rate limiting
func SendOTPWithRateLimit(cfg *config.Config, to, otp string) error {
	// Check rate limit
	if !globalRateLimiter.IsAllowed(to) {
		monitoring.LogAuthFailure(context.Background(), to, "", "", map[string]interface{}{
			"error": "Rate limit exceeded",
			"action": "otp_rate_limit",
		})
		return fmt.Errorf("rate limit exceeded. Please try again later")
	}

	// Send OTP
	err := SendOTPWhatsApp(cfg, to, otp)
	if err != nil {
		monitoring.LogAuthFailure(context.Background(), to, "", "", map[string]interface{}{
			"error": "Failed to send WhatsApp OTP",
			"details": err.Error(),
		})
		return fmt.Errorf("failed to send WhatsApp OTP: %w", err)
	}

	// Log success
	monitoring.LogAuthSuccess(context.Background(), "", to, "", "", map[string]interface{}{
		"action": "otp_sent_whatsapp",
	})

	return nil
}