package service

import (
	"context"

	"github.com/dwprz/prasorganic-auth-service/interface/cache"
	"github.com/dwprz/prasorganic-auth-service/interface/client"
	"github.com/dwprz/prasorganic-auth-service/interface/helper"
	"github.com/dwprz/prasorganic-auth-service/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/grpc"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthImpl struct {
	grpcClient     *grpc.Client
	rabbitMQClient client.RabbitMQ
	validate       *validator.Validate
	cache          cache.Authentication
	logger         *logrus.Logger
	conf           *config.Config
	helper         helper.Helper
}

func NewAuth(gc *grpc.Client, rc client.RabbitMQ, v *validator.Validate, c cache.Authentication,
	l *logrus.Logger, conf *config.Config, h helper.Helper) service.Authentication {
	return &AuthImpl{
		grpcClient:     gc,
		rabbitMQClient: rc,
		validate:       v,
		cache:          c,
		logger:         l,
		conf:           conf,
		helper:         h,
	}
}

func (a *AuthImpl) Register(ctx context.Context, data *dto.RegisterReq) (string, error) {
	if err := a.validate.Struct(data); err != nil {
		return "", err
	}

	result, err := a.grpcClient.User.FindByEmail(ctx, &user.Email{Email: data.Email})

	if err != nil {
		return "", err
	}

	if result.Data != nil {
		return "", &errors.Response{Code: 409, Message: "user already exists"}
	}

	otp := a.helper.GenerateOtp()
	data.Otp = otp

	if err := a.cache.CacheRegisterReq(ctx, data); err != nil {
		return "", err
	}

	request := &dto.VerifyRegisterReq{
		Email: data.Email,
		Otp:   otp,
	}

	go a.rabbitMQClient.Publish(ctx, "email", "otp", request)

	return data.Email, nil
}

func (a *AuthImpl) VerifyRegister(ctx context.Context, data *dto.VerifyRegisterReq) error {
	if err := a.validate.Struct(data); err != nil {
		return err
	}

	registerReq := a.cache.FindRegisterReq(ctx, data.Email)
	if registerReq == nil {
		return &errors.Response{Code: 404, Message: "register request not found"}
	}

	if registerReq.Otp != data.Otp {
		return &errors.Response{Code: 400, Message: "otp is invalid"}
	}

	encryptPwd, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	registerReq.Password = string(encryptPwd)

	req := &user.RegisterRequest{}
	copier.Copy(req, registerReq)

	if err = a.grpcClient.User.Create(ctx, req); err != nil {
		return err
	}

	return nil
}
