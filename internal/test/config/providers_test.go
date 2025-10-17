package config_test

import (
	"os"
	"testing"
	
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
)

func TestProviderManager_LoadProviders(t *testing.T) {
	// Set up test environment variables
	os.Setenv("EMAIL_PROVIDER_0_TYPE", "smtp")
	os.Setenv("EMAIL_PROVIDER_0_NAME", "test-smtp")
	os.Setenv("EMAIL_PROVIDER_0_ENABLED", "true")
	os.Setenv("EMAIL_PROVIDER_0_PRIORITY", "0")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_HOST", "smtp.gmail.com")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_PORT", "587")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_USERNAME", "test@gmail.com")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_PASSWORD", "test-password")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_FROM_EMAIL", "test@gmail.com")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_FROM_NAME", "Test Service")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_REPLY_TO", "reply@test.com")

	os.Setenv("SMS_PROVIDER_0_TYPE", "twilio")
	os.Setenv("SMS_PROVIDER_0_NAME", "test-twilio")
	os.Setenv("SMS_PROVIDER_0_ENABLED", "true")
	os.Setenv("SMS_PROVIDER_0_PRIORITY", "0")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_ACCOUNT_SID", "test-account-sid")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_AUTH_TOKEN", "test-auth-token")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("Failed to load providers: %v", err)
	}

	// Verify email provider
	if len(pm.EmailProviders) == 0 {
		t.Fatal("No email providers loaded")
	}

	emailProvider := pm.EmailProviders[0]
	if emailProvider.Name != "test-smtp" {
		t.Errorf("Expected email provider name 'test-smtp', got '%s'", emailProvider.Name)
	}
	if emailProvider.Type != "smtp" {
		t.Errorf("Expected email provider type 'smtp', got '%s'", emailProvider.Type)
	}
	if !emailProvider.Enabled {
		t.Error("Expected email provider to be enabled")
	}
	if emailProvider.Priority != 0 {
		t.Errorf("Expected email provider priority 0, got %d", emailProvider.Priority)
	}

	// Verify SMS provider
	if len(pm.SMSProviders) == 0 {
		t.Fatal("No SMS providers loaded")
	}

	smsProvider := pm.SMSProviders[0]
	if smsProvider.Name != "test-twilio" {
		t.Errorf("Expected SMS provider name 'test-twilio', got '%s'", smsProvider.Name)
	}
	if smsProvider.Type != "twilio" {
		t.Errorf("Expected SMS provider type 'twilio', got '%s'", smsProvider.Type)
	}
	if !smsProvider.Enabled {
		t.Error("Expected SMS provider to be enabled")
	}
	if smsProvider.Priority != 0 {
		t.Errorf("Expected SMS provider priority 0, got %d", smsProvider.Priority)
	}

	// Clean up
	os.Clearenv()
}

func TestProviderManager_GetPrimaryEmailProvider(t *testing.T) {
	// Set up test environment variables
	os.Setenv("EMAIL_PROVIDER_0_TYPE", "smtp")
	os.Setenv("EMAIL_PROVIDER_0_NAME", "primary-smtp")
	os.Setenv("EMAIL_PROVIDER_0_ENABLED", "true")
	os.Setenv("EMAIL_PROVIDER_0_PRIORITY", "0")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_HOST", "smtp.gmail.com")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_PORT", "587")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_USERNAME", "test@gmail.com")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_PASSWORD", "test-password")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_FROM_EMAIL", "test@gmail.com")

	os.Setenv("EMAIL_PROVIDER_1_TYPE", "smtp")
	os.Setenv("EMAIL_PROVIDER_1_NAME", "secondary-smtp")
	os.Setenv("EMAIL_PROVIDER_1_ENABLED", "true")
	os.Setenv("EMAIL_PROVIDER_1_PRIORITY", "1")
	os.Setenv("EMAIL_PROVIDER_1_SETTINGS_HOST", "smtp.gmail.com")
	os.Setenv("EMAIL_PROVIDER_1_SETTINGS_PORT", "587")
	os.Setenv("EMAIL_PROVIDER_1_SETTINGS_USERNAME", "test@gmail.com")
	os.Setenv("EMAIL_PROVIDER_1_SETTINGS_PASSWORD", "test-password")
	os.Setenv("EMAIL_PROVIDER_1_SETTINGS_FROM_EMAIL", "test@gmail.com")

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("Failed to load providers: %v", err)
	}

	// Get primary email provider
	provider, err := pm.GetPrimaryEmailProvider()
	if err != nil {
		t.Fatalf("Failed to get primary email provider: %v", err)
	}

	if provider.Name != "primary-smtp" {
		t.Errorf("Expected primary email provider name 'primary-smtp', got '%s'", provider.Name)
	}

	// Clean up
	os.Clearenv()
}

