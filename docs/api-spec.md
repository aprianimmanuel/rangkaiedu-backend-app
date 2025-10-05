# Rangkai Edu API Specification

## Overview
This document defines the REST API endpoints and data structures for the Rangkai Edu application. The API follows REST conventions and uses JSON for request/response formats.

## Authentication Endpoints

### Verify Role
Verifies if a user can log in with a specific role.

**Endpoint:** `POST /auth/verify-role`

**Request Body:**
```json
{
  "role": "string"
}
```

**Response:**
```json
{
  "success": true
}
```

**Response Codes:**
- `200 OK` - Role verification successful
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Role verification failed

### Send WhatsApp OTP
Sends an OTP via WhatsApp for authentication.

**Endpoint:** `POST /auth/whatsapp-otp/send`

**Request Body:**
```json
{
  "phone": "string",
  "role": "string"
}
```

**Response:**
```json
{
  "success": true,
  "message": "string"
}
```

**Response Codes:**
- `200 OK` - OTP sent successfully
- `400 Bad Request` - Invalid request data
- `500 Internal Server Error` - Failed to send OTP

### Verify WhatsApp OTP
Verifies the OTP sent via WhatsApp.

**Endpoint:** `POST /auth/whatsapp-otp/verify`

**Request Body:**
```json
{
  "phone": "string",
  "otp": "string",
  "role": "string"
}
```

**Response:**
```json
{
  "token": "string"
}
```

**Response Codes:**
- `200 OK` - OTP verification successful
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Invalid OTP

### Google Authentication
Authenticates user via Google.

**Endpoint:** `POST /auth/google`

**Request Body:**
```json
{
  "id_token": "string",
  "role": "string"
}
```

**Response:**
```json
{
  "token": "string"
}
```

**Response Codes:**
- `200 OK` - Authentication successful
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication failed

### Apple Authentication
Authenticates user via Apple.

**Endpoint:** `POST /auth/apple`

**Request Body:**
```json
{
  "id_token": "string",
  "role": "string"
}
```

**Response:**
```json
{
  "token": "string"
}
```

**Response Codes:**
- `200 OK` - Authentication successful
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication failed

### Register User
Registers a new user account with optional password support for hybrid authentication. Upon success, an initial OTP is sent to the user's email for verification.

**Endpoint:** `POST /register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "phone": "+628123456789",
  "role": "siswa",
  "password": "optionalSecurePassword"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "User registered successfully. OTP sent to email for verification."
}
```

**Response Codes:**
- `201 Created` - Registration successful
- `400 Bad Request` - Invalid request data (e.g., missing required fields, invalid email/phone)
- `409 Conflict` - User already exists

**Error Example:**
```json
{
  "error": "User with this email already exists",
  "code": 409
}
```

### Send OTP
Sends a one-time password (OTP) to the specified email or phone for authentication or verification purposes.

**Endpoint:** `POST /send-otp`

**Request Body:**
```json
{
  "identifier": "user@example.com",
  "type": "email"
}
```
*Note: `type` can be "email" or "phone". For phone, `identifier` should be in international format (e.g., "+628123456789").*

**Response (200 OK):**
```json
{
  "success": true,
  "message": "OTP sent successfully. Expires in 10 minutes."
}
```

**Response Codes:**
- `200 OK` - OTP sent
- `400 Bad Request` - Invalid identifier or type
- `404 Not Found` - User not found
- `500 Internal Server Error` - Failed to send OTP

**Error Example:**
```json
{
  "error": "Invalid email format",
  "code": 400
}
```

### Verify OTP
Verifies the provided OTP against the stored value for the identifier. If valid, returns a JWT token. OTP expires after 10 minutes.

**Endpoint:** `POST /verify-otp`

**Request Body:**
```json
{
  "identifier": "user@example.com",
  "otp": "123456"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "phone": "+628123456789",
    "role": "siswa"
  }
}
```

**Response Codes:**
- `200 OK` - OTP verified, JWT issued
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Invalid or expired OTP
- `404 Not Found` - User not found

**Error Example:**
```json
{
  "error": "Invalid OTP",
  "code": 401
}
```


### Deprecated Endpoints
The following endpoints are deprecated in favor of the OTP-based flow above. They are maintained for legacy support but should not be used for new implementations:

- `POST /auth/whatsapp-otp/send` - Legacy WhatsApp OTP send
- `POST /auth/whatsapp-otp/verify` - Legacy WhatsApp OTP verify
- `POST /login` - Password-based login (hybrid support available via optional password in /register, but recommend OTP)
- `POST /auth/google` - Google OAuth (consider migrating to OTP)
- `POST /auth/apple` - Apple OAuth (consider migrating to OTP)

For password-based hybrid login, use the deprecated `/login` endpoint with email and password, which returns the same JWT structure.

**Endpoint:** `POST /login` (Deprecated)

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securePassword"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "role": "siswa"
  }
}
```

**Note:** Passwords are hashed with bcrypt. Use OTP flow for enhanced security.

## Authorization

The Rangkai Edu API implements role-based access control (RBAC) using JWT tokens with embedded role claims. All protected endpoints require a valid JWT token in the Authorization header.

### JWT Token Structure
JWT tokens issued by the authentication endpoints contain the following claims:

```json
{
  "sub": "user-uuid",
  "email": "user@example.com",
  "phone": "+6281234567890",
  "role": "student",
  "iat": 1729000000,
  "exp": 1729086400
}
```

The `role` claim determines the user's permissions and can be one of:
- `admin` - System administrators with full access
- `teacher` - Educators with access to teaching resources
- `student` - Learners with access to course materials
- `parent` - Guardians with access to student progress information

### Middleware Protection
API endpoints are protected using custom middleware:

1. **AuthRequired**: Validates JWT tokens and extracts user information
2. **RoleRequired**: Checks if the authenticated user has a specific role
3. **RolesRequired**: Allows access to users with any of the specified roles

### Protected Routes Example
```go
// Routes requiring authentication
protected := r.Group("/api/protected")
protected.Use(AuthRequired())

// Routes requiring specific roles
admin := protected.Group("/admin")
admin.Use(RoleRequired(RoleAdmin))

// Routes allowing multiple roles
content := protected.Group("/content")
content.Use(RolesRequired(RoleAdmin, RoleTeacher))
```

### Authorization Error Responses
- `401 Unauthorized` - Missing or invalid JWT token
- `403 Forbidden` - Valid token but insufficient role permissions

## Data Structures

### User
Represents a user in the system.

```json
{
  "id": "string",
  "name": "string",
  "role": "string",
  "email": "string",
  "phone": "string"
}
```

### Token
Represents a JWT token.

```json
{
  "token": "string"
}
```

### Error Response
Standard error response format.

```json
{
  "error": "string",
  "message": "string"
}
```

## Authentication
All API endpoints (except authentication endpoints) require a valid JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

## Rate Limiting
API requests are subject to rate limiting to prevent abuse. Exceeding the rate limit will result in a `429 Too Many Requests` response.

## Error Handling
The API uses standard HTTP status codes to indicate the success or failure of requests. All error responses include a JSON object with error details.