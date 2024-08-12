package util

import (
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/handler"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/server"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/mock/delivery"
	"google.golang.org/grpc"
)

func InitGrpcServerTest(os service.Otp) *server.Grpc {
	otpHandler := handler.NewOtpGrpc(os)
	unaryResponseInterceptor := interceptor.NewUnaryResponse()

	grpcServer := server.NewGrpc(config.Conf.CurrentApp.GrpcPort, otpHandler, unaryResponseInterceptor)
	return grpcServer
}

func InitGrpcClientTest(ugdm *delivery.UserGrpcMock) *client.Grpc {
	userGrpcConn := new(grpc.ClientConn)
	grpcClient := client.NewGrpc(ugdm, userGrpcConn)

	return grpcClient
}
