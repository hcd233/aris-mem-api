package model

import "github.com/hcd233/aris-mem-api/internal/common/enum"

// TodoItem 待办事项数据库模型
//
//	author centonhuang
//	update 2024-06-22 09:36:22
type TodoItem struct {
	BaseModel
	ID       uint                  `json:"id" gorm:"column:id;primary_key;auto_increment;comment:用户ID"`
	Name     string                `json:"name" gorm:"column:name;not null;comment:用户名"`
	Summary  string                `json:"summary" gorm:"column:summary;not null;comment:摘要"`
	Content  string                `json:"content" gorm:"column:content;not null;comment:内容"`
	Status   enum.TodoItemStatus   `json:"status" gorm:"column:status;not null;comment:状态"`
	Priority enum.TodoItemPriority `json:"priority" gorm:"column:priority;not null;comment:优先级"`
}
