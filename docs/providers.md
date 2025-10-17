
# Provider Configuration System

## Overview

The provider configuration system is a comprehensive solution for managing multiple email and SMS providers in the Rangkai Edu backend. It supports multiple provider types, credential encryption, health monitoring, and failover capabilities.

## Features

- **Multiple Provider Support**: Configure multiple email and SMS providers
- **Priority-based Selection**: Primary and fallback providers based on priority
- **Credential Encryption**: Secure storage of sensitive credentials
- **Health Monitoring**: Automatic health checks for all providers
- **Failover Support**: Automatic switching to backup providers
- **CLI Management**: Command-line interface for provider management
- **Backward Compatibility**: Maintains compatibility with existing configurations

## Supported Providers

### Email Providers

1. **SMTP**
   - Standard SMTP protocol
   - Supports TLS and StartTLS
   - Customizable headers and settings

2. **SendGrid**
   - SendGrid API integration
   - Template support
   - Advanced analytics

3. **Gmail**
   - Gmail OAuth2 integration
   - Google Workspace support
   - Refresh token management

### SMS Providers

1. **Twilio**
   - Twilio API integration
   - Global SMS delivery
   - Advanced features

2. **Amazon SNS**
   - AWS SNS integration
   - Regional SMS delivery
   - Cost-effective solution

## Configuration Structure

### Environment Variables

The system uses environment variables to configure providers. Each provider is numbered and prefixed with either `EMAIL_PROVIDER_` or `SMS_PROVIDER_`.

#### Email Provider Configuration

```bash
# Primary SMTP Provider
EMAIL_PROVIDER_0_TYPE=smtp
EMAIL_PROVIDER_0_NAME=primary-smtp
EMAIL_PROVIDER_0_ENABLED=true
EMAIL_PROVIDER_0_PRIORITY=0
EMAIL_PROVIDER_0_SETTINGS_HOST=smtp.gmail.com
EMAIL_PROVIDER_0_SETTINGS_PORT=587
EMAIL_PROVIDER_0_SETTINGS_USERNAME=your_email@gmail.com
EMAIL_PROVIDER_0_SETTINGS_PASSWORD=your_app_password
EMAIL_PROVIDER_0_SETTINGS_FROM_EMAIL=your_email@gmail.com
EMAIL_PROVIDER_0_SETTINGS_FROM_NAME=Rangkai Edu
EMAIL_PROVIDER_0_SETTINGS_USE_TLS=true
EMAIL_PROVIDER_0_SETTINGS_USE_STARTTLS=true

# Alternative SendGrid Provider
EMAIL_PROVIDER_1_TYPE=sendgrid
EMAIL_PROVIDER_1_NAME=sendgrid-fallback
EMAIL_PROVIDER_1_ENABLED=false
EMAIL_PROVIDER_1_PRIORITY=1
EMAIL_PROVIDER_1_SETTINGS_API_KEY=your_sendgrid_api_key
EMAIL_PROVIDER_1_SETTINGS_FROM_EMAIL=noreply@yourdomain.com
EMAIL_PROVIDER_1_SETTINGS_FROM_NAME=Rangkai Edu
```

#### SMS Provider Configuration

```bash
# Primary Twilio Provider
SMS_PROVIDER_0_TYPE=twilio
SMS_PROVIDER_0_NAME=primary-twilio
SMS_PROVIDER_0_ENABLED=true
SMS_PROVIDER_0_PRIORITY=0
SMS_PROVIDER_0_SETTINGS_ACCOUNT_SID=your_twilio_account_sid
SMS_PROVIDER_0_SETTINGS_AUTH_TOKEN=your_twilio_auth_token
SMS_PROVIDER_0_SETTINGS_FROM_NUMBER=+1234567890

# Alternative SNS Provider
SMS_PROVIDER_1_TYPE=sns
SMS_PROVIDER_1_NAME=sns-fallback
SMS_PROVIDER_1_ENABLED=false
SMS_PROVIDER_1_PRIORITY=1
SMS_PROVIDER_1_SETTINGS_ACCESS_KEY_ID=your_aws_access_key_id
SMS_PROVIDER_1_SETTINGS_SECRET_ACCESS_KEY=your_aws_secret_access_key
SMS_PROVIDER_1_SETTINGS_REGION=us-east-1
SMS_PROVIDER_1_SETTINGS_FROM_NUMBER=+1234567890
```

