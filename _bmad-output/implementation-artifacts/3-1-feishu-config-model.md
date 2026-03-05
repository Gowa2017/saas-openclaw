# Story 3.1: 飞书配置数据模型

Status: ready-for-dev

## Story

As a 开发者,
I want 创建飞书配置数据模型,
so that 可以存储用户的飞书应用配置。

## Acceptance Criteria

1. **AC1: feishu_configs 表创建成功**
   - **Given** PostgreSQL 数据库已配置
   - **When** 执行数据库迁移
   - **Then** 创建 `feishu_configs` 表（表名 snake_case）
   - **And** 表包含 ID、TenantID、AppID、AppSecret、Status、CreatedAt、UpdatedAt 列（列名 PascalCase）
   - **And** TenantID 列有外键约束关联到 tenants 表

2. **AC2: AppSecret 加密存储**
   - **Given** 用户输入飞书 App Secret
   - **When** 保存配置到数据库
   - **Then** AppSecret 使用 AES-256 加密存储
   - **And** 加密密钥从环境变量读取
   - **And** 支持加密和解密函数

3. **AC3: 一租户一配置约束**
   - **Given** 租户已有飞书配置
   - **When** 尝试创建新配置
   - **Then** 数据库唯一约束阻止重复创建
   - **And** TenantID 列有唯一索引

4. **AC4: 数据库索引优化**
   - **Given** 表已创建
   - **When** 检查索引配置
   - **Then** TenantID 列有唯一索引
   - **And** Status 列有普通索引
   - **And** CreatedAt 列有索引支持排序查询

5. **AC5: Go 领域模型定义**
   - **Given** 项目已初始化
   - **When** 创建领域模型文件
   - **Then** 在 `internal/domain/config/` 定义 FeishuConfig 结构体
   - **And** 结构体字段使用 `db` tag 映射数据库列
   - **And** 结构体字段使用 `json` tag 定义 JSON 序列化
   - **And** AppSecret 字段不暴露给 JSON (`json:"-"`)

6. **AC6: Repository 层实现**
   - **Given** 领域模型已定义
   - **When** 创建 Repository 文件
   - **Then** 在 `internal/repository/feishu_config.go` 实现 FeishuConfigRepository 接口
   - **And** 支持 Create、GetByTenantID、Update、Delete 方法
   - **And** 保存时自动加密 AppSecret
   - **And** 读取时自动解密 AppSecret

## Tasks / Subtasks

- [ ] Task 1: 创建飞书配置领域模型 (AC: 5)
  - [ ] 1.1 创建 `internal/domain/config/feishu_config.go` 定义 FeishuConfig 结构体
  - [ ] 1.2 定义配置状态枚举 (StatusActive, StatusInactive, StatusPending)
  - [ ] 1.3 添加 `json` 和 `db` 结构体标签
  - [ ] 1.4 AppSecret 字段使用 `json:"-"` 标签

- [ ] Task 2: 实现 AES-256 加密工具 (AC: 2)
  - [ ] 2.1 创建 `pkg/encryption/aes.go` 加密工具包
  - [ ] 2.2 实现 Encrypt(plaintext string) (ciphertext string, error) 函数
  - [ ] 2.3 实现 Decrypt(ciphertext string) (plaintext string, error) 函数
  - [ ] 2.4 从环境变量读取加密密钥 (ENCRYPTION_KEY)
  - [ ] 2.5 密钥长度验证（必须为 32 字节）

- [ ] Task 3: 创建数据库迁移脚本 (AC: 1, 3, 4)
  - [ ] 3.1 创建 `migrations/004_create_feishu_configs.up.sql`
  - [ ] 3.2 创建 `migrations/004_create_feishu_configs.down.sql`
  - [ ] 3.3 添加 TenantID 唯一索引
  - [ ] 3.4 添加 Status 和 CreatedAt 索引
  - [ ] 3.5 添加外键约束关联 tenants 表

- [ ] Task 4: 实现 Repository 层 (AC: 6)
  - [ ] 4.1 创建 `internal/repository/feishu_config.go`
  - [ ] 4.2 实现 NewFeishuConfigRepository 构造函数
  - [ ] 4.3 实现 Create 方法（加密 AppSecret 后保存）
  - [ ] 4.4 实现 GetByTenantID 方法（解密 AppSecret 后返回）
  - [ ] 4.5 实现 Update 方法（加密 AppSecret 后更新）
  - [ ] 4.6 实现 Delete 方法

