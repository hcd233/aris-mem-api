package model

// ArticleTag 文章标签关联表
//
//	@author centonhuang
//	@update 2026-01-29 00:15:41
type ArticleTag struct {
	BaseModel
	ArticleID uint
	TagID     uint
}
