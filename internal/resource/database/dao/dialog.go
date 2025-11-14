package dao

import (
	"github.com/hcd233/aris-mem-api/internal/resource/database/model"
)

// UserDAO 用户DAO
//
//	author centonhuang
//	update 2024-10-17 02:30:24
type DialogDAO struct {
	baseDAO[model.Dialog]
}
