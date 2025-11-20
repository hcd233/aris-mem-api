package handler

import (
	"context"

	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/service"
	"github.com/hcd233/aris-mem-api/internal/util"
)

// TodoItemHandler 用户处理器
//
//	author centonhuang
//	update 2025-01-04 15:56:20
type TodoItemHandler interface {
	HandleCreateTodoItems(ctx context.Context, req *dto.CreateTodoItemsReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
	HandleListTodoItems(ctx context.Context, req *dto.ListTodoItemsReq) (*dto.HTTPResponse[*dto.ListTodoItemsRsp], error)
	HandleUpdateTodoItem(ctx context.Context, req *dto.UpdateTodoItemReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
	HandleDeleteTodoItem(ctx context.Context, req *dto.DeleteTodoItemReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
}

type todoItemHandler struct {
	svc service.TodoItemService
}

// NewTodoItemHandler 创建待办事项处理器
//
//	return TodoItemHandler
//	author centonhuang
//	update 2024-12-08 16:59:38
func NewTodoItemHandler() TodoItemHandler {
	return &todoItemHandler{
		svc: service.NewTodoItemService(),
	}
}

func (h *todoItemHandler) HandleCreateTodoItems(ctx context.Context, req *dto.CreateTodoItemsReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
	return util.WrapHTTPResponse(h.svc.CreateTodoItems(ctx, req))
}

func (h *todoItemHandler) HandleListTodoItems(ctx context.Context, req *dto.ListTodoItemsReq) (*dto.HTTPResponse[*dto.ListTodoItemsRsp], error) {
	return util.WrapHTTPResponse(h.svc.ListTodoItems(ctx, req))
}

func (h *todoItemHandler) HandleUpdateTodoItem(ctx context.Context, req *dto.UpdateTodoItemReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
	return util.WrapHTTPResponse(h.svc.UpdateTodoItem(ctx, req))
}

func (h *todoItemHandler) HandleDeleteTodoItem(ctx context.Context, req *dto.DeleteTodoItemReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
	return util.WrapHTTPResponse(h.svc.DeleteTodoItem(ctx, req))
}
