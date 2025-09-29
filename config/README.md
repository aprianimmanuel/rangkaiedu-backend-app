# Config

This directory contains all the configuration management files for the Rangkai Edu backend.

Configuration files are responsible for:
- Environment-specific settings
- Database connection parameters
- API keys and secrets
- Feature flags
- Application settings

## Database Configuration

The database configuration is managed through the `config` package which supports loading settings from both environment variables and `.env` files.

### Configuration Parameters

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_NAME` - Database name (required)
- `DB_USER` - Database username (required)
- `DB_PASSWORD` - Database password (default: empty)
- `DB_SSLMODE` - SSL mode for database connection (default: disable)

### Usage

```go
import "github.com/aprianimmanuel/backend-app/config"

// Load configuration
dbConfig, err := config.LoadConfig()
if err != nil {
    log.Fatal("Failed to load config: ", err)
}

// Build connection string for pgx driver
dsn := dbConfig.BuildDSN()
```

### Environment Files

For development, you can create a `.env` file in the project root with the following format:

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=rangkai_edu
DB_USER=your_username
DB_PASSWORD=your_password
DB_SSLMODE=disable
```

In production, it's recommended to set these as environment variables instead of using a `.env` file.