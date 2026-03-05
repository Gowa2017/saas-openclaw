import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { FeishuConfig } from '@/types/models'

export const useConfigStore = defineStore('config', () => {
  const feishuConfig = ref<FeishuConfig | null>(null)

  function setFeishuConfig(config: FeishuConfig) {
    feishuConfig.value = config
  }

  function clearFeishuConfig() {
    feishuConfig.value = null
  }

  return {
    feishuConfig,
    setFeishuConfig,
    clearFeishuConfig,
  }
})
