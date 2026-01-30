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
//	@update 2026-01-31 18:00:00
type ImageHandler interface {
	HandleUploadImage(ctx context.Context, req *dto.UploadImageReq) (*dto.HTTPResponse[*dto.UploadImageRsp], error)
	HandleGetCredential(ctx context.Context, req *dto.EmptyReq) (*dto.HTTPResponse[*dto.GetCosTempCredentialRsp], error)
}

type imageHandler struct {
	svc service.ImageService
}

// NewImageHandler 创建图片处理器
//
//	return ImageHandler
//	author centonhuang
//	@update 2026-01-31 18:00:00
func NewImageHandler() ImageHandler {
	return &imageHandler{
		svc: service.NewImageService(),
	}
}

func (h *imageHandler) HandleUploadImage(ctx context.Context, req *dto.UploadImageReq) (*dto.HTTPResponse[*dto.UploadImageRsp], error) {
	return util.WrapHTTPResponse(h.svc.UploadImage(ctx, req))
}

// HandleGetCredential 处理获取COS临时密钥请求
//
//	@param ctx context.Context
//	@param req *dto.GetCosTempCredentialReq
//	@return *dto.HTTPResponse[*dto.GetCosTempCredentialRsp]
//	@return error
//	author centonhuang
//	@update 2026-01-31 18:00:00
func (h *imageHandler) HandleGetCredential(ctx context.Context, req *dto.EmptyReq) (*dto.HTTPResponse[*dto.GetCosTempCredentialRsp], error) {
	return util.WrapHTTPResponse(h.svc.GetCosTempCredential(ctx, req))
}
