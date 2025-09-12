package database

import (
	"fmt"
	"log"

	"gin-mysql-api/pkg/config"
)

// Manager 数据库管理器
type Manager struct {
	config *config.Config
}

// NewManager 创建新的数据库管理器
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		config: cfg,
	}
}

// Initialize 初始化所有数据库连接
func (m *Manager) Initialize() error {
	// 初始化 MySQL 连接
	if err := InitMySQL(&m.config.Database); err != nil {
		return fmt.Errorf("failed to initialize MySQL: %w", err)
	}
	log.Println("MySQL 数据库连接成功")

	// 初始化 Redis 连接
	if err := InitRedis(&m.config.Redis); err != nil {
		return fmt.Errorf("failed to initialize Redis: %w", err)
	}
	log.Println("Redis 连接成功")

	// 初始化数据库表结构和数据
	if err := InitDatabase(DB); err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}
	log.Println("数据库初始化完成")

	return nil
}

// Close 关闭所有数据库连接
func (m *Manager) Close() error {
	var errors []error

	// 关闭 MySQL 连接
	if err := CloseDB(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close MySQL: %w", err))
	}

	// 关闭 Redis 连接
	if err := CloseRedis(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close Redis: %w", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors occurred while closing databases: %v", errors)
	}

	log.Println("所有数据库连接已关闭")
	return nil
}

// HealthCheck 检查数据库连接健康状态
func (m *Manager) HealthCheck() error {
	// 检查 MySQL 连接
	if DB == nil {
		return fmt.Errorf("MySQL connection is nil")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get MySQL connection: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("MySQL ping failed: %w", err)
	}

	// 检查 Redis 连接
	if RedisClient == nil {
		return fmt.Errorf("Redis connection is nil")
	}

	if err := RedisClient.Ping(RedisClient.Context()).Err(); err != nil {
		return fmt.Errorf("Redis ping failed: %w", err)
	}

	return nil
}