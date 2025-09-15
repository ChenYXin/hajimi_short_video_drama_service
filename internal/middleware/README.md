# 中间件系统使用指南

本项目实现了完整的中间件系统，提供错误处理、日志记录、安全防护、监控指标等功能。

## 中间件组件

### 1. 错误处理中间件 (`error.go`)

提供全局错误处理和参数验证错误处理：

- **ErrorHandler**: 全局 panic 恢复和错误处理
- **ValidationErrorHandler**: 参数验证错误处理
- **NotFoundHandler**: 404 错误处理
- **MethodNotAllowedHandler**: 405 错误处理

```go
// 使用示例
router.Use(middleware.ErrorHandler())
router.NoRoute(middleware.NotFoundHandler())
router.NoMethod(middleware.MethodNotAllowedHandler())
```

### 2. 日志记录中间件 (`logger.go`)

提供请求日志记录功能：

- **Logger**: 基础日志记录（JSON 格式）
- **DetailedLogger**: 详细日志记录（包含请求/响应体）
- **RequestIDMiddleware**: 请求 ID 生成和传递

```go
// 基础日志
router.Use(middleware.Logger())

// 详细日志（开发环境）
config := middleware.LoggerConfig{
    LogRequestBody:  true,
    LogResponseBody: false,
    MaxBodySize:     1024 * 10, // 10KB
}
router.Use(middleware.DetailedLogger(config))
```

### 3. CORS 中间件 (`cors.go`)

处理跨域资源共享：

```go
config := middleware.CORSConfig{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:           86400,
}
router.Use(middleware.CORS(config))
```

### 4. 安全中间件 (`security.go`)

提供多种安全防护功能：

- **Security**: 安全头设置（XSS、CSRF、HSTS 等）
- **SimpleRateLimit**: 简单限流
- **IPWhitelist**: IP 白名单
- **UserAgentFilter**: User-Agent 过滤
- **RequestSizeLimit**: 请求大小限制
- **Timeout**: 请求超时

```go
// 安全头
securityConfig := middleware.SecurityConfig{
    XSSProtection:      "1; mode=block",
    ContentTypeNosniff: true,
    XFrameOptions:      "DENY",
    HSTSMaxAge:         31536000,
}
router.Use(middleware.Security(securityConfig))

// 限流
rateLimitConfig := middleware.RateLimitConfig{
    MaxRequests: 100,
    WindowSize:  60, // 1 minute
}
router.Use(middleware.SimpleRateLimit(rateLimitConfig))
```


### 6. 中间件管理器 (`manager.go`)

统一管理所有中间件：

```go
manager := middleware.NewManager(config)

// 开发环境
manager.SetupDevelopmentMiddlewares(router)

// 生产环境
manager.SetupProductionMiddlewares(router)
```


## 日志格式

### 基础日志格式

```json
{
  "timestamp": "2023-12-01T10:00:00Z",
  "status": 200,
  "latency": "10.5ms",
  "client_ip": "192.168.1.1",
  "method": "GET",
  "path": "/api/dramas",
  "user_agent": "Mozilla/5.0...",
  "body_size": 1024,
  "request_id": "1234567890",
  "user_id": 123,
  "username": "testuser"
}
```

### 详细日志格式

```json
{
  "timestamp": "2023-12-01T10:00:00Z",
  "status": 200,
  "latency": "10.5ms",
  "latency_ms": 10.5,
  "client_ip": "192.168.1.1",
  "method": "POST",
  "path": "/api/auth/login",
  "user_agent": "Mozilla/5.0...",
  "body_size": 1024,
  "headers": {
    "Content-Type": "application/json",
    "Authorization": "[REDACTED]"
  },
  "request_body": {
    "email": "user@example.com",
    "password": "[REDACTED]"
  },
  "response_body": {
    "success": true,
    "message": "登录成功"
  }
}
```

## 错误响应格式

所有错误都使用统一的响应格式：

```json
{
  "success": false,
  "message": "错误描述",
  "error": "详细错误信息或错误列表"
}
```

### 常见错误类型

- **400 Bad Request**: 参数验证失败
- **401 Unauthorized**: 未认证
- **403 Forbidden**: 权限不足或访问被拒绝
- **404 Not Found**: 资源不存在
- **405 Method Not Allowed**: 请求方法不被允许
- **413 Request Entity Too Large**: 请求体过大
- **429 Too Many Requests**: 请求过于频繁
- **500 Internal Server Error**: 服务器内部错误

## 安全特性

### 1. 安全头设置

- **X-XSS-Protection**: XSS 攻击防护
- **X-Content-Type-Options**: MIME 类型嗅探防护
- **X-Frame-Options**: 点击劫持防护
- **Strict-Transport-Security**: HTTPS 强制
- **Content-Security-Policy**: 内容安全策略
- **Referrer-Policy**: 引用策略

### 2. 访问控制

- **CORS**: 跨域资源共享控制
- **IP 白名单**: 限制访问来源
- **User-Agent 过滤**: 阻止恶意爬虫
- **限流**: 防止 DDoS 攻击

### 3. 数据保护

- **请求大小限制**: 防止大文件攻击
- **敏感信息过滤**: 日志中自动过滤密码等敏感信息
- **请求超时**: 防止慢速攻击

## 性能优化

### 1. 日志优化

- 生产环境跳过健康检查等路径的日志
- 限制请求/响应体的记录大小
- 使用 JSON 格式便于日志分析


### 3. 中间件顺序

中间件按以下顺序执行以获得最佳性能：

1. 错误处理（最先）
2. 请求 ID
3. 日志记录
4. 安全头
5. CORS
7. 限流
8. 请求大小限制
9. 超时控制

## 配置示例

### 开发环境配置

```yaml
server:
  mode: debug
  
logging:
  level: debug
  
cors:
  allow_origins: ["http://localhost:3000"]
  allow_credentials: true
```

### 生产环境配置

```yaml
server:
  mode: release
  allowed_ips: ["10.0.0.0/8"]
  
logging:
  level: info
  
security:
  hsts_max_age: 31536000
  csp: "default-src 'self'"
```

## 最佳实践

1. **错误处理**: 始终使用统一的错误响应格式
2. **日志记录**: 记录足够的信息用于调试，但避免记录敏感信息
3. **安全防护**: 根据应用需求配置适当的安全策略
5. **限流策略**: 根据系统容量设置合理的限流参数
6. **中间件顺序**: 按照性能和功能需求合理安排中间件顺序
7. **环境区分**: 开发和生产环境使用不同的中间件配置