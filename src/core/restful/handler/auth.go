package handler

import (
	"context"
	"encoding/base64"
	"time"
	"github.com/dwprz/prasorganic-auth-service/interface/helper"
	"github.com/dwprz/prasorganic-auth-service/interface/service"
	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthRestful struct {
	authService service.Authentication
	logger      *logrus.Logger
	helper      helper.Helper
}

func NewAuthRestful(as service.Authentication, l *logrus.Logger, h helper.Helper) *AuthRestful {
	return &AuthRestful{
		authService: as,
		logger:      l,
		helper:      h,
	}
}

func (a *AuthRestful) Register(c *fiber.Ctx) error {
	defer a.helper.HandlePanic("auth handler panic (register)", c, a.logger)

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
		Path:     "/",
		Expires:  time.Now().Add(30 * time.Minute),
	})

	return c.Status(200).JSON(fiber.Map{"data": "register request successfully"})
}

func (a *AuthRestful) VerifyRegister(c *fiber.Ctx) error {
	defer a.helper.HandlePanic("auth handler panic (verify register)", c, a.logger)

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

	c.ClearCookie("pending_register")

	return c.Status(200).JSON(fiber.Map{"data": "verify register successfully"})
}