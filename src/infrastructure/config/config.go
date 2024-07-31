package config

import "github.com/sirupsen/logrus"

type currentApp struct {
	RestfulAddress string
	GrpcPort       string
}

type apiGateway struct {
	BaseUrl           string
	BasicAuth         string
	BasicAuthUsername string
	BasicAuthPassword string
}

type redis struct {
	AddrNode1 string
	AddrNode2 string
	AddrNode3 string
	AddrNode4 string
	AddrNode5 string
	AddrNode6 string
	Password  string
}

type rabbitMQEmailService struct {
	DSN string
}

type Config struct {
	CurrentApp           *currentApp
	Redis                *redis
	ApiGateway           *apiGateway
	RabbitMQEmailService *rabbitMQEmailService
}

func New(appStatus string, logger *logrus.Logger) *Config {
	var config *Config

	if appStatus == "DEVELOPMENT" {

		config = setUpForDevelopment(logger)
		return config
	}

	config = setUpForNonDevelopment(appStatus, logger)
	return config
}
