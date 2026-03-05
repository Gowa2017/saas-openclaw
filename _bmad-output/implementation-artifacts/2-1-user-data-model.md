# Story 2.1: 用户数据模型与数据库表

Status: done

## Story

As a 开发者,
I want 创建用户数据模型和数据库表,
so that 系统可以存储和管理用户信息。

## Acceptance Criteria

1. **AC1: tenant_users 表创建成功**
   - **Given** PostgreSQL 数据库已配置
   - **When** 执行数据库迁移
   - **Then** 创建 `tenant_users` 表（表名 snake_case）
   - **And** 表包含 ID、TenantID、Name、Email、Role、CreatedAt、UpdatedAt 列（列名 PascalCase）
   - **And** TenantID 列有外键约束关联到 tenants 表

2. **AC2: admin_users 表创建成功**
   - **Given** PostgreSQL 数据库已配置
   - **When** 执行数据库迁移
   - **Then** 创建 `admin_users` 表（表名 snake_case）
   - **And** 表包含 ID、Username、PasswordHash、Name、Email、Role、CreatedAt、UpdatedAt 列
   - **And** PasswordHash 列存储 bcrypt 加密的密码

3. **AC3: 数据库索引优化**
   - **Given** 表已创建
   - **When** 检查索引配置
   - **Then** tenant_users 表有 TenantID 列索引
   - **And** tenant_users 表有 Email 列唯一索引
   - **And** admin_users 表有 Username 列唯一索引
   - **And** admin_users 表有 Email 列唯一索引

4. **AC4: Go 领域模型定义**
   - **Given** 项目已初始化
   - **When** 创建领域模型文件
   - **Then** 在 `internal/domain/user/` 定义 TenantUser 结构体
   - **And** 在 `internal/domain/user/` 定义 AdminUser 结构体
   - **And** 结构体字段使用 `db` tag 映射数据库列
   - **And** 结构体字段使用 `json` tag 定义 JSON 序列化

5. **AC5: Repository 层实现**
   - **Given** 领域模型已定义
   - **When** 创建 Repository 文件
   - **Then** 在 `internal/repository/user.go` 实现 UserRepository 接口
   - **And** 支持 Create、GetByID、GetByEmail、Update、Delete 方法
   - **And** 使用 sqlx 进行数据库操作

## Tasks / Subtasks

- [x] Task 1: 创建租户表 (AC: 1)
  - [x] 1.1 创建 `internal/domain/tenant/tenant.go` 定义 Tenant 结构体
  - [x] 1.2 创建数据库迁移脚本 `migrations/000001_create_tenants.up.sql`
  - [x] 1.3 创建回滚脚本 `migrations/000001_create_tenants.down.sql`

- [x] Task 2: 创建用户领域模型 (AC: 4)
  - [x] 2.1 创建 `internal/domain/user/tenant_user.go` 定义 TenantUser 结构体
  - [x] 2.2 创建 `internal/domain/user/admin_user.go` 定义 AdminUser 结构体
  - [x] 2.3 定义用户角色枚举 (RoleTenantUser, RoleTenantAdmin)
  - [x] 2.4 添加 `json` 和 `db` 结构体标签

- [x] Task 3: 创建数据库迁移脚本 (AC: 1, 2, 3)
  - [x] 3.1 创建 `migrations/000002_create_tenant_users.up.sql`
  - [x] 3.2 创建 `migrations/000002_create_tenant_users.down.sql`
  - [x] 3.3 创建 `migrations/000003_create_admin_users.up.sql`
  - [x] 3.4 创建 `migrations/000003_create_admin_users.down.sql`
  - [x] 3.5 添加索引定义

- [x] Task 4: 实现 Repository 层 (AC: 5)
  - [x] 4.1 创建 `internal/repository/tenant.go`
  - [x] 4.2 创建 `internal/repository/tenant_user.go`
  - [x] 4.3 创建 `internal/repository/admin_user.go`
  - [x] 4.4 实现创建用户方法
  - [x] 4.5 实现查询用户方法 (ByID, ByEmail, ByTenantID)
  - [x] 4.6 实现更新用户方法
  - [x] 4.7 实现删除用户方法

