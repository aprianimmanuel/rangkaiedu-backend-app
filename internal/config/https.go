// Package config provides HTTPS configuration and TLS utilities
package config

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"strings"
)

// HTTPSConfig represents the basic HTTPS configuration
type HTTPSConfig struct {
	// Enabled enables HTTPS support
	Enabled bool `json:"enabled" mapstructure:"enabled"`
	
	// Port is the HTTPS port to listen on
	Port int `json:"port" mapstructure:"port"`
	
	// HTTPPort is the HTTP port for redirection (0 to disable)
	HTTPPort int `json:"http_port" mapstructure:"http_port"`
	
	// Certificate configuration
	Certificate CertificateConfig `json:"certificate" mapstructure:"certificate"`
	
	// TLS configuration
	TLS TLSConfig `json:"tls" mapstructure:"tls"`
	
	// HSTS configuration
	HSTS HSTSConfig `json:"hsts" mapstructure:"hsts"`
	
	// Redirect configuration
	Redirect RedirectConfig `json:"redirect" mapstructure:"redirect"`
}

// CertificateConfig represents certificate management configuration
type CertificateConfig struct {
	// Type of certificate management (file, letsencrypt)
	Type string `json:"type" mapstructure:"type"`
	
	// File-based certificate configuration
	File FileCertConfig `json:"file" mapstructure:"file"`
	
	// Let's Encrypt configuration
	LetsEncrypt LetsEncryptConfig `json:"lets_encrypt" mapstructure:"lets_encrypt"`
}

// FileCertConfig represents file-based certificate configuration
type FileCertConfig struct {
	// CertFile is the path to the certificate file
	CertFile string `json:"cert_file" mapstructure:"cert_file"`
	
	// KeyFile is the path to the private key file
	KeyFile string `json:"key_file" mapstructure:"key_file"`
	
	// CAFile is the path to the CA certificate file (optional)
	CAFile string `json:"ca_file" mapstructure:"ca_file"`
}

// LetsEncryptConfig represents Let's Encrypt configuration
type LetsEncryptConfig struct {
	// Email for Let's Encrypt account
	Email string `json:"email" mapstructure:"email"`
	
	// Domains to obtain certificates for
	Domains []string `json:"domains" mapstructure:"domains"`
	
	// Cache directory for certificates
	CacheDir string `json:"cache_dir" mapstructure:"cache_dir"`
}

// TLSConfig represents TLS configuration
type TLSConfig struct {
	// MinVersion is the minimum TLS version (TLS12, TLS13)
	MinVersion string `json:"min_version" mapstructure:"min_version"`
	
	// CipherSuites is a list of allowed cipher suites
	CipherSuites []string `json:"cipher_suites" mapstructure:"cipher_suites"`
	
	// PreferServerCipherSuites prefers server's cipher suite order
	PreferServerCipherSuites bool `json:"prefer_server_cipher_suites" mapstructure:"prefer_server_cipher_suites"`
}

// HSTSConfig represents HSTS configuration
type HSTSConfig struct {
	// Enabled enables HSTS
	Enabled bool `json:"enabled" mapstructure:"enabled"`
	
	// MaxAge is the max-age directive in seconds
	MaxAge int `json:"max_age" mapstructure:"max_age"`
	
	// IncludeSubDomains includes subdomains in HSTS policy
	IncludeSubDomains bool `json:"include_sub_domains" mapstructure:"include_sub_domains"`
	
	// Preload adds preload directive to HSTS header
	Preload bool `json:"preload" mapstructure:"preload"`
}

// RedirectConfig represents HTTP to HTTPS redirect configuration
type RedirectConfig struct {
	// Enabled enables HTTP to HTTPS redirection
	Enabled bool `json:"enabled" mapstructure:"enabled"`
	
	// StatusCode is the HTTP status code for redirects (301, 302, 307, 308)
	StatusCode int `json:"status_code" mapstructure:"status_code"`
	
	// Permanent enables permanent redirect (301)
	Permanent bool `json:"permanent" mapstructure:"permanent"`
}

