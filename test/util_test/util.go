package util_test

import (
	"context"
	"github.com/dwprz/prasorganic-auth-service/src/common/config"
)

type UtilTest interface {
	DeleteUser()
}

type UtilTestImpl struct {
	Ctx  context.Context
	Conf *config.Config
}

func NewUtilTest(ctx context.Context, conf *config.Config) UtilTest {
	return &UtilTestImpl{
		Ctx:  ctx,
		Conf: conf,
	}
}

func (util *UtilTestImpl) DeleteUser() {
	// redis := database.NewRedis(util.Conf)

}
