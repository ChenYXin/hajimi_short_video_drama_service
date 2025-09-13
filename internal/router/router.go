package router

import (
	"html/template"

	"gin-mysql-api/internal/handler"
	"gin-mysql-api/internal/middleware"
	"gin-mysql-api/internal/service"
	templateFuncs "gin-mysql-api/internal/template"
	"gin-mysql-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	
	// Prometheus 指标中间件
	r.engine.Use(middleware.PrometheusMetrics())
	
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
	
	// Prometheus 指标端点
	r.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

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
	r.engine.Static("/static", "./web/static")
	
	// 设置 HTML 模板
	tmpl := template.New("").Funcs(templateFuncs.GetFuncMap())
	
	// 手动加载所有模板文件
	tmpl = template.Must(tmpl.ParseGlob("web/templates/auth/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("web/templates/admin/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("web/templates/layout/*.html"))
	
	r.engine.SetHTMLTemplate(tmpl)
	
	// Web 管理界面路由
	r.setupWebRoutes()
}

// setupWebRoutes 设置 Web 管理界面路由
func (r *Router) setupWebRoutes() {
	// 创建 Web 处理器
	webHandler := handler.NewWebHandler(r.services.AdminService, r.services.UserService, r.services.DramaService)

	// 管理员登录页面（无需认证）
	r.engine.GET("/admin/login", webHandler.LoginPage)
	r.engine.POST("/admin/login", webHandler.Login)
	r.engine.GET("/admin/logout", webHandler.Logout)

	// 管理员界面（需要认证）
	admin := r.engine.Group("/admin")
	admin.Use(r.webAuthMiddleware())
	{
		admin.GET("/", func(c *gin.Context) {
			c.Redirect(302, "/admin/dashboard")
		})
		admin.GET("/dashboard", webHandler.Dashboard)
		admin.GET("/dramas", webHandler.DramasPage)
		admin.GET("/episodes", webHandler.EpisodesPage)
		admin.GET("/users", webHandler.UsersPage)
	}
}

// webAuthMiddleware Web 认证中间件
func (r *Router) webAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 cookie 获取 token
		token, err := c.Cookie("admin_token")
		if err != nil || token == "" {
			c.Redirect(302, "/admin/login")
			c.Abort()
			return
		}

		// 验证 token
		claims, err := r.jwtManager.VerifyToken(token)
		if err != nil || claims.Role != "admin" {
			c.Redirect(302, "/admin/login")
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}