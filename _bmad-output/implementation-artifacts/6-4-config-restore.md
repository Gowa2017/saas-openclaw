# Story 6.4: 用户配置恢复功能

Status: ready-for-dev

## Story

As a 系统,
I want 在实例重新部署时恢复用户配置,
so that 用户无需重新配置。

## Acceptance Criteria

1. **AC1: 自动检测恢复场景**
   - **Given** 用户需要重新部署实例
   - **When** 触发实例重新部署
   - **Then** 系统检测是否存在备份配置
   - **And** 自动触发配置恢复流程
   - **And** 记录恢复触发日志

2. **AC2: 飞书配置恢复**
   - **Given** 存在用户的飞书配置备份
   - **When** 执行配置恢复
   - **Then** 从备份中恢复飞书 App ID
   - **And** 从备份中恢复飞书 App Secret
   - **And** 配置正确写入新实例环境变量

3. **AC3: 配置恢复到新实例**
   - **Given** 新实例已创建
   - **When** 执行配置恢复
   - **Then** 配置恢复到新实例的环境变量
   - **And** 实例使用恢复的配置启动
   - **And** 恢复操作在部署过程中完成

4. **AC4: 恢复失败处理**
   - **Given** 配置恢复过程中发生错误
   - **When** 检测到恢复失败
   - **Then** 通知用户手动配置
   - **And** 提供手动配置指引
   - **And** 不影响实例正常部署

5. **AC5: 恢复操作日志**
   - **Given** 恢复操作执行
   - **When** 操作状态变更
   - **Then** 记录恢复操作日志
   - **And** 日志包含恢复时间、状态、详情
   - **And** 管理员可查看恢复日志

## Tasks / Subtasks

- [ ] Task 1: 创建配置恢复数据模型 (AC: 5)
  - [ ] 1.1 创建 `internal/domain/backup/config_restore.go` 定义 ConfigRestore 结构体
  - [ ] 1.2 定义恢复状态枚举
  - [ ] 1.3 创建数据库迁移脚本 `migrations/xxx_create_config_restores.up.sql`
  - [ ] 1.4 创建回滚脚本 `migrations/xxx_create_config_restores.down.sql`

- [ ] Task 2: 实现备份文件解密服务 (AC: 2)
  - [ ] 2.1 创建 `internal/services/backup/restore.go` 恢复服务
  - [ ] 2.2 实现备份文件读取功能
  - [ ] 2.3 实现备份文件解密功能
  - [ ] 2.4 实现备份数据解析功能

- [ ] Task 3: 实现飞书配置恢复逻辑 (AC: 2)
  - [ ] 3.1 创建 `RestoreFeishuConfig` 方法
  - [ ] 3.2 从备份中提取飞书配置
  - [ ] 3.3 验证配置完整性
  - [ ] 3.4 写入新实例环境变量

- [ ] Task 4: 集成部署流程 (AC: 1, 3)
  - [ ] 4.1 修改部署服务添加恢复检测
  - [ ] 4.2 在部署流程中集成配置恢复
  - [ ] 4.3 实现配置恢复时机控制
  - [ ] 4.4 实现恢复与部署的协调

- [ ] Task 5: 实现恢复失败处理 (AC: 4)
  - [ ] 5.1 实现恢复失败检测
  - [ ] 5.2 实现用户通知服务
  - [ ] 5.3 创建手动配置指引模板
  - [ ] 5.4 实现恢复失败不影响部署流程

- [ ] Task 6: 实现恢复日志记录 (AC: 5)
  - [ ] 6.1 创建 `internal/repository/config_restore.go`
  - [ ] 6.2 实现恢复记录 CRUD
  - [ ] 6.3 实现恢复日志记录
  - [ ] 6.4 实现日志查询接口

- [ ] Task 7: 编写单元测试和集成测试
  - [ ] 7.1 编写备份解密测试
  - [ ] 7.2 编写配置恢复测试
  - [ ] 7.3 编写恢复失败处理测试
  - [ ] 7.4 编写部署集成测试

## Dev Notes

### 架构模式与约束

**必须遵循的 Clean Architecture 原则 [Source: architecture.md]:**

1. **依赖方向**:
   - `domain` 层定义 ConfigRestore 实体
   - `services` 层实现恢复业务逻辑
   - `repository` 层实现恢复记录存储

