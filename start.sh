#!/bin/bash

BASE_DIR="$(cd "$(dirname "$0")" && pwd)"
NGINX_CONF="$BASE_DIR/deploy-nginx.conf"
BACKEND_BIN="$BASE_DIR/backend/orderfood-server"
BACKEND_PID_FILE="/tmp/orderfood-backend.pid"

echo "🍽️  点餐系统部署脚本"
echo "===================="
echo ""

# ==================== 环境检查 ====================

check_env() {
    local has_error=0

    if ! command -v go &> /dev/null; then
        echo "❌ 未检测到 Go，请先安装 Go 1.21+"
        has_error=1
    fi

    if ! command -v node &> /dev/null; then
        echo "❌ 未检测到 Node.js，请先安装 Node.js 16+"
        has_error=1
    fi

    if ! command -v nginx &> /dev/null; then
        echo "❌ 未检测到 Nginx，请先安装 Nginx"
        echo "   Ubuntu/Debian: sudo apt install nginx"
        echo "   CentOS/RHEL:   sudo yum install nginx"
        has_error=1
    fi

    if [ $has_error -eq 1 ]; then
        exit 1
    fi

    echo "✅ 环境检查通过 (Go: $(go version | awk '{print $3}'), Node: $(node -v), Nginx)"
    echo ""
}

# ==================== 构建前端 ====================

build_mobile() {
    echo "📱 构建移动端前端..."
    cd "$BASE_DIR/frontend-mobile"
    echo "  📦 安装依赖..."
    npm install || { echo "❌ 移动端依赖安装失败"; return 1; }
    npm run build || { echo "❌ 移动端构建失败"; return 1; }
    echo "✅ 移动端构建完成 → frontend-mobile/dist"
}

build_admin() {
    echo "💻 构建管理后台前端..."
    cd "$BASE_DIR/frontend-admin"
    echo "  📦 安装依赖..."
    npm install || { echo "❌ 管理后台依赖安装失败"; return 1; }
    npm run build || { echo "❌ 管理后台构建失败"; return 1; }
    echo "✅ 管理后台构建完成 → frontend-admin/dist"
}

build_frontend() {
    build_mobile
    echo ""
    build_admin
    cd "$BASE_DIR"
}

# ==================== 构建后端 ====================

build_backend() {
    echo "📦 编译后端服务..."
    cd "$BASE_DIR/backend"
    if [ ! -d "uploads" ]; then
        mkdir -p uploads
    fi
    go build -o orderfood-server main.go || { echo "❌ 后端编译失败"; return 1; }
    echo "✅ 后端编译完成 → backend/orderfood-server"
    cd "$BASE_DIR"
}

# ==================== 服务管理 ====================

stop_backend() {
    if [ -f "$BACKEND_PID_FILE" ]; then
        local pid
        pid=$(cat "$BACKEND_PID_FILE")
        if kill -0 "$pid" 2>/dev/null; then
            echo "🔄 停止后端服务 (PID: $pid)..."
            kill "$pid" 2>/dev/null
            sleep 1
        fi
        rm -f "$BACKEND_PID_FILE"
    fi
    # 同时杀掉可能残留的进程
    pkill -f "orderfood-server" 2>/dev/null || true
}

start_backend() {
    stop_backend
    echo "🚀 启动后端服务..."
    cd "$BASE_DIR/backend"
    nohup ./orderfood-server > /tmp/orderfood-backend.log 2>&1 &
    echo $! > "$BACKEND_PID_FILE"
    echo "✅ 后端服务已启动 (PID: $(cat "$BACKEND_PID_FILE"))"
    cd "$BASE_DIR"
}

stop_nginx() {
    # 先停掉我们自己的 nginx 实例
    if [ -f /run/nginx-orderfood.pid ]; then
        local pid
        pid=$(cat /run/nginx-orderfood.pid)
        if kill -0 "$pid" 2>/dev/null; then
            kill "$pid" 2>/dev/null || true
            sleep 1
        fi
        rm -f /run/nginx-orderfood.pid
    fi
    # 停掉系统默认 nginx（释放 80 端口）
    if command -v systemctl &> /dev/null && systemctl is-active nginx &> /dev/null; then
        echo "🔄 停止系统 Nginx 服务..."
        systemctl stop nginx 2>/dev/null || true
    fi
    # 确保没有残留的 nginx 进程占用 80 端口
    if ss -tlnp 2>/dev/null | grep -q ':80 .*nginx'; then
        echo "🔄 清理占用 80 端口的 Nginx 进程..."
        pkill nginx 2>/dev/null || true
        sleep 1
    fi
}

