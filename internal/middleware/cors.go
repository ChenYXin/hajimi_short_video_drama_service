package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig 默认 CORS 配置
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"X-Request-ID",
		},
		ExposeHeaders: []string{
			"X-Request-ID",
		},
		AllowCredentials: false,
		MaxAge:           86400, // 24 hours
	}
}

// CORS 跨域资源共享中间件
func CORS(config ...CORSConfig) gin.HandlerFunc {
	conf := DefaultCORSConfig()
	if len(config) > 0 {
		conf = config[0]
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// 检查是否允许该来源
		if len(conf.AllowOrigins) == 1 && conf.AllowOrigins[0] == "*" {
			c.Header("Access-Control-Allow-Origin", "*")
		} else if isOriginAllowed(origin, conf.AllowOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// 设置其他 CORS 头
		c.Header("Access-Control-Allow-Methods", strings.Join(conf.AllowMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(conf.AllowHeaders, ", "))
		
		if len(conf.ExposeHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(conf.ExposeHeaders, ", "))
		}
		
		if conf.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		
		if conf.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", string(rune(conf.MaxAge)))
		}

		// 处理预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// isOriginAllowed 检查来源是否被允许
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
		// 支持通配符匹配
		if strings.Contains(allowed, "*") {
			// 简单的通配符匹配，支持 *.example.com 格式
			if strings.HasPrefix(allowed, "*.") {
				domain := strings.TrimPrefix(allowed, "*.")
				if strings.HasSuffix(origin, "."+domain) || strings.HasSuffix(origin, "://"+domain) {
					return true
				}
			} else if strings.HasSuffix(allowed, "*") {
				prefix := strings.TrimSuffix(allowed, "*")
				if strings.HasPrefix(origin, prefix) {
					return true
				}
			}
		}
	}
	return false
}