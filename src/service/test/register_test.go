package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/mock/cache"
	"github.com/dwprz/prasorganic-auth-service/src/mock/delivery"
	svcmock "github.com/dwprz/prasorganic-auth-service/src/mock/service"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	serviceimpl "github.com/dwprz/prasorganic-auth-service/src/service"
	"github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

// go test -v ./src/service/test/... -count=1 -p=1
// go test -run ^TestService_Register$  -v ./src/service/test/ -count=1

type RegisterTestSuite struct {
	suite.Suite
	userGrpcDelivery *delivery.UserGrpcMock
	authService      service.Auth
	otpService       *svcmock.OtpMock
	authCache        *cache.AuthMock
}

func (r *RegisterTestSuite) SetupSuite() {
	// mock
	r.userGrpcDelivery = delivery.NewUserGrpcMock()
	userGrpcConn := new(grpc.ClientConn)

	grpcClient := client.NewGrpc(r.userGrpcDelivery, userGrpcConn)

	// mock
	r.authCache = cache.NewAuthMock()
	r.otpService = svcmock.NewOtpMock()

	r.authService = serviceimpl.NewAuth(grpcClient, r.otpService, r.authCache)
}

func (r *RegisterTestSuite) Test_Success() {
	req := &dto.RegisterReq{
		Email:    "johndoe123@gmail.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	r.userGrpcDelivery.Mock.On("FindByEmail", mock.Anything, req.Email).Return(&user.FindUserRes{Data: nil}, nil)
	r.otpService.Mock.On("Send", mock.Anything, req.Email).Return(nil)

	email, err := r.authService.Register(context.Background(), req)

	assert.NoError(r.T(), err)
	assert.Equal(r.T(), req.Email, email)
}

func (r *RegisterTestSuite) Test_AlreadyExists() {
	req := &dto.RegisterReq{
		Email:    "existeduser@gmail.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	r.userGrpcDelivery.Mock.On("FindByEmail", mock.Anything, req.Email).Return(&user.FindUserRes{Data: new(user.User)}, nil)

	email, err := r.authService.Register(context.Background(), req)
	errorRes, ok := err.(*errors.Response)

	assert.Equal(r.T(), true, ok)
	assert.Equal(r.T(), 409, errorRes.HttpCode)
	assert.Equal(r.T(), "", email)
}

func TestService_Register(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}
