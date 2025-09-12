package service

import (
	"fmt"
	"time"

	"gin-mysql-api/internal/models"
	"gin-mysql-api/internal/repository"
)

// DramaService 短剧服务接口
type DramaService interface {
	GetDramas(page, pageSize int, genre string) (*models.PaginatedDramas, error)
	GetDramaByID(id uint) (*models.Drama, error)
	GetDramaWithEpisodes(id uint) (*models.Drama, error)
	GetEpisodesByDramaID(dramaID uint, page, pageSize int) (*models.PaginatedEpisodes, error)
	GetEpisodeByID(id uint) (*models.Episode, error)
	IncrementDramaViewCount(dramaID uint) error
	IncrementEpisodeViewCount(episodeID uint) error
	SearchDramas(keyword string, page, pageSize int) (*models.PaginatedDramas, error)
	GetPopularDramas(page, pageSize int) (*models.PaginatedDramas, error)
}

// dramaService 短剧服务实现
type dramaService struct {
	dramaRepo   repository.DramaRepository
	episodeRepo repository.EpisodeRepository
	cacheService CacheService
}

// NewDramaService 创建新的短剧服务
func NewDramaService(
	dramaRepo repository.DramaRepository,
	episodeRepo repository.EpisodeRepository,
	cacheService CacheService,
) DramaService {
	return &dramaService{
		dramaRepo:    dramaRepo,
		episodeRepo:  episodeRepo,
		cacheService: cacheService,
	}
}

// GetDramas 获取短剧列表
func (s *dramaService) GetDramas(page, pageSize int, genre string) (*models.PaginatedDramas, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("dramas:page:%d:size:%d:genre:%s", page, pageSize, genre)
	var cachedResult models.PaginatedDramas
	if s.cacheService != nil {
		err := s.cacheService.GetJSON(cacheKey, &cachedResult)
		if err == nil {
			return &cachedResult, nil
		}
	}

	offset := (page - 1) * pageSize
	var dramas []models.Drama
	var total int64
	var err error

	if genre != "" {
		dramas, total, err = s.dramaRepo.GetByGenre(genre, offset, pageSize)
	} else {
		dramas, total, err = s.dramaRepo.GetActiveList(offset, pageSize)
	}

	if err != nil {
		return nil, fmt.Errorf("获取短剧列表失败: %w", err)
	}

	totalPages := (int(total) + pageSize - 1) / pageSize

	result := &models.PaginatedDramas{
		Dramas:      dramas,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}

	// 缓存结果
	if s.cacheService != nil {
		s.cacheService.SetJSON(cacheKey, result, 5*time.Minute)
	}

	return result, nil
}

// GetDramaByID 根据ID获取短剧
func (s *dramaService) GetDramaByID(id uint) (*models.Drama, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("drama:%d", id)
	var cachedDrama models.Drama
	if s.cacheService != nil {
		err := s.cacheService.GetJSON(cacheKey, &cachedDrama)
		if err == nil {
			return &cachedDrama, nil
		}
	}

	drama, err := s.dramaRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("短剧不存在: %w", err)
	}

	// 缓存结果
	if s.cacheService != nil {
		s.cacheService.SetJSON(cacheKey, drama, 10*time.Minute)
	}

	return drama, nil
}

// GetDramaWithEpisodes 获取短剧及其剧集
func (s *dramaService) GetDramaWithEpisodes(id uint) (*models.Drama, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("drama_with_episodes:%d", id)
	var cachedDrama models.Drama
	if s.cacheService != nil {
		err := s.cacheService.GetJSON(cacheKey, &cachedDrama)
		if err == nil {
			return &cachedDrama, nil
		}
	}

	drama, err := s.dramaRepo.GetByIDWithEpisodes(id)
	if err != nil {
		return nil, fmt.Errorf("短剧不存在: %w", err)
	}

	// 缓存结果
	if s.cacheService != nil {
		s.cacheService.SetJSON(cacheKey, drama, 10*time.Minute)
	}

	return drama, nil
}

