package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/handler"
	"github.com/hcd233/aris-mem-api/internal/middleware"
)

func initTodoItemRouter(todoItemGroup *huma.Group) {
	todoItemHandler := handler.NewTodoItemHandler()

	todoItemGroup.UseMiddleware(middleware.JwtMiddleware())

	// 创建待办事项
	huma.Register(todoItemGroup, huma.Operation{
		OperationID: "createTodoItems",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "CreateTodoItems",
		Description: "Create todo items",
		Tags:        []string{"todoItem"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, todoItemHandler.HandleCreateTodoItems)

	// 获取待办事项列表
	huma.Register(todoItemGroup, huma.Operation{
		OperationID: "listTodoItems",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "ListTodoItems",
		Description: "List todo items",
		Tags:        []string{"todoItem"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, todoItemHandler.HandleListTodoItems)
}
