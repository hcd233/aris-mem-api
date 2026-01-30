package model

// Tag 标签
//
//	@author centonhuang
//	@update 2026-01-28 21:56:31
type Tag struct {
	BaseModel
	DeletedAt int64  `json:"deleted_at" gorm:"column:deleted_at;uniqueIndex:uidx_name_deleted_at;comment:删除时间，默认为0"`
	ID        uint   `json:"id" gorm:"column:id;primary_key;auto_increment;comment:ID"`
	Name      string `json:"name" gorm:"column:name;not null;uniqueIndex:uidx_name_deleted_at;comment:名称"`
	Slug      string `json:"slug" gorm:"column:slug;not null;default:'';comment:简写"`
	Views     uint   `json:"views" gorm:"column:views;not null;default:0;comment:浏览数"`
}
