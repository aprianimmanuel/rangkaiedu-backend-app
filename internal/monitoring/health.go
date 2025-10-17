package monitoring

import (
	"context"
	"fmt"
	"sync"
	"time"
)


// HealthCheckResult represents the result of a health check
type HealthCheckResult struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        HealthCheckType        `json:"type"`
	Status      HealthStatus           `json:"status"`
	Message     string                 `json:"message"`
	Details     map[string]interface{} `json:"details"`
	Timestamp   time.Time              `json:"timestamp"`
	Duration    time.Duration          `json:"duration"`
	Context     context.Context        `json:"-"`
	Checks      []*HealthCheckResult   `json:"checks,omitempty"`
}

// HealthCheck defines the interface for health checks
type HealthCheck interface {
	// Check performs the health check
	Check(ctx context.Context) (*HealthCheckResult, error)
	
	// GetName returns the health check name
	GetName() string
	
	// GetType returns the health check type
	GetType() HealthCheckType
	
	// GetInterval returns the health check interval
	GetInterval() time.Duration
	
	// GetTimeout returns the health check timeout
	GetTimeout() time.Duration
	
	// GetRetries returns the number of retries
	GetRetries() int
	
	// GetEnabled returns whether the health check is enabled
	GetEnabled() bool
	
	// SetEnabled sets whether the health check is enabled
	SetEnabled(enabled bool)
	
	// GetLastResult returns the last health check result
	GetLastResult() *HealthCheckResult
	
	// GetHistory returns the health check history
	GetHistory(limit int) []*HealthCheckResult
}

// HealthCheckConfig represents the configuration for a health check
type HealthCheckConfig struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        HealthCheckType   `json:"type"`
	Enabled     bool              `json:"enabled"`
	Interval    time.Duration     `json:"interval"`
	Timeout     time.Duration     `json:"timeout"`
	Retries     int               `json:"retries"`
	Threshold   HealthStatus      `json:"threshold"`
	Labels      map[string]string `json:"labels"`
	Config      map[string]interface{} `json:"config"`
}

// HealthChecker defines the interface for health check management
type HealthChecker interface {
	// RegisterHealthCheck registers a health check
	RegisterHealthCheck(check HealthCheck) error
	
	// UnregisterHealthCheck unregisters a health check
	UnregisterHealthCheck(id string) error
	
	// GetHealthCheck returns a health check
	GetHealthCheck(id string) (HealthCheck, bool)
	
	// GetHealthChecks returns all health checks
	GetHealthChecks() []HealthCheck
	
	// GetHealthChecksByType returns health checks by type
	GetHealthChecksByType(checkType HealthCheckType) []HealthCheck
	
	// GetHealthChecksByLabel returns health checks by label
	GetHealthChecksByLabel(key, value string) []HealthCheck
	
	// RunHealthCheck runs a specific health check
	RunHealthCheck(ctx context.Context, id string) (*HealthCheckResult, error)
	
	// RunAllHealthChecks runs all health checks
	RunAllHealthChecks(ctx context.Context) (*HealthCheckResult, error)
	
	// GetSystemHealth returns the overall system health
	GetSystemHealth(ctx context.Context) (*HealthCheckResult, error)
	
	// GetHealthHistory returns health check history
	GetHealthHistory(limit int) []*HealthCheckResult
	
	// GetHealthStats returns health check statistics
	GetHealthStats() HealthStats
	
	// GetHealthChecksCount returns the number of health checks
	GetHealthChecksCount() int
	
	// Start starts the health checker
	Start(ctx context.Context) error
	
	// Stop stops the health checker
	Stop() error
	
	// IsRunning returns whether the health checker is running
	IsRunning() bool
}

// HealthStats represents health check statistics
type HealthStats struct {
	TotalChecks      int                    `json:"total_checks"`
	HealthyChecks    int                    `json:"healthy_checks"`
	DegradedChecks   int                    `json:"degraded_checks"`
	UnhealthyChecks  int                    `json:"unhealthy_checks"`
	UnknownChecks    int                    `json:"unknown_checks"`
	ChecksByType     map[HealthCheckType]int `json:"checks_by_type"`
	LastCheckTime    *time.Time             `json:"last_check_time"`
	AverageDuration  time.Duration          `json:"average_duration"`
	SuccessRate      float64                `json:"success_rate"`
}

