import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { nextTick } from 'vue'
import App from './App.vue'

// 创建测试用路由
const createTestRouter = () =>
  createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', redirect: '/dashboard' },
      { path: '/login', name: 'login', component: { template: '<div>Login</div>' } },
      { path: '/dashboard', name: 'dashboard', component: { template: '<div>Dashboard</div>' } },
    ],
  })

// Mock Naive UI components
vi.stubGlobal('n-config-provider', { template: '<div><slot /></div>' })
vi.stubGlobal('n-message-provider', { template: '<div><slot /></div>' })
vi.stubGlobal('n-layout', { template: '<div class="app-layout"><slot /></div>' })
vi.stubGlobal('n-layout-sider', { template: '<div><slot /></div>' })
vi.stubGlobal('n-layout-header', { template: '<div><slot /></div>' })
vi.stubGlobal('n-layout-content', { template: '<div><slot /></div>' })
vi.stubGlobal('n-menu', { template: '<div></div>' })
vi.stubGlobal('n-button', { template: '<button><slot /></button>' })

describe('App.vue', () => {
  it('should mount successfully', async () => {
    setActivePinia(createPinia())
    const router = createTestRouter()
    router.push('/')
    await router.isReady()

    const wrapper = mount(App, {
      global: {
        plugins: [router],
      },
    })

    expect(wrapper.exists()).toBe(true)
  })

  it('should show layout for non-login pages', async () => {
    setActivePinia(createPinia())
    const router = createTestRouter()
    router.push('/dashboard')
    await router.isReady()
    await nextTick()

    const wrapper = mount(App, {
      global: {
        plugins: [router],
      },
    })

    expect(wrapper.find('.app-layout').exists()).toBe(true)
  })

  it('should not show layout for login page', async () => {
    setActivePinia(createPinia())
    const router = createTestRouter()
    router.push('/login')
    await router.isReady()
    await nextTick()

    const wrapper = mount(App, {
      global: {
        plugins: [router],
      },
    })

    // 登录页不显示 layout
    expect(wrapper.find('.app-layout').exists()).toBe(false)
  })
})
