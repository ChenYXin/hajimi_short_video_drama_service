package testutil

import (
	"time"

	"gin-mysql-api/pkg/config"
)

// GetTestConfig 获取测试配置
func GetTestConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 1800,
			Mode: "test",
		},
		Database: config.DatabaseConfig{
			Host:            "localhost",
			Port:            3306,
			Username:        "test",
			Password:        "test",
			DBName:          "hajimi_test",
			Charset:         "utf8mb4",
			ParseTime:       true,
			Loc:             "Local",
			MaxIdleConns:    5,
			MaxOpenConns:    10,
			ConnMaxLifetime: 300 * time.Second,
		},
		Redis: config.RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       1, // 使用测试数据库
			PoolSize: 5,
		},
		JWT: config.JWTConfig{
			Secret:     "test-secret-key",
			Expiration: 1 * time.Hour,
		},
		Upload: config.UploadConfig{
			MaxSize:      10, // 10MB
			AllowedTypes: []string{"jpg", "jpeg", "png", "gif", "mp4", "avi", "mov"},
			UploadPath:   "./test_uploads",
		},
		Logging: config.LoggingConfig{
			Level:    "debug",
			Format:   "json",
			Output:   "stdout",
			Filename: "",
		},
		Prometheus: config.PrometheusConfig{
			Enabled: false,
			Path:    "/metrics",
			Port:    9090,
		},
	}
}