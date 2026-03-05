# Story 4.1: 实例数据模型

Status: ready-for-dev

## Story

As a 开发者,
I want 创建 OpenClaw 实例数据模型,
so that 可以存储实例的状态和配置信息。

## Acceptance Criteria

1. **AC1: openclaw_instances 表创建成功**
   - **Given** PostgreSQL 数据库已配置
   - **When** 执行数据库迁移
   - **Then** 创建 `openclaw_instances` 表（表名 snake_case）
   - **And** 表包含 ID、TenantID、Name、Status、ContainerID、DeployLog、CreatedAt、UpdatedAt 列（列名 PascalCase）
   - **And** TenantID 列有外键约束关联到 tenants 表

2. **AC2: 实例状态枚举定义正确**
   - **Given** 数据库表已创建
   - **When** 检查 Status 列定义
   - **Then** Status 包含：pending、deploying、running、stopped、error 状态
   - **And** 默认状态为 pending

3. **AC3: 一个租户可以有多个实例**
   - **Given** 租户已存在
   - **When** 检查外键约束
   - **Then** TenantID 外键允许一个租户关联多条实例记录
   - **And** 删除租户时级联删除相关实例

4. **AC4: 数据库索引优化查询性能**
   - **Given** 表已创建
   - **When** 检查索引配置
   - **Then** openclaw_instances 表有 TenantID 列索引
   - **And** openclaw_instances 表有 Status 列索引
   - **And** openclaw_instances 表有 CreatedAt 列索引

5. **AC5: Go 领域模型定义**
   - **Given** 项目已初始化
   - **When** 创建领域模型文件
   - **Then** 在 `internal/domain/instance/` 定义 Instance 结构体
   - **And** 定义 InstanceStatus 枚举类型
   - **And** 结构体字段使用 `db` tag 映射数据库列
   - **And** 结构体字段使用 `json` tag 定义 JSON 序列化

6. **AC6: Repository 层实现**
   - **Given** 领域模型已定义
   - **When** 创建 Repository 文件
   - **Then** 在 `internal/repository/instance.go` 实现 InstanceRepository 接口
   - **And** 支持 Create、GetByID、GetByTenantID、Update、Delete 方法
   - **And** 使用 sqlx 进行数据库操作

## Tasks / Subtasks

- [ ] Task 1: 创建实例领域模型 (AC: 2, 5)
  - [ ] 1.1 创建 `internal/domain/instance/instance.go` 定义 Instance 结构体
  - [ ] 1.2 创建 `internal/domain/instance/status.go` 定义 InstanceStatus 枚举
  - [ ] 1.3 定义状态常量：InstanceStatusPending、InstanceStatusDeploying、InstanceStatusRunning、InstanceStatusStopped、InstanceStatusError
  - [ ] 1.4 添加 `json` 和 `db` 结构体标签

- [ ] Task 2: 创建数据库迁移脚本 (AC: 1, 3, 4)
  - [ ] 2.1 创建 `migrations/004_create_openclaw_instances.up.sql`
  - [ ] 2.2 创建 `migrations/004_create_openclaw_instances.down.sql`
  - [ ] 2.3 定义 Status 列为 ENUM 类型或 VARCHAR
  - [ ] 2.4 添加外键约束关联 tenants 表
  - [ ] 2.5 添加索引定义

- [ ] Task 3: 实现 Repository 层 (AC: 6)
  - [ ] 3.1 创建 `internal/repository/instance.go`
  - [ ] 3.2 实现 InstanceRepository 结构体和构造函数
  - [ ] 3.3 实现 Create 方法
  - [ ] 3.4 实现 GetByID 方法
  - [ ] 3.5 实现 GetByTenantID 方法
  - [ ] 3.6 实现 GetByStatus 方法
  - [ ] 3.7 实现 Update 方法
  - [ ] 3.8 实现 Delete 方法

- [ ] Task 4: 编写单元测试 (AC: 1-6)
  - [ ] 4.1 编写 `instance_test.go` 测试领域模型
  - [ ] 4.2 编写 `instance_repository_test.go` 测试 Repository
  - [ ] 4.3 测试状态枚举转换
  - [ ] 4.4 测试 CRUD 操作
  - [ ] 4.5 测试索引查询性能

## Dev Notes

### 架构模式与约束

**必须遵循的 Clean Architecture 原则 [Source: architecture.md]:**

1. **依赖方向**:
   - `domain` 层定义实体和接口（无外部依赖）
   - `repository` 层实现数据访问（依赖 domain 和 infrastructure/database）