- [x] Task 5: 编写单元测试 (AC: 1-5)
  - [x] 5.1 编写 `tenant_test.go` 测试租户领域模型
  - [x] 5.2 编写 `tenant_user_test.go` 测试领域模型
  - [x] 5.3 编写 `admin_user_test.go` 测试领域模型
  - [x] 5.4 编写 `tenant_user_repository_test.go` 测试 Repository
  - [x] 5.5 编写 `admin_user_repository_test.go` 测试 Repository

## Dev Notes

### 架构模式与约束

**必须遵循的 Clean Architecture 原则 [Source: architecture.md]:**

1. **依赖方向**:
   - `domain` 层定义实体和接口（无外部依赖）
   - `repository` 层实现数据访问（依赖 domain 和 infrastructure/database）

2. **命名约定 [Source: architecture.md#Naming Patterns]:**
   - 表名: `snake_case` (例: `tenant_users`, `admin_users`)
   - 列名: `PascalCase` (例: `TenantID`, `CreatedAt`, `UpdatedAt`)
   - 外键: `table_id` (例: `tenant_id`)
   - API 端点: 复数资源名 (例: `/users`)

3. **API 响应格式 [Source: architecture.md#Format Patterns]:**
   - 统一包装器: `{ data: {...}, error: null, meta: {...} }`

### 现有项目状态

**数据库连接已完成 [Source: 1-3-postgresql-database-config.md]:**

```
backend/
├── internal/
│   ├── infrastructure/
│   │   ├── config/
│   │   │   └── config.go      # 数据库配置结构
│   │   └── database/
│   │       └── database.go    # 连接池已配置
│   ├── domain/
│   │   ├── user/              # ✅ 目录存在（空）
│   │   ├── tenant/            # ✅ 目录存在（空）
│   │   ├── instance/          # ✅ 目录存在（空）
│   │   └── config/            # ✅ 目录存在（空）
│   └── repository/            # ✅ 目录存在（空）
```

**数据库连接功能:**
- `Connect(cfg)` - 建立数据库连接
- `Ping(db)` - 健康检查
- `Stats(db)` - 连接池统计
- `Close(db)` - 优雅关闭
- 连接池: MaxOpenConns=100, MaxIdleConns=10

**依赖版本 [Source: backend/go.mod]:**
- Go: 1.25.0
- sqlx: v1.3.5
- lib/pq: v1.10.9
- Viper: v1.18.2

### 技术栈要求

**数据库迁移工具选择:**

推荐使用 **golang-migrate** 工具进行数据库迁移管理。

```bash
# 安装 migrate 工具
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 创建迁移文件
migrate create -ext sql -dir migrations -seq create_tenants
migrate create -ext sql -dir migrations -seq create_tenant_users
migrate create -ext sql -dir migrations -seq create_admin_users
```

**密码加密:**
- 使用 `golang.org/x/crypto/bcrypt` 进行密码哈希
- 成本因子: 默认 10

```go
import "golang.org/x/crypto/bcrypt"

// HashPassword 生成密码哈希
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### 数据库表设计

**tenants 表（先创建）:**

```sql
CREATE TABLE tenants (
    "ID" VARCHAR(36) PRIMARY KEY,  -- UUID
    "Name" VARCHAR(255) NOT NULL,
    "Status" VARCHAR(50) NOT NULL DEFAULT 'active',
    "CreatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "UpdatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tenants_status ON tenants ("Status");
```

**tenant_users 表:**

```sql
CREATE TABLE tenant_users (
    "ID" VARCHAR(36) PRIMARY KEY,  -- UUID
    "TenantID" VARCHAR(36) NOT NULL REFERENCES tenants("ID") ON DELETE CASCADE,
    "Name" VARCHAR(255) NOT NULL,
    "Email" VARCHAR(255) NOT NULL,
    "Role" VARCHAR(50) NOT NULL DEFAULT 'user',
    "CreatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "UpdatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tenant_users_tenant_id ON tenant_users ("TenantID");
CREATE UNIQUE INDEX idx_tenant_users_email ON tenant_users ("Email");
```

**admin_users 表:**

```sql
CREATE TABLE admin_users (
    "ID" VARCHAR(36) PRIMARY KEY,  -- UUID
    "Username" VARCHAR(100) NOT NULL UNIQUE,
    "PasswordHash" VARCHAR(255) NOT NULL,
    "Name" VARCHAR(255) NOT NULL,
    "Email" VARCHAR(255) NOT NULL,
    "Role" VARCHAR(50) NOT NULL DEFAULT 'admin',
    "CreatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "UpdatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_admin_users_username ON admin_users ("Username");
CREATE UNIQUE INDEX idx_admin_users_email ON admin_users ("Email");
```

### 领域模型设计

**TenantUser 结构体:**

```go
// internal/domain/user/tenant_user.go
package user

import "time"

type Role string

const (
    RoleTenantUser  Role = "user"
    RoleTenantAdmin Role = "admin"
)

type TenantUser struct {
    ID        string    `json:"id" db:"ID"`
    TenantID  string    `json:"tenantId" db:"TenantID"`
    Name      string    `json:"name" db:"Name"`
    Email     string    `json:"email" db:"Email"`
    Role      Role      `json:"role" db:"Role"`
    CreatedAt time.Time `json:"createdAt" db:"CreatedAt"`
    UpdatedAt time.Time `json:"updatedAt" db:"UpdatedAt"`
}
```

**AdminUser 结构体:**

```go
// internal/domain/user/admin_user.go
package user

import "time"

type AdminRole string

const (
    AdminRoleSuperAdmin AdminRole = "super_admin"
    AdminRoleAdmin      AdminRole = "admin"
)

type AdminUser struct {
    ID           string     `json:"id" db:"ID"`
    Username     string     `json:"username" db:"Username"`
    PasswordHash string     `json:"-" db:"PasswordHash"` // 不暴露密码哈希
    Name         string     `json:"name" db:"Name"`
    Email        string     `json:"email" db:"Email"`
    Role         AdminRole  `json:"role" db:"Role"`
    CreatedAt    time.Time  `json:"createdAt" db:"CreatedAt"`
    UpdatedAt    time.Time  `json:"updatedAt" db:"UpdatedAt"`
}
```

### Repository 接口设计

**TenantUserRepository:**

```go
// internal/repository/tenant_user.go
package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/gowa/saas-openclaw/backend/internal/domain/user"
)

type TenantUserRepository struct {
    db *sqlx.DB
}

func NewTenantUserRepository(db *sqlx.DB) *TenantUserRepository {
    return &TenantUserRepository{db: db}
}

func (r *TenantUserRepository) Create(u *user.TenantUser) error
func (r *TenantUserRepository) GetByID(id string) (*user.TenantUser, error)
func (r *TenantUserRepository) GetByEmail(email string) (*user.TenantUser, error)
func (r *TenantUserRepository) GetByTenantID(tenantID string) ([]*user.TenantUser, error)
func (r *TenantUserRepository) Update(u *user.TenantUser) error
func (r *TenantUserRepository) Delete(id string) error
```

### 项目结构规范

**新增文件位置:**

```
backend/
├── migrations/                    # 数据库迁移（新增目录）
│   ├── 000001_create_tenants.up.sql
│   ├── 000001_create_tenants.down.sql
│   ├── 000002_create_tenant_users.up.sql
│   ├── 000002_create_tenant_users.down.sql
│   ├── 000003_create_admin_users.up.sql
│   └── 000003_create_admin_users.down.sql
├── internal/
│   ├── domain/
│   │   ├── tenant/
│   │   │   └── tenant.go         # 租户实体（新增）
│   │   └── user/
│   │       ├── tenant_user.go    # 租户用户实体（新增）
│   │       ├── admin_user.go     # 管理员用户实体（新增）
│   │       └── role.go           # 角色定义（新增）
│   └── repository/
│       ├── tenant.go             # 租户仓库（新增）
│       ├── tenant_user.go        # 租户用户仓库（新增）
│       └── admin_user.go         # 管理员用户仓库（新增）
```

### 测试标准

**测试要求:**
- 测试框架: Go 原生 testing 包 + testify
- 测试覆盖率目标: ≥ 70%
- 使用 testcontainers 进行集成测试

**测试用例设计:**

| 测试场景 | 测试文件 | 测试方法 |
|---------|---------|---------|
| 领域模型创建 | `tenant_user_test.go` | 单元测试 |
| JSON 序列化 | `tenant_user_test.go` | 单元测试 |
| Repository CRUD | `tenant_user_repository_test.go` | 集成测试 |
| 唯一约束 | `tenant_user_repository_test.go` | 集成测试 |

### Project Structure Notes

**与 Epic 1 的连续性:**

1. **复用现有基础设施**:
   - 使用 Story 1.3 配置的数据库连接
   - 使用现有的 `config.DatabaseConfig` 结构

2. **遵循现有模式**:
   - Clean Architecture 分层
   - sqlx 数据库操作
   - testify 断言库

3. **新增能力**:
   - 数据库迁移管理
   - 领域模型定义
   - Repository 层实现

### 前序 Story 的学习经验

**从 Story 1.1 (后端项目初始化) 获得的经验:**

1. **配置管理**: 使用 Viper 管理配置，支持环境变量
2. **测试策略**: 使用 testify 断言库，目标覆盖率 70%
3. **代码质量**: Code Review 发现问题及时修复

**从 Story 1.3 (PostgreSQL 数据库配置) 获得的经验:**

1. **数据库连接**: 连接池已配置完成
2. **健康检查**: `/health/database` 端点已实现
3. **SSL 支持**: 支持 SSL 连接配置

### 常见问题与解决方案

**问题 1: 迁移文件命名**
- **原因**: migrate 工具要求特定命名格式
- **解决**: 使用 `{version}_{name}.up.sql` 和 `{version}_{name}.down.sql`

**问题 2: UUID 生成**
- **原因**: Go 没有内置 UUID 库
- **解决**: 使用 `github.com/google/uuid` 库

**问题 3: 外键约束失败**
- **原因**: 创建 tenant_users 前需要先创建 tenants
- **解决**: 确保迁移顺序正确，先创建 tenants 表

### 安全注意事项

1. **密码存储**: 使用 bcrypt 哈希，不要存储明文密码
2. **敏感字段**: `PasswordHash` 字段使用 `json:"-"` 不暴露给 API
3. **SQL 注入**: 使用参数化查询，不要拼接 SQL
4. **UUID**: 使用 UUID 而非自增 ID，防止 ID 枚举攻击

### References

- [Source: architecture.md#Naming Patterns] - 数据库命名约定
- [Source: architecture.md#Project Structure] - 项目目录结构
- [Source: architecture.md#API Response Formats] - 统一响应格式
- [Source: prd.md#FR1-FR4] - 用户认证与授权需求
- [Source: prd.md#NFR-S2] - 用户配置数据加密存储
- [Source: epics.md#Story 2.1] - 原始故事定义
- [Source: 1-3-postgresql-database-config.md] - 数据库配置上下文

## Dev Agent Record

### Agent Model Used

qianfan-code-latest

### Debug Log References

无错误日志

### Completion Notes List

**2026-03-05 - 实现完成**

1. **领域模型层**:
   - 创建 Tenant 结构体，包含 ID、Name、Status、CreatedAt、UpdatedAt 字段
   - 创建 TenantUser 结构体，支持租户用户管理
   - 创建 AdminUser 结构体，支持平台管理员管理
   - 定义角色枚举 Role (user, admin) 和 AdminRole (admin, super_admin)
   - 所有结构体使用 `json` 和 `db` 标签

2. **数据库迁移**:
   - 创建 000001_create_tenants 迁移（包含 status 和 name 索引）
   - 创建 000002_create_tenant_users 迁移（包含 TenantID 和 Email 索引）
   - 创建 000003_create_admin_users 迁移（包含 Username 和 Email 唯一索引）
   - 所有迁移使用 PascalCase 列名和 snake_case 表名

3. **Repository 层**:
   - 实现 TenantRepository（Create, GetByID, GetByName, List, Update, Delete）
   - 实现 TenantUserRepository（Create, GetByID, GetByEmail, GetByTenantID, Update, Delete）
   - 实现 AdminUserRepository（Create, GetByID, GetByUsername, GetByEmail, List, Update, UpdatePassword, VerifyPassword, Delete）
   - AdminUserRepository 内置 bcrypt 密码哈希支持

4. **测试覆盖**:
   - 领域模型单元测试全部通过
   - Repository 单元测试通过（集成测试使用 -short 跳过）
   - 所有代码通过 go vet 检查

5. **依赖添加**:
   - github.com/google/uuid v1.6.0

**2026-03-05 - 代码审查修复**

修复了 6 个问题（2 HIGH, 4 MEDIUM）:

1. **[HIGH] 输入验证** - 为所有领域模型添加 Validate() 方法
   - Tenant: 验证 Name 必填，Status 有效
   - TenantUser: 验证 TenantID、Name、Email 必填，Role 有效
   - AdminUser: 验证 Username、Name、Email 必填，Role 有效
   - Repository 层在 Create/Update 前调用 Validate()

2. **[HIGH] 密码验证** - AdminUserRepository.Create 现在验证 password 非空

3. **[MEDIUM] 错误处理改进**:
   - 创建 `internal/repository/errors.go` 定义标准错误类型
   - ErrNotFound, ErrValidation, ErrDuplicateKey
   - IsNotFoundError() 辅助函数

4. **[MEDIUM] Update 错误区分** - 现在正确返回 ErrNotFound 当 ID 不存在

5. **[MEDIUM] Status 验证** - 添加 Status.IsValid() 方法

6. **[MEDIUM] 分页支持** - List 方法添加 limit/offset 参数

### Senior Developer Review (AI)

**Review Date:** 2026-03-05
**Reviewer Model:** qianfan-code-latest
**Outcome:** ✅ Approved (after fixes)

#### Action Items

- [x] [HIGH] 添加输入验证 - 所有领域模型 Validate() 方法
- [x] [HIGH] 添加密码验证 - AdminUserRepository.Create password 非空检查
- [x] [MEDIUM] 创建自定义错误类型 - errors.go
- [x] [MEDIUM] 改进 Update 错误处理 - 返回 ErrNotFound
- [x] [MEDIUM] 添加 Status.IsValid() 方法
- [x] [MEDIUM] 添加分页支持 - List 方法 limit/offset

#### Review Summary

代码质量良好，架构遵循 Clean Architecture 原则。发现的问题已全部修复。

### File List

**新增文件:**
- backend/internal/domain/tenant/tenant.go
- backend/internal/domain/tenant/tenant_test.go
- backend/internal/domain/user/role.go
- backend/internal/domain/user/tenant_user.go
- backend/internal/domain/user/admin_user.go
- backend/internal/domain/user/tenant_user_test.go
- backend/internal/domain/user/admin_user_test.go
- backend/internal/repository/errors.go
- backend/internal/repository/tenant.go
- backend/internal/repository/tenant_user.go
- backend/internal/repository/admin_user.go
- backend/internal/repository/tenant_user_repository_test.go
- backend/internal/repository/admin_user_repository_test.go
- backend/migrations/000001_create_tenants.up.sql
- backend/migrations/000001_create_tenants.down.sql
- backend/migrations/000002_create_tenant_users.up.sql
- backend/migrations/000002_create_tenant_users.down.sql
- backend/migrations/000003_create_admin_users.up.sql
- backend/migrations/000003_create_admin_users.down.sql

**修改文件:**
- backend/go.mod (添加 github.com/google/uuid 依赖)
- backend/go.sum (依赖更新)
