import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/pages/dashboard/index.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/login/index.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/pages/dashboard/index.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/feishu-config',
    name: 'FeishuConfig',
    component: () => import('@/pages/feishu-config/index.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/instances',
    name: 'Instances',
    component: () => import('@/pages/instances/index.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/feishu-bot',
    name: 'FeishuBot',
    component: () => import('@/pages/feishu-bot/index.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/backup',
    name: 'Backup',
    component: () => import('@/pages/backup/index.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/monitoring',
    name: 'Monitoring',
    component: () => import('@/pages/monitoring/index.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/onboarding',
    name: 'Onboarding',
    component: () => import('@/pages/onboarding/index.vue'),
    meta: { requiresAuth: true },
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

// 路由守卫 - 认证保护
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.meta.requiresAuth !== false // 默认需要认证

  if (requiresAuth && !authStore.isAuthenticated) {
    // 需要认证但未登录，重定向到登录页
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (to.name === 'Login' && authStore.isAuthenticated) {
    // 已登录用户访问登录页，重定向到首页
    next({ name: 'Dashboard' })
  } else {
    next()
  }
})

export default router
