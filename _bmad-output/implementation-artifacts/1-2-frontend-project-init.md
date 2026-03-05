# Story 1.2: 前端项目初始化

Status: done

## Story

As a 开发者,
I want 使用 Vite + Vue 3 + TypeScript + Naive UI 初始化前端项目,
so that 可以快速搭建现代化前端应用。

## Acceptance Criteria

1. **AC1: Node.js 环境准备**
   - **Given** 开发环境已安装 Node.js 18+
   - **When** 执行 `node -v` 和 `npm -v` 命令
   - **Then** 返回 Node.js 版本信息（≥ 18.0.0）
   - **And** 返回 npm 版本信息

2. **AC2: Vite 项目结构符合规范**
   - **Given** 执行前端项目初始化命令
   - **When** 项目初始化完成
   - **Then** 前端项目结构符合架构文档规范
   - **And** 包含 `index.html` 入口文件
   - **And** 包含 `src/main.ts` 入口文件
   - **And** 包含 `src/App.vue` 根组件
   - **And** 包含 `vite.config.ts` 配置文件
   - **And** 包含 `tsconfig.json` TypeScript 配置

3. **AC3: 核心技术栈集成完整**
   - **Given** 项目已初始化
   - **When** 检查依赖配置
   - **Then** 集成以下技术栈：
     - Vue 3.x（前端框架）
     - TypeScript 5.x（类型系统）
     - Vite 5.x（构建工具）
     - Naive UI（UI 组件库）
     - Pinia（状态管理）
     - Vue Router（路由管理）
     - Tailwind CSS（样式框架）

4. **AC4: 项目目录结构符合架构规范**
   - **Given** 项目已初始化
   - **When** 检查目录结构
   - **Then** 包含以下目录：
     - `src/components/` - 组件目录（按领域划分子目录）
     - `src/composables/` - 组合式 API
     - `src/pages/` - 页面路由
     - `src/stores/` - Pinia 状态管理
     - `src/services/` - API 调用服务
     - `src/types/` - TypeScript 类型定义
     - `src/utils/` - 工具函数
     - `src/router/` - 路由配置
     - `src/assets/` - 静态资源

5. **AC5: 项目可编译运行**
   - **Given** 项目已初始化
   - **When** 执行 `npm run dev` 和 `npm run build`
   - **Then** 开发服务器启动成功
   - **And** 生产构建成功无错误
   - **And** 页面可正常访问

## Tasks / Subtasks

- [x] Task 1: 创建前端项目基础结构 (AC: 2)
  - [x] 1.1 创建项目根目录 `frontend/`
  - [x] 1.2 使用 Vite 初始化 Vue 3 + TypeScript 项目
  - [x] 1.3 配置 `vite.config.ts`（路径别名、代理等）
  - [x] 1.4 配置 `tsconfig.json` 和 `tsconfig.node.json`
  - [x] 1.5 创建 `index.html` 入口文件

- [x] Task 2: 安装和配置核心依赖 (AC: 3)
  - [x] 2.1 安装 Naive UI 组件库
  - [x] 2.2 安装 Pinia 状态管理
  - [x] 2.3 安装 Vue Router
  - [x] 2.4 安装和配置 Tailwind CSS
  - [x] 2.5 配置 `package.json` 脚本命令

- [x] Task 3: 创建目录结构 (AC: 4)
  - [x] 3.1 创建 `src/components/` 目录（包含 auth、instances、config、dashboard、ui 子目录）
  - [x] 3.2 创建 `src/composables/` 目录
  - [x] 3.3 创建 `src/pages/` 目录（包含 login、dashboard、feishu-config、instances、feishu-bot、backup、monitoring、onboarding 子目录）
  - [x] 3.4 创建 `src/stores/` 目录
  - [x] 3.5 创建 `src/services/` 目录
  - [x] 3.6 创建 `src/types/` 目录
  - [x] 3.7 创建 `src/utils/` 目录
  - [x] 3.8 创建 `src/router/` 目录
  - [x] 3.9 创建 `src/assets/` 目录

- [x] Task 4: 实现基础入口文件 (AC: 5)
  - [x] 4.1 编写 `src/main.ts`（Vue 应用初始化、Pinia、Router、Naive UI 配置）
  - [x] 4.2 编写 `src/App.vue` 根组件（基础布局结构）
  - [x] 4.3 配置 Vue Router 基础路由
  - [x] 4.4 创建 Pinia store 基础结构

