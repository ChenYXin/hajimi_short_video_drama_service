package repository

import (
	"gorm.io/gorm"
)

// Repository 仓库管理器，包含所有仓库接口
type Repository struct {
	User    UserRepository
	Drama   DramaRepository
	Episode EpisodeRepository
	Admin   AdminRepository
}

// NewRepository 创建仓库管理器实例
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:    NewUserRepository(db),
		Drama:   NewDramaRepository(db),
		Episode: NewEpisodeRepository(db),
		Admin:   NewAdminRepository(db),
	}
}