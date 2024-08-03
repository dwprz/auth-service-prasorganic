package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/interface/helper"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"time"
)

type AuthRestful struct {
	authService     service.Auth
	googleOauthConf *oauth2.Config
	logger          *logrus.Logger
	helper          helper.Helper
}

func NewAuthRestful(as service.Auth, goc *oauth2.Config, l *logrus.Logger, h helper.Helper) *AuthRestful {
	return &AuthRestful{
		authService:     as,
		googleOauthConf: goc,
		logger:          l,
		helper:          h,
	}
}

func (a *AuthRestful) Register(c *fiber.Ctx) error {
	defer a.helper.HandlePanic("auth handler panic (register)", c)

	request := new(dto.RegisterReq)

	if err := c.BodyParser(request); err != nil {
		return &errors.Response{Code: 400, Message: err.Error()}
	}

	email, err := a.authService.Register(context.Background(), request)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "pending_register",
		Value:    base64.StdEncoding.EncodeToString([]byte(email)),
		HTTPOnly: true,
		Path:     "/api/auth/register/verify",
		Expires:  time.Now().Add(30 * time.Minute),
	})

	return c.Status(200).JSON(fiber.Map{"data": "register request successfully"})
}

func (a *AuthRestful) VerifyRegister(c *fiber.Ctx) error {
	defer a.helper.HandlePanic("auth handler panic (verify register)", c)

	request := new(dto.VerifyRegisterReq)

	if err := c.BodyParser(request); err != nil {
		return &errors.Response{Code: 400, Message: err.Error()}
	}

	email, err := base64.StdEncoding.DecodeString(c.Cookies("pending_register"))
	if err != nil {
		return &errors.Response{Code: 400, Message: err.Error()}
	}

	request.Email = string(email)

	err = a.authService.VerifyRegister(context.Background(), request)
	if err != nil {
		return err
	}

	c.Cookie(a.helper.ClearCookie("pending_register", "/api/auth/register/verify")) // clear cookie

	return c.Status(200).JSON(fiber.Map{"data": "verify register successfully"})
}

func (a *AuthRestful) LoginWithGoogle(c *fiber.Ctx) error {
	oauthState, err := a.helper.GenerateOauthState()
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    oauthState,
		Path:     "/api/auth/login/google/callback",
		HTTPOnly: true,
		Expires:  time.Now().Add(5 * time.Minute),
	})

	url := a.googleOauthConf.AuthCodeURL(oauthState)

	return c.Status(fiber.StatusSeeOther).Redirect(url)
}

func (a *AuthRestful) LoginWithGoogleCallback(c *fiber.Ctx) error {
	req := c.Body()

	user := new(dto.LoginWithGoogleReq)
	err := json.Unmarshal(req, user)
	if err != nil {
		return err
	}

	userId, err := gonanoid.New()
	if err != nil {
		return err
	}

	user.UserId = userId

	accessToken, err := a.helper.GenerateAccessToken(user.UserId, user.Email, "USER")
	if err != nil {
		return err
	}

	refreshToken, err := a.helper.GenerateRefreshToken()
	if err != nil {
		return err
	}

	user.RefreshToken = refreshToken

	result, err := a.authService.LoginWithGoogle(context.Background(), user)
	if err != nil {
		return err
	}

	c.Cookie(a.helper.ClearCookie("oauth_state", "/api/auth/login/google/callback")) // clear cookie

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HTTPOnly: true,
		Expires:  time.Now().Add(1 * time.Hour),
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HTTPOnly: true,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})

	return c.Status(200).JSON(fiber.Map{"data": result})
}

func (a *AuthRestful) Login(c *fiber.Ctx) error {
	req := new(dto.LoginReq)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	res, err := a.authService.Login(context.Background(), req)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    res.Tokens.AccessToken,
		Path:     "/",
		HTTPOnly: true,
		Expires:  time.Now().Add(1 * time.Hour),
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.Tokens.RefreshToken,
		Path:     "/",
		HTTPOnly: true,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})

	return c.Status(200).JSON(fiber.Map{"data": res.Data})
}

func (a *AuthRestful) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	if _, err := a.helper.VerifyJwt(refreshToken); err != nil {
		return &errors.Response{Code: 401, Message: "refresh token is invalid"}
	}

	res, err := a.authService.RefreshToken(context.Background(), refreshToken)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    res.AccessToken,
		Path:     "/",
		HTTPOnly: true,
		Expires:  time.Now().Add(1 * time.Hour),
	})

	return c.Status(201).JSON(fiber.Map{"data": "refresh token successfully"})
}

func (a *AuthRestful) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	a.authService.SetNullRefreshToken(context.Background(), refreshToken)

	// clear cookie
	c.Cookie(a.helper.ClearCookie("refresh_token", "/"))
	c.Cookie(a.helper.ClearCookie("access_token", "/"))

	return c.Status(200).JSON(fiber.Map{"data": "logout successfully"})
}
