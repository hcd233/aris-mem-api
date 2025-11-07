package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/handler"
)

// initHealthRouter 初始化健康检查路由
//
//	@param healthGroup
//	@author centonhuang
//	@update 2025-11-07 14:59:06
func initHealthRouter(healthGroup huma.API) {
	pingHandler := handler.NewPingHandler()

	huma.Register(healthGroup, huma.Operation{
		OperationID: "healthCheck",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "HealthCheck",
		Description: "Check the server health",
		Tags:        []string{"health"},
	}, pingHandler.HandlePing)
}
