# KT Chat Development Environment Setup Guide

This guide will help you set up the KT Chat prototype development environment on your local machine.

## 🎯 Overview

The KT Chat prototype is a comprehensive chat system with the following components:
- **Go Backend**: Real-time chat API with WebSocket support
- **MySQL Database**: Persistent storage for messages and users
- **Nginx**: Reverse proxy for production-like setup
- **Docker**: Containerized development environment
- **Frontend**: Nuxt.js SPA (planned)
- **Laravel**: Admin panel for master/settings (planned)

## 📋 Prerequisites

### Required Software
- **Go 1.24.5+**: [Download from golang.org](https://golang.org/dl/)
- **Docker & Docker Compose**: [Download from docker.com](https://docker.com/)
- **Git**: [Download from git-scm.com](https://git-scm.com/)

### Optional Software (for local development)
- **MySQL 8.0+**: For local database development
- **Node.js 18+**: For frontend development (when implemented)

## 🚀 Quick Start

### 1. Clone the Repository
```bash
git clone <repository-url>
cd chat-prototype
```

### 2. Run the Test Script
```bash
./scripts/test-setup.sh
```

This script will verify your environment and provide feedback on what's working.

### 3. Start the Development Environment

#### Option A: Using Docker Compose (Recommended)
```bash
# Start all services (use the command that works on your system)
docker compose up -d    # Newer Docker CLI plugin (recommended)
# OR
docker-compose up -d    # Legacy standalone tool

# Check status
docker compose ps       # or docker-compose ps

# View logs
docker compose logs -f  # or docker-compose logs -f
```

#### Option B: Local Development
```bash
# Set up environment
cd backend
cp env.example .env
# Edit .env with your configuration

# Install dependencies
go mod tidy

# Run the application
go run cmd/server/main.go
```

### 4. Verify Installation

#### Test Health Endpoint
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "service": "ktchat-backend",
  "version": "1.0.0"
}
```

#### Test Authentication
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'
```

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the `backend` directory:

```env
# Environment
ENVIRONMENT=development

# Database
DATABASE_URL=ktchat:password@tcp(localhost:3306)/ktchat?charset=utf8mb4&parseTime=True&loc=Local

# JWT
JWT_SECRET=your-secret-key-change-in-production

# AWS Cognito
AWS_REGION=ap-northeast-1
COGNITO_USER_POOL_ID=your-user-pool-id
COGNITO_CLIENT_ID=your-client-id
COGNITO_CLIENT_SECRET=your-client-secret

# File Upload
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=52428800

# Encryption
ENCRYPTION_KEY=your-encryption-key-32-bytes-long

# Server
PORT=8080
```

### Database Setup

#### Using Docker (Automatic)
The database will be automatically set up when using Docker Compose.

#### Using Local MySQL
```bash
# Connect to MySQL
mysql -u root -p

# Run the initialization script
source scripts/init.sql
```

## 🏗️ Project Structure

```
chat-prototype/
├── backend/                    # Go backend application
│   ├── cmd/server/            # Application entry point
│   ├── internal/              # Internal packages
│   │   ├── auth/              # Authentication logic
│   │   ├── chat/              # Chat service
│   │   ├── database/          # Database operations
│   │   ├── encryption/        # Message encryption
│   │   ├── file/              # File handling
│   │   ├── models/            # Database models
│   │   └── websocket/         # WebSocket management
│   ├── pkg/                   # Public packages
│   │   ├── config/            # Configuration
│   │   ├── middleware/        # HTTP middleware
│   │   └── utils/             # Utilities
│   ├── Dockerfile             # Backend container
│   ├── Makefile               # Development commands
│   ├── go.mod                 # Go dependencies
│   └── README.md              # Backend documentation
├── frontend/                  # Nuxt.js frontend (planned)
├── laravel/                   # Laravel admin panel (planned)
├── nginx/                     # Nginx configuration
│   ├── nginx.conf             # Main nginx config
│   └── conf.d/                # Server configurations
├── scripts/                   # Utility scripts
│   ├── init.sql               # Database initialization
│   └── test-setup.sh          # Environment test script
├── docs/                      # Documentation
├── docker-compose.yml         # Development environment
└── README.md                  # Main project documentation
```

## 🛠️ Development Commands

### Using Makefile (Backend)
```bash
cd backend

# Build the application
make build

# Run the application
make run

# Run tests
make test

# Format code
make fmt

# Install dependencies
make deps

# View all commands
make help
```

### Using Docker Compose
```bash
# Start all services (use the command that works on your system)
docker compose up -d    # Newer Docker CLI plugin (recommended)
# OR
docker-compose up -d    # Legacy standalone tool

# Stop all services
docker compose down     # or docker-compose down

# View logs
docker compose logs -f  # or docker-compose logs -f

# Rebuild and restart
docker compose up -d --build  # or docker-compose up -d --build
```

## 📚 API Documentation

### Base URL
- **Local**: http://localhost:8080
- **Docker**: http://localhost:80 (via Nginx)

### Authentication Endpoints
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/refresh` - Token refresh

### Chat Endpoints
- `GET /api/v1/chat/rooms` - Get all chat rooms
- `POST /api/v1/chat/rooms` - Create a new chat room
- `GET /api/v1/chat/rooms/:roomID/messages` - Get messages for a room
- `POST /api/v1/chat/rooms/:roomID/join` - Join a chat room
- `POST /api/v1/chat/rooms/:roomID/leave` - Leave a chat room

### WebSocket
- `GET /api/v1/ws` - Real-time chat connection

### Health Check
- `GET /health` - Application health status

## 🔐 Security Features

### Message Encryption
- All chat messages are encrypted using AES-256-GCM
- Encryption keys are configurable via environment variables
- Searchable encryption implementation for message search

### Authentication
- JWT-based authentication
- AWS Cognito integration for user management
- Rate limiting on API endpoints

### File Security
- File type validation
- File size limits
- Secure file storage with unique naming

## 🧪 Testing

### Load Testing with JMeter
```bash
# Run JMeter tests (when implemented)
jmeter -n -t tests/load-test.jmx -l results.jtl
```

### API Testing
```bash
# Test health endpoint
curl http://localhost:8080/health

# Test authentication
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'
```

## 🚀 Deployment

### EC2 Deployment
The application is configured for deployment on:
- **Instance**: DEV-KTCHAT-WEB01 (54.179.141.95)
- **Document Root**: `/var/www/dev-ktchat/public`
- **Go Path**: `/usr/local/go/bin/go`
- **Port**: 8080 (Go app), 80 (Nginx)
- **Basic Auth**: ktchat / s9RnNyai

### Production Setup
1. Set `ENVIRONMENT=production`
2. Configure production database
3. Set secure encryption keys
4. Configure AWS Cognito credentials
5. Set up SSL certificates
6. Configure monitoring and logging

## 🐛 Troubleshooting

### Common Issues

#### Port Already in Use
```bash
# Check what's using the port
lsof -i :8080

# Kill the process
kill -9 <PID>
```

#### Database Connection Issues
```bash
# Check if MySQL is running
docker compose ps mysql    # or docker-compose ps mysql

# View MySQL logs
docker compose logs mysql  # or docker-compose logs mysql
```

#### Docker Issues
```bash
# Clean up Docker
docker system prune -a

# Rebuild containers
docker compose up -d --build  # or docker-compose up -d --build
```

#### Go Module Issues
```bash
cd backend
go mod tidy
go mod download
```

### Getting Help

1. Check the logs: `docker compose logs -f` (or `docker-compose logs -f`)
2. Run the test script: `./scripts/test-setup.sh`
3. Check the documentation in each component directory
4. Verify your environment variables

## 📞 Support

For questions and support:
- Check the documentation in each component directory
- Review the API documentation
- Check the deployment guide for production setup
- Run the test script to verify your environment

## 🎉 Next Steps

After setting up the development environment:

1. **Explore the API**: Test the endpoints using curl or Postman
2. **Set up AWS Cognito**: Configure user authentication
3. **Implement Frontend**: Build the Nuxt.js SPA
4. **Add Laravel Admin**: Implement master/settings functionality
5. **Load Testing**: Use JMeter to test performance
6. **Deploy to EC2**: Set up production environment

Happy coding! 🚀 