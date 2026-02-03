package model

// Comment 评论
//
//	@author centonhuang
//	@update 2026-02-03 21:56:31
type Comment struct {
	BaseModel
	ID        uint     `json:"id" gorm:"column:id;primary_key;auto_increment;comment:ID"`
	ArticleID uint     `json:"article_id" gorm:"column:article_id;not null;comment:文章ID"`
	UserID    uint     `json:"user_id" gorm:"column:user_id;not null;comment:用户ID"`
	ParentID  uint     `json:"parent_id" gorm:"column:parent_id;default:null;comment:父评论ID"`
	Content   string   `json:"content" gorm:"column:content;not null;comment:内容"`
	Images    []string `json:"images" gorm:"column:images;serializer:json;comment:封面图片URL"`
	Likes     uint     `json:"likes" gorm:"column:likes;not null;default:0;comment:点赞数"`
	Saves     uint     `json:"saves" gorm:"column:saves;not null;default:0;comment:收藏数"`
}
