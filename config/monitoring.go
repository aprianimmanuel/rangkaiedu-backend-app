package config

import (
	"time"
)

// MonitoringConfig represents the configuration for the monitoring system
type MonitoringConfig struct {
	Enabled            bool          `json:"enabled" mapstructure:"enabled"`
	Interval           time.Duration `json:"interval" mapstructure:"interval"`
	Timeout            time.Duration `json:"timeout" mapstructure:"timeout"`
	MaxHistorySize     int           `json:"max_history_size" mapstructure:"max_history_size"`
	CleanupInterval    time.Duration `json:"cleanup_interval" mapstructure:"cleanup_interval"`
	EnableMetrics      bool          `json:"enable_metrics" mapstructure:"enable_metrics"`
	EnableErrors       bool          `json:"enable_errors" mapstructure:"enable_errors"`
	EnableAlerts       bool          `json:"enable_alerts" mapstructure:"enable_alerts"`
	EnableHealthChecks bool          `json:"enable_health_checks" mapstructure:"enable_health_checks"`
	EnableSecurityLogs bool          `json:"enable_security_logs" mapstructure:"enable_security_logs"`
	MetricsConfig      *MetricsConfig `json:"metrics_config" mapstructure:"metrics_config"`
	ErrorConfig        *ErrorConfig  `json:"error_config" mapstructure:"error_config"`
	AlertConfig        *AlertConfig  `json:"alert_config" mapstructure:"alert_config"`
	HealthConfig       *HealthConfig `json:"health_config" mapstructure:"health_config"`
	SecurityConfig     *SecurityLoggingConfig `json:"security_config" mapstructure:"security_config"`
}

// MetricsConfig represents the configuration for metrics collection
type MetricsConfig struct {
	Enabled                bool          `json:"enabled" mapstructure:"enabled"`
	Interval               time.Duration `json:"interval" mapstructure:"interval"`
	Timeout                time.Duration `json:"timeout" mapstructure:"timeout"`
	PrometheusPort         int           `json:"prometheus_port" mapstructure:"prometheus_port"`
	PrometheusPath         string        `json:"prometheus_path" mapstructure:"prometheus_path"`
	EnableSystemMetrics    bool          `json:"enable_system_metrics" mapstructure:"enable_system_metrics"`
	EnableApplicationMetrics bool        `json:"enable_application_metrics" mapstructure:"enable_application_metrics"`
	EnableCustomMetrics    bool          `json:"enable_custom_metrics" mapstructure:"enable_custom_metrics"`
	EnableBusinessMetrics  bool          `json:"enable_business_metrics" mapstructure:"enable_business_metrics"`
	MaxMetricsPerSecond    int           `json:"max_metrics_per_second" mapstructure:"max_metrics_per_second"`
	MetricRetention        time.Duration `json:"metric_retention" mapstructure:"metric_retention"`
	EnableHistograms       bool          `json:"enable_histograms" mapstructure:"enable_histograms"`
	EnableSummaries        bool          `json:"enable_summaries" mapstructure:"enable_summaries"`
	EnableCounters         bool          `json:"enable_counters" mapstructure:"enable_counters"`
	EnableGauges           bool          `json:"enable_gauges" mapstructure:"enable_gauges"`
	EnableProcessMetrics   bool          `json:"enable_process_metrics" mapstructure:"enable_process_metrics"`
	EnableRuntimeMetrics   bool          `json:"enable_runtime_metrics" mapstructure:"enable_runtime_metrics"`
	EnableMemoryMetrics    bool          `json:"enable_memory_metrics" mapstructure:"enable_memory_metrics"`
	EnableCPUMetrics       bool          `json:"enable_cpu_metrics" mapstructure:"enable_cpu_metrics"`
	EnableDiskMetrics      bool          `json:"enable_disk_metrics" mapstructure:"enable_disk_metrics"`
	EnableNetworkMetrics   bool          `json:"enable_network_metrics" mapstructure:"enable_network_metrics"`
}

