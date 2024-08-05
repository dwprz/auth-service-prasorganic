package test

import (
	"context"
	"testing"
	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/common/helper"
	"github.com/dwprz/prasorganic-auth-service/src/common/logger"
	grpcapp "github.com/dwprz/prasorganic-auth-service/src/core/grpc/grpc"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	serviceinterface "github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/mock/cache"
	"github.com/dwprz/prasorganic-auth-service/src/mock/client"
	svcmock "github.com/dwprz/prasorganic-auth-service/src/mock/service"
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

// go test -p=1 -v ./src/service/test/... -count=1
// go test -run ^TestService_VerifyRegister$ -v ./src/service/test/ -count=1

type VerifyRegisterTestSuite struct {
	suite.Suite
	authService    serviceinterface.Auth
	otpService     *svcmock.OtpMock
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
	v.otpService = svcmock.NewOtpMock()

	v.authService = service.NewAuth(grpcClient, v.otpService, validator, v.authCache, v.logger, conf, helper)
}

func (v *VerifyRegisterTestSuite) Test_Success() {
	verifyRegisterReq := &dto.VerifyOtpReq{
		Email: "johndoe123@gmail.com",
		Otp:   "123456",
	}

	registerReq := &dto.RegisterReq{
		Email:    "johndoe123@gmail.com",
		FullName: "John Doe",
		Password: "$2a$10$MI6/KH0.g8NSLthw86K9we9SFhHIT1c6RStWasZHBPAxVrPelFZAu",
	}

	v.otpService.Mock.On("Verify", mock.Anything, verifyRegisterReq).Return(nil)
	v.authCache.Mock.On("FindRegisterReq", mock.Anything, verifyRegisterReq.Email).Return(registerReq)

	v.MockUserGrpcClient_Create(registerReq, nil)

	err := v.authService.VerifyRegister(context.Background(), verifyRegisterReq)
	assert.NoError(v.T(), err)
}

func (v *VerifyRegisterTestSuite) Test_Otp() {
	verifyRegisterReq := &dto.VerifyOtpReq{
		Email: "johndoe123@gmail.com",
		Otp:   "invalid otp",
	}

	errRes := &errors.Response{HttpCode: 400, Message: "otp is invalid"}
	v.otpService.Mock.On("Verify", mock.Anything, verifyRegisterReq).Return(errRes)

	err := v.authService.VerifyRegister(context.Background(), verifyRegisterReq)
	assert.Error(v.T(), err)

	errorResp, ok := err.(*errors.Response)

	assert.Equal(v.T(), true, ok)
	assert.Equal(v.T(), 400, errorResp.HttpCode)
}

func (v *VerifyRegisterTestSuite) MockUserGrpcClient_Create(data *dto.RegisterReq, returnArg error) {

	v.userGrpcClient.Mock.On("Create", mock.Anything, mock.MatchedBy(func(req *user.RegisterRequest) bool {
		return req.Email == data.Email && req.FullName == data.FullName && data.Password == req.Password
	})).Return(returnArg)
}

func TestService_VerifyRegister(t *testing.T) {
	suite.Run(t, new(VerifyRegisterTestSuite))
}
