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
	HandleGetCurUserInfo(ctx context.Context, req *dto.EmptyReq) (*protocol.HTTPResponse[*dto.GetCurUserInfoResp], error)
	HandleUpdateInfo(ctx context.Context, req *dto.UpdateUserInfoReq) (*protocol.HTTPResponse[*dto.EmptyResp], error)
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

func (h *userHandler) HandleGetCurUserInfo(ctx context.Context, req *dto.EmptyReq) (*protocol.HTTPResponse[*dto.GetCurUserInfoResp], error) {
	return util.WrapHTTPResponse(h.svc.GetCurUserInfo(ctx, req))
}

func (h *userHandler) HandleUpdateInfo(ctx context.Context, req *dto.UpdateUserInfoReq) (*protocol.HTTPResponse[*dto.EmptyResp], error) {
	return util.WrapHTTPResponse(h.svc.UpdateUserInfo(ctx, req))
}
