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
	"github.com/hcd233/aris-mem-api/internal/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ActionService action service
//
//	author centonhuang
//	update 2026-01-30 22:00:00
type ActionService interface {
	Do(ctx context.Context, req *dto.ActionReq) (rsp *dto.EmptyRsp, err error)
	Undo(ctx context.Context, req *dto.ActionReq) (rsp *dto.EmptyRsp, err error)
}

type actionService struct {
	articleDAO *dao.ArticleDAO
	commentDAO *dao.CommentDAO
	actionDAO  *dao.ActionDAO
}

// NewActionService create action service
//
//	return ActionService
//	author centonhuang
//	update 2026-02-03 22:30:00
func NewActionService() ActionService {
	return &actionService{
		articleDAO: dao.GetArticleDAO(),
		commentDAO: dao.GetCommentDAO(),
		actionDAO:  dao.GetActionDAO(),
	}
}

func (s *actionService) Do(ctx context.Context, req *dto.ActionReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	action := &model.Action{
		UserID:     userID,
		EntityType: enum.ActionEntityArticle,
		EntityID:   req.Body.EntityID,
		ActionType: req.Body.ActionType,
	}

	_, err := s.actionDAO.Get(db, action, []string{"id"})
	if err == nil {
		rsp.Error = constant.ErrDataExists
		return rsp, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("[ActionService] failed to get action", zap.Error(err), zap.Uint("entityID", req.Body.EntityID), zap.String("entityType", string(req.Body.EntityType)))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	var (
		article      *model.Article
		comment      *model.Comment
		updateFields = map[string]interface{}{}
	)
	switch req.Body.EntityType {
	case enum.ActionEntityArticle:
		article, err = s.articleDAO.Get(db, &model.Article{ID: req.Body.EntityID}, []string{"id", "likes", "saves"})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				rsp.Error = constant.ErrDataNotExists
				return rsp, nil
			}
			logger.Error("[ActionService] failed to get article", zap.Error(err), zap.Uint("articleID", req.Body.EntityID))
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}

		switch req.Body.ActionType {
		case enum.ActionTypeLike:
			updateFields["likes"] = lo.ToPtr(article.Likes + 1) // 防止0被判断出空值
		case enum.ActionTypeSave:
			updateFields["saves"] = lo.ToPtr(article.Saves + 1)
		}

	case enum.ActionEntityComment:
		comment, err = s.commentDAO.Get(db, &model.Comment{ID: req.Body.EntityID}, []string{"id", "likes", "saves"})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				rsp.Error = constant.ErrDataNotExists
				return rsp, nil
			}
			logger.Error("[ActionService] failed to get comment", zap.Error(err), zap.Uint("commentID", req.Body.EntityID))
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}

		switch req.Body.ActionType {
		case enum.ActionTypeLike:
			updateFields["likes"] = lo.ToPtr(comment.Likes + 1)
		case enum.ActionTypeSave:
			updateFields["saves"] = lo.ToPtr(comment.Saves + 1)
		}

	default:
		logger.Error("[ActionService] invalid entity type", zap.String("entityType", string(req.Body.EntityType)))
		rsp.Error = constant.ErrBadRequest
	}
	if !util.HasNonZeroValue(updateFields) {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	err = db.Transaction(func(tx *gorm.DB) (err error) {
		switch req.Body.EntityType {
		case enum.ActionEntityArticle:
			err = s.articleDAO.Update(tx, &model.Article{ID: req.Body.EntityID}, updateFields)
		case enum.ActionEntityComment:
			err = s.commentDAO.Update(tx, &model.Comment{ID: req.Body.EntityID}, updateFields)
		}
		if err != nil {
			logger.Error("[ActionService] failed to update entity", zap.Error(err), zap.Uint("entityID", req.Body.EntityID), zap.String("entityType", string(req.Body.EntityType)))
			return err
		}
		err = s.actionDAO.Create(tx, action)
		if err != nil {
			logger.Error("[ActionService] failed to create action", zap.Error(err), zap.Uint("entityID", req.Body.EntityID), zap.String("entityType", string(req.Body.EntityType)))
			return err
		}
		return nil
	})
	if err != nil {
		rsp.Error = constant.ErrInternalError
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
			if article.Likes > 0 {
				updateFields["likes"] = lo.ToPtr(article.Likes - 1)
			}
		case enum.ActionTypeSave:
			if article.Saves > 0 {
				updateFields["saves"] = lo.ToPtr(article.Saves - 1)
			}
		}

		if !util.HasNonZeroValue(updateFields) {
			rsp.Error = constant.ErrBadRequest
			return rsp, nil
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
		}
	case enum.ActionEntityComment:
		comment, err := s.commentDAO.Get(db, &model.Comment{ID: req.Body.EntityID}, []string{"id", "likes", "saves"})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				rsp.Error = constant.ErrDataNotExists
				return rsp, nil
			}
			logger.Error("[ActionService] failed to get comment", zap.Error(err), zap.Uint("commentID", req.Body.EntityID))
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}
		action := &model.Action{
			UserID:     userID,
			EntityType: enum.ActionEntityComment,
			EntityID:   req.Body.EntityID,
			ActionType: req.Body.ActionType,
		}
		action, err = s.actionDAO.Get(db, action, []string{"id"})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				rsp.Error = constant.ErrDataNotExists
				return rsp, nil
			}
			logger.Error("[ActionService] failed to get action", zap.Error(err), zap.Uint("commentID", req.Body.EntityID))
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}

		updateFields := map[string]interface{}{}

		switch req.Body.ActionType {
		case enum.ActionTypeLike:
			if comment.Likes > 0 {
				updateFields["likes"] = lo.ToPtr(comment.Likes - 1)
			}
		case enum.ActionTypeSave:
			if comment.Saves > 0 {
				updateFields["saves"] = lo.ToPtr(comment.Saves - 1)
			}
		}

		if !util.HasNonZeroValue(updateFields) {
			rsp.Error = constant.ErrBadRequest
			return rsp, nil
		}

		err = db.Transaction(func(tx *gorm.DB) error {
			err := s.actionDAO.Delete(tx, action)
			if err != nil {
				logger.Error("[ActionService] failed to delete action", zap.Error(err), zap.Uint("commentID", req.Body.EntityID))
				return err
			}
			err = s.commentDAO.Update(tx, &model.Comment{ID: req.Body.EntityID}, updateFields)
			if err != nil {
				logger.Error("[ActionService] failed to update comment", zap.Error(err), zap.Uint("commentID", req.Body.EntityID))
				return err
			}
			return nil
		})
		if err != nil {
			rsp.Error = constant.ErrInternalError
		}
	default:
		logger.Error("[ActionService] invalid entity type", zap.String("entityType", string(req.Body.EntityType)))
		rsp.Error = constant.ErrBadRequest
	}
	return rsp, nil
}
