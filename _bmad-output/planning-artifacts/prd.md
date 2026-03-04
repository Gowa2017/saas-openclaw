---
stepsCompleted: ["step-01-init", "step-02-discovery", "step-02b-vision", "step-02c-executive-summary", "step-03-success", "step-04-journeys", "step-05-domain", "step-06-innovation-skip", "step-07-project-type", "step-08-scoping", "step-09-functional", "step-10-nonfunctional", "step-11-polish", "step-12-complete"]
inputDocuments: ["technical-openclaw-container-platform-research-2026-03-04.md"]
workflowType: 'prd'
briefCount: 0
researchCount: 1
brainstormingCount: 0
projectDocsCount: 0
classification:
  projectType: SaaS B2C Platform
  domain: AI/Developer Tools + Social Integration
  complexity: high
  projectContext: greenfield
---

# Product Requirements Document - saas-openclaw

**Author:** Gowa
**Date:** 2026-03-04

---

## PRD Metadata

- **Project Name:** saas-openclaw
- **Document Version:** 1.0
- **Last Updated:** 2026-03-04
- **Status:** Complete

---

## Document History

| Date | Version | Author | Changes |
|------|---------|--------|---------|
| 2026-03-04 | 1.0 | Gowa | Initial document creation |
| 2026-03-04 | 1.1 | Gowa | Added Executive Summary |
| 2026-03-04 | 1.2 | Gowa | Added Success Criteria and Product Scope |
| 2026-03-04 | 1.3 | Gowa | Added User Journeys, Domain Requirements, SaaS Requirements |
| 2026-03-04 | 1.4 | Gowa | Added Functional and Non-Functional Requirements |

---

## Executive Summary

**产品愿景：** OpenClaw 一键部署平台，让普通用户通过飞书轻松拥有自己的 AI Agent 实例。

**目标用户：** 普通用户（非技术背景），希望使用 AI Agent 但被技术门槛挡住的人群。

**核心问题：** OpenClaw 作为爆火的开源 AI Agent 框架（2025年11月发布，4个月内获得25万+ GitHub stars），具备强大的本地 AI 能力，但部署复杂，需要技术背景（Docker、配置、环境搭建），普通用户无法使用。

**解决方案：** 提供托管式 SaaS 平台，用户通过业务平台一键启动自己的 OpenClaw 实例，在飞书中直接使用，零技术门槛，即时可用。

### What Makes This Special

**核心差异化：** 从"需要技术背景的复杂本地部署"变成"点一下就能跑起来"。对比本地部署，速度更快，操作极简。

**用户价值时刻：** 用户点击按钮后，直接在飞书中与自己的 OpenClaw 对话，无需关心容器、数据库、配置等技术细节。

**核心洞察：** 普通用户也想要 AI Agent 的强大能力，只是被技术门槛挡住了。通过容器化 + 托管 + 社交平台集成，可以普惠技术能力。

**价值主张：** 点一下就能跑起来了。

## Project Classification

| 维度 | 分类 |
|------|------|
| 项目类型 | SaaS B2C Platform（面向普通消费者的托管平台） |
| 领域 | AI/Developer Tools + 社交集成（飞书，后续扩展钉钉、企业微信） |
| 复杂度 | 高（多租户、容器编排、第三方平台集成、AI 安全） |
| 项目背景 | 绿地项目（全新产品） |

**技术架构：** Dokploy + Docker + Traefik + PostgreSQL（平台用户数据）+ Docker Volume（OpenClaw 实例数据持久化）+ 飞书 SSO 集成

---

## Success Criteria

### User Success

**核心成功时刻：** 用户第一次成功让 AI Agent 帮他完成任务

**具体指标：**
- 用户从打开应用到完成第一个 AI 任务的成功率 > 80%
- 用户在首次使用后完成至少一个有效任务（如文件操作、信息查询等）
- 用户感知延迟 < 3 秒（从发送消息到收到响应）

