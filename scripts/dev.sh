#!/bin/bash

# 短剧管理系统开发启动脚本

set -e

echo "🚀 启动短剧管理系统开发环境..."

# 检查 Docker 是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker 未运行，请先启动 Docker"
    exit 1
fi

# 启动依赖服务
echo "📦 启动 MySQL 和 Redis 服务..."
docker compose up -d mysql redis

# 等待服务健康
echo "⏳ 等待服务启动..."
sleep 5

# 检查服务状态
echo "🔍 检查服务状态..."
docker compose ps

# 构建前端（如果需要）
if [ "$1" = "--build-frontend" ] || [ ! -d "web/dist" ]; then
    echo "🔨 构建前端..."
    cd web
    npm install
    npm run build
    cd ..
fi

# 启动 Go 服务器
echo "🌟 启动 Go 服务器..."
echo "访问地址: http://localhost:1800"
echo "健康检查: http://localhost:1800/health"
echo "管理后台: http://localhost:1800/login"
echo ""
echo "按 Ctrl+C 停止服务器"

unset GOROOT
go run cmd/server/main.go
