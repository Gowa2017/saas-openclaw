import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>('')
  const userId = ref<string>('')
  const tenantId = ref<string>('')

  const isAuthenticated = computed(() => !!token.value)

  function setToken(newToken: string) {
    token.value = newToken
  }

  function setUserInfo(newUserId: string, newTenantId: string) {
    userId.value = newUserId
    tenantId.value = newTenantId
  }

  function logout() {
    token.value = ''
    userId.value = ''
    tenantId.value = ''
  }

  return {
    token,
    userId,
    tenantId,
    isAuthenticated,
    setToken,
    setUserInfo,
    logout,
  }
})
