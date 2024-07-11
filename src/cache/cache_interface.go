package cache

import (
	"context"
	"github.com/dwprz/prasorganic-auth-service/src/common/model/dto"
)

type AuthCache interface {
	CacheRegisterReq(ctx context.Context, data *dto.RegisterReq) error
}
