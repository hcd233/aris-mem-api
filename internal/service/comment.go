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

// CommentService comment service
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type CommentService interface {
	CreateComment(ctx context.Context, req *dto.CreateCommentReq) (rsp *dto.EmptyRsp, err error)
	ListComments(ctx context.Context, req *dto.ListCommentsReq) (rsp *dto.ListCommentsRsp, err error)
	DeleteComment(ctx context.Context, req *dto.DeleteCommentReq) (rsp *dto.EmptyRsp, err error)
	CountComments(ctx context.Context, req *dto.CountCommentsReq) (rsp *dto.CountRsp, err error)
}

type commentService struct {
	commentDAO      *dao.CommentDAO
	articleDAO      *dao.ArticleDAO
	userDAO         *dao.UserDAO
	actionDAO       *dao.ActionDAO
	notificationDAO *dao.NotificationDAO
}

// NewCommentService create comment service
//
//	return CommentService
//	author centonhuang
//	update 2026-02-03 22:30:00
func NewCommentService() CommentService {
	return &commentService{
		commentDAO:      dao.GetCommentDAO(),
		articleDAO:      dao.GetArticleDAO(),
		userDAO:         dao.GetUserDAO(),
		actionDAO:       dao.GetActionDAO(),
		notificationDAO: dao.GetNotificationDAO(),
	}
}

