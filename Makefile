# Makefile for gin-mysql-api

.PHONY: help build run test test-unit test-integration test-coverage clean deps lint fmt vet security

# 默认目标
help: ## 显示帮助信息
	@echo "可用的命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# 构建
build: ## 构建应用程序
	@echo "🔨 构建应用程序..."
	go build -o bin/gin-mysql-api cmd/server/main.go

# 运行
run: ## 运行应用程序
	@echo "🚀 启动应用程序..."
	go run cmd/server/main.go

# 开发模式运行
dev: ## 开发模式运行（带热重载）
	@echo "🔥 开发模式启动..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "请安装 air: go install github.com/cosmtrek/air@latest"; \
		go run cmd/server/main.go; \
	fi

# 依赖管理
deps: ## 安装和更新依赖
	@echo "📦 管理依赖..."
	go mod tidy
	go mod download

# 测试
test: ## 运行所有测试
	@echo "🧪 运行所有测试..."
	./scripts/run_tests.sh

test-unit: ## 运行单元测试
	@echo "🔬 运行单元测试..."
	go test -v ./pkg/... ./internal/... -short

test-integration: ## 运行集成测试
	@echo "🔗 运行集成测试..."
	go test -v ./... -tags=integration

test-coverage: ## 生成测试覆盖率报告
	@echo "📊 生成测试覆盖率报告..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

test-race: ## 运行竞态条件检测
	@echo "🏃 运行竞态条件检测..."
	go test -race ./...

test-bench: ## 运行基准测试
	@echo "⚡ 运行基准测试..."
	go test -bench=. -benchmem ./...

# 代码质量
lint: ## 运行代码检查
	@echo "🔍 运行代码检查..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "请安装 golangci-lint"; \
		echo "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.54.2"; \
	fi

fmt: ## 格式化代码
	@echo "✨ 格式化代码..."
	gofmt -w .
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	fi

vet: ## 运行 go vet
	@echo "🔍 运行 go vet..."
	go vet ./...

security: ## 运行安全检查
	@echo "🔒 运行安全检查..."
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "请安装 gosec: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# 数据库
db-migrate: ## 运行数据库迁移
	@echo "🗄️ 运行数据库迁移..."
	go run cmd/migrate/main.go

db-seed: ## 填充测试数据
	@echo "🌱 填充测试数据..."
	mysql -u root -p gin_mysql_api < scripts/seed_data.sql

db-reset: ## 重置数据库
	@echo "🔄 重置数据库..."
	mysql -u root -p -e "DROP DATABASE IF EXISTS gin_mysql_api; CREATE DATABASE gin_mysql_api;"
	mysql -u root -p gin_mysql_api < scripts/init_db.sql

# Docker
docker-build: ## 构建 Docker 镜像
	@echo "🐳 构建 Docker 镜像..."
	docker build -t gin-mysql-api .

docker-run: ## 运行 Docker 容器
	@echo "🐳 运行 Docker 容器..."
	docker-compose up -d

docker-stop: ## 停止 Docker 容器
	@echo "🛑 停止 Docker 容器..."
	docker-compose down

docker-logs: ## 查看 Docker 日志
	@echo "📋 查看 Docker 日志..."
	docker-compose logs -f

# 清理
clean: ## 清理构建文件
	@echo "🧹 清理构建文件..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -rf test_uploads/
	go clean -cache
	go clean -testcache

# 安装开发工具
install-tools: ## 安装开发工具
	@echo "🛠️ 安装开发工具..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install golang.org/x/tools/cmd/goimports@latest

# 生成文档
docs: ## 生成 API 文档
	@echo "📚 生成 API 文档..."
	@if command -v swag > /dev/null; then \
		swag init -g cmd/server/main.go; \
	else \
		echo "请安装 swag: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# 检查所有
check: deps fmt vet lint test ## 运行所有检查

# 发布准备
release: clean check build ## 准备发布版本
	@echo "🎉 发布准备完成!"

# 快速开始
quick-start: deps db-migrate db-seed run ## 快速开始（首次运行）

# 显示项目信息
info: ## 显示项目信息
	@echo "📋 项目信息:"
	@echo "Go 版本: $$(go version)"
	@echo "项目路径: $$(pwd)"
	@echo "Git 分支: $$(git branch --show-current 2>/dev/null || echo 'N/A')"
	@echo "Git 提交: $$(git rev-parse --short HEAD 2>/dev/null || echo 'N/A')"
	@echo "依赖数量: $$(go list -m all | wc -l)"