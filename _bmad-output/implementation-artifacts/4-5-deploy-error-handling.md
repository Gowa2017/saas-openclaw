# Story 4.5: 部署失败处理与错误详情

Status: ready-for-dev

## Story

As a 用户,
I want 在部署失败时看到详细的错误信息,
so that 可以了解原因并尝试修复。

## Acceptance Criteria

1. **AC1: 错误状态更新**
   - **Given** 部署过程中发生错误
   - **When** 部署失败
   - **Then** 状态更新为 error
   - **And** 记录错误发生时间
   - **And** 记录错误发生的阶段

2. **AC2: 详细错误日志记录**
   - **Given** 部署失败
   - **When** 记录错误信息
   - **Then** 记录完整错误堆栈
   - **And** 记录 Dokploy API 返回的错误信息
   - **And** 记录容器日志（如有）
   - **And** 错误日志持久化到数据库

3. **AC3: 用户友好的错误提示**
   - **Given** 部署失败
   - **When** 用户查看错误信息
   - **Then** 显示用户友好的错误提示
   - **And** 错误提示使用中文描述
   - **And** 显示错误代码便于排查

4. **AC4: 解决方案建议**
   - **Given** 部署失败
   - **When** 显示错误信息
   - **Then** 提供可能的解决方案
   - **And** 解决方案基于错误类型智能推荐
   - **And** 包含相关文档链接

5. **AC5: 重新部署功能**
   - **Given** 部署失败
   - **When** 用户点击"重新部署"
   - **Then** 清理失败的资源
   - **And** 重新执行部署流程
   - **And** 保留原有配置

6. **AC6: 完整错误日志查看**
   - **Given** 部署失败
   - **When** 用户点击"查看详情"
   - **Then** 显示完整错误日志
   - **And** 支持日志下载
   - **And** 敏感信息脱敏显示

## Tasks / Subtasks

- [ ] Task 1: 定义错误类型与代码 (AC: 1, 3, 4)
  - [ ] 1.1 创建 `internal/domain/instance/errors.go`
  - [ ] 1.2 定义错误代码枚举
  - [ ] 1.3 定义错误消息映射
  - [ ] 1.4 定义解决方案映射

- [ ] Task 2: 实现错误处理服务 (AC: 1, 2)
  - [ ] 2.1 创建 `internal/domain/instance/error_handler.go`
  - [ ] 2.2 实现 HandleDeployError 方法
  - [ ] 2.3 实现错误日志记录
  - [ ] 2.4 实现错误分类逻辑

- [ ] Task 3: 创建错误详情 API (AC: 3, 6)
  - [ ] 3.1 创建 GET /v1/instances/:id/error 端点
  - [ ] 3.2 定义 ErrorResponse 结构体
  - [ ] 3.3 实现敏感信息脱敏
  - [ ] 3.4 实现日志下载功能

- [ ] Task 4: 实现重新部署功能 (AC: 5)
  - [ ] 4.1 创建 POST /v1/instances/:id/redeploy 端点
  - [ ] 4.2 实现资源清理逻辑
  - [ ] 4.3 重置实例状态
  - [ ] 4.4 重新触发部署流程

- [ ] Task 5: 创建前端错误展示组件 (AC: 3, 4, 6)
  - [ ] 5.1 创建 `frontend/src/components/instances/DeployError.vue`
  - [ ] 5.2 实现错误信息显示
  - [ ] 5.3 实现解决方案展示
  - [ ] 5.4 实现日志查看/下载
  - [ ] 5.5 实现"重新部署"按钮

- [ ] Task 6: 编写测试 (AC: 1-6)
  - [ ] 6.1 测试错误分类
  - [ ] 6.2 测试错误记录
  - [ ] 6.3 测试敏感信息脱敏
  - [ ] 6.4 测试重新部署流程

## Dev Notes

### 架构模式与约束

**错误处理原则:**
1. 错误信息对开发者详细，对用户友好
2. 敏感信息不在客户端暴露
3. 所有错误有唯一代码便于追踪

### 错误类型定义

**错误代码枚举:**

