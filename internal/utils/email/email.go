
// Package email provides utilities for sending emails, specifically OTP codes,
// for the Rangkai Edu authentication system using multiple email providers.

package email

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
)

// EmailTemplate represents the email template structure
type EmailTemplate struct {
	OTP       string
	Expiry    string
	FromEmail string
	FromName  string
}

// SendOTPEmail sends an email containing the OTP code to the specified recipient
// using the configured email provider.
func SendOTPEmail(cfg *config.Config, to, otp string) error {
	// Get the primary email provider
	provider, err := cfg.GetPrimaryEmailProvider()
	if err != nil {
		// Fallback to simple SMTP if provider configuration fails
		return sendSimpleSMTPOTPEmail(cfg, to, otp)
	}

	// Decrypt credentials if needed
	if cfg.IsEncrypted(provider.Settings.FromEmail) {
		if decrypted, err := cfg.DecryptCredential(provider.Settings.FromEmail); err == nil {
			provider.Settings.FromEmail = decrypted
		}
	}

	// Set default values if not provided
	if provider.Settings.FromEmail == "" {
		provider.Settings.FromEmail = "noreply@rangkaiedu.com"
	}
	if provider.Settings.FromName == "" {
		provider.Settings.FromName = "Rangkai Edu"
	}

	// Send email based on provider type
	switch provider.Type {
	case "smtp":
		return sendSMTPOTPEmail(cfg, provider, to, otp)
	case "sendgrid":
		return sendSendGridOTPEmail(cfg, provider, to, otp)
	case "gmail":
		return sendGmailOTPEmail(cfg, provider, to, otp)
	default:
		// Fallback to simple SMTP if provider type is unsupported
		return sendSimpleSMTPOTPEmail(cfg, to, otp)
	}
}

// sendSimpleSMTPOTPEmail sends an email using simple SMTP configuration
func sendSimpleSMTPOTPEmail(cfg *config.Config, to, otp string) error {
	// Check if SMTP configuration is available
	if cfg.SMTPHost == "" || cfg.SMTPUser == "" || cfg.SMTPPassword == "" {
		return fmt.Errorf("SMTP configuration not complete")
	}

	// Create email message
	from := cfg.SMTPUser
	subject := "Your OTP Code"
	body := fmt.Sprintf("Your OTP code is: %s\n\nThis code will expire in 5 minutes.", otp)
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)

	// Send email using SMTP
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPHost)

	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// sendSMTPOTPEmail sends an email using SMTP with template
func sendSMTPOTPEmail(cfg *config.Config, provider config.EmailProviderConfig, to, otp string) error {
	// Create email template
	tmpl := template.Must(template.New("otp").Parse(otpEmailTemplate))
	emailData := EmailTemplate{
		OTP:       otp,
		Expiry:    "10 minutes",
		FromEmail: provider.Settings.FromEmail,
		FromName:  provider.Settings.FromName,
	}
	
	var body bytes.Buffer
	if err := tmpl.Execute(&body, emailData); err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	// For now, use the simple SMTP implementation
	// TODO: Implement full SMTP with proper configuration
	return sendSimpleSMTPOTPEmail(cfg, to, otp)
}

// sendSendGridOTPEmail sends an email using SendGrid
func sendSendGridOTPEmail(cfg *config.Config, provider config.EmailProviderConfig, to, otp string) error {
	// Get SendGrid-specific configuration
	// Since provider is already the correct type, we can use it directly

	// Decrypt API key if needed
	// if cfg.IsEncrypted(provider.APIKey) {
	// 	if decrypted, err := cfg.DecryptCredential(provider.APIKey); err == nil {
	// 		provider.APIKey = decrypted
	// 	}
	// }

	// TODO: Implement SendGrid API call
	// This would require the SendGrid Go SDK
	// For now, return a placeholder implementation
	return fmt.Errorf("SendGrid implementation not yet implemented")
}

