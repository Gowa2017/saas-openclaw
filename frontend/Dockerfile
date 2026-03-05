# 构建阶段
FROM node:20-alpine AS builder

WORKDIR /app

# 配置 npm 使用腾讯镜像，增加超时时间
RUN npm config set registry https://mirrors.cloud.tencent.com/npm/ && \
    npm config set fetch-timeout 600000 && \
    npm config set fetch-retry-mintimeout 20000 && \
    npm config set fetch-retry-maxtimeout 120000

# 复制 package 文件
COPY package*.json ./

# 安装依赖
RUN npm ci --prefer-offline

# 复制源代码
COPY . .

# 构建生产版本
RUN npm run build

# 生产阶段
FROM nginx:alpine

# 复制 nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 从构建阶段复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 暴露端口
EXPOSE 80

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:80/ || exit 1

# 启动 nginx
CMD ["nginx", "-g", "daemon off;"]
