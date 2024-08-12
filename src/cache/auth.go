package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dwprz/prasorganic-auth-service/src/common/log"
	"github.com/dwprz/prasorganic-auth-service/src/interface/cache"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type AuthImpl struct {
	redis  *redis.ClusterClient
}

func NewAuth(r *redis.ClusterClient) cache.Auth {
	return &AuthImpl{
		redis:  r,
	}
}

func (a *AuthImpl) CacheRegisterReq(ctx context.Context, data *dto.RegisterReq) {
	key := "register_request:" + data.Email

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "cache.AuthImpl/CacheRegisterReq", "section": "json.Marshal"}).Error(err)
		return
	}

	if _, err := a.redis.SetEx(ctx, key, jsonData, 30*time.Minute).Result(); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "cache.AuthImpl/CacheRegisterReq", "section": "redis.SetEx"}).Error(err)
		return
	}
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
		log.Logger.Errorf("error auth cache find register req (unmarshal): %+v", err.Error())
		return nil
	}

	return registerReq
}

func (a *AuthImpl) DeleteRegisterReq(ctx context.Context, email string) {
	key := "register_request:" + email

	_, err := a.redis.Del(ctx, key).Result()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "cache.AuthImpl/DeleteRegisterReq", "section": "redis.Del"}).Error(err)
	}
}