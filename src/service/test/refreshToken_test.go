package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/common/helper"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/mock/cache"
	"github.com/dwprz/prasorganic-auth-service/src/mock/delivery"
	svcmock "github.com/dwprz/prasorganic-auth-service/src/mock/service"
	serviceimpl "github.com/dwprz/prasorganic-auth-service/src/service"
	pb "github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

// go test -p=1 -v ./src/service/test/... -count=1
// go test -run ^TestService_RefreshToken$ -v ./src/service/test/ -count=1

type RefreshTokenTestSuite struct {
	suite.Suite
	userGrpcDelivery *delivery.UserGrpcMock
	authService      service.Auth
	otpService       *svcmock.OtpMock
	authCache        *cache.AuthMock
}

func (r *RefreshTokenTestSuite) SetupSuite() {
	// mock
	helper.LogJSON("apakah ini dikesekusi 1")
	r.userGrpcDelivery = delivery.NewUserGrpcMock()
	userGrpcConn := new(grpc.ClientConn)

	grpcClient := client.NewGrpc(r.userGrpcDelivery, userGrpcConn)

	// mock
	r.authCache = cache.NewAuthMock()
	r.otpService = svcmock.NewOtpMock()

	r.authService = serviceimpl.NewAuth(grpcClient, r.otpService, r.authCache)
}

func (r *RefreshTokenTestSuite) Test_Success() {
	helper.LogJSON("apakah ini dikesekusi")


	user := &pb.User{
		UserId: "hyfa_5Sq7nQcaY6ACksXP",
		Email:  "johndoe123@gmail.com",
		Role:   "USER",
	}

	refreshToken := `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.
					 eyJleHAiOjE3MjY3MDgyNzMsImlkIjoiQndfb
					 nBmeUhpRXlRY2d3LWM1eGtzIiwiaXNzIjoicH
					 Jhc29yZ2FuaWMtYXV0aC1zZXJ2aWNlIn0.Q0A
					 0veFoysALdg9qqypxY8lKxMFNDmBJR-hyEAHg
					 zBZu-NNzmRIl2MMUjQfBAjJI-KoNj7n9keH6x
					 gwXcuRTVeDfMFtaPm7Dp5ezEgqLXIOTR3o3Sk
					 DhljqYN0HC7UBZQ4R-SyP26MTgVBZQOyuymj3
					 XhFG29OABaJb12GglrtE`

	r.MockUserGrpcClient_FindByRefreshToken(refreshToken, user, nil)

	res, err := r.authService.RefreshToken(context.Background(), refreshToken)
	assert.NoError(r.T(), err)

	assert.NotEmpty(r.T(), res.AccessToken)
}

func (r *RefreshTokenTestSuite) Test_NotFound() {
	refreshToken := `not-found-refresh-token`

	r.MockUserGrpcClient_FindByRefreshToken(refreshToken, nil, nil)

	res, err := r.authService.RefreshToken(context.Background(), refreshToken)
	assert.Error(r.T(), err)

	resErr, ok := err.(*errors.Response)
	assert.True(r.T(), ok)
	assert.Equal(r.T(), resErr.HttpCode, 404)

	assert.Nil(r.T(), res)
}

func (r *RefreshTokenTestSuite) MockUserGrpcClient_FindByRefreshToken(refreshToken string, returnArg *pb.User, err error) {

	r.userGrpcDelivery.Mock.On("FindByRefreshToken", mock.Anything, &pb.RefreshToken{Token: refreshToken}).Return(&pb.FindUserRes{
		Data: returnArg,
	}, err)
}

func TestService_RefreshToken(t *testing.T) {
	suite.Run(t, new(RefreshTokenTestSuite))
}
