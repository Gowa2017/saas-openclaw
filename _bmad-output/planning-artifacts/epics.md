---
stepsCompleted: ["step-01-validate-prerequisites", "step-02-design-epics", "step-03-create-stories", "step-04-final-validation"]
workflowStatus: "completed"
inputDocuments: ["prd.md", "architecture.md", "ux-design-specification.md"]
workflowType: 'epics-and-stories'
project_name: 'saas-openclaw'
user_name: 'Gowa'
date: '2026-03-04'
---

# saas-openclaw - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for saas-openclaw, decomposing the requirements from the PRD, UX Design if it exists, and Architecture requirements into implementable stories.

## Requirements Inventory

### Functional Requirements

**用户认证与授权**
- FR1: 用户可以通过业务平台登录系统
- FR2: 用户可以授权飞书账号访问权限
- FR3: 系统可以验证用户的飞书授权状态
- FR4: 平台管理员可以登录管理后台

**飞书应用配置**
- FR5: 用户可以输入飞书应用的 App ID
- FR6: 用户可以输入飞书应用的 App Secret
- FR7: 用户可以保存飞书应用配置
- FR8: 用户可以修改已保存的飞书应用配置
- FR9: 系统可以验证飞书应用配置的有效性

**OpenClaw 实例管理**
- FR10: 用户可以启动自己的 OpenClaw 实例
- FR11: 系统可以在 3 分钟内部署 OpenClaw 实例
- FR12: 用户可以查看实例部署状态
- FR13: 系统可以通知用户部署成功或失败
- FR14: 系统可以在部署失败时显示错误原因
- FR15: 平台管理员可以查看所有用户的实例状态
- FR16: 平台管理员可以重启用户实例
- FR17: 平台管理员可以停止用户实例
- FR18: 平台管理员可以启动用户实例

**飞书机器人使用**
- FR19: 用户可以在飞书中添加 OpenClaw 机器人
- FR20: 用户可以在飞书中向 OpenClaw 机器人发送消息
- FR21: OpenClaw 机器人可以接收用户发送的消息
- FR22: OpenClaw 机器人可以向用户发送回复消息
- FR23: 系统提供飞书机器人添加教程

**数据管理**
- FR24: 系统可以备份用户配置数据
- FR25: 系统可以备份 OpenClaw 实例数据
- FR26: 平台管理员可以查看系统备份状态
- FR27: 系统可以在实例重新部署时恢复用户配置

**监控与运维**
- FR28: 平台管理员可以查看当前用户总数
- FR29: 平台管理员可以查看活跃实例数量
- FR30: 平台管理员可以查看部署成功率
- FR31: 平台管理员可以查看系统可用性
- FR32: 平台管理员可以查看实例部署日志
- FR33: 系统可以在资源不足时发送告警通知
- FR34: 平台管理员可以向用户发送系统通知

**用户引导**
- FR35: 系统可以显示飞书开放平台配置教程
- FR36: 系统可以显示飞书机器人添加教程
- FR37: 系统可以在用户首次访问时显示引导说明

### NonFunctional Requirements

**Performance**
- NFR-P1: OpenClaw 实例部署在 3 分钟内完成
- NFR-P2: 用户在飞书发送消息后，在 3 秒内收到 OpenClaw 回复
- NFR-P3: 业务平台页面加载时间 < 2 秒
- NFR-P4: 管理后台仪表盘数据查询时间 < 5 秒

**Security**
- NFR-S1: 用户间 OpenClaw 容器完全隔离，无法互相访问
- NFR-S2: 用户配置数据（App ID/Secret）加密存储
- NFR-S3: 平台管理员操作需要身份验证
- NFR-S4: 飞书 OAuth 遵循 OAuth 2.0 标准
- NFR-S5: 所有 API 通信使用 HTTPS 加密

**Scalability**
- NFR-SC1: 系统支持从 100 用户扩展到 1000 用户无需架构变更
- NFR-SC2: 新容器部署可以在 5 分钟内完成资源分配
- NFR-SC3: 数据库支持每秒 100 个并发请求
- NFR-SC4: Docker Swarm 集群可以通过添加节点水平扩展

**Reliability**
- NFR-R1: 平台整体可用性 ≥ 99%
- NFR-R2: OpenClaw 实例可访问率 ≥ 99%
- NFR-R3: 部署成功率 ≥ 95%
- NFR-R4: 数据每日自动备份，备份成功率 ≥ 99%
- NFR-R5: 飞书长连接断开后自动重连，重连时间 < 30 秒

