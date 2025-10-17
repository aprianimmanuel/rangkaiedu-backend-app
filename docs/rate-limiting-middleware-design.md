# Rate Limiting Middleware Design Document

## 1. Overview

This document outlines the design and implementation plan for a rate limiting middleware in the Rangkai Edu backend system. The rate limiting middleware will protect the API from abuse by limiting the number of requests that can be made by clients within a specified time period.

## 2. Current Middleware Analysis

### 2.1 Middleware Structure
The current middleware implementation follows the Gin framework pattern:
- Middleware functions are defined as `gin.HandlerFunc` functions
- They have access to the `*gin.Context` which allows them to:
  - Access request information (headers, IP, etc.)
  - Set values in the context for downstream handlers
  - Abort the request chain with appropriate responses
  - Continue to the next handler with `c.Next()`

### 2.2 Middleware Registration
Middleware is registered in route groups using the `.Use()` method:
```go
classes := r.Group("/api/classes")
{
    // Apply authentication middleware to all class routes
    classes.Use(middleware.AuthRequired())
    
    // Apply role-based access control
    classes.Use(middleware.RolesRequired(middleware.RoleTeacher, middleware.RoleAdmin))
}
```

### 2.3 Configuration Structure
The application uses a configuration system based on environment variables with the following patterns:
- Configuration is loaded through `config.Load()` function
- Environment variables are accessed using helper functions like `getEnv()`, `getEnvInt()`, `getEnvBool()`
- Complex configurations are structured in dedicated structs

## 3. Rate Limiting Strategies

### 3.1 IP-Based Rate Limiting
- Limits requests based on the client's IP address
- Useful for preventing abuse from specific IP addresses
- Implementation will use the `X-Forwarded-For` header when behind a proxy, falling back to `c.ClientIP()`

### 3.2 User-Based Rate Limiting
- Limits requests based on authenticated user ID
- More granular control for authenticated users
- Requires the authentication middleware to run first

### 3.3 Endpoint-Based Rate Limiting
- Limits requests to specific endpoints or endpoint patterns
- Allows different limits for different API resources
- Can be combined with IP or user-based limiting

## 4. Technical Specification

### 4.1 Core Components

#### 4.1.1 RateLimiter Interface
```go
type RateLimiter interface {
    // Allow checks if a request should be allowed
    Allow(identifier string) (bool, *RateLimitInfo, error)
    
    // Get current rate limit info for an identifier
    Get(identifier string) (*RateLimitInfo, error)
}
```

#### 4.1.2 RateLimitInfo Struct
```go
type RateLimitInfo struct {
    Limit     int64     // Maximum requests allowed
    Remaining int64     // Remaining requests
    ResetTime time.Time // When the limit resets
}
```

#### 4.1.3 RateLimitRule Struct
```go
type RateLimitRule struct {
    // Strategy defines the limiting strategy (ip, user, endpoint)
    Strategy string
    
    // Limit is the maximum number of requests
    Limit int64
    
    // Window is the time window in seconds
    Window int64
    
    // Key is used for endpoint-based limiting (e.g., "/api/auth/login")
    Key string
    
    // Roles that this rule applies to (optional)
    Roles []string
}
```

### 4.2 Implementation Options

#### 4.2.1 In-Memory Store
For single-instance deployments:
- Uses a concurrent map with mutex protection
- Simple to implement and sufficient for development
- Not suitable for production with multiple instances

#### 4.2.2 Redis Store
For production and multi-instance deployments:
- Uses Redis for distributed rate limiting state
- Provides consistency across multiple application instances
- Requires Redis configuration

### 4.3 Configuration Structure

#### 4.3.1 Environment Variables
```
# Rate Limiting Configuration
RATE_LIMIT_ENABLED=true
RATE_LIMIT_STRATEGY=redis  # or memory
RATE_LIMIT_REDIS_ADDR=localhost:6379
RATE_LIMIT_REDIS_PASSWORD=
RATE_LIMIT_REDIS_DB=0

# Default limits
RATE_LIMIT_DEFAULT_LIMIT=100
RATE_LIMIT_DEFAULT_WINDOW=3600  # 1 hour

# Specific endpoint limits
RATE_LIMIT_ENDPOINT_/api/auth/login_LIMIT=5
RATE_LIMIT_ENDPOINT_/api/auth/login_WINDOW=300  # 5 minutes
RATE_LIMIT_ENDPOINT_/api/auth/register_LIMIT=3
RATE_LIMIT_ENDPOINT_/api/auth/register_WINDOW=300  # 5 minutes
```

#### 4.3.2 Configuration Struct
```go
type RateLimitConfig struct {
    Enabled  bool
    Strategy string // "memory" or "redis"
    
    // Redis configuration
    RedisAddr     string
    RedisPassword string
    RedisDB       int
    
    // Default limits
    DefaultLimit  int64
    DefaultWindow int64
    
    // Endpoint-specific rules
    Rules []RateLimitRule
}
```

## 5. Integration with Existing System

### 5.1 Middleware Implementation
The rate limiting middleware will follow the same pattern as existing middleware:

```go
func RateLimit() gin.HandlerFunc {
    // Initialize rate limiter based on configuration
    limiter := NewRateLimiter(config.Load().RateLimit)
    
    return func(c *gin.Context) {
        // Determine identifier based on strategy
        identifier := getIdentifier(c, config.Load().RateLimit.Strategy)
        
        // Check if request is allowed
        allowed, info, err := limiter.Allow(identifier)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limiting error"})
            c.Abort()
            return
        }
        
        // Set rate limit headers
        setRateLimitHeaders(c, info)
        
        if !allowed {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded",
                "retry_after": info.ResetTime.Unix(),
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

### 5.2 Route Registration
Rate limiting middleware can be applied at different levels:

```go
// Global rate limiting
r.Use(middleware.RateLimit())

