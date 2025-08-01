# KT Chat Setup Guide

This guide provides comprehensive instructions for setting up the KT Chat prototype on your local development environment.

## 📋 Prerequisites

### Required Software
- **Docker & Docker Compose** (recommended)
- **Go 1.24.5+** (for local development)
- **PostgreSQL 15+** (for local development without Docker)
- **Node.js 22.12.0** (for frontend development)

### Optional Software
- **AWS CLI** (for AWS Cognito setup)
- **Postman** (for API testing)

## 🚀 Quick Setup with Docker (Recommended)

### 1. Clone the Repository
```bash
git clone <repository-url>
cd chat-prototype
```

### 2. Configure Environment Variables
```bash
# Copy the example environment file
cp .env.example .env

# Edit the .env file with your configuration
nano .env
```

**Example .env file:**
```env
# AWS Cognito Configuration (Optional for development)
AWS_REGION=ap-northeast-1
COGNITO_USER_POOL_ID=your-user-pool-id
COGNITO_CLIENT_ID=your-client-id
COGNITO_CLIENT_SECRET=your-client-secret

# AWS Credentials (for local development)
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key

# S3 Configuration (Optional - for file storage)
S3_BUCKET_NAME=your-s3-bucket-name
STORAGE_TYPE=local  # "local" or "s3"
```

### 3. Start the Development Environment
```bash
# Start all services
docker compose up -d

# Check service status
docker compose ps
```

### 4. Verify the Setup
```bash
# Test backend health
curl http://localhost:8080/health

# Test PostgreSQL connection
docker compose exec postgres psql -U ktchat -d ktchat -c "\dt"

# Test AWS Cognito integration (if configured)
./scripts/test-cognito.sh

# Test S3 integration (if configured)
./scripts/setup-s3.sh
```

## 📁 S3 File Storage Setup (Optional)

### Quick S3 Setup
```bash
# Run the automated S3 setup script
./scripts/setup-s3.sh
```

### Manual S3 Setup
1. **Create S3 Bucket**
   ```bash
   aws s3 mb s3://your-bucket-name --region ap-northeast-1
   ```

2. **Configure Bucket for Private Access**
   ```bash
   aws s3api put-public-access-block \
       --bucket your-bucket-name \
       --public-access-block-configuration \
       BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true
   ```

3. **Create IAM Policy**
   - Go to AWS IAM Console
   - Create policy with S3 permissions
   - Attach to your EC2 instance role or create IAM user

4. **Update Environment Variables**
   ```bash
   # Add to your .env file
   STORAGE_TYPE=s3
   S3_BUCKET_NAME=your-bucket-name
   ```

For detailed S3 setup instructions, see [S3_SETUP.md](S3_SETUP.md).

## 🔧 Manual Setup (Alternative)

### 1. Install PostgreSQL 15

**macOS:**
```bash
brew install postgresql@15
brew services start postgresql@15
```

**Ubuntu/Debian:**
```bash
sudo apt-get update
sudo apt-get install postgresql-15 postgresql-contrib-15
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

**Windows:**
Download and install from [PostgreSQL official website](https://www.postgresql.org/download/windows/)

### 2. Set Up Database
```bash
# Connect to PostgreSQL as superuser
sudo -u postgres psql

# Create database and user
CREATE DATABASE ktchat;
CREATE USER ktchat WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE ktchat TO ktchat;
GRANT ALL PRIVILEGES ON SCHEMA public TO ktchat;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO ktchat;
\q
```

### 3. Configure Backend
```bash
cd backend

# Copy environment file
cp env.example .env

# Edit .env file
nano .env
```

**Backend .env configuration:**
```env
# Environment
ENVIRONMENT=development

# Database (PostgreSQL)
DATABASE_URL=postgres://ktchat:password@localhost:5432/ktchat?sslmode=disable

# JWT
JWT_SECRET=dev-jwt-secret-key-change-in-production

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
ENCRYPTION_KEY=dev-encryption-key-32-bytes-long

# Server
PORT=8080
```

### 4. Install Dependencies and Run
```bash
# Install Go dependencies
go mod download

# Run the backend
go run cmd/server/main.go
```

## 🔐 AWS Cognito Setup (Optional)

For production-like authentication, follow the comprehensive AWS Cognito setup guide:

📖 **[AWS Cognito Setup Guide](AWS_COGNITO_SETUP.md)**

### Quick Cognito Setup
1. Create a User Pool in AWS Cognito Console
2. Create an App Client
3. Configure the environment variables
4. Test the integration

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

# Create a chat room
curl -X POST http://localhost:8080/api/v1/chat/rooms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name":"Test Room","description":"A test room"}'
```

## 📊 Database Management

### PostgreSQL Commands
```bash
# Connect to database
docker compose exec postgres psql -U ktchat -d ktchat

# List tables
\dt

# View table structure
\d users

# Check data
SELECT * FROM users LIMIT 5;

# Exit
\q
```

### Database Migrations
The application uses GORM auto-migration, which automatically creates and updates tables based on the model definitions.

```bash
# Tables are created automatically on startup
# Manual migration (if needed)
docker compose exec backend go run cmd/migrate/main.go
```

## 🔧 Development Workflow

### Backend Development
```bash
cd backend

# Install air for hot reload (optional)
go install github.com/cosmtrek/air@latest

# Run with hot reload
air

# Run tests
go test ./...

# Format code
go fmt ./...

# Build for production
go build -o main cmd/server/main.go
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

## 🐛 Troubleshooting

### Common Issues

**1. Database Connection Error**
```bash
# Check if PostgreSQL is running
docker compose ps

# Check logs
docker compose logs postgres

# Test connection manually
docker compose exec postgres psql -U ktchat -d ktchat -c "SELECT 1;"
```

**2. Backend Won't Start**
```bash
# Check logs
docker compose logs backend

# Check environment variables
docker compose exec backend env | grep DATABASE

# Verify database URL format
# Should be: postgres://ktchat:password@postgres:5432/ktchat?sslmode=disable
```

**3. AWS Cognito Issues**
```bash
# Check AWS credentials
aws sts get-caller-identity

# Verify environment variables
docker compose exec backend env | grep AWS

# Test Cognito connection
./scripts/test-cognito.sh
```

**4. Port Conflicts**
```bash
# Check what's using the ports
lsof -i :8080
lsof -i :5432
lsof -i :3000

# Stop conflicting services or change ports in docker-compose.yml
```

### Reset Everything
```bash
# Stop all services
docker compose down

# Remove volumes (WARNING: This deletes all data)
docker compose down -v

# Rebuild and start
docker compose up -d --build
```

## 📚 Additional Resources

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [AWS Cognito Documentation](https://docs.aws.amazon.com/cognito/)

## 🆘 Getting Help

If you encounter issues:

1. **Check the logs**: `docker compose logs`
2. **Review this guide**: Ensure all steps were followed correctly
3. **Test individual components**: Use the provided test scripts
4. **Check the documentation**: Review the docs/ folder
5. **Verify your environment**: Ensure all prerequisites are met

---

**Note**: This setup guide is for development purposes. For production deployment, additional security measures and configurations are required. 