package handler

import (
	"context"

	"github.com/hcd233/aris-mem-api/internal/protocol"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/hcd233/aris-mem-api/internal/service"
	"github.com/hcd233/aris-mem-api/internal/util"
)

// Oauth2Handler OAuth2处理器接口
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type Oauth2Handler interface {
	HandleLogin(ctx context.Context, req *dto.LoginReq) (*protocol.HTTPResponse[*dto.LoginResp], error)
	HandleCallback(ctx context.Context, req *dto.CallbackReq) (*protocol.HTTPResponse[*dto.CallbackRsp], error)
}

type oauth2Handler struct{}

// NewOauth2Handler 创建OAuth2处理器
//
//	return Oauth2Handler
//	author centonhuang
//	update 2025-01-05 21:00:00
func NewOauth2Handler() Oauth2Handler {
	return &oauth2Handler{}
}

// HandleLogin OAuth2登录
//
//	receiver h *oauth2Handler
//	param ctx context.Context
//	param req *dto.LoginRequest
//	return *protocol.HumaHTTPResponse[*dto.LoginResp]
//	return error
//	author centonhuang
//	update 2025-01-05 21:00:00
func (h *oauth2Handler) HandleLogin(ctx context.Context, req *dto.LoginReq) (*protocol.HTTPResponse[*dto.LoginResp], error) {
	svc := h.getService(req.Provider)
	return util.WrapHTTPResponse(svc.Login(ctx, req))
}

// HandleCallback OAuth2回调
//
//	receiver h *oauth2Handler
//	param ctx context.Context
//	param req *dto.CallbackRequest
//	return *protocol.HumaHTTPResponse[*dto.CallbackRsp]
//	return error
//	author centonhuang
//	update 2025-01-05 21:00:00
func (h *oauth2Handler) HandleCallback(ctx context.Context, req *dto.CallbackReq) (*protocol.HTTPResponse[*dto.CallbackRsp], error) {
	svc := h.getService(req.Provider)
	return util.WrapHTTPResponse(svc.Callback(ctx, req))
}

// getService 根据provider获取对应的service
//
//	receiver h *oauth2Handler
//	param provider string
//	return service.Oauth2Service
//	author centonhuang
//	update 2025-01-05 21:00:00
func (h *oauth2Handler) getService(provider string) service.Oauth2Service {
	switch provider {
	case "github":
		return service.NewGithubOauth2Service()
	case "google":
		return service.NewGoogleOauth2Service()
	// case "qq":
	// 	return service.NewQQOauth2Service()
	default:
		return service.NewGithubOauth2Service() // 默认返回 github
	}
}
