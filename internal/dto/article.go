package dto

import (
	"time"

	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/common/model"
)

// Article 文章实体
//
//	@author centonhuang
//	@update 2026-01-28 21:55:06
type Article struct {
	Title   string `json:"title" doc:"Title of the article"`
	Content string `json:"content" doc:"Content of the article"`
}

// UpdatedArticle 更新文章实体
//
//	@author centonhuang
//	@update 2026-01-28 21:55:37
type UpdatedArticle struct {
	ID     uint               `json:"id" doc:"ID of the article"`
	Status enum.ArticleStatus `json:"status" doc:"Status of the article"`
	Article
}

// DetailedArticle 详细文章实体
//
//	@author centonhuang
//	@update 2026-01-28 21:56:04
type DetailedArticle struct {
	ID          uint               `json:"id" doc:"ID of the article"`
	Slug        string             `json:"slug" doc:"Slug of the article"`
	CreatedAt   time.Time          `json:"created_at" doc:"Created time of the article"`
	UpdatedAt   time.Time          `json:"updated_at" doc:"Updated time of the article"`
	PublishedAt time.Time          `json:"published_at" doc:"Published time of the article"`
	Status      enum.ArticleStatus `json:"status" doc:"Status of the article"`
	Tags        []string           `json:"tags" doc:"Tags of the article"`
	Article
}

// CreateArticleReq 创建文章请求
//
//	@author centonhuang
//	@update 2026-01-29 10:00:00
type CreateArticleReq struct {
	Body *CreateArticleReqBody `json:"body" doc:"Request body containing fields to create"`
}

// CreateArticleReqBody 创建文章请求体
//
//	@author centonhuang
//	@update 2026-01-29 10:00:00
type CreateArticleReqBody struct {
	Article
}

// CreateArticleRsp 创建文章响应
//
//	@author centonhuang
//	@update 2026-01-29 10:00:00
type CreateArticleRsp struct {
	CommonRsp
	Article *DetailedArticle `json:"article,omitempty" doc:"Created article"`
}

// ListArticlesReq 获取文章列表请求
//
//	@author centonhuang
//	@update 2026-01-29 10:00:00
type ListArticlesReq struct {
	model.CommonParam
	SortField string `query:"sortField" enum:"id,createdAt,updatedAt,name" doc:"Sort field"`
	Tag string `query:"tagName" doc:"Filter by tag name"`
}

// ListArticlesRsp 获取文章列表响应
//
//	@author centonhuang
//	@update 2026-01-29 10:00:00
type ListArticlesRsp struct {
	CommonRsp
	Articles []*DetailedArticle `json:"articles" doc:"Articles to list"`
	PageInfo *model.PageInfo    `json:"pageInfo,omitempty" doc:"Page info"`
}

// UpdateArticleReq 更新文章请求
//
//	@author centonhuang
//	@update 2026-01-29 10:00:00
type UpdateArticleReq struct {
	Body *UpdateArticleReqBody `json:"body" doc:"Request body containing fields to update"`
}

// UpdateArticleReqBody 更新文章请求体
//
//	@author centonhuang
//	@update 2026-01-29 10:00:00
type UpdateArticleReqBody struct {
	UpdatedArticle
}

// DeleteArticleReq 删除文章请求
//
//	@author centonhuang
//	@update 2026-01-29 10:00:00
type DeleteArticleReq struct {
	ID uint `json:"id" query:"id" required:"true" minimum:"1" doc:"Unique identifier for the article"`
}
