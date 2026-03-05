# Story 4.6: 用户实例列表页面

Status: ready-for-dev

## Story

As a 用户,
I want 查看自己的实例列表和状态,
so that 管理我的 OpenClaw 实例。

## Acceptance Criteria

1. **AC1: 实例列表显示**
   - **Given** 用户已登录
   - **When** 访问实例管理页面
   - **Then** 显示用户的实例列表
   - **And** 使用卡片式布局展示
   - **And** 支持分页或滚动加载

2. **AC2: 实例卡片信息显示**
   - **Given** 实例列表已加载
   - **When** 查看实例卡片
   - **Then** 显示实例名称
   - **And** 显示实例状态
   - **And** 显示创建时间

3. **AC3: 状态颜色标识**
   - **Given** 实例卡片已显示
   - **When** 查看状态标识
   - **Then** 运行状态显示绿色
   - **And** 停止状态显示灰色
   - **And** 错误状态显示红色
   - **And** 部署中状态显示蓝色

4. **AC4: 快捷操作入口**
   - **Given** 实例卡片已显示
   - **When** 查看操作按钮
   - **Then** 提供"在飞书中使用"快捷入口
   - **And** 提供"查看详情"入口
   - **And** 根据状态显示可用操作

5. **AC5: 空状态显示**
   - **Given** 用户无实例
   - **When** 访问实例管理页面
   - **Then** 显示空状态提示
   - **And** 显示"创建实例"按钮
   - **And** 提供快速开始引导

6. **AC6: 响应式布局**
   - **Given** 用户使用不同设备
   - **When** 访问实例管理页面
   - **Then** 桌面端显示多列卡片布局
   - **And** 移动端显示单列卡片布局
   - **And** 平板端显示两列卡片布局

## Tasks / Subtasks

- [ ] Task 1: 创建实例列表 API (AC: 1)
  - [ ] 1.1 创建 GET /v1/instances 端点
  - [ ] 1.2 实现分页参数处理
  - [ ] 1.3 实现租户过滤
  - [ ] 1.4 返回实例列表和总数

- [ ] Task 2: 创建实例卡片组件 (AC: 2, 3, 4)
  - [ ] 2.1 创建 `frontend/src/components/instances/InstanceCard.vue`
  - [ ] 2.2 实现实例名称显示
  - [ ] 2.3 实现状态徽章显示
  - [ ] 2.4 实现创建时间格式化
  - [ ] 2.5 实现操作按钮

- [ ] Task 3: 创建实例列表页面 (AC: 1, 5, 6)
  - [ ] 3.1 创建 `frontend/src/pages/instances/InstanceListPage.vue`
  - [ ] 3.2 实现卡片网格布局
  - [ ] 3.3 实现空状态显示
  - [ ] 3.4 实现响应式布局
  - [ ] 3.5 实现加载状态

- [ ] Task 4: 创建实例详情页面 (AC: 4)
  - [ ] 4.1 创建 `frontend/src/pages/instances/InstanceDetailPage.vue`
  - [ ] 4.2 显示实例基本信息
  - [ ] 4.3 显示部署日志
  - [ ] 4.4 显示配置信息

- [ ] Task 5: 实现飞书快捷入口 (AC: 4)
  - [ ] 5.1 创建飞书机器人链接生成逻辑
  - [ ] 5.2 实现"在飞书中使用"按钮
  - [ ] 5.3 跳转到飞书机器人对话

- [ ] Task 6: 编写测试 (AC: 1-6)
  - [ ] 6.1 测试列表 API
  - [ ] 6.2 测试卡片组件
  - [ ] 6.3 测试空状态显示
  - [ ] 6.4 测试响应式布局

## Dev Notes

### 架构模式与约束

**命名约定 [Source: architecture.md]:**
- 组件名: `PascalCase` (例: `InstanceCard`)
- 文件名: `kebab-case` (例: `instance-card.vue`)
- API 端点: `/v1/instances` (复数资源名)

### API 设计

**实例列表 API:**

```yaml
# GET /v1/instances?page=1&pageSize=10
# 响应
{
  "data": [
    {
      "id": "inst-xxx",
      "tenantId": "tenant-xxx",
      "name": "我的 OpenClaw",
      "status": "running",
      "containerId": "container-xxx",
      "createdAt": "2026-03-05T10:00:00Z",
      "updatedAt": "2026-03-05T10:03:00Z"
    }
  ],
  "error": null,
  "meta": {
    "total": 5,
    "page": 1,
    "pageSize": 10
  }
}
```

### 状态颜色映射

```typescript
// frontend/src/types/instance.ts
export const statusColors: Record<InstanceStatus, string> = {
  pending: '#999999',    // 灰色
  deploying: '#1677FF',  // 蓝色
  running: '#52C41A',    // 绿色
  stopped: '#8C8C8C',    // 灰色
  error: '#FF4D4F',      // 红色
}

export const statusLabels: Record<InstanceStatus, string> = {
  pending: '等待中',
  deploying: '部署中',
  running: '运行中',
  stopped: '已停止',
  error: '部署失败',
}
```

### 前端组件设计

**InstanceCard.vue:**

