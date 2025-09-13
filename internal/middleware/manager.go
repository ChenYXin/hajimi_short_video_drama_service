package middleware

import (
	"time"

	"gin-mysql-api/pkg/config"

	"github.com/gin-gonic/gin"
)

// Manager 中间件管理器
type Manager struct {
	config          *config.Config
	businessMetrics *BusinessMetrics
}

// NewManager 创建中间件管理器
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		config:          cfg,
		businessMetrics: NewBusinessMetrics(),
	}
}

// GetBusinessMetrics 获取业务指标记录器
func (m *Manager) GetBusinessMetrics() *BusinessMetrics {
	return m.businessMetrics
}

// SetupMiddlewares 设置所有中间件
func (m *Manager) SetupMiddlewares(engine *gin.Engine) {
	// 错误处理中间件（最先设置）
	engine.Use(ErrorHandler())
	
	// 请求 ID 中间件
	engine.Use(RequestIDMiddleware())
	
	// 日志中间件
	if m.config.Logging.Level == "debug" {
		// 开发环境使用详细日志
		engine.Use(DetailedLogger(LoggerConfig{
			SkipPaths:       []string{"/health", "/ready", "/live", "/metrics"},
			LogRequestBody:  true,
			LogResponseBody: false,
			MaxBodySize:     1024 * 10, // 10KB
		}))
	} else {
		// 生产环境使用简单日志
		engine.Use(Logger(LoggerConfig{
			SkipPaths: []string{"/health", "/ready", "/live", "/metrics"},
		}))
	}
	
	// 安全中间件
	securityConfig := SecurityConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    true,
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000, // 1 year
		HSTSIncludeSubdomains: true,
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://code.jquery.com https://cdnjs.cloudflare.com; style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://cdnjs.cloudflare.com; font-src 'self' https://cdnjs.cloudflare.com; img-src 'self' data: https:;",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	}
	engine.Use(Security(securityConfig))
	
	// CORS 中间件
	corsConfig := CORSConfig{
		AllowOrigins: []string{"*"}, // 生产环境应该配置具体的域名
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "Authorization",
			"X-Request-ID", "X-Requested-With",
		},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           86400, // 24 hours
	}
	engine.Use(CORS(corsConfig))
	
	// Prometheus 指标中间件
	metricsConfig := MetricsConfig{
		SkipPaths:     []string{"/metrics", "/health", "/ready", "/live"},
		NormalizePath: true,
	}
	engine.Use(CustomPrometheusMetrics(metricsConfig))
	
	// 请求大小限制中间件
	maxRequestSize := int64(10 * 1024 * 1024) // 10MB
	if m.config.Upload.MaxSize > 0 {
		maxRequestSize = int64(m.config.Upload.MaxSize) * 1024 * 1024
	}
	engine.Use(RequestSizeLimit(maxRequestSize))
	
	// 限流中间件
	rateLimitConfig := RateLimitConfig{
		MaxRequests: 1000, // 每分钟最多 1000 个请求
		WindowSize:  60,   // 1 minute
		KeyFunc: func(c *gin.Context) string {
			// 对认证用户使用用户 ID，对未认证用户使用 IP
			if userID, exists := c.Get("user_id"); exists {
				return "user:" + string(rune(userID.(uint)))
			}
			return "ip:" + c.ClientIP()
		},
	}
	engine.Use(SimpleRateLimit(rateLimitConfig))
	
	// 请求超时中间件
	engine.Use(Timeout(30 * time.Second))
	
	// 404 和 405 处理
	engine.NoRoute(NotFoundHandler())
	engine.NoMethod(MethodNotAllowedHandler())
}

// SetupProductionMiddlewares 设置生产环境中间件
func (m *Manager) SetupProductionMiddlewares(engine *gin.Engine) {
	// 基础中间件
	m.SetupMiddlewares(engine)
	
	// 生产环境特有的中间件
	
	// IP 白名单（如果配置了）
	if len(m.config.Server.AllowedIPs) > 0 {
		engine.Use(IPWhitelist(m.config.Server.AllowedIPs))
	}
	
	// User-Agent 过滤
	blockedUserAgents := []string{
		"bot", "crawler", "spider", "scraper",
	}
	engine.Use(UserAgentFilter(blockedUserAgents))
	
	// 更严格的限流
	strictRateLimitConfig := RateLimitConfig{
		MaxRequests: 100, // 每分钟最多 100 个请求
		WindowSize:  60,  // 1 minute
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
	}
	engine.Use(SimpleRateLimit(strictRateLimitConfig))
}

// SetupDevelopmentMiddlewares 设置开发环境中间件
func (m *Manager) SetupDevelopmentMiddlewares(engine *gin.Engine) {
	// 基础中间件
	m.SetupMiddlewares(engine)
	
	// 开发环境特有的中间件
	
	// 更宽松的 CORS 配置
	devCorsConfig := CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:1800", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           86400,
	}
	engine.Use(CORS(devCorsConfig))
}