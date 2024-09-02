package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dwprz/prasorganic-auth-service/src/cache"
	"github.com/dwprz/prasorganic-auth-service/src/common/util"
	"github.com/dwprz/prasorganic-auth-service/src/core/broker"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-auth-service/src/service"
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

	redisDB := database.NewRedisCluster()
	defer redisDB.Close()

	authCache := cache.NewAuth(redisDB)
	otpCache := cache.NewOtp(redisDB)
	util := util.New()

	grpcClient := grpc.InitClient()
	defer grpcClient.Close()

	rabbitMQClient := broker.InitClient()
	defer rabbitMQClient.Email.Close()

	otpService := service.NewOtp(rabbitMQClient, otpCache, util)
	authService := service.NewAuth(grpcClient, otpService, authCache)

	restfulServer := restful.InitServer(authService)
	defer restfulServer.Stop()

	go restfulServer.Run()

	grpcServer := grpc.InitServer(otpService)
	defer grpcServer.Stop()

	go grpcServer.Run()

	<-closeCH
}
