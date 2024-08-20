package service

import (
	"context"

	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/common/helper"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/client"
	v "github.com/dwprz/prasorganic-auth-service/src/infrastructure/validator"
	"github.com/dwprz/prasorganic-auth-service/src/interface/cache"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"github.com/dwprz/prasorganic-auth-service/src/model/entity"
	pb "github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/jinzhu/copier"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthImpl struct {
	grpcClient *client.Grpc
	otpService service.Otp
	authCache  cache.Auth
}

func NewAuth(gc *client.Grpc, os service.Otp, ac cache.Auth) service.Auth {
	return &AuthImpl{
		grpcClient: gc,
		otpService: os,
		authCache:  ac,
	}
}

func (a *AuthImpl) Register(ctx context.Context, data *dto.RegisterReq) (string, error) {
	if err := v.Validate.Struct(data); err != nil {
		return "", err
	}

	result, err := a.grpcClient.User.FindByEmail(ctx, data.Email)

	if err != nil {
		return "", err
	}

	if result.Data != nil {
		return "", &errors.Response{HttpCode: 409, Message: "user already exists"}
	}

	encryptPwd, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	data.Password = string(encryptPwd)

	go a.otpService.Send(ctx, data.Email)
	go a.authCache.CacheRegisterReq(ctx, data)

	return data.Email, nil
}

func (a *AuthImpl) VerifyRegister(ctx context.Context, data *dto.VerifyOtpReq) error {
	if err := a.otpService.Verify(ctx, data); err != nil {
		return err
	}

	registerReq := a.authCache.FindRegisterReq(ctx, data.Email)
	if registerReq == nil {
		return &errors.Response{HttpCode: 404, Message: "register request not found"}
	}

	req := new(pb.RegisterReq)
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

	go a.authCache.DeleteRegisterReq(context.Background(), data.Email)

	return nil
}

func (a *AuthImpl) LoginWithGoogle(ctx context.Context, data *dto.LoginWithGoogleReq) (*dto.LoginWithGoogleRes, error) {
	if err := v.Validate.Struct(data); err != nil {
		return nil, err
	}

	req := new(pb.LoginWithGoogleReq)
	if err := copier.Copy(req, data); err != nil {
		return nil, err
	}

	res, err := a.grpcClient.User.Upsert(ctx, req)
	if err != nil {
		return nil, err
	}

	user := new(dto.LoginWithGoogleRes)
	if err := copier.Copy(user, res); err != nil {
		return nil, err
	}

	return user, nil
}

func (a *AuthImpl) Login(ctx context.Context, data *dto.LoginReq) (*dto.LoginRes, error) {
	if err := v.Validate.Struct(data); err != nil {
		return nil, err
	}

	res, err := a.grpcClient.User.FindByEmail(ctx, data.Email)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, &errors.Response{HttpCode: 404, Message: "there are no users that match this email"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(res.Data.Password), []byte(data.Password)); err != nil {
		return nil, &errors.Response{HttpCode: 401, Message: "password is invalid"}
	}

	accessToken, err := helper.GenerateAccessToken(res.Data.UserId, res.Data.Email, res.Data.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := helper.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	go a.grpcClient.User.AddRefreshToken(ctx, &pb.AddRefreshTokenReq{
		Email: data.Email,
		Token: refreshToken,
	})

	user := new(entity.SanitizedUser)
	if err := copier.Copy(user, res.Data); err != nil {
		return nil, err
	}

	user.CreatedAt = res.Data.CreatedAt.AsTime()
	user.UpdatedAt = res.Data.UpdatedAt.AsTime()

	return &dto.LoginRes{
		Data: user,
		Tokens: &entity.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (a *AuthImpl) RefreshToken(ctx context.Context, refreshToken string) (*entity.Tokens, error) {
	res, err := a.grpcClient.User.FindByRefreshToken(ctx, &pb.RefreshToken{
		Token: refreshToken,
	})

	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, &errors.Response{HttpCode: 404, Message: "there are no users that match this refresh token"}
	}

	accessToken, err := helper.GenerateAccessToken(res.Data.UserId, res.Data.Email, res.Data.Role)
	if err != nil {
		return nil, err
	}

	return &entity.Tokens{
		AccessToken: accessToken,
	}, nil
}

func (a *AuthImpl) SetNullRefreshToken(ctx context.Context, refreshToken string) error {
	go a.grpcClient.User.SetNullRefreshToken(ctx, refreshToken)

	return nil
}
