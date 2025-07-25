# KT Chat Backend

A Go-based backend for the KT Chat prototype featuring real-time messaging, file uploads, encryption, and AWS Cognito integration.

## 🚀 Features

- **Real-time Chat**: WebSocket-based messaging with room support
- **File Upload**: Images, PDFs, and videos with preview capabilities
- **Message Encryption**: AES-256-GCM encryption for secure messaging
- **User Authentication**: AWS Cognito integration with JWT tokens
- **Searchable Messages**: Encrypted message search functionality

## 🛠️ Technology Stack

- **Language**: Go 1.24.5
- **Framework**: Gin (HTTP framework)
- **Database**: PostgreSQL 15 with GORM
- **WebSocket**: Gorilla WebSocket
- **Authentication**: AWS Cognito + Custom JWT
- **Encryption**: AES-256-GCM
- **AWS SDK**: AWS SDK v2 for Go

## 📋 Requirements

- Go 1.24.5 or later
- PostgreSQL 15 or later
- Docker & Docker Compose (for containerized development)

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the backend directory:

```env
# Environment
ENVIRONMENT=development

# Database (PostgreSQL)
DATABASE_URL=postgres://ktchat:password@localhost:5432/ktchat?sslmode=disable

# JWT
JWT_SECRET=your-secret-key-change-in-production

# AWS Cognito (Optional)
AWS_REGION=ap-northeast-1
COGNITO_USER_POOL_ID=your-user-pool-id
COGNITO_CLIENT_ID=your-client-id
COGNITO_CLIENT_SECRET=your-client-secret

# AWS Credentials
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key

# File Upload
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=52428800

# Encryption
ENCRYPTION_KEY=your-encryption-key-32-bytes-long

# Server
PORT=8080
```

### Database Setup

1. **Install PostgreSQL 15+**
   ```bash
   # macOS
   brew install postgresql@15
   
   # Ubuntu/Debian
   sudo apt-get install postgresql-15
   ```

2. **Set up the database**
   ```bash
   # Create database and user
   sudo -u postgres psql
   CREATE DATABASE ktchat;
   CREATE USER ktchat WITH PASSWORD 'password';
   GRANT ALL PRIVILEGES ON DATABASE ktchat TO ktchat;
   GRANT ALL PRIVILEGES ON SCHEMA public TO ktchat;
   ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO ktchat;
   \q
   ```

## 🚀 Quick Start

### Using Docker (Recommended)

```bash
# Start all services
docker compose up -d

# Check status
docker compose ps

# View logs
docker compose logs backend
```

### Local Development

```bash
# Install dependencies
go mod download

# Run the application
go run cmd/server/main.go

# Or use air for hot reload
air
```

## 📚 API Documentation

### Base URL
- **Local**: http://localhost:8080
- **Docker**: http://localhost:80 (via Nginx)

### Authentication Endpoints
- `POST /api/v1/auth/login` - User login (Cognito or mock)
- `POST /api/v1/auth/register` - User registration (Cognito)
- `POST /api/v1/auth/refresh` - Token refresh (Cognito)
- `POST /api/v1/auth/confirm` - Confirm user account
- `POST /api/v1/auth/resend-confirmation` - Resend confirmation code

### Chat Endpoints
- `GET /api/v1/chat/rooms` - Get all chat rooms
- `POST /api/v1/chat/rooms` - Create a new chat room
- `GET /api/v1/chat/rooms/:roomID/messages` - Get messages for a room
- `POST /api/v1/chat/rooms/:roomID/join` - Join a chat room
- `POST /api/v1/chat/rooms/:roomID/leave` - Leave a chat room

### WebSocket
- `GET /api/v1/ws` - WebSocket connection for real-time chat

### Health Check
- `GET /health` - Health check endpoint with Cognito status

## 🧪 Testing

### Test Scripts
```bash
# Test basic functionality
./scripts/test-setup.sh

# Test AWS Cognito integration
./scripts/test-cognito.sh

# Test user confirmation flow
./scripts/test-confirmation.sh
```

### Manual Testing
```bash
# Health check
curl http://localhost:8080/health

# Login (mock auth)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'
```

## 🏗️ Project Structure

```
backend/
├── cmd/server/            # Application entry point
├── internal/              # Internal packages
│   ├── auth/              # AWS Cognito authentication
│   ├── chat/              # Chat service
│   ├── database/          # Database operations
│   ├── encryption/        # Message encryption
│   ├── file/              # File handling
│   ├── models/            # Database models
│   └── websocket/         # WebSocket management
├── pkg/                   # Public packages
│   ├── config/            # Configuration
│   ├── middleware/        # HTTP middleware
│   └── utils/             # Utilities
├── Dockerfile             # Backend container
├── Makefile               # Development commands
├── go.mod                 # Go dependencies
└── README.md              # Backend documentation
```

## 🛠️ Development Commands

### Using Makefile
```bash
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
# Start all services
docker compose up -d

# Stop all services
docker compose down

# View logs
docker compose logs -f

# Rebuild and restart
docker compose up -d --build
```

## 🔒 Security Features

- **Message Encryption**: AES-256-GCM encryption for all messages
- **JWT Tokens**: Secure token-based authentication
- **AWS Cognito**: Enterprise-grade user management
- **File Validation**: Type and size validation for uploads
- **Rate Limiting**: Basic rate limiting on API endpoints
- **CORS Protection**: Configured CORS headers

## 🚀 Deployment

### Production Checklist
- [ ] Configure AWS Cognito User Pool
- [ ] Set up IAM roles (no access keys)
- [ ] Enable HTTPS
- [ ] Configure proper CORS origins
- [ ] Set up monitoring and logging
- [ ] Enable MFA for user accounts
- [ ] Configure custom domain for Cognito

### EC2 Deployment
The application is designed to be deployed on the existing EC2 instance (DEV-KTCHAT-WEB01) with Nginx as a reverse proxy.

## 📖 Documentation

- **[AWS Cognito Setup](../docs/AWS_COGNITO_SETUP.md)** - Complete Cognito configuration guide
- **[Docker Setup](../docs/DOCKER_COMPOSE_GUIDE.md)** - Docker environment setup
- **[General Setup](../docs/SETUP_GUIDE.md)** - Basic setup instructions

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