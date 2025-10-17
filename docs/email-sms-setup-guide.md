# Email and SMS Service Setup Guide

## Overview

This guide provides detailed instructions for setting up and configuring email and SMS services in the Rangkai Edu backend application. It covers the configuration of SMTP email services and Twilio SMS services, including security best practices and troubleshooting tips.

## Prerequisites

Before setting up email and SMS services, ensure you have:

1. Access to the application's configuration files
2. Valid credentials for your chosen email provider
3. Valid credentials for your chosen SMS provider
4. Understanding of the application's environment configuration

## Email Service Configuration

### Supported Email Providers

The Rangkai Edu backend supports multiple email providers:

1. **SMTP** - Generic SMTP servers (Gmail, Outlook, custom SMTP)
2. **SendGrid** - SendGrid Email API
3. **Gmail/Google Workspace** - OAuth2 authentication with Gmail

### SMTP Configuration

#### Gmail/Google Workspace Setup

1. **Enable 2-Factor Authentication**:
   - Go to your Google Account settings
   - Navigate to Security > 2-Step Verification
   - Enable 2-Step Verification

2. **Generate App Password**:
   - Go to Security > App passwords
   - Select "Mail" and your device
   - Generate and save the app password

3. **Configuration**:
   ```env
   EMAIL_PROVIDER_0_TYPE=smtp
   EMAIL_PROVIDER_0_NAME=gmail-primary
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
   ```

#### Custom SMTP Setup

For custom SMTP providers:

```env
EMAIL_PROVIDER_0_TYPE=smtp
EMAIL_PROVIDER_0_NAME=custom-smtp
EMAIL_PROVIDER_0_ENABLED=true
EMAIL_PROVIDER_0_PRIORITY=0
EMAIL_PROVIDER_0_SETTINGS_HOST=your.smtp.server.com
EMAIL_PROVIDER_0_SETTINGS_PORT=587
EMAIL_PROVIDER_0_SETTINGS_USERNAME=your_username
EMAIL_PROVIDER_0_SETTINGS_PASSWORD=your_password
EMAIL_PROVIDER_0_SETTINGS_FROM_EMAIL=noreply@yourdomain.com
EMAIL_PROVIDER_0_SETTINGS_FROM_NAME=Your Application Name
EMAIL_PROVIDER_0_SETTINGS_USE_TLS=true
EMAIL_PROVIDER_0_SETTINGS_USE_STARTTLS=true
```

### SendGrid Configuration

1. **Create SendGrid Account**:
   - Sign up at https://sendgrid.com/
   - Complete the verification process

2. **Generate API Key**:
   - Navigate to Settings > API Keys
   - Create a new API key with "Mail Send" permissions
   - Save the API key securely

3. **Configuration**:
   ```env
   EMAIL_PROVIDER_0_TYPE=sendgrid
   EMAIL_PROVIDER_0_NAME=sendgrid-primary
   EMAIL_PROVIDER_0_ENABLED=true
   EMAIL_PROVIDER_0_PRIORITY=0
   EMAIL_PROVIDER_0_SETTINGS_API_KEY=your_sendgrid_api_key
   EMAIL_PROVIDER_0_SETTINGS_FROM_EMAIL=noreply@yourdomain.com
   EMAIL_PROVIDER_0_SETTINGS_FROM_NAME=Rangkai Edu
   ```

### Gmail OAuth2 Configuration (Advanced)

For enterprise environments requiring OAuth2:

1. **Create Google Cloud Project**:
   - Go to https://console.cloud.google.com/
   - Create a new project or select existing one

2. **Enable Gmail API**:
   - Navigate to APIs & Services > Library
   - Search for "Gmail API" and enable it

3. **Create OAuth2 Credentials**:
   - Go to APIs & Services > Credentials
   - Create OAuth2 Client ID for "Web application"
   - Add authorized redirect URIs

4. **Configuration**:
   ```env
   EMAIL_PROVIDER_0_TYPE=gmail
   EMAIL_PROVIDER_0_NAME=gmail-oauth2
   EMAIL_PROVIDER_0_ENABLED=true
   EMAIL_PROVIDER_0_PRIORITY=0
   EMAIL_PROVIDER_0_SETTINGS_CLIENT_ID=your_client_id
   EMAIL_PROVIDER_0_SETTINGS_CLIENT_SECRET=your_client_secret
   EMAIL_PROVIDER_0_SETTINGS_REFRESH_TOKEN=your_refresh_token
   EMAIL_PROVIDER_0_SETTINGS_FROM_EMAIL=your_email@yourdomain.com
   EMAIL_PROVIDER_0_SETTINGS_FROM_NAME=Rangkai Edu
   ```

## SMS Service Configuration

### Supported SMS Providers

The Rangkai Edu backend supports multiple SMS providers:

1. **Twilio** - Twilio Programmable SMS
2. **Amazon SNS** - Amazon Simple Notification Service

### Twilio Configuration

1. **Create Twilio Account**:
   - Sign up at https://www.twilio.com/
   - Complete the verification process

