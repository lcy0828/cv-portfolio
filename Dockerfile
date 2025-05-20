# 构建阶段 - 前端
FROM node:16-alpine AS frontend-builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# 构建阶段 - 后端
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app

# 安装构建 SQLite3 所需的依赖
RUN apk add --no-cache gcc musl-dev

# 启用 CGO
ENV CGO_ENABLED=1

COPY backend/ ./
RUN go mod download
RUN go build -o main .

# 运行阶段
FROM alpine:latest
WORKDIR /app

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata sqlite

# 创建存储数据的目录
RUN mkdir -p /app/data

# 复制前端构建产物
COPY --from=frontend-builder /app/dist /app/public

# 复制后端可执行文件
COPY --from=backend-builder /app/main /app/

# 设置工作目录为应用程序根目录
WORKDIR /app

# 设置环境变量
ENV GIN_MODE=release
ENV PORT=8080

# 暴露端口
EXPOSE 8080

# 运行命令
CMD ["/app/main"] 