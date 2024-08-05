package helper

import (
	"context"

	"github.com/dwprz/prasorganic-auth-service/src/model/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Helper interface {
	GenerateOauthState() (string, error)
	GenerateAccessToken(userId string, email string, role string) (string, error)
	GenerateRefreshToken() (string, error)
	GetMetadata(ctx context.Context) *entity.Metadata
	VerifyJwt(token string) (*jwt.MapClaims, error)
	HandlePanic(name string, c *fiber.Ctx)
	ClearCookie(name string, path string) *fiber.Cookie
}