// ErrorConfig represents the configuration for error handling
type ErrorConfig struct {
	Enabled                   bool          `json:"enabled" mapstructure:"enabled"`
	MaxHistorySize            int           `json:"max_history_size" mapstructure:"max_history_size"`
	CleanupInterval           time.Duration `json:"cleanup_interval" mapstructure:"cleanup_interval"`
	EnableErrorTracking       bool          `json:"enable_error_tracking" mapstructure:"enable_error_tracking"`
	EnableErrorRecovery       bool          `json:"enable_error_recovery" mapstructure:"enable_error_recovery"`
	EnableErrorAlerts         bool          `json:"enable_error_alerts" mapstructure:"enable_error_alerts"`
	EnableErrorLogging        bool          `json:"enable_error_logging" mapstructure:"enable_error_logging"`
	MaxRetries                int           `json:"max_retries" mapstructure:"max_retries"`
	RetryDelay                time.Duration `json:"retry_delay" mapstructure:"retry_delay"`
	EnableErrorClassification bool          `json:"enable_error_classification" mapstructure:"enable_error_classification"`
	EnableErrorAggregation    bool          `json:"enable_error_aggregation" mapstructure:"enable_error_aggregation"`
	EnableErrorRateLimiting   bool          `json:"enable_error_rate_limiting" mapstructure:"enable_error_rate_limiting"`
	ErrorRateLimit            int           `json:"error_rate_limit" mapstructure:"error_rate_limit"`
	ErrorRateWindow           time.Duration `json:"error_rate_window" mapstructure:"error_rate_window"`
	EnableErrorContext        bool          `json:"enable_error_context" mapstructure:"enable_error_context"`
	EnableErrorStackTraces    bool          `json:"enable_error_stack_traces" mapstructure:"enable_error_stack_traces"`
	EnableErrorMetrics        bool          `json:"enable_error_metrics" mapstructure:"enable_error_metrics"`
}

// AlertConfig represents the configuration for alerting
type AlertConfig struct {
	Enabled                    bool          `json:"enabled" mapstructure:"enabled"`
	MaxHistorySize             int           `json:"max_history_size" mapstructure:"max_history_size"`
	CleanupInterval            time.Duration `json:"cleanup_interval" mapstructure:"cleanup_interval"`
	EnableAlertRules           bool          `json:"enable_alert_rules" mapstructure:"enable_alert_rules"`
	EnableAlertProviders       bool          `json:"enable_alert_providers" mapstructure:"enable_alert_providers"`
	EnableAlertTemplates       bool          `json:"enable_alert_templates" mapstructure:"enable_alert_templates"`
	EnableAlertRouting         bool          `json:"enable_alert_routing" mapstructure:"enable_alert_routing"`
	EnableAlertAggregation     bool          `json:"enable_alert_aggregation" mapstructure:"enable_alert_aggregation"`
	EnableAlertDeduplication   bool          `json:"enable_alert_deduplication" mapstructure:"enable_alert_deduplication"`
	EnableAlertEscalation      bool          `json:"enable_alert_escalation" mapstructure:"enable_alert_escalation"`
	EnableAlertSilencing       bool          `json:"enable_alert_silencing" mapstructure:"enable_alert_silencing"`
	EnableAlertMetrics         bool          `json:"enable_alert_metrics" mapstructure:"enable_alert_metrics"`
	DefaultSeverity            AlertSeverity `json:"default_severity" mapstructure:"default_severity"`
	DefaultTimeout             time.Duration `json:"default_timeout" mapstructure:"default_timeout"`
	AutoResolveTimeout         time.Duration `json:"auto_resolve_timeout" mapstructure:"auto_resolve_timeout"`
	EnableAlertGroups          bool          `json:"enable_alert_groups" mapstructure:"enable_alert_groups"`
	MaxAlertsPerMinute         int           `json:"max_alerts_per_minute" mapstructure:"max_alerts_per_minute"`
	AlertRateLimit             int           `json:"alert_rate_limit" mapstructure:"alert_rate_limit"`
	AlertRateWindow            time.Duration `json:"alert_rate_window" mapstructure:"alert_rate_window"`
}

// HealthConfig represents the configuration for health checks
type HealthConfig struct {
	Enabled                   bool          `json:"enabled" mapstructure:"enabled"`
	Interval                  time.Duration `json:"interval" mapstructure:"interval"`
	Timeout                   time.Duration `json:"timeout" mapstructure:"timeout"`
	MaxHistorySize            int           `json:"max_history_size" mapstructure:"max_history_size"`
	CleanupInterval           time.Duration `json:"cleanup_interval" mapstructure:"cleanup_interval"`
	EnableSystemHealth        bool          `json:"enable_system_health" mapstructure:"enable_system_health"`
	EnableApplicationHealth   bool          `json:"enable_application_health" mapstructure:"enable_application_health"`
	EnableDatabaseHealth      bool          `json:"enable_database_health" mapstructure:"enable_database_health"`
	EnableNetworkHealth       bool          `json:"enable_network_health" mapstructure:"enable_network_health"`
	EnableDependencyHealth    bool          `json:"enable_dependency_health" mapstructure:"enable_dependency_health"`
	EnableCustomHealth        bool          `json:"enable_custom_health" mapstructure:"enable_custom_health"`
	EnableHealthMetrics       bool          `json:"enable_health_metrics" mapstructure:"enable_health_metrics"`
	EnableHealthLogging       bool          `json:"enable_health_logging" mapstructure:"enable_health_logging"`
	EnableHealthAlerts        bool          `json:"enable_health_alerts" mapstructure:"enable_health_alerts"`
	DefaultThreshold          HealthStatus  `json:"default_threshold" mapstructure:"default_threshold"`
	EnableHealthGroups        bool          `json:"enable_health_groups" mapstructure:"enable_health_groups"`
	MaxHealthChecksPerSecond  int           `json:"max_health_checks_per_second" mapstructure:"max_health_checks_per_second"`
	HealthRateLimit           int           `json:"health_rate_limit" mapstructure:"health_rate_limit"`
	HealthRateWindow          time.Duration `json:"health_rate_window" mapstructure:"health_rate_window"`
}

