import { describe, it, expect } from 'vitest'
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

// Naive UI 组件 stubs
const naiveUIStubs = {
  NConfigProvider: { template: '<div><slot /></div>' },
  NMessageProvider: { template: '<div><slot /></div>' },
  NLayout: { template: '<div class="app-layout"><slot /></div>' },
  NLayoutSider: { template: '<div><slot /></div>' },
  NLayoutHeader: { template: '<div><slot /></div>' },
  NLayoutContent: { template: '<div><slot /></div>' },
  NMenu: { template: '<div></div>' },
  NButton: { template: '<button><slot /></button>' },
}

describe('App.vue', () => {
  it('should mount successfully', async () => {
    setActivePinia(createPinia())
    const router = createTestRouter()
    router.push('/')
    await router.isReady()

    const wrapper = mount(App, {
      global: {
        plugins: [router],
        stubs: naiveUIStubs,
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
        stubs: naiveUIStubs,
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
        stubs: naiveUIStubs,
      },
    })

    // 登录页不显示 layout
    expect(wrapper.find('.app-layout').exists()).toBe(false)
  })
})
