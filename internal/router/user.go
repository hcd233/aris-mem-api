package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/handler"
	"github.com/hcd233/aris-mem-api/internal/middleware"
)

func initUserRouter(userGroup huma.API) {
	userHandler := handler.NewUserHandler()

	userGroup.UseMiddleware(middleware.JwtMiddleware())

	huma.Register(userGroup, huma.Operation{
		OperationID: "getCurrentUser",
		Method:      http.MethodGet,
		Path:        "/current",
		Summary:     "GetCurrentUser",
		Description: "Get the current user's detailed information, including user ID, username, email, avatar, and permission information",
		Tags:        []string{"User"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, userHandler.HandleGetCurUser)

	huma.Register(userGroup, huma.Operation{
		OperationID: "updateUser",
		Method:      http.MethodPatch,
		Path:        "",
		Summary:     "UpdateUser",
		Description: "Update the current user's information, including the username and other fields",
		Tags:        []string{"User"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
		Middlewares: huma.Middlewares{middleware.LimitUserPermissionMiddleware("updateUser", enum.PermissionUser)},
	}, userHandler.HandleUpdateUser)

	huma.Register(userGroup, huma.Operation{
		OperationID: "approveUser",
		Method:      http.MethodPost,
		Path:        "/approve",
		Summary:     "ApproveUser",
		Description: "Approve a pending user and promote them to user permission",
		Tags:        []string{"User"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
		Middlewares: huma.Middlewares{middleware.LimitUserPermissionMiddleware("approveUser", enum.PermissionAdmin)},
	}, userHandler.HandleApproveUser)

	huma.Register(userGroup, huma.Operation{
		OperationID: "listUsers",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "ListUsers",
		Description: "List users with pagination",
		Tags:        []string{"User"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
		Middlewares: huma.Middlewares{middleware.LimitUserPermissionMiddleware("listUsers", enum.PermissionAdmin)},
	}, userHandler.HandleListUsers)
}