```go
// internal/domain/instance/errors.go
package instance

import "errors"

// ErrorCode 错误代码
type ErrorCode string

const (
    // 配置相关错误 (E1xxx)
    ErrCodeFeishuConfigNotFound  ErrorCode = "E1001"
    ErrCodeFeishuConfigInvalid   ErrorCode = "E1002"
    ErrCodeFeishuAppIDInvalid    ErrorCode = "E1003"
    ErrCodeFeishuSecretInvalid   ErrorCode = "E1004"

    // Dokploy 相关错误 (E2xxx)
    ErrCodeDokployUnavailable    ErrorCode = "E2001"
    ErrCodeDokployTimeout        ErrorCode = "E2002"
    ErrCodeDokployAuthFailed     ErrorCode = "E2003"
    ErrCodeDokployCreateFailed   ErrorCode = "E2004"
    ErrCodeDokployStartFailed    ErrorCode = "E2005"

    // 资源相关错误 (E3xxx)
    ErrCodeInsufficientCPU       ErrorCode = "E3001"
    ErrCodeInsufficientMemory    ErrorCode = "E3002"
    ErrCodeInsufficientStorage   ErrorCode = "E3003"

    // 容器相关错误 (E4xxx)
    ErrCodeContainerStartFailed  ErrorCode = "E4001"
    ErrCodeContainerCrash        ErrorCode = "E4002"
    ErrCodeContainerHealthCheck  ErrorCode = "E4003"

    // 未知错误
    ErrCodeUnknown               ErrorCode = "E9999"
)

// ErrorInfo 错误信息
type ErrorInfo struct {
    Code        ErrorCode `json:"code"`
    Message     string    `json:"message"`      // 用户友好消息
    Detail      string    `json:"detail"`       // 详细错误（仅后端日志）
    Suggestion  string    `json:"suggestion"`   // 解决方案建议
    DocLink     string    `json:"docLink"`      // 文档链接
    Stage       string    `json:"stage"`        // 发生阶段
    Timestamp   string    `json:"timestamp"`
}
```

**错误消息映射:**

```go
// 错误配置
var errorMessages = map[ErrorCode]struct {
    Message    string
    Suggestion string
    DocLink    string
}{
    ErrCodeFeishuConfigNotFound: {
        Message:    "未找到飞书应用配置",
        Suggestion: "请先完成飞书应用配置，填写 App ID 和 App Secret",
        DocLink:    "/docs/feishu-config",
    },
    ErrCodeFeishuConfigInvalid: {
        Message:    "飞书应用配置无效",
        Suggestion: "请检查 App ID 和 App Secret 是否正确，确保应用状态正常",
        DocLink:    "/docs/feishu-config#validation",
    },
    ErrCodeDokployUnavailable: {
        Message:    "容器服务暂时不可用",
        Suggestion: "请稍后重试，如问题持续请联系技术支持",
        DocLink:    "/docs/troubleshooting#dokploy",
    },
    ErrCodeDokployTimeout: {
        Message:    "部署超时",
        Suggestion: "部署时间过长，请检查网络连接或联系技术支持",
        DocLink:    "/docs/troubleshooting#timeout",
    },
    ErrCodeInsufficientCPU: {
        Message:    "CPU 资源不足",
        Suggestion: "当前服务器资源紧张，请稍后重试或联系管理员",
        DocLink:    "/docs/troubleshooting#resources",
    },
    ErrCodeContainerStartFailed: {
        Message:    "容器启动失败",
        Suggestion: "请检查飞书配置是否正确，或查看详细日志了解原因",
        DocLink:    "/docs/troubleshooting#container",
    },
}
```

### 错误处理服务

**ErrorHandler:**

```go
// internal/domain/instance/error_handler.go
package instance

import (
    "context"
    "time"
    "github.com/gowa/saas-openclaw/backend/pkg/logger"
)

// ErrorHandler 错误处理器
type ErrorHandler struct {
    repo   InstanceRepository
    logger *logger.Logger
}

// HandleDeployError 处理部署错误
func (h *ErrorHandler) HandleDeployError(ctx context.Context, instanceID string, stage DeployStage, err error) *ErrorInfo {
    // 分类错误
    errorCode := h.classifyError(err)

    // 获取错误配置
    cfg, ok := errorMessages[errorCode]
    if !ok {
        cfg = errorMessages[ErrCodeUnknown]
    }

    // 构建错误信息
    errorInfo := &ErrorInfo{
        Code:       errorCode,
        Message:    cfg.Message,
        Detail:     err.Error(), // 仅记录到日志
        Suggestion: cfg.Suggestion,
        DocLink:    cfg.DocLink,
        Stage:      string(stage),
        Timestamp:  time.Now().Format(time.RFC3339),
    }

    // 记录详细日志
    h.logger.Error("deploy error",
        "instanceID", instanceID,
        "code", errorCode,
        "stage", stage,
        "detail", err.Error(),
    )

    // 更新实例状态
    h.repo.UpdateStatus(instanceID, InstanceStatusError, errorInfo.Message)

    // 更新部署日志
    h.appendDeployLog(instanceID, errorInfo)

    return errorInfo
}

// classifyError 分类错误
func (h *ErrorHandler) classifyError(err error) ErrorCode {
    // 根据错误类型返回对应代码
    switch {
    case errors.Is(err, ErrFeishuConfigNotFound):
        return ErrCodeFeishuConfigNotFound
    case errors.Is(err, ErrDokployUnavailable):
        return ErrCodeDokployUnavailable
    case errors.Is(err, context.DeadlineExceeded):
        return ErrCodeDokployTimeout
    // ... 更多错误映射
    default:
        return ErrCodeUnknown
    }
}
```

