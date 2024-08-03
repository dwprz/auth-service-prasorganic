package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Helper interface {
	GenerateOtp() (string, error)
	GenerateOauthState() (string, error)
	GenerateAccessToken(userId string, email string, role string) (string, error)
	GenerateRefreshToken() (string, error)
	VerifyJwt(token string) (*jwt.MapClaims, error)
	HandlePanic(name string, c *fiber.Ctx)
	ClearCookie(name string, path string) *fiber.Cookie
}
