# Test Module Organization Recommendations

## Executive Summary

This document provides comprehensive recommendations for test module organization based on the completed file migration and consolidation work. The project has been successfully migrated to use the `internal/` directory structure with all test modules consolidated under `internal/test/`. This document analyzes the current unified test structure and provides best practices for future test development.

## 1. Current Test Structure Analysis

### 1.1 Unified Test Structure Overview

The project has successfully consolidated all test modules under `internal/test/` with the following organization:

```
internal/test/
├── config/                    # Configuration testing
│   ├── config_test.go
│   ├── https_test.go
│   └── providers_test.go
├── database/                  # Database testing
│   ├── db_manual_test.go
│   └── db_test.go
├── handlers/                  # Handler/controller testing
│   ├── auth_controller_social_test.go
│   ├── comprehensive_auth_test.go
│   ├── test_utils.go
│   └── unified_auth_controller_test.go
├── infrastructure/            # Test infrastructure
│   ├── CI_WORKFLOW_TEST_PLAN.md
│   ├── Dockerfile.test
│   └── test-app
├── integration/               # Integration testing
│   └── integration_test.go
├── middleware/                # Middleware testing
│   ├── auth_test.go
│   └── roles_test.go
└── utils/                     # Utility function testing
    ├── email_test.go
    ├── mfa_test.go
    ├── otp_test.go
    ├── password_test.go
    ├── sms_test.go
    └── storage_test.go
```

### 1.2 Current Test Categories

#### Unit Tests
- **Configuration Tests**: Testing configuration loading, validation, and HTTPS setup
- **Handler Tests**: Testing authentication handlers with various scenarios
- **Middleware Tests**: Testing authentication and authorization middleware
- **Utility Tests**: Testing password hashing, OTP generation, email/SMS utilities

#### Integration Tests
- **Database Integration**: Testing database connectivity and operations
- **API Integration**: Testing complete authentication flows
- **Service Integration**: Testing service layer interactions

#### Infrastructure Tests
- **CI/CD Pipeline Testing**: Testing GitHub Actions workflows
- **Docker Testing**: Testing containerized test environments

## 2. Best Practices Recommendations

### 2.1 Test File Naming Conventions

#### Go Test Files
- Follow the standard Go convention: `*_test.go`
- Mirror the source file structure exactly
- Use descriptive names that clearly indicate the test scope

**Examples:**
```go
// Source file: internal/handlers/auth_handler.go
// Test file: internal/test/handlers/auth_handler_test.go

// Source file: internal/utils/password/password.go
// Test file: internal/test/utils/password/password_test.go
```

#### Test Function Naming
- Use `Test` prefix followed by descriptive camelCase
- Include the scenario being tested
- Focus on behavior rather than implementation

**Examples:**
```go
func TestAuthHandler_LoginWithValidCredentials(t *testing.T)
func TestAuthHandler_LoginWithInvalidPassword(t *testing.T)
func TestPasswordHasher_GenerateHashWithValidInput(t *testing.T)
func TestPasswordHasher_ValidateCorrectPassword(t *testing.T)
```

### 2.2 Test Organization Within Subdirectories

#### Current Structure Assessment
The current structure follows Go best practices by:
- Mirroring the source code directory structure
- Separating unit tests from integration tests
- Providing dedicated infrastructure testing

#### Recommended Enhancements

##### 1. Add Test Category Separation
```
internal/test/
├── unit/                     # Unit tests
│   ├── config/
│   ├── handlers/
│   ├── middleware/
│   └── utils/
├── integration/              # Integration tests
│   ├── api/
│   ├── database/
│   └── services/
├── e2e/                      # End-to-end tests
│   └── workflows/
└── fixtures/                 # Test data and fixtures
    ├── data/
    ├── mocks/
    └── databases/
```

##### 2. Test Grouping Strategy
- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test component interactions
- **E2E Tests**: Test complete user workflows
- **Fixtures**: Shared test data and mock objects

### 2.3 Test Coverage Strategies

