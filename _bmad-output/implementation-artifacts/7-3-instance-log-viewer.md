# Story 7.3: 实例日志查看功能

Status: ready-for-dev

## Story

As a 平台管理员,
I want 查看实例的部署和运行日志,
so that 排查问题原因。

## Acceptance Criteria

1. **AC1: 部署日志显示**
   - **Given** 管理员在实例详情页
   - **When** 点击"查看日志"
   - **Then** 显示实例的部署日志
   - **And** 日志按时间倒序排列
   - **And** 支持滚动查看历史日志

2. **AC2: 实时日志流**
   - **Given** 管理员正在查看日志
   - **When** 实例正在运行
   - **Then** 通过 WebSocket 实时推送新日志
   - **And** 新日志自动追加到底部
   - **And** 支持暂停/恢复实时推送

3. **AC3: 时间范围筛选**
   - **Given** 管理员正在查看日志
   - **When** 选择时间范围
   - **Then** 只显示指定时间范围内的日志
   - **And** 支持快捷选择（最近 1 小时、24 小时、7 天）
   - **And** 支持自定义时间范围

4. **AC4: 关键词搜索**
   - **Given** 管理员正在查看日志
   - **When** 输入搜索关键词
   - **Then** 高亮匹配的关键词
   - **And** 显示匹配数量
   - **And** 支持跳转到上/下一个匹配项

5. **AC5: 日志下载**
   - **Given** 管理员正在查看日志
   - **When** 点击"下载日志"
   - **Then** 下载当前筛选条件下的日志文件
   - **And** 文件格式为 .txt 或 .log
   - **And** 文件名包含实例 ID 和时间戳

## Tasks / Subtasks

- [ ] Task 1: 创建日志数据模型和 Repository (AC: 1, 3)
  - [ ] 1.1 创建 `internal/domain/log/` 目录
  - [ ] 1.2 定义 LogEntry 结构体（时间、级别、内容）
  - [ ] 1.3 创建 `internal/repository/log_repository.go`
  - [ ] 1.4 实现 GetDeploymentLogs() 方法
  - [ ] 1.5 实现 GetRuntimeLogs() 方法
  - [ ] 1.6 实现时间范围筛选逻辑

- [ ] Task 2: 实现 WebSocket 实时日志推送 (AC: 2)
  - [ ] 2.1 创建 `internal/api/admin/log_websocket.go`
  - [ ] 2.2 实现 WebSocket 连接升级
  - [ ] 2.3 实现日志流订阅机制
  - [ ] 2.4 实现日志实时推送
  - [ ] 2.5 实现连接管理和心跳检测
  - [ ] 2.6 实现断线重连处理

- [ ] Task 3: 创建日志查询和下载 API (AC: 1, 3, 4, 5)
  - [ ] 3.1 创建 `internal/api/admin/log_handler.go`
  - [ ] 3.2 实现 GET /v1/admin/instances/:id/logs 端点
  - [ ] 3.3 实现查询参数（时间范围、关键词搜索）
  - [ ] 3.4 实现 GET /v1/admin/instances/:id/logs/download 端点
  - [ ] 3.5 实现日志文件下载响应

- [ ] Task 4: 集成 Dokploy 日志接口 (AC: 1)
  - [ ] 4.1 在 `internal/infrastructure/dokploy/` 添加日志接口
  - [ ] 4.2 实现 GetApplicationLogs() 方法
  - [ ] 4.3 实现日志流 StreamLogs() 方法
  - [ ] 4.4 处理日志格式转换

- [ ] Task 5: 创建前端日志查看组件 (AC: 1-5)
  - [ ] 5.1 创建 `src/views/admin/InstanceLogs.vue` 页面
  - [ ] 5.2 创建 `src/components/admin/LogViewer.vue` 日志查看器组件
  - [ ] 5.3 实现 WebSocket 连接和日志接收
  - [ ] 5.4 实现时间范围选择器
  - [ ] 5.5 实现关键词搜索和高亮
  - [ ] 5.6 实现日志下载功能
  - [ ] 5.7 实现暂停/恢复实时推送

- [ ] Task 6: 添加单元测试 (AC: 1-5)
  - [ ] 6.1 创建 `log_repository_test.go`
  - [ ] 6.2 创建 `log_handler_test.go`
  - [ ] 6.3 创建 `LogViewer.spec.ts` 组件测试
  - [ ] 6.4 确保测试覆盖率 >= 70%

## Dev Notes

### 架构模式与约束

**必须遵循的架构原则：**
1. **WebSocket 通信**: 使用标准 WebSocket 协议
2. **日志格式**: 统一的日志格式标准
3. **性能考虑**: 大日志文件的流式处理

**关键架构决策 [Source: architecture.md]:**
- 实时通信: WebSocket
- 日志存储: Docker Volume + 可选的外部日志系统
- Dokploy 集成: 通过 Dokploy API 获取容器日志

### 数据模型设计

**LogEntry 结构体:**

```go
type LogEntry struct {
    Timestamp time.Time `json:"timestamp"`
    Level     string    `json:"level"`     // INFO, WARN, ERROR
    Message   string    `json:"message"`
    Source    string    `json:"source"`    // deploy, runtime
}

type LogQueryParams struct {
    InstanceID string     `json:"instanceId"`
    StartTime  *time.Time `json:"startTime"`
    EndTime    *time.Time `json:"endTime"`
    Keyword    string     `json:"keyword"`
    Source     string     `json:"source"`    // deploy, runtime, all
}

type LogListResponse struct {
    Logs       []LogEntry `json:"logs"`
    TotalCount int64      `json:"totalCount"`
    HasMore    bool       `json:"hasMore"`
}
```

