# Docker Compose Guide

This guide explains how to use Docker Compose to run the KT Chat prototype in a containerized environment.

## 🐳 Overview

The KT Chat prototype uses Docker Compose to orchestrate multiple services:

- **PostgreSQL 15**: Database for storing chat data and user information
- **Go Backend**: REST API and WebSocket server
- **Nginx**: Reverse proxy for production-like setup
- **Frontend**: Nuxt.js SPA (planned)
- **Laravel**: Admin panel (planned)

## 🚀 Quick Start

### Prerequisites
- Docker Desktop installed
- Docker Compose plugin enabled

### 1. Clone and Navigate
```bash
git clone <repository-url>
cd chat-prototype
```

### 2. Configure Environment
```bash
# Copy environment file
cp .env.example .env

# Edit with your configuration
nano .env
```

### 3. Start Services
```bash
# Start all services
docker compose up -d

# Check status
docker compose ps
```

### 4. Verify Setup
```bash
# Test backend
curl http://localhost:8080/health

# Test database
docker compose exec postgres psql -U ktchat -d ktchat -c "\dt"
```

## 📋 Service Configuration

### PostgreSQL Database
```yaml
postgres:
  image: postgres:15-alpine
  container_name: ktchat-postgres
  environment:
    POSTGRES_DB: ktchat
    POSTGRES_USER: ktchat
    POSTGRES_PASSWORD: password
    POSTGRES_INITDB_ARGS: "--encoding=UTF-8 --lc-collate=C --lc-ctype=C"
  ports:
    - "5432:5432"
  volumes:
    - postgres_data:/var/lib/postgresql/data
    - ./scripts/init-postgres.sql:/docker-entrypoint-initdb.d/init.sql
```

**Key Features:**
- Uses PostgreSQL 15 Alpine for smaller image size
- UTF-8 encoding for international character support
- Persistent data storage with named volume
- Automatic initialization script execution

### Go Backend
```yaml
backend:
  build:
    context: ./backend
    dockerfile: Dockerfile
  environment:
    - DATABASE_URL=postgres://ktchat:password@postgres:5432/ktchat?sslmode=disable
    - JWT_SECRET=dev-jwt-secret-key-change-in-production
    - AWS_REGION=ap-northeast-1
    - COGNITO_USER_POOL_ID=${COGNITO_USER_POOL_ID:-}
    - COGNITO_CLIENT_ID=${COGNITO_CLIENT_ID:-}
    - COGNITO_CLIENT_SECRET=${COGNITO_CLIENT_SECRET:-}
    - AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID:-}
    - AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY:-}
  ports:
    - "8080:8080"
  volumes:
    - ./backend/uploads:/root/uploads
  depends_on:
    - postgres
```

**Key Features:**
- Multi-stage build for optimized image size
- Environment variable configuration
- File upload volume mounting
- Service dependency management

### Nginx Reverse Proxy
```yaml
nginx:
  image: nginx:alpine
  container_name: ktchat-nginx
  ports:
    - "80:80"
    - "443:443"
  volumes:
    - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
    - ./nginx/conf.d:/etc/nginx/conf.d:ro
  depends_on:
    - backend
    - frontend
    - laravel
```

**Key Features:**
- Lightweight Alpine-based image
- Configuration file mounting
- SSL support ready
- Load balancing capabilities

## 🔧 Environment Variables

### Required Variables
```env
# Database
DATABASE_URL=postgres://ktchat:password@postgres:5432/ktchat?sslmode=disable

# JWT
JWT_SECRET=dev-jwt-secret-key-change-in-production

# AWS Cognito (Optional)
AWS_REGION=ap-northeast-1
COGNITO_USER_POOL_ID=your-user-pool-id
COGNITO_CLIENT_ID=your-client-id
COGNITO_CLIENT_SECRET=your-client-secret
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```

### Optional Variables
```env
# File Upload
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=52428800

# Encryption
ENCRYPTION_KEY=dev-encryption-key-32-bytes-long

# Server
PORT=8080
ENVIRONMENT=development
```

## 🛠️ Docker Commands

### Basic Operations
```bash
# Start all services
docker compose up -d

# Stop all services
docker compose down

# View logs
docker compose logs -f

# Rebuild and restart
docker compose up -d --build

# Check service status
docker compose ps
```

### Service-Specific Commands
```bash
# View backend logs
docker compose logs backend

# View database logs
docker compose logs postgres

# Execute commands in containers
docker compose exec backend sh
docker compose exec postgres psql -U ktchat -d ktchat

# Restart specific service
docker compose restart backend
```