2. **恢复策略:**
   - 在实例重新部署时自动恢复
   - 恢复失败不影响部署流程
   - 提供手动配置备选方案

3. **命名约定 [Source: architecture.md#Naming Patterns]:**
   - 表名: `snake_case` (例: `config_restore_records`)
   - 列名: `PascalCase` (例: `InstanceID`, `RestoreStatus`)

### 数据库表设计

**config_restore_records 表:**

```sql
CREATE TABLE config_restore_records (
    "ID" VARCHAR(36) PRIMARY KEY,  -- UUID
    "TenantID" VARCHAR(36) NOT NULL,
    "InstanceID" VARCHAR(36) NOT NULL,
    "BackupID" VARCHAR(36) NOT NULL,
    "Status" VARCHAR(50) NOT NULL,   -- pending, running, success, failed
    "RestoredConfigs" TEXT,          -- JSON: 恢复的配置列表
    "ErrorMessage" TEXT,
    "StartedAt" TIMESTAMP,
    "CompletedAt" TIMESTAMP,
    "CreatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_config_restore_records_tenant_id ON config_restore_records ("TenantID");
CREATE INDEX idx_config_restore_records_instance_id ON config_restore_records ("InstanceID");
CREATE INDEX idx_config_restore_records_status ON config_restore_records ("Status");
```

### 领域模型设计

**ConfigRestore 结构体:**

```go
// internal/domain/backup/config_restore.go
package backup

import "time"

type RestoreStatus string

const (
    RestoreStatusPending RestoreStatus = "pending"
    RestoreStatusRunning RestoreStatus = "running"
    RestoreStatusSuccess RestoreStatus = "success"
    RestoreStatusFailed  RestoreStatus = "failed"
)

type RestoredConfig struct {
    Type    string `json:"type"`    // feishu_app_id, feishu_app_secret
    Key     string `json:"key"`     // 环境变量名
    Success bool   `json:"success"` // 是否恢复成功
}

type ConfigRestore struct {
    ID              string          `json:"id" db:"ID"`
    TenantID        string          `json:"tenantId" db:"TenantID"`
    InstanceID      string          `json:"instanceId" db:"InstanceID"`
    BackupID        string          `json:"backupId" db:"BackupID"`
    Status          RestoreStatus   `json:"status" db:"Status"`
    RestoredConfigs []RestoredConfig `json:"restoredConfigs" db:"RestoredConfigs"`
    ErrorMessage    string          `json:"errorMessage" db:"ErrorMessage"`
    StartedAt       *time.Time      `json:"startedAt" db:"StartedAt"`
    CompletedAt     *time.Time      `json:"completedAt" db:"CompletedAt"`
    CreatedAt       time.Time       `json:"createdAt" db:"CreatedAt"`
}

func (r *ConfigRestore) IsSuccess() bool {
    return r.Status == RestoreStatusSuccess
}

func (r *ConfigRestore) GetFailedConfigs() []RestoredConfig {
    var failed []RestoredConfig
    for _, c := range r.RestoredConfigs {
        if !c.Success {
            failed = append(failed, c)
        }
    }
    return failed
}
```

### 恢复服务实现

```go
// internal/services/backup/restore.go
package backup

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/google/uuid"
)

type RestoreService struct {
    backupRepo     *repository.BackupRepository
    restoreRepo    *repository.ConfigRestoreRepository
    feishuRepo     *repository.FeishuConfigRepository
    instanceRepo   *repository.InstanceRepository
    dokployClient  *dokploy.Client
    cryptoService  *crypto.Service
    notifyService  *notify.Service
}

func NewRestoreService(
    backupRepo *repository.BackupRepository,
    restoreRepo *repository.ConfigRestoreRepository,
    feishuRepo *repository.FeishuConfigRepository,
    instanceRepo *repository.InstanceRepository,
    dokployClient *dokploy.Client,
    cryptoService *crypto.Service,
    notifyService *notify.Service,
) *RestoreService {
    return &RestoreService{
        backupRepo:    backupRepo,
        restoreRepo:   restoreRepo,
        feishuRepo:    feishuRepo,
        instanceRepo:  instanceRepo,
        dokployClient: dokployClient,
        cryptoService: cryptoService,
        notifyService: notifyService,
    }
}

// RestoreForRedeploy 在重新部署时恢复配置
func (s *RestoreService) RestoreForRedeploy(ctx context.Context, tenantID, instanceID, dokployAppID string) error {
    // 创建恢复记录
    restore := &domain.ConfigRestore{
        ID:         uuid.New().String(),
        TenantID:   tenantID,
        InstanceID: instanceID,
        Status:     domain.RestoreStatusPending,
        CreatedAt:  time.Now(),
    }

    if err := s.restoreRepo.Create(ctx, restore); err != nil {
        log.Printf("Failed to create restore record: %v", err)
    }

    // 查找最新的配置备份
    backup, err := s.backupRepo.GetLatestSuccessBackup(ctx, tenantID, "config")
    if err != nil {
        return s.handleRestoreError(ctx, restore, fmt.Errorf("no backup found: %w", err))
    }

    restore.BackupID = backup.ID

    // 开始恢复
    now := time.Now()
    restore.Status = domain.RestoreStatusRunning
    restore.StartedAt = &now
    s.restoreRepo.Update(ctx, restore)

    // 执行恢复
    if err := s.executeRestore(ctx, restore, backup, dokployAppID); err != nil {
        return s.handleRestoreError(ctx, restore, err)
    }

    // 恢复成功
    now = time.Now()
    restore.Status = domain.RestoreStatusSuccess
    restore.CompletedAt = &now
    s.restoreRepo.Update(ctx, restore)

    log.Printf("Config restore completed successfully for instance %s", instanceID)
    return nil
}

// executeRestore 执行实际的恢复操作
func (s *RestoreService) executeRestore(ctx context.Context, restore *domain.ConfigRestore, backup *domain.Backup, dokployAppID string) error {
    // 读取备份文件
    backupData, err := s.readBackupFile(ctx, backup.FilePath)
    if err != nil {
        return fmt.Errorf("failed to read backup file: %w", err)
    }

    // 解密备份数据
    decryptedData, err := s.cryptoService.Decrypt(backupData)
    if err != nil {
        return fmt.Errorf("failed to decrypt backup: %w", err)
    }

    // 解析备份数据
    var backupContent BackupContent
    if err := json.Unmarshal([]byte(decryptedData), &backupContent); err != nil {
        return fmt.Errorf("failed to parse backup: %w", err)
    }

    // 恢复飞书配置
    restoredConfigs := s.restoreFeishuConfig(ctx, restore.TenantID, backupContent, dokployAppID)
    restore.RestoredConfigs = restoredConfigs

    return nil
}

// restoreFeishuConfig 恢复飞书配置到新实例
func (s *RestoreService) restoreFeishuConfig(ctx context.Context, tenantID string, content BackupContent, dokployAppID string) []domain.RestoredConfig {
    var configs []domain.RestoredConfig

    for _, fc := range content.FeishuConfigs {
        if fc.TenantID != tenantID {
            continue
        }

        // 设置环境变量
        envVars := map[string]string{
            "FEISHU_APP_ID":     fc.AppID,
            "FEISHU_APP_SECRET": fc.AppSecret,
        }

        if err := s.dokployClient.SetEnvironmentVariables(ctx, dokployAppID, envVars); err != nil {
            log.Printf("Failed to set env vars: %v", err)
            configs = append(configs, domain.RestoredConfig{
                Type:    "feishu",
                Key:     "FEISHU_APP_ID/SECRET",
                Success: false,
            })
            continue
        }

        configs = append(configs, domain.RestoredConfig{
            Type:    "feishu",
            Key:     "FEISHU_APP_ID/SECRET",
            Success: true,
        })

        log.Printf("Restored feishu config for tenant %s", tenantID)
    }

    return configs
}

// handleRestoreError 处理恢复错误
func (s *RestoreService) handleRestoreError(ctx context.Context, restore *domain.ConfigRestore, err error) error {
    now := time.Now()
    restore.Status = domain.RestoreStatusFailed
    restore.ErrorMessage = err.Error()
    restore.CompletedAt = &now
    s.restoreRepo.Update(ctx, restore)

    // 通知用户手动配置
    s.notifyService.SendManualConfigNotification(ctx, restore.TenantID, restore.InstanceID)

    log.Printf("Config restore failed: %v", err)
    // 返回 nil 不影响部署流程
    return nil
}

// BackupContent 备份内容结构
type BackupContent struct {
    Version       string           `json:"version"`
    Timestamp     string           `json:"timestamp"`
    Type          string           `json:"type"`
    FeishuConfigs []FeishuConfig   `json:"feishu_configs"`
    Instances     []InstanceConfig `json:"instances"`
}

type FeishuConfig struct {
    TenantID  string `json:"tenant_id"`
    AppID     string `json:"app_id"`
    AppSecret string `json:"app_secret"`
}

type InstanceConfig struct {
    TenantID            string            `json:"tenant_id"`
    InstanceID          string            `json:"instance_id"`
    EnvironmentVariables map[string]string `json:"environment_variables"`
}
```

### 部署流程集成

```go
// 在部署服务中集成配置恢复
// internal/services/deploy/deploy.go

func (s *DeployService) DeployInstance(ctx context.Context, req *DeployRequest) (*DeployResult, error) {
    // ... 创建实例 ...

    // 检查是否需要恢复配置（重新部署场景）
    if req.IsRedeploy {
        go func() {
            if err := s.restoreService.RestoreForRedeploy(
                context.Background(),
                req.TenantID,
                instance.ID,
                instance.DokployAppID,
            ); err != nil {
                log.Printf("Config restore failed: %v", err)
            }
        }()
    }

    // ... 继续部署流程 ...
}
```

### 用户通知模板

```go
// 手动配置通知模板
const ManualConfigNotificationTemplate = `
亲爱的用户，

您的 OpenClaw 实例已成功部署，但配置恢复失败。请手动完成以下配置：

1. 登录 SaaS 平台
2. 进入"飞书配置"页面
3. 填写飞书 App ID 和 App Secret
4. 点击"保存配置"

如果需要帮助，请联系客服。

感谢您的理解！
`
```

### 环境变量设置

```go
// Dokploy 环境变量设置
func (c *Client) SetEnvironmentVariables(ctx context.Context, appID string, vars map[string]string) error {
    endpoint := fmt.Sprintf("/api/application/%s/environment", appID)

    req := struct {
        Variables map[string]string `json:"variables"`
    }{
        Variables: vars,
    }

    return c.put(ctx, endpoint, req, nil)
}
```

### 项目结构规范

**新增文件位置:**

```
backend/
├── internal/
│   ├── domain/
│   │   └── backup/
│   │       └── config_restore.go     # 配置恢复实体（新增）
│   ├── services/
│   │   ├── backup/
│   │   │   └── restore.go            # 恢复服务（新增）
│   │   └── notify/
│   │       └── notification.go       # 通知服务（新增）
│   └── repository/
│       └── config_restore.go         # 恢复仓库（新增）
└── migrations/
    ├── xxx_create_config_restores.up.sql
    └── xxx_create_config_restores.down.sql
```

### 测试标准

**测试要求:**
- 测试框架: Go 原生 testing 包 + testify
- 测试覆盖率目标: >= 70%
- 使用 mock 测试外部依赖

**测试用例设计:**

| 测试场景 | 测试文件 | 测试方法 |
|---------|---------|---------|
| 备份文件解密 | `restore_test.go` | 单元测试 |
| 配置恢复成功 | `restore_test.go` | 单元测试 |
| 恢复失败处理 | `restore_test.go` | 单元测试 |
| 部署集成 | `deploy_test.go` | 集成测试 |

### 前序依赖

**依赖以下 Story:**
- Story 4.3: 一键部署 OpenClaw 实例（已实现部署流程）
- Story 6.1: 用户配置自动备份（已实现备份创建）
- Story 6.2: OpenClaw 实例数据备份（已实现备份存储）

### 安全注意事项

1. **敏感数据处理**: 恢复的配置包含 App Secret，需要安全处理
2. **操作审计**: 记录所有恢复操作日志
3. **错误隔离**: 恢复失败不影响实例部署
4. **用户通知**: 恢复失败时及时通知用户

### 流程图

```
实例重新部署触发
       │
       ▼
  检查是否存在备份 ──────── 不存在 ──► 正常部署（无恢复）
       │
       存在
       ▼
  创建恢复记录
       │
       ▼
  读取最新备份文件
       │
       ▼
  解密备份数据
       │
       ▼
  解析配置内容
       │
       ▼
  设置实例环境变量 ──────── 失败 ──► 通知用户手动配置
       │
       成功
       ▼
  更新恢复记录
       │
       ▼
  启动实例
```

### References

- [Source: architecture.md#Backup Strategy] - 备份策略
- [Source: prd.md#FR27] - 恢复用户配置需求
- [Source: epics.md#Story 6.4] - 原始故事定义

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
