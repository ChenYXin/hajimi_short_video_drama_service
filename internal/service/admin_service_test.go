package service

import (
	"errors"
	"testing"
	"time"

	"gin-mysql-api/internal/models"
	"gin-mysql-api/pkg/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAdminRepository 模拟管理员仓库
type MockAdminRepository struct {
	mock.Mock
}

func (m *MockAdminRepository) Create(admin *models.Admin) error {
	args := m.Called(admin)
	return args.Error(0)
}

func (m *MockAdminRepository) GetByID(id uint) (*models.Admin, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Admin), args.Error(1)
}

func (m *MockAdminRepository) GetByEmail(email string) (*models.Admin, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Admin), args.Error(1)
}

func (m *MockAdminRepository) GetByUsername(username string) (*models.Admin, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Admin), args.Error(1)
}

func (m *MockAdminRepository) Update(admin *models.Admin) error {
	args := m.Called(admin)
	return args.Error(0)
}

func (m *MockAdminRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAdminRepository) List(offset, limit int) ([]models.Admin, int64, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]models.Admin), args.Get(1).(int64), args.Error(2)
}

func (m *MockAdminRepository) ExistsByEmail(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockAdminRepository) ExistsByUsername(username string) (bool, error) {
	args := m.Called(username)
	return args.Bool(0), args.Error(1)
}

// MockDramaRepository 模拟短剧仓库
type MockDramaRepository struct {
	mock.Mock
}

func (m *MockDramaRepository) Create(drama *models.Drama) error {
	args := m.Called(drama)
	return args.Error(0)
}

func (m *MockDramaRepository) GetByID(id uint) (*models.Drama, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Drama), args.Error(1)
}

func (m *MockDramaRepository) GetByIDWithEpisodes(id uint) (*models.Drama, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Drama), args.Error(1)
}

func (m *MockDramaRepository) GetList(offset, limit int, genre string) ([]models.Drama, int64, error) {
	args := m.Called(offset, limit, genre)
	return args.Get(0).([]models.Drama), args.Get(1).(int64), args.Error(2)
}

func (m *MockDramaRepository) Update(drama *models.Drama) error {
	args := m.Called(drama)
	return args.Error(0)
}

func (m *MockDramaRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDramaRepository) IncrementViewCount(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDramaRepository) GetByGenre(genre string, offset, limit int) ([]models.Drama, int64, error) {
	args := m.Called(genre, offset, limit)
	return args.Get(0).([]models.Drama), args.Get(1).(int64), args.Error(2)
}

func (m *MockDramaRepository) GetActiveList(offset, limit int) ([]models.Drama, int64, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]models.Drama), args.Get(1).(int64), args.Error(2)
}

// MockEpisodeRepository 模拟剧集仓库
type MockEpisodeRepository struct {
	mock.Mock
}

func (m *MockEpisodeRepository) Create(episode *models.Episode) error {
	args := m.Called(episode)
	return args.Error(0)
}

func (m *MockEpisodeRepository) GetByID(id uint) (*models.Episode, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Episode), args.Error(1)
}

func (m *MockEpisodeRepository) GetByIDWithDrama(id uint) (*models.Episode, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Episode), args.Error(1)
}

func (m *MockEpisodeRepository) GetByDramaID(dramaID uint) ([]models.Episode, error) {
	args := m.Called(dramaID)
	return args.Get(0).([]models.Episode), args.Error(1)
}

func (m *MockEpisodeRepository) GetByDramaIDPaginated(dramaID uint, offset, limit int) ([]models.Episode, int64, error) {
	args := m.Called(dramaID, offset, limit)
	return args.Get(0).([]models.Episode), args.Get(1).(int64), args.Error(2)
}

func (m *MockEpisodeRepository) Update(episode *models.Episode) error {
	args := m.Called(episode)
	return args.Error(0)
}

func (m *MockEpisodeRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEpisodeRepository) IncrementViewCount(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEpisodeRepository) GetMaxEpisodeNum(dramaID uint) (int, error) {
	args := m.Called(dramaID)
	return args.Int(0), args.Error(1)
}

