# Story 4.3: 一键部署 OpenClaw 实例

Status: ready-for-dev

## Story

As a 用户,
I want 点击按钮一键部署 OpenClaw 实例,
so that 可以快速开始使用 AI Agent。

## Acceptance Criteria

1. **AC1: 部署任务创建**
   - **Given** 用户已完成飞书配置
   - **When** 点击"启动我的 OpenClaw"按钮
   - **Then** 系统创建部署任务并返回任务 ID
   - **And** 实例状态初始化为 pending
   - **And** 实例关联到当前租户

2. **AC2: Dokploy 容器创建**
   - **Given** 部署任务已创建
   - **When** 调用 Dokploy API 创建容器
   - **Then** 容器使用 OpenClaw 官方镜像
   - **And** 容器名称包含租户 ID 保证唯一性
   - **And** 容器成功启动

3. **AC3: 环境变量自动配置**
   - **Given** 用户已配置飞书应用信息
   - **When** 创建容器
   - **Then** 自动注入飞书 App ID 环境变量
   - **And** 自动注入飞书 App Secret 环境变量
   - **And** 敏感信息不在日志中暴露

4. **AC4: 网络与存储配置**
   - **Given** 容器已创建
   - **When** 配置网络和存储
   - **Then** 容器分配独立网络（与其他租户隔离）
   - **And** 容器分配独立存储卷
   - **And** 存储卷持久化用户数据

5. **AC5: 部署状态更新**
   - **Given** 部署过程进行中
   - **When** 各阶段完成
   - **Then** 状态从 pending -> deploying -> running 更新
   - **And** 记录部署日志
   - **And** 更新 ContainerID 字段

6. **AC6: 部署时间要求**
   - **Given** 部署任务开始
   - **When** 部署完成
   - **Then** 部署过程在 3 分钟内完成
   - **And** 超时自动标记为 error 状态

## Tasks / Subtasks

- [ ] Task 1: 创建部署 API 端点 (AC: 1)
  - [ ] 1.1 创建 `internal/api/instance/handler.go`
  - [ ] 1.2 实现 POST /v1/instances 端点
  - [ ] 1.3 验证用户飞书配置是否存在
  - [ ] 1.4 创建实例记录（状态 pending）
  - [ ] 1.5 返回实例 ID

- [ ] Task 2: 实现部署服务层 (AC: 2, 3, 4, 5)
  - [ ] 2.1 创建 `internal/domain/instance/service.go`
  - [ ] 2.2 实现 DeployInstance 方法
  - [ ] 2.3 获取用户飞书配置
  - [ ] 2.4 构建 Dokploy 创建请求
  - [ ] 2.5 调用 Dokploy API 创建容器
  - [ ] 2.6 更新实例状态和 ContainerID

- [ ] Task 3: 实现异步部署流程 (AC: 6)
  - [ ] 3.1 使用 Goroutine 异步执行部署
  - [ ] 3.2 实现部署超时控制（3 分钟）
  - [ ] 3.3 实现部署日志记录
  - [ ] 3.4 实现部署失败回滚

- [ ] Task 4: 配置容器模板 (AC: 2, 3, 4)
  - [ ] 4.1 创建 `internal/infrastructure/dokploy/template.go`
  - [ ] 4.2 定义 OpenClaw 容器模板
  - [ ] 4.3 配置默认资源限制
  - [ ] 4.4 配置网络隔离策略
  - [ ] 4.5 配置存储卷挂载

- [ ] Task 5: 前端部署触发组件 (AC: 1)
  - [ ] 5.1 创建 `frontend/src/pages/instances/DeployButton.vue`
  - [ ] 5.2 实现"启动我的 OpenClaw"按钮
  - [ ] 5.3 处理未配置飞书的情况
  - [ ] 5.4 调用部署 API
  - [ ] 5.5 跳转到部署进度页面

- [ ] Task 6: 编写测试 (AC: 1-6)
  - [ ] 6.1 编写部署 API 测试
  - [ ] 6.2 编写部署服务测试
  - [ ] 6.3 Mock Dokploy API 响应
  - [ ] 6.4 测试超时场景
  - [ ] 6.5 测试并发部署

## Dev Notes

### 架构模式与约束

**必须遵循的 Clean Architecture 原则 [Source: architecture.md]:**

1. **分层依赖**:
   - API Handler -> Service -> Repository
   - Service 层协调业务逻辑
   - Repository 层处理数据持久化

