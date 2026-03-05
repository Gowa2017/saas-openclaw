import { computed } from 'vue'
import { useInstanceStore } from '@/stores/instance'
import { api } from '@/services/api'
import { dokployService } from '@/services/dokploy'
import type { Instance } from '@/types/models'

/**
 * 实例管理组合式 API
 */
export function useInstance() {
  const instanceStore = useInstanceStore()

  const instances = computed(() => instanceStore.instances)
  const currentInstance = computed(() => instanceStore.currentInstance)

  /**
   * 获取实例列表
   */
  async function fetchInstances(tenantId: string): Promise<Instance[]> {
    const response = await api.get<Instance[]>(`/v1/tenants/${tenantId}/instances`)

    if (response.data) {
      instanceStore.setInstances(response.data)
      return response.data
    }

    return []
  }

  /**
   * 创建新实例
   */
  async function createInstance(
    tenantId: string,
    name: string,
    envVars?: Record<string, string>
  ): Promise<Instance | null> {
    const response = await dokployService.createInstance(tenantId, { name, envVars })

    if (response.data) {
      instanceStore.setCurrentInstance(response.data)
      return response.data
    }

    return null
  }

  /**
   * 启动实例
   */
  async function startInstance(tenantId: string, instanceId: string): Promise<boolean> {
    const response = await dokployService.startInstance(tenantId, instanceId)
    return response.error === null
  }

  /**
   * 停止实例
   */
  async function stopInstance(tenantId: string, instanceId: string): Promise<boolean> {
    const response = await dokployService.stopInstance(tenantId, instanceId)
    return response.error === null
  }

  /**
   * 重启实例
   */
  async function restartInstance(tenantId: string, instanceId: string): Promise<boolean> {
    const response = await dokployService.restartInstance(tenantId, instanceId)
    return response.error === null
  }

  /**
   * 删除实例
   */
  async function deleteInstance(tenantId: string, instanceId: string): Promise<boolean> {
    const response = await dokployService.deleteInstance(tenantId, instanceId)

    if (response.error === null) {
      instanceStore.setCurrentInstance(null)
      return true
    }

    return false
  }

  /**
   * 设置当前实例
   */
  function setCurrentInstance(instance: Instance | null): void {
    instanceStore.setCurrentInstance(instance)
  }

  return {
    instances,
    currentInstance,
    fetchInstances,
    createInstance,
    startInstance,
    stopInstance,
    restartInstance,
    deleteInstance,
    setCurrentInstance,
  }
}
