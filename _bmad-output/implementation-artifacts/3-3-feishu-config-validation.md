# Story 3.3: 飞书配置验证功能

Status: ready-for-dev

## Story

As a 用户,
I want 验证飞书配置是否有效,
so that 确保配置正确后再部署实例。

## Acceptance Criteria

1. **AC1: 验证配置 API 端点**
   - **Given** 用户已认证且有飞书配置
   - **When** 发送 POST 请求到 /v1/feishu-configs/validate
   - **Then** 后端调用飞书 API 验证凭证有效性
   - **And** 返回验证结果（成功/失败）

2. **AC2: 验证成功响应**
   - **Given** 飞书配置有效
   - **When** 执行验证
   - **Then** 返回验证成功状态
   - **And** 显示"配置有效"提示
   - **And** 更新配置状态为 active

3. **AC3: 验证失败响应**
   - **Given** 飞书配置无效
   - **When** 执行验证
   - **Then** 返回验证失败状态
   - **And** 显示具体错误原因（如 App ID 不存在、Secret 错误、权限不足等）
   - **And** 更新配置状态为 inactive

4. **AC4: 验证结果缓存**
   - **Given** 验证成功
   - **When** 5 分钟内再次验证
   - **Then** 返回缓存的验证结果
   - **And** 不重复调用飞书 API

5. **AC5: 保存配置时自动验证**
   - **Given** 用户保存飞书配置
   - **When** 调用创建或更新 API
   - **Then** 自动执行配置验证
   - **And** 验证失败时返回警告信息
   - **And** 配置仍然保存成功

6. **AC6: 飞书 SDK 集成**
   - **Given** 需要调用飞书 API
   - **When** 实现验证功能
   - **Then** 使用飞书 Go SDK
   - **And** 调用获取 access_token 接口验证凭证
   - **And** 处理飞书 API 错误响应

## Tasks / Subtasks

- [ ] Task 1: 创建飞书客户端封装 (AC: 6)
  - [ ] 1.1 创建 `internal/infrastructure/feishu/client.go`
  - [ ] 1.2 实现 NewClient(appID, appSecret) 构造函数
  - [ ] 1.3 实现 GetAccessToken() 方法
  - [ ] 1.4 实现错误处理和重试机制
  - [ ] 1.5 添加飞书 SDK 依赖 (github.com/larksuite/oapi-sdk-go/v3)

- [ ] Task 2: 实现验证服务 (AC: 1, 2, 3)
  - [ ] 2.1 创建 `internal/service/feishu_config_service.go`
  - [ ] 2.2 实现 ValidateConfig(config) 方法
  - [ ] 2.3 解析飞书 API 错误并返回友好提示
  - [ ] 2.4 更新配置状态（active/inactive）

- [ ] Task 3: 实现验证结果缓存 (AC: 4)
  - [ ] 3.1 创建 `pkg/cache/memory_cache.go`
  - [ ] 3.2 实现 Set(key, value, ttl) 方法
  - [ ] 3.3 实现 Get(key) 方法
  - [ ] 3.4 设置缓存 TTL 为 5 分钟
  - [ ] 3.5 缓存 key 格式: validate:{tenantID}

- [ ] Task 4: 实现验证 API Handler (AC: 1)
  - [ ] 4.1 在 `internal/api/feishu_config/handler.go` 添加 Validate 方法
  - [ ] 4.2 注册路由 POST /v1/feishu-configs/validate
  - [ ] 4.3 返回验证结果响应

- [ ] Task 5: 集成自动验证到保存流程 (AC: 5)
  - [ ] 5.1 修改 Create Handler，保存后触发异步验证
  - [ ] 5.2 修改 Update Handler，更新后触发异步验证
  - [ ] 5.3 验证结果通过 WebSocket 或轮询获取

- [ ] Task 6: 编写单元测试 (AC: 1-6)
  - [ ] 6.1 编写 `feishu_client_test.go` 测试飞书 SDK 调用
  - [ ] 6.2 编写 `validation_service_test.go` 测试验证逻辑
  - [ ] 6.3 编写 `cache_test.go` 测试缓存功能
  - [ ] 6.4 编写 `handler_test.go` 测试验证 API

## Dev Notes

### 架构模式与约束

**必须遵循的 Clean Architecture 原则 [Source: architecture.md]:**

1. **依赖方向**:
   - `service` 层依赖 `domain` 和 `infrastructure/feishu`
   - `api` 层依赖 `service` 层
   - 飞书客户端属于基础设施层

2. **错误处理**:
   - 飞书 API 错误需要转换为友好的用户提示
   - 错误信息需要支持中文

3. **缓存策略**:
   - 使用内存缓存（可选 Redis）
   - 缓存 TTL 5 分钟

### 现有项目状态

**Story 3.1 和 3.2 已完成:**

