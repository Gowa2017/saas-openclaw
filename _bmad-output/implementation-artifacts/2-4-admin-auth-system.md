# Story 2.4: 平台管理员独立认证系统

Status: ready-for-dev

## Story

As a 平台管理员,
I want 使用独立的用户名密码登录管理后台,
so that 可以与普通用户认证体系分离。

## Acceptance Criteria

1. **AC1: 管理员登录验证**
   - **Given** 管理员账号已创建
   - **When** 管理员输入用户名密码登录
   - **Then** 验证用户名密码正确性
   - **And** 密码使用 bcrypt 加密存储

2. **AC2: 管理员 JWT Token 生成**
   - **Given** 登录验证通过
   - **When** 生成管理员 Token
   - **Then** 生成管理员 JWT Token
   - **And** Token 包含管理员角色标识
   - **And** Token 包含管理员 ID 和用户名

3. **AC3: 登录失败处理**
   - **Given** 管理员登录失败
   - **When** 验证不通过
   - **Then** 返回明确的错误信息
   - **And** 不暴露具体是用户名还是密码错误
   - **And** 记录登录失败日志

4. **AC4: 管理员账号管理**
   - **Given** 系统需要创建管理员账号
   - **When** 执行账号创建命令
   - **Then** 可以创建新的管理员账号
   - **And** 密码自动 bcrypt 加密

## Tasks / Subtasks

- [ ] Task 1: 创建管理员认证服务 (AC: 1, 2, 3)
  - [ ] 1.1 创建 `internal/service/admin_auth.go`
  - [ ] 1.2 实现 Login 方法（验证密码）
  - [ ] 1.3 实现 GenerateToken 方法
  - [ ] 1.4 实现登录失败日志记录

- [ ] Task 2: 创建管理员认证 API (AC: 1, 2, 3)
  - [ ] 2.1 创建 `internal/api/admin/auth_handler.go`
  - [ ] 2.2 实现 POST /v1/admin/auth/login 端点
  - [ ] 2.3 实现请求参数验证
  - [ ] 2.4 实现响应格式化

- [ ] Task 3: 创建管理员中间件 (AC: 2)
  - [ ] 3.1 创建 `pkg/middleware/admin_auth.go`
  - [ ] 3.2 实现管理员 Token 验证
  - [ ] 3.3 实现角色权限检查
  - [ ] 3.4 实现管理员上下文注入

- [ ] Task 4: 创建管理员 CLI 命令 (AC: 4)
  - [ ] 4.1 创建 `cmd/admin/main.go`
  - [ ] 4.2 实现 create-admin 命令
  - [ ] 4.3 实现密码输入（隐藏显示）
  - [ ] 4.4 实现密码加密存储

- [ ] Task 5: 更新 Repository (AC: 1)
  - [ ] 5.1 在 `AdminUserRepository` 添加 GetByUsername 方法
  - [ ] 5.2 实现密码验证方法

- [ ] Task 6: 编写单元测试 (AC: 1-4)
  - [ ] 6.1 编写 `admin_auth_test.go` 测试认证服务
  - [ ] 6.2 编写 `auth_handler_test.go` 测试 API
  - [ ] 6.3 编写 `admin_auth_middleware_test.go` 测试中间件
  - [ ] 6.4 编写 `admin_cli_test.go` 测试 CLI 命令

## Dev Notes

### 架构模式与约束

**必须遵循的 Clean Architecture 原则 [Source: architecture.md]:**

1. **分层设计**:
   - `internal/api/admin/` - 管理员 API 处理器
   - `internal/service/` - 认证业务逻辑
   - `internal/repository/` - 数据访问
   - `cmd/admin/` - CLI 命令

