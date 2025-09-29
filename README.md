# Rangkai Edu Backend

This is the backend service for the Rangkai Edu application, built with Go. This README provides instructions for setting up the development environment, including database configuration and migration.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Database Setup](#database-setup)
  - [Installing PostgreSQL](#installing-postgresql)
  - [Creating the Database](#creating-the-database)
  - [Configuration](#configuration)
  - [Configuration Management](#configuration-management)
- [Running Migrations](#running-migrations)
  - [Migration File Location](#migration-file-location)
  - [Executing Migrations](#executing-migrations)
- [Database Connection](#database-connection)
  - [Connection Implementation](#connection-implementation)
  - [Connection Pool](#connection-pool)
- [Development](#development)
  - [Docker Development Environment](#docker-development-environment)
- [Git Branching Strategy](#git-branching-strategy)

## Prerequisites

Before you begin, ensure you have the following installed:
- Go 1.13 or higher
- PostgreSQL 12 or higher
- Git

## Database Setup

### Installing PostgreSQL

#### Ubuntu/Debian
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

#### macOS (using Homebrew)
```bash
brew install postgresql
```

#### Windows
1. Download the installer from the [official PostgreSQL website](https://www.postgresql.org/download/windows/)
2. Run the installer and follow the setup wizard
3. Make sure to include pgAdmin and command line tools during installation

#### Starting PostgreSQL Service

##### Ubuntu/Debian
```bash
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

##### macOS
```bash
brew services start postgresql
```

##### Windows
PostgreSQL service should start automatically after installation. If not, start it from the Services application.

### Creating the Database

1. Switch to the postgres user (Linux/macOS):
   ```bash
   sudo -u postgres psql
   ```

   On Windows, use the pgAdmin tool or run:
   ```bash
   psql -U postgres
   ```

2. Create a new database:
   ```sql
   CREATE DATABASE rangkai_edu;
   ```

3. Create a new user (optional but recommended):
   ```sql
   CREATE USER rangkai_user WITH ENCRYPTED PASSWORD 'your_password';
   ```

4. Grant privileges to the user:
   ```sql
   GRANT ALL PRIVILEGES ON DATABASE rangkai_edu TO rangkai_user;
   ```

5. Enable UUID extension (required for the application):
   ```sql
   \c rangkai_edu
   CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
   ```

6. Exit PostgreSQL:
   ```sql
   \q
   ```

### Configuration

The application uses environment variables for database configuration. Create a `.env` file in the project root directory with the following content:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=rangkai_edu
DB_USER=your_username      # or rangkai_user if you created a specific user
DB_PASSWORD=your_password  # or the password you set for rangkai_user
DB_SSLMODE=disable
```

You can use the `.env.example` file in the `config` directory as a template:
```bash
cp config/.env.example .env
```

Then edit the `.env` file with your actual database credentials.

### Configuration Management

The application uses the `github.com/joho/godotenv` package to load environment variables from the `.env` file. In production, it's recommended to set these as actual environment variables instead of using a `.env` file.

The configuration is managed through the `config` package which:
- Loads settings from environment variables
- Provides default values for optional settings
- Validates required configuration parameters
- Builds PostgreSQL connection strings

Configuration parameters:
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_NAME` - Database name (required)
- `DB_USER` - Database username (required)
- `DB_PASSWORD` - Database password (default: empty)
- `DB_SSLMODE` - SSL mode for database connection (default: disable)

## Running Migrations

### Migration File Location

Database migration files are located in the `migrations/` directory. Each migration file follows a naming convention with a sequence number and descriptive name:
- `001_create_tables.sql` - Initial schema creation

### Executing Migrations

To run the migration script, you can use the `psql` command-line tool:

```bash
psql -h localhost -p 5432 -U your_username -d rangkai_edu -f migrations/001_create_tables.sql
```

Or if you're using the postgres user on Linux/macOS:
```bash
psql -h localhost -p 5432 -U postgres -d rangkai_edu -f migrations/001_create_tables.sql
```

On Windows, you might need to specify the full path:
```bash
psql -h localhost -p 5432 -U postgres -d rangkai_edu -f C:\path\to\project\migrations\001_create_tables.sql
```

The migration script will:
1. Create all necessary tables for the application
2. Set up proper relationships between tables
3. Create indexes for better query performance
4. Enable the UUID extension if not already enabled

## Database Connection

### Connection Implementation

The application uses the `pgx` driver (specifically `pgxpool`) for PostgreSQL connections. The database connection is managed through the `pkg/db` package which provides:

- A singleton pattern for database connections
- Connection pooling for efficient resource management
- Configuration loading from environment variables
- Connection validation and error handling

### Connection Pool

The database connection uses a connection pool implemented with `pgxpool` with the following settings:

- Minimum connections: 5
- Maximum connections: 20
- Maximum connection lifetime: 30 minutes
- Maximum idle time: 10 minutes
- Health check period: 1 minute
- Connection timeout: 5 seconds

Example usage:
```go
import "github.com/aprianimmanuel/backend-app/pkg/db"

// Connect to the database
database, err := db.ConnectDB()
if err != nil {
    log.Fatal("Failed to connect to database:", err)
}
defer database.Close()

// Get the connection pool
pool := database.GetPool()

// Use the pool to execute queries
var name string
err = pool.QueryRow(context.Background(), "SELECT name FROM users WHERE id = $1", 1).Scan(&name)
```

## Development

To run the application:
```bash
go run main.go
```

The server will start on port 8080 by default.

### Docker Development Environment

For easier development and consistent environments across different machines, you can use Docker and Docker Compose to run the application locally.

#### Prerequisites

Before you begin, ensure you have the following installed:
- Docker
- Docker Compose

#### Running with Docker Compose

To start the development environment with Docker Compose, run:

```bash
docker-compose up --build
```

This command will:
- Build the backend service from the Dockerfile
- Start a PostgreSQL database container
- Set up networking between services
- Mount volumes for data persistence

The application will be available at `http://localhost:8080`.

To stop the services, press `Ctrl+C` in the terminal where docker-compose is running, or run:

```bash
docker-compose down
```

## Git Branching Strategy

This project follows a Git branching strategy based on the Gitflow workflow to manage code changes and releases effectively:

### Branch Types

- **main**: The production-ready codebase. Only stable, tested code should be merged here.
- **staging**: The pre-production environment for testing and validation before releasing to production.
- **develop**: The main branch for active development. All new features and changes should be merged here first.
- **feature/***: Feature branches for developing new functionality. These branches are created from `develop` and merged back into `develop` when complete.
- **hotfix/***: Hotfix branches for urgent production fixes. These branches are created from `main` and merged back into both `main` and `develop`.

### Workflow

1. **Feature Development**:
   - Create a new feature branch from `develop`: `git checkout -b feature/feature-name develop`
   - Develop the feature and commit changes
   - Push the branch to the remote repository
   - Create a Pull Request to merge the feature branch into `develop`

2. **Release Preparation**:
   - Create a release branch from `develop`: `git checkout -b release/version-number develop`
   - Perform final testing and bug fixes
   - Merge the release branch into both `main` and `develop`
   - Tag the release on `main`

3. **Production Hotfixes**:
   - Create a hotfix branch from `main`: `git checkout -b hotfix/fix-description main`
   - Implement the fix and test thoroughly
   - Merge the hotfix branch into both `main` and `develop`
   - Tag the new release on `main`

### Branch Protection Rules

To ensure code quality and stability, the following branch protection rules are implemented:
- **main**: Requires pull request reviews, status checks, and up-to-date branches before merging
- **staging**: Requires pull request reviews and status checks before merging
- **develop**: Requires status checks before merging