package monitoring

import (
	"context"
	"fmt"
	"time"
)


// Global security logging functions

// LogSecurityEvent logs a security event globally
func LogSecurityEvent(ctx context.Context, event *SecurityEvent) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	
	return service.LogSecurityEvent(ctx, event)
}

// LogAuthFailure logs an authentication failure globally
func LogAuthFailure(ctx context.Context, username, ip, userAgent string, details map[string]interface{}) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	
	return service.LogAuthFailure(ctx, username, ip, userAgent, details)
}

// LogAuthSuccess logs an authentication success globally
func LogAuthSuccess(ctx context.Context, userID, username, ip, userAgent string, details map[string]interface{}) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	
	return service.LogAuthSuccess(ctx, userID, username, ip, userAgent, details)
}

// LogAuthLockout logs an authentication lockout globally
func LogAuthLockout(ctx context.Context, username, ip string, attempts int, details map[string]interface{}) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	
	return service.LogAuthLockout(ctx, username, ip, attempts, details)
}

// LogAuthBruteForce logs a brute force attack globally
func LogAuthBruteForce(ctx context.Context, ip string, attempts int, timeframe time.Duration, details map[string]interface{}) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	
	return service.LogAuthBruteForce(ctx, ip, attempts, timeframe, details)
}

// LogAuthSession logs an authentication session globally
func LogAuthSession(ctx context.Context, userID, sessionID, ip string, action string, details map[string]interface{}) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	
	return service.LogAuthSession(ctx, userID, sessionID, ip, action, details)
}

// LogAuthToken logs an authentication token globally
func LogAuthToken(ctx context.Context, userID, tokenID, ip string, action string, details map[string]interface{}) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	
	return service.LogAuthToken(ctx, userID, tokenID, ip, action, details)
}

// LogAuthMFA logs a multi-factor authentication event globally
func LogAuthMFA(ctx context.Context, userID, username, ip, method string, success bool, details map[string]interface{}) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	
	return service.LogAuthMFA(ctx, userID, username, ip, method, success, details)
}

// LogAuthViolation logs an authorization violation globally
func LogAuthViolation(ctx context.Context, userID, username, resource, action string, details map[string]interface{}) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	
	return service.LogAuthViolation(ctx, userID, username, resource, action, details)
}

// LogAuthThreat logs a security threat globally
func LogAuthThreat(ctx context.Context, threatType, ip, userAgent string, details map[string]interface{}) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	
	return service.LogAuthThreat(ctx, threatType, ip, userAgent, details)
}

// LogAuthAudit logs an audit event globally
func LogAuthAudit(ctx context.Context, userID, username, action, resource string, details map[string]interface{}) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	
	return service.LogAuthAudit(ctx, userID, username, action, resource, details)
}

// GetSecurityLogger returns the security logger globally
func GetSecurityLogger() SecurityEventLogger {
	service := GetService()
	if service == nil {
		return nil
	}
	
	return service.GetSecurityLogger()
}

// IsSecurityLogsEnabled returns whether security logging is enabled globally
func IsSecurityLogsEnabled() bool {
	service := GetService()
	if service == nil {
		return false
	}
	
	return service.IsSecurityLogsEnabled()
}

// GetSecurityEventCount returns the number of security events logged globally
func GetSecurityEventCount() int {
	service := GetService()
	if service == nil {
		return 0
	}
	
	if logger, ok := service.securityLogger.(interface{ GetEventCount() int }); ok {
		return logger.GetEventCount()
	}
	
	return 0
}

// GetSecurityEvents returns security events globally
func GetSecurityEvents(ctx context.Context, limit int, offset int) ([]*SecurityEvent, error) {
	service := GetService()
	if service == nil {
		return nil, fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return nil, fmt.Errorf("security logger is not enabled")
	}
	
	if getter, ok := service.securityLogger.(interface{ GetEvents(ctx context.Context, limit int, offset int) ([]*SecurityEvent, error) }); ok {
		return getter.GetEvents(ctx, limit, offset)
	}
	
	return nil, fmt.Errorf("security logger does not support event retrieval")
}

