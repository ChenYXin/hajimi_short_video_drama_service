package service

import (
	"testing"
	"time"

	"gin-mysql-api/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestDramaService_GetDramas(t *testing.T) {
	mockDramaRepo := new(MockDramaRepository)
	mockEpisodeRepo := new(MockEpisodeRepository)
	mockCacheService := new(MockCacheService)
	
	dramaService := NewDramaService(mockDramaRepo, mockEpisodeRepo, mockCacheService)

	t.Run("成功获取短剧列表", func(t *testing.T) {
		dramas := []models.Drama{
			{ID: 1, Title: "短剧1", Genre: "喜剧"},
			{ID: 2, Title: "短剧2", Genre: "爱情"},
		}

		// 设置缓存未命中
		mockCacheService.On("GetJSON", mock.AnythingOfType("string"), mock.Anything).Return(assert.AnError)
		
		// 设置仓库返回数据
		mockDramaRepo.On("GetActiveList", 0, 20).Return(dramas, int64(2), nil)
		
		// 设置缓存写入
		mockCacheService.On("SetJSON", mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)

		result, err := dramaService.GetDramas(1, 20, "")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Dramas, 2)
		assert.Equal(t, int64(2), result.Total)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 20, result.PageSize)

		mockDramaRepo.AssertExpectations(t)
		mockCacheService.AssertExpectations(t)
	})

	t.Run("按类型获取短剧列表", func(t *testing.T) {
		dramas := []models.Drama{
			{ID: 1, Title: "喜剧短剧", Genre: "喜剧"},
		}

		// 设置缓存未命中
		mockCacheService.On("GetJSON", mock.AnythingOfType("string"), mock.Anything).Return(assert.AnError)
		
		// 设置仓库返回数据
		mockDramaRepo.On("GetByGenre", "喜剧", 0, 20).Return(dramas, int64(1), nil)
		
		// 设置缓存写入
		mockCacheService.On("SetJSON", mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)

		result, err := dramaService.GetDramas(1, 20, "喜剧")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Dramas, 1)
		assert.Equal(t, "喜剧短剧", result.Dramas[0].Title)

		mockDramaRepo.AssertExpectations(t)
		mockCacheService.AssertExpectations(t)
	})
}

func TestDramaService_GetDramaByID(t *testing.T) {
	mockDramaRepo := new(MockDramaRepository)
	mockEpisodeRepo := new(MockEpisodeRepository)
	mockCacheService := new(MockCacheService)
	
	dramaService := NewDramaService(mockDramaRepo, mockEpisodeRepo, mockCacheService)

	t.Run("成功获取短剧详情", func(t *testing.T) {
		drama := &models.Drama{
			ID:    1,
			Title: "测试短剧",
			Genre: "喜剧",
		}

		// 设置缓存未命中
		mockCacheService.On("GetJSON", "drama:1", mock.Anything).Return(assert.AnError)
		
		// 设置仓库返回数据
		mockDramaRepo.On("GetByID", uint(1)).Return(drama, nil)
		
		// 设置缓存写入
		mockCacheService.On("SetJSON", "drama:1", drama, mock.AnythingOfType("time.Duration")).Return(nil)

		result, err := dramaService.GetDramaByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "测试短剧", result.Title)
		assert.Equal(t, uint(1), result.ID)

		mockDramaRepo.AssertExpectations(t)
		mockCacheService.AssertExpectations(t)
	})

	t.Run("短剧不存在", func(t *testing.T) {
		// 设置缓存未命中
		mockCacheService.On("GetJSON", "drama:999", mock.Anything).Return(assert.AnError)
		
		// 设置仓库返回错误
		mockDramaRepo.On("GetByID", uint(999)).Return((*models.Drama)(nil), assert.AnError)

		result, err := dramaService.GetDramaByID(999)

		assert.Error(t, err)
		assert.Nil(t, result)

		mockDramaRepo.AssertExpectations(t)
		mockCacheService.AssertExpectations(t)
	})
}

func TestDramaService_IncrementDramaViewCount(t *testing.T) {
	mockDramaRepo := new(MockDramaRepository)
	mockEpisodeRepo := new(MockEpisodeRepository)
	mockCacheService := new(MockCacheService)
	
	dramaService := NewDramaService(mockDramaRepo, mockEpisodeRepo, mockCacheService)

	t.Run("成功增加观看次数", func(t *testing.T) {
		// 设置仓库更新观看次数
		mockDramaRepo.On("IncrementViewCount", uint(1)).Return(nil)
		
		// 设置缓存清除
		mockCacheService.On("Delete", "drama:1").Return(nil)
		mockCacheService.On("Delete", "drama_with_episodes:1").Return(nil)

		err := dramaService.IncrementDramaViewCount(1)

		assert.NoError(t, err)

		mockDramaRepo.AssertExpectations(t)
		mockCacheService.AssertExpectations(t)
	})
}