### Business Success

**核心指标：** 用户数量

**时间线目标：**
| 时间节点 | 目标用户数 |
|----------|-----------|
| 第一个月 | 100 注册用户 |
| 第三个月 | 待定（根据 MVP 表现调整） |
| 第十二个月 | 待定（根据增长情况规划） |

**增长来源：** 口碑传播

### Technical Success

**核心指标：** 部署速度 + 可访问性

**具体标准：**
- 部署速度：用户点击"启动"后，OpenClaw 实例在 **3 分钟内**就绪
- 部署成功率 > 95%
- 实例可访问率 > 99%
- 平台整体可用性 > 99%

### Measurable Outcomes

| 指标 | 目标值 | 测量方式 |
|------|--------|----------|
| 首月注册用户 | ≥ 100 | 数据库统计 |
| 首次任务成功率 | > 80% | 用户行为追踪 |
| 部署启动时间 | < 3 分钟 | 系统日志 |
| 实例可访问率 | > 99% | 监控系统 |
| 平台可用性 | > 99% | 监控系统 |

---

## Product Scope

### MVP - Minimum Viable Product

**核心功能：**
1. 业务平台登录/授权
2. 一键部署 OpenClaw 实例
3. 飞书应用配置（输入 App ID 和 Secret）
4. 在飞书中与 OpenClaw 对话

**MVP 成功标准：** 100 个用户成功部署并使用

### Growth Features (Post-MVP)

- 钉钉/企业微信集成
- OpenClaw 配置选项（模型选择、插件管理等）
- 用户使用数据分析
- 备份/恢复功能
- 用户反馈和支持系统

### Vision (Future)

- 多平台 AI Agent 支持（不止 OpenClaw）
- 企业版功能（团队协作、权限管理）
- AI Agent 应用市场
- 付费订阅模式（按月付费，当用户规模达到一定水平后）

---

## User Journeys

### Journey 1: 主要用户 - 成功路径（李明首次开通）

**人物：** 李明，32岁，售前解决方案工程师

**背景：** 李明需要 AI 帮助编写方案，但不会部署 OpenClaw。

**旅程流程：**

1. **业务平台登录** - 李明登录公司业务平台，进入"AI 服务"模块
2. **填写飞书配置** - 按照教程在飞书开放平台创建应用，获取 App ID 和 Secret
3. **一键部署** - 点击"启动我的 OpenClaw"，等待 2-3 分钟
4. **飞书添加机器人** - 按照教程在飞书中搜索并添加 OpenClaw 机器人
5. **首次对话** - 发送消息，收到 OpenClaw 回复，成功完成第一个任务

**需求揭示：** 业务平台登录、飞书配置引导、一键部署、部署进度反馈、飞书机器人添加教程

### Journey 2: 主要用户 - 边缘场景（配置错误恢复）

**人物：** 李明

**场景：** 李明第一次配置时，Secret 复制少了字符，部署失败。

**旅程流程：**

1. **部署失败** - 系统显示"部署失败，请检查配置后重新申请开通"
2. **查看错误详情** - 点击查看具体错误原因
3. **重新配置** - 返回配置页面，已填写的值保留
4. **修正并重试** - 修正错误后重新部署，成功

**需求揭示：** 错误详情显示、配置保留、重新部署功能

### Journey 3: 平台管理员 - 日常监控与故障处理

**人物：** 王运维，28岁，SaaS 平台运维工程师

**职责：** 监控平台运行状态、处理用户问题

**旅程流程：**

1. **查看仪表盘** - 查看用户数、活跃实例、部署成功率、系统可用性
2. **检查告警** - 发现实例重启失败告警
3. **定位问题** - 查看日志，发现资源不足
4. **处理问题** - 添加工作节点，重新部署失败实例
5. **通知用户** - 发送系统通知告知用户实例已恢复

**需求揭示：** 管理仪表盘、告警系统、日志查看、实例管理、用户通知

