# Database Setup Guide

This document explains how to configure and use different databases with Commercify.

## Supported Databases

Commercify supports two database backends:

- **SQLite**: Recommended for local development, testing, and small deployments
- **PostgreSQL**: Recommended for production and larger deployments

## Quick Start

### SQLite (Recommended for Development)

The fastest way to get started:

```bash
# Copy local development environment
cp .env.local .env

# Start the application (database will be created automatically)
make dev-sqlite
```

This will:
- Create a SQLite database file (`commercify.db`) in the project root
- Run all database migrations automatically
- Start the API server

### PostgreSQL (Production)

For production or when you need a full-featured database:

```bash
# Start PostgreSQL with Docker
make db-start

# Setup environment and run migrations
make dev-setup

# Start the application
go run cmd/api/main.go
```

## Environment Configuration

The database is configured via environment variables. You can either:

1. Use the provided environment templates:
   - `.env.local` - SQLite configuration
   - `.env.example` - PostgreSQL configuration template
   - `.env.production` - Production PostgreSQL template

2. Set environment variables directly:

### SQLite Configuration

```bash
DB_DRIVER=sqlite
DB_NAME=commercify.db
DB_DEBUG=false
```

### PostgreSQL Configuration

```bash
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=commercify
DB_SSL_MODE=disable
DB_DEBUG=false
```

## Make Commands

The project includes helpful Make commands for database management:

### SQLite Commands

```bash
make dev-sqlite              # Start application with SQLite
make dev-setup-sqlite         # Setup SQLite environment
make dev-reset-sqlite         # Reset SQLite database
```

### PostgreSQL Commands

```bash
make dev-postgres             # Start application with PostgreSQL
make dev-setup               # Setup PostgreSQL environment (start DB, migrate, seed)
make dev-reset               # Reset PostgreSQL environment

# Database container management
make db-start                # Start PostgreSQL container
make db-stop                 # Stop PostgreSQL container
make db-restart              # Restart PostgreSQL container
make db-logs                 # View database logs
make db-clean                # Stop and remove database container and volumes
```

### Migration Commands

```bash
make migrate-up              # Run pending migrations
make migrate-down            # Rollback last migration
make migrate-status          # Show migration status
make seed-data              # Seed database with sample data
```

## Docker Setup

### SQLite with Docker

```bash
# Run with SQLite in Docker
make run-docker-sqlite

# Stop SQLite Docker setup
make stop-docker-sqlite
```

### PostgreSQL with Docker

```bash
# Run full stack with PostgreSQL
make run-docker

# Stop PostgreSQL Docker setup
make stop-docker
```

## Database Files and Cleanup

### SQLite

- Database file: `commercify.db` (created in project root)
- To reset: Delete the file or run `make dev-reset-sqlite`
- Backup: Copy the `commercify.db` file

### PostgreSQL

- Data persisted in Docker volume: `commercify_postgres_data`
- To reset: Run `make dev-reset` or `make db-clean`
- Backup: Use PostgreSQL backup tools (`pg_dump`)

## Switching Between Databases

You can easily switch between databases by changing your environment configuration:

1. **SQLite to PostgreSQL**:
   ```bash
   cp .env.example .env
   # Edit .env to set DB_DRIVER=postgres and configure connection
   make db-start
   ```

2. **PostgreSQL to SQLite**:
   ```bash
   cp .env.local .env
   make db-stop  # Stop PostgreSQL if running
   ```

## Production Deployment

For production deployments:

1. Use PostgreSQL as the database backend
2. Configure environment variables securely
3. Use SSL mode for database connections (`DB_SSL_MODE=require`)
4. Set strong passwords and limit database access
5. Regular backups are recommended

### Environment Template for Production

```bash
DB_DRIVER=postgres
DB_HOST=your-postgres-host
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-secure-password
DB_NAME=commercify_production
DB_SSL_MODE=require
DB_DEBUG=false
```

## Troubleshooting

### Common Issues

1. **"Database file locked" (SQLite)**
   - Ensure no other instances are running
   - Check file permissions

2. **"Connection refused" (PostgreSQL)**
   - Ensure PostgreSQL is running (`make db-start`)
   - Check connection parameters in `.env`

3. **Migration errors**
   - Check database permissions
   - Ensure database exists
   - Run `make migrate-status` to check migration state

### Getting Help

- Check the logs: `make db-logs` (PostgreSQL) or `make logs-sqlite` (SQLite Docker)
- Verify configuration: Check your `.env` file
- Reset environment: Use `make dev-reset` or `make dev-reset-sqlite`
