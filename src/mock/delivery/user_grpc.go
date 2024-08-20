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

func (u *UserGrpcMock) FindByEmail(ctx context.Context, email string) (*pb.FindUserRes, error) {
	arguments := u.Mock.Called(ctx, email)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*pb.FindUserRes), arguments.Error(1)
}

func (u *UserGrpcMock) FindByRefreshToken(ctx context.Context, data *pb.RefreshToken) (*pb.FindUserRes, error) {
	arguments := u.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*pb.FindUserRes), arguments.Error(1)
}

func (u *UserGrpcMock) Create(ctx context.Context, data *pb.RegisterReq) error {
	arguments := u.Mock.Called(ctx, data)

	return arguments.Error(0)
}

func (u *UserGrpcMock) Upsert(ctx context.Context, data *pb.LoginWithGoogleReq) (*pb.User, error) {
	arguments := u.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*pb.User), arguments.Error(1)
}

func (u *UserGrpcMock) AddRefreshToken(ctx context.Context, data *pb.AddRefreshTokenReq) error {
	arguments := u.Mock.Called(ctx, data)

	return arguments.Error(0)
}

func (u *UserGrpcMock) SetNullRefreshToken(ctx context.Context, refreshToken string) error {
	arguments := u.Mock.Called(ctx, refreshToken)

	return arguments.Error(0)
}
