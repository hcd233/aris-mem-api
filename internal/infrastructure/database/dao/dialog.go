package dao

import (
	"github.com/hcd233/aris-mem-api/internal/infrastructure/database/model"
)

// DialogDAO 对话DAO
//
//	author centonhuang
//	update 2024-10-17 02:30:24
type DialogDAO struct {
	baseDAO[model.Dialog]
}
