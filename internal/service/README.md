# Service 层使用指南

Service 层是业务逻辑的核心，负责处理复杂的业务规则、数据验证、缓存管理和跨领域的操作。

## 架构设计

```
Controller Layer
       ↓
   Service Layer  ← 业务逻辑层
       ↓
Repository Layer  ← 数据访问层
       ↓
   Database
```

## 服务组件

### 1. UserService - 用户服务

负责用户相关的业务逻辑：

- **用户注册**: 验证用户信息、检查重复、密码哈希
- **用户登录**: 验证凭据、生成 JWT token
- **用户资料管理**: 获取和更新用户信息
- **用户管理**: 管理员功能，用户列表、激活/禁用

```go
// 使用示例
userService := service.NewUserService(userRepo, jwtManager)

// 用户注册
user, err := userService.Register(models.RegisterRequest{
    Username: "testuser",
    Email:    "test@example.com",
    Password: "password123",
})

// 用户登录
response, err := userService.Login(models.LoginRequest{
    Email:    "test@example.com",
    Password: "password123",
})
```

### 2. DramaService - 短剧服务

负责短剧和剧集相关的业务逻辑：

- **短剧查询**: 分页列表、详情查询、搜索
- **剧集管理**: 获取剧集列表和详情
- **观看统计**: 增加观看次数
- **缓存管理**: 热门内容缓存

```go
// 使用示例
dramaService := service.NewDramaService(dramaRepo, episodeRepo, cacheService)

// 获取短剧列表
dramas, err := dramaService.GetDramas(1, 20, "喜剧")

// 获取短剧详情
drama, err := dramaService.GetDramaByID(1)

// 增加观看次数
err := dramaService.IncrementDramaViewCount(1)
```

### 3. AdminService - 管理服务

负责管理员相关的业务逻辑：

- **管理员认证**: 登录验证
- **内容管理**: 短剧和剧集的增删改查
- **管理员管理**: 创建管理员账户

```go
// 使用示例
adminService := service.NewAdminService(adminRepo, dramaRepo, episodeRepo, jwtManager, cacheService)

// 管理员登录
response, err := adminService.Login(models.AdminLoginRequest{
    Username: "admin",
    Password: "password",
})

// 创建短剧
drama, err := adminService.CreateDrama(models.CreateDramaRequest{
    Title:       "新短剧",
    Description: "描述",
    Genre:       "喜剧",
})
```

### 4. AuthService - 认证服务

统一的认证服务，整合用户和管理员认证：

- **用户认证**: 注册、登录
- **管理员认证**: 登录
- **Token 管理**: 刷新、验证

```go
// 使用示例
authService := service.NewAuthService(userRepo, adminRepo, jwtManager)

// 用户注册
user, err := authService.RegisterUser(registerReq)

// 用户登录
response, err := authService.LoginUser(loginReq)

// 管理员登录
response, err := authService.LoginAdmin(adminLoginReq)
```

### 5. CacheService - 缓存服务

Redis 缓存服务，提供高性能的数据缓存：

- **基础操作**: Set、Get、Delete、Exists
- **JSON 支持**: 自动序列化/反序列化
- **计数器**: 递增操作
- **模式匹配**: 批量删除

```go
// 使用示例
cacheService := service.NewCacheService(redisClient)

// 设置缓存
err := cacheService.Set("key", "value", time.Hour)

// 获取缓存
value, err := cacheService.Get("key")

// JSON 缓存
err := cacheService.SetJSON("user:1", user, time.Hour)
var user User
err := cacheService.GetJSON("user:1", &user)
```

### 6. FileService - 文件服务

文件上传和管理服务：

- **文件上传**: 支持多种文件类型
- **文件验证**: 类型和大小验证
- **文件管理**: 删除、URL 生成

```go
// 使用示例
fileService := service.NewFileService(uploadPath, baseURL, maxSize, allowedTypes)

// 上传文件
response, err := fileService.UploadFile(file, header, "avatar")

// 删除文件
err := fileService.DeleteFile("avatars/file.jpg")
```

## 服务容器

使用依赖注入容器管理所有服务：

```go
// 创建服务容器
container := service.NewContainer(config, repos, redisClient, jwtManager)

// 使用服务
user, err := container.UserService.Register(req)
dramas, err := container.DramaService.GetDramas(1, 20, "")
```

## 缓存策略

### 缓存键命名规范

- 短剧列表: `dramas:page:{page}:size:{size}:genre:{genre}`
- 短剧详情: `drama:{id}`
- 短剧含剧集: `drama_with_episodes:{id}`
- 剧集列表: `episodes:drama:{drama_id}:page:{page}:size:{size}`
- 剧集详情: `episode:{id}`
- 热门短剧: `popular_dramas:page:{page}:size:{size}`

### 缓存过期时间

- 短剧/剧集详情: 10 分钟
- 列表数据: 5 分钟
- 热门内容: 30 分钟

### 缓存失效策略

- 数据更新时自动清除相关缓存
- 使用模式匹配批量清除相关缓存

## 错误处理

### 统一错误格式

所有服务方法返回的错误都包含详细的错误信息：

```go
// 业务错误
return nil, errors.New("用户名已存在")

// 系统错误
return nil, fmt.Errorf("数据库操作失败: %w", err)
```

### 常见错误类型

- **验证错误**: 参数验证失败
- **业务错误**: 业务规则违反
- **系统错误**: 数据库、缓存等系统组件错误
- **认证错误**: 用户认证失败

## 测试策略

### 单元测试

使用 mock 对象测试业务逻辑：

```go
// 创建 mock repository
mockRepo := new(MockUserRepository)
mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

// 测试服务方法
userService := NewUserService(mockRepo, jwtManager)
user, err := userService.Register(req)

// 验证结果
assert.NoError(t, err)
mockRepo.AssertExpectations(t)
```

### 集成测试

测试服务与真实数据库和缓存的集成：

```bash
# 运行所有服务测试
go test ./internal/service/... -v

# 运行特定服务测试
go test ./internal/service/user_service_test.go -v
```

## 性能优化

### 缓存优化

1. **热点数据缓存**: 热门短剧、用户信息
2. **查询结果缓存**: 分页列表、搜索结果
3. **计算结果缓存**: 统计数据、排行榜

### 数据库优化

1. **批量操作**: 减少数据库往返次数
2. **预加载**: 使用 GORM 的 Preload 减少 N+1 查询
3. **索引优化**: 确保查询字段有适当索引

### 并发处理

1. **无状态设计**: 服务方法不依赖实例状态
2. **线程安全**: 使用线程安全的缓存客户端
3. **连接池**: 合理配置数据库和 Redis 连接池

## 最佳实践

1. **单一职责**: 每个服务专注于特定领域
2. **依赖注入**: 通过构造函数注入依赖
3. **接口设计**: 定义清晰的服务接口
4. **错误处理**: 提供详细的错误信息
5. **日志记录**: 记录关键操作和错误
6. **缓存管理**: 合理使用缓存提高性能
7. **测试覆盖**: 编写全面的单元测试