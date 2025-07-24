# AWS Cognito Setup Guide

This guide will help you set up AWS Cognito for authentication in your KT Chat application.

## Prerequisites

1. AWS Account with appropriate permissions
2. AWS CLI configured (optional but recommended)
3. Basic understanding of AWS Cognito concepts

## Step 1: Create a User Pool

### Using AWS Console

1. **Navigate to Cognito**
   - Go to AWS Console → Amazon Cognito
   - Click "Create user pool"

2. **Configure sign-in experience**
   - **Cognito user pool sign-in options**: Choose "Username"
   - **Cognito-assisted verification and confirmation**: Choose "Cognito-assisted verification and confirmation"
   - Click "Next"

3. **Configure security requirements**
   - **Password policy**: Choose "Cognito defaults" or customize
   - **Multi-factor authentication**: Choose "No MFA" for development (enable for production)
   - **User account recovery**: Choose "Self-service recovery"
   - Click "Next"

4. **Configure sign-up experience**
   - **Self-service sign-up**: Enable
   - **Cognito-assisted verification and confirmation**: Enable
   - **Verification message**: Choose "Email"
   - **Required attributes**: Select "email"
   - Click "Next"

5. **Configure message delivery**
   - **Email provider**: Choose "Send email with Cognito"
   - **From email address**: Use Cognito default or configure SES
   - Click "Next"

6. **Integrate your app**
   - **User pool name**: Enter "ktchat-user-pool" (or your preferred name)
   - **Initial app client**: Choose "Public client"
   - **App client name**: Enter "ktchat-client"
   - **Client secret**: **Important**: Choose "Generate client secret"
   - Click "Next"

7. **Review and create**
   - Review your settings
   - Click "Create user pool"

### Using AWS CLI

```bash
# Create user pool
aws cognito-idp create-user-pool \
  --pool-name "ktchat-user-pool" \
  --policies '{
    "PasswordPolicy": {
      "MinimumLength": 8,
      "RequireUppercase": true,
      "RequireLowercase": true,
      "RequireNumbers": true,
      "RequireSymbols": false
    }
  }' \
  --auto-verified-attributes email \
  --username-attributes email \
  --schema '[
    {
      "Name": "email",
      "AttributeDataType": "String",
      "Required": true,
      "Mutable": true
    }
  ]'

# Create app client
aws cognito-idp create-user-pool-client \
  --user-pool-id YOUR_USER_POOL_ID \
  --client-name "ktchat-client" \
  --generate-secret \
  --explicit-auth-flows ADMIN_NO_SRP_AUTH USER_PASSWORD_AUTH REFRESH_TOKEN_AUTH \
  --supported-identity-providers COGNITO
```

## Step 2: Configure App Client Settings

1. **In your User Pool, go to "App integration" tab**
2. **Under "App clients and analytics", click on your app client**
3. **Configure the following settings**:

   - **App client name**: ktchat-client
   - **Client secret**: Generate client secret (important!)
   - **Authentication flows**:
     - ✅ ADMIN_NO_SRP_AUTH
     - ✅ USER_PASSWORD_AUTH
     - ✅ REFRESH_TOKEN_AUTH
   - **OAuth 2.0 settings**:
     - **Allowed OAuth flows**: Authorization code grant, Implicit grant
     - **Allowed OAuth scopes**: openid, email, profile
     - **Callback URLs**: `http://localhost:3000/callback` (for development)
     - **Sign-out URLs**: `http://localhost:3000/logout` (for development)

## Step 3: Configure Environment Variables

### For Local Development

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
COGNITO_USER_POOL_ID=ap-northeast-1_xxxxxxxxx
COGNITO_CLIENT_ID=xxxxxxxxxxxxxxxxxxxxxxxxxx
COGNITO_CLIENT_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

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

### For Docker Development

Create a `.env` file in the root directory:

```env
# AWS Cognito Configuration
COGNITO_USER_POOL_ID=ap-northeast-1_xxxxxxxxx
COGNITO_CLIENT_ID=xxxxxxxxxxxxxxxxxxxxxxxxxx
COGNITO_CLIENT_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# AWS Credentials
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```

## Step 4: Configure AWS Credentials

### Option 1: Environment Variables (Recommended for Docker)

Set the environment variables as shown above.

### Option 2: AWS CLI Configuration

```bash
aws configure
```

Enter your AWS credentials when prompted.

### Option 3: IAM Role (Recommended for Production)

For production deployments, use IAM roles instead of access keys.

## Step 5: Test the Setup

### 1. Start the Application

```bash
# Using Docker Compose
docker compose up -d

# Or locally
cd backend
go run cmd/server/main.go
```

### 2. Check Health Endpoint

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "service": "ktchat-backend",
  "version": "1.0.0",
  "cognito_enabled": true
}
```

### 3. Test Registration

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "TestPassword123!"
  }'
```

### 4. Test Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "TestPassword123!"
  }'
```

## Step 6: Frontend Integration

The frontend is already configured to work with the new authentication flow. The login response will now include:

```json
{
  "token": "custom-jwt-token",
  "cognito_tokens": {
    "access_token": "cognito-access-token",
    "refresh_token": "cognito-refresh-token",
    "id_token": "cognito-id-token",
    "expires_in": 3600,
    "token_type": "Bearer"
  },
  "user": {
    "id": "user-uuid",
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

## Troubleshooting

### Common Issues

1. **"Invalid credentials" error**
   - Check that the user exists in Cognito
   - Verify the user has confirmed their email
   - Ensure the password meets the pool requirements

2. **"Failed to initialize Cognito service"**
   - Verify AWS credentials are correct
   - Check that the region matches your user pool
   - Ensure the user pool ID and client ID are correct

3. **"Invalid token" error**
   - Check that the JWT secret is consistent
   - Verify token expiration
   - Ensure the token format is correct

4. **"Registration failed" error**
   - Check that self-service sign-up is enabled
   - Verify the email format is valid
   - Ensure the password meets requirements

### Debug Mode

Enable debug logging by setting:

```env
ENVIRONMENT=development
```

The application will log detailed information about authentication attempts.

## Security Considerations

1. **Client Secret**: Always generate and use a client secret for web applications
2. **HTTPS**: Use HTTPS in production
3. **Token Storage**: Store tokens securely (httpOnly cookies for web apps)
4. **Password Policy**: Enforce strong password requirements
5. **MFA**: Enable multi-factor authentication for production
6. **Rate Limiting**: Implement rate limiting on authentication endpoints

## Production Deployment

For production deployment:

1. **Use IAM Roles** instead of access keys
2. **Enable MFA** for user accounts
3. **Configure custom domain** for Cognito
4. **Set up proper CORS** origins
5. **Use HTTPS** for all endpoints
6. **Implement proper error handling**
7. **Set up monitoring and logging**

## Next Steps

1. **Email Verification**: Configure SES for custom email templates
2. **Social Login**: Add social identity providers (Google, Facebook, etc.)
3. **Advanced Security**: Implement advanced security features
4. **User Management**: Build admin interfaces for user management
5. **Analytics**: Set up Cognito analytics for user behavior insights 