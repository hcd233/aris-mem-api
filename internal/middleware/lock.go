package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/resource/cache"
	"github.com/hcd233/aris-mem-api/internal/util"
	"go.uber.org/zap"
)

// RedisLockMiddlewareFiber Redis锁中间件 (Fiber版本)
//
//	@param serviceName string
//	@param key string
//	@param expire time.Duration
//	@return fiber.Handler
//	@author centonhuang
//	@update 2025-11-11 02:42:29
func RedisLockMiddleware(serviceName, key string, expire time.Duration) fiber.Handler {
	redis := cache.GetRedisClient()

	return func(ctx *fiber.Ctx) error {
		logger := logger.WithFCtx(ctx)

		value := ctx.Locals(key)

		lockKey := fmt.Sprintf("%s:%s:%v", serviceName, key, value)
		lockValue := uuid.New().String()

		success, err := redis.SetNX(ctx.Context(), lockKey, lockValue, expire).Result()
		if err != nil {
			logger.Error("[RedisLockMiddleware] failed to get lock", zap.Error(err))
			return util.WriteErrorResponse(ctx.Response().BodyWriter(), constant.ErrInternalError)
		}

		if !success {
			lockValue, err = redis.Get(ctx.Context(), lockKey).Result()
			if err != nil {
				logger.Error("[RedisLockMiddleware] failed to get lock info", zap.String("lockKey", lockKey), zap.Error(err))
				return util.WriteErrorResponse(ctx.Response().BodyWriter(), constant.ErrInternalError)
			}
		}

		defer func() {
			luaScript := `
			if redis.call("get", KEYS[1]) == ARGV[1] then
				return redis.call("del", KEYS[1])
			else
				return 0
			end
		`
			if err := redis.Eval(ctx.Context(), luaScript, []string{lockKey}, lockValue).Err(); err != nil {
				logger.Error("[RedisLockMiddleware] failed to release lock", zap.String("lockKey", lockKey), zap.Error(err))
			}
		}()
		return ctx.Next()
	}
}
