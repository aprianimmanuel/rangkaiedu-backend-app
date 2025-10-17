package monitoring

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Alert represents a monitoring alert
type Alert struct {
	ID          string                 `json:"id"`
	Type        AlertType              `json:"type"`
	Severity    AlertSeverity          `json:"severity"`
	Status      AlertStatus            `json:"status"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Description string                 `json:"description,omitempty"`
	Labels      map[string]string      `json:"labels,omitempty"`
	Annotations map[string]string      `json:"annotations,omitempty"`
	StartsAt    time.Time              `json:"starts_at"`
	EndsAt      *time.Time             `json:"ends_at,omitempty"`
	Duration    time.Duration          `json:"duration"`
	Context     context.Context        `json:"-"`
	Providers   []string               `json:"providers"`
	Silenced    bool                   `json:"silenced"`
	SilencedBy  string                 `json:"silenced_by,omitempty"`
	SilencedAt  *time.Time             `json:"silenced_at,omitempty"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty"`
	FiringCount int                    `json:"firing_count"`
	LastFiredAt *time.Time             `json:"last_fired_at,omitempty"`
}

// AlertStats represents alert statistics
type AlertStats struct {
	TotalAlerts       int                    `json:"total_alerts"`
	ActiveAlerts      int                    `json:"active_alerts"`
	CriticalAlerts   int                    `json:"critical_alerts"`
	ErrorAlerts      int                    `json:"error_alerts"`
	WarningAlerts    int                    `json:"warning_alerts"`
	InfoAlerts       int                    `json:"info_alerts"`
	DebugAlerts      int                    `json:"debug_alerts"`
	AlertsByType      map[AlertType]int      `json:"alerts_by_type"`
	AlertsBySeverity  map[AlertSeverity]int  `json:"alerts_by_severity"`
	AlertsByStatus    map[AlertStatus]int    `json:"alerts_by_status"`
	AlertRate         float64                `json:"alert_rate"`
	AverageDuration   time.Duration          `json:"average_duration"`
	LastAlertTime     *time.Time             `json:"last_alert_time"`
	LastAlert         *Alert                `json:"last_alert,omitempty"`
}