func TestProviderManager_GetPrimarySMSProvider(t *testing.T) {
	// Set up test environment variables
	os.Setenv("SMS_PROVIDER_0_TYPE", "twilio")
	os.Setenv("SMS_PROVIDER_0_NAME", "primary-twilio")
	os.Setenv("SMS_PROVIDER_0_ENABLED", "true")
	os.Setenv("SMS_PROVIDER_0_PRIORITY", "0")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_ACCOUNT_SID", "test-account-sid")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_AUTH_TOKEN", "test-auth-token")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")

	os.Setenv("SMS_PROVIDER_1_TYPE", "twilio")
	os.Setenv("SMS_PROVIDER_1_NAME", "secondary-twilio")
	os.Setenv("SMS_PROVIDER_1_ENABLED", "true")
	os.Setenv("SMS_PROVIDER_1_PRIORITY", "1")
	os.Setenv("SMS_PROVIDER_1_SETTINGS_ACCOUNT_SID", "test-account-sid")
	os.Setenv("SMS_PROVIDER_1_SETTINGS_AUTH_TOKEN", "test-auth-token")
	os.Setenv("SMS_PROVIDER_1_SETTINGS_FROM_NUMBER", "+1234567890")

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("Failed to load providers: %v", err)
	}

	// Get primary SMS provider
	provider, err := pm.GetPrimarySMSProvider()
	if err != nil {
		t.Fatalf("Failed to get primary SMS provider: %v", err)
	}

	if provider.Name != "primary-twilio" {
		t.Errorf("Expected primary SMS provider name 'primary-twilio', got '%s'", provider.Name)
	}

	// Clean up
	os.Clearenv()
}

func TestProviderManager_GetEmailProviderByName(t *testing.T) {
	// Set up test environment variables
	os.Setenv("EMAIL_PROVIDER_0_TYPE", "smtp")
	os.Setenv("EMAIL_PROVIDER_0_NAME", "test-smtp")
	os.Setenv("EMAIL_PROVIDER_0_ENABLED", "true")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_HOST", "smtp.gmail.com")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_PORT", "587")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_USERNAME", "test@gmail.com")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_PASSWORD", "test-password")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_FROM_EMAIL", "test@gmail.com")

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("Failed to load providers: %v", err)
	}

	// Get email provider by name
	provider, err := pm.GetEmailProviderByName("test-smtp")
	if err != nil {
		t.Fatalf("Failed to get email provider by name: %v", err)
	}

	if provider.Name != "test-smtp" {
		t.Errorf("Expected email provider name 'test-smtp', got '%s'", provider.Name)
	}

	// Test non-existent provider
	_, err = pm.GetEmailProviderByName("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent provider")
	}

	// Clean up
	os.Clearenv()
}

func TestProviderManager_GetSMSProviderByName(t *testing.T) {
	// Set up test environment variables
	os.Setenv("SMS_PROVIDER_0_TYPE", "twilio")
	os.Setenv("SMS_PROVIDER_0_NAME", "test-twilio")
	os.Setenv("SMS_PROVIDER_0_ENABLED", "true")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_ACCOUNT_SID", "test-account-sid")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_AUTH_TOKEN", "test-auth-token")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("Failed to load providers: %v", err)
	}

	// Get SMS provider by name
	provider, err := pm.GetSMSProviderByName("test-twilio")
	if err != nil {
		t.Fatalf("Failed to get SMS provider by name: %v", err)
	}

	if provider.Name != "test-twilio" {
		t.Errorf("Expected SMS provider name 'test-twilio', got '%s'", provider.Name)
	}

	// Test non-existent provider
	_, err = pm.GetSMSProviderByName("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent provider")
	}

	// Clean up
	os.Clearenv()
}

