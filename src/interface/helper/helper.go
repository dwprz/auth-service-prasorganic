package helper

import (
	"github.com/gofiber/fiber/v2"
)

type Helper interface {
	GenerateOtp() (string, error)
	GenerateOauthState() (string, error)
	GenerateAccessToken(userId string, email string, role string) (string, error)
	GenerateRefreshToken() (string, error)
	HandlePanic(name string, c *fiber.Ctx)
}
