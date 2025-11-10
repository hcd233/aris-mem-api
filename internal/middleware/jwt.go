// Package middleware 中间件
//
//	update 2024-06-22 11:05:33
package middleware

import (
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/jwt"
	"github.com/hcd233/aris-mem-api/internal/resource/database"
	"github.com/hcd233/aris-mem-api/internal/resource/database/dao"
	"github.com/hcd233/aris-mem-api/internal/util"
	"github.com/samber/lo"
)

// JwtMiddlewareHuma JWT 中间件
//
//	@return ctx huma.Context
//	@return next func(huma.Context)
//	@return func(ctx huma.Context, next func(huma.Context))
//	@author centonhuang
//	@update 2025-11-02 04:17:04
func JwtMiddlewareHuma() func(ctx huma.Context, next func(huma.Context)) {
	dao := dao.GetUserDAO()
	accessTokenSvc := jwt.GetAccessTokenSigner()

	return func(ctx huma.Context, next func(huma.Context)) {
		db := database.GetDBInstance(ctx.Context())

		tokenString := ctx.Header("Authorization")
		tokenString = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(tokenString), "Bearer "))
		if tokenString == "" {
			lo.Must0(util.WriteErrorResponse(ctx.BodyWriter(), constant.ErrUnauthorized))
			return
		}
		userID, err := accessTokenSvc.DecodeToken(tokenString)
		if err != nil {
			lo.Must0(util.WriteErrorResponse(ctx.BodyWriter(), constant.ErrUnauthorized))
			return
		}
		user, err := dao.GetByID(db, userID, []string{"id", "name", "permission"}, []string{})
		if err != nil {
			lo.Must0(util.WriteErrorResponse(ctx.BodyWriter(), constant.ErrInternalError))
			return
		}
		ctx = huma.WithValue(ctx, constant.CtxKeyUserID, user.ID)
		ctx = huma.WithValue(ctx, constant.CtxKeyUserName, user.Name)
		ctx = huma.WithValue(ctx, constant.CtxKeyPermission, user.Permission)
		next(ctx)
	}
}

// JwtMiddlewareFiber JWT 中间件 (Fiber版本)
//
//	@return fiber.Handler
//	@author centonhuang
//	@update 2025-11-10 21:00:00
func JwtMiddlewareFiber() fiber.Handler {
	dao := dao.GetUserDAO()
	accessTokenSvc := jwt.GetAccessTokenSigner()

	return func(ctx *fiber.Ctx) error {
		db := database.GetDBInstance(ctx.Context())
		tokenString := ctx.Get("Authorization")
		tokenString = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(tokenString), "Bearer "))
		if tokenString == "" {
			return util.WriteErrorResponse(ctx.Response().BodyWriter(), constant.ErrUnauthorized)
		}

		userID, err := accessTokenSvc.DecodeToken(tokenString)
		if err != nil {
			return util.WriteErrorResponse(ctx.Response().BodyWriter(), constant.ErrUnauthorized)
		}
		user, err := dao.GetByID(db, userID, []string{"id", "name", "permission"}, []string{})
		if err != nil {
			return util.WriteErrorResponse(ctx.Response().BodyWriter(), constant.ErrInternalError)
		}

		// 使用Fiber的上下文存储
		ctx.Locals(constant.CtxKeyUserID, user.ID)
		ctx.Locals(constant.CtxKeyUserName, user.Name)
		ctx.Locals(constant.CtxKeyPermission, user.Permission)

		return ctx.Next()
	}
}
