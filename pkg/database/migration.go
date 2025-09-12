package database

import (
	"fmt"

	"gin-mysql-api/internal/models"
	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	// 定义需要迁移的模型
	models := []interface{}{
		&models.User{},
		&models.Admin{},
		&models.Drama{},
		&models.Episode{},
	}

	// 执行自动迁移
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	return nil
}

// CreateIndexes 创建数据库索引
func CreateIndexes(db *gorm.DB) error {
	// 用户表索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)").Error; err != nil {
		return fmt.Errorf("failed to create username index: %w", err)
	}
	
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error; err != nil {
		return fmt.Errorf("failed to create email index: %w", err)
	}

	// 短剧表索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_dramas_title ON dramas(title)").Error; err != nil {
		return fmt.Errorf("failed to create drama title index: %w", err)
	}
	
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_dramas_category ON dramas(category)").Error; err != nil {
		return fmt.Errorf("failed to create drama category index: %w", err)
	}
	
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_dramas_status ON dramas(status)").Error; err != nil {
		return fmt.Errorf("failed to create drama status index: %w", err)
	}

	// 剧集表索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_episodes_drama_id ON episodes(drama_id)").Error; err != nil {
		return fmt.Errorf("failed to create episode drama_id index: %w", err)
	}
	
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_episodes_episode_number ON episodes(episode_number)").Error; err != nil {
		return fmt.Errorf("failed to create episode number index: %w", err)
	}

	return nil
}

// SeedData 插入种子数据
func SeedData(db *gorm.DB) error {
	// 检查是否已经有管理员用户
	var adminCount int64
	if err := db.Model(&models.Admin{}).Count(&adminCount).Error; err != nil {
		return fmt.Errorf("failed to count admins: %w", err)
	}

	// 如果没有管理员，创建默认管理员
	if adminCount == 0 {
		defaultAdmin := &models.Admin{
			Username: "admin",
			Email:    "admin@example.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			IsActive: true,
		}

		if err := db.Create(defaultAdmin).Error; err != nil {
			return fmt.Errorf("failed to create default admin: %w", err)
		}
	}

	// 检查是否已经有示例短剧
	var dramaCount int64
	if err := db.Model(&models.Drama{}).Count(&dramaCount).Error; err != nil {
		return fmt.Errorf("failed to count dramas: %w", err)
	}

	// 如果没有短剧，创建示例短剧
	if dramaCount == 0 {
		sampleDrama := &models.Drama{
			Title:       "示例短剧",
			Description: "这是一个示例短剧，用于测试系统功能",
			Genre:       "爱情",
			CoverImage:  "/static/images/sample-cover.jpg",
			Status:      "active",
			ViewCount:   0,
		}

		if err := db.Create(sampleDrama).Error; err != nil {
			return fmt.Errorf("failed to create sample drama: %w", err)
		}

		// 为示例短剧创建示例剧集
		sampleEpisode := &models.Episode{
			DramaID:    sampleDrama.ID,
			Title:      "第一集",
			EpisodeNum: 1,
			VideoURL:   "/static/videos/sample-episode.mp4",
			Duration:   300, // 5分钟
			ViewCount:  0,
		}

		if err := db.Create(sampleEpisode).Error; err != nil {
			return fmt.Errorf("failed to create sample episode: %w", err)
		}
	}

	return nil
}

// InitDatabase 初始化数据库（迁移、索引、种子数据）
func InitDatabase(db *gorm.DB) error {
	// 执行自动迁移
	if err := AutoMigrate(db); err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}

	// 创建索引
	if err := CreateIndexes(db); err != nil {
		return fmt.Errorf("create indexes failed: %w", err)
	}

	// 插入种子数据
	if err := SeedData(db); err != nil {
		return fmt.Errorf("seed data failed: %w", err)
	}

	return nil
}