```vue
<!-- frontend/src/components/instances/InstanceCard.vue -->
<template>
  <n-card class="instance-card" hoverable>
    <template #header>
      <div class="card-header">
        <span class="instance-name">{{ instance.name }}</span>
        <n-tag :color="statusColor" size="small">
          {{ statusLabel }}
        </n-tag>
      </div>
    </template>

    <div class="card-content">
      <div class="info-item">
        <span class="label">创建时间</span>
        <span class="value">{{ formattedCreatedAt }}</span>
      </div>
      <div class="info-item">
        <span class="label">实例 ID</span>
        <span class="value">{{ instance.id.slice(0, 8) }}...</span>
      </div>
    </div>

    <template #footer>
      <n-space justify="space-between">
        <n-button
          v-if="instance.status === 'running'"
          type="primary"
          size="small"
          @click="openFeishu"
        >
          <template #icon><n-icon :component="ChatbubbleOutline" /></template>
          在飞书中使用
        </n-button>
        <n-button size="small" @click="viewDetail">
          查看详情
        </n-button>
      </n-space>
    </template>
  </n-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ChatbubbleOutline } from '@vicons/ionicons5'
import { statusColors, statusLabels } from '@/types/instance'
import { formatRelativeTime } from '@/utils/date'

interface Props {
  instance: Instance
}

const props = defineProps<Props>()
const router = useRouter()

const statusColor = computed(() => {
  return { color: statusColors[props.instance.status] }
})

const statusLabel = computed(() => {
  return statusLabels[props.instance.status]
})

const formattedCreatedAt = computed(() => {
  return formatRelativeTime(props.instance.createdAt)
})

function openFeishu() {
  // 打开飞书机器人对话
  const feishuUrl = `feishu://chat?bot_id=${props.instance.containerId}`
  window.open(feishuUrl, '_blank')
}

function viewDetail() {
  router.push(`/instances/${props.instance.id}`)
}
</script>

<style scoped>
.instance-card {
  transition: transform 0.2s;
}

.instance-card:hover {
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.instance-name {
  font-weight: 600;
  font-size: 16px;
}

.card-content {
  padding: 8px 0;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 4px 0;
}

.label {
  color: #666;
  font-size: 13px;
}

.value {
  font-size: 13px;
}
</style>
```

**InstanceListPage.vue:**

```vue
<!-- frontend/src/pages/instances/InstanceListPage.vue -->
<template>
  <div class="instance-list-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>我的实例</h2>
      <n-button type="primary" @click="createInstance">
        <template #icon><n-icon :component="AddOutline" /></template>
        创建实例
      </n-button>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <n-spin size="large" />
    </div>

    <!-- 空状态 -->
    <div v-else-if="instances.length === 0" class="empty-state">
      <n-empty description="还没有实例">
        <template #extra>
          <n-button type="primary" @click="createInstance">
            创建第一个实例
          </n-button>
        </template>
      </n-empty>
    </div>

    <!-- 实例列表 -->
    <div v-else class="instance-grid">
      <instance-card
        v-for="instance in instances"
        :key="instance.id"
        :instance="instance"
      />
    </div>

    <!-- 分页 -->
    <div v-if="total > pageSize" class="pagination">
      <n-pagination
        v-model:page="currentPage"
        :page-count="Math.ceil(total / pageSize)"
        @update:page="loadInstances"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { AddOutline } from '@vicons/ionicons5'
import InstanceCard from '@/components/instances/InstanceCard.vue'
import { fetchInstances } from '@/services/api'

const router = useRouter()

const loading = ref(true)
const instances = ref<Instance[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = 10

async function loadInstances() {
  loading.value = true
  try {
    const response = await fetchInstances({
      page: currentPage.value,
      pageSize,
    })
    instances.value = response.data
    total.value = response.meta.total
  } finally {
    loading.value = false
  }
}

function createInstance() {
  router.push('/instances/create')
}

onMounted(loadInstances)
</script>

<style scoped>
.instance-list-page {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.instance-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

@media (max-width: 768px) {
  .instance-grid {
    grid-template-columns: 1fr;
  }
}

.empty-state {
  padding: 100px 0;
  text-align: center;
}

.pagination {
  margin-top: 24px;
  display: flex;
  justify-content: center;
}
</style>
```

### 响应式布局设计

**断点设置:**

| 设备 | 断点 | 列数 |
|-----|------|------|
| 移动端 | < 768px | 1 列 |
| 平板 | 768px - 1024px | 2 列 |
| 桌面 | > 1024px | 3-4 列（auto-fill） |

```css
/* 响应式网格 */
.instance-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

@media (max-width: 768px) {
  .instance-grid {
    grid-template-columns: 1fr;
  }
}
```

### 项目结构规范

**新增文件位置:**

```
backend/
├── internal/
│   └── api/
│       └── instance/
│           └── handler.go         # 列表 API（更新）

frontend/
├── src/
│   ├── components/
│   │   └── instances/
│   │       └── InstanceCard.vue   # 实例卡片（新增）
│   ├── pages/
│   │   └── instances/
│   │       ├── InstanceListPage.vue   # 列表页（新增）
│   │       └── InstanceDetailPage.vue # 详情页（新增）
│   ├── types/
│   │   └── instance.ts            # 类型定义（更新）
│   └── utils/
│       └── date.ts                # 日期工具（新增）
```

### UX 设计要点

**布局策略 [Source: ux-design-specification.md]:**
- 实例管理界面使用卡片式仪表盘布局
- Web 桌面优先，平板 P1，移动端 P2

**交互要求:**
- 卡片悬停效果
- 状态颜色醒目
- 快捷操作便捷

### 测试标准

**测试用例:**

| 场景 | 测试方法 |
|-----|---------|
| 列表加载 | TestInstanceList_Load |
| 分页功能 | TestInstanceList_Pagination |
| 空状态显示 | TestInstanceList_Empty |
| 状态颜色 | TestInstanceCard_StatusColor |
| 响应式布局 | TestInstanceList_Responsive |

### References

- [Source: architecture.md#Project Structure] - 项目目录结构
- [Source: ux-design-specification.md] - UX 设计规范
- [Source: prd.md#FR10-FR12] - 实例管理需求
- [Source: epics.md#Story 4.6] - 原始故事定义

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