2. **Get Account Credentials**:
   - Navigate to Console > Account Info
   - Copy your Account SID and Auth Token

3. **Configure Phone Number**:
   - Go to Console > Phone Numbers > Manage Numbers
   - Purchase or configure a phone number for sending SMS

4. **Configuration**:
   ```env
   SMS_PROVIDER_0_TYPE=twilio
   SMS_PROVIDER_0_NAME=twilio-primary
   SMS_PROVIDER_0_ENABLED=true
   SMS_PROVIDER_0_PRIORITY=0
   SMS_PROVIDER_0_SETTINGS_ACCOUNT_SID=your_account_sid
   SMS_PROVIDER_0_SETTINGS_AUTH_TOKEN=your_auth_token
   SMS_PROVIDER_0_SETTINGS_FROM_NUMBER=+1234567890
   ```

### Amazon SNS Configuration

1. **Create AWS Account**:
   - Sign up at https://aws.amazon.com/
   - Complete the verification process

2. **Create IAM User**:
   - Navigate to IAM > Users
   - Create a new user with programmatic access
   - Attach the "AmazonSNSFullAccess" policy or create a custom policy

3. **Get Credentials**:
   - Save the Access Key ID and Secret Access Key
   - Store them securely

4. **Configuration**:
   ```env
   SMS_PROVIDER_0_TYPE=sns
   SMS_PROVIDER_0_NAME=sns-primary
   SMS_PROVIDER_0_ENABLED=true
   SMS_PROVIDER_0_PRIORITY=0
   SMS_PROVIDER_0_SETTINGS_ACCESS_KEY_ID=your_access_key_id
   SMS_PROVIDER_0_SETTINGS_SECRET_ACCESS_KEY=your_secret_access_key
   SMS_PROVIDER_0_SETTINGS_REGION=us-east-1
   SMS_PROVIDER_0_SETTINGS_FROM_NUMBER=+1234567890
   ```

## Security Best Practices

### Credential Management

1. **Never Store Plain Text Credentials**:
   - Use encrypted storage for all sensitive values
   - In production, use secret management systems (HashiCorp Vault, AWS Secrets Manager)

2. **Environment-Specific Configuration**:
   - Use different credentials for development, staging, and production
   - Never use production credentials in non-production environments

3. **Regular Credential Rotation**:
   - Rotate credentials on a regular schedule
   - Implement automated rotation where possible

### Network Security

1. **Use Encrypted Connections**:
   - Always use TLS/SSL for SMTP connections
   - Validate SSL certificates
   - Use secure endpoints for API calls

2. **Firewall Configuration**:
   - Restrict access to email/SMS services to application servers only
   - Use VPC endpoints for AWS services

### Access Control

1. **Principle of Least Privilege**:
   - Grant only necessary permissions to service accounts
   - Use separate credentials for different services

2. **Audit Logging**:
   - Enable logging for all email/SMS activities
   - Monitor for unauthorized access attempts

## Environment Configuration

### Development Environment

For local development, use a `.env` file:

```env
# Email Configuration
EMAIL_PROVIDER_0_TYPE=smtp
EMAIL_PROVIDER_0_NAME=dev-smtp
EMAIL_PROVIDER_0_ENABLED=true
EMAIL_PROVIDER_0_PRIORITY=0
EMAIL_PROVIDER_0_SETTINGS_HOST=localhost
EMAIL_PROVIDER_0_SETTINGS_PORT=1025
EMAIL_PROVIDER_0_SETTINGS_USERNAME=
EMAIL_PROVIDER_0_SETTINGS_PASSWORD=
EMAIL_PROVIDER_0_SETTINGS_FROM_EMAIL=dev@rangkaiedu.local
EMAIL_PROVIDER_0_SETTINGS_FROM_NAME=Rangkai Edu Dev

# SMS Configuration (using mock service for development)
SMS_PROVIDER_0_TYPE=twilio
SMS_PROVIDER_0_NAME=dev-twilio
SMS_PROVIDER_0_ENABLED=true
SMS_PROVIDER_0_PRIORITY=0
SMS_PROVIDER_0_SETTINGS_ACCOUNT_SID=ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
SMS_PROVIDER_0_SETTINGS_AUTH_TOKEN=your_auth_token
SMS_PROVIDER_0_SETTINGS_FROM_NUMBER=+15005550006
```

### Production Environment

For production environments, use secure secret management:

```env
# Email Configuration
EMAIL_PROVIDER_0_TYPE=sendgrid
EMAIL_PROVIDER_0_NAME=production-sendgrid
EMAIL_PROVIDER_0_ENABLED=true
EMAIL_PROVIDER_0_PRIORITY=0
EMAIL_PROVIDER_0_SETTINGS_API_KEY=secret://SENDGRID_API_KEY
EMAIL_PROVIDER_0_SETTINGS_FROM_EMAIL=noreply@rangkaiedu.com
EMAIL_PROVIDER_0_SETTINGS_FROM_NAME=Rangkai Edu

# SMS Configuration
SMS_PROVIDER_0_TYPE=twilio
SMS_PROVIDER_0_NAME=production-twilio
SMS_PROVIDER_0_ENABLED=true
SMS_PROVIDER_0_PRIORITY=0
SMS_PROVIDER_0_SETTINGS_ACCOUNT_SID=secret://TWILIO_ACCOUNT_SID
SMS_PROVIDER_0_SETTINGS_AUTH_TOKEN=secret://TWILIO_AUTH_TOKEN
SMS_PROVIDER_0_SETTINGS_FROM_NUMBER=+1234567890
```