- [x] Task 5: 创建配置文件模板 (AC: 5)
  - [x] 5.1 创建 `.env.example` 文件
  - [x] 5.2 定义配置项：
    - `VITE_API_BASE_URL` - API 基础地址
    - `VITE_APP_TITLE` - 应用标题
  - [x] 5.3 创建 `tailwind.config.js` 配置
  - [x] 5.4 创建 `.gitignore` 文件

- [x] Task 6: 验证项目可运行 (AC: 5)
  - [x] 6.1 执行 `npm run dev` 启动开发服务器
  - [x] 6.2 执行 `npm run build` 构建生产版本
  - [x] 6.3 验证页面正常显示

## Dev Notes

### 架构模式与约束

**必须遵循的架构规范 [Source: architecture.md]:**

1. **技术栈要求:**
   - Vite 5.x（构建工具，极快热重载）
   - Vue 3 + Composition API（前端框架）
   - TypeScript 5.x（类型安全）
   - Naive UI（UI 组件库，100% TypeScript 编写）
   - Pinia（状态管理，Vue 3 官方推荐）
   - Vue Router（路由管理）
   - Tailwind CSS（样式框架）

2. **命名约定:**

| 类型 | 命名风格 | 示例 |
|------|---------|------|
| 组件名 | PascalCase | `TenantList`, `InstanceCard`, `ConfigForm` |
| 文件名 | kebab-case | `tenant-list.vue`, `instance-card.vue` |
| 函数名 | camelCase | `getTenantById`, `createInstance` |
| 变量名 | camelCase | `tenantId`, `instanceStatus` |
| 目录名 | kebab-case | `feishu-config/`, `user-guide/` |

3. **API 响应格式:**
   - 统一包装器: `{ data: {...}, error: null, meta: {...} }`

### 项目结构规范

**目录结构 [Source: architecture.md]:**

```
frontend/
├── src/
│   ├── components/         # 按类型组织
│   │   ├── auth/         # 认证组件
│   │   ├── instances/    # 实例管理组件
│   │   ├── config/      # 配置组件
│   │   ├── dashboard/    # 仪表盘组件
│   │   └── ui/          # 通用 UI 组件
│   ├── composables/        # 组合式 API
│   │   ├── useAuth.ts
│   │   ├── useTenant.ts
│   │   └── useInstance.ts
│   ├── pages/            # 页面路由
│   │   ├── login/        # 登录页
│   │   ├── dashboard/    # 仪表盘
│   │   ├── feishu-config/ # 飞书配置
│   │   ├── instances/     # 实例管理
│   │   ├── feishu-bot/    # 飞书机器人
│   │   ├── backup/       # 备份恢复
│   │   ├── monitoring/   # 监控运维
│   │   └── onboarding/   # 用户引导
│   ├── stores/            # Pinia 状态管理
│   │   ├── auth.ts
│   │   ├── tenant.ts
│   │   ├── instance.ts
│   │   └── config.ts
│   ├── services/          # API 调用服务
│   │   ├── api.ts
│   │   ├── dokploy.ts
│   │   └── feishu.ts
│   ├── types/            # TypeScript 类型定义
│   │   ├── api.ts
│   │   ├── models.ts
│   │   └── index.ts
│   ├── utils/            # 工具函数
│   ├── router/           # 路由配置
│   │   └── index.ts
│   ├── assets/           # 静态资源
│   ├── main.ts           # 应用入口
│   └── App.vue           # 根组件
├── public/                 # 公共静态文件
├── index.html
├── vite.config.ts
├── tsconfig.json
├── tsconfig.node.json
├── tailwind.config.js
├── .env.example
├── package.json
└── Dockerfile
```

### 技术栈要求

**核心依赖及版本 [Source: architecture.md]:**

| 依赖 | 用途 | 版本建议 |
|------|------|---------|
| `vue` | 前端框架 | ^3.4.x |
| `vue-router` | 路由管理 | ^4.2.x |
| `pinia` | 状态管理 | ^2.1.x |
| `naive-ui` | UI 组件库 | ^2.38.x |
| `tailwindcss` | 样式框架 | ^3.4.x |
| `typescript` | 类型系统 | ^5.3.x |
| `vite` | 构建工具 | ^5.1.x |
| `@vitejs/plugin-vue` | Vite Vue 插件 | ^5.0.x |

