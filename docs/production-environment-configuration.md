# Production Environment Configuration Specification

## Overview

This document provides a comprehensive technical specification for managing environment variables in the Rangkai Edu backend application. It details the categorization, validation rules, security considerations, and configuration structure for supporting multiple environments (development, staging, and production).

## Environment Variable Categorization

### 1. Database Configuration
These variables control the PostgreSQL database connection:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DB_HOST` | Yes | `localhost` | Database host address |
| `DB_PORT` | Yes | `5432` | Database port |
| `DB_NAME` | Yes | `rangkaiedu_dev` | Database name |
| `DB_USER` | Yes | `postgres` | Database username |
| `DB_PASSWORD` | Yes | `password` | Database password |
| `DB_SSLMODE` | No | `disable` | SSL mode for database connection |

### 2. Security Configuration
These variables control application security features:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `JWT_SECRET` | Yes | `default-secret-key-change-in-production` | Secret key for JWT token signing |

### 3. Email Service Configuration
These variables configure SMTP for email delivery:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `SMTP_HOST` | No | `` | SMTP server host |
| `SMTP_PORT` | No | `587` | SMTP server port |
| `SMTP_USER` | No | `` | SMTP username |
| `SMTP_PASS` | No | `` | SMTP password |

### 4. SMS Service Configuration
These variables configure Twilio for SMS delivery:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `TWILIO_ACCOUNT_SID` | No | `` | Twilio Account SID |
| `TWILIO_AUTH_TOKEN` | No | `` | Twilio Auth Token |
| `TWILIO_SENDER_PHONE` | No | `` | Twilio sender phone number |

### 5. OAuth Configuration
These variables configure social login providers:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `GOOGLE_CLIENT_ID` | No | `` | Google OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | No | `` | Google OAuth client secret |
| `FACEBOOK_CLIENT_ID` | No | `` | Facebook OAuth client ID |
| `FACEBOOK_CLIENT_SECRET` | No | `` | Facebook OAuth client secret |

### 6. Storage Configuration
These variables configure file storage options:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `STORAGE_PROVIDER` | No | `local` | Storage provider (`local`, `oss`) |
| `OSS_BUCKET_NAME` | No | `` | Alibaba Cloud OSS bucket name |
| `OSS_ACCESS_KEY_ID` | No | `` | Alibaba Cloud OSS access key ID |
| `OSS_ACCESS_KEY_SECRET` | No | `` | Alibaba Cloud OSS access key secret |
| `OSS_REGION` | No | `` | Alibaba Cloud OSS region |
| `OSS_ENDPOINT` | No | `` | Alibaba Cloud OSS endpoint |
| `GCS_BUCKET_NAME` | No | `` | Google Cloud Storage bucket name (deprecated) |
| `GCS_SERVICE_ACCOUNT_KEY_PATH` | No | `` | GCS service account key path (deprecated) |
| `GCS_PROJECT_ID` | No | `` | Google Cloud project ID (deprecated) |

## Required vs Optional Variables

### Required Variables (Application will fail to start without these)
- `DB_HOST`
- `DB_PORT`
- `DB_NAME`
- `DB_USER`
- `DB_PASSWORD`
- `JWT_SECRET`

### Optional Variables (Application will start but some features may be disabled)
All other variables are optional. If not provided:
- Email functionality will be disabled
- SMS functionality will be disabled
- Social login will be disabled
- Cloud storage will default to local file system

## Validation Rules

### Database Configuration Validation
1. All database connection parameters must be provided
2. `DB_PORT` must be a valid port number (1-65535)
3. Database connection is validated at startup

### Security Configuration Validation
1. `JWT_SECRET` must be at least 32 characters for production environments
2. Application will log a warning if a weak secret is detected

### Service Configuration Validation
1. Service credentials are validated when the respective service is first used
2. Invalid credentials will result in service-specific error messages but won't prevent application startup

## Security Considerations

### Sensitive Variables
The following variables contain sensitive information and must be handled securely:

1. `DB_PASSWORD` - Database password
2. `JWT_SECRET` - JWT signing secret
3. `SMTP_PASS` - SMTP password
4. `TWILIO_AUTH_TOKEN` - Twilio authentication token
5. `OSS_ACCESS_KEY_SECRET` - Alibaba Cloud OSS access key secret
6. `GCS_SERVICE_ACCOUNT_KEY_PATH` - GCS service account key path
7. `GOOGLE_CLIENT_SECRET` - Google OAuth client secret
8. `FACEBOOK_CLIENT_SECRET` - Facebook OAuth client secret

### Security Best Practices

1. **Environment-specific Storage**: Never store sensitive variables in version control
2. **Production Secrets**: Use secure secret management systems (e.g., HashiCorp Vault, AWS Secrets Manager)
3. **Access Control**: Limit access to environment configuration files
4. **Encryption at Rest**: Encrypt sensitive configuration files when stored
5. **Rotation Policy**: Implement regular credential rotation for all services
6. **Audit Logging**: Log access to sensitive configuration values
7. **Transmission Security**: Use secure channels when transmitting configuration values

### Development vs Production Guidelines

| Environment | Configuration Storage | Security Level |
|-------------|----------------------|----------------|
| Development | `.env` files | Low (sample data) |
| Staging | Environment variables/secrets management | High |
| Production | Environment variables/secrets management | Highest |

## Multi-Environment Configuration Structure

### Environment Hierarchy
```
Environment Configuration
├── Development
│   ├── Local development (.env files)
│   └── Shared development (environment variables)
├── Staging
│   ├── Pre-production testing
│   └── Environment variables/secrets management
└── Production
    ├── Live production environment
    └── Secure secrets management
