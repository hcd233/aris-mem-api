package dto

import (
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/common/model"
)

// TodoItem 待办事项实体
//
//	author centonhuang
//	update 2025-11-07 01:13:01
type TodoItem struct {
	Name     string                `json:"name" doc:"Name of the todo item"`
	Summary  string                `json:"summary" maxLength:"255" doc:"Summary of the todo item"`
	Content  string                `json:"content" maxLength:"4096" doc:"Content of the todo item"`
	Priority enum.TodoItemPriority `json:"priority" enum:"low,medium,high,urgent" doc:"Priority of the todo item"`
}

// DatailedTodoItem 详细待办事项实体
//
//	@author centonhuang
//	@update 2025-11-07 02:45:30
type DatailedTodoItem struct {
	TodoItem
	ID        uint                `json:"id" doc:"Unique identifier for the todo item"`
	Status    enum.TodoItemStatus `json:"status" enum:"pending,completed,cancelled,timeout" doc:"Status of the todo item"`
	CreatedAt string              `json:"createdAt,omitempty" doc:"Timestamp when the todo item was created"`
	UpdatedAt string              `json:"updatedAt,omitempty" doc:"Timestamp when the todo item was updated"`
}

// CreateTodoItemsReq 创建待办事项请求
//
//	@author centonhuang
//	@update 2025-11-07 01:37:12
type CreateTodoItemsReq struct {
	Body *CreateTodoItemsReqBody `json:"body" doc:"Request body containing fields to create"`
}

// CreateTodoItemsReqBody 创建待办事项请求体
//
//	@author centonhuang
//	@update 2025-11-07 01:37:12
type CreateTodoItemsReqBody struct {
	TodoItems []*TodoItem `json:"todoItems" minItems:"1" maxItems:"100" doc:"Items to create"`
}

// ListTodoItemsReq 获取待办事项列表请求
//
//	@author centonhuang
//	@update 2025-11-07 01:43:02
type ListTodoItemsReq struct {
	model.CommonParam
	SortField string `query:"sortField" enum:"id,createdAt,updatedAt" doc:"Sort field"`
}

// ListTodoItemsResp 获取待办事项列表响应
//
//	@author centonhuang
//	@update 2025-11-07 01:43:02
type ListTodoItemsResp struct {
	TodoItems []*DatailedTodoItem `json:"todoItems" doc:"Items to list"`
	PageInfo  *model.PageInfo     `json:"pageInfo" doc:"Page info"`
}
