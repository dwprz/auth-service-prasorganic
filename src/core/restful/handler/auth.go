package handler

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/common/helper"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/oauth"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type AuthRESTful struct {
	authService service.Auth
}

func NewAuthRESTful(as service.Auth) *AuthRESTful {
	return &AuthRESTful{
		authService: as,
	}
}

func (a *AuthRESTful) Register(c *fiber.Ctx) error {
	request := new(dto.RegisterReq)

	if err := c.BodyParser(request); err != nil {
		return err
	}

	email, err := a.authService.Register(c.Context(), request)
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

func (a *AuthRESTful) VerifyRegister(c *fiber.Ctx) error {
	request := new(dto.VerifyOtpReq)

	if err := c.BodyParser(request); err != nil {
		return err
	}

	email, err := base64.StdEncoding.DecodeString(c.Cookies("pending_register"))
	if err != nil {
		return err
	}

	request.Email = string(email)

	err = a.authService.VerifyRegister(c.Context(), request)
	if err != nil {
		return err
	}

	c.Cookie(helper.ClearCookie("pending_register", "/api/auth/register/verify")) // clear cookie

	return c.Status(200).JSON(fiber.Map{"data": "verify register successfully"})
}

func (a *AuthRESTful) LoginWithGoogle(c *fiber.Ctx) error {
	oauthState, err := helper.GenerateOauthState()
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

	url := oauth.GoogleConf.AuthCodeURL(oauthState)

	return c.Status(fiber.StatusSeeOther).Redirect(url)
}

func (a *AuthRESTful) LoginWithGoogleCallback(c *fiber.Ctx) error {
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

	accessToken, err := helper.GenerateAccessToken(user.UserId, user.Email, "USER")
	if err != nil {
		return err
	}

	refreshToken, err := helper.GenerateRefreshToken()
	if err != nil {
		return err
	}

	user.RefreshToken = refreshToken

	result, err := a.authService.LoginWithGoogle(c.Context(), user)
	if err != nil {
		return err
	}

	c.Cookie(helper.ClearCookie("oauth_state", "/api/auth/login/google/callback")) // clear cookie

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

func (a *AuthRESTful) Login(c *fiber.Ctx) error {
	req := new(dto.LoginReq)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	res, err := a.authService.Login(c.Context(), req)
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

func (a *AuthRESTful) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	if _, err := helper.VerifyJwt(refreshToken); err != nil {
		return &errors.Response{HttpCode: 401, Message: "refresh token is invalid"}
	}

	res, err := a.authService.RefreshToken(c.Context(), refreshToken)
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

func (a *AuthRESTful) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	a.authService.SetNullRefreshToken(c.Context(), refreshToken)

	// clear cookie
	c.Cookie(helper.ClearCookie("refresh_token", "/"))
	c.Cookie(helper.ClearCookie("access_token", "/"))

	return c.Status(200).JSON(fiber.Map{"data": "logout successfully"})
}
