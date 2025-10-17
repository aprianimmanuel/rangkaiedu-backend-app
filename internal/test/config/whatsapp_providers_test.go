
package config_test

import (
	"os"
	"testing"
	
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderManager_LoadWhatsAppProviders(t *testing.T) {
	// Set up test environment variables for WhatsApp Business provider
	os.Setenv("WHATSAPP_PROVIDER_0_TYPE", "whatsapp_business")
	os.Setenv("WHATSAPP_PROVIDER_0_NAME", "test-whatsapp-business")
	os.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "true")
	os.Setenv("WHATSAPP_PROVIDER_0_PRIORITY", "0")
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_PHONE_ID_NUMBER", "1234567890123456")
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_BUSINESS_ACCOUNT_ID", "business_account_123")
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_ACCESS_TOKEN", "test_access_token")
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_API_VERSION", "v15.0")

	// Set up test environment variables for Twilio WhatsApp provider
	os.Setenv("WHATSAPP_PROVIDER_1_TYPE", "twilio_whatsapp")
	os.Setenv("WHATSAPP_PROVIDER_1_NAME", "test-twilio-whatsapp")
	os.Setenv("WHATSAPP_PROVIDER_1_ENABLED", "true")
	os.Setenv("WHATSAPP_PROVIDER_1_PRIORITY", "1")
	os.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_FROM_NUMBER", "+1234567890")
	os.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_ACCOUNT_SID", "test_account_sid")
	os.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_AUTH_TOKEN", "test_auth_token")
	os.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_WHATSAPP_NUMBER", "+1234567890")

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("Failed to load providers: %v", err)
	}

	// Verify WhatsApp providers
	if len(pm.WhatsAppProviders) == 0 {
		t.Fatal("No WhatsApp providers loaded")
	}

	// Verify WhatsApp Business provider
	whatsappBusinessProvider := pm.WhatsAppProviders[0]
	if whatsappBusinessProvider.Name != "test-whatsapp-business" {
		t.Errorf("Expected WhatsApp Business provider name 'test-whatsapp-business', got '%s'", whatsappBusinessProvider.Name)
	}
	if whatsappBusinessProvider.Type != "whatsapp_business" {
		t.Errorf("Expected WhatsApp Business provider type 'whatsapp_business', got '%s'", whatsappBusinessProvider.Type)
	}
	if !whatsappBusinessProvider.Enabled {
		t.Error("Expected WhatsApp Business provider to be enabled")
	}
	if whatsappBusinessProvider.Priority != 0 {
		t.Errorf("Expected WhatsApp Business provider priority 0, got %d", whatsappBusinessProvider.Priority)
	}
	if whatsappBusinessProvider.Settings.FromNumber != "+1234567890" {
		t.Errorf("Expected FromNumber '+1234567890', got '%s'", whatsappBusinessProvider.Settings.FromNumber)
	}

	// Verify Twilio WhatsApp provider
	twilioWhatsAppProvider := pm.WhatsAppProviders[1]
	if twilioWhatsAppProvider.Name != "test-twilio-whatsapp" {
		t.Errorf("Expected Twilio WhatsApp provider name 'test-twilio-whatsapp', got '%s'", twilioWhatsAppProvider.Name)
	}
	if twilioWhatsAppProvider.Type != "twilio_whatsapp" {
		t.Errorf("Expected Twilio WhatsApp provider type 'twilio_whatsapp', got '%s'", twilioWhatsAppProvider.Type)
	}
	if !twilioWhatsAppProvider.Enabled {
		t.Error("Expected Twilio WhatsApp provider to be enabled")
	}
	if twilioWhatsAppProvider.Priority != 1 {
		t.Errorf("Expected Twilio WhatsApp provider priority 1, got %d", twilioWhatsAppProvider.Priority)
	}
	if twilioWhatsAppProvider.Settings.FromNumber != "+1234567890" {
		t.Errorf("Expected FromNumber '+1234567890', got '%s'", twilioWhatsAppProvider.Settings.FromNumber)
	}

	// Clean up
	os.Clearenv()
}

