package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Helper interface {
	GenerateOtp() string
	HandlePanic(name string, c *fiber.Ctx, logger *logrus.Logger)
}
