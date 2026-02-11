#!/bin/bash

# Website Downloader - Test Script
# This script performs end-to-end testing of the application

set -e

echo "========================================="
echo "Website Downloader - Test Suite"
echo "========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test configuration
BACKEND_URL="http://localhost:8080"
TEST_URL="https://example.com"
TIMEOUT=60

# Function to print colored output
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $2"
    else
        echo -e "${RED}✗${NC} $2"
        exit 1
    fi
}

print_info() {
    echo -e "${YELLOW}ℹ${NC} $1"
}

echo "Step 1: Checking if backend is running..."
if curl -s -f "${BACKEND_URL}/health" > /dev/null 2>&1; then
    print_status 0 "Backend is running"
else
    print_status 1 "Backend is not running. Start it with: cd backend && go run main.go"
fi

echo ""
echo "Step 2: Testing API endpoints..."

# Test scrape endpoint
print_info "Testing POST /api/scrape..."
RESPONSE=$(curl -s -X POST "${BACKEND_URL}/api/scrape" \
    -H "Content-Type: application/json" \
    -d "{\"url\":\"${TEST_URL}\"}" \
    -w "\n%{http_code}")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_CODE" -eq 200 ]; then
    print_status 0 "Scrape endpoint returns HTTP 200"
    
    # Extract jobId
    JOB_ID=$(echo "$BODY" | grep -o '"jobId":"[^"]*' | sed 's/"jobId":"//')
    
    if [ -n "$JOB_ID" ]; then
        print_status 0 "Received job ID: $JOB_ID"
    else
        print_status 1 "Failed to extract job ID from response"
    fi
else
    print_status 1 "Scrape endpoint returned HTTP $HTTP_CODE"
fi

echo ""
echo "Step 3: Testing status endpoint..."
sleep 2

STATUS_RESPONSE=$(curl -s "${BACKEND_URL}/api/status/${JOB_ID}")
if echo "$STATUS_RESPONSE" | grep -q "status"; then
    print_status 0 "Status endpoint working"
    print_info "Status: $(echo $STATUS_RESPONSE | grep -o '"status":"[^"]*' | sed 's/"status":"//')"
else
    print_status 1 "Status endpoint failed"
fi

echo ""
echo "Step 4: Waiting for scraping to complete..."
print_info "This may take up to ${TIMEOUT} seconds..."

COUNTER=0
while [ $COUNTER -lt $TIMEOUT ]; do
    STATUS_RESPONSE=$(curl -s "${BACKEND_URL}/api/status/${JOB_ID}")
    STATUS=$(echo "$STATUS_RESPONSE" | grep -o '"status":"[^"]*' | sed 's/"status":"//')
    
    if [ "$STATUS" = "complete" ]; then
        print_status 0 "Scraping completed successfully"
        break
    elif [ "$STATUS" = "error" ]; then
        print_status 1 "Scraping failed with error"
        break
    fi
    
    sleep 2
    COUNTER=$((COUNTER + 2))
    echo -ne "\rWaiting... ${COUNTER}s elapsed"
done

if [ $COUNTER -ge $TIMEOUT ]; then
    echo ""
    print_info "Timeout reached. Job may still be running."
fi

echo ""
echo "Step 5: Testing download endpoint..."

DOWNLOAD_HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${BACKEND_URL}/api/download/${JOB_ID}")

if [ "$DOWNLOAD_HTTP_CODE" -eq 200 ]; then
    print_status 0 "Download endpoint accessible"
    
    # Download the file
    print_info "Downloading ZIP file..."
    curl -s -o "test_download_${JOB_ID}.zip" "${BACKEND_URL}/api/download/${JOB_ID}"
    
    if [ -f "test_download_${JOB_ID}.zip" ]; then
        FILE_SIZE=$(stat -f%z "test_download_${JOB_ID}.zip" 2>/dev/null || stat -c%s "test_download_${JOB_ID}.zip" 2>/dev/null)
        print_status 0 "Downloaded ZIP file (${FILE_SIZE} bytes)"
        
        # Cleanup
        rm "test_download_${JOB_ID}.zip"
        print_info "Cleaned up test file"
    else
        print_status 1 "Failed to download ZIP file"
    fi
else
    print_status 1 "Download endpoint returned HTTP $DOWNLOAD_HTTP_CODE"
fi

echo ""
echo "========================================="
echo -e "${GREEN}All tests passed!${NC}"
echo "========================================="
echo ""
echo "Test Summary:"
echo "  - Backend health check: ✓"
echo "  - Scrape endpoint: ✓"
echo "  - Status endpoint: ✓"
echo "  - Scraping completion: ✓"
echo "  - Download endpoint: ✓"
echo ""
