# Backend Testing Framework Setup

## Overview
This document outlines the setup and configuration of the testing framework for the Rangkai Edu backend application using Go's built-in testing package.

## Go Testing Package Setup

The Go programming language comes with a powerful built-in testing package that provides support for automated testing of Go packages. The testing framework is designed to be simple and effective for unit testing, benchmarking, and example functions.

### Directory Structure
```
backend/
├── main.go
├── go.mod
├── go.sum
├── controllers/
├── models/
├── routes/
├── middleware/
├── config/
├── utils/
└── tests/
    ├── controllers/
    ├── models/
    ├── routes/
    ├── middleware/
    ├── config/
    └── utils/
```

### Test File Naming Convention
Go tests follow a specific naming convention where test files must end with `_test.go`. For example:
- `user_test.go` for testing `user.go`
- `auth_test.go` for testing `auth.go`

### Writing Tests
Tests in Go are written as functions that begin with the word `Test` followed by a capitalized word or phrase. Each test function takes a single parameter of type `*testing.T`.

Example test structure:
```go
package main

import (
    "testing"
)

func TestFunctionName(t *testing.T) {
    // Test implementation
    // Use t.Error() or t.Errorf() to report failures
    // Use t.Log() or t.Logf() for logging
}
```

### Running Tests
Tests can be run using the `go test` command:
```bash
# Run all tests in the current package
go test

# Run all tests in the project
go test ./...

# Run tests with verbose output
go test -v

# Run tests and generate coverage report
go test -cover
```

## Test Directory Structure

The backend testing framework will use the following directory structure:

1. **Unit Tests**: Tests for individual functions and methods
2. **Integration Tests**: Tests that verify the interaction between components
3. **API Tests**: Tests that verify the API endpoints

### Test Categories

#### 1. Unit Tests
Unit tests focus on testing individual functions and methods in isolation. These tests should:
- Be fast to execute
- Not depend on external services
- Use mocks or stubs for dependencies
- Test edge cases and error conditions

#### 2. Integration Tests
Integration tests verify that different components of the application work together correctly. These tests should:
- Test the interaction between components
- Use real dependencies where possible
- Verify data flow between layers

#### 3. API Tests
API tests verify that the REST endpoints work as expected. These tests should:
- Test all HTTP methods (GET, POST, PUT, DELETE)
- Verify correct status codes
- Validate response formats
- Test authentication and authorization

## Test Dependencies

The following dependencies will be used for testing:

1. **testify**: A toolkit with common assertions and mocks
   - Provides enhanced assertion capabilities
   - Includes mock object support
   - URL: https://github.com/stretchr/testify

2. **go-sqlmock**: SQL mock driver for testing
   - Mocks database operations
   - Verifies SQL queries
   - URL: https://github.com/DATA-DOG/go-sqlmock

## Test Configuration

Tests will use environment variables for configuration:
- `TEST_DATABASE_URL`: Database connection string for tests
- `TEST_JWT_SECRET`: JWT secret for testing authentication
- `TEST_ENV`: Environment identifier for tests

## Code Coverage

Code coverage will be measured using Go's built-in coverage tool:
```bash
# Generate coverage profile
go test -coverprofile=coverage.out

# Display coverage in HTML format
go tool cover -html=coverage.out
```

## Best Practices

1. **Table-Driven Tests**: Use table-driven tests for testing multiple scenarios
2. **Descriptive Test Names**: Use clear, descriptive names for test functions
3. **Setup and Teardown**: Use `TestMain` for setup and teardown when needed
4. **Parallel Testing**: Use `t.Parallel()` to run tests in parallel where possible
5. **Test Fixtures**: Use test fixtures for complex test data
6. **Mocking**: Use mocks for external dependencies

## Example Test Structure

```go
package controllers

import (
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
    // Setup
    req, _ := http.NewRequest("GET", "/users/1", nil)
    rr := httptest.NewRecorder()
    
    // Execute
    handler := http.HandlerFunc(GetUser)
    handler.ServeHTTP(rr, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, rr.Code)
    assert.Contains(t, rr.Body.String(), "user")
}
```

## Continuous Integration

Tests will be integrated into the CI/CD pipeline to ensure:
- All tests pass before merging
- Code coverage meets minimum requirements
- Performance benchmarks are maintained

## Next Steps

1. Create test directory structure
2. Set up test dependencies
3. Write sample tests for existing components
4. Configure CI/CD pipeline for automated testing