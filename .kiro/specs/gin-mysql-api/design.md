# Design Document

## Overview

短剧服务平台采用分层架构设计，基于 Gin 框架构建 RESTful API 服务和 Web 管理界面。系统使用 MySQL 作为主数据库，GORM 作为 ORM 框架，JWT 进行用户认证，支持文件上传和视频流媒体服务。

## Architecture

### 系统架构图

```
┌─────────────────┐    ┌─────────────────┐
│   Mobile App    │    │   Web Admin     │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          │ HTTP/JSON            │ HTTP/HTML
          │                      │
┌─────────▼──────────────────────▼───────┐
│            Gin Web Server              │
├────────────────────────────────────────┤
│  ┌──────────────┐  ┌─────────────────┐ │
│  │   API Routes │  │   Admin Routes  │ │
│  └──────────────┘  └─────────────────┘ │
├────────────────────────────────────────┤
│           Middleware Layer             │
│  (Auth, CORS, Logging, Error Handler)  │
├────────────────────────────────────────┤
│           Service Layer                │
│  (Business Logic & Validation)         │
├────────────────────────────────────────┤
│           Repository Layer             │
│         (Data Access Layer)            │
└─────────┬──────────────────────────────┘
          │
┌─────────▼───────┐    ┌─────────────────┐    ┌─────────────────┐
│   MySQL DB      │    │   Redis Cache   │    │   File Storage  │
│   (GORM ORM)    │    │   (go-redis)    │    │   (Local/Cloud) │
└─────────────────┘    └─────────────────┘    └─────────────────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │     Prometheus Server     │
                    │    (Metrics Collection)   │
                    └───────────────────────────┘
```

### 技术栈

- **Web Framework**: Gin (Go HTTP web framework)
- **Database**: MySQL 8.0+
- **ORM**: GORM v2
- **Cache**: Redis 7.0+
- **Authentication**: JWT (JSON Web Tokens)
- **Configuration**: Viper
- **Logging**: Logrus
- **Validation**: go-playground/validator
- **File Upload**: Multipart form handling
- **Template Engine**: HTML/template (for admin interface)
- **Monitoring**: Prometheus + Grafana
- **Metrics**: prometheus/client_golang

## Components and Interfaces

### 1. 核心模型 (Models)

```go
// User 用户模型
type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Username  string    `gorm:"uniqueIndex;size:50" json:"username"`
    Email     string    `gorm:"uniqueIndex;size:100" json:"email"`
    Password  string    `gorm:"size:255" json:"-"`
    Phone     string    `gorm:"size:20" json:"phone"`
    Avatar    string    `gorm:"size:255" json:"avatar"`
    IsActive  bool      `gorm:"default:true" json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Drama 短剧模型
type Drama struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Title       string    `gorm:"size:200;not null" json:"title"`
    Description string    `gorm:"type:text" json:"description"`
    CoverImage  string    `gorm:"size:255" json:"cover_image"`
    Director    string    `gorm:"size:100" json:"director"`
    Actors      string    `gorm:"size:500" json:"actors"`
    Genre       string    `gorm:"size:100" json:"genre"`
    Status      string    `gorm:"size:20;default:'active'" json:"status"`
    ViewCount   int64     `gorm:"default:0" json:"view_count"`
    Episodes    []Episode `gorm:"foreignKey:DramaID" json:"episodes,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Episode 剧集模型
type Episode struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    DramaID     uint      `gorm:"not null;index" json:"drama_id"`
    Title       string    `gorm:"size:200;not null" json:"title"`
    EpisodeNum  int       `gorm:"not null" json:"episode_num"`
    Duration    int       `gorm:"not null" json:"duration"` // 秒
    VideoURL    string    `gorm:"size:500" json:"video_url"`
    Thumbnail   string    `gorm:"size:255" json:"thumbnail"`
    Status      string    `gorm:"size:20;default:'active'" json:"status"`
    ViewCount   int64     `gorm:"default:0" json:"view_count"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Admin 管理员模型
type Admin struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Username  string    `gorm:"uniqueIndex;size:50" json:"username"`
    Email     string    `gorm:"uniqueIndex;size:100" json:"email"`
    Password  string    `gorm:"size:255" json:"-"`
    Role      string    `gorm:"size:20;default:'admin'" json:"role"`
    IsActive  bool      `gorm:"default:true" json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 2. 服务层接口 (Services)

```go
// UserService 用户服务接口
type UserService interface {
    Register(req RegisterRequest) (*User, error)
    Login(req LoginRequest) (string, error) // 返回 JWT token
    GetProfile(userID uint) (*User, error)
    UpdateProfile(userID uint, req UpdateProfileRequest) (*User, error)
}

// DramaService 短剧服务接口
type DramaService interface {
    GetDramas(page, pageSize int, genre string) (*PaginatedDramas, error)
    GetDramaByID(id uint) (*Drama, error)
    GetEpisodesByDramaID(dramaID uint) ([]Episode, error)
    GetEpisodeByID(id uint) (*Episode, error)
    IncrementViewCount(dramaID, episodeID uint) error
}

