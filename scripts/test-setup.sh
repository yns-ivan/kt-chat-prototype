#!/bin/bash

# KT Chat - Environment Test Script
# This script tests the development environment setup

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Testing KT Chat Development Environment${NC}"
echo "=============================================="
echo ""

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
        if [ "$3" != "" ]; then
            echo -e "${YELLOW}  Note: $3${NC}"
        fi
    fi
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# 1. Check Docker installation
echo -e "${BLUE}1. Checking Docker installation...${NC}"
if command_exists docker; then
    DOCKER_VERSION=$(docker --version)
    print_status 0 "Docker is installed: $DOCKER_VERSION"
else
    print_status 1 "Docker is not installed" "Install Docker Desktop from https://docker.com"
    exit 1
fi

# 2. Check Docker Compose
echo -e "${BLUE}2. Checking Docker Compose...${NC}"
if docker compose version >/dev/null 2>&1; then
    COMPOSE_VERSION=$(docker compose version --short)
    print_status 0 "Docker Compose is available: $COMPOSE_VERSION"
    COMPOSE_CMD="docker compose"
elif command_exists docker-compose; then
    COMPOSE_VERSION=$(docker-compose --version | cut -d' ' -f3 | cut -d',' -f1)
    print_status 0 "Docker Compose is available: $COMPOSE_VERSION"
    COMPOSE_CMD="docker-compose"
else
    print_status 1 "Docker Compose is not available" "Install Docker Compose or update Docker Desktop"
    exit 1
fi

# 3. Check Go installation
echo -e "${BLUE}3. Checking Go installation...${NC}"
if command_exists go; then
    GO_VERSION=$(go version | cut -d' ' -f3)
    print_status 0 "Go is installed: $GO_VERSION"
    
    # Check Go version
    GO_MAJOR=$(echo $GO_VERSION | cut -d'.' -f2)
    if [ "$GO_MAJOR" -ge 24 ]; then
        print_status 0 "Go version is compatible (1.24+)"
    else
        print_status 1 "Go version is too old" "Update to Go 1.24 or later"
    fi
else
    print_status 1 "Go is not installed" "Install Go from https://golang.org/dl/"
fi

# 4. Check PostgreSQL installation (for local development)
echo -e "${BLUE}4. Checking PostgreSQL installation...${NC}"
if command_exists psql; then
    PSQL_VERSION=$(psql --version | cut -d' ' -f3)
    print_status 0 "PostgreSQL client is installed: $PSQL_VERSION"
else
    print_status 0 "PostgreSQL client is not installed" "Will use Docker version for development"
fi

# 5. Check Node.js installation (for frontend)
echo -e "${BLUE}5. Checking Node.js installation...${NC}"
if command_exists node; then
    NODE_VERSION=$(node --version)
    print_status 0 "Node.js is installed: $NODE_VERSION"
    
    # Check Node.js version
    NODE_MAJOR=$(echo $NODE_VERSION | cut -d'v' -f2 | cut -d'.' -f1)
    if [ "$NODE_MAJOR" -ge 18 ]; then
        print_status 0 "Node.js version is compatible (18+)"
    else
        print_status 1 "Node.js version is too old" "Update to Node.js 18 or later"
    fi
else
    print_status 0 "Node.js is not installed" "Will use Docker version for frontend development"
fi

# 6. Check if Docker is running
echo -e "${BLUE}6. Checking if Docker is running...${NC}"
if docker info >/dev/null 2>&1; then
    print_status 0 "Docker is running"
else
    print_status 1 "Docker is not running" "Start Docker Desktop"
    exit 1
fi

# 7. Check if ports are available
echo -e "${BLUE}7. Checking if required ports are available...${NC}"

check_port() {
    local port=$1
    local service=$2
    if lsof -i :$port >/dev/null 2>&1; then
        print_status 1 "Port $port is in use" "Stop the service using port $port or change the port in docker-compose.yml"
    else
        print_status 0 "Port $port is available"
    fi
}

