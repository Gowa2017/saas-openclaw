---
stepsCompleted: [1, 2, 3, 4, 5, 6]
inputDocuments: ["prd.md", "research/technical-openclaw-container-platform-research-2026-03-04.md"]
workflowType: 'architecture'
project_name: 'saas-openclaw'
user_name: 'Gowa'
date: '2026-03-04'
---

# Architecture Decision Document

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

---

## Project Context Analysis

### Requirements Overview

**Functional Requirements Summary:**

项目包含 **37 个功能需求**，覆盖以下领域：

1. **用户认证与授权** (FR1-FR4)
   - 业务平台登录 + 飞书 OAuth 授权
   - 平台管理员独立认证体系

2. **飞书应用配置** (FR5-FR9)
   - 飞书 App ID / Secret 配置管理
   - 配置验证与修改功能

3. **OpenClaw 实例管理** (FR10-FR18)
   - 用户一键启动实例（3 分钟内部署）
   - 部署状态追踪与错误通知
   - 管理员实例生命周期管理（重启/停止/启动）

4. **飞书机器人使用** (FR19-FR23)
   - 飞书机器人添加与配置
   - 实时消息收发

5. **数据管理** (FR24-FR27)
   - 用户配置数据备份
   - OpenClaw 实例数据备份与恢复

6. **监控与运维** (FR28-FR34)
   - 管理仪表盘（用户数、活跃实例、部署成功率）
   - 告警系统与日志查看
   - 系统通知推送

7. **用户引导** (FR35-FR37)
   - 飞书开放平台配置教程
   - 飞书机器人添加教程
   - 首次访问引导

**Non-Functional Requirements:**

| NFR 类别 | 关键要求 | 架构影响 |
|---------|---------|---------|
| **性能** | 实例部署 < 3 分钟、消息响应 < 3 秒、页面加载 < 2 秒 | 需要高效的容器编排和异步消息处理 |
| **安全** | 容器级隔离、配置加密存储、HTTPS 强制 | 多层隔离策略、加密存储方案 |
| **可扩展性** | 100-1000 用户无缝扩展、数据库 100 并发 | Docker Swarm 水平扩展、连接池优化 |
| **可靠性** | 可用性 ≥ 99%、部署成功率 ≥ 95% | 高可用架构、健康检查与自动恢复 |
| **集成** | 飞书 SDK 稳定性 ≥ 99%、API 响应 < 5 秒 | 长连接重连机制、超时与熔断策略 |

**Scale & Complexity:**

- **主要领域**: 后端基础设施 + 容器编排 + 第三方集成
- **复杂度级别**: **高**
- **预估架构组件数**: 8-12 个主要组件

### Technical Constraints & Dependencies

**已确定的技术栈约束:**

| 层级 | 技术选型 | 来源 |
|------|---------|------|
| PaaS 平台 | Dokploy (Go + TypeScript) | 技术研究报告 |
| 容器编排 | Docker + Docker Swarm | 技术研究报告 |
| 动态路由 | Traefik | 技术研究报告 |
| 多租户数据库 | PostgreSQL (Database per Tenant) | 技术研究报告 |
| 备份恢复 | Dokploy S3 + Velero | 技术研究报告 |

**外部依赖:**
- 飞书开放平台（Go SDK，OAuth 2.0，长连接事件监听）
- Dokploy API（容器部署与管理）

**部署约束:**
- 所有数据存储在国内服务器（数据本地化要求）
- 平台与 OpenClaw 实例分离（实例故障不影响平台运行）

### Cross-Cutting Concerns Identified

以下横切关注点将影响多个架构组件的设计：

1. **多租户安全隔离** ⚠️
   - 影响：容器网络、数据存储、API 访问控制
   - 策略：Docker 网络隔离 + PostgreSQL Database per Tenant + JWT 租户识别

2. **容器生命周期管理**
   - 影响：部署 API、监控模块、备份系统
   - 策略：健康检查、自动重启、状态追踪

3. **实时消息路由与分发**
   - 影响：飞书集成层、OpenClaw 实例通信
   - 策略：飞书长连接 + 消息队列 + 异步处理

