#!/bin/bash

# 测试覆盖率报告脚本

set -e

echo "生成测试覆盖率报告..."

# 创建覆盖率输出目录
mkdir -p coverage

# 运行测试并生成覆盖率报告
echo "运行所有测试并收集覆盖率数据..."

# Repository层测试覆盖率
echo "Repository层测试覆盖率..."
go test ./internal/repository/... -coverprofile=coverage/repository.out -covermode=atomic || true

# Service层测试覆盖率
echo "Service层测试覆盖率..."
go test ./internal/service/... -coverprofile=coverage/service.out -covermode=atomic || true

# 中间件测试覆盖率
echo "中间件测试覆盖率..."
go test ./internal/middleware/... -coverprofile=coverage/middleware.out -covermode=atomic || true

# 工具函数测试覆盖率
echo "工具函数测试覆盖率..."
go test ./pkg/utils/... -coverprofile=coverage/utils.out -covermode=atomic || true

# 模型测试覆盖率
echo "模型测试覆盖率..."
go test ./internal/models/... -coverprofile=coverage/models.out -covermode=atomic || true

# 合并所有覆盖率文件
echo "合并覆盖率报告..."
echo "mode: atomic" > coverage/coverage.out

# 合并所有 .out 文件（跳过第一行的 mode 声明）
for file in coverage/*.out; do
    if [ "$file" != "coverage/coverage.out" ] && [ -f "$file" ]; then
        tail -n +2 "$file" >> coverage/coverage.out
    fi
done

# 生成 HTML 报告
echo "生成 HTML 覆盖率报告..."
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# 显示覆盖率统计
echo "覆盖率统计:"
go tool cover -func=coverage/coverage.out

echo ""
echo "覆盖率报告已生成:"
echo "- 文本报告: coverage/coverage.out"
echo "- HTML报告: coverage/coverage.html"
echo ""
echo "打开 HTML 报告查看详细覆盖率信息:"
echo "open coverage/coverage.html"