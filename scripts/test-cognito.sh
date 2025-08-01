#!/bin/bash

# KT Chat - AWS Cognito Integration Test Script
# This script tests the complete authentication flow with AWS Cognito

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
API_BASE_URL="http://localhost:8080"
TEST_USERNAME="testuser-$(date +%s)"
TEST_EMAIL="test-$(date +%s)@example.com"
TEST_PASSWORD="TestPassword123!"

echo -e "${BLUE}Testing AWS Cognito Integration${NC}"
echo "=================================="
echo "API Base URL: $API_BASE_URL"
echo "Test Username: $TEST_USERNAME"
echo "Test Email: $TEST_EMAIL"
echo ""

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
        exit 1
    fi
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if required commands exist
if ! command_exists curl; then
    echo -e "${RED}Error: curl is not installed${NC}"
    exit 1
fi

if ! command_exists jq; then
    echo -e "${YELLOW}Warning: jq is not installed. JSON responses will not be formatted.${NC}"
    JQ_AVAILABLE=false
else
    JQ_AVAILABLE=true
fi

# 1. Check if server is running
echo -e "${BLUE}1. Checking if server is running...${NC}"
if curl -s "$API_BASE_URL/health" > /dev/null; then
    print_status 0 "Server is running"
else
    print_status 1 "Server is not running"
fi

# 2. Check Cognito configuration
echo -e "${BLUE}2. Checking Cognito configuration...${NC}"
HEALTH_RESPONSE=$(curl -s "$API_BASE_URL/health")
if echo "$HEALTH_RESPONSE" | grep -q "cognito.*enabled.*true"; then
    print_status 0 "Cognito is enabled"
    COGNITO_ENABLED=true
elif echo "$HEALTH_RESPONSE" | grep -q "cognito.*enabled.*false"; then
    print_status 1 "Cognito is not enabled"
    COGNITO_ENABLED=false
else
    echo -e "${YELLOW}Warning: Could not determine Cognito status${NC}"
    COGNITO_ENABLED=false
fi

# 3. Test registration
echo -e "${BLUE}3. Testing user registration...${NC}"
REGISTER_RESPONSE=$(curl -s -X POST "$API_BASE_URL/api/v1/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"username\": \"$TEST_USERNAME\",
        \"email\": \"$TEST_EMAIL\",
        \"password\": \"$TEST_PASSWORD\"
    }")

if [ "$JQ_AVAILABLE" = true ]; then
    echo "Registration Response:"
    echo "$REGISTER_RESPONSE" | jq .
else
    echo "Registration Response: $REGISTER_RESPONSE"
fi

# Check if registration was successful
if echo "$REGISTER_RESPONSE" | grep -q "success.*true" || echo "$REGISTER_RESPONSE" | grep -q "User registered successfully"; then
    print_status 0 "User registration successful"
elif echo "$REGISTER_RESPONSE" | grep -q "UsernameExistsException"; then
    print_status 0 "User already exists (expected for repeated tests)"
else
    print_status 1 "User registration failed"
    echo "Response: $REGISTER_RESPONSE"
fi

# 4. Test login (should fail if user is not confirmed)
echo -e "${BLUE}4. Testing login (should fail if user not confirmed)...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
        \"username\": \"$TEST_USERNAME\",
        \"password\": \"$TEST_PASSWORD\"
    }")

if [ "$JQ_AVAILABLE" = true ]; then
    echo "Login Response:"
    echo "$LOGIN_RESPONSE" | jq .
else
    echo "Login Response: $LOGIN_RESPONSE"
fi

# Check login response
if echo "$LOGIN_RESPONSE" | grep -q "USER_NOT_CONFIRMED"; then
    print_status 0 "Login correctly failed - user not confirmed"
    USER_CONFIRMED=false
elif echo "$LOGIN_RESPONSE" | grep -q "token"; then
    print_status 0 "Login successful - user is confirmed"
    USER_CONFIRMED=true
else
    print_status 1 "Login failed unexpectedly"
    echo "Response: $LOGIN_RESPONSE"
fi

# 5. Test resend confirmation code (if user not confirmed)
if [ "$USER_CONFIRMED" = false ]; then
    echo -e "${BLUE}5. Testing resend confirmation code...${NC}"
    RESEND_RESPONSE=$(curl -s -X POST "$API_BASE_URL/api/v1/auth/resend-confirmation" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"$TEST_USERNAME\"
        }")

    if [ "$JQ_AVAILABLE" = true ]; then
        echo "Resend Response:"
        echo "$RESEND_RESPONSE" | jq .
    else
        echo "Resend Response: $RESEND_RESPONSE"
    fi

    if echo "$RESEND_RESPONSE" | grep -q "success.*true" || echo "$RESEND_RESPONSE" | grep -q "Confirmation code sent"; then
        print_status 0 "Resend confirmation code successful"
    else
        print_status 1 "Resend confirmation code failed"
        echo "Response: $RESEND_RESPONSE"
    fi
fi

echo ""
echo -e "${GREEN}✅ Cognito integration test completed!${NC}"
echo ""
echo -e "${YELLOW}📝 Notes:${NC}"
echo "- If user registration succeeded, check your email for confirmation code"
echo "- If login failed with 'USER_NOT_CONFIRMED', use the confirmation code to confirm the account"
echo "- You can test the frontend confirmation flow at: http://localhost:3000/login"
echo ""
echo -e "${BLUE}🔗 Useful Commands:${NC}"
echo "View logs: docker compose logs backend"
echo "Check database: docker compose exec postgres psql -U ktchat -d ktchat -c '\\dt'"
echo "Test health: curl http://localhost:8080/health" 