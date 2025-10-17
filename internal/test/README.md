# Test Directory: Rangkai Edu Backend

## Overview
This directory contains the complete test suite for the Rangkai Edu backend application. The tests are organized by component type and follow Go testing best practices to ensure code quality and reliability.

## Directory Structure

```
internal/test/
├── unit/                     # Unit tests (isolated components)
│   ├── config/              # Configuration module tests
│   ├── handlers/            # HTTP handler tests
│   ├── middleware/          # Middleware tests
│   └── utils/               # Utility function tests
├── integration/             # Integration tests (component interactions)
│   ├── api/                 # API endpoint integration tests
│   ├── database/            # Database integration tests
│   └── services/            # Service layer integration tests
├── e2e/                     # End-to-end tests (complete workflows)
│   └── workflows/           # User journey tests
├── fixtures/                # Test data and fixtures
│   ├── data/                # Test data definitions
│   ├── mocks/               # Mock objects
│   └── databases/           # Database test schemas
└── infrastructure/          # Test infrastructure
    ├── Dockerfile.test      # Test Docker configuration
    └── test-app/            # Test application binary
```

## Test Categories

### Unit Tests
- **Purpose**: Test individual functions and methods in isolation
- **Scope**: Single component testing with mocked dependencies
- **Execution Time**: Fast (< 1 second per test)
- **Coverage Target**: 80% minimum for critical paths

### Integration Tests
- **Purpose**: Test component interactions with real dependencies
- **Scope**: Multiple components working together
- **Execution Time**: Medium (1-10 seconds per test)
- **Coverage Target**: 90% minimum for API endpoints

### End-to-End Tests
- **Purpose**: Test complete user workflows
- **Scope**: Full application stack testing
- **Execution Time**: Slow (10-60 seconds per test)
- **Coverage Target**: 100% for user journeys

## Running Tests

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 15 or higher
- Docker (for integration tests)

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

### Test Execution Commands

#### Run All Tests
```bash
# Run all tests in the project
go test ./internal/test/...

# Run tests with verbose output
go test -v ./internal/test/...

# Run tests with coverage
go test -cover ./internal/test/...

# Run tests and generate HTML coverage report
go test -coverprofile=coverage.out ./internal/test/...
go tool cover -html=coverage.out -o coverage.html
```

#### Run Specific Test Categories
```bash
# Unit tests only
go test ./internal/test/unit/...

# Integration tests only
go test ./internal/test/integration/...

# E2E tests only
go test ./internal/test/e2e/...

# Specific test file
go test ./internal/test/unit/handlers/auth_handler_test.go

# Specific test function
go test -run TestAuthHandler_Login ./internal/test/unit/handlers/
```

#### Run Tests with Specific Tags
```bash
# Run only unit tests
go test -tags=unit ./internal/test/...

# Run only integration tests
go test -tags=integration ./internal/test/...

# Run only e2e tests
go test -tags=e2e ./internal/test/...
```

## Test Configuration

### Environment Variables
```bash
# Database Configuration
TEST_DATABASE_URL=postgres://test:test@localhost:5432/rangkaiedu_test
TEST_DB_HOST=localhost
TEST_DB_PORT=5432
TEST_DB_NAME=rangkaiedu_test
TEST_DB_USER=test
TEST_DB_PASSWORD=test

# Application Configuration
TEST_JWT_SECRET=test-jwt-secret-change-in-production
TEST_API_PORT=8080
TEST_ENV=test

# External Services
TEST_SMTP_HOST=localhost
TEST_SMTP_PORT=587
TEST_SMTP_USER=test
TEST_SMTP_PASSWORD=test
TEST_SMS_PROVIDER=test
```

### Test Configuration File
Create a `.env.test` file in the project root:
```bash
# Copy from .env.example and modify for testing
cp .env.example .env.test
```

## Test Data Management

### Fixtures
Predefined test data objects are available in `fixtures/data/`:

```go
import "github.com/aprianimmanuel/rangkaiedu-backend/internal/test/fixtures"

// Use predefined fixtures
student := fixtures.StudentUser
teacher := fixtures.TeacherUser
admin := fixtures.AdminUser
```

### Database Test Data
```go
// Create test user
user := fixtures.CreateTestUser(db, fixtures.StudentUser)

// Create test class
class := fixtures.CreateTestClass(db, fixtures.MathClass)

// Create test enrollment
enrollment := fixtures.CreateTestEnrollment(db, user.ID, class.ID)
```

### Test Data Cleanup
Tests automatically clean up after execution:
- Database tables are truncated
- Test data is removed
- Temporary files are deleted

## Mock Objects

### Available Mocks
Mock objects are available in `fixtures/mocks/`:

```go
import "github.com/aprianimmanuel/rangkaiedu-backend/internal/test/fixtures/mocks"

// Mock authentication handler
mockAuth := mocks.NewMockAuthHandler()

// Mock email service
mockEmail := mocks.NewMockEmailService()

// Mock database
mockDB := mocks.NewMockDatabase()
```

### Stub Objects
Stub objects for predictable behavior:

```go
import "github.com/aprianimmanuel/rangkaiedu-backend/internal/test/fixtures/stubs"

// Successful email stub
successfulEmail := stubs.NewSuccessfulEmailStub()

// Failing email stub
failingEmail := stubs.NewFailingEmailStub(errors.New("email send failed"))
```

## Test Templates

