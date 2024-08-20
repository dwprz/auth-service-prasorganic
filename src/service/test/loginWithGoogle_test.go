package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/mock/cache"
	"github.com/dwprz/prasorganic-auth-service/src/mock/delivery"
	svcmock "github.com/dwprz/prasorganic-auth-service/src/mock/service"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	serviceimpl "github.com/dwprz/prasorganic-auth-service/src/service"
	"github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

// go test -p=1 -v ./src/service/test/... -count=1
// go test -run ^TestService_LoginWithGoogle$ -v ./src/service/test/ -count=1

type LoginWithGoogleTestSuite struct {
	suite.Suite
	userGrpcDelivery *delivery.UserGrpcMock
	authService      service.Auth
	otpService       *svcmock.OtpMock
	authCache        *cache.AuthMock
}

func (l *LoginWithGoogleTestSuite) SetupSuite() {
	// mock
	l.userGrpcDelivery = delivery.NewUserGrpcMock()
	userGrpcConn := new(grpc.ClientConn)

	grpcClient := client.NewGrpc(l.userGrpcDelivery, userGrpcConn)

	// mock
	l.authCache = cache.NewAuthMock()
	l.otpService = svcmock.NewOtpMock()

	l.authService = serviceimpl.NewAuth(grpcClient, l.otpService, l.authCache)
}

func (l *LoginWithGoogleTestSuite) Test_Success() {
	req := &dto.LoginWithGoogleReq{
		UserId:       "hyfa_5Sq7nQcaY6ACksXP",
		Email:        "johndoe123@gmail.com",
		FullName:     "John Doe",
		PhotoProfile: "example-photo-profile",
		RefreshToken: `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.
					   eyJleHAiOjE3MjY3MDgyNzMsImlkIjoiQndfb
					   nBmeUhpRXlRY2d3LWM1eGtzIiwiaXNzIjoicH
					   Jhc29yZ2FuaWMtYXV0aC1zZXJ2aWNlIn0.Q0A
					   0veFoysALdg9qqypxY8lKxMFNDmBJR-hyEAHg
					   zBZu-NNzmRIl2MMUjQfBAjJI-KoNj7n9keH6x
					   gwXcuRTVeDfMFtaPm7Dp5ezEgqLXIOTR3o3Sk
					   DhljqYN0HC7UBZQ4R-SyP26MTgVBZQOyuymj3
					   XhFG29OABaJb12GglrtE`,
	}

	l.MockUserGrpcClient_Upsert(req)

	res, err := l.authService.LoginWithGoogle(context.Background(), req)
	assert.NoError(l.T(), err)

	assert.Equal(l.T(), res.UserId, req.UserId)
	assert.Equal(l.T(), res.Email, req.Email)
	assert.Equal(l.T(), res.FullName, req.FullName)
	assert.Equal(l.T(), res.PhotoProfile, req.PhotoProfile)
}

func (l *LoginWithGoogleTestSuite) Test_InvalidEmail() {
	req := &dto.LoginWithGoogleReq{
		UserId:       "hyfa_5Sq7nQcaY6ACksXP",
		Email:        "invalid_email",
		FullName:     "John Doe",
		PhotoProfile: "example-photo-profile",
		RefreshToken: `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.
					   eyJleHAiOjE3MjY3MDgyNzMsImlkIjoiQndfb
					   nBmeUhpRXlRY2d3LWM1eGtzIiwiaXNzIjoicH
					   Jhc29yZ2FuaWMtYXV0aC1zZXJ2aWNlIn0.Q0A
					   0veFoysALdg9qqypxY8lKxMFNDmBJR-hyEAHg
					   zBZu-NNzmRIl2MMUjQfBAjJI-KoNj7n9keH6x
					   gwXcuRTVeDfMFtaPm7Dp5ezEgqLXIOTR3o3Sk
					   DhljqYN0HC7UBZQ4R-SyP26MTgVBZQOyuymj3
					   XhFG29OABaJb12GglrtE`,
	}

	res, err := l.authService.LoginWithGoogle(context.Background(), req)
	assert.Error(l.T(), err)

	validationErr, ok := err.(validator.ValidationErrors)
	assert.True(l.T(), ok)
	assert.Equal(l.T(), validationErr[0].Field(), "Email")

	assert.Nil(l.T(), res)
}

func (l *LoginWithGoogleTestSuite) MockUserGrpcClient_Upsert(data *dto.LoginWithGoogleReq) {

	l.userGrpcDelivery.Mock.On("Upsert", mock.Anything, mock.MatchedBy(func(req *user.LoginWithGoogleReq) bool {
		return req.UserId == data.UserId && req.Email == data.Email && req.FullName == data.FullName && req.PhotoProfile == data.PhotoProfile && req.RefreshToken == data.RefreshToken
	})).Return(&user.User{
		UserId:       data.UserId,
		Email:        data.Email,
		FullName:     data.FullName,
		PhotoProfile: data.PhotoProfile,
		RefreshToken: data.RefreshToken,
	}, nil)
}

func TestService_LoginWithGoogle(t *testing.T) {
	suite.Run(t, new(LoginWithGoogleTestSuite))
}