### WebSocket 协议设计

**连接端点:**
```
ws://host/v1/admin/instances/:id/logs/stream
```

**消息格式:**

```json
// 服务端推送日志
{
  "type": "log",
  "data": {
    "timestamp": "2026-03-05T10:30:00Z",
    "level": "INFO",
    "message": "Container started successfully"
  }
}

// 客户端发送心跳
{
  "type": "ping"
}

// 服务端响应心跳
{
  "type": "pong"
}

// 客户端暂停推送
{
  "type": "pause"
}

// 客户端恢复推送
{
  "type": "resume"
}
```

### API 端点设计

**GET /v1/admin/instances/:id/logs**

**查询参数:**
| 参数 | 类型 | 描述 |
|------|------|------|
| startTime | string | 开始时间 (ISO 8601) |
| endTime | string | 结束时间 (ISO 8601) |
| keyword | string | 搜索关键词 |
| source | string | 日志来源 (deploy/runtime/all) |
| page | int | 页码 |
| pageSize | int | 每页数量 (默认 100) |

**响应示例:**
```json
{
  "data": {
    "logs": [
      {
        "timestamp": "2026-03-05T10:30:00Z",
        "level": "INFO",
        "message": "Deployment started",
        "source": "deploy"
      }
    ],
    "totalCount": 150,
    "hasMore": true
  },
  "error": null,
  "meta": {}
}
```

**GET /v1/admin/instances/:id/logs/download**

**查询参数:** 同上

**响应:**
- Content-Type: text/plain
- Content-Disposition: attachment; filename="instance-{id}-{timestamp}.log"

### Dokploy 日志集成

**Dokploy API 日志接口:**

```go
// internal/infrastructure/dokploy/log_client.go

type LogClient interface {
    // 获取应用日志
    GetApplicationLogs(appID string, opts LogOptions) ([]LogEntry, error)

    // 流式获取日志
    StreamLogs(ctx context.Context, appID string) (<-chan LogEntry, error)
}

type LogOptions struct {
    TailLines  int       // 最后 N 行
    Since      time.Time // 开始时间
    Follow     bool      // 是否持续跟踪
}
```

### 前端组件设计

**LogViewer 组件结构:**

```vue
<template>
  <div class="log-viewer">
    <!-- 工具栏 -->
    <div class="log-viewer__toolbar">
      <n-date-picker v-model:value="timeRange" type="datetimerange" />
      <n-input v-model:value="keyword" placeholder="搜索关键词">
        <template #prefix>
          <n-icon :component="SearchIcon" />
        </template>
      </n-input>
      <n-button @click="toggleFollow">
        {{ following ? '暂停' : '恢复' }}
      </n-button>
      <n-button @click="downloadLogs">下载</n-button>
    </div>

    <!-- 日志内容 -->
    <div class="log-viewer__content" ref="logContainer">
      <div v-for="log in logs" :key="log.timestamp" class="log-entry">
        <span class="log-time">{{ formatTime(log.timestamp) }}</span>
        <span :class="['log-level', log.level.toLowerCase()]">{{ log.level }}</span>
        <span class="log-message" v-html="highlightKeyword(log.message)"></span>
      </div>
    </div>

    <!-- 状态栏 -->
    <div class="log-viewer__status">
      <span>共 {{ totalCount }} 条日志</span>
      <span v-if="connected" class="connected">实时连接中</span>
      <span v-else class="disconnected">已断开</span>
    </div>
  </div>
</template>
```

**WebSocket 服务:**

```typescript
// src/services/logWebSocket.ts

export class LogWebSocket {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;

  constructor(private instanceId: string) {}

  connect(onMessage: (log: LogEntry) => void) {
    const wsUrl = `ws://${location.host}/v1/admin/instances/${this.instanceId}/logs/stream`;
    this.ws = new WebSocket(wsUrl);

    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.type === 'log') {
        onMessage(data.data);
      }
    };

    this.ws.onclose = () => this.handleDisconnect();
  }

  pause() {
    this.ws?.send(JSON.stringify({ type: 'pause' }));
  }

  resume() {
    this.ws?.send(JSON.stringify({ type: 'resume' }));
  }
}
```

### 性能优化

1. **虚拟滚动**:
   - 使用虚拟滚动处理大量日志
   - 只渲染可见区域的日志条目

2. **日志分页**:
   - 初始加载最近 100 条
   - 滚动到顶部时加载更多

3. **WebSocket 优化**:
   - 心跳检测保持连接
   - 断线自动重连
   - 消息队列缓冲

### 与其他 Story 的依赖关系

**前序依赖:**
- Story 4.1: 实例数据模型 - 需要实例 ID
- Story 4.2: Dokploy API 客户端 - 需要 Dokploy 日志接口
- Story 4.7: 管理员实例管理功能 - 需要实例详情页入口

**后续依赖:**
- Story 7.4: 告警系统 - 告警可关联日志

### 测试标准

**单元测试要求:**
- 日志查询 API 测试
- WebSocket 连接测试
- 时间范围筛选测试
- 关键词搜索测试

**E2E 测试要求:**
- 日志正常显示
- WebSocket 实时更新
- 日志下载功能

### References

- [Source: architecture.md#Real-time Communication] - WebSocket 集成
- [Source: architecture.md#Dokploy Integration] - Dokploy 日志接口
- [Source: epics.md#Story 7.3] - 原始故事定义
- [Source: prd.md#FR32] - 功能需求定义

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
