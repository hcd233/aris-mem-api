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
//	update 2026-01-29 10:00:00
type ArticleHandler interface {
	HandleCreateArticle(ctx context.Context, req *dto.CreateArticleReq) (*dto.HTTPResponse[*dto.CreateArticleRsp], error)
	HandleListArticles(ctx context.Context, req *dto.ListArticlesReq) (*dto.HTTPResponse[*dto.ListArticlesRsp], error)
	HandleUpdateArticle(ctx context.Context, req *dto.UpdateArticleReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
	HandleDeleteArticle(ctx context.Context, req *dto.DeleteArticleReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
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

func (h *articleHandler) HandleCreateArticle(ctx context.Context, req *dto.CreateArticleReq) (*dto.HTTPResponse[*dto.CreateArticleRsp], error) {
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
