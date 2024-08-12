package util

import (
	cacheimpl "github.com/dwprz/prasorganic-auth-service/src/cache"
	"github.com/dwprz/prasorganic-auth-service/src/interface/cache"
	"github.com/redis/go-redis/v9"
)

func InitCacheTest(db *redis.ClusterClient) (cache.Auth, cache.Otp) {
	otpCache := cacheimpl.NewOtp(db)
	authCache := cacheimpl.NewAuth(db)

	return authCache, otpCache
}
