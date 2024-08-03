package test

import (
	"context"
	"testing"

	serviceinterface "github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/mock/cache"
	"github.com/dwprz/prasorganic-auth-service/src/mock/client"
	"github.com/dwprz/prasorganic-auth-service/src/mock/helper"
	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/common/logger"
	grpcapp "github.com/dwprz/prasorganic-auth-service/src/core/grpc/grpc"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"github.com/dwprz/prasorganic-auth-service/src/service"
	"github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

// go test -v ./src/service/test/... -count=1 -p=1
// go test -run ^TestService_Register$  -v ./src/service/test/ -count=1

type RegisterTestSuite struct {
	suite.Suite
	authService    serviceinterface.Auth
	userGrpcClient *client.UserGrpcMock
	authCache      *cache.AuthMock
	logger         *logrus.Logger
	helper         *helper.HelperMock
}

func (r *RegisterTestSuite) SetupSuite() {
	r.logger = logger.New()
	validator := validator.New()
	conf := config.New("DEVELOPMENT", r.logger)

	// mock
	r.helper = helper.NewMock()
	// mock
	r.userGrpcClient = client.NewUserMock()
	userGrpcConn := new(grpc.ClientConn)

	grpcClient := grpcapp.NewClient(r.userGrpcClient, userGrpcConn, r.logger)

	// mock
	r.authCache = cache.NewAuthMock()
	// mock
	rabbitMQClient := client.NewRabbitMQMock()

	r.authService = service.NewAuth(grpcClient, rabbitMQClient, validator, r.authCache, r.logger, conf, r.helper)
}

func (r *RegisterTestSuite) Test_Success() {

	req := &dto.RegisterReq{
		Email:    "johndoe123@gmail.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	r.userGrpcClient.Mock.On("FindByEmail", mock.Anything, &user.Email{
		Email: req.Email,
	}).Return(&user.FindUserResponse{Data: nil}, nil)

	r.helper.Mock.On("GenerateOtp").Return("123456", nil)

	r.authCache.Mock.On("CacheRegisterReq", mock.Anything, mock.MatchedBy(func(req *dto.RegisterReq) bool {
		return req.Email == req.Email && req.FullName == req.FullName && req.Password == req.Password && req.Otp == "123456"
	})).Return(nil)

	email, err := r.authService.Register(context.Background(), req)

	assert.NoError(r.T(), err)
	assert.Equal(r.T(), req.Email, email)
}

func (r *RegisterTestSuite) Test_AlreadyExists() {

	req := &dto.RegisterReq{
		Email:    "userexisted@gmail.com",
		FullName: "John Doe",
		Password: "rahasia",
	}

	r.userGrpcClient.Mock.On("FindByEmail", mock.Anything, &user.Email{
		Email: req.Email,
	}).Return(&user.FindUserResponse{Data: new(user.User)}, nil)

	email, err := r.authService.Register(context.Background(), req)
	errorRes, ok := err.(*errors.Response)

	assert.Equal(r.T(), true, ok)
	assert.Equal(r.T(), 409, errorRes.Code)
	assert.Equal(r.T(), "", email)
}

func TestService_Register(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}
