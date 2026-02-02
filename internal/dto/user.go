// Package dto 用户DTO
package dto

import "github.com/hcd233/aris-mem-api/internal/common/model"

// User 用户实体
//
//	@author centonhuang
//	@update 2026-01-29 14:38:14
type User struct {
	ID     uint   `json:"id" doc:"Unique identifier for the user"`
	Name   string `json:"name" doc:"Display name of the user"`
	Avatar string `json:"avatar" doc:"URL or path to the user's avatar image"`
}

// UpdatedUser 用户实体
//
//	author centonhuang
//	update 2025-01-05 11:37:01
type UpdatedUser struct {
	Name   string `json:"name" doc:"Display name of the user"`
	Email  string `json:"email" doc:"Email address of the user"`
	Avatar string `json:"avatar" doc:"URL or path to the user's avatar image"`
}

// DetailedUser 显示用户实体
//
//	@author centonhuang
//	@update 2025-11-07 02:43:56
type DetailedUser struct {
	CreatedAt  string `json:"createdAt" doc:"Timestamp when the user account was created"`
	LastLogin  string `json:"lastLogin" doc:"Timestamp of the user's last login"`
	Permission string `json:"permission" doc:"Permission level of the user"`
	User
}

// GetCurUserRsp represents the response containing the current user's detailed information
//
//	author centonhuang
//	update 2025-01-04 21:00:59
type GetCurUserRsp struct {
	CommonRsp
	User *DetailedUser `json:"user" doc:"Complete user information including permissions"`
}

// UpdateUserReq represents a request to update the current user's information
//
//	author centonhuang
//	update 2025-01-04 21:19:47
type UpdateUserReq struct {
	Body *UpdateUserReqBody `json:"body" doc:"Request body containing fields to update"`
}

// UpdateUserReqBody contains the fields that can be updated for a user
//
//	author centonhuang
//	update 2025-10-31 02:33:48
type UpdateUserReqBody struct {
	User *UpdatedUser `json:"user" required:"true" doc:"User information to update"`
}

// ListUsersReq represents a request to list users with pagination
//
//	author centonhuang
//	update 2026-02-02 10:20:00
type ListUsersReq struct {
	model.CommonParam
	SortField string `query:"sortField" enum:"id,createdAt,lastLogin,name,email" doc:"Sort field"`
}

// ListUsersRsp represents a response containing a list of users
//
//	author centonhuang
//	update 2026-02-02 10:20:00
type ListUsersRsp struct {
	CommonRsp
	Users    []*DetailedUser `json:"users" doc:"Users to list"`
	PageInfo *model.PageInfo `json:"pageInfo" doc:"Page info"`
}

// ApproveUserReq represents a request to approve a pending user
//
//	author centonhuang
//	update 2026-02-02 10:00:00
type ApproveUserReq struct {
	Body *ApproveUserReqBody `json:"body" doc:"Request body containing user ID to approve"`
}

// ApproveUserReqBody contains the user ID for approval
//
//	author centonhuang
//	update 2026-02-02 10:00:00
type ApproveUserReqBody struct {
	UserID uint `json:"userId" required:"true" doc:"User ID to approve"`
}
