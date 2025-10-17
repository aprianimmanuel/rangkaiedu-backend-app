package monitoring

import (
	"context"
	"fmt"
	"time"
)



// Global metrics functions

// RecordMetric records a metric globally
func RecordMetric(name string, value float64, labels map[string]string) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.metricsCollector == nil {
		return fmt.Errorf("metrics collector is not enabled")
	}
	
	// Create a simple metric structure
	metric := map[string]interface{}{
		"name":      name,
		"value":     value,
		"labels":    labels,
		"timestamp": time.Now().Unix(),
		"type":      "gauge", // Default type
	}
	
	// Record the metric using the metrics collector
	if recorder, ok := service.metricsCollector.(interface{ RecordMetric(map[string]interface{}) error }); ok {
		return recorder.RecordMetric(metric)
	}
	
	return fmt.Errorf("metrics collector does not support recording")
}

// IncrementCounter increments a counter metric globally
func IncrementCounter(name string, value float64, labels map[string]string) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.metricsCollector == nil {
		return fmt.Errorf("metrics collector is not enabled")
	}
	
	// Create a counter metric structure
	metric := map[string]interface{}{
		"name":      name,
		"value":     value,
		"labels":    labels,
		"timestamp": time.Now().Unix(),
		"type":      "counter",
	}
	
	// Record the counter using the metrics collector
	if recorder, ok := service.metricsCollector.(interface{ RecordMetric(map[string]interface{}) error }); ok {
		return recorder.RecordMetric(metric)
	}
	
	return fmt.Errorf("metrics collector does not support counter recording")
}

// RecordHistogram records a histogram metric globally
func RecordHistogram(name string, value float64, labels map[string]string) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.metricsCollector == nil {
		return fmt.Errorf("metrics collector is not enabled")
	}
	
	// Create a histogram metric structure
	metric := map[string]interface{}{
		"name":      name,
		"value":     value,
		"labels":    labels,
		"timestamp": time.Now().Unix(),
		"type":      "histogram",
	}
	
	// Record the histogram using the metrics collector
	if recorder, ok := service.metricsCollector.(interface{ RecordMetric(map[string]interface{}) error }); ok {
		return recorder.RecordMetric(metric)
	}
	
	return fmt.Errorf("metrics collector does not support histogram recording")
}

// RecordGauge records a gauge metric globally
func RecordGauge(name string, value float64, labels map[string]string) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.metricsCollector == nil {
		return fmt.Errorf("metrics collector is not enabled")
	}
	
	// Create a gauge metric structure
	metric := map[string]interface{}{
		"name":      name,
		"value":     value,
		"labels":    labels,
		"timestamp": time.Now().Unix(),
		"type":      "gauge",
	}
	
	// Record the gauge using the metrics collector
	if recorder, ok := service.metricsCollector.(interface{ RecordMetric(map[string]interface{}) error }); ok {
		return recorder.RecordMetric(metric)
	}
	
	return fmt.Errorf("metrics collector does not support gauge recording")
}

// RecordSummary records a summary metric globally
func RecordSummary(name string, value float64, labels map[string]string) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.metricsCollector == nil {
		return fmt.Errorf("metrics collector is not enabled")
	}
	
	// Create a summary metric structure
	metric := map[string]interface{}{
		"name":      name,
		"value":     value,
		"labels":    labels,
		"timestamp": time.Now().Unix(),
		"type":      "summary",
	}
	
	// Record the summary using the metrics collector
	if recorder, ok := service.metricsCollector.(interface{ RecordMetric(map[string]interface{}) error }); ok {
		return recorder.RecordMetric(metric)
	}
	
	return fmt.Errorf("metrics collector does not support summary recording")
}

// GetMetricsCount returns the number of metrics recorded globally
func GetMetricsCount() int {
	service := GetService()
	if service == nil {
		return 0
	}
	
	return service.GetMetricsCount()
}

// IsMetricsEnabled returns whether metrics collection is enabled globally
func IsMetricsEnabled() bool {
	service := GetService()
	if service == nil {
		return false
	}
	
	return service.IsMetricsEnabled()
}

// Global error functions

// RecordError records an error globally
func RecordError(err error, contextData map[string]interface{}) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.errorHandler == nil {
		return fmt.Errorf("error handler is not enabled")
	}
	
	// DEBUG: Log the missing ctx issue
	fmt.Printf("DEBUG: RecordError - ctx parameter is missing but being used\n")
	
	// Create a monitoring error structure
	errorData := &MonitoringError{
		ID:        fmt.Sprintf("error_%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		Severity:  ErrorSeverityError, // Default severity
		Type:      ErrorTypeError,     // Default type
		Message:   err.Error(),
		Context:   contextData,
	}
	
	// Handle the error using the error handler - using context.Background() as fallback
	return service.errorHandler.HandleError(context.Background(), errorData)
}

// RecordWarning records a warning globally
func RecordWarning(message string, contextData map[string]interface{}) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.errorHandler == nil {
		return fmt.Errorf("error handler is not enabled")
	}
	
	// DEBUG: Log the missing ctx issue
	fmt.Printf("DEBUG: RecordWarning - ctx parameter is missing but being used\n")
	
	// Create a monitoring error structure for warning
	errorData := &MonitoringError{
		ID:        fmt.Sprintf("warning_%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		Severity:  ErrorSeverityWarning,
		Type:      ErrorTypeError, // Default type
		Message:   message,
		Context:   contextData,
	}
	
	// Handle the warning using the error handler - using context.Background() as fallback
	return service.errorHandler.HandleError(context.Background(), errorData)
}

