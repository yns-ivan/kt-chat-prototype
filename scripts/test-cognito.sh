#!/bin/bash

# Test script for AWS Cognito integration
# This script tests the authentication endpoints

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
API_BASE_URL="http://localhost:8080"
TEST_USERNAME="testuser-$(date +%s)"
TEST_EMAIL="test-$(date +%s)@example.com"
TEST_PASSWORD="TestPassword123!"

echo -e "${YELLOW}Testing AWS Cognito Integration${NC}"
echo "=================================="
echo "API Base URL: $API_BASE_URL"
echo "Test Username: $TEST_USERNAME"
echo "Test Email: $TEST_EMAIL"
echo ""

# Function to check if server is running
check_server() {
    echo -e "${YELLOW}1. Checking if server is running...${NC}"
    if curl -s "$API_BASE_URL/health" > /dev/null; then
        echo -e "${GREEN}✓ Server is running${NC}"
        
        # Check if Cognito is enabled
        RESPONSE=$(curl -s "$API_BASE_URL/health")
        if echo "$RESPONSE" | grep -q '"cognito_enabled":true'; then
            echo -e "${GREEN}✓ Cognito is enabled${NC}"
        else
            echo -e "${RED}✗ Cognito is not enabled${NC}"
            echo "Make sure you have configured the Cognito environment variables"
            exit 1
        fi
    else
        echo -e "${RED}✗ Server is not running${NC}"
        echo "Please start the server first:"
        echo "  docker compose up -d"
        echo "  or"
        echo "  cd backend && go run cmd/server/main.go"
        exit 1
    fi
    echo ""
}

# Function to test registration
test_registration() {
    echo -e "${YELLOW}2. Testing user registration...${NC}"
    
    RESPONSE=$(curl -s -X POST "$API_BASE_URL/api/v1/auth/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"$TEST_USERNAME\",
            \"email\": \"$TEST_EMAIL\",
            \"password\": \"$TEST_PASSWORD\"
        }")
    
    if echo "$RESPONSE" | grep -q '"message":"User registered successfully'; then
        echo -e "${GREEN}✓ Registration successful${NC}"
        echo "Note: User needs to confirm email before login"
    else
        echo -e "${RED}✗ Registration failed${NC}"
        echo "Response: $RESPONSE"
        echo ""
        echo "This might be expected if:"
        echo "- Cognito is not properly configured"
        echo "- Email already exists"
        echo "- Password doesn't meet requirements"
    fi
    echo ""
}

# Function to test login (will fail if user not confirmed)
test_login() {
    echo -e "${YELLOW}3. Testing user login...${NC}"
    
    RESPONSE=$(curl -s -X POST "$API_BASE_URL/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"$TEST_USERNAME\",
            \"password\": \"$TEST_PASSWORD\"
        }")
    
    if echo "$RESPONSE" | grep -q '"token"'; then
        echo -e "${GREEN}✓ Login successful${NC}"
        
        # Extract token for further testing
        TOKEN=$(echo "$RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
        echo "Token received: ${TOKEN:0:20}..."
        
        # Test protected endpoint
        test_protected_endpoint "$TOKEN"
    else
        echo -e "${YELLOW}⚠ Login failed (expected if user not confirmed)${NC}"
        echo "Response: $RESPONSE"
        echo ""
        echo "This is expected if:"
        echo "- User hasn't confirmed their email"
        echo "- User doesn't exist in Cognito"
        echo "- Credentials are incorrect"
    fi
    echo ""
}

# Function to test protected endpoint
test_protected_endpoint() {
    local token=$1
    echo -e "${YELLOW}4. Testing protected endpoint...${NC}"
    
    RESPONSE=$(curl -s -X GET "$API_BASE_URL/api/v1/chat/rooms" \
        -H "Authorization: Bearer $token")
    
    if echo "$RESPONSE" | grep -q '"rooms"'; then
        echo -e "${GREEN}✓ Protected endpoint accessible${NC}"
    else
        echo -e "${RED}✗ Protected endpoint failed${NC}"
        echo "Response: $RESPONSE"
    fi
    echo ""
}

# Function to test token refresh
test_token_refresh() {
    echo -e "${YELLOW}5. Testing token refresh...${NC}"
    echo -e "${YELLOW}   (This requires a valid refresh token from login)${NC}"
    echo -e "${YELLOW}   Skipping for now...${NC}"
    echo ""
}

# Function to show next steps
show_next_steps() {
    echo -e "${YELLOW}Next Steps:${NC}"
    echo "=========="
    echo ""
    echo "1. Configure AWS Cognito User Pool:"
    echo "   - Follow the guide in docs/AWS_COGNITO_SETUP.md"
    echo ""
    echo "2. Set up environment variables:"
    echo "   - Copy env.example to .env in backend directory"
    echo "   - Fill in your Cognito configuration"
    echo ""
    echo "3. For testing with confirmed users:"
    echo "   - Use AWS Console to confirm test users"
    echo "   - Or configure email verification"
    echo ""
    echo "4. Test with the frontend:"
    echo "   - Start the frontend: cd frontend && npm run dev"
    echo "   - Try logging in through the web interface"
    echo ""
}

# Main execution
main() {
    check_server
    test_registration
    test_login
    test_token_refresh
    show_next_steps
    
    echo -e "${GREEN}Test completed!${NC}"
}

# Run main function
main 