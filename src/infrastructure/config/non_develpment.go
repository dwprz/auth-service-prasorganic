package config

import (
	"context"
	"encoding/base64"
	"os"
	"strings"

	"github.com/dwprz/prasorganic-auth-service/src/common/log"
	vault "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

func setUpForNonDevelopment(appStatus string) *Config {
	defaultConf := vault.DefaultConfig()
	defaultConf.Address = os.Getenv("PRASORGANIC_CONFIG_ADDRESS")

	client, err := vault.NewClient(defaultConf)
	if err != nil {
		log.Logger.Fatalf("vault new client: %v", err)
	}

	client.SetToken(os.Getenv("PRASORGANIC_CONFIG_TOKEN"))

	mountPath := "prasorganic-secrets" + "-" + strings.ToLower(appStatus)

	authServiceSecrets, err := client.KVv2(mountPath).Get(context.Background(), "auth-service")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	apiGatewaySecrets, err := client.KVv2(mountPath).Get(context.Background(), "api-gateway")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	rabbitMQEmailServiceSecrets, err := client.KVv2(mountPath).Get(context.Background(), "rabbitmq-email-service")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	oauthSecrets, err := client.KVv2(mountPath).Get(context.Background(), "oauth")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	jwtSecrets, err := client.KVv2(mountPath).Get(context.Background(), "jwt")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	currentAppConf := new(currentApp)
	currentAppConf.RestfulAddress = authServiceSecrets.Data["RESTFUL_ADDRESS"].(string)
	currentAppConf.GrpcPort = authServiceSecrets.Data["GRPC_PORT"].(string)

	redisConf := new(redis)
	redisConf.AddrNode1 = authServiceSecrets.Data["REDIS_ADDR_NODE_1"].(string)
	redisConf.AddrNode2 = authServiceSecrets.Data["REDIS_ADDR_NODE_2"].(string)
	redisConf.AddrNode3 = authServiceSecrets.Data["REDIS_ADDR_NODE_3"].(string)
	redisConf.AddrNode4 = authServiceSecrets.Data["REDIS_ADDR_NODE_4"].(string)
	redisConf.AddrNode5 = authServiceSecrets.Data["REDIS_ADDR_NODE_5"].(string)
	redisConf.AddrNode6 = authServiceSecrets.Data["REDIS_ADDR_NODE_6"].(string)
	redisConf.Password = authServiceSecrets.Data["REDIS_PASSWORD"].(string)

	apiGatewayConf := new(apiGateway)
	apiGatewayConf.BaseUrl = apiGatewaySecrets.Data["BASE_URL"].(string)
	apiGatewayConf.BasicAuth = apiGatewaySecrets.Data["BASIC_AUTH"].(string)
	apiGatewayConf.BasicAuthPassword = apiGatewaySecrets.Data["BASIC_AUTH_USERNAME"].(string)
	apiGatewayConf.BasicAuthUsername = apiGatewaySecrets.Data["BASIC_AUTH_PASSWORD"].(string)

	rabbitMQEmailServiceConf := new(rabbitMQEmailService)
	rabbitMQEmailServiceConf.DSN = rabbitMQEmailServiceSecrets.Data["DSN"].(string)

	googleOauthConf := new(googleOauth)
	googleOauthConf.ClientId = oauthSecrets.Data["GOOGLE_CLIENT_ID"].(string)
	googleOauthConf.ClientSecret = oauthSecrets.Data["GOOGLE_CLIENT_SECRET"].(string)
	googleOauthConf.RedirectURL = oauthSecrets.Data["GOOGLE_REDIRECT_URL"].(string)

	jwtConf := new(jwt)

	jwtPrivateKey := jwtSecrets.Data["PRIVATE_KEY"].(string)
	base64Byte, err := base64.StdEncoding.DecodeString(jwtPrivateKey)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "base64.StdEncoding.DecodeString"}).Fatal(err)
	}
	jwtPrivateKey = string(base64Byte)

	jwtPublicKey := jwtSecrets.Data["Public_KEY"].(string)
	base64Byte, err = base64.StdEncoding.DecodeString(jwtPublicKey)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "base64.StdEncoding.DecodeString"}).Fatal(err)
	}
	jwtPublicKey = string(base64Byte)

	jwtConf.PrivateKey = loadRSAPrivateKey(jwtPrivateKey)
	jwtConf.PublicKey = loadRSAPublicKey(jwtPublicKey)

	return &Config{
		CurrentApp:           currentAppConf,
		Redis:                redisConf,
		ApiGateway:           apiGatewayConf,
		RabbitMQEmailService: rabbitMQEmailServiceConf,
		GoogleOauth:          googleOauthConf,
		Jwt:                  jwtConf,
	}
}