func TestProviderManager_HealthCheck(t *testing.T) {
	// Set up test environment variables
	os.Setenv("EMAIL_PROVIDER_0_TYPE", "smtp")
	os.Setenv("EMAIL_PROVIDER_0_NAME", "test-smtp")
	os.Setenv("EMAIL_PROVIDER_0_ENABLED", "true")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_HOST", "smtp.gmail.com")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_PORT", "587")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_USERNAME", "test@gmail.com")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_PASSWORD", "test-password")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_FROM_EMAIL", "test@example.com")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_FROM_NAME", "Test Service")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_REPLY_TO", "reply@test.com")

	os.Setenv("SMS_PROVIDER_0_TYPE", "twilio")
	os.Setenv("SMS_PROVIDER_0_NAME", "test-twilio")
	os.Setenv("SMS_PROVIDER_0_ENABLED", "true")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_ACCOUNT_SID", "test-account-sid")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_AUTH_TOKEN", "test-auth-token")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("Failed to load providers: %v", err)
	}

	// Perform health check
	health := pm.HealthCheck()

	if len(health) == 0 {
		t.Fatal("No health check results returned")
	}

	// Verify email provider health
	emailHealthFound := false
	for _, h := range health {
		if h.Name == "test-smtp" && h.Type == "smtp" {
			emailHealthFound = true
			if h.Status != "healthy" {
				t.Errorf("Expected email provider status 'healthy', got '%s'", h.Status)
			}
			break
		}
	}

	if !emailHealthFound {
		t.Error("Email provider health check not found")
	}

	// Verify SMS provider health
	smsHealthFound := false
	for _, h := range health {
		if h.Name == "test-twilio" && h.Type == "twilio" {
			smsHealthFound = true
			if h.Status != "healthy" {
				t.Errorf("Expected SMS provider status 'healthy', got '%s'", h.Status)
			}
			break
		}
	}

	if !smsHealthFound {
		t.Error("SMS provider health check not found")
	}

	// Clean up
	os.Clearenv()
}

func TestCredentialManager_EncryptDecrypt(t *testing.T) {
	// Create credential manager
	cm := config.NewCredentialManager()

	// Test value
	testValue := "test-secret-password"

	// Encrypt
	encrypted, err := cm.Encrypt(testValue)
	if err != nil {
		t.Fatalf("Failed to encrypt value: %v", err)
	}

	if encrypted == testValue {
		t.Error("Encrypted value should be different from original")
	}

	// Decrypt
	decrypted, err := cm.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt value: %v", err)
	}

	if decrypted != testValue {
		t.Errorf("Decrypted value '%s' does not match original '%s'", decrypted, testValue)
	}

	// Test with invalid encrypted value
	_, err = cm.Decrypt("invalid-encrypted-value")
	if err == nil {
		t.Error("Expected error for invalid encrypted value")
	}
}

func TestCredentialManager_IsEncrypted(t *testing.T) {
	// Create credential manager
	cm := config.NewCredentialManager()

	// Test with plain text
	plainText := "plain-text-value"
	if cm.IsEncrypted(plainText) {
		t.Error("Plain text should not be detected as encrypted")
	}

	// Test with encrypted value
	encrypted, err := cm.Encrypt("test-value")
	if err != nil {
		t.Fatalf("Failed to encrypt value: %v", err)
	}

	if !cm.IsEncrypted(encrypted) {
		t.Error("Encrypted value should be detected as encrypted")
	}
}

