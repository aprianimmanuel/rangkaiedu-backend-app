# Test Infrastructure

## Overview
This directory contains test infrastructure components for the Rangkai Edu backend application. The infrastructure includes Docker configurations, test runners, CI/CD pipeline configurations, and other testing tools to ensure consistent and reliable test execution across different environments.

## Directory Structure

```
internal/test/infrastructure/
├── Dockerfile.test        # Test Docker configuration
├── docker-compose.test.yml # Test Docker Compose configuration
├── test-app/             # Test application binary
├── scripts/              # Test scripts and utilities
│   ├── test-setup.sh    # Test environment setup script
│   ├── test-run.sh      # Test execution script
│   ├── test-cleanup.sh  # Test cleanup script
│   └── test-coverage.sh # Coverage generation script
├── ci/                  # CI/CD pipeline configurations
│   ├── github-actions.yml # GitHub Actions workflow
│   └── gitlab-ci.yml    # GitLab CI configuration
└── templates/           # Infrastructure templates
    ├── test-config.env  # Test environment configuration
    └── test-docker.env  # Docker environment variables
```

## Infrastructure Components

### Docker Configuration
- **Dockerfile.test**: Containerized test environment
- **docker-compose.test.yml**: Multi-container test setup
- **Test Application**: Standalone test runner binary

### Test Scripts
- **Setup Script**: Configures test environment
- **Run Script**: Executes tests with various options
- **Cleanup Script**: Cleans up test resources
- **Coverage Script**: Generates coverage reports

### CI/CD Integration
- **GitHub Actions**: Automated testing workflow
- **GitLab CI**: Alternative CI/CD configuration
- **Environment Templates**: Reusable configuration templates

## Docker Configuration

### Test Dockerfile
```dockerfile
# Dockerfile.test
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build test binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o test-runner ./internal/test

# Final image
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy test binary
COPY --from=builder /app/test-runner .
COPY --from=builder /app/.env.test ./

# Copy test configuration
COPY --from=builder /app/internal/test/infrastructure/templates/ ./templates/

# Expose port (if needed)
EXPOSE 8080

# Run tests
CMD ["./test-runner"]
```

### Test Docker Compose
```yaml
# docker-compose.test.yml
version: '3.8'

services:
  # Test database
  test-db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: testuser
      POSTGRES_PASSWORD: testpass
      POSTGRES_DB: rangkaiedu_test
    ports:
      - "5432:5432"
    volumes:
      - test-db-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U testuser -d rangkaiedu_test"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Test application
  test-app:
    build:
      context: .
      dockerfile: internal/test/infrastructure/Dockerfile.test
    depends_on:
      test-db:
        condition: service_healthy
    environment:
      - DB_HOST=test-db
      - DB_PORT=5432
      - DB_NAME=rangkaiedu_test
      - DB_USER=testuser
      - DB_PASSWORD=testpass
      - DB_SSLMODE=disable
      - JWT_SECRET=test-jwt-secret-change-in-production
      - API_PORT=8080
      - ENV=test
    volumes:
      - ./coverage:/app/coverage
    command: ["./test-runner", "-test.v", "-test.coverprofile=/app/coverage/coverage.out"]

  # Coverage report generator
  coverage-report:
    image: alpine:latest
    depends_on:
      - test-app
    volumes:
      - ./coverage:/app/coverage
    command: >
      sh -c "
        apk add --no-cache ca-certificates &&
        cd /app/coverage &&
        go tool cover -html=coverage.out -o coverage.html &&
        echo 'Coverage report generated: coverage.html'
      "

volumes:
  test-db-data:
```

## Test Scripts

### Setup Script
```bash
#!/bin/bash
# scripts/test-setup.sh

set -e

echo "Setting up test environment..."

# Create test database
echo "Creating test database..."
createdb rangkaiedu_test 2>/dev/null || echo "Database already exists"

# Run migrations
echo "Running database migrations..."
psql -d rangkaiedu_test -f migrations/001_create_tables.sql

# Install dependencies
echo "Installing dependencies..."
go mod download

# Create test environment file
echo "Creating test environment file..."
cp internal/test/infrastructure/templates/test-config.env .env.test

# Set environment variables
export TEST_DATABASE_URL="postgres://testuser:testpass@localhost:5432/rangkaiedu_test"
export TEST_JWT_SECRET="test-jwt-secret-change-in-production"
export TEST_API_PORT="8080"
export TEST_ENV="test"

echo "Test environment setup complete!"
```

