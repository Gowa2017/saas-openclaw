# Story 4.7: 管理员实例管理功能

Status: ready-for-dev

## Story

As a 平台管理员,
I want 管理所有用户的实例,
so that 可以处理问题实例。

## Acceptance Criteria

1. **AC1: 所有实例列表显示**
   - **Given** 管理员已登录管理后台
   - **When** 访问实例管理页面
   - **Then** 显示所有用户的实例列表
   - **And** 显示实例所属租户信息
   - **And** 支持分页浏览

2. **AC2: 实例筛选功能**
   - **Given** 实例列表已显示
   - **When** 使用筛选功能
   - **Then** 可以按用户筛选实例
   - **And** 可以按状态筛选实例
   - **And** 可以按时间范围筛选
   - **And** 筛选结果实时更新

3. **AC3: 实例详情查看**
   - **Given** 管理员查看实例
   - **When** 点击实例详情
   - **Then** 显示租户信息（名称、ID、联系方式）
   - **And** 显示容器信息（ID、镜像、网络）
   - **And** 显示实例日志

4. **AC4: 实例操作功能**
   - **Given** 管理员选择实例
   - **When** 执行操作
   - **Then** 提供重启操作按钮
   - **And** 提供停止操作按钮
   - **And** 提供启动操作按钮
   - **And** 操作按钮根据状态智能显示/隐藏

4. **AC5: 操作二次确认**
   - **Given** 管理员点击操作按钮
   - **When** 执行敏感操作
   - **Then** 弹出确认对话框
   - **And** 显示操作影响说明
   - **And** 需要输入确认文字

6. **AC6: 操作结果实时反馈**
   - **Given** 操作已执行
   - **When** 操作完成
   - **Then** 显示操作结果（成功/失败）
   - **And** 更新实例状态
   - **And** 记录操作日志

## Tasks / Subtasks

- [ ] Task 1: 创建管理员实例列表 API (AC: 1)
  - [ ] 1.1 创建 GET /v1/admin/instances 端点
  - [ ] 1.2 实现分页参数
  - [ ] 1.3 关联租户信息
  - [ ] 1.4 添加管理员权限验证

- [ ] Task 2: 创建实例筛选 API (AC: 2)
  - [ ] 2.1 实现按用户筛选
  - [ ] 2.2 实现按状态筛选
  - [ ] 2.3 实现按时间范围筛选
  - [ ] 2.4 实现组合筛选

- [ ] Task 3: 创建实例操作 API (AC: 4, 5, 6)
  - [ ] 3.1 创建 POST /v1/admin/instances/:id/restart 端点
  - [ ] 3.2 创建 POST /v1/admin/instances/:id/stop 端点
  - [ ] 3.3 创建 POST /v1/admin/instances/:id/start 端点
  - [ ] 3.4 实现操作日志记录

- [ ] Task 4: 创建管理员实例列表页面 (AC: 1, 2)
  - [ ] 4.1 创建 `frontend/src/pages/admin/InstanceManagePage.vue`
  - [ ] 4.2 实现实例表格展示
  - [ ] 4.3 实现筛选表单
  - [ ] 4.4 实现分页组件

- [ ] Task 5: 创建实例详情对话框 (AC: 3)
  - [ ] 5.1 创建 `frontend/src/components/admin/InstanceDetailModal.vue`
  - [ ] 5.2 显示租户信息
  - [ ] 5.3 显示容器信息
  - [ ] 5.4 显示实时日志

- [ ] Task 6: 创建操作确认对话框 (AC: 5, 6)
  - [ ] 6.1 创建 `frontend/src/components/admin/OperationConfirmModal.vue`
  - [ ] 6.2 实现确认对话框
  - [ ] 6.3 实现确认文字输入
  - [ ] 6.4 实现操作结果反馈

- [ ] Task 7: 编写测试 (AC: 1-6)
  - [ ] 7.1 测试管理员 API 权限
  - [ ] 7.2 测试筛选功能
  - [ ] 7.3 测试操作功能
  - [ ] 7.4 测试前端组件

## Dev Notes

### 架构模式与约束

**权限控制:**
- 所有管理员 API 需要验证管理员身份
- 操作日志需要记录操作人和时间

