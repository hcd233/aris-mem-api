package dao

import (
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/common/model"
	dbmodel "github.com/hcd233/aris-mem-api/internal/infrastructure/database/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// NotificationDAO Notification DAO
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type NotificationDAO struct {
	baseDAO[dbmodel.Notification]
}

// Paginate 分页查询
//
//	param dao *BaseDAO[T]
//	return Paginate
//	author centonhuang
//	update 2024-10-17 03:09:11
func (dao *NotificationDAO) Paginate(db *gorm.DB, where *dbmodel.Notification, fields []string, param *CommonParam) (data []*dbmodel.Notification, pageInfo *model.PageInfo, err error) {
	limit, offset := param.PageSize, (param.Page-1)*param.PageSize

	sql := db.Model(where).Select(fields).Where(where).Where("deleted_at = 0")

	category, ok := param.FieldValueMap["category"].(enum.NotificationCategory)
	// 如果有标签过滤，联表查询
	if ok && category != "" {
		switch category {
		case enum.NotificationCategoryLikeAndSave:
			sql = sql.Where("type IN (?)", []enum.NotificationType{enum.NotificationTypeLike, enum.NotificationTypeSave})
		case enum.NotificationCategoryCommentAndAt:
			sql = sql.Where("type IN (?)", []enum.NotificationType{enum.NotificationTypeComment, enum.NotificationTypeAt})
		}
	}

	delete(param.FieldValueMap, "category")

	for field, value := range param.FieldValueMap {
		if value == nil {
			sql = sql.Where(field + " IS NULL")
		} else {
			sql = sql.Where(field+" = ?", value)
		}
	}

	if param.Query != "" && len(param.QueryFields) > 0 {
		like := "%" + param.Query + "%"
		expressions := make([]clause.Expression, 0, len(param.QueryFields))
		for _, field := range param.QueryFields {
			if field == "" {
				continue
			}
			expressions = append(expressions, clause.Like{Column: clause.Column{Name: field}, Value: like})
		}

		if len(expressions) > 0 {
			sql = sql.Where(expressions[0])
			for _, expr := range expressions[1:] {
				sql = sql.Or(expr)
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

	err = sql.Count(&pageInfo.Total).Error
	if err != nil {
		return
	}

	err = sql.Limit(limit).Offset(offset).Find(&data).Error

	return
}