```

### Configuration Loading Priority
1. Environment variables (highest priority)
2. `.env` file values
3. Hardcoded defaults (lowest priority)

### Environment-specific Variables

#### Development Environment
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=rangkaiedu_dev
DB_USER=postgres
DB_PASSWORD=password
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=default-secret-key-change-in-production

# Storage Configuration
STORAGE_PROVIDER=local
```

#### Staging Environment
```env
# Database Configuration
DB_HOST=staging-db.rangkaiedu.com
DB_PORT=5432
DB_NAME=rangkai_edu_staging
DB_USER=staging_user
DB_PASSWORD=staging_password
DB_SSLMODE=require

# JWT Configuration
JWT_SECRET=staging_jwt_secret_change_me

# Storage Configuration
STORAGE_PROVIDER=oss
OSS_BUCKET_NAME=rangkaiedu-staging
OSS_ACCESS_KEY_ID=staging_access_key_id
OSS_ACCESS_KEY_SECRET=staging_access_key_secret
OSS_REGION=ap-southeast-1
OSS_ENDPOINT=https://oss-ap-southeast-1.aliyuncs.com
```

#### Production Environment
```env
# Database Configuration
DB_HOST=prod-db.rangkaiedu.com
DB_PORT=5432
DB_NAME=rangkai_edu
DB_USER=prod_user
DB_PASSWORD=prod_password
DB_SSLMODE=require

# JWT Configuration
JWT_SECRET=production_jwt_secret_very_long_and_secure

# Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=prod@rangkaiedu.com
SMTP_PASS=prod_smtp_password

# SMS Configuration
TWILIO_ACCOUNT_SID=prod_twilio_account_sid
TWILIO_AUTH_TOKEN=prod_twilio_auth_token
TWILIO_SENDER_PHONE=+1234567890

# OAuth Configuration
GOOGLE_CLIENT_ID=prod_google_client_id
GOOGLE_CLIENT_SECRET=prod_google_client_secret
FACEBOOK_CLIENT_ID=prod_facebook_client_id
FACEBOOK_CLIENT_SECRET=prod_facebook_client_secret

# Storage Configuration
STORAGE_PROVIDER=oss
OSS_BUCKET_NAME=rangkaiedu-prod
OSS_ACCESS_KEY_ID=prod_access_key_id
OSS_ACCESS_KEY_SECRET=prod_access_key_secret
OSS_REGION=ap-southeast-1
OSS_ENDPOINT=https://oss-ap-southeast-1.aliyuncs.com
```

## Configuration Management Strategy

### 1. Local Development
- Use `.env` files for configuration
- `.env.example` provided as template
- Sample values for non-sensitive variables
- Clear warnings for production values

### 2. CI/CD Pipeline
- Environment variables injected via GitHub Secrets
- Different configurations for development, staging, and production
- Automated validation of required variables
- Security scanning of configuration values

### 3. Deployment Environments
- Docker containers receive environment variables at runtime
- Kubernetes secrets for sensitive values
- Environment-specific ConfigMaps for non-sensitive values
- Automated health checks validate configuration

## Deployment Team Documentation

### For System Administrators

#### Required Actions for Production Deployment
1. Set all required environment variables
2. Configure secure secret storage for sensitive values
3. Validate database connectivity before application startup
4. Test all configured services (email, SMS, storage) after deployment
5. Implement monitoring for configuration-related errors

#### Monitoring and Troubleshooting
- Watch for startup failures related to missing required variables
- Monitor logs for service configuration errors
- Set up alerts for configuration validation failures
- Maintain documentation of current configuration values

#### Security Checklist
- [ ] All sensitive variables stored in secure secret management
- [ ] No hardcoded credentials in source code
- [ ] Regular credential rotation implemented
- [ ] Access logging enabled for configuration systems
- [ ] Encryption at rest for configuration files
- [ ] Secure transmission of configuration values

### For Developers

#### Local Development Setup
1. Copy `config/.env.example` to `.env`
2. Modify values as needed for your local environment
3. Ensure all required variables are set
4. Never commit `.env` files to version control

#### Testing Configuration Changes
1. Test configuration changes in development environment first
2. Validate required variable validation works correctly
3. Test both presence and absence of optional variables
4. Verify error messages are clear and helpful

#### Adding New Configuration Variables
1. Add the variable to the `Config` struct in `config/config.go`
2. Add loading logic in the `Load()` function
3. Add default values where appropriate
4. Add validation if required
5. Update `config/.env.example` with the new variable
6. Document the variable in this specification

## Implementation Details

### Configuration Loading Process
1. Application attempts to load `.env` file using `godotenv.Load()`
2. Environment variables are read and mapped to `Config` struct
3. Default values are applied for unset variables
4. Required variable validation is performed
5. Application fails to start if required variables are missing

### Error Handling
- Clear error messages for missing required variables
- Graceful degradation for missing optional services
- Detailed logging for configuration-related issues
- Health check endpoints to verify configuration status

## Future Considerations

### Planned Enhancements
1. Configuration hot-reloading without application restart
2. Centralized configuration management service
3. Configuration versioning and rollback capabilities
4. Enhanced validation for complex configuration structures

### Scalability Considerations
1. Support for configuration sharding across multiple services
2. Environment-specific feature flags
3. Dynamic configuration based on deployment context
4. Integration with service discovery mechanisms

This specification provides a comprehensive framework for managing environment variables across all deployment environments while maintaining security best practices and operational simplicity.