2. **命名约定 [Source: architecture.md#Naming Patterns]:**
   - 表名: `snake_case` (例: `openclaw_instances`)
   - 列名: `PascalCase` (例: `TenantID`, `Status`, `CreatedAt`)
   - 外键: `table_id` (例: `tenant_id`)

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
│   │   ├── instance/          # ✅ 目录存在（空）
│   │   ├── tenant/            # ✅ 目录存在（空）
│   │   └── user/              # ✅ 目录存在（空）
│   └── repository/            # ✅ 目录存在（空）
```

### 数据库表设计

**openclaw_instances 表:**

```sql
CREATE TABLE openclaw_instances (
    "ID" VARCHAR(36) PRIMARY KEY,  -- UUID
    "TenantID" VARCHAR(36) NOT NULL REFERENCES tenants("ID") ON DELETE CASCADE,
    "Name" VARCHAR(255) NOT NULL,
    "Status" VARCHAR(50) NOT NULL DEFAULT 'pending',
    "ContainerID" VARCHAR(255),
    "DeployLog" TEXT,
    "CreatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "UpdatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX idx_openclaw_instances_tenant_id ON openclaw_instances ("TenantID");
CREATE INDEX idx_openclaw_instances_status ON openclaw_instances ("Status");
CREATE INDEX idx_openclaw_instances_created_at ON openclaw_instances ("CreatedAt");

-- 状态检查约束
ALTER TABLE openclaw_instances ADD CONSTRAINT chk_status
    CHECK ("Status" IN ('pending', 'deploying', 'running', 'stopped', 'error'));
```

### 领域模型设计

**Instance 结构体:**

```go
// internal/domain/instance/instance.go
package instance

import "time"

// InstanceStatus 实例状态类型
type InstanceStatus string

const (
    InstanceStatusPending    InstanceStatus = "pending"
    InstanceStatusDeploying  InstanceStatus = "deploying"
    InstanceStatusRunning    InstanceStatus = "running"
    InstanceStatusStopped    InstanceStatus = "stopped"
    InstanceStatusError      InstanceStatus = "error"
)

// Instance OpenClaw 实例实体
type Instance struct {
    ID          string         `json:"id" db:"ID"`
    TenantID    string         `json:"tenantId" db:"TenantID"`
    Name        string         `json:"name" db:"Name"`
    Status      InstanceStatus `json:"status" db:"Status"`
    ContainerID string         `json:"containerId,omitempty" db:"ContainerID"`
    DeployLog   string         `json:"deployLog,omitempty" db:"DeployLog"`
    CreatedAt   time.Time      `json:"createdAt" db:"CreatedAt"`
    UpdatedAt   time.Time      `json:"updatedAt" db:"UpdatedAt"`
}

// IsValid 检查状态是否有效
func (s InstanceStatus) IsValid() bool {
    switch s {
    case InstanceStatusPending, InstanceStatusDeploying, InstanceStatusRunning,
         InstanceStatusStopped, InstanceStatusError:
        return true
    default:
        return false
    }
}
```

### Repository 接口设计

**InstanceRepository:**

```go
// internal/repository/instance.go
package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/gowa/saas-openclaw/backend/internal/domain/instance"
)

type InstanceRepository struct {
    db *sqlx.DB
}

func NewInstanceRepository(db *sqlx.DB) *InstanceRepository {
    return &InstanceRepository{db: db}
}

func (r *InstanceRepository) Create(i *instance.Instance) error
func (r *InstanceRepository) GetByID(id string) (*instance.Instance, error)
func (r *InstanceRepository) GetByTenantID(tenantID string) ([]*instance.Instance, error)
func (r *InstanceRepository) GetByStatus(status instance.InstanceStatus) ([]*instance.Instance, error)
func (r *InstanceRepository) Update(i *instance.Instance) error
func (r *InstanceRepository) Delete(id string) error
func (r *InstanceRepository) UpdateStatus(id string, status instance.InstanceStatus, log string) error
```

### 项目结构规范

**新增文件位置:**

```
backend/
├── migrations/
│   ├── 000004_create_openclaw_instances.up.sql
│   └── 000004_create_openclaw_instances.down.sql
├── internal/
│   ├── domain/
│   │   └── instance/
│   │       ├── instance.go       # 实例实体（新增）
│   │       └── status.go         # 状态枚举（新增）
│   └── repository/
│       └── instance.go           # 实例仓库（新增）
```

### 测试标准

**测试要求:**
- 测试框架: Go 原生 testing 包 + testify
- 测试覆盖率目标: >= 70%
- 使用 testcontainers 进行集成测试

**测试用例设计:**

| 测试场景 | 测试文件 | 测试方法 |
|---------|---------|---------|
| 领域模型创建 | `instance_test.go` | 单元测试 |
| 状态枚举验证 | `status_test.go` | 单元测试 |
| Repository CRUD | `instance_repository_test.go` | 集成测试 |
| 外键约束 | `instance_repository_test.go` | 集成测试 |

### 前序 Story 的依赖

**依赖 Story 2.1 (用户数据模型):**
- tenants 表需要先创建
- 数据库迁移顺序：tenants -> tenant_users -> admin_users -> openclaw_instances

### 性能要求

**查询优化 [Source: prd.md#NFR-P3]:**
- 实例列表查询 < 2 秒
- 使用索引加速常用查询

### References

- [Source: architecture.md#Naming Patterns] - 数据库命名约定
- [Source: architecture.md#Project Structure] - 项目目录结构
- [Source: prd.md#FR10-FR18] - OpenClaw 实例管理需求
- [Source: prd.md#NFR-P1] - 部署性能要求
- [Source: epics.md#Story 4.1] - 原始故事定义

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