func TestProviderManager_GetPrimaryWhatsAppProvider(t *testing.T) {
	// Set up test environment variables for WhatsApp Business provider (higher priority)
	os.Setenv("WHATSAPP_PROVIDER_0_TYPE", "whatsapp_business")
	os.Setenv("WHATSAPP_PROVIDER_0_NAME", "primary-whatsapp")
	os.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "true")
	os.Setenv("WHATSAPP_PROVIDER_0_PRIORITY", "0")
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_PHONE_ID_NUMBER", "1234567890123456")
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_BUSINESS_ACCOUNT_ID", "business_account_123")
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_ACCESS_TOKEN", "test_access_token")

	// Set up test environment variables for secondary WhatsApp provider
	os.Setenv("WHATSAPP_PROVIDER_1_TYPE", "twilio_whatsapp")
	os.Setenv("WHATSAPP_PROVIDER_1_NAME", "secondary-whatsapp")
	os.Setenv("WHATSAPP_PROVIDER_1_ENABLED", "true")
	os.Setenv("WHATSAPP_PROVIDER_1_PRIORITY", "1")
	os.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_FROM_NUMBER", "+1234567891")
	os.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_ACCOUNT_SID", "test_account_sid")
	os.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_AUTH_TOKEN", "test_auth_token")
	os.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_WHATSAPP_NUMBER", "+1234567891")

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("Failed to load providers: %v", err)
	}

	// Get primary WhatsApp provider
	provider, err := pm.GetPrimaryWhatsAppProvider()
	if err != nil {
		t.Fatalf("Failed to get primary WhatsApp provider: %v", err)
	}

	if provider.Name != "primary-whatsapp" {
		t.Errorf("Expected primary WhatsApp provider name 'primary-whatsapp', got '%s'", provider.Name)
	}
	if provider.Type != "whatsapp_business" {
		t.Errorf("Expected primary WhatsApp provider type 'whatsapp_business', got '%s'", provider.Type)
	}

	// Clean up
	os.Clearenv()
}

func TestProviderManager_GetWhatsAppProviderByName(t *testing.T) {
	// Set up test environment variables for WhatsApp Business provider
	os.Setenv("WHATSAPP_PROVIDER_0_TYPE", "whatsapp_business")
	os.Setenv("WHATSAPP_PROVIDER_0_NAME", "test-whatsapp-business")
	os.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "true")
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_PHONE_ID_NUMBER", "1234567890123456")
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_BUSINESS_ACCOUNT_ID", "business_account_123")
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_ACCESS_TOKEN", "test_access_token")

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("Failed to load providers: %v", err)
	}

	// Get WhatsApp provider by name
	provider, err := pm.GetWhatsAppProviderByName("test-whatsapp-business")
	if err != nil {
		t.Fatalf("Failed to get WhatsApp provider by name: %v", err)
	}

	if provider.Name != "test-whatsapp-business" {
		t.Errorf("Expected WhatsApp provider name 'test-whatsapp-business', got '%s'", provider.Name)
	}
	if provider.Type != "whatsapp_business" {
		t.Errorf("Expected WhatsApp provider type 'whatsapp_business', got '%s'", provider.Type)
	}

	// Test non-existent provider
	_, err = pm.GetWhatsAppProviderByName("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent provider")
	}

	// Clean up
	os.Clearenv()
}

func TestProviderManager_GetPrimaryWhatsAppProvider_NoProviders(t *testing.T) {
	// Create provider manager without any WhatsApp providers
	pm := config.NewProviderManager()

	// Try to get primary WhatsApp provider
	_, err := pm.GetPrimaryWhatsAppProvider()
	if err == nil {
		t.Error("Expected error when no WhatsApp providers are configured")
	}

	// Clean up
	os.Clearenv()
}

