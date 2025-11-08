package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/handler"
	"github.com/hcd233/aris-mem-api/internal/middleware"
)

func initAgentRouter(agentGroup huma.API) {
	agentHandler := handler.NewAgentHandler()

	agentGroup.UseMiddleware(middleware.JwtMiddleware())

	// sse.Register(agentGroup, huma.Operation{
	// 	OperationID: "chat",
	// 	Method:      http.MethodPost,
	// 	Path:        "/chat",
	// 	Summary:     "Chat",
	// 	Description: "Chat with Aris Mem Agent",
	// 	Tags:        []string{"agent"},
	// 	Security: []map[string][]string{
	// 		{"jwtAuth": {}},
	// 	},
	// }, map[string]any{
	// 	"SSEResponse": protocol.SSEResponse{},
	// }, agentHandler.HandleChat)

	huma.Register(agentGroup, huma.Operation{
		OperationID: "chat",
		Method:      http.MethodPost,
		Path:        "/chat",
		Summary:     "Chat",
		Description: "Chat with Aris Mem Agent",
		Tags:        []string{"agent"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, agentHandler.HandleChat)
}
