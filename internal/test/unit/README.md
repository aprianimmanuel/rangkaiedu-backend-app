# Unit Tests

## Overview
This directory contains unit tests for the Rangkai Edu backend application. Unit tests focus on testing individual functions and methods in isolation, with mocked external dependencies. These tests are designed to be fast, reliable, and provide immediate feedback during development.

## Directory Structure

```
internal/test/unit/
├── config/              # Configuration module tests
│   ├── config_test.go
│   ├── https_test.go
│   └── providers_test.go
├── handlers/            # HTTP handler tests
│   ├── auth_handler_test.go
│   ├── class_handler_test.go
│   └── user_handler_test.go
├── middleware/          # Middleware tests
│   ├── auth_test.go
│   └── roles_test.go
└── utils/               # Utility function tests
    ├── password_test.go
    ├── email_test.go
    └── otp_test.go
```

## Unit Test Guidelines

### Characteristics
- **Fast Execution**: Tests should complete in under 1 second
- **Isolated Dependencies**: Use mocks for external services
- **Single Focus**: Test one function or method at a time
- **Deterministic**: Tests should always produce the same results

### Best Practices
1. **Follow AAA Pattern**: Arrange, Act, Assert
2. **Use Descriptive Names**: Clear test names that explain the scenario
3. **Test Edge Cases**: Include boundary conditions and error scenarios
4. **Keep Tests Independent**: Each test should run independently
5. **Use Table-Driven Tests**: For multiple test scenarios

### Mocking Strategy
- **Database**: Use `sqlmock` for database operations
- **HTTP Services**: Use `httptest` for HTTP requests
- **External APIs**: Create mock implementations
- **File System**: Use in-memory file systems

## Running Unit Tests

### Basic Commands
```bash
# Run all unit tests
go test ./internal/test/unit/...

# Run with verbose output
go test -v ./internal/test/unit/...

# Run with coverage
go test -cover ./internal/test/unit/...

# Run specific test file
go test ./internal/test/unit/handlers/auth_handler_test.go

# Run specific test function
go test -run TestAuthHandler_Login ./internal/test/unit/handlers/
```

### Coverage Requirements
- **Minimum Coverage**: 80% for all unit tests
- **Critical Paths**: 95% coverage for authentication and security code
- **Utility Functions**: 90% coverage for helper functions

## Test Examples

### Handler Test Example
```go
package handlers

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestAuthHandler_LoginWithValidCredentials(t *testing.T) {
    // Setup
    router := setupTestRouter()
    mockDB := setupMockDB()
    
    // Arrange
    loginData := `{
        "email": "test@example.com",
        "password": "validPassword123"
    }`
    
    req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(loginData))
    req.Header.Set("Content-Type", "application/json")
    
    // Act
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, recorder.Code)
    
    var response map[string]string
    require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
    
    assert.NotEmpty(t, response["token"])
    assert.Equal(t, "test@example.com", response["user"])
}
```

### Utility Test Example
```go
package utils

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestPasswordHasher_HashPassword(t *testing.T) {
    tests := []struct {
        name        string
        password    string
        expectError bool
    }{
        {
            name:        "Valid password",
            password:    "ValidPass123!",
            expectError: false,
        },
        {
            name:        "Empty password",
            password:    "",
            expectError: true,
        },
        {
            name:        "Short password",
            password:    "short",
            expectError: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            hasher := NewPasswordHasher()
            
            // Act
            hash, err := hasher.HashPassword(tt.password)
            
            // Assert
            if tt.expectError {
                assert.Error(t, err)
                assert.Empty(t, hash)
            } else {
                assert.NoError(t, err)
                assert.NotEmpty(t, hash)
                assert.NotEqual(t, tt.password, hash)
            }
        })
    }
}
```

### Middleware Test Example
```go
package middleware

import (
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
    // Setup
    router := gin.New()
    router.Use(AuthRequired())
    
    router.GET("/protected", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "success"})
    })
    
    // Arrange
    token := generateValidToken()
    req := httptest.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    
    // Act
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
    // Setup
    router := gin.New()
    router.Use(AuthRequired())
    
    router.GET("/protected", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "success"})
    })
    
    // Arrange
    req := httptest.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", "Bearer invalid-token")
    
    // Act
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)
    
    // Assert
    assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}
```

## Mock Objects

