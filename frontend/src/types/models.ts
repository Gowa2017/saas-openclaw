// 用户模型
export interface User {
  id: string
  email: string
  name: string
  tenantId: string
  role: 'user' | 'admin'
  createdAt: string
  updatedAt: string
}

// 租户模型
export interface Tenant {
  id: string
  name: string
  feishuAppId?: string
  createdAt: string
  updatedAt: string
}

// OpenClaw 实例模型
export interface Instance {
  id: string
  tenantId: string
  name: string
  status: InstanceStatus
  url?: string
  dokployId?: string
  createdAt: string
  updatedAt: string
}

export type InstanceStatus = 'running' | 'stopped' | 'error' | 'pending'

// 飞书配置模型
export interface FeishuConfig {
  appId: string
  appSecret: string
  encryptKey?: string
  verificationToken?: string
}

// 备份配置模型
export interface BackupConfig {
  enabled: boolean
  schedule: string
  retention: number
  lastBackupAt?: string
}