// BaseHealthChecker provides the base implementation for health checkers
type BaseHealthChecker struct {
	mu            sync.RWMutex
	healthChecks  map[string]HealthCheck
	history       []*HealthCheckResult
	stats         HealthStats
	running       bool
	stopChan      chan struct{}
	config        *HealthCheckConfig
}
// BaseHealthCheck provides the base implementation for individual health checks
type BaseHealthCheck struct {
	config *HealthCheckConfig
	result *HealthCheckResult
}

// NewBaseHealthCheck creates a new base health check
func NewBaseHealthCheck(config *HealthCheckConfig) *BaseHealthCheck {
	return &BaseHealthCheck{
		config: config,
		result: &HealthCheckResult{
			ID:        config.ID,
			Name:      config.Name,
			Type:      config.Type,
			Status:    HealthStatusUnknown,
			Message:   "Health check not yet performed",
			Timestamp: time.Now(),
			Details:   make(map[string]interface{}),
		},
	}
}

// Check performs the health check
func (h *BaseHealthCheck) Check(ctx context.Context) (*HealthCheckResult, error) {
	startTime := time.Now()
	
	// Perform basic health check
	h.result.Timestamp = startTime
	h.result.Duration = time.Since(startTime)
	
	// For now, return a healthy status
	h.result.Status = HealthStatusHealthy
	h.result.Message = "Health check passed"
	
	return h.result, nil
}

// GetName returns the health check name
func (h *BaseHealthCheck) GetName() string {
	return h.config.Name
}

// GetType returns the health check type
func (h *BaseHealthCheck) GetType() HealthCheckType {
	return h.config.Type
}

// GetInterval returns the health check interval
func (h *BaseHealthCheck) GetInterval() time.Duration {
	return h.config.Interval
}

// GetTimeout returns the health check timeout
func (h *BaseHealthCheck) GetTimeout() time.Duration {
	return h.config.Timeout
}

// GetRetries returns the number of retries
func (h *BaseHealthCheck) GetRetries() int {
	return h.config.Retries
}

// GetEnabled returns whether the health check is enabled
func (h *BaseHealthCheck) GetEnabled() bool {
	return h.config.Enabled
}

// SetEnabled sets whether the health check is enabled
func (h *BaseHealthCheck) SetEnabled(enabled bool) {
	h.config.Enabled = enabled
}

// GetLastResult returns the last health check result
func (h *BaseHealthCheck) GetLastResult() *HealthCheckResult {
	return h.result
}

// GetHistory returns the health check history
func (h *BaseHealthCheck) GetHistory(limit int) []*HealthCheckResult {
	// Base health check doesn't maintain its own history
	// Return empty slice for now
	return []*HealthCheckResult{}
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(config *HealthCheckConfig) HealthChecker {
	return &BaseHealthChecker{
		healthChecks: make(map[string]HealthCheck),
		history:      make([]*HealthCheckResult, 0),
		stats:        HealthStats{},
		running:      false,
		stopChan:     make(chan struct{}),
		config:       config,
	}
}

// RegisterHealthCheck registers a health check
func (h *BaseHealthChecker) RegisterHealthCheck(check HealthCheck) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if _, exists := h.healthChecks[check.GetName()]; exists {
		return fmt.Errorf("health check %s already exists", check.GetName())
	}
	
	h.healthChecks[check.GetName()] = check
	
	return nil
}

// UnregisterHealthCheck unregisters a health check
func (h *BaseHealthChecker) UnregisterHealthCheck(id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if _, exists := h.healthChecks[id]; !exists {
		return fmt.Errorf("health check %s not found", id)
	}
	
	delete(h.healthChecks, id)
	
	return nil
}

// GetHealthCheck returns a health check
func (h *BaseHealthChecker) GetHealthCheck(id string) (HealthCheck, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	check, exists := h.healthChecks[id]
	return check, exists
}

// GetHealthChecks returns all health checks
func (h *BaseHealthChecker) GetHealthChecks() []HealthCheck {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	checks := make([]HealthCheck, 0, len(h.healthChecks))
	for _, check := range h.healthChecks {
		checks = append(checks, check)
	}
	
	return checks
}