// ClearSecurityEvents clears security events globally
func ClearSecurityEvents() error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	
	if clearer, ok := service.securityLogger.(interface{ ClearEvents() error }); ok {
		return clearer.ClearEvents()
	}
	
	return fmt.Errorf("security logger does not support event clearing")
}

// GetSecurityConfig returns the security configuration globally
func GetSecurityConfig() *SecurityLoggingConfig {
	service := GetService()
	if service == nil {
		return nil
	}
	
	if service.securityLogger == nil {
		return nil
	}
	
	// Try to get the config directly from the service
	config := service.config.SecurityConfig
	if config != nil {
		return config
	}
	
	// Fallback: Get the config as map[string]interface{} and convert it to SecurityLoggingConfig
	configMap := service.securityLogger.GetConfig()
	if configMap == nil {
		return nil
	}
	
	// Create a new SecurityLoggingConfig and populate it from the map
	config = &SecurityLoggingConfig{}
	
	// Convert map values to SecurityLoggingConfig fields
	if enabled, ok := configMap["enabled"].(bool); ok {
		config.Enabled = enabled
	}
	
	if level, ok := configMap["level"].(string); ok {
		config.Level = level
	}
	
	if format, ok := configMap["format"].(string); ok {
		config.Format = format
	}
	
	if bufferSize, ok := configMap["buffer_size"].(int); ok {
		config.BufferSize = bufferSize
	}
	
	if flushInterval, ok := configMap["flush_interval"].(float64); ok {
		config.FlushInterval = time.Duration(flushInterval)
	}
	
	if batchSize, ok := configMap["batch_size"].(int); ok {
		config.BatchSize = batchSize
	}
	
	if retentionDays, ok := configMap["retention_days"].(int); ok {
		config.RetentionDays = retentionDays
	}
	
	if maxFileSize, ok := configMap["max_file_size"].(int64); ok {
		config.MaxFileSize = maxFileSize
	}
	
	if maxBackups, ok := configMap["max_backups"].(int); ok {
		config.MaxBackups = maxBackups
	}
	
	if sanitizeData, ok := configMap["sanitize_data"].(bool); ok {
		config.SanitizeData = sanitizeData
	}
	
	if encryptLogs, ok := configMap["encrypt_logs"].(bool); ok {
		config.EncryptLogs = encryptLogs
	}
	
	if encryptionKey, ok := configMap["encryption_key"].(string); ok {
		config.EncryptionKey = encryptionKey
	}
	
	if asyncLogging, ok := configMap["async_logging"].(bool); ok {
		config.AsyncLogging = asyncLogging
	}
	
	if queueSize, ok := configMap["queue_size"].(int); ok {
		config.QueueSize = queueSize
	}
	
	if workerCount, ok := configMap["worker_count"].(int); ok {
		config.WorkerCount = workerCount
	}
	
	if environment, ok := configMap["environment"].(string); ok {
		config.Environment = environment
	}
	
	if debug, ok := configMap["debug"].(bool); ok {
		config.Debug = debug
	}
	
	return config
}

// SetSecurityConfig sets the security configuration globally
func SetSecurityConfig(config *SecurityLoggingConfig) error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	
	if setter, ok := service.securityLogger.(interface{ SetConfig(*SecurityLoggingConfig) error }); ok {
		return setter.SetConfig(config)
	}
	
	return fmt.Errorf("security logger does not support configuration setting")
}

// EnableSecurityLogging enables security logging globally
func EnableSecurityLogging() error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	
	if enabler, ok := service.securityLogger.(interface{ Enable() error }); ok {
		return enabler.Enable()
	}
	
	return fmt.Errorf("security logger does not support enabling")
}

// DisableSecurityLogging disables security logging globally
func DisableSecurityLogging() error {
	service := GetService()
	if service == nil {
		return fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return fmt.Errorf("security logger is not enabled")
	}
	
	if disabler, ok := service.securityLogger.(interface{ Disable() error }); ok {
		return disabler.Disable()
	}
	
	return fmt.Errorf("security logger does not support disabling")
}

