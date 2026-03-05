# Story 4.4: 部署状态追踪与进度反馈

Status: ready-for-dev

## Story

As a 用户,
I want 实时查看部署进度,
so that 知道实例部署进展。

## Acceptance Criteria

1. **AC1: 部署阶段状态返回**
   - **Given** 部署任务已创建
   - **When** 查询部署状态
   - **Then** 返回当前部署阶段
   - **And** 阶段包含：验证配置、创建实例、启动服务
   - **And** 每个阶段有明确的完成状态

2. **AC2: 预计剩余时间显示**
   - **Given** 部署任务进行中
   - **When** 查询部署状态
   - **Then** 返回预计剩余时间
   - **And** 时间基于历史部署数据估算
   - **And** 时间动态更新

3. **AC3: 轮询 API 支持**
   - **Given** 部署任务进行中
   - **When** 前端轮询部署状态 API
   - **Then** 返回最新状态和进度
   - **And** 支持配置轮询间隔（默认 2 秒）
   - **And** 部署完成后停止轮询

4. **AC4: WebSocket 实时更新**
   - **Given** 部署任务进行中
   - **When** 前端建立 WebSocket 连接
   - **Then** 服务端推送状态变更
   - **And** 推送部署日志增量
   - **And** 连接断开后自动重连

5. **AC5: 部署完成状态更新**
   - **Given** 部署流程完成
   - **When** 容器成功启动
   - **Then** 实例状态变为 running
   - **And** 记录部署完成时间
   - **And** 记录总部署耗时

6. **AC6: 前端进度组件**
   - **Given** 用户在部署进度页面
   - **When** 查看部署进度
   - **Then** 显示进度条（0-100%）
   - **And** 显示当前阶段说明
   - **And** 显示预计剩余时间
   - **And** 部署完成后显示成功动画

## Tasks / Subtasks

- [ ] Task 1: 创建部署状态 API (AC: 1, 2, 3)
  - [ ] 1.1 创建 GET /v1/instances/:id/deploy-status 端点
  - [ ] 1.2 定义 DeployStatus 响应结构体
  - [ ] 1.3 实现阶段状态计算逻辑
  - [ ] 1.4 实现预计时间估算逻辑
  - [ ] 1.5 添加轮询控制参数

- [ ] Task 2: 实现 WebSocket 实时推送 (AC: 4)
  - [ ] 2.1 创建 WebSocket 连接端点 /ws/deploy/:id
  - [ ] 2.2 实现 WebSocket 连接管理
  - [ ] 2.3 实现状态变更推送
  - [ ] 2.4 实现日志增量推送
  - [ ] 2.5 实现心跳和重连机制

- [ ] Task 3: 实现部署阶段追踪 (AC: 1, 5)
  - [ ] 3.1 创建 DeployStage 枚举
  - [ ] 3.2 在部署服务中更新阶段状态
  - [ ] 3.3 记录每个阶段的开始和结束时间
  - [ ] 3.4 计算部署总耗时

- [ ] Task 4: 创建前端进度组件 (AC: 6)
  - [ ] 4.1 创建 `frontend/src/components/instances/DeployProgress.vue`
  - [ ] 4.2 实现进度条显示
  - [ ] 4.3 实现阶段说明显示
  - [ ] 4.4 实现预计时间显示
  - [ ] 4.5 实现成功动画（使用 Lottie 或 CSS 动画）

- [ ] Task 5: 创建部署进度页面 (AC: 6)
  - [ ] 5.1 创建 `frontend/src/pages/instances/DeployProgressPage.vue`
  - [ ] 5.2 集成 WebSocket 连接
  - [ ] 5.3 实现轮询降级方案
  - [ ] 5.4 实现部署完成跳转
  - [ ] 5.5 实现错误状态显示

- [ ] Task 6: 编写测试 (AC: 1-6)
  - [ ] 6.1 编写状态 API 测试
  - [ ] 6.2 编写 WebSocket 连接测试
  - [ ] 6.3 编写前端组件测试
  - [ ] 6.4 测试轮询逻辑
  - [ ] 6.5 测试重连机制

## Dev Notes

### 架构模式与约束

**命名约定 [Source: architecture.md]:**
- API 端点: `/v1/instances/:id/deploy-status`
- WebSocket 端点: `/ws/deploy/:id`
- 组件名: `PascalCase` (例: `DeployProgress`)

