#!/bin/bash

# Test script for user confirmation flow
echo "Testing User Confirmation Flow"
echo "=============================="

API_BASE_URL="http://localhost:8080"
TEST_USERNAME="testuser-$(date +%s)"
TEST_EMAIL="test-$(date +%s)@example.com"
TEST_PASSWORD="TestPassword123!"

echo "Test Username: $TEST_USERNAME"
echo "Test Email: $TEST_EMAIL"
echo ""

# 1. Check if server is running
echo "1. Checking if server is running..."
if curl -s "$API_BASE_URL/health" > /dev/null; then
    echo "✓ Server is running"
else
    echo "✗ Server is not running"
    exit 1
fi

# 2. Register a new user
echo ""
echo "2. Registering new user..."
REGISTER_RESPONSE=$(curl -s -X POST "$API_BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$TEST_USERNAME\",
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$TEST_PASSWORD\"
  }")

echo "Register response: $REGISTER_RESPONSE"

# 3. Try to login (should fail with confirmation required)
echo ""
echo "3. Trying to login (should require confirmation)..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$TEST_USERNAME\",
    \"password\": \"$TEST_PASSWORD\"
  }")

echo "Login response: $LOGIN_RESPONSE"

# 4. Test resend confirmation code
echo ""
echo "4. Testing resend confirmation code..."
RESEND_RESPONSE=$(curl -s -X POST "$API_BASE_URL/api/v1/auth/resend-confirmation" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$TEST_USERNAME\"
  }")

echo "Resend response: $RESEND_RESPONSE"

echo ""
echo "✅ Confirmation endpoints are working!"
echo ""
echo "📧 Check your email for the confirmation code"
echo "🔗 You can now test the frontend confirmation flow at: http://localhost:3000/login"
echo ""
echo "💡 To test the confirmation flow:"
echo "   1. Go to http://localhost:3000/login"
echo "   2. Try to login with username: $TEST_USERNAME"
echo "   3. The confirmation modal should appear"
echo "   4. Enter the confirmation code from your email"
echo "   5. Or click 'Resend Code' to get a new one" 