// AlertHandler defines the interface for alert handling
type AlertHandler interface {
	// HandleAlert handles a monitoring alert
	HandleAlert(ctx context.Context, alert *Alert) error
	
	// ProcessAlert processes an alert
	ProcessAlert(ctx context.Context, alert *Alert) error
	
	// GetAlertStats returns alert statistics
	GetAlertStats() AlertStats
	
	// GetAlertCount returns the number of alerts
	GetAlertCount() int
	
	// IsEnabled returns whether alert handling is enabled
	IsEnabled() bool
	
	// GetAlerts returns all alerts
	GetAlerts(limit int, offset int) ([]*Alert, error)
	
	// GetAlertsBySeverity returns alerts by severity
	GetAlertsBySeverity(severity AlertSeverity, limit int, offset int) ([]*Alert, error)
	
	// GetAlertsByStatus returns alerts by status
	GetAlertsByStatus(status AlertStatus, limit int, offset int) ([]*Alert, error)
	
	// GetAlertsByTimeRange returns alerts by time range
	GetAlertsByTimeRange(startTime, endTime time.Time, limit int, offset int) ([]*Alert, error)
	
	// GetAlertsByType returns alerts by type
	GetAlertsByType(alertType AlertType, limit int, offset int) ([]*Alert, error)
	
	// ResolveAlert resolves an alert
	ResolveAlert(id string, resolvedBy string) error
	
	// ResolveAlertWithContext resolves an alert with context
	ResolveAlertWithContext(ctx context.Context, alertID string) error
	
	// SuppressAlert suppresses an alert
	SuppressAlert(id string, silencedBy string) error
	
	// ReactivateAlert reactivates an alert
	ReactivateAlert(id string) error
	
	// AddAlertLabel adds a label to an alert
	AddAlertLabel(id string, key, value string) error
	
	// RemoveAlertLabel removes a label from an alert
	RemoveAlertLabel(id string, key string) error
	
	// GetAlertLabels returns alert labels
	GetAlertLabels(id string) (map[string]string, error)
	
	// SetAlertContext sets alert context
	SetAlertContext(id string, context context.Context) error
	
	// GetAlertContext returns alert context
	GetAlertContext(id string) (context.Context, error)
	
	// CleanupOldAlerts cleans up old alerts
	CleanupOldAlerts(retentionPeriod time.Duration)
	
	// GetAlertHistory returns alert history
	GetAlertHistory(limit int) []*Alert
	
	// GetActiveAlerts returns active alerts
	GetActiveAlerts() []*Alert
	
	// GetAlert returns a specific alert
	GetAlert(alertID string) (*Alert, bool)
	
	// IsAlertActive returns whether an alert is active
	IsAlertActive(alertID string) bool
	
	// IsAlertSilenced returns whether an alert is silenced
	IsAlertSilenced(alertID string) bool
	
	// GetAlertDuration returns the duration of an alert
	GetAlertDuration(alertID string) (time.Duration, error)
	
	// GetAlertCountByType returns alert count by type
	GetAlertCountByType(alertType AlertType) int
	
	// GetAlertCountBySeverity returns alert count by severity
	GetAlertCountBySeverity(severity AlertSeverity) int
	
	// GetAlertCountByStatus returns alert count by status
	GetAlertCountByStatus(status AlertStatus) int
	
	// GetAverageAlertDuration returns the average duration of alerts
	GetAverageAlertDuration() time.Duration
	
	// GetAlertRate returns the alert rate (alerts per minute)
	GetAlertRate() float64
}

// AlertHandlerConfig represents the configuration for alert handling
type AlertHandlerConfig struct {
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

// DefaultAlertHandlerConfig returns the default configuration for alert handling
func DefaultAlertHandlerConfig() *AlertHandlerConfig {
	return &AlertHandlerConfig{
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
		DefaultTimeout:             5 * time.Minute,
		AutoResolveTimeout:        1 * time.Hour,
		EnableAlertGroups:          true,
		MaxAlertsPerMinute:         100,
		AlertRateLimit:            1000,
		AlertRateWindow:           time.Minute,
	}
}

// BaseAlertHandler provides the base implementation for alert handlers
type BaseAlertHandler struct {
	mu            sync.RWMutex
	alerts        map[string]*Alert
	alertStats    AlertStats
	config        *AlertHandlerConfig
	running       bool
	stopChan      chan struct{}
}

// NewAlertHandler creates a new alert handler
func NewAlertHandler(maxHistorySize int) AlertHandler {
	return &BaseAlertHandler{
		alerts:      make(map[string]*Alert),
		alertStats:  AlertStats{},
		running:     false,
		stopChan:    make(chan struct{}),
		config: &AlertHandlerConfig{
			MaxHistorySize: maxHistorySize,
		},
	}
}

// HandleAlert handles a monitoring alert
func (h *BaseAlertHandler) HandleAlert(ctx context.Context, alert *Alert) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	// Generate alert ID if not provided
	if alert.ID == "" {
		alert.ID = generateAlertID()
	}
	
	// Set alert status
	alert.Status = AlertStatusFiring
	alert.StartsAt = time.Now()
	
	// Update firing count
	now := time.Now()
	if existingAlert, exists := h.alerts[alert.ID]; exists {
		existingAlert.FiringCount++
		existingAlert.LastFiredAt = &now
		alert = existingAlert
	} else {
		alert.FiringCount = 1
		alert.LastFiredAt = &now
	}
	
	// Add to alerts
	h.alerts[alert.ID] = alert
	
	// Update statistics
	h.updateAlertStats(alert)
	
	return nil
}

