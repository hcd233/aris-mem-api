package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/hcd233/aris-mem-api/internal/service"
	"github.com/hcd233/aris-mem-api/internal/util"
	"go.uber.org/zap"
)

// PingHandler 健康检查处理器
//
//	author centonhuang
//	update 2025-01-04 15:52:48
type AgentHandler interface {
	HandleChat(ctx *fiber.Ctx) error
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
func (h *agentHandler) HandleChat(ctx *fiber.Ctx) error {
	body := &dto.ChatReqBody{}
	if err := ctx.BodyParser(body); err != nil {
		logger.WithFCtx(ctx).Error("[AgentHandler] failed to parse body", zap.Error(err))
		return util.WriteErrorResponse(ctx.Response().BodyWriter(), constant.ErrBadRequest)
	}
	req := &dto.ChatReq{
		Body: body,
	}
	ch, err := h.agentService.HandleChat(ctx.Context(), req)
	return util.WrapSSEResponse(ctx, ch, err)
}
