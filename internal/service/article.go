package service

import (
	"context"
	"errors"
	"time"

	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database/dao"
	dbmodel "github.com/hcd233/aris-mem-api/internal/infrastructure/database/model"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ArticleService 文章服务
//
//	author centonhuang
//	update 2026-01-29 10:00:00
type ArticleService interface {
	CreateArticle(ctx context.Context, req *dto.CreateArticleReq) (rsp *dto.CreateArticleRsp, err error)
	ListArticles(ctx context.Context, req *dto.ListArticlesReq) (rsp *dto.ListArticlesRsp, err error)
	UpdateArticle(ctx context.Context, req *dto.UpdateArticleReq) (rsp *dto.EmptyRsp, err error)
	DeleteArticle(ctx context.Context, req *dto.DeleteArticleReq) (rsp *dto.EmptyRsp, err error)
}

type articleService struct {
	articleDAO    *dao.ArticleDAO
	tagDAO        *dao.TagDAO
	articleTagDAO *dao.ArticleTagDAO
}

// NewArticleService 创建文章服务
//
//	return ArticleService
//	author centonhuang
//	update 2026-01-29 10:00:00
func NewArticleService() ArticleService {
	return &articleService{
		articleDAO:    dao.GetArticleDAO(),
		tagDAO:        dao.GetTagDAO(),
		articleTagDAO: dao.GetArticleTagDAO(),
	}
}

// CreateArticle 创建文章
//
//	return *CreateArticleRsp
//	author centonhuang
//	update 2026-01-29 10:00:00
func (s *articleService) CreateArticle(ctx context.Context, req *dto.CreateArticleReq) (*dto.CreateArticleRsp, error) {
	rsp := &dto.CreateArticleRsp{}

	if req == nil || req.Body == nil {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	db := database.GetDBInstance(ctx)
	logger := logger.WithCtx(ctx)
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	// 生成 slug
	slug := util.GenerateSlug(req.Body.Title)

	// 提取标签
	tagNames := util.ExtractTags(req.Body.Content)

	// 创建文章
	article := &dbmodel.Article{
		UserID:  userID,
		Title:   req.Body.Title,
		Slug:    slug,
		Content: req.Body.Content,
		Status:  enum.ArticleStatusDraft,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// 创建文章
		if err := s.articleDAO.Create(tx, article); err != nil {
			return err
		}

		// 处理标签
		if len(tagNames) > 0 {
			tags := make([] *dbmodel.Tag, len(tagNames))
			for _, tagName := range tagNames {
				tag := &dbmodel.Tag{Name: tagName}
				tag, err := s.tagDAO.GetOrCreate(tx, tag, tag, []string{"id"})
				if err != nil {
					return err
				}
				tags = append(tags, tag)
			}
			err := s.articleTagDAO.BatchCreate(tx, lo.Map(tags, func(item *dbmodel.Tag, _ int) *dbmodel.ArticleTag  {
				return &dbmodel.ArticleTag{
					ArticleID: article.ID,
					TagID: item.ID,
				}
			}))
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		logger.Error("[ArticleService] failed to create article", zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	// 获取创建后的文章详情
	rsp.Article = &dto.DetailedArticle{
		ID:          article.ID,
		Slug:        article.Slug,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
		PublishedAt: article.PublishedAt,
		Status:      article.Status,
		Tags:        tagNames,
		Article: dto.Article{
			Title:   article.Title,
			Content: article.Content,
		},
	}
	return rsp, nil
}

// ListArticles 获取文章列表
//
//	return *ListArticlesRsp
//	author centonhuang
//	update 2026-01-29 10:00:00
func (s *articleService) ListArticles(ctx context.Context, req *dto.ListArticlesReq) (*dto.ListArticlesRsp, error) {
	rsp := &dto.ListArticlesRsp{}

	db := database.GetDBInstance(ctx)
	logger := logger.WithCtx(ctx)
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	commonParam := &dao.CommonParam{
		PageParam: dao.PageParam{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		QueryParam: dao.QueryParam{
			Query:       req.Query,
			QueryFields: []string{"title", "content"},
		},
		SortParam: dao.SortParam{
			Sort:      enum.SortDesc,
			SortField: "published_at",
		},
		FilterParam: dao.FilterParam{
			FieldValueMap: map[string]any{
				"tag": req.Tag,
			},
		},
	}

	articles, pageInfo, err := s.articleDAO.Paginate(db, &dbmodel.Article{UserID: userID}, []string{"id", "created_at", "updated_at", "published_at", "title", "slug", "content", "status"}, commonParam)
	if err != nil {
		logger.Error("[ArticleService] failed to list articles", zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	// 组装文章列表并获取标签信息
	rsp.Articles = lo.Map(articles, func(item *dbmodel.Article, _ int) *dto.DetailedArticle {
		var tagNames []string
		// 获取标签
		tagIDs, err := s.articleTagDAO.GetTagIDsByArticleID(db, item.ID)
		if err != nil {
			logger.Info("[ArticleService] failed to get tag ids", zap.Error(err))
			tagIDs = []uint{}
		}

		if len(tagIDs) != 0 {
			tags, err := s.tagDAO.BatchGetByIDs(db, tagIDs, []string{"id", "name"})
			if err != nil {
				logger.Info("[ArticleService] failed to get tag ids", zap.Error(err))
			}
			tagNames = lo.Map(tags, func(item *dbmodel.Tag, _ int) string {return item.Name})
		}
		
		return &dto.DetailedArticle{
			ID:          item.ID,
			Slug:        item.Slug,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
			PublishedAt: item.PublishedAt,
			Status:      item.Status,
			Tags:        tagNames,
			Article: dto.Article{
				Title:   item.Title,
				Content: item.Content,
			},
		}
	})
	rsp.PageInfo = pageInfo
	return rsp, nil
}

// UpdateArticle 更新文章
//
//	return *EmptyRsp
//	author centonhuang
//	update 2026-01-29 10:00:00
func (s *articleService) UpdateArticle(ctx context.Context, req *dto.UpdateArticleReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

	if req == nil || req.Body == nil || req.Body.ID == 0 {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)
	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	// 检查文章是否存在且属于当前用户
	article, err := s.articleDAO.Get(db, &dbmodel.Article{ID: req.Body.ID}, []string{"id", "user_id"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rsp.Error = constant.ErrDataNotExists
			return rsp, nil
		}
		logger.Error("[ArticleService] failed to get article", zap.Error(err), zap.Uint("articleID", req.Body.ID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	if article.UserID != userID {
		rsp.Error = constant.ErrNoPermission
		return rsp, nil
	}

	// 构建更新字段
	updateFields := map[string]interface{}{}
	if req.Body.Title != "" {
		updateFields["title"] = req.Body.Title
		updateFields["slug"] = util.GenerateSlug(req.Body.Title)
	}
	if req.Body.Content != "" {
		updateFields["content"] = req.Body.Content
	}
	if req.Body.Status != "" {
		updateFields["status"] = req.Body.Status
		// 如果状态变为已发布，更新发布时间
		if req.Body.Status == enum.ArticleStatusPublished {
			updateFields["published_at"] = time.Now()
		}
	}

	if !hasNonZeroValue(updateFields) {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// 更新文章
		if err := s.articleDAO.Update(tx, &dbmodel.Article{ID: req.Body.ID}, updateFields); err != nil {
			return err
		}

		// 如果内容更新，重新处理标签
		if req.Body.Content != "" {
			// 删除旧标签关联
			if err := s.articleTagDAO.DeleteByArticleID(tx, req.Body.ID); err != nil {
				return err
			}

			// 提取新标签
			tagNames := util.ExtractTags(req.Body.Content)
			if len(tagNames) > 0 {
				tags := make([] *dbmodel.Tag, len(tagNames))
				for _, tagName := range tagNames {
					tag := &dbmodel.Tag{Name: tagName}
					tag, err := s.tagDAO.GetOrCreate(tx, tag, tag, []string{"id"})
					if err != nil {
						return err
					}
					tags = append(tags, tag)
				}
				err := s.articleTagDAO.BatchCreate(tx, lo.Map(tags, func(item *dbmodel.Tag, _ int) *dbmodel.ArticleTag  {
					return &dbmodel.ArticleTag{
						ArticleID: article.ID,
						TagID: item.ID,
					}
				}))
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		logger.Error("[ArticleService] failed to update article", zap.Error(err), zap.Uint("articleID", req.Body.ID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	return rsp, nil
}

// DeleteArticle 删除文章
//
//	return *EmptyRsp
//	author centonhuang
//	update 2026-01-29 10:00:00
func (s *articleService) DeleteArticle(ctx context.Context, req *dto.DeleteArticleReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

	if req == nil || req.ID == 0 {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)
	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	// 检查文章是否存在且属于当前用户
	article, err := s.articleDAO.Get(db, &dbmodel.Article{ID: req.ID}, []string{"id", "user_id"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rsp.Error = constant.ErrDataNotExists
			return rsp, nil
		}
		logger.Error("[ArticleService] failed to get article for deletion", zap.Error(err), zap.Uint("articleID", req.ID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	if article.UserID != userID {
		rsp.Error = constant.ErrNoPermission
		return rsp, nil
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// 删除文章
		if err := s.articleDAO.Delete(tx, &dbmodel.Article{ID: req.ID}); err != nil {
			return err
		}

		// 删除标签关联
		if err := s.articleTagDAO.DeleteByArticleID(tx, req.ID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error("[ArticleService] failed to delete article", zap.Error(err), zap.Uint("articleID", req.ID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	return rsp, nil
}
