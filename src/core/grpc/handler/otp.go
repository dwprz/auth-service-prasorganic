package handler

import (
	"context"

	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	pb "github.com/dwprz/prasorganic-proto/protogen/otp"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OtpGrpcImpl struct {
	otpService service.Otp
	pb.UnimplementedOtpServiceServer
}

func NewOtpGrpc(os service.Otp) pb.OtpServiceServer {
	return &OtpGrpcImpl{
		otpService: os,
	}
}

func (o *OtpGrpcImpl) Send(ctx context.Context, data *pb.SendReq) (*emptypb.Empty, error) {
	if err := o.otpService.Send(ctx, data.Email); err != nil {
		return nil, err
	}

	return nil, nil
}

func (o *OtpGrpcImpl) Verify(ctx context.Context, data *pb.VerifyReq) (*pb.VerifyRes, error) {
	req := new(dto.VerifyOtpReq)
	if err := copier.Copy(req, data); err != nil {
		return nil, err
	}

	if err := o.otpService.Verify(ctx, req); err != nil {
		return nil, err
	}

	return &pb.VerifyRes{Valid: true}, nil
}
