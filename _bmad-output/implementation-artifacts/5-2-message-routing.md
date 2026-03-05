# Story 5.2: 消息接收与路由

Status: ready-for-dev

## Story

As a 后端开发者,
I want 接收飞书消息并路由到正确的 OpenClaw 实例,
so that 用户的消息能被正确处理。

## Acceptance Criteria

1. **AC1: 消息事件接收**
   - **Given** 飞书长连接已建立
   - **When** 用户在飞书发送消息给机器人
   - **Then** 系统接收消息事件
   - **And** 事件类型为 `im.message.receive_v1`
   - **And** 消息内容正确解析

2. **AC2: 租户识别**
   - **Given** 收到消息事件
   - **When** 解析事件中的 App ID
   - **Then** 根据 App ID 查找对应的租户配置
   - **And** 租户不存在时记录错误日志
   - **And** 返回适当的错误响应

3. **AC3: OpenClaw 实例查找**
   - **Given** 已识别租户
   - **When** 查询租户的 OpenClaw 实例
   - **Then** 返回实例的连接地址
   - **And** 实例状态为 running
   - **And** 实例不存在时返回友好提示

4. **AC4: 消息转发**
   - **Given** 找到 OpenClaw 实例
   - **When** 将消息转发到实例
   - **Then** 调用实例的消息处理 API
   - **And** 转发超时设置为 30 秒
   - **And** 记录转发日志

5. **AC5: 错误处理**
   - **Given** 消息处理过程中发生错误
   - **When** 出现以下情况之一：
     - 租户未配置飞书应用
     - OpenClaw 实例不存在或未运行
     - 消息转发失败
   - **Then** 记录详细错误日志
   - **And** 向用户发送友好错误提示
   - **And** 错误信息不暴露系统内部细节

6. **AC6: 消息处理日志**
   - **Given** 消息处理完成
   - **When** 查看日志
   - **Then** 包含消息 ID、租户 ID、处理时间、处理结果
   - **And** 日志级别正确（INFO/WARN/ERROR）
   - **And** 敏感信息已脱敏

## Tasks / Subtasks

- [ ] Task 1: 创建消息领域模型 (AC: 1)
  - [ ] 1.1 创建 `internal/domain/feishu/message.go` 消息实体
  - [ ] 1.2 创建 `internal/domain/feishu/event.go` 事件实体
  - [ ] 1.3 定义消息类型枚举（文本、图片、卡片等）
  - [ ] 1.4 实现事件解析函数

- [ ] Task 2: 实现租户识别服务 (AC: 2)
  - [ ] 2.1 创建 `internal/domain/feishu/tenant_resolver.go`
  - [ ] 2.2 实现 App ID 到租户 ID 的映射查询
  - [ ] 2.3 实现租户配置缓存（Redis 可选）
  - [ ] 2.4 添加租户不存在时的错误处理

- [ ] Task 3: 实现 OpenClaw 实例查找 (AC: 3)
  - [ ] 3.1 创建 `internal/domain/instance/finder.go`
  - [ ] 3.2 实现根据租户 ID 查找实例
  - [ ] 3.3 实现实例状态验证
  - [ ] 3.4 实现实例连接地址获取

- [ ] Task 4: 实现消息转发服务 (AC: 4)
  - [ ] 4.1 创建 `internal/infrastructure/openclaw/client.go`
  - [ ] 4.2 实现消息转发 HTTP 客户端
  - [ ] 4.3 实现请求超时控制（30 秒）
  - [ ] 4.4 实现转发重试机制

- [ ] Task 5: 实现错误处理 (AC: 5)
  - [ ] 5.1 创建错误类型定义
  - [ ] 5.2 实现错误分类（租户错误、实例错误、转发错误）
  - [ ] 5.3 实现用户友好的错误消息
  - [ ] 5.4 实现错误回复消息发送

