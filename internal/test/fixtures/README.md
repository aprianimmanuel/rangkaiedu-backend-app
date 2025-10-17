# Test Fixtures

## Overview
This directory contains test fixtures, mock objects, and test data used across the Rangkai Edu backend test suite. Fixtures provide reusable test data, mock implementations, and database schemas to ensure consistent and reliable testing across the application.

## Directory Structure

```
internal/test/fixtures/
├── data/                # Test data definitions and fixtures
│   ├── users.go         # User test data
│   ├── classes.go       # Class test data
│   ├── materials.go     # Material test data
│   └── enrollments.go   # Enrollment test data
├── mocks/               # Mock implementations
│   ├── auth_mock.go     # Authentication mock
│   ├── database_mock.go # Database mock
│   ├── email_mock.go    # Email service mock
│   └── http_mock.go     # HTTP client mock
└── databases/           # Database test schemas
    ├── schema.sql       # Test database schema
    ├── migrations.sql   # Test migrations
    └── seeds.sql        # Test data seeds
```

## Fixtures Overview

### Test Data
Predefined data structures for common test scenarios:
- **Users**: Student, teacher, admin users with various roles
- **Classes**: Course and lesson data
- **Materials**: Educational content data
- **Enrollments**: User enrollment records

### Mock Objects
Simplified implementations for testing:
- **Authentication**: Mock authentication handlers
- **Database**: Mock database connections and queries
- **Email**: Mock email service implementations
- **HTTP**: Mock HTTP clients and servers

### Database Schemas
Test-specific database configurations:
- **Schema**: Database structure for testing
- **Migrations**: Test-specific database migrations
- **Seeds**: Initial test data for databases

## Test Data

### User Fixtures
```go
// data/users.go
package data

import (
    "database/sql"
    "time"
    
    "github.com/google/uuid"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
)

// UserFixture represents a test user structure
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

// CreateTestUser creates a test user in the database
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

// Predefined user fixtures
var (
    // StudentUser represents a typical student user
    StudentUser = UserFixture{
        ID:           uuid.New(),
        Name:         "John Student",
        Email:        "student@example.com",
        Phone:        "1234567890",
        Role:         "student",
        PasswordHash: "$2a$10$hashedpassword",
        IsMFAEnabled: false,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    
    // TeacherUser represents a teacher user
    TeacherUser = UserFixture{
        ID:           uuid.New(),
        Name:         "Jane Teacher",
        Email:        "teacher@example.com",
        Phone:        "1234567891",
        Role:         "teacher",
        PasswordHash: "$2a$10$hashedpassword",
        IsMFAEnabled: false,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    
    // AdminUser represents an admin user
    AdminUser = UserFixture{
        ID:           uuid.New(),
        Name:         "Admin User",
        Email:        "admin@example.com",
        Phone:        "1234567892",
        Role:         "admin",
        PasswordHash: "$2a$10$hashedpassword",
        IsMFAEnabled: true,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    
    // SocialUser represents a social auth user
    SocialUser = UserFixture{
        ID:           uuid.New(),
        Name:         "Social User",
        Email:        "social@example.com",
        Phone:        "1234567893",
        Role:         "student",
        PasswordHash: "", // No password for social users
        IsMFAEnabled: false,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
)
```

