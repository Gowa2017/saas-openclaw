# Story 4.2: Dokploy API 客户端集成

Status: ready-for-dev

## Story

As a 后端开发者,
I want 集成 Dokploy API 客户端,
so that 可以通过 API 管理 OpenClaw 容器实例。

## Acceptance Criteria

1. **AC1: Dokploy API 客户端基础结构**
   - **Given** Dokploy 服务已部署
   - **When** 实现 Dokploy API 客户端
   - **Then** 创建 `internal/infrastructure/dokploy/client.go` 客户端结构
   - **And** 支持配置 API Base URL 和 API Token
   - **And** 使用 HTTP Client 进行 API 调用

2. **AC2: CreateApplication 接口实现**
   - **Given** Dokploy 客户端已初始化
   - **When** 调用 CreateApplication 接口
   - **Then** 成功创建容器应用
   - **And** 返回应用 ID 和容器 ID
   - **And** 支持配置镜像、环境变量、网络、存储卷

3. **AC3: GetApplication 接口实现**
   - **Given** 应用已创建
   - **When** 调用 GetApplication 接口
   - **Then** 返回应用状态信息
   - **And** 包含运行状态、健康状态、资源使用情况

4. **AC4: 生命周期管理接口实现**
   - **Given** 应用已创建
   - **When** 调用 StartApplication/StopApplication/RestartApplication 接口
   - **Then** 成功执行对应操作
   - **And** 操作结果实时返回

5. **AC5: GetLogs 接口实现**
   - **Given** 应用已创建
   - **When** 调用 GetLogs 接口
   - **Then** 返回容器日志
   - **And** 支持日志行数限制
   - **And** 支持实时日志流

6. **AC6: 超时与重试机制**
   - **Given** API 调用可能失败
   - **When** 配置超时和重试策略
   - **Then** API 调用超时设置为 5 秒
   - **And** 实现重试机制（最多 3 次）
   - **And** 重试间隔采用指数退避策略

## Tasks / Subtasks

- [ ] Task 1: 创建 Dokploy 客户端基础结构 (AC: 1)
  - [ ] 1.1 创建 `internal/infrastructure/dokploy/client.go`
  - [ ] 1.2 定义 Client 结构体和配置
  - [ ] 1.3 实现 NewClient 构造函数
  - [ ] 1.4 实现 HTTP 请求封装方法
  - [ ] 1.5 添加配置项到 `.env.example`（DOKPLOY_API_URL, DOKPLOY_API_TOKEN）

- [ ] Task 2: 实现应用创建接口 (AC: 2)
  - [ ] 2.1 创建 `internal/infrastructure/dokploy/types.go` 定义请求/响应结构
  - [ ] 2.2 实现 CreateApplication 方法
  - [ ] 2.3 实现镜像配置支持
  - [ ] 2.4 实现环境变量配置支持
  - [ ] 2.5 实现网络配置支持
  - [ ] 2.6 实现存储卷配置支持

- [ ] Task 3: 实现应用查询接口 (AC: 3)
  - [ ] 3.1 实现 GetApplication 方法
  - [ ] 3.2 定义 ApplicationStatus 结构体
  - [ ] 3.3 实现状态解析逻辑

- [ ] Task 4: 实现生命周期管理接口 (AC: 4)
  - [ ] 4.1 实现 StartApplication 方法
  - [ ] 4.2 实现 StopApplication 方法
  - [ ] 4.3 实现 RestartApplication 方法
  - [ ] 4.4 实现操作结果验证

- [ ] Task 5: 实现日志接口 (AC: 5)
  - [ ] 5.1 实现 GetLogs 方法
  - [ ] 5.2 支持日志行数限制参数
  - [ ] 5.3 实现实时日志流（可选 WebSocket）

- [ ] Task 6: 实现超时与重试机制 (AC: 6)
  - [ ] 6.1 配置 HTTP Client 超时为 5 秒
  - [ ] 6.2 实现重试中间件
  - [ ] 6.3 实现指数退避策略
  - [ ] 6.4 记录重试日志

- [ ] Task 7: 编写单元测试 (AC: 1-6)
  - [ ] 7.1 创建 `client_test.go` 测试文件
  - [ ] 7.2 使用 mock 服务器测试 API 调用
  - [ ] 7.3 测试超时和重试逻辑
  - [ ] 7.4 测试错误处理

## Dev Notes

### 架构模式与约束

**必须遵循的 Clean Architecture 原则 [Source: architecture.md]:**

1. **依赖方向**:
   - `infrastructure/dokploy` 层封装外部 API 调用
   - 通过接口定义依赖，便于测试和替换