- [ ] Task 6: 实现日志记录 (AC: 6)
  - [ ] 6.1 创建消息处理日志中间件
  - [ ] 6.2 实现敏感信息脱敏（App Secret 等）
  - [ ] 6.3 实现处理时间统计
  - [ ] 6.4 实现日志级别分类

- [ ] Task 7: 单元测试 (AC: All)
  - [ ] 7.1 编写消息解析测试
  - [ ] 7.2 编写租户识别测试
  - [ ] 7.3 编写实例查找测试
  - [ ] 7.4 编写消息转发测试
  - [ ] 7.5 编写错误处理测试
  - [ ] 7.6 确保测试覆盖率 ≥ 70%

- [ ] Task 8: 集成测试 (AC: All)
  - [ ] 8.1 编写端到端消息处理测试
  - [ ] 8.2 编写错误场景测试
  - [ ] 8.3 编写性能测试

## Dev Notes

### 消息事件格式

**飞书消息事件结构:**

```json
{
  "schema": "2.0",
  "header": {
    "event_id": "xxx",
    "event_type": "im.message.receive_v1",
    "create_time": "1608877200000",
    "token": "xxx",
    "app_id": "cli_xxx",
    "tenant_key": "xxx"
  },
  "event": {
    "sender": {
      "sender_id": {
        "union_id": "xxx",
        "user_id": "xxx"
      },
      "sender_type": "user"
    },
    "message": {
      "message_id": "xxx",
      "root_id": "xxx",
      "parent_id": "xxx",
      "create_time": "1608877200000",
      "chat_id": "xxx",
      "message_type": "text",
      "content": "{\"text\":\"hello\"}",
      "mentions": []
    }
  }
}
```

### 架构设计

**消息处理流程:**

```
飞书消息事件 -> 事件解析 -> 租户识别 -> 实例查找 -> 消息转发 -> OpenClaw 实例
                |             |             |             |
                v             v             v             v
             日志记录      日志记录      日志记录      日志记录
```

**目录结构:**

```
internal/domain/feishu/
├── message.go          # 消息实体
├── event.go            # 事件实体
├── tenant_resolver.go  # 租户识别服务
└── handler.go          # 消息处理器

internal/infrastructure/openclaw/
├── client.go           # OpenClaw API 客户端
└── config.go           # 客户端配置
```

### 关键技术点

**租户识别:**
- 从事件 header 中获取 app_id
- 查询 feishu_configs 表获取 tenant_id
- 使用缓存提高查询性能（可选 Redis）

**实例查找:**
- 查询 openclaw_instances 表
- 过滤条件：tenant_id, status = 'running'
- 返回实例的内网地址或域名

**消息转发:**
- HTTP POST 请求到 OpenClaw 实例
- 请求体包含原始消息内容
- 超时时间 30 秒（参考 NFR-P2: 3 秒响应要求）

### 性能要求

**来自 NFR:**
- NFR-P2: 用户在飞书发送消息后，在 3 秒内收到 OpenClaw 回复
- 需要优化消息处理链路各环节响应时间

### 错误场景

| 场景 | 错误处理 |
|-----|---------|
| 租户未配置飞书应用 | 返回"请先完成飞书应用配置" |
| OpenClaw 实例不存在 | 返回"请先部署 OpenClaw 实例" |
| 实例状态非 running | 返回"实例正在启动中，请稍后重试" |
| 消息转发超时 | 返回"服务暂时不可用，请稍后重试" |
| 消息转发失败 | 返回"消息处理失败，请稍后重试" |

### References

- [Source: architecture.md#Integration Layer] - 消息路由设计
- [Source: prd.md#FR20, FR21] - 消息收发功能需求
- [Source: prd.md#NFR-P2] - 响应时间要求
- [Source: epics.md#Story 5.2] - 原始故事定义
- 飞书消息事件文档: https://open.feishu.cn/document/client-docs/bot-v3/events/message-received

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
