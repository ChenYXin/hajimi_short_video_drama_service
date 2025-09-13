package service

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"gin-mysql-api/internal/middleware"
	"gin-mysql-api/pkg/config"
)

// MonitoringService 监控服务
type MonitoringService struct {
	db            *gorm.DB
	redisClient   *redis.Client
	metrics       *middleware.BusinessMetrics
	config        *config.Config
	stopChan      chan struct{}
	collectTicker *time.Ticker
}

// NewMonitoringService 创建监控服务
func NewMonitoringService(db *gorm.DB, redisClient *redis.Client, cfg *config.Config) *MonitoringService {
	return &MonitoringService{
		db:          db,
		redisClient: redisClient,
		metrics:     middleware.NewBusinessMetrics(),
		config:      cfg,
		stopChan:    make(chan struct{}),
	}
}

// Start 启动监控服务
func (s *MonitoringService) Start() {
	if !s.config.Prometheus.Enabled {
		log.Println("Prometheus监控未启用")
		return
	}

	log.Println("启动监控服务...")
	
	// 每30秒收集一次系统指标
	s.collectTicker = time.NewTicker(30 * time.Second)
	
	go s.collectMetrics()
}

// Stop 停止监控服务
func (s *MonitoringService) Stop() {
	log.Println("停止监控服务...")
	
	if s.collectTicker != nil {
		s.collectTicker.Stop()
	}
	
	close(s.stopChan)
}

// collectMetrics 收集系统指标
func (s *MonitoringService) collectMetrics() {
	for {
		select {
		case <-s.collectTicker.C:
			s.collectDatabaseMetrics()
			s.collectRedisMetrics()
		case <-s.stopChan:
			return
		}
	}
}

// collectDatabaseMetrics 收集数据库指标
func (s *MonitoringService) collectDatabaseMetrics() {
	if s.db == nil {
		return
	}

	sqlDB, err := s.db.DB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		return
	}

	stats := sqlDB.Stats()
	
	// 更新数据库连接池指标
	s.metrics.UpdateDBConnections(
		stats.OpenConnections,
		stats.InUse,
		stats.Idle,
	)

	// 记录数据库性能指标
	if stats.MaxOpenConnections > 0 {
		// 可以添加更多数据库性能指标
		log.Printf("数据库连接池状态 - 打开: %d, 使用中: %d, 空闲: %d, 最大: %d",
			stats.OpenConnections, stats.InUse, stats.Idle, stats.MaxOpenConnections)
	}
}

// collectRedisMetrics 收集Redis指标
func (s *MonitoringService) collectRedisMetrics() {
	if s.redisClient == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取Redis连接池状态
	poolStats := s.redisClient.PoolStats()
	
	// 更新Redis连接指标
	s.metrics.UpdateRedisConnections(
		int(poolStats.TotalConns-poolStats.IdleConns),
		int(poolStats.IdleConns),
	)

	// 测试Redis连接
	_, err := s.redisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis连接检查失败: %v", err)
	} else {
		log.Printf("Redis连接池状态 - 总连接: %d, 空闲: %d, 命中: %d, 未命中: %d, 超时: %d",
			poolStats.TotalConns, poolStats.IdleConns, poolStats.Hits, poolStats.Misses, poolStats.Timeouts)
	}
}

// GetMetrics 获取业务指标记录器
func (s *MonitoringService) GetMetrics() *middleware.BusinessMetrics {
	return s.metrics
}

// HealthCheck 健康检查
func (s *MonitoringService) HealthCheck() map[string]interface{} {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"services":  make(map[string]interface{}),
	}

	// 检查数据库连接
	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err != nil {
			health["services"].(map[string]interface{})["database"] = map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			}
		} else {
			err = sqlDB.Ping()
			if err != nil {
				health["services"].(map[string]interface{})["database"] = map[string]interface{}{
					"status": "error",
					"error":  err.Error(),
				}
			} else {
				stats := sqlDB.Stats()
				health["services"].(map[string]interface{})["database"] = map[string]interface{}{
					"status":      "healthy",
					"connections": stats.OpenConnections,
					"in_use":      stats.InUse,
					"idle":        stats.Idle,
				}
			}
		}
	}

	// 检查Redis连接
	if s.redisClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := s.redisClient.Ping(ctx).Result()
		if err != nil {
			health["services"].(map[string]interface{})["redis"] = map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			}
		} else {
			poolStats := s.redisClient.PoolStats()
			health["services"].(map[string]interface{})["redis"] = map[string]interface{}{
				"status":      "healthy",
				"connections": poolStats.TotalConns,
				"idle":        poolStats.IdleConns,
				"hits":        poolStats.Hits,
				"misses":      poolStats.Misses,
			}
		}
	}

	// 如果任何服务不健康，整体状态为不健康
	for _, service := range health["services"].(map[string]interface{}) {
		if serviceMap, ok := service.(map[string]interface{}); ok {
			if status, exists := serviceMap["status"]; exists && status != "healthy" {
				health["status"] = "unhealthy"
				break
			}
		}
	}

	return health
}

// GetSystemStats 获取系统统计信息
func (s *MonitoringService) GetSystemStats() map[string]interface{} {
	stats := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"database":  make(map[string]interface{}),
		"redis":     make(map[string]interface{}),
	}

	// 数据库统计
	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err == nil {
			dbStats := sqlDB.Stats()
			stats["database"] = map[string]interface{}{
				"open_connections": dbStats.OpenConnections,
				"in_use":           dbStats.InUse,
				"idle":             dbStats.Idle,
				"wait_count":       dbStats.WaitCount,
				"wait_duration":    dbStats.WaitDuration.Milliseconds(),
				"max_idle_closed":  dbStats.MaxIdleClosed,
				"max_lifetime_closed": dbStats.MaxLifetimeClosed,
			}
		}
	}

	// Redis统计
	if s.redisClient != nil {
		poolStats := s.redisClient.PoolStats()
		stats["redis"] = map[string]interface{}{
			"total_conns": poolStats.TotalConns,
			"idle_conns":  poolStats.IdleConns,
			"stale_conns": poolStats.StaleConns,
			"hits":        poolStats.Hits,
			"misses":      poolStats.Misses,
			"timeouts":    poolStats.Timeouts,
		}
	}

	return stats
}