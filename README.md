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
- **File Storage**: Local filesystem or AWS S3 storage options
- **Message Encryption**: AES-256-GCM encryption for secure messaging
- **User Authentication**: AWS Cognito integration with JWT tokens
- **Searchable Messages**: Encrypted message search functionality

### Technical Features
- **Go Backend**: Latest Go 1.24.5 with Gin framework
- **PostgreSQL Database**: High-performance database optimized for chat systems
- **WebSocket Hub**: Efficient real-time communication
- **File Management**: Upload, validation, and thumbnail generation
- **Docker Support**: Complete containerized development environment
- **AWS Integration**: Cognito User Pool for authentication

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.24.5
- **Framework**: Gin (HTTP framework)
- **Database**: PostgreSQL 15 with GORM
- **WebSocket**: Gorilla WebSocket
- **Authentication**: AWS Cognito + Custom JWT
- **Encryption**: AES-256-GCM
- **AWS SDK**: AWS SDK v2 for Go

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Reverse Proxy**: Nginx
- **Database**: PostgreSQL 15 (optimized for chat systems)
- **Deployment**: EC2 with Nginx reverse proxy

### Planned Components
- **Frontend**: Nuxt.js SPA
- **Admin Panel**: Laravel for master/settings management

## 📋 Requirements

### System Requirements
- Docker & Docker Compose
- Go 1.24.5+ (for local development)
- PostgreSQL 15+ (for local development)
- Node.js 18+ (for frontend development)

### AWS Services
- AWS Cognito User Pool
- AWS S3 Bucket (for file storage)
- EC2 Instance (already provisioned: DEV-KTCHAT-WEB01)

## 🚀 Quick Start

### Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd chat-prototype
   ```

2. **Configure AWS Services (Optional)**
   ```bash
   # Follow the setup guide for AWS Cognito
   # docs/AWS_COGNITO_SETUP.md
   
   # Follow the setup guide for S3 file storage
   # docs/S3_SETUP.md
   
   # Or use mock authentication and local storage for development
   # (No configuration needed)
   ```

3. **Start the development environment**
   ```bash
   docker compose up -d    # Newer Docker CLI plugin (recommended)
   # OR
   docker-compose up -d    # Legacy standalone tool
   ```

4. **Verify the services**
   ```bash
   # Check if all services are running
   docker compose ps
   
   # Check backend health
   curl http://localhost:8080/health
   
   # Check PostgreSQL connection
   docker compose exec postgres psql -U ktchat -d ktchat -c "\dt"
   ```

### Manual Setup (Alternative)

1. **Install PostgreSQL 15+**
   ```bash
   # macOS
   brew install postgresql@15
   
   # Ubuntu/Debian
   sudo apt-get install postgresql-15
   
   # Windows
   # Download from https://www.postgresql.org/download/windows/
   ```

2. **Set up the database**
   ```bash
   # Create database and user
   sudo -u postgres psql
   CREATE DATABASE ktchat;
   CREATE USER ktchat WITH PASSWORD 'password';
   GRANT ALL PRIVILEGES ON DATABASE ktchat TO ktchat;
   \q
   ```

3. **Configure environment variables**
   ```bash
   # Create .env file in the root directory
   cp .env.example .env
   
   # Update the DATABASE_URL for PostgreSQL
   DATABASE_URL=postgres://ktchat:password@localhost:5432/ktchat?sslmode=disable
   ```

4. **Run the backend**
   ```bash
   cd backend
   go mod download
   go run cmd/server/main.go
   ```

## 🧪 Testing

### Test AWS Cognito Integration
```bash
# Test the complete authentication flow
./scripts/test-cognito.sh

# Test user confirmation flow
./scripts/test-confirmation.sh
```

### Test Database Connection
```bash
# Connect to PostgreSQL
docker compose exec postgres psql -U ktchat -d ktchat

# List tables
\dt

# Check user data
SELECT * FROM users LIMIT 5;
```

## 📊 Database Schema

### Key Tables
- **users**: User accounts with Cognito integration
- **chat_rooms**: Chat rooms with privacy settings
- **room_participants**: User participation tracking
- **messages**: Encrypted chat messages
- **file_attachments**: File uploads with metadata

### PostgreSQL Optimizations
- **UUID Primary Keys**: Using PostgreSQL's native UUID support
- **Indexed Fields**: Optimized indexes for chat queries
- **JSON Support**: Ready for message metadata storage
- **Full-Text Search**: Prepared for message search functionality

## 🔧 Development

### Backend Development
```bash
cd backend

# Install dependencies
go mod download

# Run with hot reload (requires air)
air

# Run tests
go test ./...

# Build for production
go build -o main cmd/server/main.go
```

### Database Migrations
```bash
# GORM auto-migration is enabled
# Tables are created automatically on startup

# Manual migration (if needed)
docker compose exec backend go run cmd/migrate/main.go
```

### Frontend Development
```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build
```

## 🚀 Deployment

### Production Environment
1. **Set up PostgreSQL on EC2**
   ```bash
   # Install PostgreSQL
   sudo apt-get update
   sudo apt-get install postgresql-15
   
   # Configure for production
   sudo -u postgres psql
   CREATE DATABASE ktchat_prod;
   CREATE USER ktchat_prod WITH PASSWORD 'secure_password';
   GRANT ALL PRIVILEGES ON DATABASE ktchat_prod TO ktchat_prod;
   ```

2. **Configure environment variables**
   ```bash
   # Production .env
   DATABASE_URL=postgres://ktchat_prod:secure_password@localhost:5432/ktchat_prod?sslmode=disable
   ENVIRONMENT=production
   JWT_SECRET=your-secure-jwt-secret
   ```

3. **Deploy with Docker**
   ```bash
   docker compose -f docker-compose.prod.yml up -d
   ```

## 📚 Documentation

- [AWS Cognito Setup](docs/AWS_COGNITO_SETUP.md)
- [Docker Compose Guide](docs/DOCKER_COMPOSE_GUIDE.md)
- [Setup Guide](docs/SETUP_GUIDE.md)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is for learning and technical verification purposes.

## 🆘 Support

For issues and questions:
1. Check the documentation in the `docs/` folder
2. Review the logs: `docker compose logs`
3. Test individual components using the provided scripts

---

**Note**: This is a prototype for learning purposes. For production use, ensure proper security measures, error handling, and performance optimizations are implemented. 