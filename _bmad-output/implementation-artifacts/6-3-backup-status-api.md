# Story 6.3: 备份状态查看 API

Status: ready-for-dev

## Story

As a 平台管理员,
I want 查看系统备份状态,
so that 确认备份任务正常运行。

## Acceptance Criteria

1. **AC1: 备份任务列表 API**
   - **Given** 管理员已登录管理后台
   - **When** 调用 GET /v1/admin/backups API
   - **Then** 返回最近备份任务列表
   - **And** 每个任务显示：ID、时间、类型、状态、大小
   - **And** 支持分页查询（默认 20 条）

2. **AC2: 日期范围筛选**
   - **Given** 管理员访问备份列表
   - **When** 指定日期范围参数
   - **Then** 返回指定范围内的备份任务
   - **And** 支持开始日期和结束日期筛选
   - **And** 日期格式为 YYYY-MM-DD

3. **AC3: 备份详情查看**
   - **Given** 管理员查看备份列表
   - **When** 点击某个备份任务
   - **Then** 显示备份详细信息
   - **And** 显示备份日志
   - **And** 显示备份文件存储位置

4. **AC4: 备份存储空间使用情况**
   - **Given** 管理员访问备份管理页面
   - **When** 页面加载
   - **Then** 显示备份存储空间使用情况
   - **And** 显示总存储空间
   - **And** 显示已使用空间
   - **And** 显示使用百分比

5. **AC5: 备份状态统计**
   - **Given** 管理员访问备份管理页面
   - **When** 页面加载
   - **Then** 显示备份成功率
   - **And** 显示最近 7 天备份数量
   - **And** 显示失败备份数量

## Tasks / Subtasks

- [ ] Task 1: 创建备份 API 处理器 (AC: 1)
  - [ ] 1.1 创建 `internal/api/admin/backup_handler.go`
  - [ ] 1.2 实现 `ListBackups` 方法返回备份列表
  - [ ] 1.3 实现分页逻辑
  - [ ] 1.4 实现响应数据格式化

- [ ] Task 2: 实现日期范围筛选 (AC: 2)
  - [ ] 2.1 在备份仓库添加日期范围查询方法
  - [ ] 2.2 解析和验证日期参数
  - [ ] 2.3 实现日期范围查询逻辑

- [ ] Task 3: 实现备份详情 API (AC: 3)
  - [ ] 3.1 实现 `GetBackupDetail` 方法
  - [ ] 3.2 查询备份日志记录
  - [ ] 3.3 返回完整备份信息

- [ ] Task 4: 实现存储空间统计 (AC: 4)
  - [ ] 4.1 创建 `internal/services/backup/storage_stats.go`
  - [ ] 4.2 实现 S3 存储空间查询
  - [ ] 4.3 实现本地存储空间查询（开发环境）
  - [ ] 4.4 实现存储使用情况计算

- [ ] Task 5: 实现备份统计 API (AC: 5)
  - [ ] 5.1 实现 `GetBackupStats` 方法
  - [ ] 5.2 计算备份成功率
  - [ ] 5.3 统计最近 7 天备份数量
  - [ ] 5.4 统计失败备份

- [ ] Task 6: 添加路由和中间件 (AC: 1-5)
  - [ ] 6.1 在路由器添加备份相关路由
  - [ ] 6.2 添加管理员认证中间件
  - [ ] 6.3 添加请求日志中间件

- [ ] Task 7: 编写 API 文档和测试
  - [ ] 7.1 编写 API 文档（Swagger 注释）
  - [ ] 7.2 编写 API 单元测试
  - [ ] 7.3 编写集成测试

## Dev Notes

### 架构模式与约束

**必须遵循的 Clean Architecture 原则 [Source: architecture.md]:**

