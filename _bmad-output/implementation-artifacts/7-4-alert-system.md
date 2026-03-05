# Story 7.4: 告警系统

Status: ready-for-dev

## Story

As a 系统,
I want 在资源不足时发送告警通知,
so that 管理员可以及时处理问题。

## Acceptance Criteria

1. **AC1: 资源监控检测**
   - **Given** 系统监控服务运行中
   - **When** 检测到资源不足
   - **Then** CPU > 80% 触发告警
   - **And** 内存 > 80% 触发告警
   - **And** 磁盘 > 80% 触发告警
   - **And** 每 60 秒检测一次

2. **AC2: 告警记录创建**
   - **Given** 检测到资源不足
   - **When** 触发告警
   - **Then** 创建告警记录到数据库
   - **And** 告警包含：类型、级别、时间、详情
   - **And** 同类型告警在 5 分钟内不重复创建

3. **AC3: 告警通知发送**
   - **Given** 告警记录已创建
   - **When** 需要通知管理员
   - **Then** 发送告警通知
   - **And** 支持邮件通知
   - **And** 支持钉钉通知
   - **And** 支持企业微信通知

4. **AC4: 告警确认和处理**
   - **Given** 告警已创建
   - **When** 管理员确认告警
   - **Then** 更新告警状态为已确认
   - **And** 记录确认人和确认时间
   - **And** 支持添加处理备注

5. **AC5: 告警状态更新**
   - **Given** 告警已确认
   - **When** 资源恢复正常
   - **Then** 自动更新告警状态为已解决
   - **And** 记录解决时间

## Tasks / Subtasks

- [ ] Task 1: 创建告警数据模型 (AC: 1-5)
  - [ ] 1.1 创建 `internal/domain/alert/` 目录
  - [ ] 1.2 定义 Alert 结构体
  - [ ] 1.3 定义 AlertType 枚举（cpu, memory, disk）
  - [ ] 1.4 定义 AlertLevel 枚举（warning, critical）
  - [ ] 1.5 定义 AlertStatus 枚举（pending, acknowledged, resolved）

- [ ] Task 2: 创建告警 Repository (AC: 2, 4, 5)
  - [ ] 2.1 创建 `internal/repository/alert_repository.go`
  - [ ] 2.2 实现 CreateAlert() 方法
  - [ ] 2.3 实现 GetAlertByID() 方法
  - [ ] 2.4 实现 ListAlerts() 方法
  - [ ] 2.5 实现 UpdateAlertStatus() 方法
  - [ ] 2.6 实现 GetRecentAlertsByType() 方法（防止重复告警）
  - [ ] 2.7 创建数据库表和索引

- [ ] Task 3: 创建资源监控服务 (AC: 1)
  - [ ] 3.1 创建 `internal/domain/monitor/` 目录
  - [ ] 3.2 创建 `internal/domain/monitor/resource_monitor.go`
  - [ ] 3.3 实现 CPU 使用率采集
  - [ ] 3.4 实现内存使用率采集
  - [ ] 3.5 实现磁盘使用率采集
  - [ ] 3.6 实现定时检测逻辑（60 秒间隔）

- [ ] Task 4: 创建告警服务 (AC: 2, 5)
  - [ ] 4.1 创建 `internal/domain/alert/service.go`
  - [ ] 4.2 实现 CheckAndCreateAlert() 方法
  - [ ] 4.3 实现告警去重逻辑（5 分钟内同类型不重复）
  - [ ] 4.4 实现 ResolveAlert() 方法
  - [ ] 4.5 实现告警状态转换逻辑

- [ ] Task 5: 创建告警通知服务 (AC: 3)
  - [ ] 5.1 创建 `internal/domain/alert/notifier.go`
  - [ ] 5.2 定义 Notifier 接口
  - [ ] 5.3 实现 EmailNotifier
  - [ ] 5.4 实现 DingTalkNotifier（钉钉）
  - [ ] 5.5 实现 WeChatWorkNotifier（企业微信）
  - [ ] 5.6 创建通知配置管理

- [ ] Task 6: 创建告警 API Handler (AC: 4)
  - [ ] 6.1 创建 `internal/api/admin/alert_handler.go`
  - [ ] 6.2 实现 GET /v1/admin/alerts 端点（列表）
  - [ ] 6.3 实现 GET /v1/admin/alerts/:id 端点（详情）
  - [ ] 6.4 实现 PUT /v1/admin/alerts/:id/acknowledge 端点（确认）
  - [ ] 6.5 实现 PUT /v1/admin/alerts/:id/resolve 端点（解决）
  - [ ] 6.6 实现通知配置 API

