# Story 1.3: PostgreSQL 数据库配置

Status: ready-for-dev

## Story

As a 开发者,
I want 配置 PostgreSQL 数据库连接,
so that 后端服务可以访问数据库存储数据。

## Acceptance Criteria

1. **AC1: PostgreSQL 数据库服务可访问**
   - **Given** PostgreSQL 数据库服务已启动
   - **When** 执行数据库连接测试
   - **Then** 能够成功连接 PostgreSQL 数据库
   - **And** 返回数据库版本信息

2. **AC2: 连接池配置完整**
   - **Given** 数据库连接已建立
   - **When** 检查连接池配置
   - **Then** 支持最大连接数配置（默认 100）
   - **And** 支持最大空闲连接数配置
   - **And** 支持连接最大存活时间配置
   - **And** 支持空闲连接最大存活时间配置

3. **AC3: 环境变量配置支持**
   - **Given** 项目已初始化
   - **When** 检查配置加载逻辑
   - **Then** 支持从 .env 文件读取数据库配置
   - **And** 支持环境变量覆盖配置
   - **And** 配置项包括：DB_HOST、DB_PORT、DB_USER、DB_PASSWORD、DB_NAME、DB_SSLMODE、DB_MAX_OPEN_CONNS、DB_MAX_IDLE_CONNS

4. **AC4: 数据库连接加密传输**
   - **Given** 生产环境配置
   - **When** 配置 DB_SSLMODE=require 或 verify-full
   - **Then** 数据库连接使用 SSL/TLS 加密传输
   - **And** 支持 SSL 证书验证配置

5. **AC5: 健康检查端点**
   - **Given** 后端服务已启动
   - **When** 访问 `/health/database` 端点
   - **Then** 返回数据库连接状态
   - **And** 返回连接池统计信息

## Tasks / Subtasks

- [ ] Task 1: 更新数据库配置结构 (AC: 2, 3)
  - [ ] 1.1 在 `config.go` 中添加 `ConnMaxLifetime` 字段
  - [ ] 1.2 在 `config.go` 中添加 `ConnMaxIdleTime` 字段
  - [ ] 1.3 更新 `.env.example` 中的连接池默认值（MaxOpenConns 改为 100）
  - [ ] 1.4 添加 `DB_CONN_MAX_LIFETIME` 和 `DB_CONN_MAX_IDLE_TIME` 配置项

- [ ] Task 2: 增强数据库连接模块 (AC: 1, 2, 4)
  - [ ] 2.1 更新 `database.go` 中的 `Connect` 函数，添加连接存活时间配置
  - [ ] 2.2 添加连接健康检查函数 `Ping()`
  - [ ] 2.3 添加连接池统计函数 `Stats()`
  - [ ] 2.4 添加连接关闭函数 `Close()`
  - [ ] 2.5 支持 SSL 连接配置

- [ ] Task 3: 实现数据库健康检查 API (AC: 5)
  - [ ] 3.1 创建 `internal/api/health/handler.go`
  - [ ] 3.2 实现 `/health/database` 端点
  - [ ] 3.3 返回 JSON 格式的健康状态和连接池信息
  - [ ] 3.4 在 `main.go` 中注册健康检查路由

- [ ] Task 4: 编写单元测试 (AC: 1, 2, 5)
  - [ ] 4.1 编写 `database_test.go` 测试连接功能
  - [ ] 4.2 编写 `config_test.go` 更新测试覆盖新配置项
  - [ ] 4.3 编写健康检查 handler 测试
  - [ ] 4.4 使用 mock 或 testcontainers 进行数据库连接测试

- [ ] Task 5: 集成验证 (AC: 1, 3, 5)
  - [ ] 5.1 启动本地 PostgreSQL 数据库（Docker）
  - [ ] 5.2 配置环境变量连接数据库
  - [ ] 5.3 运行后端服务并验证连接成功
  - [ ] 5.4 验证 `/health/database` 端点返回正确信息
  - [ ] 5.5 验证 SSL 连接配置

## Dev Notes

### 架构模式与约束

**必须遵循的 Clean Architecture 原则：**

1. **依赖方向**: 数据库连接属于基础设施层 (`internal/infrastructure/database/`)
2. **接口隔离**: 在 domain 层定义数据库接口，infrastructure 层实现
3. **配置管理**: 配置结构定义在 infrastructure/config

**关键架构决策 [Source: architecture.md]:**
- 多租户策略: PostgreSQL Database per Tenant
- 数据库命名: 表名 `snake_case`，列名 `PascalCase`
- 连接池: 支持高并发访问，最大连接数 100

### 现有代码状态

**Story 1.1 已完成的基础设施 [Source: 1-1-backend-project-init.md]:**

```
backend/
├── internal/infrastructure/
│   ├── config/
│   │   └── config.go      # 配置结构已定义，需扩展
│   └── database/
│       └── database.go    # 基础连接函数已存在，需增强
├── .env.example           # 已有数据库配置模板，需更新默认值
```

