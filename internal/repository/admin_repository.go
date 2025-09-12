package repository

import (
	"errors"
	"gorm.io/gorm"
	"gin-mysql-api/internal/models"
)

// adminRepository 管理员仓库实现
type adminRepository struct {
	db *gorm.DB
}

// NewAdminRepository 创建管理员仓库实例
func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

// Create 创建管理员
func (r *adminRepository) Create(admin *models.Admin) error {
	return r.db.Create(admin).Error
}

// GetByID 根据ID获取管理员
func (r *adminRepository) GetByID(id uint) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.First(&admin, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

// GetByEmail 根据邮箱获取管理员
func (r *adminRepository) GetByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.Where("email = ?", email).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

// GetByUsername 根据用户名获取管理员
func (r *adminRepository) GetByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.Where("username = ?", username).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

// Update 更新管理员信息
func (r *adminRepository) Update(admin *models.Admin) error {
	return r.db.Save(admin).Error
}

// Delete 删除管理员（软删除）
func (r *adminRepository) Delete(id uint) error {
	return r.db.Delete(&models.Admin{}, id).Error
}

// List 获取管理员列表（分页）
func (r *adminRepository) List(offset, limit int) ([]models.Admin, int64, error) {
	var admins []models.Admin
	var total int64
	
	// 获取总数
	if err := r.db.Model(&models.Admin{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据，按创建时间倒序
	if err := r.db.Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&admins).Error; err != nil {
		return nil, 0, err
	}
	
	return admins, total, nil
}

// ExistsByEmail 检查邮箱是否已存在
func (r *adminRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Admin{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByUsername 检查用户名是否已存在
func (r *adminRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Admin{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}