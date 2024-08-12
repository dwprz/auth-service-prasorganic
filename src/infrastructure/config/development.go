package config

import (
	"os"

	"github.com/dwprz/prasorganic-auth-service/src/common/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func setUpForDevelopment() *Config {
	err := os.Chdir(os.Getenv("PRASORGANIC_AUTH_SERVICE_WORKSPACE"))
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForDevelopment", "section": "os.Chdir"}).Fatal(err)
	}

	viper := viper.New()
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForDevelopment", "section": "viper.ReadInConfig"}).Fatal(err)
	}

	currentAppConf := new(currentApp)
	currentAppConf.RestfulAddress = viper.GetString("CURRENT_APP_RESTFUL_ADDRESS")
	currentAppConf.GrpcPort = viper.GetString("CURRENT_APP_GRPC_PORT")

	redisConf := new(redis)
	redisConf.AddrNode1 = viper.GetString("REDIS_ADDR_NODE_1")
	redisConf.AddrNode2 = viper.GetString("REDIS_ADDR_NODE_2")
	redisConf.AddrNode3 = viper.GetString("REDIS_ADDR_NODE_3")
	redisConf.AddrNode4 = viper.GetString("REDIS_ADDR_NODE_4")
	redisConf.AddrNode5 = viper.GetString("REDIS_ADDR_NODE_5")
	redisConf.AddrNode6 = viper.GetString("REDIS_ADDR_NODE_6")
	redisConf.Password = viper.GetString("REDIS_PASSWORD")

	apiGatewayConf := new(apiGateway)
	apiGatewayConf.BaseUrl = viper.GetString("API_GATEWAY_BASE_URL")
	apiGatewayConf.BasicAuth = viper.GetString("API_GATEWAY_BASIC_AUTH")
	apiGatewayConf.BasicAuthUsername = viper.GetString("API_GATEWAY_BASIC_AUTH_USERNAME")
	apiGatewayConf.BasicAuthPassword = viper.GetString("API_GATEWAY_BASIC_AUTH_PASSWORD")

	rabbitMQConf := new(rabbitMQEmailService)
	rabbitMQConf.DSN = viper.GetString("RABBITMQ_EMAIL_SERVICE_DSN")

	googleOauthConf := new(googleOauth)
	googleOauthConf.ClientId = viper.GetString("GOOGLE_OAUTH_CLIENT_ID")
	googleOauthConf.ClientSecret = viper.GetString("GOOGLE_OAUTH_CLIENT_SECRET")
	googleOauthConf.RedirectURL = viper.GetString("GOOGLE_OAUTH_REDIRECT_URL")

	jwtConf := new(jwt)
	jwtConf.PrivateKey = loadRSAPrivateKey(viper.GetString("JWT_PRIVATE_KEY"))
	jwtConf.PublicKey = loadRSAPublicKey(viper.GetString("JWT_PUBLIC_KEY"))

	return &Config{
		CurrentApp:           currentAppConf,
		Redis:                redisConf,
		ApiGateway:           apiGatewayConf,
		RabbitMQEmailService: rabbitMQConf,
		GoogleOauth:          googleOauthConf,
		Jwt:                  jwtConf,
	}
}
