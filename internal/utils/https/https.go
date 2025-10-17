package https

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
	"github.com/gin-gonic/gin"
)

// Manager handles HTTPS configuration and server management
type Manager struct {
	config     *config.HTTPSConfig
	tlsConfig  *tls.Config
}

// NewManager creates a new HTTPS manager
func NewManager(cfg config.HTTPSConfig) (*Manager, error) {
	if !cfg.Enabled {
		return &Manager{config: &cfg}, nil
	}

	// Create TLS configuration
	tlsConfig, err := cfg.CreateTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create TLS config: %w", err)
	}

	return &Manager{
		config:     &cfg,
		tlsConfig:  tlsConfig,
	}, nil
}

// StartHTTPS starts the HTTPS server
func (m *Manager) StartHTTPS(handler http.Handler, port int) error {
	if !m.config.Enabled {
		return nil
	}

	// Get certificate
	cert, err := m.config.GetCertificate()
	if err != nil {
		return fmt.Errorf("failed to get certificate: %w", err)
	}

	// Configure TLS
	m.tlsConfig.Certificates = []tls.Certificate{*cert}

	// Create server
	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   handler,
		TLSConfig: m.tlsConfig,
	}

	// Start server
	log.Printf("Starting HTTPS server on port %d", port)
	return server.ListenAndServeTLS("", "")
}

// StartHTTPRedirect starts the HTTP to HTTPS redirect server
func (m *Manager) StartHTTPRedirect(handler http.Handler, httpPort, httpsPort int) error {
	if !m.config.Redirect.Enabled || httpPort <= 0 {
		return nil
	}

	// Create redirect handler
	redirectHandler := m.createRedirectHandler(httpsPort)

	// Create server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: redirectHandler,
	}

	// Start server
	log.Printf("Starting HTTP redirect server on port %d", httpPort)
	return server.ListenAndServe()
}

// Shutdown gracefully shuts down the servers
func (m *Manager) Shutdown(ctx context.Context) error {
	// This method would be used to shutdown servers, but we're not storing them
	// in this simplified implementation
	return nil
}

// createRedirectHandler creates a handler for HTTP to HTTPS redirection
func (m *Manager) createRedirectHandler(httpsPort int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the target HTTPS URL
		target := fmt.Sprintf("https://%s:%d%s", r.Host, httpsPort, r.URL.Path)
		if r.URL.RawQuery != "" {
			target += "?" + r.URL.RawQuery
		}

		// Set HSTS header if enabled
		if m.config.HSTS.Enabled {
			hstsValue := fmt.Sprintf("max-age=%d", m.config.HSTS.MaxAge)
			if m.config.HSTS.IncludeSubDomains {
				hstsValue += "; includeSubDomains"
			}
			if m.config.HSTS.Preload {
				hstsValue += "; preload"
			}
			w.Header().Set("Strict-Transport-Security", hstsValue)
		}

		// Perform redirect
		statusCode := m.config.Redirect.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusMovedPermanently
		}
		http.Redirect(w, r, target, statusCode)
	})
}

// HSTSMiddleware adds HSTS headers to responses
func HSTSMiddleware(cfg *config.HSTSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.Enabled {
			hstsValue := fmt.Sprintf("max-age=%d", cfg.MaxAge)
			if cfg.IncludeSubDomains {
				hstsValue += "; includeSubDomains"
			}
			if cfg.Preload {
				hstsValue += "; preload"
			}
			c.Header("Strict-Transport-Security", hstsValue)
		}
		c.Next()
	}
}

// IsHTTPSRequest checks if the request is over HTTPS
func IsHTTPSRequest(r *http.Request) bool {
	return r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
}

// GetScheme returns the request scheme (http or https)
func GetScheme(r *http.Request) string {
	if IsHTTPSRequest(r) {
		return "https"
	}
	return "http"
}