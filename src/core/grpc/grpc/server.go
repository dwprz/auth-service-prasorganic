package grpc

import (
	"fmt"
	"net"

	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/interceptor"
	pb "github.com/dwprz/prasorganic-proto/protogen/otp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	port                     string
	server                   *grpc.Server
	otpServiceServer         pb.OtpServiceServer
	unaryResponseInterceptor *interceptor.UnaryResponse
	logger                   *logrus.Logger
}

// this main grpc server
func NewServer(port string, uss pb.OtpServiceServer, uri *interceptor.UnaryResponse, l *logrus.Logger) *Server {
	return &Server{
		port:                     port,
		otpServiceServer:         uss,
		unaryResponseInterceptor: uri,
		logger:                   l,
	}
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		s.logger.WithFields(logrus.Fields{"location": "grpc.Server/Run", "section": "net.Listen"}).Fatal(err)
	}

	s.logger.Infof("grpc run in port: %s", s.port)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			s.unaryResponseInterceptor.Recovery,
			s.unaryResponseInterceptor.Error,
		))

	s.server = grpcServer

	pb.RegisterOtpServiceServer(grpcServer, s.otpServiceServer)

	if err := grpcServer.Serve(listener); err != nil {
		s.logger.WithFields(logrus.Fields{"location": "grpc.Server/Run", "section": "grpcServer.Serve"}).Fatal(err)
	}
}

func (s *Server) Stop() {
	s.server.Stop()
}