**Node.js 版本要求:** 18.x+（支持原生 ESM 和最新 JavaScript 特性）

### UX 设计规范

**设计系统选择 [Source: ux-design-specification.md]:**

Naive UI 已在架构中确定使用，以下是关键设计规范：

1. **色彩系统:**
   - 主色: 专业蓝 `#1677FF`
   - 成功色: 翠绿 `#52C41A`
   - 警告色: 橙黄 `#FAAD14`
   - 错误色: 红色 `#FF4D4F`

2. **字体层级:**
   - H1: 24px / 行高 32px / 字重 600
   - H2: 20px / 行高 28px / 字重 600
   - H3: 16px / 行高 24px / 字重 600
   - Body: 14px / 行高 22px / 字重 400
   - Caption: 12px / 行高 20px / 字重 400

3. **间距系统:**
   - xs: 4px
   - sm: 8px
   - md: 16px
   - lg: 24px
   - xl: 32px

4. **布局参数:**
   - 最大内容宽度: 1200px
   - 侧边导航宽度: 200px
   - 基础单位: 4px

### 测试标准

**测试要求 [Source: 参考 Story 1.1 后端测试标准]:**
- 测试框架: Vitest（Vite 原生测试框架）
- 测试覆盖率目标: ≥ 70%
- 组件测试: 使用 Vue Test Utils
- E2E 测试: 可选（Playwright）

### Project Structure Notes

**与后端项目对齐:**

前端项目应与后端项目保持一致的质量标准：
1. 配置文件模板（.env.example）
2. 清晰的目录结构
3. TypeScript 类型定义
4. 代码规范（ESLint + Prettier）

**检测到的差异及处理:**

| 方面 | 后端 (Go) | 前端 (Vue) | 说明 |
|------|----------|-----------|------|
| 包管理 | go.mod | package.json | 各自标准 |
| 配置格式 | .env | .env | 统一使用 .env |
| 测试文件 | *_test.go | *.spec.ts / *.test.ts | 各自标准 |

### 前一个 Story 的学习经验

**从 Story 1.1 (后端项目初始化) 获得的经验:**

1. **目录结构创建:**
   - 先创建所有目录，再添加文件
   - 使用占位符文件确保空目录被 Git 跟踪

2. **依赖管理:**
   - 使用最新稳定版本
   - 执行依赖整理命令（go mod tidy → npm install 后检查 package-lock.json）

3. **验证流程:**
   - 先验证编译/构建
   - 再验证运行
   - 最后验证功能

4. **文档清晰度:**
   - 包含明确的版本号
   - 提供具体的命令示例
   - 引用来源文档

### References