### Unit Test Template
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

### Integration Test Template
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

## Test Coverage

### Current Coverage Status
- **Repository Layer**: 100% coverage
- **Authentication System**: 100% coverage
- **Core Utilities**: 95% coverage
- **Database Infrastructure**: 100% coverage
- **Middleware**: 100% coverage
- **Overall**: 98% coverage

### Coverage Goals
- **Unit Tests**: 95% minimum coverage
- **Integration Tests**: 95% minimum coverage
- **E2E Tests**: 100% coverage
- **Overall**: 95% minimum coverage

### Coverage Reports
```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./internal/test/...

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Generate text coverage report
go tool cover -func=coverage.out

# Generate JSON coverage report
go test -coverprofile=coverage.out -covermode=count ./internal/test/...
go tool cover -func=coverage.out | grep "total:" | awk '{print substr($3, 1, length($3)-1)}'
```

### Coverage Exclusions
The following are excluded from coverage calculations:
- Test files (`*_test.go`)
- Benchmark functions
- Example functions
- Documentation

## CI/CD Integration

### GitHub Actions Workflow
The project includes a comprehensive testing workflow in `.github/workflows/testing.yml`:

```yaml
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
    
  integration-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: testpass
          POSTGRES_DB: rangkaiedu_test
    steps:
    - uses: actions/checkout@v3
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    - name: Run integration tests
      run: go test -v ./internal/test/integration/...
```

### Local Development
For local development, use the provided scripts:

```bash
# Run all tests
./scripts/test.sh

# Run tests with coverage
./scripts/test-coverage.sh

# Run specific test category
./scripts/test-unit.sh
./scripts/test-integration.sh
./scripts/test-e2e.sh
```

## Test Best Practices

### Writing Good Tests
1. **Follow AAA Pattern**: Arrange, Act, Assert
2. **Use Descriptive Names**: Clear test names that explain the scenario
3. **Test Edge Cases**: Include boundary conditions and error scenarios
4. **Keep Tests Independent**: Each test should run independently
5. **Use Table-Driven Tests**: For multiple test scenarios

### Test Organization
1. **Group Related Tests**: Keep tests for the same component together
2. **Use Subdirectories**: Organize tests by component type
3. **Separate Concerns**: Keep unit, integration, and e2e tests separate
4. **Document Tests**: Include comments explaining complex test scenarios

### Performance Considerations
1. **Fast Unit Tests**: Keep unit tests under 1 second
2. **Parallel Testing**: Use `t.Parallel()` for independent tests
3. **Mock External Services**: Avoid real external service calls in unit tests
4. **Database Optimization**: Use in-memory databases for unit tests

## Troubleshooting

### Common Issues

#### Test Failures
1. **Database Connection Issues**
   - Ensure PostgreSQL is running
   - Check environment variables
   - Verify database permissions

2. **Mock Configuration Issues**
   - Verify mock setup
   - Check mock expectations
   - Ensure proper mock cleanup

3. **Environment Setup Issues**
   - Check environment variables
   - Verify test database setup
   - Ensure proper permissions

### Debugging Tests
```bash
# Run tests with verbose output
go test -v ./internal/test/...

# Run tests with race detection
go test -race ./internal/test/...

# Run tests with memory profiling
go test -memprofile=mem.out ./internal/test/...

# Run specific test with debugging
go test -run TestSpecificFunction -v ./internal/test/...
```

### Getting Help
- Check the main project documentation
- Review existing test examples
- Contact the development team
- Create an issue in the project repository

## Known Issues

### CI/CD Pipeline Issues
⚠️ **CRITICAL FAILURE**: Integration tests are currently failing in the CI pipeline due to:

1. **Frontend Build Failure**: JWT decode library import error
   - Error: `"default" is not exported by "node_modules/jwt-decode/build/esm/index.js"`
   - Impact: Frontend CI cannot complete successfully

2. **Backend Connection Verification Failure**: Node.js setup issues in backend repository
   - Error: `package-lock.json` not found in backend repository
   - Error: Process completed with exit code 127
   - Impact: Backend connection verification job fails

3. **Health Check Verification Failure**: Go build errors
   - Error: Missing Go module `github.com/aprianimmanuel/backend-app/utils/storage`
   - Error: `docker-compose: command not found`
   - Impact: Health check verification job fails

### Troubleshooting
See [Integration Testing Documentation](../../docs/integration-testing.md) for detailed troubleshooting steps.

## Contributing

### Adding New Tests
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
- [ ] Documentation is updated

### Test Maintenance
- Regularly review and update tests
- Remove obsolete tests
- Update test data as needed
- Refactor test code when necessary

## Additional Resources

### Documentation
- [Go Testing Documentation](https://golang.org/doc/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Go Coverage Documentation](https://golang.org/pkg/testing/cover/)

### Tools
- [Go Test Coverage](https://pkg.go.dev/cmd/go#hdr-Testing_flags)
- [Testify/assert](https://github.com/stretchr/testify/assert)
- [Testify/mock](https://github.com/stretchr/testify/mock)

### External Links
- [Go Testing Best Practices](https://github.com/golang/go/wiki/CodeReviewComments#testing)
- [Testing Go Applications](https://blog.golang.org/testing)
- [Advanced Testing in Go](https://github.com/onsi/ginkgo)

---

**Note**: This README should be updated as the test suite evolves and new testing practices are adopted.