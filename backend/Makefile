# KT Chat Backend Makefile

.PHONY: help build run test clean deps fmt lint docker-build docker-run docker-stop logs

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build the application
build: ## Build the application
	@echo "Building KT Chat Backend..."
	go build -o bin/server cmd/server/main.go
	@echo "Build complete! Binary: bin/server"

# Run the application
run: ## Run the application
	@echo "Starting KT Chat Backend..."
	go run cmd/server/main.go

# Run tests
test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	@echo "Clean complete!"

# Install dependencies
deps: ## Install Go dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed!"

# Format code
fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatting complete!"

# Lint code
lint: ## Lint Go code
	@echo "Linting code..."
	golangci-lint run
	@echo "Linting complete!"

# Build Docker image
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t ktchat-backend .
	@echo "Docker image built!"

# Run with Docker
docker-run: ## Run with Docker
	@echo "Starting services with Docker Compose..."
	docker compose up -d
	@echo "Services started!"

# Stop Docker services
docker-stop: ## Stop Docker services
	@echo "Stopping Docker services..."
	docker compose down
	@echo "Services stopped!"

# View logs
logs: ## View application logs
	@echo "Viewing logs..."
	docker compose logs -f backend

# Database operations
db-setup: ## Set up PostgreSQL database (local development)
	@echo "Setting up PostgreSQL database..."
	@echo "Please run the following commands manually:"
	@echo "sudo -u postgres psql"
	@echo "CREATE DATABASE ktchat;"
	@echo "CREATE USER ktchat WITH PASSWORD 'password';"
	@echo "GRANT ALL PRIVILEGES ON DATABASE ktchat TO ktchat;"
	@echo "GRANT ALL PRIVILEGES ON SCHEMA public TO ktchat;"
	@echo "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO ktchat;"
	@echo "\\q"

db-connect: ## Connect to PostgreSQL database
	@echo "Connecting to PostgreSQL database..."
	docker compose exec postgres psql -U ktchat -d ktchat

db-backup: ## Backup PostgreSQL database
	@echo "Backing up PostgreSQL database..."
	docker compose exec postgres pg_dump -U ktchat ktchat > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "Backup complete!"

db-restore: ## Restore PostgreSQL database from backup
	@echo "Restoring PostgreSQL database..."
	@read -p "Enter backup file name: " backup_file; \
	docker compose exec -T postgres psql -U ktchat -d ktchat < $$backup_file
	@echo "Restore complete!"

# Development helpers
dev-setup: deps fmt ## Set up development environment
	@echo "Development environment setup complete!"

dev-run: ## Run in development mode with hot reload
	@echo "Starting development server with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Production helpers
prod-build: ## Build for production
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server cmd/server/main.go
	@echo "Production build complete!"

prod-run: ## Run in production mode
	@echo "Starting in production mode..."
	ENVIRONMENT=production ./bin/server

# Health checks
health: ## Check application health
	@echo "Checking application health..."
	@curl -s http://localhost:8080/health | jq . || echo "Health check failed"

# Database health check
db-health: ## Check database health
	@echo "Checking database health..."
	docker compose exec postgres pg_isready -U ktchat

# Full system check
system-check: health db-health ## Check system health
	@echo "System health check complete!"

# Clean everything
clean-all: clean ## Clean everything including Docker
	@echo "Cleaning everything..."
	docker compose down -v
	docker system prune -f
	@echo "Full cleanup complete!" 