**Integration**
- NFR-I1: 飞书 SDK 长连接稳定性 ≥ 99%
- NFR-I2: 飞书消息投递成功率 ≥ 99%
- NFR-I3: Dokploy API 调用响应时间 < 5 秒
- NFR-I4: 飞书 SDK 版本兼容性支持至少 2 个大版本

### Additional Requirements

**来自架构文档：**

**Starter Template（Epic 1 Story 1 必须使用）：**
- 后端：go-rest-api 模板（Clean Architecture + SOLID）
- 前端：Vite + Vue 3 + TypeScript + Naive UI

**技术栈要求：**
- RESTful API 设计：资源导向，标准 HTTP 方法，版本化 (/v1/)
- Pinia 状态管理：按模块组织 (auth, tenant, instance, config)
- GitHub Actions CI/CD：代码推送 → 自动构建 → 自动测试 → 自动部署
- PostgreSQL Database per Tenant：多租户数据隔离
- Docker 网络隔离：多租户安全隔离

**命名约定：**
- 数据库表名：snake_case（例：tenant_users, openclaw_instances）
- 数据库列名：PascalCase（例：TenantID, CreatedAt）
- API 端点：复数资源名（例：/tenants, /instances）
- 组件名：PascalCase（例：TenantList, InstanceCard）
- 文件名：kebab-case（例：tenant-list.vue, instance-card.vue）

**API 响应格式：**
- 统一包装器：`{ data: {...}, error: null, meta: {...} }`

**来自 UX 设计文档：**

**布局策略：**
- 首次部署流程：向导式聚焦布局（分步引导）
- 实例管理界面：卡片式仪表盘布局
- 平台优先级：Web 桌面优先，平板 P1，移动端 P2

**核心自定义组件：**
- 部署进度组件：显示实时部署进度和预计时间
- 飞书配置引导组件：分步图文教程
- 实例状态卡片：基于 Card 定制
- 部署成功动画：庆祝动画效果

**无障碍要求：**
- WCAG 2.1 AA 级合规
- 色彩对比度 4.5:1
- 键盘导航支持
- 最小触控区域 44x44px

**UX 模式：**
- 主要按钮：蓝色填充 #1677FF
- 成功反馈：Message.success 3秒自动消失
- 错误反馈：Message.error 需手动关闭
- 表单验证：实时校验 + 提交校验 + 后端校验

### FR Coverage Map

FR1: Epic 2 - 用户业务平台登录
FR2: Epic 2 - 飞书账号授权
FR3: Epic 2 - 飞书授权验证
FR4: Epic 2 - 管理员登录后台
FR5: Epic 3 - 输入飞书 App ID
FR6: Epic 3 - 输入飞书 App Secret
FR7: Epic 3 - 保存飞书配置
FR8: Epic 3 - 修改飞书配置
FR9: Epic 3 - 验证飞书配置有效性
FR10: Epic 4 - 启动 OpenClaw 实例
FR11: Epic 4 - 3 分钟内部署实例
FR12: Epic 4 - 查看部署状态
FR13: Epic 4 - 部署成功/失败通知
FR14: Epic 4 - 显示部署错误原因
FR15: Epic 4 - 管理员查看所有实例
FR16: Epic 4 - 管理员重启实例
FR17: Epic 4 - 管理员停止实例
FR18: Epic 4 - 管理员启动实例
FR19: Epic 5 - 飞书添加机器人
FR20: Epic 5 - 飞书发送消息
FR21: Epic 5 - 机器人接收消息
FR22: Epic 5 - 机器人发送回复
FR23: Epic 5 - 机器人添加教程
FR24: Epic 6 - 备份用户配置
FR25: Epic 6 - 备份实例数据
FR26: Epic 6 - 查看备份状态
FR27: Epic 6 - 恢复用户配置
FR28: Epic 7 - 查看用户总数
FR29: Epic 7 - 查看活跃实例数
FR30: Epic 7 - 查看部署成功率
FR31: Epic 7 - 查看系统可用性
FR32: Epic 7 - 查看部署日志
FR33: Epic 7 - 资源告警通知
FR34: Epic 7 - 发送系统通知
FR35: Epic 8 - 飞书开放平台教程
FR36: Epic 8 - 飞书机器人教程
FR37: Epic 8 - 首次访问引导

## Epic List

