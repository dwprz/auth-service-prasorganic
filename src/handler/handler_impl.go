package handler

import (
	"context"
	"encoding/base64"
	"time"
	"github.com/dwprz/prasorganic-auth-service/src/common/custom_error"
	"github.com/dwprz/prasorganic-auth-service/src/common/helper"
	"github.com/dwprz/prasorganic-auth-service/src/common/model/dto"
	"github.com/dwprz/prasorganic-auth-service/src/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthHandlerImpl struct {
	Service service.AuthService
	Logger  *logrus.Logger
}

func NewAuthHandler(service service.AuthService, logger *logrus.Logger) AuthHandler {
	return &AuthHandlerImpl{
		Service: service,
		Logger:  logger,
	}
}

func (h *AuthHandlerImpl) Register(c *fiber.Ctx) error {
	defer helper.HandlePanic("auth handler panic (register)", c, h.Logger)

	request := new(dto.RegisterReq)

	if err := c.BodyParser(request); err != nil {
		return &custom_error.ResponseError{Code: 400, Message: err.Error()}
	}

	ctx := context.Background()

	email, err := h.Service.Register(ctx, request)
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
