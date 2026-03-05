import type { ApiResponse } from '@/types'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'
const DEFAULT_TIMEOUT = 30000 // 30 秒超时

/**
 * 创建带超时的 AbortController
 */
function createTimeoutController(timeout: number): AbortController {
  const controller = new AbortController()
  setTimeout(() => controller.abort(), timeout)
  return controller
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {},
  timeout: number = DEFAULT_TIMEOUT
): Promise<ApiResponse<T>> {
  const url = `${API_BASE_URL}${endpoint}`

  const defaultHeaders: HeadersInit = {
    'Content-Type': 'application/json',
  }

  const token = localStorage.getItem('auth_token')
  if (token) {
    defaultHeaders['Authorization'] = `Bearer ${token}`
  }

  const controller = createTimeoutController(timeout)

  try {
    const response = await fetch(url, {
      ...options,
      signal: controller.signal,
      headers: {
        ...defaultHeaders,
        ...options.headers,
      },
    })

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      return {
        data: null,
        error: {
          code: String(response.status),
          message: errorData.message || response.statusText,
          details: errorData.details,
        },
      }
    }

    const data = await response.json()
    return {
      data,
      error: null,
    }
  } catch (error) {
    // 处理超时错误
    if (error instanceof Error && error.name === 'AbortError') {
      return {
        data: null,
        error: {
          code: 'TIMEOUT',
          message: '请求超时，请稍后重试',
        },
      }
    }

    return {
      data: null,
      error: {
        code: 'NETWORK_ERROR',
        message: error instanceof Error ? error.message : '网络错误',
      },
    }
  }
}

export const api = {
  get: <T>(endpoint: string, timeout?: number) =>
    request<T>(endpoint, { method: 'GET' }, timeout),
  post: <T>(endpoint: string, body?: unknown, timeout?: number) =>
    request<T>(
      endpoint,
      {
        method: 'POST',
        body: body ? JSON.stringify(body) : undefined,
      },
      timeout
    ),
  put: <T>(endpoint: string, body?: unknown, timeout?: number) =>
    request<T>(
      endpoint,
      {
        method: 'PUT',
        body: body ? JSON.stringify(body) : undefined,
      },
      timeout
    ),
  patch: <T>(endpoint: string, body?: unknown, timeout?: number) =>
    request<T>(
      endpoint,
      {
        method: 'PATCH',
        body: body ? JSON.stringify(body) : undefined,
      },
      timeout
    ),
  delete: <T>(endpoint: string, timeout?: number) =>
    request<T>(endpoint, { method: 'DELETE' }, timeout),
}