- [ ] Task 5: 编写单元测试 (AC: 1-6)
  - [ ] 5.1 编写 `feishu_config_test.go` 测试领域模型
  - [ ] 5.2 编写 `aes_test.go` 测试加密解密功能
  - [ ] 5.3 编写 `feishu_config_repository_test.go` 测试 Repository
  - [ ] 5.4 测试加密密钥缺失时的错误处理

## Dev Notes

### 架构模式与约束

**必须遵循的 Clean Architecture 原则 [Source: architecture.md]:**

1. **依赖方向**:
   - `domain` 层定义实体和接口（无外部依赖）
   - `repository` 层实现数据访问（依赖 domain 和 infrastructure/database）
   - `pkg/encryption` 是独立的工具包，可被 repository 层调用

2. **命名约定 [Source: architecture.md#Naming Patterns]:**
   - 表名: `snake_case` (例: `feishu_configs`)
   - 列名: `PascalCase` (例: `TenantID`, `AppID`, `AppSecret`)
   - 外键: `table_id` (例: `tenant_id`)

3. **安全要求 [Source: prd.md#NFR-S2]:**
   - 用户配置数据（App ID/Secret）加密存储
   - 使用 AES-256 加密算法

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
│   │   ├── user/              # ✅ 目录存在
│   │   ├── tenant/            # ✅ 目录存在
│   │   ├── instance/          # ✅ 目录存在
│   │   └── config/            # ✅ 目录存在（空）
│   └── repository/            # ✅ 目录存在
├── pkg/                       # 需要创建
│   └── encryption/            # 新建目录
```

**依赖版本 [Source: backend/go.mod]:**
- Go: 1.25.0
- sqlx: v1.3.5
- lib/pq: v1.10.9

### 技术栈要求

**AES-256 加密实现:**

```go
// pkg/encryption/aes.go
package encryption

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "io"
    "os"
)

var ErrInvalidKey = errors.New("encryption key must be 32 bytes")

// GetKey 从环境变量获取加密密钥
func GetKey() ([]byte, error) {
    key := os.Getenv("ENCRYPTION_KEY")
    if len(key) != 32 {
        return nil, ErrInvalidKey
    }
    return []byte(key), nil
}

// Encrypt 使用 AES-256-GCM 加密
func Encrypt(plaintext string) (string, error) {
    key, err := GetKey()
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 使用 AES-256-GCM 解密
func Decrypt(ciphertext string) (string, error) {
    key, err := GetKey()
    if err != nil {
        return "", err
    }

    data, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return "", errors.New("ciphertext too short")
    }

    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}
```

### 数据库表设计

**feishu_configs 表:**

```sql
CREATE TABLE feishu_configs (
    "ID" VARCHAR(36) PRIMARY KEY,  -- UUID
    "TenantID" VARCHAR(36) NOT NULL REFERENCES tenants("ID") ON DELETE CASCADE,
    "AppID" VARCHAR(255) NOT NULL,
    "AppSecret" VARCHAR(512) NOT NULL,  -- 加密后可能更长
    "Status" VARCHAR(50) NOT NULL DEFAULT 'pending',
    "CreatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "UpdatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 一租户一配置约束
CREATE UNIQUE INDEX idx_feishu_configs_tenant_id ON feishu_configs ("TenantID");

-- 状态查询优化
CREATE INDEX idx_feishu_configs_status ON feishu_configs ("Status");

-- 时间排序优化
CREATE INDEX idx_feishu_configs_created_at ON feishu_configs ("CreatedAt");
```

### 领域模型设计

**FeishuConfig 结构体:**

```go
// internal/domain/config/feishu_config.go
package config

import "time"

// ConfigStatus 飞书配置状态
type ConfigStatus string

const (
    StatusPending  ConfigStatus = "pending"   // 待验证
    StatusActive   ConfigStatus = "active"    // 已验证有效
    StatusInactive ConfigStatus = "inactive"  // 已失效
)

// FeishuConfig 飞书应用配置
type FeishuConfig struct {
    ID        string       `json:"id" db:"ID"`
    TenantID  string       `json:"tenantId" db:"TenantID"`
    AppID     string       `json:"appId" db:"AppID"`
    AppSecret string       `json:"-" db:"AppSecret"` // 不暴露给 API
    Status    ConfigStatus `json:"status" db:"Status"`
    CreatedAt time.Time    `json:"createdAt" db:"CreatedAt"`
    UpdatedAt time.Time    `json:"updatedAt" db:"UpdatedAt"`
}

// IsVerified 检查配置是否已验证
func (c *FeishuConfig) IsVerified() bool {
    return c.Status == StatusActive
}
```

### Repository 接口设计

**FeishuConfigRepository:**

```go
// internal/repository/feishu_config.go
package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/gowa/saas-openclaw/backend/internal/domain/config"
    "github.com/gowa/saas-openclaw/backend/pkg/encryption"
)

type FeishuConfigRepository struct {
    db *sqlx.DB
}

func NewFeishuConfigRepository(db *sqlx.DB) *FeishuConfigRepository {
    return &FeishuConfigRepository{db: db}
}

// Create 创建飞书配置（自动加密 AppSecret）
func (r *FeishuConfigRepository) Create(c *config.FeishuConfig) error

// GetByTenantID 根据租户 ID 获取配置（自动解密 AppSecret）
func (r *FeishuConfigRepository) GetByTenantID(tenantID string) (*config.FeishuConfig, error)

// Update 更新飞书配置（自动加密 AppSecret）
func (r *FeishuConfigRepository) Update(c *config.FeishuConfig) error

// Delete 删除飞书配置
func (r *FeishuConfigRepository) Delete(tenantID string) error
```

### 项目结构规范

**新增文件位置:**

```
backend/
├── migrations/
│   ├── 000004_create_feishu_configs.up.sql    # 新增
│   └── 000004_create_feishu_configs.down.sql  # 新增
├── internal/
│   ├── domain/
│   │   └── config/
│   │       └── feishu_config.go               # 新增
│   └── repository/
│       └── feishu_config.go                   # 新增
├── pkg/
│   └── encryption/
│       ├── aes.go                             # 新增
│       └── aes_test.go                        # 新增
└── .env.example                               # 更新：添加 ENCRYPTION_KEY
```

### 测试标准

**测试要求:**
- 测试框架: Go 原生 testing 包 + testify
- 测试覆盖率目标: >= 70%

**测试用例设计:**

| 测试场景 | 测试文件 | 测试方法 |
|---------|---------|---------|
| 领域模型创建 | `feishu_config_test.go` | 单元测试 |
| 加密解密功能 | `aes_test.go` | 单元测试 |
| Repository CRUD | `feishu_config_repository_test.go` | 集成测试 |
| 唯一约束 | `feishu_config_repository_test.go` | 集成测试 |

### 前序 Story 的学习经验

**从 Story 2.1 (用户数据模型) 获得的经验:**

1. **数据库迁移**: 使用 golang-migrate 工具管理迁移
2. **命名约定**: 表名 snake_case，列名 PascalCase
3. **外键约束**: 确保迁移顺序正确，先创建被引用的表
4. **测试策略**: 使用 testcontainers 进行集成测试

### 安全注意事项

1. **加密密钥管理**:
   - 密钥必须从环境变量读取
   - 密钥长度必须为 32 字节（AES-256）
   - 生产环境密钥应使用密钥管理服务（如 AWS KMS）

2. **敏感字段保护**:
   - AppSecret 字段使用 `json:"-"` 标签
   - 不在日志中打印 AppSecret
   - 解密后的 AppSecret 仅在内存中使用

3. **SQL 注入防护**:
   - 使用参数化查询
   - 使用 sqlx 的 Named Query

### 常见问题与解决方案

**问题 1: 加密密钥长度不正确**
- **原因**: AES-256 要求密钥长度为 32 字节
- **解决**: 启动时验证密钥长度，不正确则拒绝启动

**问题 2: 一租户多配置**
- **原因**: 未添加唯一约束
- **解决**: TenantID 列添加唯一索引

**问题 3: 外键约束失败**
- **原因**: tenants 表不存在或租户 ID 无效
- **解决**: 确保迁移顺序正确，先创建 tenants 表

### References

- [Source: architecture.md#Naming Patterns] - 数据库命名约定
- [Source: architecture.md#Project Structure] - 项目目录结构
- [Source: prd.md#FR5-FR9] - 飞书应用配置需求
- [Source: prd.md#NFR-S2] - 用户配置数据加密存储
- [Source: epics.md#Story 3.1] - 原始故事定义
- [Source: 2-1-user-data-model.md] - 用户数据模型实现参考

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
