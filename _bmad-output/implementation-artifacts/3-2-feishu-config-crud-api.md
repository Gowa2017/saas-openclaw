# Story 3.2: 飞书配置 CRUD API

Status: ready-for-dev

## Story

As a 后端开发者,
I want 实现飞书配置的增删改查 API,
so that 前端可以管理飞书配置。

## Acceptance Criteria

1. **AC1: POST /v1/feishu-configs - 创建飞书配置**
   - **Given** 用户已认证
   - **When** 发送 POST 请求到 /v1/feishu-configs
   - **Then** 验证请求体格式正确
   - **And** AppID 格式校验（以 cli_ 开头）
   - **And** 创建飞书配置记录
   - **And** 返回 201 状态码和创建的配置信息
   - **And** 返回统一的 API 响应格式

2. **AC2: GET /v1/feishu-configs - 获取当前用户的飞书配置**
   - **Given** 用户已认证且有飞书配置
   - **When** 发送 GET 请求到 /v1/feishu-configs
   - **Then** 返回当前租户的飞书配置
   - **And** AppSecret 字段不返回
   - **And** 返回 200 状态码
   - **And** 无配置时返回 404 状态码

3. **AC3: PUT /v1/feishu-configs - 更新飞书配置**
   - **Given** 用户已认证且有飞书配置
   - **When** 发送 PUT 请求到 /v1/feishu-configs
   - **Then** 更新飞书配置记录
   - **And** AppSecret 更新时重新加密
   - **And** 返回 200 状态码和更新后的配置信息
   - **And** 无配置时返回 404 状态码

4. **AC4: DELETE /v1/feishu-configs - 删除飞书配置**
   - **Given** 用户已认证且有飞书配置
   - **When** 发送 DELETE 请求到 /v1/feishu-configs
   - **Then** 删除当前租户的飞书配置
   - **And** 返回 204 状态码
   - **And** 无配置时返回 404 状态码

5. **AC5: 所有 API 需要认证**
   - **Given** API 端点已实现
   - **When** 未携带有效认证信息请求
   - **Then** 返回 401 状态码
   - **And** 返回错误信息 "未授权访问"

6. **AC6: 统一 API 响应格式**
   - **Given** 任意 API 请求
   - **When** 返回响应
   - **Then** 使用统一包装器格式
   - **And** 成功响应: `{ data: {...}, error: null, meta: {...} }`
   - **And** 错误响应: `{ data: null, error: {...}, meta: {...} }`

## Tasks / Subtasks

- [ ] Task 1: 创建 API Handler 结构 (AC: 1-6)
  - [ ] 1.1 创建 `internal/api/feishu_config/handler.go`
  - [ ] 1.2 定义 FeishuConfigHandler 结构体
  - [ ] 1.3 实现 NewFeishuConfigHandler 构造函数
  - [ ] 1.4 注入 FeishuConfigRepository 依赖

- [ ] Task 2: 实现请求/响应模型 (AC: 1, 3, 6)
  - [ ] 2.1 创建 `internal/api/feishu_config/dto.go`
  - [ ] 2.2 定义 CreateFeishuConfigRequest 结构体
  - [ ] 2.3 定义 UpdateFeishuConfigRequest 结构体
  - [ ] 2.4 定义 FeishuConfigResponse 结构体
  - [ ] 2.5 添加验证标签（validate:"required"）

- [ ] Task 3: 实现 API Handler 方法 (AC: 1-4)
  - [ ] 3.1 实现 Create 方法 (POST)
  - [ ] 3.2 实现 Get 方法 (GET)
  - [ ] 3.3 实现 Update 方法 (PUT)
  - [ ] 3.4 实现 Delete 方法 (DELETE)
  - [ ] 3.5 从上下文获取租户 ID（认证中间件注入）

- [ ] Task 4: 注册路由 (AC: 1-5)
  - [ ] 4.1 在 `internal/api/router.go` 注册路由组
  - [ ] 4.2 应用认证中间件
  - [ ] 4.3 路由路径: /v1/feishu-configs

