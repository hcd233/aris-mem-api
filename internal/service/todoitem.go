package service

import (
	"context"
	"errors"
	"reflect"
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
	"gorm.io/gorm"
)

// TodoItemService 待办事项服务
//
//	author centonhuang
//	update 2025-11-07 01:12:23
type TodoItemService interface {
	CreateTodoItems(ctx context.Context, req *dto.CreateTodoItemsReq) (rsp *dto.EmptyRsp, err error)
	ListTodoItems(ctx context.Context, req *dto.ListTodoItemsReq) (rsp *dto.ListTodoItemsRsp, err error)
	UpdateTodoItem(ctx context.Context, req *dto.UpdateTodoItemReq) (rsp *dto.EmptyRsp, err error)
	DeleteTodoItem(ctx context.Context, req *dto.DeleteTodoItemReq) (rsp *dto.EmptyRsp, err error)
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

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	todoItems := lo.Map(req.Body.TodoItems, func(item *dto.TodoItem, _ int) *model.TodoItem {
		return &model.TodoItem{
			Name:     item.Name,
			Summary:  item.Summary,
			Content:  item.Content,
			Status:   enum.TodoItemStatusPending,
			Priority: item.Priority,
			UserID:   userID,
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

	if req.SortField == "" {
		req.SortField = "id"
	}

	if req.Sort == "" {
		req.Sort = enum.SortAsc
	}

	db := database.GetDBInstance(ctx)

	logger := logger.WithCtx(ctx)

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

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

	todoItems, pageInfo, err := s.todoItemDAO.PaginateByUserID(db, userID, []string{"id", "created_at", "updated_at", "name", "summary", "content", "status", "priority"}, commonParam)
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

// UpdateTodoItem 更新待办事项
//
//	@return *dto.UpdateTodoItemRsp
//	@return error
//	@author
//	@update 2025-11-14 10:11:00
func (s *todoItemService) UpdateTodoItem(ctx context.Context, req *dto.UpdateTodoItemReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

	if req == nil || req.Body == nil || req.Body.TodoItem == nil || req.Body.TodoItem.ID == 0 {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	payload := req.Body.TodoItem

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	existing, err := s.todoItemDAO.GetByID(db, payload.ID, []string{"id", "user_id"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rsp.Error = constant.ErrDataNotExists
			return rsp, nil
		}
		logger.Error("[TodoItemService] failed to get todo item", zap.Error(err), zap.Uint("todoItemID", payload.ID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	if existing.UserID != userID {
		rsp.Error = constant.ErrNoPermission
		return rsp, nil
	}

	updateFields := map[string]interface{}{
		"name":     payload.Name,
		"summary":  payload.Summary,
		"content":  payload.Content,
		"priority": payload.Priority,
		"status":   payload.Status,
	}

	if !hasNonZeroValue(updateFields) {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	if err = s.todoItemDAO.Update(db, &model.TodoItem{ID: payload.ID}, updateFields); err != nil {
		logger.Error("[TodoItemService] failed to update todo item", zap.Error(err), zap.Uint("todoItemID", payload.ID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	return rsp, nil
}

// DeleteTodoItem 删除待办事项
//
//	@return *dto.EmptyRsp
//	@return error
//	@author
//	@update 2025-11-14 10:11:00
func (s *todoItemService) DeleteTodoItem(ctx context.Context, req *dto.DeleteTodoItemReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

	if req == nil || req.ID == 0 {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	todoItem, err := s.todoItemDAO.GetByID(db, req.ID, []string{"id", "user_id"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rsp.Error = constant.ErrDataNotExists
			return rsp, nil
		}
		logger.Error("[TodoItemService] failed to get todo item for deletion", zap.Error(err), zap.Uint("todoItemID", req.ID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	if todoItem.UserID != userID {
		rsp.Error = constant.ErrNoPermission
		return rsp, nil
	}

	if err = s.todoItemDAO.Delete(db, &model.TodoItem{ID: req.ID}); err != nil {
		logger.Error("[TodoItemService] failed to delete todo item", zap.Error(err), zap.Uint("todoItemID", req.ID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	return rsp, nil
}

func hasNonZeroValue(fields map[string]interface{}) bool {
	for _, value := range fields {
		if !reflect.ValueOf(value).IsZero() {
			return true
		}
	}
	return false
}
