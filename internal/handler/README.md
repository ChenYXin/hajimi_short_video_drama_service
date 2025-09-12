# Handler 层使用指南

Handler 层是 HTTP 请求处理的入口点，负责接收 HTTP 请求、参数验证、调用业务逻辑服务，并返回 HTTP 响应。

## 架构设计

```
HTTP Request
     ↓
Handler Layer  ← HTTP 请求处理层
     ↓
Service Layer  ← 业务逻辑层
     ↓
Repository Layer
```

## 处理器组件

### 1. BaseHandler - 基础处理器

提供所有处理器的通用功能：

- **响应格式化**: 统一的成功/错误响应格式
- **参数验证**: 请求参数的验证和错误处理
- **分页处理**: 分页参数的解析和验证
- **上下文工具**: 从 JWT 中间件获取用户信息

```go
// 使用示例
type MyHandler struct {
    *BaseHandler
    myService service.MyService
}

func (h *MyHandler) MyEndpoint(c *gin.Context) {
    // 验证请求参数
    var req MyRequest
    if err := h.ValidateRequest(c, &req); err != nil {
        h.ValidationErrorResponse(c, err)
        return
    }

    // 获取用户信息
    userID, exists := h.GetUserIDFromContext(c)
    if !exists {
        h.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
        return
    }

    // 调用业务逻辑
    result, err := h.myService.DoSomething(userID, req)
    if err != nil {
        h.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    // 返回成功响应
    h.SuccessResponseWithMessage(c, "操作成功", result)
}
```

### 2. AuthHandler - 认证处理器

处理用户和管理员的认证相关请求：

- **用户注册**: `POST /api/auth/register`
- **用户登录**: `POST /api/auth/login`
- **管理员登录**: `POST /api/auth/admin/login`
- **刷新令牌**: `POST /api/auth/refresh`

```go
// 使用示例
authHandler := handler.NewAuthHandler(authService)

// 用户注册
router.POST("/api/auth/register", authHandler.Register)

// 用户登录
router.POST("/api/auth/login", authHandler.Login)
```

### 3. UserHandler - 用户处理器

处理用户资料相关的请求：

- **获取用户资料**: `GET /api/user/profile`
- **更新用户资料**: `PUT /api/user/profile`

```go
// 使用示例
userHandler := handler.NewUserHandler(userService)

// 需要认证的用户路由
userGroup := router.Group("/api/user")
userGroup.Use(middleware.AuthMiddleware(jwtManager))
{
    userGroup.GET("/profile", userHandler.GetProfile)
    userGroup.PUT("/profile", userHandler.UpdateProfile)
}
```

### 4. DramaHandler - 短剧处理器

处理短剧和剧集相关的请求：

- **获取短剧列表**: `GET /api/dramas`
- **获取短剧详情**: `GET /api/dramas/:id`
- **获取短剧剧集**: `GET /api/dramas/:id/episodes`
- **获取剧集详情**: `GET /api/episodes/:id`
- **搜索短剧**: `GET /api/dramas/search`
- **热门短剧**: `GET /api/dramas/popular`

```go
// 使用示例
dramaHandler := handler.NewDramaHandler(dramaService)

// 公开的短剧路由
dramaGroup := router.Group("/api/dramas")
{
    dramaGroup.GET("", dramaHandler.GetDramas)
    dramaGroup.GET("/:id", dramaHandler.GetDramaByID)
    dramaGroup.GET("/search", dramaHandler.SearchDramas)
}
```

### 5. AdminHandler - 管理员处理器

处理管理员相关的请求：

- **短剧管理**: 创建、更新、删除短剧
- **剧集管理**: 创建、更新、删除剧集
- **用户管理**: 查看、激活、禁用用户

```go
// 使用示例
adminHandler := handler.NewAdminHandler(adminService, userService)

// 需要管理员权限的路由
adminGroup := router.Group("/api/admin")
adminGroup.Use(middleware.AdminAuthMiddleware(jwtManager))
{
    adminGroup.POST("/dramas", adminHandler.CreateDrama)
    adminGroup.PUT("/dramas/:id", adminHandler.UpdateDrama)
    adminGroup.DELETE("/dramas/:id", adminHandler.DeleteDrama)
}
```

### 6. FileHandler - 文件处理器

处理文件上传和管理：

- **文件上传**: `POST /api/upload`
- **文件删除**: `DELETE /api/upload`

```go
// 使用示例
fileHandler := handler.NewFileHandler(fileService)

// 需要认证的文件路由
uploadGroup := router.Group("/api/upload")
uploadGroup.Use(middleware.AuthMiddleware(jwtManager))
{
    uploadGroup.POST("", fileHandler.UploadFile)
    uploadGroup.DELETE("", fileHandler.DeleteFile)
}
```

### 7. HealthHandler - 健康检查处理器

提供系统健康检查端点：

- **健康检查**: `GET /health`
- **就绪检查**: `GET /ready`
- **存活检查**: `GET /live`

```go
// 使用示例
healthHandler := handler.NewHealthHandler()

router.GET("/health", healthHandler.HealthCheck)
router.GET("/ready", healthHandler.ReadinessCheck)
router.GET("/live", healthHandler.LivenessCheck)
```

## 响应格式

### 统一响应结构

所有 API 响应都使用统一的格式：

```go
type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   interface{} `json:"error,omitempty"`
}
```

### 成功响应示例

```json
{
    "success": true,
    "message": "操作成功",
    "data": {
        "id": 1,
        "username": "testuser",
        "email": "test@example.com"
    }
}
```

### 错误响应示例

