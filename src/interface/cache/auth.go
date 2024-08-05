package cache

import (
	"context"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
)

type Auth interface {
	CacheRegisterReq(ctx context.Context, data *dto.RegisterReq)
	FindRegisterReq(ctx context.Context, email string) *dto.RegisterReq
	DeleteRegisterReq(ctx context.Context, email string)
}