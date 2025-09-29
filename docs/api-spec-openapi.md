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
  /auth/verify-role:
    post:
      summary: Verify user role
      description: Verifies if a user can log in with a specific role
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                role:
                  type: string
                  example: "guru"
              required:
                - role
      responses:
        '200':
          description: Role verification successful
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
        '400':
          description: Invalid request data
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '401':
          description: Role verification failed
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /auth/whatsapp-otp/send:
    post:
      summary: Send WhatsApp OTP
      description: Sends an OTP via WhatsApp for authentication
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                phone:
                  type: string
                  example: "+6281234567890"
                role:
                  type: string
                  example: "guru"
              required:
                - phone
                - role
      responses:
        '200':
          description: OTP sent successfully
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  message:
                    type: string
                    example: "OTP sent successfully"
        '400':
          description: Invalid request data
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '500':
          description: Failed to send OTP
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /auth/whatsapp-otp/verify:
    post:
      summary: Verify WhatsApp OTP
      description: Verifies the OTP sent via WhatsApp
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                phone:
                  type: string
                  example: "+6281234567890"
                otp:
                  type: string
                  example: "123456"
                role:
                  type: string
                  example: "guru"
              required:
                - phone
                - otp
                - role
      responses:
        '200':
          description: OTP verification successful
          content:
            application/json:
              schema:
                type: object
                properties:
                  token:
                    type: string
                    example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
        '400':
          description: Invalid request data
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '401':
          description: Invalid OTP
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /auth/google:
    post:
      summary: Google Authentication
      description: Authenticates user via Google
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                id_token:
                  type: string
                  example: "eyJhbGciOiJSUzI1NiIsImtpZCI6IjFlOWdkazcifQ..."
                role:
                  type: string
                  example: "guru"
              required:
                - id_token
                - role
      responses:
        '200':
          description: Authentication successful
          content:
            application/json:
              schema:
                type: object
                properties:
                  token:
                    type: string
                    example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
        '400':
          description: Invalid request data
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '401':
          description: Authentication failed
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /auth/apple:
    post:
      summary: Apple Authentication
      description: Authenticates user via Apple
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                id_token:
                  type: string
                  example: "eyJraWQiOiJmaDZCczhDIiwiYWxnIjoiUlMyNTYifQ..."
                role:
                  type: string
                  example: "guru"
              required:
                - id_token
                - role
      responses:
        '200':
          description: Authentication successful
          content:
            application/json:
              schema:
                type: object
                properties:
                  token:
                    type: string
                    example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
        '400':
          description: Invalid request data
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '401':
          description: Authentication failed
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
          example: "guru"
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
        message:
          type: string
          example: "The request data is invalid"

  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

security:
  - bearerAuth: []