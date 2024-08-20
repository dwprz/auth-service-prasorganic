package util

import (
	"github.com/dwprz/prasorganic-auth-service/src/core/broker/producer"
	"github.com/dwprz/prasorganic-auth-service/src/interface/cache"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/mock/util"

	serviceimpl "github.com/dwprz/prasorganic-auth-service/src/service"
)

func InitOtpService(rc *producer.RabbitMQ, oc cache.Otp, u *util.UtilMock) service.Otp {
	otpService := serviceimpl.NewOtp(rc, oc, u)
	return otpService
}
