# 个人简历网站

## 项目介绍
这是一个基于Vue 3和Go的个人简历网站项目，包含前端展示和后台管理系统。前端使用Vue 3构建响应式界面，后端使用Go语言提供API服务，数据存储使用SQLite数据库。

## 技术栈
### 前端
- Vue 3
- Vue Router
- Axios
- GSAP (动画库)
- AOS (滚动动画库)
- Animate.css
- Vite (构建工具)

### 后端
- Go
- Gin Web框架
- JWT认证
- SQLite数据库
- bcrypt密码加密

## 功能特点
- 响应式设计，适配各种设备
- 前台访客密码验证访问
- 后台管理系统，包含以下功能：
  - 个人信息管理
  - 技能管理
  - 工作经历管理
  - 项目经验管理
  - 证书认证管理
  - 访客密码管理
- 数据库自动初始化
- JWT认证保护API

## 本地开发环境搭建

### 前提条件
- Node.js 16+
- Go 1.18+
- Git

### 前端开发
1. 克隆代码库
```bash
git clone [仓库地址]
cd [项目目录]
```

2. 安装依赖
```bash
npm install
```

3. 启动开发服务器
```bash
npm run dev
```

4. 构建生产版本
```bash
npm run build
```

### 后端开发
1. 进入后端目录
```bash
cd backend
```

2. 运行后端服务
```bash
go run main.go
```

3. 构建后端可执行文件
```bash
go build -o main
```

## Docker容器化部署

```run
# 创建并进入目录 - 前端
mkdir cv-portfolio && cd cv-portfolio
# 运行容器
docker run -p 8080:8080 -v $(pwd)/data:/app/data lcy0828/cv-portfolio:latest

```

### 创建Dockerfile
在项目根目录创建`Dockerfile`文件，内容如下：

```dockerfile
# 构建阶段 - 前端
FROM node:16-alpine AS frontend-builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# 构建阶段 - 后端
FROM golang:1.18-alpine AS backend-builder
WORKDIR /app
COPY backend/ ./
RUN go mod download
RUN go build -o main .

# 运行阶段
FROM alpine:latest
WORKDIR /app

# 创建存储数据的目录
RUN mkdir -p /app/data

# 复制前端构建产物
COPY --from=frontend-builder /app/dist /app/public

# 复制后端可执行文件
COPY --from=backend-builder /app/main /app/
COPY backend/data/ /app/data/

# 暴露端口
EXPOSE 8080

# 运行命令
CMD ["/app/main"]
```

### 创建docker-compose.yml
在项目根目录创建`docker-compose.yml`文件，内容如下：

```yaml
version: '3'

services:
  cv-portfolio:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    restart: unless-stopped
```

### 构建并运行容器
```bash
# 构建容器镜像
docker build -t cv-portfolio .

# 运行容器
docker run -p 8080:8080 -v $(pwd)/data:/app/data cv-portfolio

# 或者使用docker-compose
docker-compose up -d
```

## 访问应用
- 前端展示页面: `http://localhost:8080`
- 后台管理系统: `http://localhost:8080/admin`

## 默认账户
- 管理员用户名: `admin`
- 管理员密码: `admin123`
- 默认访客密码: `default_password`

## 注意事项
1. 数据存储在SQLite数据库中，位于`data/resume.db`
2. 容器化部署时，数据卷挂载用于持久化数据
3. 首次启动时，系统会自动初始化数据库并创建默认管理员账户和访客密码

## 开发调试
- 前端开发服务器运行在`http://localhost:5173`
- 后端API服务运行在`http://localhost:8080`
- 如需修改端口，可通过环境变量`PORT`进行设置

## 项目结构
```
.
├── backend/            # 后端Go代码
│   ├── database/       # 数据库相关代码
│   ├── handlers/       # API处理函数
│   ├── models/         # 数据模型
│   └── main.go         # 主程序入口
├── data/               # 数据存储目录
├── src/                # 前端Vue代码
│   ├── assets/         # 静态资源
│   ├── components/     # 组件
│   ├── router/         # 路由
│   ├── services/       # API服务
│   └── App.vue         # 主应用组件
├── Dockerfile          # Docker构建文件
├── docker-compose.yml  # Docker编排文件
└── package.json        # 前端依赖配置
``` 