# 数据模型文档

## 概述

本目录包含了短剧服务平台的所有数据模型定义，使用 GORM 作为 ORM 框架。

## 文件结构

- `user.go` - 用户模型
- `drama.go` - 短剧模型  
- `episode.go` - 剧集模型
- `admin.go` - 管理员模型
- `base.go` - 基础模型和分页结构
- `dto.go` - 数据传输对象（请求/响应结构）
- `validator.go` - 数据验证工具
- `init.go` - 模型初始化和迁移

## 核心模型

### User（用户）
- 用户注册、登录和个人信息管理
- 包含用户名、邮箱、手机号等基本信息
- 支持软删除

### Drama（短剧）
- 短剧基本信息，包括标题、描述、导演、演员等
- 支持分类和状态管理
- 与 Episode 模型一对多关联
- 包含观看次数统计

### Episode（剧集）
- 短剧的具体剧集信息
- 包含视频 URL、缩略图、时长等
- 与 Drama 模型多对一关联
- 支持剧集编号和观看次数统计

### Admin（管理员）
- 管理员账户信息
- 支持角色权限管理
- 用于后台管理功能

## 数据验证

所有模型都使用 `github.com/go-playground/validator/v10` 进行数据验证：

```go
// 使用示例
user := &User{
    Username: "testuser",
    Email: "invalid-email", // 这会触发验证错误
}

if errors := ValidateStruct(user); len(errors) > 0 {
    // 处理验证错误
    for _, err := range errors {
        fmt.Printf("字段 %s: %s\n", err.Field, err.Message)
    }
}
```

## 序列化方法

每个模型都提供了 `ToJSON()` 方法用于安全的序列化：

```go
// 用户序列化（自动隐藏密码字段）
user := &User{...}
jsonData := user.ToJSON()

// 短剧序列化（可选择是否包含剧集）
drama := &Drama{...}
basicData := drama.ToJSON()
withEpisodes := drama.ToJSONWithEpisodes()
```

## 数据库迁移

使用 `init.go` 中的函数进行数据库初始化：

```go
// 自动迁移所有模型
err := AutoMigrate(db)

// 创建额外索引
err = CreateIndexes(db)

// 初始化种子数据
err = SeedData(db)
```

## 分页支持

使用 `PaginationRequest` 和 `PaginationResponse` 进行分页：

```go
req := &PaginationRequest{
    Page: 1,
    PageSize: 10,
}

offset := req.GetOffset()
limit := req.GetLimit()

// 查询数据...
response := NewPaginationResponse(req.Page, req.PageSize, total, data)
```

## 关联关系

### Drama -> Episodes (一对多)
```go
// 预加载剧集
var drama Drama
db.Preload("Episodes").First(&drama, id)

// 通过短剧 ID 查询剧集
var episodes []Episode
db.Where("drama_id = ?", dramaID).Find(&episodes)
```

### Episode -> Drama (多对一)
```go
// 预加载短剧信息
var episode Episode
db.Preload("Drama").First(&episode, id)
```

## 最佳实践

1. **验证**: 始终在创建/更新前验证数据
2. **序列化**: 使用 `ToJSON()` 方法避免暴露敏感信息
3. **软删除**: 使用 GORM 的软删除功能而不是物理删除
4. **索引**: 为经常查询的字段创建索引
5. **关联**: 合理使用预加载避免 N+1 查询问题

## 状态值说明

### Drama.Status / Episode.Status
- `active`: 激活状态，正常显示
- `inactive`: 非激活状态，不显示但保留数据
- `draft`: 草稿状态，未发布

### Admin.Role
- `admin`: 普通管理员
- `super_admin`: 超级管理员，拥有所有权限