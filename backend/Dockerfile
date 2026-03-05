# 后端 Dockerfile - 多阶段构建
# Story 1.5: Docker 容器化配置

# ============================================
# 阶段 1: 构建
# ============================================
FROM golang:1.25-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache git ca-certificates tzdata

# 配置 Go 模块代理（加速国内下载）
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 以利用 Docker 缓存
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 静态编译
# CGO_ENABLED=0: 禁用 CGO，生成静态链接的二进制文件
# GOOS=linux: 目标操作系统
# -a: 强制重新构建所有包
# -installsuffix cgo: 使用不同的安装后缀
# -ldflags: 链接器参数，去除调试信息减小体积
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o /app/server \
    ./cmd/server

# ============================================
# 阶段 2: 运行
# ============================================
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/server .

# 注意: 配置通过环境变量传递，不复制 .env 文件避免默认值覆盖

# 创建非 root 用户
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

# 更改文件所有权
RUN chown -R appuser:appgroup /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动服务
CMD ["./server"]
