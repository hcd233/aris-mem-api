package dto

// ChatReq 聊天请求
//
//	@author centonhuang
//	@update 2025-11-08 04:20:42
type ChatReq struct {
	Body *ChatReqBody `json:"body" doc:"Body"`
}

// ChatReqBody 聊天请求体
//
//	@author centonhuang
//	@update 2025-11-08 04:54:55
type ChatReqBody struct {
	Message string `json:"message" doc:"Message to chat"`
}

// ChatRsp 聊天响应
//
//	@author centonhuang
//	@update 2025-11-08 04:20:42
type ChatRsp struct {
	CommonRsp
	Message string `json:"message" doc:"Message to chat"`
}
