package dto

import (
	"github.com/hcd233/aris-mem-api/internal/common/model"
)

// Tag 标签实体
//
//	@author centonhuang
//	@update 2026-01-29 10:00:00
type Tag struct {
	Name string `json:"name" doc:"Name of the tag"`
}

// DetailedTag 详细标签实体
//
//	@author centonhuang
//	@update 2026-01-29 10:00:00
type DetailedTag struct {
	ID        uint   `json:"id" doc:"Unique identifier for the tag"`
	Slug      string `json:"slug" doc:"Slug of the tag"`
	CreatedAt string `json:"createdAt" doc:"Timestamp when the tag was created"`
	UpdatedAt string `json:"updatedAt" doc:"Timestamp when the tag was updated"`
	Tag
}

// ListTagsReq 获取标签列表请求
//
//	@author centonhuang
//	@update 2026-01-29 10:00:00
type ListTagsReq struct {
	model.CommonParam
	SortField string `query:"sortField" enum:"id,createdAt,updatedAt,name" doc:"Sort field"`
}

// ListTagsRsp 获取标签列表响应
//
//	@author centonhuang
//	@update 2026-01-29 10:00:00
type ListTagsRsp struct {
	CommonRsp
	Tags     []*DetailedTag  `json:"tags" doc:"Tags to list"`
	PageInfo *model.PageInfo `json:"pageInfo" doc:"Page info"`
}
