import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/services/api'

/**
 * 认证相关组合式 API
 */
export function useAuth() {
  const authStore = useAuthStore()
  const router = useRouter()

  const isAuthenticated = computed(() => authStore.isAuthenticated)
  const userId = computed(() => authStore.userId)
  const tenantId = computed(() => authStore.tenantId)

  /**
   * 用户登录
   */
  async function login(email: string, password: string): Promise<boolean> {
    const response = await api.post<{ token: string; userId: string; tenantId: string }>(
      '/v1/auth/login',
      { email, password }
    )

    if (response.data) {
      authStore.setToken(response.data.token)
      authStore.setUserInfo(response.data.userId, response.data.tenantId)
      return true
    }

    return false
  }

  /**
   * 用户登出
   */
  async function logout(): Promise<void> {
    // 可选：调用后端登出 API
    // await api.post('/v1/auth/logout')

    authStore.logout()
    router.push('/login')
  }

  /**
   * 验证当前 token 是否有效
   */
  async function verifyToken(): Promise<boolean> {
    if (!authStore.token) return false

    const response = await api.get<{ valid: boolean }>('/v1/auth/verify')
    return response.data?.valid ?? false
  }

  return {
    isAuthenticated,
    userId,
    tenantId,
    login,
    logout,
    verifyToken,
  }
}