## Testing Configuration

### Email Service Testing

1. **Test SMTP Connection**:
   ```bash
   # Test SMTP connectivity
   telnet smtp.gmail.com 587
   ```

2. **Send Test Email**:
   ```bash
   # Use the application's test endpoint
   curl -X POST http://localhost:8080/api/v1/test/email \
     -H "Content-Type: application/json" \
     -d '{"to":"test@example.com","subject":"Test Email","body":"This is a test email"}'
   ```

### SMS Service Testing

1. **Test Twilio Credentials**:
   ```bash
   # Test Twilio API access
   curl -X GET "https://api.twilio.com/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Messages.json" \
     -u "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX:your_auth_token"
   ```

2. **Send Test SMS**:
   ```bash
   # Use the application's test endpoint
   curl -X POST http://localhost:8080/api/v1/test/sms \
     -H "Content-Type: application/json" \
     -d '{"to":"+1234567890","message":"This is a test SMS"}'
   ```

## Monitoring and Troubleshooting

### Common Issues

1. **SMTP Authentication Failures**:
   - Verify username and password
   - Check if app passwords are required
   - Ensure 2FA is properly configured

2. **SMS Delivery Failures**:
   - Verify phone number format
   - Check Twilio account balance
   - Review Twilio error codes

3. **Rate Limiting**:
   - Implement exponential backoff
   - Monitor API usage limits
   - Use multiple providers for high-volume sending

### Logging

Enable detailed logging for troubleshooting:

```env
# Enable detailed email logging
LOG_LEVEL=debug
EMAIL_DEBUG=true

# Enable detailed SMS logging
SMS_DEBUG=true
```

### Health Checks

Implement health checks for monitoring:

```bash
# Check email service health
curl http://localhost:8080/health/email

# Check SMS service health
curl http://localhost:8080/health/sms
```

## Advanced Configuration

### Multiple Provider Setup (Fallback)

Configure multiple providers for redundancy:

```env
# Primary email provider
EMAIL_PROVIDER_0_TYPE=sendgrid
EMAIL_PROVIDER_0_NAME=sendgrid-primary
EMAIL_PROVIDER_0_ENABLED=true
EMAIL_PROVIDER_0_PRIORITY=0
EMAIL_PROVIDER_0_SETTINGS_API_KEY=your_sendgrid_api_key

# Fallback email provider
EMAIL_PROVIDER_1_TYPE=smtp
EMAIL_PROVIDER_1_NAME=gmail-fallback
EMAIL_PROVIDER_1_ENABLED=true
EMAIL_PROVIDER_1_PRIORITY=1
EMAIL_PROVIDER_1_SETTINGS_HOST=smtp.gmail.com
EMAIL_PROVIDER_1_SETTINGS_PORT=587
EMAIL_PROVIDER_1_SETTINGS_USERNAME=your_email@gmail.com
EMAIL_PROVIDER_1_SETTINGS_PASSWORD=your_app_password

# Primary SMS provider
SMS_PROVIDER_0_TYPE=twilio
SMS_PROVIDER_0_NAME=twilio-primary
SMS_PROVIDER_0_ENABLED=true
SMS_PROVIDER_0_PRIORITY=0
SMS_PROVIDER_0_SETTINGS_ACCOUNT_SID=your_account_sid
SMS_PROVIDER_0_SETTINGS_AUTH_TOKEN=your_auth_token

# Fallback SMS provider
SMS_PROVIDER_1_TYPE=sns
SMS_PROVIDER_1_NAME=sns-fallback
SMS_PROVIDER_1_ENABLED=true
SMS_PROVIDER_1_PRIORITY=1
SMS_PROVIDER_1_SETTINGS_ACCESS_KEY_ID=your_access_key_id
SMS_PROVIDER_1_SETTINGS_SECRET_ACCESS_KEY=your_secret_access_key
```

### Custom Templates

Customize email and SMS templates:

```go
// Email template
const otpEmailTemplate = `
Subject: Your Rangkai Edu OTP Code

Dear User,

Your OTP code is: {{.OTP}}

This code expires in {{.ExpiryMinutes}} minutes. Do not share it with anyone.

If you didn't request this, please ignore this email.

Best regards,
The Rangkai Edu Team
`

// SMS template
const otpSMSTemplate = `Your Rangkai Edu OTP code is: {{.OTP}}. This code expires in {{.ExpiryMinutes}} minutes. Do not share it with anyone. If you didn't request this, please ignore this message.`
```

This guide provides a comprehensive overview of setting up email and SMS services in the Rangkai Edu backend application. Follow these instructions to configure secure and reliable communication services for your application.