package models

import (
	"fmt"
	"time"
	"gorm.io/gorm"
)

// Episode 剧集模型
type Episode struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	DramaID     uint           `gorm:"not null;index" json:"drama_id" validate:"required"`
	Title       string         `gorm:"size:200;not null" json:"title" validate:"required,max=200"`
	EpisodeNum  int            `gorm:"not null;index" json:"episode_num" validate:"required,min=1"`
	Duration    int            `gorm:"not null" json:"duration" validate:"required,min=1"` // 时长（秒）
	VideoURL    string         `gorm:"size:500" json:"video_url"`
	Thumbnail   string         `gorm:"size:255" json:"thumbnail"`
	Status      string         `gorm:"size:20;default:'active';index" json:"status" validate:"oneof=active inactive draft"`
	ViewCount   int64          `gorm:"default:0" json:"view_count"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联关系
	Drama Drama `gorm:"foreignKey:DramaID;constraint:OnDelete:CASCADE" json:"drama,omitempty"`
}

// TableName 指定表名
func (Episode) TableName() string {
	return "episodes"
}

// ToJSON 序列化为 JSON 响应格式
func (e *Episode) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":          e.ID,
		"drama_id":    e.DramaID,
		"title":       e.Title,
		"episode_num": e.EpisodeNum,
		"duration":    e.Duration,
		"video_url":   e.VideoURL,
		"thumbnail":   e.Thumbnail,
		"status":      e.Status,
		"view_count":  e.ViewCount,
		"created_at":  e.CreatedAt,
		"updated_at":  e.UpdatedAt,
	}
}

// ToJSONWithDrama 序列化为包含短剧信息的 JSON 格式
func (e *Episode) ToJSONWithDrama() map[string]interface{} {
	result := e.ToJSON()
	if e.Drama.ID != 0 {
		result["drama"] = e.Drama.ToJSON()
	}
	return result
}

// IncrementViewCount 增加观看次数
func (e *Episode) IncrementViewCount(tx *gorm.DB) error {
	return tx.Model(e).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// GetFormattedDuration 获取格式化的时长（分:秒）
func (e *Episode) GetFormattedDuration() string {
	minutes := e.Duration / 60
	seconds := e.Duration % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}