### Journey Requirements Summary

| 功能领域 | 具体需求 |
|----------|----------|
| 用户认证 | 业务平台登录 + 飞书授权 |
| 配置引导 | 分步教程 + 飞书开放平台链接 |
| 部署管理 | 一键部署 + 进度显示 + 失败通知 |
| 文档支持 | 飞书机器人添加教程 |
| 配置管理 | 业务平台重新配置入口 |
| 错误处理 | 错误详情显示 + 配置保留 |
| 管理功能 | 仪表盘 + 告警 + 日志 + 实例管理 + 通知 |

---

## Domain-Specific Requirements

### Compliance & Regulatory

**数据存储：** 所有数据存储在国内服务器，符合数据本地化要求

**数据隐私：** 用户与 OpenClaw 的对话内容无需特殊处理；平台仅存储用户账号、订阅状态、OpenClaw 配置

### Technical Constraints

**多租户隔离：**
- 每个用户拥有独立的 OpenClaw 容器实例
- 每个用户拥有独立的 Docker Volume 数据空间
- 用户间完全隔离，互不可见

**安全防护：**
- 容器级隔离防止跨用户攻击
- 平台与 OpenClaw 实例分离，平台故障不影响实例运行

**备份策略：**
- OpenClaw 实例数据：Dokploy S3 备份
- 用户配置：PostgreSQL 数据库备份

### Integration Requirements

**飞书企业自建应用：**
- 无需应用市场审核
- 使用飞书 Go SDK，长连接形式
- 权限：接收消息、发送消息、用户身份信息（OAuth）

**飞书用户信息处理：**
- 不存储飞书用户信息
- 仅存储飞书 App ID 和 Secret（用于服务端调用）
- 用户身份通过飞书 OAuth 实时验证

### Risk Mitigations

| 风险 | 缓解措施 |
|------|----------|
| OpenClaw 实例崩溃 | 数据自动备份，可快速恢复 |
| 平台宕机 | Dokploy 部署的实例独立运行，不受影响 |
| 恶意用户攻击平台 | 容器级隔离，用户无法访问平台核心 |
| 用户数据丢失 | 定期自动备份 + 配置备份 |

---

## SaaS Platform Specific Requirements

### Multi-Tenancy Model

**隔离策略：**
- 每个用户拥有独立的 OpenClaw 容器实例
- 每个用户拥有独立的 Docker Volume 数据空间
- 用户间完全隔离，互不可见

**数据隔离：**
- 平台用户数据：PostgreSQL 数据库隔离（每用户独立记录）
- OpenClaw 实例数据：独立容器 + 独立 Volume

**协作模型：** MVP 阶段不支持协作；未来可考虑团队/企业版功能

### Permission Model

**普通用户权限：**
- 访问自己的 OpenClaw 实例
- 在业务平台查看自己的配置
- 在业务平台重新配置自己的实例

**平台管理员权限：**
- 查看所有用户列表和状态
- 查看所有实例的部署状态和日志
- 重启/停止/启动任何用户实例
- 查看平台整体指标
- 发送系统通知给用户

### Subscription Model

**当前模式：** 完全免费，无订阅层级

**未来付费模式（规划）：** 按月订阅
- 免费版：基础功能，有限资源
- 专业版：更多资源，高级配置选项
- 团队版：多用户协作，企业级支持

### Integration Requirements

**飞书集成（MVP）：**
- 使用飞书 Go SDK
- 长连接形式（事件监听）
- 消息接收、发送
- OAuth 身份验证

**未来集成：** 钉钉、企业微信（Post-MVP）

---

## Functional Requirements

### 用户认证与授权

- FR1: 用户可以通过业务平台登录系统
- FR2: 用户可以授权飞书账号访问权限
- FR3: 系统可以验证用户的飞书授权状态
- FR4: 平台管理员可以登录管理后台

### 飞书应用配置

