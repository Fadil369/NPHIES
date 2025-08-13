#!/bin/bash

# NPHIES Platform - Integration Test Suite
# Tests end-to-end workflows across all services

set -e

echo "🏥 NPHIES Platform - Integration Test Suite"
echo "=========================================="

# Configuration
API_GATEWAY_URL=${API_GATEWAY_URL:-"http://localhost:8080"}
ELIGIBILITY_SERVICE_URL=${ELIGIBILITY_SERVICE_URL:-"http://localhost:8090"}
CLAIMS_SERVICE_URL=${CLAIMS_SERVICE_URL:-"http://localhost:8080"}
TERMINOLOGY_SERVICE_URL=${TERMINOLOGY_SERVICE_URL:-"http://localhost:8091"}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
check_service() {
    local service_name=$1
    local url=$2
    local endpoint=$3
    
    echo -n "  ✓ Checking $service_name... "
    
    if curl -s -f "$url$endpoint" > /dev/null; then
        echo -e "${GREEN}UP${NC}"
        return 0
    else
        echo -e "${RED}DOWN${NC}"
        return 1
    fi
}

test_endpoint() {
    local description=$1
    local method=$2
    local url=$3
    local data=$4
    local expected_status=$5
    
    echo -n "  ✓ $description... "
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "%{http_code}" -o /tmp/response.json "$url")
    else
        response=$(curl -s -w "%{http_code}" -o /tmp/response.json -X "$method" -H "Content-Type: application/json" -d "$data" "$url")
    fi
    
    if [ "$response" = "$expected_status" ]; then
        echo -e "${GREEN}PASS${NC}"
        return 0
    else
        echo -e "${RED}FAIL (HTTP $response)${NC}"
        return 1
    fi
}

# Test 1: Service Health Checks
echo ""
echo "1. Testing Service Health Checks..."

HEALTH_FAILED=0

check_service "API Gateway" "$API_GATEWAY_URL" "/health" || HEALTH_FAILED=1
check_service "Eligibility Service" "$ELIGIBILITY_SERVICE_URL" "/health" || HEALTH_FAILED=1
check_service "Claims Service" "$CLAIMS_SERVICE_URL" "/api/v1/claims/health" || HEALTH_FAILED=1
check_service "Terminology Service" "$TERMINOLOGY_SERVICE_URL" "/health" || HEALTH_FAILED=1

if [ $HEALTH_FAILED -eq 1 ]; then
    echo -e "${RED}❌ Some services are not healthy. Stopping tests.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ All services are healthy${NC}"

# Test 2: Terminology Service Tests
echo ""
echo "2. Testing Terminology Service..."

test_endpoint "List code systems" "GET" "$TERMINOLOGY_SERVICE_URL/api/v1/code-systems" "" "200"
test_endpoint "Lookup ICD-10 code" "GET" "$TERMINOLOGY_SERVICE_URL/api/v1/codes/lookup/icd-10/Z00.00" "" "200"

# Test 3: Eligibility Service Tests
echo ""
echo "3. Testing Eligibility Service..."

ELIGIBILITY_DATA='{
  "member_id": "12345678901",
  "payer_id": "PAY001"
}'

test_endpoint "Check eligibility" "POST" "$ELIGIBILITY_SERVICE_URL/api/v1/eligibility/check" "$ELIGIBILITY_DATA" "200"

# Test 4: Claims Service Tests
echo ""
echo "4. Testing Claims Service..."

CLAIM_DATA='{
  "provider_id": "PRV001",
  "member_id": "12345678901",
  "payer_id": "PAY001",
  "service_date": "2025-08-13T10:30:00",
  "total_amount": 150.00,
  "type": "PROFESSIONAL",
  "claim_lines": [{
    "service_code": "99213",
    "service_date": "2025-08-13T10:30:00",
    "units": 1,
    "charged_amount": 150.00,
    "place_of_service": "11"
  }],
  "diagnosis_codes": [{
    "code": "Z00.00",
    "code_type": "ICD-10",
    "is_primary": true
  }]
}'

