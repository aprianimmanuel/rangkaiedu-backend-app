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