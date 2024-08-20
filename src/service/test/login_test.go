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
	pb "github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

// go test -p=1 -v ./src/service/test/... -count=1
// go test -run ^TestService_Login$ -v ./src/service/test/ -count=1

type LoginTestSuite struct {
	suite.Suite
	userGrpcDelivery *delivery.UserGrpcMock
	authService      service.Auth
	otpService       *svcmock.OtpMock
	authCache        *cache.AuthMock
}

func (l *LoginTestSuite) SetupSuite() {
	// mock
	l.userGrpcDelivery = delivery.NewUserGrpcMock()
	userGrpcConn := new(grpc.ClientConn)

	grpcClient := client.NewGrpc(l.userGrpcDelivery, userGrpcConn)

	// mock
	l.authCache = cache.NewAuthMock()
	l.otpService = svcmock.NewOtpMock()

	l.authService = serviceimpl.NewAuth(grpcClient, l.otpService, l.authCache)
}

func (l *LoginTestSuite) Test_Success() {
	user := &pb.User{
		Email:    "johndoe123@gmail.com",
		Password: "$2a$10$3VeOlmczWiQCHN1PiwJwguCeI9b0t7rtgQwaHwAtDTQs7PSlRXCl6",
	}

	req := &dto.LoginReq{
		Email:    user.Email,
		Password: "rahasia",
	}

	l.MockUserGrpcClient_FindByEmail(req.Email, user, nil)
	l.MockUserGrpcClient_AddRefreshToken(&pb.AddRefreshTokenReq{Email: req.Email})

	res, err := l.authService.Login(context.Background(), req)
	assert.NoError(l.T(), err)

	assert.Equal(l.T(), res.Data.Email, req.Email)
	assert.NotEmpty(l.T(), res.Tokens.AccessToken)
	assert.NotEmpty(l.T(), res.Tokens.RefreshToken)
}

func (l *LoginTestSuite) Test_NotFound() {
	req := &dto.LoginReq{
		Email:    "notfound123@gmail.com",
		Password: "rahasia",
	}

	l.MockUserGrpcClient_FindByEmail(req.Email, nil, nil)

	res, err := l.authService.Login(context.Background(), req)
	assert.Error(l.T(), err)

	resErr, ok := err.(*errors.Response)
	assert.True(l.T(), ok)
	assert.Equal(l.T(), resErr.HttpCode, 404)

	assert.Nil(l.T(), res)
}

func (l *LoginTestSuite) MockUserGrpcClient_FindByEmail(email string, returnArg *pb.User, err error) {

	l.userGrpcDelivery.Mock.On("FindByEmail", mock.Anything, email).Return(&pb.FindUserRes{
		Data: returnArg,
	}, err)
}

func (l *LoginTestSuite) MockUserGrpcClient_AddRefreshToken(data *pb.AddRefreshTokenReq) {

	l.userGrpcDelivery.Mock.On("AddRefreshToken", mock.Anything, mock.MatchedBy(func(req *pb.AddRefreshTokenReq) bool {
		return req.Email == data.Email && req.Token != ""
	})).Return(nil)
}

func TestService_Login(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}
