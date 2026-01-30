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
	objdao "github.com/hcd233/aris-mem-api/internal/infrastructure/storage/obj_dao"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/util"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ArticleService 文章服务
//
//	author centonhuang
//	update 2026-01-31 10:00:00
type ArticleService interface {
	CreateArticle(ctx context.Context, req *dto.CreateArticleReq) (rsp *dto.EmptyRsp, err error)
	ListArticles(ctx context.Context, req *dto.ListArticlesReq) (rsp *dto.ListArticlesRsp, err error)
	UpdateArticle(ctx context.Context, req *dto.UpdateArticleReq) (rsp *dto.EmptyRsp, err error)
	DeleteArticle(ctx context.Context, req *dto.DeleteArticleReq) (rsp *dto.EmptyRsp, err error)
	GetArticle(ctx context.Context, req *dto.GetArticleReq) (rsp *dto.GetArticleRsp, err error)
}

type articleService struct {
	userDAO       *dao.UserDAO
	articleDAO    *dao.ArticleDAO
	tagDAO        *dao.TagDAO
	articleTagDAO *dao.ArticleTagDAO
	actionDAO     *dao.ActionDAO
	imageObjDAO   objdao.ObjDAO
}

// NewArticleService 创建文章服务
//
//	return ArticleService
//	author centonhuang
//	update 2026-01-29 16:00:00
func NewArticleService() ArticleService {
	return &articleService{
		userDAO:       dao.GetUserDAO(),
		articleDAO:    dao.GetArticleDAO(),
		tagDAO:        dao.GetTagDAO(),
		articleTagDAO: dao.GetArticleTagDAO(),
		actionDAO:     dao.GetActionDAO(),
		imageObjDAO:   objdao.GetImageObjDAO(),
	}
}

