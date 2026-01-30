package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/handler"
	"github.com/hcd233/aris-mem-api/internal/middleware"
)

func initArticleRouter(articleGroup huma.API) {
	articleHandler := handler.NewArticleHandler()

	articleGroup.UseMiddleware(middleware.JwtMiddleware(),
		middleware.LimitUserPermissionMiddleware("article", enum.PermissionUser))

	huma.Register(articleGroup, huma.Operation{
		OperationID: "createArticle",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "CreateArticle",
		Description: "Create a new article with auto-generated slug and tag extraction",
		Tags:        []string{"Article"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleHandler.HandleCreateArticle)

	huma.Register(articleGroup, huma.Operation{
		OperationID: "listArticles",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "ListArticles",
		Description: "List articles with pagination, fuzzy search, tag filter and published_at sorting",
		Tags:        []string{"Article"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleHandler.HandleListArticles)

	huma.Register(articleGroup, huma.Operation{
		OperationID: "updateArticle",
		Method:      http.MethodPatch,
		Path:        "/",
		Summary:     "UpdateArticle",
		Description: "Update article fields",
		Tags:        []string{"Article"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleHandler.HandleUpdateArticle)

	huma.Register(articleGroup, huma.Operation{
		OperationID: "deleteArticle",
		Method:      http.MethodDelete,
		Path:        "/",
		Summary:     "DeleteArticle",
		Description: "Delete article and its tag associations",
		Tags:        []string{"Article"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleHandler.HandleDeleteArticle)

	huma.Register(articleGroup, huma.Operation{
		OperationID: "getArticle",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "GetArticle",
		Description: "Get article details by slug. Non-owner can only view published articles, owner can view all.",
		Tags:        []string{"Article"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleHandler.HandleGetArticle)

	huma.Register(articleGroup, huma.Operation{
		OperationID: "uploadArticleImage",
		Method:      http.MethodPost,
		Path:        "/image",
		Summary:     "UploadArticleImage",
		Description: "Upload an image for article. Returns presigned URL of the uploaded image.",
		Tags:        []string{"Article"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, articleHandler.HandleUploadArticleImage)
}
