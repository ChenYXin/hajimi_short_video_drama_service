# 测试工具包

本包提供了项目中单元测试和集成测试所需的工具和辅助函数。

## 文件说明

### config.go
提供测试环境的配置管理：
- `GetTestConfig()`: 获取测试专用的配置，包括测试数据库、Redis、JWT等配置

### database.go
提供测试数据库的管理功能：
- `SetupTestDB()`: 设置测试数据库连接
- `CleanupTestDB()`: 清理测试数据库数据
- `TruncateTable()`: 清空指定表
- `BeginTransaction()`: 开始事务
- `RollbackTransaction()`: 回滚事务
- `CommitTransaction()`: 提交事务
- `WithTransaction()`: 在事务中执行函数

### factory.go
提供测试数据工厂，用于创建测试所需的模型实例：
- `UserFactory`: 用户数据工厂
- `DramaFactory`: 短剧数据工厂
- `EpisodeFactory`: 剧集数据工厂
- `AdminFactory`: 管理员数据工厂
- `Factory`: 所有工厂的集合

### helpers.go
提供HTTP测试辅助函数：
- `HTTPTestHelper`: HTTP测试助手，提供GET、POST、PUT、DELETE等方法
- `AssertStatusCode()`: 断言HTTP状态码
- `AssertJSONResponse()`: 断言JSON响应
- `AssertAPIResponse()`: 断言API响应格式
- `GetAuthHeader()`: 获取认证头
- 各种断言辅助函数

## 使用示例

### 基本测试设置

```go
func TestSomething(t *testing.T) {
    // 设置测试数据库
    db := testutil.SetupTestDB()
    defer func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
    }()
    
    // 清理测试数据
    testutil.CleanupTestDB(db)
    
    // 创建测试工厂
    factory := testutil.NewFactory()
    
    // 创建测试用户
    user := factory.User.CreateUser()
    
    // 进行测试...
}
```

### HTTP测试

```go
func TestHTTPEndpoint(t *testing.T) {
    router := gin.New()
    router.GET("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "success"})
    })
    
    helper := testutil.NewHTTPTestHelper(t, router)
    
    // 发送GET请求
    w := helper.GET("/test")
    
    // 断言响应
    helper.AssertStatusCode(w, 200)
    helper.AssertSuccessResponse(w, "success")
}
```

### 事务测试

```go
func TestWithTransaction(t *testing.T) {
    db := testutil.SetupTestDB()
    
    err := testutil.WithTransaction(db, func(tx *gorm.DB) error {
        // 在事务中执行操作
        user := &models.User{Username: "test"}
        return tx.Create(user).Error
    })
    
    assert.NoError(t, err)
}
```

## 测试数据工厂使用

### 创建基本数据

```go
factory := testutil.NewFactory()

// 创建用户
user := factory.User.CreateUser()

// 创建带自定义属性的用户
user := factory.User.CreateUser(func(u *models.User) {
    u.Username = "customuser"
    u.Email = "custom@example.com"
})

// 批量创建用户
users := factory.User.CreateUsers(5)
```

### 创建关联数据

```go
// 创建短剧
drama := factory.Drama.CreateDrama()

// 为短剧创建剧集
episode := factory.Episode.CreateEpisode(drama.ID)

// 批量创建剧集
episodes := factory.Episode.CreateEpisodes(drama.ID, 10)
```

## 配置说明

测试配置使用独立的数据库和Redis实例，避免与开发环境冲突：

- 数据库: `gin_mysql_api_test`
- Redis DB: `1` (而不是默认的0)
- JWT Secret: 测试专用密钥
- 文件上传: `./test_uploads` 目录

## 最佳实践

1. **测试隔离**: 每个测试用例都应该清理自己的数据
2. **使用工厂**: 使用数据工厂创建测试数据，而不是手动构造
3. **事务回滚**: 对于数据库测试，考虑使用事务回滚来保持数据库清洁
4. **Mock外部依赖**: 对于外部服务（如Redis、第三方API），使用Mock
5. **并发安全**: 确保测试可以并发运行而不互相影响

## 运行测试

```bash
# 运行所有测试
make test

# 运行单元测试
make test-unit

# 生成覆盖率报告
make test-coverage

# 运行特定包的测试
go test ./internal/repository/... -v

# 运行带竞态检测的测试
make test-race
```