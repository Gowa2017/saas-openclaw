# Story 2.5: 前端登录页面与状态管理

Status: ready-for-dev

## Story

As a 前端开发者,
I want 实现登录页面和认证状态管理,
so that 用户可以登录系统并保持登录状态。

## Acceptance Criteria

1. **AC1: 路由守卫实现**
   - **Given** 前端项目已初始化
   - **When** 用户访问需要认证的页面
   - **Then** 未登录用户重定向到登录页

2. **AC2: 登录页面实现**
   - **Given** 用户访问登录页
   - **When** 页面加载完成
   - **Then** 显示业务平台登录入口
   - **And** 提供清晰的登录指引

3. **AC3: Token 存储管理**
   - **Given** 登录成功
   - **When** 获取到 Token
   - **Then** 将 Token 存储在 localStorage
   - **And** 支持 Token 自动附加到请求头

4. **AC4: Pinia 认证状态管理**
   - **Given** 使用 Pinia 管理认证状态
   - **When** 创建 auth store
   - **Then** 包含用户信息、Token、登录状态
   - **And** 提供登录、登出、刷新等方法

5. **AC5: Token 自动刷新机制**
   - **Given** Token 即将过期
   - **When** 检测到 Token 需要刷新
   - **Then** 自动执行 Token 刷新
   - **And** 无感知更新用户会话

## Tasks / Subtasks

- [ ] Task 1: 创建认证 Store (AC: 4)
  - [ ] 1.1 创建 `src/stores/auth.ts`
  - [ ] 1.2 定义用户状态接口
  - [ ] 1.3 实现 login action
  - [ ] 1.4 实现 logout action
  - [ ] 1.5 实现 token 刷新 action

- [ ] Task 2: 创建登录页面 (AC: 2)
  - [ ] 2.1 创建 `src/pages/login/index.vue`
  - [ ] 2.2 实现业务平台登录入口 UI
  - [ ] 2.3 添加登录指引说明
  - [ ] 2.4 实现加载状态展示

- [ ] Task 3: 实现路由守卫 (AC: 1)
  - [ ] 3.1 创建 `src/router/guards/auth.ts`
  - [ ] 3.2 实现认证检查逻辑
  - [ ] 3.3 实现重定向逻辑
  - [ ] 3.4 添加路由元信息配置

- [ ] Task 4: 实现 Token 管理 (AC: 3)
  - [ ] 4.1 创建 `src/utils/token.ts`
  - [ ] 4.2 实现 saveToken 方法
  - [ ] 4.3 实现 getToken 方法
  - [ ] 4.4 实现 removeToken 方法

- [ ] Task 5: 实现 API 拦截器 (AC: 3)
  - [ ] 5.1 创建 `src/services/api.ts`
  - [ ] 5.2 实现请求拦截器（自动附加 Token）
  - [ ] 5.3 实现响应拦截器（处理 401）
  - [ ] 5.4 实现 Token 刷新逻辑

- [ ] Task 6: 编写单元测试 (AC: 1-5)
  - [ ] 6.1 编写 `auth.spec.ts` 测试 Store
  - [ ] 6.2 编写 `auth-guard.spec.ts` 测试路由守卫
  - [ ] 6.3 编写 `token.spec.ts` 测试 Token 管理
  - [ ] 6.4 编写 `api-interceptor.spec.ts` 测试拦截器

## Dev Notes

### 架构模式与约束

**必须遵循的架构原则 [Source: architecture.md]:**

1. **技术栈**:
   - Vue 3 + TypeScript
   - Pinia 状态管理
   - Vue Router 路由
   - Naive UI 组件库

