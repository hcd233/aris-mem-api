package service

import (
	"context"
	"errors"
	"time"

	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/hcd233/aris-mem-api/internal/resource/database"
	"github.com/hcd233/aris-mem-api/internal/resource/database/dao"
	"github.com/hcd233/aris-mem-api/internal/resource/database/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UserService 用户服务
//
//	author centonhuang
//	update 2025-01-04 21:04:00
type UserService interface {
	GetCurUserInfo(ctx context.Context, req *dto.EmptyReq) (rsp *dto.GetCurUserInfoResp, err error)
	UpdateUserInfo(ctx context.Context, req *dto.UpdateUserInfoReq) (rsp *dto.EmptyResp, err error)
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

// GetCurUserInfo 获取当前用户信息
//
//	receiver s *userService
//	param ctx context.Context
//	param req *protocol.GetCurUserInfoRequest
//	return rsp *protocol.GetCurUserInfoResponse
//	return err error
//	author centonhuang
//	update 2025-01-04 21:04:03
func (s *userService) GetCurUserInfo(ctx context.Context, _ *dto.EmptyReq) (rsp *dto.GetCurUserInfoResp, err error) {
	rsp = &dto.GetCurUserInfoResp{}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	user, err := s.userDAO.GetByID(db, userID, []string{"id", "name", "email", "avatar", "created_at", "last_login", "permission"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[UserService] user not found")
			return nil, constant.ErrDataNotExists
		}
		logger.Error("[UserService] failed to get user by id", zap.Error(err))
		return nil, constant.ErrInternalError
	}

	rsp.User = &dto.DetailedUser{
		ID:         user.ID,
		CreatedAt:  user.CreatedAt.Format(time.DateTime),
		LastLogin:  user.LastLogin.Format(time.DateTime),
		Permission: string(user.Permission),
		User: dto.User{
			Name:   user.Name,
			Email:  user.Email,
			Avatar: user.Avatar,
		},
	}

	logger.Info("[UserService] get cur user info",
		zap.String("email", user.Email),
		zap.String("avatar", user.Avatar),
		zap.Time("createdAt", user.CreatedAt),
		zap.Time("lastLogin", user.LastLogin),
		zap.String("permission", string(user.Permission)))

	return rsp, nil
}

func (s *userService) UpdateUserInfo(ctx context.Context, req *dto.UpdateUserInfoReq) (rsp *dto.EmptyResp, err error) {
	rsp = &dto.EmptyResp{}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	logger := logger.WithCtx(ctx)

	db := database.GetDBInstance(ctx)

	if err := s.userDAO.Update(db, &model.User{ID: userID}, map[string]interface{}{
		"name":   req.Body.User.Name,
		"email":  req.Body.User.Email,
		"avatar": req.Body.User.Avatar,
	}); err != nil {
		logger.Error("[UserService] failed to update user", zap.Error(err))
		return nil, constant.ErrInternalError
	}

	return rsp, nil
}
