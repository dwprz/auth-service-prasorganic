package middleware

import (
	"context"
	"io"
	"net/http"

	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/oauth"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) FetchGoogleUser(c *fiber.Ctx) error {
	if c.Query("state") != c.Cookies("oauth_state") {
		return &errors.Response{HttpCode: 401, Message: "invalid oauth state"}
	}

	ctx := context.Background()

	code := c.Query("code")
	token, err := oauth.GoogleConf.Exchange(ctx, code)
	if err != nil {
		return err
	}

	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return err
	}

	buff, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	c.Request().SetBodyRaw(buff)
	return c.Next()
}
