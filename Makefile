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
	go test -v ./...

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
	mysql -u root -p hajimi < scripts/seed_data.sql

db-reset: ## 重置数据库
	@echo "🔄 重置数据库..."
	mysql -u root -p -e "DROP DATABASE IF EXISTS hajimi; CREATE DATABASE hajimi;"
	mysql -u root -p hajimi < scripts/init_db.sql

# Docker
docker-build: ## 构建 Docker 镜像
	@echo "🐳 构建 Docker 镜像..."
	docker build -t gin-mysql-api .

docker-run: ## 运行 Docker 容器
	@echo "🐳 运行 Docker 容器..."
	docker compose up -d

docker-stop: ## 停止 Docker 容器
	@echo "🛑 停止 Docker 容器..."
	docker compose down

docker-logs: ## 查看 Docker 日志
	@echo "📋 查看 Docker 日志..."
	docker compose logs -f

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
# 部署相
关
deploy: ## 部署到生产环境
	@echo "🚀 部署到生产环境..."
	./scripts/deploy.sh

deploy-dev: ## 部署开发环境
	@echo "🔧 部署开发环境..."
	APP_MODE=debug LOG_LEVEL=debug docker compose up --build -d

stop: ## 停止所有服务
	@echo "🛑 停止所有服务..."
	docker compose down

logs: ## 查看应用日志
	@echo "📋 查看应用日志..."
	docker compose logs -f app

logs-all: ## 查看所有服务日志
	@echo "📋 查看所有服务日志..."
	docker compose logs -f

status: ## 查看服务状态
	@echo "📊 服务状态:"
	docker compose ps

# 容器操作
shell: ## 进入应用容器
	@echo "🐚 进入应用容器..."
	docker compose exec app sh

db-shell: ## 进入数据库容器
	@echo "🗄️ 进入数据库容器..."
	docker compose exec mysql mysql -u root -pa123d789cDE hajimi

redis-shell: ## 进入Redis容器
	@echo "🔴 进入Redis容器..."
	docker compose exec redis redis-cli

# 备份和恢复
backup-db: ## 备份数据库
	@echo "💾 备份数据库..."
	docker compose exec mysql mysqldump -u root -pa123d789cDE hajimi > backup_$(shell date +%Y%m%d_%H%M%S).sql

restore-db: ## 恢复数据库
	@echo "🔄 恢复数据库..."
	@read -p "请输入备份文件名: " file; \
	docker compose exec -T mysql mysql -u root -pa123d789cDE hajimi < $$file

# 监控
metrics: ## 查看Prometheus指标
	@echo "📊 打开Prometheus指标..."
	@echo "应用指标: http://localhost:9090/metrics"
	@echo "Prometheus: http://localhost:9091"
	@echo "Grafana: http://localhost:3000"

health: ## 检查应用健康状态
	@echo "🏥 检查应用健康状态..."
	@curl -f http://localhost:8080/health || echo "应用程序未运行"

# 完整的CI/CD流程
ci: deps fmt vet lint test build ## 运行CI流程
	@echo "✅ CI流程完成"

cd: ci deploy ## 运行CD流程
	@echo "🚀 CD流程完成"

# 故障排除
fix-docker: ## 修复Docker常见问题
	@echo "🔧 修复Docker问题..."
	docker system prune -f
	docker compose down
	docker compose pull
	docker compose up --build -d

reset: ## 重置所有服务和数据
	@echo "🔄 重置所有服务..."
	docker compose down -v
	docker system prune -f
	docker compose up --build -d

health-check: ## 健康检查
	@echo "🏥 执行健康检查..."
	@echo "检查Docker状态:"
	@docker info > /dev/null && echo "✅ Docker运行正常" || echo "❌ Docker未运行"
	@echo "检查服务状态:"
	@docker compose ps
	@echo "检查应用健康:"
	@curl -f http://localhost:8080/health > /dev/null 2>&1 && echo "✅ 应用运行正常" || echo "❌ 应用未响应"