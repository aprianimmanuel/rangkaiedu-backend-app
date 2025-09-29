# Rangkai Edu Backend - Docker Image

This document provides instructions for building and running the Docker image for the Rangkai Edu backend application.

## Building the Docker Image

To build the Docker image, run the following command from the root of the backend directory:

```bash
docker build -t rangkai-edu-backend:latest .
```

### Build Arguments

The Dockerfile supports several build arguments that can be used to customize the build process:

- `GO_VERSION` - Specifies the Go version to use for building (default: 1.24.7)
- `ALPINE_VERSION` - Specifies the Alpine version for the final image (default: latest)
- `APP_NAME` - Specifies the name of the compiled binary (default: backend-app)

Example with custom build arguments:

```bash
docker build \
  --build-arg GO_VERSION=1.24.7 \
  --build-arg ALPINE_VERSION=3.18 \
  --build-arg APP_NAME=rangkaiedu-backend \
  -t rangkai-edu-backend:latest .
```

## Running the Docker Container

To run the Docker container, use the following command:

```bash
docker run --rm -p 8080:8080 rangkai-edu-backend:latest
```

### Environment Variables

The application supports several environment variables for configuration:

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_NAME` - Database name (default: rangkaiedu_dev)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: password)
- `DB_SSLMODE` - Database SSL mode (default: disable)

Example with environment variables:

```bash
docker run --rm \
  -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_NAME=rangkaiedu_prod \
  -e DB_USER=postgres \
  -e DB_PASSWORD=secretpassword \
  rangkai-edu-backend:latest
```

### Using Docker Compose

For development and testing, you can use Docker Compose to run the backend service along with its dependencies:

```bash
docker-compose up
```

## Image Optimization

The Docker image is optimized for:

1. **Size**: Using a multi-stage build process to separate the build environment from the runtime environment
2. **Security**: Running the application as a non-root user
3. **Performance**: Caching Go modules to speed up builds
4. **Health Checks**: Including a health check endpoint for container orchestration

## Troubleshooting

### Common Issues

1. **Build failures related to dependencies**: Ensure that all dependencies in `go.mod` are correctly specified and accessible.

2. **Runtime connection issues**: When running the container, ensure that the database is accessible from within the container. If running the database on the host machine, you may need to use `host.docker.internal` as the database host.

### Debugging

To debug the container, you can run it in interactive mode:

```bash
docker run --rm -it --entrypoint /bin/sh rangkai-edu-backend:latest
```

This will give you a shell inside the container where you can inspect the file system and run commands.