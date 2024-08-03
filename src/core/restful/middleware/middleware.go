package middleware

import (
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Middleware struct {
	conf            *config.Config
	googleOauthConf *oauth2.Config
	logger          *logrus.Logger
}

func New(conf *config.Config, goc *oauth2.Config, logger *logrus.Logger) *Middleware {
	return &Middleware{
		conf:            conf,
		googleOauthConf: goc,
		logger:          logger,
	}
}
