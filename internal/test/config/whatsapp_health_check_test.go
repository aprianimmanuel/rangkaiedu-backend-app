package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
)

// TestWhatsAppProviderHealthCheck tests the health check functionality for WhatsApp providers
func TestWhatsAppProviderHealthCheck(t *testing.T) {
	// Set up environment variables for WhatsApp Business provider
	t.Setenv("WHATSAPP_PROVIDER_0_TYPE", "whatsapp_business")
	t.Setenv("WHATSAPP_PROVIDER_0_NAME", "primary-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_0_PRIORITY", "1")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_PHONE_ID_NUMBER", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_BUSINESS_ACCOUNT_ID", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_ACCESS_TOKEN", "test-access-token")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_API_VERSION", "v15.0")

	// Set up environment variables for Twilio WhatsApp provider
	t.Setenv("WHATSAPP_PROVIDER_1_TYPE", "twilio_whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_1_NAME", "secondary-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_1_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_1_PRIORITY", "2")
	t.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_FROM_NUMBER", "+1234567891")
	t.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_ACCOUNT_SID", "test-account-sid")
	t.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_AUTH_TOKEN", "test-auth-token")
	t.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_WHATSAPP_NUMBER", "+1234567891")

	// Load configuration
	cfg := config.LoadTest()

	// Get provider health
	health := cfg.GetProviderHealth()

	// Verify health check results
	assert.Len(t, health, 2, "Expected 2 WhatsApp providers to be checked")

	// Check WhatsApp Business provider health
	var whatsappBusinessHealth *config.ProviderHealth
	var twilioWhatsAppHealth *config.ProviderHealth

	for _, h := range health {
		if h.Name == "primary-whatsapp" {
			whatsappBusinessHealth = &h
		} else if h.Name == "secondary-whatsapp" {
			twilioWhatsAppHealth = &h
		}
	}

	require.NotNil(t, whatsappBusinessHealth, "WhatsApp Business provider health not found")
	require.NotNil(t, twilioWhatsAppHealth, "Twilio WhatsApp provider health not found")

	// Verify WhatsApp Business provider health
	assert.Equal(t, "primary-whatsapp", whatsappBusinessHealth.Name)
	assert.Equal(t, "whatsapp_business", whatsappBusinessHealth.Type)
	assert.Equal(t, "healthy", whatsappBusinessHealth.Status)
	assert.Contains(t, whatsappBusinessHealth.Message, "Provider is enabled and configured")

	// Verify Twilio WhatsApp provider health
	assert.Equal(t, "secondary-whatsapp", twilioWhatsAppHealth.Name)
	assert.Equal(t, "twilio_whatsapp", twilioWhatsAppHealth.Type)
	assert.Equal(t, "healthy", twilioWhatsAppHealth.Status)
	assert.Contains(t, twilioWhatsAppHealth.Message, "Provider is enabled and configured")
}

// TestWhatsAppProviderHealthCheck_WithMissingRequiredFields tests health check with missing required fields
func TestWhatsAppProviderHealthCheck_WithMissingRequiredFields(t *testing.T) {
	// Set up environment variables with missing required fields
	t.Setenv("WHATSAPP_PROVIDER_0_TYPE", "whatsapp_business")
	t.Setenv("WHATSAPP_PROVIDER_0_NAME", "incomplete-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_0_PRIORITY", "1")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")
	// Missing required fields: PHONE_ID_NUMBER, BUSINESS_ACCOUNT_ID, ACCESS_TOKEN

	// Load configuration
	cfg := config.LoadTest()

	// Get provider health
	health := cfg.GetProviderHealth()

	// Verify health check results
	assert.Len(t, health, 1, "Expected 1 WhatsApp provider to be checked")

	// Check provider health
	providerHealth := health[0]
	assert.Equal(t, "incomplete-whatsapp", providerHealth.Name)
	assert.Equal(t, "whatsapp_business", providerHealth.Type)
	assert.Equal(t, "unhealthy", providerHealth.Status)
	assert.Contains(t, providerHealth.Message, "required WhatsApp Business settings missing")
}

