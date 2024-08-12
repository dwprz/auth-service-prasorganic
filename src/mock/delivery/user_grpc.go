package delivery

import (
	"context"
	pb "github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/stretchr/testify/mock"
)

type UserGrpcMock struct {
	mock.Mock
}

func NewUserGrpcMock() *UserGrpcMock {
	return &UserGrpcMock{
		Mock: mock.Mock{},
	}
}

func (u *UserGrpcMock) FindByEmail(ctx context.Context, email string) (*pb.FindUserResponse, error) {
	arguments := u.Mock.Called(ctx, email)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*pb.FindUserResponse), arguments.Error(1)
}

func (u *UserGrpcMock) FindByRefreshToken(ctx context.Context, data *pb.RefreshToken) (*pb.FindUserResponse, error) {
	arguments := u.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*pb.FindUserResponse), arguments.Error(1)
}

func (u *UserGrpcMock) Create(ctx context.Context, data *pb.RegisterRequest) error {
	arguments := u.Mock.Called(ctx, data)

	return arguments.Error(0)
}

func (u *UserGrpcMock) Upsert(ctx context.Context, data *pb.LoginWithGoogleRequest) (*pb.User, error) {
	arguments := u.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*pb.User), arguments.Error(1)
}

func (u *UserGrpcMock) AddRefreshToken(ctx context.Context, data *pb.AddRefreshToken) error {
	arguments := u.Mock.Called(ctx, data)

	return arguments.Error(0)
}

func (u *UserGrpcMock) SetNullRefreshToken(ctx context.Context, refreshToken string) error {
	arguments := u.Mock.Called(ctx, refreshToken)

	return arguments.Error(0)
}