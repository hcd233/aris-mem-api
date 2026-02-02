package service

import (
	"context"
	"errors"
	"time"

	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database/dao"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database/model"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UserService 用户服务
//
//	author centonhuang
//	update 2025-01-04 21:04:00
type UserService interface {
	GetCurrentUser(ctx context.Context, req *dto.EmptyReq) (rsp *dto.GetCurrentUserRsp, err error)
	UpdateUser(ctx context.Context, req *dto.UpdateUserReq) (rsp *dto.EmptyRsp, err error)
	ApproveUser(ctx context.Context, req *dto.ApproveUserReq) (rsp *dto.EmptyRsp, err error)
	ListUsers(ctx context.Context, req *dto.ListUsersReq) (rsp *dto.ListUsersRsp, err error)
}

type userService struct {
	userDAO *dao.UserDAO
}

// NewUserService 创建用户服务
//
//	return UserService
//	author centonhuang
//	update 2025-01-04 21:03:45
func NewUserService() UserService {
	return &userService{
		userDAO: dao.GetUserDAO(),
	}
}

// GetCurrentUser 获取当前用户信息
//
//	@receiver s *userService
//	@param ctx context.Context
//	@param _ *dto.EmptyReq
//	@return *dto.GetCurUserRsp
//	@return error
//	@author centonhuang
//	@update 2025-11-11 04:59:13
func (s *userService) GetCurrentUser(ctx context.Context, _ *dto.EmptyReq) (rsp *dto.GetCurrentUserRsp, err error) {
	rsp = &dto.GetCurrentUserRsp{}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	user, err := s.userDAO.Get(db, &model.User{ID: userID}, []string{"id", "name", "email", "avatar", "created_at", "last_login", "permission"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[UserService] user not found")
			rsp.Error = constant.ErrDataNotExists
			return rsp, nil
		}
		logger.Error("[UserService] failed to get user by id", zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	rsp.User = &dto.DetailedUser{
		User: dto.User{
			ID:     user.ID,
			Name:   user.Name,
			Avatar: user.Avatar,
		},
		CreatedAt:  user.CreatedAt.Format(time.DateTime),
		LastLogin:  user.LastLogin.Format(time.DateTime),
		Permission: string(user.Permission),
	}

	logger.Info("[UserService] get cur user info",
		zap.String("email", user.Email),
		zap.String("avatar", user.Avatar),
		zap.Time("createdAt", user.CreatedAt),
		zap.Time("lastLogin", user.LastLogin),
		zap.String("permission", string(user.Permission)))

	return rsp, nil
}

func (s *userService) UpdateUser(ctx context.Context, req *dto.UpdateUserReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	if err := s.userDAO.Update(db, &model.User{ID: userID}, map[string]interface{}{
		"name":   req.Body.User.Name,
		"email":  req.Body.User.Email,
		"avatar": req.Body.User.Avatar,
	}); err != nil {
		logger.Error("[UserService] failed to update user", zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	return rsp, nil
}

// ApproveUser approves a pending user and promotes them to user permission
//
//	@receiver s *userService
//	@param ctx context.Context
//	@param req *dto.ApproveUserReq
//	@return *dto.EmptyRsp
//	@return error
//	@author centonhuang
//	@update 2026-02-02 10:00:00
func (s *userService) ApproveUser(ctx context.Context, req *dto.ApproveUserReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

	if req == nil || req.Body == nil || req.Body.UserID == 0 {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	user, err := s.userDAO.Get(db, &model.User{ID: req.Body.UserID}, []string{"id", "permission"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[UserService] user not found")
			rsp.Error = constant.ErrDataNotExists
			return rsp, nil
		}
		logger.Error("[UserService] failed to get user by id", zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	if user.Permission != enum.PermissionPending {
		logger.Info("[UserService] user is not pending", zap.Uint("userID", req.Body.UserID), zap.String("permission", string(user.Permission)))
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	if err := s.userDAO.Update(db, &model.User{ID: req.Body.UserID}, map[string]interface{}{
		"permission": enum.PermissionUser,
	}); err != nil {
		logger.Error("[UserService] failed to approve user", zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	return rsp, nil
}

// ListUsers lists users with pagination
//
//	@receiver s *userService
//	@param ctx context.Context
//	@param req *dto.ListUsersReq
//	@return *dto.ListUsersRsp
//	@return error
//	@author centonhuang
//	@update 2026-02-02 10:20:00
func (s *userService) ListUsers(ctx context.Context, req *dto.ListUsersReq) (*dto.ListUsersRsp, error) {
	rsp := &dto.ListUsersRsp{}

	if req.SortField == "" {
		req.SortField = "id"
	}

	if req.Sort == "" {
		req.Sort = enum.SortAsc
	}

	db := database.GetDBInstance(ctx)
	logger := logger.WithCtx(ctx)

	commonParam := &dao.CommonParam{
		PageParam: dao.PageParam{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		QueryParam: dao.QueryParam{
			Query:       req.Query,
			QueryFields: []string{"name", "email"},
		},
		SortParam: dao.SortParam{
			Sort:      req.Sort,
			SortField: strcase.ToSnake(req.SortField),
		},
	}

	users, pageInfo, err := s.userDAO.Paginate(db, &model.User{}, []string{"id", "name", "email", "avatar", "permission", "created_at", "last_login"}, commonParam)
	if err != nil {
		logger.Error("[UserService] failed to list users", zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	rsp.Users = lo.Map(users, func(item *model.User, _ int) *dto.DetailedUser {
		return &dto.DetailedUser{
			User: dto.User{
				ID:     item.ID,
				Name:   item.Name,
				Avatar: item.Avatar,
			},
			CreatedAt:  item.CreatedAt.Format(time.DateTime),
			LastLogin:  item.LastLogin.Format(time.DateTime),
			Permission: string(item.Permission),
		}
	})
	rsp.PageInfo = pageInfo
	return rsp, nil
}
