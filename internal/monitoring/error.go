package monitoring

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Error represents a monitoring error
type Error struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Severity    ErrorSeverity          `json:"severity"`
	Type        ErrorType              `json:"type"`
	Status      ErrorStatus            `json:"status"`
	Message     string                 `json:"message"`
	Stack       string                 `json:"stack,omitempty"`
	Context     map[string]interface{} `json:"context,omitempty"`
	Labels      map[string]string      `json:"labels,omitempty"`
	Retries     int                    `json:"retries,omitempty"`
	MaxRetries  int                    `json:"max_retries,omitempty"`
	NextRetry   *time.Time             `json:"next_retry,omitempty"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty"`
}

// ErrorStats represents error statistics
type ErrorStats struct {
	TotalErrors      int                    `json:"total_errors"`
	CriticalErrors   int                    `json:"critical_errors"`
	ErrorErrors      int                    `json:"error_errors"`
	WarningErrors    int                    `json:"warning_errors"`
	InfoErrors       int                    `json:"info_errors"`
	DebugErrors      int                    `json:"debug_errors"`
	OpenErrors       int                    `json:"open_errors"`
	ResolvedErrors   int                    `json:"resolved_errors"`
	SuppressedErrors int                    `json:"suppressed_errors"`
	PendingErrors    int                    `json:"pending_errors"`
	ErrorsByType     map[ErrorType]int      `json:"errors_by_type"`
	ErrorsBySeverity map[ErrorSeverity]int  `json:"errors_by_severity"`
	LastErrorTime    *time.Time             `json:"last_error_time"`
	LastError        *Error                 `json:"last_error,omitempty"`
	AverageDuration  time.Duration          `json:"average_duration"`
}

// ErrorHandler defines the interface for error handling
type ErrorHandler interface {
	// HandleError handles a monitoring error
	HandleError(ctx context.Context, err *MonitoringError) error
	
	// GetErrorStats returns error statistics
	GetErrorStats() ErrorStats
	
	// GetErrorCount returns the number of errors
	GetErrorCount() int
	
	// IsEnabled returns whether error handling is enabled
	IsEnabled() bool
	
	// GetErrors returns all errors
	GetErrors(limit int, offset int) ([]*Error, error)
	
	// GetErrorsBySeverity returns errors by severity
	GetErrorsBySeverity(severity ErrorSeverity, limit int, offset int) ([]*Error, error)
	
	// GetErrorsByType returns errors by type
	GetErrorsByType(errorType ErrorType, limit int, offset int) ([]*Error, error)
	
	// GetErrorsByStatus returns errors by status
	GetErrorsByStatus(status ErrorStatus, limit int, offset int) ([]*Error, error)
	
	// GetErrorsByTimeRange returns errors by time range
	GetErrorsByTimeRange(startTime, endTime time.Time, limit int, offset int) ([]*Error, error)
	
	// ResolveError resolves an error
	ResolveError(id string) error
	
	// ResolveErrorWithContext resolves an error with context
	ResolveErrorWithContext(ctx context.Context, errorID string) error
	
	// SuppressError suppresses an error
	SuppressError(id string) error
	
	// ReactivateError reactivates an error
	ReactivateError(id string) error
	
	// RetryError retries an error
	RetryError(id string) error
	
	// AddErrorLabel adds a label to an error
	AddErrorLabel(id string, key, value string) error
	
	// RemoveErrorLabel removes a label from an error
	RemoveErrorLabel(id string, key string) error
	
	// GetErrorLabels returns error labels
	GetErrorLabels(id string) (map[string]string, error)
	
	// SetErrorContext sets error context
	SetErrorContext(id string, context map[string]interface{}) error
	
	// GetErrorContext returns error context
	GetErrorContext(id string) (map[string]interface{}, error)
	
	// CleanupExpiredErrors cleans up expired errors
	CleanupExpiredErrors(retentionPeriod time.Duration)
	
	// GetErrorHistory returns error history
	GetErrorHistory(limit int) []*Error
	
	// GetActiveErrors returns active errors
	GetActiveErrors() []*Error
	
	// GetError returns a specific error
	GetError(errorID string) (*Error, bool)
	
	// IsErrorActive returns whether an error is active
	IsErrorActive(errorID string) bool
	
	// IsErrorSuppressed returns whether an error is suppressed
	IsErrorSuppressed(errorID string) bool
	
	// GetErrorDuration returns the duration of an error
	GetErrorDuration(errorID string) (time.Duration, error)
	
	// GetErrorCountByType returns error count by type
	GetErrorCountByType(errorType ErrorType) int
	
	// GetErrorCountBySeverity returns error count by severity
	GetErrorCountBySeverity(severity ErrorSeverity) int
	
	// GetErrorCountByStatus returns error count by status
	GetErrorCountByStatus(status ErrorStatus) int
	
	// GetAverageErrorDuration returns the average duration of errors
	GetAverageErrorDuration() time.Duration
	
	// GetErrorRate returns the error rate (errors per minute)
	GetErrorRate() float64
}

