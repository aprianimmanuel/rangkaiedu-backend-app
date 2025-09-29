# Integration Testing Procedures

## Overview
This document outlines the procedures for conducting integration testing between the frontend and backend components of the Rangkai Edu application. Integration testing ensures that the different modules of the application work together as expected.

## Test Environment Setup

### Development Environment
- **Backend**: Go application running on `localhost:8080`
- **Frontend**: React application running on `localhost:3000`
- **Database**: Development database with test data
- **Proxy**: Vite proxy configured to forward API requests from `/api` to `http://localhost:8080`

### Staging Environment
- **Backend**: Go application deployed to staging server
- **Frontend**: React application deployed to staging server
- **Database**: Staging database with test data
- **Proxy**: Nginx proxy configured to route requests

### Test Data
Test data should include:
- Sample users for each role (admin, teacher, student, parent)
- Sample classes and subjects
- Sample grades and assessments
- Sample attendance records

## End-to-End Integration Tests

### Authentication Flow Tests

#### 1. Role Verification Test
**Description**: Verify that users can select their role and the system validates it correctly.

**Steps**:
1. Navigate to login page
2. Select a role (admin, teacher, student, parent)
3. Verify that the role is accepted by the backend
4. Proceed to authentication method selection

**Expected Results**:
- Role verification endpoint returns success
- Correct authentication options are displayed based on role

#### 2. WhatsApp OTP Flow Test
**Description**: Test the complete WhatsApp OTP authentication flow.

**Steps**:
1. Enter a valid phone number
2. Request OTP via WhatsApp
3. Receive OTP (mocked in test environment)
4. Enter OTP
5. Verify authentication success

**Expected Results**:
- OTP is sent successfully
- OTP verification returns valid JWT token
- User is redirected to appropriate dashboard

#### 3. Google Authentication Test
**Description**: Test Google authentication flow.

**Steps**:
1. Click Google login button
2. Complete Google authentication (mocked in test environment)
3. Verify token is received
4. Verify user is logged in

**Expected Results**:
- Google authentication returns valid JWT token
- User is redirected to appropriate dashboard

#### 4. Apple Authentication Test
**Description**: Test Apple authentication flow.

**Steps**:
1. Click Apple login button
2. Complete Apple authentication (mocked in test environment)
3. Verify token is received
4. Verify user is logged in

**Expected Results**:
- Apple authentication returns valid JWT token
- User is redirected to appropriate dashboard

### API Integration Tests

#### 1. User Data Retrieval Test
**Description**: Verify that user data can be retrieved after authentication.

**Steps**:
1. Authenticate as a user
2. Make API request to retrieve user data
3. Verify response contains expected user information

**Expected Results**:
- API returns user data in correct format
- User data matches authenticated user

#### 2. Dashboard Data Test
**Description**: Verify that dashboard data is loaded correctly for each user role.

**Steps**:
1. Authenticate as user of each role
2. Navigate to dashboard
3. Verify that appropriate data is loaded

**Expected Results**:
- Admin dashboard shows school data
- Teacher dashboard shows class data
- Student dashboard shows grade data
- Parent dashboard shows child data

### Cross-Component Integration Tests

#### 1. Context Integration Test
**Description**: Verify that the AuthContext correctly manages user state across components.

**Steps**:
1. Authenticate user
2. Navigate between different pages
3. Verify that user state is maintained
4. Log out and verify state is cleared

**Expected Results**:
- User state is consistent across all components
- Authentication status is correctly updated
- Protected routes are properly handled

#### 2. API Integration Test
**Description**: Verify that the API utility correctly handles requests and responses.

**Steps**:
1. Make various API requests
2. Verify that requests include correct headers
3. Verify that responses are properly handled
4. Test error scenarios

**Expected Results**:
- Requests include authentication headers
- Responses are parsed correctly
- Errors are handled gracefully

## Testing Procedures

### Manual Testing Procedure

#### Pre-requisites
1. Backend server running
2. Frontend development server running
3. Test user accounts created
4. Test data populated

#### Test Execution
1. **Environment Setup**
   - Start backend server: `go run main.go`
   - Start frontend server: `npm run dev`
   - Verify both servers are running

2. **Authentication Tests**
   - Test each authentication method
   - Verify role-based access control
   - Test session management

3. **Data Flow Tests**
   - Test data retrieval from API
   - Test data submission to API
   - Verify data consistency

4. **Error Handling Tests**
   - Test network errors
   - Test authentication errors
   - Test validation errors

#### Test Documentation
- Record test results in test execution log
- Document any issues found
- Include screenshots for UI-related issues

### Automated Testing Procedure

#### Test Suite Execution
1. **Unit Tests**
   - Run backend unit tests: `go test ./...`
   - Run frontend unit tests: `npm run test`

2. **Integration Tests**
   - Run backend integration tests
   - Run frontend component tests

3. **End-to-End Tests**
   - Run Cypress tests: `npm run cypress:run`

#### Continuous Integration
- Tests run automatically on each commit
- Results reported to CI/CD system
- Deployment blocked if tests fail

## Test Data Management

### Test Database
- Separate database for testing
- Pre-populated with test data
- Reset between test runs

### Test Users
- Admin user with full permissions
- Teacher user with class management permissions
- Student user with learning access
- Parent user with child data access

### Mock Services
- Mock SMS service for OTP testing
- Mock OAuth services for social login testing
- Mock email service for notification testing

## Troubleshooting Guide

### Common Issues

#### 1. Authentication Failures
**Symptoms**: Unable to log in, token errors
**Solutions**:
- Verify backend server is running
- Check JWT secret configuration
- Verify user credentials

#### 2. API Connection Errors
**Symptoms**: 404, 500 errors, timeout errors
**Solutions**:
- Verify API endpoints are correct
- Check network connectivity
- Verify proxy configuration

#### 3. Data Consistency Issues
**Symptoms**: Missing data, incorrect data
**Solutions**:
- Verify database connections
- Check data validation rules
- Review API response formats

### Debugging Tools

#### Backend Debugging
- Use logging to trace request flow
- Use debugging tools like Delve
- Check database queries

#### Frontend Debugging
- Use browser developer tools
- Check network tab for API requests
- Use React DevTools for component state

## Performance Testing

### Response Time Testing
- Measure API response times
- Measure page load times
- Identify performance bottlenecks

### Load Testing
- Test with multiple concurrent users
- Verify system stability under load
- Monitor resource usage

## Security Testing

### Authentication Security
- Test token expiration
- Verify password requirements
- Test session management

### Data Security
- Verify data encryption
- Test input validation
- Check for injection vulnerabilities

## Reporting

### Test Results
- Document all test results
- Include pass/fail status
- Record execution time

### Issue Tracking
- Log all issues found
- Include steps to reproduce
- Assign priority levels

### Test Coverage
- Measure code coverage
- Identify untested areas
- Plan additional tests

## Next Steps

1. Implement end-to-end integration tests
2. Set up continuous integration for integration tests
3. Create test data management procedures
4. Document troubleshooting procedures
5. Establish performance benchmarks