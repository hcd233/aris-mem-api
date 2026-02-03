package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/handler"
	"github.com/hcd233/aris-mem-api/internal/middleware"
)

func initCommentRouter(commentGroup huma.API) {
	commentHandler := handler.NewCommentHandler()

	commentGroup.UseMiddleware(middleware.JwtMiddleware(),
		middleware.LimitUserPermissionMiddleware("comment", enum.PermissionUser))

	huma.Register(commentGroup, huma.Operation{
		OperationID: "createComment",
		Method:      http.MethodPost,
		Path:        "",
		Summary:     "CreateComment",
		Description: "Create a new comment on an article",
		Tags:        []string{"Comment"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, commentHandler.HandleCreateComment)

	huma.Register(commentGroup, huma.Operation{
		OperationID: "listComments",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "ListComments",
		Description: "List comments with pagination and filtering by article ID or parent comment ID",
		Tags:        []string{"Comment"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, commentHandler.HandleListComments)

	huma.Register(commentGroup, huma.Operation{
		OperationID: "deleteComment",
		Method:      http.MethodDelete,
		Path:        "",
		Summary:     "DeleteComment",
		Description: "Delete a comment (only the owner can delete)",
		Tags:        []string{"Comment"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, commentHandler.HandleDeleteComment)
}