### Class Fixtures
```go
// data/classes.go
package data

import (
    "database/sql"
    "time"
    
    "github.com/google/uuid"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
)

// ClassFixture represents a test class structure
type ClassFixture struct {
    ID          uuid.UUID
    Name        string
    Description string
    SubjectID   uuid.UUID
    TeacherID   uuid.UUID
    MaxStudents int
    IsActive    bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// CreateTestClass creates a test class in the database
func CreateTestClass(db *sql.DB, fixture ClassFixture) *models.Class {
    class := &models.Class{
        ID:          fixture.ID,
        Name:        fixture.Name,
        Description: fixture.Description,
        SubjectID:   fixture.SubjectID,
        TeacherID:   fixture.TeacherID,
        MaxStudents: fixture.MaxStudents,
        IsActive:    fixture.IsActive,
        CreatedAt:   fixture.CreatedAt,
        UpdatedAt:   fixture.UpdatedAt,
    }
    
    if err := models.CreateClass(db, class); err != nil {
        panic(err)
    }
    
    return class
}

// Predefined class fixtures
var (
    // MathClass represents a mathematics class
    MathClass = ClassFixture{
        ID:          uuid.New(),
        Name:        "Advanced Mathematics",
        Description: "Advanced mathematics course for high school students",
        SubjectID:   uuid.New(), // Should reference a real subject
        TeacherID:   TeacherUser.ID,
        MaxStudents: 30,
        IsActive:    true,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    // ScienceClass represents a science class
    ScienceClass = ClassFixture{
        ID:          uuid.New(),
        Name:        "Physics Fundamentals",
        Description: "Introduction to physics concepts and principles",
        SubjectID:   uuid.New(), // Should reference a real subject
        TeacherID:   TeacherUser.ID,
        MaxStudents: 25,
        IsActive:    true,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
)
```

### Material Fixtures
```go
// data/materials.go
package data

import (
    "database/sql"
    "time"
    
    "github.com/google/uuid"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
)

// MaterialFixture represents a test material structure
type MaterialFixture struct {
    ID          uuid.UUID
    Title       string
    Description string
    Content     string
    ClassID     uuid.UUID
    Type        string
        FileURL     sql.NullString
    Order       int
        IsPublished bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// CreateTestMaterial creates a test material in the database
func CreateTestMaterial(db *sql.DB, fixture MaterialFixture) *models.Material {
    material := &models.Material{
        ID:          fixture.ID,
        Title:       fixture.Title,
        Description: fixture.Description,
        Content:     fixture.Content,
        ClassID:     fixture.ClassID,
        Type:        fixture.Type,
        FileURL:     fixture.FileURL,
        Order:       fixture.Order,
        IsPublished: fixture.IsPublished,
        CreatedAt:   fixture.CreatedAt,
        UpdatedAt:   fixture.UpdatedAt,
    }
    
    if err := models.CreateMaterial(db, material); err != nil {
        panic(err)
    }
    
    return material
}

// Predefined material fixtures
var (
    // MathLesson1 represents the first mathematics lesson
    MathLesson1 = MaterialFixture{
        ID:          uuid.New(),
        Title:       "Algebra Basics",
        Description: "Introduction to algebraic concepts and equations",
        Content:     "Lesson content about algebra basics...",
        ClassID:     MathClass.ID,
        Type:        "lesson",
        FileURL:     sql.NullString{String: "", Valid: false},
        Order:       1,
        IsPublished: true,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    // PhysicsLesson1 represents the first physics lesson
    PhysicsLesson1 = MaterialFixture{
        ID:          uuid.New(),
        Title:       "Newton's Laws of Motion",
        Description: "Understanding the fundamental laws of motion",
        Content:     "Lesson content about Newton's laws...",
        ClassID:     ScienceClass.ID,
        Type:        "lesson",
        FileURL:     sql.NullString{String: "", Valid: false},
        Order:       1,
        IsPublished: true,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
)
```

