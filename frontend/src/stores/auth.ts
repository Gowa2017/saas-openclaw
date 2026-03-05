import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

// Token 存储 Key
const TOKEN_KEY = 'auth_token'
const USER_ID_KEY = 'user_id'
const TENANT_ID_KEY = 'tenant_id'

export const useAuthStore = defineStore('auth', () => {
  // 从 localStorage 初始化状态（持久化）
  const token = ref<string>(localStorage.getItem(TOKEN_KEY) || '')
  const userId = ref<string>(localStorage.getItem(USER_ID_KEY) || '')
  const tenantId = ref<string>(localStorage.getItem(TENANT_ID_KEY) || '')

  const isAuthenticated = computed(() => !!token.value)

  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem(TOKEN_KEY, newToken)
  }

  function setUserInfo(newUserId: string, newTenantId: string) {
    userId.value = newUserId
    tenantId.value = newTenantId
    localStorage.setItem(USER_ID_KEY, newUserId)
    localStorage.setItem(TENANT_ID_KEY, newTenantId)
  }

  function logout() {
    token.value = ''
    userId.value = ''
    tenantId.value = ''
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_ID_KEY)
    localStorage.removeItem(TENANT_ID_KEY)
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
