#!/bin/bash

# 单元测试运行脚本

set -e

echo "开始运行单元测试..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试结果统计
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 运行测试的函数
run_test() {
    local test_path=$1
    local test_name=$2
    
    echo -e "${YELLOW}运行测试: $test_name${NC}"
    
    if go test $test_path -v; then
        echo -e "${GREEN}✓ $test_name 测试通过${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}✗ $test_name 测试失败${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo ""
}

# 工具函数测试
echo "=== 工具函数测试 ==="
run_test "./pkg/utils/password_test.go ./pkg/utils/password.go" "密码工具函数"
run_test "./pkg/utils/jwt_test.go ./pkg/utils/jwt.go" "JWT工具函数"

# 模型验证测试
echo "=== 模型验证测试 ==="
run_test "./internal/models/validator_test.go ./internal/models/validator.go" "数据验证"

# Service层测试
echo "=== Service层测试 ==="
run_test "./internal/service/cache_service_test.go ./internal/service/cache_service.go" "缓存服务"

# 中间件测试
echo "=== 中间件测试 ==="
run_test "./internal/middleware/logger_test.go ./internal/middleware/logger.go" "日志中间件"
run_test "./internal/middleware/security_test.go ./internal/middleware/security.go" "安全中间件"

# 输出测试结果统计
echo "=================================="
echo "测试结果统计:"
echo -e "总测试数: $TOTAL_TESTS"
echo -e "${GREEN}通过: $PASSED_TESTS${NC}"
echo -e "${RED}失败: $FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}所有测试都通过了！${NC}"
    exit 0
else
    echo -e "${RED}有 $FAILED_TESTS 个测试失败${NC}"
    exit 1
fi