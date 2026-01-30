package service

import (
	"context"

	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database/dao"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database/model"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// TagService 标签服务
//
//	author centonhuang
//	update 2026-01-29 10:00:00
type TagService interface {
	ListTags(ctx context.Context, req *dto.ListTagsReq) (rsp *dto.ListTagsRsp, err error)
}

type tagService struct {
	tagDAO *dao.TagDAO
}

// NewTagService 创建标签服务
//
//	return TagService
//	author centonhuang
//	update 2026-01-29 10:00:00
func NewTagService() TagService {
	return &tagService{
		tagDAO: dao.GetTagDAO(),
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

	tags, pageInfo, err := s.tagDAO.Paginate(db, &model.Tag{}, []string{"id", "name", "slug", "views", "created_at", "updated_at"}, commonParam)
	if err != nil {
		logger.Error("[TagService] failed to list tags", zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	rsp.Tags = lo.Map(tags, func(item *model.Tag, _ int) *dto.DetailedTag {
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
