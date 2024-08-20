package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/mock/cache"
	"github.com/dwprz/prasorganic-auth-service/src/mock/delivery"
	svcmock "github.com/dwprz/prasorganic-auth-service/src/mock/service"
	serviceimpl "github.com/dwprz/prasorganic-auth-service/src/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

// go test -p=1 -v ./src/service/test/... -count=1
// go test -run ^TestService_SetNullRefreshToken$ -v ./src/service/test/ -count=1

type SetNullRefreshTokenTestSuite struct {
	suite.Suite
	userGrpcDelivery *delivery.UserGrpcMock
	authService      service.Auth
	otpService       *svcmock.OtpMock
	authCache        *cache.AuthMock
}

func (s *SetNullRefreshTokenTestSuite) SetupSuite() {
	// mock
	s.userGrpcDelivery = delivery.NewUserGrpcMock()
	userGrpcConn := new(grpc.ClientConn)

	grpcClient := client.NewGrpc(s.userGrpcDelivery, userGrpcConn)

	// mock
	s.authCache = cache.NewAuthMock()
	s.otpService = svcmock.NewOtpMock()

	s.authService = serviceimpl.NewAuth(grpcClient, s.otpService, s.authCache)
}

func (s *SetNullRefreshTokenTestSuite) Test_Success() {
	refreshToken := `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.
					 eyJleHAiOjE3MjY3MDgyNzMsImlkIjoiQndfb
					 nBmeUhpRXlRY2d3LWM1eGtzIiwiaXNzIjoicH
					 Jhc29yZ2FuaWMtYXV0aC1zZXJ2aWNlIn0.Q0A
					 0veFoysALdg9qqypxY8lKxMFNDmBJR-hyEAHg
					 zBZu-NNzmRIl2MMUjQfBAjJI-KoNj7n9keH6x
					 gwXcuRTVeDfMFtaPm7Dp5ezEgqLXIOTR3o3Sk
					 DhljqYN0HC7UBZQ4R-SyP26MTgVBZQOyuymj3
					 XhFG29OABaJb12GglrtE`

	s.MockUserGrpcClient_SetNullRefreshToken(refreshToken, nil)

	err := s.authService.SetNullRefreshToken(context.Background(), refreshToken)
	assert.NoError(s.T(), err)
}

func (s *SetNullRefreshTokenTestSuite) MockUserGrpcClient_SetNullRefreshToken(refreshToken string, err error) {

	s.userGrpcDelivery.Mock.On("SetNullRefreshToken", mock.Anything, refreshToken).Return(err)
}

func TestService_SetNullRefreshToken(t *testing.T) {
	suite.Run(t, new(SetNullRefreshTokenTestSuite))
}
