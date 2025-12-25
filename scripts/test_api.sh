#!/bin/bash

# Test script for the authentication API
BASE_URL="http://localhost:8080"
EMAIL="test@example.com"
PASSWORD="testpassword123"

echo "Testing Authentication API..."
echo "================================"

# Test 1: Health check
echo "1. Health check..."
curl -s "$BASE_URL/health" | jq .
echo ""

# Test 2: Request verification code
echo "2. Requesting verification code for $EMAIL..."
RESPONSE=$(curl -s -X POST "$BASE_URL/auth/request-verification" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\"}")
echo "Response: $RESPONSE"
echo ""

# Test 3: Verify code (you'll need to get the actual code from your email)
read -p "Enter the 6-digit verification code from your email: " CODE
echo "3. Verifying code $CODE..."
RESPONSE=$(curl -s -X POST "$BASE_URL/auth/verify-code" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"code\":\"$CODE\"}")
echo "Response: $RESPONSE"
echo ""

# Test 4: Set password
echo "4. Setting password..."
RESPONSE=$(curl -s -X POST "$BASE_URL/auth/set-password" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")
echo "Response: $RESPONSE"

# Extract token from response
TOKEN=$(echo $RESPONSE | jq -r '.token // empty')
echo "Token: $TOKEN"
echo ""

# Test 5: Login
echo "5. Testing login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")
echo "Login Response: $LOGIN_RESPONSE"

# Extract token from login response
LOGIN_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token // empty')
echo ""

# Test 6: Access protected route
echo "6. Accessing protected route..."
if [ ! -z "$LOGIN_TOKEN" ] && [ "$LOGIN_TOKEN" != "null" ]; then
  PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/api/profile" \
    -H "Authorization: Bearer $LOGIN_TOKEN")
  echo "Profile Response: $PROFILE_RESPONSE"
else
  echo "No token available, skipping protected route test"
fi

echo ""
echo "Testing complete!"