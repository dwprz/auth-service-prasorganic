package server

import (
	"context"

	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	pb "github.com/dwprz/prasorganic-proto/protogen/otp"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OtpGrpcImpl struct {
	logger     *logrus.Logger
	otpService service.Otp
	pb.UnimplementedOtpServiceServer
}

func NewOtpGrpc(l *logrus.Logger, os service.Otp) pb.OtpServiceServer {
	return &OtpGrpcImpl{
		logger:     l,
		otpService: os,
	}
}

func (a *OtpGrpcImpl) Send(ctx context.Context, data *pb.SendRequest) (*emptypb.Empty, error) {
	if err := a.otpService.Send(ctx, data.Email); err != nil {
		return nil, err
	}

	return nil, nil
}

func (a *OtpGrpcImpl) Verify(ctx context.Context, data *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	req := new(dto.VerifyOtpReq)
	if err := copier.Copy(req, data); err != nil {
		return nil, err
	}

	if err := a.otpService.Verify(ctx, req); err != nil {
		return nil, err
	}

	return &pb.VerifyResponse{Valid: true}, nil
}
