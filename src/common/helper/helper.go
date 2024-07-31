package helper

import (
	"fmt"
	"math/rand"
	"github.com/dwprz/prasorganic-auth-service/interface/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Helper struct{}

func New() helper.Helper {
	return &Helper{}
}

func (h *Helper) GenerateOtp() string {
	otp := rand.Intn(1000000)
	return fmt.Sprintf("%06d", otp)
}

func (h *Helper) HandlePanic(name string, c *fiber.Ctx, logger *logrus.Logger) {
	message := recover()

	if message != nil {
		logger.Errorf(name+" | %v", message)

		c.Status(500).JSON(fiber.Map{
			"errors": "sorry, internal server error try again later",
		})
	}
}
