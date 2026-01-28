package model

// Tag 标签
//
//	@author centonhuang
//	@update 2026-01-28 21:56:31
type Tag struct {
	BaseModel
	ID uint `json:"id" gorm:"column:id;primary_key;auto_increment;comment:ID"`
	Name string `json:"name" gorm:"column:name;not null;comment:名称"`
	Slug string `json:"slug" gorm:"column:slug;not null;comment:slug"`
	Description string `json:"description" gorm:"column:description;not null;comment:描述"`
}