package handler

import (
	"context"

	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/service"
	"github.com/hcd233/aris-mem-api/internal/util"
)

// ImageHandler 图片处理器
//
//	author centonhuang
//	update 2026-01-31 14:00:00
type ImageHandler interface {
	HandleUploadImage(ctx context.Context, req *dto.UploadImageReq) (*dto.HTTPResponse[*dto.UploadImageRsp], error)
}

type imageHandler struct {
	svc service.ImageService
}

// NewImageHandler 创建图片处理器
//
//	return ImageHandler
//	author centonhuang
//	update 2026-01-31 14:00:00
func NewImageHandler() ImageHandler {
	return &imageHandler{
		svc: service.NewImageService(),
	}
}

func (h *imageHandler) HandleUploadImage(ctx context.Context, req *dto.UploadImageReq) (*dto.HTTPResponse[*dto.UploadImageRsp], error) {
	return util.WrapHTTPResponse(h.svc.UploadImage(ctx, req))
}
