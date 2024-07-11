package service

import (
	"context"
	"github.com/dwprz/prasorganic-auth-service/src/common/model/dto"
)

type AuthService interface {
	Register(ctx context.Context, data *dto.RegisterReq) (string, error)
	SendOtp(ctx context.Context, otp string, email string)
}
