# Rangkai Edu API - OpenAPI Specification

```yaml
openapi: 3.0.3
info:
  title: Rangkai Edu API
  description: REST API specification for the Rangkai Edu application
  version: 1.0.0
servers:
  - url: http://localhost:8080/api
    description: Development server

paths:
  /auth/register:
    post:
      summary: Register User
      description: Registers a new user account with optional password support. Upon success, an initial OTP is sent to the user's email for verification.
      operationId: registerUser
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                  example: "John Doe"
                email:
                  type: string
                  format: email
                  example: "john.doe@example.com"
                phone:
                  type: string
                  example: "+6281234567890"
                role:
                  type: string
                  enum: [admin, teacher, student, parent]
                  example: "student"
                password:
                  type: string
                  minLength: 8
                  example: "SecurePass123!"
              required:
                - name
                - email
                - phone
                - role
      responses:
        '201':
          description: Registration successful
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: "User registered successfully. Verification OTP sent to email."
                  user:
                    type: object
                    properties:
                      id:
                        type: string
                        example: ""
                      name:
                        type: string
                        example: "John Doe"
                      email:
                        type: string
                        example: "john.doe@example.com"
                      role:
                        type: string
                        example: "student"
        '400':
          description: Invalid request data
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '409':
          description: User already exists
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '500':
          description: Failed to create user or send OTP
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /auth/login:
    post:
      summary: Login User
      description: Authenticates a user with email and optional password. If no password is provided, an OTP is sent to the user's email.
      operationId: loginUser
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                email:
                  type: string
                  format: email
                  example: "john.doe@example.com"
                password:
                  type: string
                  minLength: 8
                  example: "SecurePass123!"
              required:
                - email
      responses:
        '200':
          description: Login successful or OTP sent
          content:
            application/json:
              schema:
                oneOf:
                  - type: object
                    properties:
                      message:
                        type: string
                        example: "Login successful"
                      token:
                        type: string
                        example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
                      user:
                        $ref: '#/components/schemas/User'
                  - type: object
                    properties:
                      message:
                        type: string
                        example: "OTP sent to your email"
                      otp_required:
                        type: boolean
                        example: true
        '400':
          description: Invalid request data or social auth user attempting password login
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '401':
          description: Invalid credentials
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '404':
          description: User not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '500':
          description: Failed to generate token or send OTP
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /auth/send-otp:
    post:
      summary: Send OTP
      description: Sends a one-time password (OTP) to the specified email or phone for authentication or verification purposes.
      operationId: sendOTP
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                identifier:
                  type: string
                  example: "john.doe@example.com"
                type:
                  type: string
                  enum: [email, phone]
                  example: "email"
              required:
                - identifier
                - type
      responses:
        '200':
          description: OTP sent successfully
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: "OTP sent successfully"
        '400':
          description: Invalid identifier or type
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '404':
          description: User not found for the provided identifier
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '500':
          description: Failed to generate or send OTP
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /auth/verify-otp:
    post:
      summary: Verify OTP
      description: Verifies the provided OTP against the stored value for the identifier.
      operationId: verifyOTP
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                identifier:
                  type: string
                  example: "john.doe@example.com"
                otp:
                  type: string
                  minLength: 6
                  maxLength: 6
                  example: "123456"
              required:
                - identifier
                - otp
      responses:
        '200':
          description: OTP verified successfully
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: "OTP verified successfully"
        '400':
          description: Invalid request data or invalid OTP
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '500':
          description: Failed to verify OTP
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: string
          example: "user123"
        name:
          type: string
          example: "John Doe"
        role:
          type: string
          example: "student"
        email:
          type: string
          example: "john.doe@example.com"
        phone:
          type: string
          example: "+6281234567890"

    Token:
      type: object
      properties:
        token:
          type: string
          example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

    Error:
      type: object
      properties:
        error:
          type: string
          example: "Invalid request"

  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

security:
  - bearerAuth: []