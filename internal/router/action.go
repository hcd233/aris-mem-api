package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/handler"
	"github.com/hcd233/aris-mem-api/internal/middleware"
)

func initActionRouter(actionGroup huma.API) {
	actionHandler := handler.NewActionHandler()

	actionGroup.UseMiddleware(middleware.JwtMiddleware(),
		middleware.LimitUserPermissionMiddleware("action", enum.PermissionUser))

	huma.Register(actionGroup, huma.Operation{
		OperationID: "doAction",
		Method:      http.MethodPost,
		Path:        "/do",
		Summary:     "DoAction",
		Description: "Perform an action on an entity (like/save article)",
		Tags:        []string{"Action"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
		Middlewares: huma.Middlewares{middleware.RateLimiterMiddleware("actionDo", "", constant.PeriodActionDo, constant.LimitActionDo)},
	}, actionHandler.HandleDo)

	huma.Register(actionGroup, huma.Operation{
		OperationID: "undoAction",
		Method:      http.MethodPost,
		Path:        "/undo",
		Summary:     "UndoAction",
		Description: "Undo an action on an entity (undo like/undo save article)",
		Tags:        []string{"Action"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
		Middlewares: huma.Middlewares{middleware.RateLimiterMiddleware("actionUndo", "", constant.PeriodActionUndo, constant.LimitActionUndo)},
	}, actionHandler.HandleUndo)
}
