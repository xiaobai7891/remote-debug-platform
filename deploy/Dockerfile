# 后端构建阶段
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app
COPY backend/ .

# 安装依赖
RUN go mod download

# 构建（启用 CGO 用于 SQLite）
RUN CGO_ENABLED=1 GOOS=linux go build -o server .

# 前端构建阶段
FROM node:20-alpine AS frontend-builder

WORKDIR /app
COPY frontend/ .

RUN npm install
RUN npm run build

# 最终镜像
FROM alpine:latest

WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache libc6-compat

# 复制后端
COPY --from=backend-builder /app/server .
COPY --from=backend-builder /app/config.yaml .

# 复制前端
COPY --from=frontend-builder /app/dist ./static

# 创建数据目录
RUN mkdir -p data

# 暴露端口
EXPOSE 3000
EXPOSE 10000-60000

# 启动命令
CMD ["./server"]