### Provider Settings

#### Email Provider Settings

| Setting | Description | Required |
|---------|-------------|----------|
| `SETTINGS_HOST` | SMTP server hostname | Yes (SMTP) |
| `SETTINGS_PORT` | SMTP server port | Yes (SMTP) |
| `SETTINGS_USERNAME` | SMTP username | Yes (SMTP) |
| `SETTINGS_PASSWORD` | SMTP password | Yes (SMTP) |
| `SETTINGS_FROM_EMAIL` | From email address | Yes |
| `SETTINGS_FROM_NAME` | From name | Yes |
| `SETTINGS_REPLY_TO` | Reply-to email | No |
| `SETTINGS_USE_TLS` | Use TLS encryption | No (default: false) |
| `SETTINGS_USE_STARTTLS` | Use StartTLS | No (default: false) |

#### SMS Provider Settings

| Setting | Description | Required |
|---------|-------------|----------|
| `SETTINGS_ACCOUNT_SID` | Twilio Account SID | Yes (Twilio) |
| `SETTINGS_AUTH_TOKEN` | Twilio Auth Token | Yes (Twilio) |
| `SETTINGS_ACCESS_KEY_ID` | AWS Access Key ID | Yes (SNS) |
| `SETTINGS_SECRET_ACCESS_KEY` | AWS Secret Access Key | Yes (SNS) |
| `SETTINGS_REGION` | AWS region | Yes (SNS) |
| `SETTINGS_FROM_NUMBER` | From phone number | Yes |

## Usage

### Basic Usage

```go
// Load configuration
cfg := config.Load()

// Get primary email provider
emailProvider, err := cfg.GetPrimaryEmailProvider()
if err != nil {
    log.Fatal(err)
}

// Get primary SMS provider
smsProvider, err := cfg.GetPrimarySMSProvider()
if err != nil {
    log.Fatal(err)
}

// Send OTP email
err = email.SendOTPEmail(cfg, "user@example.com", "123456")
if err != nil {
    log.Fatal(err)
}

// Send OTP SMS
err = sms.SendOTPSMS(cfg, "+1234567890", "123456")
if err != nil {
    log.Fatal(err)
}
```

### Advanced Usage

```go
// Get specific provider by name
emailProvider, err := cfg.GetEmailProviderByName("primary-smtp")
if err != nil {
    log.Fatal(err)
}

// Check provider health
health := cfg.GetProviderHealth()
for _, h := range health {
    fmt.Printf("Provider: %s, Status: %s\n", h.Name, h.Status)
}

// Validate all providers
err = cfg.ValidateProviders()
if err != nil {
    log.Fatal(err)
}
```

## CLI Management

The CLI tool provides comprehensive management capabilities for providers.

### Installation

```bash
go build -o providers cmd/providers/main.go
```

### Commands

#### List Providers

```bash
providers -list
```

#### Check Health Status

```bash
providers -health
```

#### Validate Configurations

```bash
providers -validate
```

#### Encrypt Credentials

```bash
providers -encrypt
```

#### Decrypt Credentials

```bash
providers -decrypt
```

#### Generate Configuration Template

```bash
providers -generate -type smtp -name my-smtp
```

#### View Provider Details

```bash
providers -type smtp -name primary-smtp
```

### CLI Examples

```bash
# List all configured providers
providers -list

# Check health status of all providers
providers -health

# Validate all provider configurations
providers -validate

# Encrypt all credentials
providers -encrypt

# Decrypt all credentials
providers -decrypt

# Generate SMTP provider configuration
providers -generate -type smtp -name my-smtp

# View SMTP provider details
providers -type smtp -name primary-smtp

# View Twilio provider details
providers -type twilio -name primary-twilio
```

## Security

### Credential Encryption

The system provides automatic encryption for sensitive credentials:

```go
// Encrypt credentials
err := cfg.EncryptCredential("sensitive-value")
if err != nil {
    log.Fatal(err)
}

// Decrypt credentials
decrypted, err := cfg.DecryptCredential("encrypted-value")
if err != nil {
    log.Fatal(err)
}

// Check if value is encrypted
isEncrypted := cfg.IsEncrypted("encrypted-value")
```

### Best Practices