// sendGmailOTPEmail sends an email using Gmail OAuth2
func sendGmailOTPEmail(cfg *config.Config, provider config.EmailProviderConfig, to, otp string) error {
	// Get Gmail-specific configuration
	// Since provider is already the correct type, we can use it directly

	// Decrypt credentials if needed
	// if cfg.IsEncrypted(provider.ClientSecret) {
	// 	if decrypted, err := cfg.DecryptCredential(provider.ClientSecret); err == nil {
	// 		provider.ClientSecret = decrypted
	// 	}
	// }
	// if cfg.IsEncrypted(provider.RefreshToken) {
	// 	if decrypted, err := cfg.DecryptCredential(provider.RefreshToken); err == nil {
	// 		provider.RefreshToken = decrypted
	// 	}
	// }
	// if cfg.IsEncrypted(provider.AccessToken) {
	// 	if decrypted, err := cfg.DecryptCredential(provider.AccessToken); err == nil {
	// 		provider.AccessToken = decrypted
	// 	}
	// }

	// TODO: Implement Gmail OAuth2 API call
	// This would require the Google Go SDK
	// For now, return a placeholder implementation
	return fmt.Errorf("Gmail OAuth2 implementation not yet implemented")
}

// otpEmailTemplate is the template for OTP emails
const otpEmailTemplate = `Dear User,

Your OTP code is: {{.OTP}}

This code expires in {{.Expiry}}. Do not share it with anyone.

If you didn't request this, please ignore this email.

Best regards,
{{.FromName}}`

// SendTestEmail sends a test email to verify the email configuration
func SendTestEmail(cfg *config.Config, to string) error {
	// Get the primary email provider
	provider, err := cfg.GetPrimaryEmailProvider()
	if err != nil {
		// Fallback to simple SMTP if provider configuration fails
		return sendSimpleSMTPTestEmail(cfg, to)
	}

	// Set default values if not provided
	if provider.Settings.FromEmail == "" {
		provider.Settings.FromEmail = "noreply@rangkaiedu.com"
	}
	if provider.Settings.FromName == "" {
		provider.Settings.FromName = "Rangkai Edu"
	}

	// Create test email template
	tmpl := template.Must(template.New("test").Parse(testEmailTemplate))
	emailData := EmailTemplate{
		OTP:       "TEST123",
		Expiry:    "10 minutes",
		FromEmail: provider.Settings.FromEmail,
		FromName:  provider.Settings.FromName,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, emailData); err != nil {
		return fmt.Errorf("failed to render test email template: %w", err)
	}

	// Send test email using the same provider logic
	switch provider.Type {
	case "smtp":
		return sendSMTPOTPEmail(cfg, provider, to, "TEST123")
	case "sendgrid":
		return sendSendGridOTPEmail(cfg, provider, to, "TEST123")
	case "gmail":
		return sendGmailOTPEmail(cfg, provider, to, "TEST123")
	default:
		// Fallback to simple SMTP if provider type is unsupported
		return sendSimpleSMTPTestEmail(cfg, to)
	}
}

// sendSimpleSMTPTestEmail sends a test email using simple SMTP configuration
func sendSimpleSMTPTestEmail(cfg *config.Config, to string) error {
	// Check if SMTP configuration is available
	if cfg.SMTPHost == "" || cfg.SMTPUser == "" || cfg.SMTPPassword == "" {
		return fmt.Errorf("SMTP configuration not complete")
	}

	// Create test email message
	from := cfg.SMTPUser
	subject := "Test Email from Rangkai Edu"
	body := "This is a test email to verify your SMTP configuration."
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)

	// Send email using SMTP
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPHost)

	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send test email: %w", err)
	}

	return nil
}

// testEmailTemplate is the template for test emails
const testEmailTemplate = `Dear User,

This is a test email from Rangkai Edu to verify your email configuration.

If you receive this email, your email configuration is working correctly.

Best regards,
{{.FromName}}`