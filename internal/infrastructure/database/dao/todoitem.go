package dao

import (
	dbmodel "github.com/hcd233/aris-mem-api/internal/infrastructure/database/model"
)

// TodoItemDAO 待办事项DAO
//
//	author centonhuang
//	update 2024-10-17 02:30:24
type TodoItemDAO struct {
	baseDAO[dbmodel.TodoItem]
}
