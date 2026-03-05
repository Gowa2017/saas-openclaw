# Story 2.6: 管理员登录页面

Status: ready-for-dev

## Story

As a 平台管理员,
I want 访问独立的管理员登录页面,
so that 可以登录管理后台。

## Acceptance Criteria

1. **AC1: 管理员登录页面路由**
   - **Given** 管理员登录页面已创建
   - **When** 访问 /admin/login 路径
   - **Then** 显示管理员登录表单

2. **AC2: 登录表单实现**
   - **Given** 登录表单已渲染
   - **When** 查看表单内容
   - **Then** 显示用户名输入框
   - **And** 显示密码输入框
   - **And** 显示登录按钮

3. **AC3: 表单验证**
   - **Given** 用户填写表单
   - **When** 输入或提交时
   - **Then** 实时校验必填项
   - **And** 实时校验格式（用户名长度、密码长度）
   - **And** 显示验证错误提示

4. **AC4: 登录成功处理**
   - **Given** 登录验证通过
   - **When** 后端返回成功
   - **Then** 跳转到管理后台首页
   - **And** 存储管理员 Token

5. **AC5: 登录失败处理**
   - **Given** 登录验证失败
   - **When** 后端返回错误
   - **Then** 显示错误提示信息
   - **And** 不清空已填写的用户名

6. **AC6: 记住登录状态**
   - **Given** 用户勾选"记住我"
   - **When** 登录成功
   - **Then** 延长 Token 有效期
   - **And** 下次访问自动登录

## Tasks / Subtasks

- [ ] Task 1: 创建管理员认证 Store (AC: 4, 6)
  - [ ] 1.1 创建 `src/stores/admin-auth.ts`
  - [ ] 1.2 实现 login action
  - [ ] 1.3 实现 logout action
  - [ ] 1.4 实现 rememberMe 功能
  - [ ] 1.5 实现 autoLogin 功能

- [ ] Task 2: 创建管理员登录页面 (AC: 1, 2)
  - [ ] 2.1 创建 `src/pages/admin/login/index.vue`
  - [ ] 2.2 实现用户名输入框
  - [ ] 2.3 实现密码输入框
  - [ ] 2.4 实现"记住我"复选框
  - [ ] 2.5 实现登录按钮

- [ ] Task 3: 实现表单验证 (AC: 3)
  - [ ] 3.1 使用 Naive UI Form 组件
  - [ ] 3.2 定义验证规则
  - [ ] 3.3 实现实时校验
  - [ ] 3.4 实现提交校验

- [ ] Task 4: 实现登录逻辑 (AC: 4, 5)
  - [ ] 4.1 调用管理员登录 API
  - [ ] 4.2 处理登录成功
  - [ ] 4.3 处理登录失败
  - [ ] 4.4 实现加载状态

- [ ] Task 5: 实现管理员路由守卫 (AC: 1, 4)
  - [ ] 5.1 创建 `src/router/guards/admin-auth.ts`
  - [ ] 5.2 实现管理员认证检查
  - [ ] 5.3 实现管理员路由重定向

- [ ] Task 6: 编写单元测试 (AC: 1-6)
  - [ ] 6.1 编写 `admin-auth.spec.ts` 测试 Store
  - [ ] 6.2 编写 `admin-login.spec.ts` 测试页面
  - [ ] 6.3 编写 `admin-auth-guard.spec.ts` 测试守卫

## Dev Notes

### 架构模式与约束

**必须遵循的架构原则 [Source: architecture.md]:**

1. **技术栈**:
   - Vue 3 + TypeScript
   - Pinia 状态管理
   - Vue Router 路由
   - Naive UI 组件库

