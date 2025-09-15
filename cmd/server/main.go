package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"gin-mysql-api/internal/repository"
	"gin-mysql-api/internal/router"
	"gin-mysql-api/internal/service"
	"gin-mysql-api/pkg/config"
	"gin-mysql-api/pkg/database"
	"gin-mysql-api/pkg/utils"
)

func main() {
	log.Println("Gin MySQL API Server - 启动中...")

	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 设置日志
	setupLogging(cfg)

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 连接数据库
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 连接Redis
	redisClient, err := database.NewRedisConnection(cfg)
	if err != nil {
		log.Fatalf("Redis连接失败: %v", err)
	}

	// 初始化仓储层
	userRepo := repository.NewUserRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	dramaRepo := repository.NewDramaRepository(db)
	episodeRepo := repository.NewEpisodeRepository(db)

	// 初始化JWT管理器
	jwtManager := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiration)

	// 初始化缓存服务
	cacheService := service.NewCacheService(redisClient)

	// 初始化服务层
	userService := service.NewUserService(userRepo, jwtManager)
	adminService := service.NewAdminService(adminRepo, dramaRepo, episodeRepo, jwtManager, cacheService)
	dramaService := service.NewDramaService(dramaRepo, episodeRepo, cacheService)
	fileService := service.NewFileService(cfg.Upload.UploadPath, "http://localhost:1800", int64(cfg.Upload.MaxSize*1024*1024), cfg.Upload.AllowedTypes)
	authService := service.NewAuthService(userRepo, adminRepo, jwtManager)

	// 初始化服务容器
	serviceContainer := &service.Container{
		UserService:  userService,
		AdminService: adminService,
		DramaService: dramaService,
		FileService:  fileService,
		AuthService:  authService,
	}

	// 设置路由
	r := router.NewRouter(jwtManager, serviceContainer).Setup()

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: r,
	}

	// 启动服务器
	go func() {
		log.Printf("服务器启动在 %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 优雅关闭服务器，等待5秒钟完成现有请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("服务器强制关闭: %v", err)
	}

	// 关闭数据库连接
	sqlDB, _ := db.DB()
	sqlDB.Close()

	// 关闭Redis连接
	redisClient.Close()

	log.Println("服务器已退出")
}

// setupLogging 设置日志配置
func setupLogging(cfg *config.Config) {
	// 根据配置设置日志级别和格式
	// 这里可以集成更复杂的日志库如logrus或zap
	if cfg.Logging.Output == "file" {
		// 确保日志目录存在
		if err := os.MkdirAll("logs", 0755); err != nil {
			log.Printf("创建日志目录失败: %v", err)
		}

		// 这里可以设置日志文件输出
		log.Printf("日志配置: 级别=%s, 格式=%s, 输出=%s",
			cfg.Logging.Level, cfg.Logging.Format, cfg.Logging.Output)
	}
}