### API 设计

**错误详情 API:**

```yaml
# GET /v1/instances/:id/error
# 响应
{
  "data": {
    "code": "E1002",
    "message": "飞书应用配置无效",
    "suggestion": "请检查 App ID 和 App Secret 是否正确，确保应用状态正常",
    "docLink": "/docs/feishu-config#validation",
    "stage": "validating",
    "timestamp": "2026-03-05T10:00:30Z",
    "logs": [
      {
        "time": "10:00:28",
        "level": "info",
        "message": "开始验证飞书配置..."
      },
      {
        "time": "10:00:30",
        "level": "error",
        "message": "飞书 API 返回: invalid app_id"
      }
    ]
  },
  "error": null,
  "meta": {}
}
```

**重新部署 API:**

```yaml
# POST /v1/instances/:id/redeploy
# 响应
{
  "data": {
    "id": "inst-xxx",
    "status": "pending",
    "message": "重新部署任务已创建"
  },
  "error": null,
  "meta": {}
}
```

### 敏感信息脱敏

**脱敏规则:**

```go
// internal/domain/instance/sanitizer.go
package instance

import "regexp"

// SanitizeLog 脱敏日志
func SanitizeLog(log string) string {
    // App Secret 脱敏
    secretPattern := regexp.MustCompile(`(FEISHU_APP_SECRET=)[^\s]+`)
    log = secretPattern.ReplaceAllString(log, `$1****`)

    // Token 脱敏
    tokenPattern := regexp.MustCompile(`(token["\s:=]+)[a-zA-Z0-9_-]{20,}`)
    log = tokenPattern.ReplaceAllString(log, `$1****`)

    return log
}
```

### 前端组件设计

**DeployError.vue:**

```vue
<!-- frontend/src/components/instances/DeployError.vue -->
<template>
  <div class="deploy-error">
    <!-- 错误图标 -->
    <n-result status="error" title="部署失败">
      <template #footer>
        <div class="error-content">
          <!-- 错误信息 -->
          <n-alert type="error" :title="errorInfo.message">
            <p>错误代码: {{ errorInfo.code }}</p>
            <p>发生阶段: {{ stageName }}</p>
          </n-alert>

          <!-- 解决方案 -->
          <n-card title="解决方案" size="small">
            <p>{{ errorInfo.suggestion }}</p>
            <n-button text type="primary" @click="openDoc">
              查看文档
            </n-button>
          </n-card>

          <!-- 操作按钮 -->
          <n-space>
            <n-button type="primary" @click="handleRedeploy">
              <template #icon><n-icon :component="RefreshOutline" /></template>
              重新部署
            </n-button>
            <n-button @click="showLogDrawer = true">
              查看详细日志
            </n-button>
            <n-button @click="downloadLog">
              下载日志
            </n-button>
          </n-space>
        </div>
      </template>
    </n-result>

    <!-- 日志抽屉 -->
    <n-drawer v-model:show="showLogDrawer" width="600px">
      <n-drawer-content title="部署日志">
        <n-log :rows="20" :log="logContent" language="text" />
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { RefreshOutline } from '@vicons/ionicons5'

interface Props {
  errorInfo: ErrorInfo
  instanceId: string
}

const props = defineProps<Props>()
const emit = defineEmits(['redeploy'])

const showLogDrawer = ref(false)
const logContent = ref('')

const stageName = computed(() => {
  const stageNames: Record<string, string> = {
    validating: '验证配置',
    creating: '创建实例',
    starting: '启动服务',
  }
  return stageNames[props.errorInfo.stage] || props.errorInfo.stage
})

function handleRedeploy() {
  emit('redeploy', props.instanceId)
}

function openDoc() {
  window.open(props.errorInfo.docLink, '_blank')
}

async function downloadLog() {
  // 下载日志实现
}
</script>
```

### 项目结构规范

**新增文件位置:**

```
backend/
├── internal/
│   └── domain/
│       └── instance/
│           ├── errors.go         # 错误定义（新增）
│           ├── error_handler.go  # 错误处理（新增）
│           └── sanitizer.go      # 脱敏工具（新增）

frontend/
├── src/
│   └── components/
│       └── instances/
│           └── DeployError.vue   # 错误组件（新增）
```

### 测试标准

**测试用例:**

| 场景 | 测试方法 |
|-----|---------|
| 错误分类 | TestClassifyError |
| 错误记录 | TestHandleDeployError |
| 脱敏处理 | TestSanitizeLog |
| 重新部署 | TestRedeploy |
| 敏感信息隐藏 | TestSensitiveDataHidden |

### References

- [Source: prd.md#FR14] - 显示部署错误原因
- [Source: prd.md#NFR-S2] - 敏感信息加密存储
- [Source: epics.md#Story 4.5] - 原始故事定义

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