// Group-level rate limiting
auth := r.Group("/api/auth")
auth.Use(middleware.RateLimit())

// Specific endpoint rate limiting
auth.POST("/login", middleware.RateLimitSpecific("login"), controllers.LoginHandler)
```

### 5.3 Configuration Integration
The rate limiting configuration will be integrated into the existing config system:

```go
type Config struct {
    // ... existing fields ...
    RateLimit RateLimitConfig
}

func Load() *Config {
    cfg := &Config{
        // ... existing config loading ...
        RateLimit: loadRateLimitConfig(),
    }
    return cfg
}
```

## 6. Error Handling

### 6.1 Rate Limit Exceeded
When a rate limit is exceeded:
- Return HTTP 429 (Too Many Requests) status code
- Include a JSON error response with details
- Set appropriate rate limit headers

### 6.2 Internal Errors
When rate limiting encounters internal errors:
- Log the error for debugging
- Depending on configuration, either:
  - Fail open (allow the request) for high availability
  - Fail closed (block the request) for strict security

## 7. Monitoring and Logging

### 7.1 Headers
The middleware will set standard rate limit headers:
- `X-RateLimit-Limit`: Request limit
- `X-RateLimit-Remaining`: Remaining requests
- `X-RateLimit-Reset`: Unix timestamp when limit resets

### 7.2 Logging
- Log rate limit exceeded events for security monitoring
- Log internal errors for debugging
- Include relevant context (IP, user ID, endpoint)

### 7.3 Metrics
- Track request counts by endpoint
- Track rate limit exceeded events
- Track internal errors

## 8. Implementation Plan

### 8.1 Phase 1: Core Implementation
1. Create rate limiting interface and in-memory implementation
2. Implement configuration loading
3. Create middleware function
4. Add rate limit headers to responses

### 8.2 Phase 2: Redis Integration
1. Implement Redis-based rate limiter
2. Add Redis configuration options
3. Make strategy configurable

### 8.3 Phase 3: Advanced Features
1. Implement endpoint-specific rules
2. Add user-based limiting
3. Add monitoring and logging

### 8.4 Phase 4: Testing and Documentation
1. Write unit tests
2. Create integration tests
3. Document usage and configuration

## 9. Security Considerations

### 9.1 Bypass Prevention
- Use proper IP address detection (consider X-Forwarded-For headers)
- Prevent header spoofing
- Consider using additional identifiers for more robust limiting

### 9.2 DoS Protection
- Ensure rate limiting itself doesn't become a DoS vector
- Use efficient data structures and algorithms
- Implement proper error handling to prevent resource exhaustion

## 10. Performance Considerations

### 10.1 Efficiency
- Use efficient data structures for tracking requests
- Minimize locking in concurrent environments
- Optimize Redis operations with pipelining where appropriate

### 10.2 Scalability
- Design for horizontal scaling with Redis backend
- Consider sharding for very high request volumes
- Implement circuit breakers for external dependencies

## 11. Configuration Examples

### 11.1 Basic Configuration
```
RATE_LIMIT_ENABLED=true
RATE_LIMIT_STRATEGY=memory
RATE_LIMIT_DEFAULT_LIMIT=100
RATE_LIMIT_DEFAULT_WINDOW=3600
```

### 11.2 Advanced Configuration
```
RATE_LIMIT_ENABLED=true
RATE_LIMIT_STRATEGY=redis
RATE_LIMIT_REDIS_ADDR=redis:6379
RATE_LIMIT_REDIS_PASSWORD=secret
RATE_LIMIT_REDIS_DB=1
RATE_LIMIT_DEFAULT_LIMIT=1000
RATE_LIMIT_DEFAULT_WINDOW=3600

# Endpoint-specific rules
RATE_LIMIT_ENDPOINT_/api/auth/login_LIMIT=5
RATE_LIMIT_ENDPOINT_/api/auth/login_WINDOW=300
RATE_LIMIT_ENDPOINT_/api/auth/register_LIMIT=3
RATE_LIMIT_ENDPOINT_/api/auth/register_WINDOW=300
```

## 12. Testing Strategy

### 12.1 Unit Tests
- Test rate limiter implementations
- Test identifier extraction functions
- Test configuration loading

### 12.2 Integration Tests
- Test middleware with different strategies
- Test endpoint-specific rules
- Test interaction with other middleware

### 12.3 Load Testing
- Verify performance under load
- Test rate limit behavior
- Validate Redis integration

## 13. Deployment Considerations

### 13.1 Single Instance
- Use in-memory strategy
- No external dependencies
- Suitable for development and small deployments

### 13.2 Multi-Instance
- Use Redis strategy
- Ensure Redis is properly configured and secured
- Consider Redis clustering for high availability

## 14. Future Enhancements

### 14.1 Dynamic Configuration
- Allow runtime configuration changes
- Implement configuration reloading

### 14.2 Advanced Algorithms
- Implement token bucket algorithm
- Implement leaky bucket algorithm
- Support sliding window limits

### 14.3 Analytics
- Track rate limiting metrics
- Provide dashboards for monitoring
- Generate reports on API usage patterns