package cache

import (
	"context"

	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
)

type Otp interface {
	Cache(ctx context.Context, data *dto.SendOtpReq)
	FindByEmail(ctx context.Context, email string) *dto.SendOtpReq
	DeleteByEmail(ctx context.Context, email string)
}
