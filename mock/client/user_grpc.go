package client

import (
	"context"
	"github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/stretchr/testify/mock"
)

type UserGrpcMock struct {
	mock.Mock
}

func NewUserMock() *UserGrpcMock {
	return &UserGrpcMock{
		Mock: mock.Mock{},
	}
}

func (u *UserGrpcMock) FindByEmail(ctx context.Context, data *user.Email) (*user.FindUserResponse, error) {
	arguments := u.Mock.Called(ctx, data)

	return arguments.Get(0).(*user.FindUserResponse), arguments.Error(1)
}

func (u *UserGrpcMock) Create(ctx context.Context, data *user.RegisterRequest) error {
	arguments := u.Mock.Called(ctx, data)

	return arguments.Error(0)
}
