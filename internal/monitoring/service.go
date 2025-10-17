package monitoring

import (
	"context"
	"fmt"
	"sync"
	"time"
)


// MonitoringService represents the main monitoring service
type MonitoringService struct {
	mu               sync.RWMutex
	metricsCollector MetricsCollector
	errorHandler     ErrorHandler
	alertHandler     AlertHandler
	healthChecker    HealthChecker
	securityLogger   SecurityEventLogger
	config           *MonitoringConfig
	running          bool
	stopChan         chan struct{}
}

// MonitoringConfig represents the configuration for the monitoring service
type MonitoringConfig struct {
	Enabled            bool          `json:"enabled"`
	Interval           time.Duration `json:"interval"`
	Timeout            time.Duration `json:"timeout"`
	MaxHistorySize     int           `json:"max_history_size"`
	CleanupInterval    time.Duration `json:"cleanup_interval"`
	EnableMetrics      bool          `json:"enable_metrics"`
	EnableErrors       bool          `json:"enable_errors"`
	EnableAlerts       bool          `json:"enable_alerts"`
	EnableHealthChecks bool          `json:"enable_health_checks"`
	EnableSecurityLogs bool          `json:"enable_security_logs"`
	MetricsConfig      *MetricsConfig `json:"metrics_config"`
	ErrorConfig        *ErrorConfig  `json:"error_config"`
	AlertConfig        *AlertConfig  `json:"alert_config"`
	HealthConfig       *HealthConfig `json:"health_config"`
	SecurityConfig     *SecurityLoggingConfig `json:"security_config"`
}

// DefaultMonitoringConfig returns the default configuration for the monitoring service
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

// MetricsConfig represents the configuration for metrics collection
type MetricsConfig struct {
	Enabled                bool          `json:"enabled"`
	Interval               time.Duration `json:"interval"`
	Timeout                time.Duration `json:"timeout"`
	PrometheusPort         int           `json:"prometheus_port"`
	PrometheusPath         string        `json:"prometheus_path"`
	EnableSystemMetrics    bool          `json:"enable_system_metrics"`
	EnableApplicationMetrics bool        `json:"enable_application_metrics"`
	EnableCustomMetrics    bool          `json:"enable_custom_metrics"`
	EnableBusinessMetrics  bool          `json:"enable_business_metrics"`
	MaxMetricsPerSecond    int           `json:"max_metrics_per_second"`
	MetricRetention        time.Duration `json:"metric_retention"`
	EnableHistograms       bool          `json:"enable_histograms"`
	EnableSummaries        bool          `json:"enable_summaries"`
	EnableCounters         bool          `json:"enable_counters"`
	EnableGauges           bool          `json:"enable_gauges"`
	EnableProcessMetrics   bool          `json:"enable_process_metrics"`
	EnableRuntimeMetrics   bool          `json:"enable_runtime_metrics"`
	EnableMemoryMetrics    bool          `json:"enable_memory_metrics"`
	EnableCPUMetrics       bool          `json:"enable_cpu_metrics"`
	EnableDiskMetrics      bool          `json:"enable_disk_metrics"`
	EnableNetworkMetrics   bool          `json:"enable_network_metrics"`
}

// ErrorConfig represents the configuration for error handling
type ErrorConfig struct {
	Enabled                   bool          `json:"enabled"`
	MaxHistorySize            int           `json:"max_history_size"`
	CleanupInterval           time.Duration `json:"cleanup_interval"`
	EnableErrorTracking       bool          `json:"enable_error_tracking"`
	EnableErrorRecovery       bool          `json:"enable_error_recovery"`
	EnableErrorAlerts         bool          `json:"enable_error_alerts"`
	EnableErrorLogging        bool          `json:"enable_error_logging"`
	MaxRetries                int           `json:"max_retries"`
	RetryDelay                time.Duration `json:"retry_delay"`
	EnableErrorClassification bool          `json:"enable_error_classification"`
	EnableErrorAggregation    bool          `json:"enable_error_aggregation"`
	EnableErrorRateLimiting   bool          `json:"enable_error_rate_limiting"`
	ErrorRateLimit            int           `json:"error_rate_limit"`
	ErrorRateWindow           time.Duration `json:"error_rate_window"`
	EnableErrorContext        bool          `json:"enable_error_context"`
	EnableErrorStackTraces    bool          `json:"enable_error_stack_traces"`
	EnableErrorMetrics        bool          `json:"enable_error_metrics"`
}

// AlertConfig represents the configuration for alerting
type AlertConfig struct {
	Enabled                    bool          `json:"enabled"`
	MaxHistorySize             int           `json:"max_history_size"`
	CleanupInterval            time.Duration `json:"cleanup_interval"`
	EnableAlertRules           bool          `json:"enable_alert_rules"`
	EnableAlertProviders       bool          `json:"enable_alert_providers"`
	EnableAlertTemplates       bool          `json:"enable_alert_templates"`
	EnableAlertRouting         bool          `json:"enable_alert_routing"`
	EnableAlertAggregation     bool          `json:"enable_alert_aggregation"`
	EnableAlertDeduplication   bool          `json:"enable_alert_deduplication"`
	EnableAlertEscalation      bool          `json:"enable_alert_escalation"`
	EnableAlertSilencing       bool          `json:"enable_alert_silencing"`
	EnableAlertMetrics         bool          `json:"enable_alert_metrics"`
	DefaultSeverity            AlertSeverity `json:"default_severity"`
	DefaultTimeout             time.Duration `json:"default_timeout"`
	AutoResolveTimeout         time.Duration `json:"auto_resolve_timeout"`
	EnableAlertGroups          bool          `json:"enable_alert_groups"`
	MaxAlertsPerMinute         int           `json:"max_alerts_per_minute"`
	AlertRateLimit             int           `json:"alert_rate_limit"`
	AlertRateWindow            time.Duration `json:"alert_rate_window"`
}

// HealthConfig represents the configuration for health checks
type HealthConfig struct {
	Enabled                   bool          `json:"enabled"`
	Interval                  time.Duration `json:"interval"`
	Timeout                   time.Duration `json:"timeout"`
	MaxHistorySize            int           `json:"max_history_size"`
	CleanupInterval           time.Duration `json:"cleanup_interval"`
	EnableSystemHealth        bool          `json:"enable_system_health"`
	EnableApplicationHealth   bool          `json:"enable_application_health"`
	EnableDatabaseHealth      bool          `json:"enable_database_health"`
	EnableNetworkHealth       bool          `json:"enable_network_health"`
	EnableDependencyHealth    bool          `json:"enable_dependency_health"`
	EnableCustomHealth        bool          `json:"enable_custom_health"`
	EnableHealthMetrics       bool          `json:"enable_health_metrics"`
	EnableHealthLogging       bool          `json:"enable_health_logging"`
	EnableHealthAlerts        bool          `json:"enable_health_alerts"`
	DefaultThreshold          HealthStatus  `json:"default_threshold"`
	EnableHealthGroups        bool          `json:"enable_health_groups"`
	MaxHealthChecksPerSecond  int           `json:"max_health_checks_per_second"`
	HealthRateLimit           int           `json:"health_rate_limit"`
	HealthRateWindow          time.Duration `json:"health_rate_window"`
}

// DefaultMetricsConfig returns the default configuration for metrics collection
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

// DefaultErrorConfig returns the default configuration for error handling
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

// DefaultAlertConfig returns the default configuration for alerting
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