// MonitoringError represents a monitoring error
type MonitoringError struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Severity    ErrorSeverity          `json:"severity"`
	Type        ErrorType              `json:"type"`
	Message     string                 `json:"message"`
	Stack       string                 `json:"stack,omitempty"`
	Context     map[string]interface{} `json:"context,omitempty"`
	Labels      map[string]string      `json:"labels,omitempty"`
	Retries     int                    `json:"retries,omitempty"`
	MaxRetries  int                    `json:"max_retries,omitempty"`
	NextRetry   *time.Time             `json:"next_retry,omitempty"`
}

// ErrorHandlerConfig represents the configuration for error handling
type ErrorHandlerConfig struct {
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
	DefaultSeverity           ErrorSeverity `json:"default_severity"`
	DefaultTimeout            time.Duration `json:"default_timeout"`
	AutoResolveTimeout        time.Duration `json:"auto_resolve_timeout"`
	EnableErrorGroups         bool          `json:"enable_error_groups"`
	MaxErrorsPerMinute        int           `json:"max_errors_per_minute"`
	ErrorRateLimitTotal       int           `json:"error_rate_limit_total"`
}

// DefaultErrorHandlerConfig returns the default configuration for error handling
func DefaultErrorHandlerConfig() *ErrorHandlerConfig {
	return &ErrorHandlerConfig{
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
		DefaultSeverity:           ErrorSeverityError,
		DefaultTimeout:            5 * time.Minute,
		AutoResolveTimeout:        1 * time.Hour,
		EnableErrorGroups:         true,
		MaxErrorsPerMinute:        100,
		ErrorRateLimitTotal:       1000,
	}
}

// BaseErrorHandler provides the base implementation for error handlers
type BaseErrorHandler struct {
	mu            sync.RWMutex
	errors        map[string]*Error
	errorStats    ErrorStats
	config        *ErrorHandlerConfig
	running       bool
	stopChan      chan struct{}
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(maxHistorySize int) ErrorHandler {
	return &BaseErrorHandler{
		errors:      make(map[string]*Error),
		errorStats:  ErrorStats{},
		running:     false,
		stopChan:    make(chan struct{}),
		config: &ErrorHandlerConfig{
			MaxHistorySize: maxHistorySize,
		},
	}
}

// HandleError handles a monitoring error
func (h *BaseErrorHandler) HandleError(ctx context.Context, err *MonitoringError) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	// Create error from monitoring error
	error := &Error{
		ID:        fmt.Sprintf("error_%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		Severity:  err.Severity,
		Type:      err.Type,
		Status:    ErrorStatusOpen,
		Message:   err.Message,
		Stack:     err.Stack,
		Context:   err.Context,
		Labels:    err.Labels,
		Retries:   err.Retries,
		MaxRetries: err.MaxRetries,
		NextRetry: err.NextRetry,
	}
	
	// Add to errors
	h.errors[error.ID] = error
	
	// Trim errors if we exceed max history size
	if len(h.errors) > h.config.MaxHistorySize {
		// Remove oldest errors
		keys := make([]string, 0, len(h.errors))
		for k := range h.errors {
			keys = append(keys, k)
		}
		// Simple approach: remove first half
		for i := 0; i < len(keys)/2; i++ {
			delete(h.errors, keys[i])
		}
	}
	
	// Update statistics
	h.updateErrorStats(error)
	
	return nil
}

// GetErrorStats returns error statistics
func (h *BaseErrorHandler) GetErrorStats() ErrorStats {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.errorStats
}

// GetErrorCount returns the number of errors
func (h *BaseErrorHandler) GetErrorCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return len(h.errors)
}

