package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/hcd233/aris-mem-api/internal/service"
)

// PingHandler 健康检查处理器
//
//	author centonhuang
//	update 2025-01-04 15:52:48
type AgentHandler interface {
	HandleChat(ctx context.Context, req *dto.ChatReq, sender sse.Sender)
}

type agentHandler struct {
	agentService service.AgentService
}

// NewAgentHandler 创建健康检查处理器
//
//	update 2025-01-04 15:52:48
//	@return AgentHandler
//	@author centonhuang
//	@update 2025-11-08 04:42:42
func NewAgentHandler() AgentHandler {
	return &agentHandler{
		agentService: service.NewAgentService(),
	}
}

// HandleChat 聊天处理器
//
//	@receiver h *agentHandler
//	@param _ context.Context
//	@param req *dto.ChatReq
//	@param sender sse.Sender
//	@author centonhuang
//	@update 2025-11-08 04:42:12
func (h *agentHandler) HandleChat(ctx context.Context, req *dto.ChatReq, sender sse.Sender) {
	h.agentService.HandleChat(ctx, req, sender)
}
