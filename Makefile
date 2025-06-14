.PHONY: help db-start db-stop db-restart db-logs db-clean migrate-up migrate-down seed-data build run test clean

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Database commands
db-start: ## Start PostgreSQL database container
	docker-compose up -d postgres

db-stop: ## Stop PostgreSQL database container
	docker-compose stop postgres

db-restart: ## Restart PostgreSQL database container
	docker-compose restart postgres

db-logs: ## Show PostgreSQL database logs
	docker-compose logs -f postgres

db-clean: ## Stop and remove PostgreSQL container and volumes
	docker-compose down postgres
	docker volume rm commercify_postgres_data 2>/dev/null || true

# Migration commands
migrate-up: ## Run database migrations up
	docker-compose run --rm migrate -up

migrate-down: ## Run database migrations down
	docker-compose run --rm migrate -down

migrate-status: ## Show migration status
	docker-compose run --rm migrate -status

# Seed data
seed-data: ## Seed database with sample data
	docker-compose run --rm seed -all

# Application commands
build: ## Build the application
	go build -o bin/api ./cmd/api
	go build -o bin/migrate ./cmd/migrate
	go build -o bin/seed ./cmd/seed

run: db-start ## Run the application locally with database
	@echo "Starting database and waiting for it to be ready..."
	@sleep 3
	go run ./cmd/api

run-docker: ## Run the entire application stack with Docker
	docker-compose up -d

stop-docker: ## Stop the entire application stack
	docker-compose down

logs: ## Show application logs
	docker-compose logs -f api

# Development commands
test: ## Run tests
	go test ./...

test-verbose: ## Run tests with verbose output
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

# Database setup for development
dev-setup: db-start migrate-up seed-data ## Setup development environment (start db, migrate, seed)
	@echo "Development environment ready!"

dev-reset: db-clean db-start migrate-up seed-data ## Reset development environment
	@echo "Development environment reset!"

# Format and lint
fmt: ## Format Go code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

mod-tidy: ## Tidy Go modules
	go mod tidy
