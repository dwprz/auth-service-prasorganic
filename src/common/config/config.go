package config

import (
	"context"
	"os"
	vault "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

type AppConf struct {
	Status        string
	Address       string
	AuthSecretKey string
}

type PostgresConf struct {
	Url      string
	Dsn      string
	User     string
	Password string
}

type RedisConf struct {
	AddrNode1 string
	AddrNode2 string
	AddrNode3 string
	AddrNode4 string
	AddrNode5 string
	AddrNode6 string
	Password  string
}

type OtherConf struct {
	EmailServiceUrl string
}

type Config struct {
	App      *AppConf
	Postgres *PostgresConf
	Redis    *RedisConf
	Other    *OtherConf
}

func NewConfig(logger *logrus.Logger) *Config {
	config := vault.DefaultConfig()

	config.Address = os.Getenv("PRASORGANIC_CONFIG_ADDRESS")

	client, err := vault.NewClient(config)
	if err != nil {
		logger.Errorf("error config (new client): %+v\n", err.Error())
	}

	client.SetToken(os.Getenv("PRASORGANIC_CONFIG_TOKEN"))

	secret, err := client.KVv2("prasorganic-secrets").Get(context.Background(), "auth-service")
	if err != nil {
		logger.Errorf("error config (KVv2): %+v\n", err.Error())
	}

	appStatus := secret.Data["APP_STATUS"].(string)

	appConf := &AppConf{}
	postgresConf := &PostgresConf{}
	redisConf := &RedisConf{}
	otherConf := &OtherConf{}

	switch appStatus {
	case "DEVELOPMENT":
		appConf.Status = appStatus
		appConf.Address = secret.Data["APP_ADDRESS_DEVELOPMENT"].(string)
		appConf.AuthSecretKey = secret.Data["AUTH_SECRET_KEY_DEVELOPMENT"].(string)

		postgresConf.Url = secret.Data["POSTGRES_URL_DEVELOPMENT"].(string)
		postgresConf.Dsn = secret.Data["POSTGRES_DSN_DEVELOPMENT"].(string)
		postgresConf.User = secret.Data["POSTGRES_USER_DEVELOPMENT"].(string)
		postgresConf.Password = secret.Data["POSTGRES_PASSWORD_DEVELOPMENT"].(string)

		redisConf.AddrNode1 = secret.Data["REDIS_ADDR_NODE_1_DEVELOPMENT"].(string)
		redisConf.AddrNode2 = secret.Data["REDIS_ADDR_NODE_2_DEVELOPMENT"].(string)
		redisConf.AddrNode3 = secret.Data["REDIS_ADDR_NODE_3_DEVELOPMENT"].(string)
		redisConf.AddrNode4 = secret.Data["REDIS_ADDR_NODE_4_DEVELOPMENT"].(string)
		redisConf.AddrNode5 = secret.Data["REDIS_ADDR_NODE_5_DEVELOPMENT"].(string)
		redisConf.AddrNode6 = secret.Data["REDIS_ADDR_NODE_6_DEVELOPMENT"].(string)
		redisConf.Password = secret.Data["REDIS_PASSWORD_DEVELOPMENT"].(string)

		otherConf.EmailServiceUrl = secret.Data["EMAIL_SERVICE_URL_DEVELOPMENT"].(string)
	case "STAGING":
		appConf.Status = appStatus
		appConf.Address = secret.Data["APP_ADDRESS_STAGING"].(string)
		appConf.AuthSecretKey = secret.Data["AUTH_SECRET_KEY_STAGING"].(string)

		postgresConf.Url = secret.Data["POSTGRES_URL_STAGING"].(string)
		postgresConf.Dsn = secret.Data["POSTGRES_DSN_STAGING"].(string)
		postgresConf.User = secret.Data["POSTGRES_USER_STAGING"].(string)
		postgresConf.Password = secret.Data["POSTGRES_PASSWORD_STAGING"].(string)

		redisConf.AddrNode1 = secret.Data["REDIS_ADDR_NODE_1_STAGING"].(string)
		redisConf.AddrNode2 = secret.Data["REDIS_ADDR_NODE_2_STAGING"].(string)
		redisConf.AddrNode3 = secret.Data["REDIS_ADDR_NODE_3_STAGING"].(string)
		redisConf.AddrNode4 = secret.Data["REDIS_ADDR_NODE_4_STAGING"].(string)
		redisConf.AddrNode5 = secret.Data["REDIS_ADDR_NODE_5_STAGING"].(string)
		redisConf.AddrNode6 = secret.Data["REDIS_ADDR_NODE_6_STAGING"].(string)
		redisConf.Password = secret.Data["REDIS_PASSWORD_STAGING"].(string)

		otherConf.EmailServiceUrl = secret.Data["EMAIL_SERVICE_URL_STAGING"].(string)

	case "PRODUCTION":
		appConf.Status = appStatus
		appConf.Address = secret.Data["APP_ADDRESS_PRODUCTION"].(string)
		appConf.AuthSecretKey = secret.Data["AUTH_SECRET_KEY_PRODUCTION"].(string)

		postgresConf.Url = secret.Data["POSTGRES_URL_PRODUCTION"].(string)
		postgresConf.Dsn = secret.Data["POSTGRES_DSN_PRODUCTION"].(string)
		postgresConf.User = secret.Data["POSTGRES_USER_PRODUCTION"].(string)
		postgresConf.Password = secret.Data["POSTGRES_PASSWORD_PRODUCTION"].(string)

		redisConf.AddrNode1 = secret.Data["REDIS_ADDR_NODE_1_PRODUCTION"].(string)
		redisConf.AddrNode2 = secret.Data["REDIS_ADDR_NODE_2_PRODUCTION"].(string)
		redisConf.AddrNode3 = secret.Data["REDIS_ADDR_NODE_3_PRODUCTION"].(string)
		redisConf.AddrNode4 = secret.Data["REDIS_ADDR_NODE_4_PRODUCTION"].(string)
		redisConf.AddrNode5 = secret.Data["REDIS_ADDR_NODE_5_PRODUCTION"].(string)
		redisConf.AddrNode6 = secret.Data["REDIS_ADDR_NODE_6_PRODUCTION"].(string)
		redisConf.Password = secret.Data["REDIS_PASSWORD_PRODUCTION"].(string)

		otherConf.EmailServiceUrl = secret.Data["EMAIL_SERVICE_URL_PRODUCTION"].(string)
	}

	return &Config{
		App:      appConf,
		Postgres: postgresConf,
		Redis:    redisConf,
		Other:    otherConf,
	}
}
