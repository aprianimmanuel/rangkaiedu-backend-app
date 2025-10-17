package monitoring

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SecurityConfig represents the top-level security configuration
type SecurityConfig struct {
	Security *SecuritySettings `json:"security,omitempty" yaml:"security,omitempty"`
}

// SecuritySettings represents security-related settings
type SecuritySettings struct {
	Enabled     bool                   `json:"enabled" yaml:"enabled"`
	Environment string                 `json:"environment" yaml:"environment"`
	Debug       bool                   `json:"debug" yaml:"debug"`
	Logging     *SecurityLoggingConfig `json:"logging,omitempty" yaml:"logging,omitempty"`
	Alerting    *SecurityAlertingConfig `json:"alerting,omitempty" yaml:"alerting,omitempty"`
}

// SecurityAlertingConfig represents security alerting configuration
type SecurityAlertingConfig struct {
	Enabled  bool `json:"enabled" yaml:"enabled"`
	Level    string `json:"level" yaml:"level"`
}

// SecurityConfigManager handles security configuration loading and management
type SecurityConfigManager struct {
	configPath string
	config     *SecuritySettings
}

// NewSecurityConfigManager creates a new security configuration manager
func NewSecurityConfigManager(configPath string) *SecurityConfigManager {
	return &SecurityConfigManager{
		configPath: configPath,
		config: &SecuritySettings{
			Enabled: true,
			Logging: &SecurityLoggingConfig{
				Enabled: true,
				Level:   "medium",
				Format:  "json",
			},
		},
	}
}

// LoadConfig loads security configuration from file
func (cm *SecurityConfigManager) LoadConfig() error {
	// Check if config file exists
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// Create default config file
		if err := cm.SaveConfig(); err != nil {
			return fmt.Errorf("failed to create default security config file: %v", err)
		}
		return nil
	}

	// Read config file
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return fmt.Errorf("failed to read security config file: %v", err)
	}

	// Parse config
	var config SecurityConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse security config file: %v", err)
	}

	// Validate config
	if config.Security != nil && config.Security.Logging != nil {
		if err := config.Security.Logging.Validate(); err != nil {
			return fmt.Errorf("invalid logging configuration: %v", err)
		}
	}

	cm.config = config.Security
	return nil
}

// SaveConfig saves security configuration to file
func (cm *SecurityConfigManager) SaveConfig() error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(cm.configPath), 0755); err != nil {
		return fmt.Errorf("failed to create security config directory: %v", err)
	}

	// Wrap config in SecurityConfig struct
	config := &SecurityConfig{
		Security: cm.config,
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal security config: %v", err)
	}

	// Write config file
	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write security config file: %v", err)
	}

	return nil
}

// GetConfig returns the current security configuration
func (cm *SecurityConfigManager) GetConfig() *SecuritySettings {
	return cm.config
}

// SetConfig sets the security configuration
func (cm *SecurityConfigManager) SetConfig(config *SecuritySettings) error {
	if err := cm.validateConfig(config); err != nil {
		return fmt.Errorf("invalid security configuration: %v", err)
	}

	cm.config = config
	return nil
}

// validateConfig validates the security configuration
func (cm *SecurityConfigManager) validateConfig(config *SecuritySettings) error {
	if config == nil {
		return fmt.Errorf("security configuration is nil")
	}

	if config.Logging != nil {
		if err := config.Logging.Validate(); err != nil {
			return fmt.Errorf("invalid logging configuration: %v", err)
		}
	}

	return nil
}

// IsEnabled returns whether security logging is enabled
func (cm *SecurityConfigManager) IsEnabled() bool {
	return cm.config != nil && cm.config.Enabled
}

// GetLoggingConfig returns the logging configuration
func (cm *SecurityConfigManager) GetLoggingConfig() *SecurityLoggingConfig {
	if cm.config == nil {
		return nil
	}
	return cm.config.Logging
}

// GetEnvironment returns the environment
func (cm *SecurityConfigManager) GetEnvironment() string {
	if cm.config == nil {
		return "unknown"
	}
	return cm.config.Environment
}

// IsDebug returns whether debug mode is enabled
func (cm *SecurityConfigManager) IsDebug() bool {
	if cm.config == nil {
		return false
	}
	return cm.config.Debug
}