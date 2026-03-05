# Story 6.2: OpenClaw 实例数据备份

Status: ready-for-dev

## Story

As a 系统,
I want 自动备份 OpenClaw 实例数据,
so that 实例数据不会丢失。

## Acceptance Criteria

1. **AC1: 定时备份触发**
   - **Given** OpenClaw 实例正在运行
   - **When** 定时任务触发（每日）
   - **Then** 系统识别所有运行中的实例
   - **And** 为每个实例创建备份任务
   - **And** 任务按租户隔离执行

2. **AC2: Dokploy S3 备份接口调用**
   - **Given** 实例备份任务已创建
   - **When** 调用 Dokploy S3 备份接口
   - **Then** 正确传递实例 ID 和备份参数
   - **And** 接口调用超时设置为 30 秒
   - **And** 实现重试机制（最多 3 次）

3. **AC3: Docker Volume 数据备份**
   - **Given** 实例有 Docker Volume 数据
   - **When** 执行备份操作
   - **Then** 备份实例的 Volume 数据
   - **And** 备份文件存储到 S3 兼容存储
   - **And** 备份文件命名包含实例 ID 和时间戳

4. **AC4: 备份状态和日志记录**
   - **Given** 备份任务执行
   - **When** 任务状态变更
   - **Then** 记录备份状态（pending, running, success, failed）
   - **And** 记录详细备份日志
   - **And** 日志包含开始时间、结束时间、文件大小

5. **AC5: 备份失败告警**
   - **Given** 备份任务执行失败
   - **When** 检测到备份失败
   - **Then** 发送告警通知给管理员
   - **And** 告警包含实例 ID、租户 ID、错误信息
   - **And** 记录失败原因到数据库

## Tasks / Subtasks

- [ ] Task 1: 创建实例备份数据模型 (AC: 4)
  - [ ] 1.1 创建 `internal/domain/backup/instance_backup.go` 定义 InstanceBackup 结构体
  - [ ] 1.2 定义备份状态枚举和错误类型
  - [ ] 1.3 创建数据库迁移脚本 `migrations/xxx_create_instance_backups.up.sql`
  - [ ] 1.4 创建回滚脚本 `migrations/xxx_create_instance_backups.down.sql`

- [ ] Task 2: 扩展 Dokploy API 客户端 (AC: 2)
  - [ ] 2.1 在 `internal/infrastructure/dokploy/client.go` 添加备份接口
  - [ ] 2.2 实现 `CreateBackup` 方法调用 Dokploy 备份 API
  - [ ] 2.3 实现 `GetBackupStatus` 方法查询备份状态
  - [ ] 2.4 实现 `ListBackups` 方法获取备份列表
  - [ ] 2.5 添加接口调用的错误处理和重试逻辑

- [ ] Task 3: 实现实例备份服务 (AC: 1, 2, 3)
  - [ ] 3.1 创建 `internal/services/backup/instance_backup.go`
  - [ ] 3.2 实现 `BackupInstance` 方法备份单个实例
  - [ ] 3.3 实现 `BackupAllInstances` 方法备份所有实例
  - [ ] 3.4 实现 Volume 数据备份逻辑
  - [ ] 3.5 实现备份文件命名规范

- [ ] Task 4: 实现备份状态追踪 (AC: 4)
  - [ ] 4.1 创建 `internal/repository/instance_backup.go`
  - [ ] 4.2 实现备份记录 CRUD 方法
  - [ ] 4.3 实现备份日志记录方法
  - [ ] 4.4 实现备份状态查询方法

- [ ] Task 5: 实现备份失败告警 (AC: 5)
  - [ ] 5.1 集成告警服务发送备份失败通知
  - [ ] 5.2 实现告警消息模板
  - [ ] 5.3 实现告警记录存储

- [ ] Task 6: 集成定时任务调度 (AC: 1)
  - [ ] 6.1 在调度器中添加实例备份任务
  - [ ] 6.2 配置每日凌晨 3:00 执行（在配置备份后）
  - [ ] 6.3 实现任务执行监控

- [ ] Task 7: 编写单元测试和集成测试
  - [ ] 7.1 编写 Dokploy 备份接口调用测试
  - [ ] 7.2 编写实例备份服务单元测试
  - [ ] 7.3 编写备份状态追踪测试
  - [ ] 7.4 编写定时任务集成测试