test_endpoint "Submit claim" "POST" "$CLAIMS_SERVICE_URL/api/v1/claims/submit" "$CLAIM_DATA" "202"

# Test 5: End-to-End Workflow Test
echo ""
echo "5. Testing End-to-End Claim Workflow..."

echo "  ✓ Step 1: Check member eligibility..."
ELIGIBILITY_RESULT=$(curl -s -X POST -H "Content-Type: application/json" -d "$ELIGIBILITY_DATA" "$ELIGIBILITY_SERVICE_URL/api/v1/eligibility/check")

if echo "$ELIGIBILITY_RESULT" | grep -q '"eligible".*true'; then
    echo -e "    ${GREEN}Member is eligible${NC}"
else
    echo -e "    ${YELLOW}Member eligibility check returned: $ELIGIBILITY_RESULT${NC}"
fi

echo "  ✓ Step 2: Validate diagnosis code..."
DIAGNOSIS_VALIDATION=$(curl -s "$TERMINOLOGY_SERVICE_URL/api/v1/codes/lookup/icd-10/Z00.00")

if echo "$DIAGNOSIS_VALIDATION" | grep -q '"found".*true'; then
    echo -e "    ${GREEN}Diagnosis code is valid${NC}"
else
    echo -e "    ${YELLOW}Diagnosis code validation returned: $DIAGNOSIS_VALIDATION${NC}"
fi

echo "  ✓ Step 3: Submit claim with validated data..."
CLAIM_RESULT=$(curl -s -w "%{http_code}" -o /tmp/claim_response.json -X POST -H "Content-Type: application/json" -d "$CLAIM_DATA" "$CLAIMS_SERVICE_URL/api/v1/claims/submit")

if [ "$CLAIM_RESULT" = "202" ]; then
    echo -e "    ${GREEN}Claim submitted successfully${NC}"
    CLAIM_ID=$(cat /tmp/claim_response.json | grep -o '"claim_id":"[^"]*"' | cut -d'"' -f4)
    echo "    Claim ID: $CLAIM_ID"
    
    if [ -n "$CLAIM_ID" ]; then
        echo "  ✓ Step 4: Check claim status..."
        sleep 2
        CLAIM_STATUS=$(curl -s "$CLAIMS_SERVICE_URL/api/v1/claims/$CLAIM_ID/status")
        echo "    Claim status: $CLAIM_STATUS"
    fi
else
    echo -e "    ${RED}Claim submission failed (HTTP $CLAIM_RESULT)${NC}"
fi

# Test 6: Performance and Load Test (Basic)
echo ""
echo "6. Basic Performance Tests..."

echo "  ✓ Testing API Gateway response time..."
start_time=$(date +%s%N)
curl -s "$API_GATEWAY_URL/health" > /dev/null
end_time=$(date +%s%N)
response_time=$(((end_time - start_time) / 1000000))

if [ $response_time -lt 1000 ]; then
    echo -e "    ${GREEN}Response time: ${response_time}ms (GOOD)${NC}"
else
    echo -e "    ${YELLOW}Response time: ${response_time}ms (SLOW)${NC}"
fi

# Summary
echo ""
echo "=========================================="
echo -e "${GREEN}✅ Integration tests completed!${NC}"
echo ""
echo "📊 Test Summary:"
echo "  - Service health checks: PASSED"
echo "  - Terminology service: TESTED"
echo "  - Eligibility service: TESTED"
echo "  - Claims service: TESTED"
echo "  - End-to-end workflow: TESTED"
echo "  - Basic performance: TESTED"
echo ""
echo "🎯 Next steps:"
echo "  - Run load tests with higher concurrency"
echo "  - Test error handling and edge cases"
echo "  - Validate monitoring and alerting"
echo "  - Test disaster recovery procedures"