- [ ] Task 5: 实现统一响应工具 (AC: 6)
  - [ ] 5.1 创建 `pkg/response/response.go`
  - [ ] 5.2 实现 Success(c, data) 函数
  - [ ] 5.3 实现 Error(c, code, message) 函数
  - [ ] 5.4 实现 Created(c, data) 函数
  - [ ] 5.5 实现 NoContent(c) 函数

- [ ] Task 6: 编写单元测试 (AC: 1-6)
  - [ ] 6.1 编写 `handler_test.go` 测试 CRUD 操作
  - [ ] 6.2 编写请求验证测试
  - [ ] 6.3 编写错误处理测试
  - [ ] 6.4 编写响应格式测试

## Dev Notes

### 架构模式与约束

**必须遵循的 Clean Architecture 原则 [Source: architecture.md]:**

1. **依赖方向**:
   - `api` 层依赖 `domain` 和 `repository`
   - Handler 通过接口调用 Repository，不直接依赖具体实现

2. **API 设计规范 [Source: architecture.md#API Naming Conventions]:**
   - REST 端点: 复数资源名 (例: `/feishu-configs`)
   - 版本化: `/v1/feishu-configs`
   - 标准方法: POST (创建), GET (读取), PUT (更新), DELETE (删除)

3. **API 响应格式 [Source: architecture.md#API Response Formats]:**
   - 统一包装器: `{ data: {...}, error: null, meta: {...} }`

### 现有项目状态

**Story 3.1 已完成的数据模型 [Source: 3-1-feishu-config-model.md]:**

```
backend/
├── internal/
│   ├── domain/
│   │   └── config/
│   │       └── feishu_config.go    # ✅ 已创建
│   └── repository/
│       └── feishu_config.go        # ✅ 已创建
├── pkg/
│   └── encryption/
│       └── aes.go                  # ✅ 已创建
```

**需要新增的 API 层:**

```
backend/
├── internal/
│   └── api/
│       └── feishu_config/          # 新建目录
│           ├── handler.go          # 新增
│           ├── dto.go              # 新增
│           └── handler_test.go     # 新增
├── pkg/
│   └── response/
│       └── response.go             # 新增
```

### 技术栈要求

**核心依赖 [Source: architecture.md]:**

| 依赖 | 用途 | 版本 |
|------|------|------|
| `github.com/gin-gonic/gin` | Web 框架 | v1.9+ |
| `github.com/go-playground/validator/v10` | 请求验证 | v10+ |

**API Handler 结构:**

```go
// internal/api/feishu_config/handler.go
package feishu_config

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gowa/saas-openclaw/backend/internal/domain/config"
    "github.com/gowa/saas-openclaw/backend/internal/repository"
    "github.com/gowa/saas-openclaw/backend/pkg/response"
)

type Handler struct {
    repo *repository.FeishuConfigRepository
}

func NewHandler(repo *repository.FeishuConfigRepository) *Handler {
    return &Handler{repo: repo}
}

// Create 创建飞书配置
// POST /v1/feishu-configs
func (h *Handler) Create(c *gin.Context) {
    // 1. 绑定请求体
    // 2. 验证请求
    // 3. 从上下文获取租户 ID
    // 4. 调用 repository 创建
    // 5. 返回响应
}

// Get 获取飞书配置
// GET /v1/feishu-configs
func (h *Handler) Get(c *gin.Context) {
    // 1. 从上下文获取租户 ID
    // 2. 调用 repository 查询
    // 3. 返回响应
}

// Update 更新飞书配置
// PUT /v1/feishu-configs
func (h *Handler) Update(c *gin.Context) {
    // 1. 绑定请求体
    // 2. 验证请求
    // 3. 从上下文获取租户 ID
    // 4. 调用 repository 更新
    // 5. 返回响应
}

// Delete 删除飞书配置
// DELETE /v1/feishu-configs
func (h *Handler) Delete(c *gin.Context) {
    // 1. 从上下文获取租户 ID
    // 2. 调用 repository 删除
    // 3. 返回响应
}
```

### 请求/响应模型设计

**CreateFeishuConfigRequest:**

```go
// internal/api/feishu_config/dto.go
package feishu_config

// CreateFeishuConfigRequest 创建飞书配置请求
type CreateFeishuConfigRequest struct {
    AppID     string `json:"appId" binding:"required,startswith=cli_"`
    AppSecret string `json:"appSecret" binding:"required,min=1"`
}

// UpdateFeishuConfigRequest 更新飞书配置请求
type UpdateFeishuConfigRequest struct {
    AppID     string `json:"appId" binding:"required,startswith=cli_"`
    AppSecret string `json:"appSecret" binding:"required,min=1"`
}

// FeishuConfigResponse 飞书配置响应
type FeishuConfigResponse struct {
    ID        string `json:"id"`
    TenantID  string `json:"tenantId"`
    AppID     string `json:"appId"`
    Status    string `json:"status"`
    CreatedAt string `json:"createdAt"`
    UpdatedAt string `json:"updatedAt"`
    // 注意: AppSecret 不返回
}
```

### 统一响应格式

**Success Response:**

```go
// pkg/response/response.go
package response

import (
    "time"
    "github.com/gin-gonic/gin"
)

type Meta struct {
    Timestamp string `json:"timestamp"`
}

type SuccessResponse struct {
    Data  interface{} `json:"data"`
    Error interface{} `json:"error"`
    Meta  Meta        `json:"meta"`
}

type ErrorResponse struct {
    Data  interface{} `json:"data"`
    Error ErrorDetail `json:"error"`
    Meta  Meta        `json:"meta"`
}

type ErrorDetail struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, SuccessResponse{
        Data:  data,
        Error: nil,
        Meta:  Meta{Timestamp: time.Now().UTC().Format(time.RFC3339)},
    })
}

// Created 返回创建成功响应
func Created(c *gin.Context, data interface{}) {
    c.JSON(http.StatusCreated, SuccessResponse{
        Data:  data,
        Error: nil,
        Meta:  Meta{Timestamp: time.Now().UTC().Format(time.RFC3339)},
    })
}

// NoContent 返回无内容响应
func NoContent(c *gin.Context) {
    c.Status(http.StatusNoContent)
}

// Error 返回错误响应
func Error(c *gin.Context, statusCode int, code, message string) {
    c.JSON(statusCode, ErrorResponse{
        Data: nil,
        Error: ErrorDetail{
            Code:    code,
            Message: message,
        },
        Meta: Meta{Timestamp: time.Now().UTC().Format(time.RFC3339)},
    })
}

// NotFound 返回 404 错误
func NotFound(c *gin.Context, message string) {
    Error(c, http.StatusNotFound, "NOT_FOUND", message)
}

// Unauthorized 返回 401 错误
func Unauthorized(c *gin.Context, message string) {
    Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// BadRequest 返回 400 错误
func BadRequest(c *gin.Context, message string) {
    Error(c, http.StatusBadRequest, "BAD_REQUEST", message)
}
```

### 路由注册

**Router 配置:**

```go
// internal/api/router.go
package api

import (
    "github.com/gin-gonic/gin"
    "github.com/gowa/saas-openclaw/backend/internal/api/feishu_config"
    "github.com/gowa/saas-openclaw/backend/internal/repository"
    "github.com/gowa/saas-openclaw/backend/pkg/middleware"
)

func SetupRouter(db *sqlx.DB) *gin.Engine {
    r := gin.Default()

    // API v1
    v1 := r.Group("/v1")

    // 飞书配置 API（需要认证）
    feishuConfigHandler := feishu_config.NewHandler(
        repository.NewFeishuConfigRepository(db),
    )

    feishuConfigs := v1.Group("/feishu-configs")
    feishuConfigs.Use(middleware.Auth()) // 认证中间件
    {
        feishuConfigs.POST("", feishuConfigHandler.Create)
        feishuConfigs.GET("", feishuConfigHandler.Get)
        feishuConfigs.PUT("", feishuConfigHandler.Update)
        feishuConfigs.DELETE("", feishuConfigHandler.Delete)
    }

    return r
}
```

### 认证中间件集成

**从上下文获取租户 ID:**

```go
// pkg/middleware/auth.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/gowa/saas-openclaw/backend/pkg/response"
)

func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从 Header 获取 Token
        // 2. 验证 Token
        // 3. 解析用户信息
        // 4. 将租户 ID 存入上下文

        // 临时实现：从 Header 获取
        tenantID := c.GetHeader("X-Tenant-ID")
        if tenantID == "" {
            response.Unauthorized(c, "未授权访问")
            c.Abort()
            return
        }

        c.Set("tenantID", tenantID)
        c.Next()
    }
}
```

**Handler 中获取租户 ID:**

```go
// 在 Handler 中获取租户 ID
func (h *Handler) Get(c *gin.Context) {
    tenantID, exists := c.Get("tenantID")
    if !exists {
        response.Unauthorized(c, "未授权访问")
        return
    }

    config, err := h.repo.GetByTenantID(tenantID.(string))
    // ...
}
```

### 测试标准

**测试要求:**
- 测试框架: Go 原生 testing 包 + testify
- 测试覆盖率目标: >= 70%
- 使用 httptest 进行 HTTP 测试

**测试用例设计:**

| 测试场景 | 测试方法 | 预期结果 |
|---------|---------|---------|
| 创建配置成功 | POST 有效请求 | 201, 返回配置信息 |
| 创建配置失败（无效 AppID） | POST cli_ 开头验证失败 | 400, 错误信息 |
| 获取配置成功 | GET 已有配置 | 200, 返回配置 |
| 获取配置失败（无配置） | GET 无配置 | 404, 错误信息 |
| 更新配置成功 | PUT 有效请求 | 200, 返回更新后配置 |
| 删除配置成功 | DELETE 已有配置 | 204 |
| 未认证请求 | 无 Token 请求 | 401, 错误信息 |

### API 文档

**Swagger/OpenAPI 注释:**

```go
// Create 创建飞书配置
// @Summary 创建飞书配置
// @Description 为当前租户创建飞书应用配置
// @Tags feishu-configs
// @Accept json
// @Produce json
// @Param X-Tenant-ID header string true "租户 ID"
// @Param request body CreateFeishuConfigRequest true "飞书配置"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /v1/feishu-configs [post]
func (h *Handler) Create(c *gin.Context) {
    // ...
}
```

### 错误码定义

| 错误码 | HTTP 状态码 | 说明 |
|-------|------------|------|
| `UNAUTHORIZED` | 401 | 未授权访问 |
| `NOT_FOUND` | 404 | 资源不存在 |
| `BAD_REQUEST` | 400 | 请求参数错误 |
| `INTERNAL_ERROR` | 500 | 服务器内部错误 |
| `VALIDATION_ERROR` | 400 | 数据验证失败 |
| `DUPLICATE_CONFIG` | 409 | 配置已存在 |

### 前序 Story 的学习经验

**从 Story 3.1 (飞书配置数据模型) 获得的经验:**

1. **加密存储**: AppSecret 使用 AES-256 加密
2. **命名约定**: 表名 snake_case，列名 PascalCase
3. **一租户一配置**: TenantID 唯一约束
4. **敏感字段保护**: AppSecret 不暴露给 API

### References

- [Source: architecture.md#API Naming Conventions] - API 命名约定
- [Source: architecture.md#API Response Formats] - 统一响应格式
- [Source: prd.md#FR5-FR9] - 飞书应用配置需求
- [Source: epics.md#Story 3.2] - 原始故事定义
- [Source: 3-1-feishu-config-model.md] - 飞书配置数据模型

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
