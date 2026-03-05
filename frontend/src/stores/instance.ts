import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Instance } from '@/types/models'

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