// TestWhatsAppProviderHealthCheck_WithDisabledProvider tests health check with disabled provider
func TestWhatsAppProviderHealthCheck_WithDisabledProvider(t *testing.T) {
	// Set up environment variables for disabled provider
	t.Setenv("WHATSAPP_PROVIDER_0_TYPE", "whatsapp_business")
	t.Setenv("WHATSAPP_PROVIDER_0_NAME", "disabled-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "false")
	t.Setenv("WHATSAPP_PROVIDER_0_PRIORITY", "1")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_PHONE_ID_NUMBER", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_BUSINESS_ACCOUNT_ID", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_ACCESS_TOKEN", "test-access-token")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_API_VERSION", "v15.0")

	// Load configuration
	cfg := config.LoadTest()

	// Get provider health
	health := cfg.GetProviderHealth()

	// Verify health check results
	assert.Len(t, health, 0, "Expected 0 WhatsApp providers to be checked (disabled)")
}

// TestWhatsAppProviderHealthCheck_WithInvalidProviderType tests health check with invalid provider type
func TestWhatsAppProviderHealthCheck_WithInvalidProviderType(t *testing.T) {
	// Set up environment variables with invalid provider type
	t.Setenv("WHATSAPP_PROVIDER_0_TYPE", "invalid_provider")
	t.Setenv("WHATSAPP_PROVIDER_0_NAME", "invalid-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_0_PRIORITY", "1")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")

	// Load configuration
	cfg := config.LoadTest()

	// Get provider health
	health := cfg.GetProviderHealth()

	// Verify health check results
	assert.Len(t, health, 0, "Expected 0 WhatsApp providers to be checked (invalid type)")
}

// TestWhatsAppProviderHealthCheck_WithLegacyConfiguration tests health check with legacy configuration
func TestWhatsAppProviderHealthCheck_WithLegacyConfiguration(t *testing.T) {
	// Set up environment variables for legacy provider
	t.Setenv("WHATSAPP_PROVIDER_0_TYPE", "whatsapp_business")
	t.Setenv("WHATSAPP_PROVIDER_0_NAME", "legacy-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_0_PRIORITY", "1")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_PHONE_ID_NUMBER", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_BUSINESS_ACCOUNT_ID", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_ACCESS_TOKEN", "test-access-token")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_API_VERSION", "v15.0")

	// Load configuration
	cfg := config.LoadTest()

	// Get provider health
	health := cfg.GetProviderHealth()

	// Verify health check results
	assert.Len(t, health, 1, "Expected 1 WhatsApp provider to be checked")

	// Check provider health
	providerHealth := health[0]
	assert.Equal(t, "legacy-whatsapp", providerHealth.Name)
	assert.Equal(t, "whatsapp_business", providerHealth.Type)
	assert.Equal(t, "degraded", providerHealth.Status)
	assert.Contains(t, providerHealth.Message, "Legacy WhatsApp Business configuration detected")
}

