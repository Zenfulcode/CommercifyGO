.PHONY: help db-start db-stop db-restart db-logs db-clean seed-data build run test clean docker-build docker-build-tag docker-push docker-build-push dev-sqlite dev-postgres

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development environment setup
dev-sqlite: ## Setup local development environment with SQLite
	@echo "Setting up SQLite development environment..."
	@cp .env.local .env 2>/dev/null || true
	@echo "Environment configured for SQLite. Starting application..."
	go run ./cmd/api

# Database commands (PostgreSQL)
db-start: ## Start PostgreSQL database container
	docker compose up -d postgres

db-stop: ## Stop PostgreSQL database container
	docker compose stop postgres

db-restart: ## Restart PostgreSQL database container
	docker compose restart postgres

db-logs: ## Show PostgreSQL database logs
	docker compose logs -f postgres

db-clean: ## Stop and remove PostgreSQL container and volumes
	docker compose down postgres
	docker volume rm commercify_postgres_data 2>/dev/null || true

# Seed data
seed-data: ## Seed database with sample data
	docker compose run --rm seed -all

# Application commands
build: ## Build the application
	go build -o bin/api ./cmd/api
	go build -o bin/seed ./cmd/seed
	go build -o bin/expire-checkouts ./cmd/expire-checkouts

run: db-start ## Run the application locally with database
	@echo "Starting database and waiting for it to be ready..."
	@sleep 3
	go run ./cmd/api

run-docker: ## Run the entire application stack with Docker (PostgreSQL)
	docker compose up -d

run-docker-sqlite: ## Run the application with Docker using SQLite
	docker compose -f docker-compose.local.yml up -d

stop-docker: ## Stop the entire application stack
	docker compose down

stop-docker-sqlite: ## Stop the SQLite application stack
	docker compose -f docker-compose.local.yml down

logs: ## Show application logs
	docker compose logs -f api

logs-sqlite: ## Show SQLite application logs
	docker compose -f docker-compose.local.yml logs -f api

# Docker image commands
docker-build: ## Build Docker image
	docker build -t ghcr.io/zenfulcode/commercifygo:latest .

docker-build-tag: ## Build Docker image with specific tag (use TAG=version)
	@if [ -z "$(TAG)" ]; then echo "Error: TAG is required. Use: make docker-build-tag TAG=v1.0.0"; exit 1; fi
	docker build -t ghcr.io/zenfulcode/commercifygo:$(TAG) -t ghcr.io/zenfulcode/commercifygo:latest -t ghcr.io/zenfulcode/commercifygo:dev .

docker-push: ## Push Docker image to registry (use REGISTRY and TAG)
	@if [ -z "$(REGISTRY)" ]; then echo "Error: REGISTRY is required. Use: make docker-push REGISTRY=your-registry.com"; exit 1; fi
	@if [ -z "$(TAG)" ]; then echo "Error: TAG is required. Use: make docker-push REGISTRY=your-registry.com TAG=v1.0.0"; exit 1; fi
# docker tag $(REGISTRY)commercifygo:$(TAG) $(REGISTRY)/commercifygo:$(TAG)
# docker tag $(REGISTRY)commercifygo:latest $(REGISTRY)/commercifygo:latest
	docker push $(REGISTRY)/commercifygo:$(TAG)
	docker push $(REGISTRY)/commercifygo:latest
	docker push $(REGISTRY)/commercifygo:dev

docker-build-push: docker-build-tag docker-push ## Build and push Docker image (use REGISTRY and TAG)

docker-dev-build: ## Build Docker image for development
	docker build -t ghcr.io/zenfulcode/commercifygo:dev .

# Development commands
test: ## Run tests
	go test ./...

test-verbose: ## Run tests with verbose output
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

# Database setup commands
dev-setup: ## Setup development environment with PostgreSQL (start db, seed)
	make db-start
	@sleep 3
	make seed-data
	@echo "Development environment ready with PostgreSQL!"

dev-reset: db-clean db-start seed-data ## Reset PostgreSQL development environment
	@echo "Development environment reset with PostgreSQL!"

dev-reset-sqlite: ## Reset SQLite development environment
	@echo "Resetting SQLite development environment..."
	@rm -f commercify.db 2>/dev/null || true
	@cp .env.local .env 2>/dev/null || true
	@echo "SQLite database reset!"

# Format and lint
fmt: ## Format Go code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

mod-tidy: ## Tidy Go modules
	go mod tidy

# Maintenance commands
expire-checkouts: ## Expire old checkouts manually
	go run ./cmd/expire-checkouts
