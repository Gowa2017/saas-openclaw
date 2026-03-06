# Story 2.3: 用户信息获取与租户识别

Status: done

## Story

As a 后端开发者,
I want 实现用户信息获取和租户识别功能,
so that 可以正确识别用户所属租户并实现多租户隔离。

## Acceptance Criteria

1. **AC1: 用户信息获取**
   - **Given** JWT Token 验证通过
   - **When** 调用业务平台用户信息接口
   - **Then** 获取用户 ID、姓名、邮箱、租户 ID 等信息

2. **AC2: 租户隔离实现**
   - **Given** 获取到用户的租户 ID
   - **When** 查询数据时
   - **Then** 根据租户 ID 实现数据隔离
   - **And** 用户只能访问自己租户的数据

3. **AC3: 本地用户记录同步**
   - **Given** 从业务平台获取用户信息
   - **When** 用户首次登录或信息变更
   - **Then** 自动创建本地用户记录（如不存在）
   - **And** 自动更新本地用户记录（如已存在）

4. **AC4: 统一 API 响应格式**
   - **Given** API 处理完成
   - **When** 返回响应
   - **Then** 使用统一的 API 响应格式
   - **And** 包含 data、error、meta 字段

## Tasks / Subtasks

- [x] Task 1: 创建用户同步服务 (AC: 3)
  - [x] 1.1 创建 `internal/service/user_sync.go`
  - [x] 1.2 实现 CreateOrUpdateUser 方法
  - [x] 1.3 处理并发创建场景（乐观锁）
  - [x] 1.4 添加同步日志记录

- [x] Task 2: 实现租户上下文管理 (AC: 2)
  - [x] 2.1 创建 `pkg/contextx/tenant.go` 定义租户上下文工具
  - [x] 2.2 实现 SetTenantContext 方法
  - [x] 2.3 实现 GetTenantContext 方法
  - [x] 2.4 添加租户 ID 验证

- [x] Task 3: 实现租户隔离中间件 (AC: 2)
  - [x] 3.1 创建 `pkg/middleware/tenant.go`
  - [x] 3.2 实现租户 ID 提取
  - [x] 3.3 实现数据查询自动过滤
  - [x] 3.4 添加隔离验证日志

- [x] Task 4: 创建统一响应工具 (AC: 4)
  - [x] 4.1 创建 `pkg/response/response.go`
  - [x] 4.2 实现 Success 响应函数
  - [x] 4.3 实现 Error 响应函数
  - [x] 4.4 实现 Paginated 响应函数

- [x] Task 5: 实现 Repository 租户隔离 (AC: 2)
  - [x] 5.1 修改 `TenantUserRepository` 添加租户过滤
  - [x] 5.2 实现 GetByTenantID 方法优化
  - [x] 5.3 添加查询审计日志

- [x] Task 6: 编写单元测试 (AC: 1-4)
  - [x] 6.1 编写 `user_sync_test.go` 测试用户同步
  - [x] 6.2 编写 `tenant_context_test.go` 测试上下文管理
  - [x] 6.3 编写 `tenant_middleware_test.go` 测试隔离中间件
  - [x] 6.4 编写 `response_test.go` 测试响应格式

## Dev Notes

### 架构模式与约束

**必须遵循的 Clean Architecture 原则 [Source: architecture.md]:**

1. **分层设计**:
   - `pkg/context/` - 上下文工具（无业务逻辑）
   - `pkg/middleware/` - HTTP 中间件
   - `internal/service/` - 业务服务层

