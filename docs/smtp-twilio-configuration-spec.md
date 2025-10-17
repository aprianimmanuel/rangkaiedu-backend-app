# SMTP and Twilio Configuration Specification

## Overview

This document provides a comprehensive technical specification for configuring and managing SMTP email services and Twilio SMS services in the Rangkai Edu backend application. It details the current implementation, identifies areas for improvement, and provides recommendations for a more robust and secure configuration system.

## Current Implementation Analysis

### Email Service (SMTP)

The current email implementation is located in [`utils/email/email.go`](utils/email/email.go:1-43) and uses Go's built-in `net/smtp` package:

1. **Configuration**: SMTP settings are loaded from environment variables:
   - `SMTP_HOST`: SMTP server host
   - `SMTP_PORT`: SMTP server port (default: 587)
   - `SMTP_USER`: SMTP username
   - `SMTP_PASS`: SMTP password

2. **Implementation Details**:
   - Uses `smtp.PlainAuth` for authentication
   - Sends plain text emails with basic templating
   - Gracefully skips sending if SMTP is not configured
   - Basic error handling with wrapped errors

3. **Usage**:
   - Used for sending OTP codes during registration and authentication
   - Called from the authentication controller in [`controllers/auth_controller.go`](controllers/auth_controller.go:137)

### SMS Service (Twilio)

The current SMS implementation is located in [`utils/sms/sms.go`](utils/sms/sms.go:1-36) and uses the Twilio Go SDK:

1. **Configuration**: Twilio settings are loaded from environment variables:
   - `TWILIO_ACCOUNT_SID`: Twilio Account SID
   - `TWILIO_AUTH_TOKEN`: Twilio Auth Token
   - `TWILIO_SENDER_PHONE`: Twilio sender phone number

2. **Implementation Details**:
   - Uses official Twilio Go client library
   - Sends SMS messages with basic templating
   - Checks for Twilio API errors
   - Basic error handling with wrapped errors

3. **Usage**:
   - Used for sending OTP codes as an alternative to email
   - Called from the authentication controller in [`controllers/auth_controller.go`](controllers/auth_controller.go:233)

### Configuration Management

The configuration system is implemented in [`config/config.go`](config/config.go:1-157):

1. **Environment Variable Loading**:
   - Uses `github.com/joho/godotenv` to load `.env` files
   - Falls back to default values for unset variables
   - All configuration values are strings

2. **Credential Storage**:
   - Credentials are stored as plain environment variables
   - No encryption or additional security measures
   - No credential rotation mechanisms

3. **Validation**:
   - Basic validation for required database fields
   - Service credentials are validated when first used
   - No proactive validation of service configurations

## Technical Specification

### 1. SMTP Service Configuration

#### Configuration Structure

```go
type SMTPConfig struct {
    Provider     string `json:"provider"`     // e.g., "gmail", "sendgrid", "custom"
    Host         string `json:"host"`
    Port         int    `json:"port"`
    Username     string `json:"username"`
    Password     string `json:"password"`     // Should be encrypted/securely stored
    FromEmail    string `json:"from_email"`
    FromName     string `json:"from_name"`
    UseTLS       bool   `json:"use_tls"`
    UseStartTLS  bool   `json:"use_starttls"`
    Enabled      bool   `json:"enabled"`
}
```

#### Security Requirements

1. **Credential Management**:
   - Passwords should be encrypted at rest
   - Support for credential rotation without application restart
   - Integration with secret management systems (HashiCorp Vault, AWS Secrets Manager)

2. **Connection Security**:
   - Enforce TLS connections by default
   - Certificate validation for SMTP servers
   - Support for custom CA certificates

3. **Authentication**:
   - Support for OAuth2 authentication (for providers like Gmail)
   - Support for API key authentication (for providers like SendGrid)
   - Fallback to basic authentication

#### Error Handling

1. **Connection Errors**:
   - Retry mechanism with exponential backoff
   - Circuit breaker pattern to prevent cascading failures
   - Detailed error logging for debugging

2. **Delivery Errors**:
   - Bounce handling and notification
   - Retry for transient failures
   - Permanent failure notification

3. **Monitoring**:
   - Delivery success/failure metrics
   - Latency tracking
   - Alerting for delivery failures

### 2. Twilio Service Configuration

#### Configuration Structure

```go
type TwilioConfig struct {
    Provider        string `json:"provider"`        // Always "twilio" for now
    AccountSID      string `json:"account_sid"`
    AuthToken       string `json:"auth_token"`      // Should be encrypted/securely stored
    FromPhoneNumber string `json:"from_phone_number"`
    Enabled         bool   `json:"enabled"`
}
```

#### Security Requirements

1. **Credential Management**:
   - Auth tokens should be encrypted at rest
   - Support for credential rotation without application restart
   - Integration with secret management systems

2. **API Security**:
   - Validate SSL certificates
   - Support for custom CA certificates
   - Rate limiting to prevent abuse

#### Error Handling

1. **API Errors**:
   - Parse Twilio error codes and messages
   - Retry mechanism for transient errors
   - Circuit breaker pattern

