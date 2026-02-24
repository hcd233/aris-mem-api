package dto

import (
	"time"

	"github.com/hcd233/aris-mem-api/internal/common/model"
)

// Comment comment entity
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type Comment struct {
	ID      uint   `json:"id" doc:"ID of the comment"`
	Content string `json:"content" doc:"Content of the comment"`
}

// NotifiedComment 通知评论实体
//
//	@author centonhuang
//	@update 2026-02-03 22:30:00
type NotifiedComment struct {
	Liked          bool     `json:"liked" doc:"Whether the current user has liked the comment"`
	RepliedComment *Comment `json:"repliedComment,omitempty" doc:"Replied comment"`
	RepliedArticle *Article `json:"repliedArticle" doc:"Replied article"`
	CoverImage     string   `json:"coverImage" doc:"Cover image URL of the comment"`
	Comment
}

// ListedComment listed comment entity
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type ListedComment struct {
	ArticleID uint      `json:"articleID" doc:"Article ID of the comment"`
	ParentID  uint      `json:"parentID" doc:"Parent comment ID"`
	CreatedAt time.Time `json:"createdAt" doc:"Created time of the comment"`
	UpdatedAt time.Time `json:"updatedAt" doc:"Updated time of the comment"`
	Likes     uint      `json:"likes" doc:"Likes of the comment"`
	Saves     uint      `json:"saves" doc:"Saves of the comment"`
	Liked     bool      `json:"liked" doc:"Whether the current user has liked the comment"`
	Saved     bool      `json:"saved" doc:"Whether the current user has saved the comment"`
	Author    *User     `json:"author" doc:"Author of the comment"`
	Comment
}

// CommentFilterParam comment filter params
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type CommentFilterParam struct {
	ArticleID uint `query:"articleID" required:"true" minimum:"1" doc:"Article ID to filter"`
	ParentID  uint `query:"parentID" doc:"Parent comment ID to filter"`
}

// CreateCommentReq create comment request
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type CreateCommentReq struct {
	Body *CreateCommentReqBody `json:"body" doc:"Request body containing fields to create"`
}

// CreateCommentReqBody create comment request body
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type CreateCommentReqBody struct {
	ArticleID uint     `json:"articleID" required:"true" minimum:"1" doc:"Article ID of the comment"`
	ParentID  uint     `json:"parentID" minimum:"0" doc:"Parent comment ID"`
	Content   string   `json:"content" maxLength:"4096" required:"true" doc:"Content of the comment"`
	Images    []string `json:"images" maxItems:"9" doc:"Images of the comment"`
}

// ListCommentsReq list comments request
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type ListCommentsReq struct {
	model.CommonParam
	SortField string `query:"sortField" enum:"id,createdAt,updatedAt" doc:"Sort field"`
	CommentFilterParam
}

// ListCommentsRsp list comments response
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type ListCommentsRsp struct {
	CommonRsp
	Comments []*ListedComment `json:"comments" doc:"Comments to list"`
	PageInfo *model.PageInfo  `json:"pageInfo" doc:"Page info"`
}

// DeleteCommentReq delete comment request
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type DeleteCommentReq struct {
	ID uint `json:"id" query:"id" required:"true" minimum:"1" doc:"Unique identifier for the comment"`
}

// CountCommentsReq count comments request
//
//	author centonhuang
//	update 2026-02-12 14:50:00
type CountCommentsReq struct {
	ArticleID uint `query:"articleID" required:"true" minimum:"1" doc:"Article ID to count comments for"`
}
