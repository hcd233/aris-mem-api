package dto

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

// UploadImageReq 上传图片请求
//
//	@author centonhuang
//	@update 2026-01-31 14:00:00
type UploadImageReq struct {
	RawBody huma.MultipartFormFiles[UploadImageReqBody]
}

// UploadImageReqBody 上传图片请求体
//
//	@author centonhuang
//	@update 2026-01-31 02:05:17
type UploadImageReqBody struct {
	Image huma.FormFile `form:"image" contentType:"image/gif,image/jpeg,image/png"`
}

// UploadImageRsp 上传图片响应
//
//	@author centonhuang
//	@update 2026-01-31 14:00:00
type UploadImageRsp struct {
	CommonRsp
	ImageName string `json:"imageName" doc:"Name of the uploaded image"`
}

// ImageUploadTask 图片上传任务
//
//	@author centonhuang
//	@update 2026-01-31 16:00:00
type ImageUploadTask struct {
	Ctx       context.Context
	ImageName string
	ImageData []byte
}
