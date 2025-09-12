package repository

import (
	"errors"
	"gorm.io/gorm"
	"gin-mysql-api/internal/models"
)

// dramaRepository 短剧仓库实现
type dramaRepository struct {
	db *gorm.DB
}

// NewDramaRepository 创建短剧仓库实例
func NewDramaRepository(db *gorm.DB) DramaRepository {
	return &dramaRepository{db: db}
}

// Create 创建短剧
func (r *dramaRepository) Create(drama *models.Drama) error {
	return r.db.Create(drama).Error
}

// GetByID 根据ID获取短剧
func (r *dramaRepository) GetByID(id uint) (*models.Drama, error) {
	var drama models.Drama
	if err := r.db.First(&drama, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &drama, nil
}

// GetByIDWithEpisodes 根据ID获取短剧（包含剧集信息）
func (r *dramaRepository) GetByIDWithEpisodes(id uint) (*models.Drama, error) {
	var drama models.Drama
	if err := r.db.Preload("Episodes", "status = ?", "active").First(&drama, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &drama, nil
}

// GetList 获取短剧列表（分页，可按类型筛选）
func (r *dramaRepository) GetList(offset, limit int, genre string) ([]models.Drama, int64, error) {
	var dramas []models.Drama
	var total int64
	
	query := r.db.Model(&models.Drama{})
	
	// 如果指定了类型，添加筛选条件
	if genre != "" {
		query = query.Where("genre = ?", genre)
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据，按创建时间倒序
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&dramas).Error; err != nil {
		return nil, 0, err
	}
	
	return dramas, total, nil
}

// Update 更新短剧信息
func (r *dramaRepository) Update(drama *models.Drama) error {
	return r.db.Save(drama).Error
}

// Delete 删除短剧（软删除）
func (r *dramaRepository) Delete(id uint) error {
	return r.db.Delete(&models.Drama{}, id).Error
}

// IncrementViewCount 增加观看次数
func (r *dramaRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Drama{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// GetByGenre 根据类型获取短剧列表
func (r *dramaRepository) GetByGenre(genre string, offset, limit int) ([]models.Drama, int64, error) {
	var dramas []models.Drama
	var total int64
	
	query := r.db.Model(&models.Drama{}).Where("genre = ? AND status = ?", genre, "active")
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	if err := query.Order("view_count DESC, created_at DESC").
		Offset(offset).Limit(limit).Find(&dramas).Error; err != nil {
		return nil, 0, err
	}
	
	return dramas, total, nil
}

// GetActiveList 获取活跃状态的短剧列表
func (r *dramaRepository) GetActiveList(offset, limit int) ([]models.Drama, int64, error) {
	var dramas []models.Drama
	var total int64
	
	query := r.db.Model(&models.Drama{}).Where("status = ?", "active")
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据，按观看次数和创建时间排序
	if err := query.Order("view_count DESC, created_at DESC").
		Offset(offset).Limit(limit).Find(&dramas).Error; err != nil {
		return nil, 0, err
	}
	
	return dramas, total, nil
}