**命名约定 [Source: architecture.md]:**
- 管理员 API 前缀: `/v1/admin/`
- 组件名: `PascalCase`

### API 设计

**管理员实例列表 API:**

```yaml
# GET /v1/admin/instances?page=1&pageSize=20&userId=&status=&startTime=&endTime=
# 响应
{
  "data": [
    {
      "id": "inst-xxx",
      "tenantId": "tenant-xxx",
      "tenantName": "示例公司",
      "tenantEmail": "admin@example.com",
      "name": "我的 OpenClaw",
      "status": "running",
      "containerId": "container-xxx",
      "createdAt": "2026-03-05T10:00:00Z",
      "updatedAt": "2026-03-05T10:03:00Z"
    }
  ],
  "error": null,
  "meta": {
    "total": 100,
    "page": 1,
    "pageSize": 20
  }
}
```

**实例操作 API:**

```yaml
# POST /v1/admin/instances/:id/restart
# 请求体
{
  "reason": "服务无响应"  // 操作原因（可选）
}

# 响应
{
  "data": {
    "id": "inst-xxx",
    "status": "running",
    "operationId": "op-xxx",
    "message": "重启操作已执行"
  },
  "error": null,
  "meta": {}
}
```

**实例详情 API:**

```yaml
# GET /v1/admin/instances/:id
# 响应
{
  "data": {
    "id": "inst-xxx",
    "tenantId": "tenant-xxx",
    "name": "我的 OpenClaw",
    "status": "running",
    "containerId": "container-xxx",
    "tenant": {
      "id": "tenant-xxx",
      "name": "示例公司",
      "email": "admin@example.com",
      "contactName": "张三"
    },
    "container": {
      "id": "container-xxx",
      "image": "openclaw/openclaw:latest",
      "network": "openclaw-net-xxx",
      "cpu": "1",
      "memory": "512M"
    },
    "deployLog": "部署成功...",
    "createdAt": "2026-03-05T10:00:00Z"
  },
  "error": null,
  "meta": {}
}
```

### 后端实现

**管理员 Handler:**

```go
// internal/api/admin/instance_handler.go
package admin

import (
    "github.com/gin-gonic/gin"
    "github.com/gowa/saas-openclaw/backend/internal/domain/instance"
)

// InstanceHandler 管理员实例处理器
type InstanceHandler struct {
    instanceService *instance.Service
    adminRepo       AdminRepository
}

// ListInstances 获取所有实例列表
func (h *InstanceHandler) ListInstances(c *gin.Context) {
    // 解析筛选参数
    filter := instance.ListFilter{
        UserID:    c.Query("userId"),
        Status:    instance.InstanceStatus(c.Query("status")),
        StartTime: c.Query("startTime"),
        EndTime:   c.Query("endTime"),
        Page:      parseInt(c.Query("page"), 1),
        PageSize:  parseInt(c.Query("pageSize"), 20),
    }

    // 验证管理员权限
    adminID := c.GetString("adminID")
    if !h.hasPermission(adminID, "instance:read") {
        c.JSON(403, gin.H{"error": "无权限"})
        return
    }

    // 查询实例列表
    instances, total, err := h.instanceService.ListAll(filter)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "data": instances,
        "meta": gin.H{
            "total":    total,
            "page":     filter.Page,
            "pageSize": filter.PageSize,
        },
    })
}

// RestartInstance 重启实例
func (h *InstanceHandler) RestartInstance(c *gin.Context) {
    instanceID := c.Param("id")
    adminID := c.GetString("adminID")

    // 记录操作日志
    h.logOperation(adminID, instanceID, "restart", c.Query("reason"))

    // 执行重启
    if err := h.instanceService.Restart(instanceID); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "data": gin.H{
            "id":      instanceID,
            "status":  "running",
            "message": "重启操作已执行",
        },
    })
}
```

**操作日志记录:**

```go
// internal/domain/admin/operation_log.go
package admin

import "time"

// OperationLog 操作日志
type OperationLog struct {
    ID          string    `json:"id" db:"ID"`
    AdminID     string    `json:"adminId" db:"AdminID"`
    InstanceID  string    `json:"instanceId" db:"InstanceID"`
    Operation   string    `json:"operation" db:"Operation"`
    Reason      string    `json:"reason" db:"Reason"`
    Result      string    `json:"result" db:"Result"`
    CreatedAt   time.Time `json:"createdAt" db:"CreatedAt"`
}

// 操作类型
const (
    OpRestart = "restart"
    OpStop    = "stop"
    OpStart   = "start"
)
```

