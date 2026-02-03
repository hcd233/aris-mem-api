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

// GetCosTempCredentialRsp 获取COS临时密钥响应
//
//	@author centonhuang
//	@update 2026-01-31 18:00:00
type GetCosTempCredentialRsp struct {
	CommonRsp
	CosTempCredential *CosTempCredential `json:"cosTempCredential" doc:"COS临时密钥信息"`
}

// CosTempCredential COS临时密钥信息
//
//	@author centonhuang
//	@update 2026-01-31 18:00:00
type CosTempCredential struct {
	// SecretID string SecretId 临时密钥 SecretId
	//	@update 2026-01-31 06:14:03
	SecretID string `json:"secretId" doc:"临时密钥 SecretId"`

	// SecretKey string SecretKey 临时密钥 SecretKey
	//	@update 2026-01-31 06:14:08
	SecretKey string `json:"secretKey" doc:"临时密钥 SecretKey"`

	// SessionToken string 临时密钥 SessionToken
	//	@update 2026-01-31 06:14:10
	SessionToken string `json:"sessionToken" doc:"临时密钥 SessionToken"`

	// AppID string 应用ID
	//	@update 2026-01-31 06:14:25
	AppID string `json:"appId" doc:"应用ID"`

	// ExpiredTime int64 临时密钥过期时间戳(秒)
	//	@update 2026-01-31 06:14:13
	ExpiredTime int64 `json:"expiredTime" doc:"临时密钥过期时间戳(秒)"`
	// Expiration string ISO8601格式的过期时间
	//	@update 2026-01-31 06:14:15
	Expiration string `json:"expiration" doc:"ISO8601格式的过期时间"`
	// StartTime int64 临时密钥开始时间戳(秒)
	//	@update 2026-01-31 06:14:17
	StartTime int64 `json:"startTime" doc:"临时密钥开始时间戳(秒)"`
	// RequestID string 请求ID
	//	@update 2026-01-31 06:14:19
	RequestID string `json:"requestId" doc:"请求ID"`
	// BucketName string Bucket名称
	//	@update 2026-01-31 06:14:22
	BucketName string `json:"bucketName" doc:"Bucket名称"`
	// Region string 地域
	//	@update 2026-01-31 06:13:52
	Region string `json:"region" doc:"地域"`
}
