package main

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/dwprz/prasorganic-auth-service/src/cache"
	"github.com/dwprz/prasorganic-auth-service/src/common/helper"
	"github.com/dwprz/prasorganic-auth-service/src/common/logger"
	"github.com/dwprz/prasorganic-auth-service/src/core/broker"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/grpc"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/middleware"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/restful"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/cbreaker"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-auth-service/src/service"
	"github.com/go-playground/validator/v10"
)

func handleCloseApp(closeCH chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		close(closeCH)
	}()
}

func main() {
	closeCH := make(chan struct{})
	handleCloseApp(closeCH)

	appStatus := os.Getenv("PRASORGANIC_APP_STATUS")

	logger := logger.New()
	conf := config.New(appStatus, logger)
	redis := database.NewRedisCluster(conf)
	authCache := cache.NewAuth(redis, logger)
	validate := validator.New()
	helper := helper.New()

	cbreaker := cbreaker.New(logger)

	unaryRequestInterceptor := interceptor.NewUnaryRequest(conf)
	userGrpcClient, userGrpcConn := client.NewUserGrpc(cbreaker.UserGrpc, conf, unaryRequestInterceptor)

	grpcClient := grpc.NewClient(userGrpcClient, userGrpcConn, logger)
	defer grpcClient.Close()

	rabbitMQClient := broker.NewRabbitMQClient(conf, logger)
	defer rabbitMQClient.Close()

	authService := service.NewAuth(grpcClient, rabbitMQClient, validate, authCache, logger, conf, helper)
	authRestfulHandler := handler.NewAuthRestful(authService, logger, helper)
	middleware := middleware.New(conf, logger)

	restfulServer := restful.NewServer(authRestfulHandler, middleware, conf)
	defer restfulServer.Stop()

	go restfulServer.Run()

	<-closeCH
}
