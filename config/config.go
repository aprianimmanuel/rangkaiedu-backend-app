package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// Load loads configuration from environment variables or config file
func Load() *Config {
	cfg := &Config{
		ProviderManager: ProviderManager{
			EmailProviders: []EmailProvider{},
			SMSProviders:   []SMSProvider{},
		},
		JWT: JWTConfig{
			Secret:               "default-secret-key-change-in-production",
			AccessTokenExpiry:    24 * time.Hour,
			RefreshTokenExpiry:   168 * time.Hour, // 7 days
		},
		StorageProvider: StorageProvider{
			Type: "local",
		},
	}

	// Try to load from config file if it exists
	if _, err := os.Stat("config.json"); err == nil {
		if err := loadFromFile(cfg, "config.json"); err != nil {
			log.Printf("Warning: Failed to load config from file: %v", err)
		}
	}

	return cfg
}

// loadFromFile loads configuration from a JSON file
func loadFromFile(cfg *Config, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// Save saves configuration to a JSON file
func (c *Config) Save(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}