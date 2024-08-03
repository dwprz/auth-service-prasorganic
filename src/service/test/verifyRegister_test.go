package test

import (
	"context"
	serviceinterface "github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/mock/cache"
	"github.com/dwprz/prasorganic-auth-service/src/mock/client"
	"github.com/dwprz/prasorganic-auth-service/src/common/helper"
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
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"testing"
)

// go test -p=1 -v ./src/service/test/... -count=1
// go test -run ^TestService_VerifyRegister$ -v ./src/service/test/ -count=1

type VerifyRegisterTestSuite struct {
	suite.Suite
	authService    serviceinterface.Auth
	userGrpcClient *client.UserGrpcMock
	authCache      *cache.AuthMock
	logger         *logrus.Logger
}

func (v *VerifyRegisterTestSuite) SetupSuite() {
	v.logger = logger.New()
	validator := validator.New()
	conf := config.New("DEVELOPMENT", v.logger)
	helper := helper.New(conf, v.logger)

	// mock
	v.userGrpcClient = client.NewUserMock()
	userGrpcConn := new(grpc.ClientConn)

	grpcClient := grpcapp.NewClient(v.userGrpcClient, userGrpcConn, v.logger)

	// mock
	v.authCache = cache.NewAuthMock()

	// mock
	rabbitMQClient := client.NewRabbitMQMock()

	v.authService = service.NewAuth(grpcClient, rabbitMQClient, validator, v.authCache, v.logger, conf, helper)
}

func (v *VerifyRegisterTestSuite) Test_Success() {
	verifyRegisterReq := &dto.VerifyRegisterReq{
		Email: "johndoe123@gmail.com",
		Otp:   "123456",
	}

	registerReq := &dto.RegisterReq{
		Email:    "johndoe123@gmail.com",
		FullName: "John Doe",
		Password: "rahasia",
		Otp:      verifyRegisterReq.Otp,
	}

	v.authCache.Mock.On("FindRegisterReq", mock.Anything, verifyRegisterReq.Email).Return(registerReq)

	v.MockUserGrpcClient_Create(registerReq, "rahasia", nil)

	err := v.authService.VerifyRegister(context.Background(), verifyRegisterReq)
	assert.NoError(v.T(), err)
}

func (v *VerifyRegisterTestSuite) Test_InvalidEmail() {
	verifyRegisterReq := &dto.VerifyRegisterReq{
		Email: "123456",
		Otp:   "123456",
	}

	err := v.authService.VerifyRegister(context.Background(), verifyRegisterReq)
	assert.Error(v.T(), err)

	errorResp, ok := err.(validator.ValidationErrors)

	assert.Equal(v.T(), true, ok)
	assert.NotEmpty(v.T(), errorResp)
}

// memberikan argumen password secara langsung karena registerReq di hash method authService.VerifyRegister
func (v *VerifyRegisterTestSuite) MockUserGrpcClient_Create(data *dto.RegisterReq, password string, returnArg error) {

	v.userGrpcClient.Mock.On("Create", mock.Anything, mock.MatchedBy(func(req *user.RegisterRequest) bool {
		err := bcrypt.CompareHashAndPassword([]byte(req.Password), []byte(password))
		return req.Email == data.Email && req.FullName == data.FullName && err == nil
	})).Return(returnArg)
}

func TestService_VerifyRegister(t *testing.T) {
	suite.Run(t, new(VerifyRegisterTestSuite))
}
