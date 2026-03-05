import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/pages/dashboard/index.vue'),
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/login/index.vue'),
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/pages/dashboard/index.vue'),
  },
  {
    path: '/feishu-config',
    name: 'FeishuConfig',
    component: () => import('@/pages/feishu-config/index.vue'),
  },
  {
    path: '/instances',
    name: 'Instances',
    component: () => import('@/pages/instances/index.vue'),
  },
  {
    path: '/feishu-bot',
    name: 'FeishuBot',
    component: () => import('@/pages/feishu-bot/index.vue'),
  },
  {
    path: '/backup',
    name: 'Backup',
    component: () => import('@/pages/backup/index.vue'),
  },
  {
    path: '/monitoring',
    name: 'Monitoring',
    component: () => import('@/pages/monitoring/index.vue'),
  },
  {
    path: '/onboarding',
    name: 'Onboarding',
    component: () => import('@/pages/onboarding/index.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router