### 前端组件设计

**InstanceManagePage.vue:**

```vue
<!-- frontend/src/pages/admin/InstanceManagePage.vue -->
<template>
  <div class="instance-manage-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>实例管理</h2>
    </div>

    <!-- 筛选表单 -->
    <n-card class="filter-card" size="small">
      <n-form :model="filterForm" inline>
        <n-form-item label="用户">
          <n-input v-model:value="filterForm.userId" placeholder="用户 ID" />
        </n-form-item>
        <n-form-item label="状态">
          <n-select
            v-model:value="filterForm.status"
            :options="statusOptions"
            placeholder="选择状态"
            clearable
          />
        </n-form-item>
        <n-form-item label="时间范围">
          <n-date-picker
            v-model:value="filterForm.timeRange"
            type="daterange"
            clearable
          />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" @click="handleSearch">搜索</n-button>
          <n-button @click="handleReset">重置</n-button>
        </n-form-item>
      </n-form>
    </n-card>

    <!-- 实例表格 -->
    <n-card>
      <n-data-table
        :columns="columns"
        :data="instances"
        :loading="loading"
        :pagination="pagination"
        @update:page="handlePageChange"
      />
    </n-card>

    <!-- 详情对话框 -->
    <instance-detail-modal
      v-model:show="showDetail"
      :instance-id="selectedInstanceId"
    />

    <!-- 操作确认对话框 -->
    <operation-confirm-modal
      v-model:show="showConfirm"
      :operation="currentOperation"
      @confirm="handleConfirmOperation"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h } from 'vue'
import { NButton, NTag, NSpace } from 'naive-ui'
import InstanceDetailModal from '@/components/admin/InstanceDetailModal.vue'
import OperationConfirmModal from '@/components/admin/OperationConfirmModal.vue'
import { fetchAdminInstances, restartInstance, stopInstance, startInstance } from '@/services/admin-api'

const loading = ref(false)
const instances = ref([])
const showDetail = ref(false)
const showConfirm = ref(false)
const selectedInstanceId = ref('')
const currentOperation = ref({ type: '', instanceId: '' })

const filterForm = reactive({
  userId: '',
  status: null,
  timeRange: null,
})

const statusOptions = [
  { label: '运行中', value: 'running' },
  { label: '已停止', value: 'stopped' },
  { label: '部署中', value: 'deploying' },
  { label: '部署失败', value: 'error' },
]

const columns = [
  { title: '实例 ID', key: 'id', width: 120 },
  { title: '实例名称', key: 'name', width: 150 },
  { title: '租户', key: 'tenantName', width: 150 },
  { title: '状态', key: 'status', width: 100, render: renderStatus },
  { title: '创建时间', key: 'createdAt', width: 180 },
  {
    title: '操作',
    key: 'actions',
    width: 250,
    render: (row) => renderActions(row),
  },
]

function renderStatus(row) {
  const colors = {
    running: 'success',
    stopped: 'default',
    deploying: 'info',
    error: 'error',
  }
  return h(NTag, { type: colors[row.status] || 'default', size: 'small' }, () => row.status)
}

function renderActions(row) {
  const buttons = []

  // 查看详情
  buttons.push(
    h(NButton, { size: 'small', onClick: () => viewDetail(row.id) }, () => '详情')
  )

  // 根据状态显示操作按钮
  if (row.status === 'running') {
    buttons.push(
      h(NButton, { size: 'small', type: 'warning', onClick: () => handleOperation('restart', row.id) }, () => '重启'),
      h(NButton, { size: 'small', type: 'error', onClick: () => handleOperation('stop', row.id) }, () => '停止')
    )
  } else if (row.status === 'stopped' || row.status === 'error') {
    buttons.push(
      h(NButton, { size: 'small', type: 'success', onClick: () => handleOperation('start', row.id) }, () => '启动')
    )
  }

  return h(NSpace, null, () => buttons)
}

function viewDetail(id) {
  selectedInstanceId.value = id
  showDetail.value = true
}

function handleOperation(type, instanceId) {
  currentOperation.value = { type, instanceId }
  showConfirm.value = true
}

async function handleConfirmOperation(reason) {
  const { type, instanceId } = currentOperation.value
  try {
    if (type === 'restart') {
      await restartInstance(instanceId, reason)
    } else if (type === 'stop') {
      await stopInstance(instanceId, reason)
    } else if (type === 'start') {
      await startInstance(instanceId, reason)
    }
    loadInstances()
  } finally {
    showConfirm.value = false
  }
}

async function loadInstances() {
  loading.value = true
  try {
    const res = await fetchAdminInstances(filterForm)
    instances.value = res.data
  } finally {
    loading.value = false
  }
}
</script>
```

