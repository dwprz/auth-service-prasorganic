package grpc

import (
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/delivery"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/handler"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/server"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/cbreaker"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
)

func InitServer(os service.Otp) *server.Grpc {
	otpHandler := handler.NewOtpGrpc(os)
	unaryResponseInterceptor := interceptor.NewUnaryResponse()

	grpcServer := server.NewGrpc(config.Conf.CurrentApp.GrpcPort, otpHandler, unaryResponseInterceptor)
	return grpcServer
}

func InitClient() *client.Grpc {
	cbreakerUser := cbreaker.UserGrpc
	unaryRequestInterceptor := interceptor.NewUnaryRequest()

	userDelivery, userConn := delivery.NewUserGrpc(cbreakerUser, unaryRequestInterceptor)

	grpcClient := client.NewGrpc(userDelivery, userConn)
	return grpcClient
}
