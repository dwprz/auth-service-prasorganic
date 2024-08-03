package util

import (
	"github.com/dwprz/prasorganic-auth-service/src/mock/client"
	"github.com/dwprz/prasorganic-auth-service/src/mock/helper"
	"github.com/dwprz/prasorganic-auth-service/src/cache"
	"github.com/dwprz/prasorganic-auth-service/src/common/logger"
	grpcapp "github.com/dwprz/prasorganic-auth-service/src/core/grpc/grpc"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/middleware"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/restful"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/oauth"
	"github.com/dwprz/prasorganic-auth-service/src/service"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func NewRestfulServer(cum *client.UserGrpcMock) (*restful.Server, *redis.ClusterClient, *config.Config, *logrus.Logger, *helper.HelperMock) {
	logger := logger.New()
	conf := config.New("DEVELOPMENT", logger)
	redisDB := database.NewRedisCluster(conf)
	authCache := cache.NewAuth(redisDB, logger)
	validate := validator.New()
	helper := helper.NewMock()

	userGrpcConn := new(grpc.ClientConn)

	grpcClient := grpcapp.NewClient(cum, userGrpcConn, logger)
	rabbitMQClient := client.NewRabbitMQMock()

	authService := service.NewAuth(grpcClient, rabbitMQClient, validate, authCache, logger, conf, helper)
	googleOauthConf := oauth.NewGoogleConfig(conf, helper)

	authResfulHandler := handler.NewAuthRestful(authService, googleOauthConf, logger, helper)
	middleware := middleware.New(conf, googleOauthConf, logger)

	restfulServer := restful.NewServer(authResfulHandler, middleware, conf)

	return restfulServer, redisDB, conf, logger, helper
}
