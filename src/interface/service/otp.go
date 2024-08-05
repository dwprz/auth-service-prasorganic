package service

import (
	"context"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
)

type Otp interface {
	Send(ctx context.Context, email string) error
	Verify(ctx context.Context, data *dto.VerifyOtpReq) error
}
