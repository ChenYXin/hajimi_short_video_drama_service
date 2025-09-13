#!/bin/bash

# 部署脚本
set -e

echo "开始部署 Gin MySQL API..."

# 检查Docker和Docker Compose是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "错误: Docker Compose 未安装"
    exit 1
fi

# 设置环境变量
export COMPOSE_PROJECT_NAME=gin-mysql-api

# 创建必要的目录
mkdir -p uploads logs

# 停止现有服务
echo "停止现有服务..."
docker-compose down

# 清理旧的镜像（可选）
if [ "$1" = "--clean" ]; then
    echo "清理旧镜像..."
    docker system prune -f
    docker volume prune -f
fi

# 构建并启动服务
echo "构建并启动服务..."
docker-compose up --build -d

# 等待服务启动
echo "等待服务启动..."
sleep 30

# 检查服务状态
echo "检查服务状态..."
docker-compose ps

# 健康检查
echo "执行健康检查..."
max_attempts=30
attempt=1

while [ $attempt -le $max_attempts ]; do
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ 应用程序健康检查通过"
        break
    else
        echo "⏳ 等待应用程序启动... (尝试 $attempt/$max_attempts)"
        sleep 10
        ((attempt++))
    fi
done

if [ $attempt -gt $max_attempts ]; then
    echo "❌ 应用程序健康检查失败"
    echo "查看日志:"
    docker-compose logs app
    exit 1
fi

# 显示服务信息
echo ""
echo "🎉 部署完成!"
echo ""
echo "服务访问地址:"
echo "  - 应用程序: http://localhost:8080"
echo "  - 管理后台: http://localhost:8080/admin"
echo "  - API文档: http://localhost:8080/api/docs"
echo "  - 健康检查: http://localhost:8080/health"
echo "  - Prometheus: http://localhost:9091"
echo "  - Grafana: http://localhost:3000 (admin/admin123)"
echo ""
echo "数据库连接:"
echo "  - MySQL: localhost:3306 (root/rootpassword)"
echo "  - Redis: localhost:6379"
echo ""
echo "监控指标:"
echo "  - 应用指标: http://localhost:9090/metrics"
echo "  - MySQL指标: http://localhost:9104/metrics"
echo "  - Redis指标: http://localhost:9121/metrics"
echo "  - 系统指标: http://localhost:9100/metrics"
echo ""
echo "查看日志: docker-compose logs -f [service_name]"
echo "停止服务: docker-compose down"