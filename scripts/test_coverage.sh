#!/bin/bash

# 测试覆盖率脚本
set -e

echo "🧪 运行测试覆盖率分析..."

# 创建覆盖率输出目录
mkdir -p coverage

# 运行测试并生成覆盖率报告
echo "📊 生成覆盖率报告..."
go test -v -coverprofile=coverage/coverage.out ./...

# 检查是否生成了覆盖率文件
if [ ! -f "coverage/coverage.out" ]; then
    echo "❌ 覆盖率文件生成失败"
    exit 1
fi

# 生成HTML报告
echo "🌐 生成HTML覆盖率报告..."
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# 生成总体覆盖率统计
echo "📈 生成覆盖率统计..."
go tool cover -func=coverage/coverage.out > coverage/coverage.txt

# 显示覆盖率摘要
echo ""
echo "📋 覆盖率摘要:"
echo "=============="
total_coverage=$(go tool cover -func=coverage/coverage.out | grep total | awk '{print $3}')
echo "总覆盖率: $total_coverage"

# 按包显示覆盖率
echo ""
echo "📦 各包覆盖率:"
echo "=============="
go tool cover -func=coverage/coverage.out | grep -v total | while read line; do
    package=$(echo $line | awk '{print $1}' | cut -d'/' -f1-2)
    coverage=$(echo $line | awk '{print $3}')
    printf "%-30s %s\n" "$package" "$coverage"
done | sort -u

# 检查覆盖率阈值
coverage_threshold=70
current_coverage=$(echo $total_coverage | sed 's/%//')
current_coverage_int=$(echo $current_coverage | cut -d'.' -f1)

echo ""
if [ "$current_coverage_int" -ge "$coverage_threshold" ]; then
    echo "✅ 覆盖率达标 ($total_coverage >= ${coverage_threshold}%)"
else
    echo "⚠️  覆盖率不达标 ($total_coverage < ${coverage_threshold}%)"
fi

# 生成覆盖率徽章数据
echo ""
echo "🏷️  生成覆盖率徽章数据..."
badge_color="red"
if [ "$current_coverage_int" -ge 80 ]; then
    badge_color="brightgreen"
elif [ "$current_coverage_int" -ge 70 ]; then
    badge_color="yellow"
elif [ "$current_coverage_int" -ge 60 ]; then
    badge_color="orange"
fi

cat > coverage/badge.json << EOF
{
  "schemaVersion": 1,
  "label": "coverage",
  "message": "$total_coverage",
  "color": "$badge_color"
}
EOF

echo "📁 覆盖率文件生成完成:"
echo "  - HTML报告: coverage/coverage.html"
echo "  - 文本报告: coverage/coverage.txt"
echo "  - 原始数据: coverage/coverage.out"
echo "  - 徽章数据: coverage/badge.json"

# 如果在CI环境中，输出更多信息
if [ "$CI" = "true" ]; then
    echo ""
    echo "🤖 CI环境检测到，输出详细信息:"
    echo "================================"
    
    # 输出未覆盖的函数
    echo "❌ 未覆盖的函数:"
    go tool cover -func=coverage/coverage.out | grep "0.0%" || echo "  无未覆盖函数"
    
    # 输出覆盖率最低的文件
    echo ""
    echo "⚠️  覆盖率最低的文件 (前10个):"
    go tool cover -func=coverage/coverage.out | grep -v total | sort -k3 -n | head -10
fi

echo ""
echo "✅ 测试覆盖率分析完成!"