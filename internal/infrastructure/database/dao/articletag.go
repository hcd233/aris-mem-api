package dao

import (
	"time"

	"github.com/hcd233/aris-mem-api/internal/infrastructure/database/model"
	"gorm.io/gorm"
)

// ArticleTagDAO 文章标签关联DAO
//
//	author centonhuang
//	update 2026-01-29 10:00:00
type ArticleTagDAO struct {
	baseDAO[model.ArticleTag]
}

// DeleteByArticleID 根据文章ID软删除关联
//
//	param db *gorm.DB
//	param articleID uint
//	return error
//	author centonhuang
//	update 2026-01-29 10:00:00
func (dao *ArticleTagDAO) DeleteByArticleID(db *gorm.DB, articleID uint) error {
	return db.Model(&model.ArticleTag{}).Where("article_id = ?", articleID).Where("deleted_at = 0").Update("deleted_at", time.Now().UTC().Unix()).Error
}


// GetTagIDsByArticleID 获取文章关联的标签ID列表（仅未删除的）
//
//	param db *gorm.DB
//	param articleID uint
//	return []uint
//	return error
//	author centonhuang
//	update 2026-01-29 10:00:00
func (dao *ArticleTagDAO) GetTagIDsByArticleID(db *gorm.DB, articleID uint) ([]uint, error) {
	var articleTags []*model.ArticleTag
	err := db.Where("article_id = ? AND deleted_at = 0", articleID).Find(&articleTags).Error
	if err != nil {
		return nil, err
	}
	tagIDs := make([]uint, 0, len(articleTags))
	for _, at := range articleTags {
		tagIDs = append(tagIDs, at.TagID)
	}
	return tagIDs, nil
}

// GetArticleIDsByTagID 获取标签关联的文章ID列表（仅未删除的）
//
//	param db *gorm.DB
//	param tagID uint
//	return []uint
//	return error
//	author centonhuang
//	update 2026-01-29 10:00:00
func (dao *ArticleTagDAO) GetArticleIDsByTagID(db *gorm.DB, tagID uint) ([]uint, error) {
	var articleTags []*model.ArticleTag
	err := db.Where("tag_id = ? AND deleted_at = 0", tagID).Find(&articleTags).Error
	if err != nil {
		return nil, err
	}
	articleIDs := make([]uint, 0, len(articleTags))
	for _, at := range articleTags {
		articleIDs = append(articleIDs, at.ArticleID)
	}
	return articleIDs, nil
}

// GetArticleIDsByTagName 通过标签名称获取关联的文章ID列表（仅未删除的）
//
//	param db *gorm.DB
//	param tagName string
//	return []uint
//	return error
//	author centonhuang
//	update 2026-01-29 10:00:00
func (dao *ArticleTagDAO) GetArticleIDsByTagName(db *gorm.DB, tagName string) ([]uint, error) {
	var articleTags []*model.ArticleTag
	err := db.Table("article_tags").
		Select("article_tags.article_id").
		Joins("JOIN tags ON tags.id = article_tags.tag_id").
		Where("tags.name = ? AND tags.deleted_at = 0 AND article_tags.deleted_at = 0", tagName).
		Find(&articleTags).Error
	if err != nil {
		return nil, err
	}
	articleIDs := make([]uint, 0, len(articleTags))
	for _, at := range articleTags {
		articleIDs = append(articleIDs, at.ArticleID)
	}
	return articleIDs, nil
}