// IsEnabled returns whether error handling is enabled
func (h *BaseErrorHandler) IsEnabled() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.config.Enabled
}

// GetErrors returns all errors
func (h *BaseErrorHandler) GetErrors(limit int, offset int) ([]*Error, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	errors := make([]*Error, 0, len(h.errors))
	for _, err := range h.errors {
		errors = append(errors, err)
	}
	
	if limit <= 0 {
		limit = len(errors)
	}
	if offset < 0 {
		offset = 0
	}
	
	if offset >= len(errors) {
		return []*Error{}, nil
	}
	
	end := offset + limit
	if end > len(errors) {
		end = len(errors)
	}
	
	return errors[offset:end], nil
}

// GetErrorsBySeverity returns errors by severity
func (h *BaseErrorHandler) GetErrorsBySeverity(severity ErrorSeverity, limit int, offset int) ([]*Error, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var errors []*Error
	for _, err := range h.errors {
		if err.Severity == severity {
			errors = append(errors, err)
		}
	}
	
	if limit <= 0 {
		limit = len(errors)
	}
	if offset < 0 {
		offset = 0
	}
	
	if offset >= len(errors) {
		return []*Error{}, nil
	}
	
	end := offset + limit
	if end > len(errors) {
		end = len(errors)
	}
	
	return errors[offset:end], nil
}

// GetErrorsByType returns errors by type
func (h *BaseErrorHandler) GetErrorsByType(errorType ErrorType, limit int, offset int) ([]*Error, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var errors []*Error
	for _, err := range h.errors {
		if err.Type == errorType {
			errors = append(errors, err)
		}
	}
	
	if limit <= 0 {
		limit = len(errors)
	}
	if offset < 0 {
		offset = 0
	}
	
	if offset >= len(errors) {
		return []*Error{}, nil
	}
	
	end := offset + limit
	if end > len(errors) {
		end = len(errors)
	}
	
	return errors[offset:end], nil
}

// GetErrorsByStatus returns errors by status
func (h *BaseErrorHandler) GetErrorsByStatus(status ErrorStatus, limit int, offset int) ([]*Error, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var errors []*Error
	for _, err := range h.errors {
		if err.Status == status {
			errors = append(errors, err)
		}
	}
	
	if limit <= 0 {
		limit = len(errors)
	}
	if offset < 0 {
		offset = 0
	}
	
	if offset >= len(errors) {
		return []*Error{}, nil
	}
	
	end := offset + limit
	if end > len(errors) {
		end = len(errors)
	}
	
	return errors[offset:end], nil
}

// GetErrorsByTimeRange returns errors by time range
func (h *BaseErrorHandler) GetErrorsByTimeRange(startTime, endTime time.Time, limit int, offset int) ([]*Error, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var errors []*Error
	for _, err := range h.errors {
		if err.Timestamp.After(startTime) && err.Timestamp.Before(endTime) {
			errors = append(errors, err)
		}
	}
	
	if limit <= 0 {
		limit = len(errors)
	}
	if offset < 0 {
		offset = 0
	}
	
	if offset >= len(errors) {
		return []*Error{}, nil
	}
	
	end := offset + limit
	if end > len(errors) {
		end = len(errors)
	}
	
	return errors[offset:end], nil
}

// ResolveError resolves an error
func (h *BaseErrorHandler) ResolveError(id string) error {
	return h.ResolveErrorWithContext(context.Background(), id)
}

// ResolveErrorWithContext resolves an error with context
func (h *BaseErrorHandler) ResolveErrorWithContext(ctx context.Context, errorID string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	error, exists := h.errors[errorID]
	if !exists {
		return fmt.Errorf("error %s not found", errorID)
	}
	
	now := time.Now()
	error.Status = ErrorStatusResolved
	error.ResolvedAt = &now
	
	// Update statistics
	h.updateErrorStats(error)
	
	return nil
}