start_nginx() {
    stop_nginx
    echo "🌐 启动 Nginx (端口 80)..."
    nginx -c "$NGINX_CONF" 2>&1
    if [ $? -eq 0 ]; then
        echo "✅ Nginx 已启动"
    else
        echo "❌ Nginx 启动失败，请检查配置或端口占用"
        return 1
    fi
}

stop_all() {
    echo "🛑 停止所有服务..."
    stop_backend
    stop_nginx
    echo "✅ 所有服务已停止"
}

# ==================== 查看状态 ====================

show_status() {
    echo ""
    echo "📊 服务状态："

    # 后端
    if [ -f "$BACKEND_PID_FILE" ] && kill -0 "$(cat "$BACKEND_PID_FILE")" 2>/dev/null; then
        echo "  📡 后端服务:    运行中 (PID: $(cat "$BACKEND_PID_FILE"))"
    else
        echo "  📡 后端服务:    已停止"
    fi

    # Nginx
    if [ -f /run/nginx-orderfood.pid ] && kill -0 "$(cat /run/nginx-orderfood.pid)" 2>/dev/null; then
        echo "  🌐 Nginx:       运行中 (PID: $(cat /run/nginx-orderfood.pid))"
    else
        echo "  🌐 Nginx:       已停止"
    fi

    echo ""
}

# ==================== 查看日志 ====================

show_logs() {
    echo ""
    echo "📋 后端最近日志："
    tail -30 /tmp/orderfood-backend.log 2>/dev/null || echo "  (无日志)"
    echo ""
}

# ==================== 帮助 ====================

show_help() {
    echo ""
    echo "用法: bash start.sh <命令>"
    echo ""
    echo "命令："
    echo "  deploy        - 全量部署（构建前端+后端，启动所有服务）"
    echo "  start         - 启动所有服务（需先 deploy 过）"
    echo "  stop          - 停止所有服务"
    echo "  restart       - 重启所有服务"
    echo "  restart-back  - 仅重新编译并重启后端"
    echo "  restart-front - 仅重新构建前端并重载 Nginx"
    echo "  status        - 查看服务状态"
    echo "  logs          - 查看后端日志"
    echo "  help          - 显示此帮助"
    echo ""
    echo "访问地址："
    echo "  📱 移动端:     http://localhost"
    echo "  💻 管理后台:   http://localhost/admin/"
    echo ""
}

# ==================== 命令入口 ====================

case "${1:-help}" in
    deploy)
        check_env
        echo "🚀 开始全量部署..."
        echo ""
        build_frontend
        echo ""
        build_backend
        echo ""
        start_backend
        sleep 1
        start_nginx
        echo ""
        echo "===================="
        echo "🎉 部署完成！"
        echo ""
        echo "📱 移动端:     http://localhost"
        echo "💻 管理后台:   http://localhost/admin/"
        echo ""
        echo "🔐 管理员账户: admin / admin123"
        echo "===================="
        ;;
    start)
        start_backend
        sleep 1
        start_nginx
        echo ""
        echo "🎉 服务已启动"
        echo "📱 移动端:     http://localhost"
        echo "💻 管理后台:   http://localhost/admin/"
        ;;
    stop)
        stop_all
        ;;
    restart)
        stop_all
        sleep 1
        start_backend
        sleep 1
        start_nginx
        echo ""
        echo "🎉 服务已重启"
        ;;
    restart-back|rb)
        build_backend
        start_backend
        echo ""
        echo "🎉 后端已重启"
        ;;
    restart-front|rf)
        build_frontend
        echo ""
        echo "🔄 重载 Nginx..."
        nginx -s reload -c "$NGINX_CONF" 2>/dev/null || start_nginx
        echo "🎉 前端已更新"
        ;;
    status|s)
        show_status
        ;;
    logs|l)
        show_logs
        ;;
    help|h|--help|-h)
        show_help
        ;;
    *)
        echo "❓ 未知命令: $1"
        show_help
        exit 1
        ;;
esac
