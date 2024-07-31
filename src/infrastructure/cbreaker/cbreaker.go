package cbreaker

import (
	"github.com/sirupsen/logrus"
	"github.com/sony/gobreaker/v2"
)

type CircuitBreaker struct {
	UserGrpc *gobreaker.CircuitBreaker[any]
}

func New(logger *logrus.Logger) *CircuitBreaker {
	userGrpcCBreaker := setupForUserGrpc(logger)

	return &CircuitBreaker{
		UserGrpc: userGrpcCBreaker,
	}
}

