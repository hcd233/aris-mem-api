package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/handler"
	"github.com/hcd233/aris-mem-api/internal/middleware"
)

func initImageRouter(imageGroup huma.API) {
	imageHandler := handler.NewImageHandler()

	imageGroup.UseMiddleware(middleware.JwtMiddleware(),
		middleware.LimitUserPermissionMiddleware("image", enum.PermissionUser))

	huma.Register(imageGroup, huma.Operation{
		OperationID: "uploadImage",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "UploadImage",
		Description: "Upload an image. Returns presigned URL of the uploaded image.",
		Tags:        []string{"Image"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
		Middlewares: huma.Middlewares{middleware.RateLimiterMiddleware("imageUpload", "", constant.PeriodImageUpload, constant.LimitImageUpload)},
	}, imageHandler.HandleUploadImage)

	huma.Register(imageGroup, huma.Operation{
		OperationID: "GetCredential",
		Method:      http.MethodGet,
		Path:        "/credential",
		Summary:     "GetCredential",
		Description: "Get COS temporary credential for direct upload from frontend. Returns temporary SecretId, SecretKey and SessionToken.",
		Tags:        []string{"Image"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
		Middlewares: huma.Middlewares{middleware.RateLimiterMiddleware("cosTempCredential", "", constant.PeriodCosTempCredential, constant.LimitCosTempCredential)},
	}, imageHandler.HandleGetCredential)
}
