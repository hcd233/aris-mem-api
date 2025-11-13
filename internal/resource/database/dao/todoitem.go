package dao

import (
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/common/model"
	dbmodel "github.com/hcd233/aris-mem-api/internal/resource/database/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TodoItemDAO 待办事项DAO
//
//	author centonhuang
//	update 2024-10-17 02:30:24
type TodoItemDAO struct {
	baseDAO[dbmodel.TodoItem]
}

// PaginateByUserID 通过用户ID分页查询待办事项
//
//	param dao *BaseDAO[T]
//	return Paginate
//	author centonhuang
//	update 2024-10-17 03:09:11
//	@receiver dao *TodoItemDAO
//	@param db
//	@param userID
//	@param fields
//	@param param
//	@return data
//	@return pageInfo
//	@return err
//	@author centonhuang
//	@update 2025-11-13 19:33:00
func (dao *TodoItemDAO) PaginateByUserID(db *gorm.DB, userID uint, fields []string, param *CommonParam) (data []*dbmodel.TodoItem, pageInfo *model.PageInfo, err error) {
	limit, offset := param.PageSize, (param.Page-1)*param.PageSize

	sql := db.Select(fields).Where(dbmodel.TodoItem{UserID: userID})

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

	err = sql.Limit(limit).Offset(offset).Find(&data).Error
	if err != nil {
		return
	}

	pageInfo = &model.PageInfo{
		Page:     param.Page,
		PageSize: param.PageSize,
	}

	err = db.Model(&data).Count(&pageInfo.Total).Error

	return
}
