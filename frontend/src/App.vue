<script setup lang="ts">
import { NConfigProvider, NLayout, NLayoutHeader, NLayoutSider, NLayoutContent, NMenu, NMessageProvider } from 'naive-ui'
import { useRouter } from 'vue-router'
import type { MenuOption } from 'naive-ui'

const router = useRouter()

const menuOptions: MenuOption[] = [
  {
    label: '仪表盘',
    key: 'dashboard',
    onClick: () => router.push('/dashboard'),
  },
  {
    label: '飞书配置',
    key: 'feishu-config',
    onClick: () => router.push('/feishu-config'),
  },
  {
    label: '实例管理',
    key: 'instances',
    onClick: () => router.push('/instances'),
  },
  {
    label: '飞书机器人',
    key: 'feishu-bot',
    onClick: () => router.push('/feishu-bot'),
  },
  {
    label: '备份恢复',
    key: 'backup',
    onClick: () => router.push('/backup'),
  },
  {
    label: '监控运维',
    key: 'monitoring',
    onClick: () => router.push('/monitoring'),
  },
  {
    label: '帮助中心',
    key: 'onboarding',
    onClick: () => router.push('/onboarding'),
  },
]
</script>

<template>
  <n-config-provider>
    <n-message-provider>
      <n-layout has-sider class="app-layout">
        <n-layout-sider
          bordered
          collapse-mode="width"
          :collapsed-width="64"
          :width="200"
          :native-scrollbar="false"
          show-trigger
        >
          <div class="logo">
            <h2>OpenClaw SaaS</h2>
          </div>
          <n-menu :options="menuOptions" />
        </n-layout-sider>
        <n-layout>
          <n-layout-header bordered class="app-header">
            <n-button text @click="router.push('/login')">退出登录</n-button>
          </n-layout-header>
          <n-layout-content class="app-content">
            <router-view />
          </n-layout-content>
        </n-layout>
      </n-layout>
    </n-message-provider>
  </n-config-provider>
</template>

<style scoped>
.app-layout {
  height: 100vh;
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 64px;
  border-bottom: 1px solid #e0e0e0;
}

.logo h2 {
  margin: 0;
  font-size: 16px;
  color: #1677FF;
}

.app-header {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 0 24px;
  height: 64px;
}

.app-content {
  min-height: calc(100vh - 64px);
  background-color: #f5f5f5;
}
</style>
