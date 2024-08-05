package service

import (
	"context"

	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"github.com/dwprz/prasorganic-auth-service/src/model/entity"
)

type Auth interface {
	Register(ctx context.Context, data *dto.RegisterReq) (string, error)
	VerifyRegister(ctx context.Context, data *dto.VerifyOtpReq) error
	LoginWithGoogle(ctx context.Context, data *dto.LoginWithGoogleReq) (*dto.LoginWithGoogleRes, error)
	Login(ctx context.Context, data *dto.LoginReq) (*dto.LoginRes, error)
	RefreshToken(ctx context.Context, refreshToken string) (*entity.Tokens, error)
	SetNullRefreshToken(ctx context.Context, refreshToken string) error 
}
