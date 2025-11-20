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

// UpdatedTodoItem 更新待办事项实体
//
//	@author centonhuang
//	@update 2025-11-07 15:21:39
type UpdatedTodoItem struct {
	ID     uint                `json:"id" minimum:"1" required:"true" doc:"Unique identifier for the todo item"`
	Status enum.TodoItemStatus `json:"status,omitempty" enum:"pending,completed,cancelled,timeout" doc:"Status of the todo item"`
	TodoItem
}

// TodoItemFilterParam 待办事项过滤参数
//
//	@author centonhuang
//	@update 2025-11-20 15:10:52
type TodoItemFilterParam struct {
	Priority enum.TodoItemPriority `query:"priority" enum:"low,medium,high,urgent" doc:"Priority filter"`
	Status   enum.TodoItemStatus   `query:"status" enum:"pending,completed,cancelled" doc:"Status filter"`
}

// DatailedTodoItem 详细待办事项实体
//
//	@author centonhuang
//	@update 2025-11-07 02:45:30
type DatailedTodoItem struct {
	ID        uint                `json:"id" doc:"Unique identifier for the todo item"`
	CreatedAt string              `json:"createdAt,omitempty" doc:"Timestamp when the todo item was created"`
	UpdatedAt string              `json:"updatedAt,omitempty" doc:"Timestamp when the todo item was updated"`
	Status    enum.TodoItemStatus `json:"status" enum:"pending,completed,cancelled,timeout" doc:"Status of the todo item"`
	TodoItem
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
	TodoItems []*TodoItem `json:"todoItems" minItems:"1" maxItems:"100" required:"true" doc:"Items to create"`
}

// ListTodoItemsReq 获取待办事项列表请求
//
//	@author centonhuang
//	@update 2025-11-07 01:43:02
type ListTodoItemsReq struct {
	model.CommonParam
	SortField string `query:"sortField" enum:"id,createdAt,updatedAt" doc:"Sort field"`
	TodoItemFilterParam
}

// ListTodoItemsRsp 获取待办事项列表响应
//
//	@author centonhuang
//	@update 2025-11-07 01:43:02
type ListTodoItemsRsp struct {
	CommonRsp
	TodoItems []*DatailedTodoItem `json:"todoItems" doc:"Items to list"`
	PageInfo  *model.PageInfo     `json:"pageInfo,omitempty" doc:"Page info"`
}

// UpdateTodoItemReq 更新待办事项请求
//
//	@author centonhuang
//	@update 2025-11-14 10:11:00
type UpdateTodoItemReq struct {
	Body *UpdateTodoItemReqBody `json:"body" doc:"Request body containing fields to update"`
}

// UpdateTodoItemReqBody 更新待办事项请求体
//
//	@author centonhuang
//	@update 2025-11-14 10:11:00
type UpdateTodoItemReqBody struct {
	TodoItem *UpdatedTodoItem `json:"todoItem" required:"true" doc:"Todo item to update"`
}

// DeleteTodoItemReq 删除待办事项请求体
//
//	@author centonhuang
//	@update 2025-11-14 14:14:16
type DeleteTodoItemReq struct {
	ID uint `json:"id" query:"id" required:"true" minimum:"1" doc:"Unique identifier for the todo item"`
}
