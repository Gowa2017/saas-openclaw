# Story 6.1: 用户配置自动备份

Status: ready-for-dev

## Story

As a 系统,
I want 自动备份用户配置数据,
so that 用户数据安全有保障。

## Acceptance Criteria

1. **AC1: 配置变更触发备份**
   - **Given** 用户配置数据已存储
   - **When** 配置发生变更（飞书配置、实例配置）
   - **Then** 自动触发备份任务
   - **And** 备份任务异步执行不阻塞主流程
   - **And** 记录备份触发日志

2. **AC2: 定时备份任务**
   - **Given** 系统运行正常
   - **When** 定时任务触发（每日凌晨 2:00）
   - **Then** 执行全量备份所有用户配置
   - **And** 备份包含所有租户的飞书配置
   - **And** 备份包含所有租户的实例配置

3. **AC3: 加密存储备份文件**
   - **Given** 备份任务执行
   - **When** 生成备份文件
   - **Then** 使用 AES-256 加密备份内容
   - **And** 加密密钥从环境变量获取
   - **And** 备份文件格式为 `.backup`

4. **AC4: 备份保留策略**
   - **Given** 备份文件已生成
   - **When** 检查备份文件保留
   - **Then** 保留最近 7 天的备份文件
   - **And** 自动清理超过 7 天的备份
   - **And** 记录清理日志

5. **AC5: 备份成功率监控**
   - **Given** 备份任务执行完成
   - **When** 统计备份成功率
   - **Then** 备份成功率 >= 99%
   - **And** 备份失败时发送告警通知
   - **And** 记录备份状态到数据库

## Tasks / Subtasks

- [ ] Task 1: 创建备份数据模型 (AC: 5)
  - [ ] 1.1 创建 `internal/domain/backup/backup.go` 定义 Backup 结构体
  - [ ] 1.2 创建 `internal/domain/backup/backup_status.go` 定义备份状态枚举
  - [ ] 1.3 创建数据库迁移脚本 `migrations/xxx_create_backups.up.sql`
  - [ ] 1.4 创建回滚脚本 `migrations/xxx_create_backups.down.sql`

- [ ] Task 2: 实现加密工具 (AC: 3)
  - [ ] 2.1 创建 `pkg/crypto/aes.go` 实现 AES-256 加解密
  - [ ] 2.2 实现 Encrypt 函数用于加密数据
  - [ ] 2.3 实现 Decrypt 函数用于解密数据
  - [ ] 2.4 编写单元测试验证加解密功能

- [ ] Task 3: 实现备份存储服务 (AC: 1, 2, 3, 4)
  - [ ] 3.1 创建 `internal/services/backup/storage.go` 备份存储接口
  - [ ] 3.2 实现 S3 存储后端 `internal/services/backup/s3_storage.go`
  - [ ] 3.3 实现本地文件存储后端（开发环境） `internal/services/backup/local_storage.go`
  - [ ] 3.4 实现备份文件命名和路径规范

- [ ] Task 4: 实现配置备份服务 (AC: 1, 2)
  - [ ] 4.1 创建 `internal/services/backup/config_backup.go`
  - [ ] 4.2 实现 BackupFeishuConfigs 方法备份飞书配置
  - [ ] 4.3 实现 BackupInstanceConfigs 方法备份实例配置
  - [ ] 4.4 实现增量备份和全量备份逻辑

- [ ] Task 5: 实现定时任务调度 (AC: 2)
  - [ ] 5.1 创建 `internal/services/scheduler/scheduler.go` 定时任务调度器
  - [ ] 5.2 配置每日凌晨 2:00 执行备份任务
  - [ ] 5.3 实现任务执行日志记录
  - [ ] 5.4 实现任务失败重试机制

- [ ] Task 6: 实现备份清理策略 (AC: 4)
  - [ ] 6.1 创建 `internal/services/backup/cleanup.go`
  - [ ] 6.2 实现按保留天数清理备份文件
  - [ ] 6.3 实现清理日志记录

- [ ] Task 7: 实现备份状态监控 (AC: 5)
  - [ ] 7.1 创建 `internal/repository/backup.go` 备份记录仓库
  - [ ] 7.2 实现备份状态记录和查询方法
  - [ ] 7.3 实现备份成功率统计
  - [ ] 7.4 集成告警系统发送失败通知