### Run Script
```bash
#!/bin/bash
# scripts/test-run.sh

set -e

# Parse command line arguments
TEST_TYPE="all"
VERBOSE=false
COVERAGE=false
PARALLEL=true

while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--type)
            TEST_TYPE="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -s|--sequential)
            PARALLEL=false
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo "Options:"
            echo "  -t, --type TYPE     Test type: all, unit, integration, e2e"
            echo "  -v, --verbose       Verbose output"
            echo "  -c, --coverage      Generate coverage report"
            echo "  -s, --sequential    Run tests sequentially"
            echo "  -h, --help          Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Build test command
TEST_CMD="go test"

if [ "$VERBOSE" = true ]; then
    TEST_CMD="$TEST_CMD -v"
fi

if [ "$COVERAGE" = true ]; then
    TEST_CMD="$TEST_CMD -cover"
fi

if [ "$PARALLEL" = true ]; then
    TEST_CMD="$TEST_CMD -parallel 4"
fi

# Add test type filter
case $TEST_TYPE in
    "unit")
        TEST_CMD="$TEST_CMD ./internal/test/unit/..."
        ;;
    "integration")
        TEST_CMD="$TEST_CMD ./internal/test/integration/..."
        ;;
    "e2e")
        TEST_CMD="$TEST_CMD ./internal/test/e2e/..."
        ;;
    "all")
        TEST_CMD="$TEST_CMD ./internal/test/..."
        ;;
    *)
        echo "Unknown test type: $TEST_TYPE"
        exit 1
        ;;
esac

# Execute tests
echo "Running tests: $TEST_CMD"
eval $TEST_CMD

# Generate coverage report if requested
if [ "$COVERAGE" = true ]; then
    echo "Generating coverage report..."
    go tool cover -html=coverage.out -o coverage.html
    echo "Coverage report generated: coverage.html"
fi
```

### Cleanup Script
```bash
#!/bin/bash
# scripts/test-cleanup.sh

set -e

echo "Cleaning up test environment..."

# Drop test database
echo "Dropping test database..."
dropdb rangkaiedu_test 2>/dev/null || echo "Database not found or already dropped"

# Remove coverage files
echo "Removing coverage files..."
rm -f coverage.out coverage.html coverage.txt

# Remove test environment file
echo "Removing test environment file..."
rm -f .env.test

# Clean up Docker containers and volumes
echo "Cleaning up Docker resources..."
docker-compose -f internal/test/infrastructure/docker-compose.test.yml down -v 2>/dev/null || echo "No Docker resources to clean up"

echo "Test environment cleanup complete!"
```

### Coverage Script
```bash
#!/bin/bash
# scripts/test-coverage.sh

set -e

echo "Generating test coverage report..."

# Run tests with coverage
go test -coverprofile=coverage.out ./internal/test/...

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Generate text coverage report
go tool cover -func=coverage.out > coverage.txt

# Calculate total coverage
TOTAL_COVERAGE=$(go tool cover -func=coverage.out | grep "total:" | awk '{print substr($3, 1, length($3)-1)}')

echo "Coverage report generated:"
echo "  HTML: coverage.html"
echo "  Text: coverage.txt"
echo "  Total: $TOTAL_COVERAGE%"

# Check if coverage meets minimum requirement
MIN_COVERAGE=80
if (( $(echo "$TOTAL_COVERAGE < $MIN_COVERAGE" | bc -l) )); then
    echo "Warning: Coverage ($TOTAL_COVERAGE%) is below minimum requirement ($MIN_COVERAGE%)"
    exit 1
else
    echo "Coverage meets minimum requirement ($MIN_COVERAGE%)"
fi
```

## CI/CD Integration

### GitHub Actions Workflow
```yaml
# ci/github-actions.yml
name: Testing

on:
  push:
    branches: [ main, develop, staging ]
  pull_request:
    branches: [ main, develop, staging ]

jobs:
  test-setup:
    runs-on: ubuntu-latest
    outputs:
      test-database-url: ${{ steps.setup.outputs.test-database-url }}
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Setup test environment
      run: |
        chmod +x scripts/test-setup.sh
        ./scripts/test-setup.sh
      id: setup
      env:
        TEST_DATABASE_URL: postgres://testuser:testpass@localhost:5432/rangkaiedu_test

  unit-tests:
    runs-on: ubuntu-latest
    needs: test-setup
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Run unit tests
      run: |
        chmod +x scripts/test-run.sh
        ./scripts/test-run.sh -t unit -v -c
      env:
        TEST_DATABASE_URL: ${{ needs.test-setup.outputs.test-database-url }}

  integration-tests:
    runs-on: ubuntu-latest
    needs: test-setup
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_USER: testuser
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
    
    - name: Run integration tests
      run: |
        chmod +x scripts/test-run.sh
        ./scripts/test-run.sh -t integration -v -c
      env:
        TEST_DATABASE_URL: postgres://testuser:testpass@localhost:5432/rangkaiedu_test

  e2e-tests:
    runs-on: ubuntu-latest
    needs: test-setup
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Run e2e tests
      run: |
        chmod +x scripts/test-run.sh
        ./scripts/test-run.sh -t e2e -v -c
      env:
        TEST_DATABASE_URL: ${{ needs.test-setup.outputs.test-database-url }}

  coverage-report:
    runs-on: ubuntu-latest
    needs: [unit-tests, integration-tests, e2e-tests]
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Generate coverage report
      run: |
        chmod +x scripts/test-coverage.sh
        ./scripts/test-coverage.sh
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
        flags: unittests
        name: codecov-umbrella
        fail_ci_if_error: true

  test-cleanup:
    runs-on: ubuntu-latest
    if: always()
    needs: [unit-tests, integration-tests, e2e-tests]
    steps:
    - uses: actions/checkout@v3
    
    - name: Cleanup test environment
      run: |
        chmod +x scripts/test-cleanup.sh
        ./scripts/test-cleanup.sh
```

