# Multi-stage Dockerfile for Go backend application

# Build arguments
ARG GO_VERSION=1.24.7
ARG ALPINE_VERSION=latest
ARG APP_NAME=backend-app

# Build stage
FROM golang:${GO_VERSION}-alpine AS builder

# Set environment variables for Git operations
ENV GITHUB_TOKEN=
ENV GOPRIVATE=
ENV GOSUMDB=off

# Install git and other dependencies needed for building
RUN apk add --no-cache git gcc musl-dev

# Configure Git to use HTTPS instead of SSH
RUN git config --global url."https://github.com/".insteadOf "git@github.com:" && \
    git config --global url."https://".insteadOf "git://"

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies (cached if go.mod and go.sum unchanged)
# Using GOPROXY to ensure we can download modules even if direct access fails
ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 to create statically linked binary
# GOOS=linux to ensure linux binary
# -a -installsuffix cgo to improve build cache
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${APP_NAME} .

# Final stage
FROM alpine:${ALPINE_VERSION}

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates wget

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/${APP_NAME} .

# Copy config files if they exist
COPY --from=builder /app/config/.env.example ./config/.env.example

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/ || exit 1

# Command to run the application
ENTRYPOINT ["./backend-app"]