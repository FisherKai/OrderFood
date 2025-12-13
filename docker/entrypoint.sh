#!/bin/sh
set -e

# 生成配置文件（从环境变量）
cat > /app/config.yaml << EOF
# 服务器配置
server:
  port: 8080
  mode: ${SERVER_MODE:-release}

# 数据库配置
database:
  host: "${DB_HOST:-mysql}"
  port: ${DB_PORT:-3306}
  username: "${DB_USER:-root}"
  password: "${DB_PASSWORD:-}"
  database: "${DB_NAME:-orderfood_db}"
  charset: "utf8mb4"
  max_idle_conns: ${DB_MAX_IDLE:-10}
  max_open_conns: ${DB_MAX_OPEN:-100}

# JWT配置
jwt:
  secret: "${JWT_SECRET:-your-secret-key-change-in-production}"
  expire_hours: ${JWT_EXPIRE_HOURS:-168}

# 图片存储配置
storage:
  local_path: "/app/uploads"
  image_max_size: ${IMAGE_MAX_SIZE:-5242880}
  allowed_formats: ["jpg", "jpeg", "png", "gif", "webp"]
  thumbnail_width: 300
  thumbnail_height: 300

# 公告配置
announcement:
  auto_play_interval: 3000
  max_display_count: 10
EOF

echo "配置文件已生成:"
cat /app/config.yaml

# 等待数据库就绪
if [ -n "$DB_HOST" ]; then
    echo "等待数据库连接..."
    max_attempts=30
    attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if nc -z "$DB_HOST" "${DB_PORT:-3306}" 2>/dev/null; then
            echo "数据库已就绪"
            break
        fi
        attempt=$((attempt + 1))
        echo "等待数据库... ($attempt/$max_attempts)"
        sleep 2
    done
    
    if [ $attempt -eq $max_attempts ]; then
        echo "警告: 数据库连接超时，继续启动..."
    fi
fi

# 启动nginx（后台）
echo "启动Nginx..."
nginx

# 启动后端服务（前台）
echo "启动后端服务..."
exec /app/orderfood-server
