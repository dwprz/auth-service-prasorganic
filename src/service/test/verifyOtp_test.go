package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-auth-service/src/common/util"
	"github.com/dwprz/prasorganic-auth-service/src/core/broker/producer"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/mock/cache"
	"github.com/dwprz/prasorganic-auth-service/src/mock/delivery"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	serviceimpl "github.com/dwprz/prasorganic-auth-service/src/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// go test -p=1 -v ./src/service/test/... -count=1
// go test -run ^TestService_VerifyOtp$ -v ./src/service/test/ -count=1

type VerifyOtpTestSuite struct {
	suite.Suite
	EmailBroker *delivery.EmailBrokerMock
	otpService  service.Otp
	otpCache    *cache.OtpMock
}

func (v *VerifyOtpTestSuite) SetupSuite() {
	// mock
	v.EmailBroker = delivery.NewEmailBrokerMock()

	rabbitMQProducer := producer.NewRabbitMQ(v.EmailBroker)

	// mock
	v.otpCache = cache.NewOtpMock()
	util := util.New()

	v.otpService = serviceimpl.NewOtp(rabbitMQProducer, v.otpCache, util)
}

func (v *VerifyOtpTestSuite) Test_Success() {
	email := `johndoe123@gmail.com`

	req := &dto.VerifyOtpReq{
		Email: email,
		Otp:   "123456",
	}

	v.MockOtpCache_FindByEmail(email, &dto.SendOtpReq{
		Email: email,
		Otp:   req.Otp,
	})

	err := v.otpService.Verify(context.Background(), req)
	assert.NoError(v.T(), err)
}

func (v *VerifyOtpTestSuite) MockOtpCache_FindByEmail(email string, returnArg *dto.SendOtpReq) {

	v.otpCache.Mock.On("FindByEmail", mock.Anything, email).Return(returnArg)
}

func TestService_VerifyOtp(t *testing.T) {
	suite.Run(t, new(VerifyOtpTestSuite))
}
