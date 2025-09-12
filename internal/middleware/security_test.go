package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSecurity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("默认安全头设置", func(t *testing.T) {
		router := gin.New()
		router.Use(Security())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
		assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
		assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
		assert.Equal(t, "default-src 'self'", w.Header().Get("Content-Security-Policy"))
		assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
		
		// HSTS 头只在 HTTPS 下设置，HTTP 请求不会有
		assert.Empty(t, w.Header().Get("Strict-Transport-Security"))
	})

	t.Run("自定义安全配置", func(t *testing.T) {
		config := SecurityConfig{
			XSSProtection:         "0",
			ContentTypeNosniff:    false,
			XFrameOptions:         "SAMEORIGIN",
			ContentSecurityPolicy: "default-src 'self' 'unsafe-inline'",
			ReferrerPolicy:        "no-referrer",
		}

		router := gin.New()
		router.Use(Security(config))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "0", w.Header().Get("X-XSS-Protection"))
		assert.Empty(t, w.Header().Get("X-Content-Type-Options"))
		assert.Equal(t, "SAMEORIGIN", w.Header().Get("X-Frame-Options"))
		assert.Equal(t, "default-src 'self' 'unsafe-inline'", w.Header().Get("Content-Security-Policy"))
		assert.Equal(t, "no-referrer", w.Header().Get("Referrer-Policy"))
	})
}

func TestSimpleRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常请求通过", func(t *testing.T) {
		config := RateLimitConfig{
			MaxRequests: 5,
			WindowSize:  60,
			KeyFunc: func(c *gin.Context) string {
				return c.ClientIP()
			},
		}

		router := gin.New()
		router.Use(SimpleRateLimit(config))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		// 发送几个正常请求
		for i := 0; i < 3; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
		}
	})

	t.Run("超过限制被拒绝", func(t *testing.T) {
		config := RateLimitConfig{
			MaxRequests: 2,
			WindowSize:  60,
			KeyFunc: func(c *gin.Context) string {
				return "test-key"
			},
		}

		router := gin.New()
		router.Use(SimpleRateLimit(config))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		// 发送允许的请求数
		for i := 0; i < 2; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
		}

		// 第三个请求应该被拒绝
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 429, w.Code)
		assert.Contains(t, w.Body.String(), "请求过于频繁")
	})
}

func TestIPWhitelist(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("白名单IP允许访问", func(t *testing.T) {
		allowedIPs := []string{"127.0.0.1", "192.168.1.1"}

		router := gin.New()
		router.Use(IPWhitelist(allowedIPs))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("非白名单IP被拒绝", func(t *testing.T) {
		allowedIPs := []string{"192.168.1.1"}

		router := gin.New()
		router.Use(IPWhitelist(allowedIPs))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "访问被拒绝")
	})

	t.Run("通配符允许所有IP", func(t *testing.T) {
		allowedIPs := []string{"*"}

		router := gin.New()
		router.Use(IPWhitelist(allowedIPs))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})
}

func TestUserAgentFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常User-Agent通过", func(t *testing.T) {
		blockedPatterns := []string{"bot", "crawler", "spider"}

		router := gin.New()
		router.Use(UserAgentFilter(blockedPatterns))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("被阻止的User-Agent被拒绝", func(t *testing.T) {
		blockedPatterns := []string{"bot", "crawler", "spider"}

		router := gin.New()
		router.Use(UserAgentFilter(blockedPatterns))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("User-Agent", "Googlebot/2.1")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "访问被拒绝")
	})

	t.Run("大小写不敏感匹配", func(t *testing.T) {
		blockedPatterns := []string{"BOT"}

		router := gin.New()
		router.Use(UserAgentFilter(blockedPatterns))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("User-Agent", "testbot/1.0")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
	})
}

func TestRequestSizeLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常大小请求通过", func(t *testing.T) {
		maxSize := int64(1024) // 1KB

		router := gin.New()
		router.Use(RequestSizeLimit(maxSize))
		router.POST("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		body := strings.NewReader("small request body")
		req := httptest.NewRequest("POST", "/test", body)
		req.ContentLength = int64(len("small request body"))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("超大请求被拒绝", func(t *testing.T) {
		maxSize := int64(10) // 10 bytes

		router := gin.New()
		router.Use(RequestSizeLimit(maxSize))
		router.POST("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		body := strings.NewReader("this is a very long request body that exceeds the limit")
		req := httptest.NewRequest("POST", "/test", body)
		req.ContentLength = int64(len("this is a very long request body that exceeds the limit"))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 413, w.Code)
		assert.Contains(t, w.Body.String(), "请求体过大")
	})
}

func TestTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常请求处理", func(t *testing.T) {
		router := gin.New()
		router.Use(Timeout(time.Second))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	// 注意：测试超时场景比较复杂，因为需要模拟长时间运行的处理器
	// 在实际项目中，可能需要更复杂的测试设置
}

func TestDefaultSecurityConfig(t *testing.T) {
	config := DefaultSecurityConfig()
	
	assert.Equal(t, "1; mode=block", config.XSSProtection)
	assert.True(t, config.ContentTypeNosniff)
	assert.Equal(t, "DENY", config.XFrameOptions)
	assert.Equal(t, 31536000, config.HSTSMaxAge)
	assert.True(t, config.HSTSIncludeSubdomains)
	assert.Equal(t, "default-src 'self'", config.ContentSecurityPolicy)
	assert.Equal(t, "strict-origin-when-cross-origin", config.ReferrerPolicy)
}

func TestDefaultRateLimitConfig(t *testing.T) {
	config := DefaultRateLimitConfig()
	
	assert.Equal(t, 100, config.MaxRequests)
	assert.Equal(t, 60, config.WindowSize)
	assert.NotNil(t, config.KeyFunc)
}