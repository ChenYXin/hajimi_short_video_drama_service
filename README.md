# Gin MySQL API

短剧服务平台后端系统，基于 Gin 框架的 RESTful API 服务。

## 项目结构

```
gin-mysql-api/
├── cmd/
│   └── server/          # 应用程序入口
├── internal/
│   ├── models/          # 数据模型
│   ├── repository/      # 数据访问层
│   ├── service/         # 业务逻辑层
│   ├── handler/         # HTTP 处理器
│   └── middleware/      # 中间件
├── pkg/
│   ├── config/          # 配置管理
│   └── utils/           # 工具函数
├── web/
│   ├── templates/       # HTML 模板
│   └── static/          # 静态资源
├── configs/             # 配置文件
├── scripts/             # 脚本文件
├── tests/               # 测试文件
└── uploads/             # 上传文件目录
```

## 技术栈

- **Web Framework**: Gin
- **Database**: MySQL + GORM
- **Cache**: Redis
- **Authentication**: JWT
- **Configuration**: Viper
- **Monitoring**: Prometheus

## 快速开始

1. 复制环境配置文件：
```bash
cp .env.example .env
```

2. 修改配置文件中的数据库和Redis连接信息

3. 运行应用：
```bash
go run cmd/server/main.go
```

## 配置说明

配置文件位于 `configs/config.yaml`，支持通过环境变量覆盖配置项。

环境变量命名规则：`APP_` + 配置路径（用下划线分隔）

例如：
- `server.port` -> `APP_SERVER_PORT`
- `database.host` -> `APP_DATABASE_HOST`