func TestProviderManager_GetWhatsAppProviderByName_DisabledProvider(t *testing.T) {
	// Set up test environment variables for disabled WhatsApp provider
	os.Setenv("WHATSAPP_PROVIDER_0_TYPE", "whatsapp_business")
	os.Setenv("WHATSAPP_PROVIDER_0_NAME", "disabled-whatsapp")
	os.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "false") // Disabled
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_PHONE_ID_NUMBER", "1234567890123456")
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_BUSINESS_ACCOUNT_ID", "business_account_123")
	os.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_ACCESS_TOKEN", "test_access_token")

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("Failed to load providers: %v", err)
	}

	// Try to get disabled provider by name
	_, err = pm.GetWhatsAppProviderByName("disabled-whatsapp")
	if err == nil {
		t.Error("Expected error for disabled provider")
	}

	// Clean up
	os.Clearenv()
}

func TestWhatsAppProviderConfig_Validate(t *testing.T) {
	tests := []struct {
		name     string
		provider config.WhatsAppProviderConfig
		wantErr  bool
		errMsg   string
	}{
		{
			name: "Valid provider",
			provider: config.WhatsAppProviderConfig{
				Name:     "test-provider",
				Type:     "whatsapp_business",
				Enabled:  true,
				Priority: 0,
				Settings: config.WhatsAppSettings{
					FromNumber: "+1234567890",
				},
			},
			wantErr: false,
		},
		{
			name: "Missing name",
			provider: config.WhatsAppProviderConfig{
				Name:     "",
				Type:     "whatsapp_business",
				Enabled:  true,
				Priority: 0,
				Settings: config.WhatsAppSettings{
					FromNumber: "+1234567890",
				},
			},
			wantErr: true,
			errMsg:  "WhatsApp provider name is required",
		},
		{
			name: "Missing type",
			provider: config.WhatsAppProviderConfig{
				Name:     "test-provider",
				Type:     "",
				Enabled:  true,
				Priority: 0,
				Settings: config.WhatsAppSettings{
					FromNumber: "+1234567890",
				},
			},
			wantErr: true,
			errMsg:  "WhatsApp provider type is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.provider.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProviderManager_LoadWhatsAppProviders_InvalidType(t *testing.T) {
	// Set up test environment variables with invalid WhatsApp provider type
	os.Setenv("WHATSAPP_PROVIDER_0_TYPE", "invalid-type")
	os.Setenv("WHATSAPP_PROVIDER_0_NAME", "test-invalid")
	os.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "true")

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers - should not fail even with invalid provider
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("LoadProviders should not fail for invalid provider types: %v", err)
	}

	// Verify no WhatsApp providers were loaded
	if len(pm.WhatsAppProviders) != 0 {
		t.Errorf("Expected 0 WhatsApp providers for invalid type, got %d", len(pm.WhatsAppProviders))
	}

	// Clean up
	os.Clearenv()
}

func TestProviderManager_LoadWhatsAppProviders_MissingFields(t *testing.T) {
	// Set up test environment variables with missing required fields for WhatsApp Business
	os.Setenv("WHATSAPP_PROVIDER_0_TYPE", "whatsapp_business")
	os.Setenv("WHATSAPP_PROVIDER_0_NAME", "test-whatsapp")
	os.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "true")
	// Missing required fields like PHONE_ID_NUMBER, BUSINESS_ACCOUNT_ID, ACCESS_TOKEN

	// Create provider manager
	pm := config.NewProviderManager()

	// Load providers - should not fail even with missing fields
	err := pm.LoadProviders()
	if err != nil {
		t.Fatalf("LoadProviders should not fail for missing fields: %v", err)
	}

	// Verify no WhatsApp providers were loaded due to missing required fields
	if len(pm.WhatsAppProviders) != 0 {
		t.Errorf("Expected 0 WhatsApp providers for missing required fields, got %d", len(pm.WhatsAppProviders))
	}

	// Clean up
	os.Clearenv()
}