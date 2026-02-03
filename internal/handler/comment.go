package handler

import (
	"context"

	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/service"
	"github.com/hcd233/aris-mem-api/internal/util"
)

// CommentHandler comment handler
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type CommentHandler interface {
	HandleCreateComment(ctx context.Context, req *dto.CreateCommentReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
	HandleListComments(ctx context.Context, req *dto.ListCommentsReq) (*dto.HTTPResponse[*dto.ListCommentsRsp], error)
	HandleDeleteComment(ctx context.Context, req *dto.DeleteCommentReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
}

type commentHandler struct {
	svc service.CommentService
}

// NewCommentHandler create comment handler
//
//	return CommentHandler
//	author centonhuang
//	update 2026-02-03 22:30:00
func NewCommentHandler() CommentHandler {
	return &commentHandler{
		svc: service.NewCommentService(),
	}
}

func (h *commentHandler) HandleCreateComment(ctx context.Context, req *dto.CreateCommentReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
	return util.WrapHTTPResponse(h.svc.CreateComment(ctx, req))
}

func (h *commentHandler) HandleListComments(ctx context.Context, req *dto.ListCommentsReq) (*dto.HTTPResponse[*dto.ListCommentsRsp], error) {
	return util.WrapHTTPResponse(h.svc.ListComments(ctx, req))
}

func (h *commentHandler) HandleDeleteComment(ctx context.Context, req *dto.DeleteCommentReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
	return util.WrapHTTPResponse(h.svc.DeleteComment(ctx, req))
}
