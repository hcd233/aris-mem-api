package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/handler"
	"github.com/hcd233/aris-mem-api/internal/middleware"
)

func initUserRouter(userGroup huma.API) {
	userHandler := handler.NewUserHandler()

	userGroup.UseMiddleware(middleware.JwtMiddlewareHuma())

	// 获取当前用户信息
	huma.Register(userGroup, huma.Operation{
		OperationID: "getCurrentUser",
		Method:      http.MethodGet,
		Path:        "/current",
		Summary:     "GetCurrentUser",
		Description: "Get the current user's detailed information, including user ID, username, email, avatar, and permission information",
		Tags:        []string{"user"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, userHandler.HandleGetCurUser)

	// 更新用户信息
	huma.Register(userGroup, huma.Operation{
		OperationID: "updateUser",
		Method:      http.MethodPatch,
		Path:        "/",
		Summary:     "UpdateUser",
		Description: "Update the current user's information, including the username and other fields",
		Tags:        []string{"user"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, userHandler.HandleUpdateUser)
}
