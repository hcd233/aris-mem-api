package model

import "github.com/hcd233/aris-mem-api/internal/common/enum"

// Notification 通知
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type Notification struct {
	BaseModel
	ID         uint                        `json:"id" gorm:"column:id;primary_key;auto_increment;comment:ID"`
	SenderID   uint                        `json:"sender_id" gorm:"column:sender_id;not null;comment:发送者ID"`
	ReceiverID uint                        `json:"receiver_id" gorm:"column:receiver_id;not null;comment:接收者ID"`
	Status     enum.NotificationStatus     `json:"status" gorm:"column:status;not null;comment:状态"`
	Type       enum.NotificationType       `json:"type" gorm:"column:type;not null;comment:类型"`
	EntityType enum.NotificationEntityType `json:"entity_type" gorm:"column:entity_type;index:idx_entity_type_entity_id;not null;comment:实体类型"`
	EntityID   uint                        `json:"entity_id" gorm:"column:entity_id;index:idx_entity_type_entity_id;not null;comment:实体ID"`
}