- [ ] Task 8: 编写单元测试和集成测试
  - [ ] 8.1 编写加密工具单元测试
  - [ ] 8.2 编写备份服务单元测试
  - [ ] 8.3 编写定时任务集成测试
  - [ ] 8.4 编写清理策略测试

## Dev Notes

### 架构模式与约束

**必须遵循的 Clean Architecture 原则 [Source: architecture.md]:**

1. **依赖方向**:
   - `domain` 层定义 Backup 实体和接口（无外部依赖）
   - `services` 层实现备份业务逻辑（依赖 domain 和 infrastructure）
   - `repository` 层实现数据访问

2. **命名约定 [Source: architecture.md#Naming Patterns]:**
   - 表名: `snake_case` (例: `backup_records`)
   - 列名: `PascalCase` (例: `BackupID`, `CreatedAt`)
   - API 端点: 复数资源名 (例: `/backups`)

3. **备份策略 [Source: architecture.md]:**
   - 使用 Dokploy S3 备份
   - 备份文件加密存储
   - 保留最近 7 天的备份

### 数据库表设计

**backup_records 表:**

```sql
CREATE TABLE backup_records (
    "ID" VARCHAR(36) PRIMARY KEY,  -- UUID
    "Type" VARCHAR(50) NOT NULL,   -- config, instance
    "Status" VARCHAR(50) NOT NULL, -- pending, running, success, failed
    "FilePath" VARCHAR(500),
    "FileSize" BIGINT,
    "StartedAt" TIMESTAMP,
    "CompletedAt" TIMESTAMP,
    "ErrorMessage" TEXT,
    "CreatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_backup_records_type ON backup_records ("Type");
CREATE INDEX idx_backup_records_status ON backup_records ("Status");
CREATE INDEX idx_backup_records_created_at ON backup_records ("CreatedAt");
```

### 领域模型设计

**Backup 结构体:**

```go
// internal/domain/backup/backup.go
package backup

import "time"

type BackupType string

const (
    BackupTypeConfig   BackupType = "config"
    BackupTypeInstance BackupType = "instance"
)

type BackupStatus string

const (
    StatusPending BackupStatus = "pending"
    StatusRunning BackupStatus = "running"
    StatusSuccess BackupStatus = "success"
    StatusFailed  BackupStatus = "failed"
)

type Backup struct {
    ID           string       `json:"id" db:"ID"`
    Type         BackupType   `json:"type" db:"Type"`
    Status       BackupStatus `json:"status" db:"Status"`
    FilePath     string       `json:"filePath" db:"FilePath"`
    FileSize     int64        `json:"fileSize" db:"FileSize"`
    StartedAt    *time.Time   `json:"startedAt" db:"StartedAt"`
    CompletedAt  *time.Time   `json:"completedAt" db:"CompletedAt"`
    ErrorMessage string       `json:"errorMessage" db:"ErrorMessage"`
    CreatedAt    time.Time    `json:"createdAt" db:"CreatedAt"`
}
```

### 加密工具设计

**AES-256 加密:**

```go
// pkg/crypto/aes.go
package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "io"
)

// Encrypt 使用 AES-256-GCM 加密数据
func Encrypt(plaintext string, key []byte) (string, error) {
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

// Decrypt 使用 AES-256-GCM 解密数据
func Decrypt(ciphertext string, key []byte) (string, error) {
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

### 定时任务调度

**使用 robfig/cron 库:**

```go
// internal/services/scheduler/scheduler.go
package scheduler

import (
    "github.com/robfig/cron/v3"
    "log"
)

type Scheduler struct {
    cron *cron.Cron
}

func NewScheduler() *Scheduler {
    return &Scheduler{
        cron: cron.New(cron.WithSeconds()),
    }
}

// AddDailyBackup 添加每日备份任务（凌晨 2:00 执行）
func (s *Scheduler) AddDailyBackup(backupFunc func()) {
    // 每天凌晨 2:00 执行
    s.cron.AddFunc("0 0 2 * * *", func() {
        log.Println("Starting daily backup task...")
        backupFunc()
    })
}

func (s *Scheduler) Start() {
    s.cron.Start()
}

func (s *Scheduler) Stop() {
    s.cron.Stop()
}
```

### S3 存储配置

**环境变量配置:**

```env
# 备份存储配置
BACKUP_STORAGE_TYPE=s3
BACKUP_S3_ENDPOINT=https://s3.example.com
BACKUP_S3_BUCKET=openclaw-backups
BACKUP_S3_ACCESS_KEY=your-access-key
BACKUP_S3_SECRET_KEY=your-secret-key
BACKUP_S3_REGION=cn-east-1

# 加密密钥（32 字节用于 AES-256）
BACKUP_ENCRYPTION_KEY=your-32-byte-encryption-key-here

# 备份保留天数
BACKUP_RETENTION_DAYS=7
```

### 备份文件格式

**备份文件结构:**

```
backups/
├── config/
│   ├── 2026-03-05_020000.backup
│   ├── 2026-03-04_020000.backup
│   └── ...
└── instance/
    ├── 2026-03-05_020000.backup
    ├── 2026-03-04_020000.backup
    └── ...
```

**备份内容格式（JSON 加密前）:**

```json
{
  "version": "1.0",
  "timestamp": "2026-03-05T02:00:00Z",
  "type": "config",
  "data": {
    "feishu_configs": [
      {
        "tenant_id": "xxx",
        "app_id": "cli_xxx",
        "app_secret": "xxx"
      }
    ],
    "instance_configs": [
      {
        "tenant_id": "xxx",
        "instance_id": "xxx",
        "environment_variables": {}
      }
    ]
  }
}
```

### 项目结构规范

**新增文件位置:**

```
backend/
├── internal/
│   ├── domain/
│   │   └── backup/
│   │       ├── backup.go           # 备份实体（新增）
│   │       └── backup_status.go    # 状态枚举（新增）
│   ├── services/
│   │   ├── backup/
│   │   │   ├── storage.go          # 存储接口（新增）
│   │   │   ├── s3_storage.go       # S3 存储（新增）
│   │   │   ├── local_storage.go    # 本地存储（新增）
│   │   │   ├── config_backup.go    # 配置备份（新增）
│   │   │   └── cleanup.go          # 清理策略（新增）
│   │   └── scheduler/
│   │       └── scheduler.go        # 定时调度（新增）
│   └── repository/
│       └── backup.go               # 备份仓库（新增）
├── pkg/
│   └── crypto/
│       └── aes.go                  # 加密工具（新增）
└── migrations/
    ├── xxx_create_backups.up.sql
    └── xxx_create_backups.down.sql
```

### 测试标准

**测试要求:**
- 测试框架: Go 原生 testing 包 + testify
- 测试覆盖率目标: >= 70%
- 使用 testcontainers 进行 S3 集成测试

**测试用例设计:**

| 测试场景 | 测试文件 | 测试方法 |
|---------|---------|---------|
| 加密解密 | `aes_test.go` | 单元测试 |
| 备份创建 | `config_backup_test.go` | 单元测试 |
| 定时任务 | `scheduler_test.go` | 集成测试 |
| 清理策略 | `cleanup_test.go` | 单元测试 |

### 安全注意事项

1. **加密密钥管理**: 加密密钥从环境变量获取，不要硬编码
2. **敏感数据**: 备份文件包含敏感配置（App Secret），必须加密存储
3. **访问控制**: 备份文件存储路径需要权限控制
4. **传输安全**: S3 上传使用 HTTPS

### 前序依赖

**依赖以下 Story:**
- Story 2.1: 用户数据模型（已创建 tenant_users 表）
- Story 3.1: 飞书配置数据模型（已创建 feishu_configs 表）
- Story 4.1: 实例数据模型（已创建 openclaw_instances 表）

### References

- [Source: architecture.md#Backup Strategy] - 备份策略
- [Source: prd.md#FR24] - 备份用户配置数据需求
- [Source: prd.md#NFR-R4] - 备份成功率 >= 99%
- [Source: epics.md#Story 6.1] - 原始故事定义

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
