# Story 2.2: 业务平台 JWT Token 验证中间件

Status: ready-for-dev

## Story

As a 后端开发者,
I want 实现 JWT Token 验证中间件,
so that 可以验证来自业务平台的用户身份。

## Acceptance Criteria

1. **AC1: JWT Token 验证功能**
   - **Given** 业务平台使用 JWT Token 认证
   - **When** 前端请求携带 X-Platform-Token Header
   - **Then** 中间件验证 Token 有效性
   - **And** 解析 Token 获取用户信息

2. **AC2: 业务平台用户信息获取**
   - **Given** JWT Token 验证通过
   - **When** 调用业务平台用户信息接口
   - **Then** 获取用户详细信息（租户、权限等）
   - **And** 将用户信息存入请求上下文

3. **AC3: 错误处理**
   - **Given** Token 验证失败或用户信息获取失败
   - **When** 发生异常情况
   - **Then** 返回 401 错误
   - **And** 返回明确的错误信息

4. **AC4: 中间件集成**
   - **Given** 中间件已实现
   - **When** 注册到 HTTP 路由
   - **Then** 保护需要认证的路由
   - **And** 无需修改现有 Handler 代码

## Tasks / Subtasks

- [ ] Task 1: 创建 JWT 验证工具 (AC: 1)
  - [ ] 1.1 创建 `pkg/jwt/jwt.go` 定义 JWT 解析函数
  - [ ] 1.2 实现 Token 解析和验证逻辑
  - [ ] 1.3 定义 JWT Claims 结构体
  - [ ] 1.4 添加 Token 过期检查

- [ ] Task 2: 创建业务平台客户端 (AC: 2)
  - [ ] 2.1 创建 `internal/infrastructure/platform/client.go`
  - [ ] 2.2 实现用户信息获取接口调用
  - [ ] 2.3 定义用户信息响应结构体
  - [ ] 2.4 实现错误重试机制

- [ ] Task 3: 实现认证中间件 (AC: 1, 2, 3)
  - [ ] 3.1 创建 `pkg/middleware/auth.go`
  - [ ] 3.2 实现 Token 提取逻辑
  - [ ] 3.3 实现用户信息获取逻辑
  - [ ] 3.4 实现请求上下文注入
  - [ ] 3.5 实现错误响应处理

- [ ] Task 4: 集成中间件到路由 (AC: 4)
  - [ ] 4.1 修改 `cmd/server/main.go` 注册中间件
  - [ ] 4.2 创建需要认证的路由组
  - [ ] 4.3 创建公开路由组

- [ ] Task 5: 编写单元测试 (AC: 1-4)
  - [ ] 5.1 编写 `pkg/jwt/jwt_test.go`
  - [ ] 5.2 编写 `pkg/middleware/auth_test.go`
  - [ ] 5.3 编写 `internal/infrastructure/platform/client_test.go`
  - [ ] 5.4 使用 mock 测试业务平台接口调用

## Dev Notes

### 架构模式与约束

**必须遵循的 Clean Architecture 原则 [Source: architecture.md]:**

1. **中间件位置**:
   - `pkg/middleware/` - 放置可复用的中间件
   - `pkg/jwt/` - 放置 JWT 相关工具函数
   - `internal/infrastructure/platform/` - 放置业务平台客户端

