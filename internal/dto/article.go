package dto

import (
	"time"

	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/common/model"
)

// Article 文章基础实体（用于详情返回）
//
//	@author centonhuang
//	@update 2026-01-29 17:00:00
type Article struct {
	Title      string `json:"title" doc:"Title of the article"`
	Content    string `json:"content" doc:"Content of the article"`
	CoverImage string `json:"coverImage,omitempty" doc:"Cover image URL of the article (presigned URL)"`
}

// ListedArticle 详细文章实体
//
//	@author centonhuang
//	@update 2026-01-29 12:00:00
type ListedArticle struct {
	ID          uint      `json:"id" doc:"ID of the article"`
	Slug        string    `json:"slug" doc:"Slug of the article"`
	Title       string    `json:"title" doc:"Title of the article"`
	CoverImage  string    `json:"coverImage,omitempty" doc:"Cover image URL of the article"`
	Author      *User     `json:"author" doc:"Author of the article"`
	CreatedAt   time.Time `json:"createdAt" doc:"Created time of the article"`
	UpdatedAt   time.Time `json:"updatedAt" doc:"Updated time of the article"`
	PublishedAt time.Time `json:"publishedAt" doc:"Published time of the article"`
	Likes       uint      `json:"likes" doc:"Likes of the article"`
}

// DetailedArticle 详细文章实体
//
//	@author centonhuang
//	@update 2026-01-29 11:00:50
type DetailedArticle struct {
	ID          uint               `json:"id" doc:"ID of the article"`
	Slug        string             `json:"slug" doc:"Slug of the article"`
	Author      *User              `json:"author" doc:"Author of the article"`
	CreatedAt   time.Time          `json:"createdAt" doc:"Created time of the article"`
	UpdatedAt   time.Time          `json:"updatedAt" doc:"Updated time of the article"`
	PublishedAt time.Time          `json:"publishedAt" doc:"Published time of the article"`
	Status      enum.ArticleStatus `json:"status" doc:"Status of the article"`
	Tags        []*DetailedTag     `json:"tags" doc:"Tags of the article"`
	Likes       uint               `json:"likes" doc:"Likes of the article"`
	Saves       uint               `json:"saves" doc:"Saves of the article"`
	Views       uint               `json:"views" doc:"Views of the article"`
	Article
}

// CreateArticleReq 创建文章请求
//
//	@author centonhuang
//	@update 2026-01-29 17:00:00
type CreateArticleReq struct {
	Body *CreateArticleReqBody `json:"body" doc:"Request body containing fields to create"`
}

// CreateArticleReqBody 创建文章请求体
//
//	@author centonhuang
//	@update 2026-01-29 17:00:00
type CreateArticleReqBody struct {
	Title      string `json:"title" doc:"Title of the article"`
	Content    string `json:"content" doc:"Content of the article"`
	CoverImage []byte `json:"coverImage" maxItems:"10485760" doc:"Cover image file (max 10MB), optional"`
}

// ListArticlesReq 获取文章列表请求
//
//	@author centonhuang
//	@update 2026-01-29 10:00:00
type ListArticlesReq struct {
	model.CommonParam
	SortField string `query:"sortField" enum:"id,createdAt,updatedAt,name" doc:"Sort field"`
	TagSlug   string `query:"tagSlug" doc:"Filter by tag slug"`
}

// ListArticlesRsp 获取文章列表响应
//
//	@author centonhuang
//	@update 2026-01-29 10:00:00
type ListArticlesRsp struct {
	CommonRsp
	Articles []*ListedArticle `json:"articles" doc:"Articles to list"`
	PageInfo *model.PageInfo  `json:"pageInfo" doc:"Page info"`
}

// UpdateArticleReq 更新文章请求
//
//	@author centonhuang
//	@update 2026-01-29 17:00:00
type UpdateArticleReq struct {
	Body *UpdateArticleReqBody `json:"body" doc:"Request body containing fields to update"`
}

// UpdateArticleReqBody 更新文章请求体
//
//	@author centonhuang
//	@update 2026-01-29 17:00:00
type UpdateArticleReqBody struct {
	ID         uint               `json:"id" doc:"ID of the article"`
	Status     enum.ArticleStatus `json:"status" doc:"Status of the article"`
	Title      string             `json:"title" doc:"Title of the article"`
	Content    string             `json:"content" doc:"Content of the article"`
	CoverImage []byte             `json:"coverImage" maxLength:"10485760" doc:"Cover image file (max 10MB), optional"`
}

// DeleteArticleReq 删除文章请求
//
//	@author centonhuang
//	@update 2026-01-29 10:00:00
type DeleteArticleReq struct {
	ID uint `json:"id" query:"id" required:"true" minimum:"1" doc:"Unique identifier for the article"`
}

// GetArticleReq 通过 slug 获取文章详情请求
//
//	@author centonhuang
//	@update 2026-01-29 11:30:00
type GetArticleReq struct {
	Slug string `json:"slug" query:"slug" required:"true" doc:"Unique slug of the article"`
}

// GetArticleRsp 通过 slug 获取文章详情响应
//
//	@author centonhuang
//	@update 2026-01-29 11:30:00
type GetArticleRsp struct {
	CommonRsp
	Article *DetailedArticle `json:"article" doc:"Detailed article information"`
}
