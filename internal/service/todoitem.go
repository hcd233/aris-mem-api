package service

import (
	"context"
	"time"

	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/hcd233/aris-mem-api/internal/resource/database"
	"github.com/hcd233/aris-mem-api/internal/resource/database/dao"
	"github.com/hcd233/aris-mem-api/internal/resource/database/model"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// TodoItemService 待办事项服务
//
//	author centonhuang
//	update 2025-11-07 01:12:23
type TodoItemService interface {
	CreateTodoItems(ctx context.Context, req *dto.CreateTodoItemsReq) (rsp *dto.EmptyRsp, err error)
	ListTodoItems(ctx context.Context, req *dto.ListTodoItemsReq) (rsp *dto.ListTodoItemsRsp, err error)
	// UpdateTodoItem(ctx context.Context, req *dto.UpdateTodoItemReq) (rsp *dto.UpdateTodoItemResp, err error)
	// DeleteTodoItem(ctx context.Context, req *dto.DeleteTodoItemReq) (rsp *dto.DeleteTodoItemResp, err error)
}

type todoItemService struct {
	todoItemDAO *dao.TodoItemDAO
}

// NewTodoItemService 创建待办事项服务
//
//	return TodoItemService
//	author centonhuang
//	update 2025-11-07 01:12:23
func NewTodoItemService() TodoItemService {
	return &todoItemService{
		todoItemDAO: dao.GetTodoItemDAO(),
	}
}

// CreateTodoItems 创建待办事项
//
//	return *EmptyResp
//	author centonhuang
//	update 2025-11-07 01:12:23
func (s *todoItemService) CreateTodoItems(ctx context.Context, req *dto.CreateTodoItemsReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

	db := database.GetDBInstance(ctx)

	logger := logger.WithCtx(ctx)

	todoItems := lo.Map(req.Body.TodoItems, func(item *dto.TodoItem, _ int) *model.TodoItem {
		return &model.TodoItem{
			Name:     item.Name,
			Summary:  item.Summary,
			Content:  item.Content,
			Status:   enum.TodoItemStatusPending,
			Priority: item.Priority,
		}
	})
	err := s.todoItemDAO.BatchCreate(db, todoItems)
	if err != nil {
		logger.Error("[TodoItemService] failed to create todo items", zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	return rsp, nil
}

// ListTodoItems 获取待办事项列表
//
//	return *ListTodoItemsResp
//	author centonhuang
//	update 2025-11-07 01:12:23
func (s *todoItemService) ListTodoItems(ctx context.Context, req *dto.ListTodoItemsReq) (*dto.ListTodoItemsRsp, error) {
	rsp := &dto.ListTodoItemsRsp{}

	db := database.GetDBInstance(ctx)

	logger := logger.WithCtx(ctx)

	commonParam := &dao.CommonParam{
		PageParam: dao.PageParam{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		QueryParam: dao.QueryParam{
			Query:       req.Query,
			QueryFields: []string{"name", "summary", "content"},
		},
		SortParam: dao.SortParam{
			Sort:      req.Sort,
			SortField: strcase.ToSnake(req.SortField),
		},
	}

	todoItems, pageInfo, err := s.todoItemDAO.Paginate(db, []string{"id", "created_at", "updated_at", "name", "summary", "content", "status", "priority"}, []string{}, commonParam)
	if err != nil {
		logger.Error("[TodoItemService] failed to list todo items", zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	rsp.TodoItems = lo.Map(todoItems, func(item *model.TodoItem, _ int) *dto.DatailedTodoItem {
		return &dto.DatailedTodoItem{
			ID:        item.ID,
			Status:    item.Status,
			CreatedAt: item.CreatedAt.Format(time.RFC3339),
			UpdatedAt: item.UpdatedAt.Format(time.RFC3339),
			TodoItem: dto.TodoItem{
				Name:     item.Name,
				Summary:  item.Summary,
				Content:  item.Content,
				Priority: item.Priority,
			},
		}
	})
	rsp.PageInfo = pageInfo
	return rsp, nil
}