// AdminService 管理服务接口
type AdminService interface {
    Login(req AdminLoginRequest) (string, error)
    CreateDrama(req CreateDramaRequest) (*Drama, error)
    UpdateDrama(id uint, req UpdateDramaRequest) (*Drama, error)
    DeleteDrama(id uint) error
    CreateEpisode(req CreateEpisodeRequest) (*Episode, error)
    UpdateEpisode(id uint, req UpdateEpisodeRequest) (*Episode, error)
    DeleteEpisode(id uint) error
    UploadFile(file multipart.File, header *multipart.FileHeader) (string, error)
}

// CacheService Redis 缓存服务接口
type CacheService interface {
    Set(key string, value interface{}, expiration time.Duration) error
    Get(key string) (string, error)
    Delete(key string) error
    Exists(key string) (bool, error)
    SetJSON(key string, value interface{}, expiration time.Duration) error
    GetJSON(key string, dest interface{}) error
}

// MetricsService Prometheus 监控服务接口
type MetricsService interface {
    IncrementCounter(name string, labels map[string]string)
    RecordHistogram(name string, value float64, labels map[string]string)
    SetGauge(name string, value float64, labels map[string]string)
    RecordAPIRequest(method, path string, statusCode int, duration float64)
}
```

### 3. 存储层接口 (Repositories)

```go
// UserRepository 用户数据访问接口
type UserRepository interface {
    Create(user *User) error
    GetByID(id uint) (*User, error)
    GetByEmail(email string) (*User, error)
    GetByUsername(username string) (*User, error)
    Update(user *User) error
    Delete(id uint) error
}

// DramaRepository 短剧数据访问接口
type DramaRepository interface {
    Create(drama *Drama) error
    GetByID(id uint) (*Drama, error)
    GetList(offset, limit int, genre string) ([]Drama, int64, error)
    Update(drama *Drama) error
    Delete(id uint) error
    IncrementViewCount(id uint) error
}

// EpisodeRepository 剧集数据访问接口
type EpisodeRepository interface {
    Create(episode *Episode) error
    GetByID(id uint) (*Episode, error)
    GetByDramaID(dramaID uint) ([]Episode, error)
    Update(episode *Episode) error
    Delete(id uint) error
    IncrementViewCount(id uint) error
}
```

## Data Models

### 数据库表结构

```sql
-- 用户表
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    avatar VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 短剧表
CREATE TABLE dramas (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    cover_image VARCHAR(255),
    director VARCHAR(100),
    actors VARCHAR(500),
    genre VARCHAR(100),
    status VARCHAR(20) DEFAULT 'active',
    view_count BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_genre (genre),
    INDEX idx_status (status)
);

-- 剧集表
CREATE TABLE episodes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    drama_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(200) NOT NULL,
    episode_num INT NOT NULL,
    duration INT NOT NULL,
    video_url VARCHAR(500),
    thumbnail VARCHAR(255),
    status VARCHAR(20) DEFAULT 'active',
    view_count BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (drama_id) REFERENCES dramas(id) ON DELETE CASCADE,
    INDEX idx_drama_id (drama_id),
    INDEX idx_episode_num (episode_num)
);

-- 管理员表
CREATE TABLE admins (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'admin',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## Error Handling

### 错误类型定义

```go
// 自定义错误类型
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

// 常见错误定义
var (
    ErrUserNotFound     = &AppError{Code: 404, Message: "用户不存在"}
    ErrDramaNotFound    = &AppError{Code: 404, Message: "短剧不存在"}
    ErrEpisodeNotFound  = &AppError{Code: 404, Message: "剧集不存在"}
    ErrInvalidCredentials = &AppError{Code: 401, Message: "用户名或密码错误"}
    ErrUnauthorized     = &AppError{Code: 401, Message: "未授权访问"}
    ErrValidationFailed = &AppError{Code: 400, Message: "参数验证失败"}
    ErrInternalServer   = &AppError{Code: 500, Message: "服务器内部错误"}
)
```

### 全局错误处理中间件

```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            switch e := err.Err.(type) {
            case *AppError:
                c.JSON(e.Code, gin.H{
                    "error": e.Message,
                    "details": e.Details,
                })
            default:
                c.JSON(500, gin.H{
                    "error": "服务器内部错误",
                })
            }
        }
    }
}
```

## Testing Strategy

### 测试层级

1. **单元测试**
   - Repository 层测试（使用内存数据库）
   - Service 层测试（使用 mock repository）
   - Utility 函数测试

2. **集成测试**
   - API 端点测试
   - 数据库集成测试
   - 中间件测试

3. **端到端测试**
   - 完整用户流程测试
   - 管理员操作流程测试

### 测试工具

- **测试框架**: Go 标准库 testing
- **断言库**: testify/assert
- **Mock 生成**: gomock
- **HTTP 测试**: httptest
- **数据库测试**: SQLite 内存数据库

### 测试数据管理

```go
// 测试数据工厂
type TestDataFactory struct {
    db *gorm.DB
}

func (f *TestDataFactory) CreateUser() *User {
    user := &User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "hashedpassword",
    }
    f.db.Create(user)
    return user
}

func (f *TestDataFactory) CreateDrama() *Drama {
    drama := &Drama{
        Title:       "测试短剧",
        Description: "测试描述",
        Genre:       "喜剧",
    }
    f.db.Create(drama)
    return drama
}
```