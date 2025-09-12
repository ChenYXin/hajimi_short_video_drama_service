# JWT 认证系统使用指南

本目录包含了完整的 JWT 认证系统实现，包括 JWT token 管理和密码哈希功能。

## 功能特性

### 1. JWT Token 管理 (`jwt.go`)

- **生成 Token**: 为用户生成包含用户信息的 JWT token
- **验证 Token**: 验证 token 的有效性并解析用户信息
- **刷新 Token**: 基于现有 token 生成新的 token
- **自定义过期时间**: 支持配置 token 过期时间

### 2. 密码安全 (`password.go`)

- **密码哈希**: 使用 bcrypt 算法对密码进行安全哈希
- **密码验证**: 验证明文密码与哈希密码是否匹配
- **安全性**: 使用 bcrypt 默认成本，每次哈希结果都不同

## 使用方法

### 1. 创建 JWT 管理器

```go
import (
    "time"
    "gin-mysql-api/pkg/utils"
)

// 创建 JWT 管理器
jwtManager := utils.NewJWTManager("your-secret-key", 24*time.Hour)
```

### 2. 生成和验证 Token

```go
// 生成 token
token, err := jwtManager.GenerateToken(userID, username, role)
if err != nil {
    // 处理错误
}

// 验证 token
claims, err := jwtManager.VerifyToken(token)
if err != nil {
    // token 无效
}

// 使用解析出的用户信息
fmt.Printf("用户ID: %d, 用户名: %s, 角色: %s\n", 
    claims.UserID, claims.Username, claims.Role)
```

### 3. 密码哈希和验证

```go
// 注册时哈希密码
hashedPassword, err := utils.HashPassword("user-password")
if err != nil {
    // 处理错误
}

// 登录时验证密码
isValid := utils.VerifyPassword(hashedPassword, "user-input-password")
if !isValid {
    // 密码错误
}
```

### 4. 在 Gin 中使用认证中间件

```go
import (
    "gin-mysql-api/internal/middleware"
    "github.com/gin-gonic/gin"
)

r := gin.New()

// 需要用户认证的路由
protected := r.Group("/api")
protected.Use(middleware.AuthMiddleware(jwtManager))
{
    protected.GET("/profile", func(c *gin.Context) {
        userID := c.GetUint("user_id")
        username := c.GetString("username")
        role := c.GetString("role")
        // 处理已认证的请求
    })
}

// 需要管理员权限的路由
admin := r.Group("/admin")
admin.Use(middleware.AdminAuthMiddleware(jwtManager))
{
    admin.GET("/dashboard", func(c *gin.Context) {
        // 只有管理员可以访问
    })
}

// 可选认证的路由
public := r.Group("/public")
public.Use(middleware.OptionalAuthMiddleware(jwtManager))
{
    public.GET("/content", func(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if exists {
            // 已认证用户的逻辑
        } else {
            // 未认证用户的逻辑
        }
    })
}
```

## 中间件说明

### AuthMiddleware
- **用途**: 保护需要用户认证的路由
- **行为**: 验证 JWT token，设置用户信息到上下文
- **失败处理**: 返回 401 Unauthorized

### AdminAuthMiddleware
- **用途**: 保护需要管理员权限的路由
- **行为**: 验证 JWT token 并检查角色是否为 "admin"
- **失败处理**: 返回 401 Unauthorized 或 403 Forbidden

### OptionalAuthMiddleware
- **用途**: 可选认证，不强制要求 token
- **行为**: 如果有有效 token 则设置用户信息，否则继续执行
- **失败处理**: 不返回错误，继续执行

## 配置说明

在 `configs/config.yaml` 中配置 JWT 相关参数：

```yaml
jwt:
  secret: "your-secret-key-change-in-production"  # JWT 签名密钥
  expiration: 24  # token 过期时间（小时）
```

## 安全建议

1. **密钥安全**: 在生产环境中使用强随机密钥，不要使用默认值
2. **HTTPS**: 在生产环境中始终使用 HTTPS 传输 token
3. **过期时间**: 根据安全需求设置合适的 token 过期时间
4. **密码策略**: 实施强密码策略，要求用户使用复杂密码
5. **Token 存储**: 客户端应安全存储 token，避免 XSS 攻击

## 测试

运行测试以验证功能：

```bash
# 测试 JWT 和密码功能
go test ./pkg/utils/... -v

# 测试认证中间件
go test ./internal/middleware/... -v
```

## 错误处理

常见错误及处理方式：

- **token 格式错误**: 检查 Authorization header 格式是否为 "Bearer <token>"
- **token 过期**: 使用 RefreshToken 方法刷新 token
- **签名验证失败**: 检查密钥是否正确
- **权限不足**: 检查用户角色是否匹配路由要求