### Epic 1: 项目基础架构搭建
建立项目开发基础设施，包括前后端项目初始化、CI/CD 流水线、基础配置等，为后续功能开发提供坚实基础。
**FRs covered:** 技术基础（架构文档要求的 Starter Template、CI/CD、命名规范等）
**关键技术决策：** go-rest-api 后端模板 + Vite + Vue 3 + TypeScript + Naive UI 前端

### Epic 2: 用户认证与授权
用户可以通过业务平台登录系统并授权飞书账号，平台管理员可以独立登录管理后台。实现完整的身份认证体系，支持业务平台 JWT Token 验证和飞书 OAuth 授权。
**FRs covered:** FR1, FR2, FR3, FR4
**用户价值：** 用户可以安全登录系统，管理员可以访问管理后台

### Epic 3: 飞书应用配置
用户可以输入、保存、修改和验证飞书应用的 App ID 和 Secret 配置。提供分步图文教程引导用户在飞书开放平台创建应用并获取凭证。
**FRs covered:** FR5, FR6, FR7, FR8, FR9
**用户价值：** 用户可以轻松完成飞书应用配置，为部署 OpenClaw 实例做准备

### Epic 4: OpenClaw 实例管理
用户可以一键启动自己的 OpenClaw 实例，系统在 3 分钟内完成部署。用户可以查看部署状态和进度，系统在部署成功或失败时通知用户。平台管理员可以查看所有用户的实例状态，并执行重启、停止、启动操作。
**FRs covered:** FR10, FR11, FR12, FR13, FR14, FR15, FR16, FR17, FR18
**用户价值：** 用户可以一键部署 OpenClaw 实例，实现"点一下就能跑起来"的核心价值主张

### Epic 5: 飞书机器人集成
用户可以在飞书中添加 OpenClaw 机器人，发送消息并接收回复。系统提供飞书机器人添加教程。实现飞书 SDK 长连接集成，确保消息实时收发。
**FRs covered:** FR19, FR20, FR21, FR22, FR23
**用户价值：** 用户可以在熟悉的飞书环境中使用 OpenClaw，无需切换应用

### Epic 6: 数据备份与恢复
系统可以自动备份用户配置数据和 OpenClaw 实例数据。平台管理员可以查看备份状态。在实例重新部署时可以恢复用户配置。
**FRs covered:** FR24, FR25, FR26, FR27
**用户价值：** 用户数据安全有保障，不会因系统故障而丢失

### Epic 7: 平台监控与运维
平台管理员可以通过仪表盘查看用户总数、活跃实例数量、部署成功率、系统可用性等关键指标。可以查看实例部署日志，接收资源不足告警，并向用户发送系统通知。
**FRs covered:** FR28, FR29, FR30, FR31, FR32, FR33, FR34
**用户价值：** 管理员可以实时监控平台运行状态，及时发现和处理问题

### Epic 8: 用户引导与帮助
系统提供飞书开放平台配置教程、飞书机器人添加教程，并在用户首次访问时显示引导说明。帮助非技术用户快速上手。
**FRs covered:** FR35, FR36, FR37
**用户价值：** 用户可以获得清晰的帮助和指导，降低使用门槛

---

## Epic 1: 项目基础架构搭建

建立项目开发基础设施，包括前后端项目初始化、CI/CD 流水线、基础配置等，为后续功能开发提供坚实基础。

### Story 1.1: 后端项目初始化

As a 开发者,
I want 使用 go-rest-api 模板初始化后端项目,
So that 可以快速搭建符合 Clean Architecture 的后端服务。

**Acceptance Criteria:**

**Given** 开发环境已安装 Go 1.21+
**When** 执行后端项目初始化命令
**Then** 后端项目结构符合 go-rest-api 模板规范
**And** 包含 cmd/server/main.go 入口文件
**And** 包含 internal/api、internal/domain、internal/infrastructure 目录结构
**And** 包含 pkg 公共工具目录
**And** go.mod 文件正确配置依赖

### Story 1.2: 前端项目初始化

As a 开发者,
I want 使用 Vite + Vue 3 + TypeScript + Naive UI 初始化前端项目,
So that 可以快速搭建现代化前端应用。

**Acceptance Criteria:**

**Given** 开发环境已安装 Node.js 18+
**When** 执行前端项目初始化命令
**Then** 前端项目使用 Vite 构建工具
**And** 使用 Vue 3 + TypeScript
**And** 集成 Naive UI 组件库
**And** 配置 Pinia 状态管理
**And** 配置 Vue Router 路由
**And** 项目结构符合架构文档规范

