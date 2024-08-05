package service

import (
	"context"
	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/interface/cache"
	"github.com/dwprz/prasorganic-auth-service/src/interface/client"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/interface/util"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
)

// delete cache nya menggunakan ctx.Background() supaya tidak cancel, karena ada case context lintas server
type OtpImpl struct {
	validate       *validator.Validate
	rabbitMQClient client.RabbitMQ
	otpCache       cache.Otp
	util           util.Util
}

func NewOtp(v *validator.Validate, rc client.RabbitMQ, oc cache.Otp, u util.Util) service.Otp {
	return &OtpImpl{
		validate:       v,
		rabbitMQClient: rc,
		otpCache:       oc,
		util:           u,
	}
}

func (o *OtpImpl) Send(ctx context.Context, email string) error {
	if err := o.validate.VarCtx(ctx, email, `required,email,min=10,max=100`); err != nil {
		return err
	}

	otp, err := o.util.GenerateOtp()
	if err != nil {
		return err
	}

	sendOtpReq := &dto.SendOtpReq{Email: email, Otp: otp}

	go o.otpCache.Cache(context.Background(), sendOtpReq)
	go o.rabbitMQClient.Publish("email", "otp", sendOtpReq)

	return nil
}

func (o *OtpImpl) Verify(ctx context.Context, data *dto.VerifyOtpReq) error {
	if err := o.validate.StructCtx(ctx, data); err != nil {
		return err
	}

	sendOtpReq := o.otpCache.FindByEmail(ctx, data.Email)
	if sendOtpReq == nil || sendOtpReq.Otp != data.Otp {
		return &errors.Response{HttpCode: 400, GrpcCode: codes.InvalidArgument, Message: "otp is invalid"}
	}

	go o.otpCache.DeleteByEmail(context.Background(), data.Email)

	return nil
}
