#!/bin/bash

# API Server 使用示例

BASE_URL="http://localhost:8080"

echo "========================================="
echo "API Server 使用示例"
echo "========================================="

# 1. 健康检查
echo -e "\n1. 健康检查"
curl -X GET "$BASE_URL/health" | jq .

# 2. 就绪检查
echo -e "\n2. 就绪检查"
curl -X GET "$BASE_URL/ready" | jq .

# 3. 登录获取 Token (admin 用户)
echo -e "\n3. Admin 用户登录"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }')

echo $LOGIN_RESPONSE | jq .

# 提取 access token
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.access_token')
REFRESH_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.refresh_token')

echo "Access Token: $ACCESS_TOKEN"

# 4. 列出所有用户 (需要 admin 权限)
echo -e "\n4. 列出所有用户 (需要 admin 权限)"
curl -X GET "$BASE_URL/api/v1/users" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# 5. 获取用户详情
echo -e "\n5. 获取用户详情"
curl -X GET "$BASE_URL/api/v1/users/1" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# 6. 列出资源
echo -e "\n6. 列出资源"
curl -X GET "$BASE_URL/api/v1/resources" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# 7. 创建资源 (需要 editor 权限)
echo -e "\n7. 创建资源 (使用 editor 用户)"
EDITOR_LOGIN=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "editor",
    "password": "editor123"
  }')

EDITOR_TOKEN=$(echo $EDITOR_LOGIN | jq -r '.access_token')

curl -X POST "$BASE_URL/api/v1/resources" \
  -H "Authorization: Bearer $EDITOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Resource",
    "description": "Created via API",
    "type": "compute",
    "metadata": {
      "env": "dev",
      "region": "us-west-1"
    }
  }' | jq .

# 8. 测试权限不足的情况
echo -e "\n8. 测试权限不足 (viewer 用户尝试创建资源)"
VIEWER_LOGIN=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "viewer",
    "password": "viewer123"
  }')

VIEWER_TOKEN=$(echo $VIEWER_LOGIN | jq -r '.access_token')

curl -X POST "$BASE_URL/api/v1/resources" \
  -H "Authorization: Bearer $VIEWER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Should Fail",
    "description": "This should fail",
    "type": "compute"
  }' | jq .

# 9. 刷新 Token
echo -e "\n9. 刷新 Token"
curl -X POST "$BASE_URL/api/v1/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }" | jq .

# 10. 无效 Token 测试
echo -e "\n10. 无效 Token 测试"
curl -X GET "$BASE_URL/api/v1/users" \
  -H "Authorization: Bearer invalid-token" | jq .

echo -e "\n========================================="
echo "示例完成"
echo "========================================="
