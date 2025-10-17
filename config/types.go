package config

import (
	"time"
)

// AlertSeverity represents the severity level of an alert
type AlertSeverity string

const (
	AlertSeverityCritical AlertSeverity = "critical"
	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityInfo     AlertSeverity = "info"
	AlertSeverityDebug    AlertSeverity = "debug"
)

// HealthStatus represents the health status of a service
type HealthStatus string

const (
	HealthStatusHealthy    HealthStatus = "healthy"
	HealthStatusDegraded   HealthStatus = "degraded"
	HealthStatusUnhealthy  HealthStatus = "unhealthy"
	HealthStatusUnknown    HealthStatus = "unknown"
)

// SecurityLoggingConfig represents the configuration for security logging
type SecurityLoggingConfig struct {
	Enabled                   bool          `json:"enabled" mapstructure:"enabled"`
	MaxHistorySize            int           `json:"max_history_size" mapstructure:"max_history_size"`
	CleanupInterval           time.Duration `json:"cleanup_interval" mapstructure:"cleanup_interval"`
	EnableAuthLogging         bool          `json:"enable_auth_logging" mapstructure:"enable_auth_logging"`
	EnableSecurityMetrics     bool          `json:"enable_security_metrics" mapstructure:"enable_security_metrics"`
	EnableSecurityAlerts      bool          `json:"enable_security_alerts" mapstructure:"enable_security_alerts"`
	EnableSecurityLogging     bool          `json:"enable_security_logging" mapstructure:"enable_security_logging"`
	MaxSecurityEventsPerSecond int          `json:"max_security_events_per_second" mapstructure:"max_security_events_per_second"`
	SecurityEventRetention    time.Duration `json:"security_event_retention" mapstructure:"security_event_retention"`
	EnableIPWhitelist         bool          `json:"enable_ip_whitelist" mapstructure:"enable_ip_whitelist"`
	EnableRateLimiting        bool          `json:"enable_rate_limiting" mapstructure:"enable_rate_limiting"`
	RateLimit                 int           `json:"rate_limit" mapstructure:"rate_limit"`
	RateLimitWindow           time.Duration `json:"rate_limit_window" mapstructure:"rate_limit_window"`
}

// ProviderHealth represents the health status of a provider
type ProviderHealth struct {
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	Status  HealthStatus `json:"status"`
	Message string      `json:"message"`
}

// DefaultSecurityLoggingConfig returns the default security logging configuration
func DefaultSecurityLoggingConfig() *SecurityLoggingConfig {
	return &SecurityLoggingConfig{
		Enabled:                   true,
		MaxHistorySize:            1000,
		CleanupInterval:           1 * time.Hour,
		EnableAuthLogging:         true,
		EnableSecurityMetrics:     true,
		EnableSecurityAlerts:      true,
		EnableSecurityLogging:     true,
		MaxSecurityEventsPerSecond: 100,
		SecurityEventRetention:    24 * time.Hour,
		EnableIPWhitelist:         false,
		EnableRateLimiting:        true,
		RateLimit:                 100,
		RateLimitWindow:           time.Minute,
	}
}

// ProviderManager manages different types of providers
type ProviderManager struct {
	EmailProviders []EmailProvider
	SMSProviders   []SMSProvider
}

// EmailProvider represents an email provider configuration
type EmailProvider struct {
	Type      string                 `json:"type"`
	Name      string                 `json:"name"`
	Enabled   bool                   `json:"enabled"`
	Host      string                 `json:"host"`
	Port      int                    `json:"port"`
	Username  string                 `json:"username"`
	Password  string                 `json:"password"`
	From      string                 `json:"from"`
	Priority  int                    `json:"priority"`
	Settings  map[string]interface{} `json:"settings"`
}

// SMSProvider represents an SMS provider configuration
type SMSProvider struct {
	Type      string                 `json:"type"`
	Name      string                 `json:"name"`
	Enabled   bool                   `json:"enabled"`
	APIKey    string                 `json:"api_key"`
	Endpoint  string                 `json:"endpoint"`
	Priority  int                    `json:"priority"`
	Settings  map[string]interface{} `json:"settings"`
}

// JWTConfig represents JWT configuration
type JWTConfig struct {
	Secret        string        `json:"secret"`
	AccessTokenExpiry time.Duration `json:"access_token_expiry"`
	RefreshTokenExpiry time.Duration `json:"refresh_token_expiry"`
}

// StorageProvider represents storage provider configuration
type StorageProvider struct {
	Type            string                 `json:"type"`
	Provider        string                 `json:"provider"`
	Endpoint        string                 `json:"endpoint"`
	AccessKeyID     string                 `json:"access_key_id"`
	AccessKeySecret string                 `json:"access_key_secret"`
	BucketName      string                 `json:"bucket_name"`
	Region          string                 `json:"region"`
	Settings        map[string]interface{} `json:"settings"`
}

// GCSConfig represents Google Cloud Storage configuration
type GCSConfig struct {
	ServiceAccountKeyPath string `json:"service_account_key_path"`
	BucketName            string `json:"bucket_name"`
	ProjectID             string `json:"project_id"`
}

// OSSConfig represents Alibaba Cloud OSS configuration
type OSSConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	BucketName      string `json:"bucket_name"`
}

// EmailProviderConfig represents email provider configuration
type EmailProviderConfig struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	From      string `json:"from"`
	Enabled   bool   `json:"enabled"`
	Priority  int    `json:"priority"`
}

// SMSProviderConfig represents SMS provider configuration
type SMSProviderConfig struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	APIKey    string `json:"api_key"`
	Endpoint  string `json:"endpoint"`
	Enabled   bool   `json:"enabled"`
	Priority  int    `json:"priority"`
}

// Config represents the main configuration structure
type Config struct {
	ProviderManager     ProviderManager     `json:"provider_manager"`
	JWT                 JWTConfig           `json:"jwt"`
	StorageProvider     StorageProvider     `json:"storage_provider"`
	GCSServiceAccountKeyPath string          `json:"gcs_service_account_key_path"`
	GCSBucketName       string              `json:"gcs_bucket_name"`
	GCSProjectID        string              `json:"gcs_project_id"`
	OSSEndpoint         string              `json:"oss_endpoint"`
	OSSAccessKeyID      string              `json:"oss_access_key_id"`
	OSSAccessKeySecret  string              `json:"oss_access_key_secret"`
	OSSBucketName       string              `json:"oss_bucket_name"`
}

// GetProviderHealth returns the health status of all configured providers
func (c *Config) GetProviderHealth() []ProviderHealth {
	var health []ProviderHealth
	
	// Add email providers health
	for _, provider := range c.ProviderManager.EmailProviders {
		status := HealthStatusHealthy
		message := "Provider is healthy"
		
		if !provider.Enabled {
			status = HealthStatusDegraded
			message = "Provider is disabled"
		}
		
		health = append(health, ProviderHealth{
			Name:    provider.Name,
			Type:    "email",
			Status:  status,
			Message: message,
		})
	}
	
	// Add SMS providers health
	for _, provider := range c.ProviderManager.SMSProviders {
		status := HealthStatusHealthy
		message := "Provider is healthy"
		
		if !provider.Enabled {
			status = HealthStatusDegraded
			message = "Provider is disabled"
		}
		
		health = append(health, ProviderHealth{
			Name:    provider.Name,
			Type:    "sms",
			Status:  status,
			Message: message,
		})
	}
	
	return health
}