// ProcessAlert processes an alert
func (h *BaseAlertHandler) ProcessAlert(ctx context.Context, alert *Alert) error {
	return h.HandleAlert(ctx, alert)
}

// GetAlertStats returns alert statistics
func (h *BaseAlertHandler) GetAlertStats() AlertStats {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.alertStats
}

// GetAlertCount returns the number of alerts
func (h *BaseAlertHandler) GetAlertCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return len(h.alerts)
}

// IsEnabled returns whether alert handling is enabled
func (h *BaseAlertHandler) IsEnabled() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.config.Enabled
}

// GetAlerts returns all alerts
func (h *BaseAlertHandler) GetAlerts(limit int, offset int) ([]*Alert, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	alerts := make([]*Alert, 0, len(h.alerts))
	for _, alert := range h.alerts {
		alerts = append(alerts, alert)
	}
	
	if limit <= 0 {
		limit = len(alerts)
	}
	if offset < 0 {
		offset = 0
	}
	
	if offset >= len(alerts) {
		return []*Alert{}, nil
	}
	
	end := offset + limit
	if end > len(alerts) {
		end = len(alerts)
	}
	
	return alerts[offset:end], nil
}

// GetAlertsBySeverity returns alerts by severity
func (h *BaseAlertHandler) GetAlertsBySeverity(severity AlertSeverity, limit int, offset int) ([]*Alert, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var alerts []*Alert
	for _, alert := range h.alerts {
		if alert.Severity == severity {
			alerts = append(alerts, alert)
		}
	}
	
	if limit <= 0 {
		limit = len(alerts)
	}
	if offset < 0 {
		offset = 0
	}
	
	if offset >= len(alerts) {
		return []*Alert{}, nil
	}
	
	end := offset + limit
	if end > len(alerts) {
		end = len(alerts)
	}
	
	return alerts[offset:end], nil
}

// GetAlertsByStatus returns alerts by status
func (h *BaseAlertHandler) GetAlertsByStatus(status AlertStatus, limit int, offset int) ([]*Alert, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var alerts []*Alert
	for _, alert := range h.alerts {
		if alert.Status == status {
			alerts = append(alerts, alert)
		}
	}
	
	if limit <= 0 {
		limit = len(alerts)
	}
	if offset < 0 {
		offset = 0
	}
	
	if offset >= len(alerts) {
		return []*Alert{}, nil
	}
	
	end := offset + limit
	if end > len(alerts) {
		end = len(alerts)
	}
	
	return alerts[offset:end], nil
}

// GetAlertsByTimeRange returns alerts by time range
func (h *BaseAlertHandler) GetAlertsByTimeRange(startTime, endTime time.Time, limit int, offset int) ([]*Alert, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var alerts []*Alert
	for _, alert := range h.alerts {
		if alert.StartsAt.After(startTime) && alert.StartsAt.Before(endTime) {
			alerts = append(alerts, alert)
		}
	}
	
	if limit <= 0 {
		limit = len(alerts)
	}
	if offset < 0 {
		offset = 0
	}
	
	if offset >= len(alerts) {
		return []*Alert{}, nil
	}
	
	end := offset + limit
	if end > len(alerts) {
		end = len(alerts)
	}
	
	return alerts[offset:end], nil
}

// GetAlertsByType returns alerts by type
func (h *BaseAlertHandler) GetAlertsByType(alertType AlertType, limit int, offset int) ([]*Alert, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var alerts []*Alert
	for _, alert := range h.alerts {
		if alert.Type == alertType {
			alerts = append(alerts, alert)
		}
	}
	
	if limit <= 0 {
		limit = len(alerts)
	}
	if offset < 0 {
		offset = 0
	}
	
	if offset >= len(alerts) {
		return []*Alert{}, nil
	}
	
	end := offset + limit
	if end > len(alerts) {
		end = len(alerts)
	}
	
	return alerts[offset:end], nil
}