2. **多租户隔离策略 [Source: architecture.md#Cross-Cutting Concerns]:**
   - Docker 网络隔离
   - PostgreSQL Database per Tenant
   - JWT 租户识别

3. **API 响应格式 [Source: architecture.md#Format Patterns]:**
   - 统一包装器: `{ data: {...}, error: null, meta: {...} }`

### 现有项目状态

**依赖 Story 2.1 和 2.2 已完成:**

```
backend/
├── internal/
│   ├── domain/
│   │   ├── user/
│   │   │   ├── tenant_user.go    # 租户用户实体
│   │   │   └── admin_user.go     # 管理员用户实体
│   │   └── tenant/
│   │       └── tenant.go         # 租户实体
│   ├── repository/
│   │   ├── tenant_user.go        # 租户用户仓库
│   │   └── admin_user.go         # 管理员用户仓库
│   └── infrastructure/
│       └── platform/
│           └── client.go         # 业务平台客户端
├── pkg/
│   ├── middleware/
│   │   └── auth.go               # 认证中间件
│   └── jwt/
│       └── jwt.go                # JWT 工具
```

### 用户同步服务设计

**服务实现:**

```go
// internal/service/user_sync.go
package service

import (
    "context"
    "github.com/gowa/saas-openclaw/backend/internal/domain/user"
    "github.com/gowa/saas-openclaw/backend/internal/repository"
)

type UserSyncService struct {
    userRepo   *repository.TenantUserRepository
    tenantRepo *repository.TenantRepository
}

type SyncResult struct {
    UserID   string
    TenantID string
    IsNew    bool
}

// SyncUser 同步用户信息到本地数据库
func (s *UserSyncService) SyncUser(ctx context.Context, userInfo *platform.UserInfo) (*SyncResult, error) {
    // 1. 检查租户是否存在，不存在则创建
    tenant, err := s.ensureTenant(ctx, userInfo.TenantID)
    if err != nil {
        return nil, err
    }

    // 2. 查找或创建用户
    existingUser, err := s.userRepo.GetByEmail(ctx, userInfo.Email)
    if err != nil && !errors.Is(err, sql.ErrNoRows) {
        return nil, err
    }

    if existingUser != nil {
        // 更新现有用户
        existingUser.Name = userInfo.Name
        existingUser.Role = user.Role(userInfo.Role)
        if err := s.userRepo.Update(ctx, existingUser); err != nil {
            return nil, err
        }
        return &SyncResult{UserID: existingUser.ID, TenantID: tenant.ID, IsNew: false}, nil
    }

    // 创建新用户
    newUser := &user.TenantUser{
        ID:       uuid.New().String(),
        TenantID: tenant.ID,
        Name:     userInfo.Name,
        Email:    userInfo.Email,
        Role:     user.Role(userInfo.Role),
    }
    if err := s.userRepo.Create(ctx, newUser); err != nil {
        return nil, err
    }
    return &SyncResult{UserID: newUser.ID, TenantID: tenant.ID, IsNew: true}, nil
}
```

### 租户上下文设计

**上下文工具:**

```go
// pkg/context/tenant.go
package context

import "context"

type tenantContextKey struct{}

type TenantContext struct {
    TenantID string
    UserID   string
    Email    string
    Role     string
}

// SetTenantContext 设置租户上下文
func SetTenantContext(ctx context.Context, tc *TenantContext) context.Context {
    return context.WithValue(ctx, tenantContextKey{}, tc)
}

// GetTenantContext 获取租户上下文
func GetTenantContext(ctx context.Context) (*TenantContext, bool) {
    tc, ok := ctx.Value(tenantContextKey{}).(*TenantContext)
    return tc, ok
}

// GetTenantID 获取租户 ID
func GetTenantID(ctx context.Context) (string, bool) {
    tc, ok := GetTenantContext(ctx)
    if !ok {
        return "", false
    }
    return tc.TenantID, true
}
```

### 租户隔离中间件设计

**中间件实现:**

```go
// pkg/middleware/tenant.go
package middleware

import (
    "net/http"
    "github.com/gowa/saas-openclaw/backend/pkg/context"
)

// TenantMiddleware 租户隔离中间件
func TenantMiddleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 从认证中间件注入的上下文获取用户信息
            userCtx, ok := middleware.GetUserContext(r.Context())
            if !ok {
                respondWithError(w, http.StatusUnauthorized, "user context not found")
                return
            }

            // 设置租户上下文
            tc := &context.TenantContext{
                TenantID: userCtx.TenantID,
                UserID:   userCtx.ID,
                Email:    userCtx.Email,
                Role:     userCtx.Role,
            }
            ctx := context.SetTenantContext(r.Context(), tc)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### 统一响应工具设计

**响应格式:**

```go
// pkg/response/response.go
package response

import (
    "encoding/json"
    "net/http"
)

type APIResponse struct {
    Data  interface{} `json:"data"`
    Error *APIError   `json:"error"`
    Meta  *Meta       `json:"meta,omitempty"`
}

type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

type Meta struct {
    Total    int64 `json:"total,omitempty"`
    Page     int   `json:"page,omitempty"`
    PageSize int   `json:"pageSize,omitempty"`
}

// Success 成功响应
func Success(w http.ResponseWriter, data interface{}) {
    respond(w, http.StatusOK, &APIResponse{
        Data:  data,
        Error: nil,
    })
}

// Error 错误响应
func Error(w http.ResponseWriter, status int, code, message string) {
    respond(w, status, &APIResponse{
        Data: nil,
        Error: &APIError{
            Code:    code,
            Message: message,
        },
    })
}

// Paginated 分页响应
func Paginated(w http.ResponseWriter, data interface{}, total int64, page, pageSize int) {
    respond(w, http.StatusOK, &APIResponse{
        Data: data,
        Meta: &Meta{
            Total:    total,
            Page:     page,
            PageSize: pageSize,
        },
    })
}

func respond(w http.ResponseWriter, status int, resp *APIResponse) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(resp)
}
```

### Repository 租户隔离增强

**租户过滤查询:**

```go
// internal/repository/tenant_user.go

// GetByTenantID 获取租户下所有用户（带分页）
func (r *TenantUserRepository) GetByTenantID(ctx context.Context, tenantID string, page, pageSize int) ([]*user.TenantUser, int64, error) {
    // 验证上下文中的租户 ID
    ctxTenantID, ok := context.GetTenantID(ctx)
    if !ok || ctxTenantID != tenantID {
        return nil, 0, errors.New("tenant isolation violation")
    }

    // 查询数据
    var users []*user.TenantUser
    query := `SELECT * FROM tenant_users WHERE "TenantID" = $1 ORDER BY "CreatedAt" DESC LIMIT $2 OFFSET $3`
    err := r.db.SelectContext(ctx, &users, query, tenantID, pageSize, (page-1)*pageSize)

    // 查询总数
    var total int64
    countQuery := `SELECT COUNT(*) FROM tenant_users WHERE "TenantID" = $1`
    r.db.GetContext(ctx, &total, countQuery, tenantID)

    return users, total, err
}
```

### 项目结构规范

**新增文件位置:**

```
backend/
├── pkg/
│   ├── context/
│   │   ├── tenant.go            # 租户上下文（新增）
│   │   └── tenant_test.go       # 测试文件（新增）
│   ├── middleware/
│   │   ├── tenant.go            # 租户隔离中间件（新增）
│   │   └── tenant_test.go       # 测试文件（新增）
│   └── response/
│       ├── response.go          # 响应工具（新增）
│       └── response_test.go     # 测试文件（新增）
├── internal/
│   └── service/
│       ├── user_sync.go         # 用户同步服务（新增）
│       └── user_sync_test.go    # 测试文件（新增）
```

### 测试标准

**测试要求:**
- 测试框架: Go 原生 testing 包 + testify
- 测试覆盖率目标: >= 70%
- 使用 mock 测试外部依赖

**测试用例设计:**

| 测试场景 | 测试文件 | 测试方法 |
|---------|---------|---------|
| 用户首次同步（创建） | `user_sync_test.go` | 单元测试 |
| 用户二次同步（更新） | `user_sync_test.go` | 单元测试 |
| 并发同步冲突处理 | `user_sync_test.go` | 单元测试 |
| 租户上下文设置 | `tenant_context_test.go` | 单元测试 |
| 租户上下文获取 | `tenant_context_test.go` | 单元测试 |
| 租户隔离验证 | `tenant_middleware_test.go` | 单元测试 |
| 成功响应格式 | `response_test.go` | 单元测试 |
| 错误响应格式 | `response_test.go` | 单元测试 |
| 分页响应格式 | `response_test.go` | 单元测试 |

### 安全注意事项

1. **租户隔离验证**: 每次数据查询必须验证租户 ID
2. **上下文注入**: 确保上下文来自可信的认证中间件
3. **日志审计**: 记录跨租户访问尝试
4. **错误信息**: 不暴露敏感的租户信息

### Project Structure Notes

**与前置 Story 的连续性:**

1. **依赖 Story 2.1**:
   - 使用 TenantUser 结构体
   - 使用 TenantUserRepository 进行数据访问

2. **依赖 Story 2.2**:
   - 使用认证中间件获取的用户信息
   - 使用业务平台客户端获取用户详情

3. **为后续 Story 提供基础**:
   - 租户隔离机制供所有业务模块使用
   - 统一响应格式供所有 API 使用

### References

- [Source: architecture.md#Cross-Cutting Concerns] - 多租户安全隔离策略
- [Source: architecture.md#Format Patterns] - API 响应格式
- [Source: architecture.md#Authentication & Security] - JWT 租户识别
- [Source: epics.md#Story 2.3] - 原始故事定义
- [Source: 2-1-user-data-model.md] - 用户数据模型
- [Source: 2-2-jwt-token-middleware.md] - JWT 中间件

## Dev Agent Record

### Agent Model Used

qianfan-code-latest

### Debug Log References

无需调试日志，所有测试一次通过。

### Completion Notes List

1. **用户同步服务** - 实现了 `UserSyncService`，支持创建新用户和更新现有用户，包含租户不匹配验证。
2. **租户上下文管理** - 创建了 `pkg/contextx` 包（避免与标准库冲突），实现了 `TenantContext` 结构体和相关方法。
3. **租户隔离中间件** - 实现了 `TenantMiddleware`，从认证中间件获取用户信息并设置租户上下文。
4. **统一响应工具** - 实现了 `APIResponse` 结构体和 `Success`、`Error`、`Paginated` 等响应方法。
5. **Repository 租户隔离** - 为 `TenantUserRepository` 添加了带租户隔离验证的方法。

### File List

**新增文件:**
- `backend/internal/service/user_sync.go` - 用户同步服务（含乐观锁重试逻辑）
- `backend/internal/service/user_sync_test.go` - 用户同步服务测试（含乐观锁测试）
- `backend/pkg/contextx/tenant.go` - 租户上下文工具
- `backend/pkg/contextx/tenant_test.go` - 租户上下文测试
- `backend/pkg/middleware/tenant.go` - 租户隔离中间件
- `backend/pkg/middleware/tenant_test.go` - 租户隔离中间件测试
- `backend/pkg/response/response.go` - 统一响应工具
- `backend/pkg/response/response_test.go` - 响应工具测试
- `backend/internal/repository/tenant_isolation.go` - Repository 租户隔离方法（含乐观锁）
- `backend/internal/repository/tenant_isolation_test.go` - 租户隔离测试

**修改文件:**
- `backend/internal/domain/user/tenant_user.go` - 添加 Version 字段
- `backend/internal/repository/errors.go` - 添加 ErrOptimisticLockConflict 错误
- `backend/internal/repository/tenant_user.go` - Update 方法支持乐观锁
- `backend/go.mod` - 添加 sqlmock 和 testify/mock 依赖
- `backend/go.sum` - 依赖更新

## Change Log

| 日期 | 变更内容 |
|------|---------|
| 2026-03-05 | 完成Story 2.3 实现：用户同步服务、租户上下文、租户隔离中间件、统一响应格式、Repository租户隔离 |
| 2026-03-05 | 代码审查修复：H1错误检查方式、H2导入格式、H3中间件响应一致性、M1测试验证、添加同步日志和审计日志 |
| 2026-03-05 | 实现乐观锁：TenantUser添加Version字段、Repository Update使用版本检查、UserSyncService实现重试逻辑 |

## Review Record

### Code Review - 2026-03-05

**审查结果:** 3 CRITICAL, 3 HIGH, 3 MEDIUM, 2 LOW

**已修复的问题:**
- ✅ [H1] 脆弱的错误检查方式 - 改用 `repository.ErrNotFound` 进行错误判断
- ✅ [H2] 导入格式错误 - 修正导入格式
- ✅ [H3] 中间件错误响应不一致 - 使用统一的 `response.Error` 函数
- ✅ [M1] 测试缺少用户参数验证 - 添加字段验证断言
- ✅ [C2] 添加同步日志记录 - 在 `user_sync.go` 中添加 zap 日志
- ✅ [C3] 添加查询审计日志 - 在 `tenant_isolation.go` 中添加审计日志
- ✅ [C1] 实现乐观锁 - TenantUser添加Version字段，Update使用版本检查和自动重试