4. **备份与灾难恢复**
   - 影响：数据层、容器编排层
   - 策略：定时自动备份 + 快速恢复机制

5. **可观测性**
   - 影响：所有组件
   - 策略：统一日志、指标追踪、告警聚合

---

## Starter Template Evaluation

### Primary Technology Domain

**全栈 SaaS 应用** - 基于路径 B（独立 SaaS 平台 + Dokploy 集成），需要前后端分离架构，后端使用 Go，前端使用 Vue 3 + TypeScript。

### Starter Options Considered

**后端 Starter 评估:**

| Starter | 技术栈 | 评估结果 |
|---------|--------|---------|
| **go-rest-api** | Go + Clean Architecture + SOLID | ✅ 推荐 - 领域驱动设计，高内聚低耦合 |
| **production-go-api-template** | Go + 最少依赖 + 生产就绪 | ⚠️ 备选 - SQLite 默认需改 PostgreSQL |

**前端 Starter 评估 (Vue 3):**

| Starter | 技术栈 | 评估结果 |
|---------|--------|---------|
| **Vite + Vue 3 + TypeScript** | 现代构建工具链 | ✅ 推荐 - 极快热重载，开箱即用 |
| **Naive UI** | Vue 3 原生组件库 | ✅ 推荐 - 100% TypeScript，类型安全 |
| **Element Plus** | Element UI Vue 3 版本 | ⚠️ 备选 - 成熟但 TS 支持较弱 |
| **Ant Design Vue** | 蚂蚁设计 Vue 版 | ⚠️ 备选 - 企业级但包体积大 |

### Selected Starter: 组合方案

**选择理由:**

1. **go-rest-api** 提供 Clean Architecture 基础，适合复杂业务逻辑和多租户隔离需求
2. **Vite + Vue 3 + TypeScript** 提供现代化前端开发体验，Vue 3 Composition API 更适合复杂状态管理
3. **Naive UI** 提供 100% TypeScript 编写的组件库，类型安全性最佳，90+ 组件全覆盖

**初始化命令:**

```bash
# 后端 - 使用 go-rest-api 模板
git clone https://github.com/username/go-rest-api.git backend
cd backend
# 修改数据库配置为 PostgreSQL

# 前端 - 创建 Vite + Vue 3 + TypeScript 项目
npm create vite@latest frontend -- --template vue-ts
cd frontend
npm install
# 安装 Naive UI
npm install naive-ui
```

**Note**: 项目初始化应作为第一个实现的 Story。

---

## Core Architectural Decisions

### Decision Priority Analysis

**Critical Decisions (阻塞实现):**
- 数据验证策略
- 认证与安全架构
- API 设计模式
- 状态管理
- CI/CD 管道

**Important Decisions (塑造架构):**
- 数据库多租户策略 (Database per Tenant - 已在技术研究确定)
- 日志与监控方案
- 备份恢复策略

**Deferred Decisions (MVP 后期):**
- 缓存策略 (可后期添加 Redis)
- 消息队列优化 (可后期添加 RabbitMQ/Kafka)

### Data Architecture

**决策: 前后端双验证**
- **后端**: Go validator (结构体标签)
- **前端**: Naive UI Form Rules
- **优势**: 双重保障，即时用户反馈

### Authentication & Security

**决策: 混合认证 - 飞书 OAuth (前端) + 平台 JWT (业务平台 → SaaS)**
- **认证流程**:
  1. 用户在业务平台登录 → 获得业务平台 JWT Token
  2. 业务平台调用 SaaS 平台时携带此 Token
  3. SaaS 平台验证 Token 后向业务平台请求用户信息
  4. 获取当前登录用户的租户、用户、授权等信息
- **影响**: 租户隔离、跨系统用户同步

### API & Communication Patterns

**决策: RESTful API**
- **原因**: 与 Dokploy API 集成更自然，飞书 SDK 集成更简单
- **设计原则**: 资源导向，标准 HTTP 方法，版本化 (/v1/)

### Frontend Architecture

**决策: Pinia 状态管理**
- **原因**: Vue 3 官方推荐，TypeScript 支持优秀
- **状态结构**: 按模块组织 (auth, tenant, instance, config)

