# WhatsApp Authentication Guide

## Overview

This document provides a comprehensive guide for implementing and configuring WhatsApp authentication in the Rangkai Edu application. WhatsApp authentication allows users to securely authenticate using their WhatsApp number via one-time password (OTP) codes sent through WhatsApp messages.

## Implementation Status

**Status:** ✅ Fully Implemented  
**Last Updated:** October 2025  
**Dependencies:** WhatsApp Business API integration with rate limiting

## Features

### 1. WhatsApp OTP Login
- Users can initiate login with their email and receive OTP via WhatsApp
- Supports international phone number format (e.g., +628123456789)
- Rate limiting to prevent abuse and ensure compliance with WhatsApp policies

### 2. WhatsApp OTP Verification
- 6-digit OTP codes with 10-minute expiry
- Secure storage in database with automatic cleanup
- Verification endpoint for completing the authentication flow

### 3. Security Features
- Rate limiting to prevent spam and abuse
- Secure OTP generation using crypto/rand
- Automatic OTP expiry and cleanup
- Comprehensive security event logging

## API Endpoints

### 1. WhatsApp Login with OTP
**Endpoint:** `POST /api/auth/whatsapp/login`  
**Description:** Initiates WhatsApp OTP authentication for a user

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "optionalSecurePassword"
}
```

**Response with Password (200 OK):**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid-string",
    "name": "User Name",
    "email": "user@example.com",
    "phone": "+628123456789",
    "role": "student"
  }
}
```

**Response without Password (200 OK):**
```json
{
  "message": "OTP sent to your WhatsApp number",
  "otp_required": true
}
```

### 2. Send OTP via WhatsApp
**Endpoint:** `POST /api/auth/send-otp`  
**Description:** Sends OTP to WhatsApp number

**Request Body:**
```json
{
  "identifier": "+628123456789",
  "type": "whatsapp"
}
```

**Response (200 OK):**
```json
{
  "message": "OTP sent successfully"
}
```

### 3. Verify OTP
**Endpoint:** `POST /api/auth/verify-otp`  
**Description:** Verifies the OTP for WhatsApp authentication

**Request Body:**
```json
{
  "identifier": "+628123456789",
  "otp": "123456"
}
```

**Response (200 OK):**
```json
{
  "message": "OTP verified successfully"
}
```

## Configuration

### Environment Variables
```bash
# WhatsApp Business API Configuration
WHATSAPP_PROVIDER_0_TYPE=whatsapp_business
WHATSAPP_PROVIDER_0_NAME=primary-whatsapp
WHATSAPP_PROVIDER_0_ENABLED=true
WHATSAPP_PROVIDER_0_PRIORITY=0
WHATSAPP_PROVIDER_0_SETTINGS_PHONE_ID_NUMBER=your_phone_number_id
WHATSAPP_PROVIDER_0_SETTINGS_BUSINESS_ACCOUNT_ID=your_business_account_id
WHATSAPP_PROVIDER_0_SETTINGS_ACCESS_TOKEN=your_access_token
WHATSAPP_PROVIDER_0_SETTINGS_API_VERSION=v18.0

# Alternative WhatsApp provider (Twilio WhatsApp)
WHATSAPP_PROVIDER_1_TYPE=twilio_whatsapp
WHATSAPP_PROVIDER_1_NAME=backup-whatsapp
WHATSAPP_PROVIDER_1_ENABLED=false
WHATSAPP_PROVIDER_1_PRIORITY=1
WHATSAPP_PROVIDER_1_SETTINGS_ACCOUNT_SID=your_account_sid
WHATSAPP_PROVIDER_1_SETTINGS_AUTH_TOKEN=your_auth_token
WHATSAPP_PROVIDER_1_SETTINGS_WHATSAPP_NUMBER=your_whatsapp_number

# Rate Limiting Configuration
WHATSAPP_RATE_LIMIT_REQUESTS=10
WHATSAPP_RATE_LIMIT_WINDOW=3600
```

### WhatsApp Business API Setup
1. Create a Facebook Business Manager account
2. Set up a WhatsApp Business account and get it verified
3. Register and approve your phone number
4. Obtain API credentials (access token, phone number ID)
5. Configure the API in your environment variables
6. Ensure your phone number is verified and approved for messaging

### Twilio WhatsApp Setup
1. Create a Twilio account
2. Set up WhatsApp sandbox or get WhatsApp business verification
3. Provision a WhatsApp-enabled phone number
4. Obtain Account SID and Auth Token
5. Configure the credentials in your environment variables

## Database Schema

The OTP system uses the existing `otps` table:

