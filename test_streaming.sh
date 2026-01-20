#!/bin/bash

# Simple API Test for Customer Excel Streaming
BASE_URL="http://localhost:3000"  # Assuming your server runs on port 3000

echo "🚀 Testing Customer Excel Streaming Endpoints"
echo "=============================================="
echo ""

# Test 1: POST with JSON body (Small dataset)
echo "📊 Test 1: POST - Small Dataset (1000 records)"
echo "-----------------------------------------------"
curl -X POST "${BASE_URL}/customers/export/excel/streaming" \
  -H "Content-Type: application/json" \
  -d '{"total_data": 1000, "filename": "test_small.xlsx"}' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

sleep 2

# Test 2: GET with query parameters (Medium dataset)
echo "📊 Test 2: GET - Medium Dataset (5000 records)"
echo "----------------------------------------------"
curl "${BASE_URL}/customers/export/excel/streaming?total_data=5000&filename=test_medium.xlsx" \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

sleep 2

# Test 3: Error handling - Invalid data
echo "❌ Test 3: Error Handling - Invalid Data"
echo "----------------------------------------"
curl -X POST "${BASE_URL}/customers/export/excel/streaming" \
  -H "Content-Type: application/json" \
  -d '{"total_data": -1}' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

echo "✅ Testing completed!"
echo ""
echo "💡 To test file download:"
echo "   curl -O '${BASE_URL}/customers/download/test_small.xlsx'"
echo ""
echo "🔧 To start your server:"
echo "   go run main.go"
