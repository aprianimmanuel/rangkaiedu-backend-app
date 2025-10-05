#!/bin/sh

set -e  # Exit on error

# Create the user if not exists (this can stay in a DO block)
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  DO \$\$
  BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'rangkaiedudev1') THEN
      CREATE USER rangkaiedudev1 WITH PASSWORD '12d1q23wxm19wkc1fsdcq23';
    END IF;
  END \$\$; 
EOSQL

# Function to create database if not exists and grant privileges
create_db_if_not_exists() {
  DB="$1"
  
  # Check if database exists
  exists=$(psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -t <<-EOSQL
    SELECT 1 FROM pg_database WHERE datname = '$DB';
EOSQL
  )
  
  if [ -z "$exists" ]; then
    echo "Creating database $DB..."
    createdb --username="$POSTGRES_USER" "$DB"
  else
    echo "Database $DB already exists."
  fi
  
  # Grant privileges (runs regardless, safe if DB exists)
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    GRANT ALL PRIVILEGES ON DATABASE "$DB" TO "$POSTGRES_USER";
    GRANT ALL PRIVILEGES ON DATABASE "$DB" TO rangkaiedudev1;
EOSQL
}

# Create main and test databases
create_db_if_not_exists "rangkaiedu"
create_db_if_not_exists "rangkaiedu_test"

# After create_db_if_not_exists calls, add:
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "rangkaiedu" -f /docker-entrypoint-initdb.d/create_tables.sql
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "rangkaiedu" -f /docker-entrypoint-initdb.d/seed_data.sql

# Repeat for test DB if needed
psql -v ON_ERROR_STOP=1 --username "rangkaiedudev1" --dbname "rangkaiedu_test" -f /docker-entrypoint-initdb.d/create_tables_test.sql
psql -v ON_ERROR_STOP=1 --username "rangkaiedudev1" --dbname "rangkaiedu_test" -f /docker-entrypoint-initdb.d/seed_data_test.sql