2. **命名约定 [Source: architecture.md]:**
   - API 端点: `/v1/instances` (复数资源名)
   - 服务方法: `DeployInstance`
   - 文件名: `kebab-case`

### API 设计

**部署 API:**

```yaml
# POST /v1/instances
# 请求
{
  "name": "我的 OpenClaw"  // 可选，默认自动生成
}

# 响应
{
  "data": {
    "id": "inst-xxx",
    "tenantId": "tenant-xxx",
    "name": "我的 OpenClaw",
    "status": "pending",
    "createdAt": "2026-03-05T10:00:00Z"
  },
  "error": null,
  "meta": {}
}
```

### 部署流程设计

**部署状态机:**

```
pending -> deploying -> running
    |          |
    v          v
  error      error
```

**部署阶段:**

1. **验证配置** (约 5 秒)
   - 检查飞书配置是否存在
   - 验证飞书配置有效性

2. **创建实例记录** (约 1 秒)
   - 创建数据库记录
   - 状态设为 pending

3. **调用 Dokploy API** (约 60-120 秒)
   - 拉取镜像
   - 创建容器
   - 启动容器
   - 状态设为 deploying

4. **健康检查** (约 30 秒)
   - 检查容器运行状态
   - 状态设为 running 或 error

### 服务层设计

**InstanceService:**

```go
// internal/domain/instance/service.go
package instance

import (
    "context"
    "time"
    "github.com/gowa/saas-openclaw/backend/internal/domain/config"
    "github.com/gowa/saas-openclaw/backend/internal/infrastructure/dokploy"
)

// DeployConfig 部署配置
type DeployConfig struct {
    Timeout       time.Duration // 部署超时，默认 3 分钟
    ImageName     string        // OpenClaw 镜像名
    NetworkPrefix string        // 网络前缀
    VolumePrefix  string        // 存储卷前缀
}

// Service 实例服务
type Service struct {
    repo          InstanceRepository
    configRepo    config.Repository
    dokployClient *dokploy.Client
    deployConfig  DeployConfig
}

// DeployInstance 部署 OpenClaw 实例
func (s *Service) DeployInstance(ctx context.Context, tenantID, name string) (*Instance, error) {
    // 1. 获取飞书配置
    feishuConfig, err := s.configRepo.GetByTenantID(tenantID)
    if err != nil {
        return nil, ErrFeishuConfigNotFound
    }

    // 2. 创建实例记录
    inst := &Instance{
        TenantID: tenantID,
        Name:     name,
        Status:   InstanceStatusPending,
    }
    if err := s.repo.Create(inst); err != nil {
        return nil, err
    }

    // 3. 异步执行部署
    go s.executeDeployment(context.Background(), inst, feishuConfig)

    return inst, nil
}

// executeDeployment 执行部署
func (s *Service) executeDeployment(ctx context.Context, inst *Instance, fc *config.FeishuConfig) {
    // 设置超时
    ctx, cancel := context.WithTimeout(ctx, s.deployConfig.Timeout)
    defer cancel()

    // 更新状态为 deploying
    s.updateStatus(inst.ID, InstanceStatusDeploying, "开始部署")

    // 构建 Dokploy 请求
    req := s.buildDokployRequest(inst, fc)

    // 调用 Dokploy API
    app, err := s.dokployClient.CreateApplication(ctx, req)
    if err != nil {
        s.handleDeployError(inst.ID, err)
        return
    }

    // 更新 ContainerID
    inst.ContainerID = app.ContainerID
    s.repo.Update(inst)

    // 等待容器启动
    if err := s.waitForRunning(ctx, app.ID); err != nil {
        s.handleDeployError(inst.ID, err)
        return
    }

    // 更新状态为 running
    s.updateStatus(inst.ID, InstanceStatusRunning, "部署成功")
}
```

### 容器模板设计

**OpenClaw 容器配置:**

