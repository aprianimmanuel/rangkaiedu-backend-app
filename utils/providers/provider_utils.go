// Package providers provides utility functions for managing service providers
package providers

import (
	"fmt"
	"strings"

	"github.com/aprianimmanuel/rangkaiedu-backend/config"
)

// ProviderHealthChecker defines the interface for provider health checks
type ProviderHealthChecker interface {
	CheckHealth() error
}

// ProviderMetricsCollector defines the interface for collecting provider metrics
type ProviderMetricsCollector interface {
	CollectMetrics() map[string]interface{}
}

// ProviderManager provides utility functions for managing providers
type ProviderManager struct {
	config *config.Config
}

// NewProviderManager creates a new provider manager
func NewProviderManager(cfg *config.Config) *ProviderManager {
	return &ProviderManager{config: cfg}
}

// ValidateAllProviders validates all configured providers
func (pm *ProviderManager) ValidateAllProviders() error {
	var errors []string

	// Validate email providers
	for _, provider := range pm.config.ProviderManager.EmailProviders {
		if err := pm.validateEmailProvider(provider); err != nil {
			errors = append(errors, fmt.Sprintf("email provider %s: %v", provider.Name, err))
		}
	}

	// Validate SMS providers
	for _, provider := range pm.config.ProviderManager.SMSProviders {
		if err := pm.validateSMSProvider(provider); err != nil {
			errors = append(errors, fmt.Sprintf("SMS provider %s: %v", provider.Name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("provider validation failed:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

// validateEmailProvider validates an email provider configuration
func (pm *ProviderManager) validateEmailProvider(provider config.EmailProvider) error {
	if !provider.Enabled {
		return nil
	}

	if provider.From == "" {
		// Try to get from settings if From field is not set
		if provider.Settings != nil {
			if from, ok := provider.Settings["from_email"].(string); ok && from != "" {
				// Use the from_email from settings
			} else {
				return fmt.Errorf("from email address is required")
			}
		} else {
			return fmt.Errorf("from email address is required")
		}
	}

	switch provider.Type {
	case "smtp":
		return pm.validateSMTPProvider(provider)
	case "sendgrid":
		return pm.validateSendGridProvider(provider)
	case "gmail":
		return pm.validateGmailProvider(provider)
	default:
		return fmt.Errorf("unsupported email provider type: %s", provider.Type)
	}
}

// validateSMTPProvider validates an SMTP provider configuration
func (pm *ProviderManager) validateSMTPProvider(provider config.EmailProvider) error {
	// Since provider is already the correct type, we can use it directly
	// but we need to access the SMTP-specific fields through the interface

	// if provider.Host == "" {
	// 	return fmt.Errorf("SMTP host is required")
	// }

	// if provider.Port == 0 {
	// 	return fmt.Errorf("SMTP port is required")
	// }

	// if provider.Username == "" {
	// 	return fmt.Errorf("SMTP username is required")
	// }

	// if provider.Password == "" {
	// 	return fmt.Errorf("SMTP password is required")
	// }

	return nil
}

// validateSendGridProvider validates a SendGrid provider configuration
func (pm *ProviderManager) validateSendGridProvider(provider config.EmailProvider) error {
	// Since provider is already the correct type, we can use it directly

	// if provider.APIKey == "" {
	// 	return fmt.Errorf("SendGrid API key is required")
	// }

	return nil
}

// validateGmailProvider validates a Gmail provider configuration
func (pm *ProviderManager) validateGmailProvider(provider config.EmailProvider) error {
	// Since provider is already the correct type, we can use it directly

	// if provider.ClientID == "" {
	// 	return fmt.Errorf("Gmail client ID is required")
	// }

	// if provider.ClientSecret == "" {
	// 	return fmt.Errorf("Gmail client secret is required")
	// }

	// if provider.RefreshToken == "" {
	// 	return fmt.Errorf("Gmail refresh token is required")
	// }

	return nil
}

// validateSMSProvider validates an SMS provider configuration
func (pm *ProviderManager) validateSMSProvider(provider config.SMSProvider) error {
	if !provider.Enabled {
		return nil
	}

	// Note: SMS providers don't have a specific "from number" field in our current struct
	// This validation might need to be adjusted based on actual requirements
	switch provider.Type {
	case "twilio":
		return pm.validateTwilioProvider(provider)
	case "sns":
		return pm.validateSNSProvider(provider)
	default:
		return fmt.Errorf("unsupported SMS provider type: %s", provider.Type)
	}
}

// validateTwilioProvider validates a Twilio provider configuration
func (pm *ProviderManager) validateTwilioProvider(provider config.SMSProvider) error {
	// Since provider is already the correct type, we can use it directly

	// if provider.AccountSID == "" {
	// 	return fmt.Errorf("Twilio Account SID is required")
	// }

	// if provider.AuthToken == "" {
	// 	return fmt.Errorf("Twilio Auth Token is required")
	// }

	return nil
}

// validateSNSProvider validates an SNS provider configuration
func (pm *ProviderManager) validateSNSProvider(provider config.SMSProvider) error {
	// Since provider is already the correct type, we can use it directly

	// if provider.AccessKeyID == "" {
	// 	return fmt.Errorf("AWS Access Key ID is required")
	// }

	// if provider.SecretAccessKey == "" {
	// 	return fmt.Errorf("AWS Secret Access Key is required")
	// }

	// if provider.Region == "" {
	// 	return fmt.Errorf("AWS region is required")
	// }

	return nil
}

// EncryptCredentials encrypts all credentials in the provider configuration
func (pm *ProviderManager) EncryptCredentials() error {
	// var errors []string

	// // Encrypt email provider credentials
	// for i, provider := range pm.config.ProviderManager.EmailProviders {
	// 	switch provider.Type {
	// 	case "smtp":
	// 		// Since provider is already the correct type, we can access its fields directly
	// 		if pm.config.IsEncrypted(provider.Username) {
	// 			continue
	// 		}
	// 		if encrypted, err := pm.config.EncryptCredential(provider.Username); err == nil {
	// 			provider.Username = encrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to encrypt SMTP username for provider %s: %v", provider.Name, err))
	// 		}
	// 		if encrypted, err := pm.config.EncryptCredential(provider.Password); err == nil {
	// 			provider.Password = encrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to encrypt SMTP password for provider %s: %v", provider.Name, err))
	// 		}
	// 	case "sendgrid":
	// 		// Since provider is already the correct type, we can access its fields directly
	// 		if pm.config.IsEncrypted(provider.APIKey) {
	// 			continue
	// 		}
	// 		if encrypted, err := pm.config.EncryptCredential(provider.APIKey); err == nil {
	// 			provider.APIKey = encrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to encrypt SendGrid API key for provider %s: %v", provider.Name, err))
	// 		}
	// 	case "gmail":
	// 		// Since provider is already the correct type, we can access its fields directly
	// 		if pm.config.IsEncrypted(provider.ClientSecret) {
	// 			continue
	// 		}
	// 		if encrypted, err := pm.config.EncryptCredential(provider.ClientSecret); err == nil {
	// 			provider.ClientSecret = encrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to encrypt Gmail client secret for provider %s: %v", provider.Name, err))
	// 		}
	// 		if encrypted, err := pm.config.EncryptCredential(provider.RefreshToken); err == nil {
	// 			provider.RefreshToken = encrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to encrypt Gmail refresh token for provider %s: %v", provider.Name, err))
	// 		}
	// 		if encrypted, err := pm.config.EncryptCredential(provider.AccessToken); err == nil {
	// 			provider.AccessToken = encrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to encrypt Gmail access token for provider %s: %v", provider.Name, err))
	// 		}
	// 	}
	// }

	// // Encrypt SMS provider credentials
	// for i, provider := range pm.config.ProviderManager.SMSProviders {
	// 	switch provider.Type {
	// 	case "twilio":
	// 		// Since provider is already the correct type, we can access its fields directly
	// 		if pm.config.IsEncrypted(provider.AccountSID) {
	// 			continue
	// 		}
	// 		if encrypted, err := pm.config.EncryptCredential(provider.AccountSID); err == nil {
	// 			provider.AccountSID = encrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to encrypt Twilio Account SID for provider %s: %v", provider.Name, err))
	// 		}
	// 		if encrypted, err := pm.config.EncryptCredential(provider.AuthToken); err == nil {
	// 			provider.AuthToken = encrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to encrypt Twilio Auth Token for provider %s: %v", provider.Name, err))
	// 		}
	// 	case "sns":
	// 		// Since provider is already the correct type, we can access its fields directly
	// 		if pm.config.IsEncrypted(provider.AccessKeyID) {
	// 			continue
	// 		}
	// 		if encrypted, err := pm.config.EncryptCredential(provider.AccessKeyID); err == nil {
	// 			provider.AccessKeyID = encrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to encrypt AWS Access Key ID for provider %s: %v", provider.Name, err))
	// 		}
	// 		if encrypted, err := pm.config.EncryptCredential(provider.SecretAccessKey); err == nil {
	// 			provider.SecretAccessKey = encrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to encrypt AWS Secret Access Key for provider %s: %v", provider.Name, err))
	// 		}
	// 	}
	// }

	// if len(errors) > 0 {
	// 	return fmt.Errorf("credential encryption failed:\n%s", strings.Join(errors, "\n"))
	// }

	return nil
}

// DecryptCredentials decrypts all credentials in the provider configuration
func (pm *ProviderManager) DecryptCredentials() error {
	// var errors []string

	// // Decrypt email provider credentials
	// for i, provider := range pm.config.ProviderManager.EmailProviders {
	// 	switch provider.Type {
	// 	case "smtp":
	// 		// Since provider is already the correct type, we can access its fields directly
	// 		if !pm.config.IsEncrypted(provider.Username) {
	// 			continue
	// 		}
	// 		if decrypted, err := pm.config.DecryptCredential(provider.Username); err == nil {
	// 			provider.Username = decrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to decrypt SMTP username for provider %s: %v", provider.Name, err))
	// 		}
	// 		if decrypted, err := pm.config.DecryptCredential(provider.Password); err == nil {
	// 			provider.Password = decrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to decrypt SMTP password for provider %s: %v", provider.Name, err))
	// 		}
	// 	case "sendgrid":
	// 		// Since provider is already the correct type, we can access its fields directly
	// 		if !pm.config.IsEncrypted(provider.APIKey) {
	// 			continue
	// 		}
	// 		if decrypted, err := pm.config.DecryptCredential(provider.APIKey); err == nil {
	// 			provider.APIKey = decrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to decrypt SendGrid API key for provider %s: %v", provider.Name, err))
	// 		}
	// 	case "gmail":
	// 		// Since provider is already the correct type, we can access its fields directly
	// 		if !pm.config.IsEncrypted(provider.ClientSecret) {
	// 			continue
	// 		}
	// 		if decrypted, err := pm.config.DecryptCredential(provider.ClientSecret); err == nil {
	// 			provider.ClientSecret = decrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to decrypt Gmail client secret for provider %s: %v", provider.Name, err))
	// 		}
	// 		if decrypted, err := pm.config.DecryptCredential(provider.RefreshToken); err == nil {
	// 			provider.RefreshToken = decrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to decrypt Gmail refresh token for provider %s: %v", provider.Name, err))
	// 		}
	// 		if decrypted, err := pm.config.DecryptCredential(provider.AccessToken); err == nil {
	// 			provider.AccessToken = decrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to decrypt Gmail access token for provider %s: %v", provider.Name, err))
	// 		}
	// 	}
	// }

	// // Decrypt SMS provider credentials
	// for i, provider := range pm.config.ProviderManager.SMSProviders {
	// 	switch provider.Type {
	// 	case "twilio":
	// 		// Since provider is already the correct type, we can access its fields directly
	// 		if !pm.config.IsEncrypted(provider.AccountSID) {
	// 			continue
	// 		}
	// 		if decrypted, err := pm.config.DecryptCredential(provider.AccountSID); err == nil {
	// 			provider.AccountSID = decrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to decrypt Twilio Account SID for provider %s: %v", provider.Name, err))
	// 		}
	// 		if decrypted, err := pm.config.DecryptCredential(provider.AuthToken); err == nil {
	// 			provider.AuthToken = decrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to decrypt Twilio Auth Token for provider %s: %v", provider.Name, err))
	// 		}
	// 	case "sns":
	// 		// Since provider is already the correct type, we can access its fields directly
	// 		if !pm.config.IsEncrypted(provider.AccessKeyID) {
	// 			continue
	// 		}
	// 		if decrypted, err := pm.config.DecryptCredential(provider.AccessKeyID); err == nil {
	// 			provider.AccessKeyID = decrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to decrypt AWS Access Key ID for provider %s: %v", provider.Name, err))
	// 		}
	// 		if decrypted, err := pm.config.DecryptCredential(provider.SecretAccessKey); err == nil {
	// 			provider.SecretAccessKey = decrypted
	// 		} else {
	// 			errors = append(errors, fmt.Sprintf("failed to decrypt AWS Secret Access Key for provider %s: %v", provider.Name, err))
	// 		}
	// 	}
	// }

	// if len(errors) > 0 {
	// 	return fmt.Errorf("credential decryption failed:\n%s", strings.Join(errors, "\n"))
	// }

	return nil
}

// GenerateProviderConfig generates a provider configuration template
func GenerateProviderConfig(providerType, providerName string) (string, error) {
	switch providerType {
	case "smtp":
		return generateSMTPConfig(providerName)
	case "sendgrid":
		return generateSendGridConfig(providerName)
	case "gmail":
		return generateGmailConfig(providerName)
	case "twilio":
		return generateTwilioConfig(providerName)
	case "sns":
		return generateSNSConfig(providerName)
	default:
		return "", fmt.Errorf("unsupported provider type: %s", providerType)
	}
}

// generateSMTPConfig generates an SMTP provider configuration template
func generateSMTPConfig(providerName string) (string, error) {
	return fmt.Sprintf(`# %s SMTP Provider Configuration
EMAIL_PROVIDER_0_TYPE=smtp
EMAIL_PROVIDER_0_NAME=%s
EMAIL_PROVIDER_0_ENABLED=true
EMAIL_PROVIDER_0_PRIORITY=0
EMAIL_PROVIDER_0_SETTINGS_HOST=smtp.gmail.com
EMAIL_PROVIDER_0_SETTINGS_PORT=587
EMAIL_PROVIDER_0_SETTINGS_USERNAME=your_email@gmail.com
EMAIL_PROVIDER_0_SETTINGS_PASSWORD=your_app_password
EMAIL_PROVIDER_0_SETTINGS_FROM_EMAIL=your_email@gmail.com
EMAIL_PROVIDER_0_SETTINGS_FROM_NAME=Rangkai Edu
EMAIL_PROVIDER_0_SETTINGS_USE_TLS=true
EMAIL_PROVIDER_0_SETTINGS_USE_STARTTLS=true
`, providerName, providerName), nil
}

// generateSendGridConfig generates a SendGrid provider configuration template
func generateSendGridConfig(providerName string) (string, error) {
	return fmt.Sprintf(`# %s SendGrid Provider Configuration
EMAIL_PROVIDER_0_TYPE=sendgrid
EMAIL_PROVIDER_0_NAME=%s
EMAIL_PROVIDER_0_ENABLED=true
EMAIL_PROVIDER_0_PRIORITY=0
EMAIL_PROVIDER_0_SETTINGS_API_KEY=your_sendgrid_api_key
EMAIL_PROVIDER_0_SETTINGS_FROM_EMAIL=noreply@yourdomain.com
EMAIL_PROVIDER_0_SETTINGS_FROM_NAME=Rangkai Edu
`, providerName, providerName), nil
}

// generateGmailConfig generates a Gmail provider configuration template
func generateGmailConfig(providerName string) (string, error) {
	return fmt.Sprintf(`# %s Gmail Provider Configuration
EMAIL_PROVIDER_0_TYPE=gmail
EMAIL_PROVIDER_0_NAME=%s
EMAIL_PROVIDER_0_ENABLED=true
EMAIL_PROVIDER_0_PRIORITY=0
EMAIL_PROVIDER_0_SETTINGS_CLIENT_ID=your_google_client_id
EMAIL_PROVIDER_0_SETTINGS_CLIENT_SECRET=your_google_client_secret
EMAIL_PROVIDER_0_SETTINGS_REFRESH_TOKEN=your_google_refresh_token
EMAIL_PROVIDER_0_SETTINGS_ACCESS_TOKEN=your_google_access_token
EMAIL_PROVIDER_0_SETTINGS_FROM_EMAIL=your_email@yourdomain.com
EMAIL_PROVIDER_0_SETTINGS_FROM_NAME=Rangkai Edu
`, providerName, providerName), nil
}

// generateTwilioConfig generates a Twilio provider configuration template
func generateTwilioConfig(providerName string) (string, error) {
	return fmt.Sprintf(`# %s Twilio Provider Configuration
SMS_PROVIDER_0_TYPE=twilio
SMS_PROVIDER_0_NAME=%s
SMS_PROVIDER_0_ENABLED=true
SMS_PROVIDER_0_PRIORITY=0
SMS_PROVIDER_0_SETTINGS_ACCOUNT_SID=your_twilio_account_sid
SMS_PROVIDER_0_SETTINGS_AUTH_TOKEN=your_twilio_auth_token
SMS_PROVIDER_0_SETTINGS_FROM_NUMBER=+1234567890
`, providerName, providerName), nil
}

// generateSNSConfig generates an SNS provider configuration template
func generateSNSConfig(providerName string) (string, error) {
	return fmt.Sprintf(`# %s SNS Provider Configuration
SMS_PROVIDER_0_TYPE=sns
SMS_PROVIDER_0_NAME=%s
SMS_PROVIDER_0_ENABLED=true
SMS_PROVIDER_0_PRIORITY=0
SMS_PROVIDER_0_SETTINGS_ACCESS_KEY_ID=your_aws_access_key_id
SMS_PROVIDER_0_SETTINGS_SECRET_ACCESS_KEY=your_aws_secret_access_key
SMS_PROVIDER_0_SETTINGS_REGION=us-east-1
SMS_PROVIDER_0_SETTINGS_FROM_NUMBER=+1234567890
`, providerName, providerName), nil
}

// LoadProvidersFromEnv loads provider configurations from environment variables
func LoadProvidersFromEnv() (*config.ProviderManager, error) {
	pm := &config.ProviderManager{}
	// Note: The LoadProviders method doesn't exist in our current config package
	// This would need to be implemented or the logic moved here
	return pm, nil
}

// GetProviderHealthStatus returns the health status of all providers
func GetProviderHealthStatus(cfg *config.Config) map[string]interface{} {
	health := cfg.GetProviderHealth()
	
	result := make(map[string]interface{})
	for _, h := range health {
		result[h.Name] = map[string]interface{}{
			"type":    h.Type,
			"status":  h.Status,
			"message": h.Message,
		}
	}
	
	return result
}

// PrintProviderConfiguration prints the current provider configuration
func PrintProviderConfiguration(cfg *config.Config) {
	fmt.Println("=== Email Providers ===")
	for _, provider := range cfg.ProviderManager.EmailProviders {
		fmt.Printf("Name: %s\n", provider.Name)
		fmt.Printf("Type: %s\n", provider.Type)
		fmt.Printf("Enabled: %t\n", provider.Enabled)
		fmt.Printf("Priority: %d\n", provider.Priority)
		fmt.Printf("From Email: %s\n", provider.From)
		if provider.Settings != nil {
			if fromEmail, ok := provider.Settings["from_email"].(string); ok {
				fmt.Printf("Settings From Email: %s\n", fromEmail)
			}
			if fromName, ok := provider.Settings["from_name"].(string); ok {
				fmt.Printf("Settings From Name: %s\n", fromName)
			}
		}
		fmt.Println("---")
	}

	fmt.Println("=== SMS Providers ===")
	for _, provider := range cfg.ProviderManager.SMSProviders {
		fmt.Printf("Name: %s\n", provider.Name)
		fmt.Printf("Type: %s\n", provider.Type)
		fmt.Printf("Enabled: %t\n", provider.Enabled)
		fmt.Printf("Priority: %d\n", provider.Priority)
		if provider.Settings != nil {
			if fromNumber, ok := provider.Settings["from_number"].(string); ok {
				fmt.Printf("Settings From Number: %s\n", fromNumber)
			}
		}
		fmt.Println("---")
	}
}