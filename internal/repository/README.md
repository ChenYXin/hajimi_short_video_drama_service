# Repository 层

Repository 层负责数据访问逻辑，提供了对数据库操作的抽象接口。

## 架构设计

Repository 层采用接口与实现分离的设计模式：

- `interfaces.go`: 定义所有仓库接口
- `*_repository.go`: 各个仓库的具体实现
- `repository.go`: 仓库管理器，统一管理所有仓库实例

## 仓库接口

### UserRepository
用户数据访问接口，提供用户的 CRUD 操作：
- 创建、查询、更新、删除用户
- 根据邮箱、用户名查询用户
- 检查邮箱、用户名是否已存在
- 分页获取用户列表

### DramaRepository
短剧数据访问接口，提供短剧的 CRUD 操作：
- 创建、查询、更新、删除短剧
- 根据类型筛选短剧
- 获取活跃状态的短剧列表
- 增加观看次数
- 支持预加载剧集信息

### EpisodeRepository
剧集数据访问接口，提供剧集的 CRUD 操作：
- 创建、查询、更新、删除剧集
- 根据短剧ID获取剧集列表
- 获取最大剧集号
- 检查剧集号是否已存在
- 增加观看次数
- 支持预加载短剧信息

### AdminRepository
管理员数据访问接口，提供管理员的 CRUD 操作：
- 创建、查询、更新、删除管理员
- 根据邮箱、用户名查询管理员
- 检查邮箱、用户名是否已存在
- 分页获取管理员列表

## 使用示例

```go
package main

import (
    "gin-mysql-api/internal/repository"
    "gin-mysql-api/internal/models"
    "gorm.io/gorm"
)

func main() {
    // 假设已经初始化了数据库连接
    var db *gorm.DB
    
    // 创建仓库管理器
    repo := repository.NewRepository(db)
    
    // 使用用户仓库
    user := &models.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "hashedpassword",
    }
    
    // 创建用户
    err := repo.User.Create(user)
    if err != nil {
        // 处理错误
    }
    
    // 根据ID查询用户
    foundUser, err := repo.User.GetByID(user.ID)
    if err != nil {
        // 处理错误
    }
    
    // 检查邮箱是否已存在
    exists, err := repo.User.ExistsByEmail("test@example.com")
    if err != nil {
        // 处理错误
    }
    
    // 使用短剧仓库
    dramas, total, err := repo.Drama.GetActiveList(0, 10)
    if err != nil {
        // 处理错误
    }
    
    // 增加短剧观看次数
    err = repo.Drama.IncrementViewCount(1)
    if err != nil {
        // 处理错误
    }
}
```

## 错误处理

所有仓库方法都返回 error，调用方需要适当处理错误：

- 当记录不存在时，Get* 方法返回 `nil, nil`
- 数据库操作失败时，返回相应的 error
- 使用 `errors.Is(err, gorm.ErrRecordNotFound)` 判断记录是否不存在

## 事务支持

所有仓库方法都接受 `*gorm.DB` 参数，支持在事务中使用：

```go
// 在事务中使用仓库
err := db.Transaction(func(tx *gorm.DB) error {
    userRepo := repository.NewUserRepository(tx)
    
    // 在事务中执行操作
    err := userRepo.Create(user)
    if err != nil {
        return err // 自动回滚
    }
    
    return nil // 自动提交
})
```

## 测试

Repository 层包含完整的测试用例：
- 接口实现验证
- 方法签名测试
- 数据库操作测试（需要测试数据库）

运行测试：
```bash
go test ./internal/repository/...
```