# AGENTS.md - AI Coding Assistant Guidelines

## Project Overview

**Aris Mem API** - A RESTful API service built with Go 1.25.1 using Clean Architecture patterns.

### Tech Stack

- **Language**: Go 1.25.1 (CGO_ENABLED=0)
- **Web Framework**: Fiber v2 + Huma v2 (OpenAPI 3.1)
- **Database**: PostgreSQL (GORM v1.25)
- **Cache**: Redis (go-redis v9)
- **Object Storage**: MinIO / Tencent COS
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **CLI**: Cobra
- **Logging**: Zap + Lumberjack
- **JSON**: Sonic (bytedance)
- **AI/LLM**: CloudWeGo Eino with OpenAI

## Build & Development Commands

```bash
# Build the application
go build -o aris-mem-api

# Run the server
go run main.go server start
go run main.go server start --host 0.0.0.0 --port 8080

# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/service/...
go test ./internal/handler/...

# Run a single test function
go test -run TestFunctionName ./package/path

# Run tests with coverage
go test -cover ./...

# Download dependencies
go mod tidy

# Docker build
docker build -f docker/dockerfile -t aris-mem-api:latest .
```

## Architecture & Project Structure

```
.
├── cmd/                    # Cobra CLI commands
│   ├── server.go          # Server start/stop commands
│   ├── database.go        # Database CLI tools
│   └── client.go          # Client commands
├── internal/              # Private application code
│   ├── api/               # Fiber & Huma API setup
│   ├── handler/           # HTTP handlers (Controller layer)
│   ├── service/           # Business logic layer
│   ├── dto/               # Data Transfer Objects (Request/Response)
│   ├── router/            # Route definitions
│   ├── middleware/        # Fiber middlewares
│   ├── common/            # Shared constants, enums, models
│   │   ├── constant/      # Constants (CtxKey*, error codes)
│   │   ├── enum/          # Enumerations (Status, Type)
│   │   └── model/         # Domain models (Error, etc.)
│   ├── infrastructure/    # External services
│   │   ├── database/      # DB connections & DAOs
│   │   │   ├── dao/       # Data Access Objects
│   │   │   └── model/     # GORM models
│   │   ├── cache/         # Redis cache
│   │   ├── storage/       # Object storage
│   │   └── smtp/          # Email service
│   ├── util/              # Utility functions
│   ├── logger/            # Zap logger setup
│   ├── config/            # Configuration
│   ├── jwt/               # JWT utilities
│   ├── lock/              # Distributed locking
│   └── oauth2/            # OAuth2 implementations
├── docker/                # Docker configurations
├── script/                # Deployment scripts
└── env/                   # Environment configurations
```

### Layer Dependencies

```
Handler → Service → DAO → Database
   ↓         ↓        ↓
  DTO      DTO     Model
```

## Code Style Guidelines

### General Principles

1. **Follow standard Go conventions** (gofmt, golint)
2. **All comments in English** (except Chinese in doc comments)
3. **Use explicit error handling** - never ignore errors
4. **Prefer composition over inheritance**
5. **Keep functions small and focused**

### Naming Conventions

```go
// Packages: lowercase, single word
package handler
package article

// Interfaces: noun, ending with -er/-or or descriptive
type ArticleService interface{}
type Handler interface{}

// Structs: PascalCase for exported, camelCase for internal
type articleService struct{}  // internal
type Article struct{}        // exported

// Functions: PascalCase for exported, camelCase for internal
func NewArticleService() ArticleService {}
func handleRequest() {}

// Variables: camelCase
var articleCount int
var userID uint

// Constants: CamelCase or UPPER_SNAKE_CASE for exported
const CtxKeyUserID = "userID"
const DefaultMaxImageSize = 10 * 1024 * 1024

// Acronyms: Keep uppercase (HTTP, URL, ID, JSON, SSE)
HTTPResponse
CoverImageURL
UserID
```

### File Organization