```
backend/
├── internal/
│   ├── api/
│   │   └── feishu_config/
│   │       └── handler.go          # ✅ 已创建
│   ├── domain/
│   │   └── config/
│   │       └── feishu_config.go    # ✅ 已创建
│   ├── infrastructure/
│   │   └── database/               # ✅ 已创建
│   └── repository/
│       └── feishu_config.go        # ✅ 已创建
├── pkg/
│   ├── encryption/
│   │   └── aes.go                  # ✅ 已创建
│   └── response/
│       └── response.go             # ✅ 已创建
```

**需要新增的组件:**

```
backend/
├── internal/
│   ├── infrastructure/
│   │   └── feishu/
│   │       ├── client.go           # 新增
│   │       └── client_test.go      # 新增
│   └── service/
│       ├── feishu_config_service.go  # 新增
│       └── feishu_config_service_test.go  # 新增
├── pkg/
│   └── cache/
│       ├── memory_cache.go         # 新增
│       └── memory_cache_test.go    # 新增
```

### 技术栈要求

**飞书 Go SDK:**

```go
// go.mod 添加依赖
require github.com/larksuite/oapi-sdk-go/v3 v3.0.0
```

**飞书客户端实现:**

```go
// internal/infrastructure/feishu/client.go
package feishu

import (
    "context"
    "fmt"
    lark "github.com/larksuite/oapi-sdk-go/v3"
    "github.com/larksuite/oapi-sdk-go/v3/auth"
)

type Client struct {
    appID     string
    appSecret string
    client    *lark.Client
}

func NewClient(appID, appSecret string) *Client {
    return &Client{
        appID:     appID,
        appSecret: appSecret,
        client: lark.NewClient(appID, appSecret,
            lark.WithLogLevel(lark.LogLevelDebug),
        ),
    }
}

// ValidateConfig 验证配置有效性
func (c *Client) ValidateConfig(ctx context.Context) error {
    // 获取 tenant_access_token 来验证凭证有效性
    resp, err := c.client.Auth.TenantAccessToken.InternalGetTenantAccessToken(ctx,
        &auth.InternalGetTenantAccessTokenReq{
            AppID:     c.appID,
            AppSecret: c.appSecret,
        },
    )
    if err != nil {
        return fmt.Errorf("飞书 API 调用失败: %w", err)
    }

    if resp == nil || resp.TenantAccessToken == "" {
        return fmt.Errorf("获取访问令牌失败")
    }

    return nil
}

// GetTenantAccessToken 获取租户访问令牌
func (c *Client) GetTenantAccessToken(ctx context.Context) (string, error) {
    resp, err := c.client.Auth.TenantAccessToken.InternalGetTenantAccessToken(ctx,
        &auth.InternalGetTenantAccessTokenReq{
            AppID:     c.appID,
            AppSecret: c.appSecret,
        },
    )
    if err != nil {
        return "", err
    }
    return resp.TenantAccessToken, nil
}
```

### 验证服务实现

**服务层:**

```go
// internal/service/feishu_config_service.go
package service

import (
    "context"
    "fmt"
    "github.com/gowa/saas-openclaw/backend/internal/domain/config"
    "github.com/gowa/saas-openclaw/backend/internal/infrastructure/feishu"
    "github.com/gowa/saas-openclaw/backend/internal/repository"
    "github.com/gowa/saas-openclaw/backend/pkg/cache"
)

type FeishuConfigService struct {
    repo  *repository.FeishuConfigRepository
    cache cache.Cache
}

func NewFeishuConfigService(repo *repository.FeishuConfigRepository, cache cache.Cache) *FeishuConfigService {
    return &FeishuConfigService{repo: repo, cache: cache}
}

// ValidateConfig 验证飞书配置
func (s *FeishuConfigService) ValidateConfig(ctx context.Context, tenantID string) (*ValidationResult, error) {
    // 1. 检查缓存
    cacheKey := fmt.Sprintf("validate:%s", tenantID)
    if cached, ok := s.cache.Get(cacheKey); ok {
        return cached.(*ValidationResult), nil
    }

    // 2. 获取配置
    cfg, err := s.repo.GetByTenantID(tenantID)
    if err != nil {
        return nil, fmt.Errorf("获取配置失败: %w", err)
    }

    // 3. 调用飞书 API 验证
    client := feishu.NewClient(cfg.AppID, cfg.AppSecret)
    err = client.ValidateConfig(ctx)

    result := &ValidationResult{
        TenantID: tenantID,
    }

    if err != nil {
        result.Success = false
        result.ErrorMessage = s.parseFeishuError(err)
        cfg.Status = config.StatusInactive
    } else {
        result.Success = true
        result.Message = "配置有效"
        cfg.Status = config.StatusActive
    }

    // 4. 更新状态
    if err := s.repo.Update(cfg); err != nil {
        // 记录日志但不影响验证结果
    }

    // 5. 缓存结果
    s.cache.Set(cacheKey, result, 5*time.Minute)

    return result, nil
}

// parseFeishuError 解析飞书错误为友好提示
func (s *FeishuConfigService) parseFeishuError(err error) string {
    // 根据飞书错误码返回友好提示
    // 10003: App ID 不存在
    // 10004: App Secret 错误
    // 10009: 权限不足
    return "配置验证失败: " + err.Error()
}

type ValidationResult struct {
    TenantID     string `json:"tenantId"`
    Success      bool   `json:"success"`
    Message      string `json:"message,omitempty"`
    ErrorMessage string `json:"errorMessage,omitempty"`
}
```