- [ ] Task 7: 创建前端告警管理页面 (AC: 4, 5)
  - [ ] 7.1 创建 `src/views/admin/AlertManagement.vue` 页面
  - [ ] 7.2 创建 `src/components/admin/AlertList.vue` 列表组件
  - [ ] 7.3 创建 `src/components/admin/AlertDetail.vue` 详情组件
  - [ ] 7.4 创建 `src/components/admin/NotificationConfig.vue` 配置组件
  - [ ] 7.5 实现告警确认和备注功能

- [ ] Task 8: 添加单元测试 (AC: 1-5)
  - [ ] 8.1 创建 `alert_repository_test.go`
  - [ ] 8.2 创建 `alert_service_test.go`
  - [ ] 8.3 创建 `resource_monitor_test.go`
  - [ ] 8.4 创建 `notifier_test.go`
  - [ ] 8.5 创建前端组件测试
  - [ ] 8.6 确保测试覆盖率 >= 70%

## Dev Notes

### 架构模式与约束

**必须遵循的架构原则：**
1. **定时任务**: 使用 Go ticker 实现定时检测
2. **策略模式**: 通知渠道使用策略模式
3. **去重机制**: 防止告警风暴

**关键架构决策:**
- 监控间隔: 60 秒
- 告警阈值: CPU/内存/磁盘 > 80%
- 去重窗口: 5 分钟内同类型不重复

### 数据模型设计

**Alert 结构体:**

```go
type Alert struct {
    ID            string       `json:"id" db:"ID"`
    Type          AlertType    `json:"type" db:"Type"`
    Level         AlertLevel   `json:"level" db:"Level"`
    Status        AlertStatus  `json:"status" db:"Status"`
    Message       string       `json:"message" db:"Message"`
    Value         float64      `json:"value" db:"Value"`         // 触发值
    Threshold     float64      `json:"threshold" db:"Threshold"` // 阈值
    AcknowledgedBy *string     `json:"acknowledgedBy,omitempty" db:"AcknowledgedBy"`
    AcknowledgedAt *time.Time  `json:"acknowledgedAt,omitempty" db:"AcknowledgedAt"`
    ResolvedAt    *time.Time   `json:"resolvedAt,omitempty" db:"ResolvedAt"`
    Notes         string       `json:"notes" db:"Notes"`
    CreatedAt     time.Time    `json:"createdAt" db:"CreatedAt"`
    UpdatedAt     time.Time    `json:"updatedAt" db:"UpdatedAt"`
}

type AlertType string

const (
    AlertTypeCPU    AlertType = "cpu"
    AlertTypeMemory AlertType = "memory"
    AlertTypeDisk   AlertType = "disk"
)

type AlertLevel string

const (
    AlertLevelWarning  AlertLevel = "warning"  // 80-90%
    AlertLevelCritical AlertLevel = "critical" // > 90%
)

type AlertStatus string

const (
    AlertStatusPending      AlertStatus = "pending"
    AlertStatusAcknowledged AlertStatus = "acknowledged"
    AlertStatusResolved     AlertStatus = "resolved"
)
```

### 数据库表设计

```sql
CREATE TABLE alerts (
    ID VARCHAR(36) PRIMARY KEY,
    Type VARCHAR(20) NOT NULL,
    Level VARCHAR(20) NOT NULL,
    Status VARCHAR(20) NOT NULL DEFAULT 'pending',
    Message TEXT NOT NULL,
    Value DECIMAL(5,2) NOT NULL,
    Threshold DECIMAL(5,2) NOT NULL,
    AcknowledgedBy VARCHAR(100),
    AcknowledgedAt TIMESTAMP,
    ResolvedAt TIMESTAMP,
    Notes TEXT,
    CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX idx_alerts_status ON alerts(Status);
CREATE INDEX idx_alerts_type ON alerts(Type);
CREATE INDEX idx_alerts_created_at ON alerts(CreatedAt);
CREATE INDEX idx_alerts_type_created ON alerts(Type, CreatedAt);
```

### 资源监控实现

**系统资源采集:**

```go
// internal/domain/monitor/resource_monitor.go

type ResourceMonitor struct {
    interval time.Duration
}

type ResourceUsage struct {
    CPUUsage    float64
    MemoryUsage float64
    DiskUsage   float64
}

func (m *ResourceMonitor) Collect() (*ResourceUsage, error) {
    cpu, err := m.getCPUUsage()
    if err != nil {
        return nil, err
    }

    memory, err := m.getMemoryUsage()
    if err != nil {
        return nil, err
    }

    disk, err := m.getDiskUsage()
    if err != nil {
        return nil, err
    }

    return &ResourceUsage{
        CPUUsage:    cpu,
        MemoryUsage: memory,
        DiskUsage:   disk,
    }, nil
}

// 使用 gopsutil 库采集系统指标
func (m *ResourceMonitor) getCPUUsage() (float64, error) {
    percent, err := cpu.Percent(time.Second, false)
    if err != nil {
        return 0, err
    }
    return percent[0], nil
}
```

