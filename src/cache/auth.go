package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"github.com/dwprz/prasorganic-auth-service/interface/cache"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type AuthImpl struct {
	redis  *redis.ClusterClient
	logger *logrus.Logger
}

func NewAuth(r *redis.ClusterClient, l *logrus.Logger) cache.Authentication {
	return &AuthImpl{
		redis:  r,
		logger: l,
	}
}

func (a *AuthImpl) CacheRegisterReq(ctx context.Context, data *dto.RegisterReq) error {
	key := "register_request:" + data.Email

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cache register (marshal): %w", err)
	}

	if _, err := a.redis.SetEx(ctx, key, jsonData, 30*time.Minute).Result(); err != nil {
		return fmt.Errorf("cache register (setex): %w", err)
	}

	return nil
}

func (a *AuthImpl) FindRegisterReq(ctx context.Context, email string) *dto.RegisterReq {
	key := "register_request:" + email

	result, _ := a.redis.Get(ctx, key).Result()

	if result == "" {
		return nil
	}

	registerReq := &dto.RegisterReq{}

	err := json.Unmarshal([]byte(result), registerReq)
	if err != nil {
		a.logger.Errorf("error auth cache find register req (unmarshal): %+v", err.Error())
		return nil
	}

	return registerReq
}