### GitLab CI Configuration
```yaml
# ci/gitlab-ci.yml
stages:
  - setup
  - test
  - coverage
  - cleanup

variables:
  TEST_DATABASE_URL: "postgres://testuser:testpass@postgres:5432/rangkaiedu_test"
  GO_VERSION: "1.21"

setup-test-environment:
  stage: setup
  image: golang:$GO_VERSION
  script:
    - chmod +x scripts/test-setup.sh
    - ./scripts/test-setup.sh
  services:
    - postgres:15
  variables:
    POSTGRES_USER: testuser
    POSTGRES_PASSWORD: testpass
    POSTGRES_DB: rangkaiedu_test

unit-tests:
  stage: test
  image: golang:$GO_VERSION
  script:
    - chmod +x scripts/test-run.sh
    - ./scripts/test-run.sh -t unit -v -c
  dependencies:
    - setup-test-environment
  coverage: '/total:\s+\( statements \)\s+\d+\.\d+%/'
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage.out

integration-tests:
  stage: test
  image: golang:$GO_VERSION
  script:
    - chmod +x scripts/test-run.sh
    - ./scripts/test-run.sh -t integration -v -c
  dependencies:
    - setup-test-environment
  coverage: '/total:\s+\( statements \)\s+\d+\.\d+%/'
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage.out

e2e-tests:
  stage: test
  image: golang:$GO_VERSION
  script:
    - chmod +x scripts/test-run.sh
    - ./scripts/test-run.sh -t e2e -v -c
  dependencies:
    - setup-test-environment
  coverage: '/total:\s+\( statements \)\s+\d+\.\d+%/'
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage.out

coverage-report:
  stage: coverage
  image: golang:$GO_VERSION
  script:
    - chmod +x scripts/test-coverage.sh
    - ./scripts/test-coverage.sh
  dependencies:
    - unit-tests
    - integration-tests
    - e2e-tests
  artifacts:
    paths:
      - coverage.html
      - coverage.txt
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage.out

cleanup-test-environment:
  stage: cleanup
  image: golang:$GO_VERSION
  script:
    - chmod +x scripts/test-cleanup.sh
    - ./scripts/test-cleanup.sh
  when: always
```

## Environment Templates

### Test Configuration Template
```env
# templates/test-config.env
# Test environment configuration

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=rangkaiedu_test
DB_USER=testuser
DB_PASSWORD=testpass
DB_SSLMODE=disable

# Application Configuration
JWT_SECRET=test-jwt-secret-change-in-production
API_PORT=8080
ENV=test

# External Services
SMTP_HOST=localhost
SMTP_PORT=587
SMTP_USER=test
SMTP_PASSWORD=test
SMS_PROVIDER=test

# Test Configuration
TEST_DATABASE_URL=postgres://testuser:testpass@localhost:5432/rangkaiedu_test
TEST_JWT_SECRET=test-jwt-secret-change-in-production
TEST_API_PORT=8080
TEST_ENV=test

# Coverage Configuration
TEST_COVERAGE_MIN=80
TEST_COVERAGE_MODE=count
TEST_COVERAGE_PROFILE=coverage.out
```

### Docker Environment Template
```env
# templates/test-docker.env
# Docker-specific environment variables

# Database Configuration
DB_HOST=db
DB_PORT=5432
DB_NAME=rangkaiedu_test
DB_USER=testuser
DB_PASSWORD=testpass
DB_SSLMODE=disable

# Application Configuration
JWT_SECRET=test-jwt-secret-change-in-production
API_PORT=8080
ENV=test

# External Services
SMTP_HOST=mail
SMTP_PORT=587
SMTP_USER=test
SMTP_PASSWORD=test
SMS_PROVIDER=twilio

# Test Configuration
TEST_DATABASE_URL=postgres://testuser:testpass@db:5432/rangkaiedu_test
TEST_JWT_SECRET=test-jwt-secret-change-in-production
TEST_API_PORT=8080
TEST_ENV=test
```