### 内存缓存实现

**缓存工具:**

```go
// pkg/cache/memory_cache.go
package cache

import (
    "sync"
    "time"
)

type item struct {
    value     interface{}
    expiredAt time.Time
}

type MemoryCache struct {
    items map[string]*item
    mu    sync.RWMutex
}

func NewMemoryCache() *MemoryCache {
    return &MemoryCache{
        items: make(map[string]*item),
    }
}

func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.items[key] = &item{
        value:     value,
        expiredAt: time.Now().Add(ttl),
    }
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    item, ok := c.items[key]
    if !ok || time.Now().After(item.expiredAt) {
        return nil, false
    }

    return item.value, true
}

func (c *MemoryCache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.items, key)
}

// Cleanup 清理过期缓存
func (c *MemoryCache) Cleanup() {
    c.mu.Lock()
    defer c.mu.Unlock()

    now := time.Now()
    for key, item := range c.items {
        if now.After(item.expiredAt) {
            delete(c.items, key)
        }
    }
}
```

### API Handler 实现

**验证 API:**

```go
// internal/api/feishu_config/handler.go 新增方法

// Validate 验证飞书配置
// POST /v1/feishu-configs/validate
func (h *Handler) Validate(c *gin.Context) {
    tenantID, _ := c.Get("tenantID")

    result, err := h.service.ValidateConfig(c.Request.Context(), tenantID.(string))
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "VALIDATION_ERROR", err.Error())
        return
    }

    response.Success(c, result)
}
```

**路由注册:**

```go
// 更新路由注册
feishuConfigs := v1.Group("/feishu-configs")
feishuConfigs.Use(middleware.Auth())
{
    feishuConfigs.POST("", feishuConfigHandler.Create)
    feishuConfigs.GET("", feishuConfigHandler.Get)
    feishuConfigs.PUT("", feishuConfigHandler.Update)
    feishuConfigs.DELETE("", feishuConfigHandler.Delete)
    feishuConfigs.POST("/validate", feishuConfigHandler.Validate)  // 新增
}
```

### 自动验证集成

**保存后自动验证:**

```go
// Create 方法修改
func (h *Handler) Create(c *gin.Context) {
    // ... 创建配置逻辑

    // 异步验证
    go func() {
        ctx := context.Background()
        h.service.ValidateConfig(ctx, cfg.TenantID)
    }()

    response.Created(c, toResponse(cfg))
}
```

### 飞书 API 错误码映射

| 飞书错误码 | HTTP 状态码 | 友好提示 |
|-----------|------------|---------|
| 10003 | 400 | App ID 不存在，请检查是否正确 |
| 10004 | 400 | App Secret 错误，请检查是否正确 |
| 10009 | 403 | 权限不足，请检查应用权限配置 |
| 10013 | 400 | 应用已被禁用，请联系管理员 |
| 10014 | 400 | 应用未发布，请先发布应用 |
| -1 | 500 | 飞书服务暂时不可用，请稍后重试 |

### 测试标准

**测试要求:**
- 测试框架: Go 原生 testing 包 + testify
- 测试覆盖率目标: >= 70%
- Mock 飞书 API 调用

**测试用例设计:**

| 测试场景 | 测试方法 | 预期结果 |
|---------|---------|---------|
| 验证成功 | Mock 成功响应 | success: true |
| 验证失败（App ID 错误） | Mock 错误响应 | success: false, 错误提示 |
| 缓存命中 | 5 分钟内重复验证 | 返回缓存结果 |
| 缓存过期 | 5 分钟后验证 | 重新调用 API |

### 前序 Story 的学习经验

**从 Story 3.2 (飞书配置 CRUD API) 获得的经验:**

1. **统一响应格式**: 使用 pkg/response 包
2. **认证中间件**: 从上下文获取租户 ID
3. **错误处理**: 返回友好的错误信息
4. **路由注册**: 遵循 RESTful 规范

### References

- [Source: architecture.md#Integration Requirements] - 飞书 SDK 集成要求
- [Source: prd.md#FR9] - 验证飞书应用配置有效性
- [Source: prd.md#NFR-I1] - 飞书 SDK 长连接稳定性
- [Source: epics.md#Story 3.3] - 原始故事定义
- [Source: 3-2-feishu-config-crud-api.md] - CRUD API 实现参考
- [飞书开放平台文档](https://open.feishu.cn/document/server-docs/getting-started/getting-started-2)

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
