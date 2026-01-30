package service

import (
	"context"
	"errors"

	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database/dao"
	dbmodel "github.com/hcd233/aris-mem-api/internal/infrastructure/database/model"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TagService 标签服务
//
//	author centonhuang
//	@update 2026-01-29 10:00:00
type TagService interface {
	ListTags(ctx context.Context, req *dto.ListTagsReq) (rsp *dto.ListTagsRsp, err error)
	DeleteTag(ctx context.Context, req *dto.DeleteTagReq) (rsp *dto.EmptyRsp, err error)
}

type tagService struct {
	tagDAO        *dao.TagDAO
	articleTagDAO *dao.ArticleTagDAO
}

// NewTagService 创建标签服务
//
//	return TagService
//	author centonhuang
//	@update 2026-01-29 10:00:00
func NewTagService() TagService {
	return &tagService{
		tagDAO:        dao.GetTagDAO(),
		articleTagDAO: dao.GetArticleTagDAO(),
	}
}

// ListTags 获取标签列表
//
//	return *ListTagsRsp
//	author centonhuang
//	update 2026-01-29 10:00:00
func (s *tagService) ListTags(ctx context.Context, req *dto.ListTagsReq) (*dto.ListTagsRsp, error) {
	rsp := &dto.ListTagsRsp{}

	if req.SortField == "" {
		req.SortField = "id"
	}

	if req.Sort == "" {
		req.Sort = enum.SortAsc
	}

	db := database.GetDBInstance(ctx)

	logger := logger.WithCtx(ctx)

	commonParam := &dao.CommonParam{
		PageParam: dao.PageParam{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		QueryParam: dao.QueryParam{
			Query:       req.Query,
			QueryFields: []string{"name"},
		},
		SortParam: dao.SortParam{
			Sort:      req.Sort,
			SortField: strcase.ToSnake(req.SortField),
		},
	}

	tags, pageInfo, err := s.tagDAO.Paginate(db, &dbmodel.Tag{}, []string{"id", "name", "slug", "views", "created_at", "updated_at"}, commonParam)
	if err != nil {
		logger.Error("[TagService] failed to list tags", zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	rsp.Tags = lo.Map(tags, func(item *dbmodel.Tag, _ int) *dto.DetailedTag {
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
	})
	rsp.PageInfo = pageInfo
	return rsp, nil
}

// DeleteTag 删除标签
//
//	return *EmptyRsp
//	author centonhuang
//	@update 2026-01-31 10:00:00
func (s *tagService) DeleteTag(ctx context.Context, req *dto.DeleteTagReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

	if req == nil || req.ID == 0 {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)
	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	// 检查标签是否存在
	tag, err := s.tagDAO.Get(db, &dbmodel.Tag{ID: req.ID}, []string{"id"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rsp.Error = constant.ErrDataNotExists
			return rsp, nil
		}
		logger.Error("[TagService] failed to get tag for deletion", zap.Error(err), zap.Uint("tagID", req.ID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	_ = tag

	// 检查当前用户是否有权限删除（只有管理员可以删除标签）
	permission := ctx.Value(constant.CtxKeyPermission).(enum.Permission)
	if permission != enum.PermissionAdmin {
		rsp.Error = constant.ErrNoPermission
		return rsp, nil
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// 删除标签
		if err := s.tagDAO.Delete(tx, &dbmodel.Tag{ID: req.ID}); err != nil {
			return err
		}

		// 删除标签关联
		if err := s.articleTagDAO.DeleteByTagID(tx, req.ID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Error("[TagService] failed to delete tag", zap.Error(err), zap.Uint("tagID", req.ID), zap.Uint("userID", userID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	return rsp, nil
}
