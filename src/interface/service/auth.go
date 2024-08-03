package service

import (
	"context"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
)

type Auth interface {
	Register(ctx context.Context, data *dto.RegisterReq) (string, error)
	VerifyRegister(ctx context.Context, data *dto.VerifyRegisterReq) error
	LoginWithGoogle(ctx context.Context, data *dto.LoginWithGoogleReq) (*dto.LoginWithGoogleRes, error)
	Login(ctx context.Context, data *dto.LoginReq) (*dto.LoginRes, error)
}
