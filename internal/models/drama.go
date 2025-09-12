package models

import (
	"time"
	"gorm.io/gorm"
)

// Drama 短剧模型
type Drama struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:200;not null" json:"title" validate:"required,max=200"`
	Description string         `gorm:"type:text" json:"description"`
	CoverImage  string         `gorm:"size:255" json:"cover_image"`
	Director    string         `gorm:"size:100" json:"director"`
	Actors      string         `gorm:"size:500" json:"actors"`
	Genre       string         `gorm:"size:100;index" json:"genre" validate:"required"`
	Status      string         `gorm:"size:20;default:'active';index" json:"status" validate:"oneof=active inactive draft"`
	ViewCount   int64          `gorm:"default:0" json:"view_count"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联关系
	Episodes []Episode `gorm:"foreignKey:DramaID;constraint:OnDelete:CASCADE" json:"episodes,omitempty"`
}

// TableName 指定表名
func (Drama) TableName() string {
	return "dramas"
}

// ToJSON 序列化为 JSON 响应格式
func (d *Drama) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":          d.ID,
		"title":       d.Title,
		"description": d.Description,
		"cover_image": d.CoverImage,
		"director":    d.Director,
		"actors":      d.Actors,
		"genre":       d.Genre,
		"status":      d.Status,
		"view_count":  d.ViewCount,
		"created_at":  d.CreatedAt,
		"updated_at":  d.UpdatedAt,
	}
}

// ToJSONWithEpisodes 序列化为包含剧集信息的 JSON 格式
func (d *Drama) ToJSONWithEpisodes() map[string]interface{} {
	result := d.ToJSON()
	
	episodes := make([]map[string]interface{}, len(d.Episodes))
	for i, episode := range d.Episodes {
		episodes[i] = episode.ToJSON()
	}
	result["episodes"] = episodes
	
	return result
}

// IncrementViewCount 增加观看次数
func (d *Drama) IncrementViewCount(tx *gorm.DB) error {
	return tx.Model(d).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}