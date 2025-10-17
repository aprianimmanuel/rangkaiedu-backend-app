package monitoring

import (
	"time"
)

// AlertSeverity represents the severity level of an alert
type AlertSeverity string

const (
	AlertSeverityDebug    AlertSeverity = "debug"
	AlertSeverityInfo     AlertSeverity = "info"
	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityError    AlertSeverity = "error"
	AlertSeverityCritical AlertSeverity = "critical"
)

// AlertStatus represents the status of an alert
type AlertStatus string

const (
	AlertStatusActive    AlertStatus = "active"
	AlertStatusFiring    AlertStatus = "firing"
	AlertStatusResolved  AlertStatus = "resolved"
	AlertStatusSuppressed AlertStatus = "suppressed"
	AlertStatusPending   AlertStatus = "pending"
)

// AlertType represents the type of an alert
type AlertType string

const (
	AlertTypeSecurity    AlertType = "security"
	AlertTypePerformance AlertType = "performance"
	AlertTypeBusiness    AlertType = "business"
	AlertTypeSystem      AlertType = "system"
	AlertTypeCustom      AlertType = "custom"
)

// AlertProviderType represents the type of an alert provider
type AlertProviderType string

const (
	AlertProviderTypeSlack      AlertProviderType = "slack"
	AlertProviderTypeEmail      AlertProviderType = "email"
	AlertProviderTypeWebhook    AlertProviderType = "webhook"
	AlertProviderTypePagerDuty  AlertProviderType = "pagerduty"
	AlertProviderTypeSMS        AlertProviderType = "sms"
	AlertProviderTypePushbullet AlertProviderType = "pushbullet"
)

// ErrorSeverity represents the severity level of an error
type ErrorSeverity string

const (
	ErrorSeverityDebug    ErrorSeverity = "debug"
	ErrorSeverityInfo     ErrorSeverity = "info"
	ErrorSeverityWarning  ErrorSeverity = "warning"
	ErrorSeverityError    ErrorSeverity = "error"
	ErrorSeverityCritical ErrorSeverity = "critical"
)

// ErrorStatus represents the status of an error
type ErrorStatus string

const (
	ErrorStatusOpen      ErrorStatus = "open"
	ErrorStatusResolved  ErrorStatus = "resolved"
	ErrorStatusSuppressed ErrorStatus = "suppressed"
	ErrorStatusPending   ErrorStatus = "pending"
)

// ErrorType represents the type of an error
type ErrorType string

const (
	ErrorTypeError     ErrorType = "error"
	ErrorTypeErrorType ErrorType = "error_type"
	ErrorTypeErrorKind ErrorType = "error_kind"
)

// ErrorCategory represents the category of an error
type ErrorCategory string

const (
	ErrorCategorySystem    ErrorCategory = "system"
	ErrorCategoryNetwork   ErrorCategory = "network"
	ErrorCategoryDatabase  ErrorCategory = "database"
	ErrorCategoryBusiness  ErrorCategory = "business"
	ErrorCategorySecurity  ErrorCategory = "security"
	ErrorCategoryUnknown   ErrorCategory = "unknown"
)

// SecurityAlertSeverity represents the severity level of security alerts
type SecurityAlertSeverity string

const (
	SecurityAlertSeverityLow      SecurityAlertSeverity = "low"
	SecurityAlertSeverityMedium   SecurityAlertSeverity = "medium"
	SecurityAlertSeverityHigh     SecurityAlertSeverity = "high"
	SecurityAlertSeverityCritical SecurityAlertSeverity = "critical"
)

// RetryPolicy represents retry policy for operations
type RetryPolicy struct {
	MaxRetries    int           `json:"max_retries"`
	InitialDelay  time.Duration `json:"initial_delay"`
	MaxDelay      time.Duration `json:"max_delay"`
	BackoffFactor float64       `json:"backoff_factor"`
}

// AlertTemplates represents alert templates
type AlertTemplates struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Slack   string `json:"slack"`
	Email   string `json:"email"`
	Webhook string `json:"webhook"`
}

// HealthCheckType represents the type of health check
type HealthCheckType string

const (
	HealthCheckTypeHTTP    HealthCheckType = "http"
	HealthCheckTypeTCP     HealthCheckType = "tcp"
	HealthCheckTypeDatabase HealthCheckType = "database"
	HealthCheckTypeCustom  HealthCheckType = "custom"
	HealthCheckTypeService HealthCheckType = "service"
)

// HealthStatus represents the status of a health check
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnknown   HealthStatus = "unknown"
	HealthStatusPending   HealthStatus = "pending"
)

// AlertRule represents an alert rule
type AlertRule struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Enabled     bool                   `json:"enabled"`
	Conditions  []AlertCondition       `json:"conditions"`
	Actions     []AlertAction          `json:"actions"`
	Labels      map[string]string      `json:"labels"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// AlertCondition represents an alert condition
type AlertCondition struct {
	Type     string      `json:"type"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
	Target   string      `json:"target"`
}

// AlertAction represents an alert action
type AlertAction struct {
	Type    string                 `json:"type"`
	Config  map[string]interface{} `json:"config"`
	Enabled bool                   `json:"enabled"`
}

// AlertProvider represents an alert provider
type AlertProvider struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Config      map[string]interface{} `json:"config"`
	Enabled     bool                   `json:"enabled"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// HealthCheckTypeSystem represents the system health check type
const HealthCheckTypeSystem HealthCheckType = "system"
