# Docker Compose Configuration for Rangkai Edu

This document explains how to use the Docker Compose setup for local development of the Rangkai Edu application.

## Overview

The Docker Compose configuration includes three services:
1. **PostgreSQL Database** - For data persistence
2. **Backend Service** - Go application API
3. **Frontend Service** - React/Vite application

## Prerequisites

- Docker Engine 20.10+ installed
- Docker Compose v2+ installed
- At least 4GB of available RAM

## Quick Start

1. Clone the repository (if not already done):
   ```bash
   git clone <repository-url>
   cd rangkaiedu
   ```

2. Navigate to the backend directory:
   ```bash
   cd backend
   ```

3. Start all services:
   ```bash
   docker-compose up
   ```

4. Access the applications:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Database: localhost:5432 (PostgreSQL)

## Configuration

### Environment Variables

The `.env` file contains the basic configuration for the services. You can modify these values as needed:

```env
# Database Configuration
DB_NAME=rangkai_edu
DB_USER=postgres
DB_PASSWORD=postgres123

# Backend Configuration
BACKEND_PORT=8080

# Frontend Configuration
FRONTEND_PORT=3000
```

### Service Details

#### PostgreSQL Database
- Image: `postgres:15-alpine`
- Port: 5432 (mapped to host)
- Data is persisted in a Docker volume named `postgres_data`
- Database initialization scripts from the `migrations/` directory are automatically executed

#### Backend Service
- Built from the Dockerfile in the current directory
- Port: 8080 (mapped to host)
- Depends on the database service
- Environment variables are passed to configure database connection

#### Frontend Service
- Built from the Dockerfile in the `../frontend-app` directory
- Port: 3000 (mapped to host)
- Environment variables configure the backend API URL

## Common Commands

### Start Services

Start all services in the foreground:
```bash
docker-compose up
```

Start all services in the background:
```bash
docker-compose up -d
```

### Stop Services

Stop all services:
```bash
docker-compose down
```

Stop all services and remove volumes (WARNING: This will delete the database data):
```bash
docker-compose down -v
```

### View Logs

View logs for all services:
```bash
docker-compose logs
```

View logs for a specific service:
```bash
docker-compose logs backend
```

### Execute Commands in Containers

Access the database container:
```bash
docker-compose exec db psql -U postgres
```

Access the backend container:
```bash
docker-compose exec backend sh
```

Access the frontend container:
```bash
docker-compose exec frontend sh
```

## Development Workflow

### Backend Development

For backend development with hot reloading:
1. Modify the docker-compose.yml to use volume mounting for Go source files
2. Use a development-specific Dockerfile that supports hot reloading

### Frontend Development

For frontend development with hot reloading:
1. The docker-compose.yml already mounts the source code directory
2. Changes to the frontend code will be reflected immediately

## Troubleshooting

### Common Issues

1. **Port already in use**
   - Solution: Stop services using those ports or change the port mappings in docker-compose.yml

2. **Database connection failed**
   - Solution: Check that the database service is running and the environment variables are correct

3. **Permission denied errors**
   - Solution: Ensure Docker has the necessary permissions to access the project directories

### Useful Debugging Commands

Check service status:
```bash
docker-compose ps
```

View resource usage:
```bash
docker stats
```

Rebuild services:
```bash
docker-compose build
```

## Security Notes

- The database password is stored in the `.env` file which should not be committed to version control
- In production, use Docker secrets or a more secure method to manage credentials
- The PostgreSQL container uses `POSTGRES_HOST_AUTH_METHOD: trust` for development convenience, which should not be used in production

## Customization

You can customize the Docker Compose setup by modifying:
- `docker-compose.yml` - Service definitions
- `.env` - Environment variables
- `migrations/` - Database initialization scripts

For different environments (staging, production), create separate docker-compose.override.yml files.