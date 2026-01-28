package dto

import (
	"time"

	"github.com/hcd233/aris-mem-api/internal/common/enum"
)

// Article 文章实体
//
//	@author centonhuang
//	@update 2026-01-28 21:55:06
type Article struct {
	Title string `json:"title" doc:"Title of the article"`
	Content string `json:"content" doc:"Content of the article"`
}

// UpdatedArticle 更新文章实体
//
//	@author centonhuang 
//	@update 2026-01-28 21:55:37 
type UpdatedArticle struct {
	ID uint `json:"id" doc:"ID of the article"`
	Status enum.ArticleStatus `json:"status" doc:"Status of the article"`
	Article
}

// DetailedArticle 详细文章实体
//
//	@author centonhuang 
//	@update 2026-01-28 21:56:04 
type DetailedArticle struct {
	ID uint `json:"id" doc:"ID of the article"`
	CreatedAt time.Time `json:"created_at" doc:"Created time of the article"`
	UpdatedAt time.Time `json:"updated_at" doc:"Updated time of the article"`
	PublishedAt time.Time `json:"published_at" doc:"Published time of the article"`
	Status enum.ArticleStatus `json:"status" doc:"Status of the article"`
	Article
}