# AGENTS.md - AI 编码助手指南

## 项目概述

**Aris Mem API** - 基于 Go 1.25.1 和整洁架构模式构建的 RESTful API 服务。

**重要提示：AGENTS.md中必须使用中文。后续所有更新、修改和新功能开发都必须遵循此规定。**

### 技术栈

- **语言**: Go 1.25.1 (CGO_ENABLED=0)
- **Web 框架**: Fiber v2 + Huma v2 (OpenAPI 3.1)
- **数据库**: PostgreSQL (GORM v1.25)
- **缓存**: Redis (go-redis v9)
- **对象存储**: MinIO / 腾讯云 COS
- **认证**: JWT (golang-jwt/jwt/v5)
- **CLI**: Cobra
- **日志**: Zap + Lumberjack
- **JSON**: Sonic (bytedance)
- **AI/LLM**: CloudWeGo Eino + OpenAI

## 构建与开发命令

```bash
# 构建应用程序
go build -o aris-mem-api

# 运行服务器
go run main.go server start
go run main.go server start --host 0.0.0.0 --port 8080

# 运行所有测试
go test ./...

# 运行指定包的测试
go test ./internal/service/...
go test ./internal/handler/...

# 运行单个测试函数
go test -run TestFunctionName ./package/path

# 运行测试并生成覆盖率报告
go test -cover ./...

# 下载依赖
go mod tidy

# Docker 构建
docker build -f docker/dockerfile -t aris-mem-api:latest .
```

## 架构与项目结构

```
.
├── cmd/                    # Cobra CLI 命令
│   ├── server.go          # 服务器启动/停止命令
│   ├── database.go        # 数据库 CLI 工具
│   └── client.go          # 客户端命令
├── internal/              # 私有应用代码
│   ├── api/               # Fiber 和 Huma API 设置
│   ├── handler/           # HTTP 处理器（控制层）
│   ├── service/           # 业务逻辑层
│   ├── dto/               # 数据传输对象（请求/响应）
│   ├── router/            # 路由定义
│   ├── middleware/        # Fiber 中间件
│   ├── common/            # 共享常量、枚举、模型
│   │   ├── constant/      # 常量（CtxKey*、错误码）
│   │   ├── enum/          # 枚举（Status、Type）
│   │   └── model/         # 领域模型（Error 等）
│   ├── infrastructure/    # 外部服务
│   │   ├── database/      # 数据库连接和 DAO
│   │   │   ├── dao/       # 数据访问对象
│   │   │   └── model/     # GORM 模型
│   │   ├── cache/         # Redis 缓存
│   │   ├── storage/       # 对象存储
│   │   └── smtp/          # 邮件服务
│   ├── util/              # 工具函数
│   ├── logger/            # Zap 日志设置
│   ├── config/            # 配置
│   ├── jwt/               # JWT 工具
│   ├── lock/              # 分布式锁
│   └── oauth2/            # OAuth2 实现
├── docker/                # Docker 配置
├── script/                # 部署脚本
└── env/                   # 环境配置
```

### 分层依赖关系

```
Handler → Service → DAO → Database
   ↓         ↓        ↓
  DTO      DTO     Model
```

## 代码风格指南

### 一般原则

1. **遵循标准 Go 约定**（gofmt、golint）
2. **所有注释使用中文**（要求后续更新必须使用中文）
3. **显式错误处理** - 永不忽略错误
4. **优先使用组合而非继承**
5. **保持函数小而专注**

### 命名规范

```go
// 包名：小写，单个单词
package handler
package article

// 接口名：名词，以 -er/-or 结尾或描述性名称
type ArticleService interface{}
type Handler interface{}

// 结构体名：导出用大驼峰，内部用小驼峰
type articleService struct{}  // 内部使用
type Article struct{}        // 导出使用

// 函数名：导出用大驼峰，内部用小驼峰
func NewArticleService() ArticleService {}
func handleRequest() {}

// 变量名：小驼峰
var articleCount int
var userID uint

// 常量名：导出用驼峰或大写下划线
const CtxKeyUserID = "userID"
const DefaultMaxImageSize = 10 * 1024 * 1024

// 缩写保持大写（HTTP、URL、ID、JSON、SSE）
HTTPResponse
CoverImageURL
UserID
```

