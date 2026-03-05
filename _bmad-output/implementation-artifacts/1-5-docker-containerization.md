# Story 1.5: Docker 容器化配置

Status: done

## Story

As a 开发者,
I want 为前后端服务配置 Dockerfile,
so that 服务可以通过 Docker 容器部署。

## Acceptance Criteria

1. **AC1: 后端 Dockerfile 多阶段构建**
   - **Given** 后端项目已初始化
   - **When** 编写后端 Dockerfile 文件
   - **Then** 使用多阶段构建优化镜像大小
   - **And** 构建阶段使用 golang:alpine
   - **And** 运行阶段使用 alpine 基础镜像
   - **And** 最终镜像大小 < 50MB (实际: 38.6MB)

2. **AC2: 前端 Dockerfile 配置**
   - **Given** 前端项目已初始化
   - **When** 编写前端 Dockerfile 文件
   - **Then** 使用 nginx 静态服务
   - **And** 多阶段构建优化镜像大小
   - **And** nginx 配置支持 Vue Router History 模式

3. **AC3: docker-compose.yml 本地开发环境**
   - **Given** 前后端 Dockerfile 已配置
   - **When** 创建 docker-compose.yml
   - **Then** 配置后端服务 (backend)
   - **And** 配置前端服务 (frontend)
   - **And** 配置 PostgreSQL 数据库服务
   - **And** 配置服务间网络通信
   - **And** 配置数据卷持久化

4. **AC4: .dockerignore 文件配置**
   - **Given** 前后端项目已初始化
   - **When** 创建 .dockerignore 文件
   - **Then** 后端排除 go.sum 缓存、测试文件等
   - **And** 前端排除 node_modules、构建产物等
   - **And** 减少构建上下文大小

5. **AC5: 镜像构建验证**
   - **Given** 所有 Dockerfile 配置完成
   - **When** 执行 docker build 命令
   - **Then** 后端镜像构建成功
   - **And** 前端镜像构建成功
   - **And** docker-compose up 启动成功

## Tasks / Subtasks

- [x] Task 1: 创建后端 Dockerfile (AC: 1)
  - [x] 1.1 创建 `backend/Dockerfile` 多阶段构建文件
  - [x] 1.2 配置构建阶段 (golang:alpine)
  - [x] 1.3 配置运行阶段 (alpine)
  - [x] 1.4 配置 CGO_ENABLED=0 静态编译
  - [x] 1.5 暴露端口 8080

- [x] Task 2: 创建/验证前端 Dockerfile (AC: 2)
  - [x] 2.1 验证现有 `frontend/Dockerfile` 符合要求
  - [x] 2.2 验证 `frontend/nginx.conf` 配置正确
  - [x] 2.3 更新 Dockerfile 配置（添加 npm 镜像）

- [x] Task 3: 创建 docker-compose.yml (AC: 3)
  - [x] 3.1 在项目根目录创建 `docker-compose.yml`
  - [x] 3.2 配置 backend 服务（依赖 postgres）
  - [x] 3.3 配置 frontend 服务（依赖 backend）
  - [x] 3.4 配置 postgres 服务（数据卷持久化）
  - [x] 3.5 配置网络隔离
  - [x] 3.6 配置环境变量传递

- [x] Task 4: 创建 .dockerignore 文件 (AC: 4)
  - [x] 4.1 创建 `backend/.dockerignore`
  - [x] 4.2 创建 `frontend/.dockerignore`
  - [x] 4.3 根目录 `.dockerignore`（不需要，各服务独立配置）

- [x] Task 5: 验证构建和运行 (AC: 5)
  - [x] 5.1 验证后端镜像构建 `docker build -t backend ./backend`
  - [x] 5.2 验证前端镜像构建 `docker build -t frontend ./frontend`
  - [x] 5.3 验证 docker-compose 启动 `docker-compose up -d`
  - [x] 5.4 验证服务健康检查
  - [x] 5.5 验证服务间通信

## Dev Notes

### 架构模式与约束

**必须遵循的 Docker 配置原则：**

1. **多阶段构建**: 减少最终镜像大小，分离构建和运行环境
2. **最小化基础镜像**: 使用 alpine 镜像减少攻击面
3. **安全最佳实践**: 不在镜像中硬编码敏感信息
4. **层优化**: 合并 RUN 命令减少层数

**关键架构决策 [Source: architecture.md]:**
- 容器编排: **Docker + Docker Swarm**
- PaaS 平台: **Dokploy (Go + TypeScript)**
- 多租户数据库: **PostgreSQL (Database per Tenant)**
- 部署约束: 所有数据存储在国内服务器

### 现有项目状态

**后端项目 [Source: 1-1-backend-project-init.md]:**

```
backend/
├── cmd/server/main.go      # 入口文件
├── internal/               # Clean Architecture 结构
│   ├── api/               # REST API 处理器
│   ├── domain/            # 领域模型
│   ├── infrastructure/    # 基础设施
│   └── repository/        # 数据访问层
├── pkg/                    # 共享工具
├── go.mod                  # Go 1.25.0
├── go.sum
├── .env.example
└── .gitignore
```

**关键依赖版本 [Source: backend/go.mod]:**
- Go: 1.25.0
- Gin: v1.9.1
- sqlx: v1.3.5
- Viper: v1.18.2
- Zap: v1.26.0

**前端项目 [Source: 1-2-frontend-project-init.md]:**

```
frontend/
├── src/                    # Vue 3 + TypeScript
├── public/
├── package.json            # npm
├── vite.config.ts          # Vite 构建工具
├── tsconfig.json           # TypeScript 配置
├── Dockerfile              # ✅ 已存在
├── nginx.conf              # ✅ 已存在
└── .env.example
```

