# Story 3.4: 飞书配置前端页面

Status: ready-for-dev

## Story

As a 用户,
I want 在前端页面配置飞书应用信息,
so that 可以完成 OpenClaw 实例部署前的准备工作。

## Acceptance Criteria

1. **AC1: 飞书配置页面路由**
   - **Given** 前端项目已初始化
   - **When** 访问 /feishu-config 路径
   - **Then** 显示飞书配置页面
   - **And** 未登录用户重定向到登录页
   - **And** 页面布局符合 UX 设计规范

2. **AC2: App ID 输入框**
   - **Given** 用户在飞书配置页面
   - **When** 查看 App ID 输入区域
   - **Then** 显示标签"App ID"
   - **And** 显示占位符提示"请输入飞书应用 App ID"
   - **And** 实时校验格式（以 cli_ 开头）
   - **And** 格式错误显示"App ID 必须以 cli_ 开头"

3. **AC3: App Secret 输入框**
   - **Given** 用户在飞书配置页面
   - **When** 查看 App Secret 输入区域
   - **Then** 显示为密码输入框（内容不可见）
   - **And** 显示标签"App Secret"
   - **And** 显示占位符提示"请输入飞书应用 App Secret"
   - **And** 提供显示/隐藏切换按钮

4. **AC4: 验证配置按钮**
   - **Given** 用户已输入 App ID 和 Secret
   - **When** 点击"验证配置"按钮
   - **Then** 按钮显示加载状态
   - **And** 调用验证 API
   - **And** 验证成功显示"配置有效"提示
   - **And** 验证失败显示具体错误原因

5. **AC5: 保存配置按钮**
   - **Given** 用户已输入完整配置信息
   - **When** 点击"保存配置"按钮
   - **Then** 按钮显示加载状态
   - **And** 调用创建/更新 API
   - **And** 保存成功显示成功提示
   - **And** 保存失败显示错误信息

6. **AC6: 已有配置状态显示**
   - **Given** 用户已有飞书配置
   - **When** 访问飞书配置页面
   - **Then** 自动填充已有配置信息
   - **And** 显示配置状态（已验证/未验证）
   - **And** 显示最后更新时间
   - **And** App Secret 显示为占位符（不显示实际值）

## Tasks / Subtasks

- [ ] Task 1: 创建页面路由配置 (AC: 1)
  - [ ] 1.1 在 `src/router/index.ts` 添加飞书配置路由
  - [ ] 1.2 配置路由守卫（认证检查）
  - [ ] 1.3 设置页面标题和 meta 信息

- [ ] Task 2: 创建飞书配置页面组件 (AC: 1-6)
  - [ ] 2.1 创建 `src/pages/feishu-config/index.vue` 页面组件
  - [ ] 2.2 实现页面布局（使用 Naive UI 布局组件）
  - [ ] 2.3 添加页面标题和说明文字

- [ ] Task 3: 创建配置表单组件 (AC: 2, 3)
  - [ ] 3.1 创建 `src/components/config/FeishuConfigForm.vue` 表单组件
  - [ ] 3.2 实现 App ID 输入框（带格式校验）
  - [ ] 3.3 实现 App Secret 密码输入框（带显示/隐藏切换）
  - [ ] 3.4 实现实时表单验证
  - [ ] 3.5 添加表单验证规则

- [ ] Task 4: 创建 API 服务 (AC: 4, 5)
  - [ ] 4.1 创建 `src/services/feishu-config.ts` API 服务
  - [ ] 4.2 实现 getFeishuConfig() 方法
  - [ ] 4.3 实现 createFeishuConfig(data) 方法
  - [ ] 4.4 实现 updateFeishuConfig(data) 方法
  - [ ] 4.5 实现 validateFeishuConfig() 方法

- [ ] Task 5: 创建 Pinia Store (AC: 4, 5, 6)
  - [ ] 5.1 创建 `src/stores/feishu-config.ts` Store
  - [ ] 5.2 定义状态（config, loading, validationStatus）
  - [ ] 5.3 实现 fetchConfig action
  - [ ] 5.4 实现 saveConfig action
  - [ ] 5.5 实现 validateConfig action

- [ ] Task 6: 实现配置状态显示 (AC: 6)
  - [ ] 6.1 创建 `src/components/config/ConfigStatusCard.vue` 状态卡片组件
  - [ ] 6.2 显示配置验证状态
  - [ ] 6.3 显示最后更新时间
  - [ ] 6.4 显示 App Secret 占位符

