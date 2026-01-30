package dao

import (
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	dbmodel "github.com/hcd233/aris-mem-api/internal/infrastructure/database/model"
	"gorm.io/gorm"
)

// ActionDAO user action DAO
//
//	author centonhuang
//	update 2026-01-30 21:00:00
type ActionDAO struct {
	baseDAO[dbmodel.Action]
}

// BatchGetByUserIDAndActionType 根据用户ID、动作类型和实体ID列表批量查询动作
//
//	@receiver dao *ActionDAO
//	@param db
//	@param userID
//	@param actionType
//	@param entityIDs
//	@param fields
//	@return actions
//	@return err
//	@author centonhuang
//	@update 2026-01-30 17:04:18
func (dao *ActionDAO) BatchGetByUserIDAndActionType(db *gorm.DB, userID uint, actionType enum.ActionType, entityType enum.ActionEntityType, entityIDs []uint, fields []string) (data []*dbmodel.Action, err error) {
	if len(entityIDs) == 0 {
		return nil, nil
	}
	sql := db.Select(fields)
	sql = sql.Where("user_id = ? AND action_type = ?", userID, actionType)
	sql = sql.Where("entity_type = ? AND entity_id IN ?", entityType, entityIDs)
	sql = sql.Where("deleted_at = 0")
	err = sql.Find(&data).Error
	return
}
