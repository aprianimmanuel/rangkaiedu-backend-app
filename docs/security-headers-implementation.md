# Security Headers Implementation Specification

## Overview

This document outlines the technical specification for implementing comprehensive security headers in the Rangkai Edu backend application. The implementation will enhance the application's security posture by adding industry-standard HTTP security headers to all HTTP responses.

## Current State Analysis

### Existing Security Implementations

1. **HSTS (HTTP Strict Transport Security)**:
   - Partially implemented in `utils/https/https.go`
   - Configurable through `config/https.go`
   - Applied via middleware in `main.go`

2. **Authentication & Authorization**:
   - JWT-based authentication in `middleware/auth.go`
   - Role-based access control in `middleware/roles.go`

3. **HTTPS Enforcement**:
   - Configurable HTTPS settings in `config/https.go`
   - Automatic HTTP to HTTPS redirection

### Identified Gaps

1. **Missing Security Headers**:
   - Content Security Policy (CSP)
   - X-Frame-Options
   - X-XSS-Protection
   - X-Content-Type-Options
   - Referrer-Policy
   - Permissions-Policy

2. **Configuration Limitations**:
   - No centralized configuration for security headers
   - No environment-specific security header settings
   - No validation of security header values

## Security Headers Implementation Specification

### 1. Content Security Policy (CSP)

**Header**: `Content-Security-Policy`

**Purpose**: Mitigate cross-site scripting (XSS) and data injection attacks by specifying allowed sources of content.

**Default Configuration**:
```
default-src 'self';
script-src 'self' 'unsafe-inline' 'unsafe-eval';
style-src 'self' 'unsafe-inline';
img-src 'self' data: https:;
font-src 'self';
connect-src 'self';
media-src 'self';
object-src 'none';
child-src 'self';
frame-ancestors 'none';
base-uri 'self';
form-action 'self';
```

### 2. X-Frame-Options

**Header**: `X-Frame-Options`

**Purpose**: Prevent clickjacking attacks by controlling whether the browser should allow the page to be displayed in a frame.

**Default Configuration**: `DENY`

### 3. X-XSS-Protection

**Header**: `X-XSS-Protection`

**Purpose**: Enable XSS filtering in browsers that support it.

**Default Configuration**: `1; mode=block`

### 4. X-Content-Type-Options

**Header**: `X-Content-Type-Options`

**Purpose**: Prevent MIME type sniffing by instructing the browser to use the declared content type.

**Default Configuration**: `nosniff`

### 5. Referrer-Policy

**Header**: `Referrer-Policy`

**Purpose**: Control how much referrer information is included with requests.

**Default Configuration**: `strict-origin-when-cross-origin`

### 6. Permissions-Policy

**Header**: `Permissions-Policy`

**Purpose**: Control which features and APIs can be used in the browser.

**Default Configuration**:
```
geolocation=(), microphone=(), camera=(), fullscreen=(self), payment=()
```

## Configuration Structure

### SecurityHeadersConfig

```go
// SecurityHeadersConfig represents the configuration for security headers
type SecurityHeadersConfig struct {
    // Enabled enables security headers
    Enabled bool `json:"enabled" mapstructure:"enabled"`
    
    // ContentSecurityPolicy configuration
    ContentSecurityPolicy CSPConfig `json:"content_security_policy" mapstructure:"content_security_policy"`
    
    // XFrameOptions configuration
    XFrameOptions XFrameOptionsConfig `json:"x_frame_options" mapstructure:"x_frame_options"`
    
    // XXSSProtection configuration
    XXSSProtection XXSSProtectionConfig `json:"x_xss_protection" mapstructure:"x_xss_protection"`
    
    // XContentTypeOptions configuration
    XContentTypeOptions XContentTypeOptionsConfig `json:"x_content_type_options" mapstructure:"x_content_type_options"`
    
    // ReferrerPolicy configuration
    ReferrerPolicy ReferrerPolicyConfig `json:"referrer_policy" mapstructure:"referrer_policy"`
    
    // PermissionsPolicy configuration
    PermissionsPolicy PermissionsPolicyConfig `json:"permissions_policy" mapstructure:"permissions_policy"`
}
```

### CSPConfig

```go
// CSPConfig represents Content Security Policy configuration
type CSPConfig struct {
    // Enabled enables CSP
    Enabled bool `json:"enabled" mapstructure:"enabled"`
    
    // Policy is the CSP policy string
    Policy string `json:"policy" mapstructure:"policy"`
    
    // ReportOnly enables report-only mode
    ReportOnly bool `json:"report_only" mapstructure:"report_only"`
    
    // ReportURI is the URI to send violation reports to
    ReportURI string `json:"report_uri" mapstructure:"report_uri"`
}
```

### XFrameOptionsConfig

```go
// XFrameOptionsConfig represents X-Frame-Options configuration
type XFrameOptionsConfig struct {
    // Enabled enables X-Frame-Options
    Enabled bool `json:"enabled" mapstructure:"enabled"`
    
    // Value is the X-Frame-Options value (DENY, SAMEORIGIN)
    Value string `json:"value" mapstructure:"value"`
}
```

