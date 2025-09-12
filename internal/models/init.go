package models

import (
	"gorm.io/gorm"
)

// AllModels 返回所有需要迁移的模型
func AllModels() []interface{} {
	return []interface{}{
		&User{},
		&Drama{},
		&Episode{},
		&Admin{},
	}
}

// AutoMigrate 自动迁移所有模型
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(AllModels()...)
}

// CreateIndexes 创建额外的索引
func CreateIndexes(db *gorm.DB) error {
	// 为 dramas 表创建复合索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_dramas_genre_status ON dramas(genre, status)").Error; err != nil {
		return err
	}
	
	// 为 episodes 表创建复合索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_episodes_drama_episode ON episodes(drama_id, episode_num)").Error; err != nil {
		return err
	}
	
	return nil
}

// SeedData 初始化种子数据
func SeedData(db *gorm.DB) error {
	// 创建默认管理员账户
	var adminCount int64
	db.Model(&Admin{}).Count(&adminCount)
	
	if adminCount == 0 {
		defaultAdmin := &Admin{
			Username: "admin",
			Email:    "admin@example.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     "super_admin",
			IsActive: true,
		}
		
		if err := db.Create(defaultAdmin).Error; err != nil {
			return err
		}
	}
	
	return nil
}