### Story 1.3: PostgreSQL 数据库配置

As a 开发者,
I want 配置 PostgreSQL 数据库连接,
So that 后端服务可以访问数据库存储数据。

**Acceptance Criteria:**

**Given** PostgreSQL 数据库服务已启动
**When** 配置数据库连接参数
**Then** 后端服务成功连接 PostgreSQL 数据库
**And** 支持连接池配置（最大连接数 100）
**And** 支持 .env 环境变量配置
**And** 数据库连接使用加密传输

### Story 1.4: GitHub Actions CI/CD 流水线

As a 开发者,
I want 配置 GitHub Actions CI/CD 流水线,
So that 代码推送后自动执行构建、测试和部署。

**Acceptance Criteria:**

**Given** 代码仓库已在 GitHub 创建
**When** 推送代码到 main 分支
**Then** 自动触发 CI 流水线
**And** 后端 CI 执行 go build、go test
**And** 前端 CI 执行 npm build、npm test
**And** 构建失败时发送通知
**And** 构建产物自动上传

### Story 1.5: Docker 容器化配置

As a 开发者,
I want 为前后端服务配置 Dockerfile,
So that 服务可以通过 Docker 容器部署。

**Acceptance Criteria:**

**Given** 前后端项目已初始化
**When** 编写 Dockerfile 文件
**Then** 后端 Dockerfile 使用多阶段构建优化镜像大小
**And** 前端 Dockerfile 使用 nginx 静态服务
**And** docker-compose.yml 配置本地开发环境
**And** 包含 .dockerignore 文件排除不必要的文件

---

## Epic 2: 用户认证与授权

用户可以通过业务平台登录系统并授权飞书账号，平台管理员可以独立登录管理后台。实现完整的身份认证体系，支持业务平台 JWT Token 验证和飞书 OAuth 授权。

### Story 2.1: 用户数据模型与数据库表

As a 开发者,
I want 创建用户数据模型和数据库表,
So that 系统可以存储和管理用户信息。

**Acceptance Criteria:**

**Given** PostgreSQL 数据库已配置
**When** 执行数据库迁移
**Then** 创建 tenant_users 表（表名 snake_case）
**And** 表包含 ID、TenantID、Name、Email、Role、CreatedAt、UpdatedAt 列（列名 PascalCase）
**And** 创建 admin_users 表用于平台管理员
**And** 创建数据库索引优化查询性能

### Story 2.2: 业务平台 JWT Token 验证中间件

As a 后端开发者,
I want 实现 JWT Token 验证中间件,
So that 可以验证来自业务平台的用户身份。

**Acceptance Criteria:**

**Given** 业务平台使用 JWT Token 认证
**When** 前端请求携带 X-Platform-Token Header
**Then** 中间件验证 Token 有效性
**And** 解析 Token 获取用户信息
**And** 调用业务平台接口获取用户详细信息（租户、权限等）
**And** 将用户信息存入请求上下文
**And** Token 无效时返回 401 错误

### Story 2.3: 用户信息获取与租户识别

As a 后端开发者,
I want 实现用户信息获取和租户识别功能,
So that 可以正确识别用户所属租户并实现多租户隔离。

**Acceptance Criteria:**

**Given** JWT Token 验证通过
**When** 调用业务平台用户信息接口
**Then** 获取用户 ID、姓名、邮箱、租户 ID 等信息
**And** 根据租户 ID 实现数据隔离
**And** 自动创建或更新本地用户记录
**And** 返回统一的 API 响应格式

### Story 2.4: 平台管理员独立认证系统

As a 平台管理员,
I want 使用独立的用户名密码登录管理后台,
So that 可以与普通用户认证体系分离。

**Acceptance Criteria:**

**Given** 管理员账号已创建
**When** 管理员输入用户名密码登录
**Then** 验证用户名密码正确性
**And** 生成管理员 JWT Token
**And** Token 包含管理员角色标识
**And** 密码使用 bcrypt 加密存储
**And** 登录失败返回明确的错误信息

### Story 2.5: 前端登录页面与状态管理

As a 前端开发者,
I want 实现登录页面和认证状态管理,
So that 用户可以登录系统并保持登录状态。

**Acceptance Criteria:**

**Given** 前端项目已初始化
**When** 用户访问需要认证的页面
**Then** 未登录用户重定向到登录页
**And** 登录页显示业务平台登录入口
**And** 登录成功后将 Token 存储在 localStorage
**And** 使用 Pinia 管理认证状态（auth store）
**And** 支持自动刷新 Token 机制

