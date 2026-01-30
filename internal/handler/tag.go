package handler

import (
	"context"

	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/service"
	"github.com/hcd233/aris-mem-api/internal/util"
)

// TagHandler 标签处理器
//
//	author centonhuang
//	@update 2026-01-29 10:00:00
type TagHandler interface {
	HandleListTags(ctx context.Context, req *dto.ListTagsReq) (*dto.HTTPResponse[*dto.ListTagsRsp], error)
	HandleDeleteTag(ctx context.Context, req *dto.DeleteTagReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
}

type tagHandler struct {
	svc service.TagService
}

// NewTagHandler 创建标签处理器
//
//	return TagHandler
//	author centonhuang
//	update 2026-01-29 10:00:00
func NewTagHandler() TagHandler {
	return &tagHandler{
		svc: service.NewTagService(),
	}
}

func (h *tagHandler) HandleListTags(ctx context.Context, req *dto.ListTagsReq) (*dto.HTTPResponse[*dto.ListTagsRsp], error) {
	return util.WrapHTTPResponse(h.svc.ListTags(ctx, req))
}

func (h *tagHandler) HandleDeleteTag(ctx context.Context, req *dto.DeleteTagReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
	return util.WrapHTTPResponse(h.svc.DeleteTag(ctx, req))
}
