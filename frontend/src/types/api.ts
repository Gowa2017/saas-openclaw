// API 响应格式
export interface ApiResponse<T> {
  data: T | null
  error: ApiError | null
  meta?: {
    total?: number
    page?: number
    pageSize?: number
  }
}

export interface ApiError {
  code: string
  message: string
  details?: Record<string, unknown>
}
