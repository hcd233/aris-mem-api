package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-mem-api/internal/handler"
	"github.com/hcd233/aris-mem-api/internal/middleware"
)

// initAgentRouter 初始化Agent路由
//
//	@param agentGroup
//	@author centonhuang
//	@update 2025-11-10 18:39:15
func initAgentRouter(agentGroup fiber.Router) {
	agentHandler := handler.NewAgentHandler()

	agentGroup.Use(middleware.JwtMiddlewareFiber())

	agentGroup.Post("chat",
		// middleware.RedisLockMiddleware("agentChat", constant.CtxKeyUserID, constant.AgentChatLockExpire),
		agentHandler.HandleChat)
}