### Enrollment Fixtures
```go
// data/enrollments.go
package data

import (
    "database/sql"
    "time"
    
    "github.com/google/uuid"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/models"
)

// EnrollmentFixture represents a test enrollment structure
type EnrollmentFixture struct {
    ID         uuid.UUID
    UserID     uuid.UUID
    ClassID    uuid.UUID
    Status     string
    EnrolledAt time.Time
    CompletedAt sql.NullTime
}

// CreateTestEnrollment creates a test enrollment in the database
func CreateTestEnrollment(db *sql.DB, fixture EnrollmentFixture) *models.StudentEnrollment {
    enrollment := &models.StudentEnrollment{
        ID:         fixture.ID,
        UserID:     fixture.UserID,
        ClassID:    fixture.ClassID,
        Status:     fixture.Status,
        EnrolledAt: fixture.EnrolledAt,
        CompletedAt: fixture.CompletedAt,
    }
    
    if err := models.CreateStudentEnrollment(db, enrollment); err != nil {
        panic(err)
    }
    
    return enrollment
}

// Predefined enrollment fixtures
var (
    // ActiveEnrollment represents an active class enrollment
    ActiveEnrollment = EnrollmentFixture{
        ID:         uuid.New(),
        UserID:     StudentUser.ID,
        ClassID:    MathClass.ID,
        Status:     "active",
        EnrolledAt: time.Now().Add(-7 * 24 * time.Hour), // Enrolled 7 days ago
        CompletedAt: sql.NullTime{Valid: false},
    }
    
    // CompletedEnrollment represents a completed class enrollment
    CompletedEnrollment = EnrollmentFixture{
        ID:         uuid.New(),
        UserID:     StudentUser.ID,
        ClassID:    ScienceClass.ID,
        Status:     "completed",
        EnrolledAt: time.Now().Add(-30 * 24 * time.Hour), // Enrolled 30 days ago
        CompletedAt: sql.NullTime{Time: time.Now().Add(-7 * 24 * time.Hour), Valid: true}, // Completed 7 days ago
    }
)
```

## Mock Objects

### Authentication Mock
```go
// mocks/auth_mock.go
package mocks

import (
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/mock"
)

// MockAuthHandler mocks the authentication handler
type MockAuthHandler struct {
    mock.Mock
}

// Login mocks the login method
func (m *MockAuthHandler) Login(c *gin.Context) {
    m.Called(c)
}

// Register mocks the register method
func (m *MockAuthHandler) Register(c *gin.Context) {
    m.Called(c)
}

// SendOTP mocks the send OTP method
func (m *MockAuthHandler) SendOTP(c *gin.Context) {
    m.Called(c)
}

// VerifyOTP mocks the verify OTP method
func (m *MockAuthHandler) VerifyOTP(c *gin.Context) {
    m.Called(c)
}

// NewMockAuthHandler creates a new mock auth handler
func NewMockAuthHandler() *MockAuthHandler {
    return &MockAuthHandler{}
}
```

### Database Mock
```go
// mocks/database_mock.go
package mocks

import (
    "database/sql"
    "errors"
    
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/mock"
)

// MockDatabase mocks the database connection
type MockDatabase struct {
    mock.Mock
}

// Query mocks the query method
func (m *MockDatabase) Query(query string, args ...interface{}) (*sql.Rows, error) {
    argsMock := m.Called(query, args)
    return argsMock.Get(0).(*sql.Rows), argsMock.Error(1)
}

// QueryRow mocks the query row method
func (m *MockDatabase) QueryRow(query string, args ...interface{}) *sql.Row {
    argsMock := m.Called(query, args)
    return argsMock.Get(0).(*sql.Row)
}

// Exec mocks the exec method
func (m *MockDatabase) Exec(query string, args ...interface{}) (sql.Result, error) {
    argsMock := m.Called(query, args)
    return argsMock.Get(0).(sql.Result), argsMock.Error(1)
}

// NewMockDatabase creates a new mock database
func NewMockDatabase() (*sql.DB, sqlmock.Sqlmock) {
    db, mock, err := sqlmock.New()
    if err != nil {
        panic(err)
    }
    
    return db, mock
}
```

