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
	SMTPHost       string
	SMTPPort       string
	SMTPUser       string
	SMTPPass       string
	TWILIOAccountSID   string
	TWILIOAuthToken    string
	TWILIOSenderPhone  string
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
}

// Load loads the configuration from environment variables.
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBName:         getEnv("DB_NAME", "rangkaiedu_dev"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "password"), // Change this in production
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		JWTSecret:      getEnv("JWT_SECRET", "default-secret-key-change-in-production"),
		SMTPHost:       getEnv("SMTP_HOST", ""),
		SMTPPort:       getEnv("SMTP_PORT", "587"),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPass:       getEnv("SMTP_PASS", ""),
		TWILIOAccountSID:   getEnv("TWILIO_ACCOUNT_SID", ""),
		TWILIOAuthToken:    getEnv("TWILIO_AUTH_TOKEN", ""),
		TWILIOSenderPhone:  getEnv("TWILIO_SENDER_PHONE", ""),
		// Storage configuration
		StorageProvider:      getEnv("STORAGE_PROVIDER", "local"),
		// OSS configuration
		OSSBucketName:        getEnv("OSS_BUCKET_NAME", ""),
		OSSAccessKeyID:       getEnv("OSS_ACCESS_KEY_ID", ""),
		OSSAccessKeySecret:   getEnv("OSS_ACCESS_KEY_SECRET", ""),
		OSSRegion:            getEnv("OSS_REGION", ""),
		OSSEndpoint:          getEnv("OSS_ENDPOINT", ""),
		// GCS configuration (deprecated)
		GCSBucketName:            getEnv("GCS_BUCKET_NAME", ""),
		GCSServiceAccountKeyPath: getEnv("GCS_SERVICE_ACCOUNT_KEY_PATH", ""),
		GCSProjectID:             getEnv("GCS_PROJECT_ID", ""),
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
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "7110"),
		DBName:         getEnv("DB_NAME", "rangkaiedu_test"),
		DBUser:         getEnv("DB_USER", "rangkaiedudev1"),
		DBPassword:     getEnv("DB_PASSWORD", "12d1q23wxm19wkc1fsdcq23"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		JWTSecret:      getEnv("JWT_SECRET", "test-jwt-secret-change-in-production"),
		SMTPHost:       getEnv("SMTP_HOST", ""),
		SMTPPort:       getEnv("SMTP_PORT", "587"),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPass:       getEnv("SMTP_PASS", ""),
		TWILIOAccountSID:   getEnv("TWILIO_ACCOUNT_SID", ""),
		TWILIOAuthToken:    getEnv("TWILIO_AUTH_TOKEN", ""),
		TWILIOSenderPhone:  getEnv("TWILIO_SENDER_PHONE", ""),
		// Storage configuration
		StorageProvider:      getEnv("STORAGE_PROVIDER", "local"),
		// OSS configuration
		OSSBucketName:        getEnv("OSS_BUCKET_NAME", ""),
		OSSAccessKeyID:       getEnv("OSS_ACCESS_KEY_ID", ""),
		OSSAccessKeySecret:   getEnv("OSS_ACCESS_KEY_SECRET", ""),
		OSSRegion:            getEnv("OSS_REGION", ""),
		OSSEndpoint:          getEnv("OSS_ENDPOINT", ""),
		// GCS configuration (deprecated)
		GCSBucketName:            getEnv("GCS_BUCKET_NAME", ""),
		GCSServiceAccountKeyPath: getEnv("GCS_SERVICE_ACCOUNT_KEY_PATH", ""),
		GCSProjectID:             getEnv("GCS_PROJECT_ID", ""),
	}

	// Validate required fields
	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" || cfg.DBUser == "" {
		log.Fatal("Missing required database configuration. Please set DB_HOST, DB_PORT, DB_NAME, DB_USER, and DB_PASSWORD in .env.test")
	}

	return cfg
}

// getEnv returns the value of the environment variable named by the key.
// If the variable is not found, it returns the provided default value.
func getEnv(key, defaultValue string) string {
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