1. **API 响应格式 [Source: architecture.md#Format Patterns]:**
   - 统一包装器: `{ data: {...}, error: null, meta: {...} }`
   - 分页响应包含 `total`, `page`, `pageSize`

2. **API 命名约定 [Source: architecture.md#Naming Patterns]:**
   - REST 端点: 复数资源名 (例: `/backups`)
   - 查询参数: snake_case (例: `?start_date=2026-03-01`)

3. **认证要求:**
   - 管理员 API 需要管理员认证中间件
   - JWT Token 包含管理员角色标识

### API 设计

**备份列表 API:**

```
GET /v1/admin/backups
Query Parameters:
  - page: int (default: 1)
  - page_size: int (default: 20, max: 100)
  - type: string (config, instance)
  - status: string (pending, running, success, failed)
  - start_date: string (YYYY-MM-DD)
  - end_date: string (YYYY-MM-DD)

Response:
{
  "data": {
    "items": [
      {
        "id": "uuid",
        "type": "config",
        "status": "success",
        "fileSize": 1024,
        "startedAt": "2026-03-05T02:00:00Z",
        "completedAt": "2026-03-05T02:01:00Z",
        "createdAt": "2026-03-05T02:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "pageSize": 20
  },
  "error": null,
  "meta": {}
}
```

**备份详情 API:**

```
GET /v1/admin/backups/:id

Response:
{
  "data": {
    "id": "uuid",
    "type": "config",
    "status": "success",
    "filePath": "backups/config/2026-03-05_020000.backup",
    "fileSize": 1024,
    "startedAt": "2026-03-05T02:00:00Z",
    "completedAt": "2026-03-05T02:01:00Z",
    "logs": [
      {
        "time": "2026-03-05T02:00:00Z",
        "level": "info",
        "message": "Backup started"
      }
    ],
    "createdAt": "2026-03-05T02:00:00Z"
  },
  "error": null,
  "meta": {}
}
```

**备份统计 API:**

```
GET /v1/admin/backups/stats

Response:
{
  "data": {
    "successRate": 99.5,
    "totalBackups": 200,
    "successBackups": 199,
    "failedBackups": 1,
    "last7DaysCount": 14,
    "storageUsage": {
      "total": 107374182400,  // 100GB in bytes
      "used": 1073741824,     // 1GB in bytes
      "usedPercentage": 1.0
    }
  },
  "error": null,
  "meta": {}
}
```

### API 处理器实现

```go
// internal/api/admin/backup_handler.go
package admin

import (
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
)

type BackupHandler struct {
    backupService *backup.Service
    backupRepo    *repository.BackupRepository
}

func NewBackupHandler(backupService *backup.Service, backupRepo *repository.BackupRepository) *BackupHandler {
    return &BackupHandler{
        backupService: backupService,
        backupRepo:    backupRepo,
    }
}

// ListBackups 获取备份列表
// @Summary 获取备份列表
// @Tags admin/backups
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param type query string false "备份类型" Enums(config, instance)
// @Param status query string false "备份状态" Enums(pending, running, success, failed)
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} Response{data=BackupListResponse}
// @Router /v1/admin/backups [get]
func (h *BackupHandler) ListBackups(c *gin.Context) {
    // 解析分页参数
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    if pageSize > 100 {
        pageSize = 100
    }

    // 解析筛选参数
    filter := &repository.BackupFilter{
        Type:   c.Query("type"),
        Status: c.Query("status"),
    }

    // 解析日期范围
    if startDate := c.Query("start_date"); startDate != "" {
        if t, err := time.Parse("2006-01-02", startDate); err == nil {
            filter.StartDate = &t
        }
    }
    if endDate := c.Query("end_date"); endDate != "" {
        if t, err := time.Parse("2006-01-02", endDate); err == nil {
            filter.EndDate = &t
        }
    }

    // 查询备份列表
    items, total, err := h.backupRepo.List(c.Request.Context(), filter, page, pageSize)
    if err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Error: &Error{Message: "Failed to list backups"},
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Data: BackupListResponse{
            Items:    items,
            Total:    total,
            Page:     page,
            PageSize: pageSize,
        },
    })
}

// GetBackupDetail 获取备份详情
// @Summary 获取备份详情
// @Tags admin/backups
// @Security BearerAuth
// @Param id path string true "备份 ID"
// @Success 200 {object} Response{data=BackupDetailResponse}
// @Router /v1/admin/backups/{id} [get]
func (h *BackupHandler) GetBackupDetail(c *gin.Context) {
    id := c.Param("id")

    backup, err := h.backupRepo.GetByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, Response{
            Error: &Error{Message: "Backup not found"},
        })
        return
    }

    // 获取备份日志
    logs, _ := h.backupRepo.GetLogs(c.Request.Context(), id)

    c.JSON(http.StatusOK, Response{
        Data: BackupDetailResponse{
            Backup: backup,
            Logs:   logs,
        },
    })
}

// GetBackupStats 获取备份统计
// @Summary 获取备份统计
// @Tags admin/backups
// @Security BearerAuth
// @Success 200 {object} Response{data=BackupStatsResponse}
// @Router /v1/admin/backups/stats [get]
func (h *BackupHandler) GetBackupStats(c *gin.Context) {
    stats, err := h.backupService.GetStats(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Error: &Error{Message: "Failed to get backup stats"},
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Data: stats,
    })
}
```

### 存储空间统计实现

```go
// internal/services/backup/storage_stats.go
package backup

import (
    "context"
)

type StorageStats struct {
    Total          int64   `json:"total"`
    Used           int64   `json:"used"`
    UsedPercentage float64 `json:"usedPercentage"`
}

type StorageStatsService struct {
    s3Client *s3.Client
    bucket   string
}

func NewStorageStatsService(s3Client *s3.Client, bucket string) *StorageStatsService {
    return &StorageStatsService{
        s3Client: s3Client,
        bucket:   bucket,
    }
}

// GetStorageStats 获取存储空间使用情况
func (s *StorageStatsService) GetStorageStats(ctx context.Context) (*StorageStats, error) {
    // 列出所有备份文件并计算总大小
    var totalSize int64

    input := &s3.ListObjectsV2Input{
        Bucket: aws.String(s.bucket),
    }

    paginator := s3.NewListObjectsV2Paginator(s.s3Client, input)
    for paginator.HasMorePages() {
        page, err := paginator.NextPage(ctx)
        if err != nil {
            return nil, err
        }

        for _, obj := range page.Contents {
            totalSize += *obj.Size
        }
    }

    // 存储配额（从配置获取）
    storageQuota := int64(100 * 1024 * 1024 * 1024) // 100GB

    return &StorageStats{
        Total:          storageQuota,
        Used:           totalSize,
        UsedPercentage: float64(totalSize) / float64(storageQuota) * 100,
    }, nil
}
```

### 备份统计实现

```go
// internal/services/backup/stats.go
package backup

import (
    "context"
    "time"
)

type BackupStats struct {
    SuccessRate    float64       `json:"successRate"`
    TotalBackups   int64         `json:"totalBackups"`
    SuccessBackups int64         `json:"successBackups"`
    FailedBackups  int64         `json:"failedBackups"`
    Last7DaysCount int64         `json:"last7DaysCount"`
    StorageUsage   *StorageStats `json:"storageUsage"`
}

func (s *Service) GetStats(ctx context.Context) (*BackupStats, error) {
    // 计算最近 7 天的日期范围
    now := time.Now()
    sevenDaysAgo := now.AddDate(0, 0, -7)

    // 统计备份数量
    total, _ := s.backupRepo.Count(ctx, nil)
    success, _ := s.backupRepo.Count(ctx, &repository.BackupFilter{Status: "success"})
    failed, _ := s.backupRepo.Count(ctx, &repository.BackupFilter{Status: "failed"})
    last7Days, _ := s.backupRepo.CountByDateRange(ctx, sevenDaysAgo, now)

    // 计算成功率
    var successRate float64
    if total > 0 {
        successRate = float64(success) / float64(total) * 100
    }

    // 获取存储使用情况
    storageStats, _ := s.storageStatsService.GetStorageStats(ctx)

    return &BackupStats{
        SuccessRate:    successRate,
        TotalBackups:   total,
        SuccessBackups: success,
        FailedBackups:  failed,
        Last7DaysCount: last7Days,
        StorageUsage:   storageStats,
    }, nil
}
```

### 路由配置

```go
// 在路由器中添加备份路由
func (r *Router) SetupAdminRoutes(adminHandler *admin.Handler) {
    admin := r.Group("/v1/admin")
    admin.Use(middleware.AdminAuth())

    // 备份管理
    backups := admin.Group("/backups")
    {
        backups.GET("", adminHandler.Backup.ListBackups)
        backups.GET("/stats", adminHandler.Backup.GetBackupStats)
        backups.GET("/:id", adminHandler.Backup.GetBackupDetail)
    }
}
```

### 项目结构规范

**新增文件位置:**

```
backend/
├── internal/
│   ├── api/
│   │   └── admin/
│   │       └── backup_handler.go    # 备份 API 处理器（新增）
│   └── services/
│       └── backup/
│           ├── stats.go             # 备份统计服务（新增）
│           └── storage_stats.go     # 存储统计服务（新增）
```

### 测试标准

**测试要求:**
- 测试框架: Go 原生 testing 包 + testify
- 测试覆盖率目标: >= 70%
- API 测试使用 httptest

**测试用例设计:**

| 测试场景 | 测试文件 | 测试方法 |
|---------|---------|---------|
| 备份列表 | `backup_handler_test.go` | API 测试 |
| 日期筛选 | `backup_handler_test.go` | API 测试 |
| 备份详情 | `backup_handler_test.go` | API 测试 |
| 备份统计 | `stats_test.go` | 单元测试 |
| 存储统计 | `storage_stats_test.go` | 单元测试 |

### 前序依赖

**依赖以下 Story:**
- Story 2.4: 平台管理员独立认证系统（已实现管理员认证）
- Story 6.1: 用户配置自动备份（已创建 backup_records 表）
- Story 6.2: OpenClaw 实例数据备份（已创建 instance_backup_records 表）

### 安全注意事项

1. **认证要求**: 所有 API 需要管理员认证
2. **权限验证**: 验证 JWT Token 包含管理员角色
3. **数据脱敏**: 备份文件路径不暴露敏感信息
4. **访问日志**: 记录所有 API 访问日志

### References

- [Source: architecture.md#API Response Formats] - API 响应格式
- [Source: prd.md#FR26] - 查看系统备份状态需求
- [Source: epics.md#Story 6.3] - 原始故事定义

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
