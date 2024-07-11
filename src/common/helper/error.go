package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func HandlePanic(name string, c *fiber.Ctx, logger *logrus.Logger) {
	message := recover()

	if message != nil {
		logger.Errorf(name+" | %v", message)

		c.Status(500).JSON(fiber.Map{
			"errors": "sorry, internal server error try again later",
		})
	}
}
