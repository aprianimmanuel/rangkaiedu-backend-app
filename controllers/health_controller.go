package controllers

import (
	"context"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/aprianimmanuel/rangkaiedu-backend/utils/db"
)

// HealthCheckResponse represents the response structure for health check endpoint
type HealthCheckResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version"`
	Uptime    time.Duration          `json:"uptime"`
	Checks    map[string]interface{} `json:"checks"`
	Info      map[string]interface{} `json:"info,omitempty"`
}

// HealthCheckDetailedResponse represents a detailed health check response
type HealthCheckDetailedResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version"`
	Uptime    time.Duration          `json:"uptime"`
	Checks    map[string]interface{} `json:"checks"`
	Info      map[string]interface{} `json:"info"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// HealthCheckRequest represents a request for specific health check
type HealthCheckRequest struct {
	CheckType string `json:"check_type"` // "basic", "database", "full"
}

// HealthCheck performs a comprehensive health check of the application
func HealthCheck(c *gin.Context) {
	startTime := time.Now()
	
	// Basic application health check
	appStatus := "healthy"
	appChecks := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"uptime":    time.Since(startTime),
		"version":   "1.0.0",
		"go_version": runtime.Version(),
		"goroutines": runtime.NumGoroutine(),
		"memory_usage": getMemoryUsage(),
	}

	// Database health check
	dbStatus, dbInfo := checkDatabaseHealth()
	appChecks["database"] = map[string]interface{}{
		"status": dbStatus,
		"info":   dbInfo,
	}

	// Overall status determination
	if dbStatus != "healthy" {
		appStatus = "degraded"
	}

	response := HealthCheckResponse{
		Status:    appStatus,
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    time.Since(startTime),
		Checks:    appChecks,
	}

	c.JSON(http.StatusOK, response)
}

// HealthCheckDetailed performs a detailed health check with more information
func HealthCheckDetailed(c *gin.Context) {
	startTime := time.Now()
	
	// Basic application health check
	appStatus := "healthy"
	appChecks := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"uptime":    time.Since(startTime),
		"version":   "1.0.0",
		"go_version": runtime.Version(),
		"goroutines": runtime.NumGoroutine(),
		"memory_usage": getMemoryUsage(),
		"cpu_usage": getCPUUsage(),
	}

	// Database health check
	dbStatus, dbInfo := checkDatabaseHealth()
	appChecks["database"] = map[string]interface{}{
		"status": dbStatus,
		"info":   dbInfo,
		"details": getDatabaseDetails(),
	}

	// Overall status determination
	if dbStatus != "healthy" {
		appStatus = "degraded"
	}

	// Additional system information
	systemInfo := map[string]interface{}{
		"host":       getHostname(),
		"architecture": runtime.GOARCH,
		"os":         runtime.GOOS,
		"max_procs":  runtime.GOMAXPROCS(0),
		"num_cpu":    runtime.NumCPU(),
	}

	response := HealthCheckDetailedResponse{
		Status:    appStatus,
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    time.Since(startTime),
		Checks:    appChecks,
		Info:      systemInfo,
		Details:   map[string]interface{}{
			"environment": getEnvironment(),
			"build_time":  getBuildTime(),
		},
	}

	c.JSON(http.StatusOK, response)
}

// HealthCheckDatabase performs a database-specific health check
func HealthCheckDatabase(c *gin.Context) {
	status, info := checkDatabaseHealth()
	
	response := map[string]interface{}{
		"status": status,
		"timestamp": time.Now(),
		"info":   info,
		"details": getDatabaseDetails(),
	}

	var statusCode int
	switch status {
	case "healthy":
		statusCode = http.StatusOK
	case "degraded":
		statusCode = http.StatusPartialContent
	default:
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// HealthCheckBasic performs a basic health check (fast response)
func HealthCheckBasic(c *gin.Context) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"version":   "1.0.0",
		"uptime":    time.Since(time.Now()),
	}

	c.JSON(http.StatusOK, response)
}

// checkDatabaseHealth performs a comprehensive database health check
func checkDatabaseHealth() (string, map[string]interface{}) {
	// Use the global database connection pool
	pool := db.GetDB()
	if pool == nil {
		log.Printf("Database pool is not initialized")
		return "unhealthy", map[string]interface{}{
			"error": "Database pool is not initialized",
			"details": "The database connection pool has not been initialized",
		}
	}

	// Create a context with timeout for database connection check
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test database connection with a simple query
	var result string
	err := pool.QueryRowContext(ctx, "SELECT NOW()").Scan(&result)
	if err != nil {
		log.Printf("Database query failed: %v", err)
		return "unhealthy", map[string]interface{}{
			"error": "Database query failed",
			"details": err.Error(),
		}
	}

	// For sql.DB, we can't get detailed stats like with pgxpool.Pool
	// But we can check if the connection is working
	return "healthy", map[string]interface{}{
		"query_result": result,
	}
}

// getDatabaseDetails returns detailed database information
func getDatabaseDetails() map[string]interface{} {
	return map[string]interface{}{
		"max_connections": "not_available",
		"min_connections": "not_available",
		"max_conn_lifetime": "not_available",
		"max_conn_idle_time": "not_available",
		"health_check_period": "not_available",
		"database_host": "localhost",
		"database_port": "5432",
		"database_name": "rangkaiedu_dev",
	}
}

// getMemoryUsage returns current memory usage information
func getMemoryUsage() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return map[string]interface{}{
		"allocated_mb":     bToMB(m.Alloc),
		"total_allocated_mb": bToMB(m.TotalAlloc),
		"system_memory_mb": bToMB(m.Sys),
		"gc_runs":          m.NumGC,
		"gc_pause_total_ns": m.PauseTotalNs,
		"gc_cpu_fraction":  m.GCCPUFraction,
	}
}

// getCPUUsage returns CPU usage information (placeholder - would need actual implementation)
func getCPUUsage() map[string]interface{} {
	return map[string]interface{}{
		"usage_percent": 0.0, // Placeholder - would need actual CPU usage calculation
		"cores":         runtime.NumCPU(),
	}
}

// getHostname returns the hostname (placeholder - would need actual implementation)
func getHostname() string {
	// Placeholder - would need actual hostname implementation
	return "unknown"
}

// getEnvironment returns the current environment
func getEnvironment() string {
	// Since Config doesn't have Environment field, we'll return a default value
	// In a real implementation, you might want to add this to the Config struct
	return "development"
}

// getBuildTime returns the build time (placeholder - would need actual implementation)
func getBuildTime() string {
	// Placeholder - would need actual build time implementation
	return "unknown"
}

// Helper functions

// bToMB converts bytes to megabytes
func bToMB(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

// maskDatabasePassword masks the password in the database URL
func maskDatabasePassword(dsn string) string {
	// Simple implementation - replace password with *****
	// In a real implementation, you would parse the URL and mask the password
	if len(dsn) > 10 {
		return dsn[:8] + "*****"
	}
	return "*****"
}

// extractDatabaseHost extracts the host from the database URL
func extractDatabaseHost(dsn string) string {
	// Placeholder - would need actual URL parsing
	return "localhost"
}

// extractDatabasePort extracts the port from the database URL
func extractDatabasePort(dsn string) string {
	// Placeholder - would need actual URL parsing
	return "5432"
}

// extractDatabaseName extracts the database name from the database URL
func extractDatabaseName(dsn string) string {
	// Placeholder - would need actual URL parsing
	return "rangkai_edu"
}