### Email Mock
```go
// mocks/email_mock.go
package mocks

import (
    "errors"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/utils/email"
    "github.com/stretchr/testify/mock"
)

// MockEmailService mocks the email service
type MockEmailService struct {
    mock.Mock
}

// SendOTPEmail mocks the send OTP email method
func (m *MockEmailService) SendOTPEmail(cfg config.Config, to, otp string) error {
    args := m.Called(cfg, to, otp)
    return args.Error(0)
}

// SendWelcomeEmail mocks the send welcome email method
func (m *MockEmailService) SendWelcomeEmail(cfg config.Config, to, name string) error {
    args := m.Called(cfg, to, name)
    return args.Error(0)
}

// NewMockEmailService creates a new mock email service
func NewMockEmailService() *MockEmailService {
    return &MockEmailService{}
}

// EmailStub provides predictable email service behavior
type EmailStub struct {
    SendError error
    SendCount int
}

// SendOTPEmail implements the email service interface
func (s *EmailStub) SendOTPEmail(cfg config.Config, to, otp string) error {
    s.SendCount++
    return s.SendError
}

// NewSuccessfulEmailStub creates a stub that always succeeds
func NewSuccessfulEmailStub() *EmailStub {
    return &EmailStub{SendError: nil}
}

// NewFailingEmailStub creates a stub that always fails
func NewFailingEmailStub(err error) *EmailStub {
    return &EmailStub{SendError: err}
}
```

### HTTP Mock
```go
// mocks/http_mock.go
package mocks

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    
    "github.com/gin-gonic/gin"
)

// MockHTTPClient mocks the HTTP client
type MockHTTPClient struct {
    DoFunc func(req *http.Request) (*http.Response, error)
}

// Do implements the HTTP client interface
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    if m.DoFunc != nil {
        return m.DoFunc(req)
    }
    
    // Default mock response
    return &http.Response{
        StatusCode: http.StatusOK,
        Body:       bytes.NewBufferString(`{"status": "ok"}`),
    }, nil
}

// NewMockHTTPClient creates a new mock HTTP client
func NewMockHTTPClient() *MockHTTPClient {
    return &MockHTTPClient{}
}

// MockHTTPServer creates a mock HTTP server for testing
func MockHTTPServer(handler http.Handler) *httptest.Server {
    return httptest.NewServer(handler)
}

// MockGinRouter creates a mock Gin router for testing
func MockGinRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    return gin.New()
}

// MakeRequest makes an HTTP request to a Gin router
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
```

## Database Schemas

### Test Schema
```sql
-- databases/schema.sql
-- Test database schema for Rangkai Edu backend

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(50) NOT NULL DEFAULT 'student',
    password_hash VARCHAR(255),
    google_id VARCHAR(255),
    facebook_id VARCHAR(255),
    is_mfa_enabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- OTPs table
CREATE TABLE otps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    identifier VARCHAR(255) NOT NULL,
    otp VARCHAR(6) NOT NULL,
    type VARCHAR(50) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Subjects table
CREATE TABLE subjects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Classes table
CREATE TABLE classes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    subject_id UUID REFERENCES subjects(id),
    teacher_id UUID REFERENCES users(id),
    max_students INTEGER DEFAULT 30,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Materials table
CREATE TABLE materials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    content TEXT,
    class_id UUID REFERENCES classes(id),
    type VARCHAR(50) NOT NULL,
    file_url VARCHAR(500),
    "order" INTEGER DEFAULT 0,
    is_published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Enrollments table
CREATE TABLE student_enrollments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    class_id UUID REFERENCES classes(id),
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    enrolled_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(user_id, class_id)
);

-- Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_otps_identifier ON otps(identifier);
CREATE INDEX idx_otps_expires_at ON otps(expires_at);
CREATE INDEX idx_classes_subject_id ON classes(subject_id);
CREATE INDEX idx_classes_teacher_id ON classes(teacher_id);
CREATE INDEX idx_materials_class_id ON materials(class_id);
CREATE INDEX idx_materials_type ON materials(type);
CREATE INDEX idx_enrollments_user_id ON student_enrollments(user_id);
CREATE INDEX idx_enrollments_class_id ON student_enrollments(class_id);
CREATE INDEX idx_enrollments_status ON student_enrollments(status);
```