## Dev Notes

### 架构模式与约束

**必须遵循的 Clean Architecture 原则 [Source: architecture.md]:**

1. **依赖方向**:
   - `domain` 层定义 InstanceBackup 实体
   - `infrastructure/dokploy` 层封装 Dokploy API 调用
   - `services` 层实现备份业务逻辑

2. **备份策略 [Source: architecture.md]:**
   - 使用 Dokploy S3 备份接口
   - 备份 Docker Volume 数据
   - 备份到 S3 兼容存储

3. **命名约定 [Source: architecture.md#Naming Patterns]:**
   - 表名: `snake_case` (例: `instance_backup_records`)
   - 列名: `PascalCase` (例: `InstanceID`, `BackupStatus`)

### 数据库表设计

**instance_backup_records 表:**

```sql
CREATE TABLE instance_backup_records (
    "ID" VARCHAR(36) PRIMARY KEY,  -- UUID
    "InstanceID" VARCHAR(36) NOT NULL,
    "TenantID" VARCHAR(36) NOT NULL,
    "DokployBackupID" VARCHAR(100),  -- Dokploy 返回的备份 ID
    "Status" VARCHAR(50) NOT NULL,   -- pending, running, success, failed
    "BackupType" VARCHAR(50) NOT NULL, -- full, volume
    "StoragePath" VARCHAR(500),
    "FileSize" BIGINT,
    "StartedAt" TIMESTAMP,
    "CompletedAt" TIMESTAMP,
    "ErrorMessage" TEXT,
    "CreatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_instance_backup_records_instance_id ON instance_backup_records ("InstanceID");
CREATE INDEX idx_instance_backup_records_tenant_id ON instance_backup_records ("TenantID");
CREATE INDEX idx_instance_backup_records_status ON instance_backup_records ("Status");
CREATE INDEX idx_instance_backup_records_created_at ON instance_backup_records ("CreatedAt");
```

### 领域模型设计

**InstanceBackup 结构体:**

```go
// internal/domain/backup/instance_backup.go
package backup

import "time"

type BackupType string

const (
    BackupTypeFull   BackupType = "full"
    BackupTypeVolume BackupType = "volume"
)

type InstanceBackup struct {
    ID              string       `json:"id" db:"ID"`
    InstanceID      string       `json:"instanceId" db:"InstanceID"`
    TenantID        string       `json:"tenantId" db:"TenantID"`
    DokployBackupID string       `json:"dokployBackupId" db:"DokployBackupID"`
    Status          BackupStatus `json:"status" db:"Status"`
    BackupType      BackupType   `json:"backupType" db:"BackupType"`
    StoragePath     string       `json:"storagePath" db:"StoragePath"`
    FileSize        int64        `json:"fileSize" db:"FileSize"`
    StartedAt       *time.Time   `json:"startedAt" db:"StartedAt"`
    CompletedAt     *time.Time   `json:"completedAt" db:"CompletedAt"`
    ErrorMessage    string       `json:"errorMessage" db:"ErrorMessage"`
    CreatedAt       time.Time    `json:"createdAt" db:"CreatedAt"`
}
```

### Dokploy 备份 API 集成

**Dokploy API 客户端扩展:**

```go
// internal/infrastructure/dokploy/backup.go
package dokploy

import (
    "context"
    "fmt"
    "time"
)

type BackupRequest struct {
    AppID      string `json:"appId"`
    BackupType string `json:"backupType"` // "volume", "database"
    Prefix     string `json:"prefix"`     // 备份文件前缀
}

type BackupResponse struct {
    BackupID   string `json:"backupId"`
    Status     string `json:"status"`
    StorageKey string `json:"storageKey"`
    CreatedAt  string `json:"createdAt"`
}

type BackupStatusResponse struct {
    BackupID    string `json:"backupId"`
    Status      string `json:"status"` // "pending", "running", "completed", "failed"
    Progress    int    `json:"progress"`
    FileSize    int64  `json:"fileSize"`
    StoragePath string `json:"storagePath"`
    Error       string `json:"error"`
}

// CreateBackup 调用 Dokploy 创建备份
func (c *Client) CreateBackup(ctx context.Context, req *BackupRequest) (*BackupResponse, error) {
    endpoint := fmt.Sprintf("/api/application/%s/backup", req.AppID)

    var resp BackupResponse
    if err := c.post(ctx, endpoint, req, &resp); err != nil {
        return nil, fmt.Errorf("failed to create backup: %w", err)
    }

    return &resp, nil
}

// GetBackupStatus 查询备份状态
func (c *Client) GetBackupStatus(ctx context.Context, appID, backupID string) (*BackupStatusResponse, error) {
    endpoint := fmt.Sprintf("/api/application/%s/backup/%s", appID, backupID)

    var resp BackupStatusResponse
    if err := c.get(ctx, endpoint, &resp); err != nil {
        return nil, fmt.Errorf("failed to get backup status: %w", err)
    }

    return &resp, nil
}

// ListBackups 获取备份列表
func (c *Client) ListBackups(ctx context.Context, appID string) ([]*BackupResponse, error) {
    endpoint := fmt.Sprintf("/api/application/%s/backups", appID)

    var resp []*BackupResponse
    if err := c.get(ctx, endpoint, &resp); err != nil {
        return nil, fmt.Errorf("failed to list backups: %w", err)
    }

    return resp, nil
}
```

### 实例备份服务实现

```go
// internal/services/backup/instance_backup.go
package backup

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/google/uuid"
)

type InstanceBackupService struct {
    dokployClient  *dokploy.Client
    backupRepo     *repository.InstanceBackupRepository
    instanceRepo   *repository.InstanceRepository
    alertService   *alert.Service
}

func NewInstanceBackupService(
    dokployClient *dokploy.Client,
    backupRepo *repository.InstanceBackupRepository,
    instanceRepo *repository.InstanceRepository,
    alertService *alert.Service,
) *InstanceBackupService {
    return &InstanceBackupService{
        dokployClient:  dokployClient,
        backupRepo:     backupRepo,
        instanceRepo:   instanceRepo,
        alertService:   alertService,
    }
}

// BackupAllInstances 备份所有运行中的实例
func (s *InstanceBackupService) BackupAllInstances(ctx context.Context) error {
    // 获取所有运行中的实例
    instances, err := s.instanceRepo.GetByStatus(ctx, "running")
    if err != nil {
        return fmt.Errorf("failed to get running instances: %w", err)
    }

    log.Printf("Found %d running instances to backup", len(instances))

    for _, instance := range instances {
        go func(inst *domain.Instance) {
            if err := s.BackupInstance(context.Background(), inst); err != nil {
                log.Printf("Failed to backup instance %s: %v", inst.ID, err)
            }
        }(instance)
    }

    return nil
}

// BackupInstance 备份单个实例
func (s *InstanceBackupService) BackupInstance(ctx context.Context, instance *domain.Instance) error {
    // 创建备份记录
    backup := &domain.InstanceBackup{
        ID:         uuid.New().String(),
        InstanceID: instance.ID,
        TenantID:   instance.TenantID,
        Status:     domain.StatusPending,
        BackupType: domain.BackupTypeVolume,
        CreatedAt:  time.Now(),
    }

    if err := s.backupRepo.Create(ctx, backup); err != nil {
        return fmt.Errorf("failed to create backup record: %w", err)
    }

    // 更新状态为运行中
    now := time.Now()
    backup.Status = domain.StatusRunning
    backup.StartedAt = &now
    s.backupRepo.Update(ctx, backup)

    // 调用 Dokploy API 创建备份
    req := &dokploy.BackupRequest{
        AppID:      instance.DokployAppID,
        BackupType: "volume",
        Prefix:     fmt.Sprintf("tenant-%s-instance-%s", instance.TenantID, instance.ID),
    }

    resp, err := s.dokployClient.CreateBackup(ctx, req)
    if err != nil {
        return s.handleBackupError(ctx, backup, err)
    }

    backup.DokployBackupID = resp.BackupID

    // 轮询备份状态
    if err := s.pollBackupStatus(ctx, backup); err != nil {
        return s.handleBackupError(ctx, backup, err)
    }

    return nil
}

// pollBackupStatus 轮询备份状态直到完成
func (s *InstanceBackupService) pollBackupStatus(ctx context.Context, backup *domain.InstanceBackup) error {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    timeout := time.After(30 * time.Minute)

    for {
        select {
        case <-ticker.C:
            status, err := s.dokployClient.GetBackupStatus(ctx, backup.InstanceID, backup.DokployBackupID)
            if err != nil {
                log.Printf("Failed to get backup status: %v", err)
                continue
            }

            if status.Status == "completed" {
                now := time.Now()
                backup.Status = domain.StatusSuccess
                backup.StoragePath = status.StoragePath
                backup.FileSize = status.FileSize
                backup.CompletedAt = &now
                return s.backupRepo.Update(ctx, backup)
            }

            if status.Status == "failed" {
                backup.Status = domain.StatusFailed
                backup.ErrorMessage = status.Error
                return s.backupRepo.Update(ctx, backup)
            }

        case <-timeout:
            return fmt.Errorf("backup timeout after 30 minutes")
        }
    }
}

// handleBackupError 处理备份错误
func (s *InstanceBackupService) handleBackupError(ctx context.Context, backup *domain.InstanceBackup, err error) error {
    backup.Status = domain.StatusFailed
    backup.ErrorMessage = err.Error()
    now := time.Now()
    backup.CompletedAt = &now
    s.backupRepo.Update(ctx, backup)

    // 发送告警
    s.alertService.SendBackupFailureAlert(ctx, backup)

    return err
}
```

### 环境变量配置

```env
# Dokploy 备份配置
DOKPLOY_BACKUP_TIMEOUT=30m
DOKPLOY_BACKUP_RETRY_COUNT=3
DOKPLOY_BACKUP_POLL_INTERVAL=5s

# S3 备份存储配置
BACKUP_S3_BUCKET=openclaw-instance-backups
BACKUP_S3_PREFIX=backups/instances
```

### 备份文件命名规范

```
backups/instances/
├── tenant-{tenant_id}/
│   ├── instance-{instance_id}/
│   │   ├── volume-{timestamp}.tar.gz
│   │   └── volume-{timestamp}.tar.gz
│   └── ...
└── ...
```

### 定时任务配置

```go
// 在调度器中添加实例备份任务
func (s *Scheduler) AddInstanceBackup(backupFunc func()) {
    // 每天凌晨 3:00 执行（在配置备份后 1 小时）
    s.cron.AddFunc("0 0 3 * * *", func() {
        log.Println("Starting instance backup task...")
        backupFunc()
    })
}
```

### 项目结构规范

**新增文件位置:**

```
backend/
├── internal/
│   ├── domain/
│   │   └── backup/
│   │       └── instance_backup.go    # 实例备份实体（新增）
│   ├── infrastructure/
│   │   └── dokploy/
│   │       └── backup.go             # 备份 API（新增）
│   ├── services/
│   │   └── backup/
│   │       └── instance_backup.go    # 实例备份服务（新增）
│   └── repository/
│       └── instance_backup.go        # 备份仓库（新增）
└── migrations/
    ├── xxx_create_instance_backups.up.sql
    └── xxx_create_instance_backups.down.sql
```

### 测试标准

**测试要求:**
- 测试框架: Go 原生 testing 包 + testify
- 测试覆盖率目标: >= 70%
- 使用 mock 测试 Dokploy API 调用

**测试用例设计:**

| 测试场景 | 测试文件 | 测试方法 |
|---------|---------|---------|
| Dokploy API 调用 | `backup_test.go` | 单元测试（mock） |
| 备份状态轮询 | `instance_backup_test.go` | 单元测试 |
| 备份失败处理 | `instance_backup_test.go` | 单元测试 |
| 定时任务执行 | `scheduler_test.go` | 集成测试 |

### 前序依赖

**依赖以下 Story:**
- Story 4.1: 实例数据模型（已创建 openclaw_instances 表）
- Story 4.2: Dokploy API 客户端集成（已实现基础 API）
- Story 6.1: 用户配置自动备份（已实现备份基础设施）

### 安全注意事项

1. **租户隔离**: 备份文件按租户 ID 分隔存储
2. **访问控制**: 只有管理员可以触发备份操作
3. **错误处理**: 备份失败不影响实例正常运行
4. **超时控制**: 备份任务设置超时防止无限等待

### References

- [Source: architecture.md#Backup Strategy] - 备份策略
- [Source: prd.md#FR25] - 备份 OpenClaw 实例数据需求
- [Source: prd.md#NFR-R4] - 备份成功率 >= 99%
- [Source: epics.md#Story 6.2] - 原始故事定义

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