- [Source: architecture.md#Starter Template Evaluation] - Vite + Vue 3 + TypeScript + Naive UI 选择理由
- [Source: architecture.md#Project Structure & Boundaries] - 完整项目目录结构
- [Source: architecture.md#Implementation Patterns & Consistency Rules] - 命名约定和模式
- [Source: architecture.md#Frontend Architecture] - Pinia 状态管理决策
- [Source: ux-design-specification.md#Design System Foundation] - Naive UI 配置和定制策略
- [Source: ux-design-specification.md#Visual Design Foundation] - 色彩、字体、间距规范
- [Source: ux-design-specification.md#Responsive Design & Accessibility] - 响应式和无障碍要求
- [Source: prd.md#Technical Constraints] - 技术栈约束
- [Source: epics.md#Story 1.2] - 原始故事定义
- [Source: 1-1-backend-project-init.md] - 前一个 Story 的实现经验和模式

## Dev Agent Record

### Agent Model Used

qianfan-code-latest

### Debug Log References

- 修复 TypeScript 编译错误：auth.ts 中参数名与 ref 变量名冲突
- 修复 Tailwind CSS 版本：从 v4 降级到 v3.4.x 以匹配配置文件格式

### Completion Notes List

- ✅ 使用 Vite 脚手架创建 Vue 3 + TypeScript 项目
- ✅ 安装并配置 Naive UI、Pinia、Vue Router、Tailwind CSS
- ✅ 创建符合架构规范的完整目录结构
- ✅ 实现基础布局框架（侧边导航、顶部栏、内容区）
- ✅ 创建所有页面的基础组件
- ✅ 配置 Vue Router 路由懒加载
- ✅ 创建 Pinia stores（auth、tenant、instance、config）
- ✅ 创建 TypeScript 类型定义
- ✅ 创建 API 服务基础封装
- ✅ 验证开发服务器和生产构建均成功

### File List

**新增文件:**
- frontend/index.html
- frontend/package.json
- frontend/vite.config.ts
- frontend/tsconfig.json
- frontend/tsconfig.app.json
- frontend/tsconfig.node.json
- frontend/tailwind.config.js
- frontend/postcss.config.js
- frontend/.env.example
- frontend/.gitignore
- frontend/src/main.ts
- frontend/src/App.vue
- frontend/src/style.css
- frontend/src/router/index.ts
- frontend/src/stores/auth.ts
- frontend/src/stores/tenant.ts
- frontend/src/stores/instance.ts
- frontend/src/stores/config.ts
- frontend/src/types/api.ts
- frontend/src/types/models.ts
- frontend/src/types/index.ts
- frontend/src/services/api.ts
- frontend/src/pages/login/index.vue
- frontend/src/pages/dashboard/index.vue
- frontend/src/pages/feishu-config/index.vue
- frontend/src/pages/instances/index.vue
- frontend/src/pages/feishu-bot/index.vue
- frontend/src/pages/backup/index.vue
- frontend/src/pages/monitoring/index.vue
- frontend/src/pages/onboarding/index.vue

**修改文件:**
- 无（全新项目）

## Senior Developer Review (AI)

### Review Date: 2026-03-05

### Reviewer: Gowa (via AI Code Review Agent)

### Issues Found & Fixed

#### 🔴 CRITICAL Issues (4)

1. **App.vue 布局缺陷** - 登录页面显示完整布局
   - 修复：添加条件渲染，登录页不显示侧边栏和顶部导航

2. **缺少路由守卫** - 无认证保护
   - 修复：添加 `beforeEach` 守卫验证认证状态

3. **登录功能未实现** - 仅有 TODO 占位符
   - 修复：实现基础登录逻辑，添加表单验证

4. **类型定义重复** - 违反 DRY 原则
   - 修复：stores 中的类型定义改为从 `types/models.ts` 导入

#### 🟡 MEDIUM Issues (7)

5. **缺少测试配置** - 添加 vitest.config.ts 和测试依赖
6. **缺少 ESLint/Prettier** - 添加 .eslintrc.cjs 和 .prettierrc
7. **缺少 Dockerfile** - 添加 Dockerfile 和 nginx.conf
8. **Token 未持久化** - 添加 localStorage 持久化
9. **API 缺少超时设置** - 添加 30 秒超时机制
10. **缺少 dokploy.ts 和 feishu.ts 服务** - 已添加
11. **缺少 composables** - 添加 useAuth.ts、useTenant.ts、useInstance.ts

### Files Added During Review

- frontend/vitest.config.ts
- frontend/.eslintrc.cjs
- frontend/.prettierrc
- frontend/Dockerfile
- frontend/nginx.conf
- frontend/src/services/dokploy.ts
- frontend/src/services/feishu.ts
- frontend/src/composables/useAuth.ts
- frontend/src/composables/useTenant.ts
- frontend/src/composables/useInstance.ts

### Files Modified During Review

- frontend/src/App.vue - 条件布局渲染
- frontend/src/router/index.ts - 添加路由守卫
- frontend/src/pages/login/index.vue - 实现登录逻辑
- frontend/src/stores/auth.ts - Token 持久化
- frontend/src/stores/tenant.ts - 移除重复类型
- frontend/src/stores/instance.ts - 移除重复类型
- frontend/src/stores/config.ts - 移除重复类型
- frontend/src/services/api.ts - 添加超时机制
- frontend/package.json - 添加测试和 lint 依赖

### Files Deleted During Review

- frontend/src/components/HelloWorld.vue - 未使用的示例组件

### Build Verification

- ✅ TypeScript 类型检查通过
- ✅ 生产构建成功
- ⚠️ 构建 chunk 1.4MB（建议后续优化 manualChunks）

### Review Outcome: APPROVED

所有 CRITICAL 和 MEDIUM 问题已修复，项目可进入 `done` 状态。
