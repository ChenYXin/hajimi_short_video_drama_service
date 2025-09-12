package models

import (
	"time"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username" validate:"required,min=3,max=50"`
	Email     string         `gorm:"uniqueIndex;size:100;not null" json:"email" validate:"required,email"`
	Password  string         `gorm:"size:255;not null" json:"-" validate:"required,min=6"`
	Phone     string         `gorm:"size:20" json:"phone" validate:"omitempty,len=11"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// ToJSON 序列化为 JSON 响应格式（隐藏敏感信息）
func (u *User) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"username":   u.Username,
		"email":      u.Email,
		"phone":      u.Phone,
		"avatar":     u.Avatar,
		"is_active":  u.IsActive,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	}
}

// BeforeCreate GORM 钩子：创建前处理
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 这里可以添加密码加密等逻辑
	return nil
}