# Rangkai Edu Database Migration Guide

## Overview
This guide documents the migration strategy for the Rangkai Edu PostgreSQL database using the Goose migration tool. It covers setup, execution, best practices, and rollback procedures. The initial migrations (001_create_tables.sql and 002_seed_data.sql) implement Phase 1 schema design, including 12 tables (10 core + 2 junctions), relationships, constraints, indexes, and sample seed data for development/testing.

## Prerequisites
- PostgreSQL 14+ installed and running (e.g., local dev DB named `rangkaiedu`).
- Go 1.21+ (updated from 1.19 for Goose v3 compatibility).
- Project dependencies: Run `go mod tidy` after setup.
- Environment variables: Set `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD` in `.env` for connection strings.

## Installing Goose
Goose is the migration tool for managing schema changes. Install the CLI:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

**Note**: If installation fails (e.g., "cannot load cmp" error due to toolchain issues), try:
- `go clean -modcache`
- Ensure Go 1.21+ with `go version`
- Update Go if needed: Download from golang.org
- Alternative: Use v2 if v3 issues persist: `go install github.com/pressly/goose/cmd/goose@latest`

Add to `go.mod` if integrating programmatically:
```
require github.com/pressly/goose/v3 v3.25.0
```

## Migration Directory Structure
Migrations are stored in `migrations/` as SQL files named `XXX_description.sql` (e.g., 001_create_tables.sql).
- **Up Section**: `-- +goose Up` followed by forward changes (CREATE TABLE, INSERT, etc.).
- **Down Section**: `-- +goose Down` followed by reverse changes (DROP TABLE, DELETE, etc.).
- Files run in numerical order on `goose up`.

Current migrations:
- `001_create_tables.sql`: Creates extensions (pgcrypto for UUIDs), tables with FKs/constraints/triggers, indexes (individual/composite per Phase 1).
- `002_seed_data.sql`: Inserts sample data (1 admin, 2 teachers, 3 students, 2 parents, 2 schools, 4 classes, 5 subjects; populates junctions).

## Running Migrations
Use Goose CLI with a PostgreSQL connection string. Example for local dev:

```bash
# Apply all pending migrations (up)
goose -dir migrations postgres "host=localhost port=5432 dbname=rangkaiedu user=postgres password=yourpass sslmode=disable" up

# Rollback last migration (down)
goose -dir migrations postgres "host=localhost port=5432 dbname=rangkaiedu user=postgres password=yourpass sslmode=disable" down

# Rollback to specific version (e.g., before 002)
goose -dir migrations postgres "host=localhost port=5432 dbname=rangkaiedu user=postgres password=yourpass sslmode=disable" down-to 001

# Status: Show applied/pending migrations
goose -dir migrations postgres "host=localhost port=5432 dbname=rangkaiedu user=postgres password=yourpass sslmode=disable" status

# Create new migration skeleton
goose -dir migrations create new_feature.sql
```

For production, use environment-specific connection strings (e.g., via config). Always backup DB before running.

## Up/Down Procedures
- **Up**: Applies forward changes sequentially. Goose tracks version in `goose_db_version` table (auto-created).
- **Down**: Reverses last migration or to a version. Ensure down sections are complete (e.g., DROP before DELETE to avoid FK errors).
- Compatibility: Scripts use PostgreSQL-specific syntax (e.g., gen_random_uuid() from pgcrypto). Test on matching version.

## Guidelines for Future Changes
### Naming Conventions
- Prefix: Sequential number (001, 002, ...).
- Description: Lowercase with underscores (e.g., add_grades_table.sql).
- Sections: Always include `-- +goose Up` and `-- +goose Down` with complete reversals.

### Testing Migrations
1. **Syntax Validation**: Run `psql -f migrations/XXX.sql` on a test DB to check errors.
2. **Local Testing**:
   - Create test DB: `createdb rangkaiedu_test`
   - Run `goose up` on test DB.
   - Verify: Query tables (e.g., `SELECT * FROM users;`), check indexes (`\di`), relationships (joins).
   - Run `goose down` and confirm clean state.
3. **Integration**: After migration, test app connectivity (Phase 3).
4. **Edge Cases**: Test with data (e.g., FK violations), large datasets for performance.
5. **CI/CD**: Integrate into pipelines (e.g., GitHub Actions: run goose up on staging).

### Rollback Strategies
- **Immediate Rollback**: Use `goose down` for recent changes.
- **Partial Rollback**: `down-to N` for multi-step issues.
- **Manual Intervention**: If down fails, edit DB directly (e.g., DROP failed objects), then re-run.
- **Backup First**: Always `pg_dump` before up/down.
- **Version Control**: Commit migrations; never alter applied ones—create fix-up migrations instead.
- **Production**: Use blue-green deployments; rollback to previous release if needed.

## Confirmation of Initial Migration Testing
The initial migrations (001 and 002) have been created following Phase 1 schema exactly:
- **Extensions**: pgcrypto enabled for UUID generation.
- **Tables**: 12 tables with UUID PKs, timestamps, enums (CHECK), FKs (ON DELETE CASCADE/SET NULL), UNIQUEs.
- **Indexes**: 12+ as planned (e.g., idx_users_email, composite idx_class_teachers_class_subject).
- **Triggers**: Generic `update_updated_at_column()` for tables with updated_at (users, schools, etc.).
- **Seed Data**: Diverse samples for testing (e.g., teachers assigned to classes/subjects, students linked to parents/classes).
- **No Custom Functions**: Only triggers for auditing; no stored procs needed yet.
- **Validity**: Scripts are PostgreSQL-compatible, follow security (no raw passwords), performance (indexes for joins). Syntactically valid; recommend running on local DB via above commands for full verification.

For issues, check Goose logs or PostgreSQL error messages. Future phases (e.g., Phase 3 connectivity) will integrate these migrations.

**Version**: 1.0 | **Date**: 2025-09-28