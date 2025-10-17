package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration.
type Config struct {
	DBHost         string
	DBPort         string
	DBName         string
	DBUser         string
	DBPassword     string
	DBSSLMode      string // "disable" for local dev, "require" for production
	JWTSecret      string
	// OAuth configuration
	GoogleClientID     string
	GoogleClientSecret string
	FacebookClientID   string
	FacebookClientSecret string
	// Storage configuration
	StorageProvider      string
	// OSS configuration
	OSSBucketName        string
	OSSAccessKeyID       string
	OSSAccessKeySecret   string
	OSSRegion            string
	OSSEndpoint          string
	// GCS configuration (deprecated)
	GCSBucketName             string
	GCSServiceAccountKeyPath  string
	GCSProjectID              string
	
	// Provider management
	ProviderManager *ProviderManager
	CredentialManager *CredentialManager
	
	// HTTPS configuration
	HTTPS HTTPSConfig
	
	// SMTP configuration for email sending
	SMTPHost string
	SMTPUser string
	SMTPPassword string
	SMTPPort string
}

// Load loads the configuration from environment variables.
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := &Config{
		DBHost:         getConfigEnv("DB_HOST", "localhost"),
		DBPort:         getConfigEnv("DB_PORT", "5432"),
		DBName:         getConfigEnv("DB_NAME", "rangkaiedu_dev"),
		DBUser:         getConfigEnv("DB_USER", "postgres"),
		DBPassword:     getConfigEnv("DB_PASSWORD", "password"), // Change this in production
		DBSSLMode:      getConfigEnv("DB_SSLMODE", "disable"),
		JWTSecret:      getConfigEnv("JWT_SECRET", "default-secret-key-change-in-production"),
		// OAuth configuration
		GoogleClientID:     getConfigEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getConfigEnv("GOOGLE_CLIENT_SECRET", ""),
		FacebookClientID:   getConfigEnv("FACEBOOK_CLIENT_ID", ""),
		FacebookClientSecret: getConfigEnv("FACEBOOK_CLIENT_SECRET", ""),
		// Storage configuration
		StorageProvider:      getConfigEnv("STORAGE_PROVIDER", "local"),
		// OSS configuration
		OSSBucketName:        getConfigEnv("OSS_BUCKET_NAME", ""),
		OSSAccessKeyID:       getConfigEnv("OSS_ACCESS_KEY_ID", ""),
		OSSAccessKeySecret:   getConfigEnv("OSS_ACCESS_KEY_SECRET", ""),
		OSSRegion:            getConfigEnv("OSS_REGION", ""),
		OSSEndpoint:          getConfigEnv("OSS_ENDPOINT", ""),
		// GCS configuration (deprecated)
		GCSBucketName:            getConfigEnv("GCS_BUCKET_NAME", ""),
		GCSServiceAccountKeyPath: getConfigEnv("GCS_SERVICE_ACCOUNT_KEY_PATH", ""),
		GCSProjectID:             getConfigEnv("GCS_PROJECT_ID", ""),
		
		// Initialize provider management
		ProviderManager:     NewProviderManager(),
		CredentialManager:   NewCredentialManager(),
		
		// Initialize HTTPS configuration
		HTTPS: LoadHTTPSConfig(),
		
		// Initialize SMTP configuration
		SMTPHost:     getConfigEnv("SMTP_HOST", ""),
		SMTPUser:     getConfigEnv("SMTP_USER", ""),
		SMTPPassword: getConfigEnv("SMTP_PASSWORD", ""),
		SMTPPort:     getConfigEnv("SMTP_PORT", "587"),
	}

	// Load provider configurations
	if err := cfg.ProviderManager.LoadProviders(); err != nil {
		log.Printf("Warning: Failed to load provider configurations: %v", err)
	}

	// Validate required fields
	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" || cfg.DBUser == "" {
		log.Fatal("Missing required database configuration. Please set DB_HOST, DB_PORT, DB_NAME, DB_USER, and DB_PASSWORD in .env")
	}

	return cfg
}

