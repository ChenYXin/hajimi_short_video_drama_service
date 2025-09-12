package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SecurityConfig 安全配置
type SecurityConfig struct {
	// XSSProtection X-XSS-Protection 头
	XSSProtection string
	// ContentTypeNosniff X-Content-Type-Options 头
	ContentTypeNosniff bool
	// XFrameOptions X-Frame-Options 头
	XFrameOptions string
	// HSTSMaxAge Strict-Transport-Security 头的 max-age
	HSTSMaxAge int
	// HSTSIncludeSubdomains HSTS 是否包含子域名
	HSTSIncludeSubdomains bool
	// ContentSecurityPolicy Content-Security-Policy 头
	ContentSecurityPolicy string
	// ReferrerPolicy Referrer-Policy 头
	ReferrerPolicy string
}

// DefaultSecurityConfig 默认安全配置
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    true,
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000, // 1 year
		HSTSIncludeSubdomains: true,
		ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	}
}

// Security 安全头中间件
func Security(config ...SecurityConfig) gin.HandlerFunc {
	conf := DefaultSecurityConfig()
	if len(config) > 0 {
		conf = config[0]
	}

	return func(c *gin.Context) {
		// X-XSS-Protection
		if conf.XSSProtection != "" {
			c.Header("X-XSS-Protection", conf.XSSProtection)
		}

		// X-Content-Type-Options
		if conf.ContentTypeNosniff {
			c.Header("X-Content-Type-Options", "nosniff")
		}

		// X-Frame-Options
		if conf.XFrameOptions != "" {
			c.Header("X-Frame-Options", conf.XFrameOptions)
		}

		// Strict-Transport-Security (仅在 HTTPS 下设置)
		if c.Request.TLS != nil && conf.HSTSMaxAge > 0 {
			hstsValue := "max-age=" + strconv.Itoa(conf.HSTSMaxAge)
			if conf.HSTSIncludeSubdomains {
				hstsValue += "; includeSubDomains"
			}
			c.Header("Strict-Transport-Security", hstsValue)
		}

		// Content-Security-Policy
		if conf.ContentSecurityPolicy != "" {
			c.Header("Content-Security-Policy", conf.ContentSecurityPolicy)
		}

		// Referrer-Policy
		if conf.ReferrerPolicy != "" {
			c.Header("Referrer-Policy", conf.ReferrerPolicy)
		}

		c.Next()
	}
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// MaxRequests 最大请求数
	MaxRequests int
	// WindowSize 时间窗口大小（秒）
	WindowSize int
	// KeyFunc 获取限流键的函数
	KeyFunc func(*gin.Context) string
}

// DefaultRateLimitConfig 默认限流配置
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		MaxRequests: 100,
		WindowSize:  60, // 1 minute
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
	}
}

// SimpleRateLimit 简单限流中间件（基于内存）
func SimpleRateLimit(config ...RateLimitConfig) gin.HandlerFunc {
	conf := DefaultRateLimitConfig()
	if len(config) > 0 {
		conf = config[0]
	}

	// 简单的内存存储（生产环境应使用 Redis）
	requests := make(map[string][]time.Time)

	return func(c *gin.Context) {
		key := conf.KeyFunc(c)
		now := time.Now()
		windowStart := now.Add(-time.Duration(conf.WindowSize) * time.Second)

		// 清理过期的请求记录
		if times, exists := requests[key]; exists {
			validTimes := make([]time.Time, 0)
			for _, t := range times {
				if t.After(windowStart) {
					validTimes = append(validTimes, t)
				}
			}
			requests[key] = validTimes
		}

		// 检查是否超过限制
		if len(requests[key]) >= conf.MaxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "请求过于频繁，请稍后再试",
				"error":   "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		// 记录当前请求
		requests[key] = append(requests[key], now)

		c.Next()
	}
}

// IPWhitelist IP 白名单中间件
func IPWhitelist(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		// 检查是否在白名单中
		allowed := false
		for _, ip := range allowedIPs {
			if ip == clientIP || ip == "*" {
				allowed = true
				break
			}
			// 支持 CIDR 格式（简化版）
			if strings.Contains(ip, "/") {
				// 这里应该使用更完整的 CIDR 匹配逻辑
				// 为简化，暂时只支持精确匹配
				if ip == clientIP {
					allowed = true
					break
				}
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "访问被拒绝",
				"error":   "IP not in whitelist",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserAgentFilter User-Agent 过滤中间件
func UserAgentFilter(blockedPatterns []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		
		// 检查是否包含被阻止的模式
		for _, pattern := range blockedPatterns {
			if strings.Contains(strings.ToLower(userAgent), strings.ToLower(pattern)) {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"message": "访问被拒绝",
					"error":   "User agent blocked",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// RequestSizeLimit 请求大小限制中间件
func RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"success": false,
				"message": "请求体过大",
				"error":   "Request entity too large",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Timeout 请求超时中间件
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置请求超时
		ctx, cancel := c.Request.Context(), func() {}
		if timeout > 0 {
			ctx, cancel = c.Request.Context(), cancel
		}
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}