2. **命名约定 [Source: architecture.md#Naming Patterns]:**
   - 组件名: PascalCase (例: `AdminLoginPage`)
   - 文件名: kebab-case (例: `admin-login.vue`)
   - 路由路径: `/admin/login`

3. **API 端点**:
   - POST `/v1/admin/auth/login`

### 现有项目状态

**依赖 Story 2.4 管理员认证系统和 Story 2.5 前端登录:**

```
frontend/
├── src/
│   ├── stores/
│   │   └── auth.ts              # 租户认证 Store
│   ├── pages/
│   │   └── login/
│   │       └── index.vue        # 租户登录页面
│   ├── router/
│   │   └── guards/
│   │       └── auth.ts          # 租户认证守卫
│   ├── services/
│   │   ├── api.ts               # API 客户端
│   │   └── auth.ts              # 认证服务
│   └── utils/
│       └── token.ts             # Token 工具
```

### 管理员认证 Store 设计

**Store 实现:**

```typescript
// src/stores/admin-auth.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { AdminUser } from '@/types/models'
import { getAdminToken, saveAdminToken, removeAdminToken } from '@/utils/admin-token'
import { adminLogin as adminLoginApi } from '@/services/admin-auth'

export const useAdminAuthStore = defineStore('adminAuth', () => {
  // State
  const admin = ref<AdminUser | null>(null)
  const token = ref<string | null>(getAdminToken())
  const isLoading = ref(false)
  const rememberMe = ref(false)

  // Getters
  const isAuthenticated = computed(() => !!token.value && !!admin.value)
  const adminId = computed(() => admin.value?.id)
  const adminRole = computed(() => admin.value?.role)

  // Actions
  async function login(username: string, password: string, remember: boolean) {
    isLoading.value = true
    rememberMe.value = remember

    try {
      const response = await adminLoginApi({ username, password })
      token.value = response.token
      admin.value = response.admin

      // 根据记住登录设置存储策略
      const expiresIn = remember ? 7 * 24 * 60 * 60 : 24 * 60 * 60 // 7天 or 1天
      saveAdminToken(response.token, expiresIn)

      return { success: true }
    } catch (error: any) {
      return {
        success: false,
        message: error.response?.data?.error?.message || '登录失败，请重试'
      }
    } finally {
      isLoading.value = false
    }
  }

  function logout() {
    token.value = null
    admin.value = null
    removeAdminToken()
  }

  function autoLogin() {
    const savedToken = getAdminToken()
    if (savedToken) {
      token.value = savedToken
      // 可以在这里调用 API 验证 Token 并获取用户信息
    }
  }

  return {
    // State
    admin,
    token,
    isLoading,
    rememberMe,
    // Getters
    isAuthenticated,
    adminId,
    adminRole,
    // Actions
    login,
    logout,
    autoLogin,
  }
})
```

### 管理员类型定义

**类型定义:**

```typescript
// src/types/models.ts (添加)
export interface AdminUser {
  id: string
  username: string
  name: string
  email: string
  role: 'admin' | 'super_admin'
  createdAt: string
  updatedAt: string
}

export interface AdminLoginRequest {
  username: string
  password: string
}

export interface AdminLoginResponse {
  token: string
  admin: AdminUser
  expiresAt: number
}
```

### 管理员登录页面设计

**页面实现:**

```vue
<!-- src/pages/admin/login/index.vue -->
<template>
  <div class="admin-login-page">
    <n-card class="login-card" title="管理员登录">
      <n-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-placement="left"
        label-width="auto"
      >
        <n-form-item path="username" label="用户名">
          <n-input
            v-model:value="formData.username"
            placeholder="请输入用户名"
            :disabled="isLoading"
            @keyup.enter="handleLogin"
          />
        </n-form-item>

        <n-form-item path="password" label="密码">
          <n-input
            v-model:value="formData.password"
            type="password"
            placeholder="请输入密码"
            show-password-on="click"
            :disabled="isLoading"
            @keyup.enter="handleLogin"
          />
        </n-form-item>

        <n-form-item path="rememberMe">
          <n-checkbox v-model:checked="formData.rememberMe">
            记住我
          </n-checkbox>
        </n-form-item>

        <n-form-item>
          <n-button
            type="primary"
            block
            :loading="isLoading"
            @click="handleLogin"
          >
            登录
          </n-button>
        </n-form-item>
      </n-form>

      <n-alert v-if="errorMessage" type="error" :title="errorMessage" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useMessage } from 'naive-ui'
import { useAdminAuthStore } from '@/stores/admin-auth'
import type { FormInst, FormRules } from 'naive-ui'

const router = useRouter()
const route = useRoute()
const message = useMessage()
const adminAuthStore = useAdminAuthStore()

const formRef = ref<FormInst | null>(null)
const errorMessage = ref('')

const formData = reactive({
  username: '',
  password: '',
  rememberMe: false,
})

const rules: FormRules = {
  username: [
    {
      required: true,
      message: '请输入用户名',
      trigger: ['blur', 'input'],
    },
    {
      min: 3,
      max: 50,
      message: '用户名长度应为 3-50 个字符',
      trigger: ['blur', 'input'],
    },
  ],
  password: [
    {
      required: true,
      message: '请输入密码',
      trigger: ['blur', 'input'],
    },
    {
      min: 6,
      message: '密码长度至少 6 个字符',
      trigger: ['blur', 'input'],
    },
  ],
}

const isLoading = ref(false)

async function handleLogin() {
  errorMessage.value = ''

  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  isLoading.value = true

  const result = await adminAuthStore.login(
    formData.username,
    formData.password,
    formData.rememberMe
  )

  isLoading.value = false

  if (result.success) {
    message.success('登录成功')
    const redirect = (route.query.redirect as string) || '/admin'
    router.push(redirect)
  } else {
    errorMessage.value = result.message || '登录失败'
  }
}
</script>

<style scoped>
.admin-login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
}

.login-card {
  width: 400px;
  max-width: 90%;
}

.login-card :deep(.n-card-header) {
  text-align: center;
  font-size: 20px;
  font-weight: 600;
}
</style>
```

### 管理员 Token 工具

**工具实现:**

```typescript
// src/utils/admin-token.ts
const ADMIN_TOKEN_KEY = 'openclaw_admin_token'
const ADMIN_TOKEN_EXPIRY_KEY = 'openclaw_admin_token_expiry'
const ADMIN_REMEMBER_KEY = 'openclaw_admin_remember'

export function saveAdminToken(token: string, expiresIn: number): void {
  localStorage.setItem(ADMIN_TOKEN_KEY, token)
  const expiresAt = Date.now() + expiresIn * 1000
  localStorage.setItem(ADMIN_TOKEN_EXPIRY_KEY, expiresAt.toString())
}

export function getAdminToken(): string | null {
  const token = localStorage.getItem(ADMIN_TOKEN_KEY)
  const expiry = localStorage.getItem(ADMIN_TOKEN_EXPIRY_KEY)

  if (!token || !expiry) return null

  // 检查是否过期
  if (Date.now() > parseInt(expiry)) {
    removeAdminToken()
    return null
  }

  return token
}

export function removeAdminToken(): void {
  localStorage.removeItem(ADMIN_TOKEN_KEY)
  localStorage.removeItem(ADMIN_TOKEN_EXPIRY_KEY)
  localStorage.removeItem(ADMIN_REMEMBER_KEY)
}

export function setRememberMe(value: boolean): void {
  localStorage.setItem(ADMIN_REMEMBER_KEY, value.toString())
}

export function getRememberMe(): boolean {
  return localStorage.getItem(ADMIN_REMEMBER_KEY) === 'true'
}
```

### 管理员 API 服务

**服务实现:**

```typescript
// src/services/admin-auth.ts
import api from './api'
import type { AdminLoginRequest, AdminLoginResponse } from '@/types/models'

export async function adminLogin(data: AdminLoginRequest): Promise<AdminLoginResponse> {
  return api.post('/admin/auth/login', data)
}

export async function adminLogout(): Promise<void> {
  return api.post('/admin/auth/logout')
}
```

### 管理员路由守卫

**守卫实现:**

```typescript
// src/router/guards/admin-auth.ts
import type { Router } from 'vue-router'
import { useAdminAuthStore } from '@/stores/admin-auth'
import { getAdminToken } from '@/utils/admin-token'

export function setupAdminAuthGuard(router: Router) {
  router.beforeEach((to, from, next) => {
    const isAdminRoute = to.path.startsWith('/admin')
    const isAdminLogin = to.path === '/admin/login'
    const token = getAdminToken()

    // 管理员路由需要认证
    if (isAdminRoute && !isAdminLogin && !token) {
      next({
        path: '/admin/login',
        query: { redirect: to.fullPath },
      })
      return
    }

    // 已登录管理员访问登录页
    if (isAdminLogin && token) {
      next({ path: '/admin' })
      return
    }

    next()
  })
}
```

**路由配置:**

```typescript
// src/router/index.ts (添加)
import { setupAdminAuthGuard } from './guards/admin-auth'

const routes = [
  // ... 现有路由
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('@/pages/admin/login/index.vue'),
    meta: { requiresAuth: false, isAdmin: true },
  },
  {
    path: '/admin',
    name: 'AdminDashboard',
    component: () => import('@/pages/admin/dashboard/index.vue'),
    meta: { requiresAuth: true, isAdmin: true },
  },
]

setupAdminAuthGuard(router)
```

### 项目结构规范

**新增文件位置:**

```
frontend/
├── src/
│   ├── stores/
│   │   ├── auth.ts              # 租户认证 Store
│   │   └── admin-auth.ts        # 管理员认证 Store（新增）
│   ├── pages/
│   │   └── admin/
│   │       └── login/
│   │           └── index.vue    # 管理员登录页面（新增）
│   ├── router/
│   │   └── guards/
│   │       ├── auth.ts          # 租户认证守卫
│   │       └── admin-auth.ts    # 管理员认证守卫（新增）
│   ├── services/
│   │   ├── api.ts               # API 客户端
│   │   ├── auth.ts              # 租户认证服务
│   │   └── admin-auth.ts        # 管理员认证服务（新增）
│   ├── utils/
│   │   ├── token.ts             # 租户 Token 工具
│   │   └── admin-token.ts       # 管理员 Token 工具（新增）
│   └── types/
│       └── models.ts            # 类型定义（修改）
```

### 测试标准

**测试要求:**
- 测试框架: Vitest
- 测试覆盖率目标: >= 70%
- 组件测试: Vue Test Utils

**测试用例设计:**

| 测试场景 | 测试文件 | 测试方法 |
|---------|---------|---------|
| Store 初始状态 | `admin-auth.spec.ts` | 单元测试 |
| 登录成功 | `admin-auth.spec.ts` | 单元测试 |
| 登录失败 | `admin-auth.spec.ts` | 单元测试 |
| 记住登录 | `admin-auth.spec.ts` | 单元测试 |
| 自动登录 | `admin-auth.spec.ts` | 单元测试 |
| 表单验证 | `admin-login.spec.ts` | 组件测试 |
| 登录按钮交互 | `admin-login.spec.ts` | 组件测试 |
| 路由守卫 | `admin-auth-guard.spec.ts` | 单元测试 |

### UX 设计规范

**遵循 UX 设计规范 [Source: ux-design-specification.md]:**

1. **表单样式**:
   - 输入框最小高度: 36px
   - 标签位置: 左侧
   - 错误提示: 红色文字

2. **按钮样式**:
   - 主要按钮: 蓝色填充 #1677FF
   - 加载状态: 显示 loading 动画

3. **反馈机制**:
   - 成功: Message.success 3秒消失
   - 错误: Alert 显示，需手动关闭

4. **无障碍**:
   - 标签关联输入框
   - 键盘导航支持
   - Enter 键提交

### 安全注意事项

1. **密码输入**: 使用密码类型输入框，支持显示/隐藏切换
2. **Token 存储**: 根据记住登录设置不同的有效期
3. **自动登出**: Token 过期自动登出
4. **XSS 防护**: 不在 URL 中传递敏感信息

### Project Structure Notes

**与前置 Story 的连续性:**

1. **依赖 Story 2.4**:
   - 使用后端管理员认证 API
   - 使用管理员 JWT Token 结构

2. **复用 Story 2.5 的基础设施**:
   - API 客户端
   - 路由守卫模式
   - Store 设计模式

3. **与租户认证分离**:
   - 独立的管理员 Store
   - 独立的管理员 Token 存储
   - 独立的管理员路由守卫

### References

- [Source: architecture.md#Frontend Architecture] - Pinia 状态管理
- [Source: architecture.md#Naming Patterns] - 命名约定
- [Source: ux-design-specification.md] - UX 设计规范
- [Source: epics.md#Story 2.6] - 原始故事定义
- [Source: 2-4-admin-auth-system.md] - 管理员认证系统
- [Source: 2-5-frontend-login-page.md] - 前端登录页面

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
