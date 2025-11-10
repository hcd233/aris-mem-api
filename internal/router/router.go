// Package router 路由
package router

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-mem-api/internal/api"
)

// RegisterDocsRouter 注册文档路由
//
//	@author centonhuang
//	@update 2025-11-10 17:26:08
func RegisterDocsRouter() {
	app := api.GetFiberApp()
	app.Get("/docs", func(c *fiber.Ctx) error {
		html := `<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/openapi.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`
		return c.Type("html").SendString(html)
	})
}

// RegisterAPIRouter 注册API路由
//
//	@author centonhuang
//	@update 2025-11-10 17:26:08
func RegisterAPIRouter() {
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

// RegisterSSERouter 注册SSE路由
//
//	@author centonhuang
//	@update 2025-11-10 18:39:27
func RegisterSSERouter() {
	app := api.GetFiberApp()

	sseGroup := app.Group("/sse")
	v1Group := sseGroup.Group("/v1")

	initSSEHealthRouter(app)

	agentGroup := v1Group.Group("/agent")
	initAgentRouter(agentGroup)
}