### Infrastructure & Deployment

**决策: GitHub Actions CI/CD**
- **原因**: GitHub 原生集成，免费稳定，YAML 配置简单
- **流程**: 代码推送 → 自动构建 → 自动测试 → 自动部署

### Decision Impact Analysis

**Implementation Sequence:**
1. 初始化项目 (Starter Template)
2. 数据验证配置
3. 认证与安全实现
4. RESTful API 设计
5. 前端状态管理配置
6. CI/CD 管道搭建
7. Dokploy 集成
8. 多租户隔离实现
9. 飞书 SDK 集成
10. 备份恢复配置

**Cross-Component Dependencies:**
- 数据验证依赖认证 (用户信息来源)
- 认证影响 API 安全 (中间件层)
- 前端状态依赖 API 调用
- CI/CD 依赖所有组件
---

## Implementation Patterns & Consistency Rules

### Pattern Categories Defined

**Critical Conflict Points Identified:**
- 命名冲突（数据库、API、代码）
- 结构冲突（项目组织、文件结构）
- 格式冲突（API 响应、JSON 格式）
- 通信冲突（事件、状态、动作）
- 过程冲突（加载状态、错误处理）

### Naming Patterns

**Database Naming Conventions:**
- 表命名: snake_case (例: `tenant_users`, `openclaw_instances`)
- 列命名: PascalCase (例: `TenantID`, `CreatedAt`, `UpdatedAt`)
- 外键格式: table_id (例: `tenant_id`, `user_id`)

**API Naming Conventions:**
- REST 端点: 复数资源名 (例: `/tenants`, `/users`, `/instances`)
- 路由参数: {id} (例: `/tenants/:id`, `/instances/:id`)
- 查询参数: snake_case (例: `?tenant_id=xxx&status=active`)
- Header 命名: X-Custom-Header (例: `X-Platform-Token`, `X-Tenant-ID`)

**Code Naming Conventions:**
- 组件名: PascalCase (例: `TenantList`, `InstanceCard`, `ConfigForm`)
- 文件名: kebab-case (例: `tenant-list.vue`, `instance-card.vue`)
- 函数名: camelCase (例: `getTenantById`, `createInstance`, `updateConfig`)
- 变量名: camelCase (例: `tenantId`, `instanceStatus`, `configValue`)

### Structure Patterns

**Project Organization:**
- 按类型组织: `src/components`, `src/composables`, `src/utils`, `src/stores`, `src/pages`

**File Structure Patterns:**
- 配置文件: 根目录 (`.env`, `docker-compose.yml`, `vite.config.ts`)

### Format Patterns

**API Response Formats:**
- 统一包装器: `{ data: {...}, error: null, meta: {...} }`

### Communication Patterns

**Event System Patterns:**
- 事件命名: domain.event (例: `tenant.created`, `instance.deployed`)

### Process Patterns

**Error Handling Patterns:**
- 全局错误处理 + Vue 3 错误边界

### Enforcement Guidelines

**All AI Agents MUST:**
- 遵循上述命名约定（数据库表名 snake_case，列名 PascalCase）
- 遵循 API 命名约定（复数资源名，路由参数 {id}）
- 遵循代码命名约定（组件 PascalCase，文件 kebab-case，函数/变量 camelCase）
- 所有 API 响应使用统一包装器格式
- 事件命名使用 domain.event 风格

**Pattern Enforcement:**
- 代码审查时检查命名约定
- Lint 规则强制执行命名规范
---

## Project Structure & Boundaries

### Complete Project Directory Structure

