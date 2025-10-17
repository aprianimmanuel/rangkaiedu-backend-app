# Integration Tests

## Overview
This directory contains integration tests for the Rangkai Edu backend application. Integration tests verify that different components of the application work together correctly with real dependencies. These tests use actual databases, HTTP services, and external dependencies to ensure complete functionality.

## Directory Structure

```
internal/test/integration/
├── api/                 # API endpoint integration tests
│   ├── auth_test.go
│   ├── users_test.go
│   └── classes_test.go
├── database/            # Database integration tests
│   ├── connection_test.go
│   ├── migrations_test.go
│   └── queries_test.go
└── services/            # Service layer integration tests
    ├── auth_service_test.go
    ├── email_service_test.go
    └── notification_service_test.go
```

## Integration Test Guidelines

### Characteristics
- **Real Dependencies**: Use actual database, HTTP services, and external APIs
- **Component Interaction**: Test how different components work together
- **End-to-End Flows**: Test complete business processes
- **Environment Specific**: Require proper test environment setup

### Best Practices
1. **Use Real Data**: Test with realistic test data
2. **Test Complete Flows**: Verify entire user journeys
3. **Environment Setup**: Ensure proper test environment isolation
4. **Data Cleanup**: Clean up test data after each test
5. **Error Handling**: Test error scenarios and recovery

### Test Categories

#### API Integration Tests
- Test complete API endpoints
- Verify HTTP status codes and responses
- Test authentication and authorization
- Validate request/response formats

#### Database Integration Tests
- Test database connectivity and operations
- Verify data persistence and retrieval
- Test transaction management
- Validate database constraints

#### Service Integration Tests
- Test service layer interactions
- Verify external service integrations
- Test business logic implementation
- Validate error handling and recovery

## Running Integration Tests

### Prerequisites
- PostgreSQL 15+ running on localhost:5432
- Test database `rangkaiedu_test` created
- Environment variables configured

### Basic Commands
```bash
# Run all integration tests
go test ./internal/test/integration/...

# Run with verbose output
go test -v ./internal/test/integration/...

# Run with coverage
go test -cover ./internal/test/integration/...

# Run specific test file
go test ./internal/test/integration/api/auth_test.go

# Run specific test function
go test -run TestAuthFlow ./internal/test/integration/api/
```

### Environment Setup
```bash
# Create test database
createdb rangkaiedu_test

# Run migrations
psql -d rangkaiedu_test -f migrations/001_create_tables.sql

# Set environment variables
export TEST_DATABASE_URL="postgres://test:test@localhost:5432/rangkaiedu_test"
export TEST_JWT_SECRET="test-jwt-secret-change-in-production"
```

### Coverage Requirements
- **Minimum Coverage**: 90% for all integration tests
- **API Endpoints**: 100% coverage for all REST endpoints
- **Database Operations**: 95% coverage for CRUD operations
- **Service Integration**: 90% coverage for service interactions

## Test Examples

### API Integration Test Example
```go
package api

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/handlers"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
)

func TestAuthAPI_RegisterUser(t *testing.T) {
    // Setup test database
    db := setupTestDatabase(t)
    defer cleanupTestDatabase(t, db)
    
    // Setup test server
    router := setupTestServer(db)
    
    // Arrange
    registerData := map[string]interface{}{
        "name":     "John Doe",
        "email":    "john@example.com",
        "phone":    "1234567890",
        "password": "StrongPass123!",
        "role":     "student",
    }
    
    jsonData, _ := json.Marshal(registerData)
    req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    // Act
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)
    
    // Assert
    assert.Equal(t, http.StatusCreated, recorder.Code)
    
    var response map[string]interface{}
    require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
    
    // Verify user was created in database
    var user models.User
    err := db.QueryRow("SELECT id, name, email, role FROM users WHERE email = $1", "john@example.com").Scan(
        &user.ID, &user.Name, &user.Email, &user.Role)
    require.NoError(t, err)
    assert.Equal(t, "John Doe", user.Name)
    assert.Equal(t, "student", user.Role)
}
```

