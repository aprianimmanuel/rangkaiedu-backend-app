# Rangkai Edu API Specification

## Overview
This document defines the REST API endpoints and data structures for the Rangkai Edu application. The API follows REST conventions and uses JSON for request/response formats.

## Implementation Status
**Status:** ✅ Partially Implemented
**Last Updated:** October 2025
**Note:** This document reflects the actual implemented API endpoints based on test results. Some endpoints from the original design have been deprecated or modified based on implementation decisions.

## Authentication Endpoints

### Register User
Registers a new user account with optional password support. Upon success, an initial OTP is sent to the user's email for verification.

**Endpoint:** `POST /api/auth/register`
**Implementation Status:** ✅ Implemented and Tested

**Request Body:**
```json
{
  "name": "string",
  "email": "user@example.com",
  "phone": "+628123456789",
  "role": "student",
  "password": "optionalSecurePassword"
}
```

**Response (201 Created):**
```json
{
  "message": "User registered successfully. Verification OTP sent to email.",
  "user": {
    "id": "",
    "name": "string",
    "email": "user@example.com",
    "role": "student"
  }
}
```

**Response Codes:**
- `201 Created` - Registration successful
- `400 Bad Request` - Invalid request data (e.g., missing required fields, invalid email/phone)
- `409 Conflict` - User already exists
- `500 Internal Server Error` - Failed to create user or send OTP

**Error Example:**
```json
{
  "error": "Email already exists"
}
```

### Login User
Authenticates a user with email and optional password. If no password is provided, an OTP is sent to the user's email.

**Endpoint:** `POST /api/auth/login`
**Implementation Status:** ✅ Implemented and Tested

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "optionalSecurePassword"
}
```

**Response with Password (200 OK):**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid-string",
    "name": "User Name",
    "email": "user@example.com",
    "phone": "+628123456789",
    "role": "student"
  }
}
```

**Response without Password (200 OK):**
```json
{
  "message": "OTP sent to your email",
  "otp_required": true
}
```

**Response Codes:**
- `200 OK` - Login successful or OTP sent
- `400 Bad Request` - Invalid request data or social auth user attempting password login
- `401 Unauthorized` - Invalid credentials
- `404 Not Found` - User not found
- `500 Internal Server Error` - Failed to generate token or send OTP

**Error Example:**
```json
{
  "error": "Invalid credentials"
}
```

### Login with WhatsApp OTP
Authenticates a user with email and optional WhatsApp OTP. If no password is provided, an OTP is sent to the user's WhatsApp number.

**Endpoint:** `POST /api/auth/whatsapp/login`
**Implementation Status:** ✅ Implemented and Tested

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "optionalSecurePassword"
}
```

**Response with Password (200 OK):**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid-string",
    "name": "User Name",
    "email": "user@example.com",
    "phone": "+628123456789",
    "role": "student"
  }
}
```

**Response without Password (200 OK):**
```json
{
  "message": "OTP sent to your WhatsApp number",
  "otp_required": true
}
```

**Response Codes:**
- `200 OK` - Login successful or OTP sent
- `400 Bad Request` - Invalid request data or social auth user attempting password login
- `401 Unauthorized` - Invalid credentials
- `404 Not Found` - User not found
- `500 Internal Server Error` - Failed to generate token or send OTP

**Error Example:**
```json
{
  "error": "Invalid credentials"
}
```

### Send OTP
Sends a one-time password (OTP) to the specified email, phone, or WhatsApp for authentication or verification purposes.

**Endpoint:** `POST /api/auth/send-otp`
**Implementation Status:** ✅ Implemented and Tested

**Request Body:**
```json
{
  "identifier": "user@example.com",
  "type": "email"
}
```
*Note: `type` can be "email", "phone", or "whatsapp". For phone/WhatsApp, `identifier` should be in international format (e.g., "+628123456789").*

**Response (200 OK):**
```json
{
  "message": "OTP sent successfully"
}
```

**Response Codes:**
- `200 OK` - OTP sent
- `400 Bad Request` - Invalid identifier or type
- `404 Not Found` - User not found for the provided identifier
- `500 Internal Server Error` - Failed to generate or send OTP

**Error Example:**
```json
{
  "error": "User not found for the provided identifier"
}
```

**WhatsApp-Specific Error Example:**
```json
{
  "error": "Invalid phone number format. Use international format (e.g., +628123456789)"
}
```

**Rate Limit Error Example:**
```json
{
  "error": "Rate limit exceeded. Please try again later."
}
```

### Verify OTP
Verifies the provided OTP against the stored value for the identifier.

**Endpoint:** `POST /api/auth/verify-otp`
**Implementation Status:** ✅ Implemented and Tested

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
  "message": "OTP verified successfully"
}
```

**Response Codes:**
- `200 OK` - OTP verified successfully
- `400 Bad Request` - Invalid request data or invalid OTP
- `500 Internal Server Error` - Failed to verify OTP

**Error Example:**
```json
{
  "error": "Invalid OTP"
}
```

## Authorization

The Rangkai Edu API implements role-based access control (RBAC) using JWT tokens with embedded role claims. All protected endpoints require a valid JWT token in the Authorization header.

### JWT Token Structure
JWT tokens issued by the authentication endpoints contain the following claims:

```json
{
  "user_id": "uuid-string",
  "email": "user@example.com",
  "role": "student",
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
  "error": "string"
}
```

## Authentication
All API endpoints (except authentication endpoints) require a valid JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

## Rate Limiting
API requests are subject to rate limiting to prevent abuse. Exceeding the rate limit will result in a `429 Too Many Requests` response.

### WhatsApp Authentication Rate Limiting
WhatsApp authentication endpoints have specific rate limiting to comply with WhatsApp Business API policies:
- Maximum 10 OTP requests per hour per phone number
- Rate limit exceeded responses include a `Retry-After` header indicating when the limit resets

## Error Handling
The API uses standard HTTP status codes to indicate the success or failure of requests. All error responses include a JSON object with error details.