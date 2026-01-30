package service

import (
	"context"
	"errors"

	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database/dao"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database/model"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ActionService action service
//
//	author centonhuang
//	update 2026-01-30 21:00:00
type ActionService interface {
	Do(ctx context.Context, req *dto.ActionReq) (rsp *dto.EmptyRsp, err error)
	Undo(ctx context.Context, req *dto.ActionReq) (rsp *dto.EmptyRsp, err error)
}

type actionService struct {
	articleDAO *dao.ArticleDAO
	actionDAO  *dao.ActionDAO
}

// NewActionService create action service
//
//	return ActionService
//	author centonhuang
//	update 2026-01-30 21:00:00
func NewActionService() ActionService {
	return &actionService{
		articleDAO: dao.GetArticleDAO(),
		actionDAO:  dao.GetActionDAO(),
	}
}

func (s *actionService) Do(ctx context.Context, req *dto.ActionReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	switch req.Body.EntityType {
	case enum.ActionEntityArticle:
		article, err := s.articleDAO.Get(db, &model.Article{ID: req.Body.EntityID}, []string{"id", "likes", "saves"})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				rsp.Error = constant.ErrDataNotExists
				return rsp, nil
			}
			logger.Error("[ActionService] failed to get article", zap.Error(err), zap.Uint("articleID", req.Body.EntityID))
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}
		action := &model.Action{
			UserID:     userID,
			EntityType: enum.ActionEntityArticle,
			EntityID:   req.Body.EntityID,
			ActionType: req.Body.ActionType,
		}

		_, err = s.actionDAO.Get(db, action, []string{"id"})
		if err == nil {
			rsp.Error = constant.ErrDataExists
			return rsp, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[ActionService] failed to get action", zap.Error(err), zap.Uint("articleID", req.Body.EntityID))
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}

		updateFields := map[string]interface{}{}

		switch req.Body.ActionType {
		case enum.ActionTypeLike:
			updateFields["likes"] = article.Likes + 1
		case enum.ActionTypeSave:
			updateFields["saves"] = article.Saves + 1
		}

		err = db.Transaction(func(tx *gorm.DB) error {
			err := s.actionDAO.Create(tx, action)
			if err != nil {
				logger.Error("[ActionService] failed to create action", zap.Error(err), zap.Uint("articleID", req.Body.EntityID))
				return err
			}
			err = s.articleDAO.Update(tx, &model.Article{ID: req.Body.EntityID}, updateFields)
			if err != nil {
				logger.Error("[ActionService] failed to update article", zap.Error(err), zap.Uint("articleID", req.Body.EntityID))
				return err
			}
			return nil
		})
		if err != nil {
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}
	case enum.ActionEntityComment:
		logger.Info("[ActionService] comment action not implemented", zap.String("entityType", string(req.Body.EntityType)))
		rsp.Error = constant.ErrNoImplement
		return rsp, nil
	default:
		logger.Error("[ActionService] invalid entity type", zap.String("entityType", string(req.Body.EntityType)))
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}
	return rsp, nil
}

func (s *actionService) Undo(ctx context.Context, req *dto.ActionReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	switch req.Body.EntityType {
	case enum.ActionEntityArticle:
		article, err := s.articleDAO.Get(db, &model.Article{ID: req.Body.EntityID}, []string{"id", "likes", "saves"})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				rsp.Error = constant.ErrDataNotExists
				return rsp, nil
			}
			logger.Error("[ActionService] failed to get article", zap.Error(err), zap.Uint("articleID", req.Body.EntityID))
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}
		action := &model.Action{
			UserID:     userID,
			EntityType: enum.ActionEntityArticle,
			EntityID:   req.Body.EntityID,
			ActionType: req.Body.ActionType,
		}
		action, err = s.actionDAO.Get(db, action, []string{"id"})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				rsp.Error = constant.ErrDataNotExists
				return rsp, nil
			}
			logger.Error("[ActionService] failed to get action", zap.Error(err), zap.Uint("articleID", req.Body.EntityID))
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}

		updateFields := map[string]interface{}{}

		switch req.Body.ActionType {
		case enum.ActionTypeLike:
			updateFields["likes"] = article.Likes - 1
		case enum.ActionTypeSave:
			updateFields["saves"] = article.Saves - 1
		}

		err = db.Transaction(func(tx *gorm.DB) error {
			err := s.actionDAO.Delete(tx, action)
			if err != nil {
				logger.Error("[ActionService] failed to delete action", zap.Error(err), zap.Uint("articleID", req.Body.EntityID))
				return err
			}
			err = s.articleDAO.Update(tx, &model.Article{ID: req.Body.EntityID}, updateFields)
			if err != nil {
				logger.Error("[ActionService] failed to update article", zap.Error(err), zap.Uint("articleID", req.Body.EntityID))
				return err
			}
			return nil
		})
		if err != nil {
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}
		return rsp, nil
	case enum.ActionEntityComment:
		logger.Info("[ActionService] comment action not implemented", zap.String("entityType", string(req.Body.EntityType)))
		rsp.Error = constant.ErrNoImplement
		return rsp, nil
	default:
		logger.Error("[ActionService] invalid entity type", zap.String("entityType", string(req.Body.EntityType)))
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}
	return rsp, nil
}
