package database

import (
	"github.com/dwprz/prasorganic-auth-service/src/common/config"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewRedis(conf *config.Config) *redis.ClusterClient {

	db := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			conf.Redis.AddrNode1,
			conf.Redis.AddrNode2,
			conf.Redis.AddrNode3,
			conf.Redis.AddrNode4,
			conf.Redis.AddrNode5,
			conf.Redis.AddrNode6,
		},
		Password:     conf.Redis.Password,
		DialTimeout:  20 * time.Second,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	})

	return db
}
