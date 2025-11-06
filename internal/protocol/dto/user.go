// Package dto 用户DTO
package dto

// User 用户实体
//
//	author centonhuang
//	update 2025-01-05 11:37:01
type User struct {
	Name   string `json:"name" doc:"Display name of the user"`
	Email  string `json:"email,omitempty" doc:"Email address of the user"`
	Avatar string `json:"avatar" doc:"URL or path to the user's avatar image"`
}

// DetailedUser 显示用户实体
//
//	@author centonhuang
//	@update 2025-11-07 02:43:56
type DetailedUser struct {
	User
	ID         uint   `json:"id" doc:"Unique identifier for the user"`
	CreatedAt  string `json:"createdAt,omitempty" doc:"Timestamp when the user account was created"`
	LastLogin  string `json:"lastLogin,omitempty" doc:"Timestamp of the user's last login"`
	Permission string `json:"permission,omitempty" doc:"Permission level of the user"`
}

// GetCurUserInfoReq represents a request to get the current authenticated user's information
//
//	author centonhuang
//	update 2025-01-04 21:00:54
type GetCurUserInfoReq struct {
	UserID uint `json:"userID" doc:"User ID extracted from the JWT token"`
}

// GetCurUserInfoResp represents the response containing the current user's detailed information
//
//	author centonhuang
//	update 2025-01-04 21:00:59
type GetCurUserInfoResp struct {
	User *DetailedUser `json:"user" doc:"Complete user information including permissions"`
}

// GetUserInfoReq represents a request to get a specific user's public information
//
//	author centonhuang
//	update 2025-01-04 21:19:41
type GetUserInfoReq struct {
	UserID uint `json:"userID" path:"userID" doc:"Unique identifier of the user to retrieve"`
}

// GetUserInfoResp represents the response containing a user's public information
//
//	author centonhuang
//	update 2025-01-04 21:19:44
type GetUserInfoResp struct {
	User *DetailedUser `json:"user" doc:"Public user information"`
}

// UpdateUserInfoReq represents a request to update the current user's information
//
//	author centonhuang
//	update 2025-01-04 21:19:47
type UpdateUserInfoReq struct {
	Body *UpdateUserInfoReqBody `json:"body" doc:"Request body containing fields to update"`
}

// UpdateUserInfoReqBody contains the fields that can be updated for a user
//
//	author centonhuang
//	update 2025-10-31 02:33:48
type UpdateUserInfoReqBody struct {
	UserName string `json:"userName" doc:"New display name for the user"`
}
