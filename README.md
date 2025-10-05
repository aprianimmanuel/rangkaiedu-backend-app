# Rangkai Edu Backend

This is the backend service for the Rangkai Edu application, built with Go. This README provides comprehensive instructions for setting up the development environment, including Docker/Docker Compose usage, database configuration, migration, cloud storage setup (Alibaba Cloud OSS), and information about our Git branching strategy and CI/CD pipeline.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Setup for New Developers](#environment-setup-for-new-developers)
  - [Docker and Docker Compose Installation](#docker-and-docker-compose-installation)
- [Docker/Docker Compose Usage](#dockerdocker-compose-usage)
  - [Backend Docker Image](#backend-docker-image)
  - [Frontend Docker Image](#frontend-docker-image)
  - [Docker Compose for Local Development](#docker-compose-for-local-development)
  - [Environment Variable Configuration](#environment-variable-configuration)
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
  - [Running the Application Locally](#running-the-application-locally)
  - [Docker Development Environment](#docker-development-environment)
- [Git Branching Strategy](#git-branching-strategy)
  - [Branch Structure](#branch-structure)
  - [Workflow Diagrams](#workflow-diagrams)
  - [Workflows](#workflows)
  - [Merge Strategies and Code Review Process](#merge-strategies-and-code-review-process)
- [CI/CD Pipeline](#cicd-pipeline)
  - [GitHub Actions Workflows](#github-actions-workflows)
  - [Trigger Events](#trigger-events)
  - [Deployment Process](#deployment-process)
- [Cloud Storage](#cloud-storage)
  - [Alibaba Cloud OSS Setup](#alibaba-cloud-oss-setup)
  - [Local Storage](#local-storage)

## Prerequisites

Before you begin, ensure you have the following installed:
- Go 1.24 or higher
- PostgreSQL 15 or higher
- Git
- Docker 20.10+ (for containerized development)
- Docker Compose v2+ (for containerized development)

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
- `STORAGE_PROVIDER` - Storage provider to use (`local` or `oss`) (default: local)
- `OSS_BUCKET_NAME` - Name of the OSS bucket (required for OSS)
- `OSS_ACCESS_KEY_ID` - Access key ID for OSS access (required for OSS)
- `OSS_ACCESS_KEY_SECRET` - Access key secret for OSS access (required for OSS)
- `OSS_REGION` - Alibaba Cloud region (required for OSS)
- `OSS_ENDPOINT` - OSS endpoint URL (required for OSS)

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

## Environment Setup for New Developers

### Prerequisites

Before you begin, ensure you have the following installed:
- Go 1.24 or higher
- PostgreSQL 15 or higher
- Git
- Docker 20.10+ (for containerized development)
- Docker Compose v2+ (for containerized development)

### Step-by-Step Setup Instructions

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd rangkaiedu/backend
   ```

2. **Set up environment variables**:
   Copy the example environment file and modify it with your settings:
   ```bash
   cp config/.env.example .env
   ```
   Edit the `.env` file with your actual database credentials.

3. **Update Go dependencies** (if needed):
   ```bash
   go mod tidy
   ```

4. **Option 1: Run with Docker Compose (Recommended for new developers)**:
   ```bash
   docker-compose up --build
   ```
   This will start the backend service, PostgreSQL database, and frontend service.

5. **Option 2: Run locally without Docker**:
   - Install PostgreSQL and create the database as described in the [Database Setup](#database-setup) section
   - Run the application:
     ```bash
     go run main.go
     ```

### Common Troubleshooting Tips

1. **Port already in use**: If you encounter port conflicts, modify the ports in `docker-compose.yml` or ensure no other services are using ports 8080, 3000, or 5432.

2. **Database connection issues**:
   - Ensure the database service is running
   - Verify your environment variables in the `.env` file
   - Check that the PostgreSQL user has the correct permissions

3. **Docker build failures**:
   - Ensure Docker daemon is running
   - Check Docker has necessary permissions to access project directories
   - Verify all dependencies in `go.mod` are correctly specified

4. **Frontend not connecting to backend**:
   - Check that both services are running
   - Verify the `VITE_BACKEND_URL` environment variable is correctly set

## Docker/Docker Compose Usage

### Backend Docker Image

To build and run the backend Docker image:

1. **Building the Docker Image**:
   ```bash
   docker build -t rangkai-edu-backend:latest .
   ```

2. **Running the Docker Container**:
   ```bash
   docker run --rm -p 8080:8080 rangkai-edu-backend:latest
   ```

3. **Running with Environment Variables**:
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

### Frontend Docker Image

To build and run the frontend Docker image:

1. **Navigate to the frontend directory**:
   ```bash
   cd ../frontend-app
   ```

2. **Building the Docker Image**:
   ```bash
   docker build -t rangkaiedu-frontend:latest .
   ```

3. **Running the Docker Container**:
   ```bash
   docker run -d -p 8080:80 --name rangkaiedu-frontend rangkaiedu-frontend:latest
   ```

### Docker Compose for Local Development

For easier development and consistent environments across different machines, you can use Docker Compose to run the complete application stack.

#### Running with Docker Compose

To start the development environment with Docker Compose, run:

```bash
docker-compose up --build
```

This command will:
- Build the backend service from the Dockerfile
- Build the frontend service from the frontend Dockerfile
- Start a PostgreSQL database container
- Set up networking between services
- Mount volumes for data persistence and development

The applications will be available at:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Database: localhost:5432 (PostgreSQL)

To stop the services, press `Ctrl+C` in the terminal where docker-compose is running, or run:

```bash
docker-compose down
```

To stop services and remove volumes (WARNING: This will delete the database data):

```bash
docker-compose down -v
```

### Environment Variable Configuration

The application uses environment variables for configuration. These can be set in multiple ways:

1. **Using .env file**:
   Create a `.env` file in the project root directory based on `config/.env.example`:
   ```env
   # Database Configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_NAME=rangkai_edu
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_SSLMODE=disable
   
   # Storage Configuration
   STORAGE_PROVIDER=local
   # For Alibaba Cloud OSS:
   # OSS_BUCKET_NAME=your-bucket-name
   # OSS_ACCESS_KEY_ID=your-access-key-id
   # OSS_ACCESS_KEY_SECRET=your-access-key-secret
   # OSS_REGION=ap-southeast-1
   # OSS_ENDPOINT=https://oss-ap-southeast-1.aliyuncs.com
   ```

2. **Docker Compose environment variables**:
   The `docker-compose.yml` file uses environment variables defined in the `.env` file.

3. **Docker container environment variables**:
   When running containers directly, you can pass environment variables using the `-e` flag as shown in the examples above.

## Development

### Running the Application Locally

To run the application locally without Docker:

```bash
go run main.go
```

The server will start on port 8080 by default.

### Running Tests

To run the tests:

```bash
go test ./...
```

To run tests with coverage:

```bash
go test -cover ./...
```

To run tests with coverage and generate an HTML report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Docker Development Environment

For containerized development, see the [Docker Compose for Local Development](#docker-compose-for-local-development) section above.

## Git Branching Strategy

This project follows an enhanced Git branching strategy based on the Gitflow workflow to manage code changes and releases effectively.

### Branch Structure

1. **main**: The production-ready codebase
   - Contains code that is currently running in production
   - Only accepts merges from `staging` branch via pull requests
   - Protected with strict rules

2. **staging**: The pre-production testing environment
   - Contains code that is ready for production but undergoing final testing
   - Only accepts merges from `develop` branch via pull requests
   - Protected with moderate rules

3. **develop**: The ongoing development branch
   - Integration branch for features
   - Accepts merges from feature branches
   - Less restrictive than main and staging

4. **Supporting Branches**:
   - **feature/***: Feature branches for developing new functionality (branched from `develop`)
   - **release/***: Release branches for release preparation (branched from `develop`)
   - **hotfix/***: Hotfix branches for urgent production fixes (branched from `main`)

### Workflow Diagrams

```
main       o --------------------------------------o------------------------o
           |                                       |                        |
           |                                       |                        |
staging    | o------------------o------------------|------------------------|---o
           | |                  |                  |                        |   |
           | |                  |                  |                        |   |
develop    | | o---o---o---o----|------------------|------------------------|---|---o
           | | |   |   |   |    |                  |                        |   |   |
           | | |   |   |   |    |                  |                        |   |   |
feature    | | o---o   |   o----o                  |                        |   |   o
           | |         |                          |                        |   |
           | |         o--------------------------o                        |   |
           | |                                                            |   |
release    | o------------------------------------------------------------o   |
           |                                                                  |
           |                                                                  |
hotfix     o------------------------------------------------------------------o
```

### Workflows

#### Feature Development Workflow

1. Create a feature branch from `develop`
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/JIRA-123-short-description
   ```

2. Develop the feature and commit changes
   ```bash
   git add .
   git commit -m "feat: implement user authentication"
   ```

3. Push the feature branch to remote
   ```bash
   git push origin feature/JIRA-123-short-description
   ```

4. Create a Pull Request from feature branch to `develop`
5. After review and approval, merge the Pull Request
6. Delete the feature branch

#### Release Workflow

1. Create a release branch from `develop`
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b release/v1.2.0
   ```

2. Perform release preparations (version bump, final testing, etc.)
3. Push the release branch to remote
   ```bash
   git push origin release/v1.2.0
   ```

4. Create a Pull Request from release branch to `staging`
5. After testing on staging environment, create a Pull Request from release branch to `main`
6. After review and approval, merge both Pull Requests
7. Create a release tag on `main`
   ```bash
   git checkout main
   git pull origin main
   git tag -a v1.2.0 -m "Release version 1.2.0"
   git push origin v1.2.0
   ```

8. Merge release branch back to `develop` to incorporate any changes
9. Delete the release branch

#### Hotfix Workflow

1. Create a hotfix branch from `main`
   ```bash
   git checkout main
   git pull origin main
   git checkout -b hotfix/JIRA-456-critical-bug-fix
   ```

2. Implement the hotfix
3. Push the hotfix branch to remote
   ```bash
   git push origin hotfix/JIRA-456-critical-bug-fix
   ```

4. Create a Pull Request from hotfix branch to `main`
5. After review and approval, merge the Pull Request
6. Create a release tag on `main`
   ```bash
   git checkout main
   git pull origin main
   git tag -a v1.1.1 -m "Hotfix version 1.1.1"
   git push origin v1.1.1
   ```

7. Create Pull Requests to merge the hotfix into `staging` and `develop`
8. After review and approval, merge both Pull Requests
9. Delete the hotfix branch

### Merge Strategies and Code Review Process

#### Merge Strategies

1. **Merge Commit**:
   - Used for: Release and hotfix branches
   - Preserves complete history and chronological order
   - Creates a merge commit

2. **Squash and Merge**:
   - Used for: Feature branches
   - Creates a single commit with a summary of all changes
   - Keeps history clean and linear
   - Recommended for most feature branches

3. **Rebase and Merge**:
   - Used for: When maintaining a linear history is important
   - Reapplies commits on top of the target branch
   - Creates a linear history without merge commits

#### Code Review Process

All pull requests must be reviewed before merging:

- Feature branches to `develop`: At least 1 reviewer
- Release branches to `staging`: At least 1 reviewer
- Release branches to `main`: At least 2 reviewers
- Hotfix branches to `main`: At least 2 reviewers

Review guidelines cover:
- Code quality and adherence to project standards
- Security considerations (no hardcoded credentials, proper input validation)
- Test coverage and meaningful test cases
- Documentation updates where necessary

## CI/CD Pipeline

The project uses GitHub Actions for Continuous Integration and Continuous Deployment.

### GitHub Actions Workflows

1. **CI Pipeline** (`.github/workflows/ci.yml`):
   - Runs on push or pull request to `develop`, `staging`, or `main` branches
   - Executes backend and frontend builds, tests, and security scans
   - Builds and scans Docker images for vulnerabilities

2. **CD Pipeline** (`.github/workflows/cd.yml`):
   - Triggers after successful CI workflow completion
   - Deploys to staging environment on successful `develop` branch builds
   - Deploys to production environment on successful `main` branch builds

3. **Security Scanning** (`.github/workflows/security.yml`):
   - Runs daily at 2 AM UTC or manually via workflow dispatch
   - Performs security scans on both backend and frontend code
   - Checks for dependency vulnerabilities

### Trigger Events

- **CI Pipeline**:
  - Push to `develop`, `staging`, or `main` branches
  - Pull requests to `develop`, `staging`, or `main` branches

- **CD Pipeline**:
  - Successful completion of CI workflow

- **Security Scanning**:
  - Scheduled daily at 2 AM UTC
  - Manual trigger via workflow dispatch

### Deployment Process

#### Staging Deployment
- Triggered when CI workflow completes successfully on `develop` branch
- Deploys to staging environment at https://staging.rangkaiedu.com
- Environment variables are configured for staging

#### Production Deployment
- Triggered when CI workflow completes successfully on `main` branch
- Deploys to production environment at https://rangkaiedu.com
- Environment variables are configured for production

The deployment process includes:
1. Deploying Docker images to the target infrastructure
2. Running database migrations if needed
3. Performing health checks to ensure services are running properly
4. Notifying the team of deployment status

## Cloud Storage

The application supports multiple storage providers for teaching materials:
- Local file system storage (default for development)
- Alibaba Cloud OSS (recommended for production)

### Alibaba Cloud OSS Setup

To use Alibaba Cloud OSS for storing teaching materials:

1. **Set up infrastructure using Terraform**:
   Follow the instructions in `terraform/README.md` to provision an OSS bucket and RAM user.

2. **Configure environment variables**:
   Update your `.env` file with the OSS configuration:
   ```env
   STORAGE_PROVIDER=oss
   OSS_BUCKET_NAME=your-bucket-name
   OSS_ACCESS_KEY_ID=your-access-key-id
   OSS_ACCESS_KEY_SECRET=your-access-key-secret
   OSS_REGION=ap-southeast-1
   OSS_ENDPOINT=https://oss-ap-southeast-1.aliyuncs.com
   ```

3. **Run the application**:
   Start the application using Docker Compose or run it locally.

### Local Storage

For development purposes, you can use local file system storage:

1. **Configure environment variables**:
   ```env
   STORAGE_PROVIDER=local
   ```

2. **Run the application**:
   Files will be stored in the `./uploads` directory.