1. **Use Environment Variables**: Never hardcode credentials in source code
2. **Enable Encryption**: Always encrypt sensitive credentials
3. **Use Secure Storage**: Store encrypted credentials in secure locations
4. **Regular Rotation**: Rotate credentials regularly
5. **Monitor Access**: Monitor access to sensitive configurations

## Health Monitoring

The system provides automatic health monitoring for all providers:

```go
// Get health status
health := cfg.GetProviderHealth()

// Health status structure
type ProviderHealth struct {
    Name      string    // Provider name
    Type      string    // Provider type (email/sms)
    Status    string    // Status (healthy/unhealthy)
    Message   string    // Status message
    LastTest  time.Time // Last test time
}
```

### Health Check Implementation

```go
// Custom health checker
type CustomHealthChecker struct {
    provider config.EmailProviderConfig
}

func (c *CustomHealthChecker) CheckHealth() error {
    // Implement custom health check logic
    return nil
}

// Register health checker
pm := providers.NewProviderManager(cfg)
pm.RegisterHealthChecker("custom", &CustomHealthChecker{})
```

## Failover Configuration

The system supports automatic failover to backup providers:

```bash
# Primary provider
EMAIL_PROVIDER_0_TYPE=smtp
EMAIL_PROVIDER_0_NAME=primary
EMAIL_PROVIDER_0_ENABLED=true
EMAIL_PROVIDER_0_PRIORITY=0

# Fallback provider
EMAIL_PROVIDER_1_TYPE=sendgrid
EMAIL_PROVIDER_1_NAME=fallback
EMAIL_PROVIDER_1_ENABLED=true
EMAIL_PROVIDER_1_PRIORITY=1
```

### Failover Logic

1. **Priority-based Selection**: Lower priority numbers are selected first
2. **Health-based Selection**: Only healthy providers are selected
3. **Automatic Switching**: System automatically switches to next available provider
4. **Error Handling**: Errors are logged but don't stop the process

## Migration Guide

### From Legacy Configuration

The system maintains backward compatibility with legacy configurations:

```bash
# Legacy SMTP configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASS=your_app_password

# Legacy Twilio configuration
TWILIO_ACCOUNT_SID=your_twilio_account_sid
TWILIO_AUTH_TOKEN=your_twilio_auth_token
TWILIOSENDER_PHONE=+1234567890
```

### Migration Steps

1. **Backup Existing Configuration**: Create a backup of existing environment variables
2. **Add New Provider Variables**: Add new provider configuration variables
3. **Test New Configuration**: Test the new configuration with the CLI tool
4. **Enable Encryption**: Enable credential encryption for security
5. **Update Code**: Update code to use new provider methods
6. **Remove Legacy Variables**: Remove legacy variables after successful migration

### Migration Example

```bash
# Step 1: Backup existing configuration
cp .env .env.backup

# Step 2: Add new provider configuration
echo "EMAIL_PROVIDER_0_TYPE=smtp" >> .env
echo "EMAIL_PROVIDER_0_NAME=primary-smtp" >> .env
echo "EMAIL_PROVIDER_0_ENABLED=true" >> .env
echo "EMAIL_PROVIDER_0_SETTINGS_HOST=smtp.gmail.com" >> .env
echo "EMAIL_PROVIDER_0_SETTINGS_USERNAME=your_email@gmail.com" >> .env
echo "EMAIL_PROVIDER_0_SETTINGS_PASSWORD=your_app_password" >> .env

# Step 3: Test new configuration
providers -validate

# Step 4: Enable encryption
providers -encrypt

# Step 5: Update code (see examples above)
# Step 6: Remove legacy variables
sed -i '/SMTP_/d' .env
sed -i '/TWILIO_/d' .env
```

## Troubleshooting

### Common Issues

1. **Provider Not Found**
   - Check if provider is enabled
   - Verify provider name and type
   - Check environment variable configuration

2. **Authentication Failed**
   - Verify credentials
   - Check if credentials are encrypted
   - Test credentials with provider's test tools

3. **Connection Failed**
   - Check network connectivity
   - Verify provider server status
   - Check firewall settings

4. **Configuration Validation Failed**
   - Check required fields
   - Verify data types
   - Check for syntax errors

### Debug Mode

Enable debug mode for detailed logging:

```bash
export PROVIDER_DEBUG=true
providers -health
```

