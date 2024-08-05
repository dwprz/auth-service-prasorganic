package cache

import (
	"context"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"github.com/stretchr/testify/mock"
)

type AuthMock struct {
	mock.Mock
}

func NewAuthMock() *AuthMock {
	return &AuthMock{
		Mock: mock.Mock{},
	}
}

func (a *AuthMock) CacheRegisterReq(ctx context.Context, data *dto.RegisterReq) {}

func (a *AuthMock) FindRegisterReq(ctx context.Context, email string) *dto.RegisterReq {
	arguments := a.Mock.Called(ctx, email)

	return arguments.Get(0).(*dto.RegisterReq)
}

func (a *AuthMock) DeleteRegisterReq(ctx context.Context, email string) {}