// GetHealthChecksByType returns health checks by type
func (h *BaseHealthChecker) GetHealthChecksByType(checkType HealthCheckType) []HealthCheck {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var checks []HealthCheck
	for _, check := range h.healthChecks {
		if check.GetType() == checkType {
			checks = append(checks, check)
		}
	}
	
	return checks
}

// GetHealthChecksByLabel returns health checks by label
func (h *BaseHealthChecker) GetHealthChecksByLabel(key, value string) []HealthCheck {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var checks []HealthCheck
	for _, check := range h.healthChecks {
		// This would require the HealthCheck interface to have a GetLabels method
		// For now, return all checks
		checks = append(checks, check)
	}
	
	return checks
}

// RunHealthCheck runs a specific health check
func (h *BaseHealthChecker) RunHealthCheck(ctx context.Context, id string) (*HealthCheckResult, error) {
	h.mu.RLock()
	check, exists := h.healthChecks[id]
	h.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("health check %s not found", id)
	}
	
	if !check.GetEnabled() {
		return &HealthCheckResult{
			ID:        id,
			Name:      check.GetName(),
			Type:      check.GetType(),
			Status:    HealthStatusUnknown,
			Message:   "health check is disabled",
			Timestamp: time.Now(),
		}, nil
	}
	
	result, err := check.Check(ctx)
	if err != nil {
		result = &HealthCheckResult{
			ID:        id,
			Name:      check.GetName(),
			Type:      check.GetType(),
			Status:    HealthStatusUnhealthy,
			Message:   fmt.Sprintf("health check failed: %v", err),
			Timestamp: time.Now(),
		}
	}
	
	// Add to history
	h.addToHistory(result)
	
	// Update statistics
	h.updateHealthStats(result)
	
	return result, nil
}

// RunAllHealthChecks runs all health checks
func (h *BaseHealthChecker) RunAllHealthChecks(ctx context.Context) (*HealthCheckResult, error) {
	startTime := time.Now()
	
	// Create parent result
	parentResult := &HealthCheckResult{
		ID:        "system",
		Name:      "System Health",
		Type:      HealthCheckTypeSystem,
		Status:    HealthStatusHealthy,
		Message:   "All health checks completed",
		Timestamp: startTime,
		Checks:    make([]*HealthCheckResult, 0),
	}
	
	// Run all health checks
	for _, check := range h.GetHealthChecks() {
		if !check.GetEnabled() {
			continue
		}
		
		result, err := h.RunHealthCheck(ctx, check.GetName())
		if err != nil {
			// Log error but continue with other checks
			continue
		}
		
		parentResult.Checks = append(parentResult.Checks, result)
		
		// Update parent status based on child results
		if result.Status == HealthStatusUnhealthy {
			parentResult.Status = HealthStatusUnhealthy
			parentResult.Message = "One or more health checks failed"
		} else if result.Status == HealthStatusDegraded && parentResult.Status != HealthStatusUnhealthy {
			parentResult.Status = HealthStatusDegraded
			parentResult.Message = "One or more health checks are degraded"
		}
	}
	
	parentResult.Duration = time.Since(startTime)
	
	// Add to history
	h.addToHistory(parentResult)
	
	// Update statistics
	h.updateHealthStats(parentResult)
	
	return parentResult, nil
}

// GetSystemHealth returns the overall system health
func (h *BaseHealthChecker) GetSystemHealth(ctx context.Context) (*HealthCheckResult, error) {
	return h.RunAllHealthChecks(ctx)
}

// GetHealthHistory returns health check history
func (h *BaseHealthChecker) GetHealthHistory(limit int) []*HealthCheckResult {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	if limit <= 0 || limit > len(h.history) {
		limit = len(h.history)
	}
	
	return h.history[len(h.history)-limit:]
}

// GetHealthStats returns health check statistics
func (h *BaseHealthChecker) GetHealthStats() HealthStats {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.stats
}

// Start starts the health checker
func (h *BaseHealthChecker) Start(ctx context.Context) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if h.running {
		return fmt.Errorf("health checker is already running")
	}
	
	h.running = true
	
	// Start health check goroutine
	go h.healthCheckLoop(ctx)
	
	return nil
}

// Stop stops the health checker
func (h *BaseHealthChecker) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if !h.running {
		return fmt.Errorf("health checker is not running")
	}
	
	h.running = false
	close(h.stopChan)
	
	return nil
}