2. **命名约定 [Source: architecture.md#Naming Patterns]:**
   - API 端点: `/v1/admin/auth/login`
   - 密码字段: `PasswordHash`（bcrypt 加密）

3. **API 响应格式 [Source: architecture.md#Format Patterns]:**
   - 统一包装器: `{ data: {...}, error: null, meta: {...} }`

### 现有项目状态

**依赖 Story 2.1 用户数据模型已完成:**

```
backend/
├── internal/
│   ├── domain/
│   │   └── user/
│   │       └── admin_user.go     # 管理员用户实体
│   ├── repository/
│   │   └── admin_user.go         # 管理员用户仓库
│   └── infrastructure/
│       └── database/
│           └── database.go       # 数据库连接
├── pkg/
│   ├── jwt/
│   │   └── jwt.go                # JWT 工具
│   ├── middleware/
│   │   └── auth.go               # 租户认证中间件
│   └── response/
│       └── response.go           # 响应工具
```

### 管理员认证服务设计

**服务实现:**

```go
// internal/service/admin_auth.go
package service

import (
    "context"
    "errors"
    "time"
    "github.com/gowa/saas-openclaw/backend/internal/domain/user"
    "github.com/gowa/saas-openclaw/backend/internal/repository"
    "github.com/gowa/saas-openclaw/backend/pkg/jwt"
    "golang.org/x/crypto/bcrypt"
)

type AdminAuthService struct {
    adminRepo  *repository.AdminUserRepository
    jwtSecret  string
    tokenExp   time.Duration
}

type LoginRequest struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
    Token     string `json:"token"`
    ExpiresAt int64  `json:"expiresAt"`
    AdminID   string `json:"adminId"`
    Username  string `json:"username"`
    Role      string `json:"role"`
}

// Login 管理员登录
func (s *AdminAuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
    // 1. 查找用户
    admin, err := s.adminRepo.GetByUsername(ctx, req.Username)
    if err != nil {
        // 记录失败日志（不暴露是用户名还是密码错误）
        logLoginFailure(req.Username, "user not found")
        return nil, errors.New("invalid credentials")
    }

    // 2. 验证密码
    if err := bcrypt.CompareHashAndPassword(
        []byte(admin.PasswordHash),
        []byte(req.Password),
    ); err != nil {
        logLoginFailure(req.Username, "invalid password")
        return nil, errors.New("invalid credentials")
    }

    // 3. 生成 Token
    expiresAt := time.Now().Add(s.tokenExp)
    claims := &jwt.AdminClaims{
        AdminID:  admin.ID,
        Username: admin.Username,
        Role:     string(admin.Role),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token, err := jwt.GenerateToken(claims, s.jwtSecret)
    if err != nil {
        return nil, errors.New("failed to generate token")
    }

    // 4. 返回响应
    return &LoginResponse{
        Token:     token,
        ExpiresAt: expiresAt.Unix(),
        AdminID:   admin.ID,
        Username:  admin.Username,
        Role:      string(admin.Role),
    }, nil
}
```

### 管理员 JWT Claims 设计

**Token Claims 结构:**

```go
// pkg/jwt/admin_claims.go
package jwt

import "github.com/golang-jwt/jwt/v5"

type AdminClaims struct {
    AdminID  string `json:"admin_id"`
    Username string `json:"username"`
    Role     string `json:"role"` // "admin" 或 "super_admin"
    jwt.RegisteredClaims
}

// GenerateAdminToken 生成管理员 Token
func GenerateAdminToken(claims *AdminClaims, secret string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

// ValidateAdminToken 验证管理员 Token
func ValidateAdminToken(tokenString, secret string) (*AdminClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{},
        func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*AdminClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, errors.New("invalid admin token")
}
```

### 管理员认证 API 设计

**登录 API:**

```go
// internal/api/admin/auth_handler.go
package admin

import (
    "net/http"
    "github.com/gowa/saas-openclaw/backend/internal/service"
    "github.com/gowa/saas-openclaw/backend/pkg/response"
)

type AuthHandler struct {
    authService *service.AdminAuthService
}

// Login 处理管理员登录
// POST /v1/admin/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req service.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
        return
    }

    // 验证请求参数
    if req.Username == "" || req.Password == "" {
        response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Username and password are required")
        return
    }

    // 执行登录
    resp, err := h.authService.Login(r.Context(), &req)
    if err != nil {
        response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid credentials")
        return
    }

    response.Success(w, resp)
}
```

### 管理员中间件设计

**中间件实现:**

```go
// pkg/middleware/admin_auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"
    "github.com/gowa/saas-openclaw/backend/pkg/jwt"
    "github.com/gowa/saas-openclaw/backend/pkg/response"
)

type AdminContextKey struct{}

type AdminContext struct {
    AdminID  string
    Username string
    Role     string
}

// AdminAuthMiddleware 管理员认证中间件
func AdminAuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 1. 提取 Token
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization header")
                return
            }

            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            if tokenString == authHeader {
                response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization format")
                return
            }

            // 2. 验证 Token
            claims, err := jwt.ValidateAdminToken(tokenString, jwtSecret)
            if err != nil {
                response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
                return
            }

            // 3. 注入上下文
            ctx := context.WithValue(r.Context(), AdminContextKey{}, &AdminContext{
                AdminID:  claims.AdminID,
                Username: claims.Username,
                Role:     claims.Role,
            })

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// RequireRole 角色检查中间件
func RequireRole(roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            adminCtx, ok := r.Context().Value(AdminContextKey{}).(*AdminContext)
            if !ok {
                response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Admin context not found")
                return
            }

            for _, role := range roles {
                if adminCtx.Role == role {
                    next.ServeHTTP(w, r)
                    return
                }
            }

            response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
        })
    }
}
```

### 管理员 CLI 命令设计

**创建管理员命令:**

```go
// cmd/admin/main.go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "github.com/gowa/saas-openclaw/backend/internal/repository"
    "github.com/gowa/saas-openclaw/backend/internal/domain/user"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "golang.org/x/term"
)

func main() {
    if len(os.Args) < 2 {
        printUsage()
        os.Exit(1)
    }

    switch os.Args[1] {
    case "create-admin":
        createAdmin()
    default:
        printUsage()
        os.Exit(1)
    }
}

func createAdmin() {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Enter username: ")
    username, _ := reader.ReadString('\n')
    username = strings.TrimSpace(username)

    fmt.Print("Enter password: ")
    passwordBytes, _ := term.ReadPassword(int(os.Stdin.Fd()))
    fmt.Println()

    fmt.Print("Enter name: ")
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)

    fmt.Print("Enter email: ")
    email, _ := reader.ReadString('\n')
    email = strings.TrimSpace(email)

    // 加密密码
    hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
    if err != nil {
        fmt.Printf("Error hashing password: %v\n", err)
        os.Exit(1)
    }

    // 创建管理员
    admin := &user.AdminUser{
        ID:           uuid.New().String(),
        Username:     username,
        PasswordHash: string(hashedPassword),
        Name:         name,
        Email:        email,
        Role:         user.AdminRoleAdmin,
    }

    // 保存到数据库
    // ...

    fmt.Printf("Admin user '%s' created successfully!\n", username)
}

func printUsage() {
    fmt.Println("Usage: admin <command>")
    fmt.Println("Commands:")
    fmt.Println("  create-admin  Create a new admin user")
}
```

### 项目结构规范

**新增文件位置:**

```
backend/
├── cmd/
│   └── admin/
│       └── main.go              # 管理员 CLI（新增）
├── internal/
│   ├── api/
│   │   └── admin/
│   │       ├── auth_handler.go  # 认证 API（新增）
│   │       └── auth_handler_test.go
│   └── service/
│       ├── admin_auth.go        # 管理员认证服务（新增）
│       └── admin_auth_test.go
├── pkg/
│   ├── jwt/
│   │   ├── admin_claims.go      # 管理员 Token（新增）
│   │   └── admin_claims_test.go
│   └── middleware/
│       ├── admin_auth.go        # 管理员中间件（新增）
│       └── admin_auth_test.go
```

### 测试标准

**测试要求:**
- 测试框架: Go 原生 testing 包 + testify
- 测试覆盖率目标: >= 70%
- 使用 mock 测试数据库操作

**测试用例设计:**

| 测试场景 | 测试文件 | 测试方法 |
|---------|---------|---------|
| 登录成功 | `admin_auth_test.go` | 单元测试 |
| 用户名不存在 | `admin_auth_test.go` | 单元测试 |
| 密码错误 | `admin_auth_test.go` | 单元测试 |
| Token 生成 | `admin_claims_test.go` | 单元测试 |
| Token 验证 | `admin_claims_test.go` | 单元测试 |
| 中间件认证 | `admin_auth_middleware_test.go` | 单元测试 |
| 角色检查 | `admin_auth_middleware_test.go` | 单元测试 |
| API 登录 | `auth_handler_test.go` | 集成测试 |

### 配置项

**需要的环境变量:**

```env
# .env.example
ADMIN_JWT_SECRET=your-admin-jwt-secret-key
ADMIN_TOKEN_EXPIRATION=24h
```

### 安全注意事项

1. **密码存储**: 必须使用 bcrypt 加密，成本因子 >= 10
2. **错误信息**: 登录失败不暴露是用户名还是密码错误
3. **登录日志**: 记录登录失败尝试，可用于检测暴力破解
4. **Token 有效期**: 建议 24 小时，支持刷新机制
5. **角色权限**: 区分 admin 和 super_admin 角色

### Project Structure Notes

**与前置 Story 的连续性:**

1. **依赖 Story 2.1**:
   - 使用 AdminUser 结构体
   - 使用 AdminUserRepository

2. **复用 Story 2.2 的基础设施**:
   - JWT 验证工具
   - 统一响应格式

3. **与租户认证分离**:
   - 独立的管理员认证流程
   - 独立的 Token Claims 结构
   - 独立的中间件

### References

- [Source: architecture.md#Authentication & Security] - 混合认证架构
- [Source: architecture.md#Format Patterns] - API 响应格式
- [Source: prd.md#FR4] - 平台管理员登录管理后台
- [Source: prd.md#NFR-S3] - 平台管理员操作需要身份验证
- [Source: epics.md#Story 2.4] - 原始故事定义
- [Source: 2-1-user-data-model.md] - 用户数据模型

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