// SuppressError suppresses an error
func (h *BaseErrorHandler) SuppressError(id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	error, exists := h.errors[id]
	if !exists {
		return fmt.Errorf("error %s not found", id)
	}
	
	error.Status = ErrorStatusSuppressed
	
	// Update statistics
	h.updateErrorStats(error)
	
	return nil
}

// ReactivateError reactivates an error
func (h *BaseErrorHandler) ReactivateError(id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	error, exists := h.errors[id]
	if !exists {
		return fmt.Errorf("error %s not found", id)
	}
	
	error.Status = ErrorStatusOpen
	
	// Update statistics
	h.updateErrorStats(error)
	
	return nil
}

// RetryError retries an error
func (h *BaseErrorHandler) RetryError(id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	error, exists := h.errors[id]
	if !exists {
		return fmt.Errorf("error %s not found", id)
	}
	
	if error.Retries < error.MaxRetries {
		error.Retries++
		nextRetry := time.Now().Add(h.config.RetryDelay)
		error.NextRetry = &nextRetry
		error.Status = ErrorStatusPending
		h.updateErrorStats(error)
		return nil
	}
	return fmt.Errorf("error %s has reached max retries", id)
}

// AddErrorLabel adds a label to an error
func (h *BaseErrorHandler) AddErrorLabel(id string, key, value string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	error, exists := h.errors[id]
	if !exists {
		return fmt.Errorf("error %s not found", id)
	}
	
	if error.Labels == nil {
		error.Labels = make(map[string]string)
	}
	error.Labels[key] = value
	
	return nil
}

// RemoveErrorLabel removes a label from an error
func (h *BaseErrorHandler) RemoveErrorLabel(id string, key string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	error, exists := h.errors[id]
	if !exists {
		return fmt.Errorf("error %s not found", id)
	}
	
	if error.Labels != nil {
		delete(error.Labels, key)
	}
	
	return nil
}

// GetErrorLabels returns error labels
func (h *BaseErrorHandler) GetErrorLabels(id string) (map[string]string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	error, exists := h.errors[id]
	if !exists {
		return nil, fmt.Errorf("error %s not found", id)
	}
	
	return error.Labels, nil
}

// SetErrorContext sets error context
func (h *BaseErrorHandler) SetErrorContext(id string, context map[string]interface{}) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	error, exists := h.errors[id]
	if !exists {
		return fmt.Errorf("error %s not found", id)
	}
	
	error.Context = context
	
	return nil
}

// GetErrorContext returns error context
func (h *BaseErrorHandler) GetErrorContext(id string) (map[string]interface{}, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	error, exists := h.errors[id]
	if !exists {
		return nil, fmt.Errorf("error %s not found", id)
	}
	
	return error.Context, nil
}

// CleanupExpiredErrors cleans up expired errors
func (h *BaseErrorHandler) CleanupExpiredErrors(retentionPeriod time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	// Remove resolved errors older than retention period
	cutoff := time.Now().Add(-retentionPeriod)
	validErrors := make(map[string]*Error)
	
	for id, error := range h.errors {
		if error.Status != ErrorStatusResolved || error.ResolvedAt == nil || error.ResolvedAt.After(cutoff) {
			validErrors[id] = error
		}
	}
	
	h.errors = validErrors
	h.updateAllErrorStats()
}

// GetErrorHistory returns error history
func (h *BaseErrorHandler) GetErrorHistory(limit int) []*Error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	errors := make([]*Error, 0, len(h.errors))
	for _, err := range h.errors {
		errors = append(errors, err)
	}
	
	// Sort errors by timestamp (newest first)
	// This is a simplified sort - in production you'd use a proper sorting algorithm
	if len(errors) > limit {
		errors = errors[len(errors)-limit:]
	}
	
	return errors
}

// GetActiveErrors returns active errors
func (h *BaseErrorHandler) GetActiveErrors() []*Error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var activeErrors []*Error
	for _, err := range h.errors {
		if err.Status == ErrorStatusOpen || err.Status == ErrorStatusPending {
			activeErrors = append(activeErrors, err)
		}
	}
	
	return activeErrors
}

// GetError returns a specific error
func (h *BaseErrorHandler) GetError(errorID string) (*Error, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	error, exists := h.errors[errorID]
	return error, exists
}

