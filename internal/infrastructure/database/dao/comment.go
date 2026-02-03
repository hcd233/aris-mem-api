package dao

import (
	dbmodel "github.com/hcd233/aris-mem-api/internal/infrastructure/database/model"
)

// CommentDAO comment DAO
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type CommentDAO struct {
	baseDAO[dbmodel.Comment]
}