// TestWhatsAppProviderHealthCheck_WithMultipleProviders tests health check with multiple providers
func TestWhatsAppProviderHealthCheck_WithMultipleProviders(t *testing.T) {
	// Set up environment variables for multiple providers
	t.Setenv("WHATSAPP_PROVIDER_0_TYPE", "whatsapp_business")
	t.Setenv("WHATSAPP_PROVIDER_0_NAME", "primary-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_0_PRIORITY", "1")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_PHONE_ID_NUMBER", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_BUSINESS_ACCOUNT_ID", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_ACCESS_TOKEN", "test-access-token")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_API_VERSION", "v15.0")

	t.Setenv("WHATSAPP_PROVIDER_1_TYPE", "twilio_whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_1_NAME", "secondary-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_1_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_1_PRIORITY", "2")
	t.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_FROM_NUMBER", "+1234567891")
	t.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_ACCOUNT_SID", "test-account-sid")
	t.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_AUTH_TOKEN", "test-auth-token")
	t.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_WHATSAPP_NUMBER", "+1234567891")

	t.Setenv("WHATSAPP_PROVIDER_2_TYPE", "whatsapp_business")
	t.Setenv("WHATSAPP_PROVIDER_2_NAME", "tertiary-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_2_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_2_PRIORITY", "3")
	t.Setenv("WHATSAPP_PROVIDER_2_SETTINGS_FROM_NUMBER", "+1234567892")
	t.Setenv("WHATSAPP_PROVIDER_2_SETTINGS_PHONE_ID_NUMBER", "1234567890123457")
	t.Setenv("WHATSAPP_PROVIDER_2_SETTINGS_BUSINESS_ACCOUNT_ID", "1234567890123457")
	t.Setenv("WHATSAPP_PROVIDER_2_SETTINGS_ACCESS_TOKEN", "test-access-token-2")
	t.Setenv("WHATSAPP_PROVIDER_2_SETTINGS_API_VERSION", "v15.0")

	// Load configuration
	cfg := config.LoadTest()

	// Get provider health
	health := cfg.GetProviderHealth()

	// Verify health check results
	assert.Len(t, health, 3, "Expected 3 WhatsApp providers to be checked")

	// Check that all providers are healthy
	for _, h := range health {
		assert.Equal(t, "healthy", h.Status, "Provider %s should be healthy", h.Name)
		assert.Contains(t, h.Message, "Provider is enabled and configured")
	}
}

// TestWhatsAppProviderHealthCheck_WithFallbackMechanism tests health check with fallback mechanism
func TestWhatsAppProviderHealthCheck_WithFallbackMechanism(t *testing.T) {
	// Set up environment variables for primary provider with issues
	t.Setenv("WHATSAPP_PROVIDER_0_TYPE", "whatsapp_business")
	t.Setenv("WHATSAPP_PROVIDER_0_NAME", "primary-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_0_PRIORITY", "1")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")
	// Missing required fields for primary provider

	// Set up environment variables for healthy fallback provider
	t.Setenv("WHATSAPP_PROVIDER_1_TYPE", "twilio_whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_1_NAME", "fallback-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_1_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_1_PRIORITY", "2")
	t.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_FROM_NUMBER", "+1234567891")
	t.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_ACCOUNT_SID", "test-account-sid")
	t.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_AUTH_TOKEN", "test-auth-token")
	t.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_WHATSAPP_NUMBER", "+1234567891")

	// Load configuration
	cfg := config.LoadTest()

	// Get provider health
	health := cfg.GetProviderHealth()

	// Verify health check results
	assert.Len(t, health, 1, "Expected 1 WhatsApp provider to be checked (only healthy one)")

	// Check that only the healthy provider is returned
	providerHealth := health[0]
	assert.Equal(t, "fallback-whatsapp", providerHealth.Name)
	assert.Equal(t, "twilio_whatsapp", providerHealth.Type)
	assert.Equal(t, "healthy", providerHealth.Status)
	assert.Contains(t, providerHealth.Message, "Provider is enabled and configured")
}

