# Story 5.1: 飞书 SDK 长连接集成

Status: ready-for-dev

## Story

As a 后端开发者,
I want 集成飞书 Go SDK 长连接,
so that 可以实时接收飞书消息事件。

## Acceptance Criteria

1. **AC1: 飞书长连接服务启动**
   - **Given** 飞书应用配置已完成
   - **When** 启动飞书长连接服务
   - **Then** 成功连接飞书服务器
   - **And** 连接状态记录到日志

2. **AC2: 消息接收事件订阅**
   - **Given** 长连接已建立
   - **When** 飞书推送消息事件
   - **Then** 系统成功接收事件
   - **And** 事件被正确解析和处理

3. **AC3: 自动重连机制**
   - **Given** 长连接断开
   - **When** 检测到连接断开
   - **Then** 在 30 秒内自动重连
   - **And** 重连成功后恢复消息接收
   - **And** 记录重连日志

4. **AC4: 连接状态监控**
   - **Given** 长连接服务运行中
   - **When** 查询连接状态
   - **Then** 返回当前连接状态（connected/disconnected/reconnecting）
   - **And** 返回最后心跳时间
   - **And** 返回重连次数统计

5. **AC5: 多租户消息路由支持**
   - **Given** 系统有多个租户
   - **When** 接收到消息事件
   - **Then** 根据 App ID 识别租户
   - **And** 消息路由到正确的处理队列

## Tasks / Subtasks

- [ ] Task 1: 添加飞书 SDK 依赖 (AC: 1, 2)
  - [ ] 1.1 添加 `github.com/larksuite/oapi-sdk-go/v3` 依赖到 go.mod
  - [ ] 1.2 创建 `internal/infrastructure/feishu/` 目录
  - [ ] 1.3 创建 `client.go` 飞书客户端封装

- [ ] Task 2: 实现长连接服务 (AC: 1, 2, 3)
  - [ ] 2.1 创建 `internal/infrastructure/feishu/websocket/` 目录
  - [ ] 2.2 实现 WebSocket 连接管理器
  - [ ] 2.3 实现消息事件处理器
  - [ ] 2.4 实现心跳保活机制
  - [ ] 2.5 实现自动重连逻辑（30 秒内重连）

- [ ] Task 3: 实现连接状态监控 (AC: 4)
  - [ ] 3.1 创建 `internal/infrastructure/feishu/monitor/` 目录
  - [ ] 3.2 实现连接状态追踪
  - [ ] 3.3 实现心跳时间记录
  - [ ] 3.4 实现重连计数器
  - [ ] 3.5 添加状态查询 API

- [ ] Task 4: 实现多租户路由 (AC: 5)
  - [ ] 4.1 创建 `internal/domain/feishu/message.go` 消息领域模型
  - [ ] 4.2 创建 `internal/infrastructure/feishu/router/` 目录
  - [ ] 4.3 实现消息路由器
  - [ ] 4.4 实现租户识别逻辑（根据 App ID）
  - [ ] 4.5 实现消息队列分发

- [ ] Task 5: 日志记录与监控 (AC: 1, 3, 4)
  - [ ] 5.1 添加连接事件日志
  - [ ] 5.2 添加消息接收日志
  - [ ] 5.3 添加重连事件日志
  - [ ] 5.4 添加错误日志

- [ ] Task 6: 单元测试 (AC: All)
  - [ ] 6.1 编写客户端测试
  - [ ] 6.2 编写连接管理器测试
  - [ ] 6.3 编写路由器测试
  - [ ] 6.4 编写监控器测试
  - [ ] 6.5 确保测试覆盖率 ≥ 70%

- [ ] Task 7: 集成测试 (AC: All)
  - [ ] 7.1 编写端到端连接测试
  - [ ] 7.2 编写重连场景测试
  - [ ] 7.3 编写多租户路由测试

## Dev Notes

### 技术选型

**飞书 Go SDK:**
- 使用官方 SDK: `github.com/larksuite/oapi-sdk-go/v3`
- 长连接模式: WebSocket 连接
- 支持事件订阅: 消息接收、机器人事件等

### 架构设计

**目录结构:**

```
internal/infrastructure/feishu/
├── client.go              # 飞书客户端封装
├── config.go              # 飞书配置
├── websocket/
│   ├── connection.go      # WebSocket 连接管理
│   ├── handler.go         # 消息处理器
│   └── reconnect.go       # 重连逻辑
├── router/
│   ├── router.go          # 消息路由器
│   └── tenant_resolver.go # 租户识别器
└── monitor/
    ├── status.go          # 连接状态
    └── metrics.go         # 监控指标
```

### 关键技术点

**长连接重连策略:**
- 检测断开: 心跳超时或连接错误
- 重连间隔: 指数退避（1s, 2s, 4s, 8s, 16s），最大 30 秒
- 重连次数: 无限重试直到成功
- 状态通知: 重连开始/成功/失败事件

**多租户路由:**
- 每个 App ID 对应一个租户
- 消息事件中包含 App ID
- 路由表: App ID -> 租户 ID 映射
- 支持 App ID 动态注册和注销

### 性能要求

**来自 NFR:**
- NFR-R5: 飞书长连接断开后自动重连，重连时间 < 30 秒
- NFR-I1: 飞书 SDK 长连接稳定性 ≥ 99%
- NFR-I2: 飞书消息投递成功率 ≥ 99%
- NFR-I4: 飞书 SDK 版本兼容性支持至少 2 个大版本

### References

- [Source: architecture.md#Integration Layer] - 飞书 SDK 集成要求
- [Source: prd.md#NFR-R5] - 重连时间要求
- [Source: prd.md#NFR-I1] - 长连接稳定性要求
- [Source: epics.md#Story 5.1] - 原始故事定义
- 飞书开放平台文档: https://open.feishu.cn/document/client-docs/bot-v3/events/overview

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