## Usage Examples

### Running Tests Locally
```bash
# Setup test environment
./scripts/test-setup.sh

# Run all tests
./scripts/test-run.sh

# Run unit tests with coverage
./scripts/test-run.sh -t unit -c

# Run integration tests with verbose output
./scripts/test-run.sh -t integration -v

# Generate coverage report
./scripts/test-coverage.sh

# Cleanup test environment
./scripts/test-cleanup.sh
```

### Running Tests with Docker
```bash
# Build test image
docker build -f internal/test/infrastructure/Dockerfile.test -t rangkaiedu-test .

# Run tests in Docker container
docker run --rm -v $(pwd)/coverage:/app/coverage rangkaiedu-test

# Run tests with Docker Compose
docker-compose -f internal/test/infrastructure/docker-compose.test.yml up

# Run tests with Docker Compose in detached mode
docker-compose -f internal/test/infrastructure/docker-compose.test.yml up -d

# View test logs
docker-compose -f internal/test/infrastructure/docker-compose.test.yml logs -f test-app

# Cleanup Docker resources
docker-compose -f internal/test/infrastructure/docker-compose.test.yml down -v
```

### Running Tests in CI/CD
```bash
# GitHub Actions
git push origin feature-branch

# GitLab CI
git push origin feature-branch

# Manual CI trigger
curl -X POST -F "token=<TOKEN>" -F "ref=main" https://gitlab.com/api/v4/projects/<PROJECT_ID>/trigger/pipeline
```

## Best Practices

### Docker Configuration
1. **Use Multi-stage Builds**: Reduce image size
2. **Optimize Layer Caching**: Order Dockerfile instructions efficiently
3. **Use Specific Versions**: Pin exact versions for dependencies
4. **Minimize Image Size**: Use lightweight base images

### Test Scripts
1. **Make Scripts Executable**: Use `chmod +x` for script files
2. **Handle Errors Gracefully**: Use `set -e` for error handling
3. **Provide Clear Output**: Use descriptive messages
4. **Support Command Line Options**: Use `getopts` for argument parsing

### CI/CD Integration
1. **Use Matrix Builds**: Test across multiple environments
2. **Cache Dependencies**: Use Go module caching
3. **Parallel Execution**: Run tests in parallel when possible
4. **Generate Reports**: Create coverage and quality reports

### Environment Management
1. **Use Environment Variables**: Configure tests via environment
2. **Template Configuration**: Use templates for different environments
3. **Secure Secrets**: Use CI/CD secret management
4. **Version Configuration**: Track environment configuration versions

## Troubleshooting

### Common Issues

#### Docker Build Failures
- Check Dockerfile syntax
- Verify base image availability
- Ensure proper file permissions
- Check disk space

#### Test Execution Failures
- Verify test environment setup
- Check database connectivity
- Ensure proper dependencies
- Review test configuration

#### CI/CD Pipeline Issues
- Check workflow syntax
- Verify secret configuration
- Ensure proper permissions
- Review resource limits

### Debugging Tips
```bash
# Run tests with verbose output
./scripts/test-run.sh -v

# Run tests with race detection
go test -race ./internal/test/...

# Run tests with memory profiling
go test -memprofile=mem.out ./internal/test/...

# Check Docker container logs
docker logs <container-id>

# Check Docker network connectivity
docker network inspect <network-name>
```

## Contributing

### Adding New Infrastructure Components
1. Follow existing naming conventions
2. Add comprehensive documentation
3. Include usage examples
4. Ensure compatibility with existing infrastructure

### Updating CI/CD Configuration
1. Test changes in a feature branch
2. Verify pipeline execution
3. Update documentation
4. Communicate changes to the team

### Improving Test Scripts
1. Add error handling
2. Improve user experience
3. Add new features as needed
4. Keep scripts maintainable

## Additional Resources

### Documentation
- [Docker Documentation](https://docs.docker.com/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GitLab CI Documentation](https://docs.gitlab.com/ee/ci/)
- [Go Testing Documentation](https://golang.org/doc/testing/)

### Tools
- [Docker Compose](https://docs.docker.com/compose/)
- [Codecov](https://codecov.io/)
- [Go Test Coverage](https://pkg.go.dev/cmd/go#hdr-Testing_flags)
- [Testify](https://github.com/stretchr/testify)

### External Links
- [Docker Best Practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [CI/CD Best Practices](https://www.atlassian.com/devops/cicd)
- [Testing Infrastructure](https://martinfowler.com/bliki/TestInfrastructure.html)

---

**Note**: This README should be updated as the test infrastructure evolves and new practices are adopted.