### 通知服务接口

**Notifier 接口:**

```go
// internal/domain/alert/notifier.go

type Notifier interface {
    Send(alert *Alert) error
    GetType() string
}

type NotificationConfig struct {
    Email      *EmailConfig      `json:"email"`
    DingTalk   *DingTalkConfig   `json:"dingTalk"`
    WeChatWork *WeChatWorkConfig `json:"weChatWork"`
}

type EmailConfig struct {
    Enabled  bool     `json:"enabled"`
    SMTP     string   `json:"smtp"`
    Port     int      `json:"port"`
    Username string   `json:"username"`
    Password string   `json:"password"`
    To       []string `json:"to"`
}

type DingTalkConfig struct {
    Enabled  bool   `json:"enabled"`
    Webhook  string `json:"webhook"`
    Secret   string `json:"secret"`
}

type WeChatWorkConfig struct {
    Enabled bool   `json:"enabled"`
    Webhook string `json:"webhook"`
}
```

**钉钉通知实现:**

```go
func (n *DingTalkNotifier) Send(alert *Alert) error {
    message := map[string]interface{}{
        "msgtype": "markdown",
        "markdown": map[string]string{
            "title": fmt.Sprintf("[%s] 系统告警", strings.ToUpper(string(alert.Level))),
            "text": fmt.Sprintf(
                "### 系统告警通知\n\n"+
                    "- **类型**: %s\n"+
                    "- **级别**: %s\n"+
                    "- **当前值**: %.2f%%\n"+
                    "- **阈值**: %.2f%%\n"+
                    "- **时间**: %s\n"+
                    "- **详情**: %s\n",
                alert.Type,
                alert.Level,
                alert.Value,
                alert.Threshold,
                alert.CreatedAt.Format("2006-01-02 15:04:05"),
                alert.Message,
            ),
        },
    }

    // 发送 HTTP 请求到钉钉 Webhook
    // ...
}
```

### API 端点设计

**告警列表:**

```
GET /v1/admin/alerts?status=pending&type=cpu&page=1&pageSize=20
```

**响应:**
```json
{
  "data": {
    "alerts": [...],
    "total": 50,
    "page": 1,
    "pageSize": 20
  }
}
```

**确认告警:**

```
PUT /v1/admin/alerts/:id/acknowledge
```

**请求体:**
```json
{
  "notes": "已检查，是由于定时任务导致的 CPU 峰值"
}
```

### 前端组件设计

**告警列表页面:**

```vue
<template>
  <div class="alert-management">
    <!-- 筛选栏 -->
    <n-space>
      <n-select v-model:value="filter.status" :options="statusOptions" />
      <n-select v-model:value="filter.type" :options="typeOptions" />
    </n-space>

    <!-- 告警列表 -->
    <n-data-table :columns="columns" :data="alerts" />

    <!-- 详情弹窗 -->
    <n-modal v-model:show="showDetail">
      <AlertDetail :alert="selectedAlert" @acknowledge="handleAcknowledge" />
    </n-modal>
  </div>
</template>
```

### 告警流程

```
┌─────────────────┐
│  资源监控服务   │
│  (60s 定时检测) │
└────────┬────────┘
         │ 检测到资源 > 80%
         ▼
┌─────────────────┐
│  告警去重检查   │
│  (5 分钟窗口)   │
└────────┬────────┘
         │ 未在窗口内
         ▼
┌─────────────────┐
│  创建告警记录   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  发送通知       │
│  (邮件/钉钉/微信)│
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  管理员确认     │
└────────┬────────┘
         │ 资源恢复
         ▼
┌─────────────────┐
│  自动解决       │
└─────────────────┘
```

### 与其他 Story 的依赖关系

**前序依赖:**
- Story 2.4: 平台管理员独立认证系统 - 需要管理员信息
- Story 4.2: Dokploy API 客户端 - 可复用 HTTP 客户端

**后续依赖:**
- Story 7.1: 管理仪表盘数据 API - 需要告警数据
- Story 7.2: 管理仪表盘前端页面 - 显示告警列表

### 测试标准

**单元测试要求:**
- 资源监控采集测试
- 告警去重逻辑测试
- 通知发送测试（mock 外部服务）

**集成测试要求:**
- 告警完整流程测试
- WebSocket 实时更新测试

### References

- [Source: architecture.md#Monitoring] - 监控系统设计
- [Source: epics.md#Story 7.4] - 原始故事定义
- [Source: prd.md#FR33] - 功能需求定义

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
