package checkpoint

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/hcd233/aris-mem-api/internal/resource/cache"
	"github.com/redis/go-redis/v9"
)

type redisCheckPointStore struct {
	redis *redis.Client
}

// NewRedisCheckPointStore 创建Redis检查点存储
//
//	@return compose.CheckPointStore
//	@author centonhuang
//	@update 2025-11-12 15:14:58
func NewRedisCheckPointStore() compose.CheckPointStore {
	return &redisCheckPointStore{redis: cache.GetRedisClient()}
}

func (s *redisCheckPointStore) Get(ctx context.Context, checkPointID string) ([]byte, bool, error) {
	data, err := s.redis.Get(ctx, checkPointID).Result()
	if err != nil {
		return nil, false, err
	}
	if data == "" {
		return nil, false, nil
	}
	return []byte(data), true, nil
}

func (s *redisCheckPointStore) Set(ctx context.Context, checkPointID string, checkPoint []byte) error {
	return s.redis.Set(ctx, checkPointID, checkPoint, 0).Err()
}
