# Rangkai Edu - User Management & Authentication Module
**Comprehensive Project Plan Summary**

## Overview
This document provides a comprehensive overview of the User Management & Authentication module implementation for the Rangkai Edu platform, covering tasks T1.4 (Backend Authentication), T1.5 (Frontend Implementation), and T1.6 (Backend Authorization). The implementation follows a logical sequence: building the backend API first, then the frontend UI that consumes it.

## Detailed Task Breakdown

### T1.4: Backend API for Authentication
**Objective**: Create the core API endpoints for user authentication with optional password handling, OTP for email/phone login, and JWT token generation.

#### Subtasks:
1. **T1.4.1: Research and Implement Password Hashing using bcrypt (Optional)**
   - **Description**: Research and implement optional password hashing using `bcrypt` in Go for hybrid security.
   - **Assigned Mode**:
     - First: Project Research (to provide a brief on why bcrypt is a strong choice)
     - Then: Code (Go specialist) to create a Go helper package for hashing and verifying passwords.
   - **Acceptance Criteria**: [x] Password hashing utility package created with proper functions, bcrypt library properly integrated, helper functions for hashing and verification implemented, security documentation created.

2. **T1.4.2: Build User Registration Handler with Initial OTP**
   - **Description**: Build the user registration handler (`/register`) that accepts new user data (email, whatsapp_number/mobile_number), optionally hashes the password using bcrypt, saves to database, and sends initial OTP to email.
   - **Assigned Mode**: Code (Go specialist)
   - **Acceptance Criteria**: [x] `/register` endpoint accepts POST requests and creates new users, optional passwords hashed before storage, request validation prevents invalid data, duplicate registration handled, initial OTP sent.

3. **T1.4.3 & T1.4.4: Build OTP-Based Authentication Handlers**
   - **Description**: Build handlers for `/send-otp` and `/verify-otp` that generate/send OTP via email/SMS, verify against database, and generate JWT upon success. Deprecate `/login` for password.
   - **Assigned Mode**: Code (Go specialist)
   - **Acceptance Criteria**: [x] `/send-otp` and `/verify-otp` endpoints work for email/phone, valid OTP generates JWT with role, invalid/expired return errors, JWT includes role.

### T1.5: Frontend Implementation & Integration
**Objective**: Develop the frontend UI components and integrate them with the backend authentication API.

#### Subtasks:
1. **T1.5.1: LoginPage UI Development**
   - **Description**: Develop the UI for the `LoginPage.jsx`, including the form for login and integration with the `/login` API endpoint.
   - **Assigned Mode**: Code (React specialist)
   - **Acceptance Criteria**: LoginPage component renders correctly with styled form, form validation prevents submission of invalid data, successful integration with backend `/login` endpoint, proper error handling and user feedback, responsive design works on mobile and desktop.

2. **T1.5.2: Client-Side Authentication State Management**
   - **Description**: Implement client-side state management for authentication, storing the received JWT securely and creating a global context to share the user's authentication status.
   - **Assigned Mode**: Code (React specialist)
   - **Acceptance Criteria**: AuthContext provides authentication state to all components, JWT tokens stored securely with appropriate storage mechanism, login/logout functions work correctly, token validation and expiration handled properly, authentication state persists across page refreshes.

3. **T1.5.3: ProtectedRoute Component**
   - **Description**: Create a `ProtectedRoute.jsx` component that checks if a user is logged in and has the correct role before rendering a protected page, redirecting unauthenticated users to the login page.
   - **Assigned Mode**: Code (React specialist)
   - **Acceptance Criteria**: ProtectedRoute component prevents access to unauthenticated users, role-based authorization works correctly for different user types, unauthenticated users redirected to login page, protected routes render correctly for authorized users, proper error handling for authorization failures.

### T1.6: Backend API for Authorization
**Objective**: Implement role-based access control and enhance JWT tokens with role claims for secure route protection.

#### Subtasks:
1. **T1.6.1: JWT with Role Claim**
   - **Description**: The JWT generated in T1.4.4 must include the user's `role` (e.g., 'guru', 'siswa', 'admin') as a claim in its payload.
   - **Assigned Mode**: Code (Go specialist)
   - **Acceptance Criteria**: JWT tokens include user role in payload, role information is properly encoded and secured, token verification correctly extracts role information, tokens generated for all user roles, documentation of JWT structure with role claims.

2. **T1.6.2 & T1.6.3: Authentication Middleware**
   - **Description**: Create a custom authentication middleware in Gin that inspects the JWT from incoming requests, validates it, and extracts the user's role to grant or deny access to specific API routes.
   - **Assigned Mode**: Code (Go specialist)
   - **Acceptance Criteria**: Authentication middleware validates JWT tokens correctly, user role extracted from tokens for authorization, middleware protects routes based on role requirements, proper error responses for invalid/missing tokens, comprehensive test coverage for middleware functions, documentation for middleware usage.

## Dependencies
- **Phase 2 (T1.6) depends on Phase 1 (T1.4)**: Authorization middleware and role-based JWT claims require the authentication endpoints to be functional first.
- **Phase 3 (T1.5) depends on Phase 1 & 2 (T1.4 & T1.6)**: The frontend implementation requires functional backend authentication and authorization APIs to integrate with.
- **All tasks depend on T1.2 Database Implementation**: Authentication and authorization require database connectivity and user table structures.

## Acceptance Criteria by Phase

### Phase 1 & 2 Completion (T1.4 & T1.6):
"Backend provides fully functional `/register`, `/send-otp`, and `/verify-otp` endpoints. The `/verify-otp` endpoint returns a JWT containing a `role` claim. A middleware exists that can protect routes based on this role."