### Database Operations
```bash
# Connect to PostgreSQL
docker compose exec postgres psql -U ktchat -d ktchat

# Backup database
docker compose exec postgres pg_dump -U ktchat ktchat > backup.sql

# Restore database
docker compose exec -T postgres psql -U ktchat -d ktchat < backup.sql

# Reset database (WARNING: Deletes all data)
docker compose down -v
docker compose up -d
```

## 📊 Monitoring and Debugging

### Health Checks
```bash
# Backend health
curl http://localhost:8080/health

# Database connection
docker compose exec postgres pg_isready -U ktchat

# Service status
docker compose ps
```

### Log Analysis
```bash
# Follow all logs
docker compose logs -f

# Follow specific service
docker compose logs -f backend

# Search logs
docker compose logs | grep "ERROR"

# Export logs
docker compose logs > logs.txt
```

### Resource Usage
```bash
# Container resource usage
docker stats

# Disk usage
docker system df

# Clean up unused resources
docker system prune
```

## 🔒 Security Considerations

### Development Environment
- Default passwords are used for development
- SSL is disabled for local development
- All services are exposed on localhost

### Production Recommendations
```yaml
# Production docker-compose.yml example
version: '3.8'
services:
  postgres:
    environment:
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}  # Use secrets
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - internal  # Internal network only
  
  backend:
    environment:
      - DATABASE_URL=${DATABASE_URL}
      - JWT_SECRET=${JWT_SECRET}
    networks:
      - internal
      - external
  
  nginx:
    networks:
      - external
    volumes:
      - ./ssl:/etc/nginx/ssl:ro  # SSL certificates

networks:
  internal:
    internal: true
  external:
    driver: bridge

volumes:
  postgres_data:
    driver: local
```

## 🐛 Troubleshooting

### Common Issues

**1. Port Conflicts**
```bash
# Check what's using the ports
lsof -i :8080
lsof -i :5432
lsof -i :80

# Change ports in docker-compose.yml if needed
ports:
  - "8081:8080"  # Use different host port
```

**2. Database Connection Issues**
```bash
# Check if PostgreSQL is running
docker compose ps postgres

# Check database logs
docker compose logs postgres

# Test connection manually
docker compose exec postgres psql -U ktchat -d ktchat -c "SELECT 1;"
```

**3. Backend Won't Start**
```bash
# Check backend logs
docker compose logs backend

# Check environment variables
docker compose exec backend env | grep DATABASE

# Verify database URL format
# Should be: postgres://ktchat:password@postgres:5432/ktchat?sslmode=disable
```

**4. Permission Issues**
```bash
# Fix file permissions
sudo chown -R $USER:$USER .

# Fix Docker volume permissions
docker compose down
sudo chown -R 1001:1001 ./backend/uploads
docker compose up -d
```

### Reset Everything
```bash
# Stop and remove everything
docker compose down -v

# Remove all images
docker rmi $(docker images -q)

# Clean up volumes
docker volume prune

# Rebuild from scratch
docker compose up -d --build
```

## 📈 Performance Optimization

### Resource Limits
```yaml
services:
  postgres:
    deploy:
      resources:
        limits:
          memory: 1G
        reservations:
          memory: 512M
  
  backend:
    deploy:
      resources:
        limits:
          memory: 512M
        reservations:
          memory: 256M
```

### Database Optimization
```yaml
postgres:
  environment:
    POSTGRES_INITDB_ARGS: "--encoding=UTF-8 --lc-collate=C --lc-ctype=C"
    POSTGRES_SHARED_BUFFERS: 256MB
    POSTGRES_EFFECTIVE_CACHE_SIZE: 1GB
    POSTGRES_WORK_MEM: 4MB
    POSTGRES_MAINTENANCE_WORK_MEM: 64MB
```

### Caching
```yaml
nginx:
  volumes:
    - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
    - nginx_cache:/var/cache/nginx  # Add caching volume

volumes:
  nginx_cache:
```

## 🔄 CI/CD Integration

### GitHub Actions Example
```yaml
name: Deploy to Production
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Deploy to server
        run: |
          ssh user@server "cd /path/to/app && git pull && docker compose up -d --build"
```

## 📚 Additional Resources

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [PostgreSQL Docker Image](https://hub.docker.com/_/postgres)
- [Nginx Docker Image](https://hub.docker.com/_/nginx)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)

---

**Note**: This guide covers development setup. For production deployment, additional security measures, monitoring, and backup strategies should be implemented. 