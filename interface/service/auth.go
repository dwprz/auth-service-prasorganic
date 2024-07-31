package service

import (
	"context"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
)

type Authentication interface {
	Register(ctx context.Context, data *dto.RegisterReq) (string, error)
	VerifyRegister(ctx context.Context, data *dto.VerifyRegisterReq) error
}
