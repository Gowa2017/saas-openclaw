import { api } from './api'
import type { ApiResponse } from '@/types'
import type { FeishuConfig } from '@/types/models'

/**
 * 飞书集成服务
 * 封装飞书开放平台 API 调用
 */

export interface FeishuAuthResult {
  accessToken: string
  refreshToken: string
  expiresIn: number
}

export interface FeishuUserInfo {
  openId: string
  unionId: string
  name: string
  avatarUrl?: string
  email?: string
}

export interface FeishuBotConfig {
  name: string
  description?: string
  avatarUrl?: string
}

export const feishuService = {
  /**
   * 验证飞书应用配置
   */
  async verifyConfig(config: FeishuConfig): Promise<ApiResponse<{ valid: boolean }>> {
    return api.post('/v1/feishu/verify', config)
  },

  /**
   * 获取飞书 OAuth 授权 URL
   */
  async getAuthUrl(redirectUri: string): Promise<ApiResponse<{ url: string }>> {
    return api.get(`/v1/feishu/auth/url?redirect_uri=${encodeURIComponent(redirectUri)}`)
  },

  /**
   * 处理飞书 OAuth 回调
   */
  async handleAuthCallback(code: string): Promise<ApiResponse<FeishuAuthResult>> {
    return api.post('/v1/feishu/auth/callback', { code })
  },

  /**
   * 获取飞书用户信息
   */
  async getUserInfo(): Promise<ApiResponse<FeishuUserInfo>> {
    return api.get('/v1/feishu/user/info')
  },

  /**
   * 配置飞书机器人
   */
  async configureBot(
    tenantId: string,
    config: FeishuBotConfig
  ): Promise<ApiResponse<void>> {
    return api.post(`/v1/tenants/${tenantId}/feishu/bot`, config)
  },

  /**
   * 获取飞书机器人状态
   */
  async getBotStatus(tenantId: string): Promise<ApiResponse<{ online: boolean }>> {
    return api.get(`/v1/tenants/${tenantId}/feishu/bot/status`)
  },

  /**
   * 发送飞书消息
   */
  async sendMessage(
    tenantId: string,
    openId: string,
    content: string
  ): Promise<ApiResponse<{ messageId: string }>> {
    return api.post(`/v1/tenants/${tenantId}/feishu/messages`, {
      openId,
      content,
    })
  },
}
