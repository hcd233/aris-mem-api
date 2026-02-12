package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/handler"
	"github.com/hcd233/aris-mem-api/internal/middleware"
)

func initNotificationRouter(notificationGroup huma.API) {
	notificationHandler := handler.NewNotificationHandler()

	notificationGroup.UseMiddleware(middleware.JwtMiddleware(),
		middleware.LimitUserPermissionMiddleware("notification", enum.PermissionUser))

	huma.Register(notificationGroup, huma.Operation{
		OperationID: "listNotifications",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "ListNotifications",
		Description: "List user notifications with pagination, status filter, and type filter (comment/like/save)",
		Tags:        []string{"Notification"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, notificationHandler.HandleListNotifications)

	huma.Register(notificationGroup, huma.Operation{
		OperationID: "ackNotification",
		Method:      http.MethodPatch,
		Path:        "/ack",
		Summary:     "AckNotification",
		Description: "Acknowledge notification by marking it as read",
		Tags:        []string{"Notification"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, notificationHandler.HandleAckNotification)

	huma.Register(notificationGroup, huma.Operation{
		OperationID: "countNotifications",
		Method:      http.MethodGet,
		Path:        "/count",
		Summary:     "CountNotifications",
		Description: "Count user notifications with optional status and type filters",
		Tags:        []string{"Notification"},
		Security: []map[string][]string{
			{"jwtAuth": {}},
		},
	}, notificationHandler.HandleCountNotifications)
}
