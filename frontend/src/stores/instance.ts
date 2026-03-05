import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Instance {
  id: string
  tenantId: string
  name: string
  status: 'running' | 'stopped' | 'error' | 'pending'
  url?: string
  createdAt: string
  updatedAt: string
}

export const useInstanceStore = defineStore('instance', () => {
  const instances = ref<Instance[]>([])
  const currentInstance = ref<Instance | null>(null)

  function setInstances(instanceList: Instance[]) {
    instances.value = instanceList
  }

  function setCurrentInstance(instance: Instance | null) {
    currentInstance.value = instance
  }

  return {
    instances,
    currentInstance,
    setInstances,
    setCurrentInstance,
  }
})