### Story 2.6: 管理员登录页面

As a 平台管理员,
I want 访问独立的管理员登录页面,
So that 可以登录管理后台。

**Acceptance Criteria:**

**Given** 管理员登录页面已创建
**When** 访问 /admin/login 路径
**Then** 显示管理员登录表单（用户名、密码）
**And** 表单有实时校验（必填、格式）
**And** 登录成功后跳转到管理后台首页
**And** 登录失败显示错误提示
**And** 支持记住登录状态功能

---

## Epic 3: 飞书应用配置

用户可以输入、保存、修改和验证飞书应用的 App ID 和 Secret 配置。提供分步图文教程引导用户在飞书开放平台创建应用并获取凭证。

### Story 3.1: 飞书配置数据模型

As a 开发者,
I want 创建飞书配置数据模型,
So that 可以存储用户的飞书应用配置。

**Acceptance Criteria:**

**Given** 用户数据模型已创建
**When** 创建 feishu_configs 表
**Then** 表包含 ID、TenantID、AppID、AppSecret、Status、CreatedAt、UpdatedAt 列
**And** AppSecret 列使用加密存储（AES-256）
**And** 一个租户只有一条飞书配置记录
**And** 创建索引优化查询性能

### Story 3.2: 飞书配置 CRUD API

As a 后端开发者,
I want 实现飞书配置的增删改查 API,
So that 前端可以管理飞书配置。

**Acceptance Criteria:**

**Given** 飞书配置数据模型已创建
**When** 实现以下 API 端点
**Then** POST /v1/feishu-configs - 创建飞书配置
**And** GET /v1/feishu-configs - 获取当前用户的飞书配置
**And** PUT /v1/feishu-configs - 更新飞书配置
**And** DELETE /v1/feishu-configs - 删除飞书配置
**And** 所有 API 需要认证
**And** 返回统一的 API 响应格式

### Story 3.3: 飞书配置验证功能

As a 用户,
I want 验证飞书配置是否有效,
So that 确保配置正确后再部署实例。

**Acceptance Criteria:**

**Given** 用户已输入飞书 App ID 和 Secret
**When** 点击"验证配置"按钮
**Then** 后端调用飞书 API 验证凭证有效性
**And** 验证成功显示"配置有效"提示
**And** 验证失败显示具体错误原因（如 App ID 不存在、Secret 错误等）
**And** 验证结果缓存 5 分钟
**And** 保存配置时自动执行验证

### Story 3.4: 飞书配置前端页面

As a 用户,
I want 在前端页面配置飞书应用信息,
So that 可以完成 OpenClaw 实例部署前的准备工作。

**Acceptance Criteria:**

**Given** 用户已登录系统
**When** 访问飞书配置页面
**Then** 显示 App ID 和 App Secret 输入框
**And** App ID 格式实时校验（cli_ 开头）
**And** App Secret 显示为密码输入框
**And** 提供"验证配置"按钮
**And** 提供"保存配置"按钮
**And** 已有配置时自动填充并显示状态

### Story 3.5: 飞书开放平台配置教程组件

As a 用户,
I want 查看飞书开放平台配置教程,
So that 可以了解如何创建飞书应用并获取 App ID 和 Secret。

**Acceptance Criteria:**

**Given** 用户在飞书配置页面
**When** 查看"如何获取配置"区域
**Then** 显示分步教程（4 步）
**And** 步骤 1：跳转到飞书开放平台链接
**And** 步骤 2：创建企业自建应用截图指引
**And** 步骤 3：获取 App ID/Secret 位置截图（高亮标注）
**And** 步骤 4：返回平台填写指引
**And** 支持展开/收起截图

---

## Epic 4: OpenClaw 实例管理

用户可以一键启动自己的 OpenClaw 实例，系统在 3 分钟内完成部署。用户可以查看部署状态和进度，系统在部署成功或失败时通知用户。平台管理员可以查看所有用户的实例状态，并执行重启、停止、启动操作。

### Story 4.1: 实例数据模型

As a 开发者,
I want 创建 OpenClaw 实例数据模型,
So that 可以存储实例的状态和配置信息。

**Acceptance Criteria:**