```json
{
    "success": false,
    "message": "参数验证失败",
    "error": [
        "username 是必填字段",
        "email 必须是有效的邮箱地址"
    ]
}
```

## 参数验证

### 验证标签

使用 `validator` 包的标签进行参数验证：

```go
type RegisterRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
    Phone    string `json:"phone" validate:"omitempty,len=11"`
}
```

### 常用验证标签

- `required`: 必填字段
- `email`: 邮箱格式
- `min=n`: 最小长度
- `max=n`: 最大长度
- `len=n`: 固定长度
- `oneof=a b c`: 枚举值
- `omitempty`: 可选字段

## 分页处理

### 分页参数

- `page`: 页码，默认为 1
- `page_size`: 每页数量，默认为 20，最大为 100

### 分页响应

```go
type PaginatedResponse struct {
    Data        []interface{} `json:"data"`
    Total       int64         `json:"total"`
    Page        int           `json:"page"`
    PageSize    int           `json:"page_size"`
    TotalPages  int           `json:"total_pages"`
    HasNext     bool          `json:"has_next"`
    HasPrevious bool          `json:"has_previous"`
}
```

## 错误处理

### HTTP 状态码使用

- `200 OK`: 成功
- `400 Bad Request`: 参数错误、业务逻辑错误
- `401 Unauthorized`: 未认证
- `403 Forbidden`: 权限不足
- `404 Not Found`: 资源不存在
- `500 Internal Server Error`: 服务器内部错误

### 错误处理最佳实践

```go
func (h *MyHandler) MyEndpoint(c *gin.Context) {
    // 1. 参数验证错误
    if err := h.ValidateRequest(c, &req); err != nil {
        h.ValidationErrorResponse(c, err) // 400 + 详细验证错误
        return
    }

    // 2. 认证错误
    userID, exists := h.GetUserIDFromContext(c)
    if !exists {
        h.ErrorResponse(c, http.StatusUnauthorized, "用户未认证") // 401
        return
    }

    // 3. 业务逻辑错误
    result, err := h.service.DoSomething(req)
    if err != nil {
        h.ErrorResponse(c, http.StatusBadRequest, err.Error()) // 400
        return
    }

    // 4. 资源不存在
    if result == nil {
        h.ErrorResponse(c, http.StatusNotFound, "资源不存在") // 404
        return
    }

    h.SuccessResponse(c, result)
}
```

## 中间件集成

### JWT 认证中间件

```go
// 需要用户认证
userGroup.Use(middleware.AuthMiddleware(jwtManager))

// 需要管理员权限
adminGroup.Use(middleware.AdminAuthMiddleware(jwtManager))

// 可选认证
publicGroup.Use(middleware.OptionalAuthMiddleware(jwtManager))
```

### 获取认证信息

```go
func (h *MyHandler) MyEndpoint(c *gin.Context) {
    // 获取用户ID
    userID, exists := h.GetUserIDFromContext(c)
    
    // 获取用户角色
    role, exists := h.GetUserRoleFromContext(c)
    
    // 使用认证信息
    if role == "admin" {
        // 管理员逻辑
    }
}
```

## 测试策略

### 单元测试

使用 mock 服务测试处理器逻辑：

```go
func TestMyHandler(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    mockService := new(MockMyService)
    handler := NewMyHandler(mockService)
    
    // 设置 mock 期望
    mockService.On("DoSomething", mock.Anything).Return(result, nil)
    
    // 创建测试请求
    req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(reqBody))
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = req
    
    // 执行处理器
    handler.MyEndpoint(c)
    
    // 验证结果
    assert.Equal(t, http.StatusOK, w.Code)
    mockService.AssertExpectations(t)
}
```

### 集成测试

测试完整的 HTTP 请求流程：

```bash
# 运行所有处理器测试
go test ./internal/handler/... -v

# 运行特定处理器测试
go test ./internal/handler/auth_handler_test.go -v
```

## 性能优化

### 响应缓存

对于不经常变化的数据，可以在处理器层添加缓存：

```go
func (h *DramaHandler) GetDramas(c *gin.Context) {
    // 检查缓存
    cacheKey := fmt.Sprintf("dramas:%s", c.Request.URL.RawQuery)
    if cached := h.getFromCache(cacheKey); cached != nil {
        h.SuccessResponse(c, cached)
        return
    }
    
    // 获取数据
    dramas, err := h.dramaService.GetDramas(page, pageSize, genre)
    if err != nil {
        h.ErrorResponse(c, http.StatusInternalServerError, "获取失败")
        return
    }
    
    // 缓存结果
    h.setToCache(cacheKey, dramas, 5*time.Minute)
    
    h.SuccessResponse(c, dramas)
}
```

### 异步处理

对于不影响响应的操作，可以异步执行：

```go
func (h *DramaHandler) GetDramaByID(c *gin.Context) {
    drama, err := h.dramaService.GetDramaByID(id)
    if err != nil {
        h.ErrorResponse(c, http.StatusNotFound, "短剧不存在")
        return
    }
    
    // 异步增加观看次数
    go h.dramaService.IncrementViewCount(id)
    
    h.SuccessResponse(c, drama)
}
```

## 最佳实践

1. **单一职责**: 每个处理器专注于特定的业务领域
2. **统一响应**: 使用统一的响应格式和错误处理
3. **参数验证**: 在处理器层进行完整的参数验证
4. **错误处理**: 提供清晰的错误信息和适当的 HTTP 状态码
5. **日志记录**: 记录关键操作和错误信息
6. **性能优化**: 合理使用缓存和异步处理
7. **测试覆盖**: 编写全面的单元测试和集成测试
8. **文档注释**: 使用 Swagger 注释生成 API 文档