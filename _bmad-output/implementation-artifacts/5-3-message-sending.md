# Story 5.3: 消息发送功能

Status: ready-for-dev

## Story

As a 后端开发者,
I want 实现飞书消息发送功能,
so that 可以将 OpenClaw 的回复发送给用户。

## Acceptance Criteria

1. **AC1: 文本消息发送**
   - **Given** OpenClaw 实例生成了回复
   - **When** 调用消息发送接口发送文本消息
   - **Then** 使用飞书 API 发送文本消息给用户
   - **And** 消息格式符合飞书规范
   - **And** 消息成功送达

2. **AC2: 卡片消息发送**
   - **Given** 需要发送结构化回复
   - **When** 调用消息发送接口发送卡片消息
   - **Then** 使用飞书 API 发送卡片消息
   - **And** 卡片模板正确渲染
   - **And** 卡片交互元素正常工作

3. **AC3: 消息发送重试**
   - **Given** 消息发送失败
   - **When** 检测到发送失败
   - **Then** 自动重试发送
   - **And** 最多重试 3 次
   - **And** 重试间隔采用指数退避（1s, 2s, 4s）
   - **And** 记录重试日志

4. **AC4: 发送状态记录**
   - **Given** 消息发送完成（成功或失败）
   - **When** 记录发送状态
   - **Then** 记录消息 ID、发送时间、发送状态
   - **And** 记录失败原因（如果失败）
   - **And** 日志持久化存储

5. **AC5: 多租户支持**
   - **Given** 系统有多个租户
   - **When** 发送消息
   - **Then** 使用正确租户的飞书凭证
   - **And** 使用租户对应的 App ID 和 App Secret
   - **And** 凭证从加密存储中获取

6. **AC6: API 封装**
   - **Given** 消息发送功能实现
   - **When** 其他模块调用
   - **Then** 提供简洁的 API 接口
   - **And** 接口支持发送文本、卡片消息
   - **And** 接口返回发送结果

## Tasks / Subtasks

- [ ] Task 1: 创建消息发送服务 (AC: 1, 2)
  - [ ] 1.1 创建 `internal/infrastructure/feishu/message/sender.go`
  - [ ] 1.2 实现文本消息发送方法
  - [ ] 1.3 实现卡片消息发送方法
  - [ ] 1.4 实现消息内容构建器

- [ ] Task 2: 实现飞书 API 调用 (AC: 1, 2)
  - [ ] 2.1 封装飞书消息发送 API
  - [ ] 2.2 实现 tenant_access_token 获取
  - [ ] 2.3 实现请求签名
  - [ ] 2.4 实现响应解析

- [ ] Task 3: 实现重试机制 (AC: 3)
  - [ ] 3.1 创建 `pkg/retry/retry.go` 重试工具
  - [ ] 3.2 实现指数退避算法
  - [ ] 3.3 实现最大重试次数控制（3 次）
  - [ ] 3.4 实现可重试错误判断

- [ ] Task 4: 实现发送状态记录 (AC: 4)
  - [ ] 4.1 创建 `internal/domain/feishu/message_log.go` 消息日志实体
  - [ ] 4.2 创建消息日志存储接口
  - [ ] 4.3 实现日志持久化（数据库）
  - [ ] 4.4 实现日志查询功能

- [ ] Task 5: 实现多租户凭证管理 (AC: 5)
  - [ ] 5.1 创建 `internal/infrastructure/feishu/credential.go`
  - [ ] 5.2 实现凭证获取和缓存
  - [ ] 5.3 实现凭证解密（App Secret）
  - [ ] 5.4 实现 access_token 缓存和刷新

- [ ] Task 6: 实现消息 API 接口 (AC: 6)
  - [ ] 6.1 创建 `internal/domain/feishu/message_service.go`
  - [ ] 6.2 定义消息服务接口
  - [ ] 6.3 实现 SendTextMessage 方法
  - [ ] 6.4 实现 SendCardMessage 方法
  - [ ] 6.5 实现发送结果返回

- [ ] Task 7: 单元测试 (AC: All)
  - [ ] 7.1 编写消息发送服务测试
  - [ ] 7.2 编写重试机制测试
  - [ ] 7.3 编写凭证管理测试
  - [ ] 7.4 编写消息 API 接口测试
  - [ ] 7.5 Mock 飞书 API 响应
  - [ ] 7.6 确保测试覆盖率 ≥ 70%

- [ ] Task 8: 集成测试 (AC: All)
  - [ ] 8.1 编写端到端消息发送测试
  - [ ] 8.2 编写重试场景测试
  - [ ] 8.3 编写多租户场景测试

## Dev Notes

### 消息发送 API

**飞书消息发送接口:**

```
POST https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=chat_id
Authorization: Bearer {tenant_access_token}
Content-Type: application/json

{
  "receive_id": "oc_xxx",
  "msg_type": "text",
  "content": "{\"text\":\"hello world\"}"
}
```

**消息类型:**
- `text`: 文本消息
- `post`: 富文本消息
- `interactive`: 卡片消息
- `image`: 图片消息
- `file`: 文件消息

### 卡片消息示例

```json
{
  "type": "template",
  "data": {
    "template_id": "AAqk*****",
    "template_variable": {
      "title": "OpenClaw 回复",
      "content": "这是回复内容..."
    }
  }
}
```

### 架构设计

**目录结构:**

```
internal/infrastructure/feishu/
├── message/
│   ├── sender.go           # 消息发送器
│   ├── builder.go          # 消息构建器
│   └── card.go             # 卡片模板
├── credential.go           # 凭证管理
└── token.go                # Token 管理

internal/domain/feishu/
├── message_service.go      # 消息服务接口
└── message_log.go          # 消息日志实体

pkg/retry/
└── retry.go                # 重试工具
```

### 关键技术点

**Access Token 获取:**

```
POST https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal

{
  "app_id": "cli_xxx",
  "app_secret": "xxx"
}
```

- Token 有效期 2 小时
- 建议缓存 Token，过期前刷新
- 每个 App ID 独立的 Token

**重试策略:**
- 可重试错误: 网络超时、5xx 错误、限流
- 不可重试错误: 4xx 错误（参数错误、权限不足等）
- 重试间隔: 1s, 2s, 4s（指数退避）
- 最大重试次数: 3 次

**凭证加密存储:**
- App Secret 使用 AES-256 加密存储
- 加密密钥从环境变量获取
- 解密后在内存中使用，不持久化

### 性能要求

**来自 NFR:**
- NFR-P2: 用户在飞书发送消息后，在 3 秒内收到 OpenClaw 回复
- NFR-I2: 飞书消息投递成功率 ≥ 99%

### 错误码处理

| 错误码 | 说明 | 处理方式 |
|-------|------|---------|
| 400 | 参数错误 | 不重试，记录日志 |
| 401 | Token 过期 | 刷新 Token 后重试 |
| 403 | 权限不足 | 不重试，通知管理员 |
| 429 | 限流 | 等待后重试 |
| 500+ | 服务端错误 | 重试 |

### References

- [Source: architecture.md#Integration Layer] - 飞书消息发送设计
- [Source: prd.md#FR22] - 机器人发送回复需求
- [Source: prd.md#NFR-I2] - 消息投递成功率要求
- [Source: epics.md#Story 5.3] - 原始故事定义
- 飞书消息发送文档: https://open.feishu.cn/document/client-docs/bot-v3/messages/bot_send_messages

## Dev Agent Record

### Agent Model Used

(待填写)

### Debug Log References

(待填写)

### Completion Notes List

(待填写)

### File List

(待填写)

## Senior Developer Review (AI)

(待填写)