check_port 8080 "Backend API"
check_port 5432 "PostgreSQL"
check_port 3000 "Frontend"
check_port 80 "Nginx"

# 8. Check if .env file exists
echo -e "${BLUE}8. Checking environment configuration...${NC}"
if [ -f ".env" ]; then
    print_status 0 ".env file exists"
    
    # Check for required variables
    if grep -q "COGNITO_USER_POOL_ID" .env; then
        print_status 0 "AWS Cognito configuration found"
    else
        print_status 0 "AWS Cognito not configured" "Will use mock authentication for development"
    fi
else
    print_status 1 ".env file not found" "Copy .env.example to .env and configure your settings"
fi

# 9. Test Docker Compose configuration
echo -e "${BLUE}9. Testing Docker Compose configuration...${NC}"
if [ -f "docker-compose.yml" ]; then
    print_status 0 "docker-compose.yml found"
    
    # Validate docker-compose.yml
    if $COMPOSE_CMD config >/dev/null 2>&1; then
        print_status 0 "Docker Compose configuration is valid"
    else
        print_status 1 "Docker Compose configuration is invalid"
        exit 1
    fi
else
    print_status 1 "docker-compose.yml not found"
    exit 1
fi

# 10. Check backend dependencies
echo -e "${BLUE}10. Checking backend dependencies...${NC}"
if [ -f "backend/go.mod" ]; then
    print_status 0 "Go module file found"
    
    # Check if dependencies are downloaded
    if [ -f "backend/go.sum" ]; then
        print_status 0 "Go dependencies are available"
    else
        print_status 0 "Go dependencies need to be downloaded" "Run 'cd backend && go mod download'"
    fi
else
    print_status 1 "Go module file not found"
fi

echo ""
echo -e "${GREEN}✅ Environment test completed!${NC}"
echo ""

# Summary and next steps
echo -e "${BLUE}📋 Summary:${NC}"
echo "=========="

if command_exists docker && docker info >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Docker environment is ready${NC}"
else
    echo -e "${RED}✗ Docker environment needs setup${NC}"
fi

if [ -f ".env" ]; then
    echo -e "${GREEN}✓ Environment configuration is ready${NC}"
else
    echo -e "${YELLOW}⚠ Environment configuration needs setup${NC}"
fi

echo ""
echo -e "${BLUE}🚀 Next Steps:${NC}"
echo "=============="

if [ -f ".env" ]; then
    echo "1. Start the development environment:"
    echo "   $COMPOSE_CMD up -d"
    echo ""
    echo "2. Check if services are running:"
    echo "   $COMPOSE_CMD ps"
    echo ""
    echo "3. Test the backend:"
    echo "   curl http://localhost:8080/health"
    echo ""
    echo "4. Test AWS Cognito integration:"
    echo "   ./scripts/test-cognito.sh"
else
    echo "1. Set up environment configuration:"
    echo "   cp .env.example .env"
    echo "   # Edit .env with your settings"
    echo ""
    echo "2. Then run the steps above"
fi

echo ""
echo -e "${BLUE}🔗 Useful Commands:${NC}"
echo "=================="
echo "View logs: $COMPOSE_CMD logs -f"
echo "Stop services: $COMPOSE_CMD down"
echo "Rebuild: $COMPOSE_CMD up -d --build"
echo "Database: $COMPOSE_CMD exec postgres psql -U ktchat -d ktchat"
echo "Backend shell: $COMPOSE_CMD exec backend sh"
echo ""
echo -e "${YELLOW}📚 Documentation:${NC}"
echo "================="
echo "Setup Guide: docs/SETUP_GUIDE.md"
echo "Docker Guide: docs/DOCKER_COMPOSE_GUIDE.md"
echo "AWS Cognito: docs/AWS_COGNITO_SETUP.md" 