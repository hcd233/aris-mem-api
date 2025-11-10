package middleware

import (
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/resource/cache"
	"github.com/hcd233/aris-mem-api/internal/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// RedisLockMiddleware Redis锁中间件
//
//	@param serviceName string
//	@param key string
//	@param expire time.Duration
//	@return fiber.Handler
//	@author centonhuang
//	@update 2025-11-11 04:52:25
func RedisLockMiddleware(serviceName, key string, expire time.Duration) func(ctx huma.Context, next func(huma.Context)) {
	redis := cache.GetRedisClient()

	return func(ctx huma.Context, next func(huma.Context)) {
		logger := logger.WithCtx(ctx.Context())

		value := ctx.Context().Value(key)
		if value == nil {
			logger.Error("[RedisLockMiddleware] value is nil", zap.String("key", key))
			lo.Must0(util.WriteErrorResponse(ctx.BodyWriter(), constant.ErrInternalError))
			return
		}

		lockKey := fmt.Sprintf("%s:%s:%v", serviceName, key, value)
		lockValue := uuid.New().String()

		success, err := redis.SetNX(ctx.Context(), lockKey, lockValue, expire).Result()
		if err != nil {
			logger.Error("[RedisLockMiddleware] failed to get lock", zap.Error(err))
			lo.Must0(util.WriteErrorResponse(ctx.BodyWriter(), constant.ErrInternalError))
			return
		}

		if !success {
			lockValue, err = redis.Get(ctx.Context(), lockKey).Result()
			if err != nil {
				logger.Error("[RedisLockMiddleware] failed to get lock info", zap.String("lockKey", lockKey), zap.Error(err))
				lo.Must0(util.WriteErrorResponse(ctx.BodyWriter(), constant.ErrInternalError))
				return
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
		next(ctx)
	}
}
