package repository

import (
	"errors"
	"gin-mysql-api/internal/models"

	"gorm.io/gorm"
)

// episodeRepository 剧集仓库实现
type episodeRepository struct {
	db *gorm.DB
}

// NewEpisodeRepository 创建剧集仓库实例
func NewEpisodeRepository(db *gorm.DB) EpisodeRepository {
	return &episodeRepository{db: db}
}

// Create 创建剧集
func (r *episodeRepository) Create(episode *models.Episode) error {
	return r.db.Create(episode).Error
}

// GetByID 根据ID获取剧集
func (r *episodeRepository) GetByID(id uint) (*models.Episode, error) {
	var episode models.Episode
	if err := r.db.First(&episode, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &episode, nil
}

// GetByIDWithDrama 根据ID获取剧集（包含短剧信息）
func (r *episodeRepository) GetByIDWithDrama(id uint) (*models.Episode, error) {
	var episode models.Episode
	if err := r.db.Preload("Drama").First(&episode, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &episode, nil
}

// GetByDramaID 根据短剧ID获取所有剧集
func (r *episodeRepository) GetByDramaID(dramaID uint) ([]models.Episode, error) {
	var episodes []models.Episode
	if err := r.db.Where("drama_id = ? AND status = ?", dramaID, "published").
		Order("episode_num ASC").Find(&episodes).Error; err != nil {
		return nil, err
	}
	return episodes, nil
}

// GetByDramaIDPaginated 根据短剧ID获取剧集列表（分页）
func (r *episodeRepository) GetByDramaIDPaginated(dramaID uint, offset, limit int) ([]models.Episode, int64, error) {
	var episodes []models.Episode
	var total int64

	query := r.db.Model(&models.Episode{}).Where("drama_id = ?", dramaID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据，按剧集号排序
	if err := query.Order("episode_num ASC").
		Offset(offset).Limit(limit).Find(&episodes).Error; err != nil {
		return nil, 0, err
	}

	return episodes, total, nil
}

// GetList 获取所有剧集列表（分页）
func (r *episodeRepository) GetList(offset, limit int) ([]models.Episode, int64, error) {
	var episodes []models.Episode
	var total int64

	query := r.db.Model(&models.Episode{}).Preload("Drama")

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据，按创建时间倒序
	if err := query.Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&episodes).Error; err != nil {
		return nil, 0, err
	}

	return episodes, total, nil
}

// Update 更新剧集信息
func (r *episodeRepository) Update(episode *models.Episode) error {
	return r.db.Save(episode).Error
}

// Delete 删除剧集（软删除）
func (r *episodeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Episode{}, id).Error
}

// IncrementViewCount 增加观看次数
func (r *episodeRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Episode{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// GetMaxEpisodeNum 获取指定短剧的最大剧集号
func (r *episodeRepository) GetMaxEpisodeNum(dramaID uint) (int, error) {
	var maxEpisodeNum int
	if err := r.db.Model(&models.Episode{}).
		Where("drama_id = ?", dramaID).
		Select("COALESCE(MAX(episode_num), 0)").
		Scan(&maxEpisodeNum).Error; err != nil {
		return 0, err
	}
	return maxEpisodeNum, nil
}

// ExistsByDramaIDAndEpisodeNum 检查指定短剧的剧集号是否已存在
func (r *episodeRepository) ExistsByDramaIDAndEpisodeNum(dramaID uint, episodeNum int) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Episode{}).
		Where("drama_id = ? AND episode_num = ?", dramaID, episodeNum).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
