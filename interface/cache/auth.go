package cache

import (
	"context"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
)

type Authentication interface {
	CacheRegisterReq(ctx context.Context, data *dto.RegisterReq) error
	FindRegisterReq(ctx context.Context, email string) *dto.RegisterReq
}