- [ ] Task 7: 编写组件测试 (AC: 1-6)
  - [ ] 7.1 编写 `FeishuConfigForm.test.ts` 表单组件测试
  - [ ] 7.2 编写 `feishu-config.test.ts` Store 测试
  - [ ] 7.3 编写表单验证测试
  - [ ] 7.4 编写 API 服务测试

## Dev Notes

### 架构模式与约束

**必须遵循的前端架构 [Source: architecture.md]:**

1. **项目结构**:
   - `src/pages/` - 页面组件
   - `src/components/` - 可复用组件
   - `src/stores/` - Pinia 状态管理
   - `src/services/` - API 调用服务

2. **命名约定 [Source: architecture.md#Naming Patterns]:**
   - 组件名: PascalCase (例: `FeishuConfigForm`)
   - 文件名: kebab-case (例: `feishu-config-form.vue`)
   - 函数/变量: camelCase (例: `fetchConfig`, `validationStatus`)

3. **UI 组件库 [Source: architecture.md]:**
   - 使用 Naive UI 组件库
   - 遵循 Naive UI 设计规范

### 现有项目状态

**前端项目结构 [Source: 1-2-frontend-project-init.md]:**

```
frontend/
├── src/
│   ├── components/
│   │   ├── auth/              # ✅ 已存在
│   │   ├── instances/         # ✅ 已存在
│   │   ├── config/            # ✅ 已存在（空）
│   │   ├── dashboard/         # ✅ 已存在
│   │   └── ui/                # ✅ 已存在
│   ├── composables/           # ✅ 已存在
│   ├── pages/                 # ✅ 已存在
│   │   ├── login/             # ✅ 已存在
│   │   ├── dashboard/         # ✅ 已存在
│   │   ├── feishu-config/     # ✅ 已存在（空）
│   │   └── ...
│   ├── stores/                # ✅ 已存在
│   │   ├── auth.ts            # ✅ 已存在
│   │   └── config.ts          # ✅ 已存在（空）
│   ├── services/              # ✅ 已存在
│   ├── types/                 # ✅ 已存在
│   └── utils/                 # ✅ 已存在
```

### 技术栈要求

**核心依赖:**

| 依赖 | 用途 | 版本 |
|------|------|------|
| Vue | 前端框架 | 3.x |
| TypeScript | 类型支持 | 5.x |
| Naive UI | 组件库 | 2.x |
| Pinia | 状态管理 | 2.x |
| Vue Router | 路由 | 4.x |
| Axios | HTTP 客户端 | 1.x |

### 页面组件实现

**飞书配置页面:**

```vue
<!-- src/pages/feishu-config/index.vue -->
<template>
  <div class="feishu-config-page">
    <n-page-header title="飞书应用配置" subtitle="配置您的飞书应用信息">
      <template #extra>
        <n-button type="primary" @click="handleSave" :loading="saving">
          保存配置
        </n-button>
      </template>
    </n-page-header>

    <n-space vertical size="large" class="content">
      <!-- 配置状态卡片 -->
      <ConfigStatusCard
        v-if="configStore.config"
        :config="configStore.config"
      />

      <!-- 配置表单 -->
      <FeishuConfigForm
        v-model:appId="formData.appId"
        v-model:appSecret="formData.appSecret"
        :loading="configStore.loading"
        @validate="handleValidate"
      />

      <!-- 配置教程（Story 3.5 实现） -->
      <FeishuConfigTutorial />
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NPageHeader, NButton, NSpace, useMessage } from 'naive-ui'
import { useFeishuConfigStore } from '@/stores/feishu-config'
import FeishuConfigForm from '@/components/config/FeishuConfigForm.vue'
import ConfigStatusCard from '@/components/config/ConfigStatusCard.vue'
import FeishuConfigTutorial from '@/components/config/FeishuConfigTutorial.vue'

const message = useMessage()
const configStore = useFeishuConfigStore()

const formData = ref({
  appId: '',
  appSecret: ''
})
const saving = ref(false)

onMounted(async () => {
  await configStore.fetchConfig()
  if (configStore.config) {
    formData.value.appId = configStore.config.appId
    // AppSecret 不返回，保持为空
  }
})

const handleValidate = async () => {
  try {
    const result = await configStore.validateConfig()
    if (result.success) {
      message.success('配置验证成功')
    } else {
      message.error(result.errorMessage || '配置验证失败')
    }
  } catch (error) {
    message.error('验证失败，请稍后重试')
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await configStore.saveConfig(formData.value)
    message.success('配置保存成功')
  } catch (error) {
    message.error('保存失败，请稍后重试')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.feishu-config-page {
  padding: 24px;
  max-width: 800px;
  margin: 0 auto;
}

.content {
  margin-top: 24px;
}
</style>
```

### 表单组件实现

**飞书配置表单:**

```vue
<!-- src/components/config/FeishuConfigForm.vue -->
<template>
  <n-card title="飞书应用凭证" class="config-form-card">
    <n-form
      ref="formRef"
      :model="modelValue"
      :rules="rules"
      label-placement="left"
      label-width="120"
    >
      <n-form-item label="App ID" path="appId">
        <n-input
          :value="modelValue.appId"
          @update:value="updateAppId"
          placeholder="请输入飞书应用 App ID"
          :disabled="loading"
        />
        <template #feedback>
          <span v-if="appIdError" class="error-text">{{ appIdError }}</span>
          <span v-else class="hint-text">App ID 以 cli_ 开头</span>
        </template>
      </n-form-item>

      <n-form-item label="App Secret" path="appSecret">
        <n-input
          :value="modelValue.appSecret"
          @update:value="updateAppSecret"
          type="password"
          show-password-on="click"
          placeholder="请输入飞书应用 App Secret"
          :disabled="loading"
        />
      </n-form-item>

      <n-form-item :show-label="false">
        <n-space>
          <n-button
            type="primary"
            ghost
            @click="handleValidate"
            :loading="validating"
            :disabled="!isFormValid"
          >
            验证配置
          </n-button>
        </n-space>
      </n-form-item>
    </n-form>
  </n-card>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { NCard, NForm, NFormItem, NInput, NButton, NSpace } from 'naive-ui'

interface Props {
  appId: string
  appSecret: string
  loading?: boolean
}

interface Emits {
  (e: 'update:appId', value: string): void
  (e: 'update:appSecret', value: string): void
  (e: 'validate'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const modelValue = computed(() => ({
  appId: props.appId,
  appSecret: props.appSecret
}))

const validating = ref(false)
const appIdError = ref('')

const rules = {
  appId: [
    { required: true, message: '请输入 App ID' },
    {
      pattern: /^cli_/,
      message: 'App ID 必须以 cli_ 开头',
      trigger: ['blur', 'input']
    }
  ],
  appSecret: [
    { required: true, message: '请输入 App Secret' }
  ]
}

const isFormValid = computed(() => {
  return props.appId.startsWith('cli_') && props.appSecret.length > 0
})

const updateAppId = (value: string) => {
  emit('update:appId', value)
  // 实时校验
  if (value && !value.startsWith('cli_')) {
    appIdError.value = 'App ID 必须以 cli_ 开头'
  } else {
    appIdError.value = ''
  }
}

const updateAppSecret = (value: string) => {
  emit('update:appSecret', value)
}

const handleValidate = () => {
  emit('validate')
}
</script>

<style scoped>
.config-form-card {
  margin-bottom: 16px;
}

.error-text {
  color: #d03050;
}

.hint-text {
  color: #999;
}
</style>
```

### Pinia Store 实现

**飞书配置状态管理:**

```typescript
// src/stores/feishu-config.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as feishuConfigApi from '@/services/feishu-config'

export interface FeishuConfig {
  id: string
  tenantId: string
  appId: string
  status: 'pending' | 'active' | 'inactive'
  createdAt: string
  updatedAt: string
}

export interface ValidationResult {
  success: boolean
  message?: string
  errorMessage?: string
}

export const useFeishuConfigStore = defineStore('feishu-config', () => {
  const config = ref<FeishuConfig | null>(null)
  const loading = ref(false)
  const validationStatus = ref<ValidationResult | null>(null)

  const isConfigured = computed(() => !!config.value)
  const isActive = computed(() => config.value?.status === 'active')

  async function fetchConfig() {
    loading.value = true
    try {
      const response = await feishuConfigApi.getFeishuConfig()
      config.value = response.data
    } catch (error: any) {
      if (error.response?.status !== 404) {
        throw error
      }
      config.value = null
    } finally {
      loading.value = false
    }
  }

  async function saveConfig(data: { appId: string; appSecret: string }) {
    loading.value = true
    try {
      if (config.value) {
        const response = await feishuConfigApi.updateFeishuConfig(data)
        config.value = response.data
      } else {
        const response = await feishuConfigApi.createFeishuConfig(data)
        config.value = response.data
      }
    } finally {
      loading.value = false
    }
  }

  async function validateConfig() {
    const response = await feishuConfigApi.validateFeishuConfig()
    validationStatus.value = response.data
    return response.data
  }

  function reset() {
    config.value = null
    validationStatus.value = null
  }

  return {
    config,
    loading,
    validationStatus,
    isConfigured,
    isActive,
    fetchConfig,
    saveConfig,
    validateConfig,
    reset
  }
})
```

### API 服务实现

**飞书配置 API 服务:**

```typescript
// src/services/feishu-config.ts
import request from './api'
import type { ApiResponse } from '@/types/api'

export interface FeishuConfig {
  id: string
  tenantId: string
  appId: string
  status: string
  createdAt: string
  updatedAt: string
}

export interface CreateFeishuConfigRequest {
  appId: string
  appSecret: string
}

export interface ValidationResult {
  success: boolean
  message?: string
  errorMessage?: string
}

export function getFeishuConfig(): Promise<ApiResponse<FeishuConfig>> {
  return request.get('/v1/feishu-configs')
}

export function createFeishuConfig(
  data: CreateFeishuConfigRequest
): Promise<ApiResponse<FeishuConfig>> {
  return request.post('/v1/feishu-configs', data)
}

export function updateFeishuConfig(
  data: CreateFeishuConfigRequest
): Promise<ApiResponse<FeishuConfig>> {
  return request.put('/v1/feishu-configs', data)
}

export function deleteFeishuConfig(): Promise<void> {
  return request.delete('/v1/feishu-configs')
}

export function validateFeishuConfig(): Promise<ApiResponse<ValidationResult>> {
  return request.post('/v1/feishu-configs/validate')
}
```

### 路由配置

**Vue Router 配置:**

```typescript
// src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  // ... 其他路由
  {
    path: '/feishu-config',
    name: 'feishu-config',
    component: () => import('@/pages/feishu-config/index.vue'),
    meta: {
      requiresAuth: true,
      title: '飞书应用配置'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'login', query: { redirect: to.fullPath } })
  } else {
    next()
  }
})

export default router
```

### UX 设计规范

**来自 UX 设计文档的要求:**

| 元素 | 规范 |
|------|------|
| 主要按钮 | 蓝色填充 #1677FF |
| 次要按钮 | 蓝色边框（ghost 样式） |
| 成功反馈 | Message.success 3秒自动消失 |
| 错误反馈 | Message.error 需手动关闭 |
| 表单验证 | 实时校验 + 提交校验 |
| 加载状态 | 按钮显示 loading 图标 |
| 最小触控区域 | 44x44px |

### 测试标准

**测试要求:**
- 测试框架: Vitest + Vue Test Utils
- 测试覆盖率目标: >= 70%

**测试用例设计:**

| 测试场景 | 测试方法 | 预期结果 |
|---------|---------|---------|
| 表单渲染 | mount 组件 | 显示输入框和按钮 |
| App ID 格式验证 | 输入错误格式 | 显示错误提示 |
| App Secret 隐藏 | 查看输入框 | 内容不可见 |
| 验证按钮点击 | 点击验证 | 调用验证 API |
| 保存按钮点击 | 点击保存 | 调用保存 API |
| 已有配置加载 | 有配置数据 | 自动填充表单 |

### 前序 Story 的学习经验

**从 Story 3.3 (飞书配置验证功能) 获得的经验:**

1. **验证 API**: POST /v1/feishu-configs/validate
2. **验证结果缓存**: 前端可考虑显示"上次验证时间"
3. **错误提示**: 需要显示友好的错误信息

### References

- [Source: architecture.md#Frontend Architecture] - 前端架构
- [Source: architecture.md#Naming Patterns] - 命名约定
- [Source: ux-design-specification.md] - UX 设计规范
- [Source: prd.md#FR5-FR8] - 飞书配置需求
- [Source: epics.md#Story 3.4] - 原始故事定义
- [Source: 3-3-feishu-config-validation.md] - 验证 API 接口

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