// DefaultMonitoringConfig returns the default monitoring configuration
func DefaultMonitoringConfig() *MonitoringConfig {
	return &MonitoringConfig{
		Enabled:            true,
		Interval:           30 * time.Second,
		Timeout:            10 * time.Second,
		MaxHistorySize:     1000,
		CleanupInterval:    1 * time.Hour,
		EnableMetrics:      true,
		EnableErrors:       true,
		EnableAlerts:       true,
		EnableHealthChecks: true,
		EnableSecurityLogs: true,
		MetricsConfig:      DefaultMetricsConfig(),
		ErrorConfig:        DefaultErrorConfig(),
		AlertConfig:        DefaultAlertConfig(),
		HealthConfig:       DefaultHealthConfig(),
		SecurityConfig:     DefaultSecurityLoggingConfig(),
	}
}

// DefaultMetricsConfig returns the default metrics configuration
func DefaultMetricsConfig() *MetricsConfig {
	return &MetricsConfig{
		Enabled:                true,
		Interval:               15 * time.Second,
		Timeout:                5 * time.Second,
		PrometheusPort:         8080,
		PrometheusPath:         "/metrics",
		EnableSystemMetrics:    true,
		EnableApplicationMetrics: true,
		EnableCustomMetrics:    true,
		EnableBusinessMetrics:  true,
		MaxMetricsPerSecond:    100,
		MetricRetention:        24 * time.Hour,
		EnableHistograms:       true,
		EnableSummaries:        true,
		EnableCounters:         true,
		EnableGauges:           true,
		EnableProcessMetrics:   true,
		EnableRuntimeMetrics:   true,
		EnableMemoryMetrics:    true,
		EnableCPUMetrics:       true,
		EnableDiskMetrics:      true,
		EnableNetworkMetrics:   true,
	}
}

// DefaultErrorConfig returns the default error configuration
func DefaultErrorConfig() *ErrorConfig {
	return &ErrorConfig{
		Enabled:                   true,
		MaxHistorySize:            1000,
		CleanupInterval:           1 * time.Hour,
		EnableErrorTracking:       true,
		EnableErrorRecovery:       true,
		EnableErrorAlerts:         true,
		EnableErrorLogging:        true,
		MaxRetries:                3,
		RetryDelay:                5 * time.Second,
		EnableErrorClassification: true,
		EnableErrorAggregation:    true,
		EnableErrorRateLimiting:   true,
		ErrorRateLimit:            100,
		ErrorRateWindow:           time.Minute,
		EnableErrorContext:        true,
		EnableErrorStackTraces:    true,
		EnableErrorMetrics:        true,
	}
}

// DefaultAlertConfig returns the default alert configuration
func DefaultAlertConfig() *AlertConfig {
	return &AlertConfig{
		Enabled:                    true,
		MaxHistorySize:             1000,
		CleanupInterval:            1 * time.Hour,
		EnableAlertRules:           true,
		EnableAlertProviders:       true,
		EnableAlertTemplates:       true,
		EnableAlertRouting:         true,
		EnableAlertAggregation:     true,
		EnableAlertDeduplication:   true,
		EnableAlertEscalation:      true,
		EnableAlertSilencing:       true,
		EnableAlertMetrics:         true,
		DefaultSeverity:            AlertSeverityWarning,
		DefaultTimeout:             1 * time.Hour,
		AutoResolveTimeout:         2 * time.Hour,
		EnableAlertGroups:          true,
		MaxAlertsPerMinute:         10,
		AlertRateLimit:             100,
		AlertRateWindow:            time.Minute,
	}
}

// DefaultHealthConfig returns the default health configuration
func DefaultHealthConfig() *HealthConfig {
	return &HealthConfig{
		Enabled:                   true,
		Interval:                  30 * time.Second,
		Timeout:                   10 * time.Second,
		MaxHistorySize:            1000,
		CleanupInterval:           1 * time.Hour,
		EnableSystemHealth:        true,
		EnableApplicationHealth:   true,
		EnableDatabaseHealth:      true,
		EnableNetworkHealth:       true,
		EnableDependencyHealth:    true,
		EnableCustomHealth:        true,
		EnableHealthMetrics:       true,
		EnableHealthLogging:       true,
		EnableHealthAlerts:        true,
		DefaultThreshold:          HealthStatusHealthy,
		EnableHealthGroups:        true,
		MaxHealthChecksPerSecond:  10,
		HealthRateLimit:           100,
		HealthRateWindow:          time.Minute,
	}
}