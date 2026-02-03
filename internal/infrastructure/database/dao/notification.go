package dao

import (
	dbmodel "github.com/hcd233/aris-mem-api/internal/infrastructure/database/model"
)

// NotificationDAO Notification DAO
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type NotificationDAO struct {
	baseDAO[dbmodel.Notification]
}
