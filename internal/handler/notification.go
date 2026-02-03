package handler

import (
	"context"

	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/service"
	"github.com/hcd233/aris-mem-api/internal/util"
)

// NotificationHandler Notification handler interface
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type NotificationHandler interface {
	HandleListNotifications(ctx context.Context, req *dto.ListNotificationsReq) (*dto.HTTPResponse[*dto.ListNotificationsRsp], error)
	HandleAckNotification(ctx context.Context, req *dto.AckNotificationReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
}

type notificationHandler struct {
	svc service.NotificationService
}

// NewNotificationHandler Create notification handler
//
//	return NotificationHandler
//	author centonhuang
//	update 2026-02-03 22:30:00
func NewNotificationHandler() NotificationHandler {
	return &notificationHandler{
		svc: service.NewNotificationService(),
	}
}

func (h *notificationHandler) HandleListNotifications(ctx context.Context, req *dto.ListNotificationsReq) (*dto.HTTPResponse[*dto.ListNotificationsRsp], error) {
	return util.WrapHTTPResponse(h.svc.ListNotifications(ctx, req))
}

func (h *notificationHandler) HandleAckNotification(ctx context.Context, req *dto.AckNotificationReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
	return util.WrapHTTPResponse(h.svc.AckNotification(ctx, req))
}