### 文件组织

```go
// 1. 包声明带文档注释
// Package handler 文章处理器
package handler

// 2. 导入分组：标准库、第三方库、内部包
import (
    "context"
    "errors"
    
    "github.com/google/uuid"
    "github.com/samber/lo"
    "go.uber.org/zap"
    
    "github.com/hcd233/aris-mem-api/internal/dto"
    "github.com/hcd233/aris-mem-api/internal/service"
)

// 3. 接口定义（公共 API）
// ArticleHandler 文章处理器接口
//
//	author centonhuang
//	update 2026-01-29 14:00:00
type ArticleHandler interface {
    HandleCreateArticle(ctx context.Context, req *dto.CreateArticleReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
}

// 4. 结构体定义
// articleHandler 文章处理器实现
type articleHandler struct {
    svc service.ArticleService
}

// 5. 构造函数
// NewArticleHandler 创建文章处理器
//
//	return ArticleHandler
//	author centonhuang
//	update 2026-01-29 10:00:00
func NewArticleHandler() ArticleHandler {
    return &articleHandler{svc: service.NewArticleService()}
}

// 6. 方法实现
func (h *articleHandler) HandleCreateArticle(ctx context.Context, req *dto.CreateArticleReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
    return util.WrapHTTPResponse(h.svc.CreateArticle(ctx, req))
}
```

### 文档风格

```go
// 函数/接口文档注释：
// FunctionName 简要描述
//
//	详细描述（如需要）
//	@param paramName 描述
//	@return 描述
//	author centonhuang
//	update YYYY-MM-DD HH:MM:SS

// 结构体文档注释：
// StructName 描述
//
//	author centonhuang
//	update YYYY-MM-DD HH:MM:SS

type Article struct {
    ID     uint   `json:"id" gorm:"column:id;primary_key" doc:"文章 ID"`
    Title  string `json:"title" doc:"文章标题"`
}
```

## 错误处理模式

### 自定义错误模型

```go
// 错误定义在 internal/common/constant/error.go
var (
    ErrInternalError     = model.NewError(10000, "InternalError")
    ErrUnauthorized      = model.NewError(10001, "Unauthorized")
    ErrNoPermission      = model.NewError(10002, "NoPermission")
    ErrDataNotExists     = model.NewError(10003, "DataNotExists")
    ErrDataExists        = model.NewError(10004, "DataExists")
    ErrTooManyRequests   = model.NewError(10005, "TooManyRequests")
    ErrBadRequest        = model.NewError(10006, "BadRequest")
    ErrInsufficientQuota = model.NewError(10007, "InsufficientQuota")
    ErrNoImplement       = model.NewError(10008, "NoImplement")
    ErrInvalidFile       = model.NewError(10009, "InvalidFile")
)
```

### Service 层错误处理

```go
func (s *articleService) CreateArticle(ctx context.Context, req *dto.CreateArticleReq) (*dto.EmptyRsp, error) {
    rsp := &dto.EmptyRsp{}
    
    // 验证输入
    if req == nil || req.Body == nil {
        rsp.Error = constant.ErrBadRequest
        return rsp, nil  // 返回 nil 错误，错误在响应中
    }
    
    // 业务逻辑错误放入 rsp.Error
    if !isValid {
        rsp.Error = constant.ErrDataExists
        return rsp, nil
    }
    
    // 系统错误作为 Go error 返回
    db := database.GetDBInstance(ctx)
    if err := db.Create(&article).Error; err != nil {
        logger.Error("[ArticleService] failed to create article", zap.Error(err))
        return nil, err  // 返回系统错误
    }
    
    return rsp, nil
}
```

### Handler 错误模式

```go
// Handler 使用 util.WrapHTTPResponse 包装 Service 响应
func (h *articleHandler) HandleCreateArticle(ctx context.Context, req *dto.CreateArticleReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
    return util.WrapHTTPResponse(h.svc.CreateArticle(ctx, req))
}
```

