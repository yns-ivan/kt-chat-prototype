# Makefile for KT Chat Backend

.PHONY: help build run test clean docker-build docker-run

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	@echo "Running application..."
	go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t ktchat-backend .

# Run with Docker Compose
docker-run:
	@echo "Starting services with Docker Compose..."
	../scripts/docker-compose-wrapper.sh up -d

# Stop Docker Compose
docker-stop:
	@echo "Stopping Docker Compose services..."
	../scripts/docker-compose-wrapper.sh down

# View logs
logs:
	@echo "Viewing logs..."
	../scripts/docker-compose-wrapper.sh logs -f

# Database setup
db-setup:
	@echo "Setting up database..."
	mysql -u root -p < ../scripts/init.sql

# Generate encryption key
gen-key:
	@echo "Generating encryption key..."
	openssl rand -base64 32

# Health check
health:
	@echo "Checking application health..."
	curl -f http://localhost:8080/health || echo "Application not running"

# Development setup
dev-setup: deps build
	@echo "Development setup complete!"
	@echo "Run 'make run' to start the application"
	@echo "Run 'make docker-run' to start with Docker Compose" 