// ResolveAlert resolves an alert
func (h *BaseAlertHandler) ResolveAlert(id string, resolvedBy string) error {
	return h.ResolveAlertWithContext(context.Background(), id)
}

// ResolveAlertWithContext resolves an alert with context
func (h *BaseAlertHandler) ResolveAlertWithContext(ctx context.Context, alertID string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	alert, exists := h.alerts[alertID]
	if !exists {
		return fmt.Errorf("alert %s not found", alertID)
	}
	
	now := time.Now()
	alert.Status = AlertStatusResolved
	alert.ResolvedAt = &now
	alert.EndsAt = &now
	
	// Update statistics
	h.updateAlertStats(alert)
	
	return nil
}

// SuppressAlert suppresses an alert
func (h *BaseAlertHandler) SuppressAlert(id string, silencedBy string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	alert, exists := h.alerts[id]
	if !exists {
		return fmt.Errorf("alert %s not found", id)
	}
	
	alert.Silenced = true
	alert.SilencedBy = silencedBy
	now := time.Now()
	alert.SilencedAt = &now
	
	return nil
}

// ReactivateAlert reactivates an alert
func (h *BaseAlertHandler) ReactivateAlert(id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	alert, exists := h.alerts[id]
	if !exists {
		return fmt.Errorf("alert %s not found", id)
	}
	
	alert.Silenced = false
	alert.SilencedBy = ""
	alert.SilencedAt = nil
	
	return nil
}

// AddAlertLabel adds a label to an alert
func (h *BaseAlertHandler) AddAlertLabel(id string, key, value string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	alert, exists := h.alerts[id]
	if !exists {
		return fmt.Errorf("alert %s not found", id)
	}
	
	if alert.Labels == nil {
		alert.Labels = make(map[string]string)
	}
	alert.Labels[key] = value
	
	return nil
}

// RemoveAlertLabel removes a label from an alert
func (h *BaseAlertHandler) RemoveAlertLabel(id string, key string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	alert, exists := h.alerts[id]
	if !exists {
		return fmt.Errorf("alert %s not found", id)
	}
	
	if alert.Labels != nil {
		delete(alert.Labels, key)
	}
	
	return nil
}

// GetAlertLabels returns alert labels
func (h *BaseAlertHandler) GetAlertLabels(id string) (map[string]string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	alert, exists := h.alerts[id]
	if !exists {
		return nil, fmt.Errorf("alert %s not found", id)
	}
	
	return alert.Labels, nil
}

// SetAlertContext sets alert context
func (h *BaseAlertHandler) SetAlertContext(id string, context context.Context) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	alert, exists := h.alerts[id]
	if !exists {
		return fmt.Errorf("alert %s not found", id)
	}
	
	alert.Context = context
	
	return nil
}

// GetAlertContext returns alert context
func (h *BaseAlertHandler) GetAlertContext(id string) (context.Context, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	alert, exists := h.alerts[id]
	if !exists {
		return nil, fmt.Errorf("alert %s not found", id)
	}
	
	return alert.Context, nil
}

// CleanupOldAlerts cleans up old alerts
func (h *BaseAlertHandler) CleanupOldAlerts(retentionPeriod time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	now := time.Now()
	validAlerts := make(map[string]*Alert)
	
	for id, alert := range h.alerts {
		if now.Sub(alert.StartsAt) <= retentionPeriod {
			validAlerts[id] = alert
		}
	}
	
	h.alerts = validAlerts
}

// GetAlertHistory returns alert history
func (h *BaseAlertHandler) GetAlertHistory(limit int) []*Alert {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	alerts := make([]*Alert, 0, len(h.alerts))
	for _, alert := range h.alerts {
		alerts = append(alerts, alert)
	}
	
	// Sort alerts by start time (newest first)
	// This is a simplified sort - in production you'd use a proper sorting algorithm
	if len(alerts) > limit {
		alerts = alerts[len(alerts)-limit:]
	}
	
	return alerts
}

