package model

import (
	"time"

	"github.com/hcd233/aris-mem-api/internal/common/enum"
)

// Article 文章
//
//	author centonhuang
//	update 2025-11-13 19:10:00
//	@author centonhuang
//	@update 2026-01-29 12:00:00
type Article struct {
	BaseModel
	ID          uint               `json:"id" gorm:"column:id;primary_key;auto_increment;comment:用户ID"`
	UserID      uint               `json:"user_id" gorm:"column:user_id;not null;comment:用户ID"`
	Title       string             `json:"title" gorm:"column:title;not null;comment:标题"`
	Slug        string             `json:"slug" gorm:"column:slug;not null;comment:slug"`
	Content     string             `json:"content" gorm:"column:content;not null;comment:内容"`
	CoverImage  string             `json:"cover_image" gorm:"column:cover_image;comment:封面图片URL"`
	PublishedAt time.Time          `json:"published_at" gorm:"column:published_at;not null;comment:发布时间"`
	Status      enum.ArticleStatus `json:"status" gorm:"column:status;not null;comment:状态"`
}