// DefaultHealthConfig returns the default configuration for health checks
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


// PrometheusConfig represents Prometheus configuration
type PrometheusConfig struct {
	Address   string `json:"address"`
	Port      int    `json:"port"`
	Namespace string `json:"namespace"`
	Subsystem string `json:"subsystem"`
	Path      string `json:"path"`
}

// MetricsInfo represents information about metrics
type MetricsInfo struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Tags        map[string]string      `json:"tags"`
	Enabled     bool                   `json:"enabled"`
	LastUpdated time.Time              `json:"last_updated"`
	Config      map[string]interface{} `json:"config"`
}

// SecurityEvent represents a security event
type SecurityEvent struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Type        string                 `json:"type"`
	UserID      string                 `json:"user_id,omitempty"`
	IP          string                 `json:"ip,omitempty"`
	Severity    string                 `json:"severity"`
	Resource    string                 `json:"resource"`
	Action      string                 `json:"action"`
	Status      string                 `json:"status"`
	Description string                 `json:"description"`
	Labels      map[string]string      `json:"labels"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

// MetricsCollector defines the interface for metrics collection
type MetricsCollector interface {
	RecordMetric(metric map[string]interface{}) error
	GetMetrics() map[string]interface{}
	GetMetricsStats() map[string]interface{}
	GetMetricsCount() int
	IsEnabled() bool
}

// SecurityEventLogger defines the interface for security event logging
type SecurityEventLogger interface {
	LogEvent(ctx context.Context, event *SecurityEvent) error
	GetEvents(ctx context.Context, limit int, offset int) ([]*SecurityEvent, error)
	GetEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffset(ctx context.Context, startTime, endTime time.Time, eventType, userID, ip, severity, resource, action, status string, labelKey, labelValue string, sortBy string, sortOrder string, page, pageSize, limit, offset int) ([]*SecurityEvent, int, error)
	GetEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrder(ctx context.Context, startTime, endTime time.Time, eventType, userID, ip, severity, resource, action, status string, labelKey, labelValue string, sortBy string, sortOrder string, page, pageSize, limit, offset int, order string) ([]*SecurityEvent, int, error)
	GetEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrderAndDirection(ctx context.Context, startTime, endTime time.Time, eventType, userID, ip, severity, resource, action, status string, labelKey, labelValue string, sortBy string, sortOrder string, page, pageSize, limit, offset int, order, direction string) ([]*SecurityEvent, int, error)
	GetEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrderAndDirectionAndPage(ctx context.Context, startTime, endTime time.Time, eventType, userID, ip, severity, resource, action, status string, labelKey, labelValue string, sortBy string, sortOrder string, page, pageSize, limit, offset int, order, direction string) ([]*SecurityEvent, int, error)
	GetConfig() map[string]interface{}
	Close() error
	
	// Authentication event logging methods
	LogAuthFailure(ctx context.Context, username, ip, userAgent string, details map[string]interface{}) error
	LogAuthSuccess(ctx context.Context, userID, username, ip, userAgent string, details map[string]interface{}) error
	LogAuthLockout(ctx context.Context, username, ip string, attempts int, details map[string]interface{}) error
	LogAuthBruteForce(ctx context.Context, ip string, attempts int, timeframe time.Duration, details map[string]interface{}) error
	LogAuthSession(ctx context.Context, userID, sessionID, ip string, action string, details map[string]interface{}) error
	LogAuthToken(ctx context.Context, userID, tokenID, ip string, action string, details map[string]interface{}) error
	LogAuthMFA(ctx context.Context, userID, username, ip, method string, success bool, details map[string]interface{}) error
	LogAuthViolation(ctx context.Context, userID, username, resource, action string, details map[string]interface{}) error
	LogAuthThreat(ctx context.Context, threatType, ip, userAgent string, details map[string]interface{}) error
	LogAuthAudit(ctx context.Context, userID, username, action, resource string, details map[string]interface{}) error
}

// Validate validates the monitoring configuration
func (c *MonitoringConfig) Validate() error {
	if c.Interval <= 0 {
		return fmt.Errorf("interval must be positive")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	if c.MaxHistorySize <= 0 {
		return fmt.Errorf("max_history_size must be positive")
	}
	if c.CleanupInterval <= 0 {
		return fmt.Errorf("cleanup_interval must be positive")
	}
	return nil
}

// Validate validates the metrics configuration
func (c *MetricsConfig) Validate() error {
	if c.Interval <= 0 {
		return fmt.Errorf("interval must be positive")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	if c.PrometheusPort <= 0 {
		return fmt.Errorf("prometheus_port must be positive")
	}
	if c.MaxMetricsPerSecond <= 0 {
		return fmt.Errorf("max_metrics_per_second must be positive")
	}
	if c.MetricRetention <= 0 {
		return fmt.Errorf("metric_retention must be positive")
	}
	return nil
}

// Validate validates the error configuration
func (c *ErrorConfig) Validate() error {
	if c.MaxHistorySize <= 0 {
		return fmt.Errorf("max_history_size must be positive")
	}
	if c.CleanupInterval <= 0 {
		return fmt.Errorf("cleanup_interval must be positive")
	}
	if c.MaxRetries < 0 {
		return fmt.Errorf("max_retries must be non-negative")
	}
	if c.RetryDelay < 0 {
		return fmt.Errorf("retry_delay must be non-negative")
	}
	if c.ErrorRateLimit <= 0 {
		return fmt.Errorf("error_rate_limit must be positive")
	}
	if c.ErrorRateWindow <= 0 {
		return fmt.Errorf("error_rate_window must be positive")
	}
	return nil
}

// Validate validates the alert configuration
func (c *AlertConfig) Validate() error {
	if c.MaxHistorySize <= 0 {
		return fmt.Errorf("max_history_size must be positive")
	}
	if c.CleanupInterval <= 0 {
		return fmt.Errorf("cleanup_interval must be positive")
	}
	if c.DefaultTimeout <= 0 {
		return fmt.Errorf("default_timeout must be positive")
	}
	if c.AutoResolveTimeout <= 0 {
		return fmt.Errorf("auto_resolve_timeout must be positive")
	}
	if c.MaxAlertsPerMinute <= 0 {
		return fmt.Errorf("max_alerts_per_minute must be positive")
	}
	if c.AlertRateLimit <= 0 {
		return fmt.Errorf("alert_rate_limit must be positive")
	}
	if c.AlertRateWindow <= 0 {
		return fmt.Errorf("alert_rate_window must be positive")
	}
	return nil
}

// Validate validates the health configuration
func (c *HealthConfig) Validate() error {
	if c.Interval <= 0 {
		return fmt.Errorf("interval must be positive")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	if c.MaxHistorySize <= 0 {
		return fmt.Errorf("max_history_size must be positive")
	}
	if c.CleanupInterval <= 0 {
		return fmt.Errorf("cleanup_interval must be positive")
	}
	if c.MaxHealthChecksPerSecond <= 0 {
		return fmt.Errorf("max_health_checks_per_second must be positive")
	}
	if c.HealthRateLimit <= 0 {
		return fmt.Errorf("health_rate_limit must be positive")
	}
	if c.HealthRateWindow <= 0 {
		return fmt.Errorf("health_rate_window must be positive")
	}
	return nil
}

// Merge merges the current configuration with another configuration
func (c *MonitoringConfig) Merge(other *MonitoringConfig) *MonitoringConfig {
	if other == nil {
		return c
	}

	merged := &MonitoringConfig{
		Enabled:            other.Enabled || c.Enabled,
		Interval:           other.Interval,
		Timeout:            other.Timeout,
		MaxHistorySize:     other.MaxHistorySize,
		CleanupInterval:    other.CleanupInterval,
		EnableMetrics:      other.EnableMetrics || c.EnableMetrics,
		EnableErrors:       other.EnableErrors || c.EnableErrors,
		EnableAlerts:       other.EnableAlerts || c.EnableAlerts,
		EnableHealthChecks: other.EnableHealthChecks || c.EnableHealthChecks,
		EnableSecurityLogs: other.EnableSecurityLogs || c.EnableSecurityLogs,
	}

	if other.MetricsConfig != nil {
		merged.MetricsConfig = other.MetricsConfig
	} else if c.MetricsConfig != nil {
		merged.MetricsConfig = c.MetricsConfig
	}

	if other.ErrorConfig != nil {
		merged.ErrorConfig = other.ErrorConfig
	} else if c.ErrorConfig != nil {
		merged.ErrorConfig = c.ErrorConfig
	}

	if other.AlertConfig != nil {
		merged.AlertConfig = other.AlertConfig
	} else if c.AlertConfig != nil {
		merged.AlertConfig = c.AlertConfig
	}

	if other.HealthConfig != nil {
		merged.HealthConfig = other.HealthConfig
	} else if c.HealthConfig != nil {
		merged.HealthConfig = c.HealthConfig
	}

	if other.SecurityConfig != nil {
		merged.SecurityConfig = other.SecurityConfig
	} else if c.SecurityConfig != nil {
		merged.SecurityConfig = c.SecurityConfig
	}

	return merged
}

// Clone creates a copy of the configuration
func (c *MonitoringConfig) Clone() *MonitoringConfig {
	clone := &MonitoringConfig{
		Enabled:            c.Enabled,
		Interval:           c.Interval,
		Timeout:            c.Timeout,
		MaxHistorySize:     c.MaxHistorySize,
		CleanupInterval:    c.CleanupInterval,
		EnableMetrics:      c.EnableMetrics,
		EnableErrors:       c.EnableErrors,
		EnableAlerts:       c.EnableAlerts,
		EnableHealthChecks: c.EnableHealthChecks,
		EnableSecurityLogs: c.EnableSecurityLogs,
	}

	if c.MetricsConfig != nil {
		clone.MetricsConfig = c.MetricsConfig
	}

	if c.ErrorConfig != nil {
		clone.ErrorConfig = c.ErrorConfig
	}

	if c.AlertConfig != nil {
		clone.AlertConfig = c.AlertConfig
	}

	if c.HealthConfig != nil {
		clone.HealthConfig = c.HealthConfig
	}

	if c.SecurityConfig != nil {
		clone.SecurityConfig = c.SecurityConfig
	}

	return clone
}

// NewMonitoringService creates a new monitoring service
func NewMonitoringService(config *MonitoringConfig) (*MonitoringService, error) {
	if config == nil {
		config = DefaultMonitoringConfig()
	}
	
	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}
	
	service := &MonitoringService{
		config:  config,
		stopChan: make(chan struct{}),
	}
	
	// Initialize components
	if err := service.initializeComponents(); err != nil {
		return nil, fmt.Errorf("failed to initialize monitoring components: %v", err)
	}
	
	return service, nil
}

// initializeComponents initializes all monitoring components
func (s *MonitoringService) initializeComponents() error {
	// Initialize metrics collector
	if s.config.EnableMetrics && s.config.MetricsConfig != nil {
		metricsCollector, err := NewMetricsCollector("prometheus", &PrometheusConfig{
			Address:   "0.0.0.0",
			Port:      s.config.MetricsConfig.PrometheusPort,
			Namespace: "app",
			Subsystem: "monitoring",
			Path:      s.config.MetricsConfig.PrometheusPath,
		}, MetricsInfo{
			Name:        "monitoring",
			Type:        "prometheus",
			Version:     "1.0.0",
			Description: "Application monitoring metrics",
			Tags:        map[string]string{"service": "monitoring"},
			Enabled:     true,
			LastUpdated: time.Now(),
		})
		if err != nil {
			return fmt.Errorf("failed to create metrics collector: %v", err)
		}
		s.metricsCollector = metricsCollector
	}
	
	// Initialize error handler
	if s.config.EnableErrors && s.config.ErrorConfig != nil {
		s.errorHandler = NewErrorHandler(s.config.ErrorConfig.MaxHistorySize)
	}
	
	// Initialize alert handler
	if s.config.EnableAlerts && s.config.AlertConfig != nil {
		s.alertHandler = NewAlertHandler(s.config.AlertConfig.MaxHistorySize)
	}
	
	// Initialize health checker
	if s.config.EnableHealthChecks && s.config.HealthConfig != nil {
		healthConfig := &HealthCheckConfig{
			ID:          "system",
			Name:        "System Health",
			Type:        HealthCheckTypeSystem,
			Enabled:     true,
			Interval:    s.config.HealthConfig.Interval,
			Timeout:     s.config.HealthConfig.Timeout,
			Retries:     3,
			Threshold:   HealthStatusHealthy,
			Labels:      map[string]string{"service": "monitoring"},
			Config:      make(map[string]interface{}),
		}
		s.healthChecker = NewHealthChecker(healthConfig)
	}
	
	// Initialize security logger
	if s.config.EnableSecurityLogs && s.config.SecurityConfig != nil {
		securityLogger, err := NewSecurityEventLogger(s.config.SecurityConfig)
		if err != nil {
			return fmt.Errorf("failed to create security event logger: %v", err)
		}
		s.securityLogger = securityLogger
	}
	
	return nil
}

// Start starts the monitoring service
func (s *MonitoringService) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.running {
		return fmt.Errorf("monitoring service is already running")
	}
	
	s.running = true
	
	// Start all components
	if s.metricsCollector != nil {
		if err := s.startMetricsCollection(ctx); err != nil {
			return fmt.Errorf("failed to start metrics collection: %v", err)
		}
	}
	
	if s.healthChecker != nil {
		if err := s.healthChecker.Start(ctx); err != nil {
			return fmt.Errorf("failed to start health checker: %v", err)
		}
	}
	
	// Start cleanup goroutine
	go s.cleanupLoop(ctx)
	
	return nil
}

// Stop stops the monitoring service
func (s *MonitoringService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if !s.running {
		return fmt.Errorf("monitoring service is not running")
	}
	
	s.running = false
	close(s.stopChan)
	
	// Stop all components
	if s.healthChecker != nil {
		if err := s.healthChecker.Stop(); err != nil {
			return fmt.Errorf("failed to stop health checker: %v", err)
		}
	}
	
	// Close security logger
	if s.securityLogger != nil {
		if err := s.securityLogger.Close(); err != nil {
			return fmt.Errorf("failed to close security logger: %v", err)
		}
	}
	
	return nil
}

// IsRunning returns whether the monitoring service is running
func (s *MonitoringService) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return s.running
}

// GetMetrics returns metrics data
func (s *MonitoringService) GetMetrics() map[string]interface{} {
	if s.metricsCollector == nil {
		return make(map[string]interface{})
	}
	
	return s.metricsCollector.GetMetrics()
}

// GetErrorStats returns error statistics
func (s *MonitoringService) GetErrorStats() ErrorStats {
	if s.errorHandler == nil {
		return ErrorStats{}
	}
	
	return s.errorHandler.GetErrorStats()
}

// GetAlertStats returns alert statistics
func (s *MonitoringService) GetAlertStats() AlertStats {
	if s.alertHandler == nil {
		return AlertStats{}
	}
	
	return s.alertHandler.GetAlertStats()
}

// GetHealthStats returns health check statistics
func (s *MonitoringService) GetHealthStats() HealthStats {
	if s.healthChecker == nil {
		return HealthStats{}
	}
	
	return s.healthChecker.GetHealthStats()
}

// GetSystemHealth returns the overall system health
func (s *MonitoringService) GetSystemHealth(ctx context.Context) (*HealthCheckResult, error) {
	if s.healthChecker == nil {
		return nil, fmt.Errorf("health checker is not enabled")
	}
	
	return s.healthChecker.GetSystemHealth(ctx)
}

// GetMonitoringStats returns overall monitoring statistics
func (s *MonitoringService) GetMonitoringStats() map[string]interface{} {
	stats := make(map[string]interface{})
	
	stats["running"] = s.IsRunning()
	stats["config"] = s.config
	
	if s.metricsCollector != nil {
		stats["metrics"] = s.metricsCollector.GetMetricsStats()
	}
	
	if s.errorHandler != nil {
		stats["errors"] = s.errorHandler.GetErrorStats()
	}
	
	if s.alertHandler != nil {
		stats["alerts"] = s.alertHandler.GetAlertStats()
	}
	
	if s.healthChecker != nil {
		stats["health"] = s.healthChecker.GetHealthStats()
	}
	
	if s.securityLogger != nil {
		stats["security"] = s.securityLogger.GetConfig()
	}
	
	return stats
}

// HandleError handles a monitoring error
func (s *MonitoringService) HandleError(ctx context.Context, err *MonitoringError) error {
	if s.errorHandler == nil {
		return fmt.Errorf("error handler is not enabled")
	}
	
	return s.errorHandler.HandleError(ctx, err)
}

// ProcessAlert processes an alert
func (s *MonitoringService) ProcessAlert(ctx context.Context, alert *Alert) error {
	if s.alertHandler == nil {
		return fmt.Errorf("alert handler is not enabled")
	}
	
	return s.alertHandler.ProcessAlert(ctx, alert)
}

// RunHealthCheck runs a health check
func (s *MonitoringService) RunHealthCheck(ctx context.Context, id string) (*HealthCheckResult, error) {
	if s.healthChecker == nil {
		return nil, fmt.Errorf("health checker is not enabled")
	}
	
	return s.healthChecker.RunHealthCheck(ctx, id)
}

// RegisterHealthCheck registers a health check
func (s *MonitoringService) RegisterHealthCheck(check HealthCheck) error {
	if s.healthChecker == nil {
		return fmt.Errorf("health checker is not enabled")
	}
	
	return s.healthChecker.RegisterHealthCheck(check)
}


// GetMetricsCollector returns the metrics collector
func (s *MonitoringService) GetMetricsCollector() MetricsCollector {
	return s.metricsCollector
}

// GetErrorHandler returns the error handler
func (s *MonitoringService) GetErrorHandler() ErrorHandler {
	return s.errorHandler
}

// GetAlertHandler returns the alert handler
func (s *MonitoringService) GetAlertHandler() AlertHandler {
	return s.alertHandler
}

// GetHealthChecker returns the health checker
func (s *MonitoringService) GetHealthChecker() HealthChecker {
	return s.healthChecker
}

// SetMetricsCollector sets the metrics collector
func (s *MonitoringService) SetMetricsCollector(collector MetricsCollector) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.metricsCollector = collector
}

// SetErrorHandler sets the error handler
func (s *MonitoringService) SetErrorHandler(handler ErrorHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.errorHandler = handler
}

// SetAlertHandler sets the alert handler
func (s *MonitoringService) SetAlertHandler(handler AlertHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.alertHandler = handler
}

// SetHealthChecker sets the health checker
func (s *MonitoringService) SetHealthChecker(checker HealthChecker) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.healthChecker = checker
}

// GetConfig returns the monitoring configuration
func (s *MonitoringService) GetConfig() *MonitoringConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return s.config
}

// SetConfig sets the monitoring configuration
func (s *MonitoringService) SetConfig(config *MonitoringConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.config = config
}

// EnableMetrics enables metrics collection
func (s *MonitoringService) EnableMetrics() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.config.EnableMetrics = true
}

// DisableMetrics disables metrics collection
func (s *MonitoringService) DisableMetrics() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.config.EnableMetrics = false
}

// EnableErrors enables error handling
func (s *MonitoringService) EnableErrors() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.config.EnableErrors = true
}

// DisableErrors disables error handling
func (s *MonitoringService) DisableErrors() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.config.EnableErrors = false
}

// EnableAlerts enables alerting
func (s *MonitoringService) EnableAlerts() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.config.EnableAlerts = true
}

// DisableAlerts disables alerting
func (s *MonitoringService) DisableAlerts() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.config.EnableAlerts = false
}

// EnableHealthChecks enables health checks
func (s *MonitoringService) EnableHealthChecks() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.config.EnableHealthChecks = true
}

// DisableHealthChecks disables health checks
func (s *MonitoringService) DisableHealthChecks() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.config.EnableHealthChecks = false
}

// EnableSecurityLogs enables security logging
func (s *MonitoringService) EnableSecurityLogs() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.config.EnableSecurityLogs = true
}

// DisableSecurityLogs disables security logging
func (s *MonitoringService) DisableSecurityLogs() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.config.EnableSecurityLogs = false
}

// IsMetricsEnabled returns whether metrics collection is enabled
func (s *MonitoringService) IsMetricsEnabled() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return s.config.EnableMetrics
}

// IsErrorsEnabled returns whether error handling is enabled
func (s *MonitoringService) IsErrorsEnabled() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return s.config.EnableErrors
}

// IsAlertsEnabled returns whether alerting is enabled
func (s *MonitoringService) IsAlertsEnabled() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return s.config.EnableAlerts
}

// IsHealthChecksEnabled returns whether health checks are enabled
func (s *MonitoringService) IsHealthChecksEnabled() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return s.config.EnableHealthChecks
}

// IsSecurityLogsEnabled returns whether security logging is enabled
func (s *MonitoringService) IsSecurityLogsEnabled() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return s.config.EnableSecurityLogs
}

// GetEnabledComponents returns the list of enabled components
func (s *MonitoringService) GetEnabledComponents() []string {
	var components []string
	
	if s.IsMetricsEnabled() {
		components = append(components, "metrics")
	}
	
	if s.IsErrorsEnabled() {
		components = append(components, "errors")
	}
	
	if s.IsAlertsEnabled() {
		components = append(components, "alerts")
	}
	
	if s.IsHealthChecksEnabled() {
		components = append(components, "health_checks")
	}
	
	if s.IsSecurityLogsEnabled() {
		components = append(components, "security_logs")
	}
	
	return components
}

// GetDisabledComponents returns the list of disabled components
func (s *MonitoringService) GetDisabledComponents() []string {
	var components []string
	
	if !s.IsMetricsEnabled() {
		components = append(components, "metrics")
	}
	
	if !s.IsErrorsEnabled() {
		components = append(components, "errors")
	}
	
	if !s.IsAlertsEnabled() {
		components = append(components, "alerts")
	}
	
	if !s.IsHealthChecksEnabled() {
		components = append(components, "health_checks")
	}
	
	if !s.IsSecurityLogsEnabled() {
		components = append(components, "security_logs")
	}
	
	return components
}

// EnableAllComponents enables all monitoring components
func (s *MonitoringService) EnableAllComponents() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.config.EnableMetrics = true
	s.config.EnableErrors = true
	s.config.EnableAlerts = true
	s.config.EnableHealthChecks = true
	s.config.EnableSecurityLogs = true
}

// DisableAllComponents disables all monitoring components
func (s *MonitoringService) DisableAllComponents() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.config.EnableMetrics = false
	s.config.EnableErrors = false
	s.config.EnableAlerts = false
	s.config.EnableHealthChecks = false
	s.config.EnableSecurityLogs = false
}

// startMetricsCollection starts metrics collection
func (s *MonitoringService) startMetricsCollection(ctx context.Context) error {
	if s.metricsCollector == nil {
		return fmt.Errorf("metrics collector is not initialized")
	}
	
	// Start metrics collection in a goroutine
	go func() {
		ticker := time.NewTicker(s.config.Interval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				// Instead of calling Collect directly, we'll record a simple metric
				metric := map[string]interface{}{
					"name":      "monitoring_collection",
					"value":     1.0,
					"timestamp": time.Now().Unix(),
				}
				if err := s.metricsCollector.RecordMetric(metric); err != nil {
					// Log error
					continue
				}
			case <-s.stopChan:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
	
	return nil
}

// cleanupLoop runs the cleanup loop
func (s *MonitoringService) cleanupLoop(ctx context.Context) {
	ticker := time.NewTicker(s.config.CleanupInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			// Cleanup old data
			s.cleanupOldData()
		case <-s.stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}

// cleanupOldData cleans up old monitoring data
func (s *MonitoringService) cleanupOldData() {
	// Cleanup old error history
	if s.errorHandler != nil {
		// This would require the error handler to have a cleanup method
		// For now, we'll just call it if it exists
		if cleaner, ok := s.errorHandler.(interface{ CleanupExpiredErrors() }); ok {
			cleaner.CleanupExpiredErrors()
		}
	}
	
	// Cleanup old alert history
	if s.alertHandler != nil {
		// This would require the alert handler to have a cleanup method
		// For now, we'll just call it if it exists
		if cleaner, ok := s.alertHandler.(interface{ CleanupOldAlerts(time.Duration) }); ok {
			cleaner.CleanupOldAlerts(24 * time.Hour * 7) // Keep 7 days of alerts
		}
	}
	
	// Cleanup old health check history
	if s.healthChecker != nil {
		// This would require the health checker to have a cleanup method
		// For now, we'll just call it if it exists
		if cleaner, ok := s.healthChecker.(interface{ ClearHealthHistory() }); ok {
			cleaner.ClearHealthHistory()
		}
	}
}

// GetMetricsCount returns the number of metrics
func (s *MonitoringService) GetMetricsCount() int {
	if s.metricsCollector == nil {
		return 0
	}
	
	if counter, ok := s.metricsCollector.(interface{ GetMetricsCount() int }); ok {
		return counter.GetMetricsCount()
	}
	
	return 0
}

// GetErrorCount returns the number of errors
func (s *MonitoringService) GetErrorCount() int {
	if s.errorHandler == nil {
		return 0
	}
	
	return s.errorHandler.GetErrorStats().TotalErrors
}

// GetAlertCount returns the number of alerts
func (s *MonitoringService) GetAlertCount() int {
	if s.alertHandler == nil {
		return 0
	}
	
	return s.alertHandler.GetAlertStats().TotalAlerts
}

// GetHealthCheckCount returns the number of health checks
func (s *MonitoringService) GetHealthCheckCount() int {
	if s.healthChecker == nil {
		return 0
	}
	
	return s.healthChecker.GetHealthChecksCount()
}

// GetTotalCount returns the total number of monitoring items
func (s *MonitoringService) GetTotalCount() int {
	return s.GetMetricsCount() + s.GetErrorCount() + s.GetAlertCount() + s.GetHealthCheckCount()
}

// GetMemoryUsage returns the memory usage of monitoring components
func (s *MonitoringService) GetMemoryUsage() int64 {
	var total int64
	
	if s.metricsCollector != nil {
		if counter, ok := s.metricsCollector.(interface{ GetMetricsMemoryUsage() int64 }); ok {
			total += counter.GetMetricsMemoryUsage()
		}
	}
	
	// Add memory usage for other components
	// This would require each component to have a GetMemoryUsage method
	
	return total
}

// GetVersion returns the monitoring service version
func (s *MonitoringService) GetVersion() string {
	return "1.0.0"
}

// GetInfo returns monitoring service information
func (s *MonitoringService) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"version":     s.GetVersion(),
		"running":     s.IsRunning(),
		"enabled":     s.config.Enabled,
		"components":  s.GetEnabledComponents(),
		"total_count": s.GetTotalCount(),
		"memory_usage": s.GetMemoryUsage(),
		"config":      s.config,
		"security_enabled": s.IsSecurityLogsEnabled(),
	}
}

// ValidateConfig validates the monitoring configuration
func (s *MonitoringService) ValidateConfig() error {
	if s.config == nil {
		return fmt.Errorf("configuration is nil")
	}
	
	if s.config.Interval <= 0 {
		return fmt.Errorf("interval must be positive")
	}
	
	if s.config.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	
	if s.config.MaxHistorySize <= 0 {
		return fmt.Errorf("max history size must be positive")
	}
	
	if s.config.CleanupInterval <= 0 {
		return fmt.Errorf("cleanup interval must be positive")
	}
	
	return nil
}

// ReloadConfig reloads the monitoring configuration
func (s *MonitoringService) ReloadConfig(config *MonitoringConfig) error {
	if err := s.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	oldConfig := s.config
	s.config = config
	
	// Restart service if it was running
	if s.running {
		if err := s.Stop(); err != nil {
			s.config = oldConfig
			return fmt.Errorf("failed to stop monitoring service: %v", err)
		}
		
		ctx := context.Background()
		if err := s.Start(ctx); err != nil {
			s.config = oldConfig
			return fmt.Errorf("failed to restart monitoring service: %v", err)
		}
	}
	
	return nil
}

// GetConfigValidationErrors returns configuration validation errors
func (s *MonitoringService) GetConfigValidationErrors() []string {
	var errors []string
	
	if s.config == nil {
		errors = append(errors, "configuration is nil")
		return errors
	}
	
	if s.config.Interval <= 0 {
		errors = append(errors, "interval must be positive")
	}
	
	if s.config.Timeout <= 0 {
		errors = append(errors, "timeout must be positive")
	}
	
	if s.config.MaxHistorySize <= 0 {
		errors = append(errors, "max history size must be positive")
	}
	
	if s.config.CleanupInterval <= 0 {
		errors = append(errors, "cleanup interval must be positive")
	}
	
	return errors
}

// IsConfigValid returns whether the configuration is valid
func (s *MonitoringService) IsConfigValid() bool {
	return len(s.GetConfigValidationErrors()) == 0
}

// GetHealthStatus returns the health status of the monitoring service
func (s *MonitoringService) GetHealthStatus() HealthStatus {
	if !s.IsRunning() {
		return HealthStatusUnhealthy
	}
	
	// Check if all enabled components are healthy
	if s.IsMetricsEnabled() && s.metricsCollector != nil {
		if !s.metricsCollector.IsEnabled() {
			return HealthStatusDegraded
		}
	}
	
	if s.IsErrorsEnabled() && s.errorHandler != nil {
		// Check error handler health
		if stats := s.errorHandler.GetErrorStats(); stats.TotalErrors > 100 {
			return HealthStatusDegraded
		}
	}
	
	if s.IsAlertsEnabled() && s.alertHandler != nil {
		// Check alert handler health
		if stats := s.alertHandler.GetAlertStats(); stats.ActiveAlerts > 50 {
			return HealthStatusDegraded
		}
	}
	
	if s.IsHealthChecksEnabled() && s.healthChecker != nil {
		// Check health checker health
		if !s.healthChecker.IsRunning() {
			return HealthStatusDegraded
		}
	}
	
	if s.IsSecurityLogsEnabled() && s.securityLogger != nil {
		// Check security logger health
		if config := s.securityLogger.GetConfig(); config != nil {
			// Check if the security logger is active instead of checking config.Enabled
			// DEBUG: Log interface type and method availability
			fmt.Printf("DEBUG: Security logger type: %T\n", s.securityLogger)
			fmt.Printf("DEBUG: Security logger interface: %T\n", (interface{})(s.securityLogger))
			
			// Check if IsActive method exists via type assertion
			if activeLogger, ok := s.securityLogger.(interface{ IsActive() bool }); ok {
				fmt.Printf("DEBUG: IsActive method found via type assertion\n")
				if !activeLogger.IsActive() {
					return HealthStatusDegraded
				}
			} else {
				fmt.Printf("DEBUG: IsActive method NOT found via type assertion\n")
				// Fallback: check config.Enabled directly
				if enabledConfig, ok := config["enabled"].(bool); ok && !enabledConfig {
					return HealthStatusDegraded
				}
			}
		}
	}
	
	return HealthStatusHealthy
}

// GetHealthMessage returns the health message for the monitoring service
func (s *MonitoringService) GetHealthMessage() string {
	status := s.GetHealthStatus()
	switch status {
	case HealthStatusHealthy:
		return "Monitoring service is healthy"
	case HealthStatusDegraded:
		return "Monitoring service is degraded"
	case HealthStatusUnhealthy:
		return "Monitoring service is unhealthy"
	case HealthStatusUnknown:
		return "Monitoring service health is unknown"
	default:
		return "Unknown monitoring service health status"
	}
}

// GetSecurityLogger returns the security event logger
func (s *MonitoringService) GetSecurityLogger() SecurityEventLogger {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.securityLogger
}

// SetSecurityLogger sets the security event logger
func (s *MonitoringService) SetSecurityLogger(logger SecurityEventLogger) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.securityLogger = logger
}

// LogSecurityEvent logs a security event
func (s *MonitoringService) LogSecurityEvent(ctx context.Context, event *SecurityEvent) error {
	if s.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	return s.securityLogger.LogEvent(ctx, event)
}

// LogAuthFailure logs an authentication failure event
func (s *MonitoringService) LogAuthFailure(ctx context.Context, username, ip, userAgent string, details map[string]interface{}) error {
	if s.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	return s.securityLogger.LogAuthFailure(ctx, username, ip, userAgent, details)
}

// LogAuthSuccess logs an authentication success event
func (s *MonitoringService) LogAuthSuccess(ctx context.Context, userID, username, ip, userAgent string, details map[string]interface{}) error {
	if s.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	return s.securityLogger.LogAuthSuccess(ctx, userID, username, ip, userAgent, details)
}

// LogAuthLockout logs an authentication lockout event
func (s *MonitoringService) LogAuthLockout(ctx context.Context, username, ip string, attempts int, details map[string]interface{}) error {
	if s.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	return s.securityLogger.LogAuthLockout(ctx, username, ip, attempts, details)
}

// LogAuthBruteForce logs a brute force attack event
func (s *MonitoringService) LogAuthBruteForce(ctx context.Context, ip string, attempts int, timeframe time.Duration, details map[string]interface{}) error {
	if s.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	return s.securityLogger.LogAuthBruteForce(ctx, ip, attempts, timeframe, details)
}

// LogAuthSession logs an authentication session event
func (s *MonitoringService) LogAuthSession(ctx context.Context, userID, sessionID, ip string, action string, details map[string]interface{}) error {
	if s.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	return s.securityLogger.LogAuthSession(ctx, userID, sessionID, ip, action, details)
}

// LogAuthToken logs an authentication token event
func (s *MonitoringService) LogAuthToken(ctx context.Context, userID, tokenID, ip string, action string, details map[string]interface{}) error {
	if s.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	return s.securityLogger.LogAuthToken(ctx, userID, tokenID, ip, action, details)
}

// LogAuthMFA logs a multi-factor authentication event
func (s *MonitoringService) LogAuthMFA(ctx context.Context, userID, username, ip, method string, success bool, details map[string]interface{}) error {
	if s.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	return s.securityLogger.LogAuthMFA(ctx, userID, username, ip, method, success, details)
}

// LogAuthViolation logs an authorization violation event
func (s *MonitoringService) LogAuthViolation(ctx context.Context, userID, username, resource, action string, details map[string]interface{}) error {
	if s.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	return s.securityLogger.LogAuthViolation(ctx, userID, username, resource, action, details)
}

// LogAuthThreat logs a security threat event
func (s *MonitoringService) LogAuthThreat(ctx context.Context, threatType, ip, userAgent string, details map[string]interface{}) error {
	if s.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	return s.securityLogger.LogAuthThreat(ctx, threatType, ip, userAgent, details)
}

// LogAuthAudit logs an audit event
func (s *MonitoringService) LogAuthAudit(ctx context.Context, userID, username, action, resource string, details map[string]interface{}) error {
	if s.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	return s.securityLogger.LogAuthAudit(ctx, userID, username, action, resource, details)
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(collectorType string, config *PrometheusConfig, info MetricsInfo) (MetricsCollector, error) {
	switch collectorType {
	case "prometheus":
		return NewPrometheusMetricsCollector(config, info)
	default:
		return NewDefaultMetricsCollector(info)
	}
}

// NewSecurityEventLogger creates a new security event logger
func NewSecurityEventLogger(config *SecurityLoggingConfig) (SecurityEventLogger, error) {
	return NewDefaultSecurityEventLogger(config)
}

// NewPrometheusMetricsCollector creates a Prometheus metrics collector
func NewPrometheusMetricsCollector(config *PrometheusConfig, info MetricsInfo) (MetricsCollector, error) {
	return &prometheusMetricsCollector{
		config: config,
		info:   info,
		metrics: make(map[string]interface{}),
	}, nil
}

// NewDefaultMetricsCollector creates a default metrics collector
func NewDefaultMetricsCollector(info MetricsInfo) (MetricsCollector, error) {
	return &defaultMetricsCollector{
		info:    info,
		metrics: make(map[string]interface{}),
	}, nil
}

// NewDefaultSecurityEventLogger creates a default security event logger
func NewDefaultSecurityEventLogger(config *SecurityLoggingConfig) (SecurityEventLogger, error) {
	return &defaultSecurityEventLogger{
		config:  config,
		events:  make([]*SecurityEvent, 0),
		configs: make(map[string]interface{}),
	}, nil
}

// prometheusMetricsCollector implements MetricsCollector for Prometheus
type prometheusMetricsCollector struct {
	config  *PrometheusConfig
	info    MetricsInfo
	metrics map[string]interface{}
	mu      sync.RWMutex
}

func (p *prometheusMetricsCollector) RecordMetric(metric map[string]interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.metrics[fmt.Sprintf("%s_%s", p.info.Name, metric["name"].(string))] = metric
	return nil
}

func (p *prometheusMetricsCollector) GetMetrics() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	result := make(map[string]interface{})
	for k, v := range p.metrics {
		result[k] = v
	}
	return result
}

func (p *prometheusMetricsCollector) GetMetricsStats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	return map[string]interface{}{
		"total_metrics": len(p.metrics),
		"collector_type": "prometheus",
		"namespace":     p.info.Name,
		"enabled":       p.info.Enabled,
	}
}

func (p *prometheusMetricsCollector) GetMetricsCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	return len(p.metrics)
}

func (p *prometheusMetricsCollector) IsEnabled() bool {
	return p.info.Enabled
}

// defaultMetricsCollector implements MetricsCollector with default behavior
type defaultMetricsCollector struct {
	info    MetricsInfo
	metrics map[string]interface{}
	mu      sync.RWMutex
}

func (d *defaultMetricsCollector) RecordMetric(metric map[string]interface{}) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	d.metrics[fmt.Sprintf("%s_%s", d.info.Name, metric["name"].(string))] = metric
	return nil
}

func (d *defaultMetricsCollector) GetMetrics() map[string]interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	result := make(map[string]interface{})
	for k, v := range d.metrics {
		result[k] = v
	}
	return result
}

func (d *defaultMetricsCollector) GetMetricsStats() map[string]interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	return map[string]interface{}{
		"total_metrics": len(d.metrics),
		"collector_type": "default",
		"namespace":     d.info.Name,
		"enabled":       d.info.Enabled,
	}
}

func (d *defaultMetricsCollector) GetMetricsCount() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	return len(d.metrics)
}

func (d *defaultMetricsCollector) IsEnabled() bool {
	return d.info.Enabled
}

// defaultSecurityEventLogger implements SecurityEventLogger with default behavior
type defaultSecurityEventLogger struct {
	config  *SecurityLoggingConfig
	events  []*SecurityEvent
	configs map[string]interface{}
	mu      sync.RWMutex
}

func (d *defaultSecurityEventLogger) LogEvent(ctx context.Context, event *SecurityEvent) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	event.ID = fmt.Sprintf("event_%d", time.Now().UnixNano())
	event.Timestamp = time.Now()
	d.events = append(d.events, event)
	return nil
}

func (d *defaultSecurityEventLogger) GetEvents(ctx context.Context, limit int, offset int) ([]*SecurityEvent, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if limit <= 0 {
		limit = len(d.events)
	}
	if offset < 0 {
		offset = 0
	}
	
	if offset >= len(d.events) {
		return []*SecurityEvent{}, nil
	}
	
	end := offset + limit
	if end > len(d.events) {
		end = len(d.events)
	}
	
	return d.events[offset:end], nil
}

func (d *defaultSecurityEventLogger) GetEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffset(ctx context.Context, startTime, endTime time.Time, eventType, userID, ip, severity, resource, action, status string, labelKey, labelValue string, sortBy string, sortOrder string, page, pageSize, limit, offset int) ([]*SecurityEvent, int, error) {
	// Implementation would filter and sort events based on parameters
	// For now, return all events
	events, err := d.GetEvents(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return events, len(d.events), nil
}

func (d *defaultSecurityEventLogger) GetEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrder(ctx context.Context, startTime, endTime time.Time, eventType, userID, ip, severity, resource, action, status string, labelKey, labelValue string, sortBy string, sortOrder string, page, pageSize, limit, offset int, order string) ([]*SecurityEvent, int, error) {
	// Implementation would filter and sort events based on parameters
	// For now, return all events
	events, err := d.GetEvents(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return events, len(d.events), nil
}

func (d *defaultSecurityEventLogger) GetEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrderAndDirection(ctx context.Context, startTime, endTime time.Time, eventType, userID, ip, severity, resource, action, status string, labelKey, labelValue string, sortBy string, sortOrder string, page, pageSize, limit, offset int, order, direction string) ([]*SecurityEvent, int, error) {
	// Implementation would filter and sort events based on parameters
	// For now, return all events
	events, err := d.GetEvents(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return events, len(d.events), nil
}

func (d *defaultSecurityEventLogger) GetEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrderAndDirectionAndPage(ctx context.Context, startTime, endTime time.Time, eventType, userID, ip, severity, resource, action, status string, labelKey, labelValue string, sortBy string, sortOrder string, page, pageSize, limit, offset int, order, direction string) ([]*SecurityEvent, int, error) {
	// Implementation would filter and sort events based on parameters
	// For now, return all events
	events, err := d.GetEvents(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return events, len(d.events), nil
}

func (d *defaultSecurityEventLogger) GetConfig() map[string]interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	config := make(map[string]interface{})
	for k, v := range d.configs {
		config[k] = v
	}
	return config
}

func (d *defaultSecurityEventLogger) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	d.events = nil
	d.configs = nil
	return nil
}

// LogAuthFailure logs an authentication failure event
func (d *defaultSecurityEventLogger) LogAuthFailure(ctx context.Context, username, ip, userAgent string, details map[string]interface{}) error {
	event := &SecurityEvent{
		Type:        "auth_failure",
		UserID:      username,
		IP:          ip,
		Severity:    "warning",
		Resource:    "authentication",
		Action:      "login",
		Status:      "failed",
		Description: fmt.Sprintf("Authentication failure for user %s from IP %s", username, ip),
		Labels: map[string]string{
			"username":   username,
			"ip":         ip,
			"user_agent": userAgent,
		},
		Context: details,
	}
	return d.LogEvent(ctx, event)
}

// LogAuthSuccess logs an authentication success event
func (d *defaultSecurityEventLogger) LogAuthSuccess(ctx context.Context, userID, username, ip, userAgent string, details map[string]interface{}) error {
	event := &SecurityEvent{
		Type:        "auth_success",
		UserID:      userID,
		IP:          ip,
		Severity:    "info",
		Resource:    "authentication",
		Action:      "login",
		Status:      "success",
		Description: fmt.Sprintf("Authentication success for user %s from IP %s", username, ip),
		Labels: map[string]string{
			"username":   username,
			"ip":         ip,
			"user_agent": userAgent,
		},
		Context: details,
	}
	return d.LogEvent(ctx, event)
}

// LogAuthLockout logs an authentication lockout event
func (d *defaultSecurityEventLogger) LogAuthLockout(ctx context.Context, username, ip string, attempts int, details map[string]interface{}) error {
	event := &SecurityEvent{
		Type:        "auth_lockout",
		UserID:      username,
		IP:          ip,
		Severity:    "critical",
		Resource:    "authentication",
		Action:      "lockout",
		Status:      "locked",
		Description: fmt.Sprintf("User %s locked out after %d failed attempts from IP %s", username, attempts, ip),
		Labels: map[string]string{
			"username": username,
			"ip":       ip,
			"attempts": fmt.Sprintf("%d", attempts),
		},
		Context: details,
	}
	return d.LogEvent(ctx, event)
}

// LogAuthBruteForce logs a brute force attack event
func (d *defaultSecurityEventLogger) LogAuthBruteForce(ctx context.Context, ip string, attempts int, timeframe time.Duration, details map[string]interface{}) error {
	event := &SecurityEvent{
		Type:        "auth_bruteforce",
		IP:          ip,
		Severity:    "critical",
		Resource:    "authentication",
		Action:      "bruteforce",
		Status:      "detected",
		Description: fmt.Sprintf("Brute force attack detected from IP %s with %d attempts in %v", ip, attempts, timeframe),
		Labels: map[string]string{
			"ip":       ip,
			"attempts": fmt.Sprintf("%d", attempts),
			"timeframe": timeframe.String(),
		},
		Context: details,
	}
	return d.LogEvent(ctx, event)
}

// LogAuthSession logs an authentication session event
func (d *defaultSecurityEventLogger) LogAuthSession(ctx context.Context, userID, sessionID, ip string, action string, details map[string]interface{}) error {
	event := &SecurityEvent{
		Type:        "auth_session",
		UserID:      userID,
		IP:          ip,
		Severity:    "info",
		Resource:    "authentication",
		Action:      action,
		Status:      "active",
		Description: fmt.Sprintf("Session %s for user %s from IP %s", action, userID, ip),
		Labels: map[string]string{
			"user_id":    userID,
			"session_id": sessionID,
			"ip":         ip,
			"action":     action,
		},
		Context: details,
	}
	return d.LogEvent(ctx, event)
}

// LogAuthToken logs an authentication token event
func (d *defaultSecurityEventLogger) LogAuthToken(ctx context.Context, userID, tokenID, ip string, action string, details map[string]interface{}) error {
	event := &SecurityEvent{
		Type:        "auth_token",
		UserID:      userID,
		IP:          ip,
		Severity:    "info",
		Resource:    "authentication",
		Action:      action,
		Status:      "active",
		Description: fmt.Sprintf("Token %s for user %s from IP %s", action, userID, ip),
		Labels: map[string]string{
			"user_id": userID,
			"token_id": tokenID,
			"ip":      ip,
			"action":  action,
		},
		Context: details,
	}
	return d.LogEvent(ctx, event)
}

// LogAuthMFA logs a multi-factor authentication event
func (d *defaultSecurityEventLogger) LogAuthMFA(ctx context.Context, userID, username, ip, method string, success bool, details map[string]interface{}) error {
	status := "failed"
	severity := "warning"
	if success {
		status = "success"
		severity = "info"
	}
	
	event := &SecurityEvent{
		Type:        "auth_mfa",
		UserID:      userID,
		IP:          ip,
		Severity:    severity,
		Resource:    "authentication",
		Action:      "mfa",
		Status:      status,
		Description: fmt.Sprintf("MFA %s for user %s using %s from IP %s", status, username, method, ip),
		Labels: map[string]string{
			"user_id":  userID,
			"username": username,
			"ip":       ip,
			"method":   method,
			"success":  fmt.Sprintf("%t", success),
		},
		Context: details,
	}
	return d.LogEvent(ctx, event)
}

// LogAuthViolation logs an authorization violation event
func (d *defaultSecurityEventLogger) LogAuthViolation(ctx context.Context, userID, username, resource, action string, details map[string]interface{}) error {
	event := &SecurityEvent{
		Type:        "auth_violation",
		UserID:      userID,
		Severity:    "critical",
		Resource:    resource,
		Action:      action,
		Status:      "denied",
		Description: fmt.Sprintf("Authorization violation by user %s on resource %s with action %s", username, resource, action),
		Labels: map[string]string{
			"user_id":  userID,
			"username": username,
			"resource": resource,
			"action":   action,
		},
		Context: details,
	}
	return d.LogEvent(ctx, event)
}

// LogAuthThreat logs a security threat event
func (d *defaultSecurityEventLogger) LogAuthThreat(ctx context.Context, threatType, ip, userAgent string, details map[string]interface{}) error {
	event := &SecurityEvent{
		Type:        "auth_threat",
		IP:          ip,
		Severity:    "critical",
		Resource:    "authentication",
		Action:      "threat",
		Status:      "detected",
		Description: fmt.Sprintf("Security threat %s detected from IP %s", threatType, ip),
		Labels: map[string]string{
			"threat_type": threatType,
			"ip":          ip,
			"user_agent":  userAgent,
		},
		Context: details,
	}
	return d.LogEvent(ctx, event)
}

// LogAuthAudit logs an audit event
func (d *defaultSecurityEventLogger) LogAuthAudit(ctx context.Context, userID, username, action, resource string, details map[string]interface{}) error {
	event := &SecurityEvent{
		Type:        "auth_audit",
		UserID:      userID,
		Severity:    "info",
		Resource:    resource,
		Action:      action,
		Status:      "logged",
		Description: fmt.Sprintf("Audit log for user %s on resource %s with action %s", username, resource, action),
		Labels: map[string]string{
			"user_id":  userID,
			"username": username,
			"resource": resource,
			"action":   action,
		},
		Context: details,
	}
	return d.LogEvent(ctx, event)
}

// GetEventCount returns the number of security events
func (d *defaultSecurityEventLogger) GetEventCount() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.events)
}

// ClearEvents clears all security events
func (d *defaultSecurityEventLogger) ClearEvents() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.events = make([]*SecurityEvent, 0)
	return nil
}

// SetConfig sets the security logging configuration
func (d *defaultSecurityEventLogger) SetConfig(config *SecurityLoggingConfig) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.config = config
	return nil
}

// Enable enables the security logger
func (d *defaultSecurityEventLogger) Enable() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.config != nil {
		d.config.Enabled = true
	}
	return nil
}

// Disable disables the security logger
func (d *defaultSecurityEventLogger) Disable() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.config != nil {
		d.config.Enabled = false
	}
	return nil
}

// IsActive returns whether the security logger is active
func (d *defaultSecurityEventLogger) IsActive() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.config != nil {
		return d.config.Enabled
	}
	return false
}

// GetEventsByMultipleFilters returns security events filtered by multiple criteria
func (d *defaultSecurityEventLogger) GetEventsByMultipleFilters(ctx context.Context, filters map[string]interface{}, limit int, offset int) ([]*SecurityEvent, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if limit <= 0 {
		limit = len(d.events)
	}
	if offset < 0 {
		offset = 0
	}
	
	if offset >= len(d.events) {
		return []*SecurityEvent{}, nil
	}
	
	// Filter events based on the provided filters
	var filteredEvents []*SecurityEvent
	for _, event := range d.events {
		match := true
		
		for key, value := range filters {
			switch key {
			case "type":
				if event.Type != value {
					match = false
					break
				}
			case "user_id":
				if event.UserID != value {
					match = false
					break
				}
			case "ip":
				if event.IP != value {
					match = false
					break
				}
			case "severity":
				if event.Severity != value {
					match = false
					break
				}
			case "resource":
				if event.Resource != value {
					match = false
					break
				}
			case "action":
				if event.Action != value {
					match = false
					break
				}
			case "status":
				if event.Status != value {
					match = false
					break
				}
			case "label":
				// Handle label filtering
				if labelMap, ok := value.(map[string]string); ok {
					for labelKey, labelValue := range labelMap {
						if event.Labels[labelKey] != labelValue {
							match = false
							break
						}
					}
				}
			}
			
			if !match {
				break
			}
		}
		
		if match {
			filteredEvents = append(filteredEvents, event)
		}
	}
	
	// Apply limit and offset
	if offset < len(filteredEvents) {
		end := offset + limit
		if end > len(filteredEvents) {
			end = len(filteredEvents)
		}
		if limit > 0 {
			return filteredEvents[offset:end], nil
		}
		return filteredEvents[offset:], nil
	}
	
	return []*SecurityEvent{}, nil
}
