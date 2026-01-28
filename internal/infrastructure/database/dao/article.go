package dao

import (
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/common/model"
	dbmodel "github.com/hcd233/aris-mem-api/internal/infrastructure/database/model"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ArticleDAO 文章DAO
//
//	author centonhuang
//	update 2026-01-29 10:00:00
type ArticleDAO struct {
	baseDAO[dbmodel.Article]
}

// Paginate 文章分页查询（支持联表查询标签）
//
//	param db *gorm.DB
//	param where *dbmodel.Article
//	param fields []string
//	param param *CommonParam
//	param tagName string
//	return data []*dbmodel.Article
//	return pageInfo *model.PageInfo
//	return err error
//	author centonhuang
//	update 2026-01-29 10:00:00
func (dao *ArticleDAO) Paginate(db *gorm.DB, where *dbmodel.Article, fields []string, param *CommonParam) (data []*dbmodel.Article, pageInfo *model.PageInfo, err error) {
	limit, offset := param.PageSize, (param.Page-1)*param.PageSize

	// 构建基础查询
	sql := db.Model(where).Select(lo.Map(fields, func (item string, _ int) string {
		return "articles."+item
	})).Where(where).Where("articles.deleted_at = 0")

	tag, ok := param.FieldValueMap["tag"].(string)
	// 如果有标签过滤，联表查询
	if ok && tag != "" {
		sql = sql.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Joins("JOIN tags ON tags.id = article_tags.tag_id").
			Where("tags.name = ? AND tags.deleted_at = 0 AND article_tags.deleted_at = 0", tag)
	}

	delete(param.FieldValueMap, "tag")

	for field, value := range param.FieldValueMap {
		sql = sql.Where("articles." + field + " = ?", value)
	}

	// 模糊搜索
	if param.Query != "" && len(param.QueryFields) > 0 {
		like := "%" + param.Query + "%"
		for i, field := range param.QueryFields {
			if field == "" {
				continue
			}
			if i == 0 {
				sql = sql.Where("articles."+field+" LIKE ?", like)
			} else {
				sql = sql.Or("articles."+field+" LIKE ?", like)
			}
		}
	}

	if param.Sort != "" && param.SortField != "" {
		sql = sql.Order(clause.OrderByColumn{Column: clause.Column{Name: param.SortField}, Desc: param.Sort == enum.SortDesc})
	}

	pageInfo = &model.PageInfo{
		Page:     param.Page,
		PageSize: param.PageSize,
	}

	// 计数
	if err = sql.Count(&pageInfo.Total).Error; err != nil {
		return
	}

	// 查询数据
	err = sql.Limit(limit).Offset(offset).Find(&data).Error
	return
}