// IsRunning returns whether the health checker is running
func (h *BaseHealthChecker) IsRunning() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.running
}

// healthCheckLoop runs the health check loop
func (h *BaseHealthChecker) healthCheckLoop(ctx context.Context) {
	ticker := time.NewTicker(h.config.Interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			// Run all health checks
			_, err := h.RunAllHealthChecks(ctx)
			if err != nil {
				// Log error
				continue
			}
		case <-h.stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}

// addToHistory adds a health check result to history
func (h *BaseHealthChecker) addToHistory(result *HealthCheckResult) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.history = append(h.history, result)
	
	// Trim history if it gets too large
	if len(h.history) > 1000 {
		h.history = h.history[1:]
	}
}

// updateHealthStats updates health check statistics
func (h *BaseHealthChecker) updateHealthStats(result *HealthCheckResult) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	// Initialize maps if they are nil
	if h.stats.ChecksByType == nil {
		h.stats.ChecksByType = make(map[HealthCheckType]int)
	}
	
	// Update counts
	h.stats.TotalChecks++
	h.stats.ChecksByType[result.Type]++
	
	// Update status counts
	switch result.Status {
	case HealthStatusHealthy:
		h.stats.HealthyChecks++
	case HealthStatusDegraded:
		h.stats.DegradedChecks++
	case HealthStatusUnhealthy:
		h.stats.UnhealthyChecks++
	case HealthStatusUnknown:
		h.stats.UnknownChecks++
	}
	
	// Update last check time
	h.stats.LastCheckTime = &result.Timestamp
	
	// Calculate success rate
	h.stats.SuccessRate = float64(h.stats.HealthyChecks) / float64(h.stats.TotalChecks)
	
	// Update average duration
	if h.stats.AverageDuration == 0 {
		h.stats.AverageDuration = result.Duration
	} else {
		h.stats.AverageDuration = (h.stats.AverageDuration + result.Duration) / 2
	}
}

// GetHealthCountByType returns health check count by type
func (h *BaseHealthChecker) GetHealthCountByType(checkType HealthCheckType) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.stats.ChecksByType[checkType]
}

// GetHealthCountByStatus returns health check count by status
func (h *BaseHealthChecker) GetHealthCountByStatus(status HealthStatus) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	switch status {
	case HealthStatusHealthy:
		return h.stats.HealthyChecks
	case HealthStatusDegraded:
		return h.stats.DegradedChecks
	case HealthStatusUnhealthy:
		return h.stats.UnhealthyChecks
	case HealthStatusUnknown:
		return h.stats.UnknownChecks
	default:
		return 0
	}
}

// GetAverageHealthDuration returns the average duration of health checks
func (h *BaseHealthChecker) GetAverageHealthDuration() time.Duration {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.stats.AverageDuration
}

// GetHealthSuccessRate returns the health check success rate
func (h *BaseHealthChecker) GetHealthSuccessRate() float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.stats.SuccessRate
}

// GetLastHealthCheckTime returns the last health check time
func (h *BaseHealthChecker) GetLastHealthCheckTime() *time.Time {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.stats.LastCheckTime
}

// GetHealthHistorySize returns the current health history size
func (h *BaseHealthChecker) GetHealthHistorySize() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return len(h.history)
}

// ClearHealthHistory clears health check history
func (h *BaseHealthChecker) ClearHealthHistory() {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.history = make([]*HealthCheckResult, 0)
}

// GetHealthChecksByStatus returns health checks by status
func (h *BaseHealthChecker) GetHealthChecksByStatus(status HealthStatus) []HealthCheck {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var checks []HealthCheck
	for _, check := range h.healthChecks {
		if check.GetLastResult() != nil && check.GetLastResult().Status == status {
			checks = append(checks, check)
		}
	}
	
	return checks
}

// GetHealthChecksEnabled returns enabled health checks
func (h *BaseHealthChecker) GetHealthChecksEnabled() []HealthCheck {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var checks []HealthCheck
	for _, check := range h.healthChecks {
		if check.GetEnabled() {
			checks = append(checks, check)
		}
	}
	
	return checks
}

// GetHealthChecksDisabled returns disabled health checks
func (h *BaseHealthChecker) GetHealthChecksDisabled() []HealthCheck {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var checks []HealthCheck
	for _, check := range h.healthChecks {
		if !check.GetEnabled() {
			checks = append(checks, check)
		}
	}
	
	return checks
}