// CreateComment create a comment
//
//	return *EmptyRsp
//	author centonhuang
//	update 2026-02-03 22:30:00
func (s *commentService) CreateComment(ctx context.Context, req *dto.CreateCommentReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

	if req == nil || req.Body == nil {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	db := database.GetDBInstance(ctx)
	logger := logger.WithCtx(ctx)
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	// Check if article exists
	article, err := s.articleDAO.Get(db, &dbmodel.Article{ID: req.Body.ArticleID}, []string{"id", "user_id"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rsp.Error = constant.ErrDataNotExists
			return rsp, nil
		}
		logger.Error("[CommentService] failed to get article", zap.Error(err), zap.Uint("articleID", req.Body.ArticleID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	notifications := []*dbmodel.Notification{
		{
			SenderID:   userID,
			ReceiverID: article.UserID,
			Type:       enum.NotificationTypeComment,
			Status:     enum.NotificationStatusUnread,
		},
	}

	// If parent comment ID is provided, check if it exists
	if req.Body.ParentID > 0 {
		comment, err := s.commentDAO.Get(db, &dbmodel.Comment{ID: req.Body.ParentID}, []string{"id", "user_id"})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				rsp.Error = constant.ErrDataNotExists
				return rsp, nil
			}
			logger.Error("[CommentService] failed to get parent comment", zap.Error(err), zap.Uint("parentID", req.Body.ParentID))
			rsp.Error = constant.ErrInternalError
			return rsp, nil
		}
		notifications = append(notifications, &dbmodel.Notification{
			SenderID:   userID,
			ReceiverID: comment.UserID,
			Type:       enum.NotificationTypeComment,
			Status:     enum.NotificationStatusUnread,
		})

	}

	// Create comment
	comment := &dbmodel.Comment{
		ArticleID: req.Body.ArticleID,
		UserID:    userID,
		ParentID:  req.Body.ParentID,
		Content:   req.Body.Content,
	}

	err = s.commentDAO.Create(db, comment)
	if err != nil {
		logger.Error("[CommentService] failed to create comment", zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		for _, notification := range notifications {
			notification.EntityType = enum.NotificationEntityTypeComment
			notification.EntityID = comment.ID
			err = s.notificationDAO.Create(tx, notification)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.Warn("[CommentService] failed to create notifications", zap.Error(err))
	}
	return rsp, nil
}

// ListComments list comments with pagination
//
//	return *ListCommentsRsp
//	author centonhuang
//	update 2026-02-03 22:30:00
func (s *commentService) ListComments(ctx context.Context, req *dto.ListCommentsReq) (*dto.ListCommentsRsp, error) {
	rsp := &dto.ListCommentsRsp{}

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
			QueryFields: []string{"content"},
		},
		SortParam: dao.SortParam{
			Sort:      req.Sort,
			SortField: strcase.ToSnake(req.SortField),
		},
		FilterParam: dao.FilterParam{
			FieldValueMap: map[string]any{
				"article_id": req.ArticleID,
				"parent_id":  nil,
			},
		},
	}

	if req.ParentID > 0 {
		commonParam.FilterParam.FieldValueMap["parent_id"] = lo.ToPtr(req.ParentID)
	}

	comments, pageInfo, err := s.commentDAO.Paginate(db, &dbmodel.Comment{}, []string{
		"id", "article_id", "user_id", "parent_id", "content",
		"likes", "saves", "created_at", "updated_at",
	}, commonParam)
	if err != nil {
		logger.Error("[CommentService] failed to paginate comments", zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	// Get user info
	userIDs := lo.Uniq(lo.Map(comments, func(item *dbmodel.Comment, _ int) uint {
		return item.UserID
	}))
	users, err := s.userDAO.BatchGetByIDs(db, userIDs, []string{"id", "name", "avatar"})
	if err != nil {
		logger.Error("[CommentService] failed to get users", zap.Error(err), zap.Uints("userIDs", userIDs))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	userIDUserMap := lo.SliceToMap(users, func(item *dbmodel.User) (uint, *dbmodel.User) {
		return item.ID, item
	})

	// Get like/save status for current user
	commentIDs := lo.Uniq(lo.Map(comments, func(item *dbmodel.Comment, _ int) uint {
		return item.ID
	}))

	likeActions, err := s.actionDAO.BatchGetByUserIDAndActionType(db, userID, "like", "comment", commentIDs, []string{"id", "entity_id"})
	if err != nil {
		logger.Error("[CommentService] failed to get like actions", zap.Error(err), zap.Uints("commentIDs", commentIDs))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	saveActions, err := s.actionDAO.BatchGetByUserIDAndActionType(db, userID, "save", "comment", commentIDs, []string{"id", "entity_id"})
	if err != nil {
		logger.Error("[CommentService] failed to get save actions", zap.Error(err), zap.Uints("commentIDs", commentIDs))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	likedCommentIDSet := lo.SliceToMap(likeActions, func(item *dbmodel.Action) (uint, struct{}) {
		return item.EntityID, struct{}{}
	})

	savedCommentIDSet := lo.SliceToMap(saveActions, func(item *dbmodel.Action) (uint, struct{}) {
		return item.EntityID, struct{}{}
	})

	// Assemble response
	rsp.Comments = lo.Map(comments, func(item *dbmodel.Comment, _ int) *dto.ListedComment {
		user := userIDUserMap[item.UserID]
		_, liked := likedCommentIDSet[item.ID]
		_, saved := savedCommentIDSet[item.ID]

		return &dto.ListedComment{
			ArticleID: item.ArticleID,
			ParentID:  item.ParentID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			Likes:     item.Likes,
			Saves:     item.Saves,
			Liked:     liked,
			Saved:     saved,
			Author: &dto.User{
				ID:     user.ID,
				Name:   user.Name,
				Avatar: user.Avatar,
			},
			Comment: dto.Comment{
				ID:      item.ID,
				Content: item.Content,
			},
		}
	})
	rsp.PageInfo = pageInfo
	return rsp, nil
}

// DeleteComment delete a comment
//
//	return *EmptyRsp
//	author centonhuang
//	update 2026-02-03 22:30:00
func (s *commentService) DeleteComment(ctx context.Context, req *dto.DeleteCommentReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

	if req == nil || req.ID == 0 {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)
	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	// Check if comment exists and belongs to current user
	comment, err := s.commentDAO.Get(db, &dbmodel.Comment{ID: req.ID}, []string{"id", "user_id"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rsp.Error = constant.ErrDataNotExists
			return rsp, nil
		}
		logger.Error("[CommentService] failed to get comment for deletion", zap.Error(err), zap.Uint("commentID", req.ID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	if comment.UserID != userID {
		logger.Error("[CommentService] user not allowed to delete comment", zap.Uint("commentID", req.ID), zap.Uint("userID", userID))
		rsp.Error = constant.ErrNoPermission
		return rsp, nil
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		err := s.commentDAO.Delete(tx, &dbmodel.Comment{ID: req.ID})
		if err != nil {
			return err
		}
		err = s.commentDAO.Delete(tx, &dbmodel.Comment{ParentID: req.ID})
		return err
	})
	if err != nil {
		logger.Error("[CommentService] failed to delete comment", zap.Error(err), zap.Uint("commentID", req.ID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	return rsp, nil
}

// CountComments count comments for an article
//
//	return *CountCommentsRsp
//	author centonhuang
//	update 2026-02-12 14:50:00
func (s *commentService) CountComments(ctx context.Context, req *dto.CountCommentsReq) (*dto.CountRsp, error) {
	rsp := &dto.CountRsp{}

	if req == nil || req.ArticleID == 0 {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	db := database.GetDBInstance(ctx)
	logger := logger.WithCtx(ctx)

	// Count comments for the article
	count, err := s.commentDAO.Count(db, &dbmodel.Comment{ArticleID: req.ArticleID}, nil)
	if err != nil {
		logger.Error("[CommentService] failed to count comments", zap.Error(err), zap.Uint("articleID", req.ArticleID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	rsp.Count = count
	return rsp, nil
}