// GetActiveAlerts returns active alerts
func (h *BaseAlertHandler) GetActiveAlerts() []*Alert {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var activeAlerts []*Alert
	for _, alert := range h.alerts {
		if alert.Status == AlertStatusFiring && !alert.Silenced {
			activeAlerts = append(activeAlerts, alert)
		}
	}
	
	return activeAlerts
}

// GetAlert returns a specific alert
func (h *BaseAlertHandler) GetAlert(alertID string) (*Alert, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	alert, exists := h.alerts[alertID]
	return alert, exists
}

// IsAlertActive returns whether an alert is active
func (h *BaseAlertHandler) IsAlertActive(alertID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	alert, exists := h.alerts[alertID]
	if !exists {
		return false
	}
	
	return alert.Status == AlertStatusFiring && !alert.Silenced
}

// IsAlertSilenced returns whether an alert is silenced
func (h *BaseAlertHandler) IsAlertSilenced(alertID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	alert, exists := h.alerts[alertID]
	if !exists {
		return false
	}
	
	return alert.Silenced
}

// GetAlertDuration returns the duration of an alert
func (h *BaseAlertHandler) GetAlertDuration(alertID string) (time.Duration, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	alert, exists := h.alerts[alertID]
	if !exists {
		return 0, fmt.Errorf("alert %s not found", alertID)
	}
	
	if alert.EndsAt != nil {
		return alert.EndsAt.Sub(alert.StartsAt), nil
	}
	
	return time.Since(alert.StartsAt), nil
}

// GetAlertCountByType returns alert count by type
func (h *BaseAlertHandler) GetAlertCountByType(alertType AlertType) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.alertStats.AlertsByType[alertType]
}

// GetAlertCountBySeverity returns alert count by severity
func (h *BaseAlertHandler) GetAlertCountBySeverity(severity AlertSeverity) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.alertStats.AlertsBySeverity[severity]
}

// GetAlertCountByStatus returns alert count by status
func (h *BaseAlertHandler) GetAlertCountByStatus(status AlertStatus) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.alertStats.AlertsByStatus[status]
}

// GetAverageAlertDuration returns the average duration of alerts
func (h *BaseAlertHandler) GetAverageAlertDuration() time.Duration {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.alertStats.AverageDuration
}

// GetAlertRate returns the alert rate (alerts per minute)
func (h *BaseAlertHandler) GetAlertRate() float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.alertStats.AlertRate
}

// updateAlertStats updates alert statistics
func (h *BaseAlertHandler) updateAlertStats(alert *Alert) {
	// Initialize maps if they are nil
	if h.alertStats.AlertsByType == nil {
		h.alertStats.AlertsByType = make(map[AlertType]int)
	}
	if h.alertStats.AlertsBySeverity == nil {
		h.alertStats.AlertsBySeverity = make(map[AlertSeverity]int)
	}
	if h.alertStats.AlertsByStatus == nil {
		h.alertStats.AlertsByStatus = make(map[AlertStatus]int)
	}
	
	// Update counts
	h.alertStats.TotalAlerts++
	h.alertStats.AlertsByType[alert.Type]++
	h.alertStats.AlertsBySeverity[alert.Severity]++
	h.alertStats.AlertsByStatus[alert.Status]++
	
	// Update active alerts count
	if alert.Status == AlertStatusFiring && !alert.Silenced {
		h.alertStats.ActiveAlerts++
	} else {
		h.alertStats.ActiveAlerts--
	}
	
	// Update last alert time
	h.alertStats.LastAlertTime = &alert.StartsAt
	h.alertStats.LastAlert = alert
	
	// Calculate alert rate
	h.alertStats.AlertRate = float64(h.alertStats.TotalAlerts) / time.Since(alert.StartsAt).Minutes()
}

// generateAlertID generates a unique alert ID
func generateAlertID() string {
	return fmt.Sprintf("alert_%d", time.Now().UnixNano())
}