#!/bin/bash

# Simple CSV Upload Test Script
# This script tests the CSV upload functionality to verify the fix

BASE_URL="http://localhost:8086"
API_KEY="oms-dev-key-2025"

echo "🧪 Testing CSV Upload Fix"
echo "=========================="

# Test CSV Upload
echo ""
echo "Testing CSV Upload with correct field name 'file'"
response=$(curl -s -X POST "$BASE_URL/api/v1/orders/upload" \
  -H "X-API-Key: $API_KEY" \
  -F "file=@test_working_inventory.csv")

echo "Response: $response"

# Check if the response contains success indicators
if echo "$response" | grep -q "CSV file uploaded and queued for processing successfully"; then
    echo "✅ CSV Upload Test PASSED!"
else
    echo "❌ CSV Upload Test FAILED!"
    echo "Expected success message not found in response"
fi

echo ""
echo "Test Complete!" 