// GetEpisodesByDramaID 获取短剧的剧集列表
func (s *dramaService) GetEpisodesByDramaID(dramaID uint, page, pageSize int) (*models.PaginatedEpisodes, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 检查短剧是否存在
	_, err := s.dramaRepo.GetByID(dramaID)
	if err != nil {
		return nil, fmt.Errorf("短剧不存在: %w", err)
	}

	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("episodes:drama:%d:page:%d:size:%d", dramaID, page, pageSize)
	var cachedResult models.PaginatedEpisodes
	if s.cacheService != nil {
		err := s.cacheService.GetJSON(cacheKey, &cachedResult)
		if err == nil {
			return &cachedResult, nil
		}
	}

	offset := (page - 1) * pageSize
	episodes, total, err := s.episodeRepo.GetByDramaIDPaginated(dramaID, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取剧集列表失败: %w", err)
	}

	totalPages := (int(total) + pageSize - 1) / pageSize

	result := &models.PaginatedEpisodes{
		Episodes:    episodes,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}

	// 缓存结果
	if s.cacheService != nil {
		s.cacheService.SetJSON(cacheKey, result, 5*time.Minute)
	}

	return result, nil
}

// GetEpisodeByID 根据ID获取剧集
func (s *dramaService) GetEpisodeByID(id uint) (*models.Episode, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("episode:%d", id)
	var cachedEpisode models.Episode
	if s.cacheService != nil {
		err := s.cacheService.GetJSON(cacheKey, &cachedEpisode)
		if err == nil {
			return &cachedEpisode, nil
		}
	}

	episode, err := s.episodeRepo.GetByIDWithDrama(id)
	if err != nil {
		return nil, fmt.Errorf("剧集不存在: %w", err)
	}

	// 缓存结果
	if s.cacheService != nil {
		s.cacheService.SetJSON(cacheKey, episode, 10*time.Minute)
	}

	return episode, nil
}

// IncrementDramaViewCount 增加短剧观看次数
func (s *dramaService) IncrementDramaViewCount(dramaID uint) error {
	err := s.dramaRepo.IncrementViewCount(dramaID)
	if err != nil {
		return fmt.Errorf("更新短剧观看次数失败: %w", err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.Delete(fmt.Sprintf("drama:%d", dramaID))
		s.cacheService.Delete(fmt.Sprintf("drama_with_episodes:%d", dramaID))
		// 清除列表缓存（简单处理，实际可以更精细）
		s.clearDramaListCache()
	}

	return nil
}

// IncrementEpisodeViewCount 增加剧集观看次数
func (s *dramaService) IncrementEpisodeViewCount(episodeID uint) error {
	err := s.episodeRepo.IncrementViewCount(episodeID)
	if err != nil {
		return fmt.Errorf("更新剧集观看次数失败: %w", err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.Delete(fmt.Sprintf("episode:%d", episodeID))
		// 获取剧集信息以清除相关缓存
		episode, err := s.episodeRepo.GetByID(episodeID)
		if err == nil {
			s.cacheService.Delete(fmt.Sprintf("drama_with_episodes:%d", episode.DramaID))
		}
	}

	return nil
}

// SearchDramas 搜索短剧
func (s *dramaService) SearchDramas(keyword string, page, pageSize int) (*models.PaginatedDramas, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 这里简化实现，实际应该在 repository 层实现搜索功能
	// 暂时返回所有短剧
	return s.GetDramas(page, pageSize, "")
}

// GetPopularDramas 获取热门短剧
func (s *dramaService) GetPopularDramas(page, pageSize int) (*models.PaginatedDramas, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("popular_dramas:page:%d:size:%d", page, pageSize)
	var cachedResult models.PaginatedDramas
	if s.cacheService != nil {
		err := s.cacheService.GetJSON(cacheKey, &cachedResult)
		if err == nil {
			return &cachedResult, nil
		}
	}

	// 这里简化实现，实际应该按观看次数排序
	result, err := s.GetDramas(page, pageSize, "")
	if err != nil {
		return nil, err
	}

	// 缓存结果（热门内容缓存时间更长）
	if s.cacheService != nil {
		s.cacheService.SetJSON(cacheKey, result, 30*time.Minute)
	}

	return result, nil
}

// clearDramaListCache 清除短剧列表缓存
func (s *dramaService) clearDramaListCache() {
	// 这里简化处理，实际应该更精细地管理缓存
	// 可以使用缓存标签或者模式匹配来批量清除
}