## DTO 模式

### 请求/响应结构

```go
// 请求包装器
type CreateArticleReq struct {
    Body *CreateArticleReqBody `json:"body" doc:"请求体"`
}

type CreateArticleReqBody struct {
    Title      string `json:"title" doc:"文章标题"`
    Content    string `json:"content" doc:"文章内容"`
    CoverImage string `json:"coverImage" maxLength:"15000000" doc:"封面图片 base64"`
}

// 响应包装器使用泛型
type HTTPResponse[BodyT any] struct {
    Body BodyT `json:"data"`
}

type EmptyRsp struct {
    Error *model.Error `json:"error,omitempty"`
}

// 列表项使用 List 前缀
type ListArticlesRsp struct {
    Items []*ListedArticle `json:"items" doc:"文章列表"`
    Total int64            `json:"total" doc:"总数"`
}

// 详情视图使用 Detailed 前缀
type GetArticleRsp struct {
    Article *DetailedArticle `json:"article" doc:"文章详情"`
}
```

## 数据库模式

### GORM 模型结构

```go
// BaseModel 包含公共字段
type BaseModel struct {
    CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
    UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// 实体模型
type Article struct {
    BaseModel
    DeletedAt   int64              `json:"deleted_at" gorm:"column:deleted_at;uniqueIndex:uidx_slug_deleted_at"`
    ID          uint               `json:"id" gorm:"column:id;primary_key;auto_increment"`
    UserID      uint               `json:"user_id" gorm:"column:user_id;not null"`
    Title       string             `json:"title" gorm:"column:title;not null"`
    Slug        string             `json:"slug" gorm:"column:slug;not null;uniqueIndex:uidx_slug_deleted_at"`
    Content     string             `json:"content" gorm:"column:content;not null"`
    Status      enum.ArticleStatus `json:"status" gorm:"column:status;not null"`
    // ... 字段
}

// 如需可覆盖表名
func (Article) TableName() string {
    return "article"
}
```

### DAO 模式

```go
// DAO 每实体单例
type ArticleDAO struct{}

var (
    articleDAO     *ArticleDAO
    articleDAOOnce sync.Once
)

// GetArticleDAO 返回单例实例
func GetArticleDAO() *ArticleDAO {
    articleDAOOnce.Do(func() {
        articleDAO = &ArticleDAO{}
    })
    return articleDAO
}

// DAO 方法
type ArticleQuery struct {
    ID       *uint
    Slug     *string
    UserID   *uint
    Status   *enum.ArticleStatus
    Page     int
    PageSize int
}

func (dao *ArticleDAO) GetByQuery(ctx context.Context, query *ArticleQuery) (*Article, error) {
    db := database.GetDBInstance(ctx)
    // 动态构建查询
    if query.ID != nil {
        db = db.Where("id = ?", *query.ID)
    }
    // ... 更多条件
    var article Article
    if err := db.First(&article).Error; err != nil {
        return nil, err
    }
    return &article, nil
}
```

## 中间件模式

```go
// 中间件返回 fiber.Handler
//
//	@author centonhuang
//	@update YYYY-MM-DD HH:MM:SS
func MiddlewareName() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // 预处理
        
        err := c.Next()
        
        // 后处理
        
        return err
    }
}

// 基于配置的中间件
func RecoverMiddleware() fiber.Handler {
    return recover.New(recover.Config{
        EnableStackTrace: true,
        StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
            logger.WithFCtx(c).Error("[Panic Recovery]",
                zap.Any("error", e),
                zap.ByteString("stack", debug.Stack()))
        },
    })
}
```

## Context 值

```go
// Context 键定义在 internal/common/constant/ctx.go
const (
    CtxKeyUserID     = "userID"
    CtxKeyUserName   = "userName"
    CtxKeyPermission = "permission"
    CtxKeyTraceID    = "traceID"
    CtxKeyLimiter    = "limiter"
)

// 在 Service 中使用
userID := ctx.Value(constant.CtxKeyUserID).(uint)
traceID := ctx.Value(constant.CtxKeyTraceID).(string)
```