### Test Migrations
```sql
-- databases/migrations.sql
-- Test-specific migrations

-- Clean up test data
TRUNCATE TABLE users, otps, classes, materials, student_enrollments RESTART IDENTITY CASCADE;

-- Insert test subjects
INSERT INTO subjects (name, code, description) VALUES
('Mathematics', 'MATH', 'Mathematics courses for various levels'),
('Science', 'SCI', 'Science courses including physics, chemistry, and biology'),
('English', 'ENG', 'English language and literature courses'),
('History', 'HIST', 'History courses for different periods and regions');

-- Insert test users
INSERT INTO users (name, email, phone, role, password_hash, is_mfa_enabled) VALUES
('John Student', 'student@example.com', '1234567890', 'student', '$2a$10$hashedpassword', false),
('Jane Teacher', 'teacher@example.com', '1234567891', 'teacher', '$2a$10$hashedpassword', false),
('Admin User', 'admin@example.com', '1234567892', 'admin', '$2a$10$hashedpassword', true);

-- Insert test classes
INSERT INTO classes (name, description, subject_id, teacher_id, max_students, is_active) VALUES
('Advanced Mathematics', 'Advanced mathematics course for high school students', 
 (SELECT id FROM subjects WHERE code = 'MATH'), 
 (SELECT id FROM users WHERE email = 'teacher@example.com'), 30, true),
('Physics Fundamentals', 'Introduction to physics concepts and principles', 
 (SELECT id FROM subjects WHERE code = 'SCI'), 
 (SELECT id FROM users WHERE email = 'teacher@example.com'), 25, true);

-- Insert test materials
INSERT INTO materials (title, description, content, class_id, type, "order", is_published) VALUES
('Algebra Basics', 'Introduction to algebraic concepts and equations', 
 'Lesson content about algebra basics...', 
 (SELECT id FROM classes WHERE name = 'Advanced Mathematics'), 'lesson', 1, true),
('Newton''s Laws of Motion', 'Understanding the fundamental laws of motion', 
 'Lesson content about Newton''s laws...', 
 (SELECT id FROM classes WHERE name = 'Physics Fundamentals'), 'lesson', 1, true);

-- Insert test enrollments
INSERT INTO student_enrollments (user_id, class_id, status, enrolled_at) VALUES
((SELECT id FROM users WHERE email = 'student@example.com'), 
 (SELECT id FROM classes WHERE name = 'Advanced Mathematics'), 'active', CURRENT_TIMESTAMP - INTERVAL '7 days'),
((SELECT id FROM users WHERE email = 'student@example.com'), 
 (SELECT id FROM classes WHERE name = 'Physics Fundamentals'), 'completed', CURRENT_TIMESTAMP - INTERVAL '30 days');
```

## Usage Examples

### Using Test Data
```go
package tests

import (
    "testing"
    
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/test/fixtures/data"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/test/fixtures/mocks"
)

func TestUserRegistration(t *testing.T) {
    // Setup mock database
    db, mock := mocks.NewMockDatabase()
    
    // Use predefined fixtures
    student := data.StudentUser
    teacher := data.TeacherUser
    
    // Create test user
    user := data.CreateTestUser(db, student)
    
    // Assert user creation
    assert.Equal(t, student.Name, user.Name)
    assert.Equal(t, student.Email, user.Email)
}
```

### Using Mock Objects
```go
package tests

import (
    "testing"
    
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/test/fixtures/mocks"
    "github.com/gin-gonic/gin"
)

func TestAuthHandler(t *testing.T) {
    // Setup mock
    mockAuth := mocks.NewMockAuthHandler()
    
    // Setup test router
    router := mocks.MockGinRouter()
    
    // Register mock handler
    router.POST("/api/auth/login", func(c *gin.Context) {
        mockAuth.Login(c)
    })
    
    // Make test request
    recorder := mocks.MakeRequest(router, "POST", "/api/auth/login", map[string]interface{}{
        "email":    "test@example.com",
        "password": "password123",
    })
    
    // Assert response
    assert.Equal(t, http.StatusOK, recorder.Code)
    
    // Verify mock was called
    mockAuth.AssertExpectations(t)
}
```

