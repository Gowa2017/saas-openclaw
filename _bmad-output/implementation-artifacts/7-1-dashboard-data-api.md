# Story 7.1: 管理仪表盘数据 API

Status: ready-for-dev

## Story

As a 后端开发者,
I want 实现管理仪表盘数据 API,
so that 前端可以展示平台运营数据。

## Acceptance Criteria

1. **AC1: 用户总数统计**
   - **Given** 系统已运行
   - **When** 调用 GET /v1/admin/dashboard API
   - **Then** 返回用户总数
   - **And** 返回用户总数环比变化（与上周对比）

2. **AC2: 活跃实例数量统计**
   - **Given** 系统已运行
   - **When** 调用 GET /v1/admin/dashboard API
   - **Then** 返回活跃实例数量（状态为 running）
   - **And** 返回实例总数

3. **AC3: 部署成功率统计**
   - **Given** 系统已运行
   - **When** 调用 GET /v1/admin/dashboard API
   - **Then** 返回最近 7 天的部署成功率
   - **And** 成功率 = 成功部署次数 / 总部署次数 * 100%
   - **And** 返回部署总数和成功数

4. **AC4: 系统可用性统计**
   - **Given** 系统已运行
   - **When** 调用 GET /v1/admin/dashboard API
   - **Then** 返回最近 7 天的系统可用性
   - **And** 可用性 = 正常运行时间 / 总时间 * 100%

5. **AC5: 最近告警列表**
   - **Given** 系统已运行
   - **When** 调用 GET /v1/admin/dashboard API
   - **Then** 返回最近 10 条告警记录
   - **And** 每条告警包含：类型、级别、时间、状态

6. **AC6: API 性能要求**
   - **Given** API 已实现
   - **When** 调用 GET /v1/admin/dashboard API
   - **Then** API 响应时间 < 5 秒

## Tasks / Subtasks

- [ ] Task 1: 创建仪表盘数据模型 (AC: 1-5)
  - [ ] 1.1 创建 `internal/domain/dashboard/` 目录
  - [ ] 1.2 定义 DashboardStats 结构体（用户总数、活跃实例、成功率、可用性）
  - [ ] 1.3 定义 AlertSummary 结构体（类型、级别、时间、状态）
  - [ ] 1.4 定义 DashboardResponse 结构体（统一响应格式）

- [ ] Task 2: 创建仪表盘 Repository (AC: 1-5)
  - [ ] 2.1 创建 `internal/repository/dashboard_repository.go`
  - [ ] 2.2 实现 GetUserTotalCount() 方法
  - [ ] 2.3 实现 GetActiveInstanceCount() 方法
  - [ ] 2.4 实现 GetDeploymentSuccessRate() 方法（最近 7 天）
  - [ ] 2.5 实现 GetSystemAvailability() 方法（最近 7 天）
  - [ ] 2.6 实现 GetRecentAlerts() 方法（最近 10 条）

- [ ] Task 3: 创建仪表盘 Service (AC: 1-6)
  - [ ] 3.1 创建 `internal/domain/dashboard/service.go`
  - [ ] 3.2 实现 GetDashboardStats() 方法
  - [ ] 3.3 实现计算逻辑和聚合查询优化
  - [ ] 3.4 添加缓存机制（Redis 可选）

- [ ] Task 4: 创建仪表盘 API Handler (AC: 1-6)
  - [ ] 4.1 创建 `internal/api/admin/dashboard_handler.go`
  - [ ] 4.2 实现 GET /v1/admin/dashboard 端点
  - [ ] 4.3 添加管理员认证中间件
  - [ ] 4.4 实现错误处理和日志记录
  - [ ] 4.5 添加 API 响应时间监控

- [ ] Task 5: 添加单元测试 (AC: 1-6)
  - [ ] 5.1 创建 `dashboard_repository_test.go`
  - [ ] 5.2 创建 `dashboard_service_test.go`
  - [ ] 5.3 创建 `dashboard_handler_test.go`
  - [ ] 5.4 确保测试覆盖率 >= 70%

- [ ] Task 6: API 性能优化 (AC: 6)
  - [ ] 6.1 添加数据库查询索引
  - [ ] 6.2 实现查询结果缓存（可选 Redis）
  - [ ] 6.3 性能测试验证响应时间 < 5 秒

## Dev Notes

### 架构模式与约束

**必须遵循的 Clean Architecture 原则：**
1. **依赖方向**: Handler → Service → Repository
2. **领域层纯净**: `internal/domain/dashboard/` 不依赖外部框架
3. **API 设计**: RESTful 风格，版本化 `/v1/`

**关键架构决策 [Source: architecture.md]:**
- API 响应格式: `{ data: {...}, error: null, meta: {...} }`
- 命名约定: API 端点使用复数资源名
- 认证方式: 管理员 JWT Token 验证