```go
// 1. Package declaration with doc comment
// Package handler Article handlers
package handler

// 2. Imports grouped: stdlib, third-party, internal
import (
    "context"
    "errors"
    
    "github.com/google/uuid"
    "github.com/samber/lo"
    "go.uber.org/zap"
    
    "github.com/hcd233/aris-mem-api/internal/dto"
    "github.com/hcd233/aris-mem-api/internal/service"
)

// 3. Interface definitions (public API)
// ArticleHandler Article handler interface
//
//	author centonhuang
//	update 2026-01-29 14:00:00
type ArticleHandler interface {
    HandleCreateArticle(ctx context.Context, req *dto.CreateArticleReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
}

// 4. Struct definitions
// articleHandler Article handler implementation
type articleHandler struct {
    svc service.ArticleService
}

// 5. Constructor functions
// NewArticleHandler Create new article handler
//
//	return ArticleHandler
//	author centonhuang
//	update 2026-01-29 10:00:00
func NewArticleHandler() ArticleHandler {
    return &articleHandler{svc: service.NewArticleService()}
}

// 6. Method implementations
func (h *articleHandler) HandleCreateArticle(ctx context.Context, req *dto.CreateArticleReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
    return util.WrapHTTPResponse(h.svc.CreateArticle(ctx, req))
}
```

### Documentation Style

```go
// Function/Interface doc comments:
// FunctionName Brief description
//
//	Detailed description if needed
//	@param paramName description
//	@return description
//	author centonhuang
//	update YYYY-MM-DD HH:MM:SS

// Struct doc comments:
// StructName Description
//
//	author centonhuang
//	update YYYY-MM-DD HH:MM:SS

type Article struct {
    ID     uint   `json:"id" gorm:"column:id;primary_key" doc:"Article ID"`
    Title  string `json:"title" doc:"Title of the article"`
}
```

## Error Handling Patterns

### Custom Error Model

```go
// Errors are defined in internal/common/constant/error.go
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

### Service Layer Error Handling

```go
func (s *articleService) CreateArticle(ctx context.Context, req *dto.CreateArticleReq) (*dto.EmptyRsp, error) {
    rsp := &dto.EmptyRsp{}
    
    // Validate input
    if req == nil || req.Body == nil {
        rsp.Error = constant.ErrBadRequest
        return rsp, nil  // Return nil error, error is in response
    }
    
    // Business logic errors go in rsp.Error
    if !isValid {
        rsp.Error = constant.ErrDataExists
        return rsp, nil
    }
    
    // System errors return as Go error
    db := database.GetDBInstance(ctx)
    if err := db.Create(&article).Error; err != nil {
        logger.Error("[ArticleService] failed to create", zap.Error(err))
        return nil, err  // Return system error
    }
    
    return rsp, nil
}
```

### Handler Error Pattern

```go
// Handlers use util.WrapHTTPResponse to wrap service responses
func (h *articleHandler) HandleCreateArticle(ctx context.Context, req *dto.CreateArticleReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
    return util.WrapHTTPResponse(h.svc.CreateArticle(ctx, req))
}
```

## DTO Patterns

### Request/Response Structure

```go
// Request wrapper
type CreateArticleReq struct {
    Body *CreateArticleReqBody `json:"body" doc:"Request body"`
}

type CreateArticleReqBody struct {
    Title      string `json:"title" doc:"Title of the article"`
    Content    string `json:"content" doc:"Content of the article"`
    CoverImage string `json:"coverImage" maxLength:"15000000" doc:"Cover image base64"`
}

// Response wrapper using generics
type HTTPResponse[BodyT any] struct {
    Body BodyT `json:"data"`
}

type EmptyRsp struct {
    Error *model.Error `json:"error,omitempty"`
}

// Listed items use List prefix
type ListArticlesRsp struct {
    Items []*ListedArticle `json:"items" doc:"List of articles"`
    Total int64            `json:"total" doc:"Total count"`
}