### Database Integration Test Example
```go
package database

import (
    "database/sql"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
)

func TestUserDatabase_CRUDOperations(t *testing.T) {
    // Setup test database
    db := setupTestDatabase(t)
    defer cleanupTestDatabase(t, db)
    
    // Test Create
    t.Run("CreateUser", func(t *testing.T) {
        user := &models.User{
            Name:         "Test User",
            Email:        sql.NullString{String: "test@example.com", Valid: true},
            Phone:        "1234567890",
            Role:         "student",
            PasswordHash: "hashed-password",
            IsMFAEnabled: false,
        }
        
        err := models.CreateUser(db, user)
        require.NoError(t, err)
        assert.NotEmpty(t, user.ID)
    })
    
    // Test Read
    t.Run("GetUser", func(t *testing.T) {
        user, err := models.GetUserByEmail(db, "test@example.com")
        require.NoError(t, err)
        assert.Equal(t, "Test User", user.Name)
        assert.Equal(t, "student", user.Role)
    })
    
    // Test Update
    t.Run("UpdateUser", func(t *testing.T) {
        user, err := models.GetUserByEmail(db, "test@example.com")
        require.NoError(t, err)
        
        user.Name = "Updated User"
        user.Role = "teacher"
        
        err = models.UpdateUser(db, user)
        require.NoError(t, err)
        
        updatedUser, err := models.GetUserByEmail(db, "test@example.com")
        require.NoError(t, err)
        assert.Equal(t, "Updated User", updatedUser.Name)
        assert.Equal(t, "teacher", updatedUser.Role)
    })
    
    // Test Delete
    t.Run("DeleteUser", func(t *testing.T) {
        user, err := models.GetUserByEmail(db, "test@example.com")
        require.NoError(t, err)
        
        err = models.DeleteUser(db, user.ID)
        require.NoError(t, err)
        
        _, err = models.GetUserByEmail(db, "test@example.com")
        assert.Error(t, err)
    })
}
```

### Service Integration Test Example
```go
package services

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/services"
)

func TestAuthService_LoginFlow(t *testing.T) {
    // Setup test database
    db := setupTestDatabase(t)
    defer cleanupTestDatabase(t, db)
    
    // Create test user
    user := createTestUser(t, db, "test@example.com", "password123")
    
    // Setup auth service
    authService := services.NewAuthService(db)
    
    t.Run("SuccessfulLogin", func(t *testing.T) {
        // Act
        token, err := authService.Login("test@example.com", "password123")
        
        // Assert
        require.NoError(t, err)
        assert.NotEmpty(t, token)
        
        // Verify token is valid
        claims, err := authService.ValidateToken(token)
        require.NoError(t, err)
        assert.Equal(t, user.ID, claims["sub"])
        assert.Equal(t, "test@example.com", claims["email"])
    })
    
    t.Run("InvalidPassword", func(t *testing.T) {
        // Act
        token, err := authService.Login("test@example.com", "wrongpassword")
        
        // Assert
        assert.Error(t, err)
        assert.Empty(t, token)
        assert.Contains(t, err.Error(), "invalid credentials")
    })
    
    t.Run("UserNotFound", func(t *testing.T) {
        // Act
        token, err := authService.Login("nonexistent@example.com", "password123")
        
        // Assert
        assert.Error(t, err)
        assert.Empty(t, token)
        assert.Contains(t, err.Error(), "user not found")
    })
}
```

## Test Data Management

### Test Database Setup
```go
package testutils

import (
    "database/sql"
    "os"
    "testing"
    
    _ "github.com/jackc/pgx/v5/stdlib"
)

func SetupTestDatabase(t *testing.T) *sql.DB {
    t.Helper()
    
    // Connect to test database
    db, err := sql.Open("pgx", os.Getenv("TEST_DATABASE_URL"))
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }
    
    // Run migrations
    if err := runMigrations(db); err != nil {
        t.Fatalf("Failed to run migrations: %v", err)
    }
    
    return db
}

func CleanupTestDatabase(t *testing.T, db *sql.DB) {
    t.Helper()
    
    // Truncate all tables
    if _, err := db.Exec("TRUNCATE TABLE users, otps, classes, materials, enrollments RESTART IDENTITY CASCADE"); err != nil {
        t.Logf("Warning: Failed to cleanup test database: %v", err)
    }
    
    db.Close()
}
```

### Test Data Fixtures
```go
package fixtures

import (
    "database/sql"
    "time"
    
    "github.com/google/uuid"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
)

type UserFixture struct {
    ID           uuid.UUID
    Name         string
    Email        string
    Phone        string
    Role         string
    PasswordHash string
    IsMFAEnabled bool
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

func CreateTestUser(db *sql.DB, fixture UserFixture) *models.User {
    user := &models.User{
        ID:           fixture.ID,
        Name:         fixture.Name,
        Email:        sql.NullString{String: fixture.Email, Valid: true},
        Phone:        fixture.Phone,
        Role:         fixture.Role,
        PasswordHash: fixture.PasswordHash,
        IsMFAEnabled: fixture.IsMFAEnabled,
        CreatedAt:    fixture.CreatedAt,
        UpdatedAt:    fixture.UpdatedAt,
    }
    
    if err := models.CreateUser(db, user); err != nil {
        panic(err)
    }
    
    return user
}

// Predefined fixtures
var (
    StudentUser = UserFixture{
        ID:           uuid.New(),
        Name:         "Test Student",
        Email:        "student@example.com",
        Phone:        "1234567890",
        Role:         "student",
        PasswordHash: "$2a$10$hashedpassword",
        IsMFAEnabled: false,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    
    TeacherUser = UserFixture{
        ID:           uuid.New(),
        Name:         "Test Teacher",
        Email:        "teacher@example.com",
        Phone:        "1234567891",
        Role:         "teacher",
        PasswordHash: "$2a$10$hashedpassword",
        IsMFAEnabled: false,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
)
```

