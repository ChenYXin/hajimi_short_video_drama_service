package models

import (
	"time"
	"gorm.io/gorm"
)

// Admin 管理员模型
type Admin struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username" validate:"required,min=3,max=50"`
	Email     string         `gorm:"uniqueIndex;size:100;not null" json:"email" validate:"required,email"`
	Password  string         `gorm:"size:255;not null" json:"-" validate:"required,min=6"`
	Role      string         `gorm:"size:20;default:'admin'" json:"role" validate:"oneof=admin super_admin"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Admin) TableName() string {
	return "admins"
}

// ToJSON 序列化为 JSON 响应格式（隐藏敏感信息）
func (a *Admin) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":         a.ID,
		"username":   a.Username,
		"email":      a.Email,
		"role":       a.Role,
		"is_active":  a.IsActive,
		"created_at": a.CreatedAt,
		"updated_at": a.UpdatedAt,
	}
}

// BeforeCreate GORM 钩子：创建前处理
func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	// 这里可以添加密码加密等逻辑
	return nil
}