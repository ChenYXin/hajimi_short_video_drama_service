package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP 请求总数
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// HTTP 请求持续时间
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	// HTTP 请求大小
	httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Size of HTTP requests in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"method", "path"},
	)

	// HTTP 响应大小
	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"method", "path", "status"},
	)

	// 当前活跃连接数
	httpActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_connections",
			Help: "Number of active HTTP connections",
		},
	)

	// 业务指标
	userRegistrations = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "user_registrations_total",
			Help: "Total number of user registrations",
		},
	)

	userLogins = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "user_logins_total",
			Help: "Total number of user logins",
		},
	)

	dramaViews = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "drama_views_total",
			Help: "Total number of drama views",
		},
		[]string{"drama_id"},
	)

	episodeViews = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "episode_views_total",
			Help: "Total number of episode views",
		},
		[]string{"episode_id", "drama_id"},
	)

	fileUploads = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "file_uploads_total",
			Help: "Total number of file uploads",
		},
		[]string{"type", "status"},
	)

	// 数据库连接池指标
	dbConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "database_connections",
			Help: "Number of database connections",
		},
		[]string{"state"}, // open, in_use, idle
	)

	// Redis 连接指标
	redisConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redis_connections",
			Help: "Number of Redis connections",
		},
		[]string{"state"}, // active, idle
	)

	// 缓存命中率
	cacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_type"},
	)

	cacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_type"},
	)
)

// PrometheusMetrics Prometheus 指标中间件
func PrometheusMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// 增加活跃连接数
		httpActiveConnections.Inc()
		defer httpActiveConnections.Dec()

		// 记录请求大小
		if c.Request.ContentLength > 0 {
			httpRequestSize.WithLabelValues(c.Request.Method, path).Observe(float64(c.Request.ContentLength))
		}

		// 处理请求
		c.Next()

		// 计算请求持续时间
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		// 记录指标
		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path, status).Observe(duration)
		httpResponseSize.WithLabelValues(c.Request.Method, path, status).Observe(float64(c.Writer.Size()))
	}
}

// BusinessMetrics 业务指标记录器
type BusinessMetrics struct{}

// NewBusinessMetrics 创建业务指标记录器
func NewBusinessMetrics() *BusinessMetrics {
	return &BusinessMetrics{}
}

// RecordUserRegistration 记录用户注册
func (m *BusinessMetrics) RecordUserRegistration() {
	userRegistrations.Inc()
}

// RecordUserLogin 记录用户登录
func (m *BusinessMetrics) RecordUserLogin() {
	userLogins.Inc()
}

// RecordDramaView 记录短剧观看
func (m *BusinessMetrics) RecordDramaView(dramaID string) {
	dramaViews.WithLabelValues(dramaID).Inc()
}

// RecordEpisodeView 记录剧集观看
func (m *BusinessMetrics) RecordEpisodeView(episodeID, dramaID string) {
	episodeViews.WithLabelValues(episodeID, dramaID).Inc()
}

// RecordFileUpload 记录文件上传
func (m *BusinessMetrics) RecordFileUpload(fileType, status string) {
	fileUploads.WithLabelValues(fileType, status).Inc()
}

// RecordCacheHit 记录缓存命中
func (m *BusinessMetrics) RecordCacheHit(cacheType string) {
	cacheHits.WithLabelValues(cacheType).Inc()
}

// RecordCacheMiss 记录缓存未命中
func (m *BusinessMetrics) RecordCacheMiss(cacheType string) {
	cacheMisses.WithLabelValues(cacheType).Inc()
}

// UpdateDBConnections 更新数据库连接数
func (m *BusinessMetrics) UpdateDBConnections(open, inUse, idle int) {
	dbConnections.WithLabelValues("open").Set(float64(open))
	dbConnections.WithLabelValues("in_use").Set(float64(inUse))
	dbConnections.WithLabelValues("idle").Set(float64(idle))
}

// UpdateRedisConnections 更新 Redis 连接数
func (m *BusinessMetrics) UpdateRedisConnections(active, idle int) {
	redisConnections.WithLabelValues("active").Set(float64(active))
	redisConnections.WithLabelValues("idle").Set(float64(idle))
}

// MetricsConfig 指标配置
type MetricsConfig struct {
	// SkipPaths 跳过指标记录的路径
	SkipPaths []string
	// NormalizePath 是否标准化路径（将参数替换为占位符）
	NormalizePath bool
}

// DefaultMetricsConfig 默认指标配置
func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		SkipPaths: []string{
			"/metrics",
			"/health",
			"/ready",
			"/live",
		},
		NormalizePath: true,
	}
}

// CustomPrometheusMetrics 自定义 Prometheus 指标中间件
func CustomPrometheusMetrics(config ...MetricsConfig) gin.HandlerFunc {
	conf := DefaultMetricsConfig()
	if len(config) > 0 {
		conf = config[0]
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 检查是否跳过此路径
		for _, skipPath := range conf.SkipPaths {
			if path == skipPath {
				c.Next()
				return
			}
		}

		start := time.Now()
		
		// 标准化路径
		if conf.NormalizePath {
			if fullPath := c.FullPath(); fullPath != "" {
				path = fullPath
			}
		}

		// 增加活跃连接数
		httpActiveConnections.Inc()
		defer httpActiveConnections.Dec()

		// 记录请求大小
		if c.Request.ContentLength > 0 {
			httpRequestSize.WithLabelValues(c.Request.Method, path).Observe(float64(c.Request.ContentLength))
		}

		// 处理请求
		c.Next()

		// 计算请求持续时间
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		// 记录指标
		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path, status).Observe(duration)
		httpResponseSize.WithLabelValues(c.Request.Method, path, status).Observe(float64(c.Writer.Size()))
	}
}