package monitoring

import (
	"fmt"
	"time"
)

// SecurityLoggingConfig represents the configuration for security event logging
type SecurityLoggingConfig struct {
	// Basic configuration
	Enabled          bool          `json:"enabled" yaml:"enabled"`
	Level            string        `json:"level" yaml:"level"` // low, medium, high, critical
	Format           string        `json:"format" yaml:"format"` // json, text, csv
	
	// Output destinations
	FileOutput       *FileOutputConfig   `json:"file_output,omitempty" yaml:"file_output,omitempty"`
	ConsoleOutput    *ConsoleOutputConfig `json:"console_output,omitempty" yaml:"console_output,omitempty"`
	HTTPOutput       *HTTPOutputConfig  `json:"http_output,omitempty" yaml:"http_output,omitempty"`
	SyslogOutput     *SyslogOutputConfig `json:"syslog_output,omitempty" yaml:"syslog_output,omitempty"`
	
	// Buffering and batching
	BufferSize       int           `json:"buffer_size" yaml:"buffer_size"`
	FlushInterval    time.Duration `json:"flush_interval" yaml:"flush_interval"`
	BatchSize        int           `json:"batch_size" yaml:"batch_size"`
	
	// Filtering
	IncludeEvents    []string      `json:"include_events,omitempty" yaml:"include_events,omitempty"`
	ExcludeEvents    []string      `json:"exclude_events,omitempty" yaml:"exclude_events,omitempty"`
	IncludeUsers     []string      `json:"include_users,omitempty" yaml:"include_users,omitempty"`
	ExcludeUsers     []string      `json:"exclude_users,omitempty" yaml:"exclude_users,omitempty"`
	IncludeIPs       []string      `json:"include_ips,omitempty" yaml:"include_ips,omitempty"`
	ExcludeIPs       []string      `json:"exclude_ips,omitempty" yaml:"exclude_ips,omitempty"`
	
	// Sensitive data handling
	SanitizeData     bool          `json:"sanitize_data" yaml:"sanitize_data"`
	MaskPatterns     []string      `json:"mask_patterns,omitempty" yaml:"mask_patterns,omitempty"`
	
	// Retention and rotation
	RetentionDays    int           `json:"retention_days" yaml:"retention_days"`
	MaxFileSize      int64         `json:"max_file_size" yaml:"max_file_size"`
	MaxBackups       int           `json:"max_backups" yaml:"max_backups"`
	
	// Security settings
	EncryptLogs      bool          `json:"encrypt_logs" yaml:"encrypt_logs"`
	EncryptionKey    string        `json:"encryption_key,omitempty" yaml:"encryption_key,omitempty"`
	SignLogs         bool          `json:"sign_logs" yaml:"sign_logs"`
	PrivateKeyPath   string        `json:"private_key_path,omitempty" yaml:"private_key_path,omitempty"`
	PublicKeyPath    string        `json:"public_key_path,omitempty" yaml:"public_key_path,omitempty"`
	
	// Performance settings
	AsyncLogging     bool          `json:"async_logging" yaml:"async_logging"`
	QueueSize        int           `json:"queue_size" yaml:"queue_size"`
	WorkerCount      int           `json:"worker_count" yaml:"worker_count"`
	
	// Environment-specific settings
	Environment      string        `json:"environment" yaml:"environment"`
	Debug            bool          `json:"debug" yaml:"debug"`
}

// FileOutputConfig represents file output configuration
type FileOutputConfig struct {
	Path            string        `json:"path" yaml:"path"`
	Format          string        `json:"format" yaml:"format"`
	Rotate          bool          `json:"rotate" yaml:"rotate"`
	Compress        bool          `json:"compress" yaml:"compress"`
	Permissions     int           `json:"permissions" yaml:"permissions"`
}

