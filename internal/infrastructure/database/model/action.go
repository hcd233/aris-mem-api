package model

import (
	"github.com/hcd233/aris-mem-api/internal/common/enum"
)

// Action user action unified model
//
//	author centonhuang
//	update 2026-01-30 21:00:00
type Action struct {
	BaseModel
	ID         uint                  `json:"id" gorm:"column:id;primary_key;auto_increment;comment:ID"`
	DeletedAt  int64                 `json:"deleted_at" gorm:"column:deleted_at;default:0;uniqueIndex:uidx_user_entity_action_deleted_at,priority:5;comment:删除时间，默认为0"`
	UserID     uint                  `json:"user_id" gorm:"column:user_id;not null;uniqueIndex:uidx_user_entity_action_deleted_at,priority:1;index:idx_user_action_deleted_at,priority:1;index:idx_entity_action_deleted_at,priority:4;comment:用户ID"`
	EntityType enum.ActionEntityType `json:"entity_type" gorm:"column:entity_type;type:varchar(32);not null;uniqueIndex:uidx_user_entity_action_deleted_at,priority:2;index:idx_user_action_deleted_at,priority:2;index:idx_entity_action_deleted_at,priority:1;comment:实体类型"`
	EntityID   uint                  `json:"entity_id" gorm:"column:entity_id;not null;uniqueIndex:uidx_user_entity_action_deleted_at,priority:3;index:idx_user_action_deleted_at,priority:3;index:idx_entity_action_deleted_at,priority:2;comment:实体ID"`
	ActionType enum.ActionType       `json:"action_type" gorm:"column:action_type;type:varchar(32);not null;uniqueIndex:uidx_user_entity_action_deleted_at,priority:4;index:idx_user_action_deleted_at,priority:4;index:idx_entity_action_deleted_at,priority:3;comment:动作类型"`
}