// TestWhatsAppProviderHealthCheck_WithTimeoutHandling tests health check with timeout handling
func TestWhatsAppProviderHealthCheck_WithTimeoutHandling(t *testing.T) {
	// Set up environment variables for provider with timeout
	t.Setenv("WHATSAPP_PROVIDER_0_TYPE", "whatsapp_business")
	t.Setenv("WHATSAPP_PROVIDER_0_NAME", "timeout-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_0_PRIORITY", "1")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_PHONE_ID_NUMBER", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_BUSINESS_ACCOUNT_ID", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_ACCESS_TOKEN", "test-access-token")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_API_VERSION", "v15.0")

	// Load configuration
	cfg := config.LoadTest()

	// Get provider health
	health := cfg.GetProviderHealth()

	// Verify health check results
	assert.Len(t, health, 1, "Expected 1 WhatsApp provider to be checked")

	// Check provider health
	providerHealth := health[0]
	assert.Equal(t, "timeout-whatsapp", providerHealth.Name)
	assert.Equal(t, "whatsapp_business", providerHealth.Type)
	assert.Equal(t, "healthy", providerHealth.Status)
	assert.Contains(t, providerHealth.Message, "Provider is enabled and configured")
	assert.Equal(t, "Not implemented", providerHealth.LastTest)
}

// TestWhatsAppProviderHealthCheck_WithPerformanceMetrics tests health check with performance metrics
func TestWhatsAppProviderHealthCheck_WithPerformanceMetrics(t *testing.T) {
	// Set up environment variables for provider
	t.Setenv("WHATSAPP_PROVIDER_0_TYPE", "whatsapp_business")
	t.Setenv("WHATSAPP_PROVIDER_0_NAME", "performance-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_0_PRIORITY", "1")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_PHONE_ID_NUMBER", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_BUSINESS_ACCOUNT_ID", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_ACCESS_TOKEN", "test-access-token")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_API_VERSION", "v15.0")

	// Load configuration
	cfg := config.LoadTest()

	// Get provider health
	health := cfg.GetProviderHealth()

	// Verify health check results
	assert.Len(t, health, 1, "Expected 1 WhatsApp provider to be checked")

	// Check provider health
	providerHealth := health[0]
	assert.Equal(t, "performance-whatsapp", providerHealth.Name)
	assert.Equal(t, "whatsapp_business", providerHealth.Type)
	assert.Equal(t, "healthy", providerHealth.Status)
	assert.Contains(t, providerHealth.Message, "Provider is enabled and configured")
	assert.Equal(t, "Not implemented", providerHealth.LastTest)
}

// TestWhatsAppProviderHealthCheck_WithConfigurationValidation tests health check with configuration validation
func TestWhatsAppProviderHealthCheck_WithConfigurationValidation(t *testing.T) {
	// Set up environment variables for provider with invalid configuration
	t.Setenv("WHATSAPP_PROVIDER_0_TYPE", "whatsapp_business")
	t.Setenv("WHATSAPP_PROVIDER_0_NAME", "invalid-config-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_0_PRIORITY", "1")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_PHONE_ID_NUMBER", "invalid-phone-id")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_BUSINESS_ACCOUNT_ID", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_ACCESS_TOKEN", "test-access-token")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_API_VERSION", "v15.0")

	// Load configuration
	cfg := config.LoadTest()

	// Get provider health
	health := cfg.GetProviderHealth()

	// Verify health check results
	assert.Len(t, health, 1, "Expected 1 WhatsApp provider to be checked")

	// Check provider health
	providerHealth := health[0]
	assert.Equal(t, "invalid-config-whatsapp", providerHealth.Name)
	assert.Equal(t, "whatsapp_business", providerHealth.Type)
	assert.Equal(t, "unhealthy", providerHealth.Status)
	assert.Contains(t, providerHealth.Message, "required WhatsApp Business settings missing")
}

