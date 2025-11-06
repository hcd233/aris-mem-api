// Package router 路由
package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/api"
	"github.com/hcd233/aris-mem-api/internal/handler"
)

// RegisterRouter 注册路由
//
//	param app *fiber.App
//	author centonhuang
//	update 2025-01-04 15:32:40
func RegisterRouter() {
	pingService := handler.NewPingHandler()

	api := api.GetHumaAPI()

	apiGroup := huma.NewGroup(api, "/api")

	v1Group := huma.NewGroup(apiGroup, "/v1")

	userGroup := huma.NewGroup(v1Group, "/user")
	initUserRouter(userGroup)

	tokenGroup := huma.NewGroup(v1Group, "/token")
	initTokenRouter(tokenGroup)

	oauth2Group := huma.NewGroup(v1Group, "/oauth2")
	initOauth2Router(oauth2Group)

	huma.Register(api, huma.Operation{
		OperationID: "ping",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Ping",
		Description: "Check service if available.",
		Tags:        []string{"ping"},
	}, pingService.HandlePing)
}
