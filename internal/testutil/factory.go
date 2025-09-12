package testutil

import (
	"time"

	"gin-mysql-api/internal/models"
	"gin-mysql-api/pkg/utils"
)

// UserFactory 用户工厂
type UserFactory struct{}

// NewUserFactory 创建用户工厂
func NewUserFactory() *UserFactory {
	return &UserFactory{}
}

// CreateUser 创建测试用户
func (f *UserFactory) CreateUser(overrides ...func(*models.User)) *models.User {
	hashedPassword, _ := utils.HashPassword("password123")
	
	user := &models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  hashedPassword,
		Phone:     "12345678901",
		Avatar:    "",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 应用覆盖函数
	for _, override := range overrides {
		override(user)
	}

	return user
}

// CreateUsers 批量创建测试用户
func (f *UserFactory) CreateUsers(count int) []*models.User {
	users := make([]*models.User, count)
	for i := 0; i < count; i++ {
		users[i] = f.CreateUser(func(u *models.User) {
			u.Username = "testuser" + string(rune(i+'0'))
			u.Email = "test" + string(rune(i+'0')) + "@example.com"
		})
	}
	return users
}

// DramaFactory 短剧工厂
type DramaFactory struct{}

// NewDramaFactory 创建短剧工厂
func NewDramaFactory() *DramaFactory {
	return &DramaFactory{}
}

// CreateDrama 创建测试短剧
func (f *DramaFactory) CreateDrama(overrides ...func(*models.Drama)) *models.Drama {
	drama := &models.Drama{
		Title:       "测试短剧",
		Description: "这是一个测试短剧的描述",
		CoverImage:  "https://example.com/cover.jpg",
		Director:    "测试导演",
		Actors:      "测试演员1, 测试演员2",
		Genre:       "喜剧",
		Status:      "active",
		ViewCount:   0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 应用覆盖函数
	for _, override := range overrides {
		override(drama)
	}

	return drama
}

// CreateDramas 批量创建测试短剧
func (f *DramaFactory) CreateDramas(count int) []*models.Drama {
	dramas := make([]*models.Drama, count)
	for i := 0; i < count; i++ {
		dramas[i] = f.CreateDrama(func(d *models.Drama) {
			d.Title = "测试短剧" + string(rune(i+'0'))
		})
	}
	return dramas
}

// EpisodeFactory 剧集工厂
type EpisodeFactory struct{}

// NewEpisodeFactory 创建剧集工厂
func NewEpisodeFactory() *EpisodeFactory {
	return &EpisodeFactory{}
}

// CreateEpisode 创建测试剧集
func (f *EpisodeFactory) CreateEpisode(dramaID uint, overrides ...func(*models.Episode)) *models.Episode {
	episode := &models.Episode{
		DramaID:     dramaID,
		Title:       "测试剧集",
		EpisodeNum:  1,
		Duration:    30, // 30分钟
		VideoURL:    "https://example.com/video.mp4",
		Thumbnail:   "https://example.com/thumbnail.jpg",
		Status:      "active",
		ViewCount:   0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 应用覆盖函数
	for _, override := range overrides {
		override(episode)
	}

	return episode
}

// CreateEpisodes 批量创建测试剧集
func (f *EpisodeFactory) CreateEpisodes(dramaID uint, count int) []*models.Episode {
	episodes := make([]*models.Episode, count)
	for i := 0; i < count; i++ {
		episodes[i] = f.CreateEpisode(dramaID, func(e *models.Episode) {
			e.Title = "测试剧集" + string(rune(i+'0'))
			e.EpisodeNum = i + 1
		})
	}
	return episodes
}

// AdminFactory 管理员工厂
type AdminFactory struct{}

// NewAdminFactory 创建管理员工厂
func NewAdminFactory() *AdminFactory {
	return &AdminFactory{}
}

// CreateAdmin 创建测试管理员
func (f *AdminFactory) CreateAdmin(overrides ...func(*models.Admin)) *models.Admin {
	hashedPassword, _ := utils.HashPassword("admin123")
	
	admin := &models.Admin{
		Username:  "testadmin",
		Email:     "admin@example.com",
		Password:  hashedPassword,
		Role:      "admin",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 应用覆盖函数
	for _, override := range overrides {
		override(admin)
	}

	return admin
}

// Factory 测试工厂集合
type Factory struct {
	User    *UserFactory
	Drama   *DramaFactory
	Episode *EpisodeFactory
	Admin   *AdminFactory
}

// NewFactory 创建测试工厂集合
func NewFactory() *Factory {
	return &Factory{
		User:    NewUserFactory(),
		Drama:   NewDramaFactory(),
		Episode: NewEpisodeFactory(),
		Admin:   NewAdminFactory(),
	}
}