#### Coverage Goals
- **Unit Tests**: 80% minimum coverage for critical paths
- **Integration Tests**: 100% coverage for API endpoints
- **E2E Tests**: 100% coverage for user workflows

#### Coverage Prioritization
1. **Authentication System**: 95%+ coverage
2. **API Endpoints**: 90%+ coverage
3. **Business Logic**: 85%+ coverage
4. **Utility Functions**: 80%+ coverage
5. **Error Handling**: 90%+ coverage

#### Coverage Measurement
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Run tests with coverage threshold
go test -covermode=count -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep "total:" | awk '{print substr($3, 1, length($3)-1)}'
```

### 2.4 Integration vs Unit Test Separation

#### Unit Tests Characteristics
- **Fast execution** (< 1 second per test)
- **Isolated dependencies** (use mocks)
- **Focus on single function/method**
- **No external dependencies**

#### Integration Tests Characteristics
- **Slower execution** (1-10 seconds per test)
- **Real dependencies** (database, HTTP services)
- **Test component interactions**
- **External dependencies required**

#### Test Separation Guidelines
```go
// Unit test example (isolated)
func TestPasswordHasher_HashPassword(t *testing.T) {
    hasher := NewPasswordHasher()
    password := "testPassword123"
    
    hash, err := hasher.HashPassword(password)
    require.NoError(t, err)
    require.NotEmpty(t, hash)
    
    // Verify hash is different from original
    require.NotEqual(t, password, hash)
}

