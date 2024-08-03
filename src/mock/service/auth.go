package service

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

func (a *AuthMock) Register(ctx context.Context, data *dto.RegisterReq) (string, error) {
	arguments := a.Mock.Called(ctx, data)

	return arguments.Get(0).(string), arguments.Error(1)
}

func (a *AuthMock) VerifyRegister(ctx context.Context, data *dto.VerifyRegisterReq) error {
	arguments := a.Mock.Called(ctx, data)

	return arguments.Error(0)
}

func (a *AuthMock) LoginWithGoogle(ctx context.Context, data *dto.LoginWithGoogleReq) (*dto.LoginWithGoogleRes, error) {
	arguments := a.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*dto.LoginWithGoogleRes), arguments.Error(1)
}

func (a *AuthMock) Login(ctx context.Context, data *dto.LoginReq) (*dto.LoginRes, error) {
	arguments := a.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*dto.LoginRes), arguments.Error(1)
}