**Given** PostgreSQL 数据库已配置
**When** 创建 openclaw_instances 表
**Then** 表包含 ID、TenantID、Name、Status、ContainerID、DeployLog、CreatedAt、UpdatedAt 列
**And** Status 包含：pending、deploying、running、stopped、error 状态
**And** 一个租户可以有多个实例
**And** 创建索引优化查询性能

### Story 4.2: Dokploy API 客户端集成

As a 后端开发者,
I want 集成 Dokploy API 客户端,
So that 可以通过 API 管理 OpenClaw 容器实例。

**Acceptance Criteria:**

**Given** Dokploy 服务已部署
**When** 实现 Dokploy API 客户端
**Then** 支持 CreateApplication 接口创建容器
**And** 支持 GetApplication 接口查询容器状态
**And** 支持 StartApplication/StopApplication/RestartApplication 接口
**And** 支持 GetLogs 接口获取容器日志
**And** API 调用超时设置为 5 秒
**And** 实现重试机制（最多 3 次）

### Story 4.3: 一键部署 OpenClaw 实例

As a 用户,
I want 点击按钮一键部署 OpenClaw 实例,
So that 可以快速开始使用 AI Agent。

**Acceptance Criteria:**

**Given** 用户已完成飞书配置
**When** 点击"启动我的 OpenClaw"按钮
**Then** 系统创建部署任务并返回任务 ID
**And** 调用 Dokploy API 创建容器
**And** 容器使用 OpenClaw 官方镜像
**And** 自动配置环境变量（飞书 App ID/Secret）
**And** 容器分配独立网络和存储卷
**And** 部署过程在 3 分钟内完成

### Story 4.4: 部署状态追踪与进度反馈

As a 用户,
I want 实时查看部署进度,
So that 知道实例部署进展。

**Acceptance Criteria:**

**Given** 部署任务已创建
**When** 查询部署状态
**Then** 返回当前部署阶段（验证配置、创建实例、启动服务）
**And** 返回预计剩余时间
**And** 支持轮询或 WebSocket 实时更新
**And** 部署完成后状态变为 running
**And** 前端显示部署进度组件（进度条 + 阶段说明）

### Story 4.5: 部署失败处理与错误详情

As a 用户,
I want 在部署失败时看到详细的错误信息,
So that 可以了解原因并尝试修复。

**Acceptance Criteria:**

**Given** 部署过程中发生错误
**When** 部署失败
**Then** 状态更新为 error
**And** 记录详细错误日志
**And** 显示用户友好的错误提示
**And** 提供可能的解决方案（如"请检查飞书配置"）
**And** 提供"重新部署"按钮
**And** 支持查看完整错误日志

### Story 4.6: 用户实例列表页面

As a 用户,
I want 查看自己的实例列表和状态,
So that 管理我的 OpenClaw 实例。

**Acceptance Criteria:**

**Given** 用户已登录
**When** 访问实例管理页面
**Then** 显示用户的实例列表（卡片式布局）
**And** 每个实例卡片显示：名称、状态、创建时间
**And** 状态使用颜色标识（运行-绿色、停止-灰色、错误-红色）
**And** 提供"在飞书中使用"快捷入口
**And** 提供"查看详情"入口
**And** 无实例时显示空状态和"创建实例"按钮

### Story 4.7: 管理员实例管理功能

As a 平台管理员,
I want 管理所有用户的实例,
So that 可以处理问题实例。

**Acceptance Criteria:**

**Given** 管理员已登录管理后台
**When** 访问实例管理页面
**Then** 显示所有用户的实例列表
**And** 可以按用户、状态筛选实例
**And** 可以查看实例详情（租户信息、容器信息、日志）
**And** 提供重启/停止/启动操作按钮
**And** 操作需要二次确认
**And** 操作结果实时反馈

---

## Epic 5: 飞书机器人集成

用户可以在飞书中添加 OpenClaw 机器人，发送消息并接收回复。系统提供飞书机器人添加教程。实现飞书 SDK 长连接集成，确保消息实时收发。

### Story 5.1: 飞书 SDK 长连接集成

As a 后端开发者,
I want 集成飞书 Go SDK 长连接,
So that 可以实时接收飞书消息事件。

**Acceptance Criteria:**

**Given** 飞书应用配置已完成
**When** 启动飞书长连接服务
**Then** 成功连接飞书服务器
**And** 订阅消息接收事件
**And** 支持自动重连（断开后 30 秒内重连）
**And** 连接状态监控和日志记录
**And** 支持多租户消息路由

### Story 5.2: 消息接收与路由

