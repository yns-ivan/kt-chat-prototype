#!/bin/bash

# Docker Compose Wrapper Script
# This script automatically detects and uses the correct Docker Compose command

# Function to detect Docker Compose command
detect_docker_compose() {
    if command -v docker compose &> /dev/null; then
        echo "docker compose"
    elif command -v docker-compose &> /dev/null; then
        echo "docker-compose"
    else
        echo "error"
    fi
}

# Get the Docker Compose command
DOCKER_COMPOSE_CMD=$(detect_docker_compose)

# Check if Docker Compose is available
if [ "$DOCKER_COMPOSE_CMD" = "error" ]; then
    echo "❌ Error: Docker Compose is not installed"
    echo "Please install Docker Compose or the Docker CLI plugin"
    exit 1
fi

# Display which command is being used
echo "🔧 Using: $DOCKER_COMPOSE_CMD"

# Execute the command with all arguments
exec $DOCKER_COMPOSE_CMD "$@" 