2. **命名约定 [Source: architecture.md#Naming Patterns]:**
   - 组件名: PascalCase (例: `LoginPage`)
   - 文件名: kebab-case (例: `login-page.vue`)
   - 函数名: camelCase (例: `getToken`, `setAuthHeader`)

3. **状态管理结构 [Source: architecture.md#Frontend Architecture]:**
   - 按模块组织: `auth`, `tenant`, `instance`, `config`

### 现有项目状态

**依赖 Story 1.2 前端项目初始化已完成:**

```
frontend/
├── src/
│   ├── components/           # 组件目录
│   ├── composables/          # 组合式 API
│   ├── pages/               # 页面路由
│   ├── stores/              # Pinia 状态管理（待实现）
│   ├── services/            # API 调用服务（待实现）
│   ├── types/               # TypeScript 类型定义
│   ├── utils/               # 工具函数（待实现）
│   ├── router/              # 路由配置
│   └── assets/              # 静态资源
```

### 认证 Store 设计

**Store 实现:**

```typescript
// src/stores/auth.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types/models'
import { getToken, saveToken, removeToken } from '@/utils/token'
import { login as loginApi, refreshToken as refreshTokenApi } from '@/services/auth'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const token = ref<string | null>(getToken())
  const isLoading = ref(false)

  // Getters
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const userId = computed(() => user.value?.id)
  const tenantId = computed(() => user.value?.tenantId)

  // Actions
  async function login(platformToken: string) {
    isLoading.value = true
    try {
      const response = await loginApi(platformToken)
      token.value = response.token
      user.value = response.user
      saveToken(response.token)
      return true
    } catch (error) {
      console.error('Login failed:', error)
      return false
    } finally {
      isLoading.value = false
    }
  }

  async function logout() {
    token.value = null
    user.value = null
    removeToken()
  }

  async function refreshToken() {
    if (!token.value) return false
    try {
      const response = await refreshTokenApi(token.value)
      token.value = response.token
      saveToken(response.token)
      return true
    } catch (error) {
      logout()
      return false
    }
  }

  function setUser(userData: User) {
    user.value = userData
  }

  return {
    // State
    user,
    token,
    isLoading,
    // Getters
    isAuthenticated,
    userId,
    tenantId,
    // Actions
    login,
    logout,
    refreshToken,
    setUser,
  }
})
```

### 用户类型定义

**类型定义:**

```typescript
// src/types/models.ts
export interface User {
  id: string
  tenantId: string
  name: string
  email: string
  role: 'user' | 'admin'
  createdAt: string
  updatedAt: string
}

export interface AuthResponse {
  token: string
  user: User
  expiresAt: number
}
```

### 登录页面设计

**页面实现:**

```vue
<!-- src/pages/login/index.vue -->
<template>
  <div class="login-page">
    <n-card class="login-card" title="登录">
      <div class="login-content">
        <n-space vertical align="center" :size="24">
          <n-icon size="64" color="#1677FF">
            <LogInOutline />
          </n-icon>

          <n-h3>欢迎使用 OpenClaw SaaS 平台</n-h3>

          <n-p depth="3">
            请通过业务平台登录，获取访问权限
          </n-p>

          <n-button
            type="primary"
            size="large"
            block
            :loading="authStore.isLoading"
            @click="handleLogin"
          >
            <template #icon>
              <n-icon><LogInOutline /></n-icon>
            </template>
            业务平台登录
          </n-button>

          <n-collapse>
            <n-collapse-item title="如何登录？" name="help">
              <n-list>
                <n-list-item>
                  1. 点击"业务平台登录"按钮
                </n-list-item>
                <n-list-item>
                  2. 跳转到业务平台完成认证
                </n-list-item>
                <n-list-item>
                  3. 自动返回本平台并完成登录
                </n-list-item>
              </n-list>
            </n-collapse-item>
          </n-collapse>
        </n-space>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { LogInOutline } from '@vicons/ionicons5'

const authStore = useAuthStore()
const route = useRoute()

const handleLogin = async () => {
  // 模拟从 URL 参数获取 platform token
  const platformToken = route.query.token as string
  if (platformToken) {
    const success = await authStore.login(platformToken)
    if (success) {
      // 跳转到首页或原目标页面
      const redirect = route.query.redirect as string || '/'
      navigateTo(redirect)
    }
  } else {
    // 跳转到业务平台登录页
    window.location.href = import.meta.env.VITE_PLATFORM_LOGIN_URL
  }
}
</script>

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
  max-width: 90%;
}

.login-content {
  padding: 24px 0;
}
</style>
```

### 路由守卫设计

**守卫实现:**

```typescript
// src/router/guards/auth.ts
import type { Router } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { getToken } from '@/utils/token'

export function setupAuthGuard(router: Router) {
  router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()
    const requiresAuth = to.meta.requiresAuth
    const token = getToken()

    // 页面需要认证
    if (requiresAuth && !token) {
      next({
        path: '/login',
        query: { redirect: to.fullPath },
      })
      return
    }

    // 已登录用户访问登录页
    if (to.path === '/login' && token) {
      next({ path: '/' })
      return
    }

    next()
  })
}
```

**路由配置:**

```typescript
// src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router'
import { setupAuthGuard } from './guards/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/login/index.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('@/pages/dashboard/index.vue'),
    meta: { requiresAuth: true },
  },
  // ... 其他路由
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

setupAuthGuard(router)

export default router
```

### Token 管理工具

**工具实现:**

```typescript
// src/utils/token.ts
const TOKEN_KEY = 'openclaw_token'
const TOKEN_EXPIRY_KEY = 'openclaw_token_expiry'

export function saveToken(token: string, expiresAt?: number): void {
  localStorage.setItem(TOKEN_KEY, token)
  if (expiresAt) {
    localStorage.setItem(TOKEN_EXPIRY_KEY, expiresAt.toString())
  }
}

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function removeToken(): void {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(TOKEN_EXPIRY_KEY)
}

export function isTokenExpired(): boolean {
  const expiry = localStorage.getItem(TOKEN_EXPIRY_KEY)
  if (!expiry) return false
  return Date.now() > parseInt(expiry) * 1000
}

export function isTokenExpiringSoon(thresholdMinutes: number = 5): boolean {
  const expiry = localStorage.getItem(TOKEN_EXPIRY_KEY)
  if (!expiry) return true
  const expiresAt = parseInt(expiry) * 1000
  const threshold = thresholdMinutes * 60 * 1000
  return Date.now() + threshold > expiresAt
}
```

### API 拦截器设计

**拦截器实现:**

```typescript
// src/services/api.ts
import axios from 'axios'
import { getToken, removeToken, isTokenExpiringSoon } from '@/utils/token'
import { useAuthStore } from '@/stores/auth'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 30000,
})

// 请求拦截器
api.interceptors.request.use(
  async (config) => {
    const token = getToken()
    if (token) {
      // 检查 Token 是否即将过期
      if (isTokenExpiringSoon()) {
        const authStore = useAuthStore()
        await authStore.refreshToken()
      }
      config.headers['X-Platform-Token'] = token
    }
    return config
  },
  (error) => Promise.reject(error)
)

// 响应拦截器
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      removeToken()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api
```

### 项目结构规范

**新增文件位置:**

```
frontend/
├── src/
│   ├── stores/
│   │   ├── auth.ts              # 认证状态（新增）
│   │   └── auth.spec.ts         # 测试文件（新增）
│   ├── pages/
│   │   └── login/
│   │       └── index.vue        # 登录页面（新增）
│   ├── router/
│   │   ├── index.ts             # 路由配置（修改）
│   │   └── guards/
│   │       ├── auth.ts          # 认证守卫（新增）
│   │       └── auth.spec.ts     # 测试文件（新增）
│   ├── services/
│   │   ├── api.ts               # API 客户端（新增）
│   │   └── auth.ts              # 认证服务（新增）
│   ├── utils/
│   │   ├── token.ts             # Token 工具（新增）
│   │   └── token.spec.ts        # 测试文件（新增）
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
| Store 初始状态 | `auth.spec.ts` | 单元测试 |
| 登录成功 | `auth.spec.ts` | 单元测试 |
| 登录失败 | `auth.spec.ts` | 单元测试 |
| 登出 | `auth.spec.ts` | 单元测试 |
| Token 刷新 | `auth.spec.ts` | 单元测试 |
| 路由守卫重定向 | `auth-guard.spec.ts` | 单元测试 |
| Token 存储 | `token.spec.ts` | 单元测试 |
| Token 过期检查 | `token.spec.ts` | 单元测试 |

### UX 设计规范

**遵循 UX 设计规范 [Source: ux-design-specification.md]:**

1. **按钮样式**:
   - 主要按钮: 蓝色填充 #1677FF
   - 按钮最小高度: 36px

2. **反馈机制**:
   - 成功反馈: Message.success 3秒自动消失
   - 错误反馈: Message.error 需手动关闭

3. **表单验证**:
   - 实时校验 + 提交校验

4. **无障碍**:
   - WCAG 2.1 AA 级合规
   - 色彩对比度 4.5:1
   - 键盘导航支持

### 配置项

**需要的环境变量:**

```env
# .env.example
VITE_API_BASE_URL=http://localhost:8080/v1
VITE_PLATFORM_LOGIN_URL=https://platform.example.com/login
```

### Project Structure Notes

**与前置 Story 的连续性:**

1. **依赖 Story 1.2**:
   - 使用 Vue 3 + TypeScript 项目结构
   - 使用 Naive UI 组件库
   - 使用 Pinia 状态管理

2. **为后续 Story 提供基础**:
   - 认证状态供所有需要认证的页面使用
   - 路由守卫保护所有需要认证的路由

### References

- [Source: architecture.md#Frontend Architecture] - Pinia 状态管理
- [Source: architecture.md#Naming Patterns] - 命名约定
- [Source: ux-design-specification.md] - UX 设计规范
- [Source: epics.md#Story 2.5] - 原始故事定义
- [Source: 1-2-frontend-project-init.md] - 前端项目上下文

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
