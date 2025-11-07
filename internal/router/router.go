// Package router 路由
package router

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/api"
)

// RegisterRouter 注册路由
//
//	param app *fiber.App
//	author centonhuang
//	update 2025-01-04 15:32:40
func RegisterRouter() {
	api := api.GetHumaAPI()

	apiGroup := huma.NewGroup(api, "/api")

	v1Group := huma.NewGroup(apiGroup, "/v1")

	initHealthRouter(api)

	userGroup := huma.NewGroup(v1Group, "/user")
	initUserRouter(userGroup)

	todoItemGroup := huma.NewGroup(v1Group, "/todoItem")
	initTodoItemRouter(todoItemGroup)

	tokenGroup := huma.NewGroup(v1Group, "/token")
	initTokenRouter(tokenGroup)

	oauth2Group := huma.NewGroup(v1Group, "/oauth2")
	initOauth2Router(oauth2Group)
}