**OperationConfirmModal.vue:**

```vue
<!-- frontend/src/components/admin/OperationConfirmModal.vue -->
<template>
  <n-modal v-model:show="show" preset="dialog" title="确认操作">
    <n-alert type="warning" title="警告">
      您即将对实例执行 <strong>{{ operationTitle }}</strong> 操作。
      此操作可能会影响用户正常使用，请确认。
    </n-alert>

    <n-form ref="formRef" :model="form" :rules="rules" style="margin-top: 16px">
      <n-form-item label="操作原因" path="reason">
        <n-input
          v-model:value="form.reason"
          type="textarea"
          placeholder="请输入操作原因（必填）"
        />
      </n-form-item>
      <n-form-item label="确认文字" path="confirmText">
        <n-input
          v-model:value="form.confirmText"
          placeholder="请输入 '确认' 以继续"
        />
      </n-form-item>
    </n-form>

    <template #action>
      <n-space justify="end">
        <n-button @click="$emit('update:show', false)">取消</n-button>
        <n-button
          type="primary"
          :disabled="!canConfirm"
          @click="handleConfirm"
        >
          确认执行
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed, reactive } from 'vue'

interface Props {
  show: boolean
  operation: { type: string; instanceId: string }
}

const props = defineProps<Props>()
const emit = defineEmits(['update:show', 'confirm'])

const form = reactive({
  reason: '',
  confirmText: '',
})

const rules = {
  reason: { required: true, message: '请输入操作原因' },
}

const operationTitle = computed(() => {
  const titles = {
    restart: '重启实例',
    stop: '停止实例',
    start: '启动实例',
  }
  return titles[props.operation.type] || '操作'
})

const canConfirm = computed(() => {
  return form.reason && form.confirmText === '确认'
})

function handleConfirm() {
  emit('confirm', form.reason)
}
</script>
```

### 项目结构规范

**新增文件位置:**

```
backend/
├── internal/
│   ├── api/
│   │   └── admin/
│   │       └── instance_handler.go  # 管理员实例 API（新增）
│   └── domain/
│       └── admin/
│           └── operation_log.go     # 操作日志（新增）

frontend/
├── src/
│   ├── pages/
│   │   └── admin/
│   │       └── InstanceManagePage.vue  # 管理页面（新增）
│   ├── components/
│   │   └── admin/
│   │       ├── InstanceDetailModal.vue  # 详情对话框（新增）
│   │       └── OperationConfirmModal.vue # 确认对话框（新增）
│   └── services/
│       └── admin-api.ts            # 管理员 API（新增）
```

### 权限控制

**管理员权限验证中间件:**

```go
// pkg/middleware/admin_auth.go
package middleware

// AdminAuth 管理员认证中间件
func AdminAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "未授权"})
            c.Abort()
            return
        }

        // 验证 JWT 并获取管理员信息
        adminID, role, err := validateAdminToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "无效的令牌"})
            c.Abort()
            return
        }

        c.Set("adminID", adminID)
        c.Set("adminRole", role)
        c.Next()
    }
}
```

### 测试标准

**测试用例:**

| 场景 | 测试方法 |
|-----|---------|
| 列表加载 | TestAdminInstanceList |
| 权限验证 | TestAdminAuth |
| 筛选功能 | TestInstanceFilter |
| 重启操作 | TestRestartInstance |
| 操作日志 | TestOperationLog |

### References

- [Source: prd.md#FR15-FR18] - 管理员实例管理需求
- [Source: prd.md#NFR-S3] - 管理员操作验证要求
- [Source: epics.md#Story 4.7] - 原始故事定义

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
