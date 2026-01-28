package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/handler"
	"github.com/hcd233/aris-mem-api/internal/middleware"
)

func initTagRouter(tagGroup huma.API) {
	tagHandler := handler.NewTagHandler()

	tagGroup.UseMiddleware(middleware.JwtMiddleware(),
		middleware.LimitUserPermissionMiddleware("tag", enum.PermissionUser))

	huma.Register(tagGroup, huma.Operation{
		OperationID: "listTags",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "ListTags",
		Description: "List tags with pagination and fuzzy search",
		Tags:        []string{"Tag"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, tagHandler.HandleListTags)
}