// LoadTest loads the configuration from environment variables for testing.
func LoadTest() *Config {
	// Load .env.test file if it exists
	if err := godotenv.Load(".env.test"); err != nil {
		log.Println("No .env.test file found")
	}

	cfg := &Config{
		DBHost:         getConfigEnv("DB_HOST", "localhost"),
		DBPort:         getConfigEnv("DB_PORT", "5432"),
		DBName:         getConfigEnv("DB_NAME", "rangkaiedu_test"),
		DBUser:         getConfigEnv("DB_USER", "rangkaiedudev1"),
		DBPassword:     getConfigEnv("DB_PASSWORD", "12d1q23wxm19wkc1fsdcq23"),
		DBSSLMode:      getConfigEnv("DB_SSLMODE", "disable"),
		JWTSecret:      getConfigEnv("JWT_SECRET", "test-jwt-secret-change-in-production"),
		// OAuth configuration
		GoogleClientID:     getConfigEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getConfigEnv("GOOGLE_CLIENT_SECRET", ""),
		FacebookClientID:   getConfigEnv("FACEBOOK_CLIENT_ID", ""),
		FacebookClientSecret: getConfigEnv("FACEBOOK_CLIENT_SECRET", ""),
		// Storage configuration
		StorageProvider:      getConfigEnv("STORAGE_PROVIDER", "local"),
		// OSS configuration
		OSSBucketName:        getConfigEnv("OSS_BUCKET_NAME", ""),
		OSSAccessKeyID:       getConfigEnv("OSS_ACCESS_KEY_ID", ""),
		OSSAccessKeySecret:   getConfigEnv("OSS_ACCESS_KEY_SECRET", ""),
		OSSRegion:            getConfigEnv("OSS_REGION", ""),
		OSSEndpoint:          getConfigEnv("OSS_ENDPOINT", ""),
		// GCS configuration (deprecated)
		GCSBucketName:            getConfigEnv("GCS_BUCKET_NAME", ""),
		GCSServiceAccountKeyPath: getConfigEnv("GCS_SERVICE_ACCOUNT_KEY_PATH", ""),
		GCSProjectID:             getConfigEnv("GCS_PROJECT_ID", ""),
		
		// Initialize provider management
		ProviderManager:     NewProviderManager(),
		CredentialManager:   NewCredentialManager(),
		
		// Initialize HTTPS configuration
		HTTPS: LoadHTTPSConfig(),
		
		// Initialize SMTP configuration
		SMTPHost:     getConfigEnv("SMTP_HOST", ""),
		SMTPUser:     getConfigEnv("SMTP_USER", ""),
		SMTPPassword: getConfigEnv("SMTP_PASSWORD", ""),
		SMTPPort:     getConfigEnv("SMTP_PORT", "587"),
	}

	// Load provider configurations
	if err := cfg.ProviderManager.LoadProviders(); err != nil {
		log.Printf("Warning: Failed to load provider configurations: %v", err)
	}

	// Validate required fields
	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" || cfg.DBUser == "" {
		log.Fatal("Missing required database configuration. Please set DB_HOST, DB_PORT, DB_NAME, DB_USER, and DB_PASSWORD in .env.test")
	}

	return cfg
}

// getEnv returns the value of the environment variable named by the key.
// If the variable is not found, it returns the provided default value.
func getConfigEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// DSN returns the PostgreSQL DSN string.
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBName, c.DBUser, c.DBPassword, c.DBSSLMode)
}

// GetPrimaryEmailProvider returns the primary email provider configuration
func (c *Config) GetPrimaryEmailProvider() (EmailProviderConfig, error) {
	return c.ProviderManager.GetPrimaryEmailProvider()
}

// GetPrimarySMSProvider returns the primary SMS provider configuration
func (c *Config) GetPrimarySMSProvider() (SMSProviderConfig, error) {
	return c.ProviderManager.GetPrimarySMSProvider()
}

// GetEmailProviderByName returns a specific email provider by name
func (c *Config) GetEmailProviderByName(name string) (EmailProviderConfig, error) {
	return c.ProviderManager.GetEmailProviderByName(name)
}

// GetSMSProviderByName returns a specific SMS provider by name
func (c *Config) GetSMSProviderByName(name string) (SMSProviderConfig, error) {
	return c.ProviderManager.GetSMSProviderByName(name)
}

// GetPrimaryWhatsAppProvider returns the primary WhatsApp provider configuration
func (c *Config) GetPrimaryWhatsAppProvider() (WhatsAppProviderConfig, error) {
	return c.ProviderManager.GetPrimaryWhatsAppProvider()
}

// GetWhatsAppProviderByName returns a specific WhatsApp provider by name
func (c *Config) GetWhatsAppProviderByName(name string) (WhatsAppProviderConfig, error) {
	return c.ProviderManager.GetWhatsAppProviderByName(name)
}

// ValidateProviders validates all provider configurations
func (c *Config) ValidateProviders() error {
	return c.ProviderManager.Validate()
}

// GetProviderHealth returns health status of all providers
func (c *Config) GetProviderHealth() []ProviderHealth {
	return c.ProviderManager.HealthCheck()
}

// EncryptCredential encrypts a credential value
func (c *Config) EncryptCredential(value string) (string, error) {
	return c.CredentialManager.Encrypt(value)
}

// DecryptCredential decrypts a credential value
func (c *Config) DecryptCredential(value string) (string, error) {
	return c.CredentialManager.Decrypt(value)
}

// IsEncrypted checks if a value is encrypted
func (c *Config) IsEncrypted(value string) bool {
	return c.CredentialManager.IsEncrypted(value)
}