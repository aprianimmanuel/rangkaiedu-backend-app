# Database Package

This package provides database connection functionality using PostgreSQL with connection pooling.

## Overview

The `db` package implements a singleton pattern for database connections using `pgxpool` from the `pgx` driver. It provides a single connection pool that can be used throughout the application.

## Usage

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
if err != nil {
    log.Fatal("Query failed:", err)
}
```

## Example

To run the example:

```bash
go run examples/db_example.go
```

Note: The example will fail to connect to a database because it uses dummy credentials. To run it successfully, you need to set the appropriate environment variables for a real database.

## Configuration

The database connection uses the configuration package to get connection parameters from environment variables:

- `DB_HOST`: Database host (default: "localhost")
- `DB_PORT`: Database port (default: "5432")
- `DB_NAME`: Database name (required)
- `DB_USER`: Database username (required)
- `DB_PASSWORD`: Database password (required)
- `DB_SSLMODE`: SSL mode (default: "disable")

## Connection Pool Settings

The connection pool is configured with the following settings:

- Minimum connections: 5
- Maximum connections: 20
- Maximum connection lifetime: 30 minutes
- Maximum idle time: 10 minutes
- Health check period: 1 minute
- Connection timeout: 5 seconds

## Functions

### ConnectDB() (*DB, error)

Creates a new database connection pool using the configuration from environment variables.

Returns:
- `*DB`: Database connection instance
- `error`: Error if connection fails

### Close()

Closes the database connection pool.

### GetPool() *pgxpool.Pool

Returns the underlying pgxpool.Pool for executing database queries.