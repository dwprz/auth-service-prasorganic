package server

import (
	"fmt"
	"net"

	"github.com/dwprz/prasorganic-auth-service/src/common/log"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/interceptor"
	pb "github.com/dwprz/prasorganic-proto/protogen/otp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Grpc struct {
	port                     string
	server                   *grpc.Server
	otpHandler               pb.OtpServiceServer
	unaryResponseInterceptor *interceptor.UnaryResponse
}

// this main grpc server
func NewGrpc(port string, otpHandler pb.OtpServiceServer, uri *interceptor.UnaryResponse) *Grpc {
	return &Grpc{
		port:                     port,
		otpHandler:               otpHandler,
		unaryResponseInterceptor: uri,
	}
}

func (g *Grpc) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", g.port))
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "server.Grpc/Run", "section": "net.Listen"}).Fatal(err)
	}

	log.Logger.Infof("grpc run in port: %s", g.port)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			g.unaryResponseInterceptor.Recovery,
			g.unaryResponseInterceptor.Error,
		))

	g.server = grpcServer

	pb.RegisterOtpServiceServer(grpcServer, g.otpHandler)

	if err := grpcServer.Serve(listener); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "server.Grpc/Run", "section": "grpcServer.Serve"}).Fatal(err)
	}
}

func (g *Grpc) Stop() {
	g.server.Stop()
}
