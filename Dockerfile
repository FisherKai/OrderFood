# ==================== 阶段1: 构建前端 ====================
FROM node:18-alpine AS frontend-builder

WORKDIR /app

# 构建管理后台
COPY frontend-admin/package*.json ./frontend-admin/
RUN cd frontend-admin && npm install

COPY frontend-admin ./frontend-admin
RUN cd frontend-admin && npm run build

# 构建移动端
COPY frontend-mobile/package*.json ./frontend-mobile/
RUN cd frontend-mobile && npm install

COPY frontend-mobile ./frontend-mobile
RUN cd frontend-mobile && npm run build

# ==================== 阶段2: 构建后端 ====================
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制go mod文件
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 复制源码并构建
COPY backend .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o orderfood-server .

# ==================== 阶段3: 最终运行镜像 ====================
FROM alpine:3.19

WORKDIR /app

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata nginx netcat-openbsd

# 设置时区和字符集
ENV TZ=Asia/Shanghai
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8

# 从构建阶段复制后端二进制文件
COPY --from=backend-builder /app/orderfood-server /app/

# 从构建阶段复制前端静态文件
COPY --from=frontend-builder /app/frontend-admin/dist /app/static/admin
COPY --from=frontend-builder /app/frontend-mobile/dist /app/static/mobile

# 复制nginx配置
COPY docker/nginx.conf /etc/nginx/nginx.conf

# 复制启动脚本
COPY docker/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# 创建上传目录
RUN mkdir -p /app/uploads

# 暴露端口
EXPOSE 80 8080

# 启动脚本
ENTRYPOINT ["/app/entrypoint.sh"]
