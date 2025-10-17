# Security Logging Configuration Guide

## Table of Contents
1. [Overview](#overview)
2. [Configuration Structure](#configuration-structure)
3. [Security Logging Configuration](#security-logging-configuration)
4. [Environment-Specific Setup](#environment-specific-setup)
5. [Output Destinations](#output-destinations)
6. [Security Considerations](#security-considerations)
7. [Best Practices](#best-practices)

## Overview

This guide provides detailed instructions for configuring security event logging in the RangkaiEdu backend system. The security logging system is designed to capture, process, and store security-relevant events while maintaining performance and compliance requirements.

## Configuration Structure

The security configuration is organized in a hierarchical structure:

```json
{
  "security": {
    "enabled": true,
    "environment": "production",
    "debug": false,
    "logging": { /* Logging configuration */ },
    "alerting": { /* Alerting configuration */ },
    "authentication": { /* Authentication configuration */ },
    "authorization": { /* Authorization configuration */ },
    "session": { /* Session management configuration */ },
    "rate_limiting": { /* Rate limiting configuration */ },
    "ip_filtering": { /* IP filtering configuration */ },
    "threat_detection": { /* Threat detection configuration */ },
    "compliance": { /* Compliance configuration */ },
    "encryption": { /* Encryption configuration */ },
    "audit": { /* Audit configuration */ },
    "performance": { /* Performance configuration */ },
    "retention": { /* Retention configuration */ },
    "monitoring": { /* Monitoring configuration */ }
  }
}
```

## Security Logging Configuration

### Basic Configuration

The core logging configuration includes:

```json
{
  "level": "medium",
  "format": "json",
  "buffer_size": 1000,
  "flush_interval": "5s",
  "batch_size": 100,
  "async_logging": true,
  "queue_size": 10000,
  "worker_count": 3
}
```

**Configuration Options:**

- `level`: Minimum severity level to log (`low`, `medium`, `high`, `critical`)
- `format`: Log format (`json`, `text`, `csv`)
- `buffer_size`: Number of events to buffer in memory
- `flush_interval`: How often to flush buffered events
- `batch_size`: Number of events to process in each batch
- `async_logging`: Whether to log asynchronously (recommended for performance)
- `queue_size`: Size of the event queue
- `worker_count`: Number of worker goroutines for processing events

### Event Filtering

Control which events are logged using include/exclude filters:

```json
{
  "include_events": [
    "auth_failure",
    "auth_lockout",
    "auth_brute_force",
    "auth_breach",
    "auth_threat"
  ],
  "exclude_events": [
    "auth_success"
  ],
  "include_users": [
    "admin",
    "root"
  ],
  "exclude_users": [
    "test"
  ],
  "include_ips": [
    "192.168.0.0/16"
  ],
  "exclude_ips": [
    "127.0.0.1"
  ]
}
```

### Data Sanitization

Protect sensitive information in logs:

```json
{
  "sanitize_data": true,
  "mask_patterns": [
    "password",
    "token",
    "secret",
    "key",
    "authorization",
    "cookie"
  ],
  "retention_days": 90,
  "max_file_size": 104857600,
  "max_backups": 10,
  "encrypt_logs": true,
  "encryption_key": "base64-encoded-encryption-key",
  "sign_logs": true,
  "private_key_path": "/etc/rangkaiedu/security/private.key",
  "public_key_path": "/etc/rangkaiedu/security/public.key"
}
```

## Environment-Specific Setup

### Development Environment

For development, use minimal logging with console output:

```json
{
  "security": {
    "enabled": true,
    "environment": "development",
    "debug": true,
    "logging": {
      "level": "low",
      "format": "text",
      "console_output": {
        "format": "text",
        "color": true,
        "structured": false
      },
      "async_logging": false,
      "sanitize_data": false
    }
  }
}
```

### Production Environment

For production, use comprehensive logging with multiple outputs:

```json
{
  "security": {
    "enabled": true,
    "environment": "production",
    "debug": false,
    "logging": {
      "level": "medium",
      "format": "json",
      "file_output": {
        "path": "/var/log/rangkaiedu/security.log",
        "format": "json",
        "rotate": true,
        "compress": true,
        "permissions": 640
      },
      "http_output": {
        "url": "https://logging-service.example.com/api/v1/security-events",
        "method": "POST",
        "headers": {
          "Authorization": "Bearer secret-token"
        },
        "timeout": "30s"
      },
      "async_logging": true,
      "sanitize_data": true,
      "encrypt_logs": true,
      "sign_logs": true
    }
  }
}
```

## Output Destinations

### File Output

Log to local files with rotation:

```json
{
  "file_output": {
    "path": "/var/log/rangkaiedu/security.log",
    "format": "json",
    "rotate": true,
    "compress": true,
    "permissions": 640
  }
}
```

### Console Output

Log to standard output:

```json
{
  "console_output": {
    "format": "text",
    "color": true,
    "structured": false
  }
}
```

### HTTP Output

Send logs to external services:

```json
{
  "http_output": {
    "url": "https://logging-service.example.com/api/v1/security-events",
    "method": "POST",
    "headers": {
      "Authorization": "Bearer secret-token",
      "Content-Type": "application/json"
    },
    "timeout": "30s",
    "retry_count": 3,
    "retry_delay": "5s"
  }
}
```

### Syslog Output

Send logs to syslog:

```json
{
  "syslog_output": {
    "network": "tcp",
    "address": "syslog.example.com:514",
    "tag": "rangkaiedu-security",
    "facility": "local0",
    "severity": "info"
  }
}
```

## Security Considerations

### Log Encryption

Enable encryption for sensitive environments:

```json
{
  "encrypt_logs": true,
  "encryption_key": "base64-encoded-encryption-key",
  "sign_logs": true,
  "private_key_path": "/etc/rangkaiedu/security/private.key",
  "public_key_path": "/etc/rangkaiedu/security/public.key"
}
```

### Access Control

Restrict access to log files:
- Set appropriate file permissions (e.g., 640)
- Use dedicated log directories with restricted access
- Implement audit logging for log access

### Log Integrity

Ensure log integrity through:
- Digital signatures for log entries
- Hash chains to detect tampering
- Regular integrity checks

## Best Practices

### 1. Performance Optimization

- Use asynchronous logging in production
- Configure appropriate buffer sizes
- Implement connection pooling for external outputs
- Monitor queue depths and processing rates

### 2. Data Protection

- Always sanitize sensitive data
- Encrypt logs when stored or transmitted
- Implement proper access controls
- Regularly rotate encryption keys

### 3. Compliance

- Retain logs according to regulatory requirements
- Implement audit trails for log access
- Use structured logging formats for analysis
- Regular compliance reporting

### 4. Monitoring

- Monitor log processing performance
- Set up alerts for logging failures
- Track security event volumes and patterns
- Implement dashboards for security metrics

### 5. Maintenance

- Regular log rotation and cleanup
- Monitor storage capacity
- Update threat detection rules
- Review and update configurations regularly

## Troubleshooting

### Common Issues

1. **Log Processing Delays**
   - Check queue depths
   - Increase worker count if needed
   - Optimize external service connections

2. **Logging Failures**
   - Verify output destination connectivity
   - Check file permissions
   - Review authentication credentials

3. **Performance Issues**
   - Reduce logging level
   - Optimize buffer sizes
   - Implement circuit breakers

### Monitoring Metrics

Key metrics to monitor:
- Event processing rate
- Queue depth
- Error rates
- Output latency
- Resource usage

## Conclusion

This configuration guide provides a comprehensive framework for setting up security event logging in the RangkaiEdu backend system. By following these guidelines, you can ensure that your security logging implementation is both effective and efficient while meeting compliance requirements.