### Database Mock
```go
package mocks

import (
    "database/sql"
    "errors"
    
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/mock"
)

func SetupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create mock database: %v", err)
    }
    
    return db, mock
}

func MockUserQuery(mock sqlmock.Sqlmock, email string, user *models.User) {
    mock.ExpectQuery("SELECT id, name, email, password_hash FROM users WHERE email = \\$1").
        WithArgs(email).
        WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password_hash"}).
            AddRow("user-id", "Test User", "test@example.com", "hashed-password"))
}
```

### HTTP Mock
```go
package mocks

import (
    "net/http"
    "net/http/httptest"
    
    "github.com/gin-gonic/gin"
)

func SetupTestRouter() *gin.Engine {
    // Set Gin to test mode
    gin.SetMode(gin.TestMode)
    
    // Create router
    router := gin.New()
    
    // Add routes
    router.POST("/api/auth/login", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"token": "mock-token"})
    })
    
    return router
}
```

## Test Data

### Test User Data
```go
package data

import (
    "github.com/google/uuid"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
)

var TestUser = models.User{
    ID:           uuid.New(),
    Name:         "Test User",
    Email:        "test@example.com",
    Phone:        "1234567890",
    Role:         "student",
    PasswordHash: "hashed-password",
    IsMFAEnabled: false,
}

var TestAdmin = models.User{
    ID:           uuid.New(),
    Name:         "Admin User",
    Email:        "admin@example.com",
    Phone:        "1234567891",
    Role:         "admin",
    PasswordHash: "hashed-password",
    IsMFAEnabled: false,
}
```

## Test Utilities

### Test Setup Helpers
```go
package testutils

import (
    "database/sql"
    "testing"
    
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/require"
)

func SetupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create mock database: %v", err)
    }
    
    return db, mock
}

func RequireNoError(t *testing.T, err error) {
    require.NoError(t, err, "Unexpected error occurred")
}

func RequireError(t *testing.T, err error, expectedMessage string) {
    require.Error(t, err, "Expected error but got none")
    require.Contains(t, err.Error(), expectedMessage, "Error message mismatch")
}
```

## Performance Considerations

### Fast Execution
- Keep tests under 1 second
- Use in-memory databases
- Mock external services
- Avoid file I/O operations

### Parallel Testing
```go
func TestParallelExample(t *testing.T) {
    t.Parallel() // Run this test in parallel with others
    
    // Test implementation
    assert.True(t, true)
}
```

### Memory Management
- Clean up resources after each test
- Use defer statements for cleanup
- Avoid memory leaks in tests

## Troubleshooting

### Common Issues

#### Mock Setup Problems
- Verify mock expectations are set correctly
- Check that all expected queries are called
- Ensure proper cleanup of mock objects

#### Test Dependencies
- Keep tests independent of each other
- Avoid shared state between tests
- Use proper setup and teardown

#### Performance Issues
- Identify slow tests with `-v` flag
- Use `go test -run` to isolate slow tests
- Optimize mock configurations

### Debugging Tests
```bash
# Run tests with verbose output
go test -v ./internal/test/unit/...

# Run tests with race detection
go test -race ./internal/test/unit/...

# Run specific test with debugging
go test -run TestSpecificFunction -v ./internal/test/unit/...

# Run tests with memory profiling
go test -memprofile=mem.out ./internal/test/unit/...
go tool pprof mem.out
```

## Contributing

### Adding New Unit Tests
1. Follow the existing naming conventions
2. Use the provided test templates
3. Include both positive and negative test cases
4. Ensure proper cleanup of test data
5. Update coverage reports

### Code Review Checklist
- [ ] Tests follow naming conventions
- [ ] Tests are properly organized
- [ ] Mock objects are used appropriately
- [ ] Test data is properly managed
- [ ] Coverage requirements are met
- [ ] Tests are fast and reliable

## Additional Resources

### Documentation
- [Go Testing Documentation](https://golang.org/doc/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Go Mock Documentation](https://github.com/DATA-DOG/go-sqlmock)

### Tools
- [Go Test Coverage](https://pkg.go.dev/cmd/go#hdr-Testing_flags)
- [Testify/assert](https://github.com/stretchr/testify/assert)
- [Testify/mock](https://github.com/stretchr/testify/mock)

---

**Note**: This README should be updated as the unit test suite evolves and new testing practices are adopted.