func (m *MockEpisodeRepository) ExistsByDramaIDAndEpisodeNum(dramaID uint, episodeNum int) (bool, error) {
	args := m.Called(dramaID, episodeNum)
	return args.Bool(0), args.Error(1)
}

// MockCacheService 模拟缓存服务
type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) Set(key string, value interface{}, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheService) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheService) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockCacheService) Exists(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheService) SetJSON(key string, value interface{}, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetJSON(key string, dest interface{}) error {
	args := m.Called(key, dest)
	return args.Error(0)
}

func (m *MockCacheService) DeletePattern(pattern string) error {
	args := m.Called(pattern)
	return args.Error(0)
}

func (m *MockCacheService) Increment(key string) (int64, error) {
	args := m.Called(key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCacheService) Expire(key string, expiration time.Duration) error {
	args := m.Called(key, expiration)
	return args.Error(0)
}

func TestAdminService_Login(t *testing.T) {
	mockAdminRepo := new(MockAdminRepository)
	mockDramaRepo := new(MockDramaRepository)
	mockEpisodeRepo := new(MockEpisodeRepository)
	mockCacheService := new(MockCacheService)
	jwtManager := utils.NewJWTManager("test-secret", time.Hour)
	
	adminService := NewAdminService(mockAdminRepo, mockDramaRepo, mockEpisodeRepo, jwtManager, mockCacheService)

	t.Run("成功登录", func(t *testing.T) {
		req := models.AdminLoginRequest{
			Username: "admin",
			Password: "password123",
		}

		hashedPassword, _ := utils.HashPassword(req.Password)
		admin := &models.Admin{
			ID:       1,
			Username: req.Username,
			Password: hashedPassword,
			IsActive: true,
		}

		mockAdminRepo.On("GetByUsername", req.Username).Return(admin, nil)

		response, err := adminService.Login(req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.Token)
		assert.NotNil(t, response.User)

		mockAdminRepo.AssertExpectations(t)
	})

	t.Run("管理员不存在", func(t *testing.T) {
		mockAdminRepo := new(MockAdminRepository)
		adminService := NewAdminService(mockAdminRepo, mockDramaRepo, mockEpisodeRepo, jwtManager, mockCacheService)
		
		req := models.AdminLoginRequest{
			Username: "nonexistent",
			Password: "password123",
		}

		mockAdminRepo.On("GetByUsername", req.Username).Return((*models.Admin)(nil), errors.New("管理员不存在"))

		response, err := adminService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "用户名或密码错误")

		mockAdminRepo.AssertExpectations(t)
	})

	t.Run("管理员已被禁用", func(t *testing.T) {
		mockAdminRepo := new(MockAdminRepository)
		adminService := NewAdminService(mockAdminRepo, mockDramaRepo, mockEpisodeRepo, jwtManager, mockCacheService)
		
		req := models.AdminLoginRequest{
			Username: "admin",
			Password: "password123",
		}

		hashedPassword, _ := utils.HashPassword(req.Password)
		admin := &models.Admin{
			ID:       1,
			Username: req.Username,
			Password: hashedPassword,
			IsActive: false, // 管理员被禁用
		}

		mockAdminRepo.On("GetByUsername", req.Username).Return(admin, nil)

		response, err := adminService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "管理员账户已被禁用")

		mockAdminRepo.AssertExpectations(t)
	})
}

func TestAdminService_CreateDrama(t *testing.T) {
	mockAdminRepo := new(MockAdminRepository)
	mockDramaRepo := new(MockDramaRepository)
	mockEpisodeRepo := new(MockEpisodeRepository)
	mockCacheService := new(MockCacheService)
	jwtManager := utils.NewJWTManager("test-secret", time.Hour)
	
	adminService := NewAdminService(mockAdminRepo, mockDramaRepo, mockEpisodeRepo, jwtManager, mockCacheService)

	t.Run("成功创建短剧", func(t *testing.T) {
		req := models.CreateDramaRequest{
			Title:       "测试短剧",
			Description: "测试描述",
			Genre:       "喜剧",
		}

		mockDramaRepo.On("Create", mock.AnythingOfType("*models.Drama")).Return(nil)
		mockCacheService.On("DeletePattern", "dramas:*").Return(nil)
		mockCacheService.On("DeletePattern", "popular_dramas:*").Return(nil)

		drama, err := adminService.CreateDrama(req)

		assert.NoError(t, err)
		assert.NotNil(t, drama)
		assert.Equal(t, req.Title, drama.Title)
		assert.Equal(t, req.Description, drama.Description)
		assert.Equal(t, "active", drama.Status) // 默认状态

		mockDramaRepo.AssertExpectations(t)
		mockCacheService.AssertExpectations(t)
	})
}