### XXSSProtectionConfig

```go
// XXSSProtectionConfig represents X-XSS-Protection configuration
type XXSSProtectionConfig struct {
    // Enabled enables X-XSS-Protection
    Enabled bool `json:"enabled" mapstructure:"enabled"`
    
    // Value is the X-XSS-Protection value
    Value string `json:"value" mapstructure:"value"`
}
```

### XContentTypeOptionsConfig

```go
// XContentTypeOptionsConfig represents X-Content-Type-Options configuration
type XContentTypeOptionsConfig struct {
    // Enabled enables X-Content-Type-Options
    Enabled bool `json:"enabled" mapstructure:"enabled"`
    
    // Value is the X-Content-Type-Options value
    Value string `json:"value" mapstructure:"value"`
}
```

### ReferrerPolicyConfig

```go
// ReferrerPolicyConfig represents Referrer-Policy configuration
type ReferrerPolicyConfig struct {
    // Enabled enables Referrer-Policy
    Enabled bool `json:"enabled" mapstructure:"enabled"`
    
    // Value is the Referrer-Policy value
    Value string `json:"value" mapstructure:"value"`
}
```

### PermissionsPolicyConfig

```go
// PermissionsPolicyConfig represents Permissions-Policy configuration
type PermissionsPolicyConfig struct {
    // Enabled enables Permissions-Policy
    Enabled bool `json:"enabled" mapstructure:"enabled"`
    
    // Value is the Permissions-Policy value
    Value string `json:"value" mapstructure:"value"`
}
```

## Default Configuration Values

```go
// GetDefaultSecurityHeadersConfig returns the default security headers configuration
func GetDefaultSecurityHeadersConfig() SecurityHeadersConfig {
    return SecurityHeadersConfig{
        Enabled: true,
        ContentSecurityPolicy: CSPConfig{
            Enabled: true,
            Policy: "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; media-src 'self'; object-src 'none'; child-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self';",
            ReportOnly: false,
        },
        XFrameOptions: XFrameOptionsConfig{
            Enabled: true,
            Value:   "DENY",
        },
        XXSSProtection: XXSSProtectionConfig{
            Enabled: true,
            Value:   "1; mode=block",
        },
        XContentTypeOptions: XContentTypeOptionsConfig{
            Enabled: true,
            Value:   "nosniff",
        },
        ReferrerPolicy: ReferrerPolicyConfig{
            Enabled: true,
            Value:   "strict-origin-when-cross-origin",
        },
        PermissionsPolicy: PermissionsPolicyConfig{
            Enabled: true,
            Value:   "geolocation=(), microphone=(), camera=(), fullscreen=(self), payment=()",
        },
    }
}
```

## Integration with Existing Middleware System

### Security Headers Middleware

A new middleware function will be created to apply security headers to all responses:

```go
// SecurityHeadersMiddleware adds security headers to responses
func SecurityHeadersMiddleware(cfg *config.SecurityHeadersConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        if cfg.Enabled {
            // Apply Content Security Policy
            if cfg.ContentSecurityPolicy.Enabled {
                headerName := "Content-Security-Policy"
                if cfg.ContentSecurityPolicy.ReportOnly {
                    headerName = "Content-Security-Policy-Report-Only"
                }
                c.Header(headerName, cfg.ContentSecurityPolicy.Policy)
                if cfg.ContentSecurityPolicy.ReportURI != "" {
                    c.Header("Content-Security-Policy-Report-Only", cfg.ContentSecurityPolicy.Policy+"; report-uri "+cfg.ContentSecurityPolicy.ReportURI)
                }
            }
            
            // Apply X-Frame-Options
            if cfg.XFrameOptions.Enabled {
                c.Header("X-Frame-Options", cfg.XFrameOptions.Value)
            }
            
            // Apply X-XSS-Protection
            if cfg.XXSSProtection.Enabled {
                c.Header("X-XSS-Protection", cfg.XXSSProtection.Value)
            }
            
            // Apply X-Content-Type-Options
            if cfg.XContentTypeOptions.Enabled {
                c.Header("X-Content-Type-Options", cfg.XContentTypeOptions.Value)
            }
            
            // Apply Referrer-Policy
            if cfg.ReferrerPolicy.Enabled {
                c.Header("Referrer-Policy", cfg.ReferrerPolicy.Value)
            }
            
            // Apply Permissions-Policy
            if cfg.PermissionsPolicy.Enabled {
                c.Header("Permissions-Policy", cfg.PermissionsPolicy.Value)
            }
        }
        
        c.Next()
    }
}
```

### Integration Points

1. **Main Application Setup**:
   - Add security headers middleware to the Gin engine in `main.go`
   - Place after HSTS middleware but before route definitions

2. **Configuration Loading**:
   - Extend `config.Load()` to include security headers configuration
   - Add environment variable parsing for security headers

3. **HTTPS Integration**:
   - Coordinate with existing HSTS implementation
   - Ensure consistent security header application across HTTP/HTTPS

## Environment-Specific Security Configurations