2. **Delivery Errors**:
   - Handle undelivered messages
   - Bounce handling and notification
   - Retry for transient failures

3. **Monitoring**:
   - Delivery success/failure metrics
   - Latency tracking
   - Cost tracking (Twilio charges per message)

### 3. Multi-Provider Support

#### Email Providers

Support for multiple email providers with different configuration requirements:

1. **Gmail/Google Workspace**:
   - OAuth2 authentication
   - App passwords for legacy support

2. **SendGrid**:
   - API key authentication
   - Web API integration

3. **Amazon SES**:
   - AWS credentials
   - Region-specific endpoints

4. **Custom SMTP**:
   - Full SMTP configuration
   - TLS/STARTTLS support

#### SMS Providers

Support for multiple SMS providers:

1. **Twilio** (current implementation):
   - Account SID and Auth Token
   - Phone number management

2. **Amazon SNS**:
   - AWS credentials
   - Region-specific endpoints

3. **Nexmo/Vonage**:
   - API key and secret
   - Branding and sender ID support

### 4. Secure Credential Management

#### Encryption at Rest

1. **Local Development**:
   - Encrypt sensitive values in `.env` files
   - Use development-specific encryption keys

2. **Production**:
   - Integration with enterprise secret management systems
   - Automatic decryption at runtime
   - Key rotation support

#### Credential Rotation

1. **Hot Reloading**:
   - Reload configuration without application restart
   - Graceful transition between old and new credentials

2. **Rotation Notifications**:
   - Alert when credentials are nearing expiration
   - Automated rotation where possible

### 5. Environment-Specific Configuration

#### Configuration Hierarchy

1. **Local Development**:
   - Encrypted `.env` files
   - Sample configurations in documentation
   - Disabled services by default

2. **Staging**:
   - Environment variables
   - Test credentials
   - Limited functionality

3. **Production**:
   - Secret management systems
   - Real credentials
   - Full functionality with monitoring

### 6. Monitoring and Logging

#### Email Service Monitoring

1. **Delivery Metrics**:
   - Success rate tracking
   - Delivery latency
   - Bounce rates

2. **Error Tracking**:
   - SMTP error categorization
   - Retry attempt logging
   - Failure notifications

#### SMS Service Monitoring

1. **Delivery Metrics**:
   - Success rate tracking
   - Delivery latency
   - Cost tracking

2. **Error Tracking**:
   - Twilio error code categorization
   - Retry attempt logging
   - Failure notifications

### 7. Configuration Validation

#### Schema Validation

1. **Required Fields**:
   - Validate all required configuration fields
   - Clear error messages for missing values

2. **Format Validation**:
   - Email format validation
   - Phone number format validation
   - URL format validation

#### Service Validation

1. **Connectivity Testing**:
   - Test SMTP connection at startup
   - Test Twilio API connectivity at startup
   - Graceful degradation for failed services

2. **Credential Validation**:
   - Validate credentials without sending messages
   - Alert on invalid credentials

## Implementation Plan

### Phase 1: Configuration Structure Enhancement

1. **Refactor Configuration Loading**:
   - Create provider-specific configuration structures
   - Implement configuration validation
   - Add support for multiple providers

2. **Secure Credential Storage**:
   - Implement encryption for sensitive values
   - Add support for secret management systems
   - Create credential rotation mechanisms

### Phase 2: Service Implementation Enhancement

1. **Email Service Improvements**:
   - Add support for multiple email providers
   - Implement retry mechanisms
   - Add circuit breaker pattern

2. **SMS Service Improvements**:
   - Add support for multiple SMS providers
   - Implement retry mechanisms
   - Add circuit breaker pattern

### Phase 3: Monitoring and Observability

1. **Metrics Collection**:
   - Implement delivery metrics
   - Add latency tracking
   - Create cost tracking for SMS

2. **Error Handling**:
   - Enhance error categorization
   - Implement alerting mechanisms
   - Add detailed logging

## Security Considerations

### Credential Security

1. **Storage**:
   - Never store plain text credentials in version control
   - Encrypt credentials at rest
   - Use secure secret management systems in production

2. **Transmission**:
   - Always use encrypted connections
   - Validate SSL certificates
   - Support for custom CA certificates

### Access Control

1. **Configuration Access**:
   - Limit access to configuration files
   - Audit logging for configuration access
   - Role-based access control

2. **Runtime Security**:
   - Validate configuration changes
   - Prevent unauthorized configuration updates
   - Secure credential injection

## Future Considerations

### Scalability

1. **High Availability**:
   - Support for multiple SMTP servers
   - Load balancing between providers
   - Failover mechanisms

2. **Performance**:
   - Connection pooling
   - Asynchronous message sending
   - Batch processing capabilities

### Extensibility

1. **Plugin Architecture**:
   - Support for custom email providers
   - Support for custom SMS providers
   - Easy integration of new services

2. **Configuration Management**:
   - Dynamic configuration updates
   - Configuration versioning
   - Rollback capabilities