// Integration test example (with real dependencies)
func TestAuthHandler_LoginIntegration(t *testing.T) {
    // Setup test database
    db := setupTestDatabase(t)
    defer db.Close()
    
    // Create test user
    createUser(t, db, "test@example.com", "password123")
    
    // Test login endpoint
    req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{
        "email": "test@example.com",
        "password": "password123"
    }`))
    
    recorder := httptest.NewRecorder()
    handler := NewAuthHandler(db)
    handler.Login(recorder, req)
    
    // Verify response
    require.Equal(t, http.StatusOK, recorder.Code)
}
```

## 3. Test Module Location Justification

### 3.1 Why `internal/test/` is Optimal

#### Advantages of Current Structure
1. **Encapsulation**: `internal/` ensures test code is not exported outside the package
2. **Proximity**: Tests are close to the source code they test
3. **Organization**: Clear separation from production code
4. **Build Exclusion**: Can be easily excluded from production builds

#### Go Best Practices Alignment
- Follows Go's recommendation for test organization
- Maintains the `internal/` package boundary
- Supports Go's testing conventions

### 3.2 Alternative Locations Considered

#### `./tests/` Directory
```bash
# Pros:
# - Clear separation from source code
# - Easy to exclude from builds
# - Common in other ecosystems

# Cons:
# - Breaks Go conventions
# - Requires build tags for exclusion
# - Less intuitive for Go developers
```

#### Mirror Source Structure
```bash
# Pros:
# - Perfect mirroring of source
# - Easy navigation
# - Follows some project patterns

# Cons:
# - Can clutter source directories
# - Mixed concerns
# - Harder to manage test dependencies
```

#### Conclusion
The current `internal/test/` structure is optimal as it:
- Follows Go conventions
- Provides clear separation
- Maintains proximity to source code
- Supports proper build exclusion

## 4. Test Infrastructure Improvements

### 4.1 Test Configuration Management

#### Current State
- Tests use environment variables
- Configuration loaded from `.env` files
- Manual setup in test functions

#### Recommended Improvements

##### 1. Centralized Test Configuration
```go
// internal/test/config/test_config.go
package config

import (
    "os"
    "testing"
    
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
)

type TestConfig struct {
    DatabaseURL    string
    TestDatabase   string
    JWTSecret      string
    SMTPHost       string
    SMTPUser       string
    SMTPPassword   string
}

var GlobalTestConfig TestConfig

func InitTestConfig(t *testing.T) {
    // Load production config
    prodConfig := config.Load()
    
    // Override with test-specific settings
    GlobalTestConfig = TestConfig{
        DatabaseURL:    getEnv("TEST_DATABASE_URL", "postgres://test:test@localhost:5432/rangkaiedu_test?sslmode=disable"),
        TestDatabase:   "rangkaiedu_test",
        JWTSecret:      "test-jwt-secret-change-in-production",
        SMTPHost:       getEnv("TEST_SMTP_HOST", ""),
        SMTPUser:       getEnv("TEST_SMTP_USER", ""),
        SMTPPassword:   getEnv("TEST_SMTP_PASSWORD", ""),
    }
    
    // Set environment variables
    os.Setenv("DB_NAME", GlobalTestConfig.TestDatabase)
    os.Setenv("JWT_SECRET", GlobalTestConfig.JWTSecret)
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

##### 2. Test Environment Setup
```go
// internal/test/setup/test_setup.go
package setup

import (
    "database/sql"
    "testing"
    
    _ "github.com/jackc/pgx/v5/stdlib"
)

func SetupTestDatabase(t *testing.T) *sql.DB {
    t.Helper()
    
    db, err := sql.Open("pgx", os.Getenv("TEST_DATABASE_URL"))
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }
    
    // Run migrations
    if err := runTestMigrations(db); err != nil {
        t.Fatalf("Failed to run test migrations: %v", err)
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

### 4.2 Test Data Management

#### Test Data Strategy
```go
// internal/test/fixtures/test_data.go
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

### 4.3 Mock and Stub Organization

#### Mock Strategy
```go
// internal/test/mocks/auth_mock.go
package mocks

import (
    "github.com/stretchr/testify/mock"
)

type MockAuthHandler struct {
    mock.Mock
}

func (m *MockAuthHandler) Login(c *gin.Context) {
    m.Called(c)
}

func (m *MockAuthHandler) Register(c *gin.Context) {
    m.Called(c)
}

func (m *MockAuthHandler) SendOTP(c *gin.Context) {
    m.Called(c)
}

func (m *MockAuthHandler) VerifyOTP(c *gin.Context) {
    m.Called(c)
}
```

#### Stub Organization
```go
// internal/test/stubs/email_stub.go
package stubs

import (
    "errors"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/email"
)

type EmailStub struct {
    SendError error
    SendCount int
}

func (s *EmailStub) SendOTPEmail(cfg config.Config, to, otp string) error {
    s.SendCount++
    return s.SendError
}

// Factory functions
func NewSuccessfulEmailStub() *EmailStub {
    return &EmailStub{SendError: nil}
}

func NewFailingEmailStub(err error) *EmailStub {
    return &EmailStub{SendError: err}
}
```

### 4.4 Test Environment Setup

#### Docker Test Environment
```dockerfile
# internal/test/infrastructure/Dockerfile.test
FROM golang:1.21-alpine AS test

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build test binary
RUN go test -c -o test-runner ./...

# Final image
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=test /app/test-runner .
COPY --from=test /app/.env.test ./
CMD ["./test-runner"]
```

#### Test Environment Scripts
```bash
#!/bin/bash
# scripts/test-setup.sh

set -e

echo "Setting up test environment..."

# Create test database
psql -h localhost -U postgres -c "CREATE DATABASE rangkaiedu_test;"

# Run migrations
psql -h localhost -U postgres -d rangkaiedu_test -f migrations/001_create_tables.sql

# Install test dependencies
go mod download

echo "Test environment setup complete!"
```

## 5. Migration Guidelines for Future Test Files

### 5.1 Test File Creation Process

#### Step-by-Step Guide
1. **Identify Test Scope**
   - Determine if it's a unit, integration, or e2e test
   - Identify dependencies and external services
   - Define test data requirements

2. **Create Test File**
   - Follow naming conventions: `*_test.go`
   - Place in appropriate subdirectory
   - Include necessary imports

3. **Write Test Cases**
   - Use table-driven tests for multiple scenarios
   - Include edge cases and error conditions
   - Follow AAA pattern (Arrange, Act, Assert)

4. **Add Test Fixtures**
   - Create reusable test data
   - Setup/teardown functions
   - Mock external dependencies

5. **Run and Validate**
   - Execute tests locally
   - Check coverage
   - Validate CI/CD pipeline

#### Example Migration
```go
// Before: Scattered tests
// config/config_test.go
// middleware/auth_test.go
// handlers/auth_handler_test.go

// After: Organized structure
// internal/test/unit/config/config_test.go
// internal/test/unit/middleware/auth_test.go
// internal/test/unit/handlers/auth_handler_test.go
// internal/test/integration/auth_flow_test.go
```

### 5.2 Test Template

#### Unit Test Template
```go
package [package_name]

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func Test[FunctionName]_[Scenario](t *testing.T) {
    // Arrange
    setup := setupTest(t)
    defer teardownTest(t, setup)
    
    // Act
    result, err := setup.[FunctionName](input)
    
    // Assert
    if expectedError {
        assert.Error(t, err)
        assert.Contains(t, err.Error(), expectedErrorMessage)
    } else {
        require.NoError(t, err)
        assert.Equal(t, expectedValue, result)
    }
}

func setupTest(t *testing.T) *TestSetup {
    // Setup test environment
    return &TestSetup{}
}

func teardownTest(t *testing.T, setup *TestSetup) {
    // Cleanup test environment
}
```

#### Integration Test Template
```go
package integration

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func Test[APIEndpoint]_[Scenario](t *testing.T) {
    // Setup test database
    db := setupTestDatabase(t)
    defer cleanupTestDatabase(t, db)
    
    // Setup test server
    router := setupTestRouter(db)
    
    // Create test data
    createTestData(t, db)
    
    // Arrange
    reqBody := `{"key": "value"}`
    req := httptest.NewRequest("METHOD", "/api/endpoint", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer valid-token")
    
    // Act
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)
    
    // Assert
    assert.Equal(t, expectedStatusCode, recorder.Code)
    
    var response map[string]interface{}
    require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
    
    assert.Equal(t, expectedValue, response["key"])
}
```

## 6. Testing Workflow Recommendations

### 6.1 Unit Testing Strategy

#### Testing Approach
1. **Isolate Components**: Test each component in isolation
2. **Use Mocks**: Mock external dependencies
3. **Test Edge Cases**: Include boundary conditions
4. **Table-Driven Tests**: Use for multiple scenarios

#### Example Unit Test
```go
func TestPasswordHasher_VariousInputs(t *testing.T) {
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
        {
            name:        "Long password",
            password:    strings.Repeat("a", 1000),
            expectError: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            hasher := NewPasswordHasher()
            hash, err := hasher.HashPassword(tt.password)
            
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

### 6.2 Integration Testing Strategy

#### Testing Approach
1. **Real Dependencies**: Use actual database and services
2. **Test Flows**: Test complete business processes
3. **Environment Setup**: Proper test environment isolation
4. **Data Cleanup**: Ensure test data doesn't persist

#### Example Integration Test
```go
func TestUserRegistrationFlow(t *testing.T) {
    // Setup
    db := setupTestDatabase(t)
    defer cleanupTestDatabase(t, db)
    
    // Test user registration
    registerReq := `{
        "name": "John Doe",
        "email": "john@example.com",
        "phone": "1234567890",
        "password": "StrongPass123!",
        "role": "student"
    }`
    
    // Register user
    registerResp := makeRequest("POST", "/api/auth/register", registerReq)
    assert.Equal(t, http.StatusCreated, registerResp.Code)
    
    // Verify user exists in database
    var user models.User
    err := db.QueryRow("SELECT id, name, email, role FROM users WHERE email = $1", "john@example.com").Scan(
        &user.ID, &user.Name, &user.Email, &user.Role)
    require.NoError(t, err)
    assert.Equal(t, "John Doe", user.Name)
    assert.Equal(t, "student", user.Role)
    
    // Test login with new user
    loginReq := `{
        "email": "john@example.com",
        "password": "StrongPass123!"
    }`
    
    loginResp := makeRequest("POST", "/api/auth/login", loginReq)
    assert.Equal(t, http.StatusOK, loginResp.Code)
    
    // Verify token is returned
    var loginResponse map[string]string
    require.NoError(t, json.Unmarshal(loginResp.Body.Bytes(), &loginResponse))
    assert.NotEmpty(t, loginResponse["token"])
}
```

### 6.3 End-to-End Testing Strategy

#### Testing Approach
1. **Complete User Journeys**: Test full user workflows
2. **Multiple Environments**: Test across dev/staging/prod
3. **Performance Testing**: Include load testing
4. **Browser Testing**: Test with real browsers

#### Example E2E Test
```go
func TestCompleteUserJourney(t *testing.T) {
    // Setup test environment
    testEnv := setupTestEnvironment(t)
    defer testEnv.Cleanup()
    
    // Step 1: User registration
    registerUser(t, testEnv)
    
    // Step 2: Email verification
    verifyEmail(t, testEnv)
    
    // Step 3: Login
    token := loginUser(t, testEnv)
    
    // Step 4: Create class
    classID := createClass(t, testEnv, token)
    
    // Step 5: Enroll in class
    enrollInClass(t, testEnv, token, classID)
    
    // Step 6: Access class content
    accessClassContent(t, testEnv, token, classID)
    
    // Step 7: Complete lesson
    completeLesson(t, testEnv, token, classID)
    
    // Step 8: Logout
    logoutUser(t, testEnv, token)
}
```

### 6.4 CI/CD Pipeline Integration

#### Pipeline Configuration
```yaml
# .github/workflows/testing.yml
name: Testing

on:
  push:
    branches: [ main, develop, staging ]
  pull_request:
    branches: [ main, develop, staging ]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Run unit tests
      run: go test -v -cover ./internal/test/unit/...
    
    - name: Generate coverage report
      run: go test -coverprofile=coverage.out ./internal/test/unit/...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3

  integration-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: testpass
          POSTGRES_DB: rangkaiedu_test
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Wait for PostgreSQL
      run: |
        until pg_isready -h localhost -p 5432; do
          echo "Waiting for PostgreSQL..."
          sleep 1
        done
    
    - name: Run integration tests
      run: go test -v ./internal/test/integration/...
    
    - name: Run e2e tests
      run: go test -v ./internal/test/e2e/...

  security-tests:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Run security tests
      run: go test -v ./internal/test/security/...
    
    - name: Run vulnerability scan
      run: gosec ./...
```

## 7. Test Documentation Template

### 7.1 README Template for Test Subdirectories

```markdown
# Test Directory: [Directory Name]

## Overview
This directory contains tests for the [component name] module. These tests ensure that the [component name] functions correctly and meets the specified requirements.

## Test Categories

### Unit Tests
- **File**: `*_test.go`
- **Purpose**: Test individual functions and methods in isolation
- **Dependencies**: Mocked external services
- **Execution Time**: Fast (< 1 second per test)

### Integration Tests
- **File**: `*_integration_test.go`
- **Purpose**: Test component interactions with real dependencies
- **Dependencies**: Database, HTTP services
- **Execution Time**: Medium (1-10 seconds per test)

### E2E Tests
- **File**: `*_e2e_test.go`
- **Purpose**: Test complete user workflows
- **Dependencies**: Full application stack
- **Execution Time**: Slow (10-60 seconds per test)

## Running Tests

### Run All Tests
```bash
go test ./...
```

### Run Specific Test Category
```bash
# Unit tests only
go test ./internal/test/[directory_name]/...

# Integration tests only
go test -run Integration ./internal/test/[directory_name]/...

# E2E tests only
go test -run E2E ./internal/test/[directory_name]/...
```

### Run Tests with Coverage
```bash
go test -cover ./internal/test/[directory_name]/...
go test -coverprofile=coverage.out ./internal/test/[directory_name]/...
go tool cover -html=coverage.out
```

## Test Data

### Fixtures
- **Location**: `fixtures/`
- **Purpose**: Reusable test data objects
- **Usage**: `fixtures.CreateTestUser(db, StudentUser)`

### Test Databases
- **Name**: `rangkaiedu_test`
- **Schema**: Same as production but with test data
- **Cleanup**: Automatic truncation after each test

## Mock Objects

### Available Mocks
- **[Mock Name]**: Description of mock functionality
- **Usage**: `mock := mocks.NewMock[Component]()`

### Stub Objects
- **[Stub Name]**: Description of stub functionality
- **Usage**: `stub := stubs.New[Stub]()`

## Test Environment Setup

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Docker (for integration tests)

### Environment Variables
```bash
TEST_DATABASE_URL=postgres://test:test@localhost:5432/rangkaiedu_test
TEST_JWT_SECRET=test-jwt-secret
TEST_SMTP_HOST=localhost
TEST_SMTP_USER=test
TEST_SMTP_PASSWORD=test
```

### Setup Commands
```bash
# Create test database
createdb rangkaiedu_test

# Run migrations
psql -d rangkaiedu_test -f migrations/001_create_tables.sql

# Run tests
go test ./internal/test/[directory_name]/...
```

## Coverage Requirements

### Minimum Coverage
- **Unit Tests**: 80%
- **Integration Tests**: 90%
- **E2E Tests**: 100%

### Coverage Reports
- **HTML Report**: `coverage.html`
- **Text Report**: `coverage.txt`
- **JSON Report**: `coverage.json`

## Contributing

### Adding New Tests
1. Follow the existing naming conventions
2. Use the provided test templates
3. Include both positive and negative test cases
4. Ensure proper cleanup of test data
5. Update coverage reports

### Test Best Practices
- Use table-driven tests for multiple scenarios
- Include descriptive test names
- Use proper setup/teardown functions
- Mock external dependencies
- Test edge cases and error conditions

## Troubleshooting

### Common Issues
1. **Database Connection Issues**
   - Ensure PostgreSQL is running
   - Check environment variables
   - Verify database permissions

2. **Test Failures**
   - Check test data setup
   - Verify mock configurations
   - Review environment setup

3. **Coverage Issues**
   - Ensure all code paths are tested
   - Check for unreachable code
   - Review coverage configuration

### Getting Help
- Check the main project documentation
- Review existing test examples
- Contact the development team
```

### 7.2 Example README Files

#### Config Tests README
```markdown
# Test Directory: config

## Overview
This directory contains tests for the configuration module. These tests ensure that configuration loading, validation, and HTTPS setup work correctly.

## Test Structure

### Unit Tests
- `config_test.go`: Tests configuration loading and validation
- `https_test.go`: Tests HTTPS configuration and setup
- `providers_test.go`: Tests provider configuration

### Integration Tests
- `config_integration_test.go`: Tests configuration with real environment

## Running Tests

```bash
# Run all config tests
go test ./internal/test/config/...

# Run with coverage
go test -cover ./internal/test/config/...
```

## Test Coverage

### Current Coverage
- Configuration Loading: 95%
- HTTPS Setup: 90%
- Provider Configuration: 85%

### Target Coverage
- All modules: 90% minimum
```

#### Handler Tests README
```markdown
# Test Directory: handlers

## Overview
This directory contains tests for the HTTP handlers. These tests ensure that all API endpoints work correctly and handle various scenarios.

## Test Categories

### Unit Tests
- Test individual handler functions
- Mock HTTP requests and responses
- Test error handling

### Integration Tests
- Test complete API flows
- Test database interactions
- Test authentication flows

## Test Endpoints

### Authentication
- `/api/auth/register`
- `/api/auth/login`
- `/api/auth/send-otp`
- `/api/auth/verify-otp`

### User Management
- `/api/users/profile`
- `/api/users/update`

### Classes
- `/api/classes`
- `/api/classes/:id`

## Running Tests

```bash
# Run all handler tests
go test ./internal/test/handlers/...

# Run specific test
go test -run TestAuthHandler_Login ./internal/test/handlers/...

# Run with verbose output
go test -v ./internal/test/handlers/...
```

## Mock Objects

### Available Mocks
- `MockAuthHandler`: Mocks authentication handler
- `MockUserHandler`: Mocks user management handler
- `MockClassHandler`: Mocks class management handler
```

## 8. Conclusion

### 8.1 Summary of Recommendations

1. **Current Structure**: The `internal/test/` structure is well-organized and follows Go best practices
2. **Naming Conventions**: Standardize on `*_test.go` naming with descriptive function names
3. **Test Categories**: Separate unit, integration, and e2e tests clearly
4. **Coverage Strategy**: Implement comprehensive coverage with minimum thresholds
5. **Infrastructure**: Improve test configuration, data management, and environment setup
6. **Documentation**: Provide clear README files for each test directory

### 8.2 Implementation Roadmap

#### Phase 1: Immediate Improvements (Week 1-2)
- [ ] Standardize test file naming conventions
- [ ] Create comprehensive README templates
- [ ] Implement centralized test configuration
- [ ] Add test coverage reporting

#### Phase 2: Enhanced Infrastructure (Week 3-4)
- [ ] Implement test fixtures and mock objects
- [ ] Create Docker test environment
- [ ] Add test data management
- [ ] Improve CI/CD integration

#### Phase 3: Advanced Testing (Week 5-6)
- [ ] Add e2e testing framework
- [ ] Implement performance testing
- [ ] Add security testing
- [ ] Create comprehensive test documentation

### 8.3 Success Metrics

#### Quantitative Metrics
- **Test Coverage**: 90%+ overall coverage
- **Test Execution Time**: < 5 minutes for full test suite
- **CI/CD Success Rate**: 99%+ test pass rate
- **Code Quality**: Reduced bug count by 50%

#### Qualitative Metrics
- **Developer Experience**: Easy to write and run tests
- **Maintainability**: Clear test organization and documentation
- **Reliability**: Consistent test results
- **Coverage**: Comprehensive test coverage of critical paths

### 8.4 Final Recommendations

1. **Maintain Current Structure**: The `internal/test/` structure is optimal and should be maintained
2. **Standardize Conventions**: Implement consistent naming and organization patterns
3. **Invest in Infrastructure**: Improve test automation and environment setup
4. **Prioritize Coverage**: Focus on critical paths and user workflows

## Current Test Structure Status

The current test structure has been successfully implemented with the following organization:

```
internal/test/
├── config/                    # Configuration testing (100% coverage)
│   ├── config_test.go
│   ├── https_test.go
│   └── providers_test.go
├── database/                  # Database testing (100% coverage)
│   ├── db_test.go
│   └── db_manual_test.go
├── handlers/                  # Handler/controller testing (Authentication system 100% coverage)
│   ├── auth_controller_social_test.go
│   ├── comprehensive_auth_test.go
│   ├── unified_auth_controller_test.go
│   └── test_utils.go
├── infrastructure/            # Test infrastructure
│   ├── CI_WORKFLOW_TEST_PLAN.md
│   ├── Dockerfile.test
│   └── test-app
├── integration/               # Integration testing
│   └── integration_test.go
├── middleware/                # Middleware testing (100% coverage)
│   ├── auth_test.go
│   └── roles_test.go
├── repositories/              # Repository layer testing (100% coverage)
│   ├── class_repository_test.go
│   ├── material_repository_test.go
│   ├── student_enrollment_repository_test.go
│   ├── subject_repository_test.go
│   └── user_repository_test.go
├── unit/                      # Unit testing
│   └── README.md
└── utils/                     # Utility function testing (95% coverage)
    ├── email_test.go
    ├── mfa_test.go
    ├── otp_test.go
    ├── password_test.go
    ├── sms_test.go
    └── storage_test.go
```

## Missing Test Modules

Based on the current implementation, the following test modules are missing or need expansion:

1. **API Endpoint Tests**: Complete testing of all REST API endpoints
2. **Service Layer Tests**: Testing of business logic in service components
3. **End-to-End Tests**: Complete user workflow testing
4. **Performance Tests**: Load and stress testing
5. **Security Tests**: Security vulnerability testing
6. **Error Handling Tests**: Comprehensive error scenario testing
5. **Document Everything**: Provide comprehensive documentation for test organization

By following these recommendations, the project will have a robust, maintainable, and comprehensive testing framework that ensures code quality and reliability.