// IsErrorActive returns whether an error is active
func (h *BaseErrorHandler) IsErrorActive(errorID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	error, exists := h.errors[errorID]
	if !exists {
		return false
	}
	
	return error.Status == ErrorStatusOpen || error.Status == ErrorStatusPending
}

// IsErrorSuppressed returns whether an error is suppressed
func (h *BaseErrorHandler) IsErrorSuppressed(errorID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	error, exists := h.errors[errorID]
	if !exists {
		return false
	}
	
	return error.Status == ErrorStatusSuppressed
}

// GetErrorDuration returns the duration of an error
func (h *BaseErrorHandler) GetErrorDuration(errorID string) (time.Duration, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	error, exists := h.errors[errorID]
	if !exists {
		return 0, fmt.Errorf("error %s not found", errorID)
	}
	
	if error.ResolvedAt != nil {
		return error.ResolvedAt.Sub(error.Timestamp), nil
	}
	
	return time.Since(error.Timestamp), nil
}

// GetErrorCountByType returns error count by type
func (h *BaseErrorHandler) GetErrorCountByType(errorType ErrorType) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.errorStats.ErrorsByType[errorType]
}

// GetErrorCountBySeverity returns error count by severity
func (h *BaseErrorHandler) GetErrorCountBySeverity(severity ErrorSeverity) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.errorStats.ErrorsBySeverity[severity]
}

// GetErrorCountByStatus returns error count by status
func (h *BaseErrorHandler) GetErrorCountByStatus(status ErrorStatus) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	switch status {
	case ErrorStatusOpen:
		return h.errorStats.OpenErrors
	case ErrorStatusResolved:
		return h.errorStats.ResolvedErrors
	case ErrorStatusSuppressed:
		return h.errorStats.SuppressedErrors
	case ErrorStatusPending:
		return h.errorStats.PendingErrors
	default:
		return 0
	}
}

// GetAverageErrorDuration returns the average duration of errors
func (h *BaseErrorHandler) GetAverageErrorDuration() time.Duration {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.errorStats.AverageDuration
}

// GetErrorRate returns the error rate (errors per minute)
func (h *BaseErrorHandler) GetErrorRate() float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	if h.errorStats.LastError == nil {
		return 0
	}
	
	return float64(h.errorStats.TotalErrors) / time.Since(h.errorStats.LastError.Timestamp).Minutes()
}

// updateErrorStats updates error statistics
func (h *BaseErrorHandler) updateErrorStats(err *Error) {
	// Initialize maps if they are nil
	if h.errorStats.ErrorsByType == nil {
		h.errorStats.ErrorsByType = make(map[ErrorType]int)
	}
	if h.errorStats.ErrorsBySeverity == nil {
		h.errorStats.ErrorsBySeverity = make(map[ErrorSeverity]int)
	}
	
	// Update counts
	h.errorStats.TotalErrors++
	h.errorStats.ErrorsByType[err.Type]++
	h.errorStats.ErrorsBySeverity[err.Severity]++
	
	// Update status counts
	switch err.Status {
	case ErrorStatusOpen:
		h.errorStats.OpenErrors++
	case ErrorStatusResolved:
		h.errorStats.ResolvedErrors++
	case ErrorStatusSuppressed:
		h.errorStats.SuppressedErrors++
	case ErrorStatusPending:
		h.errorStats.PendingErrors++
	}
	
	// Update severity counts
	switch err.Severity {
	case ErrorSeverityCritical:
		h.errorStats.CriticalErrors++
	case ErrorSeverityError:
		h.errorStats.ErrorErrors++
	case ErrorSeverityWarning:
		h.errorStats.WarningErrors++
	case ErrorSeverityInfo:
		h.errorStats.InfoErrors++
	case ErrorSeverityDebug:
		h.errorStats.DebugErrors++
	}
	
	// Update last error time
	h.errorStats.LastErrorTime = &err.Timestamp
	h.errorStats.LastError = err
}

// updateAllErrorStats updates all error statistics
func (h *BaseErrorHandler) updateAllErrorStats() {
	// Reset stats
	h.errorStats = ErrorStats{
		ErrorsByType:     make(map[ErrorType]int),
		ErrorsBySeverity: make(map[ErrorSeverity]int),
	}
	
	// Recalculate stats
	for _, err := range h.errors {
		h.updateErrorStats(err)
	}
}