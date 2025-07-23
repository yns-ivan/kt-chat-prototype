# KT Chat Backend

A Go-based backend for the KT Chat prototype system with real-time messaging, file uploads, and encryption.

## Features

- **Real-time Chat**: WebSocket-based real-time messaging
- **File Upload**: Support for images, PDFs, and videos with preview
- **Message Encryption**: AES-256-GCM encryption for chat messages
- **AWS Cognito Integration**: User authentication via AWS Cognito
- **MySQL Database**: Persistent storage with GORM
- **RESTful API**: Clean API design with Gin framework

## Prerequisites

- Go 1.24.5 or later
- MySQL 8.0 or later
- AWS Cognito User Pool (for production)

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── auth/               # Authentication logic
│   ├── chat/               # Chat service and handlers
│   ├── database/           # Database initialization and migrations
│   ├── encryption/         # Message encryption/decryption
│   ├── file/               # File upload and management
│   ├── models/             # Database models
│   └── websocket/          # WebSocket hub and client management
├── pkg/
│   ├── config/             # Configuration management
│   ├── middleware/         # HTTP middleware
│   └── utils/              # Utility functions
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies checksum
├── env.example             # Environment variables example
└── README.md               # This file
```

## Setup

1. **Clone the repository and navigate to backend directory**
   ```bash
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

4. **Set up MySQL database**
   ```sql
   CREATE DATABASE ktchat CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   CREATE USER 'ktchat'@'localhost' IDENTIFIED BY 'password';
   GRANT ALL PRIVILEGES ON ktchat.* TO 'ktchat'@'localhost';
   FLUSH PRIVILEGES;
   ```

5. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENVIRONMENT` | Application environment | `development` |
| `DATABASE_URL` | MySQL connection string | `ktchat:password@tcp(localhost:3306)/ktchat?charset=utf8mb4&parseTime=True&loc=Local` |
| `JWT_SECRET` | JWT signing secret | `your-secret-key-change-in-production` |
| `AWS_REGION` | AWS region for Cognito | `ap-northeast-1` |
| `COGNITO_USER_POOL_ID` | AWS Cognito User Pool ID | - |
| `COGNITO_CLIENT_ID` | AWS Cognito Client ID | - |
| `COGNITO_CLIENT_SECRET` | AWS Cognito Client Secret | - |
| `UPLOAD_PATH` | File upload directory | `./uploads` |
| `MAX_FILE_SIZE` | Maximum file size in bytes | `52428800` (50MB) |
| `ENCRYPTION_KEY` | Encryption key (32 bytes) | `your-encryption-key-32-bytes-long` |
| `PORT` | Server port | `8080` |

## API Endpoints

### Authentication

- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/refresh` - Token refresh

### Chat Rooms

- `GET /api/v1/chat/rooms` - Get all chat rooms
- `POST /api/v1/chat/rooms` - Create a new chat room
- `GET /api/v1/chat/rooms/:roomID/messages` - Get messages for a room
- `POST /api/v1/chat/rooms/:roomID/join` - Join a chat room
- `POST /api/v1/chat/rooms/:roomID/leave` - Leave a chat room

### WebSocket

- `GET /api/v1/ws` - WebSocket connection for real-time chat

### Health Check

- `GET /health` - Health check endpoint

## WebSocket Protocol

### Connection
Connect to `/api/v1/ws?user_id=<user_id>&username=<username>&room_id=<room_id>`

### Message Format
```json
{
  "type": "message",
  "room_id": "room-123",
  "user_id": "user-123",
  "username": "john_doe",
  "content": "Hello, world!",
  "timestamp": "2024-01-01T12:00:00Z",
  "files": [
    {
      "id": "file-123",
      "file_name": "image.jpg",
      "file_type": "image",
      "file_size": 1024
    }
  ]
}
```

## Encryption

Messages are encrypted using AES-256-GCM before being stored in the database. The encryption key is configurable via the `ENCRYPTION_KEY` environment variable.

### Searchable Encryption

For message search functionality, a simplified searchable encryption scheme is implemented. In production, consider using more sophisticated approaches like:

- Encrypted Bloom filters
- Searchable symmetric encryption (SSE)
- Homomorphic encryption

## File Upload

Supported file types:
- **Images**: JPG, JPEG, PNG, GIF, BMP, WebP
- **Documents**: PDF
- **Videos**: MP4, AVI, MOV, WMV, FLV, WebM
- **Audio**: MP3, WAV, OGG, AAC

Files are stored in the configured upload directory with unique filenames.

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o bin/server cmd/server/main.go
```

### Docker Support
```bash
docker build -t ktchat-backend .
docker run -p 8080:8080 ktchat-backend
```

## Deployment

### EC2 Deployment
The application is configured to run on EC2 instance `DEV-KTCHAT-WEB01` (54.179.141.95) with:
- Nginx reverse proxy on port 80
- Go application on port 8080
- Basic auth: `ktchat` / `s9RnNyai`

### Environment Setup
1. Set `ENVIRONMENT=production`
2. Configure production database URL
3. Set secure JWT and encryption keys
4. Configure AWS Cognito credentials

## Performance Considerations

- Database connection pooling is configured
- Rate limiting middleware is implemented
- WebSocket connections are managed efficiently
- File uploads are validated and size-limited

## Security

- JWT-based authentication
- Message encryption at rest
- CORS configuration
- Input validation and sanitization
- File type validation

## Monitoring

- Health check endpoint available
- Request logging middleware
- Error handling and logging

## Contributing

1. Follow Go coding standards
2. Add tests for new features
3. Update documentation
4. Use conventional commit messages 