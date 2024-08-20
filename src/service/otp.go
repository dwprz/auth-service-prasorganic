package service

import (
	"context"

	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/core/broker/producer"
	v "github.com/dwprz/prasorganic-auth-service/src/infrastructure/validator"
	"github.com/dwprz/prasorganic-auth-service/src/interface/cache"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/interface/util"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"google.golang.org/grpc/codes"
)

// caching nya menggunakan ctx.Background() supaya tidak cancel, karena ada case context lintas server
type OtpImpl struct {
	rabbitMQProducer *producer.RabbitMQ
	otpCache         cache.Otp
	util             util.Util
}

func NewOtp(rp *producer.RabbitMQ, oc cache.Otp, u util.Util) service.Otp {
	return &OtpImpl{
		rabbitMQProducer: rp,
		otpCache:         oc,
		util:             u,
	}
}

func (o *OtpImpl) Send(ctx context.Context, email string) error {
	if err := v.Validate.Var(email, `required,email,min=10,max=100`); err != nil {
		return err
	}

	otp, err := o.util.GenerateOtp()
	if err != nil {
		return err
	}

	sendOtpReq := &dto.SendOtpReq{Email: email, Otp: otp}

	go o.otpCache.Cache(context.Background(), sendOtpReq)
	go o.rabbitMQProducer.Email.Publish("email", "otp", sendOtpReq)

	return nil
}

func (o *OtpImpl) Verify(ctx context.Context, data *dto.VerifyOtpReq) error {
	if err := v.Validate.Struct(data); err != nil {
		return err
	}

	sendOtpReq := o.otpCache.FindByEmail(ctx, data.Email)
	if sendOtpReq == nil || sendOtpReq.Otp != data.Otp {
		return &errors.Response{HttpCode: 400, GrpcCode: codes.InvalidArgument, Message: "otp is invalid"}
	}

	go o.otpCache.DeleteByEmail(context.Background(), data.Email)

	return nil
}
