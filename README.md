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
- **User Authentication**: AWS Cognito integration with JWT tokens
- **Searchable Messages**: Encrypted message search functionality

### Technical Features
- **Go Backend**: Latest Go 1.24.5 with Gin framework
- **MySQL Database**: Persistent storage with GORM
- **WebSocket Hub**: Efficient real-time communication
- **File Management**: Upload, validation, and thumbnail generation
- **Docker Support**: Complete containerized development environment
- **AWS Integration**: Cognito User Pool for authentication

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.24.5
- **Framework**: Gin (HTTP framework)
- **Database**: MySQL 8.0 with GORM
- **WebSocket**: Gorilla WebSocket
- **Authentication**: AWS Cognito + Custom JWT
- **Encryption**: AES-256-GCM
- **AWS SDK**: AWS SDK v2 for Go

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

2. **Configure AWS Cognito (Optional)**
   ```bash
   # Follow the setup guide for AWS Cognito
   # docs/AWS_COGNITO_SETUP.md
   
   # Or use mock authentication for development
   # (No configuration needed)
   ```

3. **Start the development environment**
   ```bash
   docker compose up -d    # Newer Docker CLI plugin (recommended)
   # OR
   docker-compose up -d    # Legacy standalone tool
   ```

4. **Access the application**
   - Backend API: http://localhost:8080
   - Health Check: http://localhost:8080/health
   - Nginx Proxy: http://localhost:80

5. **Test the setup**
   ```bash
   # Run the test script
   ./scripts/test-cognito.sh
   ```

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

# AWS Cognito (Optional - for production)
AWS_REGION=ap-northeast-1
COGNITO_USER_POOL_ID=your-user-pool-id
COGNITO_CLIENT_ID=your-client-id
COGNITO_CLIENT_SECRET=your-client-secret

# AWS Credentials (for local development)
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

### AWS Cognito Setup

For production authentication, follow the comprehensive setup guide:

📖 **[AWS Cognito Setup Guide](docs/AWS_COGNITO_SETUP.md)**

The guide includes:
- Step-by-step User Pool creation
- App client configuration
- Environment variable setup
- Testing procedures
- Troubleshooting tips

## 📚 API Documentation

### Base URL
- **Local**: http://localhost:8080
- **Docker**: http://localhost:80 (via Nginx)

### Authentication Endpoints
- `POST /api/v1/auth/login` - User login (Cognito or mock)
- `POST /api/v1/auth/register` - User registration (Cognito)
- `POST /api/v1/auth/refresh` - Token refresh (Cognito)

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

## 🔐 Authentication Flow

### With AWS Cognito (Production)
1. User registers via `/api/v1/auth/register`
2. User confirms email (Cognito sends confirmation email)
3. User logs in via `/api/v1/auth/login`
4. Backend validates with Cognito and returns custom JWT
5. Frontend uses custom JWT for API calls

### Without AWS Cognito (Development)
1. Application falls back to mock authentication
2. Use username: `admin`, password: `password`
3. Backend generates custom JWT tokens
4. Same API flow as production

## 🧪 Testing

### Test Scripts
```bash
# Test Cognito integration
./scripts/test-cognito.sh

# Test basic functionality
./scripts/test-setup.sh
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
chat-prototype/
├── backend/                    # Go backend application
│   ├── cmd/server/            # Application entry point
│   ├── internal/              # Internal packages
│   │   ├── auth/              # AWS Cognito authentication
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
│   ├── test-setup.sh          # Environment test script
│   └── test-cognito.sh        # Cognito integration test
├── docs/                      # Documentation
│   ├── AWS_COGNITO_SETUP.md   # Cognito setup guide
│   ├── DOCKER_COMPOSE_GUIDE.md # Docker setup guide
│   └── SETUP_GUIDE.md         # General setup guide
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

- **[AWS Cognito Setup](docs/AWS_COGNITO_SETUP.md)** - Complete Cognito configuration guide
- **[Docker Setup](docs/DOCKER_COMPOSE_GUIDE.md)** - Docker environment setup
- **[General Setup](docs/SETUP_GUIDE.md)** - Basic setup instructions
- **[Backend README](backend/README.md)** - Backend-specific documentation

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is for learning and technical verification purposes.

## 🆘 Support

For issues and questions:
1. Check the documentation in the `docs/` folder
2. Review the troubleshooting sections
3. Check the test scripts for examples
4. Verify your AWS Cognito configuration 