func TestConfig_LoadWithProviders(t *testing.T) {
	// Set up test environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USER", "test_user")
	os.Setenv("DB_PASSWORD", "test_password")

	os.Setenv("EMAIL_PROVIDER_0_TYPE", "smtp")
	os.Setenv("EMAIL_PROVIDER_0_NAME", "test-smtp")
	os.Setenv("EMAIL_PROVIDER_0_ENABLED", "true")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_HOST", "smtp.gmail.com")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_PORT", "587")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_USERNAME", "test@gmail.com")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_PASSWORD", "test-password")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_FROM_EMAIL", "test@gmail.com")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_FROM_NAME", "Test Service")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_REPLY_TO", "reply@test.com")

	os.Setenv("SMS_PROVIDER_0_TYPE", "twilio")
	os.Setenv("SMS_PROVIDER_0_NAME", "test-twilio")
	os.Setenv("SMS_PROVIDER_0_ENABLED", "true")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_ACCOUNT_SID", "test-account-sid")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_AUTH_TOKEN", "test-auth-token")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")

	// Load configuration
	cfg := config.Load()

	// Verify database configuration
	if cfg.DBHost != "localhost" {
		t.Errorf("Expected DB_HOST 'localhost', got '%s'", cfg.DBHost)
	}
	if cfg.DBName != "test_db" {
		t.Errorf("Expected DB_NAME 'test_db', got '%s'", cfg.DBName)
	}

	// Verify provider manager
	if cfg.ProviderManager == nil {
		t.Fatal("Provider manager not initialized")
	}

	// Verify credential manager
	if cfg.CredentialManager == nil {
		t.Fatal("Credential manager not initialized")
	}

	// Verify email provider
	emailProvider, err := cfg.GetPrimaryEmailProvider()
	if err != nil {
		t.Fatalf("Failed to get primary email provider: %v", err)
	}

	if emailProvider.Name != "test-smtp" {
		t.Errorf("Expected email provider name 'test-smtp', got '%s'", emailProvider.Name)
	}

	// Verify SMS provider
	smsProvider, err := cfg.GetPrimarySMSProvider()
	if err != nil {
		t.Fatalf("Failed to get primary SMS provider: %v", err)
	}

	if smsProvider.Name != "test-twilio" {
		t.Errorf("Expected SMS provider name 'test-twilio', got '%s'", smsProvider.Name)
	}

	// Clean up
	os.Clearenv()
}

func TestProviderManager_LoadProviders_InvalidEmailType(t *testing.T) {
	// Set up test environment variables with invalid email provider type
	os.Setenv("EMAIL_PROVIDER_0_TYPE", "invalid-type")
	os.Setenv("EMAIL_PROVIDER_0_NAME", "test-invalid")
	os.Setenv("EMAIL_PROVIDER_0_ENABLED", "true")

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers - should not fail even with invalid provider
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("LoadProviders should not fail for invalid provider types: %v", err)
	}

	// Verify no email providers were loaded
	if len(pm.EmailProviders) != 0 {
		t.Errorf("Expected 0 email providers for invalid type, got %d", len(pm.EmailProviders))
	}

	// Clean up
	os.Clearenv()
}

func TestProviderManager_LoadProviders_InvalidSMSType(t *testing.T) {
	// Set up test environment variables with invalid SMS provider type
	os.Setenv("SMS_PROVIDER_0_TYPE", "invalid-type")
	os.Setenv("SMS_PROVIDER_0_NAME", "test-invalid")
	os.Setenv("SMS_PROVIDER_0_ENABLED", "true")

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers - should not fail even with invalid provider
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("LoadProviders should not fail for invalid provider types: %v", err)
	}

	// Verify no SMS providers were loaded
	if len(pm.SMSProviders) != 0 {
		t.Errorf("Expected 0 SMS providers for invalid type, got %d", len(pm.SMSProviders))
	}

	// Clean up
	os.Clearenv()
}

func TestProviderManager_LoadProviders_MissingEmailFields(t *testing.T) {
	// Set up test environment variables with missing required fields
	os.Setenv("EMAIL_PROVIDER_0_TYPE", "smtp")
	os.Setenv("EMAIL_PROVIDER_0_NAME", "test-smtp")
	os.Setenv("EMAIL_PROVIDER_0_ENABLED", "true")
	// Missing required fields like HOST, USERNAME, PASSWORD

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers - should not fail even with missing fields
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("LoadProviders should not fail for missing fields: %v", err)
	}

	// Verify no email providers were loaded due to missing required fields
	if len(pm.EmailProviders) != 0 {
		t.Errorf("Expected 0 email providers for missing required fields, got %d", len(pm.EmailProviders))
	}

	// Clean up
	os.Clearenv()
}

func TestProviderManager_LoadProviders_MissingSMSFields(t *testing.T) {
	// Set up test environment variables with missing required fields
	os.Setenv("SMS_PROVIDER_0_TYPE", "twilio")
	os.Setenv("SMS_PROVIDER_0_NAME", "test-twilio")
	os.Setenv("SMS_PROVIDER_0_ENABLED", "true")
	// Missing required fields like ACCOUNT_SID, AUTH_TOKEN

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers - should not fail even with missing fields
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("LoadProviders should not fail for missing fields: %v", err)
	}

	// Verify no SMS providers were loaded due to missing required fields
	if len(pm.SMSProviders) != 0 {
		t.Errorf("Expected 0 SMS providers for missing required fields, got %d", len(pm.SMSProviders))
	}

	// Clean up
	os.Clearenv()
}