- FR5: 用户可以输入飞书应用的 App ID
- FR6: 用户可以输入飞书应用的 App Secret
- FR7: 用户可以保存飞书应用配置
- FR8: 用户可以修改已保存的飞书应用配置
- FR9: 系统可以验证飞书应用配置的有效性

### OpenClaw 实例管理

- FR10: 用户可以启动自己的 OpenClaw 实例
- FR11: 系统可以在 3 分钟内部署 OpenClaw 实例
- FR12: 用户可以查看实例部署状态
- FR13: 系统可以通知用户部署成功或失败
- FR14: 系统可以在部署失败时显示错误原因
- FR15: 平台管理员可以查看所有用户的实例状态
- FR16: 平台管理员可以重启用户实例
- FR17: 平台管理员可以停止用户实例
- FR18: 平台管理员可以启动用户实例

### 飞书机器人使用

- FR19: 用户可以在飞书中添加 OpenClaw 机器人
- FR20: 用户可以在飞书中向 OpenClaw 机器人发送消息
- FR21: OpenClaw 机器人可以接收用户发送的消息
- FR22: OpenClaw 机器人可以向用户发送回复消息
- FR23: 系统提供飞书机器人添加教程

### 数据管理

- FR24: 系统可以备份用户配置数据
- FR25: 系统可以备份 OpenClaw 实例数据
- FR26: 平台管理员可以查看系统备份状态
- FR27: 系统可以在实例重新部署时恢复用户配置

### 监控与运维

- FR28: 平台管理员可以查看当前用户总数
- FR29: 平台管理员可以查看活跃实例数量
- FR30: 平台管理员可以查看部署成功率
- FR31: 平台管理员可以查看系统可用性
- FR32: 平台管理员可以查看实例部署日志
- FR33: 系统可以在资源不足时发送告警通知
- FR34: 平台管理员可以向用户发送系统通知

### 用户引导

- FR35: 系统可以显示飞书开放平台配置教程
- FR36: 系统可以显示飞书机器人添加教程
- FR37: 系统可以在用户首次访问时显示引导说明

---

## Non-Functional Requirements

### Performance

- NFR-P1: OpenClaw 实例部署在 3 分钟内完成
- NFR-P2: 用户在飞书发送消息后，在 3 秒内收到 OpenClaw 回复
- NFR-P3: 业务平台页面加载时间 < 2 秒
- NFR-P4: 管理后台仪表盘数据查询时间 < 5 秒

### Security

- NFR-S1: 用户间 OpenClaw 容器完全隔离，无法互相访问
- NFR-S2: 用户配置数据（App ID/Secret）加密存储
- NFR-S3: 平台管理员操作需要身份验证
- NFR-S4: 飞书 OAuth 遵循 OAuth 2.0 标准
- NFR-S5: 所有 API 通信使用 HTTPS 加密

### Scalability

- NFR-SC1: 系统支持从 100 用户扩展到 1000 用户无需架构变更
- NFR-SC2: 新容器部署可以在 5 分钟内完成资源分配
- NFR-SC3: 数据库支持每秒 100 个并发请求
- NFR-SC4: Docker Swarm 集群可以通过添加节点水平扩展

### Reliability

- NFR-R1: 平台整体可用性 ≥ 99%
- NFR-R2: OpenClaw 实例可访问率 ≥ 99%
- NFR-R3: 部署成功率 ≥ 95%
- NFR-R4: 数据每日自动备份，备份成功率 ≥ 99%
- NFR-R5: 飞书长连接断开后自动重连，重连时间 < 30 秒

### Integration

- NFR-I1: 飞书 SDK 长连接稳定性 ≥ 99%
- NFR-I2: 飞书消息投递成功率 ≥ 99%
- NFR-I3: Dokploy API 调用响应时间 < 5 秒
- NFR-I4: 飞书 SDK 版本兼容性支持至少 2 个大版本
