# API Server - 云原生权限管理系统

一个基于云原生架构设计的 API Server，实现了完整的认证和授权机制。

## 特性

### 核心功能
- ✅ **JWT 认证**: 使用 JWT Token 进行无状态认证
- ✅ **RBAC 授权**: 基于角色的访问控制
- ✅ **云原生设计**: 遵循 12-Factor App 原则
- ✅ **容器化**: 支持 Docker 和 Kubernetes 部署
- ✅ **健康检查**: 提供存活和就绪探针
- ✅ **优雅关闭**: 支持优雅关闭机制
- ✅ **结构化日志**: JSON 格式的结构化日志
- ✅ **水平扩展**: 无状态设计，支持水平扩展

### 安全特性
- 🔒 JWT Token 短期有效（15分钟）
- 🔒 Refresh Token 长期有效（7天）
- 🔒 环境变量配置敏感信息
- 🔒 细粒度的权限控制
- 🔒 中间件架构确保安全

## 快速开始

### 前置要求
- Go 1.21+
- Docker (可选)
- Kubernetes (可选)

### 本地运行

1. 克隆仓库
```bash
git clone https://github.com/jason0730/claude-code-demo.git
cd claude-code-demo
```

2. 安装依赖
```bash
make deps
```

3. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 文件，修改配置
```

4. 运行应用
```bash
make run
# 或直接运行
go run cmd/api-server/main.go
```

5. 测试 API
```bash
# 健康检查
curl http://localhost:8080/health

# 登录获取 Token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 使用 Docker

1. 构建镜像
```bash
make docker-build
```

2. 运行容器
```bash
make docker-run
```

3. 停止容器
```bash
make docker-stop
```

### 部署到 Kubernetes

1. 构建并推送镜像
```bash
docker build -t your-registry/api-server:v1.0.0 .
docker push your-registry/api-server:v1.0.0
```

2. 更新镜像配置
```bash
# 编辑 deployments/kubernetes/kustomization.yaml
# 修改 image 配置
```

3. 部署
```bash
make k8s-deploy
```

4. 查看状态
```bash
kubectl get pods -n api-server
kubectl get svc -n api-server
```

5. 查看日志
```bash
make k8s-logs
```

## API 文档

### 端点概览

#### 系统端点（无需认证）
- `GET /health` - 健康检查（存活探针）
- `GET /ready` - 就绪检查（就绪探针）
- `GET /metrics` - Prometheus 指标

#### 认证端点（无需认证）
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/refresh` - 刷新 Token

#### 用户端点（需要认证）
- `GET /api/v1/users` - 列出所有用户（需要 admin 角色）
- `GET /api/v1/users/{id}` - 获取用户详情（需要 admin 或 user 角色）

#### 资源端点（需要认证）
- `GET /api/v1/resources` - 列出资源（需要 viewer 权限）
- `POST /api/v1/resources` - 创建资源（需要 editor 权限）

### 示例请求

详细的 API 使用示例请参考 `examples/api_examples.sh`

#### 登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

响应:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900,
  "token_type": "Bearer"
}
```

#### 访问受保护的端点
```bash
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 角色和权限

### 预定义角色

| 角色 | 权限 | 说明 |
|------|------|------|
| admin | 所有权限 | 管理员，拥有完全访问权限 |
| editor | resource:read, resource:write, resource:list, user:read | 可以创建和修改资源 |
| viewer | resource:read, resource:list | 只能查看资源 |
| user | user:read, resource:read | 普通用户，可以查看自己的信息 |

### 测试用户

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | admin123 | admin |
| editor | editor123 | editor |
| viewer | viewer123 | viewer |

**注意**: 这些是示例用户，仅用于测试。生产环境请使用真实的用户管理系统。

## 配置

### 环境变量

所有配置通过环境变量管理，支持的环境变量：

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| SERVER_HOST | 0.0.0.0 | 服务器监听地址 |
| SERVER_PORT | 8080 | 服务器监听端口 |
| JWT_SECRET | - | JWT 签名密钥 |
| JWT_EXPIRATION | 15m | JWT 过期时间 |
| REFRESH_EXPIRATION | 168h | 刷新令牌过期时间 |
| LOG_LEVEL | info | 日志级别 |
| LOG_FORMAT | json | 日志格式 |

详细配置请参考 `.env.example`

## 架构设计

项目采用标准的云原生架构设计，详细架构说明请参考 [ARCHITECTURE.md](ARCHITECTURE.md)

### 目录结构
```
.
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

## 开发

### 运行测试
```bash
make test
```

### 代码格式化
```bash
make fmt
```

### 代码检查
```bash
make lint
```

### 构建
```bash
make build
```

## 生产部署建议

1. **使用 HTTPS/TLS**: 在生产环境必须使用 HTTPS
2. **使用密钥管理服务**: 不要在代码中硬编码密钥，使用 Vault 等密钥管理系统
3. **使用 RSA 密钥对**: 生产环境建议使用 RS256 算法替代 HS256
4. **启用速率限制**: 在 Ingress 或 API Gateway 层面启用速率限制
5. **监控和告警**: 集成 Prometheus 和 Grafana 进行监控
6. **日志聚合**: 使用 ELK 或其他日志聚合系统
7. **定期更新依赖**: 保持依赖库的更新以修复安全漏洞

## 监控

### Prometheus 指标

访问 `/metrics` 端点获取 Prometheus 格式的指标数据。

### 健康检查

- 存活探针: `GET /health`
- 就绪探针: `GET /ready`

## 故障排查

### 常见问题

1. **Token 过期**
   - 使用 refresh token 获取新的 access token

2. **权限不足**
   - 检查用户角色和端点所需权限

3. **无法连接**
   - 检查服务是否正常运行
   - 检查端口是否被占用
   - 检查防火墙设置

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 联系方式

如有问题，请提交 Issue 或联系维护者。