As a 后端开发者,
I want 接收飞书消息并路由到正确的 OpenClaw 实例,
So that 用户的消息能被正确处理。

**Acceptance Criteria:**

**Given** 飞书长连接已建立
**When** 用户在飞书发送消息给机器人
**Then** 系统接收消息事件
**And** 根据 App ID 识别租户
**And** 查找租户对应的 OpenClaw 实例
**And** 将消息转发到 OpenClaw 实例
**And** 记录消息处理日志

### Story 5.3: 消息发送功能

As a 后端开发者,
I want 实现飞书消息发送功能,
So that 可以将 OpenClaw 的回复发送给用户。

**Acceptance Criteria:**

**Given** OpenClaw 实例生成了回复
**When** 调用消息发送接口
**Then** 使用飞书 API 发送消息给用户
**And** 消息格式正确（文本/卡片）
**And** 发送失败时自动重试（最多 3 次）
**And** 记录发送状态和日志

### Story 5.4: 飞书机器人添加教程页面

As a 用户,
I want 查看如何在飞书中添加机器人的教程,
So that 可以在飞书中开始使用 OpenClaw。

**Acceptance Criteria:**

**Given** 实例已部署成功
**When** 查看"如何在飞书中使用"教程
**Then** 显示分步教程（3 步）
**And** 步骤 1：在飞书搜索机器人名称
**And** 步骤 2：添加机器人到会话
**And** 步骤 3：发送测试消息验证
**And** 每步包含截图指引
**And** 提供"在飞书中打开"快捷按钮

---

## Epic 6: 数据备份与恢复

系统可以自动备份用户配置数据和 OpenClaw 实例数据。平台管理员可以查看备份状态。在实例重新部署时可以恢复用户配置。

### Story 6.1: 用户配置自动备份

As a 系统,
I want 自动备份用户配置数据,
So that 用户数据安全有保障。

**Acceptance Criteria:**

**Given** 用户配置数据已存储
**When** 配置发生变更或定时任务触发
**Then** 自动备份用户配置到备份存储
**And** 备份包含飞书配置、实例配置
**And** 使用加密存储备份文件
**And** 保留最近 7 天的备份
**And** 备份成功率 ≥ 99%

### Story 6.2: OpenClaw 实例数据备份

As a 系统,
I want 自动备份 OpenClaw 实例数据,
So that 实例数据不会丢失。

**Acceptance Criteria:**

**Given** OpenClaw 实例正在运行
**When** 定时任务触发（每日）
**Then** 调用 Dokploy S3 备份接口
**And** 备份实例的 Docker Volume 数据
**And** 备份文件存储到 S3 兼容存储
**And** 记录备份状态和日志
**And** 备份失败时发送告警

### Story 6.3: 备份状态查看 API

As a 平台管理员,
I want 查看系统备份状态,
So that 确认备份任务正常运行。

**Acceptance Criteria:**

**Given** 管理员已登录管理后台
**When** 访问备份管理页面
**Then** 显示最近备份任务列表
**And** 显示每个任务的：时间、类型、状态、大小
**And** 可以按日期范围筛选
**And** 可以查看备份详情和日志
**And** 显示备份存储空间使用情况

### Story 6.4: 用户配置恢复功能

As a 系统,
I want 在实例重新部署时恢复用户配置,
So that 用户无需重新配置。

**Acceptance Criteria:**

**Given** 用户需要重新部署实例
**When** 触发实例重新部署
**Then** 自动恢复用户的飞书配置
**And** 配置恢复到新实例
**And** 恢复失败时通知用户手动配置
**And** 记录恢复操作日志

---

## Epic 7: 平台监控与运维

平台管理员可以通过仪表盘查看用户总数、活跃实例数量、部署成功率、系统可用性等关键指标。可以查看实例部署日志，接收资源不足告警，并向用户发送系统通知。

### Story 7.1: 管理仪表盘数据 API

As a 后端开发者,
I want 实现管理仪表盘数据 API,
So that 前端可以展示平台运营数据。

**Acceptance Criteria:**

**Given** 系统已运行
**When** 调用 GET /v1/admin/dashboard API
**Then** 返回用户总数
**And** 返回活跃实例数量
**And** 返回部署成功率（最近 7 天）
**And** 返回系统可用性（最近 7 天）
**And** 返回最近告警列表
**And** API 响应时间 < 5 秒

### Story 7.2: 管理仪表盘前端页面