### Log Analysis

Check logs for detailed error information:

```bash
# View provider logs
tail -f logs/providers.log

# View error logs
grep "ERROR" logs/providers.log
```

## Testing

### Unit Tests

```bash
# Run all tests
go test ./config/...

# Run specific test
go test ./config/ -run TestProviderManager_LoadProviders

# Run with coverage
go test ./config/ -cover
```

### Integration Tests

```bash
# Test email providers
go test ./utils/email/...

# Test SMS providers
go test ./utils/sms/...

# Test provider management
go test ./utils/providers/...
```

### Performance Testing

```bash
# Test provider performance
go test ./config/ -bench=.

# Test concurrent access
go test ./config/ -run TestConcurrentAccess
```

## API Reference

### Configuration Methods

- `Load() *Config`: Load configuration from environment variables
- `LoadTest() *Config`: Load test configuration
- `GetPrimaryEmailProvider() (EmailProviderConfig, error)`: Get primary email provider
- `GetPrimarySMSProvider() (SMSProviderConfig, error)`: Get primary SMS provider
- `GetEmailProviderByName(name string) (EmailProviderConfig, error)`: Get email provider by name
- `GetSMSProviderByName(name string) (SMSProviderConfig, error)`: Get SMS provider by name
- `ValidateProviders() error`: Validate all provider configurations
- `GetProviderHealth() []ProviderHealth`: Get provider health status

###

### Provider Configuration Methods

- `NewProviderManager() *ProviderManager`: Create new provider manager
- `LoadProviders() error`: Load providers from environment variables
- `GetEmailProviders() []EmailProviderConfig`: Get all email providers
- `GetSMSProviders() []SMSProviderConfig`: Get all SMS providers
- `HealthCheck() []ProviderHealth`: Perform health check on all providers
- `RegisterHealthChecker(name string, checker ProviderHealthChecker)`: Register custom health checker

### Utility Functions

- `GenerateProviderConfig(providerType, providerName string) (string, error)`: Generate provider configuration template
- `LoadProvidersFromEnv() (*ProviderManager, error)`: Load providers from environment
- `GetProviderHealthStatus(cfg *config.Config) map[string]interface{}`: Get provider health status
- `PrintProviderConfiguration(cfg *config.Config)`: Print provider configuration

### Provider Types

#### EmailProviderConfig

```go
type EmailProviderConfig interface {
    GetName() string
    GetType() string
    IsEnabled() bool
    GetPriority() int
    GetSettings() EmailProviderSettings
}
```

#### SMSProviderConfig

```go
type SMSProviderConfig interface {
    GetName() string
    GetType() string
    IsEnabled() bool
    GetPriority() int
    GetSettings() SMSProviderSettings
}
```

#### ProviderHealth

```go
type ProviderHealth struct {
    Name      string    // Provider name
    Type      string    // Provider type (email/sms)
    Status    string    // Status (healthy/unhealthy)
    Message   string    // Status message
    LastTest  time.Time // Last test time
}
```

## Contributing

### Development Setup

```bash
# Clone repository
git clone <repository-url>
cd rangkaiedu-backend

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Build CLI tool
go build -o providers cmd/providers/main.go
```

### Code Style

- Follow Go standard formatting
- Use meaningful variable names
- Add comments for complex logic
- Write unit tests for new features
- Update documentation for changes

### Testing Guidelines

1. **Unit Tests**: Test individual components in isolation
2. **Integration Tests**: Test component interactions
3. **End-to-End Tests**: Test complete workflows
4. **Performance Tests**: Test under load conditions

### Pull Request Process

1. Create feature branch from `main`
2. Write tests for new functionality
3. Ensure all tests pass
4. Update documentation
5. Submit pull request with detailed description
6. Address review comments

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Support

For support and questions:

1. Check the troubleshooting section
2. Review existing issues
3. Create new issue with detailed description
4. Contact development team

## Changelog

### Version 1.0.0

- Initial release of provider configuration system
- Support for multiple email and SMS providers
- Credential encryption and management
- Health monitoring and failover
- CLI management tool
- Comprehensive documentation

### Version 0.9.0

- Beta release
- Basic provider management
- Health monitoring
- CLI tool development

### Version 0.8.0

- Alpha release
- Core provider management system
- Basic configuration structure
- Initial documentation