**现有 database.go 功能:**
- `Connect(cfg *config.DatabaseConfig)` - 建立数据库连接
- 设置 `MaxOpenConns` 和 `MaxIdleConns`
- 使用 sqlx 和 lib/pq 驱动

**需要增强的功能:**
- `ConnMaxLifetime` - 连接最大存活时间
- `ConnMaxIdleTime` - 空闲连接最大存活时间
- `Ping()` - 健康检查
- `Stats()` - 连接池统计
- `Close()` - 优雅关闭

### 技术栈要求

**核心依赖 [Source: architecture.md]:**

| 依赖 | 用途 | 当前版本 |
|------|------|---------|
| `github.com/lib/pq` | PostgreSQL 驱动 | v1.11.2 |
| `github.com/jmoiron/sqlx` | SQL 扩展 | v1.4.0 |
| `github.com/spf13/viper` | 配置管理 | v1.21.0 |

**连接池最佳实践:**

| 参数 | 建议值 | 说明 |
|------|-------|------|
| `MaxOpenConns` | 100 | 最大连接数，根据 AC 要求 |
| `MaxIdleConns` | 10 | 空闲连接数，约为最大连接数的 10% |
| `ConnMaxLifetime` | 30m | 连接最大存活时间，避免长时间连接 |
| `ConnMaxIdleTime` | 10m | 空闲连接最大存活时间 |

**SSL 模式说明:**

| 模式 | 说明 | 适用场景 |
|------|------|---------|
| `disable` | 不使用 SSL | 本地开发 |
| `require` | 使用 SSL，但不验证证书 | 测试环境 |
| `verify-full` | 使用 SSL 并验证证书 | 生产环境 |

### 项目结构规范

**新增文件位置:**

```
backend/
├── internal/
│   ├── api/
│   │   └── health/
│   │       └── handler.go    # 健康检查 handler（新增）
│   └── infrastructure/
│       ├── config/
│       │   └── config.go     # 更新：添加连接时间配置
│       └── database/
│           ├── database.go   # 更新：增强连接功能
│           └── database_test.go  # 新增：测试文件
```

### API 设计规范

**健康检查 API [Source: architecture.md]:**

```
GET /health/database

Response:
{
  "data": {
    "status": "healthy",
    "database": {
      "connected": true,
      "version": "PostgreSQL 15.x",
      "pool_stats": {
        "open_connections": 5,
        "in_use": 2,
        "idle": 3,
        "max_open": 100,
        "max_idle": 10
      }
    }
  },
  "error": null,
  "meta": {
    "timestamp": "2026-03-05T10:00:00Z"
  }
}
```

### 测试标准

**测试要求 [Source: 1-1-backend-project-init.md]:**
- 测试框架: Go 原生 testing 包
- 断言库: 可选 stretchr/testify
- 测试覆盖率目标: ≥ 70%
- 测试方式: 使用 Docker testcontainers 或 mock

**测试用例设计:**

| 测试场景 | 测试文件 | 测试方法 |
|---------|---------|---------|
| 配置加载 | `config_test.go` | 单元测试 |
| 连接建立 | `database_test.go` | 集成测试 (testcontainers) |
| 连接池配置 | `database_test.go` | 单元测试 |
| 健康检查 | `handler_test.go` | HTTP 测试 |

### Project Structure Notes

**与 Story 1.1 的连续性:**

Story 1.1 已创建数据库基础设施的基本框架，本 Story 需要在此基础上：

1. **复用现有代码**:
   - `config.go` 中的 `DatabaseConfig` 结构
   - `database.go` 中的 `Connect` 函数
   - `.env.example` 配置模板

2. **增量开发**:
   - 扩展配置结构
   - 增强连接功能
   - 添加健康检查 API

3. **质量标准**:
   - 保持 70%+ 测试覆盖率
   - 遵循命名约定
   - 保持代码风格一致

### 前一个 Story 的学习经验

**从 Story 1.1 (后端项目初始化) 获得的经验:**

1. **配置管理:**
   - 使用 Viper 的 `SetDefault` 设置默认值
   - 支持 `.env` 文件和环境变量
   - 配置结构体清晰分离

2. **测试策略:**
   - 使用 testify 断言库提高可读性
   - 测试覆盖率目标 70%
   - 每个模块有对应的 `_test.go` 文件

3. **代码质量:**
   - Code Review 发现配置默认值问题
   - 确保所有配置项都有合理的默认值
   - 错误处理要完善

### References

- [Source: architecture.md#Data Architecture] - 多租户 Database per Tenant 策略
- [Source: architecture.md#Technical Constraints] - PostgreSQL 技术栈约束
- [Source: architecture.md#Implementation Patterns] - 数据库命名约定
- [Source: architecture.md#API Response Formats] - 统一响应格式
- [Source: prd.md#NFR-SC3] - 数据库并发要求：100 并发请求
- [Source: epics.md#Story 1.3] - 原始故事定义
- [Source: 1-1-backend-project-init.md] - 前一个 Story 的实现经验和代码基础

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List