2. **命名约定 [Source: architecture.md#Naming Patterns]:**
   - Header 命名: X-Platform-Token
   - 函数名: camelCase (例: `ValidateToken`, `GetUserInfo`)
   - 变量名: camelCase

3. **API 响应格式 [Source: architecture.md#Format Patterns]:**
   - 统一包装器: `{ data: {...}, error: null, meta: {...} }`
   - 错误响应: `{ data: null, error: { code: "UNAUTHORIZED", message: "..." }, meta: null }`

### 现有项目状态

**依赖 Story 2.1 用户数据模型已完成:**

```
backend/
├── internal/
│   ├── domain/
│   │   ├── user/
│   │   │   ├── tenant_user.go    # 租户用户实体
│   │   │   └── admin_user.go     # 管理员用户实体
│   ├── repository/
│   │   ├── tenant_user.go        # 租户用户仓库
│   │   └── admin_user.go         # 管理员用户仓库
│   └── infrastructure/
│       └── database/
│           └── database.go       # 数据库连接
├── pkg/
│   └── middleware/               # 待创建
```

### JWT Token 验证设计

**Token Header 格式:**
```
X-Platform-Token: <jwt_token>
```

**JWT Claims 结构:**

```go
// pkg/jwt/jwt.go
package jwt

import "github.com/golang-jwt/jwt/v5"

type PlatformClaims struct {
    UserID   string `json:"user_id"`
    Email    string `json:"email"`
    TenantID string `json:"tenant_id"`
    jwt.RegisteredClaims
}
```

**验证逻辑:**

```go
func ValidateToken(tokenString, secret string) (*PlatformClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &PlatformClaims{},
        func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*PlatformClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, errors.New("invalid token")
}
```

### 业务平台客户端设计

**用户信息接口调用:**

```go
// internal/infrastructure/platform/client.go
package platform

import (
    "context"
    "net/http"
)

type Client struct {
    baseURL    string
    httpClient *http.Client
}

type UserInfo struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    TenantID string `json:"tenantId"`
    Role     string `json:"role"`
}

func (c *Client) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
    // 调用业务平台接口获取用户信息
}
```

### 认证中间件设计

**中间件实现:**

```go
// pkg/middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"
)

type contextKey string

const (
    UserContextKey contextKey = "user"
)

type UserContext struct {
    ID       string
    Name     string
    Email    string
    TenantID string
    Role     string
}

func AuthMiddleware(jwtSecret string, platformClient *platform.Client) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 1. 提取 Token
            token := extractToken(r)
            if token == "" {
                respondWithError(w, http.StatusUnauthorized, "missing token")
                return
            }

            // 2. 验证 Token
            claims, err := jwt.ValidateToken(token, jwtSecret)
            if err != nil {
                respondWithError(w, http.StatusUnauthorized, "invalid token")
                return
            }

            // 3. 获取用户信息
            userInfo, err := platformClient.GetUserInfo(r.Context(), token)
            if err != nil {
                respondWithError(w, http.StatusUnauthorized, "failed to get user info")
                return
            }

            // 4. 注入上下文
            ctx := context.WithValue(r.Context(), UserContextKey, &UserContext{
                ID:       userInfo.ID,
                Name:     userInfo.Name,
                Email:    userInfo.Email,
                TenantID: userInfo.TenantID,
                Role:     userInfo.Role,
            })

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### 项目结构规范

**新增文件位置:**

```
backend/
├── pkg/
│   ├── jwt/
│   │   ├── jwt.go              # JWT 工具函数（新增）
│   │   └── jwt_test.go         # 测试文件（新增）
│   └── middleware/
│       ├── auth.go             # 认证中间件（新增）
│       └── auth_test.go        # 测试文件（新增）
├── internal/
│   └── infrastructure/
│       └── platform/
│           ├── client.go       # 业务平台客户端（新增）
│           └── client_test.go  # 测试文件（新增）
```

### 测试标准

**测试要求:**
- 测试框架: Go 原生 testing 包 + testify
- 测试覆盖率目标: >= 70%
- 使用 httptest 测试 HTTP 中间件

**测试用例设计:**

| 测试场景 | 测试文件 | 测试方法 |
|---------|---------|---------|
| Token 解析成功 | `jwt_test.go` | 单元测试 |
| Token 过期 | `jwt_test.go` | 单元测试 |
| Token 格式错误 | `jwt_test.go` | 单元测试 |
| 中间件正常流程 | `auth_test.go` | 单元测试 |
| 缺少 Token | `auth_test.go` | 单元测试 |
| Token 无效 | `auth_test.go` | 单元测试 |
| 业务平台调用失败 | `auth_test.go` | 单元测试 |

### 配置项

**需要的环境变量:**

```env
# .env.example
JWT_SECRET=your-jwt-secret-key
PLATFORM_BASE_URL=https://platform.example.com
PLATFORM_TIMEOUT=5s
```

### Project Structure Notes

**与 Story 2.1 的连续性:**

1. **复用用户领域模型**:
   - 使用 `TenantUser` 结构体定义
   - 使用 `TenantUserRepository` 查询用户

2. **依赖关系**:
   - Story 2.1 提供用户数据模型和仓库
   - 本 Story 提供认证中间件

3. **后续 Story**:
   - Story 2.3 将使用本中间件获取的用户信息

### References

- [Source: architecture.md#Authentication & Security] - 混合认证架构
- [Source: architecture.md#Naming Patterns] - Header 命名约定
- [Source: architecture.md#Format Patterns] - API 响应格式
- [Source: epics.md#Story 2.2] - 原始故事定义
- [Source: 2-1-user-data-model.md] - 用户数据模型上下文

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
