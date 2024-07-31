package middleware

import (
	"github.com/dwprz/prasorganic-auth-service/src/common/errors"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Middleware struct {
	conf   *config.Config
	logger *logrus.Logger
}

func New(conf *config.Config, logger *logrus.Logger) *Middleware {
	return &Middleware{
		conf:   conf,
		logger: logger,
	}
}

func (m *Middleware) Error(c *fiber.Ctx, err error) error {
	m.logger.WithFields(logrus.Fields{
		"host":     c.Hostname(),
		"ip":       c.IP(),
		"protocol": c.Protocol(),
		"location": c.OriginalURL(),
		"method":   c.Method(),
		"from":     "error middleware",
	}).Error(err.Error())

	if validationError, ok := err.(validator.ValidationErrors); ok {

		return c.Status(400).JSON(fiber.Map{
			"errors": map[string]any{
				"field":       validationError[0].Field(),
				"description": validationError[0].Error(),
			},
		})
	}

	if responseError, ok := err.(*errors.Response); ok {
		return c.Status(responseError.Code).JSON(fiber.Map{
			"errors": responseError.Message,
		})
	}

	return c.Status(500).JSON(fiber.Map{
		"errors": "sorry, internal server error try again later",
	})
}