// EnableHealthCheck enables a health check
func (h *BaseHealthChecker) EnableHealthCheck(id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	check, exists := h.healthChecks[id]
	if !exists {
		return fmt.Errorf("health check %s not found", id)
	}
	
	check.SetEnabled(true)
	
	return nil
}

// DisableHealthCheck disables a health check
func (h *BaseHealthChecker) DisableHealthCheck(id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	check, exists := h.healthChecks[id]
	if !exists {
		return fmt.Errorf("health check %s not found", id)
	}
	
	check.SetEnabled(false)
	
	return nil
}

// SetHealthCheckInterval sets the interval for a health check
func (h *BaseHealthChecker) SetHealthCheckInterval(id string, interval time.Duration) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if _, exists := h.healthChecks[id]; !exists {
		return fmt.Errorf("health check %s not found", id)
	}
	
	// This would require the HealthCheck interface to have a SetInterval method
	// For now, we'll just update the config
	h.config.Interval = interval
	
	return nil
}

// SetHealthCheckTimeout sets the timeout for a health check
func (h *BaseHealthChecker) SetHealthCheckTimeout(id string, timeout time.Duration) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if _, exists := h.healthChecks[id]; !exists {
		return fmt.Errorf("health check %s not found", id)
	}
	
	// This would require the HealthCheck interface to have a SetTimeout method
	// For now, we'll just update the config
	h.config.Timeout = timeout
	
	return nil
}

// SetHealthCheckRetries sets the number of retries for a health check
func (h *BaseHealthChecker) SetHealthCheckRetries(id string, retries int) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if _, exists := h.healthChecks[id]; !exists {
		return fmt.Errorf("health check %s not found", id)
	}
	
	// This would require the HealthCheck interface to have a SetRetries method
	// For now, we'll just update the config
	h.config.Retries = retries
	
	return nil
}

// GetHealthCheckConfig returns the health check configuration
func (h *BaseHealthChecker) GetHealthCheckConfig() *HealthCheckConfig {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return h.config
}

// SetHealthCheckConfig sets the health check configuration
func (h *BaseHealthChecker) SetHealthCheckConfig(config *HealthCheckConfig) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.config = config
}

// GetHealthChecksCount returns the number of health checks
func (h *BaseHealthChecker) GetHealthChecksCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return len(h.healthChecks)
}

// IsSystemHealthy returns whether the system is healthy
func (h *BaseHealthChecker) IsSystemHealthy() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	// Check if there are any unhealthy checks
	for _, check := range h.healthChecks {
		if check.GetEnabled() && check.GetLastResult() != nil {
			if check.GetLastResult().Status == HealthStatusUnhealthy {
				return false
			}
		}
	}
	
	return true
}

// GetSystemHealthStatus returns the system health status
func (h *BaseHealthChecker) GetSystemHealthStatus() HealthStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	// Check for any unhealthy checks
	for _, check := range h.healthChecks {
		if check.GetEnabled() && check.GetLastResult() != nil {
			if check.GetLastResult().Status == HealthStatusUnhealthy {
				return HealthStatusUnhealthy
			}
		}
	}
	
	// Check for any degraded checks
	for _, check := range h.healthChecks {
		if check.GetEnabled() && check.GetLastResult() != nil {
			if check.GetLastResult().Status == HealthStatusDegraded {
				return HealthStatusDegraded
			}
		}
	}
	
	// Check for any unknown checks
	for _, check := range h.healthChecks {
		if check.GetEnabled() && check.GetLastResult() != nil {
			if check.GetLastResult().Status == HealthStatusUnknown {
				return HealthStatusUnknown
			}
		}
	}
	
	return HealthStatusHealthy
}

// GetSystemHealthMessage returns the system health message
func (h *BaseHealthChecker) GetSystemHealthMessage() string {
	status := h.GetSystemHealthStatus()
	switch status {
	case HealthStatusHealthy:
		return "All systems are healthy"
	case HealthStatusDegraded:
		return "Some systems are degraded"
	case HealthStatusUnhealthy:
		return "Some systems are unhealthy"
	case HealthStatusUnknown:
		return "System health is unknown"
	default:
		return "Unknown system health status"
	}
}