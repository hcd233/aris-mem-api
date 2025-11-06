# Go Backend Template / Go 后端模板

[English](#english) | [中文](#中文)

---

<a name="english"></a>
## English

### 📖 Introduction

A production-ready Go backend template built with modern technologies and best practices. This template provides a robust foundation for building RESTful APIs with authentication, database management, object storage, and more.

### ✨ Features

- 🚀 **High Performance**: Built with [Fiber](https://gofiber.io/) framework and [Sonic](https://github.com/bytedance/sonic) JSON serialization
- 🔐 **Authentication**: JWT-based authentication with access and refresh tokens
- 🌐 **OAuth2 Integration**: Support for GitHub and Google OAuth2 login
- 💾 **Database**: PostgreSQL with GORM ORM
- 📦 **Object Storage**: Support for both MinIO and Tencent COS
- 🔴 **Caching**: Redis integration for high-performance caching
- 🤖 **AI Integration**: OpenAI client integration
- 📝 **API Documentation**: Auto-generated Swagger documentation
- 🔒 **Middleware**: Comprehensive middleware stack including:
  - JWT authentication
  - CORS
  - Rate limiting
  - Request logging with trace ID
  - Compression
  - Recovery from panics
  - Permission validation
- 🎯 **Project Structure**: Clean architecture with separation of concerns
- 🐳 **Docker Support**: Complete Docker Compose setup for easy deployment
- ⏰ **Scheduled Tasks**: Cron job support
- 📊 **Profiling**: Built-in fgprof for performance profiling

### 🛠️ Tech Stack

- **Framework**: Fiber v2
- **Language**: Go 1.25.1
- **Database**: PostgreSQL (with GORM)
- **Cache**: Redis
- **Object Storage**: MinIO / Tencent COS
- **Authentication**: JWT, OAuth2 (GitHub, Google)
- **API Docs**: Swagger/OpenAPI
- **CLI**: Cobra
- **Configuration**: Viper
- **Logging**: Zap with Lumberjack rotation
- **JSON**: Sonic (high performance)

### 📁 Project Structure

```
aris-mem-api/
├── cmd/                    # Command line interface
│   ├── server.go          # Server start command
│   ├── database.go        # Database management commands
│   └── root.go            # Root command
├── internal/
│   ├── auth/              # JWT authentication logic
│   ├── config/            # Configuration management
│   ├── constant/          # Constants
│   ├── cron/              # Scheduled tasks
│   ├── handler/           # HTTP request handlers
│   ├── logger/            # Logging utilities
│   ├── middleware/        # HTTP middlewares
│   ├── protocol/          # Request/response protocols
│   ├── resource/          # External resource integrations
│   │   ├── cache/         # Redis cache
│   │   ├── database/      # PostgreSQL + GORM
│   │   ├── llm/           # OpenAI client
│   │   └── storage/       # Object storage (MinIO/COS)
│   ├── router/            # Route definitions
│   ├── service/           # Business logic
│   └── util/              # Utility functions
├── docker/                # Docker configuration files
├── env/                   # Environment variable templates
├── docs/                  # Swagger documentation
└── main.go               # Application entry point
```

### 🚀 Quick Start

#### Prerequisites

- Go 1.25.1 or higher
- Docker and Docker Compose (for containerized setup)
- PostgreSQL (if running locally)
- Redis (if running locally)

#### 1. Clone the Repository

```bash
git clone https://github.com/hcd233/aris-mem-api.git
cd aris-mem-api
```

#### 2. Configure Environment Variables

Copy the environment template and modify as needed:

```bash
cp env/api.env.template env/api.env
# Edit env/api.env with your configurations
```

Key configurations to set:
- Database credentials (`POSTGRES_*`)
- Redis connection (`REDIS_*`)
- JWT secrets (`JWT_ACCESS_TOKEN_SECRET`, `JWT_REFRESH_TOKEN_SECRET`)
- OAuth2 credentials (if using OAuth2 login)
- Object storage credentials (MinIO or COS)
- OpenAI API key (if using AI features)

#### 3. Run with Docker Compose (Recommended)

Create required volumes:
```bash
docker volume create postgresql-data
docker volume create redis-data
docker volume create minio-data
```

Start all services:
```bash
docker compose -f docker/docker-compose.yml up -d
```

This will start:
- PostgreSQL database
- Redis cache
- MinIO object storage
- The API server (accessible at http://localhost:8170)

#### 4. Run Locally

Install dependencies:
```bash
go mod download
```

Run database migration:
```bash
go run main.go database migrate
```

Start the server:
```bash
go run main.go server start --host localhost --port 8080
```

### 📚 API Documentation

Once the server is running, access the Swagger documentation at:
```
http://localhost:8080/swagger/
```

### 🔑 Available Commands

```bash
# Start the server
go run main.go server start [--host HOST] [--port PORT]

# Database migration
go run main.go database migrate

# Object storage management (if applicable)
go run main.go object [subcommand]
```

### 🔐 Authentication

The API supports multiple authentication methods:

1. **OAuth2**: Login via GitHub or Google
   - `GET /v1/oauth2/github/login`
   - `GET /v1/oauth2/google/login`

2. **JWT Tokens**: After OAuth2 login, obtain access/refresh tokens
   - `POST /v1/token/refresh` - Refresh access token

### 🛡️ API Endpoints

- `GET /` - Health check
- `GET /swagger/*` - API documentation
- `GET /v1/oauth2/{provider}/login` - OAuth2 login
- `GET /v1/oauth2/{provider}/callback` - OAuth2 callback
- `POST /v1/token/refresh` - Refresh JWT token
- `GET /v1/user/current` - Get current user info (requires auth)
- `GET /v1/user/{userID}` - Get user info by ID (requires auth)
- `PATCH /v1/user` - Update user info (requires auth)

### 🔧 Development

Build the binary:
```bash
go build -o aris-mem-api main.go
```

Run tests (if available):
```bash
go test ./...
```

Generate Swagger docs:
```bash
swag init
```

### 📝 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | 8080 |
| `READ_TIMEOUT` | Read timeout in seconds | 10 |
| `WRITE_TIMEOUT` | Write timeout in seconds | 10 |
| `LOG_LEVEL` | Logging level | INFO |
| `POSTGRES_*` | PostgreSQL connection settings | - |
| `REDIS_*` | Redis connection settings | - |
| `JWT_ACCESS_TOKEN_EXPIRED` | Access token expiry | 12h |
| `JWT_REFRESH_TOKEN_EXPIRED` | Refresh token expiry | 168h |
| `OAUTH2_*` | OAuth2 provider settings | - |
| `MINIO_*` | MinIO storage settings | - |
| `COS_*` | Tencent COS storage settings | - |
| `OPENAI_*` | OpenAI API settings | - |

### 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

### 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

<a name="中文"></a>
## 中文

### 📖 简介

一个基于现代技术栈和最佳实践构建的生产级 Go 后端模板。该模板为构建具有身份验证、数据库管理、对象存储等功能的 RESTful API 提供了坚实的基础。

### ✨ 特性

- 🚀 **高性能**: 使用 [Fiber](https://gofiber.io/) 框架和 [Sonic](https://github.com/bytedance/sonic) JSON 序列化
- 🔐 **身份验证**: 基于 JWT 的身份验证,支持访问令牌和刷新令牌
- 🌐 **OAuth2 集成**: 支持 GitHub 和 Google OAuth2 登录
- 💾 **数据库**: PostgreSQL 配合 GORM ORM
- 📦 **对象存储**: 支持 MinIO 和腾讯云 COS
- 🔴 **缓存**: Redis 集成,提供高性能缓存
- 🤖 **AI 集成**: OpenAI 客户端集成
- 📝 **API 文档**: 自动生成的 Swagger 文档
- 🔒 **中间件**: 完善的中间件栈,包括:
  - JWT 身份验证
  - CORS 跨域处理
  - 请求限流
  - 带追踪 ID 的请求日志
  - 响应压缩
  - Panic 恢复
  - 权限验证
- 🎯 **项目结构**: 清晰的架构设计,关注点分离
- 🐳 **Docker 支持**: 完整的 Docker Compose 配置,便于部署
- ⏰ **定时任务**: Cron 定时任务支持
- 📊 **性能分析**: 内置 fgprof 性能分析工具

### 🛠️ 技术栈

- **框架**: Fiber v2
- **语言**: Go 1.25.1
- **数据库**: PostgreSQL (使用 GORM)
- **缓存**: Redis
- **对象存储**: MinIO / 腾讯云 COS
- **身份验证**: JWT, OAuth2 (GitHub, Google)
- **API 文档**: Swagger/OpenAPI
- **CLI**: Cobra
- **配置管理**: Viper
- **日志**: Zap 配合 Lumberjack 日志轮转
- **JSON**: Sonic (高性能)

### 📁 项目结构

```
aris-mem-api/
├── cmd/                    # 命令行接口
│   ├── server.go          # 服务器启动命令
│   ├── database.go        # 数据库管理命令
│   └── root.go            # 根命令
├── internal/
│   ├── auth/              # JWT 身份验证逻辑
│   ├── config/            # 配置管理
│   ├── constant/          # 常量定义
│   ├── cron/              # 定时任务
│   ├── handler/           # HTTP 请求处理器
│   ├── logger/            # 日志工具
│   ├── middleware/        # HTTP 中间件
│   ├── protocol/          # 请求/响应协议
│   ├── resource/          # 外部资源集成
│   │   ├── cache/         # Redis 缓存
│   │   ├── database/      # PostgreSQL + GORM
│   │   ├── llm/           # OpenAI 客户端
│   │   └── storage/       # 对象存储 (MinIO/COS)
│   ├── router/            # 路由定义
│   ├── service/           # 业务逻辑
│   └── util/              # 工具函数
├── docker/                # Docker 配置文件
├── env/                   # 环境变量模板
├── docs/                  # Swagger 文档
└── main.go               # 应用程序入口
```

### 🚀 快速开始

#### 前置要求

- Go 1.25.1 或更高版本
- Docker 和 Docker Compose (用于容器化部署)
- PostgreSQL (如果本地运行)
- Redis (如果本地运行)

#### 1. 克隆仓库

```bash
git clone https://github.com/hcd233/aris-mem-api.git
cd aris-mem-api
```

#### 2. 配置环境变量

复制环境变量模板并根据需要修改:

```bash
cp env/api.env.template env/api.env
# 编辑 env/api.env 填入你的配置
```

需要配置的关键项:
- 数据库凭据 (`POSTGRES_*`)
- Redis 连接 (`REDIS_*`)
- JWT 密钥 (`JWT_ACCESS_TOKEN_SECRET`, `JWT_REFRESH_TOKEN_SECRET`)
- OAuth2 凭据 (如果使用 OAuth2 登录)
- 对象存储凭据 (MinIO 或 COS)
- OpenAI API 密钥 (如果使用 AI 功能)

#### 3. 使用 Docker Compose 运行 (推荐)

创建所需的数据卷:
```bash
docker volume create postgresql-data
docker volume create redis-data
docker volume create minio-data
```

启动所有服务:
```bash
docker compose -f docker/docker-compose.yml up -d
```

这将启动:
- PostgreSQL 数据库
- Redis 缓存
- MinIO 对象存储
- API 服务器 (访问地址: http://localhost:8170)

#### 4. 本地运行

安装依赖:
```bash
go mod download
```

运行数据库迁移:
```bash
go run main.go database migrate
```

启动服务器:
```bash
go run main.go server start --host localhost --port 8080
```

### 📚 API 文档

服务器运行后,访问 Swagger 文档:
```
http://localhost:8080/swagger/
```

### 🔑 可用命令

```bash
# 启动服务器
go run main.go server start [--host HOST] [--port PORT]

# 数据库迁移
go run main.go database migrate

# 对象存储管理 (如果适用)
go run main.go object [subcommand]
```

### 🔐 身份验证

API 支持多种身份验证方式:

1. **OAuth2**: 通过 GitHub 或 Google 登录
   - `GET /v1/oauth2/github/login`
   - `GET /v1/oauth2/google/login`

2. **JWT 令牌**: OAuth2 登录后获取访问/刷新令牌
   - `POST /v1/token/refresh` - 刷新访问令牌

### 🛡️ API 端点

- `GET /` - 健康检查
- `GET /swagger/*` - API 文档
- `GET /v1/oauth2/{provider}/login` - OAuth2 登录
- `GET /v1/oauth2/{provider}/callback` - OAuth2 回调
- `POST /v1/token/refresh` - 刷新 JWT 令牌
- `GET /v1/user/current` - 获取当前用户信息 (需要认证)
- `GET /v1/user/{userID}` - 根据 ID 获取用户信息 (需要认证)
- `PATCH /v1/user` - 更新用户信息 (需要认证)

### 🔧 开发

构建二进制文件:
```bash
go build -o aris-mem-api main.go
```

运行测试 (如果有):
```bash
go test ./...
```

生成 Swagger 文档:
```bash
swag init
```

### 📝 环境变量

| 变量 | 描述 | 默认值 |
|------|------|--------|
| `PORT` | 服务器端口 | 8080 |
| `READ_TIMEOUT` | 读取超时时间(秒) | 10 |
| `WRITE_TIMEOUT` | 写入超时时间(秒) | 10 |
| `LOG_LEVEL` | 日志级别 | INFO |
| `POSTGRES_*` | PostgreSQL 连接设置 | - |
| `REDIS_*` | Redis 连接设置 | - |
| `JWT_ACCESS_TOKEN_EXPIRED` | 访问令牌过期时间 | 12h |
| `JWT_REFRESH_TOKEN_EXPIRED` | 刷新令牌过期时间 | 168h |
| `OAUTH2_*` | OAuth2 提供商设置 | - |
| `MINIO_*` | MinIO 存储设置 | - |
| `COS_*` | 腾讯云 COS 存储设置 | - |
| `OPENAI_*` | OpenAI API 设置 | - |

### 📄 许可证

本项目采用 Apache License 2.0 许可证 - 详情请见 [LICENSE](LICENSE) 文件。

### 🤝 贡献

欢迎贡献! 请随时提交 Pull Request。

---

**Made with ❤️ by hcd233**
