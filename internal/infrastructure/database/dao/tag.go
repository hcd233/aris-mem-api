package dao

import (
	dbmodel "github.com/hcd233/aris-mem-api/internal/infrastructure/database/model"
)

// TagDAO 标签DAO
//
//	author centonhuang
//	update 2026-01-29 10:00:00
type TagDAO struct {
	baseDAO[dbmodel.Tag]
}
