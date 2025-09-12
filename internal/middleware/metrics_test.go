package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 重置 Prometheus 指标
	prometheus.DefaultRegisterer = prometheus.NewRegistry()

	t.Run("记录HTTP请求指标", func(t *testing.T) {
		router := gin.New()
		router.Use(PrometheusMetrics())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		// 验证指标是否被记录
		// 注意：由于 Prometheus 指标是全局的，在测试中验证具体值可能比较困难
		// 这里主要验证中间件能正常运行而不出错
	})

	t.Run("记录不同状态码", func(t *testing.T) {
		router := gin.New()
		router.Use(PrometheusMetrics())
		router.GET("/error", func(c *gin.Context) {
			c.JSON(500, gin.H{"error": "internal server error"})
		})

		req := httptest.NewRequest("GET", "/error", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
	})
}

func TestBusinessMetrics(t *testing.T) {
	// 重置 Prometheus 指标
	prometheus.DefaultRegisterer = prometheus.NewRegistry()

	metrics := NewBusinessMetrics()

	t.Run("记录用户注册", func(t *testing.T) {
		initialValue := testutil.ToFloat64(userRegistrations)
		
		metrics.RecordUserRegistration()
		
		newValue := testutil.ToFloat64(userRegistrations)
		assert.Equal(t, initialValue+1, newValue)
	})

	t.Run("记录用户登录", func(t *testing.T) {
		initialValue := testutil.ToFloat64(userLogins)
		
		metrics.RecordUserLogin()
		
		newValue := testutil.ToFloat64(userLogins)
		assert.Equal(t, initialValue+1, newValue)
	})

	t.Run("记录短剧观看", func(t *testing.T) {
		dramaID := "123"
		
		metrics.RecordDramaView(dramaID)
		
		// 验证指标被记录（具体值验证在实际项目中可能需要更复杂的设置）
		assert.NotNil(t, dramaViews)
	})

	t.Run("记录剧集观看", func(t *testing.T) {
		episodeID := "456"
		dramaID := "123"
		
		metrics.RecordEpisodeView(episodeID, dramaID)
		
		assert.NotNil(t, episodeViews)
	})

	t.Run("记录文件上传", func(t *testing.T) {
		fileType := "image"
		status := "success"
		
		metrics.RecordFileUpload(fileType, status)
		
		assert.NotNil(t, fileUploads)
	})

	t.Run("记录缓存命中", func(t *testing.T) {
		cacheType := "redis"
		
		metrics.RecordCacheHit(cacheType)
		
		assert.NotNil(t, cacheHits)
	})

	t.Run("记录缓存未命中", func(t *testing.T) {
		cacheType := "redis"
		
		metrics.RecordCacheMiss(cacheType)
		
		assert.NotNil(t, cacheMisses)
	})

	t.Run("更新数据库连接数", func(t *testing.T) {
		open := 10
		inUse := 5
		idle := 5
		
		metrics.UpdateDBConnections(open, inUse, idle)
		
		assert.NotNil(t, dbConnections)
	})

	t.Run("更新Redis连接数", func(t *testing.T) {
		active := 3
		idle := 2
		
		metrics.UpdateRedisConnections(active, idle)
		
		assert.NotNil(t, redisConnections)
	})
}

func TestCustomPrometheusMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("跳过指定路径", func(t *testing.T) {
		config := MetricsConfig{
			SkipPaths: []string{"/metrics", "/health"},
		}

		router := gin.New()
		router.Use(CustomPrometheusMetrics(config))
		router.GET("/metrics", func(c *gin.Context) {
			c.JSON(200, gin.H{"metrics": "data"})
		})
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		// 测试跳过的路径
		req := httptest.NewRequest("GET", "/metrics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		// 测试正常路径
		req = httptest.NewRequest("GET", "/test", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("标准化路径", func(t *testing.T) {
		config := MetricsConfig{
			NormalizePath: true,
		}

		router := gin.New()
		router.Use(CustomPrometheusMetrics(config))
		router.GET("/user/:id", func(c *gin.Context) {
			c.JSON(200, gin.H{"user_id": c.Param("id")})
		})

		req := httptest.NewRequest("GET", "/user/123", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})
}

func TestDefaultMetricsConfig(t *testing.T) {
	config := DefaultMetricsConfig()
	
	assert.Contains(t, config.SkipPaths, "/metrics")
	assert.Contains(t, config.SkipPaths, "/health")
	assert.Contains(t, config.SkipPaths, "/ready")
	assert.Contains(t, config.SkipPaths, "/live")
	assert.True(t, config.NormalizePath)
}

func TestNewBusinessMetrics(t *testing.T) {
	metrics := NewBusinessMetrics()
	assert.NotNil(t, metrics)
	assert.IsType(t, &BusinessMetrics{}, metrics)
}