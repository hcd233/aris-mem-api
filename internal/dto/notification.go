package dto

import (
	"time"

	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/common/model"
)

// ListedNotification Listed notification entity
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type ListedNotification struct {
	ID        uint                    `json:"id" doc:"ID of the notification"`
	Content   string                  `json:"content" doc:"Content of the notification"`
	Sender    *User                   `json:"sender" doc:"Sender of the notification"`
	Status    enum.NotificationStatus `json:"status" doc:"Status of the notification"`
	Type      enum.NotificationType   `json:"type" doc:"Type of the notification"`
	Article   *NotifiedArticle        `json:"article,omitempty" doc:"Article of the notification"`
	Comment   *NotifiedComment        `json:"comment,omitempty" doc:"Comment of the notification"`
	CreatedAt time.Time               `json:"createdAt" doc:"Created time of the notification"`
}

// ListNotificationsReq List notifications request
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type ListNotificationsReq struct {
	model.CommonParam
	SortField string `query:"sortField" enum:"id,createdAt" doc:"Sort field"`
	NotificationFilterParam
}

// NotificationFilterParam Notification filter parameters
//
//	author centonhuang
//	update 2026-02-12 16:00:00
type NotificationFilterParam struct {
	Status   enum.NotificationStatus   `query:"status" enum:"unread,read" doc:"Filter by status (unread/read), empty for all"`
	Category enum.NotificationCategory `query:"category" enum:"like_and_save,comment_and_at" doc:"Filter by category (like_and_save/comment_and_at), empty for all"`
}

// ListNotificationsRsp List notifications response
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type ListNotificationsRsp struct {
	CommonRsp
	Notifications []*ListedNotification `json:"notifications" doc:"Notifications to list"`
	PageInfo      *model.PageInfo       `json:"pageInfo" doc:"Page info"`
}

// AckNotificationReq Acknowledge notification request (mark as read)
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type AckNotificationReq struct {
	ID uint `json:"id" query:"id" required:"true" minimum:"1" doc:"ID of the notification to acknowledge"`
}

// CountNotificationsReq Count notifications request
//
//	@author centonhuang
//	@update 2026-02-12 19:14:48
type CountNotificationsReq struct {
	Status   enum.NotificationStatus   `query:"status" enum:"unread,read" doc:"Filter by status (unread/read), empty for all"`
	Category enum.NotificationCategory `query:"category" enum:"like_and_save,comment_and_at" doc:"Filter by category (like_and_save/comment_and_at), empty for all"`
}
