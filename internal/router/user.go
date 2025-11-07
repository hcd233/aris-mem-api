package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/handler"
	"github.com/hcd233/aris-mem-api/internal/middleware"
)

func initUserRouter(userGroup huma.API) {
	userHandler := handler.NewUserHandler()

	userGroup.UseMiddleware(middleware.JwtMiddleware())

	// 获取当前用户信息
	huma.Register(userGroup, huma.Operation{
		OperationID: "getCurrentUserInfo",
		Method:      http.MethodGet,
		Path:        "/current",
		Summary:     "GetCurrentUserInfo",
		Description: "Get the current user's detailed information, including user ID, username, email, avatar, and permission information",
		Tags:        []string{"user"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, userHandler.HandleGetCurUserInfo)

	// 更新用户信息
	huma.Register(userGroup, huma.Operation{
		OperationID: "updateUserInfo",
		Method:      http.MethodPatch,
		Path:        "/",
		Summary:     "UpdateUserInfo",
		Description: "Update the current user's information, including the username and other fields",
		Tags:        []string{"user"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, userHandler.HandleUpdateInfo)
}
