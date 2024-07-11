package middleware

import (
	"github.com/dwprz/prasorganic-auth-service/src/common/config"
	"github.com/dwprz/prasorganic-auth-service/src/common/custom_error"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type MiddlewareImpl struct {
	Conf   *config.Config
	Logger *logrus.Logger
}

func NewMiddleware(conf *config.Config, logger *logrus.Logger) Middleware {
	return &MiddlewareImpl{
		Conf:   conf,
		Logger: logger,
	}
}

func (m *MiddlewareImpl) ErrorMiddleware(c *fiber.Ctx, err error) error {
	m.Logger.Errorf("log error middleware | %s", err.Error())

	if validationError, ok := err.(*custom_error.ValidationError); ok {
		return c.Status(400).JSON(fiber.Map{
			"errors": validationError.Message,
		})
	}

	if responseError, ok := err.(*custom_error.ResponseError); ok {
		return c.Status(responseError.Code).JSON(fiber.Map{
			"errors": responseError.Message,
		})
	}

	return c.Status(500).JSON(fiber.Map{
		"errors": "sorry, internal server error try again later",
	})
}

func (m *MiddlewareImpl) AuthMiddleware(c *fiber.Ctx) error {
	if c.Get("AUTHORIZATION") != m.Conf.App.AuthSecretKey {
		return c.Status(401).JSON(fiber.Map{
			"errors": "invalid authorization key",
		})
	}

	return c.Next()
}
