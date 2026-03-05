# Story 1.5: Docker 容器化配置

Status: ready-for-dev

## Story

As a 开发者,
I want 为前后端服务配置 Dockerfile,
so that 服务可以通过 Docker 容器部署。

## Acceptance Criteria

1. **AC1: 后端 Dockerfile 多阶段构建**
   - **Given** 后端项目已初始化
   - **When** 编写后端 Dockerfile 文件
   - **Then** 使用多阶段构建优化镜像大小
   - **And** 构建阶段使用 golang:1.25-alpine
   - **And** 运行阶段使用 alpine 基础镜像
   - **And** 最终镜像大小 < 50MB

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

- [ ] Task 1: 创建后端 Dockerfile (AC: 1)
  - [ ] 1.1 创建 `backend/Dockerfile` 多阶段构建文件
  - [ ] 1.2 配置构建阶段 (golang:1.25-alpine)
  - [ ] 1.3 配置运行阶段 (alpine)
  - [ ] 1.4 配置 CGO_ENABLED=0 静态编译
  - [ ] 1.5 暴露端口 8080

- [ ] Task 2: 创建/验证前端 Dockerfile (AC: 2)
  - [ ] 2.1 验证现有 `frontend/Dockerfile` 符合要求
  - [ ] 2.2 验证 `frontend/nginx.conf` 配置正确
  - [ ] 2.3 如需修改，更新 Dockerfile 配置

- [ ] Task 3: 创建 docker-compose.yml (AC: 3)
  - [ ] 3.1 在项目根目录创建 `docker-compose.yml`
  - [ ] 3.2 配置 backend 服务（依赖 postgres）
  - [ ] 3.3 配置 frontend 服务（依赖 backend）
  - [ ] 3.4 配置 postgres 服务（数据卷持久化）
  - [ ] 3.5 配置网络隔离
  - [ ] 3.6 配置环境变量传递

- [ ] Task 4: 创建 .dockerignore 文件 (AC: 4)
  - [ ] 4.1 创建 `backend/.dockerignore`
  - [ ] 4.2 创建 `frontend/.dockerignore`
  - [ ] 4.3 创建根目录 `.dockerignore`（如需要）

- [ ] Task 5: 验证构建和运行 (AC: 5)
  - [ ] 5.1 验证后端镜像构建 `docker build -t backend ./backend`
  - [ ] 5.2 验证前端镜像构建 `docker build -t frontend ./frontend`
  - [ ] 5.3 验证 docker-compose 启动 `docker-compose up -d`
  - [ ] 5.4 验证服务健康检查
  - [ ] 5.5 验证服务间通信

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
| Go 构建 | golang:1.25-alpine | 官方 Alpine 镜像，体积小 |
| Go 运行 | alpine:latest | 最小化运行环境 |
| Node 构建 | node:20-alpine | LTS 版本，Alpine 体积小 |
| 前端运行 | nginx:alpine | 静态服务，配置简单 |
| PostgreSQL | postgres:16-alpine | 官方 Alpine 镜像 |

**多阶段构建优化策略:**

```dockerfile
# 后端 Dockerfile 模板
# 阶段 1: 构建
FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# 阶段 2: 运行
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env
EXPOSE 8080
CMD ["./main"]
```

### docker-compose.yml 设计

**服务架构:**

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: openclaw
      POSTGRES_USER: openclaw
      POSTGRES_PASSWORD: openclaw_dev
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U openclaw"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - backend-network

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=openclaw
      - DB_USER=openclaw
      - DB_PASSWORD=openclaw_dev
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - backend-network
      - frontend-network

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - frontend-network

volumes:
  postgres_data:

networks:
  backend-network:
    driver: bridge
  frontend-network:
    driver: bridge
```

### .dockerignore 配置

**后端 .dockerignore:**
```
# Git
.git
.gitignore

# Documentation
*.md
README.md

# Test files
*_test.go
coverage.out

# Development files
.vscode
.idea

# Binary
backend
main

# Environment
.env
.env.local
```

**前端 .dockerignore:**
```
# Dependencies
node_modules

# Build output
dist

# Git
.git
.gitignore

# Documentation
*.md
README.md

# Development files
.vscode
.idea
*.log

# Test files
coverage
.nyc_output

# Environment
.env
.env.local
```

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

### 测试标准

**Docker 构建验证清单:**

| 验证项 | 命令 | 预期结果 |
|--------|------|---------|
| 后端镜像构建 | `docker build -t backend ./backend` | 构建成功，镜像 < 50MB |
| 前端镜像构建 | `docker build -t frontend ./frontend` | 构建成功 |
| Compose 启动 | `docker-compose up -d` | 所有服务启动成功 |
| 后端健康检查 | `curl http://localhost:8080/health` | 返回 200 OK |
| 前端访问 | `curl http://localhost:80` | 返回 HTML 页面 |
| 数据库连接 | `docker exec -it <postgres> psql -U openclaw` | 连接成功 |

### Project Structure Notes

**与 Story 1.1、1.2、1.3、1.4 的连续性:**

1. **复用现有项目结构**:
   - 后端: `backend/` 目录，Go 1.25.0，Clean Architecture
   - 前端: `frontend/` 目录，Node 20，Vue 3
   - CI/CD: `.github/workflows/` 已配置

2. **增量开发**:
   - 创建后端 Dockerfile
   - 创建 docker-compose.yml
   - 创建 .dockerignore 文件

3. **与 CI/CD 集成**:
   - Docker 镜像构建可与 GitHub Actions 集成
   - 为后续部署到 Dokploy 做准备

### 前序 Story 的学习经验

**从 Story 1.1 (后端项目初始化) 获得的经验:**

1. **配置管理**: 使用 Viper 管理配置，支持环境变量覆盖
2. **健康检查**: `/health` 端点已实现，可用于 Docker healthcheck
3. **测试覆盖率**: 94.7% 测试覆盖率

**从 Story 1.2 (前端项目初始化) 获得的经验:**

1. **构建工具**: Vite 提供快速构建，输出到 `dist/` 目录
2. **Nginx 配置**: 已配置好 API 代理和静态资源缓存
3. **Dockerfile**: 已实现多阶段构建

**从 Story 1.3 (PostgreSQL 数据库配置) 获得的经验:**

1. **数据库配置**: 连接池已配置，支持环境变量
2. **.env.example**: 已有数据库配置模板
3. **健康检查**: 数据库健康检查端点已实现

**从 Story 1.4 (GitHub Actions CI/CD) 获得的经验:**

1. **CI 流程**: 后端和前端 CI 已配置
2. **构建产物**: 前端 dist 已上传到 Artifacts
3. **测试自动化**: CI 中自动运行测试

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

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
