package middleware

import (
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// RecoverMiddleware 恢复中间件
//
//	@return fiber.Handler
//	@author centonhuang
//	@update 2025-08-18 20:21:14
func RecoverMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				logger.WithFCtx(c).Error("[Panic Recovery] recovered panic", zap.Any("error", r), zap.ByteString("stack", debug.Stack()))
				rsp := dto.CommonRsp{
					Error: constant.ErrInternalError,
				}
				lo.Must0(c.JSON(lo.Must1(util.WrapHTTPResponse(rsp, nil))))
			}
		}()
		return c.Next()
	}
}
