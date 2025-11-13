package dto

import "github.com/danielgtaylor/huma/v2"

// ChatReq 聊天请求
//
//	@author centonhuang
//	@update 2025-11-08 04:20:42
type ChatReq struct {
	RawBody huma.MultipartFormFiles[ChatReqBody]
}

// ChatReqBody 聊天请求体
//
//	@author centonhuang
//	@update 2025-11-08 04:54:55
type ChatReqBody struct {
	Content string          `form:"content" doc:"Content to chat"`
	Audio   huma.FormFile   `form:"audio" contentType:"audio/wav" doc:"Audio file"`
	Images  []huma.FormFile `form:"images" contentType:"image/png,image/jpeg" doc:"Images"`
}
