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

type OtpImpl struct {
	redis  *redis.ClusterClient
}

func NewOtp(r *redis.ClusterClient) cache.Otp {
	return &OtpImpl{
		redis:  r,
	}
}

func (o *OtpImpl) Cache(ctx context.Context, data *dto.SendOtpReq) {
	key := "otp:" + data.Email

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "cache.OtpImpl/Cache", "section": "json.Marshal"}).Error(err)
		return
	}

	if _, err := o.redis.SetEx(ctx, key, jsonData, 30*time.Minute).Result(); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "cache.OtpImpl/Cache", "section": "redis.SetEx"}).Error(err)
		return
	}
}

func (o *OtpImpl) FindByEmail(ctx context.Context, email string) *dto.SendOtpReq {
	key := "otp:" + email

	result, err := o.redis.Get(ctx, key).Result()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "cache.OtpImpl/FindByEmail", "section": "redis.Get"}).Error(err)
		return nil
	}

	if result == "" {
		return nil
	}

	sendOtpReq := new(dto.SendOtpReq)

	err = json.Unmarshal([]byte(result), sendOtpReq)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "cache.OtpImpl/FindByEmail", "section": "json.Unmarshal"}).Error(err)
		return nil
	}

	return sendOtpReq
}

func (o *OtpImpl) DeleteByEmail(ctx context.Context, email string) {
	key := "otp:" + email

	_, err := o.redis.Del(ctx, key).Result()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "cache.OtpImpl/DeleteByEmail", "section": "redis.Del"}).Error(err)
	}
}