// ConsoleOutputConfig represents console output configuration
type ConsoleOutputConfig struct {
	Format          string        `json:"format" yaml:"format"`
	Color           bool          `json:"color" yaml:"color"`
	Structured      bool          `json:"structured" yaml:"structured"`
}

// HTTPOutputConfig represents HTTP output configuration
type HTTPOutputConfig struct {
	URL             string        `json:"url" yaml:"url"`
	Method          string        `json:"method" yaml:"method"`
	Headers         map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	Timeout         time.Duration `json:"timeout" yaml:"timeout"`
	RetryCount      int           `json:"retry_count" yaml:"retry_count"`
	RetryDelay      time.Duration `json:"retry_delay" yaml:"retry_delay"`
	BasicAuth       *BasicAuthConfig `json:"basic_auth,omitempty" yaml:"basic_auth,omitempty"`
	BearerToken     string        `json:"bearer_token,omitempty" yaml:"bearer_token,omitempty"`
}

// BasicAuthConfig represents basic authentication configuration
type BasicAuthConfig struct {
	Username        string        `json:"username" yaml:"username"`
	Password        string        `json:"password" yaml:"password"`
}

// SyslogOutputConfig represents syslog output configuration
type SyslogOutputConfig struct {
	Network         string        `json:"network" yaml:"network"` // tcp, udp, unix
	Address         string        `json:"address" yaml:"address"`
	Tag             string        `json:"tag" yaml:"tag"`
	Facility        string        `json:"facility" yaml:"facility"`
	Severity        string        `json:"severity" yaml:"severity"`
	Format          string        `json:"format" yaml:"format"` // json, text, csv
}

// DefaultSecurityLoggingConfig returns the default configuration for security event logging
func DefaultSecurityLoggingConfig() *SecurityLoggingConfig {
	return &SecurityLoggingConfig{
		Enabled:       true,
		Level:         "medium",
		Format:        "json",
		BufferSize:    1000,
		FlushInterval: 5 * time.Second,
		BatchSize:     100,
		AsyncLogging:  true,
		QueueSize:     10000,
		WorkerCount:   3,
		Environment:   "development",
		Debug:         false,
		ConsoleOutput: &ConsoleOutputConfig{
			Format:     "json",
			Color:      true,
			Structured: true,
		},
	}
}

// Validate validates the security logging configuration
func (c *SecurityLoggingConfig) Validate() error {
	// Basic validation
	if c.Level == "" {
		c.Level = "medium"
	}
	
	if c.Format == "" {
		c.Format = "json"
	}
	
	if c.BufferSize <= 0 {
		c.BufferSize = 1000
	}
	
	if c.FlushInterval <= 0 {
		c.FlushInterval = 5 * time.Second
	}
	
	if c.BatchSize <= 0 {
		c.BatchSize = 100
	}
	
	if c.QueueSize <= 0 {
		c.QueueSize = 10000
	}
	
	if c.WorkerCount <= 0 {
		c.WorkerCount = 3
	}
	
	// Validate file output configuration
	if c.FileOutput != nil && c.FileOutput.Path == "" {
		// If file output is configured, path is required
		return fmt.Errorf("file output path is required")
	}
	
	// Validate HTTP output configuration
	if c.HTTPOutput != nil && c.HTTPOutput.URL == "" {
		// If HTTP output is configured, URL is required
		return fmt.Errorf("HTTP output URL is required")
	}
	
	// Validate syslog output configuration
	if c.SyslogOutput != nil && c.SyslogOutput.Address == "" {
		// If syslog output is configured, address is required
		return fmt.Errorf("syslog output address is required")
	}
	
	// Validate encryption settings
	if c.EncryptLogs && c.EncryptionKey == "" {
		return fmt.Errorf("encryption key is required when encrypt_logs is enabled")
	}
	
	if c.SignLogs && (c.PrivateKeyPath == "" || c.PublicKeyPath == "") {
		return fmt.Errorf("private and public key paths are required when sign_logs is enabled")
	}
	
	return nil
}