```
saas-openclaw/
├── backend/                      # Go 后端服务
│   ├── cmd/
│   │   └── server/
│   │       └── main.go      # 应用入口
│   ├── internal/
│   │   ├── api/              # REST API 处理器
│   │   ├── domain/           # 领域模型和业务逻辑
│   │   │   ├── tenant/     # 租户领域
│   │   │   ├── user/       # 用户领域
│   │   │   ├── instance/   # 实例领域
│   │   │   └── config/    # 配置领域
│   │   ├── infrastructure/    # 基础设施
│   │   │   ├── config/    # 配置加载
│   │   │   ├── database/   # PostgreSQL 连接
│   │   │   └── dokploy/    # Dokploy 客户端
│   │   └── repository/      # 数据访问层
│   │       ├── tenant.go
│   │       ├── user.go
│   │       └── instance.go
│   ├── pkg/                   # 共享工具
│   │   ├── validator/       # 数据验证
│   │   ├── middleware/     # 中间件（认证、日志等）
│   │   └── logger/         # 日志工具
│   ├── go.mod
│   ├── go.sum
│   ├── .env.example
│   └── Dockerfile
│
├── frontend/                     # Vue 3 前端应用
│   ├── src/
│   │   ├── components/         # 按类型组织
│   │   │   ├── auth/         # 认证组件
│   │   │   ├── instances/    # 实例管理组件
│   │   │   ├── config/      # 配置组件
│   │   │   ├── dashboard/    # 仪表盘组件
│   │   │   └── ui/          # 通用 UI 组件
│   │   ├── composables/        # 组合式 API
│   │   │   ├── useAuth.ts
│   │   │   ├── useTenant.ts
│   │   │   └── useInstance.ts
│   │   ├── pages/            # 页面路由
│   │   │   ├── login/        # 登录页
│   │   │   ├── dashboard/    # 仪表盘
│   │   │   ├── feishu-config/ # 飞书配置
│   │   │   ├── instances/     # 实例管理
│   │   │   ├── feishu-bot/    # 飞书机器人
│   │   │   ├── backup/       # 备份恢复
│   │   │   ├── monitoring/   # 监控运维
│   │   │   └── onboarding/   # 用户引导
│   │   ├── stores/            # Pinia 状态管理
│   │   │   ├── auth.ts
│   │   │   ├── tenant.ts
│   │   │   ├── instance.ts
│   │   │   └── config.ts
│   │   ├── services/          # API 调用服务
│   │   │   ├── api.ts
│   │   │   ├── dokploy.ts
│   │   │   └── feishu.ts
│   │   ├── types/            # TypeScript 类型定义
│   │   │   ├── api.ts
│   │   │   ├── models.ts
│   │   │   └── index.ts
│   │   ├── utils/            # 工具函数
│   │   ├── router/           # 路由配置
│   │   └── assets/           # 静态资源
│   ├── public/                 # 公共静态文件
│   ├── index.html
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── tsconfig.node.json
│   ├── tailwind.config.js
│   ├── .env.example
│   ├── package.json
│   └── Dockerfile
│
├── .github/
│   └── workflows/
│       ├── backend-ci.yml     # 后端 CI/CD
│       └── frontend-ci.yml    # 前端 CI/CD
│
├── docker-compose.yml
└── README.md
```

### Architectural Boundaries

**API Boundaries:**
- 外部 API: `/v1/tenants`, `/v1/instances`, `/v1/users`
- 内部 API: 前端服务调用（业务平台用户信息接口）
- Dokploy API: Dokploy 容器管理接口

**Component Boundaries:**
- 前端组件通过 API 服务获取数据
- 组件间通过 Pinia 状态共享数据

**Service Boundaries:**
- 前端服务 (`src/services/`) 封装外部 API 调用
- 后端服务 (`internal/api/`) 处理 HTTP 请求和响应
- 仓储层 (`internal/repository/`) 处理数据库操作

### Requirements to Structure Mapping

| 功能领域 | 后端位置 | 前端位置 |
|---------|----------|----------|
| 认证授权 | `internal/api/auth/` | `src/components/auth/`, `src/stores/auth.ts` |
| 飞书配置 | `internal/domain/config/` | `src/pages/feishu-config/` |
| 实例管理 | `internal/api/instance/` | `src/pages/instances/` |
| 飞书机器人 | `internal/domain/instance/` | `src/pages/feishu-bot/` |
| 数据管理 | `internal/repository/` | `src/pages/backup/` |
| 监控运维 | `internal/api/monitoring/` | `src/pages/monitoring/` |
- PR 模板中检查模式一致性