### Using Database Fixtures
```go
package tests

import (
    "database/sql"
    "testing"
    
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/test/fixtures/databases"
    "github.com/aprianimmanuel/rangkaiedu-backend/internal/test/fixtures/data"
)

func TestDatabaseIntegration(t *testing.T) {
    // Setup test database
    db := databases.SetupTestDatabase(t)
    defer databases.CleanupTestDatabase(t, db)
    
    // Create test data
    student := data.CreateTestUser(db, data.StudentUser)
    teacher := data.CreateTestUser(db, data.TeacherUser)
    class := data.CreateTestClass(db, data.MathClass)
    
    // Create enrollment
    enrollment := data.CreateTestEnrollment(db, data.ActiveEnrollment)
    
    // Test database operations
    // ... test implementation
}
```

## Best Practices

### Using Fixtures
1. **Reuse Predefined Fixtures**: Use existing fixtures when possible
2. **Customize When Needed**: Extend fixtures for specific test scenarios
3. **Keep Fixtures Simple**: Avoid complex fixture logic
4. **Document Fixtures**: Add comments explaining fixture purposes

### Using Mocks
1. **Mock External Dependencies**: Use mocks for databases, APIs, and external services
2. **Set Clear Expectations**: Define exactly what mocks should return
3. **Verify Mock Calls**: Ensure mocks are called as expected
4. **Clean Up Mocks**: Properly dispose of mock objects

### Database Testing
1. **Use Test Databases**: Never use production databases for testing
2. **Clean Up Data**: Truncate tables after each test
3. **Use Transactions**: For tests that need to rollback changes
4. **Optimize Performance**: Use proper indexing and connection pooling

## Contributing

### Adding New Fixtures
1. Follow the existing naming conventions
2. Add comprehensive documentation
3. Include usage examples
4. Ensure fixtures are reusable

### Adding New Mocks
1. Implement all required methods
2. Set clear default behaviors
3. Add helper functions for common scenarios
4. Document mock usage

### Database Schema Changes
1. Update schema files
2. Add migration scripts
3. Update test data as needed
4. Ensure backward compatibility

## Troubleshooting

### Common Issues

#### Fixture Not Found
- Check fixture names and paths
- Verify fixture imports
- Ensure fixtures are properly defined

#### Mock Configuration Problems
- Verify mock setup
- Check mock expectations
- Ensure proper mock cleanup

#### Database Connection Issues
- Check database configuration
- Verify connection strings
- Ensure proper database permissions

### Debugging Tips
```bash
# Run tests with verbose output
go test -v ./internal/test/...

# Run specific test with debugging
go test -run TestSpecificFunction -v ./internal/test/...

# Check fixture definitions
go doc github.com/aprianimmanuel/rangkaiedu-backend/internal/test/fixtures/data

# Verify mock implementations
go doc github.com/aprianimmanuel/rangkaiedu-backend/internal/test/fixtures/mocks
```

## Additional Resources

### Documentation
- [Go Testing Documentation](https://golang.org/doc/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [SQLMock Documentation](https://github.com/DATA-DOG/go-sqlmock)

### Tools
- [Go Test Coverage](https://pkg.go.dev/cmd/go#hdr-Testing_flags)
- [Testify/assert](https://github.com/stretchr/testify/assert)
- [pgAdmin](https://www.pgadmin.org/) for database management

### External Links
- [Test Data Management](https://martinfowler.com/bliki/TestData.html)
- [Mock Objects Pattern](https://martinfowler.com/bliki/MockObject.html)
- [Database Testing Strategies](https://testing.googleblog.com/2017/07/testing-database-applications.html)

---

**Note**: This README should be updated as the fixtures evolve and new testing practices are adopted.