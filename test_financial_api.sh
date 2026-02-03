#!/bin/bash

# 财务管理模块API测试脚本

BASE_URL="http://localhost:5001"

echo "========================================"
echo "财务管理模块 API 测试"
echo "========================================"
echo ""

# 测试1: 登录（获取token）
echo "测试1: 用户登录"
echo "----------------------------------------"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')

echo "登录响应: $LOGIN_RESPONSE"

# 如果登录失败，尝试其他用户
if [[ $LOGIN_RESPONSE == *"错误"* ]]; then
  echo "尝试使用测试用户..."
  LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"test","password":"test123"}')
  echo "登录响应: $LOGIN_RESPONSE"
fi

# 提取token（假设返回格式包含token）
TOKEN=$(echo $LOGIN_RESPONSE | grep -oP '(?<="token":")[^"]*')

if [ -z "$TOKEN" ]; then
  echo "⚠️  未获取到token，可能需要手动登录"
  echo "继续测试未认证的接口..."
  echo ""
else
  echo "✅ Token已获取: ${TOKEN:0:20}..."
  echo ""
fi

# 测试2: 获取收款列表（无token - 应返回401）
echo "测试2: 获取收款列表（无token）"
echo "----------------------------------------"
RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET "$BASE_URL/api/payment-collections?page=1&page_size=10")
HTTP_CODE=$(echo "$RESPONSE" | grep -oP '(?<=HTTP_CODE:)[0-9]+')
BODY=$(echo "$RESPONSE" | sed 's/HTTP_CODE:[0-9]*//')

echo "HTTP状态码: $HTTP_CODE"
echo "响应内容: $BODY"
echo ""

# 测试3: 获取收款列表（有token）
if [ ! -z "$TOKEN" ]; then
  echo "测试3: 获取收款列表（有token）"
  echo "----------------------------------------"
  RESPONSE=$(curl -s -X GET "$BASE_URL/api/payment-collections?page=1&page_size=10" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

  echo "响应: $RESPONSE"
  echo ""
fi

# 测试4: 新增收款（模拟数据）
if [ ! -z "$TOKEN" ]; then
  echo "测试4: 新增收款"
  echo "----------------------------------------"
  CREATE_DATA='{
    "order_id": 1,
    "student_id": 1,
    "payment_scenario": 1,
    "payment_method": 0,
    "payment_amount": 100.00,
    "payer": "测试用户",
    "payee_entity": 0,
    "merchant_order": "TEST001"
  }'

  RESPONSE=$(curl -s -X POST "$BASE_URL/api/payment-collections" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "$CREATE_DATA")

  echo "响应: $RESPONSE"
  echo ""

  # 提取收款ID
  PAYMENT_ID=$(echo $RESPONSE | grep -oP '(?<="id":)[0-9]+')

  if [ ! -z "$PAYMENT_ID" ]; then
    echo "✅ 收款创建成功，ID: $PAYMENT_ID"
    echo ""

    # 测试5: 确认到账
    echo "测试5: 确认收款到账（ID: $PAYMENT_ID）"
    echo "----------------------------------------"
    RESPONSE=$(curl -s -X PUT "$BASE_URL/api/payment-collections/$PAYMENT_ID/confirm" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json")

    echo "响应: $RESPONSE"
    echo ""

    # 测试6: 查询分账明细
    echo "测试6: 查询分账明细"
    echo "----------------------------------------"
    RESPONSE=$(curl -s -X GET "$BASE_URL/api/separate-accounts?payment_id=$PAYMENT_ID" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json")

    echo "响应: $RESPONSE"
    echo ""
  fi
fi

# 测试7: 删除收款
if [ ! -z "$TOKEN" ] && [ ! -z "$PAYMENT_ID" ]; then
  echo "测试7: 删除收款（ID: $PAYMENT_ID）"
  echo "----------------------------------------"
  RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/payment-collections/$PAYMENT_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

  echo "响应: $RESPONSE"
  echo ""
fi

echo "========================================"
echo "测试完成"
echo "========================================"
