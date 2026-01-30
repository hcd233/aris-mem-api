package handler

import (
	"context"

	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/service"
	"github.com/hcd233/aris-mem-api/internal/util"
)

// ActionHandler action handler
//
//	author centonhuang
//	update 2026-01-30 21:00:00
type ActionHandler interface {
	HandleDo(ctx context.Context, req *dto.ActionReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
	HandleUndo(ctx context.Context, req *dto.ActionReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
}

type actionHandler struct {
	svc service.ActionService
}

// NewActionHandler create action handler
//
//	return ActionHandler
//	author centonhuang
//	update 2026-01-30 21:00:00
func NewActionHandler() ActionHandler {
	return &actionHandler{
		svc: service.NewActionService(),
	}
}

func (h *actionHandler) HandleDo(ctx context.Context, req *dto.ActionReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
	return util.WrapHTTPResponse(h.svc.Do(ctx, req))
}

func (h *actionHandler) HandleUndo(ctx context.Context, req *dto.ActionReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
	return util.WrapHTTPResponse(h.svc.Undo(ctx, req))
}
