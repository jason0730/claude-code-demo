# API Server 架构设计

## 概述
这是一个云原生的 API Server，实现了标准的认证和授权机制。

## 核心特性

### 1. 认证机制
- **JWT (JSON Web Token)**: 使用 JWT 进行无状态认证
- 支持 Token 过期和刷新
- 使用 RS256 算法签名（公钥/私钥）

### 2. 授权机制
- **RBAC (Role-Based Access Control)**: 基于角色的访问控制
- 支持多角色分配
- 细粒度的权限控制

### 3. 云原生特性
- 12-Factor App 原则
- 环境变量配置
- 健康检查端点 (/health, /ready)
- 优雅关闭 (Graceful Shutdown)
- 结构化日志
- 容器化支持
- Kubernetes 就绪

## 架构组件

```
├── cmd/
│   └── api-server/          # 主程序入口
├── internal/
│   ├── config/              # 配置管理
│   ├── auth/                # 认证模块
│   │   ├── jwt/             # JWT 实现
│   │   └── middleware/      # 认证中间件
│   ├── authz/               # 授权模块
│   │   ├── rbac/            # RBAC 实现
│   │   └── middleware/      # 授权中间件
│   ├── handler/             # HTTP 处理器
│   └── model/               # 数据模型
├── deployments/
│   ├── kubernetes/          # K8s 部署配置
│   └── docker/              # Docker 配置
└── examples/                # 使用示例
```

## API 端点设计

### 认证端点
- `POST /api/v1/auth/login` - 用户登录，获取 JWT Token
- `POST /api/v1/auth/refresh` - 刷新 Token

### 业务端点（需要认证和授权）
- `GET /api/v1/users` - 列出用户（需要 admin 角色）
- `GET /api/v1/users/:id` - 获取用户详情（需要 admin 或 user 角色）
- `POST /api/v1/resources` - 创建资源（需要 editor 角色）
- `GET /api/v1/resources` - 列出资源（需要 viewer 角色）

### 系统端点
- `GET /health` - 健康检查（存活探针）
- `GET /ready` - 就绪检查（就绪探针）
- `GET /metrics` - Prometheus 指标

## 权限模型

### 角色定义
- **admin**: 管理员，拥有所有权限
- **editor**: 编辑者，可以创建和修改资源
- **viewer**: 查看者，只能查看资源
- **user**: 普通用户，可以查看自己的信息

### 权限检查流程
1. 请求到达 -> 认证中间件验证 JWT
2. 提取用户信息和角色
3. 授权中间件检查角色权限
4. 执行业务逻辑

## 安全考虑
- 使用 HTTPS/TLS
- JWT Token 短期有效（15分钟）
- Refresh Token 长期有效（7天）
- 敏感配置使用 Secret 管理
- 请求限流和防 DDoS
- CORS 配置
- 安全头设置

## 可观测性
- 结构化日志（JSON格式）
- 请求追踪 ID
- Prometheus 指标导出
- 错误追踪

## 部署方式
- Docker 容器
- Kubernetes Deployment
- 支持水平扩展
- 无状态设计
