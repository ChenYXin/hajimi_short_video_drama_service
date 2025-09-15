package router

import (
	"gin-mysql-api/internal/handler"
	"gin-mysql-api/internal/middleware"
	"gin-mysql-api/internal/service"
	"gin-mysql-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Router 路由配置
type Router struct {
	engine     *gin.Engine
	jwtManager *utils.JWTManager
	services   *service.Container
}

// NewRouter 创建新的路由器
func NewRouter(jwtManager *utils.JWTManager, services *service.Container) *Router {
	engine := gin.New()

	return &Router{
		engine:     engine,
		jwtManager: jwtManager,
		services:   services,
	}
}

// Setup 设置路由
func (r *Router) Setup() *gin.Engine {
	// 设置中间件
	r.setupMiddleware()

	// 设置路由
	r.setupRoutes()

	return r.engine
}

// setupMiddleware 设置中间件
func (r *Router) setupMiddleware() {
	// 错误处理中间件
	r.engine.Use(middleware.ErrorHandler())

	// 请求 ID 中间件
	r.engine.Use(middleware.RequestIDMiddleware())

	// 日志中间件
	r.engine.Use(middleware.Logger())

	// 安全中间件
	r.engine.Use(middleware.Security())

	// CORS 中间件
	corsConfig := middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "Authorization",
			"X-Request-ID",
		},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           86400,
	}
	r.engine.Use(middleware.CORS(corsConfig))

	// 设置CSP头部，允许加载外部图片和视频
	r.engine.Use(func(c *gin.Context) {
		c.Header("Content-Security-Policy", "default-src 'self'; img-src 'self' data: https: http:; media-src 'self' https: http:; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; font-src 'self' data: https:;")
		c.Next()
	})

	// 请求大小限制中间件
	r.engine.Use(middleware.RequestSizeLimit(10 * 1024 * 1024)) // 10MB

	// 简单限流中间件
	r.engine.Use(middleware.SimpleRateLimit())

	// 404 和 405 处理
	r.engine.NoRoute(middleware.NotFoundHandler())
	r.engine.NoMethod(middleware.MethodNotAllowedHandler())
}

// setupRoutes 设置路由
func (r *Router) setupRoutes() {
	// 创建处理器
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(r.services.AuthService)
	userHandler := handler.NewUserHandler(r.services.UserService)
	dramaHandler := handler.NewDramaHandler(r.services.DramaService)
	adminHandler := handler.NewAdminHandler(r.services.AdminService, r.services.UserService)
	fileHandler := handler.NewFileHandler(r.services.FileService)

	// 健康检查路由
	r.engine.GET("/health", healthHandler.HealthCheck)
	r.engine.GET("/ready", healthHandler.ReadinessCheck)
	r.engine.GET("/live", healthHandler.LivenessCheck)

	// API 路由组
	api := r.engine.Group("/api")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/admin/login", authHandler.AdminLogin)

			// 需要认证的认证路由
			authProtected := auth.Group("")
			authProtected.Use(middleware.AuthMiddleware(r.jwtManager))
			{
				authProtected.POST("/refresh", authHandler.RefreshToken)
			}
		}

		// 用户路由
		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
		}

		// 短剧路由（公开）
		dramas := api.Group("/dramas")
		{
			dramas.GET("", dramaHandler.GetDramas)
			dramas.GET("/search", dramaHandler.SearchDramas)
			dramas.GET("/popular", dramaHandler.GetPopularDramas)
			dramas.GET("/:id", dramaHandler.GetDramaByID)
			dramas.GET("/:id/episodes", dramaHandler.GetDramaWithEpisodes)
			dramas.GET("/:id/episodes/list", dramaHandler.GetEpisodesByDramaID)
		}

		// 剧集路由（公开）
		episodes := api.Group("/episodes")
		{
			episodes.GET("/:id", dramaHandler.GetEpisodeByID)
		}

		// 文件上传路由
		upload := api.Group("/upload")
		upload.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			upload.POST("", fileHandler.UploadFile)
			upload.DELETE("", fileHandler.DeleteFile)
		}

		// 管理员路由
		admin := api.Group("/admin")
		admin.Use(middleware.AdminAuthMiddleware(r.jwtManager))
		{
			// 短剧管理
			adminDramas := admin.Group("/dramas")
			{
				adminDramas.GET("", adminHandler.GetDramaList)
				adminDramas.POST("", adminHandler.CreateDrama)
				adminDramas.PUT("/:id", adminHandler.UpdateDrama)
				adminDramas.DELETE("/:id", adminHandler.DeleteDrama)
				adminDramas.GET("/:drama_id/episodes", adminHandler.GetEpisodeList)
			}

			// 剧集管理
			adminEpisodes := admin.Group("/episodes")
			{
				adminEpisodes.GET("", adminHandler.GetAllEpisodeList)
				adminEpisodes.POST("", adminHandler.CreateEpisode)
				adminEpisodes.PUT("/:id", adminHandler.UpdateEpisode)
				adminEpisodes.DELETE("/:id", adminHandler.DeleteEpisode)
			}

			// 用户管理
			adminUsers := admin.Group("/users")
			{
				adminUsers.GET("", adminHandler.GetUserList)
				adminUsers.POST("/:id/activate", adminHandler.ActivateUser)
				adminUsers.POST("/:id/deactivate", adminHandler.DeactivateUser)
			}
		}
	}

	// 静态文件服务
	r.engine.Static("/uploads", "./uploads")

	// Vue 前端静态文件服务
	r.engine.Static("/assets", "./web/dist/assets")
	r.engine.StaticFile("/favicon.ico", "./web/dist/favicon.ico")

	// Vue SPA 路由支持
	r.setupSPARoutes()
}

// setupSPARoutes 设置 Vue SPA 路由
func (r *Router) setupSPARoutes() {
	// 管理员 API 路由
	adminAPI := r.engine.Group("/admin/api")
	{
		// 认证路由
		authHandler := handler.NewAuthHandler(r.services.AuthService)
		auth := adminAPI.Group("/auth")
		{
			auth.POST("/login", authHandler.AdminLogin)
			auth.POST("/logout", func(c *gin.Context) {
				c.JSON(200, gin.H{"success": true, "message": "退出成功"})
			})

			// 需要认证的路由
			authProtected := auth.Group("")
			authProtected.Use(middleware.AuthMiddleware(r.jwtManager))
			{
				authProtected.GET("/me", func(c *gin.Context) {
					userID := c.GetUint("user_id")
					username := c.GetString("username")
					role := c.GetString("role")

					c.JSON(200, gin.H{
						"success": true,
						"data": gin.H{
							"id":       userID,
							"username": username,
							"role":     role,
						},
					})
				})
			}
		}

		// 其他管理 API 路由可以在这里添加
		// 例如：统计数据、用户管理、短剧管理等
	}

	// SPA 路由处理 - 所有未匹配的路由都返回 index.html
	r.engine.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 如果是 API 请求，返回 404
		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(404, gin.H{"error": "API endpoint not found"})
			return
		}
		if len(path) >= 11 && path[:11] == "/admin/api" {
			c.JSON(404, gin.H{"error": "Admin API endpoint not found"})
			return
		}

		// 其他路由返回 Vue SPA 的 index.html
		c.File("./web/dist/index.html")
	})
}
