package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dwprz/prasorganic-auth-service/src/common/model/dto"
	"github.com/redis/go-redis/v9"
)

type AuthCacheImpl struct {
	Redis *redis.ClusterClient
}

func NewAuthCache(redis *redis.ClusterClient) AuthCache {
	return &AuthCacheImpl{
		Redis: redis,
	}
}

func (cache *AuthCacheImpl) CacheRegisterReq(ctx context.Context, data *dto.RegisterReq) error {
	key := "register_request:" + data.Email

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cache register (marshal): %w", err)
	}

	if _, err := cache.Redis.SetEx(ctx, key, jsonData, 30*time.Minute).Result(); err != nil {
		return fmt.Errorf("cache register (setex): %w", err)
	}

	return nil
}
