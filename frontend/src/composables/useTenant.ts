import { computed } from 'vue'
import { useTenantStore } from '@/stores/tenant'
import { api } from '@/services/api'
import type { Tenant } from '@/types/models'

/**
 * 租户相关组合式 API
 */
export function useTenant() {
  const tenantStore = useTenantStore()

  const currentTenant = computed(() => tenantStore.currentTenant)
  const tenants = computed(() => tenantStore.tenants)

  /**
   * 获取当前租户信息
   */
  async function fetchCurrentTenant(): Promise<Tenant | null> {
    const response = await api.get<Tenant>('/v1/tenants/current')

    if (response.data) {
      tenantStore.setCurrentTenant(response.data)
      return response.data
    }

    return null
  }

  /**
   * 获取租户列表
   */
  async function fetchTenants(): Promise<Tenant[]> {
    const response = await api.get<Tenant[]>('/v1/tenants')

    if (response.data) {
      tenantStore.setTenants(response.data)
      return response.data
    }

    return []
  }

  /**
   * 更新租户信息
   */
  async function updateTenant(tenantId: string, data: Partial<Tenant>): Promise<boolean> {
    const response = await api.patch<Tenant>(`/v1/tenants/${tenantId}`, data)

    if (response.data) {
      tenantStore.setCurrentTenant(response.data)
      return true
    }

    return false
  }

  /**
   * 清除当前租户
   */
  function clearCurrentTenant(): void {
    tenantStore.setCurrentTenant(null)
  }

  return {
    currentTenant,
    tenants,
    fetchCurrentTenant,
    fetchTenants,
    updateTenant,
    clearCurrentTenant,
  }
}