// IsSecurityLoggingActive returns whether security logging is active globally
func IsSecurityLoggingActive() bool {
	service := GetService()
	if service == nil {
		return false
	}
	
	if service.securityLogger == nil {
		return false
	}
	
	if checker, ok := service.securityLogger.(interface{ IsActive() bool }); ok {
		return checker.IsActive()
	}
	
	return false
}
// GetSecurityEventsByMultipleFilters returns security events filtered by multiple criteria globally
func GetSecurityEventsByMultipleFilters(ctx context.Context, filters map[string]interface{}, limit int, offset int) ([]*SecurityEvent, error) {
	service := GetService()
	if service == nil {
		return nil, fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return nil, fmt.Errorf("security logger is not enabled")
	}
	
	// Try to use the logger's GetEventsByMultipleFilters method if it exists
	if getter, ok := service.securityLogger.(interface{ GetEventsByMultipleFilters(ctx context.Context, filters map[string]interface{}, limit int, offset int) ([]*SecurityEvent, error) }); ok {
		return getter.GetEventsByMultipleFilters(ctx, filters, limit, offset)
	}
	
	// Fallback: get all events and filter them manually
	events, err := GetSecurityEvents(ctx, 0, 0) // Get all events
	if err != nil {
		return nil, err
	}
	
	// Filter events based on the provided filters
	var filteredEvents []*SecurityEvent
	for _, event := range events {
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

// GetSecurityEventsByType returns security events filtered by type globally
func GetSecurityEventsByType(ctx context.Context, eventType string, limit int, offset int) ([]*SecurityEvent, error) {
	filters := map[string]interface{}{
		"type": eventType,
	}
	return GetSecurityEventsByMultipleFilters(ctx, filters, limit, offset)
}

// GetSecurityEventsByType returns security




// GetSecurityEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffset returns security events filtered by time range and all criteria with pagination, sorting, limit, and offset globally
func GetSecurityEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffset(ctx context.Context, startTime, endTime time.Time, eventType, userID, ip, severity, resource, action, status string, labelKey, labelValue string, sortBy string, sortOrder string, page, pageSize, limit, offset int) ([]*SecurityEvent, int, error) {
	service := GetService()
	if service == nil {
		return nil, 0, fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return nil, 0, fmt.Errorf("security logger is not enabled")
	}
	
	if getter, ok := service.securityLogger.(interface{ GetEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffset(ctx context.Context, startTime, endTime time.Time, eventType, userID, ip, severity, resource, action, status string, labelKey, labelValue string, sortBy string, sortOrder string, page, pageSize, limit, offset int) ([]*SecurityEvent, int, error) }); ok {
		return getter.GetEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffset(ctx, startTime, endTime, eventType, userID, ip, severity, resource, action, status, labelKey, labelValue, sortBy, sortOrder, page, pageSize, limit, offset)
	}
	
	// Fallback: use multiple filters with time range and pagination
	filters := make(map[string]interface{})
	if eventType != "" {
		filters["type"] = eventType
	}
	if userID != "" {
		filters["user_id"] = userID
	}
	if ip != "" {
		filters["ip"] = ip
	}
	if severity != "" {
		filters["severity"] = severity
	}
	if resource != "" {
		filters["resource"] = resource
	}
	if action != "" {
		filters["action"] = action
	}
	if status != "" {
		filters["status"] = status
	}
	if labelKey != "" && labelValue != "" {
		filters["label"] = map[string]string{labelKey: labelValue}
	}
	
	// Apply offset
	if offset > 0 {
		events, err := GetSecurityEventsByMultipleFilters(ctx, filters, pageSize, offset)
		if err != nil {
			return nil, 0, err
		}
		
		// Sort events
		switch sortBy {
		case "timestamp":
			if sortOrder == "desc" {
				// Sort by timestamp descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Timestamp.Before(events[j].Timestamp) {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by timestamp ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Timestamp.After(events[j].Timestamp) {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "severity":
			if sortOrder == "desc" {
				// Sort by severity descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Severity < events[j].Severity {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by severity ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Severity > events[j].Severity {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "type":
			if sortOrder == "desc" {
				// Sort by type descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Type < events[j].Type {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by type ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Type > events[j].Type {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "user_id":
			if sortOrder == "desc" {
				// Sort by user_id descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].UserID < events[j].UserID {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by user_id ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].UserID > events[j].UserID {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "ip":
			if sortOrder == "desc" {
				// Sort by ip descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].IP < events[j].IP {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by ip ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].IP > events[j].IP {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "resource":
			if sortOrder == "desc" {
				// Sort by resource descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Resource < events[j].Resource {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by resource ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Resource > events[j].Resource {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "action":
			if sortOrder == "desc" {
				// Sort by action descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Action < events[j].Action {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by action ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Action > events[j].Action {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "status":
			if sortOrder == "desc" {
				// Sort by status descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Status < events[j].Status {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by status ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Status > events[j].Status {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		}
		
		// Apply limit
		if limit > 0 && len(events) > limit {
			events = events[:limit]
		}
		
		// Get total count for pagination
		total := GetSecurityEventCount()
		
		return events, total, nil
	}
	
	// Get total count for pagination
	total := GetSecurityEventCount()
	
	return nil, total, nil
}

// GetSecurityEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrder returns security events filtered by time range and all criteria with pagination, sorting, limit, offset, and order globally
func GetSecurityEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrder(ctx context.Context, startTime, endTime time.Time, eventType, userID, ip, severity, resource, action, status string, labelKey, labelValue string, sortBy string, sortOrder string, page, pageSize, limit, offset int, order string) ([]*SecurityEvent, int, error) {
	service := GetService()
	if service == nil {
		return nil, 0, fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return nil, 0, fmt.Errorf("security logger is not enabled")
	}
	
	if getter, ok := service.securityLogger.(interface{ GetEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrder(ctx context.Context, startTime, endTime time.Time, eventType, userID, ip, severity, resource, action, status string, labelKey, labelValue string, sortBy string, sortOrder string, page, pageSize, limit, offset int, order string) ([]*SecurityEvent, int, error) }); ok {
		return getter.GetEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrder(ctx, startTime, endTime, eventType, userID, ip, severity, resource, action, status, labelKey, labelValue, sortBy, sortOrder, page, pageSize, limit, offset, order)
	}
	
	// Fallback: use multiple filters with time range and pagination
	filters := make(map[string]interface{})
	if eventType != "" {
		filters["type"] = eventType
	}
	if userID != "" {
		filters["user_id"] = userID
	}
	if ip != "" {
		filters["ip"] = ip
	}
	if severity != "" {
		filters["severity"] = severity
	}
	if resource != "" {
		filters["resource"] = resource
	}
	if action != "" {
		filters["action"] = action
	}
	if status != "" {
		filters["status"] = status
	}
	if labelKey != "" && labelValue != "" {
		filters["label"] = map[string]string{labelKey: labelValue}
	}
	
	// Apply offset
	if offset > 0 {
		events, err := GetSecurityEventsByMultipleFilters(ctx, filters, pageSize, offset)
		if err != nil {
			return nil, 0, err
		}
		
		// Sort events
		switch sortBy {
		case "timestamp":
			if sortOrder == "desc" {
				// Sort by timestamp descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Timestamp.Before(events[j].Timestamp) {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by timestamp ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Timestamp.After(events[j].Timestamp) {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "severity":
			if sortOrder == "desc" {
				// Sort by severity descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Severity < events[j].Severity {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by severity ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Severity > events[j].Severity {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "type":
			if sortOrder == "desc" {
				// Sort by type descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Type < events[j].Type {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by type ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Type > events[j].Type {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "user_id":
			if sortOrder == "desc" {
				// Sort by user_id descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].UserID < events[j].UserID {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by user_id ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].UserID > events[j].UserID {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "ip":
			if sortOrder == "desc" {
				// Sort by ip descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].IP < events[j].IP {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by ip ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].IP > events[j].IP {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "resource":
			if sortOrder == "desc" {
				// Sort by resource descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Resource < events[j].Resource {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by resource ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Resource > events[j].Resource {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "action":
			if sortOrder == "desc" {
				// Sort by action descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Action < events[j].Action {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by action ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Action > events[j].Action {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "status":
			if sortOrder == "desc" {
				// Sort by status descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Status < events[j].Status {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by status ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Status > events[j].Status {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		}
		
		// Apply limit
		if limit > 0 && len(events) > limit {
			events = events[:limit]
		}
		
		// Get total count for pagination
		total := GetSecurityEventCount()
		
		return events, total, nil
	}
	
	// Get total count for pagination
	total := GetSecurityEventCount()
	
	return nil, total, nil
}

// GetSecurityEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrderAndDirection returns security events filtered by time range and all criteria with pagination, sorting, limit, offset, order, and direction globally
func GetSecurityEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrderAndDirection(ctx context.Context, startTime, endTime time.Time, eventType, userID, ip, severity, resource, action, status string, labelKey, labelValue string, sortBy string, sortOrder string, page, pageSize, limit, offset int, order, direction string) ([]*SecurityEvent, int, error) {
	service := GetService()
	if service == nil {
		return nil, 0, fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return nil, 0, fmt.Errorf("security logger is not enabled")
	}
	
	if getter, ok := service.securityLogger.(interface{ GetEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrderAndDirection(ctx context.Context, startTime, endTime time.Time, eventType, userID, ip, severity, resource, action, status string, labelKey, labelValue string, sortBy string, sortOrder string, page, pageSize, limit, offset int, order, direction string) ([]*SecurityEvent, int, error) }); ok {
		return getter.GetEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrderAndDirection(ctx, startTime, endTime, eventType, userID, ip, severity, resource, action, status, labelKey, labelValue, sortBy, sortOrder, page, pageSize, limit, offset, order, direction)
	}
	
	// Fallback: use multiple filters with time range and pagination
	filters := make(map[string]interface{})
	if eventType != "" {
		filters["type"] = eventType
	}
	if userID != "" {
		filters["user_id"] = userID
	}
	if ip != "" {
		filters["ip"] = ip
	}
	if severity != "" {
		filters["severity"] = severity
	}
	if resource != "" {
		filters["resource"] = resource
	}
	if action != "" {
		filters["action"] = action
	}
	if status != "" {
		filters["status"] = status
	}
	if labelKey != "" && labelValue != "" {
		filters["label"] = map[string]string{labelKey: labelValue}
	}
	
	// Apply offset
	if offset > 0 {
		events, err := GetSecurityEventsByMultipleFilters(ctx, filters, pageSize, offset)
		if err != nil {
			return nil, 0, err
		}
		
		// Sort events
		switch sortBy {
		case "timestamp":
			if sortOrder == "desc" {
				// Sort by timestamp descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Timestamp.Before(events[j].Timestamp) {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by timestamp ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Timestamp.After(events[j].Timestamp) {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "severity":
			if sortOrder == "desc" {
				// Sort by severity descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Severity < events[j].Severity {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by severity ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Severity > events[j].Severity {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "type":
			if sortOrder == "desc" {
				// Sort by type descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Type < events[j].Type {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by type ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Type > events[j].Type {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "user_id":
			if sortOrder == "desc" {
				// Sort by user_id descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].UserID < events[j].UserID {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by user_id ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].UserID > events[j].UserID {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "ip":
			if sortOrder == "desc" {
				// Sort by ip descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].IP < events[j].IP {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by ip ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].IP > events[j].IP {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "resource":
			if sortOrder == "desc" {
				// Sort by resource descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Resource < events[j].Resource {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by resource ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Resource > events[j].Resource {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "action":
			if sortOrder == "desc" {
				// Sort by action descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Action < events[j].Action {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by action ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Action > events[j].Action {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "status":
			if sortOrder == "desc" {
				// Sort by status descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Status < events[j].Status {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by status ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Status > events[j].Status {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		}
		
		// Apply limit
		if limit > 0 && len(events) > limit {
			events = events[:limit]
		}
		
		// Get total count for pagination
		total := GetSecurityEventCount()
		
		return events, total, nil
	}
	
	// Get total count for pagination
	total := GetSecurityEventCount()
	
	return nil, total, nil
}

// GetSecurityEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrderAndDirectionAndPage returns security events filtered by time range and all criteria with pagination, sorting, limit, offset, order, direction, and page globally
func GetSecurityEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrderAndDirectionAndPage(ctx context.Context, startTime, endTime time.Time, eventType, userID, ip, severity, resource, action, status string, labelKey, labelValue string, sortBy string, sortOrder string, page, pageSize, limit, offset int, order, direction string) ([]*SecurityEvent, int, error) {
	service := GetService()
	if service == nil {
		return nil, 0, fmt.Errorf("monitoring service is not initialized")
	}
	
	if service.securityLogger == nil {
		return nil, 0, fmt.Errorf("security logger is not enabled")
	}
	
	if getter, ok := service.securityLogger.(interface{ GetEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrderAndDirectionAndPage(ctx context.Context, startTime, endTime time.Time, eventType, userID, ip, severity, resource, action, status string, labelKey, labelValue string, sortBy string, sortOrder string, page, pageSize, limit, offset int, order, direction string) ([]*SecurityEvent, int, error) }); ok {
		return getter.GetEventsByTimeRangeAndAllFiltersWithPaginationAndSortingAndLimitAndOffsetAndOrderAndDirectionAndPage(ctx, startTime, endTime, eventType, userID, ip, severity, resource, action, status, labelKey, labelValue, sortBy, sortOrder, page, pageSize, limit, offset, order, direction)
	}
	
	// Fallback: use multiple filters with time range and pagination
	filters := make(map[string]interface{})
	if eventType != "" {
		filters["type"] = eventType
	}
	if userID != "" {
		filters["user_id"] = userID
	}
	if ip != "" {
		filters["ip"] = ip
	}
	if severity != "" {
		filters["severity"] = severity
	}
	if resource != "" {
		filters["resource"] = resource
	}
	if action != "" {
		filters["action"] = action
	}
	if status != "" {
		filters["status"] = status
	}
	if labelKey != "" && labelValue != "" {
		filters["label"] = map[string]string{labelKey: labelValue}
	}
	
	// Apply offset
	if offset > 0 {
		events, err := GetSecurityEventsByMultipleFilters(ctx, filters, pageSize, offset)
		if err != nil {
			return nil, 0, err
		}
		
		// Sort events
		switch sortBy {
		case "timestamp":
			if sortOrder == "desc" {
				// Sort by timestamp descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Timestamp.Before(events[j].Timestamp) {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by timestamp ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Timestamp.After(events[j].Timestamp) {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "severity":
			if sortOrder == "desc" {
				// Sort by severity descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Severity < events[j].Severity {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by severity ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Severity > events[j].Severity {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "type":
			if sortOrder == "desc" {
				// Sort by type descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Type < events[j].Type {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by type ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Type > events[j].Type {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "user_id":
			if sortOrder == "desc" {
				// Sort by user_id descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].UserID < events[j].UserID {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by user_id ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].UserID > events[j].UserID {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "ip":
			if sortOrder == "desc" {
				// Sort by ip descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].IP < events[j].IP {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by ip ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].IP > events[j].IP {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "resource":
			if sortOrder == "desc" {
				// Sort by resource descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Resource < events[j].Resource {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by resource ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Resource > events[j].Resource {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "action":
			if sortOrder == "desc" {
				// Sort by action descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Action < events[j].Action {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by action ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Action > events[j].Action {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		case "status":
			if sortOrder == "desc" {
				// Sort by status descending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Status < events[j].Status {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			} else {
				// Sort by status ascending
				for i := 0; i < len(events)-1; i++ {
					for j := i + 1; j < len(events); j++ {
						if events[i].Status > events[j].Status {
							events[i], events[j] = events[j], events[i]
						}
					}
				}
			}
		}
		
		// Apply limit
		if limit > 0 && len(events) > limit {
			events = events[:limit]
		}
		
		// Get total count for pagination
		total := GetSecurityEventCount()
		
		return events, total, nil
	}
	
	// Get total count for pagination
	total := GetSecurityEventCount()
	
	return nil, total, nil
}

