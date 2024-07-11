package middleware

import "github.com/gofiber/fiber/v2"

type Middleware interface {
	ErrorMiddleware(ctx *fiber.Ctx, err error) error
	AuthMiddleware(ctx *fiber.Ctx) error
}