// GetDefaultHTTPSConfig returns the default HTTPS configuration
func GetDefaultHTTPSConfig() HTTPSConfig {
	return HTTPSConfig{
		Enabled: false,
		Port:    8443,
		HTTPPort: 8080,
		Certificate: CertificateConfig{
			Type: "file",
			File: FileCertConfig{
				CertFile: "./certs/server.crt",
				KeyFile:  "./certs/server.key",
				CAFile:  "",
			},
		},
		TLS: TLSConfig{
			MinVersion:             "TLS12",
			CipherSuites:           []string{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"},
			PreferServerCipherSuites: true,
		},
		HSTS: HSTSConfig{
			Enabled:            true,
			MaxAge:             31536000, // 1 year
			IncludeSubDomains: true,
			Preload:           false,
		},
		Redirect: RedirectConfig{
			Enabled:    true,
			StatusCode:  301,
			Permanent:  true,
		},
	}
}

// LoadHTTPSConfig loads HTTPS configuration from environment variables
func LoadHTTPSConfig() HTTPSConfig {
	config := GetDefaultHTTPSConfig()
	
	// Load from environment variables
	if enabled := getEnvBool("HTTPS_ENABLED", false); enabled {
		config.Enabled = true
	}
	
	if port := getEnvInt("HTTPS_PORT", 8443); port > 0 {
		config.Port = port
	}
	
	if httpPort := getEnvInt("HTTP_PORT", 8080); httpPort > 0 {
		config.HTTPPort = httpPort
	}
	
	// Certificate configuration
	if certType := getProviderEnv("HTTPS_CERTIFICATE_TYPE", "file"); certType != "" {
		config.Certificate.Type = certType
	}
	
	// File-based certificate configuration
	if certFile := getProviderEnv("HTTPS_CERTIFICATE_FILE_CERT_FILE", "./certs/server.crt"); certFile != "" {
		config.Certificate.File.CertFile = certFile
	}
	
	if keyFile := getProviderEnv("HTTPS_CERTIFICATE_FILE_KEY_FILE", "./certs/server.key"); keyFile != "" {
		config.Certificate.File.KeyFile = keyFile
	}
	
	if caFile := getProviderEnv("HTTPS_CERTIFICATE_FILE_CA_FILE", ""); caFile != "" {
		config.Certificate.File.CAFile = caFile
	}
	
	// Let's Encrypt configuration
	if email := getProviderEnv("HTTPS_CERTIFICATE_LETSENCRYPT_EMAIL", ""); email != "" {
		config.Certificate.LetsEncrypt.Email = email
	}
	
	if domains := getProviderEnv("HTTPS_CERTIFICATE_LETSENCRYPT_DOMAINS", ""); domains != "" {
		config.Certificate.LetsEncrypt.Domains = strings.Split(domains, ",")
	}
	
	if cacheDir := getProviderEnv("HTTPS_CERTIFICATE_LETSENCRYPT_CACHE_DIR", "/var/cache/letsencrypt"); cacheDir != "" {
		config.Certificate.LetsEncrypt.CacheDir = cacheDir
	}
	
	// TLS configuration
	if minVersion := getProviderEnv("HTTPS_TLS_MIN_VERSION", "TLS12"); minVersion != "" {
		config.TLS.MinVersion = minVersion
	}
	
	if cipherSuites := getProviderEnv("HTTPS_TLS_CIPHER_SUITES", ""); cipherSuites != "" {
		config.TLS.CipherSuites = strings.Split(cipherSuites, ",")
	}
	
	if preferServerCipherSuites := getEnvBool("HTTPS_TLS_PREFER_SERVER_CIPHER_SUITES", true); preferServerCipherSuites {
		config.TLS.PreferServerCipherSuites = true
	}
	
	// HSTS configuration
	if hstsEnabled := getEnvBool("HTTPS_HSTS_ENABLED", true); hstsEnabled {
		config.HSTS.Enabled = true
	}
	
	if maxAge := getEnvInt("HTTPS_HSTS_MAX_AGE", 31536000); maxAge > 0 {
		config.HSTS.MaxAge = maxAge
	}
	
	if includeSubDomains := getEnvBool("HTTPS_HSTS_INCLUDE_SUBDOMAINS", true); includeSubDomains {
		config.HSTS.IncludeSubDomains = true
	}
	
	if preload := getEnvBool("HTTPS_HSTS_PRELOAD", false); preload {
		config.HSTS.Preload = true
	}
	
	// Redirect configuration
	if redirectEnabled := getEnvBool("HTTPS_REDIRECT_ENABLED", true); redirectEnabled {
		config.Redirect.Enabled = true
	}
	
	if statusCode := getEnvInt("HTTPS_REDIRECT_STATUS_CODE", 301); statusCode > 0 {
		config.Redirect.StatusCode = statusCode
	}
	
	if permanent := getEnvBool("HTTPS_REDIRECT_PERMANENT", true); permanent {
		config.Redirect.Permanent = true
	}
	
	return config
}

// CreateTLSConfig creates a TLS configuration from the HTTPS config
func (c *HTTPSConfig) CreateTLSConfig() (*tls.Config, error) {
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: c.TLS.PreferServerCipherSuites,
	}
	
	// Set minimum TLS version
	switch c.TLS.MinVersion {
	case "TLS12":
		tlsConfig.MinVersion = tls.VersionTLS12
	case "TLS13":
		tlsConfig.MinVersion = tls.VersionTLS13
	default:
		tlsConfig.MinVersion = tls.VersionTLS12
	}
	
	// Set cipher suites
	if len(c.TLS.CipherSuites) > 0 {
		cipherSuites := make([]uint16, 0, len(c.TLS.CipherSuites))
		for _, cipherName := range c.TLS.CipherSuites {
			if cipher, ok := cipherSuiteMap[cipherName]; ok {
				cipherSuites = append(cipherSuites, cipher)
			}
		}
		tlsConfig.CipherSuites = cipherSuites
	}
	
	// Load CA certificates if specified
	if c.Certificate.File.CAFile != "" {
		caCert, err := ioutil.ReadFile(c.Certificate.File.CAFile)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}
	
	return tlsConfig, nil
}

// GetCertificate gets the certificate based on the configuration
func (c *HTTPSConfig) GetCertificate() (*tls.Certificate, error) {
	if c.Certificate.Type == "file" {
		cert, err := tls.LoadX509KeyPair(c.Certificate.File.CertFile, c.Certificate.File.KeyFile)
		if err != nil {
			return nil, err
		}
		return &cert, nil
	}
	
	// For now, return nil for other certificate types
	// This will be extended in future implementations
	return nil, nil
}

// Cipher suite mapping
var cipherSuiteMap = map[string]uint16{
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256": tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384": tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384": tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305": tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305": tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	"TLS_RSA_WITH_AES_128_GCM_SHA256": tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_RSA_WITH_AES_256_GCM_SHA384": tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
}