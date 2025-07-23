#!/bin/bash

# KT Chat Development Environment Test Script

echo "🧪 Testing KT Chat Development Environment"
echo "=========================================="

# Check if Go is installed
echo "📋 Checking Go installation..."
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo "✅ Go is installed: $GO_VERSION"
else
    echo "❌ Go is not installed"
    exit 1
fi

# Check if Docker is installed
echo "📋 Checking Docker installation..."
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
    echo "✅ Docker is installed: $DOCKER_VERSION"
else
    echo "❌ Docker is not installed"
fi

# Check if Docker Compose is installed
echo "📋 Checking Docker Compose installation..."
if command -v docker compose &> /dev/null; then
    COMPOSE_VERSION=$(docker compose --version | awk '{print $3}' | sed 's/,//')
    echo "✅ Docker Compose (new) is installed: $COMPOSE_VERSION"
    DOCKER_COMPOSE_CMD="docker compose"
elif command -v docker-compose &> /dev/null; then
    COMPOSE_VERSION=$(docker-compose --version | awk '{print $3}' | sed 's/,//')
    echo "✅ Docker Compose (legacy) is installed: $COMPOSE_VERSION"
    DOCKER_COMPOSE_CMD="docker-compose"
else
    echo "❌ Docker Compose is not installed"
    DOCKER_COMPOSE_CMD=""
fi

# Check if MySQL is installed (for local development)
echo "📋 Checking MySQL installation..."
if command -v mysql &> /dev/null; then
    echo "✅ MySQL is installed"
else
    echo "⚠️  MySQL is not installed (will use Docker version)"
fi

# Test Go application build
echo "📋 Testing Go application build..."
cd backend
if go build cmd/server/main.go; then
    echo "✅ Go application builds successfully"
else
    echo "❌ Go application build failed"
    exit 1
fi

# Test Docker build
echo "📋 Testing Docker build..."
if docker build -t ktchat-backend-test . &> /dev/null; then
    echo "✅ Docker build successful"
    docker rmi ktchat-backend-test &> /dev/null
else
    echo "❌ Docker build failed"
fi

# Test Docker Compose
echo "📋 Testing Docker Compose..."
cd ..
if [ -n "$DOCKER_COMPOSE_CMD" ]; then
    if $DOCKER_COMPOSE_CMD config &> /dev/null; then
        echo "✅ Docker Compose configuration is valid"
    else
        echo "❌ Docker Compose configuration is invalid"
    fi
else
    echo "❌ Docker Compose not available for testing"
fi

# Check environment file
echo "📋 Checking environment configuration..."
if [ -f "backend/.env" ]; then
    echo "✅ Environment file exists"
else
    echo "⚠️  Environment file not found (copy from env.example)"
fi

# Check project structure
echo "📋 Checking project structure..."
REQUIRED_DIRS=("backend" "frontend" "laravel" "nginx" "scripts" "docs")
for dir in "${REQUIRED_DIRS[@]}"; do
    if [ -d "$dir" ]; then
        echo "✅ $dir directory exists"
    else
        echo "⚠️  $dir directory missing"
    fi
done

echo ""
echo "🎉 Environment test completed!"
echo ""
echo "Next steps:"
echo "1. Set up your .env file: cp backend/env.example backend/.env"
if [ -n "$DOCKER_COMPOSE_CMD" ]; then
    echo "2. Start the development environment: $DOCKER_COMPOSE_CMD up -d"
else
    echo "2. Start the development environment: docker compose up -d (or docker-compose up -d)"
fi
echo "3. Or run locally: cd backend && make run"
echo ""
echo "For more information, see the README.md files in each directory." 