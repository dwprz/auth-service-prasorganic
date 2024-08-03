package oauth

import (
	"github.com/dwprz/prasorganic-auth-service/src/interface/helper"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewGoogleConfig(conf *config.Config, h helper.Helper) *oauth2.Config {
	googleConf := &oauth2.Config{
		ClientID:     conf.GoogleOauth.ClientId,
		ClientSecret: conf.GoogleOauth.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		RedirectURL: conf.GoogleOauth.RedirectURL,
		Endpoint:    google.Endpoint,
	}

	return googleConf
}