func TestAdminService_CreateEpisode(t *testing.T) {
	mockAdminRepo := new(MockAdminRepository)
	mockDramaRepo := new(MockDramaRepository)
	mockEpisodeRepo := new(MockEpisodeRepository)
	mockCacheService := new(MockCacheService)
	jwtManager := utils.NewJWTManager("test-secret", time.Hour)
	
	adminService := NewAdminService(mockAdminRepo, mockDramaRepo, mockEpisodeRepo, jwtManager, mockCacheService)

	t.Run("成功创建剧集", func(t *testing.T) {
		req := models.CreateEpisodeRequest{
			DramaID:    1,
			Title:      "测试剧集",
			EpisodeNum: 1,
			Duration:   30,
		}

		drama := &models.Drama{ID: 1, Title: "测试短剧"}

		mockDramaRepo.On("GetByID", req.DramaID).Return(drama, nil)
		mockEpisodeRepo.On("ExistsByDramaIDAndEpisodeNum", req.DramaID, req.EpisodeNum).Return(false, nil)
		mockEpisodeRepo.On("Create", mock.AnythingOfType("*models.Episode")).Return(nil)
		mockCacheService.On("Delete", "drama_with_episodes:1").Return(nil)
		mockCacheService.On("DeletePattern", "episodes:drama:1:*").Return(nil)

		episode, err := adminService.CreateEpisode(req)

		assert.NoError(t, err)
		assert.NotNil(t, episode)
		assert.Equal(t, req.Title, episode.Title)
		assert.Equal(t, req.DramaID, episode.DramaID)
		assert.Equal(t, "active", episode.Status) // 默认状态

		mockDramaRepo.AssertExpectations(t)
		mockEpisodeRepo.AssertExpectations(t)
		mockCacheService.AssertExpectations(t)
	})

	t.Run("短剧不存在", func(t *testing.T) {
		mockDramaRepo := new(MockDramaRepository)
		mockEpisodeRepo := new(MockEpisodeRepository)
		adminService := NewAdminService(mockAdminRepo, mockDramaRepo, mockEpisodeRepo, jwtManager, mockCacheService)
		
		req := models.CreateEpisodeRequest{
			DramaID:    999,
			Title:      "测试剧集",
			EpisodeNum: 1,
		}

		mockDramaRepo.On("GetByID", req.DramaID).Return((*models.Drama)(nil), errors.New("短剧不存在"))

		episode, err := adminService.CreateEpisode(req)

		assert.Error(t, err)
		assert.Nil(t, episode)
		assert.Contains(t, err.Error(), "短剧不存在")

		mockDramaRepo.AssertExpectations(t)
	})

	t.Run("剧集编号已存在", func(t *testing.T) {
		mockDramaRepo := new(MockDramaRepository)
		mockEpisodeRepo := new(MockEpisodeRepository)
		adminService := NewAdminService(mockAdminRepo, mockDramaRepo, mockEpisodeRepo, jwtManager, mockCacheService)
		
		req := models.CreateEpisodeRequest{
			DramaID:    1,
			Title:      "测试剧集",
			EpisodeNum: 1,
		}

		drama := &models.Drama{ID: 1, Title: "测试短剧"}

		mockDramaRepo.On("GetByID", req.DramaID).Return(drama, nil)
		mockEpisodeRepo.On("ExistsByDramaIDAndEpisodeNum", req.DramaID, req.EpisodeNum).Return(true, nil)

		episode, err := adminService.CreateEpisode(req)

		assert.Error(t, err)
		assert.Nil(t, episode)
		assert.Contains(t, err.Error(), "该剧集编号已存在")

		mockDramaRepo.AssertExpectations(t)
		mockEpisodeRepo.AssertExpectations(t)
	})
}