package middleware

import (
	"github.com/bytedance/sonic"
	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/samber/lo"
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
			rsp := &dto.CommonRsp{Error: constant.ErrNoPermission}
			_ = lo.Must1(ctx.BodyWriter().Write(lo.Must1(sonic.Marshal(rsp))))
			return
		}

		if permission.GetLevel() < requiredPermission.GetLevel() {
			logger.WithCtx(ctx.Context()).Info("[LimitUserPermissionMiddleware] permission denied",
				zap.String("serviceName", serviceName),
				zap.String("requiredPermission", string(requiredPermission)),
				zap.String("permission", string(permission)))
			rsp := &dto.CommonRsp{Error: constant.ErrNoPermission}
			_ = lo.Must1(ctx.BodyWriter().Write(lo.Must1(sonic.Marshal(rsp))))
			return
		}

		next(ctx)
	}
}
