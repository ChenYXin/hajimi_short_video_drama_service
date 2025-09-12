package service

import (
	"gin-mysql-api/internal/repository"
	"gin-mysql-api/pkg/config"
	"gin-mysql-api/pkg/utils"

	"github.com/go-redis/redis/v8"
)

// Container 服务容器
type Container struct {
	UserService  UserService
	DramaService DramaService
	AdminService AdminService
	AuthService  AuthService
	CacheService CacheService
	FileService  FileService
}

// NewContainer 创建新的服务容器
func NewContainer(
	cfg *config.Config,
	repos *repository.Repository,
	redisClient *redis.Client,
	jwtManager *utils.JWTManager,
) *Container {
	// 创建缓存服务
	cacheService := NewCacheService(redisClient)

	// 创建文件服务
	fileService := NewFileService(
		cfg.Upload.UploadPath,
		"http://localhost:8080", // 这里应该从配置中获取
		int64(cfg.Upload.MaxSize)*1024*1024, // 转换为字节
		cfg.Upload.AllowedTypes,
	)

	// 创建用户服务
	userService := NewUserService(repos.User, jwtManager)

	// 创建短剧服务
	dramaService := NewDramaService(repos.Drama, repos.Episode, cacheService)

	// 创建管理服务
	adminService := NewAdminService(
		repos.Admin,
		repos.Drama,
		repos.Episode,
		jwtManager,
		cacheService,
	)

	// 创建认证服务
	authService := NewAuthService(repos.User, repos.Admin, jwtManager)

	return &Container{
		UserService:  userService,
		DramaService: dramaService,
		AdminService: adminService,
		AuthService:  authService,
		CacheService: cacheService,
		FileService:  fileService,
	}
}