```go
// internal/infrastructure/dokploy/template.go
package dokploy

import "fmt"

// OpenClawTemplate OpenClaw 容器模板
type OpenClawTemplate struct {
    BaseImage    string
    NetworkName  string
    VolumeName   string
    ResourceLimits *ResourceConfig
}

// DefaultOpenClawTemplate 默认模板
var DefaultOpenClawTemplate = OpenClawTemplate{
    BaseImage: "openclaw/openclaw:latest",
    ResourceLimits: &ResourceConfig{
        CPU:    "1",
        Memory: "512M",
    },
}

// BuildCreateRequest 构建创建请求
func (t *OpenClawTemplate) BuildCreateRequest(tenantID string, appID, appSecret string) CreateApplicationRequest {
    return CreateApplicationRequest{
        Name:     fmt.Sprintf("openclaw-%s", tenantID),
        Image:    t.BaseImage,
        Env: []string{
            fmt.Sprintf("FEISHU_APP_ID=%s", appID),
            fmt.Sprintf("FEISHU_APP_SECRET=%s", appSecret),
        },
        Networks: []string{fmt.Sprintf("openclaw-net-%s", tenantID)},
        Volumes: []VolumeConfig{
            {
                Source: fmt.Sprintf("openclaw-data-%s", tenantID),
                Target: "/app/data",
            },
        },
        Resources: t.ResourceLimits,
    }
}
```

### 前端组件设计

**DeployButton.vue:**

```vue
<!-- frontend/src/pages/instances/DeployButton.vue -->
<template>
  <n-button
    type="primary"
    size="large"
    :loading="deploying"
    :disabled="!canDeploy"
    @click="handleDeploy"
  >
    <template #icon>
      <n-icon><RocketOutline /></n-icon>
    </template>
    启动我的 OpenClaw
  </n-button>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { useConfigStore } from '@/stores/config'
import { deployInstance } from '@/services/api'

const router = useRouter()
const message = useMessage()
const configStore = useConfigStore()

const deploying = ref(false)

const canDeploy = computed(() => {
  return configStore.hasValidConfig
})

async function handleDeploy() {
  if (!canDeploy.value) {
    message.warning('请先完成飞书配置')
    return
  }

  deploying.value = true
  try {
    const instance = await deployInstance()
    message.success('部署任务已创建')
    router.push(`/instances/${instance.id}/deploy-progress`)
  } catch (error) {
    message.error('部署失败：' + error.message)
  } finally {
    deploying.value = false
  }
}
</script>
```

### 性能要求

**部署时间要求 [Source: prd.md#NFR-P1]:**
- 实例部署在 3 分钟内完成
- 每个阶段有明确的时间限制

**并发部署:**
- 支持多个租户同时部署
- 使用 Goroutine 异步执行

### 错误处理

**部署错误类型:**

```go
var (
    ErrFeishuConfigNotFound = errors.New("飞书配置未找到")
    ErrFeishuConfigInvalid  = errors.New("飞书配置无效")
    ErrDeployTimeout        = errors.New("部署超时")
    ErrDokployUnavailable   = errors.New("Dokploy 服务不可用")
    ErrInsufficientResource = errors.New("资源不足")
)
```

### 安全考虑

**敏感信息保护:**
- App Secret 不在日志中打印
- 环境变量通过安全方式传递
- 容器间网络隔离

### 项目结构规范

**新增文件位置:**

```
backend/
├── internal/
│   ├── api/
│   │   └── instance/
│   │       └── handler.go        # 实例 API（新增）
│   ├── domain/
│   │   └── instance/
│   │       └── service.go        # 实例服务（新增）
│   └── infrastructure/
│       └── dokploy/
│           └── template.go       # 容器模板（新增）

frontend/
├── src/
│   ├── pages/
│   │   └── instances/
│   │       └── DeployButton.vue  # 部署按钮组件（新增）
│   └── services/
│       └── api.ts                # API 服务（更新）
```

### 测试标准

**测试要求:**
- 测试框架: Go testing + testify
- 前端测试: Vitest + Vue Test Utils
- 测试覆盖率目标: >= 70%

**关键测试场景:**

| 场景 | 测试方法 |
|-----|---------|
| 正常部署流程 | TestDeployInstance_Success |
| 飞书配置缺失 | TestDeployInstance_NoConfig |
| 部署超时 | TestDeployInstance_Timeout |
| Dokploy 错误 | TestDeployInstance_DokployError |
| 并发部署 | TestDeployInstance_Concurrent |

### References

- [Source: architecture.md#API Boundaries] - API 边界设计
- [Source: prd.md#FR10-FR11] - 实例部署需求
- [Source: prd.md#NFR-P1] - 部署时间要求
- [Source: prd.md#NFR-S1] - 容器隔离要求
- [Source: epics.md#Story 4.3] - 原始故事定义

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
