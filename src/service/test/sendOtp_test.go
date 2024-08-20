package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-auth-service/src/common/util"
	"github.com/dwprz/prasorganic-auth-service/src/core/broker/producer"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/mock/cache"
	"github.com/dwprz/prasorganic-auth-service/src/mock/delivery"
	serviceimpl "github.com/dwprz/prasorganic-auth-service/src/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// go test -p=1 -v ./src/service/test/... -count=1
// go test -run ^TestService_SendOtp$ -v ./src/service/test/ -count=1

type SendOtpTestSuite struct {
	suite.Suite
	EmailBroker *delivery.EmailBrokerMock
	otpService  service.Otp
	otpCache    *cache.OtpMock
}

func (s *SendOtpTestSuite) SetupSuite() {
	// mock
	s.EmailBroker = delivery.NewEmailBrokerMock()

	rabbitMQProducer := producer.NewRabbitMQ(s.EmailBroker)

	// mock
	s.otpCache = cache.NewOtpMock()
	util := util.New()

	s.otpService = serviceimpl.NewOtp(rabbitMQProducer, s.otpCache, util)
}

func (s *SendOtpTestSuite) Test_Success() {
	email := `johndoe123@gmail.com`

	err := s.otpService.Send(context.Background(), email)
	assert.NoError(s.T(), err)
}

func TestService_SendOtp(t *testing.T) {
	suite.Run(t, new(SendOtpTestSuite))
}
