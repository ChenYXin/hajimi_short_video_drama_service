# 详细部署指南

> 💡 **快速开始**: 如果你只想快速体验项目，请查看 [README.md](README.md) 的快速部署部分。

本文档提供详细的部署配置和故障排除指南。

## 📋 目录

- [环境要求](#环境要求)
- [详细部署步骤](#详细部署步骤)
- [高级配置](#高级配置)
- [监控配置](#监控配置)
- [故障排除](#故障排除)

## 环境要求

### 基础要求

- Docker 20.10+
- Docker Compose 2.0+
- Go 1.21+ (仅开发环境)
- Make (可选，用于简化命令)

### 系统要求

- **最小配置**: 2 CPU, 4GB RAM, 20GB 磁盘
- **推荐配置**: 4 CPU, 8GB RAM, 50GB 磁盘

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd gin-mysql-api
```

### 2. 一键部署

```bash
# 使用 Make (推荐)
make deploy

# 或直接运行脚本
./scripts/deploy.sh
```

### 3. 验证部署

```bash
# 检查服务状态
make status

# 健康检查
make health

# 查看日志
make logs
```

## 开发环境部署

### 方式一：本地开发

```bash
# 1. 启动依赖服务
docker-compose -f docker-compose.dev.yml up -d mysql redis

# 2. 安装依赖
make deps

# 3. 运行数据库迁移
make db-migrate

# 4. 启动应用
make dev
```

### 方式二：Docker 开发环境

```bash
# 启动完整开发环境
make deploy-dev
```

### 开发环境配置

开发环境使用 `configs/config.yaml` 配置文件，主要特点：

- 数据库: localhost:3306
- Redis: localhost:6379
- 日志级别: debug
- 热重载支持

## 生产环境部署

### 1. 环境准备

```bash
# 创建必要目录
mkdir -p uploads logs

# 设置权限
chmod 755 uploads logs
```

### 2. 配置应用

#### 方式一：使用配置文件
```bash
# 复制配置模板
cp configs/config.example.yaml configs/production.yaml

# 编辑生产环境配置
vim configs/production.yaml
```

#### 方式二：使用环境变量
```bash
# 设置关键环境变量
export APP_SERVER_MODE=release
export APP_DATABASE_PASSWORD=your-secure-password
export APP_JWT_SECRET=your-jwt-secret-key
export APP_LOGGING_LEVEL=info
```

详细的环境变量配置说明请参考 [configs/ENV_VARIABLES.md](configs/ENV_VARIABLES.md)。

### 3. 部署应用

```bash
# 完整部署
make deploy

# 或分步部署
make docker-build
docker-compose up -d
```

### 4. 配置反向代理 (可选)

#### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:1800;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /metrics {
        proxy_pass http://localhost:9090;
        # 限制访问
        allow 127.0.0.1;
        deny all;
    }
}
```

## 监控配置

### Prometheus

- **访问地址**: http://localhost:9091
- **配置文件**: `configs/prometheus.yml`
- **告警规则**: `configs/alert_rules.yml`

### Grafana

- **访问地址**: http://localhost:3000
- **默认账号**: admin/admin123
- **数据源**: Prometheus (http://prometheus:9090)

### 监控指标

#### 应用指标

- HTTP 请求数量和响应时间
- 数据库连接池状态
- Redis 连接状态
- 业务指标 (用户注册、短剧观看等)

#### 系统指标

- CPU 使用率
- 内存使用率
- 磁盘使用率
- 网络流量

### 告警配置

主要告警规则：

- HTTP 错误率 > 10%
- 响应时间 > 1秒
- 数据库连接使用率 > 80%
- 应用程序宕机
- 系统资源使用率过高

## 服务管理

### 常用命令

```bash
# 查看服务状态
make status

# 查看日志
make logs          # 应用日志
make logs-all      # 所有服务日志

# 重启服务
docker-compose restart app

# 停止服务
make stop

# 更新应用
git pull
make deploy
```

### 数据库管理

```bash
# 进入数据库
make db-shell

# 备份数据库
make backup-db

# 恢复数据库
make restore-db

# 运行迁移
make db-migrate
```

### 容器管理

```bash
# 进入应用容器
make shell

# 进入数据库容器
make db-shell

# 进入 Redis 容器
make redis-shell

# 查看容器资源使用
docker stats
```

## 配置说明

### 环境变量

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| `APP_SERVER_HOST` | 服务器监听地址 | 0.0.0.0 |
| `APP_SERVER_PORT` | 服务器端口 | 1800 |
| `APP_DATABASE_HOST` | 数据库地址 | mysql |
| `APP_DATABASE_PASSWORD` | 数据库密码 | rootpassword |
| `APP_REDIS_HOST` | Redis 地址 | redis |
| `APP_JWT_SECRET` | JWT 密钥 | - |

### 配置文件

- `configs/config.yaml`: 开发环境配置
- `configs/production.yaml`: 生产环境配置
- `configs/test.yaml`: 测试环境配置

## 性能优化

### 数据库优化

```yaml
database:
  maxIdleConns: 20      # 最大空闲连接
  maxOpenConns: 200     # 最大打开连接
  connMaxLifetime: 7200 # 连接最大生存时间(秒)
```

### Redis 优化

```yaml
redis:
  poolSize: 20  # 连接池大小
```

### 应用优化

- 启用 Gzip 压缩
- 配置静态文件缓存
- 使用 CDN 加速静态资源

## 安全配置

### 1. 数据库安全

- 使用强密码
- 限制数据库访问IP
- 定期备份数据

### 2. 应用安全

- 更改默认JWT密钥
- 启用HTTPS
- 配置防火墙规则

### 3. 容器安全

- 使用非root用户运行
- 定期更新基础镜像
- 扫描安全漏洞

## 故障排除

### 常见问题

#### 1. 应用启动失败

```bash
# 查看详细日志
docker-compose logs app

# 检查配置
docker-compose config

# 检查端口占用
netstat -tlnp | grep :1800
```

#### 2. 数据库连接失败

```bash
# 检查数据库状态
docker-compose ps mysql

# 测试数据库连接
docker-compose exec mysql mysqladmin ping -h localhost -u root -p

# 查看数据库日志
docker-compose logs mysql
```

#### 3. Redis 连接失败

```bash
# 检查 Redis 状态
docker-compose ps redis

# 测试 Redis 连接
docker-compose exec redis redis-cli ping

# 查看 Redis 日志
docker-compose logs redis
```

#### 4. 监控服务异常

```bash
# 检查 Prometheus
curl http://localhost:9091/-/healthy

# 检查 Grafana
curl http://localhost:3000/api/health

# 重启监控服务
docker-compose restart prometheus grafana
```

### 日志分析

#### 应用日志位置

- 容器内: `/app/logs/app.log`
- 主机上: `./logs/app.log`

#### 日志级别

- `debug`: 调试信息
- `info`: 一般信息
- `warn`: 警告信息
- `error`: 错误信息

### 性能问题排查

```bash
# 查看系统资源使用
docker stats

# 查看应用性能指标
curl http://localhost:9090/metrics

# 数据库性能分析
docker-compose exec mysql mysql -u root -p -e "SHOW PROCESSLIST;"
```

## 升级指南

### 应用升级

```bash
# 1. 备份数据
make backup-db

# 2. 拉取最新代码
git pull

# 3. 重新部署
make deploy

# 4. 验证升级
make health
```

### 数据库升级

```bash
# 运行数据库迁移
make db-migrate

# 验证数据完整性
make db-shell
# 在数据库中执行验证查询
```

## 联系支持

如果遇到问题，请：

1. 查看本文档的故障排除部分
2. 检查项目的 Issues 页面
3. 提交新的 Issue 并包含详细的错误信息和日志