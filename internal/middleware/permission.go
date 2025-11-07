package middleware

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/api"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/util"
	"go.uber.org/zap"
)

// LimitUserPermissionMiddleware 限制用户权限中间件
//
//	@param serviceName string
//	@param requiredPermission model.Permission
//	@return ctx huma.Context
//	@return next func(huma.Context)
//	@return func(ctx huma.Context, next func(huma.Context))
//	@author centonhuang
//	@update 2025-11-02 04:16:51
func LimitUserPermissionMiddleware(serviceName string, requiredPermission enum.Permission) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		permission, ok := ctx.Context().Value(constant.CtxKeyPermission).(enum.Permission)
		if !ok {
			_, err := util.WrapHTTPResponse[any](nil, constant.ErrNoPermission)
			huma.WriteErr(api.GetHumaAPI(), ctx, err.GetStatus(), err.Error(), err)
			return
		}

		if permission.GetLevel() < requiredPermission.GetLevel() {
			logger.WithCtx(ctx.Context()).Info("[LimitUserPermissionMiddleware] permission denied",
				zap.String("serviceName", serviceName),
				zap.String("requiredPermission", string(requiredPermission)),
				zap.String("permission", string(permission)))
			_, err := util.WrapHTTPResponse[any](nil, constant.ErrNoPermission)
			huma.WriteErr(api.GetHumaAPI(), ctx, err.GetStatus(), err.Error(), err)
			return
		}

		next(ctx)
	}
}