### Phase 3 Completion (T1.5):
"Users can log in via the UI. Pages wrapped with `ProtectedRoute` are inaccessible to unauthenticated users and redirect them to the login page."

## Follow-up Documentation Tasks

### API Documentation Task:
**Description**: Create API documentation for the authentication endpoints, updated for OTP, including request/response examples.
**Assigned Mode**: Documentation Writer
**Deliverables**:
1. Document `/register` endpoint with request/response examples (optional password, initial OTP)
2. Document `/send-otp` and `/verify-otp` endpoints with request/response examples
3. Document deprecated `/login` endpoint
4. Explain JWT token structure and usage
5. Include error response formats and codes
6. Document how JWT tokens include role claims
7. Document middleware usage for protecting routes
8. Explain role-based access control implementation
9. Include examples of protected endpoints and required roles

### Frontend Usage Documentation Task:
**Description**: Explain how the `ProtectedRoute` component should be used in the frontend `README`.
**Assigned Mode**: Documentation Writer
**Deliverables**:
1. How to use the AuthContext for authentication state management
2. How to implement ProtectedRoute for protecting components
3. Best practices for handling JWT tokens in the frontend
4. Example usage of authentication components with code snippets

## Implementation Sequence
1. Begin with T1.4.1 (Research and Implementation of bcrypt)
2. Proceed with T1.4.2 and T1.4.3/T1.4.4 (Registration and Login handlers)
3. Continue with T1.6.1 and T1.6.2/T1.6.3 (JWT role claims and middleware)
4. Finally implement T1.5.1, T1.5.2, and T1.5.3 (Frontend components)

This comprehensive plan ensures a secure, well-structured authentication and authorization system with clear separation of concerns between backend and frontend implementations, following the established patterns of the Rangkai Edu project.

## OTP Flow and Examples

### Authentication Flow Overview
The module implements a primary OTP-based authentication flow for enhanced security, with optional hybrid password support for legacy compatibility. The `/login` endpoint is deprecated in favor of OTP; new users should use the OTP flow exclusively.

**OTP Flow Steps:**
1. **Registration**: POST to `/register` creates a user (optional password hashed with bcrypt) and sends an initial OTP to email.
2. **OTP Request**: POST to `/send-otp` with email/phone to generate and send a 6-digit OTP (expires in 10 minutes).
3. **Verification**: POST to `/verify-otp` with identifier and OTP; success returns JWT with user claims including role.
4. **Protected Access**: Use JWT in `Authorization: Bearer <token>` header; middleware validates token and checks role.

**Hybrid Password Support**: During registration, an optional password can be provided (hashed with bcrypt). Legacy `/login` uses email/password but is deprecated—migrate to OTP for better UX and security (no password management needed).

**Key Security Features:**
- OTPs stored hashed in database, expire 10min.
- JWT signed with HS256, includes role for authorization.
- Rate limiting on OTP sends to prevent abuse.
- Error responses standardized: `{"error": "description", "code": HTTP_status}`.

### API Endpoint Examples

#### 1. User Registration (`/register`)
**Purpose**: Create new user, optional password, initial OTP sent.

**Request Example:**
```json
{
  "email": "siswa@example.com",
  "phone": "+6281234567890",
  "role": "siswa",
  "password": "MySecurePass123"  // Optional; hashed if provided
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "User registered. Please check your email for OTP."
}
```

**Error Example (409 Conflict):**
```json
{
  "error": "User with email already exists",
  "code": 409
}
```

#### 2. Send OTP (`/send-otp`)
**Purpose**: Request OTP for login/verification.

**Request Example (Email):**
```json
{
  "identifier": "siswa@example.com",
  "type": "email"
}
```

**Request Example (Phone):**
```json
{
  "identifier": "+6281234567890",
  "type": "phone"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "OTP sent successfully. Valid for 10 minutes."
}
```

**Error Example (400):**
```json
{
  "error": "Invalid phone format",
  "code": 400
}
```

#### 3. Verify OTP (`/verify-otp`)
**Purpose**: Validate OTP, issue JWT.

**Request Example:**
```json
{
  "identifier": "siswa@example.com",
  "otp": "123456"
}
```

**Success Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwiZW1haWwiOiJzaXN3YUBleGFtcGxlLmNvbSIsInBob25lIjoiKzYyODEyMzQ1Njc4OTAiLCJyb2xlIjoic2lzd2EiLCJpYXQiOjE3MjkwMDAwMDAsImV4cCI6MTcyOTA4NjQwMH0.signature",
  "user": {
    "id": 1,
    "email": "siswa@example.com",
    "phone": "+6281234567890",
    "role": "siswa"
  }
}
```

**Error Example (401):**
```json
{
  "error": "Invalid or expired OTP",
  "code": 401
}
```

#### 4. Deprecated: Password Login (`/login`)
**Purpose**: Legacy password-based auth (use OTP instead).

**Request Example:**
```json
{
  "email": "siswa@example.com",
  "password": "MySecurePass123"
}
```

**Success Response (200):** Same as `/verify-otp`.

**Note**: Only works if password was provided during registration. Deprecated; OTP preferred for passwordless security.

### JWT Structure and Usage
JWT payload example (decoded):
```json
{
  "sub": "1",
  "email": "siswa@example.com",
  "phone": "+6281234567890",
  "role": "siswa",
  "iat": 1729000000,
  "exp": 1729086400
}
```
- Signed with secret key (HS256).
- Use in headers for protected routes.
- Middleware extracts `role` for access control (e.g., admin-only endpoints).

### Error Handling
All errors return JSON with `error` (description) and `code` (HTTP status). Common codes: 400 (bad request), 401 (unauth), 409 (conflict), 500 (server error).

This integration completes the documentation for the authentication module, addressing the follow-up tasks.