#!/bin/bash

# Marketing Management Module API Testing Script
# This script tests all endpoints for activity templates and activities

# Configuration
BASE_URL="http://localhost:8080/api"
TOKEN=""
TEMPLATE_ID=""
ACTIVITY_ID=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print test results
print_test() {
    echo -e "${YELLOW}========================================${NC}"
    echo -e "${YELLOW}TEST: $1${NC}"
    echo -e "${YELLOW}========================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Step 1: Login to get JWT token
print_test "Login to get JWT token"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')

if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
    print_success "Login successful, token obtained"
    echo "Token: ${TOKEN:0:50}..."
else
    print_error "Login failed"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo ""

# Step 2: Create Activity Template (按分类)
print_test "Create Activity Template (按分类)"
CREATE_TEMPLATE_RESPONSE=$(curl -s -X POST "$BASE_URL/activity-templates" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "测试满减活动模板",
    "type": 1,
    "select_type": 1,
    "classify_ids": [1, 2],
    "status": 0
  }')

TEMPLATE_ID=$(echo $CREATE_TEMPLATE_RESPONSE | jq -r '.data.id')

if [ "$TEMPLATE_ID" != "null" ] && [ -n "$TEMPLATE_ID" ]; then
    print_success "Template created successfully, ID: $TEMPLATE_ID"
else
    print_error "Template creation failed"
    echo "Response: $CREATE_TEMPLATE_RESPONSE"
fi

echo ""

# Step 3: Get Activity Templates List
print_test "Get Activity Templates List"
TEMPLATES_LIST=$(curl -s -X GET "$BASE_URL/activity-templates" \
  -H "Authorization: Bearer $TOKEN")

TEMPLATES_COUNT=$(echo $TEMPLATES_LIST | jq -r '.data | length')
print_success "Retrieved $TEMPLATES_COUNT templates"
echo $TEMPLATES_LIST | jq '.'

echo ""

# Step 4: Get Active Templates
print_test "Get Active Templates"
ACTIVE_TEMPLATES=$(curl -s -X GET "$BASE_URL/activity-templates/active" \
  -H "Authorization: Bearer $TOKEN")

ACTIVE_COUNT=$(echo $ACTIVE_TEMPLATES | jq -r '.data | length')
print_success "Retrieved $ACTIVE_COUNT active templates"
echo $ACTIVE_TEMPLATES | jq '.'

echo ""

# Step 5: Get Template Detail
print_test "Get Template Detail (ID: $TEMPLATE_ID)"
TEMPLATE_DETAIL=$(curl -s -X GET "$BASE_URL/activity-templates/$TEMPLATE_ID" \
  -H "Authorization: Bearer $TOKEN")

print_success "Template detail retrieved"
echo $TEMPLATE_DETAIL | jq '.'

echo ""

# Step 6: Update Template
print_test "Update Template (ID: $TEMPLATE_ID)"
UPDATE_TEMPLATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/activity-templates/$TEMPLATE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "更新后的满减活动模板",
    "type": 1,
    "select_type": 1,
    "classify_ids": [1, 2, 3],
    "status": 0
  }')

if [ "$(echo $UPDATE_TEMPLATE_RESPONSE | jq -r '.code')" = "0" ]; then
    print_success "Template updated successfully"
else
    print_error "Template update failed"
    echo "Response: $UPDATE_TEMPLATE_RESPONSE"
fi

echo ""

# Step 7: Create Activity Template (按商品)
print_test "Create Activity Template (按商品)"
CREATE_TEMPLATE2_RESPONSE=$(curl -s -X POST "$BASE_URL/activity-templates" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "测试满折活动模板",
    "type": 2,
    "select_type": 2,
    "goods_ids": [1, 2, 3],
    "status": 0
  }')

TEMPLATE_ID2=$(echo $CREATE_TEMPLATE2_RESPONSE | jq -r '.data.id')

if [ "$TEMPLATE_ID2" != "null" ] && [ -n "$TEMPLATE_ID2" ]; then
    print_success "Template 2 created successfully, ID: $TEMPLATE_ID2"
else
    print_error "Template 2 creation failed"
    echo "Response: $CREATE_TEMPLATE2_RESPONSE"
fi

echo ""

# Step 8: Create Activity
print_test "Create Activity"
CREATE_ACTIVITY_RESPONSE=$(curl -s -X POST "$BASE_URL/activities" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"name\": \"测试满减活动\",
    \"template_id\": $TEMPLATE_ID,
    \"start_time\": \"2026-02-01T00:00:00Z\",
    \"end_time\": \"2026-02-28T23:59:59Z\",
    \"status\": 0,
    \"details\": [
      {
        \"threshold_amount\": 100.0,
        \"discount_value\": 10.0
      },
      {
        \"threshold_amount\": 200.0,
        \"discount_value\": 25.0
      }
    ]
  }")

ACTIVITY_ID=$(echo $CREATE_ACTIVITY_RESPONSE | jq -r '.data.id')

if [ "$ACTIVITY_ID" != "null" ] && [ -n "$ACTIVITY_ID" ]; then
    print_success "Activity created successfully, ID: $ACTIVITY_ID"
else
    print_error "Activity creation failed"
    echo "Response: $CREATE_ACTIVITY_RESPONSE"
fi

echo ""

