package config_test

import (
	"os"
	"testing"
	
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
)

func TestLoadHTTPSConfig(t *testing.T) {
	// Test default configuration
	cfg := config.LoadHTTPSConfig()
	
	if cfg.Enabled != false {
		t.Errorf("Expected default HTTPS enabled to be false, got %v", cfg.Enabled)
	}
	
	if cfg.Port != 8443 {
		t.Errorf("Expected default HTTPS port to be 8443, got %d", cfg.Port)
	}
	
	if cfg.HTTPPort != 8080 {
		t.Errorf("Expected default HTTP port to be 8080, got %d", cfg.HTTPPort)
	}
	
	// Test with environment variables
	os.Setenv("HTTPS_ENABLED", "true")
	os.Setenv("HTTPS_PORT", "443")
	os.Setenv("HTTP_PORT", "80")
	os.Setenv("HTTPS_CERTIFICATE_FILE_CERT_FILE", "/test/cert.pem")
	os.Setenv("HTTPS_CERTIFICATE_FILE_KEY_FILE", "/test/key.pem")
	
	// Reload configuration
	cfg = config.LoadHTTPSConfig()
	
	if cfg.Enabled != true {
		t.Errorf("Expected HTTPS enabled to be true, got %v", cfg.Enabled)
	}
	
	if cfg.Port != 443 {
		t.Errorf("Expected HTTPS port to be 443, got %d", cfg.Port)
	}
	
	if cfg.HTTPPort != 80 {
		t.Errorf("Expected HTTP port to be 80, got %d", cfg.HTTPPort)
	}
	
	if cfg.Certificate.File.CertFile != "/test/cert.pem" {
		t.Errorf("Expected cert file to be /test/cert.pem, got %s", cfg.Certificate.File.CertFile)
	}
	
	if cfg.Certificate.File.KeyFile != "/test/key.pem" {
		t.Errorf("Expected key file to be /test/key.pem, got %s", cfg.Certificate.File.KeyFile)
	}
	
	// Clean up environment variables
	os.Unsetenv("HTTPS_ENABLED")
	os.Unsetenv("HTTPS_PORT")
	os.Unsetenv("HTTP_PORT")
	os.Unsetenv("HTTPS_CERTIFICATE_FILE_CERT_FILE")
	os.Unsetenv("HTTPS_CERTIFICATE_FILE_KEY_FILE")
}

func TestCreateTLSConfig(t *testing.T) {
	config := config.GetDefaultHTTPSConfig()
	
	tlsConfig, err := config.CreateTLSConfig()
	if err != nil {
		t.Fatalf("Failed to create TLS config: %v", err)
	}
	
	if tlsConfig.MinVersion != 0x0303 { // TLS 1.2
		t.Errorf("Expected minimum TLS version to be TLS 1.2, got %x", tlsConfig.MinVersion)
	}
	
	if !tlsConfig.PreferServerCipherSuites {
		t.Error("Expected PreferServerCipherSuites to be true")
	}
}

func TestHSTSConfig(t *testing.T) {
	config := config.GetDefaultHTTPSConfig()
	
	if !config.HSTS.Enabled {
		t.Error("Expected HSTS to be enabled by default")
	}
	
	if config.HSTS.MaxAge != 31536000 {
		t.Errorf("Expected HSTS max age to be 31536000, got %d", config.HSTS.MaxAge)
	}
	
	if !config.HSTS.IncludeSubDomains {
		t.Error("Expected HSTS include subdomains to be true by default")
	}
}