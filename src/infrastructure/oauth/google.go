package oauth

import (
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleConf *oauth2.Config

func init() {
	GoogleConf = &oauth2.Config{
		ClientID:     config.Conf.GoogleOauth.ClientId,
		ClientSecret: config.Conf.GoogleOauth.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		RedirectURL: config.Conf.GoogleOauth.RedirectURL,
		Endpoint:    google.Endpoint,
	}
}