func TestProviderManager_HealthCheck_DegradedLegacy(t *testing.T) {
	// Set up legacy environment variables
	os.Setenv("SMTP_HOST", "smtp.gmail.com")
	os.Setenv("SMTP_USER", "test@gmail.com")
	os.Setenv("SMTP_PASS", "test-password")
	os.Setenv("TWILIO_ACCOUNT_SID", "test-account-sid")
	os.Setenv("TWILIO_AUTH_TOKEN", "test-auth-token")
	os.Setenv("TWILIO_SENDER_PHONE", "+1234567890")

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("Failed to load providers: %v", err)
	}

	// Perform health check
	health := pm.HealthCheck()

	if len(health) < 2 {
		t.Fatalf("Expected at least 2 health check results for legacy providers, got %d", len(health))
	}

	// Verify legacy email provider health
	emailHealthFound := false
	for _, h := range health {
		if h.Name == "legacy-smtp" && h.Type == "smtp" {
			emailHealthFound = true
			if h.Status != "degraded" {
				t.Errorf("Expected legacy email provider status 'degraded', got '%s'", h.Status)
			}
			break
		}
	}

	if !emailHealthFound {
		t.Error("Legacy email provider health check not found")
	}

	// Verify legacy SMS provider health
	smsHealthFound := false
	for _, h := range health {
		if h.Name == "legacy-twilio" && h.Type == "twilio" {
			smsHealthFound = true
			if h.Status != "degraded" {
				t.Errorf("Expected legacy SMS provider status 'degraded', got '%s'", h.Status)
			}
			break
		}
	}

	if !smsHealthFound {
		t.Error("Legacy SMS provider health check not found")
	}

	// Clean up
	os.Clearenv()
}

func TestProviderManager_HealthCheck_UnhealthyMissingFields(t *testing.T) {
	// Set up test environment variables with missing required fields for health
	os.Setenv("EMAIL_PROVIDER_0_TYPE", "smtp")
	os.Setenv("EMAIL_PROVIDER_0_NAME", "test-smtp")
	os.Setenv("EMAIL_PROVIDER_0_ENABLED", "true")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_HOST", "smtp.gmail.com")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_PORT", "587")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_USERNAME", "test@gmail.com")
	os.Setenv("EMAIL_PROVIDER_0_SETTINGS_PASSWORD", "test-password")
	// Missing FROM_EMAIL which is required for health check

	os.Setenv("SMS_PROVIDER_0_TYPE", "twilio")
	os.Setenv("SMS_PROVIDER_0_NAME", "test-twilio")
	os.Setenv("SMS_PROVIDER_0_ENABLED", "true")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_ACCOUNT_SID", "test-account-sid")
	os.Setenv("SMS_PROVIDER_0_SETTINGS_AUTH_TOKEN", "test-auth-token")
	// Missing FROM_NUMBER which is required for health check

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("Failed to load providers: %v", err)
	}

	// Perform health check
	health := pm.HealthCheck()

	if len(health) < 2 {
		t.Fatalf("Expected at least 2 health check results, got %d", len(health))
	}

	// Verify email provider health is unhealthy
	emailHealthFound := false
	for _, h := range health {
		if h.Name == "test-smtp" && h.Type == "smtp" {
			emailHealthFound = true
			if h.Status != "unhealthy" {
				t.Errorf("Expected email provider status 'unhealthy', got '%s'", h.Status)
			}
			if h.Message != "From email address is required" {
				t.Errorf("Expected email provider message 'From email address is required', got '%s'", h.Message)
			}
			break
		}
	}

	if !emailHealthFound {
		t.Error("Email provider health check not found")
	}

	// Verify SMS provider health is unhealthy
	smsHealthFound := false
	for _, h := range health {
		if h.Name == "test-twilio" && h.Type == "twilio" {
			smsHealthFound = true
			if h.Status != "unhealthy" {
				t.Errorf("Expected SMS provider status 'unhealthy', got '%s'", h.Status)
			}
			if h.Message != "From phone number is required" {
				t.Errorf("Expected SMS provider message 'From phone number is required', got '%s'", h.Message)
			}
			break
		}
	}

	if !smsHealthFound {
		t.Error("SMS provider health check not found")
	}

	// Clean up
	os.Clearenv()
}