// Detailed view uses Detailed prefix
type GetArticleRsp struct {
    Article *DetailedArticle `json:"article" doc:"Article details"`
}
```

## Database Patterns

### GORM Model Structure

```go
// BaseModel contains common fields
type BaseModel struct {
    CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
    UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// Entity models
type Article struct {
    BaseModel
    DeletedAt   int64              `json:"deleted_at" gorm:"column:deleted_at;uniqueIndex:uidx_slug_deleted_at"`
    ID          uint               `json:"id" gorm:"column:id;primary_key;auto_increment"`
    UserID      uint               `json:"user_id" gorm:"column:user_id;not null"`
    Title       string             `json:"title" gorm:"column:title;not null"`
    Slug        string             `json:"slug" gorm:"column:slug;not null;uniqueIndex:uidx_slug_deleted_at"`
    Content     string             `json:"content" gorm:"column:content;not null"`
    Status      enum.ArticleStatus `json:"status" gorm:"column:status;not null"`
    // ... fields
}

// TableName override if needed
func (Article) TableName() string {
    return "article"
}
```

### DAO Pattern

```go
// DAO is singleton per entity
type ArticleDAO struct{}

var (
    articleDAO     *ArticleDAO
    articleDAOOnce sync.Once
)

// GetArticleDAO returns singleton instance
func GetArticleDAO() *ArticleDAO {
    articleDAOOnce.Do(func() {
        articleDAO = &ArticleDAO{}
    })
    return articleDAO
}

// DAO methods
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
    // Build query dynamically
    if query.ID != nil {
        db = db.Where("id = ?", *query.ID)
    }
    // ... more conditions
    var article Article
    if err := db.First(&article).Error; err != nil {
        return nil, err
    }
    return &article, nil
}
```

## Middleware Patterns

```go
// Middleware returns fiber.Handler
//
//	@author centonhuang
//	@update YYYY-MM-DD HH:MM:SS
func MiddlewareName() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Pre-processing
        
        err := c.Next()
        
        // Post-processing
        
        return err
    }
}

// Configuration-based middleware
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

## Context Values

```go
// Context keys defined in internal/common/constant/ctx.go
const (
    CtxKeyUserID     = "userID"
    CtxKeyUserName   = "userName"
    CtxKeyPermission = "permission"
    CtxKeyTraceID    = "traceID"
    CtxKeyLimiter    = "limiter"
)

// Usage in services
userID := ctx.Value(constant.CtxKeyUserID).(uint)
traceID := ctx.Value(constant.CtxKeyTraceID).(string)
```

## Logger Usage

```go
// Get logger with context
logger := logger.WithCtx(ctx)
logger := logger.WithFCtx(c) // Fiber context

// Log levels
logger.Debug("[Service] debug message", zap.String("key", value))
logger.Info("[Service] info message", zap.Uint("userID", userID))
logger.Warn("[Service] warning", zap.Error(err))
logger.Error("[Service] error occurred", zap.Error(err), zap.Stack("stack"))

// Naming convention: [ComponentName] message
```

## Testing Guidelines

```go
// Test files: *_test.go in same package
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

// Table-driven tests
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

## Common Utilities

```go
// String utilities
util.GenerateSlug(title string) string           // Generate URL-friendly slug
util.ExtractTags(content string) []string        // Extract #tags from content

// Image utilities
util.DecodeBase64OrDataURL(data string) ([]byte, string, error)
util.ConvertImageToJPEG(data []byte, mimeType string) ([]byte, error)

// Response wrapper
util.WrapHTTPResponse[rspT any](rsp rspT, err error) (*dto.HTTPResponse[rspT], error)

// Map utilities
util.MergeMaps(maps ...map[string]any) map[string]any
```

## Key Conventions Summary

1. **Always use struct tags**: `json:"fieldName" doc:"Description"`
2. **Pointer receivers** for methods that modify state
3. **Value receivers** for methods that only read
4. **Context first** parameter for all service/handler methods
5. **Return response + error**, not just error
6. **Use samber/lo** for functional programming helpers
7. **Use sonic** for JSON operations (fastest in Go)
8. **Never use init()** for business logic (only for config setup)
9. **Prefer dependency injection** via constructors
10. **Check all errors** - use `lo.Must0()` or `lo.Must1()` only when safe

## Docker Deployment

```bash
# Build
docker build -f docker/dockerfile -t aris-mem-api:latest .

# Run with env file
docker run -d -p 8080:8080 \
  --env-file env/api.env \
  --name aris-mem-api \
  aris-mem-api:latest \
  /app/aris-mem-api server start --host 0.0.0.0 --port 8080
```
