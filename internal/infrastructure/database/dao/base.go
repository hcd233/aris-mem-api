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

// BatchGetByIDs 根据ID列表批量查询数据
//
//	param db *gorm.DB
//	param ids []uint
//	param fields []string
//	return data []*ModelT
//	return err error
//	author centonhuang
//	update 2026-01-29 10:00:00
func (dao *baseDAO[ModelT]) BatchGetByIDs(db *gorm.DB, ids []uint, fields []string) (data []*ModelT, err error) {
	if len(ids) == 0 {
		return []*ModelT{}, nil
	}
	err = db.Select(fields).Where("id IN ?", ids).Where("deleted_at = 0").Find(&data).Error
	return
}

// Update 使用ID更新数据
//
//	param dao *BaseDAO[T]
//	return Update
//	author centonhuang
//	update 2026-01-30 22:00:00
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
func (dao *baseDAO[ModelT]) Delete(db *gorm.DB, where *ModelT) (err error) {
	err = db.Where(where).Update("deleted_at", time.Now().UTC().Unix()).Error
	return
}

func (dao *baseDAO[ModelT]) BatchDelete(db *gorm.DB, data *[]ModelT) (err error) {
	err = db.Model(data).Update("deleted_at", time.Now().UTC().Unix()).Error
	return
}

// func GetByID 使用ID查询指定数据
//
//	param dao *BaseDAO[T]
//	return GetByID
//	author centonhuang
//	update 2024-10-17 03:06:57
//	@param dao
//	@return Get
//	@author centonhuang
//	@update 2025-11-14 16:05:03
func (dao *baseDAO[ModelT]) Get(db *gorm.DB, where *ModelT, fields []string) (data *ModelT, err error) {
	err = db.Select(fields).Where(where).Where("deleted_at = 0").First(&data).Error
	return
}

// GetOrCreate 获取或创建数据
//
//	param db *gorm.DB
//	param where *ModelT
//	param fields []string
//	return data *ModelT
//	return err error
//	author centonhuang
//	update 2026-01-29 10:00:00
func (dao *baseDAO[ModelT]) GetOrCreate(db *gorm.DB, where *ModelT, createData *ModelT, fields []string) (data *ModelT, err error) {
	// 先尝试获取
	data, err = dao.Get(db, where, fields)
	if err == nil {
		return
	}

	// 如果不存在则创建
	if err == gorm.ErrRecordNotFound {
		if err = db.Create(createData).Error; err != nil {
			return nil, err
		}
		return createData, nil
	}

	return
}

// Paginate 分页查询
//
//	param dao *BaseDAO[T]
//	return Paginate
//	author centonhuang
//	update 2024-10-17 03:09:11
func (dao *baseDAO[ModelT]) Paginate(db *gorm.DB, where *ModelT, fields []string, param *CommonParam) (data []*ModelT, pageInfo *model.PageInfo, err error) {
	limit, offset := param.PageSize, (param.Page-1)*param.PageSize

	sql := db.Model(where).Select(fields).Where(where).Where("deleted_at = 0")

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

// func Count 统计
//
//	@param dao
//	@return Count
//	@author centonhuang
//	@update 2026-02-12 14:48:14
func (dao *baseDAO[ModelT]) Count(db *gorm.DB, where *ModelT) (count int64, err error) {
	err = db.Model(where).Where("deleted_at = 0").Count(&count).Error
	return
}
