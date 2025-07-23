# KT Chat Prototype

A comprehensive chat system prototype built for learning and technical verification purposes, featuring real-time messaging, file uploads, encryption, and AWS Cognito integration.

## 🎯 Purpose

This project serves as a learning platform for:
- Go language development and environment setup
- Chat system performance validation
- AWS Cognito integration understanding
- File upload and preview functionality
- Encryption logic implementation and testing

## 🏗️ Architecture

```
chat-prototype/
├── backend/           # Go backend application
├── frontend/          # Nuxt.js SPA (planned)
├── laravel/           # Laravel for master/settings (planned)
├── nginx/             # Reverse proxy configuration
├── scripts/           # Database and deployment scripts
├── docs/              # Documentation
└── docker-compose.yml # Development environment
```

## 🚀 Features

### Core Functionality
- **Real-time Chat**: WebSocket-based messaging with room support
- **File Upload**: Images, PDFs, and videos with preview capabilities
- **Message Encryption**: AES-256-GCM encryption for secure messaging
- **User Authentication**: AWS Cognito integration
- **Searchable Messages**: Encrypted message search functionality

### Technical Features
- **Go Backend**: Latest Go 1.24.5 with Gin framework
- **MySQL Database**: Persistent storage with GORM
- **WebSocket Hub**: Efficient real-time communication
- **File Management**: Upload, validation, and thumbnail generation
- **Docker Support**: Complete containerized development environment

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.24.5
- **Framework**: Gin (HTTP framework)
- **Database**: MySQL 8.0 with GORM
- **WebSocket**: Gorilla WebSocket
- **Authentication**: AWS Cognito + JWT
- **Encryption**: AES-256-GCM

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Reverse Proxy**: Nginx
- **Database**: MySQL 8.0
- **Deployment**: EC2 with Nginx reverse proxy

### Planned Components
- **Frontend**: Nuxt.js SPA
- **Admin Panel**: Laravel for master/settings management

## 📋 Requirements

### System Requirements
- Docker & Docker Compose
- Go 1.24.5+ (for local development)
- MySQL 8.0+ (for local development)
- Node.js 18+ (for frontend development)

### AWS Services
- AWS Cognito User Pool
- EC2 Instance (already provisioned: DEV-KTCHAT-WEB01)

## 🚀 Quick Start

### Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd chat-prototype
   ```

2. **Start the development environment**
   ```bash
   docker compose up -d    # Newer Docker CLI plugin (recommended)
   # OR
   docker-compose up -d    # Legacy standalone tool
   ```

3. **Access the application**
   - Backend API: http://localhost:8080
   - Health Check: http://localhost:8080/health
   - Nginx Proxy: http://localhost:80

### Local Development

1. **Set up the backend**
   ```bash
   cd backend
   cp env.example .env
   # Edit .env with your configuration
   go mod tidy
   go run cmd/server/main.go
   ```

2. **Set up MySQL**
   ```bash
   # Run the database initialization script
   mysql -u root -p < scripts/init.sql
   ```

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the backend directory:

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

## 📚 API Documentation

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

## 📊 Performance Considerations

### Load Testing
- Target: 500 concurrent users
- JMeter test scripts included
- Performance monitoring and optimization

### Scalability
- Database connection pooling
- WebSocket connection management
- Efficient message broadcasting
- File upload optimization

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

## 🧪 Testing

### Load Testing with JMeter
```bash
# Run JMeter tests
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

## 📁 Project Structure

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
│   ├── go.mod                 # Go dependencies
│   └── README.md              # Backend documentation
├── frontend/                  # Nuxt.js frontend (planned)
├── laravel/                   # Laravel admin panel (planned)
├── nginx/                     # Nginx configuration
│   ├── nginx.conf             # Main nginx config
│   └── conf.d/                # Server configurations
├── scripts/                   # Utility scripts
│   └── init.sql               # Database initialization
├── docs/                      # Documentation
├── docker-compose.yml         # Development environment
└── README.md                  # This file
```

## 🤝 Contributing

1. Follow Go coding standards
2. Add tests for new features
3. Update documentation
4. Use conventional commit messages
5. Ensure all tests pass before submitting

## 📝 Development Notes

### Go Development Environment
- Latest Go version (1.24.5)
- Standard Go project layout
- Comprehensive error handling
- Proper logging and monitoring

### Database Design
- Normalized schema for scalability
- Proper indexing for performance
- Foreign key relationships
- Audit trails (created_at, updated_at)

### Security Considerations
- Input validation and sanitization
- SQL injection prevention
- XSS protection
- CSRF protection
- Rate limiting

## 📞 Support

For questions and support:
- Check the documentation in each component directory
- Review the API documentation
- Check the deployment guide for production setup

## 📄 License

This project is for learning and technical verification purposes. 