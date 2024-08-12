package util

import (
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-auth-service/src/interface/cache"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	serviceimpl "github.com/dwprz/prasorganic-auth-service/src/service"
)


func InitAuthServiceTest(gc *client.Grpc, os service.Otp, ac cache.Auth) (service.Auth) {
	authService := serviceimpl.NewAuth(gc, os, ac)
	return authService
}
