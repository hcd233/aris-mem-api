package dto

import "mime/multipart"

// UploadImageReq 上传图片请求
//
//	@author centonhuang
//	@update 2026-01-31 14:00:00
type UploadImageReq struct {
	RawBody *multipart.FileHeader
}

// UploadImageRsp 上传图片响应
//
//	@author centonhuang
//	@update 2026-01-31 14:00:00
type UploadImageRsp struct {
	CommonRsp
	ImageName string `json:"imageName" doc:"Name of the uploaded image"`
}