// CreateArticle 创建文章
//
//	return *CreateArticleRsp
//	author centonhuang
//	update 2026-01-31 10:00:00
func (s *articleService) CreateArticle(ctx context.Context, req *dto.CreateArticleReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

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
	tagNames = lo.Uniq(tagNames)

	for _, image := range req.Body.Images {
		exists, err := s.imageObjDAO.CheckObject(ctx, userID, image)
		if err != nil {
			logger.Error("[ArticleService] failed to check image", zap.Error(err), zap.Uint("userID", userID), zap.String("image", image))
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}
		if !exists {
			logger.Error("[ArticleService] image not exists", zap.Uint("userID", userID), zap.String("image", image))
			rsp.Error = constant.ErrDataNotExists
			return rsp, nil
		}
	}

	// 创建文章
	article := &dbmodel.Article{
		UserID:  userID,
		Title:   req.Body.Title,
		Slug:    slug,
		Content: req.Body.Content,
		Images:  req.Body.Images,
		Status:  enum.ArticleStatusPublished,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// 创建文章
		if err := s.articleDAO.Create(tx, article); err != nil {
			return err
		}

		// 处理标签
		if len(tagNames) > 0 {
			tags := make([]*dbmodel.Tag, 0, len(tagNames))
			for _, tagName := range tagNames {
				where := &dbmodel.Tag{Name: tagName}
				tag := &dbmodel.Tag{Name: tagName, Slug: util.GenerateSlug(tagName)}
				tag, err := s.tagDAO.GetOrCreate(tx, where, tag, []string{"id"})
				if err != nil {
					return err
				}
				tags = append(tags, tag)
			}
			err := s.articleTagDAO.BatchCreate(tx, lo.Map(tags, func(item *dbmodel.Tag, _ int) *dbmodel.ArticleTag {
				return &dbmodel.ArticleTag{
					ArticleID: article.ID,
					TagID:     item.ID,
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
			Sort:      req.Sort,
			SortField: strcase.ToSnake(req.SortField),
		},
		FilterParam: dao.FilterParam{
			FieldValueMap: map[string]any{},
		},
	}

	if req.TagSlug != "" {
		_, err := s.tagDAO.Get(db, &dbmodel.Tag{Slug: req.TagSlug}, []string{"id"})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				rsp.Error = constant.ErrDataNotExists
				return rsp, nil
			}
			logger.Error("[ArticleService] failed to get tag", zap.Error(err), zap.String("tagSlug", req.TagSlug))
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}
		commonParam.FilterParam.FieldValueMap["tag"] = req.TagSlug
	}

	if req.UserID != userID {
		commonParam.FilterParam.FieldValueMap["status"] = enum.ArticleStatusPublished
	}

	if req.UserID != 0 {
		commonParam.FilterParam.FieldValueMap["user_id"] = req.UserID
	}

	articles, pageInfo, err := s.articleDAO.Paginate(db, &dbmodel.Article{}, []string{
		"id", "created_at", "updated_at", "published_at",
		"user_id", "title", "slug", "content", "cover_image", "status", "likes",
	}, commonParam)
	if err != nil {
		logger.Error("[ArticleService] failed to paginate articles", zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	userIDs := lo.Map(articles, func(item *dbmodel.Article, _ int) uint {
		return item.UserID
	})
	users, err := s.userDAO.BatchGetByIDs(db, userIDs, []string{"id", "name", "avatar"})
	if err != nil {
		logger.Error("[ArticleService] failed to get users", zap.Error(err), zap.Uints("userIDs", userIDs))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	userIDUserMap := lo.SliceToMap(users, func(item *dbmodel.User) (uint, *dbmodel.User) {
		return item.ID, item
	})

	articleIDs := lo.Map(articles, func(item *dbmodel.Article, _ int) uint {
		return item.ID
	})

	actions, err := s.actionDAO.BatchGetByUserIDAndActionType(db, userID, enum.ActionTypeLike, enum.ActionEntityArticle, articleIDs, []string{"id", "entity_id"})
	if err != nil {
		logger.Error("[ArticleService] failed to get like actions", zap.Error(err), zap.Uints("articleIDs", articleIDs))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	likedArticleIDSet := lo.SliceToMap(actions, func(item *dbmodel.Action) (uint, struct{}) {
		return item.EntityID, struct{}{}
	})

	// 组装文章列表并获取标签信息
	rsp.Articles = lo.Map(articles, func(item *dbmodel.Article, _ int) *dto.ListedArticle {
		user := userIDUserMap[item.UserID]

		// 为封面图片生成预签名 URL
		coverImage := ""
		if len(item.Images) > 0 {
			presignedURL, err := s.imageObjDAO.PresignObject(ctx, item.UserID, item.Images[0])
			if err != nil {
				logger.Warn("[ArticleService] failed to generate presigned URL for cover image",
					zap.Error(err),
					zap.Uint("userID", item.UserID),
					zap.String("coverImage", item.Images[0]))
			}
			coverImage = presignedURL.String()
			coverImage = util.ToThumbnailURL(coverImage)
		}

		_, liked := likedArticleIDSet[item.ID]

		return &dto.ListedArticle{
			ID:          item.ID,
			Slug:        item.Slug,
			CoverImage:  coverImage,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
			PublishedAt: item.PublishedAt,
			Likes:       item.Likes,
			Liked:       liked,
			Article: dto.Article{
				Title: item.Title,
				Author: &dto.User{
					ID:     user.ID,
					Name:   user.Name,
					Avatar: user.Avatar,
				},
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
//	update 2026-01-29 17:00:00
func (s *articleService) UpdateArticle(ctx context.Context, req *dto.UpdateArticleReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

	if req == nil || req.Body == nil || req.Body.ID == 0 {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)
	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	// 检查文章是否存在且属于当前用户，同时获取旧的封面图片
	article, err := s.articleDAO.Get(db, &dbmodel.Article{ID: req.Body.ID}, []string{"id", "user_id", "cover_image"})
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
	if req.Body.Images != nil {
		updateFields["images"] = req.Body.Images
	}

	if !util.HasNonZeroValue(updateFields) {
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
				tags := make([]*dbmodel.Tag, 0, len(tagNames))
				for _, tagName := range tagNames {
					where := &dbmodel.Tag{Name: tagName}
					tag := &dbmodel.Tag{Name: tagName, Slug: util.GenerateSlug(tagName)}
					tag, err := s.tagDAO.GetOrCreate(tx, where, tag, []string{"id"})
					if err != nil {
						return err
					}
					tags = append(tags, tag)
				}
				err := s.articleTagDAO.BatchCreate(tx, lo.Map(tags, func(item *dbmodel.Tag, _ int) *dbmodel.ArticleTag {
					return &dbmodel.ArticleTag{
						ArticleID: req.Body.ID,
						TagID:     item.ID,
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
//	update 2026-01-29 16:00:00
func (s *articleService) DeleteArticle(ctx context.Context, req *dto.DeleteArticleReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

	if req == nil || req.ID == 0 {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)
	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	// 检查文章是否存在且属于当前用户，同时获取封面图片信息
	article, err := s.articleDAO.Get(db, &dbmodel.Article{ID: req.ID}, []string{"id", "user_id", "cover_image"})
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

	// 删除成功后，同步删除COS中的封面图片
	if len(article.Images) > 0 {
		err := s.imageObjDAO.DeleteObjects(ctx, article.UserID, article.Images)
		if err != nil {
			// Log the error but don't fail the request since article is already deleted
			logger.Warn("[ArticleService] failed to delete cover image from storage",
				zap.Error(err),
				zap.Uint("articleID", req.ID),
				zap.Strings("images", article.Images))
		}
	}

	return rsp, nil
}

// GetArticle 通过 slug 获取文章详情
//
//	return *GetArticleRsp
//	author centonhuang
//	update 2026-01-29 11:30:00
func (s *articleService) GetArticle(ctx context.Context, req *dto.GetArticleReq) (*dto.GetArticleRsp, error) {
	rsp := &dto.GetArticleRsp{}

	if req == nil || req.Slug == "" {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	db := database.GetDBInstance(ctx)
	logger := logger.WithCtx(ctx)
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	user, err := s.userDAO.Get(db, &dbmodel.User{ID: userID}, []string{"id", "name", "avatar"})
	if err != nil {
		logger.Error("[ArticleService] failed to get user", zap.Error(err), zap.Uint("userID", userID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	// 查询文章
	article, err := s.articleDAO.Get(db, &dbmodel.Article{Slug: req.Slug}, []string{
		"id", "user_id", "title", "slug", "content", "cover_image", "status",
		"created_at", "updated_at", "published_at",
		"likes", "saves", "views",
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rsp.Error = constant.ErrDataNotExists
			return rsp, nil
		}
		logger.Error("[ArticleService] failed to get article by slug", zap.Error(err), zap.String("slug", req.Slug))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	// 权限检查：如果不是本人，只能查看已发布的文章
	if article.UserID != userID && article.Status != enum.ArticleStatusPublished {
		logger.Info("[ArticleService] user not allowed to access article", zap.Uint("articleID", article.ID), zap.String("slug", req.Slug), zap.Uint("articleUserID", article.UserID), zap.String("status", string(article.Status)))
		rsp.Error = constant.ErrNoPermission
		return rsp, nil
	}

	// 获取文章标签
	tagIDs, err := s.articleTagDAO.GetTagIDsByArticleID(db, article.ID)
	if err != nil {
		logger.Error("[ArticleService] failed to get tag IDs", zap.Error(err), zap.Uint("articleID", article.ID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	tags, err := s.tagDAO.BatchGetByIDs(db, tagIDs, []string{"id", "name", "slug", "views", "created_at", "updated_at"})
	if err != nil {
		logger.Error("[ArticleService] failed to get tags", zap.Error(err), zap.Uints("tagIDs", tagIDs))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	// Generate presigned URL for cover image
	images := make([]string, 0, len(article.Images))
	for _, image := range article.Images {
		presignedURL, err := s.imageObjDAO.PresignObject(ctx, article.UserID, image)
		if err != nil {
			logger.Warn("[ArticleService] failed to generate presigned URL for image", zap.Error(err), zap.Uint("userID", article.UserID), zap.String("image", image))
		}
		images = append(images, presignedURL.String())
	}

	liked, saved := true, true
	_, err = s.actionDAO.Get(db, &dbmodel.Action{UserID: userID, EntityType: enum.ActionEntityArticle, EntityID: article.ID, ActionType: enum.ActionTypeLike}, []string{"id"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			liked = false
		} else {
			logger.Error("[ArticleService] failed to get like action", zap.Error(err), zap.Uint("userID", userID), zap.Uint("articleID", article.ID))
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}
	}
	_, err = s.actionDAO.Get(db, &dbmodel.Action{UserID: userID, EntityType: enum.ActionEntityArticle, EntityID: article.ID, ActionType: enum.ActionTypeSave}, []string{"id"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			saved = false
		} else {
			logger.Error("[ArticleService] failed to get save action", zap.Error(err), zap.Uint("userID", userID), zap.Uint("articleID", article.ID))
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}
	}

	// 组装响应
	rsp.Article = &dto.DetailedArticle{
		ID:          article.ID,
		Slug:        article.Slug,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
		PublishedAt: article.PublishedAt,
		Content:     article.Content,
		Images:      images,
		Status:      article.Status,
		Likes:       article.Likes,
		Saves:       article.Saves,
		Views:       article.Views,
		Liked:       liked,
		Saved:       saved,
		Tags: lo.Map(tags, func(item *dbmodel.Tag, _ int) *dto.DetailedTag {
			return &dto.DetailedTag{
				ID:        item.ID,
				Slug:      item.Slug,
				Views:     item.Views,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
				Tag: dto.Tag{
					Name: item.Name,
				},
			}
		}),
		Article: dto.Article{
			Title: article.Title,
			Author: &dto.User{
				ID:     user.ID,
				Name:   user.Name,
				Avatar: user.Avatar,
			},
		},
	}

	return rsp, nil
}
