package handler

import (
	"context"

	"github.com/hcd233/aris-mem-api/internal/protocol"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/hcd233/aris-mem-api/internal/service"
	"github.com/hcd233/aris-mem-api/internal/util"
)

// UserHandler 用户处理器
//
//	author centonhuang
//	update 2025-01-04 15:56:20
type UserHandler interface {
	HandleGetCurUser(ctx context.Context, req *dto.EmptyReq) (*protocol.HTTPResponse[*dto.GetCurUserRsp], error)
	HandleUpdateUser(ctx context.Context, req *dto.UpdateUserReq) (*protocol.HTTPResponse[*dto.EmptyRsp], error)
}

type userHandler struct {
	svc service.UserService
}

// NewUserHandler 创建用户处理器
//
//	return UserHandler
//	author centonhuang
//	update 2024-12-08 16:59:38
func NewUserHandler() UserHandler {
	return &userHandler{
		svc: service.NewUserService(),
	}
}

func (h *userHandler) HandleGetCurUser(ctx context.Context, req *dto.EmptyReq) (*protocol.HTTPResponse[*dto.GetCurUserRsp], error) {
	return util.WrapHTTPResponse(h.svc.GetCurUser(ctx, req))
}

func (h *userHandler) HandleUpdateUser(ctx context.Context, req *dto.UpdateUserReq) (*protocol.HTTPResponse[*dto.EmptyRsp], error) {
	return util.WrapHTTPResponse(h.svc.UpdateUser(ctx, req))
}
