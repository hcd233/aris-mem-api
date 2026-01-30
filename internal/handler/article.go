package handler

import (
	"context"

	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/service"
	"github.com/hcd233/aris-mem-api/internal/util"
)

// ArticleHandler 文章处理器
//
//	author centonhuang
//	update 2026-01-31 10:00:00
type ArticleHandler interface {
	HandleCreateArticle(ctx context.Context, req *dto.CreateArticleReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
	HandleListArticles(ctx context.Context, req *dto.ListArticlesReq) (*dto.HTTPResponse[*dto.ListArticlesRsp], error)
	HandleUpdateArticle(ctx context.Context, req *dto.UpdateArticleReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
	HandleDeleteArticle(ctx context.Context, req *dto.DeleteArticleReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
	HandleGetArticle(ctx context.Context, req *dto.GetArticleReq) (*dto.HTTPResponse[*dto.GetArticleRsp], error)
	HandleUploadArticleImage(ctx context.Context, req *dto.UploadArticleImageReq) (*dto.HTTPResponse[*dto.UploadArticleImageRsp], error)
}

type articleHandler struct {
	svc service.ArticleService
}

// NewArticleHandler 创建文章处理器
//
//	return ArticleHandler
//	author centonhuang
//	update 2026-01-29 10:00:00
func NewArticleHandler() ArticleHandler {
	return &articleHandler{
		svc: service.NewArticleService(),
	}
}

func (h *articleHandler) HandleCreateArticle(ctx context.Context, req *dto.CreateArticleReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
	return util.WrapHTTPResponse(h.svc.CreateArticle(ctx, req))
}

func (h *articleHandler) HandleListArticles(ctx context.Context, req *dto.ListArticlesReq) (*dto.HTTPResponse[*dto.ListArticlesRsp], error) {
	return util.WrapHTTPResponse(h.svc.ListArticles(ctx, req))
}

func (h *articleHandler) HandleUpdateArticle(ctx context.Context, req *dto.UpdateArticleReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
	return util.WrapHTTPResponse(h.svc.UpdateArticle(ctx, req))
}

func (h *articleHandler) HandleDeleteArticle(ctx context.Context, req *dto.DeleteArticleReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
	return util.WrapHTTPResponse(h.svc.DeleteArticle(ctx, req))
}

func (h *articleHandler) HandleGetArticle(ctx context.Context, req *dto.GetArticleReq) (*dto.HTTPResponse[*dto.GetArticleRsp], error) {
	return util.WrapHTTPResponse(h.svc.GetArticle(ctx, req))
}

func (h *articleHandler) HandleUploadArticleImage(ctx context.Context, req *dto.UploadArticleImageReq) (*dto.HTTPResponse[*dto.UploadArticleImageRsp], error) {
	return util.WrapHTTPResponse(h.svc.UploadArticleImage(ctx, req))
}