## 日志使用

```go
// 获取带 context 的日志
logger := logger.WithCtx(ctx)
logger := logger.WithFCtx(c) // Fiber context

// 日志级别
logger.Debug("[Service] Debug", zap.String("key", value))
logger.Info("[Service] Info", zap.Uint("userID", userID))
logger.Warn("[Service] Warn", zap.Error(err))
logger.Error("[Service] Error", zap.Error(err), zap.Stack("stack"))

// 命名约定：[组件名] 消息
```

## 测试指南

```go
// 测试文件：*_test.go 在同一包中
package service

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func TestArticleService_CreateArticle(t *testing.T) {
    // Arrange
    svc := NewArticleService()
    ctx := context.Background()
    req := &dto.CreateArticleReq{...}
    
    // Act
    rsp, err := svc.CreateArticle(ctx, req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, rsp)
    assert.Nil(t, rsp.Error)
}

// 表驱动测试
func TestArticleService_ValidateArticle(t *testing.T) {
    tests := []struct {
        name    string
        article *Article
        wantErr bool
    }{
        {"valid", &Article{Title: "Test"}, false},
        {"empty title", &Article{Title: ""}, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateArticle(tt.article)
            assert.Equal(t, tt.wantErr, err != nil)
        })
    }
}
```

## 常用工具

```go
// 字符串工具
util.GenerateSlug(title string) string           // 生成 URL 友好 slug
util.ExtractTags(content string) []string        // 从内容提取 #标签

// 图片工具
util.DecodeBase64OrDataURL(data string) ([]byte, string, error)
util.ConvertImageToJPEG(data []byte, mimeType string) ([]byte, error)

// 响应包装器
util.WrapHTTPResponse[rspT any](rsp rspT, err error) (*dto.HTTPResponse[rspT], error)

// Map 工具
util.MergeMaps(maps ...map[string]any) map[string]any
```

## 关键约定总结

1. **始终使用结构体标签**：`json:"fieldName" doc:"描述"`
2. **修改状态的方法使用指针接收器**
3. **只读方法使用值接收器**
4. **所有 Service/Handler 方法的第一个参数是 Context**
5. **返回响应 + 错误，而非仅返回错误**
6. **使用 samber/lo** 进行函数式编程
7. **使用 sonic** 进行 JSON 操作（Go 中最快）
8. **永远不要使用 init()** 处理业务逻辑（仅用于配置设置）
9. **优先通过构造函数进行依赖注入**
10. **检查所有错误** - 仅在安全时使用 `lo.Must0()` 或 `lo.Must1()`

## Docker 部署

```bash
# 构建
docker build -f docker/dockerfile -t aris-mem-api:latest .

# 使用环境文件运行
docker run -d -p 8080:8080 \
  --env-file env/api.env \
  --name aris-mem-api \
  aris-mem-api:latest \
  /app/aris-mem-api server start --host 0.0.0.0 --port 8080
```

---

## Bug修复日志

每次解决bug后必须在此记录，格式如下：

### [2026-01-30] Undo点赞/收藏时计数未减少
- **现象**: 在撤销(unlike/unsave)文章时，文章的点赞数和收藏数没有减少
- **原因**: 
  1. `Update` 方法使用 `lo.Filter` 过滤掉值为零值的字段
  2. 当 `likes` 或 `saves` 从 1 减到 0 时，值 0 被当作零值过滤掉，导致不更新
  3. 使用 `uint` 类型无法区分"不更新"和"更新为0"两种情况
- **解决方案**: 
  1. 在 `action.go` 中使用 `*uint` 指针类型传递更新值
  2. nil 指针表示不更新，指向 0 的指针表示更新为 0
  3. 保留 `Update` 方法的零值过滤逻辑，通过指针类型正确区分
- **文件位置**: 
  - `internal/service/action.go` (Do/Undo 方法中使用 *uint)
  - `internal/infrastructure/database/dao/base.go` (Update 方法)