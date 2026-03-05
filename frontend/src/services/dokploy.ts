import { api } from './api'
import type { ApiResponse } from '@/types'
import type { Instance } from '@/types/models'

/**
 * Dokploy 容器管理服务
 * 封装 Dokploy API 调用
 */

export interface DokployInstanceConfig {
  name: string
  envVars?: Record<string, string>
  resources?: {
    cpu?: number
    memory?: number
  }
}

export interface DokployInstanceStatus {
  id: string
  status: 'running' | 'stopped' | 'pending' | 'error'
  url?: string
  createdAt: string
  updatedAt: string
}

export const dokployService = {
  /**
   * 创建 OpenClaw 实例
   */
  async createInstance(
    tenantId: string,
    config: DokployInstanceConfig
  ): Promise<ApiResponse<Instance>> {
    return api.post(`/v1/tenants/${tenantId}/instances`, config)
  },

  /**
   * 获取实例状态
   */
  async getInstanceStatus(
    tenantId: string,
    instanceId: string
  ): Promise<ApiResponse<DokployInstanceStatus>> {
    return api.get(`/v1/tenants/${tenantId}/instances/${instanceId}/status`)
  },

  /**
   * 启动实例
   */
  async startInstance(
    tenantId: string,
    instanceId: string
  ): Promise<ApiResponse<void>> {
    return api.post(`/v1/tenants/${tenantId}/instances/${instanceId}/start`)
  },

  /**
   * 停止实例
   */
  async stopInstance(
    tenantId: string,
    instanceId: string
  ): Promise<ApiResponse<void>> {
    return api.post(`/v1/tenants/${tenantId}/instances/${instanceId}/stop`)
  },

  /**
   * 重启实例
   */
  async restartInstance(
    tenantId: string,
    instanceId: string
  ): Promise<ApiResponse<void>> {
    return api.post(`/v1/tenants/${tenantId}/instances/${instanceId}/restart`)
  },

  /**
   * 删除实例
   */
  async deleteInstance(
    tenantId: string,
    instanceId: string
  ): Promise<ApiResponse<void>> {
    return api.delete(`/v1/tenants/${tenantId}/instances/${instanceId}`)
  },

  /**
   * 获取实例日志
   */
  async getInstanceLogs(
    tenantId: string,
    instanceId: string,
    lines: number = 100
  ): Promise<ApiResponse<string[]>> {
    return api.get(`/v1/tenants/${tenantId}/instances/${instanceId}/logs?lines=${lines}`)
  },
}
