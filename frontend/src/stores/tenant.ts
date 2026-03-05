import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Tenant {
  id: string
  name: string
  feishuAppId?: string
  createdAt: string
  updatedAt: string
}

export const useTenantStore = defineStore('tenant', () => {
  const currentTenant = ref<Tenant | null>(null)
  const tenants = ref<Tenant[]>([])

  function setCurrentTenant(tenant: Tenant | null) {
    currentTenant.value = tenant
  }

  function setTenants(tenantList: Tenant[]) {
    tenants.value = tenantList
  }

  return {
    currentTenant,
    tenants,
    setCurrentTenant,
    setTenants,
  }
})