## Test Utilities

### Test Server Setup
```go
package testutils

import (
    "database/sql"
    
    "github.com/gin-gonic/gin"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/handlers"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/routes"
)

func SetupTestServer(db *sql.DB) *gin.Engine {
    // Set Gin to test mode
    gin.SetMode(gin.TestMode)
    
    // Create router
    router := gin.New()
    
    // Setup routes
    routes.SetupRoutes(router, db)
    
    return router
}
```

### HTTP Test Helpers
```go
package testutils

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    
    "github.com/gin-gonic/gin"
)

func MakeRequest(router *gin.Engine, method, path string, payload interface{}) *httptest.ResponseRecorder {
    var req *http.Request
    
    if payload != nil {
        jsonData, _ := json.Marshal(payload)
        req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonData))
        req.Header.Set("Content-Type", "application/json")
    } else {
        req = httptest.NewRequest(method, path, nil)
    }
    
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)
    
    return recorder
}

func AssertJSONResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedStatus int, expectedKey string, expectedValue interface{}) {
    t.Helper()
    
    // Check status code
    if recorder.Code != expectedStatus {
        t.Errorf("Expected status %d, got %d", expectedStatus, recorder.Code)
    }
    
    // Parse JSON response
    var response map[string]interface{}
    if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse JSON response: %v", err)
    }
    
    // Check specific field
    if actualValue, exists := response[expectedKey]; exists {
        if actualValue != expectedValue {
            t.Errorf("Expected %s = %v, got %v", expectedKey, expectedValue, actualValue)
        }
    } else {
        t.Errorf("Expected key '%s' not found in response", expectedKey)
    }
}
```

## Performance Considerations

### Test Execution Time
- Integration tests should complete in under 10 seconds
- Use proper indexing for database queries
- Optimize test data creation
- Avoid unnecessary external API calls

### Database Optimization
- Use connection pooling
- Optimize query performance
- Use proper indexing
- Clean up test data efficiently

### Memory Management
- Close database connections properly
- Clean up resources after each test
- Avoid memory leaks in test code

## Troubleshooting

### Common Issues

#### Database Connection Problems
- Ensure PostgreSQL is running
- Check database connection string
- Verify database permissions
- Check if test database exists

#### Test Data Issues
- Ensure proper test data cleanup
- Check for duplicate test data
- Verify test data integrity
- Handle foreign key constraints properly

#### Environment Setup Issues
- Check environment variables
- Verify service dependencies
- Ensure proper port availability
- Check network connectivity

### Debugging Tests
```bash
# Run tests with verbose output
go test -v ./internal/test/integration/...

# Run tests with race detection
go test -race ./internal/test/integration/...

# Run specific test with debugging
go test -run TestSpecificFunction -v ./internal/test/integration/...

# Run tests with database logging
go test -v ./internal/test/integration/... -args -db-log=true

# Run tests with slow threshold
go test -v ./internal/test/integration/... -args -slow-threshold=5s
```

## Contributing

### Adding New Integration Tests
1. Follow the existing naming conventions
2. Use the provided test templates
3. Include complete test scenarios
4. Ensure proper test data management
5. Update coverage reports

### Code Review Checklist
- [ ] Tests follow naming conventions
- [ ] Tests use real dependencies appropriately
- [ ] Test data is properly managed
- [ ] Tests cover complete scenarios
- [ ] Coverage requirements are met
- [ ] Tests are properly documented

## Additional Resources

### Documentation
- [Go Testing Documentation](https://golang.org/doc/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [PostgreSQL Testing](https://www.postgresql.org/docs/current/testing.html)

### Tools
- [Go Test Coverage](https://pkg.go.dev/cmd/go#hdr-Testing_flags)
- [Testify/assert](https://github.com/stretchr/testify/assert)
- [pgAdmin](https://www.pgadmin.org/) for database management

### External Links
- [Integration Testing Best Practices](https://martinfowler.com/bliki/IntegrationTest.html)
- [Database Testing Strategies](https://testing.googleblog.com/2017/07/testing-database-applications.html)
- [API Testing Guidelines](https://restfulapi.net/testing-apis/)

---

**Note**: This README should be updated as the integration test suite evolves and new testing practices are adopted.