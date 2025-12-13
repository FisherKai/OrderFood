#!/bin/bash

echo "🍽️  点餐系统启动脚本"
echo "===================="
echo ""

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "❌ 未检测到Go，请先安装Go 1.21+"
    exit 1
fi

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo "❌ 未检测到Node.js，请先安装Node.js 16+"
    exit 1
fi

echo "✅ 环境检查通过"
echo ""

# 启动后端
echo "📦 启动后端服务..."
cd backend
if [ ! -d "uploads" ]; then
    mkdir uploads
fi

# 后台运行
go run main.go &
BACKEND_PID=$!
echo "✅ 后端服务已启动 (PID: $BACKEND_PID)"
cd ..

# 等待后端启动
sleep 2

# 启动移动端前端
echo ""
echo "📱 启动移动端前端..."
cd frontend-mobile

# 检查依赖
if [ ! -d "node_modules" ]; then
    echo "📦 安装移动端依赖..."
    npm install
fi

npm run dev &
MOBILE_PID=$!
echo "✅ 移动端前端已启动 (PID: $MOBILE_PID)"
cd ..

# 启动PC端管理后台
echo ""
echo "💻 启动PC端管理后台..."
cd frontend-admin

# 检查依赖
if [ ! -d "node_modules" ]; then
    echo "📦 安装管理后台依赖..."
    npm install
fi

npm run dev &
ADMIN_PID=$!
echo "✅ PC端管理后台已启动 (PID: $ADMIN_PID)"
cd ..

echo ""
echo "===================="
echo "🎉 所有服务已启动！"
echo ""
echo "📡 后端地址: http://localhost:8080"
echo "📱 移动端地址: http://localhost:3000"
echo "💻 管理后台地址: http://localhost:3001"
echo ""
echo "🔐 管理员账户: admin / admin123"
echo ""
echo "按 Ctrl+C 停止所有服务"
echo "===================="

# 等待用户中断
wait