### API 设计

**部署状态 API:**

```yaml
# GET /v1/instances/:id/deploy-status
# 响应
{
  "data": {
    "instanceId": "inst-xxx",
    "status": "deploying",
    "stage": {
      "current": "starting_service",
      "progress": 66,
      "name": "启动服务"
    },
    "stages": [
      { "name": "验证配置", "status": "completed", "duration": "5s" },
      { "name": "创建实例", "status": "completed", "duration": "60s" },
      { "name": "启动服务", "status": "in_progress", "duration": null }
    ],
    "estimatedTimeRemaining": "30s",
    "logs": [
      { "time": "10:00:00", "message": "开始验证配置" },
      { "time": "10:00:05", "message": "配置验证通过" },
      { "time": "10:00:05", "message": "开始创建容器..." }
    ]
  },
  "error": null,
  "meta": {}
}
```

### 部署阶段定义

**DeployStage 枚举:**

```go
// internal/domain/instance/deploy_stage.go
package instance

// DeployStage 部署阶段
type DeployStage string

const (
    DeployStageValidating  DeployStage = "validating"    // 验证配置
    DeployStageCreating    DeployStage = "creating"      // 创建实例
    DeployStageStarting    DeployStage = "starting"      // 启动服务
    DeployStageCompleted   DeployStage = "completed"     // 已完成
)

// StageInfo 阶段信息
type StageInfo struct {
    Name     string       `json:"name"`
    Status   string       `json:"status"` // pending, in_progress, completed, failed
    Duration string       `json:"duration,omitempty"`
    Progress int          `json:"progress"` // 0-100
}

// DeployStatus 部署状态
type DeployStatus struct {
    InstanceID             string       `json:"instanceId"`
    Status                 InstanceStatus `json:"status"`
    CurrentStage           DeployStage  `json:"currentStage"`
    Stages                 []StageInfo  `json:"stages"`
    EstimatedTimeRemaining string       `json:"estimatedTimeRemaining"`
    Logs                   []LogEntry   `json:"logs"`
}

// LogEntry 日志条目
type LogEntry struct {
    Time    string `json:"time"`
    Message string `json:"message"`
    Level   string `json:"level"` // info, warn, error
}
```

### WebSocket 设计

**WebSocket 消息格式:**

```go
// WebSocket 消息类型
type WSMessageType string

const (
    WSMessageTypeStatus WSMessageType = "status"
    WSMessageTypeLog    WSMessageType = "log"
    WSMessageTypeError  WSMessageType = "error"
    WSMessageTypePong   WSMessageType = "pong"
)

// WSMessage WebSocket 消息
type WSMessage struct {
    Type    WSMessageType `json:"type"`
    Payload interface{}   `json:"payload"`
}
```

**WebSocket Handler:**

```go
// internal/api/instance/websocket.go
package instance

import (
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

// HandleDeployWS 处理部署 WebSocket 连接
func (h *Handler) HandleDeployWS(c *gin.Context) {
    instanceID := c.Param("id")

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    // 订阅部署事件
    sub := h.deployEventBus.Subscribe(instanceID)
    defer h.deployEventBus.Unsubscribe(instanceID, sub)

    // 发送当前状态
    status := h.service.GetDeployStatus(instanceID)
    conn.WriteJSON(WSMessage{Type: WSMessageTypeStatus, Payload: status})

    // 监听事件并推送
    for {
        select {
        case event := <-sub:
            conn.WriteJSON(event)
        case <-c.Done():
            return
        }
    }
}
```

### 前端组件设计

**DeployProgress.vue:**