### Configuration by Environment

Security headers can be configured differently for each environment:

1. **Development**:
   - Relaxed CSP for easier debugging
   - Report-only mode for CSP violations
   - Detailed error reporting

2. **Staging**:
   - Production-like security headers
   - CSP reporting enabled
   - Monitoring of security violations

3. **Production**:
   - Strictest security headers
   - No report-only modes
   - Minimal error information disclosure

### Environment Variable Configuration

Security headers will be configurable through environment variables:

```
# Security Headers Configuration
SECURITY_HEADERS_ENABLED=true

# Content Security Policy
SECURITY_HEADERS_CSP_ENABLED=true
SECURITY_HEADERS_CSP_POLICY="default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data: https:;"
SECURITY_HEADERS_CSP_REPORT_ONLY=false
SECURITY_HEADERS_CSP_REPORT_URI=""

# X-Frame-Options
SECURITY_HEADERS_X_FRAME_OPTIONS_ENABLED=true
SECURITY_HEADERS_X_FRAME_OPTIONS_VALUE="DENY"

# X-XSS-Protection
SECURITY_HEADERS_X_XSS_PROTECTION_ENABLED=true
SECURITY_HEADERS_X_XSS_PROTECTION_VALUE="1; mode=block"

# X-Content-Type-Options
SECURITY_HEADERS_X_CONTENT_TYPE_OPTIONS_ENABLED=true
SECURITY_HEADERS_X_CONTENT_TYPE_OPTIONS_VALUE="nosniff"

# Referrer-Policy
SECURITY_HEADERS_REFERRER_POLICY_ENABLED=true
SECURITY_HEADERS_REFERRER_POLICY_VALUE="strict-origin-when-cross-origin"

# Permissions-Policy
SECURITY_HEADERS_PERMISSIONS_POLICY_ENABLED=true
SECURITY_HEADERS_PERMISSIONS_POLICY_VALUE="geolocation=(), microphone=(), camera=()"
```

## Error Handling Strategy

### Configuration Validation

1. **Validation at Startup**:
   - Validate security header configuration values
   - Log warnings for potentially unsafe configurations
   - Fail fast for critical configuration errors

2. **Runtime Error Handling**:
   - Graceful degradation if security headers cannot be applied
   - Logging of security header application failures
   - Monitoring of security header violations

### Error Recovery

1. **Fallback Mechanisms**:
   - Use default security headers if configuration is invalid
   - Continue serving requests even if security headers fail to apply
   - Alert administrators of security header configuration issues

2. **Monitoring and Alerting**:
   - Log security header violations
   - Monitor CSP violation reports
   - Alert on misconfigured security headers

## Implementation Plan

### Phase 1: Core Implementation

1. Create security headers configuration structures
2. Implement security headers middleware
3. Integrate with existing configuration system
4. Add environment variable support

### Phase 2: Testing and Validation

1. Unit tests for security headers middleware
2. Integration tests with existing middleware
3. Security header validation tests
4. Environment-specific configuration tests

### Phase 3: Documentation and Deployment

1. Update documentation with security headers configuration
2. Create deployment guide for different environments
3. Monitor security headers in production
4. Gather feedback and iterate

## Security Considerations

### Potential Risks

1. **Overly Restrictive Headers**:
   - May break legitimate functionality
   - Requires thorough testing

2. **Misconfiguration**:
   - Could weaken security instead of strengthening it
   - Requires validation and monitoring

3. **Browser Compatibility**:
   - Some headers may not be supported by all browsers
   - Requires fallback mechanisms

### Mitigation Strategies

1. **Gradual Rollout**:
   - Start with report-only modes
   - Monitor violations before enforcing

2. **Comprehensive Testing**:
   - Test across different browsers
   - Validate with security scanning tools

3. **Monitoring and Alerting**:
   - Implement violation reporting
   - Set up alerts for security issues

## Performance Impact

### Expected Impact

1. **Minimal Response Time Increase**:
   - Additional headers add negligible overhead
   - No complex processing required

2. **Memory Usage**:
   - Configuration structures are lightweight
   - Minimal impact on application memory

### Optimization Considerations

1. **Header Caching**:
   - Pre-compute header values where possible
   - Avoid repeated string operations

2. **Conditional Application**:
   - Skip header application for static assets if needed
   - Optimize for high-traffic scenarios

## Monitoring and Maintenance

### Key Metrics

1. **Security Header Compliance**:
   - Percentage of responses with proper security headers
   - Violation reports from browsers

2. **Configuration Health**:
   - Validity of security header configurations
   - Environment-specific configuration accuracy

### Maintenance Tasks

1. **Regular Updates**:
   - Review and update security header values
   - Stay current with security best practices

2. **Security Audits**:
   - Periodic review of security header effectiveness
   - Update based on new threats and vulnerabilities

## Conclusion

This security headers implementation will significantly enhance the security posture of the Rangkai Edu backend application by adding industry-standard HTTP security headers. The implementation is designed to be flexible, configurable, and maintainable while providing strong security protections against common web vulnerabilities.