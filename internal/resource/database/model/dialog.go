package model

import (
	"github.com/cloudwego/eino/schema"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
)

// Dialog 对话
//
//	author centonhuang
//	update 2025-11-13 19:10:00
type Dialog struct {
	BaseModel
	ID             uint              `json:"id" gorm:"column:id;primary_key;auto_increment;comment:用户ID"`
	UserID         uint              `json:"user_id" gorm:"column:user_id;not null;comment:用户ID"`
	InputMessages  []*schema.Message `json:"input_messages" gorm:"column:input_messages;not null;serializer:json;comment:输入消息"`
	OutputMessages []*schema.Message `json:"output_messages" gorm:"column:output_messages;not null;serializer:json;comment:输出消息"`
	Status         enum.DialogStatus `json:"status" gorm:"column:status;not null;comment:状态"`
}
