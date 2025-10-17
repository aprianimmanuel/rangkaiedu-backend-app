package monitoring

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)


// ConfigManager handles configuration loading and management
type ConfigManager struct {
	configPath string
	config     *MonitoringConfig
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(configPath string) *ConfigManager {
	return &ConfigManager{
		configPath: configPath,
		config:     DefaultMonitoringConfig(),
	}
}

// LoadConfig loads configuration from file
func (cm *ConfigManager) LoadConfig() error {
	// Check if config file exists
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// Create default config file
		if err := cm.SaveConfig(); err != nil {
			return fmt.Errorf("failed to create default config file: %v", err)
		}
		return nil
	}

	// Read config file
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse config
	var config MonitoringConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	// Validate config
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	cm.config = &config
	return nil
}

// SaveConfig saves configuration to file
func (cm *ConfigManager) SaveConfig() error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(cm.configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write config file
	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// GetConfig returns the current configuration
func (cm *ConfigManager) GetConfig() *MonitoringConfig {
	return cm.config
}

// SetConfig sets the configuration
func (cm *ConfigManager) SetConfig(config *MonitoringConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	cm.config = config
	return nil
}

// ReloadConfig reloads configuration from file
func (cm *ConfigManager) ReloadConfig() error {
	return cm.LoadConfig()
}

// WatchConfig watches for configuration changes and reloads automatically
func (cm *ConfigManager) WatchConfig(callback func(*MonitoringConfig)) error {
	// This would use file system watching to detect changes
	// For now, we'll just provide a placeholder implementation
	return nil
}

// GetConfigPath returns the configuration file path
func (cm *ConfigManager) GetConfigPath() string {
	return cm.configPath
}

// SetConfigPath sets the configuration file path
func (cm *ConfigManager) SetConfigPath(path string) {
	cm.configPath = path
}

// ValidateConfig validates the current configuration
func (cm *ConfigManager) ValidateConfig() error {
	return cm.config.Validate()
}

// IsConfigValid returns whether the current configuration is valid
func (cm *ConfigManager) IsConfigValid() bool {
	return cm.config.Validate() == nil
}

// GetConfigValidationErrors returns validation errors for the current configuration
func (cm *ConfigManager) GetConfigValidationErrors() []string {
	var errors []string
	
	if cm.config == nil {
		errors = append(errors, "configuration is nil")
		return errors
	}
	
	if cm.config.Interval <= 0 {
		errors = append(errors, "interval must be positive")
	}
	
	if cm.config.Timeout <= 0 {
		errors = append(errors, "timeout must be positive")
	}
	
	if cm.config.MaxHistorySize <= 0 {
		errors = append(errors, "max history size must be positive")
	}
	
	if cm.config.CleanupInterval <= 0 {
		errors = append(errors, "cleanup interval must be positive")
	}
	
	// Validate sub-configurations
	if cm.config.MetricsConfig != nil {
		if err := cm.config.MetricsConfig.Validate(); err != nil {
			errors = append(errors, fmt.Sprintf("metrics config: %v", err))
		}
	}
	
	if cm.config.ErrorConfig != nil {
		if err := cm.config.ErrorConfig.Validate(); err != nil {
			errors = append(errors, fmt.Sprintf("error config: %v", err))
		}
	}
	
	if cm.config.AlertConfig != nil {
		if err := cm.config.AlertConfig.Validate(); err != nil {
			errors = append(errors, fmt.Sprintf("alert config: %v", err))
		}
	}
	
	if cm.config.HealthConfig != nil {
		if err := cm.config.HealthConfig.Validate(); err != nil {
			errors = append(errors, fmt.Sprintf("health config: %v", err))
		}
	}
	
	return errors
}

// MergeConfig merges the current configuration with another configuration
func (cm *ConfigManager) MergeConfig(other *MonitoringConfig) error {
	if other == nil {
		return fmt.Errorf("cannot merge with nil configuration")
	}

	merged := cm.config.Merge(other)
	if err := merged.Validate(); err != nil {
		return fmt.Errorf("merged configuration is invalid: %v", err)
	}

	cm.config = merged
	return nil
}

// CloneConfig creates a copy of the current configuration
func (cm *ConfigManager) CloneConfig() *MonitoringConfig {
	return cm.config.Clone()
}

// ResetConfig resets the configuration to default values
func (cm *ConfigManager) ResetConfig() error {
	cm.config = DefaultMonitoringConfig()
	return nil
}

// EnableComponent enables a specific monitoring component
func (cm *ConfigManager) EnableComponent(component string) error {
	switch component {
	case "metrics":
		cm.config.EnableMetrics = true
	case "errors":
		cm.config.EnableErrors = true
	case "alerts":
		cm.config.EnableAlerts = true
	case "health_checks":
		cm.config.EnableHealthChecks = true
	default:
		return fmt.Errorf("unknown component: %s", component)
	}
	return nil
}

// DisableComponent disables a specific monitoring component
func (cm *ConfigManager) DisableComponent(component string) error {
	switch component {
	case "metrics":
		cm.config.EnableMetrics = false
	case "errors":
		cm.config.EnableErrors = false
	case "alerts":
		cm.config.EnableAlerts = false
	case "health_checks":
		cm.config.EnableHealthChecks = false
	default:
		return fmt.Errorf("unknown component: %s", component)
	}
	return nil
}

// IsComponentEnabled returns whether a specific component is enabled
func (cm *ConfigManager) IsComponentEnabled(component string) bool {
	switch component {
	case "metrics":
		return cm.config.EnableMetrics
	case "errors":
		return cm.config.EnableErrors
	case "alerts":
		return cm.config.EnableAlerts
	case "health_checks":
		return cm.config.EnableHealthChecks
	default:
		return false
	}
}

// GetEnabledComponents returns the list of enabled components
func (cm *ConfigManager) GetEnabledComponents() []string {
	var components []string

	if cm.config.EnableMetrics {
		components = append(components, "metrics")
	}
	if cm.config.EnableErrors {
		components = append(components, "errors")
	}
	if cm.config.EnableAlerts {
		components = append(components, "alerts")
	}
	if cm.config.EnableHealthChecks {
		components = append(components, "health_checks")
	}

	return components
}

// GetDisabledComponents returns the list of disabled components
func (cm *ConfigManager) GetDisabledComponents() []string {
	var components []string

	if !cm.config.EnableMetrics {
		components = append(components, "metrics")
	}
	if !cm.config.EnableErrors {
		components = append(components, "errors")
	}
	if !cm.config.EnableAlerts {
		components = append(components, "alerts")
	}
	if !cm.config.EnableHealthChecks {
		components = append(components, "health_checks")
	}

	return components
}

// EnableAllComponents enables all monitoring components
func (cm *ConfigManager) EnableAllComponents() {
	cm.config.EnableMetrics = true
	cm.config.EnableErrors = true
	cm.config.EnableAlerts = true
	cm.config.EnableHealthChecks = true
}

// DisableAllComponents disables all monitoring components
func (cm *ConfigManager) DisableAllComponents() {
	cm.config.EnableMetrics = false
	cm.config.EnableErrors = false
	cm.config.EnableAlerts = false
	cm.config.EnableHealthChecks = false
}

// GetConfigSummary returns a summary of the current configuration
func (cm *ConfigManager) GetConfigSummary() map[string]interface{} {
	return map[string]interface{}{
		"enabled":            cm.config.Enabled,
		"interval":           cm.config.Interval.String(),
		"timeout":            cm.config.Timeout.String(),
		"max_history_size":   cm.config.MaxHistorySize,
		"cleanup_interval":   cm.config.CleanupInterval.String(),
		"enable_metrics":     cm.config.EnableMetrics,
		"enable_errors":      cm.config.EnableErrors,
		"enable_alerts":      cm.config.EnableAlerts,
		"enable_health_checks": cm.config.EnableHealthChecks,
		"components_enabled": cm.GetEnabledComponents(),
		"components_disabled": cm.GetDisabledComponents(),
		"config_path":        cm.configPath,
		"is_valid":           cm.IsConfigValid(),
	}
}

// ExportConfig exports the configuration to a format
func (cm *ConfigManager) ExportConfig(format string) ([]byte, error) {
	switch format {
	case "json":
		return json.MarshalIndent(cm.config, "", "  ")
	case "yaml":
		// This would require a YAML library
		return nil, fmt.Errorf("YAML export not implemented")
	case "toml":
		// This would require a TOML library
		return nil, fmt.Errorf("TOML export not implemented")
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// ImportConfig imports configuration from data
func (cm *ConfigManager) ImportConfig(data []byte, format string) error {
	switch format {
	case "json":
		var config MonitoringConfig
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse JSON: %v", err)
		}
		return cm.SetConfig(&config)
	case "yaml":
		// This would require a YAML library
		return fmt.Errorf("YAML import not implemented")
	case "toml":
		// This would require a TOML library
		return fmt.Errorf("TOML import not implemented")
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// GetEnvironmentConfig loads configuration from environment variables
func (cm *ConfigManager) GetEnvironmentConfig() *MonitoringConfig {
	config := DefaultMonitoringConfig()

	// Load from environment variables
	if enabled := os.Getenv("MONITORING_ENABLED"); enabled != "" {
		config.Enabled = enabled == "true"
	}

	if interval := os.Getenv("MONITORING_INTERVAL"); interval != "" {
		if d, err := time.ParseDuration(interval); err == nil {
			config.Interval = d
		}
	}

	if timeout := os.Getenv("MONITORING_TIMEOUT"); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			config.Timeout = d
		}
	}

	if maxHistory := os.Getenv("MONITORING_MAX_HISTORY_SIZE"); maxHistory != "" {
		if i, err := parseInt(maxHistory); err == nil {
			config.MaxHistorySize = i
		}
	}

	if cleanupInterval := os.Getenv("MONITORING_CLEANUP_INTERVAL"); cleanupInterval != "" {
		if d, err := time.ParseDuration(cleanupInterval); err == nil {
			config.CleanupInterval = d
		}
	}

	if metricsEnabled := os.Getenv("MONITORING_ENABLE_METRICS"); metricsEnabled != "" {
		config.EnableMetrics = metricsEnabled == "true"
	}

	if errorsEnabled := os.Getenv("MONITORING_ENABLE_ERRORS"); errorsEnabled != "" {
		config.EnableErrors = errorsEnabled == "true"
	}

	if alertsEnabled := os.Getenv("MONITORING_ENABLE_ALERTS"); alertsEnabled != "" {
		config.EnableAlerts = alertsEnabled == "true"
	}

	if healthChecksEnabled := os.Getenv("MONITORING_ENABLE_HEALTH_CHECKS"); healthChecksEnabled != "" {
		config.EnableHealthChecks = healthChecksEnabled == "true"
	}

	// Load metrics configuration
	if metricsPort := os.Getenv("MONITORING_METRICS_PORT"); metricsPort != "" {
		if i, err := parseInt(metricsPort); err == nil {
			if config.MetricsConfig == nil {
				config.MetricsConfig = DefaultMetricsConfig()
			}
			config.MetricsConfig.PrometheusPort = i
		}
	}

	if metricsPath := os.Getenv("MONITORING_METRICS_PATH"); metricsPath != "" {
		if config.MetricsConfig == nil {
			config.MetricsConfig = DefaultMetricsConfig()
		}
		config.MetricsConfig.PrometheusPath = metricsPath
	}

	// Load error configuration
	if errorRateLimit := os.Getenv("MONITORING_ERROR_RATE_LIMIT"); errorRateLimit != "" {
		if i, err := parseInt(errorRateLimit); err == nil {
			if config.ErrorConfig == nil {
				config.ErrorConfig = DefaultErrorConfig()
			}
			config.ErrorConfig.ErrorRateLimit = i
		}
	}

	// Load alert configuration
	if alertRateLimit := os.Getenv("MONITORING_ALERT_RATE_LIMIT"); alertRateLimit != "" {
		if i, err := parseInt(alertRateLimit); err == nil {
			if config.AlertConfig == nil {
				config.AlertConfig = DefaultAlertConfig()
			}
			config.AlertConfig.AlertRateLimit = i
		}
	}

	// Load health configuration
	if healthRateLimit := os.Getenv("MONITORING_HEALTH_RATE_LIMIT"); healthRateLimit != "" {
		if i, err := parseInt(healthRateLimit); err == nil {
			if config.HealthConfig == nil {
				config.HealthConfig = DefaultHealthConfig()
			}
			config.HealthConfig.HealthRateLimit = i
		}
	}

	return config
}

// parseInt is a helper function to parse string to int
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

// GetConfigTemplate returns a configuration template
func GetConfigTemplate() string {
	return `{
  "enabled": true,
  "interval": "30s",
  "timeout": "10s",
  "max_history_size": 1000,
  "cleanup_interval": "1h",
  "enable_metrics": true,
  "enable_errors": true,
  "enable_alerts": true,
  "enable_health_checks": true,
  "metrics_config": {
    "enabled": true,
    "interval": "15s",
    "timeout": "5s",
    "prometheus_port": 8080,
    "prometheus_path": "/metrics",
    "enable_system_metrics": true,
    "enable_application_metrics": true,
    "enable_custom_metrics": true,
    "enable_business_metrics": true,
    "max_metrics_per_second": 100,
    "metric_retention": "24h",
    "enable_histograms": true,
    "enable_summaries": true,
    "enable_counters": true,
    "enable_gauges": true,
    "enable_process_metrics": true,
    "enable_runtime_metrics": true,
    "enable_memory_metrics": true,
    "enable_cpu_metrics": true,
    "enable_disk_metrics": true,
    "enable_network_metrics": true
  },
  "error_config": {
    "enabled": true,
    "max_history_size": 1000,
    "cleanup_interval": "1h",
    "enable_error_tracking": true,
    "enable_error_recovery": true,
    "enable_error_alerts": true,
    "enable_error_logging": true,
    "max_retries": 3,
    "retry_delay": "5s",
    "enable_error_classification": true,
    "enable_error_aggregation": true,
    "enable_error_rate_limiting": true,
    "error_rate_limit": 100,
    "error_rate_window": "1m",
    "enable_error_context": true,
    "enable_error_stack_traces": true,
    "enable_error_metrics": true
  },
  "alert_config": {
    "enabled": true,
    "max_history_size": 1000,
    "cleanup_interval": "1h",
    "enable_alert_rules": true,
    "enable_alert_providers": true,
    "enable_alert_templates": true,
    "enable_alert_routing": true,
    "enable_alert_aggregation": true,
    "enable_alert_deduplication": true,
    "enable_alert_escalation": true,
    "enable_alert_silencing": true,
    "enable_alert_metrics": true,
    "default_severity": "warning",
    "default_timeout": "1h",
    "auto_resolve_timeout": "2h",
    "enable_alert_groups": true,
    "max_alerts_per_minute": 10,
    "alert_rate_limit": 100,
    "alert_rate_window": "1m"
  },
  "health_config": {
    "enabled": true,
    "interval": "30s",
    "timeout": "10s",
    "max_history_size": 1000,
    "cleanup_interval": "1h",
    "enable_system_health": true,
    "enable_application_health": true,
    "enable_database_health": true,
    "enable_network_health": true,
    "enable_dependency_health": true,
    "enable_custom_health": true,
    "enable_health_metrics": true,
    "enable_health_logging": true,
    "enable_health_alerts": true,
    "default_threshold": "healthy",
    "enable_health_groups": true,
    "max_health_checks_per_second": 10,
    "health_rate_limit": 100,
    "health_rate_window": "1m"
  }
}`
}