package service

import (
	"context"

	"github.com/dwprz/prasorganic-auth-service/src/interface/cache"
	"github.com/dwprz/prasorganic-auth-service/src/interface/client"
	"github.com/dwprz/prasorganic-auth-service/src/interface/helper"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/grpc"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	pb "github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthImpl struct {
	grpcClient     *grpc.Client
	rabbitMQClient client.RabbitMQ
	validate       *validator.Validate
	cache          cache.Auth
	logger         *logrus.Logger
	conf           *config.Config
	helper         helper.Helper
}

func NewAuth(gc *grpc.Client, rc client.RabbitMQ, v *validator.Validate, c cache.Auth,
	l *logrus.Logger, conf *config.Config, h helper.Helper) service.Auth {
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

	result, err := a.grpcClient.User.FindByEmail(ctx, &pb.Email{Email: data.Email})

	if err != nil {
		return "", err
	}

	if result.Data != nil {
		return "", &errors.Response{Code: 409, Message: "user already exists"}
	}

	otp, err := a.helper.GenerateOtp()
	if err != nil {
		return "", err
	}

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

	req := new(pb.RegisterRequest)
	if err := copier.Copy(req, registerReq); err != nil {
		return err
	}

	userId, err := gonanoid.New()
	if err != nil {
		return err
	}

	req.UserId = userId

	if err = a.grpcClient.User.Create(ctx, req); err != nil {
		return err
	}

	return nil
}

func (a *AuthImpl) LoginWithGoogle(ctx context.Context, data *dto.LoginWithGoogleReq) (*dto.LoginRes, error) {
	if err := a.validate.Struct(data); err != nil {
		return nil, err
	}

	req := new(pb.LoginWithGoogleRequest)
	if err := copier.Copy(req, data); err != nil {
		return nil, err
	}

	res, err := a.grpcClient.User.Upsert(ctx, req)
	if err != nil {
		return nil, err
	}

	user := new(dto.LoginRes)
	if err := copier.Copy(user, res); err != nil {
		return nil, err
	}

	return user, nil
}