2. **命名约定 [Source: architecture.md#Naming Patterns]:**
   - 文件名: `kebab-case` (例: `dokploy-client.go`)
   - 结构体: `PascalCase` (例: `DokployClient`)
   - 方法名: `PascalCase` (导出方法)

3. **错误处理**:
   - 定义明确的错误类型
   - 包装错误信息便于调试

### Dokploy API 参考

**API 端点 [Source:技术研究报告]:**

| 操作 | 端点 | 方法 |
|-----|------|------|
| 创建应用 | `/api/application` | POST |
| 获取应用 | `/api/application/{id}` | GET |
| 启动应用 | `/api/application/{id}/start` | POST |
| 停止应用 | `/api/application/{id}/stop` | POST |
| 重启应用 | `/api/application/{id}/restart` | POST |
| 获取日志 | `/api/application/{id}/logs` | GET |

**请求示例:**

```json
// POST /api/application
{
  "name": "openclaw-{tenant_id}",
  "image": "openclaw/openclaw:latest",
  "env": [
    "FEISHU_APP_ID=cli_xxx",
    "FEISHU_APP_SECRET=xxx"
  ],
  "networks": ["openclaw-network"],
  "volumes": [
    {
      "source": "openclaw-data-{tenant_id}",
      "target": "/app/data"
    }
  ]
}
```

### 客户端设计

**Client 结构体:**

```go
// internal/infrastructure/dokploy/client.go
package dokploy

import (
    "context"
    "net/http"
    "time"
)

// Config Dokploy 客户端配置
type Config struct {
    BaseURL   string
    APIToken  string
    Timeout   time.Duration // 默认 5 秒
    MaxRetry  int           // 默认 3 次
}

// Client Dokploy API 客户端
type Client struct {
    config     Config
    httpClient *http.Client
}

// NewClient 创建新的 Dokploy 客户端
func NewClient(cfg Config) *Client {
    if cfg.Timeout == 0 {
        cfg.Timeout = 5 * time.Second
    }
    if cfg.MaxRetry == 0 {
        cfg.MaxRetry = 3
    }

    return &Client{
        config: cfg,
        httpClient: &http.Client{
            Timeout: cfg.Timeout,
        },
    }
}
```

**请求/响应类型:**

```go
// internal/infrastructure/dokploy/types.go
package dokploy

// CreateApplicationRequest 创建应用请求
type CreateApplicationRequest struct {
    Name       string            `json:"name"`
    Image      string            `json:"image"`
    Env        []string          `json:"env,omitempty"`
    Networks   []string          `json:"networks,omitempty"`
    Volumes    []VolumeConfig    `json:"volumes,omitempty"`
    Resources  *ResourceConfig   `json:"resources,omitempty"`
}

// VolumeConfig 存储卷配置
type VolumeConfig struct {
    Source string `json:"source"`
    Target string `json:"target"`
}

// ResourceConfig 资源配置
type ResourceConfig struct {
    CPU    string `json:"cpu,omitempty"`
    Memory string `json:"memory,omitempty"`
}

// Application 应用信息
type Application struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Status      string `json:"status"`
    ContainerID string `json:"containerId"`
    Image       string `json:"image"`
    CreatedAt   string `json:"createdAt"`
}

// ApplicationStatus 应用状态
type ApplicationStatus struct {
    Status    string `json:"status"` // running, stopped, error
    Health    string `json:"health"` // healthy, unhealthy
    CPUUsage  string `json:"cpuUsage"`
    MemUsage  string `json:"memUsage"`
}

// LogsResponse 日志响应
type LogsResponse struct {
    Logs []string `json:"logs"`
}
```

### 重试机制实现

**指数退避策略:**

```go
// internal/infrastructure/dokploy/retry.go
package dokploy

import (
    "context"
    "math"
    "time"
)

// RetryConfig 重试配置
type RetryConfig struct {
    MaxRetry  int
    BaseDelay time.Duration
    MaxDelay  time.Duration
}

// DefaultRetryConfig 默认重试配置
var DefaultRetryConfig = RetryConfig{
    MaxRetry:  3,
    BaseDelay: 100 * time.Millisecond,
    MaxDelay:  5 * time.Second,
}

// calculateDelay 计算重试延迟（指数退避）
func calculateDelay(attempt int, cfg RetryConfig) time.Duration {
    delay := time.Duration(math.Pow(2, float64(attempt))) * cfg.BaseDelay
    if delay > cfg.MaxDelay {
        delay = cfg.MaxDelay
    }
    return delay
}
```

### 项目结构规范

**新增文件位置:**

```
backend/
├── internal/
│   └── infrastructure/
│       └── dokploy/
│           ├── client.go        # 客户端主文件（新增）
│           ├── types.go         # 类型定义（新增）
│           ├── retry.go         # 重试逻辑（新增）
│           └── client_test.go   # 客户端测试（新增）
```

### 性能要求

**API 调用性能 [Source: prd.md#NFR-I3]:**
- Dokploy API 调用响应时间 < 5 秒
- 连接超时: 3 秒
- 读写超时: 5 秒

### 错误处理

**错误类型定义:**

```go
// internal/infrastructure/dokploy/errors.go
package dokploy

import "errors"

var (
    ErrTimeout      = errors.New("dokploy: request timeout")
    ErrUnauthorized = errors.New("dokploy: unauthorized")
    ErrNotFound     = errors.New("dokploy: application not found")
    ErrMaxRetry     = errors.New("dokploy: max retry exceeded")
)

// APIError Dokploy API 错误
type APIError struct {
    Code    int
    Message string
}
```

### 测试标准

**测试要求:**
- 测试框架: Go 原生 testing 包 + testify
- 测试覆盖率目标: >= 70%
- 使用 httptest 模拟 Dokploy API 服务器

**测试用例设计:**

| 测试场景 | 测试方法 |
|---------|---------|
| 创建应用成功 | TestCreateApplication_Success |
| 创建应用失败 | TestCreateApplication_Error |
| 获取应用状态 | TestGetApplication_Status |
| 启动/停止应用 | TestLifecycle_Operations |
| 超时处理 | TestTimeout_Handling |
| 重试机制 | TestRetry_Logic |

### 配置项

**环境变量:**

```bash
# .env.example
DOKPLOY_API_URL=http://dokploy.local:3000
DOKPLOY_API_TOKEN=your-api-token
DOKPLOY_TIMEOUT=5s
DOKPLOY_MAX_RETRY=3
```

### References

- [Source: architecture.md#Project Structure] - 项目目录结构
- [Source: prd.md#FR10-FR18] - OpenClaw 实例管理需求
- [Source: prd.md#NFR-I3] - Dokploy API 响应时间要求
- [Source: epics.md#Story 4.2] - 原始故事定义
- [Source: technical-openclaw-container-platform-research] - Dokploy API 文档

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
