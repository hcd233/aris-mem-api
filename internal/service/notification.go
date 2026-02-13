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
	objdao "github.com/hcd233/aris-mem-api/internal/infrastructure/storage/obj_dao"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/util"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// NotificationService Notification service interface
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type NotificationService interface {
	ListNotifications(ctx context.Context, req *dto.ListNotificationsReq) (rsp *dto.ListNotificationsRsp, err error)
	AckNotification(ctx context.Context, req *dto.AckNotificationReq) (rsp *dto.EmptyRsp, err error)
	CountNotifications(ctx context.Context, req *dto.CountNotificationsReq) (rsp *dto.CountRsp, err error)
}

type notificationService struct {
	userDAO         *dao.UserDAO
	articleDAO      *dao.ArticleDAO
	commentDAO      *dao.CommentDAO
	notificationDAO *dao.NotificationDAO
	imageObjDAO     objdao.ObjDAO
}

// NewNotificationService Create notification service
//
//	return NotificationService
//	author centonhuang
//	update 2026-02-03 22:30:00
func NewNotificationService() NotificationService {
	return &notificationService{
		notificationDAO: dao.GetNotificationDAO(),
		userDAO:         dao.GetUserDAO(),
		articleDAO:      dao.GetArticleDAO(),
		commentDAO:      dao.GetCommentDAO(),
		imageObjDAO:     objdao.GetImageObjDAO(),
	}
}

