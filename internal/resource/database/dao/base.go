// Package dao DAO
//
//	update 2024-10-17 02:31:49
package dao

import (
	"reflect"
	"time"

	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/common/model"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// baseDAO 基础DAO
//
//	author centonhuang
//	update 2024-10-17 02:32:22
type baseDAO[ModelT interface{}] struct{}

// Create 创建数据
//
//	param dao *BaseDAO[T]
//	return Create
//	author centonhuang
//	update 2024-10-17 02:51:49
func (dao *baseDAO[ModelT]) Create(db *gorm.DB, data *ModelT) (err error) {
	err = db.Create(&data).Error
	return
}

// BatchCreate 批量创建数据
//
//	@param dao *baseDAO[ModelT]
//	@return BatchCreate
//	@author centonhuang
//	@update 2025-11-07 01:57:42
func (dao *baseDAO[ModelT]) BatchCreate(db *gorm.DB, data []*ModelT) (err error) {
	err = db.Create(&data).Error
	return
}

// Update 使用ID更新数据
//
//	param dao *BaseDAO[T]
//	return Update
//	author centonhuang
//	update 2024-10-17 02:52:18
func (dao *baseDAO[ModelT]) Update(db *gorm.DB, data *ModelT, info map[string]interface{}) (err error) {
	updateAtField := "updated_at"
	info[updateAtField] = time.Now().UTC()

	sql := db.Model(data)
	selectFields := lo.Filter(lo.Keys(info), func(item string, _ int) bool {
		return !reflect.ValueOf(info[item]).IsZero()
	})
	sql = sql.Select(selectFields)
	err = sql.Updates(info).Error
	return
}

// Delete 删除
//
//	param dao *BaseDAO[T]
//	return Delete
//	author centonhuang
//	update 2024-10-17 02:52:33
func (dao *baseDAO[ModelT]) Delete(db *gorm.DB, data *ModelT) (err error) {
	err = db.Delete(&data).Error
	return
}

func (dao *baseDAO[ModelT]) BatchDelete(db *gorm.DB, data *[]ModelT) (err error) {
	err = db.Delete(&data).Error
	return
}

// GetByID 使用ID查询指定数据
//
//	param dao *BaseDAO[T]
//	return GetByID
//	author centonhuang
//	update 2024-10-17 03:06:57
func (dao *baseDAO[ModelT]) GetByID(db *gorm.DB, id uint, fields []string, preloads []string) (data *ModelT, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}

	err = sql.Where("id = ?", id).First(&data).Error
	return
}

// BatchGetByIDs 批量使用ID查询指定数据
//
//	param dao *baseDAO[T]
//	return BatchGetByIDs
//	author centonhuang
//	update 2024-11-03 07:34:47
func (dao *baseDAO[ModelT]) BatchGetByIDs(db *gorm.DB, ids []uint, fields []string, preloads []string) (data *[]ModelT, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where("id IN ?", ids).Find(&data).Error
	return
}

// Paginate 分页查询
//
//	param dao *BaseDAO[T]
//	return Paginate
//	author centonhuang
//	update 2024-10-17 03:09:11
func (dao *baseDAO[ModelT]) Paginate(db *gorm.DB, fields []string, preloads []string, param *CommonParam) (data []*ModelT, pageInfo *model.PageInfo, err error) {
	limit, offset := param.PageSize, (param.Page-1)*param.PageSize

	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
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