### 数据模型设计

**DashboardStats 结构体:**

```go
type DashboardStats struct {
    UserStats      UserStats      `json:"userStats"`
    InstanceStats  InstanceStats  `json:"instanceStats"`
    DeployStats    DeployStats    `json:"deployStats"`
    SystemStats    SystemStats    `json:"systemStats"`
    RecentAlerts   []AlertSummary `json:"recentAlerts"`
}

type UserStats struct {
    TotalCount     int64   `json:"totalCount"`
    WeekOverWeek   float64 `json:"weekOverWeek"` // 环比变化百分比
}

type InstanceStats struct {
    ActiveCount  int64 `json:"activeCount"`  // running 状态
    TotalCount   int64 `json:"totalCount"`
}

type DeployStats struct {
    SuccessRate    float64 `json:"successRate"`    // 7 天成功率
    TotalCount     int64   `json:"totalCount"`     // 7 天总部署
    SuccessCount   int64   `json:"successCount"`   // 7 天成功数
}

type SystemStats struct {
    Availability   float64 `json:"availability"`   // 7 天可用性
    UptimeSeconds  int64   `json:"uptimeSeconds"`  // 系统运行时间
}

type AlertSummary struct {
    ID        string    `json:"id"`
    Type      string    `json:"type"`      // cpu, memory, disk
    Level     string    `json:"level"`     // warning, critical
    Message   string    `json:"message"`
    Status    string    `json:"status"`    // pending, resolved
    CreatedAt time.Time `json:"createdAt"`
}
```

### API 端点设计

**GET /v1/admin/dashboard**

**请求头:**
```
Authorization: Bearer <admin_jwt_token>
```

**响应示例:**
```json
{
  "data": {
    "userStats": {
      "totalCount": 150,
      "weekOverWeek": 12.5
    },
    "instanceStats": {
      "activeCount": 120,
      "totalCount": 135
    },
    "deployStats": {
      "successRate": 96.5,
      "totalCount": 200,
      "successCount": 193
    },
    "systemStats": {
      "availability": 99.2,
      "uptimeSeconds": 604800
    },
    "recentAlerts": [
      {
        "id": "alert-001",
        "type": "cpu",
        "level": "warning",
        "message": "CPU 使用率超过 80%",
        "status": "pending",
        "createdAt": "2026-03-05T10:30:00Z"
      }
    ]
  },
  "error": null,
  "meta": {
    "timestamp": "2026-03-05T12:00:00Z"
  }
}
```

### 数据库查询优化

**必要的索引:**

```sql
-- 用户表索引
CREATE INDEX idx_tenant_users_created_at ON tenant_users(CreatedAt);

-- 实例表索引
CREATE INDEX idx_openclaw_instances_status ON openclaw_instances(Status);
CREATE INDEX idx_openclaw_instances_created_at ON openclaw_instances(CreatedAt);

-- 部署记录表索引（如存在）
CREATE INDEX idx_deployments_created_at ON deployments(CreatedAt);
CREATE INDEX idx_deployments_status ON deployments(Status);

-- 告警表索引
CREATE INDEX idx_alerts_created_at ON alerts(CreatedAt);
CREATE INDEX idx_alerts_status ON alerts(Status);
```

### 性能优化策略

1. **查询优化**:
   - 使用聚合查询减少数据库往返
   - 添加必要索引
   - 避免 N+1 查询

2. **缓存策略** (可选):
   - 使用 Redis 缓存仪表盘数据
   - 缓存 TTL: 1 分钟
   - 缓存键: `dashboard:stats`

3. **并发查询**:
   - 使用 goroutine 并行查询多个统计数据
   - 使用 sync.WaitGroup 等待所有查询完成

### 与其他 Story 的依赖关系

**前序依赖:**
- Story 2.1: 用户数据模型 - 需要用户表进行统计
- Story 4.1: 实例数据模型 - 需要实例表进行统计
- Story 7.4: 告警系统 - 需要告警数据（可先 mock）

**后续依赖:**
- Story 7.2: 管理仪表盘前端页面 - 需要此 API

### 测试标准

**单元测试要求:**
- Repository 测试: 使用 mock 或测试数据库
- Service 测试: 覆盖所有业务逻辑
- Handler 测试: 覆盖成功和错误场景

**性能测试要求:**
- API 响应时间 < 5 秒
- 并发 100 请求无阻塞

### References

- [Source: architecture.md#API Design] - RESTful API 设计规范
- [Source: architecture.md#Database] - 数据库命名约定
- [Source: epics.md#Story 7.1] - 原始故事定义
- [Source: prd.md#FR28-FR31] - 功能需求定义
- [Source: prd.md#NFR-P4] - API 性能要求

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