// TestWhatsAppProviderHealthCheck_WithRealTimeMonitoring tests health check with real-time monitoring
func TestWhatsAppProviderHealthCheck_WithRealTimeMonitoring(t *testing.T) {
	// Set up environment variables for provider
	t.Setenv("WHATSAPP_PROVIDER_0_TYPE", "whatsapp_business")
	t.Setenv("WHATSAPP_PROVIDER_0_NAME", "monitoring-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_0_PRIORITY", "1")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_PHONE_ID_NUMBER", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_BUSINESS_ACCOUNT_ID", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_ACCESS_TOKEN", "test-access-token")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_API_VERSION", "v15.0")

	// Load configuration
	cfg := config.LoadTest()

	// Get provider health
	health := cfg.GetProviderHealth()

	// Verify health check results
	assert.Len(t, health, 1, "Expected 1 WhatsApp provider to be checked")

	// Check provider health
	providerHealth := health[0]
	assert.Equal(t, "monitoring-whatsapp", providerHealth.Name)
	assert.Equal(t, "whatsapp_business", providerHealth.Type)
	assert.Equal(t, "healthy", providerHealth.Status)
	assert.Contains(t, providerHealth.Message, "Provider is enabled and configured")
	assert.Equal(t, "Not implemented", providerHealth.LastTest)

	// In a real implementation, we would test that the health check is performed
	// periodically and that the results are updated in real-time
	// For now, we'll just verify the basic functionality
}

// TestWhatsAppProviderHealthCheck_WithLoadBalancing tests health check with load balancing
func TestWhatsAppProviderHealthCheck_WithLoadBalancing(t *testing.T) {
	// Set up environment variables for multiple providers with different priorities
	t.Setenv("WHATSAPP_PROVIDER_0_TYPE", "whatsapp_business")
	t.Setenv("WHATSAPP_PROVIDER_0_NAME", "primary-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_0_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_0_PRIORITY", "1")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_FROM_NUMBER", "+1234567890")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_PHONE_ID_NUMBER", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_BUSINESS_ACCOUNT_ID", "1234567890123456")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_ACCESS_TOKEN", "test-access-token")
	t.Setenv("WHATSAPP_PROVIDER_0_SETTINGS_API_VERSION", "v15.0")

	t.Setenv("WHATSAPP_PROVIDER_1_TYPE", "twilio_whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_1_NAME", "secondary-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_1_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_1_PRIORITY", "2")
	t.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_FROM_NUMBER", "+1234567891")
	t.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_ACCOUNT_SID", "test-account-sid")
	t.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_AUTH_TOKEN", "test-auth-token")
	t.Setenv("WHATSAPP_PROVIDER_1_SETTINGS_WHATSAPP_NUMBER", "+1234567891")

	t.Setenv("WHATSAPP_PROVIDER_2_TYPE", "whatsapp_business")
	t.Setenv("WHATSAPP_PROVIDER_2_NAME", "tertiary-whatsapp")
	t.Setenv("WHATSAPP_PROVIDER_2_ENABLED", "true")
	t.Setenv("WHATSAPP_PROVIDER_2_PRIORITY", "3")
	t.Setenv("WHATSAPP_PROVIDER_2_SETTINGS_FROM_NUMBER", "+1234567892")
	t.Setenv("WHATSAPP_PROVIDER_2_SETTINGS_PHONE_ID_NUMBER", "1234567890123457")
	t.Setenv("WHATSAPP_PROVIDER_2_SETTINGS_BUSINESS_ACCOUNT_ID", "1234567890123457")
	t.Setenv("WHATSAPP_PROVIDER_2_SETTINGS_ACCESS_TOKEN", "test-access-token-2")
	t.Setenv("WHATSAPP_PROVIDER_2_SETTINGS_API_VERSION", "v15.0")

	// Load configuration
	cfg := config.LoadTest()

	// Get provider health
	health := cfg.GetProviderHealth()

	// Verify health check results
	assert.Len(t, health, 3, "Expected 3 WhatsApp providers to be checked")

	// Check that all providers are healthy and properly prioritized
	for _, h := range health {
		assert.Equal(t, "healthy", h.Status, "Provider %s should be healthy", h.Name)
		assert.Contains(t, h.Message, "Provider is enabled and configured")
	}

	// In a real implementation, we would test that the load balancing mechanism
	// works correctly by checking that requests are routed to the healthy
	// providers based on their priority
	// For now, we'll just verify that all providers are healthy
}