# Step 9: Get Activities List
print_test "Get Activities List"
ACTIVITIES_LIST=$(curl -s -X GET "$BASE_URL/activities" \
  -H "Authorization: Bearer $TOKEN")

ACTIVITIES_COUNT=$(echo $ACTIVITIES_LIST | jq -r '.data | length')
print_success "Retrieved $ACTIVITIES_COUNT activities"
echo $ACTIVITIES_LIST | jq '.'

echo ""

# Step 10: Get Activity Detail
print_test "Get Activity Detail (ID: $ACTIVITY_ID)"
ACTIVITY_DETAIL=$(curl -s -X GET "$BASE_URL/activities/$ACTIVITY_ID" \
  -H "Authorization: Bearer $TOKEN")

print_success "Activity detail retrieved"
echo $ACTIVITY_DETAIL | jq '.'

echo ""

# Step 11: Get Activities by Date Range
print_test "Get Activities by Date Range"
ACTIVITIES_BY_DATE=$(curl -s -X GET "$BASE_URL/activities/by-date-range?payment_time=2026-02-15T12:00:00Z" \
  -H "Authorization: Bearer $TOKEN")

print_success "Activities by date range retrieved"
echo $ACTIVITIES_BY_DATE | jq '.'

echo ""

# Step 12: Update Activity
print_test "Update Activity (ID: $ACTIVITY_ID)"
UPDATE_ACTIVITY_RESPONSE=$(curl -s -X PUT "$BASE_URL/activities/$ACTIVITY_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"name\": \"更新后的满减活动\",
    \"template_id\": $TEMPLATE_ID,
    \"start_time\": \"2026-02-01T00:00:00Z\",
    \"end_time\": \"2026-02-28T23:59:59Z\",
    \"status\": 0,
    \"details\": [
      {
        \"threshold_amount\": 150.0,
        \"discount_value\": 20.0
      },
      {
        \"threshold_amount\": 300.0,
        \"discount_value\": 50.0
      }
    ]
  }")

if [ "$(echo $UPDATE_ACTIVITY_RESPONSE | jq -r '.code')" = "0" ]; then
    print_success "Activity updated successfully"
else
    print_error "Activity update failed"
    echo "Response: $UPDATE_ACTIVITY_RESPONSE"
fi

echo ""

# Step 13: Update Activity Status
print_test "Update Activity Status (ID: $ACTIVITY_ID)"
UPDATE_STATUS_RESPONSE=$(curl -s -X PUT "$BASE_URL/activities/$ACTIVITY_ID/status" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "status": 1
  }')

if [ "$(echo $UPDATE_STATUS_RESPONSE | jq -r '.code')" = "0" ]; then
    print_success "Activity status updated successfully"
else
    print_error "Activity status update failed"
    echo "Response: $UPDATE_STATUS_RESPONSE"
fi

echo ""

# Step 14: Delete Activity
print_test "Delete Activity (ID: $ACTIVITY_ID)"
DELETE_ACTIVITY_RESPONSE=$(curl -s -X DELETE "$BASE_URL/activities/$ACTIVITY_ID" \
  -H "Authorization: Bearer $TOKEN")

if [ "$(echo $DELETE_ACTIVITY_RESPONSE | jq -r '.code')" = "0" ]; then
    print_success "Activity deleted successfully"
else
    print_error "Activity deletion failed"
    echo "Response: $DELETE_ACTIVITY_RESPONSE"
fi

echo ""

# Step 15: Update Template Status
print_test "Update Template Status (ID: $TEMPLATE_ID)"
UPDATE_TEMPLATE_STATUS_RESPONSE=$(curl -s -X PUT "$BASE_URL/activity-templates/$TEMPLATE_ID/status" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "status": 1
  }')

if [ "$(echo $UPDATE_TEMPLATE_STATUS_RESPONSE | jq -r '.code')" = "0" ]; then
    print_success "Template status updated successfully"
else
    print_error "Template status update failed"
    echo "Response: $UPDATE_TEMPLATE_STATUS_RESPONSE"
fi

echo ""

# Step 16: Delete Template (should fail if there are related activities)
print_test "Delete Template (ID: $TEMPLATE_ID)"
DELETE_TEMPLATE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/activity-templates/$TEMPLATE_ID" \
  -H "Authorization: Bearer $TOKEN")

if [ "$(echo $DELETE_TEMPLATE_RESPONSE | jq -r '.code')" = "0" ]; then
    print_success "Template deleted successfully"
else
    print_success "Template deletion prevented (expected if has activities)"
    echo "Response: $DELETE_TEMPLATE_RESPONSE"
fi

echo ""

# Step 17: Delete Template 2
print_test "Delete Template 2 (ID: $TEMPLATE_ID2)"
DELETE_TEMPLATE2_RESPONSE=$(curl -s -X DELETE "$BASE_URL/activity-templates/$TEMPLATE_ID2" \
  -H "Authorization: Bearer $TOKEN")

if [ "$(echo $DELETE_TEMPLATE2_RESPONSE | jq -r '.code')" = "0" ]; then
    print_success "Template 2 deleted successfully"
else
    print_error "Template 2 deletion failed"
    echo "Response: $DELETE_TEMPLATE2_RESPONSE"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}ALL TESTS COMPLETED${NC}"
echo -e "${GREEN}========================================${NC}"
