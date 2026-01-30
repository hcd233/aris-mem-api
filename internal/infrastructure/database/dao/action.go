package dao

import (
	dbmodel "github.com/hcd233/aris-mem-api/internal/infrastructure/database/model"
)

// ActionDAO user action DAO
//
//	author centonhuang
//	update 2026-01-30 21:00:00
type ActionDAO struct {
	baseDAO[dbmodel.Action]
}