**关键依赖版本:**
- Node.js: 20 LTS
- Vue: 3.x
- TypeScript: 5.x
- Vite: 5.x
- Naive UI: 最新

**现有前端 Dockerfile 分析 [Source: frontend/Dockerfile]:**
- ✅ 已使用多阶段构建
- ✅ 构建阶段: node:20-alpine
- ✅ 生产阶段: nginx:alpine
- ✅ nginx 配置支持 Vue Router History 模式
- ✅ 配置了 API 代理 (proxy_pass http://backend:8080)
- ✅ 静态资源缓存配置

### 技术栈要求

**Docker 镜像版本选择:**

| 组件 | 基础镜像 | 原因 |
|------|---------|------|
| Go 构建 | golang:alpine | 官方 Alpine 镜像，体积小 |
| Go 运行 | alpine:latest | 最小化运行环境 |
| Node 构建 | node:20-alpine | LTS 版本，Alpine 体积小 |
| 前端运行 | nginx:alpine | 静态服务，配置简单 |
| PostgreSQL | postgres:16-alpine | 官方 Alpine 镜像 |

### 项目结构规范

**新增文件位置:**

```
saas-openclaw/
├── backend/
│   ├── Dockerfile           # 后端 Dockerfile（新增）
│   └── .dockerignore        # 后端 .dockerignore（新增）
├── frontend/
│   ├── Dockerfile           # 前端 Dockerfile（已存在 ✅）
│   ├── nginx.conf           # nginx 配置（已存在 ✅）
│   └── .dockerignore        # 前端 .dockerignore（新增）
├── docker-compose.yml       # 本地开发环境（新增）
└── .dockerignore            # 根目录 .dockerignore（可选）
```

### 常见问题与解决方案

**问题 1: 后端镜像体积过大**
- **原因**: 使用 golang 完整镜像，包含编译工具
- **解决**: 使用多阶段构建，只保留编译后的二进制文件

**问题 2: 前端容器内 API 调用失败**
- **原因**: nginx 代理配置错误
- **解决**: 确保 nginx.conf 中 `proxy_pass` 指向正确的后端服务名

**问题 3: 数据库连接失败**
- **原因**: 后端启动时数据库未就绪
- **解决**: 使用 `depends_on` + `healthcheck` 确保启动顺序

**问题 4: 环境变量未传递**
- **原因**: docker-compose.yml 未正确配置环境变量
- **解决**: 使用 `environment` 或 `env_file` 传递配置

**问题 5: 国内网络构建慢**
- **原因**: 访问国外镜像源超时
- **解决**: 配置 Go 模块代理 (goproxy.cn) 和 npm 镜像 (mirrors.cloud.tencent.com)

### 安全注意事项

1. **不要在 Dockerfile 中硬编码敏感信息**
2. **使用 .dockerignore 排除 .env 文件**
3. **生产环境使用 Docker Secrets 管理密钥**
4. **定期更新基础镜像版本**

### References

- [Source: architecture.md#Infrastructure & Deployment] - Docker 容器编排决策
- [Source: architecture.md#Project Structure & Boundaries] - 项目目录结构
- [Source: epics.md#Story 1.5] - 原始故事定义
- [Source: prd.md#NFR-S1] - 用户间容器完全隔离
- [Source: prd.md#NFR-SC2] - 新容器部署 5 分钟内完成资源分配
- [Source: 1-1-backend-project-init.md] - 后端项目实现经验和结构
- [Source: 1-2-frontend-project-init.md] - 前端项目实现经验
- [Source: 1-3-postgresql-database-config.md] - 数据库配置相关上下文
- [Source: 1-4-github-actions-cicd.md] - CI/CD 配置相关上下文
- [Source: frontend/Dockerfile] - 现有前端 Dockerfile
- [Source: frontend/nginx.conf] - 现有 nginx 配置

## Dev Agent Record

### Agent Model Used

qianfan-code-latest

### Debug Log References

- 网络问题：国内访问国外镜像源超时，已配置 goproxy.cn 和腾讯 npm 镜像
- 环境变量问题：Viper 需要显式 BindEnv 才能正确读取环境变量
- 本地进程干扰：本地运行的 main 进程占用 8080 端口，干扰测试

### Completion Notes List

1. **后端 Dockerfile**: 创建了多阶段构建的 Dockerfile，使用 golang:1.24-alpine 构建，最终镜像大小 38.6MB
2. **前端 Dockerfile**: 验证并更新了现有配置，添加了 npm 镜像和 HEALTHCHECK
3. **docker-compose.yml**: 创建了完整的本地开发环境配置，使用 env_file 替代硬编码密码
4. **.dockerignore**: 为前后端创建了 .dockerignore 文件，减少构建上下文
5. **网络优化**: 配置了 Go 模块代理和 npm 镜像以解决国内网络问题
6. **配置修复**: 修复了 Viper 环境变量绑定问题，确保配置正确传递
7. **安全改进**: 创建 .env.example，移除 docker-compose.yml 中的硬编码密码

### Change Log

| 日期 | 变更 | 作者 |
|------|------|------|
| 2026-03-05 | 初始实现：创建 Dockerfile、docker-compose.yml、.dockerignore | Dev Agent |
| 2026-03-05 | Code Review 修复：Go 版本固定、HEALTHCHECK、env_file、代码清理 | Code Review Agent |

### File List

**新增文件:**
- backend/Dockerfile
- backend/.dockerignore
- frontend/.dockerignore
- docker-compose.yml
- .env.example

**修改文件:**
- frontend/Dockerfile (添加 npm 镜像配置、HEALTHCHECK)
- backend/internal/infrastructure/config/config.go (添加 BindEnv 调用，重构为循环)