// GetErrorCount returns the number of errors recorded globally
func GetErrorCount() int {
	service := GetService()
	if service == nil {
		return 0
	}
	
	return service.GetErrorCount()
}

// IsErrorsEnabled returns whether error handling is enabled globally
func IsErrorsEnabled() bool {
	service := GetService()
	if service == nil {
		return false
	}
	
	return service.IsErrorsEnabled()
}

// Global alert functions

// CreateAlert creates an alert globally
func CreateAlert(name, message string, severity AlertSeverity, contextData map[string]interface{}) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.alertHandler == nil {
		return fmt.Errorf("alert handler is not enabled")
	}
	
	// DEBUG: Log the missing ctx issue
	fmt.Printf("DEBUG: CreateAlert - ctx parameter is missing but being used\n")
	
	// Create an alert structure
	alert := &Alert{
		ID:        fmt.Sprintf("alert_%d", time.Now().UnixNano()),
		Title:     name,
		Message:   message,
		Severity:  severity,
		StartsAt:  time.Now(),
		Status:    AlertStatusFiring,
	}
	
	// Handle the alert using the alert handler - using context.Background() as fallback
	return service.alertHandler.HandleAlert(context.Background(), alert)
}

// ResolveAlert resolves an alert globally
func ResolveAlert(alertID string) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.alertHandler == nil {
		return fmt.Errorf("alert handler is not enabled")
	}
	
	// Resolve the alert using the alert handler
	return service.alertHandler.ResolveAlert(alertID, "global")
}

// GetAlertCount returns the number of alerts globally
func GetAlertCount() int {
	service := GetService()
	if service == nil {
		return 0
	}
	
	return service.GetAlertCount()
}

// IsAlertsEnabled returns whether alerting is enabled globally
func IsAlertsEnabled() bool {
	service := GetService()
	if service == nil {
		return false
	}
	
	return service.IsAlertsEnabled()
}

// Global health check functions

// RegisterHealthCheck registers a health check globally
func RegisterHealthCheck(name string, checkFunc func() HealthStatus) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.healthChecker == nil {
		return fmt.Errorf("health checker is not enabled")
	}
	
	// Create a simple health check that implements the HealthCheck interface
	healthCheck := &simpleHealthCheck{
		id:   name,
		name: name,
		fn:   checkFunc,
	}
	
	// Register the health check using the health checker
	return service.healthChecker.RegisterHealthCheck(healthCheck)
}

// simpleHealthCheck implements the HealthCheck interface for global health checks
type simpleHealthCheck struct {
	id      string
	name    string
	fn      func() HealthStatus
	enabled bool
	result  *HealthCheckResult
}

func (s *simpleHealthCheck) Check(ctx context.Context) (*HealthCheckResult, error) {
	status := s.fn()
	result := &HealthCheckResult{
		ID:        s.id,
		Name:      s.name,
		Type:      HealthCheckTypeCustom,
		Status:    status,
		Message:   fmt.Sprintf("Health check %s completed", s.name),
		Timestamp: time.Now(),
	}
	s.result = result
	return result, nil
}

func (s *simpleHealthCheck) GetName() string {
	return s.name
}

func (s *simpleHealthCheck) GetType() HealthCheckType {
	return HealthCheckTypeCustom
}

func (s *simpleHealthCheck) GetInterval() time.Duration {
	return 30 * time.Second
}

func (s *simpleHealthCheck) GetTimeout() time.Duration {
	return 10 * time.Second
}

func (s *simpleHealthCheck) GetRetries() int {
	return 3
}

func (s *simpleHealthCheck) GetEnabled() bool {
	return s.enabled
}

func (s *simpleHealthCheck) SetEnabled(enabled bool) {
	s.enabled = enabled
}

func (s *simpleHealthCheck) GetLastResult() *HealthCheckResult {
	return s.result
}

func (s *simpleHealthCheck) GetHistory(limit int) []*HealthCheckResult {
	// Simple implementation that returns the last result if available
	if s.result != nil {
		return []*HealthCheckResult{s.result}
	}
	return []*HealthCheckResult{}
}

// RunHealthCheck runs a health check globally
func RunHealthCheck(name string) (HealthStatus, error) {
	service := GetService()
	if service == nil {
		return HealthStatusUnknown, fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.healthChecker == nil {
		return HealthStatusUnknown, fmt.Errorf("health checker is not enabled")
	}
	
	// Run the health check using the health checker
	result, err := service.healthChecker.RunHealthCheck(context.Background(), name)
	if err != nil {
		return HealthStatusUnknown, err
	}
	
	return result.Status, nil
}

// GetHealthCheckCount returns the number of health checks globally
func GetHealthCheckCount() int {
	service := GetService()
	if service == nil {
		return 0
	}
	
	return service.GetHealthCheckCount()
}

// IsHealthChecksEnabled returns whether health checks are enabled globally
func IsHealthChecksEnabled() bool {
	service := GetService()
	if service == nil {
		return false
	}
	
	return service.IsHealthChecksEnabled()
}

// GetHealthStatus returns the overall health status globally
func GetHealthStatus() HealthStatus {
	service := GetService()
	if service == nil {
		return HealthStatusUnhealthy
	}
	
	return service.GetHealthStatus()
}

// GetHealthMessage returns the health message globally
func GetHealthMessage() string {
	service := GetService()
	if service == nil {
		return "Monitoring service is not initialized"
	}
	
	return service.GetHealthMessage()
}