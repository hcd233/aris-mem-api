package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/handler"
	"github.com/hcd233/aris-mem-api/internal/middleware"
)

// initAgentRouter 初始化Agent路由
//
//	@param agentGroup
//	@author centonhuang
//	@update 2025-11-10 18:39:15
func initAgentRouter(agentGroup huma.API) {
	agentHandler := handler.NewAgentHandler()

	agentGroup.UseMiddleware(middleware.JwtMiddleware(),
		middleware.LimitUserPermissionMiddleware("agent", enum.PermissionUser))

	huma.Register(agentGroup, huma.Operation{
		OperationID: "chat",
		Method:      http.MethodPost,
		Path:        "/chat",
		Summary:     "Chat",
		Description: "Chat with the agent",
		Tags:        []string{"Agent"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, agentHandler.HandleChat)
}