// ListNotifications List user notifications with pagination
//
//	return *ListNotificationsRsp
//	author centonhuang
//	update 2026-02-03 22:30:00
func (s *notificationService) ListNotifications(ctx context.Context, req *dto.ListNotificationsReq) (*dto.ListNotificationsRsp, error) {
	rsp := &dto.ListNotificationsRsp{}

	db := database.GetDBInstance(ctx)
	logger := logger.WithCtx(ctx)
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	if req.SortField == "" {
		req.SortField = "createdAt"
	}
	if req.Sort == "" {
		req.Sort = enum.SortDesc
	}

	where := &dbmodel.Notification{
		ReceiverID: userID,
		Status:     req.Status,
	}

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
			FieldValueMap: map[string]any{},
		},
	}

	if req.Category != "" {
		commonParam.FilterParam.FieldValueMap["category"] = req.Category
	}

	notifications, pageInfo, err := s.notificationDAO.Paginate(db, where, []string{
		"id", "status", "created_at", "type", "entity_type", "entity_id", "sender_id", "receiver_id",
	}, commonParam)
	if err != nil {
		logger.Error("[NotificationService] failed to paginate notifications", zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	userIDs := lo.Uniq(lo.Map(notifications, func(item *dbmodel.Notification, _ int) uint {
		return item.SenderID
	}))

	users, err := s.userDAO.BatchGetByIDs(db, userIDs, []string{"id", "name", "avatar"})
	if err != nil {
		logger.Error("[NotificationService] failed to get users", zap.Error(err), zap.Uints("userIDs", userIDs))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	userIDUserMap := lo.SliceToMap(users, func(item *dbmodel.User) (uint, *dbmodel.User) {
		return item.ID, item
	})

	entityTypeNotificationMap := lo.GroupBy(notifications, func(item *dbmodel.Notification) enum.NotificationEntityType {
		return item.EntityType
	})

	articleIDs := lo.Uniq(lo.Map(entityTypeNotificationMap[enum.NotificationEntityTypeArticle], func(item *dbmodel.Notification, _ int) uint {
		return item.EntityID
	}))

	commentIDs := lo.Uniq(lo.Map(entityTypeNotificationMap[enum.NotificationEntityTypeComment], func(item *dbmodel.Notification, _ int) uint {
		return item.EntityID
	}))

	articles, err := s.articleDAO.BatchGetByIDs(db, articleIDs, []string{"id", "title", "slug", "user_id", "images"})
	if err != nil {
		logger.Error("[NotificationService] failed to get articles", zap.Error(err), zap.Uints("articleIDs", articleIDs))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	comments, err := s.commentDAO.BatchGetByIDs(db, commentIDs, []string{"id", "content", "article_id", "parent_id"})
	if err != nil {
		logger.Error("[NotificationService] failed to get comments", zap.Error(err), zap.Uints("commentIDs", commentIDs))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	articleIDs = lo.Uniq(lo.Map(comments, func(item *dbmodel.Comment, _ int) uint {
		return item.ArticleID
	}))

	parentCommentIDs := lo.Uniq(lo.Map(comments, func(item *dbmodel.Comment, _ int) uint {
		return item.ParentID
	}))

	parentComments, err := s.commentDAO.BatchGetByIDs(db, parentCommentIDs, []string{"id", "content"})
	if err != nil {
		logger.Error("[NotificationService] failed to get parent comments", zap.Error(err), zap.Uints("parentCommentIDs", parentCommentIDs))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	commentArticles, err := s.articleDAO.BatchGetByIDs(db, articleIDs, []string{"id", "images", "slug", "title"})
	if err != nil {
		logger.Error("[NotificationService] failed to get comment articles", zap.Error(err), zap.Uints("articleIDs", articleIDs))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	articleIDArticleMap := lo.SliceToMap(articles, func(item *dbmodel.Article) (uint, *dbmodel.Article) {
		return item.ID, item
	})
	commentIDCommentMap := lo.SliceToMap(comments, func(item *dbmodel.Comment) (uint, *dbmodel.Comment) {
		return item.ID, item
	})
	parentCommentIDCommentMap := lo.SliceToMap(parentComments, func(item *dbmodel.Comment) (uint, *dbmodel.Comment) {
		return item.ID, item
	})
	articleIDCommentArticleMap := lo.SliceToMap(commentArticles, func(item *dbmodel.Article) (uint, *dbmodel.Article) {
		return item.ID, item
	})

	rsp.Notifications = lo.Map(notifications, func(item *dbmodel.Notification, _ int) *dto.ListedNotification {
		var (
			notifiedArticle *dto.NotifiedArticle
			notifiedComment *dto.NotifiedComment
		)
		sender := userIDUserMap[item.SenderID]

		listedNotification := &dto.ListedNotification{
			ID: item.ID,
			Sender: &dto.User{
				ID:     sender.ID,
				Name:   sender.Name,
				Avatar: sender.Avatar,
			},
			Status:    item.Status,
			Type:      item.Type,
			CreatedAt: item.CreatedAt,
			Article:   notifiedArticle,
			Comment:   notifiedComment,
		}

		switch item.EntityType {
		case enum.NotificationEntityTypeArticle:
			article, ok := articleIDArticleMap[item.EntityID]
			if !ok {
				logger.Warn("[NotificationService] article not found for notification", zap.Uint("notificationID", item.ID), zap.Uint("articleID", item.EntityID))
				return listedNotification
			}
			coverImage := ""
			if len(article.Images) > 0 {
				presignedURL, err := s.imageObjDAO.PresignObject(ctx, article.UserID, article.Images[0])
				if err != nil {
					logger.Warn("[ArticleService] failed to generate presigned URL for cover image",
						zap.Error(err),
						zap.String("coverImage", article.Images[0]))
				} else {
					coverImage = presignedURL.String()
					coverImage = util.ToThumbnailURL(coverImage)
				}
			}
			listedNotification.Article = &dto.NotifiedArticle{
				CoverImage: coverImage,
				Article: dto.Article{
					ID:    article.ID,
					Slug:  article.Slug,
					Title: article.Title,
				},
			}
		case enum.NotificationEntityTypeComment:
			comment, ok := commentIDCommentMap[item.EntityID]
			if !ok {
				logger.Warn("[NotificationService] comment not found for notification", zap.Uint("notificationID", item.ID), zap.Uint("commentID", item.EntityID))
				return listedNotification
			}
			parentComment, ok := parentCommentIDCommentMap[comment.ParentID]
			if !ok {
				logger.Warn("[NotificationService] parent comment not found for comment notification", zap.Uint("notificationID", item.ID), zap.Uint("parentCommentID", comment.ParentID))
				return listedNotification
			}
			article, ok := articleIDCommentArticleMap[comment.ArticleID]
			if !ok {
				logger.Warn("[NotificationService] article not found for comment notification", zap.Uint("notificationID", item.ID), zap.Uint("articleID", comment.ArticleID))
				return listedNotification
			}
			coverImage := ""
			if len(article.Images) > 0 {
				presignedURL, err := s.imageObjDAO.PresignObject(ctx, article.UserID, article.Images[0])
				if err != nil {
					logger.Warn("[ArticleService] failed to generate presigned URL for cover image",
						zap.Error(err),
						zap.String("coverImage", article.Images[0]))
				} else {
					coverImage = presignedURL.String()
					coverImage = util.ToThumbnailURL(coverImage)
				}
			}
			listedNotification.Comment = &dto.NotifiedComment{
				Comment: dto.Comment{
					ID:      comment.ID,
					Content: comment.Content,
				},
				RepliedComment: &dto.Comment{
					ID:      parentComment.ID,
					Content: parentComment.Content,
				},
				RepliedArticle: &dto.Article{
					ID:    article.ID,
					Slug:  article.Slug,
					Title: article.Title,
				},
				CoverImage: coverImage,
			}
		}
		return listedNotification
	})
	rsp.PageInfo = pageInfo

	return rsp, nil
}

// AckNotification Acknowledge notification (mark as read)
//
//	return *EmptyRsp
//	author centonhuang
//	update 2026-02-03 22:30:00
func (s *notificationService) AckNotification(ctx context.Context, req *dto.AckNotificationReq) (*dto.EmptyRsp, error) {
	rsp := &dto.EmptyRsp{}

	if req == nil || req.ID == 0 {
		rsp.Error = constant.ErrBadRequest
		return rsp, nil
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)
	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	// Check if notification exists and belongs to current user
	notification, err := s.notificationDAO.Get(db, &dbmodel.Notification{ID: req.ID}, []string{"id", "receiver_id", "status"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rsp.Error = constant.ErrDataNotExists
			return rsp, nil
		}
		logger.Error("[NotificationService] failed to get notification", zap.Error(err), zap.Uint("notificationID", req.ID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	// Check permission
	if notification.ReceiverID != userID {
		logger.Error("[NotificationService] user not allowed to acknowledge notification", zap.Uint("notificationID", req.ID), zap.Uint("userID", userID))
		rsp.Error = constant.ErrNoPermission
		return rsp, nil
	}

	// Already read, no need to update
	if notification.Status == enum.NotificationStatusRead {
		return rsp, nil
	}

	// Update status to read
	updateFields := map[string]interface{}{
		"status": enum.NotificationStatusRead,
	}

	if err := s.notificationDAO.Update(db, &dbmodel.Notification{ID: req.ID}, updateFields); err != nil {
		logger.Error("[NotificationService] failed to update notification status", zap.Error(err), zap.Uint("notificationID", req.ID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	return rsp, nil
}

// CountNotifications Count user notifications
//
//	return *CountNotificationsRsp
//	author centonhuang
//	update 2026-02-12 15:00:00
//	@receiver s *notificationService
//	@param ctx
//	@param req
//	@return *dto.CountRsp
//	@return error
//	@author centonhuang
//	@update 2026-02-12 18:39:42
func (s *notificationService) CountNotifications(ctx context.Context, req *dto.CountNotificationsReq) (*dto.CountRsp, error) {
	rsp := &dto.CountRsp{}

	db := database.GetDBInstance(ctx)
	logger := logger.WithCtx(ctx)
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	// Build where condition
	where := &dbmodel.Notification{
		ReceiverID: userID,
		Status:     req.Status,
	}

	param := &dao.FilterParam{
		FieldValueMap: map[string]any{
			"category": req.Category,
		},
	}

	// Count notifications for the user
	count, err := s.notificationDAO.Count(db, where, param)
	if err != nil {
		logger.Error("[NotificationService] failed to count notifications", zap.Error(err), zap.Uint("userID", userID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	rsp.Count = count
	return rsp, nil
}
