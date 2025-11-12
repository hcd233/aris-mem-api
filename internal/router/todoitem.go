package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/handler"
	"github.com/hcd233/aris-mem-api/internal/middleware"
)

func initTodoItemRouter(todoItemGroup huma.API) {
	todoItemHandler := handler.NewTodoItemHandler()

	todoItemGroup.UseMiddleware(middleware.JwtMiddleware(),
		middleware.LimitUserPermissionMiddleware("todoItem", enum.PermissionUser))

	huma.Register(todoItemGroup, huma.Operation{
		OperationID: "createTodoItems",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "CreateTodoItems",
		Description: "Create todo items",
		Tags:        []string{"TodoItem"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, todoItemHandler.HandleCreateTodoItems)

	huma.Register(todoItemGroup, huma.Operation{
		OperationID: "listTodoItems",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "ListTodoItems",
		Description: "List todo items",
		Tags:        []string{"TodoItem"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, todoItemHandler.HandleListTodoItems)
}