```sql
CREATE TABLE IF NOT EXISTS otps (
    id SERIAL PRIMARY KEY,
    identifier VARCHAR(255) NOT NULL,  -- WhatsApp phone number
    otp VARCHAR(6) NOT NULL,           -- 6-digit OTP code
    expiry TIMESTAMP NOT NULL,         -- OTP expiry time
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(identifier, otp)
);
```

## Security Considerations

### 1. Rate Limiting
- Maximum 10 OTP requests per hour per phone number
- Prevents abuse and spam
- Complies with WhatsApp Business API policies

### 2. OTP Security
- 6-digit codes with 10-minute expiry
- One-time use (deleted after verification)
- Secure generation using crypto/rand

### 3. Data Protection
- Phone numbers stored in database
- No sensitive information in logs
- Automatic cleanup of expired OTPs

### 4. Authentication Flow
- JWT tokens with role-based access
- Secure session management
- Comprehensive security event logging

## Error Handling

### Common Error Responses

**Invalid Phone Number Format:**
```json
{
  "error": "Invalid phone number format. Use international format (e.g., +628123456789)"
}
```

**Rate Limit Exceeded:**
```json
{
  "error": "Rate limit exceeded. Please try again later."
}
```

**OTP Expired:**
```json
{
  "error": "OTP expired. Please request a new one."
}
```

**Invalid OTP:**
```json
{
  "error": "Invalid OTP. Please check and try again."
}
```

## Integration Guide

### Frontend Integration

1. **WhatsApp Login Flow:**
   ```javascript
   // Initiate WhatsApp login
   const response = await fetch('/api/auth/whatsapp/login', {
     method: 'POST',
     headers: { 'Content-Type': 'application/json' },
     body: JSON.stringify({ email: user.email })
   });
   
   const data = await response.json();
   
   if (data.otp_required) {
     // Show OTP input UI
     // Redirect to OTP verification
   }
   ```

2. **OTP Verification:**
   ```javascript
   // Verify OTP
   const response = await fetch('/api/auth/verify-otp', {
     method: 'POST',
     headers: { 'Content-Type': 'application/json' },
     body: JSON.stringify({
       identifier: phoneNumber,
       otp: userEnteredOTP
     })
   });
   ```

### Backend Integration

The WhatsApp OTP system integrates seamlessly with the existing authentication middleware:

```go
// Use existing authentication middleware
router.Use(middleware.AuthRequired())

// Role-based access control
router.Use(middleware.RoleRequired("teacher"))
```

## Testing

### Unit Tests
- OTP generation and validation
- WhatsApp message sending
- Rate limiting functionality
- Security event logging

### Integration Tests
- End-to-end authentication flow
- Database operations
- API endpoint validation
- Error handling scenarios

### Security Tests
- Rate limiting enforcement
- OTP expiry handling
- Authentication bypass attempts
- Data validation

## Monitoring and Logging

### Security Events
- OTP generation attempts
- Successful OTP verifications
- Failed authentication attempts
- Rate limit violations

### Performance Metrics
- OTP delivery success rate
- Authentication success rate
- API response times
- WhatsApp API usage

## Troubleshooting

### Common Issues

1. **OTP Not Received:**
   - Check WhatsApp Business API configuration
   - Verify phone number format
   - Check rate limits
   - Review logs for errors

2. **Authentication Failures:**
   - Verify OTP expiry
   - Check OTP format (6 digits)
   - Validate phone number in database
   - Review security logs

3. **Rate Limit Issues:**
   - Monitor OTP request frequency
   - Check client-side rate limiting
   - Review WhatsApp API quotas

### Debug Mode
Enable debug logging for troubleshooting:
```bash
DEBUG=whatsapp:* npm start
```

## Compliance

### WhatsApp Business API Policies
- Follow message template requirements
- Respect rate limits and quotas
- Include proper opt-out mechanisms
- Maintain user consent records

### Data Protection
- Comply with data privacy regulations
- Secure storage of user data
- Proper data retention policies
- User consent management

## Future Enhancements

### Planned Features
- Multi-language WhatsApp templates
- Voice OTP support
- Biometric authentication fallback
- Advanced analytics dashboard

### Performance Improvements
- OTP caching for better performance
- Asynchronous message delivery
- Database optimization
- Caching layer

## Support

For issues or questions regarding WhatsApp OTP authentication:
1. Check this documentation
2. Review troubleshooting section
3. Check application logs
4. Contact development team

## Changelog

### Version 1.0.0 (October 2025)
- Initial implementation of WhatsApp OTP authentication
- Rate limiting and security features
- API endpoints and documentation
- Integration with existing authentication system