```vue
<!-- frontend/src/components/instances/DeployProgress.vue -->
<template>
  <div class="deploy-progress">
    <!-- 进度条 -->
    <n-progress
      type="line"
      :percentage="progress"
      :status="progressStatus"
      :height="24"
      :border-radius="4"
    />

    <!-- 当前阶段 -->
    <div class="stage-info">
      <n-icon :component="stageIcon" />
      <span>{{ currentStageName }}</span>
      <span v-if="estimatedTime" class="time-remaining">
        预计剩余 {{ estimatedTime }}
      </span>
    </div>

    <!-- 阶段列表 -->
    <div class="stages-list">
      <div
        v-for="stage in stages"
        :key="stage.name"
        :class="['stage-item', stage.status]"
      >
        <n-icon :component="getStageIcon(stage.status)" />
        <span>{{ stage.name }}</span>
        <span v-if="stage.duration">{{ stage.duration }}</span>
      </div>
    </div>

    <!-- 成功动画 -->
    <div v-if="isCompleted" class="success-animation">
      <lottie-player
        src="/animations/success.json"
        autoplay
      />
      <p>部署成功！</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  CheckCircleOutline,
  TimeOutline,
  SyncOutline,
  CloseCircleOutline
} from '@vicons/ionicons5'

interface Props {
  progress: number
  currentStage: string
  stages: StageInfo[]
  estimatedTime?: string
  status: 'deploying' | 'running' | 'error'
}

const props = defineProps<Props>()

const isCompleted = computed(() => props.status === 'running')
const progressStatus = computed(() => {
  if (props.status === 'error') return 'error'
  if (isCompleted.value) return 'success'
  return 'default'
})

const currentStageName = computed(() => {
  const stage = props.stages.find(s => s.status === 'in_progress')
  return stage?.name || '准备中...'
})
</script>
```

**WebSocket Composable:**

```typescript
// frontend/src/composables/useDeployWebSocket.ts
import { ref, onUnmounted } from 'vue'
import { useMessage } from 'naive-ui'

export function useDeployWebSocket(instanceId: string) {
  const message = useMessage()
  const status = ref<DeployStatus | null>(null)
  const connected = ref(false)

  let ws: WebSocket | null = null
  let reconnectTimer: number | null = null

  function connect() {
    const wsUrl = `${WS_BASE_URL}/ws/deploy/${instanceId}`
    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      connected.value = true
      console.log('WebSocket connected')
    }

    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data)
      if (msg.type === 'status') {
        status.value = msg.payload
      } else if (msg.type === 'log') {
        // 处理增量日志
      }
    }

    ws.onclose = () => {
      connected.value = false
      // 自动重连
      reconnectTimer = setTimeout(connect, 3000)
    }

    ws.onerror = () => {
      message.error('WebSocket 连接失败')
    }
  }

  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
    }
    ws?.close()
  }

  onUnmounted(disconnect)

  return {
    status,
    connected,
    connect,
    disconnect
  }
}
```

### 进度计算逻辑

**阶段进度映射:**

| 阶段 | 进度范围 | 预计时间 |
|-----|---------|---------|
| 验证配置 | 0-15% | 5 秒 |
| 创建实例 | 15-80% | 60-120 秒 |
| 启动服务 | 80-100% | 30 秒 |

```go
// 计算进度百分比
func calculateProgress(stage DeployStage, stageProgress float64) int {
    stageWeights := map[DeployStage]struct{ start, end float64 }{
        DeployStageValidating: {0, 15},
        DeployStageCreating:   {15, 80},
        DeployStageStarting:   {80, 100},
    }

    w := stageWeights[stage]
    return int(w.start + (w.end-w.start)*stageProgress)
}
```

### 项目结构规范

**新增文件位置:**

```
backend/
├── internal/
│   ├── api/
│   │   └── instance/
│   │       ├── handler.go         # 状态 API（更新）
│   │       └── websocket.go       # WebSocket 处理（新增）
│   └── domain/
│       └── instance/
│           ├── deploy_stage.go    # 阶段定义（新增）
│           └── service.go         # 服务层（更新）

frontend/
├── src/
│   ├── components/
│   │   └── instances/
│   │       └── DeployProgress.vue # 进度组件（新增）
│   ├── pages/
│   │   └── instances/
│   │       └── DeployProgressPage.vue # 进度页面（新增）
│   └── composables/
│       └── useDeployWebSocket.ts  # WebSocket composable（新增）
```

### 测试标准

**测试要求:**
- 测试覆盖率目标: >= 70%
- WebSocket 测试使用 gorilla/websocket test server

**关键测试场景:**

| 场景 | 测试方法 |
|-----|---------|
| 状态查询 | TestGetDeployStatus |
| WebSocket 连接 | TestWebSocket_Connect |
| 状态推送 | TestWebSocket_StatusPush |
| 断线重连 | TestWebSocket_Reconnect |
| 进度计算 | TestCalculateProgress |

### References

- [Source: prd.md#FR12] - 查看实例部署状态
- [Source: prd.md#FR13] - 部署成功/失败通知
- [Source: epics.md#Story 4.4] - 原始故事定义

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