As a 平台管理员,
I want 在仪表盘查看平台运营概览,
So that 快速了解平台运行状况。

**Acceptance Criteria:**

**Given** 管理员已登录管理后台
**When** 访问仪表盘页面
**Then** 显示用户总数卡片（数字 + 环比变化）
**And** 显示活跃实例数卡片
**And** 显示部署成功率卡片（进度条）
**And** 显示系统可用性卡片
**And** 显示最近告警列表（可点击查看详情）
**And** 页面加载时间 < 2 秒

### Story 7.3: 实例日志查看功能

As a 平台管理员,
I want 查看实例的部署和运行日志,
So that 排查问题原因。

**Acceptance Criteria:**

**Given** 管理员在实例详情页
**When** 点击"查看日志"
**Then** 显示实例的部署日志
**And** 支持实时日志流（WebSocket）
**And** 支持按时间范围筛选
**And** 支持关键词搜索
**And** 支持下载日志文件

### Story 7.4: 告警系统

As a 系统,
I want 在资源不足时发送告警通知,
So that 管理员可以及时处理问题。

**Acceptance Criteria:**

**Given** 系统监控服务运行中
**When** 检测到资源不足（CPU > 80%、内存 > 80%、磁盘 > 80%）
**Then** 创建告警记录
**And** 发送告警通知给管理员
**And** 告警包含：类型、级别、时间、详情
**And** 支持告警确认和处理状态更新
**And** 告警通知支持邮件/钉钉/企业微信

### Story 7.5: 系统通知功能

As a 平台管理员,
I want 向用户发送系统通知,
So that 可以告知用户重要信息。

**Acceptance Criteria:**

**Given** 管理员在通知管理页面
**When** 创建并发送通知
**Then** 可以选择通知对象（全部用户/指定用户）
**And** 可以编辑通知标题和内容
**And** 用户在业务平台看到通知
**And** 支持通知已读/未读状态
**And** 记录通知发送日志

---

## Epic 8: 用户引导与帮助

系统提供飞书开放平台配置教程、飞书机器人添加教程，并在用户首次访问时显示引导说明。帮助非技术用户快速上手。

### Story 8.1: 飞书开放平台配置教程页面

As a 用户,
I want 查看详细的飞书开放平台配置教程,
So that 可以独立完成飞书应用创建。

**Acceptance Criteria:**

**Given** 用户访问帮助中心
**When** 点击"飞书开放平台配置教程"
**Then** 显示完整的图文教程
**And** 教程包含 5 个步骤
**And** 每步包含截图和文字说明
**And** 步骤 1：注册飞书开放平台账号
**And** 步骤 2：创建企业自建应用
**And** 步骤 3：配置应用权限
**And** 步骤 4：获取 App ID 和 Secret
**And** 步骤 5：返回平台填写配置
**And** 支持打印和分享

### Story 8.2: 飞书机器人添加教程页面

As a 用户,
I want 查看如何在飞书中添加 OpenClaw 机器人的教程,
So that 可以在飞书中开始使用。

**Acceptance Criteria:**

**Given** 用户访问帮助中心
**When** 点击"飞书机器人添加教程"
**Then** 显示完整的图文教程
**And** 教程包含 3 个步骤
**And** 步骤 1：在飞书搜索机器人
**And** 步骤 2：添加到单聊/群聊
**And** 步骤 3：发送消息测试
**And** 包含常见问题解答（FAQ）
**And** 提供"遇到问题"反馈入口

### Story 8.3: 首次访问引导

As a 新用户,
I want 在首次访问时看到引导说明,
So that 快速了解产品功能和使用流程。

**Acceptance Criteria:**

**Given** 用户首次登录系统
**When** 进入系统首页
**Then** 显示欢迎引导弹窗
**And** 引导包含：产品介绍、核心功能、使用流程
**And** 提供"开始使用"按钮跳转到配置页
**And** 提供"稍后再看"按钮关闭引导
**And** 用户可以选择"不再显示"
**And** 引导记录存储在用户偏好中

### Story 8.4: 帮助中心页面

As a 用户,
I want 访问帮助中心获取各类帮助,
So that 遇到问题时可以自助解决。

**Acceptance Criteria:**

**Given** 用户在任意页面
**When** 点击"帮助"入口
**Then** 显示帮助中心页面
**And** 包含快速入门指南
**And** 包含常见问题（FAQ）
**And** 包含视频教程链接
**And** 包含联系客服入口
**And** 支持关键词搜索帮助内容
